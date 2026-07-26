package policy

import (
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestRequireOwnershipFindsUnownedAndOverlappingNodes(t *testing.T) {
	document := Document{
		Version: Version,
		Areas: []Area{
			{ID: "api", Paths: []string{"src/api"}, Kinds: []string{"file"}},
			{ID: "api-shadow", Paths: []string{"src/api"}, Kinds: []string{"file"}},
			{ID: "storage", Paths: []string{"src/storage"}, Kinds: []string{"file"}},
		},
		Rules: []Rule{
			{
				ID:          "source-ownership",
				Type:        RuleRequireOwnership,
				ScopePaths:  []string{"src"},
				SourceKinds: []string{"file"},
			},
		},
	}
	if err := document.NormalizeAndValidate(); err != nil {
		t.Fatal(err)
	}

	graph := arcana.Graph{
		Sources: []arcana.Node{
			{NodeID: 1, Kind: "file", Path: "src/api/api.go", Name: "api.go"},
			{NodeID: 2, Kind: "file", Path: "src/shared/shared.go", Name: "shared.go"},
			{NodeID: 3, Kind: "file", Path: "src/storage/store.go", Name: "store.go"},
			{NodeID: 4, Kind: "function", Path: "src/shared/shared.go", Name: "Shared"},
		},
		Outgoing: map[uint32][]arcana.Relationship{},
	}

	diagnostics := Evaluate(document, graph)
	if len(diagnostics) != 1 {
		t.Fatalf("expected one grouped diagnostic, got %+v", diagnostics)
	}
	if len(diagnostics[0].Evidence) != 2 {
		t.Fatalf("expected two ownership findings, got %+v", diagnostics[0].Evidence)
	}
	issues := map[string]bool{}
	for _, evidence := range diagnostics[0].Evidence {
		issues[evidence.Issue] = true
	}
	if !issues["unowned"] || !issues["multiple_owners"] {
		t.Fatalf("unexpected ownership issues: %v", issues)
	}
	prefixes := SourcePrefixes(document)
	if len(prefixes) != 1 || prefixes[0] != "src" {
		t.Fatalf("unexpected ownership prefixes: %v", prefixes)
	}
	if relationshipPrefixes := RelationshipSourcePrefixes(document); len(relationshipPrefixes) != 0 {
		t.Fatalf("ownership-only policy should not load relationships: %v", relationshipPrefixes)
	}
}

func TestRequireOwnershipCanAllowOneFailureClass(t *testing.T) {
	document := Document{
		Version: Version,
		Areas: []Area{
			{ID: "one", Paths: []string{"src/owned"}},
			{ID: "two", Paths: []string{"src/owned"}},
		},
		Rules: []Rule{
			{
				ID:            "ownership",
				Type:          RuleRequireOwnership,
				ScopePaths:    []string{"src"},
				SourceKinds:   []string{"file"},
				AllowOverlaps: true,
			},
		},
	}
	if err := document.NormalizeAndValidate(); err != nil {
		t.Fatal(err)
	}
	graph := arcana.Graph{Sources: []arcana.Node{
		{NodeID: 1, Kind: "file", Path: "src/owned/a.go"},
		{NodeID: 2, Kind: "file", Path: "src/unowned/b.go"},
	}}
	diagnostics := Evaluate(document, graph)
	if len(diagnostics) != 1 || len(diagnostics[0].Evidence) != 1 || diagnostics[0].Evidence[0].Issue != "unowned" {
		t.Fatalf("unexpected diagnostics: %+v", diagnostics)
	}
}
