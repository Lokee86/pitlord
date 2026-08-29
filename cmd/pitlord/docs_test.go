package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestDocsRequiresChangedFrom(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"docs"}, &stdout, &stderr); code != 2 {
		t.Fatalf("expected usage failure, got %d", code)
	}
	if !strings.Contains(stderr.String(), "--changed-from is required") {
		t.Fatalf("unexpected error: %s", stderr.String())
	}
}
