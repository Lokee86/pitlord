package schema

import (
	"encoding/json"
	"testing"
)

func TestEmbeddedSchemasAreValidJSON(t *testing.T) {
	for _, kind := range []string{"policy", "baseline"} {
		data, err := Read(kind)
		if err != nil {
			t.Fatal(err)
		}
		var decoded map[string]any
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("%s schema is invalid JSON: %v", kind, err)
		}
		if decoded["$schema"] == nil {
			t.Fatalf("%s schema has no $schema", kind)
		}
	}
}

func TestReadRejectsUnknownSchema(t *testing.T) {
	if _, err := Read("unknown"); err == nil {
		t.Fatal("expected unknown schema to fail")
	}
}
