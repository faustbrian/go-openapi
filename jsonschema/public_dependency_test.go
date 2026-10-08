package jsonschema_test

import (
	"context"
	"errors"
	"testing"

	canonical "github.com/faustbrian/go-json-schema/v2"
	openapischema "github.com/faustbrian/go-openapi/v2/jsonschema"
)

func TestCompilerUsesPublicUnicodeGrammar(t *testing.T) {
	for _, dialect := range []openapischema.Dialect{
		openapischema.DialectOAS30, openapischema.DialectOAS31, openapischema.DialectOAS32,
	} {
		t.Run(string(dialect), func(t *testing.T) {
			compiler, err := openapischema.NewCompiler(dialect)
			if err != nil {
				t.Fatal(err)
			}
			compiled, err := compiler.Compile(t.Context(), mustValue(t,
				`{"type":"string","pattern":"^\\p{Script=Greek}+$"}`))
			if err != nil {
				t.Fatalf("standard Unicode grammar compile: %v", err)
			}
			for _, test := range []struct {
				instance string
				valid    bool
			}{
				{`"αβ"`, true}, {`"AB"`, false},
			} {
				result, err := compiled.Validate(t.Context(), []byte(test.instance))
				if err != nil || result.Valid != test.valid {
					t.Fatalf("Unicode validity = %t, error = %v; want %t", result.Valid, err, test.valid)
				}
			}
			if _, err := compiler.Compile(t.Context(), mustValue(t,
				`{"pattern":"\\p{Greek}"}`)); !errors.Is(err, canonical.ErrInvalidSchema) {
				t.Fatalf("nonstandard Unicode grammar error = %v", err)
			}
			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			if _, err := compiler.Compile(ctx, mustValue(t, `{"type":"string"}`)); !errors.Is(err, context.Canceled) {
				t.Fatalf("canceled compile error = %v", err)
			}
		})
	}
}

// OpenAPI does not implicitly assert formats. An explicitly configured
// canonical compiler yields the same public Schema alias with format checks.
func TestPublicSchemaAliasUsesExplicitURITemplateGrammar(t *testing.T) {
	compiler, err := canonical.NewCompiler(canonical.WithFormatAssertion())
	if err != nil {
		t.Fatal(err)
	}
	var compiled *openapischema.Schema
	compiled, err = compiler.Compile(t.Context(), []byte(`{"format":"uri-template"}`))
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		instance string
		valid    bool
	}{
		{`"{var:1}"`, true}, {`"{var*}"`, true}, {`"{var:1*}"`, false},
	} {
		result, err := compiled.Validate(t.Context(), []byte(test.instance))
		if err != nil || result.Valid != test.valid {
			t.Fatalf("URI-template validity = %t, error = %v; want %t", result.Valid, err, test.valid)
		}
	}
}
