package scan

import (
	"context"
	"fmt"
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestUnobservedOutcomeAnalyzerUsesNormalizedFactsAcrossLanguages(t *testing.T) {
	capabilities := "semantic-capabilities:%s:calls,source-spans,outcome-obligations"
	nodes := []arcana.Node{
		{NodeID: 1, Kind: "protocol", Path: "src/lib.rs", Name: fmt.Sprintf(capabilities, "rust")},
		{NodeID: 2, Kind: "protocol", Path: "src/lib.rs", Name: "outcome-operation:rust:fallible", Identity: "rust-discarded", Span: &arcana.Span{Path: "src/lib.rs", StartLine: 4, StartColumn: 5, EndLine: 4, EndColumn: 15}},
		{NodeID: 3, Kind: "protocol", Path: "src/app.ts", Name: fmt.Sprintf(capabilities, "typescript")},
		{NodeID: 4, Kind: "protocol", Path: "src/app.ts", Name: "outcome-operation:typescript:async", Identity: "ts-observed"},
		{NodeID: 5, Kind: "protocol", Path: "src/app.ts", Name: "outcome-action:consume"},
		{NodeID: 6, Kind: "protocol", Path: "src/app.py", Name: fmt.Sprintf(capabilities, "python")},
		{NodeID: 7, Kind: "protocol", Path: "src/app.py", Name: "outcome-operation:python:async", Identity: "python-discarded", Span: &arcana.Span{Path: "src/app.py", StartLine: 7, StartColumn: 5, EndLine: 7, EndColumn: 11}},
	}
	graph := arcana.Graph{Sources: nodes, Outgoing: map[uint32][]arcana.Relationship{
		4: {{Relation: "contains", Node: nodes[4]}},
	}}
	result, err := analyzeWithAnalyzers(context.Background(), AnalyzerContext{Graph: graph}, []Analyzer{unobservedOutcomeAnalyzer{}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Findings) != 2 {
		t.Fatalf("expected Rust and Python discarded outcomes, got %+v", result.Findings)
	}
	languages := map[string]bool{}
	for _, finding := range result.Findings {
		languages[finding.Language] = true
		if finding.RuleID != AnalyzerUnobservedOutcome || finding.Location == nil {
			t.Fatalf("unexpected finding: %+v", finding)
		}
	}
	if !languages["rust"] || !languages["python"] || languages["typescript"] {
		t.Fatalf("unexpected outcome findings: %+v", result.Findings)
	}
}
