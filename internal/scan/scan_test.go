package scan

import "testing"

func TestRunCreatesPolicyFreeEmptyResult(t *testing.T) {
	result, err := Run(Input{SnapshotPath: "snapshot", PathPrefix: "."})
	if err != nil {
		t.Fatal(err)
	}
	if result.Schema != Schema {
		t.Fatalf("unexpected schema %q", result.Schema)
	}
	if result.Scope.Kind != "repository" || result.Scope.Path != "." {
		t.Fatalf("unexpected scope: %+v", result.Scope)
	}
	if result.Findings == nil || len(result.Findings) != 0 {
		t.Fatalf("expected non-nil empty findings, got %#v", result.Findings)
	}
	if result.Summary.FindingCount != 0 {
		t.Fatalf("unexpected summary: %+v", result.Summary)
	}
}

func TestRunNormalizesRepositoryRelativeScope(t *testing.T) {
	result, err := Run(Input{SnapshotPath: "snapshot", PathPrefix: "./src/service/"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Scope.Path != "src/service" {
		t.Fatalf("unexpected normalized path %q", result.Scope.Path)
	}
}

func TestRunRejectsMissingSnapshotAndEscapingScope(t *testing.T) {
	if _, err := Run(Input{}); err == nil {
		t.Fatal("expected missing snapshot failure")
	}
	if _, err := Run(Input{SnapshotPath: "snapshot", PathPrefix: "../outside"}); err == nil {
		t.Fatal("expected escaping scope failure")
	}
}

func TestFinalizeOrdersFindingsAndRecalculatesSummary(t *testing.T) {
	result := finalize(Result{Findings: []Finding{
		{ID: "b", Detector: "hub", Disposition: DispositionAdvisory, Severity: SeverityWarning, Scope: Scope{Path: "z"}},
		{ID: "a", Detector: "cycle", Disposition: DispositionGuard, Severity: SeverityWarning, Scope: Scope{Path: "a"}},
		{ID: "c", Detector: "hub", Disposition: DispositionAdvisory, Severity: SeverityCritical, Scope: Scope{Path: "a"}},
	}})
	if result.Findings[0].ID != "a" || result.Findings[1].ID != "c" || result.Findings[2].ID != "b" {
		t.Fatalf("unexpected finding order: %#v", result.Findings)
	}
	if result.Summary.FindingCount != 3 || result.Summary.Guard != 1 || result.Summary.Advisory != 2 || result.Summary.Critical != 1 || result.Summary.Warnings != 2 {
		t.Fatalf("unexpected summary: %+v", result.Summary)
	}
	for _, finding := range result.Findings {
		if finding.Evidence == nil {
			t.Fatalf("finding evidence should be non-nil: %#v", finding)
		}
	}
}
