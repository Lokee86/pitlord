package scan

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestSwallowedErrorAnalyzerUsesNormalizedSemanticFactsAcrossLanguages(t *testing.T) {
	graph := semanticFixtureGraph()
	result, err := analyzeWithAnalyzers(context.Background(), AnalyzerContext{Graph: graph}, []Analyzer{swallowedErrorAnalyzer{}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Findings) != 3 {
		t.Fatalf("expected Rust, TypeScript, and Python swallowed errors, got %+v", result.Findings)
	}
	languages := map[string]bool{}
	for _, finding := range result.Findings {
		languages[finding.Language] = true
		if finding.RuleID != AnalyzerSwallowedError || finding.Location == nil {
			t.Fatalf("unexpected finding: %+v", finding)
		}
	}
	if !languages["rust"] || !languages["typescript"] || !languages["python"] {
		t.Fatalf("missing cross-language findings: %+v", result.Findings)
	}
}

func TestSwallowedErrorAnalyzerRejectsMissingCapabilities(t *testing.T) {
	_, err := analyzeWithAnalyzers(context.Background(), AnalyzerContext{Graph: arcana.Graph{}}, []Analyzer{swallowedErrorAnalyzer{}})
	if err == nil || !strings.Contains(err.Error(), "requires semantic capabilities") {
		t.Fatalf("unexpected capability error: %v", err)
	}
}

func TestBuiltInSemanticAnalyzerRequiresDeclaredCapabilities(t *testing.T) {
	analyzers, err := BuiltInAnalyzers("semantic")
	if err != nil {
		t.Fatal(err)
	}
	if len(analyzers) != 2 || !AnalyzersRequireGraph(analyzers) {
		t.Fatalf("unexpected semantic analyzer registry: %+v", analyzers)
	}
	metadata := analyzers[0].Metadata()
	if metadata.ID != AnalyzerSwallowedError || len(metadata.Capabilities) != len(swallowedErrorCapabilities) {
		t.Fatalf("unexpected semantic analyzer metadata: %+v", metadata)
	}
	outcomeMetadata := analyzers[1].Metadata()
	if outcomeMetadata.ID != AnalyzerUnobservedOutcome || len(outcomeMetadata.Capabilities) != len(unobservedOutcomeCapabilities) {
		t.Fatalf("unexpected outcome analyzer metadata: %+v", outcomeMetadata)
	}
}

func semanticFixtureGraph() arcana.Graph {
	capabilities := "semantic-capabilities:%s:control-flow,error-handling,calls,source-spans"
	nodes := []arcana.Node{
		{NodeID: 1, Kind: "protocol", Path: "src/lib.rs", Name: fmt.Sprintf(capabilities, "rust")},
		{NodeID: 2, Kind: "protocol", Path: "src/lib.rs", Name: "error-handler:rust", Identity: "rust-swallowed", Span: &arcana.Span{Path: "src/lib.rs", StartLine: 4, StartColumn: 9, EndLine: 4, EndColumn: 20}},
		{NodeID: 3, Kind: "protocol", Path: "src/handled.rs", Name: "error-handler:rust", Identity: "rust-handled"},
		{NodeID: 4, Kind: "protocol", Path: "src/handled.rs", Name: fmt.Sprintf(capabilities, "rust")},
		{NodeID: 5, Kind: "protocol", Path: "src/handled.rs", Name: "error-action:record"},
		{NodeID: 6, Kind: "protocol", Path: "src/app.ts", Name: fmt.Sprintf(capabilities, "typescript")},
		{NodeID: 7, Kind: "protocol", Path: "src/app.ts", Name: "error-handler:typescript", Identity: "ts-swallowed", Span: &arcana.Span{Path: "src/app.ts", StartLine: 7, StartColumn: 18, EndLine: 7, EndColumn: 28}},
		{NodeID: 8, Kind: "protocol", Path: "src/app.py", Name: fmt.Sprintf(capabilities, "python")},
		{NodeID: 9, Kind: "protocol", Path: "src/app.py", Name: "error-handler:python", Identity: "python-swallowed", Span: &arcana.Span{Path: "src/app.py", StartLine: 6, StartColumn: 5, EndLine: 7, EndColumn: 13}},
	}
	return arcana.Graph{Sources: nodes, Outgoing: map[uint32][]arcana.Relationship{
		3: {{Relation: "contains", Node: nodes[4]}},
	}}
}
