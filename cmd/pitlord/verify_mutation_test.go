package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestVerifyMutationRequiresManifestAndBaseline(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"verify-mutation"}, &stdout, &stderr)
	if code != 2 || !strings.Contains(stderr.String(), "required") {
		t.Fatalf("unexpected result %d: %s", code, stderr.String())
	}
}
