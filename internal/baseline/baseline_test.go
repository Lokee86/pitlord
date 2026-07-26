package baseline

import (
	"path/filepath"
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
	"github.com/Lokee86/pitlord/internal/policy"
)

func TestBuildWriteLoadAndApply(t *testing.T) {
	target := arcana.Node{Identity: "target", Path: "storage/write.go"}
	known := policy.Evidence{
		Issue:    "forbidden_dependency",
		Source:   arcana.Node{Identity: "source", Path: "api/create.go"},
		Relation: "calls",
		Target:   &target,
	}
	newEvidence := policy.Evidence{
		Issue:  "unowned",
		Source: arcana.Node{Identity: "new", Path: "shared/new.go"},
	}
	diagnostics := []policy.Diagnostic{
		{RuleID: "dependency", Evidence: []policy.Evidence{known}},
		{RuleID: "ownership", Evidence: []policy.Evidence{newEvidence}},
	}

	document := Build(diagnostics[:1])
	if len(document.Entries) != 1 {
		t.Fatalf("expected one baseline entry, got %d", len(document.Entries))
	}
	path := filepath.Join(t.TempDir(), "pitlord.baseline.json")
	if err := Write(path, document); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	filtered, suppressed := Apply(loaded, diagnostics)
	if suppressed != 1 {
		t.Fatalf("expected one suppressed finding, got %d", suppressed)
	}
	if len(filtered) != 1 || filtered[0].RuleID != "ownership" {
		t.Fatalf("unexpected filtered diagnostics: %+v", filtered)
	}
}

func TestFingerprintPrefersStableNodeIdentity(t *testing.T) {
	first := policy.Evidence{Issue: "unowned", Source: arcana.Node{Identity: "stable", Path: "old/path.go"}}
	second := policy.Evidence{Issue: "unowned", Source: arcana.Node{Identity: "stable", Path: "new/path.go"}}
	if Fingerprint("ownership", first) != Fingerprint("ownership", second) {
		t.Fatal("renamed node with stable identity should retain its fingerprint")
	}
}
