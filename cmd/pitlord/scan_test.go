package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
	"github.com/Lokee86/pitlord/internal/scan"
)

type commandScanLoader struct {
	graph arcana.Graph
}

func (loader commandScanLoader) LoadGraphWithOptions(context.Context, string, arcana.LoadOptions) (arcana.Graph, error) {
	return loader.graph, nil
}

func TestScanRunsWithoutPolicy(t *testing.T) {
	snapshotPath := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := runScanWithLoader(
		[]string{"--snapshot", snapshotPath, "--format", "json"},
		&stdout,
		&stderr,
		commandScanLoader{graph: arcana.Graph{Outgoing: map[uint32][]arcana.Relationship{}}},
	)
	if code != 0 {
		t.Fatalf("expected success, got %d: %s", code, stderr.String())
	}
	var result scan.Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Schema != scan.Schema || result.Findings == nil {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestScanResolvesRepositorySnapshot(t *testing.T) {
	repo := t.TempDir()
	digest := "abc123"
	snapshotPath := filepath.Join(repo, ".arcana", "snapshots", digest)
	if err := os.MkdirAll(snapshotPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, ".arcana", "CURRENT"), []byte("sha256:"+digest+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if code := runScanWithLoader(
		[]string{"--repo", repo},
		&stdout,
		&stderr,
		commandScanLoader{graph: arcana.Graph{Outgoing: map[uint32][]arcana.Relationship{}}},
	); code != 0 {
		t.Fatalf("expected success, got %d: %s", code, stderr.String())
	}
}
