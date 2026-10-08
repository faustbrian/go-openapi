package main

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestKinRetainsDocumentSemanticsAcrossReload(t *testing.T) {
	fixtures := map[string]string{
		"json": `{"openapi":"3.0.3","info":{"title":"sentinel","version":"1"},"paths":{"/widgets":{"get":{"operationId":"readWidgets","x-sentinel":"operation-retained","responses":{"200":{"description":"found","content":{"application/json":{"schema":{"$ref":"#/components/schemas/Widget"}}}},"404":{"description":"missing"}}}}},"components":{"schemas":{"Widget":{"type":"object","required":["code"],"properties":{"code":{"type":"string","default":"007"}}}}}}`,
		"yaml": `openapi: 3.0.3
info:
  title: sentinel
  version: "1"
paths:
  /widgets:
    get:
      operationId: readWidgets
      x-sentinel: operation-retained
      responses:
        "200":
          description: found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Widget'
        "404":
          description: missing
components:
  schemas:
    Widget:
      type: object
      required: [code]
      properties:
        code:
          type: string
          default: "007"
`,
	}
	for name, fixture := range fixtures {
		t.Run(name, func(t *testing.T) {
			loader := openapi3.NewLoader()
			loader.IsExternalRefsAllowed = false
			document, err := loader.LoadFromData([]byte(fixture))
			if err != nil {
				t.Fatal(err)
			}
			for pass := 0; pass < 2; pass++ {
				if err := document.Validate(context.Background()); err != nil {
					t.Fatal(err)
				}
				rendered, err := json.Marshal(document)
				if err != nil {
					t.Fatal(err)
				}
				var decoded map[string]any
				if err := json.Unmarshal(rendered, &decoded); err != nil {
					t.Fatal(err)
				}
				assertRetainedValue(t, decoded, "readWidgets", "paths", "/widgets", "get", "operationId")
				assertRetainedValue(t, decoded, "operation-retained", "paths", "/widgets", "get", "x-sentinel")
				assertRetainedValue(t, decoded, "found", "paths", "/widgets", "get", "responses", "200", "description")
				assertRetainedValue(t, decoded, "missing", "paths", "/widgets", "get", "responses", "404", "description")
				assertRetainedValue(t, decoded, "#/components/schemas/Widget", "paths", "/widgets", "get", "responses", "200", "content", "application/json", "schema", "$ref")
				assertRetainedValue(t, decoded, "object", "components", "schemas", "Widget", "type")
				assertRetainedValue(t, decoded, []any{"code"}, "components", "schemas", "Widget", "required")
				assertRetainedValue(t, decoded, "string", "components", "schemas", "Widget", "properties", "code", "type")
				assertRetainedValue(t, decoded, "007", "components", "schemas", "Widget", "properties", "code", "default")
				document, err = loader.LoadFromData(rendered)
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestKinRetainsBooleanSchemasAcrossReload(t *testing.T) {
	const input = `{"openapi":"3.1.0","info":{"title":"boolean schemas","version":"1"},"paths":{},"components":{"schemas":{"Anything":true,"Nothing":false}}}`
	loader := openapi3.NewLoader()
	document, err := loader.LoadFromData([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	for pass := 0; pass < 2; pass++ {
		if err := document.Validate(context.Background()); err != nil {
			t.Fatal(err)
		}
		rendered, err := json.Marshal(document)
		if err != nil {
			t.Fatal(err)
		}
		var decoded map[string]any
		if err := json.Unmarshal(rendered, &decoded); err != nil {
			t.Fatal(err)
		}
		assertRetainedValue(t, decoded, true, "components", "schemas", "Anything")
		assertRetainedValue(t, decoded, false, "components", "schemas", "Nothing")
		document, err = loader.LoadFromData(rendered)
		if err != nil {
			t.Fatal(err)
		}
	}
	input30 := strings.Replace(input, "3.1.0", "3.0.3", 1)
	if got := runKin("boolean-30.json", []byte(input30)); got != (result{parse: "accept", model: "accept", validate: "reject", roundtrip: "reject"}) {
		t.Fatalf("OpenAPI 3.0 boolean-schema result = %#v", got)
	}
}

func TestKinRejectsMalformedAndExternalReferences(t *testing.T) {
	if got := runKin("malformed.json", []byte(`{"openapi":`)); got != (result{parse: "reject", model: "reject", validate: "na", roundtrip: "na"}) {
		t.Fatalf("malformed result = %#v", got)
	}
	t.Chdir(t.TempDir())
	schemaPath := "widget.json"
	if err := os.WriteFile(schemaPath, []byte(`{"type":"object"}`), 0600); err != nil {
		t.Fatal(err)
	}
	input := []byte(`{"openapi":"3.0.3","info":{"title":"external ref","version":"1"},"paths":{},"components":{"schemas":{"Widget":{"$ref":"widget.json"}}}}`)
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	document, err := loader.LoadFromData(input)
	if err != nil {
		t.Fatalf("enabled-reference control: %v", err)
	}
	if err := document.Validate(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got := runKin("external.json", input); got != (result{parse: "reject", model: "reject", validate: "na", roundtrip: "na"}) {
		t.Fatalf("external-reference result = %#v", got)
	}
}
