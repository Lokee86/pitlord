package scan

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const AnalyzerClippy = "clippy"

type commandRunner interface {
	Run(context.Context, string, string, ...string) ([]byte, []byte, error)
}

type execCommandRunner struct{}

func (execCommandRunner) Run(ctx context.Context, cwd, name string, args ...string) ([]byte, []byte, error) {
	command := exec.CommandContext(ctx, name, args...)
	command.Dir = cwd
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}

type clippyAnalyzer struct {
	runner commandRunner
}

func NewClippyAnalyzer() Analyzer {
	return clippyAnalyzer{runner: execCommandRunner{}}
}

func (clippyAnalyzer) Metadata() AnalyzerMetadata {
	return AnalyzerMetadata{ID: AnalyzerClippy, Category: "lint", Language: "rust"}
}

func (analyzer clippyAnalyzer) Analyze(ctx context.Context, input AnalyzerContext) ([]Finding, error) {
	repositoryRoot := strings.TrimSpace(input.RepositoryRoot)
	if repositoryRoot == "" {
		repositoryRoot = "."
	}
	manifest, err := resolveCargoManifest(repositoryRoot, input.PathPrefix)
	if err != nil {
		return nil, err
	}
	runner := analyzer.runner
	if runner == nil {
		runner = execCommandRunner{}
	}
	stdout, stderr, err := runner.Run(ctx, filepath.Dir(manifest), "cargo",
		"clippy", "--manifest-path", manifest, "--message-format=json", "--all-targets", "--no-deps")
	if err != nil {
		message := strings.TrimSpace(string(stderr))
		if message == "" {
			message = err.Error()
		}
		return nil, fmt.Errorf("cargo clippy failed: %s", message)
	}
	return parseClippyOutput(stdout, repositoryRoot, filepath.Dir(manifest), input.PathPrefix)
}

func resolveCargoManifest(repositoryRoot, pathPrefix string) (string, error) {
	root, err := filepath.Abs(repositoryRoot)
	if err != nil {
		return "", fmt.Errorf("resolve repository root: %w", err)
	}
	start := root
	if prefix := normalizedUnitPath(pathPrefix); prefix != "" && prefix != "." {
		start = filepath.Join(root, filepath.FromSlash(prefix))
		if info, statErr := os.Stat(start); statErr == nil && !info.IsDir() {
			start = filepath.Dir(start)
		}
	}
	for directory := filepath.Clean(start); ; directory = filepath.Dir(directory) {
		manifest := filepath.Join(directory, "Cargo.toml")
		if info, statErr := os.Stat(manifest); statErr == nil && !info.IsDir() {
			return manifest, nil
		}
		if samePath(directory, root) {
			break
		}
		parent := filepath.Dir(directory)
		if parent == directory || !pathWithinRoot(root, parent) {
			break
		}
	}
	return "", fmt.Errorf("Cargo.toml was not found in %s", repositoryRoot)
}

func pathWithinRoot(root, candidate string) bool {
	relative, err := filepath.Rel(root, candidate)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

func samePath(left, right string) bool {
	return strings.EqualFold(filepath.Clean(left), filepath.Clean(right))
}
