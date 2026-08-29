package scan

import "github.com/Lokee86/pitlord/internal/arcana"

func dependencyDepthTestGraph(paths []string) arcana.Graph {
	graph := arcana.Graph{Outgoing: map[uint32][]arcana.Relationship{}}
	for index, filePath := range paths {
		fileID := uint32(index + 1)
		graph.Sources = append(graph.Sources,
			arcana.Node{NodeID: fileID, Kind: "file", Path: filePath, Name: "file"},
			arcana.Node{NodeID: 10_000 + fileID, Kind: "symbol", Path: filePath, Name: "symbol"},
		)
	}
	return graph
}

func addDepthDependency(graph *arcana.Graph, paths []string, source, target int, relation string) {
	sourceID := uint32(10_001 + source)
	targetID := uint32(10_001 + target)
	graph.Outgoing[sourceID] = append(graph.Outgoing[sourceID], arcana.Relationship{
		Relation: relation,
		Node:     arcana.Node{NodeID: targetID, Kind: "symbol", Path: paths[target], Name: "target"},
	})
}
