package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
	"github.com/Lokee86/pitlord/internal/policy"
)

func TestWriteSARIFEmitsOneResultPerEvidence(t *testing.T) {
	target := arcana.Node{Identity: "target", Path: "storage/write.go", Name: "Write"}
	value := policy.Report{
		Snapshot:     "snapshot",
		PolicySource: "pitlord.json",
		Diagnostics: []policy.Diagnostic{
			{
				RuleID:   "api-must-not-access-storage",
				Message:  "api must not access storage",
				Severity: "error",
				Evidence: []policy.Evidence{
					{
						Issue:    "forbidden_dependency",
						Source:   arcana.Node{Identity: "source", Path: "api/create.go", Name: "Create", Span: &arcana.Span{Path: "api/create.go", StartLine: 7, StartColumn: 2}},
						Relation: "calls",
						Target:   &target,
					},
				},
			},
		},
	}
	var output bytes.Buffer
	if err := WriteSARIF(&output, value, "0.1.0-test"); err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(output.Bytes(), &decoded); err != nil {
		t.Fatal(err)
	}
	runs := decoded["runs"].([]any)
	run := runs[0].(map[string]any)
	results := run["results"].([]any)
	if len(results) != 1 {
		t.Fatalf("expected one result, got %d", len(results))
	}
	result := results[0].(map[string]any)
	if result["ruleId"] != "api-must-not-access-storage" {
		t.Fatalf("unexpected rule id: %v", result["ruleId"])
	}
	fingerprints := result["partialFingerprints"].(map[string]any)
	if fingerprints["pitlordEvidenceFingerprint/v1"] == "" {
		t.Fatal("expected stable evidence fingerprint")
	}
}
