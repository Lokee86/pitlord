package scan

import (
	"context"
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

type fakeGraphLoader struct {
	graph   arcana.Graph
	options arcana.LoadOptions
}

func (loader *fakeGraphLoader) LoadGraphWithOptions(_ context.Context, _ string, options arcana.LoadOptions) (arcana.Graph, error) {
	loader.options = options
	return loader.graph, nil
}

func TestEngineRunsWithoutPolicyAndScopesGraphLoad(t *testing.T) {
	loader := &fakeGraphLoader{graph: arcana.Graph{Outgoing: map[uint32][]arcana.Relationship{}}}
	result, err := (Engine{Loader: loader}).Run(context.Background(), Input{SnapshotPath: "snapshot", PathPrefix: "./src/service/"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Schema != Schema || result.Scope.Path != "src/service" {
		t.Fatalf("unexpected result: %+v", result)
	}
	if len(loader.options.SourcePrefixes) != 1 || loader.options.SourcePrefixes[0] != "src/service" {
		t.Fatalf("unexpected graph load scope: %+v", loader.options)
	}
	if result.Findings == nil {
		t.Fatal("findings must be non-nil")
	}
}

func TestEngineRejectsMissingSnapshotEscapingScopeAndLoader(t *testing.T) {
	loader := &fakeGraphLoader{}
	if _, err := (Engine{Loader: loader}).Run(context.Background(), Input{}); err == nil {
		t.Fatal("expected missing snapshot failure")
	}
	if _, err := (Engine{Loader: loader}).Run(context.Background(), Input{SnapshotPath: "snapshot", PathPrefix: "../outside"}); err == nil {
		t.Fatal("expected escaping scope failure")
	}
	if _, err := (Engine{}).Run(context.Background(), Input{SnapshotPath: "snapshot"}); err == nil {
		t.Fatal("expected missing loader failure")
	}
}

func TestFinalizeOrdersFindingsAndRecalculatesSummary(t *testing.T) {
	result := finalize(Result{Findings: []Finding{
		{ID: "b", Detector: "hub", Disposition: DispositionAdvisory, Severity: SeverityWarning, Scope: Scope{Path: "z"}},
		{ID: "a", Detector: "cycle", Disposition: DispositionGuard, Severity: SeverityWarning, Scope: Scope{Path: "a"}},
		{ID: "c", Detector: "hub", Disposition: DispositionAdvisory, Severity: SeverityCritical, Scope: Scope{Path: "a"}},
	}})
	if result.Findings[0].ID != "a" || result.Findings[1].ID != "c" || result.Findings[2].ID != "b" {
		t.Fatalf("unexpected finding order: %#v", result.Findings)
	}
	if result.Summary.FindingCount != 3 || result.Summary.Guard != 1 || result.Summary.Advisory != 2 || result.Summary.Critical != 1 || result.Summary.Warnings != 2 {
		t.Fatalf("unexpected summary: %+v", result.Summary)
	}
	for _, finding := range result.Findings {
		if finding.Evidence == nil {
			t.Fatalf("finding evidence should be non-nil: %#v", finding)
		}
	}
}
