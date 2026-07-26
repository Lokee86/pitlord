package snapshot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Resolve(repo, explicit string) (string, error) {
	if explicit != "" {
		path, err := filepath.Abs(explicit)
		if err != nil {
			return "", fmt.Errorf("resolve snapshot path: %w", err)
		}
		if err := requireDirectory(path); err != nil {
			return "", err
		}
		return path, nil
	}
	if repo == "" {
		repo = "."
	}
	repoPath, err := filepath.Abs(repo)
	if err != nil {
		return "", fmt.Errorf("resolve repository path: %w", err)
	}
	currentPath := filepath.Join(repoPath, ".arcana", "CURRENT")
	data, err := os.ReadFile(currentPath)
	if err != nil {
		return "", fmt.Errorf("read Arcana CURRENT: %w", err)
	}
	identity := strings.TrimSpace(string(data))
	const prefix = "sha256:"
	if !strings.HasPrefix(identity, prefix) || len(identity) == len(prefix) {
		return "", fmt.Errorf("invalid Arcana CURRENT identity %q", identity)
	}
	digest := strings.TrimPrefix(identity, prefix)
	path := filepath.Join(repoPath, ".arcana", "snapshots", digest)
	if err := requireDirectory(path); err != nil {
		return "", err
	}
	return path, nil
}

func requireDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("open Arcana snapshot %q: %w", path, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("Arcana snapshot %q is not a directory", path)
	}
	return nil
}
