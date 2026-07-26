package policy

import (
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestAnalyzeAreasReportsCoverageDependenciesAndCycles(t *testing.T) {
	document := Document{
		Version: Version,
		Areas: []Area{
			{ID: "api", Paths: []string{"src/api"}},
			{ID: "storage", Paths: []string{"src/storage"}},
		},
		Rules: []Rule{{ID: "ownership", Type: RuleRequireOwnership, ScopePaths: []string{"src"}}},
	}
	if err := document.NormalizeAndValidate(); err != nil {
		t.Fatal(err)
	}
	apiFile := arcana.Node{NodeID: 1, Kind: "file", Path: "src/api/api.go", Name: "api.go"}
	apiFunction := arcana.Node{NodeID: 2, Kind: "function", Path: "src/api/api.go", Name: "Create"}
	storageFile := arcana.Node{NodeID: 3, Kind: "file", Path: "src/storage/store.go", Name: "store.go"}
	storageFunction := arcana.Node{NodeID: 4, Kind: "function", Path: "src/storage/store.go", Name: "Write"}
	sharedFile := arcana.Node{NodeID: 5, Kind: "file", Path: "src/shared.go", Name: "shared.go"}
	external := arcana.Node{NodeID: 6, Kind: "function", Path: "vendor/log.go", Name: "Log"}
	graph := arcana.Graph{
		Sources: []arcana.Node{apiFile, apiFunction, storageFile, storageFunction, sharedFile},
		Outgoing: map[uint32][]arcana.Relationship{
			2: {
				{Relation: "calls", Node: storageFunction},
				{Relation: "calls", Node: external},
			},
			4: {{Relation: "calls", Node: apiFunction}},
		},
		Relationships: 3,
	}

	analysis := AnalyzeAreas(document, graph, AreaAnalysisOptions{ScopePaths: []string{"src"}})
	if analysis.Ownership.CheckedNodes != 3 || analysis.Ownership.UnownedNodes != 1 {
		t.Fatalf("unexpected ownership: %+v", analysis.Ownership)
	}
	if len(analysis.Dependencies) != 3 {
		t.Fatalf("expected api->storage, api->external, storage->api, got %+v", analysis.Dependencies)
	}
	if len(analysis.Cycles) != 1 || len(analysis.Cycles[0]) != 2 {
		t.Fatalf("expected api/storage cycle, got %v", analysis.Cycles)
	}
	if analysis.Summary.RelationshipsScanned != 3 {
		t.Fatalf("unexpected summary: %+v", analysis.Summary)
	}
}

func TestAreaPathsAreUniqueAndSorted(t *testing.T) {
	document := Document{Areas: []Area{
		{ID: "b", Paths: []string{"src/b", "src/common"}},
		{ID: "a", Paths: []string{"src/a", "src/common"}},
	}}
	paths := AreaPaths(document)
	if len(paths) != 3 || paths[0] != "src/a" || paths[2] != "src/common" {
		t.Fatalf("unexpected paths: %v", paths)
	}
}
