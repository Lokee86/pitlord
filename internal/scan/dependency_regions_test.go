package scan

import (
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestSemanticDependencyRegionsPreferLanguageNamespaceOverDirectory(t *testing.T) {
	graph := arcana.Graph{
		Sources: []arcana.Node{
			{NodeID: 1, Key: "file-a", Kind: "file", Path: "src/Api/A.cs", Name: "A.cs"},
			{NodeID: 2, Key: "file-b", Kind: "file", Path: "src/Utilities/B.cs", Name: "B.cs"},
			{NodeID: 3, Key: "namespace", Kind: "namespace", Path: "src/Api/A.cs", Name: "Product"},
		},
		Outgoing: map[uint32][]arcana.Relationship{
			3: {
				{Relation: "contains", Node: arcana.Node{NodeID: 4, Kind: "type", Path: "src/Api/A.cs", Span: &arcana.Span{Path: "src/Api/A.cs"}}},
				{Relation: "contains", Node: arcana.Node{NodeID: 5, Kind: "type", Path: "src/Utilities/B.cs", Span: &arcana.Span{Path: "src/Utilities/B.cs"}}},
			},
		},
	}

	regions := semanticDependencyRegions(graph, repositoryFilePaths(graph.Sources))
	a := dependencyRegionForFile("src/Api/A.cs", regions)
	b := dependencyRegionForFile("src/Utilities/B.cs", regions)
	if a.id != b.id {
		t.Fatalf("same language namespace should define one architectural region: %#v %#v", a, b)
	}
	if a.label != "namespace Product" {
		t.Fatalf("unexpected semantic region label: %#v", a)
	}
}

func TestSemanticDependencyRegionsUseMultiFileModuleOwnership(t *testing.T) {
	graph := arcana.Graph{
		Sources: []arcana.Node{
			{NodeID: 1, Key: "file-a", Kind: "file", Path: "internal/net/a.go", Name: "a.go"},
			{NodeID: 2, Key: "file-b", Kind: "file", Path: "internal/net/b.go", Name: "b.go"},
			{NodeID: 3, Key: "module", Kind: "module", Path: "internal/net", Name: "net"},
		},
		Outgoing: map[uint32][]arcana.Relationship{
			3: {
				{Relation: "defines", Node: arcana.Node{NodeID: 4, Kind: "function", Path: "internal/net/a.go", Span: &arcana.Span{Path: "internal/net/a.go"}}},
				{Relation: "defines", Node: arcana.Node{NodeID: 5, Kind: "function", Path: "internal/net/b.go", Span: &arcana.Span{Path: "internal/net/b.go"}}},
			},
		},
	}

	regions := semanticDependencyRegions(graph, repositoryFilePaths(graph.Sources))
	a := dependencyRegionForFile("internal/net/a.go", regions)
	b := dependencyRegionForFile("internal/net/b.go", regions)
	if a.id != b.id || a.label != "module internal/net" {
		t.Fatalf("multi-file module should define one architectural region: %#v %#v", a, b)
	}
}
