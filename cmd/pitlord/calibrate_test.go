package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

type fixedRevisionResolver struct {
	revision string
	diffHash string
}

func (resolver fixedRevisionResolver) Resolve(string) (string, error) {
	return resolver.revision, nil
}

func (resolver fixedRevisionResolver) ResolveWorktreeDiffSHA256(string) (string, error) {
	return resolver.diffHash, nil
}

func TestCalibrateFailsOnLabelledMismatchWhenRequested(t *testing.T) {
	dir := t.TempDir()
	reference := filepath.Join(dir, "reference.json")
	if err := os.WriteFile(reference, []byte(`{
		"schema":"pitlord.calibration.v1",
		"corpus":"fixture",
		"source_revision":"abc123",
		"detector":"dependency-pressure",
		"expectations":[{
			"id":"pressure",
			"path":"src/pressure.go",
			"class":"pressure",
			"finding":"required"
		}]
	}`), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := runCalibrateWithDependencies(
		[]string{
			"--repo", dir,
			"--reference", reference,
			"--snapshot", dir,
			"--format", "json",
			"--fail-on-mismatch",
		},
		&stdout,
		&stderr,
		commandScanLoader{graph: arcana.Graph{Outgoing: map[uint32][]arcana.Relationship{}}},
		fixedRevisionResolver{revision: "abc123"},
	)
	if code != 1 {
		t.Fatalf("expected mismatch exit 1, got %d: %s", code, stderr.String())
	}
	if !bytes.Contains(stdout.Bytes(), []byte(`"false_negative": 1`)) {
		t.Fatalf("unexpected output: %s", stdout.String())
	}
}

func TestCalibrateAcceptsPinnedWorktreeDiff(t *testing.T) {
	dir := t.TempDir()
	reference := filepath.Join(dir, "reference.json")
	diffHash := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	if err := os.WriteFile(reference, []byte(`{
		"schema":"pitlord.calibration.v1",
		"corpus":"fixture",
		"source_revision":"abc123",
		"worktree_diff_sha256":"`+diffHash+`",
		"detector":"dependency-pressure",
		"expectations":[{"id":"clean","path_prefix":".","class":"clean","finding":"absent"}]
	}`), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := runCalibrateWithDependencies(
		[]string{"--repo", dir, "--reference", reference, "--snapshot", dir},
		&stdout,
		&stderr,
		commandScanLoader{graph: arcana.Graph{Outgoing: map[uint32][]arcana.Relationship{}}},
		fixedRevisionResolver{revision: "abc123", diffHash: diffHash},
	)
	if code != 0 {
		t.Fatalf("expected pinned worktree calibration to pass, got code=%d stderr=%s", code, stderr.String())
	}
}

func TestCalibrateRejectsWorktreeDiffMismatch(t *testing.T) {
	dir := t.TempDir()
	reference := filepath.Join(dir, "reference.json")
	if err := os.WriteFile(reference, []byte(`{
		"schema":"pitlord.calibration.v1",
		"corpus":"fixture",
		"source_revision":"abc123",
		"worktree_diff_sha256":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"detector":"dependency-pressure",
		"expectations":[{"id":"clean","path_prefix":".","class":"clean","finding":"absent"}]
	}`), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := runCalibrateWithDependencies(
		[]string{"--repo", dir, "--reference", reference, "--snapshot", dir},
		&stdout,
		&stderr,
		commandScanLoader{graph: arcana.Graph{Outgoing: map[uint32][]arcana.Relationship{}}},
		fixedRevisionResolver{revision: "abc123", diffHash: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},
	)
	if code != 2 || !bytes.Contains(stderr.Bytes(), []byte("worktree diff mismatch")) {
		t.Fatalf("expected worktree diff mismatch, got code=%d stderr=%s", code, stderr.String())
	}
}

func TestCalibrateRunsTargetedSemanticAnalyzer(t *testing.T) {
	dir := t.TempDir()
	reference := filepath.Join(dir, "reference.json")
	if err := os.WriteFile(reference, []byte(`{
		"schema":"pitlord.calibration.v1",
		"corpus":"fixture",
		"source_revision":"abc123",
		"analyzer":"swallowed-error",
		"rule_id":"swallowed-error",
		"language":"python",
		"require_fully_labelled":true,
		"expectations":[{
			"id":"python-swallowed",
			"path":"src/app.py",
			"location":{"start_line":12,"start_column":5,"end_line":13,"end_column":9},
			"class":"swallowed-error",
			"finding":"required"
		}]
	}`), 0o644); err != nil {
		t.Fatal(err)
	}
	graph := arcana.Graph{
		Sources: []arcana.Node{
			{NodeID: 1, Kind: "protocol", Path: "src/app.py", Name: "semantic-capabilities:python:calls,control-flow,error-handling,source-spans"},
			{NodeID: 2, Kind: "protocol", Path: "src/app.py", Name: "error-handler:python", Identity: "handler", Span: &arcana.Span{Path: "src/app.py", StartLine: 12, StartColumn: 5, EndLine: 13, EndColumn: 9}},
		},
		Outgoing: map[uint32][]arcana.Relationship{},
	}
	var stdout, stderr bytes.Buffer
	code := runCalibrateWithDependencies(
		[]string{"--repo", dir, "--reference", reference, "--snapshot", dir, "--format", "json", "--fail-on-mismatch"},
		&stdout,
		&stderr,
		commandScanLoader{graph: graph},
		fixedRevisionResolver{revision: "abc123"},
	)
	if code != 0 {
		t.Fatalf("expected semantic calibration to pass, got code=%d stderr=%s", code, stderr.String())
	}
	if !bytes.Contains(stdout.Bytes(), []byte(`"true_positive": 1`)) || !bytes.Contains(stdout.Bytes(), []byte(`"language": "python"`)) {
		t.Fatalf("unexpected semantic calibration output: %s", stdout.String())
	}
}

func TestCalibrateRejectsRevisionMismatch(t *testing.T) {
	dir := t.TempDir()
	reference := filepath.Join(dir, "reference.json")
	if err := os.WriteFile(reference, []byte(`{
		"schema":"pitlord.calibration.v1",
		"corpus":"fixture",
		"source_revision":"expected",
		"detector":"dependency-pressure",
		"expectations":[{"id":"clean","path":"src/a.go","class":"clean","finding":"absent"}]
	}`), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := runCalibrateWithDependencies(
		[]string{"--repo", dir, "--reference", reference, "--snapshot", dir},
		&stdout,
		&stderr,
		commandScanLoader{graph: arcana.Graph{Outgoing: map[uint32][]arcana.Relationship{}}},
		fixedRevisionResolver{revision: "actual"},
	)
	if code != 2 || !bytes.Contains(stderr.Bytes(), []byte("revision mismatch")) {
		t.Fatalf("expected revision mismatch, got code=%d stderr=%s", code, stderr.String())
	}
}
