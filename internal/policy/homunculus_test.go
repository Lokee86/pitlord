package policy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFromHomunculusBuildsPolicyAndExpectations(t *testing.T) {
	path := filepath.Join(t.TempDir(), "homunculus.manifest.json")
	manifest := `{
  "version": 1,
  "architecture": {
    "areas": ["api", "service", "storage"],
    "forbidden": ["api -> storage", "storage -> service"]
  },
  "mutations": [
    {"expected_diagnostics": ["api-must-not-access-storage"]}
  ]
}`
	if err := os.WriteFile(path, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	document, expected, err := FromHomunculus(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(document.Areas) != 3 {
		t.Fatalf("expected three areas, got %d", len(document.Areas))
	}
	if len(document.Rules) != 2 {
		t.Fatalf("expected two rules, got %d", len(document.Rules))
	}
	if document.Rules[0].ID != "api-must-not-access-storage" {
		t.Fatalf("unexpected rule id %q", document.Rules[0].ID)
	}
	if len(document.Rules[0].FromAreas) != 1 || document.Rules[0].FromAreas[0] != "api" {
		t.Fatalf("unexpected source areas: %v", document.Rules[0].FromAreas)
	}
	if len(expected) != 1 || expected[0] != "api-must-not-access-storage" {
		t.Fatalf("unexpected expected diagnostics: %v", expected)
	}
}

func TestFromHomunculusRejectsUnknownAreas(t *testing.T) {
	path := filepath.Join(t.TempDir(), "homunculus.manifest.json")
	manifest := `{
  "version": 1,
  "architecture": {
    "areas": ["api", "service"],
    "forbidden": ["api -> storage"]
  }
}`
	if err := os.WriteFile(path, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := FromHomunculus(path); err == nil {
		t.Fatal("expected unknown area validation to fail")
	}
}
