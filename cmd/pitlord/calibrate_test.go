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
}

func (resolver fixedRevisionResolver) Resolve(string) (string, error) {
	return resolver.revision, nil
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
