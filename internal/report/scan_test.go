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
	if !strings.Contains(output.String(), "no generalized architecture findings") {
		t.Fatalf("unexpected output: %s", output.String())
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
