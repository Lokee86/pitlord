package scan

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

type Input struct {
	SnapshotPath string
	PathPrefix   string
}

func Run(input Input) (Result, error) {
	if strings.TrimSpace(input.SnapshotPath) == "" {
		return Result{}, fmt.Errorf("Arcana snapshot is required")
	}
	pathPrefix, err := normalizePathPrefix(input.PathPrefix)
	if err != nil {
		return Result{}, err
	}
	return finalize(Result{
		Schema: Schema,
		Scope: Scope{
			Kind: "repository",
			Path: pathPrefix,
		},
		Findings: []Finding{},
	}), nil
}

func normalizePathPrefix(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" || value == "." {
		return ".", nil
	}
	if filepath.IsAbs(value) {
		return "", fmt.Errorf("scan path prefix must be repository-relative")
	}
	value = filepath.ToSlash(value)
	cleaned := path.Clean(value)
	if cleaned == "." {
		return ".", nil
	}
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return "", fmt.Errorf("scan path prefix must not escape the repository")
	}
	return strings.TrimPrefix(cleaned, "./"), nil
}
