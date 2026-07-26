package policy

import (
	"reflect"
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestForbidAreaCyclesDetectsStronglyConnectedAreas(t *testing.T) {
	document := cycleDocument()
	if err := document.NormalizeAndValidate(); err != nil {
		t.Fatal(err)
	}

	a := arcana.Node{NodeID: 1, Kind: "function", Path: "src/a/a.go", Name: "A"}
	b := arcana.Node{NodeID: 2, Kind: "function", Path: "src/b/b.go", Name: "B"}
	c := arcana.Node{NodeID: 3, Kind: "function", Path: "src/c/c.go", Name: "C"}
	graph := arcana.Graph{
		Sources: []arcana.Node{a, b, c},
		Outgoing: map[uint32][]arcana.Relationship{
			1: {{Relation: "calls", Node: b}},
			2: {{Relation: "calls", Node: c}},
			3: {{Relation: "calls", Node: a}},
		},
		Relationships: 3,
	}

	diagnostics := Evaluate(document, graph)
	if len(diagnostics) != 1 {
		t.Fatalf("expected one cycle diagnostic, got %+v", diagnostics)
	}
	if len(diagnostics[0].Evidence) != 3 {
		t.Fatalf("expected one representative edge per area pair, got %+v", diagnostics[0].Evidence)
	}
	for _, evidence := range diagnostics[0].Evidence {
		if evidence.Issue != "area_cycle" {
			t.Fatalf("unexpected issue %q", evidence.Issue)
		}
		if !reflect.DeepEqual(evidence.Areas, []string{"a", "b", "c"}) {
			t.Fatalf("unexpected component areas: %v", evidence.Areas)
		}
		if evidence.SourceArea == "" || evidence.TargetArea == "" || evidence.Target == nil {
			t.Fatalf("incomplete cycle evidence: %+v", evidence)
		}
	}
}

func TestForbidAreaCyclesAllowsAcyclicAreaGraph(t *testing.T) {
	document := cycleDocument()
	if err := document.NormalizeAndValidate(); err != nil {
		t.Fatal(err)
	}
	a := arcana.Node{NodeID: 1, Kind: "function", Path: "src/a/a.go", Name: "A"}
	b := arcana.Node{NodeID: 2, Kind: "function", Path: "src/b/b.go", Name: "B"}
	c := arcana.Node{NodeID: 3, Kind: "function", Path: "src/c/c.go", Name: "C"}
	graph := arcana.Graph{
		Sources: []arcana.Node{a, b, c},
		Outgoing: map[uint32][]arcana.Relationship{
			1: {{Relation: "calls", Node: b}},
			2: {{Relation: "calls", Node: c}},
		},
		Relationships: 2,
	}
	if diagnostics := Evaluate(document, graph); len(diagnostics) != 0 {
		t.Fatalf("expected acyclic graph to pass, got %+v", diagnostics)
	}
}

func TestForbidAreaCyclesUsesDeterministicRepresentativeEdges(t *testing.T) {
	document := Document{
		Version: Version,
		Areas: []Area{
			{ID: "a", Paths: []string{"src/a"}},
			{ID: "b", Paths: []string{"src/b"}},
		},
		Rules: []Rule{{ID: "no-cycles", Type: RuleForbidAreaCycles, Relations: []string{"calls"}}},
	}
	if err := document.NormalizeAndValidate(); err != nil {
		t.Fatal(err)
	}
	aLate := arcana.Node{NodeID: 1, Kind: "function", Path: "src/a/z.go", Name: "Late"}
	aEarly := arcana.Node{NodeID: 2, Kind: "function", Path: "src/a/a.go", Name: "Early"}
	b := arcana.Node{NodeID: 3, Kind: "function", Path: "src/b/b.go", Name: "B"}
	graph := arcana.Graph{
		Sources: []arcana.Node{aLate, aEarly, b},
		Outgoing: map[uint32][]arcana.Relationship{
			1: {{Relation: "calls", Node: b}},
			2: {{Relation: "calls", Node: b}},
			3: {{Relation: "calls", Node: aEarly}},
		},
	}
	diagnostics := Evaluate(document, graph)
	if len(diagnostics) != 1 || len(diagnostics[0].Evidence) != 2 {
		t.Fatalf("unexpected diagnostics: %+v", diagnostics)
	}
	for _, evidence := range diagnostics[0].Evidence {
		if evidence.SourceArea == "a" && evidence.Source.NodeID != aEarly.NodeID {
			t.Fatalf("expected stable earliest a-to-b representative, got %+v", evidence)
		}
	}
}

func TestAreaCyclePrefixesLoadSelectedAreasAndRelationships(t *testing.T) {
	document := cycleDocument()
	if err := document.NormalizeAndValidate(); err != nil {
		t.Fatal(err)
	}
	expected := []string{"src/a", "src/b", "src/c"}
	if prefixes := SourcePrefixes(document); !reflect.DeepEqual(prefixes, expected) {
		t.Fatalf("unexpected source prefixes: %v", prefixes)
	}
	if prefixes := RelationshipSourcePrefixes(document); !reflect.DeepEqual(prefixes, expected) {
		t.Fatalf("unexpected relationship prefixes: %v", prefixes)
	}
}

func cycleDocument() Document {
	return Document{
		Version: Version,
		Areas: []Area{
			{ID: "a", Paths: []string{"src/a"}},
			{ID: "b", Paths: []string{"src/b"}},
			{ID: "c", Paths: []string{"src/c"}},
		},
		Rules: []Rule{
			{ID: "no-area-cycles", Type: RuleForbidAreaCycles, Relations: []string{"calls"}},
		},
	}
}
