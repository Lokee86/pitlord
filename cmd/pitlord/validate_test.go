package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidatePolicy(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pitlord.json")
	content := `{
  "version": 1,
  "areas": [
    {"id": "api", "paths": ["api"]},
    {"id": "storage", "paths": ["storage"]}
  ],
  "rules": [
    {"id": "rule", "from_areas": ["api"], "to_areas": ["storage"]}
  ]
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{"validate", "--policy", path}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected success, got %d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "2 area(s), 1 rule(s)") {
		t.Fatalf("unexpected output: %s", stdout.String())
	}
}

func TestValidateRequiresOneInput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"validate"}, &stdout, &stderr)
	if code != 2 || !strings.Contains(stderr.String(), "exactly one") {
		t.Fatalf("unexpected result %d: %s", code, stderr.String())
	}
}
