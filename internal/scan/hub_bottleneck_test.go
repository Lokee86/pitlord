package scan

import (
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestHubBottleneckFindsCentralCoordinationHubs(t *testing.T) {
	graph := hubHeavyCalibrationGraph()
	addCalibrationNamespaceRegions(&graph, 8)

	findings := detectHubBottlenecks(graph, ".")
	if len(findings) != 4 {
		t.Fatalf("expected four central bottlenecks, got %#v", findings)
	}
	for _, finding := range findings {
		if finding.Detector != DetectorHubBottleneck || finding.Scope.Kind != "file" {
			t.Fatalf("unexpected bottleneck finding: %+v", finding)
		}
		if finding.Disposition != DispositionAdvisory {
			t.Fatalf("bottleneck must begin advisory: %+v", finding)
		}
	}
}

func TestHubBottleneckKeepsPureSharedLeafQuiet(t *testing.T) {
	graph := dependencyTestGraph(32)
	addCalibrationNamespaceRegions(&graph, 8)
	for source := uint32(9); source <= 32; source++ {
		addDependency(&graph, source, 1, "references")
	}

	if findings := detectHubBottlenecks(graph, "."); len(findings) != 0 {
		t.Fatalf("high-fan-in leaf abstraction should remain quiet: %#v", findings)
	}
}

func TestHubBottleneckKeepsCompositionRootQuiet(t *testing.T) {
	graph := dependencyTestGraph(32)
	addCalibrationNamespaceRegions(&graph, 8)
	for target := uint32(2); target <= 24; target++ {
		addDependency(&graph, 1, target, "depends-on")
	}
	addDependency(&graph, 2, 1, "references")

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
				t.Fatalf("%s topology should not create isolated coordination bottlenecks: %#v", control.name, findings)
			}
		})
	}
}

func TestHubBottleneckEntangledTopologyIsMeasuredExplicitly(t *testing.T) {
	graph := entangledCalibrationGraph()
	addCalibrationNamespaceRegions(&graph, 8)
	findings := detectHubBottlenecks(graph, ".")
	if len(findings) != 4 {
		t.Fatalf("entangled topology has four central coordination hubs, got %#v", findings)
	}
}
