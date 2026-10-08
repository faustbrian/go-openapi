package main

import (
	"encoding/json"
	"testing"

	"github.com/go-openapi/loads"
)

func TestLoadsRetainsSwaggerSemanticsAcrossReload(t *testing.T) {
	t.Parallel()
	const input = `{"swagger":"2.0","info":{"title":"sentinel","version":"1"},"paths":{"/widgets":{"get":{"operationId":"readWidgets","x-sentinel":"operation-retained","responses":{"200":{"description":"found","schema":{"$ref":"#/definitions/Widget"}},"404":{"description":"missing"}}}}},"definitions":{"Widget":{"type":"object","required":["code"],"properties":{"code":{"type":"string","default":"007"}}}}}`
	document, err := loads.Analyzed(json.RawMessage(input), "2.0")
	if err != nil {
		t.Fatal(err)
	}
	for pass := 0; pass < 2; pass++ {
		rendered, err := json.Marshal(document.Spec())
		if err != nil {
			t.Fatal(err)
		}
		var got struct {
			Paths map[string]struct {
				Get struct {
					OperationID string `json:"operationId"`
					Sentinel    string `json:"x-sentinel"`
					Responses   map[string]struct {
						Description string
						Schema      struct {
							Ref string `json:"$ref"`
						}
					}
				}
			}
			Definitions map[string]struct {
				Type       string
				Required   []string
				Properties map[string]struct {
					Type    string
					Default any
				}
			}
		}
		if err := json.Unmarshal(rendered, &got); err != nil {
			t.Fatal(err)
		}
		operation := got.Paths["/widgets"].Get
		if operation.OperationID != "readWidgets" || operation.Sentinel != "operation-retained" {
			t.Fatalf("operation lost: %#v", operation)
		}
		if len(operation.Responses) != 2 || operation.Responses["200"].Description != "found" || operation.Responses["404"].Description != "missing" || operation.Responses["200"].Schema.Ref != "#/definitions/Widget" {
			t.Fatalf("response data lost: %#v", operation.Responses)
		}
		definition := got.Definitions["Widget"]
		if definition.Type != "object" || len(definition.Required) != 1 || definition.Required[0] != "code" || definition.Properties["code"].Type != "string" || definition.Properties["code"].Default != "007" {
			t.Fatalf("definition data lost: %#v", definition)
		}
		document, err = loads.Analyzed(rendered, "2.0")
		if err != nil {
			t.Fatalf("reload: %v", err)
		}
	}
	if _, err := loads.Analyzed(json.RawMessage(`{"swagger":`), "2.0"); err == nil {
		t.Fatal("malformed JSON accepted")
	}
}
