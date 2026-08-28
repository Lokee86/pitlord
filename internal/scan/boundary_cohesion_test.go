package scan

import (
	"fmt"
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestBoundaryCohesionFindsOutwardFacingWeakRegion(t *testing.T) {
	graph := dependencyTestGraph(32)
	addCalibrationNamespaceRegions(&graph, 8)

	for region := uint32(0); region < 4; region++ {
		start := region*8 + 1
		for offset := uint32(0); offset < 8; offset++ {
			source := start + offset
			addDependency(&graph, source, start+(offset+1)%8, "references")
			addDependency(&graph, source, start+(offset+2)%8, "depends-on")
		}
	}
	for source := uint32(1); source <= 8; source++ {
		graph.Outgoing[dependencySymbolID(source)] = nil
		for region := uint32(1); region < 4; region++ {
			addDependency(&graph, source, region*8+source, "depends-on")
		}
	}
	addDependency(&graph, 1, 2, "references")

	findings := detectBoundaryCohesion(graph, ".")
	if len(findings) != 1 {
		t.Fatalf("expected one boundary/cohesion finding, got %#v", findings)
	}
	finding := findings[0]
	if finding.Detector != DetectorBoundaryCohesion || finding.Scope.Name != "namespace Region01" {
		t.Fatalf("unexpected finding: %+v", finding)
	}
	if finding.Scope.Kind != "architecture-region" || finding.Disposition != DispositionAdvisory {
		t.Fatalf("unexpected judgment: %+v", finding)
	}
}

func TestBoundaryCohesionKeepsIncomingSharedHubRegionQuiet(t *testing.T) {
	graph := directionalizeCalibrationReferences(hubHeavyCalibrationGraph())
	addCalibrationNamespaceRegions(&graph, 8)
	if findings := detectBoundaryCohesion(graph, "."); len(findings) != 0 {
		t.Fatalf("incoming-heavy shared hub topology belongs to bottleneck analysis, got %#v", findings)
	}
}

func TestBoundaryCohesionCalibrationAcrossArcanaTopologyFamilies(t *testing.T) {
	tests := []struct {
		name  string
		graph arcana.Graph
	}{
		{name: "modular", graph: modularCalibrationGraph()},
		{name: "entangled", graph: entangledCalibrationGraph()},
		{name: "hub-heavy", graph: hubHeavyCalibrationGraph()},
		{name: "layered", graph: layeredCalibrationGraph()},
		{name: "dense-subsystem", graph: denseSubsystemCalibrationGraph()},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			graph := directionalizeCalibrationReferences(test.graph)
			addCalibrationNamespaceRegions(&graph, 8)
			if findings := detectBoundaryCohesion(graph, "."); len(findings) != 0 {
				t.Fatalf("family-wide topology should not create a peer region outlier: %#v", findings)
			}
		})
	}
}

func addCalibrationNamespaceRegions(graph *arcana.Graph, groupSize uint32) {
	fileCount := uint32(0)
	for _, node := range graph.Sources {
		if node.Kind == "file" && node.NodeID > fileCount {
			fileCount = node.NodeID
		}
	}
	region := uint32(0)
	for start := uint32(1); start <= fileCount; start += groupSize {
		region++
		containerID := uint32(50_000) + region
		container := arcana.Node{NodeID: containerID, Key: filePath(start), Kind: "namespace", Name: calibrationRegionName(region)}
		graph.Sources = append(graph.Sources, container)
		end := min(start+groupSize-1, fileCount)
		for fileID := start; fileID <= end; fileID++ {
			graph.Outgoing[containerID] = append(graph.Outgoing[containerID], arcana.Relationship{
				Relation: "contains",
				Node:     arcana.Node{NodeID: uint32(60_000) + fileID, Kind: "type", Path: filePath(fileID), Span: &arcana.Span{Path: filePath(fileID)}},
			})
		}
	}
}

func calibrationRegionName(region uint32) string {
	return fmt.Sprintf("Region%02d", region)
}
