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

func TestEvaluateRootPrefixMatchesAllRepositoryFindings(t *testing.T) {
	reference := Reference{
		Schema:         Schema,
		Corpus:         "fixture",
		SourceRevision: "abc123",
		Detector:       "symbol-intermediary-bypass",
		Expectations: []Expectation{
			{ID: "clean", PathPrefix: ".", Class: "clean", Finding: FindingAbsent},
		},
	}
	result := Evaluate(reference, scan.Result{Findings: []scan.Finding{
		{Detector: "symbol-intermediary-bypass", Severity: scan.SeverityWarning, Scope: scan.Scope{Path: "src/nested/file.go"}},
	}})
	if result.Summary.FalsePositive != 1 || result.Summary.UnlabelledFindings != 0 {
		t.Fatalf("root prefix did not cover repository finding: %+v", result.Summary)
	}
}

func TestEvaluateFiltersAnalyzerRuleLanguageAndExactLocation(t *testing.T) {
	reference := Reference{
		Schema:         Schema,
		Corpus:         "fixture",
		SourceRevision: "abc123",
		Analyzer:       "swallowed-error",
		RuleID:         "swallowed-error",
		Language:       "python",
		Expectations: []Expectation{
			{
				ID:       "python-swallowed",
				Path:     "src/app.py",
				Location: &LocationExpectation{StartLine: 12, StartColumn: 5, EndLine: 13, EndColumn: 9},
				Class:    "swallowed-error",
				Finding:  FindingRequired,
			},
		},
	}
	result := Evaluate(reference, scan.Result{Findings: []scan.Finding{
		{
			Analyzer: "swallowed-error", RuleID: "swallowed-error", Detector: "swallowed-error", Language: "python",
			Severity: scan.SeverityWarning, Scope: scan.Scope{Path: "src/app.py"},
			Location: &scan.SourceSpan{Path: "src/app.py", StartLine: 12, StartColumn: 5, EndLine: 13, EndColumn: 9},
		},
		{
			Analyzer: "swallowed-error", RuleID: "swallowed-error", Detector: "swallowed-error", Language: "rust",
			Severity: scan.SeverityWarning, Scope: scan.Scope{Path: "src/app.py"},
			Location: &scan.SourceSpan{Path: "src/app.py", StartLine: 12, StartColumn: 5, EndLine: 13, EndColumn: 9},
		},
		{
			Analyzer: "swallowed-error", RuleID: "swallowed-error", Detector: "swallowed-error", Language: "python",
			Severity: scan.SeverityWarning, Scope: scan.Scope{Path: "src/app.py"},
			Location: &scan.SourceSpan{Path: "src/app.py", StartLine: 20, StartColumn: 5, EndLine: 21, EndColumn: 9},
		},
	}})
	if result.Summary.TruePositive != 1 || result.Summary.UnlabelledFindings != 1 {
		t.Fatalf("unexpected semantic score: %+v", result.Summary)
	}
	if result.Analyzer != "swallowed-error" || result.Language != "python" {
		t.Fatalf("semantic selectors not preserved: %+v", result)
	}
}

func TestEvaluateStrictReferenceTreatsUnlabelledFindingAsMismatch(t *testing.T) {
	reference := Reference{
		Schema:               Schema,
		Corpus:               "fixture",
		SourceRevision:       "abc123",
		Analyzer:             "swallowed-error",
		Language:             "python",
		RequireFullyLabelled: true,
		Expectations: []Expectation{
			{ID: "known", Path: "src/app.py", Class: "swallowed-error", Finding: FindingRequired},
		},
	}
	result := Evaluate(reference, scan.Result{Findings: []scan.Finding{
		{Analyzer: "swallowed-error", Detector: "swallowed-error", Language: "python", Scope: scan.Scope{Path: "src/app.py"}},
		{Analyzer: "swallowed-error", Detector: "swallowed-error", Language: "python", Scope: scan.Scope{Path: "src/other.py"}},
	}})
	if result.Summary.UnlabelledFindings != 1 || !result.HasMismatch() {
		t.Fatalf("strict reference should fail on unlabelled finding: %+v", result.Summary)
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
