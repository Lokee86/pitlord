package scan

import (
	"strings"
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestDependencyPressureFindsLanguageNeutralFileHub(t *testing.T) {
	graph := dependencyTestGraph(10)
	for target := uint32(2); target <= 10; target++ {
		addDependency(&graph, 1, target, []string{"imports", "calls", "implements"}[int(target)%3])
	}
	addDependency(&graph, 1, 8, "contains")
	for source := uint32(2); source < 10; source++ {
		addDependency(&graph, source, source+1, "references")
	}

	findings := detectDependencyPressure(graph, ".")
	if len(findings) != 1 {
		t.Fatalf("expected one hub finding, got %#v", findings)
	}
	finding := findings[0]
	if finding.Detector != DetectorDependencyPressure || finding.Scope.Path != "src/file01.any" {
		t.Fatalf("unexpected finding: %+v", finding)
	}
	if finding.Disposition != DispositionAdvisory || finding.Severity != SeverityHigh {
		t.Fatalf("unexpected judgment: %+v", finding)
	}
	if !strings.Contains(finding.RequiredOutcome, "below") || len(finding.Evidence) != 4 {
		t.Fatalf("finding is not mechanically actionable: %+v", finding)
	}
}

func TestDependencyPressureKeepsBalancedGraphQuiet(t *testing.T) {
	graph := dependencyTestGraph(10)
	for source := uint32(1); source <= 10; source++ {
		target := source + 1
		if target > 10 {
			target = 1
		}
		addDependency(&graph, source, target, "depends-on")
	}
	if findings := detectDependencyPressure(graph, "."); len(findings) != 0 {
		t.Fatalf("balanced graph should not produce hubs: %#v", findings)
	}
}

func TestDependencyPressureDeduplicatesSymbolEdgesByFile(t *testing.T) {
	graph := dependencyTestGraph(8)
	for repeat := 0; repeat < 20; repeat++ {
		addDependency(&graph, 1, 2, "calls")
	}
	for source := uint32(2); source < 8; source++ {
		addDependency(&graph, source, source+1, "references")
	}
	if findings := detectDependencyPressure(graph, "."); len(findings) != 0 {
		t.Fatalf("repeated symbol edges to one file must not look like broad coupling: %#v", findings)
	}
}

func TestDependencyPressureIgnoresRepositoryStructureNodes(t *testing.T) {
	graph := dependencyTestGraph(8)
	graph.Sources = append(graph.Sources, arcana.Node{NodeID: 99, Kind: "directory", Path: "src", Name: "src"})
	for target := uint32(1); target <= 8; target++ {
		graph.Outgoing[99] = append(graph.Outgoing[99], arcana.Relationship{
			Relation: "depends-on",
			Node:     arcana.Node{NodeID: target, Kind: "symbol", Path: filePath(target), Name: "target"},
		})
	}
	for source := uint32(1); source < 8; source++ {
		addDependency(&graph, source, source+1, "references")
	}
	if findings := detectDependencyPressure(graph, "."); len(findings) != 0 {
		t.Fatalf("directory/repository nodes must not become code dependency units: %#v", findings)
	}
}

func dependencyTestGraph(fileCount uint32) arcana.Graph {
	graph := arcana.Graph{Outgoing: map[uint32][]arcana.Relationship{}}
	kinds := []string{"function", "type", "module", "symbol"}
	for id := uint32(1); id <= fileCount; id++ {
		graph.Sources = append(graph.Sources, arcana.Node{
			NodeID: id,
			Kind:   kinds[int(id)%len(kinds)],
			Path:   filePath(id),
			Name:   "node",
		})
	}
	return graph
}

func addDependency(graph *arcana.Graph, source, target uint32, relation string) {
	graph.Outgoing[source] = append(graph.Outgoing[source], arcana.Relationship{
		Relation: relation,
		Node:     arcana.Node{NodeID: target, Path: filePath(target), Kind: "symbol", Name: "target"},
	})
}

func filePath(id uint32) string {
	return "src/file" + twoDigits(id) + ".any"
}

func twoDigits(value uint32) string {
	if value < 10 {
		return "0" + string(rune('0'+value))
	}
	return "10"
}
