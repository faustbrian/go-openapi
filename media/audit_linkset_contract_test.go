package media_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/faustbrian/go-openapi/v2/jsonvalue"
	"github.com/faustbrian/go-openapi/v2/media"
)

func TestLinksetAnchorAppliesIndependentOfMemberOrder(t *testing.T) {
	const expected = `</one>; rel="next"; anchor="https://example.test/base"`
	for _, raw := range []string{
		`{"linkset":[{"anchor":"https://example.test/base","next":[{"href":"/one"}]}]}`,
		`{"linkset":[{"next":[{"href":"/one"}],"anchor":"https://example.test/base"}]}`,
	} {
		value := mustMediaValue(t, raw)
		text, err := media.SerializeLinkset(value, 1, 1024)
		if text != expected || err != nil {
			t.Fatalf("ordered context = %q, %v", text, err)
		}
	}
}

func TestLinksetPropagatesFailedInputRepresentation(t *testing.T) {
	// One valid ASCII title exceeds the existing default representation budget.
	// The explicit application limit is larger, so it cannot mask marshal failure.
	title, err := jsonvalue.String(strings.Repeat("x", jsonvalue.DefaultMarshalLimits().MaxBytes))
	if err != nil {
		t.Fatal(err)
	}
	href, err := jsonvalue.String("/one")
	if err != nil {
		t.Fatal(err)
	}
	target, err := jsonvalue.Object([]jsonvalue.Member{{Name: "href", Value: href}, {Name: "title", Value: title}})
	if err != nil {
		t.Fatal(err)
	}
	targets, err := jsonvalue.Array([]jsonvalue.Value{target})
	if err != nil {
		t.Fatal(err)
	}
	context, err := jsonvalue.Object([]jsonvalue.Member{{Name: "next", Value: targets}})
	if err != nil {
		t.Fatal(err)
	}
	contexts, err := jsonvalue.Array([]jsonvalue.Value{context})
	if err != nil {
		t.Fatal(err)
	}
	value, err := jsonvalue.Object([]jsonvalue.Member{{Name: "linkset", Value: contexts}})
	if err != nil {
		t.Fatal(err)
	}
	if encoded, err := value.MarshalJSON(); encoded != nil || !errors.Is(err, jsonvalue.ErrMarshalLimit) {
		t.Fatalf("independent marshal control = %d bytes, %v", len(encoded), err)
	}
	if err := media.ValidateLinksetJSON(value, 1, 4*jsonvalue.DefaultMarshalLimits().MaxBytes); !errors.Is(err, jsonvalue.ErrMarshalLimit) {
		t.Fatalf("input admission lost marshal failure: %v", err)
	}
	if text, err := media.SerializeLinkset(value, 1, 4*jsonvalue.DefaultMarshalLimits().MaxBytes); text != "" || !errors.Is(err, jsonvalue.ErrMarshalLimit) {
		t.Fatalf("failed serialization published %d bytes, %v", len(text), err)
	}
}

func TestLinksetRejectsLateInvalidAnchorWithoutPartialText(t *testing.T) {
	value := mustMediaValue(t, `{"linkset":[{"next":[{"href":"/one"}],"anchor":"bad path"}]}`)
	if text, err := media.SerializeLinkset(value, 1, 1024); text != "" || !errors.Is(err, media.ErrInvalidLinkset) {
		t.Fatalf("invalid late anchor published text: %q, %v", text, err)
	}
}
