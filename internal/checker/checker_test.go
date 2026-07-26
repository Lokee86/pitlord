package checker

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestRunRepositoryOnlyPolicyDoesNotRequireArcanaSnapshot(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "bad.txt"), []byte("forbidden\n"), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}
	policyPath := filepath.Join(root, "pitlord.json")
	policy := `{
  "version": 1,
  "rules": [
    {
      "id": "content",
      "type": "forbid_content",
      "literal": "forbidden",
      "include_paths": ["*.txt"]
    }
  ]
}
`
	if err := os.WriteFile(policyPath, []byte(policy), 0o644); err != nil {
		t.Fatalf("write policy: %v", err)
	}

	result, err := Run(context.Background(), Request{
		Repository:    root,
		PolicyPath:    policyPath,
		ArcanaCommand: "arcana-command-that-must-not-run",
	})
	if err != nil {
		t.Fatalf("run repository-only policy: %v", err)
	}
	if result.Report.Snapshot != "" {
		t.Fatalf("repository-only report unexpectedly resolved a snapshot: %q", result.Report.Snapshot)
	}
	if len(result.Report.Diagnostics) != 1 || result.Report.Diagnostics[0].RuleID != "content" {
		t.Fatalf("unexpected diagnostics: %#v", result.Report.Diagnostics)
	}
	if result.Report.Summary.SourceNodesScanned != 0 || result.Report.Summary.RelationshipsScanned != 0 {
		t.Fatalf("repository-only policy unexpectedly loaded graph data: %#v", result.Report.Summary)
	}
}
