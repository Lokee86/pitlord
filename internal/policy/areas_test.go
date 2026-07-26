package policy

import (
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestEvaluateDependencyRuleUsesNamedAreas(t *testing.T) {
	document := Document{
		Version: Version,
		Areas: []Area{
			{ID: "api", Paths: []string{"src/api"}},
			{ID: "storage", Paths: []string{"src/storage"}},
		},
		Rules: []Rule{
			{
				ID:        "api-must-not-access-storage",
				FromAreas: []string{"api"},
				ToAreas:   []string{"storage"},
				Relations: []string{"calls"},
			},
		},
	}
	if err := document.NormalizeAndValidate(); err != nil {
		t.Fatal(err)
	}

	source := arcana.Node{NodeID: 1, Kind: "function", Path: "src/api/create.go", Name: "Create"}
	target := arcana.Node{NodeID: 2, Kind: "function", Path: "src/storage/write.go", Name: "Write"}
	graph := arcana.Graph{
		Sources: []arcana.Node{source},
		Outgoing: map[uint32][]arcana.Relationship{
			1: {{Relation: "calls", Node: target}},
		},
		Relationships: 1,
	}

	prefixes := SourcePrefixes(document)
	if len(prefixes) != 1 || prefixes[0] != "src/api" {
		t.Fatalf("unexpected source prefixes: %v", prefixes)
	}
	diagnostics := Evaluate(document, graph)
	if len(diagnostics) != 1 || len(diagnostics[0].Evidence) != 1 {
		t.Fatalf("unexpected diagnostics: %+v", diagnostics)
	}
	if diagnostics[0].Evidence[0].Target == nil || diagnostics[0].Evidence[0].Target.NodeID != target.NodeID {
		t.Fatalf("unexpected target evidence: %+v", diagnostics[0].Evidence[0])
	}
}

func TestAreaExclusionsAndKindsAreApplied(t *testing.T) {
	document := Document{
		Version: Version,
		Areas: []Area{
			{ID: "api", Paths: []string{"src/api"}, ExcludePaths: []string{"src/api/generated"}, Kinds: []string{"function"}},
			{ID: "storage", Paths: []string{"src/storage"}, Kinds: []string{"function"}},
		},
		Rules: []Rule{{ID: "rule", FromAreas: []string{"api"}, ToAreas: []string{"storage"}, Relations: []string{"calls"}}},
	}
	if err := document.NormalizeAndValidate(); err != nil {
		t.Fatal(err)
	}
	target := arcana.Node{NodeID: 10, Kind: "function", Path: "src/storage/write.go", Name: "Write"}
	sources := []arcana.Node{
		{NodeID: 1, Kind: "function", Path: "src/api/live.go", Name: "Live"},
		{NodeID: 2, Kind: "function", Path: "src/api/generated/client.go", Name: "Generated"},
		{NodeID: 3, Kind: "file", Path: "src/api/live.go", Name: "api"},
	}
	graph := arcana.Graph{Sources: sources, Outgoing: map[uint32][]arcana.Relationship{}}
	for _, source := range sources {
		graph.Outgoing[source.NodeID] = []arcana.Relationship{{Relation: "calls", Node: target}}
	}
	diagnostics := Evaluate(document, graph)
	if len(diagnostics) != 1 || len(diagnostics[0].Evidence) != 1 {
		t.Fatalf("expected one filtered edge, got %+v", diagnostics)
	}
	if diagnostics[0].Evidence[0].Source.NodeID != 1 {
		t.Fatalf("unexpected source: %+v", diagnostics[0].Evidence[0].Source)
	}
}
