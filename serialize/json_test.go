package serialize_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"slices"
	"strings"
	"testing"

	openapi "github.com/faustbrian/go-openapi"
	"github.com/faustbrian/go-openapi/jsonvalue"
	"github.com/faustbrian/go-openapi/parse"
	"github.com/faustbrian/go-openapi/serialize"
)

func TestJSONStringEscapingMatchesEncodingJSON(t *testing.T) {
	t.Parallel()

	for _, text := range []string{
		"", "plain", "quoted\"\\value", "\x00\b\f\n\r\t\x1f",
		"<tag>&value", "café 世界 😀", "line\u2028paragraph\u2029",
	} {
		value, err := jsonvalue.String(text)
		if err != nil {
			t.Fatal(err)
		}
		raw, err := json.Marshal(text)
		if err != nil {
			t.Fatal(err)
		}
		var scalar bytes.Buffer
		if err := serialize.JSON(context.Background(), &scalar, value, serialize.DefaultOptions()); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(scalar.Bytes(), raw) {
			t.Fatalf("string %q = %q, want %q", text, scalar.Bytes(), raw)
		}

		tail, _ := jsonvalue.String("tail")
		members := []jsonvalue.Member{{Name: text, Value: value}, {Name: "tail", Value: tail}}
		object, err := jsonvalue.Object(members)
		if err != nil {
			t.Fatal(err)
		}
		before, _ := object.Members()
		for _, mode := range []serialize.Mode{serialize.Preserve, serialize.Canonical, serialize.Preserve} {
			options := serialize.DefaultOptions()
			options.Mode = mode
			ordered := slices.Clone(members)
			if mode == serialize.Canonical {
				slices.SortFunc(ordered, func(left, right jsonvalue.Member) int {
					return strings.Compare(left.Name, right.Name)
				})
			}
			var want bytes.Buffer
			want.WriteByte('{')
			for index, member := range ordered {
				if index != 0 {
					want.WriteByte(',')
				}
				key, _ := json.Marshal(member.Name)
				memberText, _ := member.Value.Text()
				encoded, _ := json.Marshal(memberText)
				want.Write(key)
				want.WriteByte(':')
				want.Write(encoded)
			}
			want.WriteByte('}')
			var output bytes.Buffer
			if err := serialize.JSON(context.Background(), &output, object, options); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(output.Bytes(), want.Bytes()) {
				t.Fatalf("mode %d key/value %q = %q, want %q", mode, text, output.Bytes(), want.Bytes())
			}
			after, _ := object.Members()
			if !reflect.DeepEqual(before, after) {
				t.Fatal("serialization changed source members")
			}
		}
	}
}

func TestJSONPreservesOrderAndExactNumbers(t *testing.T) {
	t.Parallel()

	document := mustDocument(t, `{
		"openapi":"3.2.0",
		"x-number":-0.0e+2,
		"info":{"version":"1","title":"API"},
		"paths":{}
	}`)
	var output bytes.Buffer
	if err := serialize.JSON(
		context.Background(),
		&output,
		document,
		serialize.DefaultOptions(),
	); err != nil {
		t.Fatal(err)
	}
	if output.String() != `{"openapi":"3.2.0","x-number":-0.0e+2,"info":{"version":"1","title":"API"},"paths":{}}` {
		t.Fatalf("preserving JSON = %s", output.String())
	}
}

func TestJSONCanonicalizesObjectOrderRecursively(t *testing.T) {
	t.Parallel()

	document := mustDocument(t, `{
		"paths":{},
		"info":{"version":"1","title":"API"},
		"openapi":"3.2.0"
	}`)
	options := serialize.DefaultOptions()
	options.Mode = serialize.Canonical
	var output bytes.Buffer
	if err := serialize.JSON(context.Background(), &output, document, options); err != nil {
		t.Fatal(err)
	}
	if output.String() != `{"info":{"title":"API","version":"1"},"openapi":"3.2.0","paths":{}}` {
		t.Fatalf("canonical JSON = %s", output.String())
	}
}

func TestJSONEnforcesOutputLimitAndCancellation(t *testing.T) {
	t.Parallel()

	document := mustDocument(t, `{
		"openapi":"3.2.0",
		"info":{"title":"API","version":"1"},
		"paths":{}
	}`)
	options := serialize.DefaultOptions()
	options.MaxBytes = 8
	var output bytes.Buffer
	if err := serialize.JSON(
		context.Background(), &output, document, options,
	); !errors.Is(err, serialize.ErrLimitExceeded) {
		t.Fatalf("limit error = %v", err)
	}
	if output.Len() > options.MaxBytes {
		t.Fatalf("wrote %d bytes past limit %d", output.Len(), options.MaxBytes)
	}
	options = serialize.DefaultOptions()
	options.MaxNodes = 1
	if err := serialize.JSON(
		context.Background(), &bytes.Buffer{}, document, options,
	); !errors.Is(err, serialize.ErrLimitExceeded) {
		t.Fatalf("node limit error = %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := serialize.JSON(ctx, &bytes.Buffer{}, document, serialize.DefaultOptions()); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation error = %v", err)
	}
}

func TestJSONRejectsInvalidInputsAndPropagatesWriterErrors(t *testing.T) {
	t.Parallel()

	document := mustDocument(t, `{
		"openapi":"3.2.0",
		"info":{"title":"API","version":"1"},
		"paths":{}
	}`)
	//nolint:staticcheck // This assertion verifies the nil-context contract.
	if err := serialize.JSON(
		//lint:ignore SA1012 This assertion verifies the nil-context contract.
		nil, &bytes.Buffer{}, document, serialize.DefaultOptions(),
	); err == nil {
		t.Fatal("nil context was accepted")
	}
	if err := serialize.JSON(
		context.Background(), nil, document, serialize.DefaultOptions(),
	); err == nil {
		t.Fatal("nil writer was accepted")
	}
	if err := serialize.JSON(
		context.Background(), &bytes.Buffer{}, nil, serialize.DefaultOptions(),
	); err == nil {
		t.Fatal("nil source was accepted")
	}
	sentinel := errors.New("writer failed")
	if err := serialize.JSON(
		context.Background(),
		failingWriter{err: sentinel},
		document,
		serialize.DefaultOptions(),
	); !errors.Is(err, sentinel) {
		t.Fatalf("writer error = %v", err)
	}
}

type failingWriter struct {
	err error
}

func (writer failingWriter) Write([]byte) (int, error) {
	return 0, writer.err
}

func mustDocument(t *testing.T, raw string) openapi.Document {
	t.Helper()
	document, err := openapi.ParseJSON(
		context.Background(),
		strings.NewReader(raw),
		parse.DefaultLimits(),
	)
	if err != nil {
		t.Fatal(err)
	}
	return document
}
