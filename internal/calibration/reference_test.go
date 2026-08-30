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
