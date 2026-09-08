package main

import (
	"bytes"
	"os"
	"path/filepath"
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

func TestPublicCommandsAreDocumented(t *testing.T) {
	commands := []string{
		"check", "baseline", "validate", "schema", "verify-mutation", "diff",
		"analyze", "scan", "docs", "calibrate", "inspect", "generate", "version",
	}
	for _, relative := range []string{"../../README.md", "../../docs/COMMANDS.md"} {
		data, err := os.ReadFile(filepath.Clean(relative))
		if err != nil {
			t.Fatal(err)
		}
		content := string(data)
		for _, command := range commands {
			if !strings.Contains(content, "`"+command+"`") && !strings.Contains(content, "pitlord "+command) {
				t.Errorf("%s does not document public command %q", relative, command)
			}
		}
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
