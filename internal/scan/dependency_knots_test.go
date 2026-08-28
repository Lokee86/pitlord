package scan

import (
	"strings"
	"testing"
)

func TestDependencyKnotsFindsCrossRegionCycle(t *testing.T) {
	units := []dependencyUnit{
		{path: "api/handler.any", incoming: stringSetFrom("storage/store.any"), outgoing: stringSetFrom("service/usecase.any")},
		{path: "service/usecase.any", incoming: stringSetFrom("api/handler.any"), outgoing: stringSetFrom("storage/store.any")},
		{path: "storage/store.any", incoming: stringSetFrom("service/usecase.any"), outgoing: stringSetFrom("api/handler.any")},
		{path: "model/value.any", incoming: map[string]struct{}{}, outgoing: map[string]struct{}{}},
	}
	findings := detectDependencyKnotsFromUnits(units, ".")
	if len(findings) != 1 {
		t.Fatalf("expected one dependency knot, got %#v", findings)
	}
	finding := findings[0]
	if finding.Detector != DetectorDependencyKnots || finding.Scope.Path != "api/handler.any" {
		t.Fatalf("unexpected finding: %#v", finding)
	}
	if finding.Severity != SeverityHigh || len(finding.Evidence) != 3 {
		t.Fatalf("unexpected knot judgment: %#v", finding)
	}
	if !strings.Contains(finding.RequiredOutcome, "no longer mutually reachable") {
		t.Fatalf("finding is not mechanically actionable: %#v", finding)
	}
}

func TestDependencyKnotsIgnoresGenericReferenceCycles(t *testing.T) {
	graph := dependencyTestGraph(2)
	addDependency(&graph, 1, 2, "references")
	addDependency(&graph, 2, 1, "references")
	if findings := detectDependencyKnots(graph, "."); len(findings) != 0 {
		t.Fatalf("generic references must not define architecture dependency knots: %#v", findings)
	}

	graph = dependencyTestGraph(2)
	addDependency(&graph, 1, 2, "calls")
	addDependency(&graph, 2, 1, "calls")
	findings := detectDependencyKnots(graph, ".")
	if len(findings) != 1 || findings[0].Severity != SeverityWarning {
		t.Fatalf("directional two-file cycle should remain a warning: %#v", findings)
	}
}

func TestDependencyKnotsIgnoresLocalCycleWhenRepositoryHasBoundaries(t *testing.T) {
	units := []dependencyUnit{
		{path: "parser/a.any", incoming: stringSetFrom("parser/b.any"), outgoing: stringSetFrom("parser/b.any")},
		{path: "parser/b.any", incoming: stringSetFrom("parser/a.any"), outgoing: stringSetFrom("parser/a.any")},
		{path: "other/c.any", incoming: map[string]struct{}{}, outgoing: map[string]struct{}{}},
	}
	if findings := detectDependencyKnotsFromUnits(units, "."); len(findings) != 0 {
		t.Fatalf("same-region implementation cycle should not be promoted to an architectural knot: %#v", findings)
	}
}

func TestDependencyKnotsSingleRegionFallbackRequiresDensity(t *testing.T) {
	dense := []dependencyUnit{
		{path: "src/a.any", incoming: stringSetFrom("src/b.any", "src/c.any"), outgoing: stringSetFrom("src/b.any", "src/c.any")},
		{path: "src/b.any", incoming: stringSetFrom("src/a.any", "src/c.any"), outgoing: stringSetFrom("src/a.any", "src/c.any")},
		{path: "src/c.any", incoming: stringSetFrom("src/a.any", "src/b.any"), outgoing: stringSetFrom("src/a.any", "src/b.any")},
	}
	if findings := detectDependencyKnotsFromUnits(dense, "src"); len(findings) != 1 {
		t.Fatalf("dense single-region knot should remain visible in a narrowed scope: %#v", findings)
	}

	sparse := make([]dependencyUnit, 20)
	for index := range sparse {
		current := "src/file" + twoDigits(index) + ".any"
		next := "src/file" + twoDigits((index+1)%len(sparse)) + ".any"
		previous := "src/file" + twoDigits((index-1+len(sparse))%len(sparse)) + ".any"
		sparse[index] = dependencyUnit{path: current, incoming: stringSetFrom(previous), outgoing: stringSetFrom(next)}
	}
	if findings := detectDependencyKnotsFromUnits(sparse, "src"); len(findings) != 0 {
		t.Fatalf("sparse single-region ring should stay quiet: %#v", findings)
	}
}

func TestStronglyConnectedDependencyComponentsDeterministic(t *testing.T) {
	units := []dependencyUnit{
		{path: "c", incoming: stringSetFrom("b"), outgoing: stringSetFrom("a")},
		{path: "a", incoming: stringSetFrom("c"), outgoing: stringSetFrom("b")},
		{path: "b", incoming: stringSetFrom("a"), outgoing: stringSetFrom("c")},
		{path: "d", incoming: map[string]struct{}{}, outgoing: map[string]struct{}{}},
	}
	components := stronglyConnectedDependencyComponents(units)
	if len(components) != 2 || strings.Join(components[0], ",") != "a,b,c" || strings.Join(components[1], ",") != "d" {
		t.Fatalf("unexpected components: %#v", components)
	}
}

func stringSetFrom(values ...string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		result[value] = struct{}{}
	}
	return result
}

func twoDigits(value int) string {
	return string([]byte{'0' + byte(value/10), '0' + byte(value%10)})
}
