package docguard

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func LoadDataset(ctx context.Context, repository, docsRoot, command, arcanaCommand string) (Dataset, error) {
	resolved, err := resolveDDocsCommand(command)
	if err != nil {
		return Dataset{}, err
	}
	args := []string{"codemaps", "export"}
	if strings.TrimSpace(docsRoot) != "" {
		args = append(args, "--root", docsRoot)
	}
	process := exec.CommandContext(ctx, resolved, args...)
	process.Dir = repository
	process.Env = os.Environ()
	if strings.TrimSpace(arcanaCommand) != "" {
		process.Env = append(process.Env, "DDOCS_ARCANA_COMMAND="+arcanaCommand)
	}
	output, err := process.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return Dataset{}, fmt.Errorf("ddocs codemap export failed: %s", strings.TrimSpace(string(exitErr.Stderr)))
		}
		return Dataset{}, fmt.Errorf("run ddocs codemap export: %w", err)
	}
	var dataset Dataset
	if err := json.Unmarshal(output, &dataset); err != nil {
		return Dataset{}, fmt.Errorf("decode ddocs codemap dataset: %w", err)
	}
	if dataset.SchemaVersion != DatasetSchemaVersion {
		return Dataset{}, fmt.Errorf("unsupported ddocs codemap dataset schema %d", dataset.SchemaVersion)
	}
	return dataset, nil
}

func resolveDDocsCommand(requested string) (string, error) {
	if requested = strings.TrimSpace(requested); requested != "" {
		if path, err := exec.LookPath(requested); err == nil {
			return path, nil
		}
		if info, err := os.Stat(requested); err == nil && !info.IsDir() {
			return filepath.Abs(requested)
		}
		return "", fmt.Errorf("Demon Docs executable %q was not found", requested)
	}
	if executable, err := os.Executable(); err == nil {
		name := "ddocs"
		if runtime.GOOS == "windows" {
			name += ".exe"
		}
		candidate := filepath.Join(filepath.Dir(executable), name)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, nil
		}
	}
	if path, err := exec.LookPath("ddocs"); err == nil {
		return path, nil
	}
	return "", fmt.Errorf("Demon Docs executable was not found; install ddocs or use --ddocs")
}
