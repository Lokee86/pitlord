package docguard

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func LoadChanges(ctx context.Context, repository, changedFrom string) ([]Change, error) {
	if strings.TrimSpace(changedFrom) == "" {
		return nil, fmt.Errorf("changed-from revision is required")
	}
	mergeBaseOutput, err := exec.CommandContext(ctx, "git", "-C", repository, "merge-base", changedFrom, "HEAD").Output()
	if err != nil {
		return nil, fmt.Errorf("resolve Git merge base from %q: %w", changedFrom, err)
	}
	mergeBase := strings.TrimSpace(string(mergeBaseOutput))
	if mergeBase == "" {
		return nil, fmt.Errorf("resolve Git merge base from %q: empty result", changedFrom)
	}
	args := []string{"-C", repository, "diff", "--name-status", "-z", "--find-renames", mergeBase, "--"}
	output, err := exec.CommandContext(ctx, "git", args...).Output()
	if err != nil {
		return nil, fmt.Errorf("read Git changes from %q: %w", changedFrom, err)
	}
	changes := parseChanges(output)
	untracked, err := exec.CommandContext(ctx, "git", "-C", repository, "ls-files", "--others", "--exclude-standard", "-z").Output()
	if err != nil {
		return nil, fmt.Errorf("read untracked Git paths: %w", err)
	}
	for _, path := range strings.Split(string(untracked), "\x00") {
		if path = strings.TrimSpace(path); path != "" {
			changes = append(changes, Change{Status: "A", Path: path})
		}
	}
	return changes, nil
}

func parseChanges(output []byte) []Change {
	parts := strings.Split(string(output), "\x00")
	changes := make([]Change, 0)
	for index := 0; index < len(parts); {
		status := strings.TrimSpace(parts[index])
		index++
		if status == "" || index >= len(parts) {
			continue
		}
		if strings.HasPrefix(status, "R") || strings.HasPrefix(status, "C") {
			if index+1 >= len(parts) {
				break
			}
			oldPath, newPath := parts[index], parts[index+1]
			index += 2
			changes = append(changes, Change{Status: status, Path: newPath, OldPath: oldPath})
			continue
		}
		path := parts[index]
		index++
		change := Change{Status: status, Path: path}
		if strings.HasPrefix(status, "D") {
			change.OldPath = path
		}
		changes = append(changes, change)
	}
	return changes
}
