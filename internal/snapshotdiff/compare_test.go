package snapshotdiff

import (
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
	"github.com/Lokee86/pitlord/internal/policy"
)

func TestCompareClassifiesEvidenceByFingerprint(t *testing.T) {
	persistent := policy.Evidence{Issue: "unowned", Source: arcana.Node{Identity: "persistent", Path: "src/persistent.go"}}
	resolved := policy.Evidence{Issue: "unowned", Source: arcana.Node{Identity: "resolved", Path: "src/resolved.go"}}
	introduced := policy.Evidence{Issue: "unowned", Source: arcana.Node{Identity: "introduced", Path: "src/introduced.go"}}
	before := []policy.Diagnostic{{RuleID: "ownership", Message: "ownership", Severity: "error", Evidence: []policy.Evidence{persistent, resolved}}}
	after := []policy.Diagnostic{{RuleID: "ownership", Message: "ownership", Severity: "error", Evidence: []policy.Evidence{persistent, introduced}}}

	result := Compare("pitlord.json", "before", "after", before, after)
	if result.Summary.Introduced != 1 || result.Summary.Resolved != 1 || result.Summary.Persistent != 1 {
		t.Fatalf("unexpected summary: %+v", result.Summary)
	}
	if result.Introduced[0].Evidence.Source.Identity != "introduced" {
		t.Fatalf("unexpected introduced finding: %+v", result.Introduced[0])
	}
	if result.Resolved[0].Evidence.Source.Identity != "resolved" {
		t.Fatalf("unexpected resolved finding: %+v", result.Resolved[0])
	}
}

func TestDiagnosticsGroupsFindingsByRule(t *testing.T) {
	findings := []Finding{
		{RuleID: "rule", Message: "message", Severity: "error", Evidence: policy.Evidence{Issue: "unowned", Source: arcana.Node{NodeID: 1}}},
		{RuleID: "rule", Message: "message", Severity: "error", Evidence: policy.Evidence{Issue: "unowned", Source: arcana.Node{NodeID: 2}}},
	}
	diagnostics := Diagnostics(findings)
	if len(diagnostics) != 1 || len(diagnostics[0].Evidence) != 2 {
		t.Fatalf("unexpected diagnostics: %+v", diagnostics)
	}
}
