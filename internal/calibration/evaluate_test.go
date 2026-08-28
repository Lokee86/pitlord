package calibration

import (
	"testing"

	"github.com/Lokee86/pitlord/internal/scan"
)

func TestEvaluateScoresRequiredAbsentAllowedAndUnlabelled(t *testing.T) {
	reference := Reference{
		Schema:         Schema,
		Corpus:         "fixture",
		SourceRevision: "abc123",
		Detector:       "dependency-pressure",
		Expectations: []Expectation{
			{ID: "required", Path: "src/pressure.go", Class: "pressure", Finding: FindingRequired, MinSeverity: scan.SeverityWarning},
			{ID: "absent", Path: "src/clean.go", Class: "clean", Finding: FindingAbsent},
			{ID: "allowed", PathPrefix: "src/watch", Class: "watch", Finding: FindingAllowed, MaxSeverity: scan.SeverityWarning},
		},
	}
	result := Evaluate(reference, scan.Result{Findings: []scan.Finding{
		{Detector: "dependency-pressure", Severity: scan.SeverityHigh, Scope: scan.Scope{Path: "src/pressure.go"}},
		{Detector: "dependency-pressure", Severity: scan.SeverityWarning, Scope: scan.Scope{Path: "src/clean.go"}},
		{Detector: "dependency-pressure", Severity: scan.SeverityWarning, Scope: scan.Scope{Path: "src/watch/item.go"}},
		{Detector: "dependency-pressure", Severity: scan.SeverityWarning, Scope: scan.Scope{Path: "src/unlabelled.go"}},
	}})

	if result.Summary.TruePositive != 1 || result.Summary.FalsePositive != 1 {
		t.Fatalf("unexpected score: %+v", result.Summary)
	}
	if result.Summary.FalseNegative != 0 || result.Summary.TrueNegative != 0 {
		t.Fatalf("unexpected negative score: %+v", result.Summary)
	}
	if result.Summary.UnlabelledFindings != 1 {
		t.Fatalf("unlabelled findings = %d, want 1", result.Summary.UnlabelledFindings)
	}
	if result.Summary.SeverityMismatches != 0 {
		t.Fatalf("severity mismatches = %d, want 0", result.Summary.SeverityMismatches)
	}
}

func TestEvaluateReportsFalseNegativeAndSeverityMismatch(t *testing.T) {
	reference := Reference{
		Schema:         Schema,
		Corpus:         "fixture",
		SourceRevision: "abc123",
		Detector:       "dependency-pressure",
		Expectations: []Expectation{
			{ID: "missing", Path: "src/missing.go", Class: "pressure", Finding: FindingRequired},
			{ID: "too-high", Path: "src/watch.go", Class: "watch", Finding: FindingAllowed, MaxSeverity: scan.SeverityWarning},
		},
	}
	result := Evaluate(reference, scan.Result{Findings: []scan.Finding{
		{Detector: "dependency-pressure", Severity: scan.SeverityHigh, Scope: scan.Scope{Path: "src/watch.go"}},
	}})

	if result.Summary.FalseNegative != 1 {
		t.Fatalf("false negatives = %d, want 1", result.Summary.FalseNegative)
	}
	if result.Summary.SeverityMismatches != 1 {
		t.Fatalf("severity mismatches = %d, want 1", result.Summary.SeverityMismatches)
	}
	if !result.HasMismatch() {
		t.Fatal("expected mismatch")
	}
}
