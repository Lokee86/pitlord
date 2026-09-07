package arcana

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type lexiconConfiguration struct {
	AdapterRoot string `json:"adapter_root"`
}

// ResolveCommand finds the Arcana executable associated with the current
// Lexicon + Arcana installation before falling back to PATH.
func ResolveCommand(repositoryRoot, requested string) (string, error) {
	executable, _ := os.Executable()
	return resolveCommand(repositoryRoot, requested, executable)
}

func resolveCommand(repositoryRoot, requested, executable string) (string, error) {
	if repositoryRoot = strings.TrimSpace(repositoryRoot); repositoryRoot == "" {
		repositoryRoot = "."
	}
	if requested = strings.TrimSpace(requested); requested != "" {
		return requested, nil
	}
	if configured := strings.TrimSpace(os.Getenv("PITLORD_ARCANA_COMMAND")); configured != "" {
		if command := resolveConfiguredCommand(repositoryRoot, configured); command != "" {
			return command, nil
		}
		return "", fmt.Errorf("PITLORD_ARCANA_COMMAND %q does not resolve to an executable", configured)
	}
	if command := commandFromLexiconConfiguration(repositoryRoot); command != "" {
		return command, nil
	}
	if executable != "" {
		if command := commandInDirectory(filepath.Dir(executable)); command != "" {
			return command, nil
		}
	}
	if lexicon, err := exec.LookPath("lexicon"); err == nil {
		if command := commandInDirectory(filepath.Dir(lexicon)); command != "" {
			return command, nil
		}
	}
	if command := commandFromDevelopmentCheckout(repositoryRoot); command != "" {
		return command, nil
	}
	if command, err := exec.LookPath("arcana"); err == nil {
		return command, nil
	}
	return "", fmt.Errorf("Arcana executable was not found; install the Lexicon + Arcana bundle, put arcana on PATH, or use --arcana")
}

func commandFromLexiconConfiguration(repositoryRoot string) string {
	data, err := os.ReadFile(filepath.Join(repositoryRoot, ".lexicon", "config.json"))
	if err != nil {
		return ""
	}
	var configuration lexiconConfiguration
	if json.Unmarshal(data, &configuration) != nil {
		return ""
	}
	adapterRoot := strings.TrimSpace(configuration.AdapterRoot)
	if adapterRoot == "" {
		return ""
	}
	if !filepath.IsAbs(adapterRoot) {
		adapterRoot = filepath.Join(repositoryRoot, adapterRoot)
	}
	root := filepath.Clean(filepath.Dir(adapterRoot))
	for depth := 0; depth < 3; depth++ {
		if command := commandUnderRoot(root); command != "" {
			return command
		}
		parent := filepath.Dir(root)
		if parent == root {
			break
		}
		root = parent
	}
	return ""
}

func commandFromDevelopmentCheckout(repositoryRoot string) string {
	root, err := filepath.Abs(repositoryRoot)
	if err != nil {
		root = filepath.Clean(repositoryRoot)
	}
	for depth := 0; depth < 4; depth++ {
		if command := commandUnderRoot(filepath.Join(root, "lexicon-arcana")); command != "" {
			return command
		}
		parent := filepath.Dir(root)
		if parent == root {
			break
		}
		root = parent
	}
	return ""
}

func commandUnderRoot(root string) string {
	for _, directory := range []string{
		root,
		filepath.Join(root, "bin"),
		filepath.Join(root, "arcana", "target", "release"),
		filepath.Join(root, "arcana", "target", "debug"),
	} {
		if command := commandInDirectory(directory); command != "" {
			return command
		}
	}
	return ""
}

func commandInDirectory(directory string) string {
	for _, name := range commandNames("arcana") {
		candidate := filepath.Join(directory, name)
		if regularFile(candidate) {
			absolute, err := filepath.Abs(candidate)
			if err == nil {
				return absolute
			}
			return filepath.Clean(candidate)
		}
	}
	return ""
}

func resolveConfiguredCommand(repositoryRoot, command string) string {
	command = strings.TrimSpace(command)
	if command == "" {
		return ""
	}
	if filepath.IsAbs(command) {
		if regularFile(command) {
			return filepath.Clean(command)
		}
		return ""
	}
	if candidate := filepath.Join(repositoryRoot, command); regularFile(candidate) {
		absolute, err := filepath.Abs(candidate)
		if err == nil {
			return absolute
		}
		return filepath.Clean(candidate)
	}
	if path, err := exec.LookPath(command); err == nil {
		return path
	}
	return ""
}

func commandNames(name string) []string {
	if runtime.GOOS == "windows" {
		return []string{name + ".exe", name}
	}
	return []string{name}
}

func regularFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
