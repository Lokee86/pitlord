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
