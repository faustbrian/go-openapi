package expression_test

import (
	"errors"
	"testing"

	"github.com/faustbrian/go-openapi/v2/expression"
	"github.com/faustbrian/go-openapi/v2/jsonvalue"
)

func TestTemplatePropagatesNonStringMarshalFailure(t *testing.T) {
	template, err := expression.ParseTemplate("before-{$request.body}-after")
	if err != nil {
		t.Fatal(err)
	}
	value := jsonvalue.Null()
	for range jsonvalue.DefaultMarshalLimits().MaxDepth {
		value, err = jsonvalue.Array([]jsonvalue.Value{value})
		if err != nil {
			t.Fatal(err)
		}
	}
	if encoded, err := value.MarshalJSON(); encoded != nil || !errors.Is(err, jsonvalue.ErrMarshalLimit) {
		t.Fatalf("independent marshal control = %q, %v", encoded, err)
	}
	text, err := template.Evaluate(expression.Context{Request: expression.Message{Body: value}}, expression.EvaluationOptions{})
	if text != "" || !errors.Is(err, jsonvalue.ErrMarshalLimit) {
		t.Fatalf("failed interpolation = %q, %v", text, err)
	}
	text, err = template.Evaluate(expression.Context{Request: expression.Message{Body: jsonvalue.Null()}}, expression.EvaluationOptions{})
	if text != "before-null-after" || err != nil {
		t.Fatalf("ordinary interpolation = %q, %v", text, err)
	}
}
