package mutation

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMutationManifest(t *testing.T) {
	path := filepath.Join(t.TempDir(), "architecture-mutation.json")
	content := `{
  "version": 1,
  "mutation_id": "mutation-1",
  "kind": "call-bypass",
  "language": "go",
  "summary": "bypass",
  "source_repository": "repository",
  "source_was_dirty": false,
  "base_commit": "abc123",
  "worktree": "worktree",
  "expected_added_relationships": [
    {"source": "api/api.go::CreateUser", "relation": "calls", "target": "storage/users.go::InsertUser"}
  ],
  "expected_removed_relationships": []
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	manifest, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if manifest.MutationID != "mutation-1" || len(manifest.ExpectedAddedRelationships) != 1 {
		t.Fatalf("unexpected manifest: %+v", manifest)
	}
}

func TestManifestRejectsUnqualifiedRelationship(t *testing.T) {
	manifest := Manifest{
		Version:          1,
		MutationID:       "mutation",
		Kind:             "call-bypass",
		Language:         "go",
		SourceRepository: "repository",
		BaseCommit:       "abc",
		Worktree:         "worktree",
		ExpectedAddedRelationships: []Relationship{
			{Source: "CreateUser", Relation: "calls", Target: "InsertUser"},
		},
	}
	if err := manifest.Validate(); err == nil {
		t.Fatal("expected unqualified relationship to fail")
	}
}
