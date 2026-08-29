package docguard

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestLoadChangesIncludesWorkingTreeChanges(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not available")
	}
	root := t.TempDir()
	runGit(t, root, "init", "-q")
	runGit(t, root, "config", "user.email", "pitlord@example.invalid")
	runGit(t, root, "config", "user.name", "Pitlord Test")
	path := filepath.Join(root, "source.go")
	if err := os.WriteFile(path, []byte("package sample\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", "source.go")
	runGit(t, root, "commit", "-q", "-m", "base")
	if err := os.WriteFile(path, []byte("package sample\n\nvar changed = true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "new.go"), []byte("package sample\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	changes, err := LoadChanges(context.Background(), root, "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 2 {
		t.Fatalf("unexpected working-tree changes: %#v", changes)
	}
	if changes[0].Status != "M" || changes[0].Path != "source.go" {
		t.Fatalf("unexpected tracked change: %#v", changes[0])
	}
	if changes[1].Status != "A" || changes[1].Path != "new.go" {
		t.Fatalf("unexpected untracked change: %#v", changes[1])
	}
}

func runGit(t *testing.T, root string, args ...string) {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", root}, args...)...)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v: %s", args, err, output)
	}
}
