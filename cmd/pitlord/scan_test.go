package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Lokee86/pitlord/internal/scan"
)

func TestScanRunsWithoutPolicy(t *testing.T) {
	snapshotPath := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := run([]string{"scan", "--snapshot", snapshotPath, "--format", "json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected success, got %d: %s", code, stderr.String())
	}
	var result scan.Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Schema != scan.Schema || result.Findings == nil || len(result.Findings) != 0 {
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
	if code := run([]string{"scan", "--repo", repo}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected success, got %d: %s", code, stderr.String())
	}
}
