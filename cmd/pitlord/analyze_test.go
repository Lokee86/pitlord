package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestAnalyzeRequiresPolicy(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"analyze"}, &stdout, &stderr)
	if code != 2 || !strings.Contains(stderr.String(), "--policy is required") {
		t.Fatalf("unexpected result %d: %s", code, stderr.String())
	}
}

func TestUniqueStringsSortsAndDeduplicates(t *testing.T) {
	values := uniqueStrings([]string{"b", "a", "b", ""})
	if len(values) != 2 || values[0] != "a" || values[1] != "b" {
		t.Fatalf("unexpected values: %v", values)
	}
}
