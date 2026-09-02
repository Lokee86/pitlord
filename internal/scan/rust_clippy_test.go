package scan

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeCommandRunner struct {
	stdout []byte
	stderr []byte
	err    error
	cwd    string
	name   string
	args   []string
}

func (runner *fakeCommandRunner) Run(_ context.Context, cwd, name string, args ...string) ([]byte, []byte, error) {
	runner.cwd = cwd
	runner.name = name
	runner.args = append([]string(nil), args...)
	return runner.stdout, runner.stderr, runner.err
}

func TestClippyAnalyzerNormalizesAndDeduplicatesDiagnostics(t *testing.T) {
	repo := t.TempDir()
	if err := os.WriteFile(filepath.Join(repo, "Cargo.toml"), []byte("[package]\nname='fixture'\nversion='0.1.0'\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	message := `{"reason":"compiler-message","message":{"code":{"code":"clippy::clone_on_copy"},"level":"warning","message":"using clone on a Copy type","spans":[{"file_name":"src/main.rs","line_start":3,"line_end":3,"column_start":18,"column_end":31,"is_primary":true,"suggested_replacement":null,"suggestion_applicability":null}],"children":[{"code":null,"level":"help","message":"try removing the clone call","spans":[{"file_name":"src/main.rs","line_start":3,"line_end":3,"column_start":18,"column_end":31,"is_primary":true,"suggested_replacement":"value","suggestion_applicability":"MachineApplicable"}],"children":[]}]}}`
	runner := &fakeCommandRunner{stdout: []byte(message + "\n" + message + "\n" + `{"reason":"build-finished"}` + "\n")}
	analyzer := clippyAnalyzer{runner: runner}
	findings, err := analyzer.Analyze(context.Background(), AnalyzerContext{RepositoryRoot: repo, PathPrefix: "."})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected one deduplicated finding, got %+v", findings)
	}
	finding := findings[0]
	if finding.RuleID != "clippy::clone_on_copy" || finding.Scope.Path != "src/main.rs" || finding.Location == nil || finding.Location.StartLine != 3 {
		t.Fatalf("unexpected finding: %+v", finding)
	}
	if finding.SuggestedFix == nil || finding.SuggestedFix.Applicability != ApplicabilityMachine || len(finding.SuggestedFix.Edits) != 1 || finding.SuggestedFix.Edits[0].Replacement != "value" {
		t.Fatalf("unexpected suggestion: %+v", finding.SuggestedFix)
	}
	if runner.name != "cargo" || !strings.Contains(strings.Join(runner.args, " "), "clippy") || runner.cwd != repo {
		t.Fatalf("unexpected command: cwd=%q name=%q args=%v", runner.cwd, runner.name, runner.args)
	}
}

func TestClippyAnalyzerFiltersPathPrefix(t *testing.T) {
	repo := t.TempDir()
	if err := os.WriteFile(filepath.Join(repo, "Cargo.toml"), []byte("[workspace]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	output := []byte(`{"reason":"compiler-message","message":{"code":{"code":"clippy::needless_return"},"level":"warning","message":"unneeded return","spans":[{"file_name":"src/main.rs","line_start":1,"line_end":1,"column_start":1,"column_end":7,"is_primary":true}],"children":[]}}` + "\n")
	findings, err := (clippyAnalyzer{runner: &fakeCommandRunner{stdout: output}}).Analyze(context.Background(), AnalyzerContext{RepositoryRoot: repo, PathPrefix: "tests"})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("unexpected findings outside prefix: %+v", findings)
	}
}

func TestClippyAnalyzerPropagatesToolFailure(t *testing.T) {
	repo := t.TempDir()
	if err := os.WriteFile(filepath.Join(repo, "Cargo.toml"), []byte("[workspace]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := (clippyAnalyzer{runner: &fakeCommandRunner{stderr: []byte("component unavailable"), err: errors.New("exit 1")}}).Analyze(context.Background(), AnalyzerContext{RepositoryRoot: repo})
	if err == nil || !strings.Contains(err.Error(), "component unavailable") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuiltInClippyAnalyzerDoesNotRequireGraph(t *testing.T) {
	analyzers, err := BuiltInAnalyzers("clippy")
	if err != nil {
		t.Fatal(err)
	}
	if len(analyzers) != 1 || AnalyzersRequireGraph(analyzers) {
		t.Fatalf("unexpected built-in analyzers: %+v", analyzers)
	}
}
