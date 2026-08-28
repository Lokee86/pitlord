package scan

import (
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestDependencyKnotsCalibrationAcrossArcanaTopologyFamilies(t *testing.T) {
	tests := []struct {
		name      string
		graph     arcana.Graph
		wantKnots int
	}{
		{name: "modular", graph: modularCalibrationGraph(), wantKnots: 0},
		{name: "entangled", graph: directionalizeCalibrationReferences(entangledCalibrationGraph()), wantKnots: 1},
		{name: "hub-heavy", graph: hubHeavyCalibrationGraph(), wantKnots: 0},
		{name: "layered", graph: layeredCalibrationGraph(), wantKnots: 0},
		{name: "dense-subsystem", graph: denseSubsystemCalibrationGraph(), wantKnots: 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			findings := detectDependencyKnots(test.graph, ".")
			if len(findings) != test.wantKnots {
				t.Fatalf("%s topology: expected %d dependency-knot findings, got %#v", test.name, test.wantKnots, findings)
			}
		})
	}
}
