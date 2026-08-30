package calibration

import (
	"strings"
	"testing"
)

type fixedWorktreeDiffResolver struct {
	hash string
}

func (resolver fixedWorktreeDiffResolver) ResolveWorktreeDiffSHA256(string) (string, error) {
	return resolver.hash, nil
}

func TestVerifyWorktreeDiffAcceptsMatchingHash(t *testing.T) {
	hash := strings.Repeat("a", 64)
	if err := VerifyWorktreeDiff(fixedWorktreeDiffResolver{hash: hash}, ".", hash); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyWorktreeDiffRejectsMismatch(t *testing.T) {
	err := VerifyWorktreeDiff(fixedWorktreeDiffResolver{hash: strings.Repeat("b", 64)}, ".", strings.Repeat("a", 64))
	if err == nil || !strings.Contains(err.Error(), "worktree diff mismatch") {
		t.Fatalf("expected worktree diff mismatch, got %v", err)
	}
}
