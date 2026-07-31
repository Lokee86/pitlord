package policy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestRepositoryRulesEvaluateContentAndPaths(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "client/scripts/nested/bad.gd", "safe\nforbidden token\n")
	writeTestFile(t, root, "client/scripts/generated/allowed.gd", "forbidden token\n")
	writeTestFile(t, root, "required/control.go", "package required\n")
	writeTestFile(t, root, "forbidden/debug_old.go", "package forbidden\n")

	document := Document{
		Version: Version,
		Rules: []Rule{
			{
				ID:           "content",
				Type:         RuleForbidContent,
				Literal:      "forbidden token",
				IncludePaths: []string{"client/scripts/**/*.gd"},
				ExcludePaths: []string{"client/scripts/generated/**"},
			},
			{ID: "required", Type: RuleRequirePath, Path: "required/control.go"},
			{ID: "forbidden", Type: RuleForbidPath, Path: "forbidden/debug_*.go"},
		},
	}
	if err := document.NormalizeAndValidate(); err != nil {
		t.Fatalf("normalize policy: %v", err)
	}

	diagnostics := Evaluate(document, arcana.Graph{}, root)
	if len(diagnostics) != 2 {
		t.Fatalf("expected 2 diagnostics, got %#v", diagnostics)
	}
	if diagnostics[0].RuleID != "content" || diagnostics[1].RuleID != "forbidden" {
		t.Fatalf("unexpected diagnostics: %#v", diagnostics)
	}
	content := diagnostics[0].Evidence
	if len(content) != 1 || content[0].Source.Path != "client/scripts/nested/bad.gd" {
		t.Fatalf("unexpected content evidence: %#v", content)
	}
	if content[0].Source.Span == nil || content[0].Source.Span.StartLine != 2 {
		t.Fatalf("unexpected content span: %#v", content[0].Source.Span)
	}
	if diagnostics[1].Evidence[0].Source.Path != "forbidden/debug_old.go" {
		t.Fatalf("unexpected path evidence: %#v", diagnostics[1].Evidence)
	}
}

func TestRepositoryRulesReportMissingPathAndRegex(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "auth/store.gd", "var value = JSON.stringify({\"token\": token})\n")

	document := Document{
		Version: Version,
		Rules: []Rule{
			{
				ID:           "regex",
				Type:         RuleForbidContent,
				Regex:        `JSON\.stringify\([^\n]*token`,
				IncludePaths: []string{"auth/*.gd"},
			},
			{ID: "missing", Type: RuleRequirePath, Path: "required/*.go"},
		},
	}
	if err := document.NormalizeAndValidate(); err != nil {
		t.Fatalf("normalize policy: %v", err)
	}

	diagnostics := Evaluate(document, arcana.Graph{}, root)
	if len(diagnostics) != 2 {
		t.Fatalf("expected 2 diagnostics, got %#v", diagnostics)
	}
	if diagnostics[0].RuleID != "missing" || diagnostics[0].Evidence[0].Issue != "missing_path" {
		t.Fatalf("unexpected missing-path diagnostic: %#v", diagnostics[0])
	}
	if diagnostics[1].RuleID != "regex" || diagnostics[1].Evidence[0].Source.Path != "auth/store.gd" {
		t.Fatalf("unexpected regex diagnostic: %#v", diagnostics[1])
	}
}

func TestRepositoryRulesSkipNestedWorktreesAndCaches(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, ".worktrees/branch/bad.go", "forbidden\n")
	writeTestFile(t, root, ".workingtrees/legacy/bad.go", "forbidden\n")
	writeTestFile(t, root, "node_modules/package/bad.go", "forbidden\n")
	writeTestFile(t, root, ".grimoire/knowledge/bad.go", "forbidden\n")
	writeTestFile(t, root, ".arcana/snapshots/bad.go", "forbidden\n")
	writeTestFile(t, root, "source/good.go", "safe\n")

	document := Document{
		Version: Version,
		Rules: []Rule{{
			ID:           "content",
			Type:         RuleForbidContent,
			Literal:      "forbidden",
			IncludePaths: []string{"**/*.go"},
		}},
	}
	if err := document.NormalizeAndValidate(); err != nil {
		t.Fatalf("normalize policy: %v", err)
	}
	if diagnostics := Evaluate(document, arcana.Graph{}, root); len(diagnostics) != 0 {
		t.Fatalf("ignored directories produced diagnostics: %#v", diagnostics)
	}
}

func TestRequireContentReportsEachSelectedFileWithoutMatch(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "AGENTS.md", "Documentation is part of the implementation.\n")
	writeTestFile(t, root, "nested/AGENTS.md", "No documentation rule yet.\n")

	document := Document{
		Version: Version,
		Rules: []Rule{{
			ID:           "required-content",
			Type:         RuleRequireContent,
			Literal:      "Documentation is part of the implementation",
			IncludePaths: []string{"**/AGENTS.md"},
		}},
	}
	if err := document.NormalizeAndValidate(); err != nil {
		t.Fatalf("normalize policy: %v", err)
	}

	diagnostics := Evaluate(document, arcana.Graph{}, root)
	if len(diagnostics) != 1 {
		t.Fatalf("expected 1 diagnostic, got %#v", diagnostics)
	}
	if diagnostics[0].RuleID != "required-content" || len(diagnostics[0].Evidence) != 1 {
		t.Fatalf("unexpected diagnostic: %#v", diagnostics[0])
	}
	if diagnostics[0].Evidence[0].Issue != "missing_content" || diagnostics[0].Evidence[0].Source.Path != "nested/AGENTS.md" {
		t.Fatalf("unexpected evidence: %#v", diagnostics[0].Evidence)
	}
}

func TestRequireContentReportsWhenNoSelectedFileExists(t *testing.T) {
	root := t.TempDir()
	document := Document{
		Version: Version,
		Rules: []Rule{{
			ID:           "required-content",
			Type:         RuleRequireContent,
			Literal:      "required phrase",
			IncludePaths: []string{"docs/**/*.md"},
		}},
	}
	if err := document.NormalizeAndValidate(); err != nil {
		t.Fatalf("normalize policy: %v", err)
	}

	diagnostics := Evaluate(document, arcana.Graph{}, root)
	if len(diagnostics) != 1 || diagnostics[0].Evidence[0].Source.Path != "docs/**/*.md" {
		t.Fatalf("unexpected diagnostics: %#v", diagnostics)
	}
}

func writeTestFile(t *testing.T, root, relative, contents string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
}
