package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Lokee86/pitlord/internal/scan"
)

func TestWriteScanTextReportsEmptyResult(t *testing.T) {
	result := scan.Result{Schema: scan.Schema, Scope: scan.Scope{Kind: "repository", Path: "."}, Findings: []scan.Finding{}}
	var output bytes.Buffer
	if err := WriteScanText(&output, result); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "no findings") {
		t.Fatalf("unexpected output: %s", output.String())
	}
}

func TestWriteScanTextIncludesSourceDiagnosticMetadata(t *testing.T) {
	result := scan.Result{Schema: scan.Schema, Scope: scan.Scope{Kind: "repository", Path: "."}, Findings: []scan.Finding{{
		ID:          "lint-1",
		RuleID:      "needless-clone",
		Analyzer:    "rust-lints",
		Detector:    "rust-lints",
		Disposition: scan.DispositionAdvisory,
		Severity:    scan.SeverityWarning,
		Scope:       scan.Scope{Kind: "file", Path: "src/lib.rs"},
		Location:    &scan.SourceSpan{Path: "src/lib.rs", StartLine: 7, StartColumn: 3},
		Summary:     "Clone is unnecessary",
		SuggestedFix: &scan.SuggestedFix{
			Message:       "remove the clone",
			Applicability: scan.ApplicabilityMachine,
		},
	}}}
	var output bytes.Buffer
	if err := WriteScanText(&output, result); err != nil {
		t.Fatal(err)
	}
	text := output.String()
	if !strings.Contains(text, "needless-clone") || !strings.Contains(text, "src/lib.rs:7:3") || !strings.Contains(text, "Fix (machine-applicable): remove the clone") {
		t.Fatalf("unexpected output: %s", text)
	}
}

func TestWriteScanJSONPreservesEmptyFindingsArray(t *testing.T) {
	result := scan.Result{Schema: scan.Schema, Scope: scan.Scope{Kind: "repository", Path: "."}, Findings: []scan.Finding{}}
	var output bytes.Buffer
	if err := WriteScanJSON(&output, result); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), `"findings": []`) {
		t.Fatalf("unexpected output: %s", output.String())
	}
}
