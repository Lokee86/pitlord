package arcana

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestResolveCommandPreservesExplicitOverride(t *testing.T) {
	t.Setenv("GRIMOIRE_ARCANA_COMMAND", "")
	got, err := resolveCommand(t.TempDir(), "custom-arcana", "")
	if err != nil {
		t.Fatal(err)
	}
	if got != "custom-arcana" {
		t.Fatalf("explicit command = %q, want custom-arcana", got)
	}
}

func TestResolveCommandUsesRepositoryProviderConfiguration(t *testing.T) {
	repository := t.TempDir()
	provider := writeArcanaExecutable(t, t.TempDir())
	writeProviderConfig(t, repository, provider)
	t.Setenv("GRIMOIRE_ARCANA_COMMAND", "")
	t.Setenv("GRIMOIRE_HOME", "")
	t.Setenv("PATH", t.TempDir())

	got, err := resolveCommand(repository, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if got != provider {
		t.Fatalf("configured command = %q, want %q", got, provider)
	}
}

func TestResolveCommandFindsInstalledArcanaFromLexiconAdapterRoot(t *testing.T) {
	installRoot := t.TempDir()
	provider := writeArcanaExecutable(t, installRoot)
	adapterRoot := filepath.Join(installRoot, "adapters")
	if err := os.MkdirAll(adapterRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	repository := t.TempDir()
	writeLexiconConfig(t, repository, adapterRoot)
	t.Setenv("GRIMOIRE_ARCANA_COMMAND", "")
	t.Setenv("GRIMOIRE_HOME", "")
	t.Setenv("PATH", t.TempDir())

	got, err := resolveCommand(repository, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if got != provider {
		t.Fatalf("installed command = %q, want %q", got, provider)
	}
}

func TestResolveCommandFindsBuildArcanaFromLexiconAdapterRoot(t *testing.T) {
	grimoireRoot := t.TempDir()
	provider := writeArcanaExecutable(t, filepath.Join(grimoireRoot, "build", "bin"))
	adapterRoot := filepath.Join(grimoireRoot, "build", "adapters")
	if err := os.MkdirAll(adapterRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	repository := t.TempDir()
	writeLexiconConfig(t, repository, adapterRoot)
	t.Setenv("GRIMOIRE_ARCANA_COMMAND", "")
	t.Setenv("GRIMOIRE_HOME", "")
	t.Setenv("PATH", t.TempDir())

	got, err := resolveCommand(repository, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if got != provider {
		t.Fatalf("build command = %q, want %q", got, provider)
	}
}

func TestResolveCommandUsesEnvironmentOverride(t *testing.T) {
	provider := writeArcanaExecutable(t, t.TempDir())
	t.Setenv("GRIMOIRE_ARCANA_COMMAND", provider)

	got, err := resolveCommand(t.TempDir(), "", "")
	if err != nil {
		t.Fatal(err)
	}
	if got != provider {
		t.Fatalf("environment command = %q, want %q", got, provider)
	}
}

func TestResolveCommandReportsMissingProvider(t *testing.T) {
	t.Setenv("GRIMOIRE_ARCANA_COMMAND", "")
	t.Setenv("GRIMOIRE_HOME", "")
	t.Setenv("PATH", t.TempDir())
	if _, err := resolveCommand(t.TempDir(), "", ""); err == nil {
		t.Fatal("expected missing Arcana error")
	}
}

func writeProviderConfig(t *testing.T, repository, command string) {
	t.Helper()
	directory := filepath.Join(repository, ".grimoire")
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(providerConfiguration{Version: providerConfigVersion, ArcanaCommand: command})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "providers.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
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

func writeArcanaExecutable(t *testing.T, directory string) string {
	t.Helper()
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	name := "arcana"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	path := filepath.Join(directory, name)
	if err := os.WriteFile(path, []byte("arcana"), 0o755); err != nil {
		t.Fatal(err)
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	return absolute
}
