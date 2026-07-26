package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"version"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected success, got %d: %s", code, stderr.String())
	}
	if strings.TrimSpace(stdout.String()) != version {
		t.Fatalf("unexpected version output %q", stdout.String())
	}
}

func TestSplitCommaList(t *testing.T) {
	values := splitCommaList(" calls, imports ,, references ")
	if len(values) != 3 || values[0] != "calls" || values[2] != "references" {
		t.Fatalf("unexpected values: %v", values)
	}
}

func TestCheckRequiresOnePolicySource(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"check"}, &stdout, &stderr); code != 2 {
		t.Fatalf("expected usage failure, got %d", code)
	}
	if !strings.Contains(stderr.String(), "exactly one") {
		t.Fatalf("unexpected error: %s", stderr.String())
	}
}
