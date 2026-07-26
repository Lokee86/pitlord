package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestDiffRequiresBeforeSnapshotAndPolicy(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"diff"}, &stdout, &stderr)
	if code != 2 || !strings.Contains(stderr.String(), "required") {
		t.Fatalf("unexpected result %d: %s", code, stderr.String())
	}
}

func TestDiffRejectsUnknownFormatBeforeEvaluation(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run(
		[]string{"diff", "--before-snapshot", "before", "--policy", "policy.json", "--format", "xml"},
		&stdout,
		&stderr,
	)
	if code != 2 || !strings.Contains(stderr.String(), "unsupported format") {
		t.Fatalf("unexpected result %d: %s", code, stderr.String())
	}
}
