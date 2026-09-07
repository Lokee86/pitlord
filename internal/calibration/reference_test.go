package calibration

import (
	"strings"
	"testing"
)

func TestDecodeReferenceRejectsAmbiguousPathMatcher(t *testing.T) {
	_, err := DecodeReference(strings.NewReader(`{
		"schema":"pitlord.calibration.v1",
		"corpus":"fixture",
		"source_revision":"abc",
		"detector":"dependency-pressure",
		"expectations":[{
			"id":"bad",
			"path":"src/a.go",
			"path_prefix":"src",
			"class":"clean",
			"finding":"absent"
		}]
	}`))
	if err == nil || !strings.Contains(err.Error(), "exactly one") {
		t.Fatalf("expected path matcher validation error, got %v", err)
	}
}

func TestDecodeReferenceAcceptsRootPrefixAndWorktreeDiffHash(t *testing.T) {
	reference, err := DecodeReference(strings.NewReader(`{
		"schema":"pitlord.calibration.v1",
		"corpus":"fixture",
		"source_revision":"abc",
		"worktree_diff_sha256":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"detector":"symbol-intermediary-bypass",
		"expectations":[{"id":"clean","path_prefix":".","class":"clean","finding":"absent"}]
	}`))
	if err != nil {
		t.Fatal(err)
	}
	if reference.WorktreeDiffSHA256 == "" || reference.Expectations[0].PathPrefix != "." {
		t.Fatalf("unexpected reference: %+v", reference)
	}
}

func TestDecodeReferenceRejectsInvalidWorktreeDiffHash(t *testing.T) {
	_, err := DecodeReference(strings.NewReader(`{
		"schema":"pitlord.calibration.v1",
		"corpus":"fixture",
		"source_revision":"abc",
		"worktree_diff_sha256":"not-a-sha256",
		"detector":"symbol-intermediary-bypass",
		"expectations":[{"id":"clean","path_prefix":".","class":"clean","finding":"absent"}]
	}`))
	if err == nil || !strings.Contains(err.Error(), "worktree_diff_sha256") {
		t.Fatalf("expected worktree diff hash validation error, got %v", err)
	}
}

func TestDecodeReferenceAcceptsAnalyzerRuleLanguageAndLocation(t *testing.T) {
	reference, err := DecodeReference(strings.NewReader(`{
		"schema":"pitlord.calibration.v1",
		"corpus":"fixture",
		"source_revision":"abc",
		"analyzer":"swallowed-error",
		"rule_id":"swallowed-error",
		"language":"python",
		"expectations":[{
			"id":"semantic-hit",
			"path":"src/app.py",
			"location":{"start_line":12,"start_column":5,"end_line":13,"end_column":9},
			"class":"swallowed-error",
			"finding":"required"
		}]
	}`))
	if err != nil {
		t.Fatal(err)
	}
	if reference.Analyzer != "swallowed-error" || reference.RuleID != "swallowed-error" || reference.Language != "python" {
		t.Fatalf("unexpected semantic selector: %+v", reference)
	}
	if reference.Expectations[0].Location == nil || reference.Expectations[0].Location.StartLine != 12 {
		t.Fatalf("unexpected location selector: %+v", reference.Expectations[0])
	}
}

func TestDecodeReferenceRejectsDetectorAndAnalyzerTogether(t *testing.T) {
	_, err := DecodeReference(strings.NewReader(`{
		"schema":"pitlord.calibration.v1",
		"corpus":"fixture",
		"source_revision":"abc",
		"detector":"dependency-pressure",
		"analyzer":"swallowed-error",
		"expectations":[{"id":"x","path":"src/a.go","class":"clean","finding":"absent"}]
	}`))
	if err == nil || !strings.Contains(err.Error(), "exactly one of detector or analyzer") {
		t.Fatalf("expected target validation error, got %v", err)
	}
}

func TestDecodeReferenceRejectsLocationWithPrefix(t *testing.T) {
	_, err := DecodeReference(strings.NewReader(`{
		"schema":"pitlord.calibration.v1",
		"corpus":"fixture",
		"source_revision":"abc",
		"analyzer":"swallowed-error",
		"expectations":[{
			"id":"bad-location",
			"path_prefix":"src",
			"location":{"start_line":12},
			"class":"clean",
			"finding":"absent"
		}]
	}`))
	if err == nil || !strings.Contains(err.Error(), "location requires exact path") {
		t.Fatalf("expected location validation error, got %v", err)
	}
}

func TestDecodeReferenceRejectsUnknownField(t *testing.T) {
	_, err := DecodeReference(strings.NewReader(`{
		"schema":"pitlord.calibration.v1",
		"corpus":"fixture",
		"source_revision":"abc",
		"detector":"dependency-pressure",
		"unexpected":true,
		"expectations":[{"id":"x","path":"src/a.go","class":"clean","finding":"absent"}]
	}`))
	if err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("expected strict decode error, got %v", err)
	}
}
