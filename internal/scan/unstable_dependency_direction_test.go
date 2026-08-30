package scan

import (
	"strings"
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestUnstableDependencyDirectionFindsStableToUnstableRegionDependency(t *testing.T) {
	graph := stableToUnstableDirectionGraph(true)
	findings := detectUnstableDependencyDirection(graph, ".")
	if len(findings) != 1 {
		t.Fatalf("expected one unstable-direction finding, got %#v", findings)
	}
	finding := findings[0]
	if finding.Detector != DetectorUnstableDependencyDirection || finding.Scope.Kind != "architecture-dependency" {
		t.Fatalf("unexpected finding: %+v", finding)
	}
	if finding.Scope.Name != "namespace Region01 -> namespace Region02" || finding.Disposition != DispositionAdvisory {
		t.Fatalf("unexpected direction judgment: %+v", finding)
	}
	if !strings.Contains(finding.Evidence[1].Message, "2 incoming regions and 1 outgoing") ||
		!strings.Contains(finding.Evidence[2].Message, "1 incoming regions and 2 outgoing") {
		t.Fatalf("finding does not expose stability evidence: %+v", finding)
	}
}

func TestUnstableDependencyDirectionKeepsUnstableToStableDependencyQuiet(t *testing.T) {
	graph := dependencyTestGraph(36)
	addCalibrationNamespaceRegions(&graph, 6)
	addDependency(&graph, 1, 7, "depends-on")
	addDependency(&graph, 2, 8, "imports")
	addDependency(&graph, 3, 13, "depends-on")
	addDependency(&graph, 19, 1, "depends-on")
	addDependency(&graph, 25, 7, "depends-on")
	addDependency(&graph, 7, 31, "depends-on")

	if findings := detectUnstableDependencyDirection(graph, "."); len(findings) != 0 {
		t.Fatalf("dependency from outgoing-heavy to incoming-heavy region should remain quiet: %#v", findings)
	}
}

func TestUnstableDependencyDirectionRequiresArchitectureLevelBoundarySupport(t *testing.T) {
	graph := stableToUnstableDirectionGraph(false)
	if findings := detectUnstableDependencyDirection(graph, "."); len(findings) != 0 {
		t.Fatalf("one source file should not establish an architecture-level direction violation: %#v", findings)
	}
}

func TestUnstableDependencyDirectionDefersBidirectionalRegionPairsToCycleAnalysis(t *testing.T) {
	graph := stableToUnstableDirectionGraph(true)
	addDependency(&graph, 7, 1, "depends-on")
	if findings := detectUnstableDependencyDirection(graph, "."); len(findings) != 0 {
		t.Fatalf("bidirectional region pair should be deferred to cycle/coupling analysis: %#v", findings)
	}
}

func TestUnstableDependencyDirectionDefersLongerRegionCyclesToKnotAnalysis(t *testing.T) {
	graph := stableToUnstableDirectionGraph(true)
	addDependency(&graph, 25, 1, "depends-on")
	if findings := detectUnstableDependencyDirection(graph, "."); len(findings) != 0 {
		t.Fatalf("dependency participating in a multi-region cycle should be deferred to knot analysis: %#v", findings)
	}
}

func TestUnstableDependencyDirectionKeepsHealthyLayeredTopologyQuiet(t *testing.T) {
	graph := layeredCalibrationGraph()
	addCalibrationNamespaceRegions(&graph, 8)
	if findings := detectUnstableDependencyDirection(graph, "."); len(findings) != 0 {
		t.Fatalf("healthy dependency layering should point toward equal or greater stability: %#v", findings)
	}
}

func TestAnalyzeGraphIncludesUnstableDependencyDirection(t *testing.T) {
	result := analyzeGraph(".", stableToUnstableDirectionGraph(true))
	for _, finding := range result.Findings {
		if finding.Detector == DetectorUnstableDependencyDirection {
			return
		}
	}
	t.Fatalf("generalized scan did not include %s: %#v", DetectorUnstableDependencyDirection, result.Findings)
}

func stableToUnstableDirectionGraph(multipleBoundarySources bool) arcana.Graph {
	graph := dependencyTestGraph(36)
	addCalibrationNamespaceRegions(&graph, 6)

	addDependency(&graph, 1, 7, "depends-on")
	if multipleBoundarySources {
		addDependency(&graph, 2, 8, "imports")
	}
	addDependency(&graph, 13, 1, "depends-on")
	addDependency(&graph, 19, 2, "depends-on")
	addDependency(&graph, 7, 25, "depends-on")
	addDependency(&graph, 8, 31, "depends-on")
	return graph
}
