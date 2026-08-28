package scan

import (
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestImpactBlastRadiusFindsTransitiveAmplification(t *testing.T) {
	graph := dependencyTestGraph(64)
	addCalibrationNamespaceRegions(&graph, 8)
	for source := uint32(2); source <= 64; source++ {
		addDependency(&graph, source, source/2, "depends-on")
	}

	findings := detectImpactBlastRadius(graph, ".")
	if len(findings) == 0 {
		t.Fatal("expected transitive impact hotspots")
	}
	seenRoot := false
	for _, finding := range findings {
		if finding.Detector != DetectorImpactBlastRadius || finding.Disposition != DispositionAdvisory {
			t.Fatalf("unexpected impact finding: %+v", finding)
		}
		if finding.Scope.Path == filePath(1) {
			seenRoot = true
		}
	}
	if !seenRoot {
		t.Fatalf("expected low-direct-fan-in foundation %s, got %#v", filePath(1), findings)
	}
}

func TestImpactBlastRadiusKeepsPureDirectSharedLeafQuiet(t *testing.T) {
	graph := dependencyTestGraph(64)
	addCalibrationNamespaceRegions(&graph, 8)
	for source := uint32(2); source <= 40; source++ {
		addDependency(&graph, source, 1, "references")
	}
	if findings := detectImpactBlastRadius(graph, "."); len(findings) != 0 {
		t.Fatalf("direct fan-in without transitive amplification belongs to centrality analysis: %#v", findings)
	}
}

func TestImpactBlastRadiusFindsLayeredFoundation(t *testing.T) {
	graph := layeredCalibrationGraph()
	addCalibrationNamespaceRegions(&graph, 8)
	findings := detectImpactBlastRadius(graph, ".")
	if len(findings) == 0 {
		t.Fatal("deep layered topology should expose foundational blast radius")
	}
	for _, finding := range findings {
		if finding.Detector != DetectorImpactBlastRadius {
			t.Fatalf("unexpected detector: %+v", finding)
		}
	}
}

func TestImpactBlastRadiusKeepsFamilyWideReachabilityQuiet(t *testing.T) {
	controls := []struct {
		name  string
		graph func() arcana.Graph
	}{
		{name: "modular", graph: modularCalibrationGraph},
		{name: "entangled", graph: entangledCalibrationGraph},
		{name: "hub-heavy", graph: hubHeavyCalibrationGraph},
		{name: "dense-subsystem", graph: denseSubsystemCalibrationGraph},
	}
	for _, control := range controls {
		t.Run(control.name, func(t *testing.T) {
			graph := control.graph()
			addCalibrationNamespaceRegions(&graph, 8)
			if findings := detectImpactBlastRadius(graph, "."); len(findings) != 0 {
				t.Fatalf("family-wide reachability should not be isolated as a hotspot: %#v", findings)
			}
		})
	}
}
