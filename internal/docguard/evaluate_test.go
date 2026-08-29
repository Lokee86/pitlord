package docguard

import "testing"

func TestEvaluateRequiresPerFileCodemapCoverage(t *testing.T) {
	dataset := Dataset{SchemaVersion: 1, Entries: []DatasetEntry{
		{
			Entry:      Entry{DocumentPath: "docs/runtime.md", Kind: "file"},
			Resolution: Resolution{ResolvedPath: "src/a.go"},
		},
		{
			Entry:      Entry{DocumentPath: "docs/storage.md", Kind: "glob"},
			Resolution: Resolution{Matches: []TargetMatch{{Path: "src/b.go"}}},
		},
		{
			Entry:      Entry{DocumentPath: "docs/service.md", Kind: "symbol"},
			Resolution: Resolution{SemanticNode: &SemanticNode{Path: "src/service.go"}},
		},
		{
			Entry:      Entry{DocumentPath: "docs/package.md", Kind: "directory"},
			Resolution: Resolution{ResolvedPath: "src/pkg"},
		},
	}}

	report := Evaluate([]string{"src/a.go", "src/b.go", "src/service.go", "src/pkg/c.go"}, dataset, nil, "base")
	if report.CodeFiles != 4 || report.MappedCodeFiles != 3 {
		t.Fatalf("unexpected coverage: %#v", report)
	}
	if len(report.Findings) != 1 || report.Findings[0].Code != "missing_codemap" || report.Findings[0].CodePath != "src/pkg/c.go" {
		t.Fatalf("unexpected findings: %#v", report.Findings)
	}
}

func TestEvaluateChangedUnmappedCodeReportsCoverageOnce(t *testing.T) {
	report := Evaluate(
		[]string{"src/unmapped.go"},
		Dataset{SchemaVersion: 1},
		[]Change{{Status: "M", Path: "src/unmapped.go"}},
		"base",
	)
	if report.ChangedCode != 1 {
		t.Fatalf("expected one changed code file, got %#v", report)
	}
	if len(report.Findings) != 1 || report.Findings[0].Code != "missing_codemap" {
		t.Fatalf("expected one coverage finding, got %#v", report.Findings)
	}
}

func TestEvaluateChangedCodeRequiresMappedDocumentChange(t *testing.T) {
	dataset := Dataset{SchemaVersion: 1, Entries: []DatasetEntry{{
		Entry:      Entry{DocumentPath: "docs/runtime.md", Kind: "file"},
		Resolution: Resolution{ResolvedPath: "src/runtime.go"},
	}}}
	changes := []Change{{Status: "M", Path: "src/runtime.go"}}

	report := Evaluate([]string{"src/runtime.go"}, dataset, changes, "base")
	if len(report.Findings) != 1 || report.Findings[0].Code != "documentation_not_changed" {
		t.Fatalf("expected stale documentation finding, got %#v", report.Findings)
	}
	if len(report.Findings[0].Documents) != 1 || report.Findings[0].Documents[0] != "docs/runtime.md" {
		t.Fatalf("unexpected expected documents: %#v", report.Findings[0].Documents)
	}
}

func TestEvaluateChangedMappedDocumentSatisfiesExpectation(t *testing.T) {
	dataset := Dataset{SchemaVersion: 1, Entries: []DatasetEntry{{
		Entry:      Entry{DocumentPath: "docs/runtime.md", Kind: "file"},
		Resolution: Resolution{ResolvedPath: "src/runtime.go"},
	}}}
	changes := []Change{
		{Status: "M", Path: "src/runtime.go"},
		{Status: "M", Path: "docs/runtime.md"},
	}

	report := Evaluate([]string{"src/runtime.go"}, dataset, changes, "base")
	if len(report.Findings) != 0 {
		t.Fatalf("expected guard to pass, got %#v", report.Findings)
	}
	if report.ChangedCode != 1 || report.ChangedDocs != 1 {
		t.Fatalf("unexpected change counts: %#v", report)
	}
}

func TestParseChangesPreservesRenameDestination(t *testing.T) {
	changes := parseChanges([]byte("R100\x00src/old.go\x00src/new.go\x00M\x00docs/runtime.md\x00"))
	if len(changes) != 2 {
		t.Fatalf("unexpected changes: %#v", changes)
	}
	if changes[0].Status != "R100" || changes[0].OldPath != "src/old.go" || changes[0].Path != "src/new.go" {
		t.Fatalf("unexpected rename: %#v", changes[0])
	}
	if changes[1].Status != "M" || changes[1].Path != "docs/runtime.md" {
		t.Fatalf("unexpected modification: %#v", changes[1])
	}
}
