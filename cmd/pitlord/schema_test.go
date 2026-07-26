package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestSchemaCommandWritesJSON(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"schema", "--kind", "policy"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected success, got %d: %s", code, stderr.String())
	}
	var decoded map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid schema output: %v", err)
	}
	if decoded["title"] != "Pitlord Policy v1" {
		t.Fatalf("unexpected schema title: %v", decoded["title"])
	}
}
