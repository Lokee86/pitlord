package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveCurrentSnapshot(t *testing.T) {
	repo := t.TempDir()
	digest := "0123456789abcdef"
	snapshotPath := filepath.Join(repo, ".arcana", "snapshots", digest)
	if err := os.MkdirAll(snapshotPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, ".arcana", "CURRENT"), []byte("sha256:"+digest+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	resolved, err := Resolve(repo, "")
	if err != nil {
		t.Fatal(err)
	}
	if resolved != snapshotPath {
		t.Fatalf("expected %q, got %q", snapshotPath, resolved)
	}
}
