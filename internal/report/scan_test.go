package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Lokee86/pitlord/internal/scan"
)

func TestWriteScanTextReportsEmptyResult(t *testing.T) {
	result, err := scan.Run(scan.Input{SnapshotPath: "snapshot"})
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := WriteScanText(&output, result); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "no generalized architecture findings") {
		t.Fatalf("unexpected output: %s", output.String())
	}
}

func TestWriteScanJSONPreservesEmptyFindingsArray(t *testing.T) {
	result, err := scan.Run(scan.Input{SnapshotPath: "snapshot"})
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := WriteScanJSON(&output, result); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), `"findings": []`) {
		t.Fatalf("unexpected output: %s", output.String())
	}
}
