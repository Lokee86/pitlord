package calibration

import (
	"crypto/sha256"
	"fmt"
	"os/exec"
	"strings"
)

type RevisionResolver interface {
	Resolve(repo string) (string, error)
}

type WorktreeDiffResolver interface {
	ResolveWorktreeDiffSHA256(repo string) (string, error)
}

type GitRevisionResolver struct {
	Command string
}

func (resolver GitRevisionResolver) Resolve(repo string) (string, error) {
	command := resolver.command()
	output, err := exec.Command(command, "-C", repo, "rev-parse", "HEAD").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("resolve corpus revision: %w: %s", err, strings.TrimSpace(string(output)))
	}
	revision := strings.TrimSpace(string(output))
	if revision == "" {
		return "", fmt.Errorf("resolve corpus revision: git returned an empty revision")
	}
	return revision, nil
}

func (resolver GitRevisionResolver) ResolveWorktreeDiffSHA256(repo string) (string, error) {
	output, err := exec.Command(
		resolver.command(), "-C", repo, "diff", "--binary", "--full-index", "--no-ext-diff", "HEAD", "--",
	).Output()
	if err != nil {
		return "", fmt.Errorf("resolve corpus worktree diff: %w", err)
	}
	digest := sha256.Sum256(output)
	return fmt.Sprintf("%x", digest[:]), nil
}

func (resolver GitRevisionResolver) command() string {
	command := strings.TrimSpace(resolver.Command)
	if command == "" {
		return "git"
	}
	return command
}

func VerifyRevision(resolver RevisionResolver, repo, expected string) error {
	actual, err := resolver.Resolve(repo)
	if err != nil {
		return err
	}
	if !strings.EqualFold(strings.TrimSpace(actual), strings.TrimSpace(expected)) {
		return fmt.Errorf("calibration corpus revision mismatch: expected %s, got %s", expected, actual)
	}
	return nil
}

func VerifyWorktreeDiff(resolver WorktreeDiffResolver, repo, expected string) error {
	actual, err := resolver.ResolveWorktreeDiffSHA256(repo)
	if err != nil {
		return err
	}
	if !strings.EqualFold(strings.TrimSpace(actual), strings.TrimSpace(expected)) {
		return fmt.Errorf("calibration corpus worktree diff mismatch: expected %s, got %s", expected, actual)
	}
	return nil
}
