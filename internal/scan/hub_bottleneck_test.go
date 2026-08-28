package scan

import (
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestHubBottleneckFindsManyToManyBehavioralWaists(t *testing.T) {
	graph := behavioralHubCalibrationGraph()
	addCalibrationNamespaceRegions(&graph, 8)

	findings := detectHubBottlenecks(graph, ".")
	if len(findings) != 4 {
		t.Fatalf("expected four behavioral bottlenecks, got %#v", findings)
	}
	for _, finding := range findings {
		if finding.Detector != DetectorHubBottleneck || finding.Scope.Kind != "file" {
			t.Fatalf("unexpected bottleneck finding: %+v", finding)
		}
		if finding.Disposition != DispositionAdvisory {
			t.Fatalf("bottleneck must remain advisory: %+v", finding)
		}
	}
}

func TestHubBottleneckKeepsReferenceHubsQuiet(t *testing.T) {
	for _, graph := range []arcana.Graph{hubHeavyCalibrationGraph(), entangledCalibrationGraph()} {
		addCalibrationNamespaceRegions(&graph, 8)
		if findings := detectHubBottlenecks(graph, "."); len(findings) != 0 {
			t.Fatalf("reference-driven hub centrality is not a behavioral bottleneck: %#v", findings)
		}
	}
}

func TestHubBottleneckKeepsPureSharedLeafQuiet(t *testing.T) {
	graph := dependencyTestGraph(32)
	addCalibrationNamespaceRegions(&graph, 8)
	for source := uint32(9); source <= 32; source++ {
		addDependency(&graph, source, 1, "calls")
	}

	if findings := detectHubBottlenecks(graph, "."); len(findings) != 0 {
		t.Fatalf("high-fan-in leaf abstraction should remain quiet: %#v", findings)
	}
}

func TestHubBottleneckKeepsNarrowBehavioralFacadeQuiet(t *testing.T) {
	graph := dependencyTestGraph(64)
	addCalibrationNamespaceRegions(&graph, 8)
	for source := uint32(9); source <= 64; source++ {
		addDependency(&graph, source, 1, "calls")
	}
	for target := uint32(2); target <= 8; target++ {
		addDependency(&graph, 1, target, "calls")
	}

	if findings := detectHubBottlenecks(graph, "."); len(findings) != 0 {
		t.Fatalf("cross-cutting facade with narrow outgoing ownership should remain quiet: %#v", findings)
	}
}

func TestHubBottleneckKeepsCompositionRootQuiet(t *testing.T) {
	graph := dependencyTestGraph(32)
	addCalibrationNamespaceRegions(&graph, 8)
	for target := uint32(2); target <= 24; target++ {
		addDependency(&graph, 1, target, "calls")
	}
	addDependency(&graph, 2, 1, "calls")

	if findings := detectHubBottlenecks(graph, "."); len(findings) != 0 {
		t.Fatalf("high-fan-out low-fan-in composition root belongs to dependency pressure: %#v", findings)
	}
}

func TestHubBottleneckKeepsHealthyTopologyFamiliesQuiet(t *testing.T) {
	controls := []struct {
		name  string
		graph arcana.Graph
	}{
		{name: "modular", graph: modularCalibrationGraph()},
		{name: "layered", graph: layeredCalibrationGraph()},
		{name: "dense-subsystem", graph: denseSubsystemCalibrationGraph()},
	}
	for _, control := range controls {
		t.Run(control.name, func(t *testing.T) {
			graph := control.graph
			addCalibrationNamespaceRegions(&graph, 8)
			if findings := detectHubBottlenecks(graph, "."); len(findings) != 0 {
				t.Fatalf("%s topology should remain quiet: %#v", control.name, findings)
			}
		})
	}
}

func behavioralHubCalibrationGraph() arcana.Graph {
	graph := dependencyTestGraph(64)
	for target := uint32(5); target <= 64; target++ {
		hub := (target-5)%4 + 1
		addDependency(&graph, hub, target, "calls")
		addDependency(&graph, target, hub, "calls")
		next := target + 1
		if next > 64 {
			next = 5
		}
		addDependency(&graph, target, next, "depends-on")
	}
	return graph
}
