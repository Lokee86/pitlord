package calibration

import (
	"fmt"
	"os/exec"
	"strings"
)

type RevisionResolver interface {
	Resolve(repo string) (string, error)
}

type GitRevisionResolver struct {
	Command string
}

func (resolver GitRevisionResolver) Resolve(repo string) (string, error) {
	command := strings.TrimSpace(resolver.Command)
	if command == "" {
		command = "git"
	}
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
