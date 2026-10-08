package main

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"testing"

	openapi "github.com/faustbrian/go-openapi"
	"github.com/faustbrian/go-openapi/parse"
	"github.com/faustbrian/go-openapi/serialize"
	"github.com/pb33f/libopenapi"
)

func TestLibopenapiRoundTripRetainsDocumentSemantics(t *testing.T) {
	fixtures := map[string]string{
		"json": `{
  "openapi": "3.1.0",
  "info": {"title": "001", "version": "1.0"},
  "paths": {"/widgets": {"get": {
    "operationId": "fetchWidgets",
    "responses": {
      "200": {"description": "kept response", "content": {
        "application/json": {"schema": {"$ref": "#/components/schemas/Result"}}
      }},
      "201": {"description": "second response"}
    }
  }}},
  "components": {"schemas": {
    "Result": {"type": "object", "properties": {
      "marker": {"type": "string", "default": "0123"},
      "count": {"type": "integer"}
    }},
    "Always": true
  }}
}`,
		"yaml": `openapi: 3.1.0
info:
  title: "001"
  version: "1.0"
paths:
  /widgets:
    get:
      operationId: fetchWidgets
      responses:
        "200":
          description: kept response
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Result'
        "201":
          description: second response
components:
  schemas:
    Result:
      type: object
      properties:
        marker:
          type: string
          default: "0123"
        count:
          type: integer
    Always: true
`,
	}
	for name, fixture := range fixtures {
		t.Run(name, func(t *testing.T) {
			document, err := libopenapi.NewDocument([]byte(fixture))
			if err != nil {
				t.Fatalf("NewDocument() error = %v", err)
			}
			if _, err := document.BuildV3Model(); err != nil {
				t.Fatalf("BuildV3Model() error = %v", err)
			}
			rendered, err := document.Render()
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}
			reloaded, err := libopenapi.NewDocument(rendered)
			if err != nil {
				t.Fatalf("reload error = %v", err)
			}
			if _, err := reloaded.BuildV3Model(); err != nil {
				t.Fatalf("reloaded BuildV3Model() error = %v", err)
			}

			// Decode through the independently published owned parser, not the
			// peer's model or renderer, and compare explicit input sentinels.
			owned, err := openapi.ParseYAML(
				context.Background(), bytes.NewReader(rendered), parse.DefaultLimits(),
			)
			if err != nil {
				t.Fatalf("owned ParseYAML() error = %v", err)
			}
			var encoded bytes.Buffer
			if err := serialize.JSON(
				context.Background(), &encoded, owned, serialize.DefaultOptions(),
			); err != nil {
				t.Fatalf("owned serialize.JSON() error = %v", err)
			}
			var decoded map[string]any
			if err := json.Unmarshal(encoded.Bytes(), &decoded); err != nil {
				t.Fatalf("independent JSON decode error = %v", err)
			}
			assertRetainedValue(t, decoded, "001", "info", "title")
			assertRetainedValue(t, decoded, "fetchWidgets", "paths", "/widgets", "get", "operationId")
			assertRetainedValue(t, decoded, "kept response", "paths", "/widgets", "get", "responses", "200", "description")
			assertRetainedValue(t, decoded, "second response", "paths", "/widgets", "get", "responses", "201", "description")
			assertRetainedValue(t, decoded, "#/components/schemas/Result", "paths", "/widgets", "get", "responses", "200", "content", "application/json", "schema", "$ref")
			assertRetainedValue(t, decoded, "object", "components", "schemas", "Result", "type")
			assertRetainedValue(t, decoded, "0123", "components", "schemas", "Result", "properties", "marker", "default")
			assertRetainedValue(t, decoded, "integer", "components", "schemas", "Result", "properties", "count", "type")
			// Both peer versions normalize the permissive boolean schema to an
			// empty schema. Assert its unrestricted semantics, not its spelling.
			assertRetainedValue(t, decoded, map[string]any{}, "components", "schemas", "Always")
		})
	}
}

func assertRetainedValue(t *testing.T, document map[string]any, want any, path ...string) {
	t.Helper()
	var value any = document
	for _, key := range path {
		object, ok := value.(map[string]any)
		if !ok {
			t.Fatalf("value at %v is not an object before %q", path, key)
		}
		value, ok = object[key]
		if !ok {
			t.Fatalf("missing value at %v", path)
		}
	}
	if !reflect.DeepEqual(value, want) {
		t.Fatalf("value at %v = %#v, want %#v", path, value, want)
	}
}
