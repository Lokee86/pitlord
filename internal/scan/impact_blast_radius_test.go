package scan

import (
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestImpactBlastRadiusFindsIndependentTransitiveAmplification(t *testing.T) {
	graph := balancedImpactBranchGraph()
	findings := detectImpactBlastRadius(graph, ".")
	if len(findings) == 0 {
		t.Fatal("expected independent transitive impact hotspot")
	}
	for _, finding := range findings {
		if finding.Detector != DetectorImpactBlastRadius || finding.Disposition != DispositionAdvisory {
			t.Fatalf("unexpected impact finding: %+v", finding)
		}
		if finding.Scope.Path == filePath(1) {
			return
		}
	}
	t.Fatalf("expected balanced three-branch foundation %s, got %#v", filePath(1), findings)
}

func TestImpactBlastRadiusDoesNotTransferGatewayImpactToFoundation(t *testing.T) {
	graph := singleGatewayImpactGraph()
	findings := detectImpactBlastRadius(graph, ".")
	seenGateway := false
	for _, finding := range findings {
		if finding.Scope.Path == filePath(1) {
			t.Fatalf("gateway blast radius must not transfer to foundation %s: %#v", filePath(1), findings)
		}
		if finding.Scope.Path == filePath(2) {
			seenGateway = true
		}
	}
	if !seenGateway {
		t.Fatalf("expected actual three-branch gateway %s to remain visible: %#v", filePath(2), findings)
	}
}

func TestImpactBlastRadiusKeepsDominantBranchAmplificationQuiet(t *testing.T) {
	graph := dominantImpactBranchGraph()
	if findings := detectImpactBlastRadius(graph, "."); len(findings) != 0 {
		t.Fatalf("one dominant branch with two small side branches is not independent broad impact: %#v", findings)
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

func TestImpactBlastRadiusKeepsHealthyLayeredFoundationQuiet(t *testing.T) {
	graph := layeredCalibrationGraph()
	addCalibrationNamespaceRegions(&graph, 8)
	if findings := detectImpactBlastRadius(graph, "."); len(findings) != 0 {
		t.Fatalf("ordinary layered reuse should not be isolated as a blast-radius hotspot: %#v", findings)
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

func balancedImpactBranchGraph() arcana.Graph {
	graph := dependencyTestGraph(64)
	addCalibrationNamespaceRegions(&graph, 8)
	addDependency(&graph, 2, 1, "depends-on")
	addDependency(&graph, 3, 1, "depends-on")
	addDependency(&graph, 4, 1, "depends-on")
	addImpactChain(&graph, 2, 5, 23)
	addImpactChain(&graph, 3, 24, 43)
	addImpactChain(&graph, 4, 44, 64)
	return graph
}

func singleGatewayImpactGraph() arcana.Graph {
	graph := dependencyTestGraph(64)
	addCalibrationNamespaceRegions(&graph, 8)
	addDependency(&graph, 2, 1, "depends-on")
	addDependency(&graph, 3, 2, "depends-on")
	addDependency(&graph, 4, 2, "depends-on")
	addDependency(&graph, 5, 2, "depends-on")
	addImpactChain(&graph, 3, 6, 24)
	addImpactChain(&graph, 4, 25, 44)
	addImpactChain(&graph, 5, 45, 64)
	return graph
}

func dominantImpactBranchGraph() arcana.Graph {
	graph := dependencyTestGraph(64)
	addCalibrationNamespaceRegions(&graph, 8)
	addDependency(&graph, 2, 1, "depends-on")
	addDependency(&graph, 3, 1, "depends-on")
	addDependency(&graph, 4, 1, "depends-on")
	addImpactChain(&graph, 2, 5, 55)
	addImpactChain(&graph, 3, 56, 59)
	addImpactChain(&graph, 4, 60, 64)
	return graph
}

func addImpactChain(graph *arcana.Graph, root, start, end uint32) {
	parent := root
	for node := start; node <= end; node++ {
		addDependency(graph, node, parent, "depends-on")
		parent = node
	}
}
