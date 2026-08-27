package scan

import (
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestDependencyPressureCalibrationAcrossArcanaTopologyFamilies(t *testing.T) {
	tests := []struct {
		name      string
		graph     arcana.Graph
		wantHubs  []uint32
		wantQuiet bool
	}{
		{name: "modular", graph: modularCalibrationGraph(), wantQuiet: true},
		{name: "entangled", graph: entangledCalibrationGraph(), wantHubs: []uint32{1, 2, 3, 4}},
		{name: "hub-heavy", graph: hubHeavyCalibrationGraph(), wantHubs: []uint32{1, 2, 3, 4}},
		{name: "layered", graph: layeredCalibrationGraph(), wantQuiet: true},
		{name: "dense-subsystem", graph: denseSubsystemCalibrationGraph(), wantQuiet: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			findings := detectDependencyPressure(test.graph, ".")
			if test.wantQuiet {
				if len(findings) != 0 {
					t.Fatalf("%s topology should not be classified as isolated hub pressure: %#v", test.name, findings)
				}
				return
			}
			if len(findings) != len(test.wantHubs) {
				t.Fatalf("%s topology: expected %d hub findings, got %#v", test.name, len(test.wantHubs), findings)
			}
			seen := make(map[string]struct{}, len(findings))
			for _, finding := range findings {
				seen[finding.Scope.Path] = struct{}{}
			}
			for _, hub := range test.wantHubs {
				if _, ok := seen[filePath(hub)]; !ok {
					t.Fatalf("%s topology: missing expected hub %s in %#v", test.name, filePath(hub), findings)
				}
			}
		})
	}
}

// These fixtures preserve the structural signatures and default 64-node scale
// of Arcana's deterministic topology families. They intentionally use Pitlord's
// normalized graph model rather than importing language or parser concepts.
func modularCalibrationGraph() arcana.Graph {
	graph := dependencyTestGraph(64)
	for source := uint32(1); source <= 64; source++ {
		clusterStart := ((source-1)/8)*8 + 1
		local := (source - clusterStart) % 8
		for step := uint32(1); step <= 3; step++ {
			target := clusterStart + (local+step)%8
			addDependency(&graph, source, target, "depends-on")
		}
		if source%4 == 0 {
			addDependency(&graph, source, ((source+7)%64)+1, "imports")
		}
	}
	return graph
}

func hubHeavyCalibrationGraph() arcana.Graph {
	graph := dependencyTestGraph(64)
	for target := uint32(5); target <= 64; target++ {
		hub := (target-5)%4 + 1
		addDependency(&graph, hub, target, "calls")
		addDependency(&graph, target, hub, "references")
		next := target + 1
		if next > 64 {
			next = 5
		}
		addDependency(&graph, target, next, "references")
	}
	return graph
}

func entangledCalibrationGraph() arcana.Graph {
	graph := dependencyTestGraph(64)
	for target := uint32(5); target <= 64; target++ {
		hub := (target-5)%4 + 1
		addDependency(&graph, hub, target, "calls")
		addDependency(&graph, target, hub, "references")
		for _, offset := range []uint32{1, 7, 13} {
			other := 5 + ((target - 5 + offset) % 60)
			addDependency(&graph, target, other, "depends-on")
		}
	}
	return graph
}

func layeredCalibrationGraph() arcana.Graph {
	graph := dependencyTestGraph(64)
	for layer := uint32(0); layer < 7; layer++ {
		for position := uint32(0); position < 8; position++ {
			source := layer*8 + position + 1
			for offset := uint32(0); offset < 3; offset++ {
				target := (layer+1)*8 + (position+offset)%8 + 1
				addDependency(&graph, source, target, "depends-on")
			}
		}
	}
	return graph
}

func denseSubsystemCalibrationGraph() arcana.Graph {
	graph := dependencyTestGraph(64)
	for source := uint32(1); source <= 16; source++ {
		for target := uint32(1); target <= 16; target++ {
			if source != target {
				addDependency(&graph, source, target, "calls")
			}
		}
		addDependency(&graph, source, 16+source, "references")
	}
	for source := uint32(17); source <= 64; source++ {
		next := source + 1
		if next > 64 {
			next = 17
		}
		addDependency(&graph, source, next, "references")
	}
	return graph
}
