package arcana

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestResolveCommandPreservesExplicitOverride(t *testing.T) {
	clearArcanaEnvironment(t)
	got, err := resolveCommand(t.TempDir(), "custom-arcana", "")
	if err != nil {
		t.Fatal(err)
	}
	if got != "custom-arcana" {
		t.Fatalf("explicit command = %q, want custom-arcana", got)
	}
}

func TestResolveCommandUsesPitlordEnvironmentOverride(t *testing.T) {
	provider := writeExecutable(t, t.TempDir(), "arcana")
	clearArcanaEnvironment(t)
	t.Setenv("PITLORD_ARCANA_COMMAND", provider)

	got, err := resolveCommand(t.TempDir(), "", "")
	if err != nil {
		t.Fatal(err)
	}
	if got != provider {
		t.Fatalf("environment command = %q, want %q", got, provider)
	}
}

func TestResolveCommandRejectsInvalidPitlordEnvironmentOverride(t *testing.T) {
	clearArcanaEnvironment(t)
	t.Setenv("PITLORD_ARCANA_COMMAND", filepath.Join(t.TempDir(), "missing-arcana"))
	_, err := resolveCommand(t.TempDir(), "", "")
	if err == nil || !strings.Contains(err.Error(), "PITLORD_ARCANA_COMMAND") {
		t.Fatalf("expected explicit environment failure, got %v", err)
	}
}

func TestResolveCommandFindsInstalledArcanaFromLexiconAdapterRoot(t *testing.T) {
	installRoot := t.TempDir()
	provider := writeExecutable(t, installRoot, "arcana")
	adapterRoot := filepath.Join(installRoot, "adapters")
	if err := os.MkdirAll(adapterRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	repository := t.TempDir()
	writeLexiconConfig(t, repository, adapterRoot)
	clearArcanaEnvironment(t)

	got, err := resolveCommand(repository, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if got != provider {
		t.Fatalf("installed command = %q, want %q", got, provider)
	}
}

func TestResolveCommandFindsSourceArcanaFromLexiconAdapterRoot(t *testing.T) {
	checkout := t.TempDir()
	provider := writeExecutable(t, filepath.Join(checkout, "arcana", "target", "release"), "arcana")
	adapterRoot := filepath.Join(checkout, "lexicon", "adapters")
	if err := os.MkdirAll(adapterRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	repository := t.TempDir()
	writeLexiconConfig(t, repository, adapterRoot)
	clearArcanaEnvironment(t)

	got, err := resolveCommand(repository, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if got != provider {
		t.Fatalf("source command = %q, want %q", got, provider)
	}
}

func TestResolveCommandFindsArcanaBesideLexiconOnPath(t *testing.T) {
	bin := t.TempDir()
	_ = writeExecutable(t, bin, "lexicon")
	provider := writeExecutable(t, bin, "arcana")
	clearArcanaEnvironment(t)
	t.Setenv("PATH", bin)

	got, err := resolveCommand(t.TempDir(), "", "")
	if err != nil {
		t.Fatal(err)
	}
	if got != provider {
		t.Fatalf("PATH-adjacent command = %q, want %q", got, provider)
	}
}

func TestResolveCommandFindsSiblingDevelopmentCheckout(t *testing.T) {
	workspace := t.TempDir()
	repository := filepath.Join(workspace, "project")
	if err := os.MkdirAll(repository, 0o755); err != nil {
		t.Fatal(err)
	}
	provider := writeExecutable(t, filepath.Join(workspace, "lexicon-arcana", "arcana", "target", "release"), "arcana")
	clearArcanaEnvironment(t)

	got, err := resolveCommand(repository, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if got != provider {
		t.Fatalf("development command = %q, want %q", got, provider)
	}
}

func TestResolveCommandDoesNotUseRetiredGrimoireEnvironment(t *testing.T) {
	legacy := writeExecutable(t, t.TempDir(), "arcana")
	clearArcanaEnvironment(t)
	t.Setenv("GRIMOIRE_ARCANA_COMMAND", legacy)
	_, err := resolveCommand(t.TempDir(), "", "")
	if err == nil {
		t.Fatal("retired Grimoire environment unexpectedly resolved Arcana")
	}
}

func TestResolveCommandReportsCurrentInstallationGuidance(t *testing.T) {
	clearArcanaEnvironment(t)
	_, err := resolveCommand(t.TempDir(), "", "")
	if err == nil || !strings.Contains(err.Error(), "Lexicon + Arcana bundle") || strings.Contains(err.Error(), "Grimoire") {
		t.Fatalf("unexpected missing Arcana error: %v", err)
	}
}

func clearArcanaEnvironment(t *testing.T) {
	t.Helper()
	t.Setenv("PITLORD_ARCANA_COMMAND", "")
	t.Setenv("GRIMOIRE_ARCANA_COMMAND", "")
	t.Setenv("GRIMOIRE_HOME", "")
	t.Setenv("PATH", t.TempDir())
}

func writeLexiconConfig(t *testing.T, repository, adapterRoot string) {
	t.Helper()
	directory := filepath.Join(repository, ".lexicon")
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(lexiconConfiguration{AdapterRoot: adapterRoot})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "config.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeExecutable(t *testing.T, directory, name string) string {
	t.Helper()
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	path := filepath.Join(directory, name)
	if err := os.WriteFile(path, []byte(name), 0o755); err != nil {
		t.Fatal(err)
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	return absolute
}
