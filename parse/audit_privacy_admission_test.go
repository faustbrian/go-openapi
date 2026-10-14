package parse_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/faustbrian/go-openapi/v2/jsonvalue"
	"github.com/faustbrian/go-openapi/v2/parse"
)

type auditFailingReader struct{ failure error }

func (reader auditFailingReader) Read([]byte) (int, error) { return 0, reader.failure }

func TestParserDefaultPrivacyRetainsAuthorizedReaderCause(t *testing.T) {
	for _, decoder := range []struct {
		name string
		run  func(context.Context, io.Reader, parse.Limits) (jsonvalue.Value, error)
	}{{"JSON", parse.JSON}, {"YAML", parse.YAML}} {
		t.Run(decoder.name, func(t *testing.T) {
			cause := errors.New("caller-private-diagnostic")
			value, err := decoder.run(t.Context(), auditFailingReader{failure: cause}, parse.DefaultLimits())
			var diagnostic *parse.Error
			if value.Kind() != jsonvalue.InvalidKind || !errors.Is(err, cause) || !errors.As(err, &diagnostic) || !errors.Is(diagnostic.Cause, cause) {
				t.Fatalf("reader failure lost classification or inspectable cause: %v", err)
			}
			if decoder.name == "JSON" && diagnostic.Cause != cause {
				t.Fatal("JSON reader cause identity changed")
			}
			if strings.Contains(err.Error(), cause.Error()) || diagnostic.Code != "reader_failed" {
				t.Fatal("package default diagnostic exposed caller cause")
			}
		})
	}
}

func TestParserSemanticAdmissionRetainsExactAcceptedLimits(t *testing.T) {
	for _, decoder := range []struct {
		name string
		run  func(context.Context, io.Reader, parse.Limits) (jsonvalue.Value, error)
	}{{"JSON", parse.JSON}, {"YAML", parse.YAML}} {
		for _, boundary := range []struct {
			name, accepted, rejected string
			set                      func(*parse.Limits)
		}{
			{"bytes", "0", "[]", func(limits *parse.Limits) { limits.MaxBytes = 1 }},
			{"tokens", "null", "[0]", func(limits *parse.Limits) { limits.MaxTokens = 1 }},
			{"values", "null", "[0]", func(limits *parse.Limits) { limits.MaxTotalValues = 1 }},
			{"depth", "null", "[0]", func(limits *parse.Limits) { limits.MaxDepth = 1 }},
			{"members", `{"a":0}`, `{"a":0,"b":0}`, func(limits *parse.Limits) { limits.MaxObjectMembers = 1 }},
			{"items", "[0]", "[0,0]", func(limits *parse.Limits) { limits.MaxArrayItems = 1 }},
			{"scalar", `"a"`, `"ab"`, func(limits *parse.Limits) { limits.MaxScalarBytes = 1 }},
		} {
			t.Run(decoder.name+"/"+boundary.name, func(t *testing.T) {
				limits := parse.DefaultLimits()
				boundary.set(&limits)
				value, err := decoder.run(t.Context(), strings.NewReader(boundary.accepted), limits)
				if value.Kind() == jsonvalue.InvalidKind || err != nil {
					t.Fatalf("exact accepted input: %v", err)
				}
				value, err = decoder.run(t.Context(), strings.NewReader(boundary.rejected), limits)
				if value.Kind() != jsonvalue.InvalidKind || !errors.Is(err, parse.ErrLimitExceeded) {
					t.Fatalf("excess input published semantic value or lost limit: %v", err)
				}
			})
		}
	}
}

func TestMalformedDecoderDefaultPrivacyKeepsInspectableCause(t *testing.T) {
	for _, decoder := range []struct {
		name, source string
		run          func(context.Context, io.Reader, parse.Limits) (jsonvalue.Value, error)
	}{
		{"JSON", `{"x":"caller-private-diagnostic","y":?}`, parse.JSON},
		{"YAML", "x: *caller-private-diagnostic", parse.YAML},
	} {
		t.Run(decoder.name, func(t *testing.T) {
			value, err := decoder.run(t.Context(), strings.NewReader(decoder.source), parse.DefaultLimits())
			var diagnostic *parse.Error
			if value.Kind() != jsonvalue.InvalidKind || !errors.As(err, &diagnostic) || diagnostic.Cause == nil {
				t.Fatalf("decoder failure lost zero-value or cause contract: %v", err)
			}
			if strings.Contains(err.Error(), "caller-private-diagnostic") {
				t.Fatal("default parser diagnostic exposed input text")
			}
		})
	}
}
