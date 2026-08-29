package scan

import (
	"fmt"
	"strings"
	"testing"
)

func TestDependencyDepthFindsCrossRegionSerialChain(t *testing.T) {
	paths := []string{
		"app/entry.go", "api/handler.go", "service/service.go",
		"domain/workflow.go", "storage/repository.go",
	}
	graph := dependencyDepthTestGraph(paths)
	for index := 0; index < len(paths)-1; index++ {
		relation := []string{"calls", "depends-on", "extends"}[index%3]
		addDepthDependency(&graph, paths, index, index+1, relation)
	}

	findings := detectDependencyDepth(graph, ".")
	if len(findings) != 1 {
		t.Fatalf("expected one depth finding, got %#v", findings)
	}
	finding := findings[0]
	if finding.Detector != DetectorDependencyDepth || finding.Scope.Path != paths[0] {
		t.Fatalf("unexpected finding: %+v", finding)
	}
	if finding.Disposition != DispositionAdvisory || finding.Severity != SeverityWarning {
		t.Fatalf("unexpected judgment: %+v", finding)
	}
	if !strings.Contains(finding.Summary, "4 dependency hops") || len(finding.Evidence) != 3 {
		t.Fatalf("finding lacks chain evidence: %+v", finding)
	}
}

func TestDependencyDepthAllowsShallowSideDependencies(t *testing.T) {
	paths := []string{
		"app/entry.go", "api/handler.go", "service/service.go", "domain/workflow.go",
		"storage/repository.go", "transport/client.go", "foundation/core.go",
		"utility/a.go", "utility/b.go", "utility/c.go",
	}
	graph := dependencyDepthTestGraph(paths)
	for index := 0; index < 6; index++ {
		addDepthDependency(&graph, paths, index, index+1, "depends-on")
		if index < 3 {
			addDepthDependency(&graph, paths, index, 7+index, "calls")
		}
	}
	findings := detectDependencyDepth(graph, ".")
	if len(findings) != 1 {
		t.Fatalf("shallow side dependencies should preserve a dominant deep chain: %#v", findings)
	}
	if !strings.Contains(findings[0].Evidence[2].Message, "depth-dominant") {
		t.Fatalf("expected dominant-branch evidence, got %+v", findings[0])
	}
}

func TestDependencyDepthPromotesVeryDeepCrossRegionChain(t *testing.T) {
	paths := make([]string, 10)
	for index := range paths {
		paths[index] = fmt.Sprintf("layer%d/file.go", index)
	}
	graph := dependencyDepthTestGraph(paths)
	for index := 0; index < len(paths)-1; index++ {
		addDepthDependency(&graph, paths, index, index+1, "depends-on")
	}
	findings := detectDependencyDepth(graph, ".")
	if len(findings) != 1 || findings[0].Severity != SeverityHigh {
		t.Fatalf("expected one high-severity deep chain, got %#v", findings)
	}
}

func TestDependencyDepthKeepsLongSingleRegionChainQuiet(t *testing.T) {
	paths := make([]string, 10)
	for index := range paths {
		paths[index] = fmt.Sprintf("service/file%d.go", index)
	}
	graph := dependencyDepthTestGraph(paths)
	for index := 0; index < len(paths)-1; index++ {
		addDepthDependency(&graph, paths, index, index+1, "depends-on")
	}
	if findings := detectDependencyDepth(graph, "."); len(findings) != 0 {
		t.Fatalf("single-region implementation depth should remain quiet: %#v", findings)
	}
}

func TestDependencyDepthKeepsWideLayeringQuiet(t *testing.T) {
	paths := make([]string, 0, 21)
	for layer := 0; layer < 7; layer++ {
		for file := 0; file < 3; file++ {
			paths = append(paths, fmt.Sprintf("layer%d/file%d.go", layer, file))
		}
	}
	graph := dependencyDepthTestGraph(paths)
	for layer := 0; layer < 6; layer++ {
		for file := 0; file < 3; file++ {
			source := layer*3 + file
			addDepthDependency(&graph, paths, source, (layer+1)*3+file, "depends-on")
			addDepthDependency(&graph, paths, source, (layer+1)*3+(file+1)%3, "depends-on")
		}
	}
	if findings := detectDependencyDepth(graph, "."); len(findings) != 0 {
		t.Fatalf("wide intentional layering should not look like a fragile serial chain: %#v", findings)
	}
}

func TestDependencyDepthKeepsShortCrossRegionChainQuiet(t *testing.T) {
	paths := []string{"a/a.go", "b/b.go", "c/c.go", "d/d.go"}
	graph := dependencyDepthTestGraph(paths)
	for index := 0; index < len(paths)-1; index++ {
		addDepthDependency(&graph, paths, index, index+1, "depends-on")
	}
	if findings := detectDependencyDepth(graph, "."); len(findings) != 0 {
		t.Fatalf("three-hop chain should remain below threshold: %#v", findings)
	}
}

func TestDependencyDepthRestartsAtSharedConvergenceBoundary(t *testing.T) {
	paths := []string{
		"caller-a/a.go", "caller-b/b.go", "shared/root.go", "layer1/a.go", "layer2/b.go",
		"layer3/c.go", "layer4/d.go", "layer5/e.go", "layer6/f.go",
	}
	graph := dependencyDepthTestGraph(paths)
	addDepthDependency(&graph, paths, 0, 2, "depends-on")
	addDepthDependency(&graph, paths, 1, 2, "depends-on")
	for index := 2; index < len(paths)-1; index++ {
		addDepthDependency(&graph, paths, index, index+1, "depends-on")
	}
	findings := detectDependencyDepth(graph, ".")
	if len(findings) != 1 || findings[0].Scope.Path != "shared/root.go" {
		t.Fatalf("shared convergence should own downstream depth instead of callers: %#v", findings)
	}
}

func TestDependencyDepthStopsAtHighlyReusedCentralFoundation(t *testing.T) {
	paths := []string{
		"a/a.go", "b/b.go", "c/c.go", "foundation/core.go",
		"foundation/internal.go", "foundation/storage.go",
	}
	for index := 0; index < 24; index++ {
		paths = append(paths, fmt.Sprintf("caller%d/use.go", index))
	}
	graph := dependencyDepthTestGraph(paths)
	for index := 0; index < 5; index++ {
		addDepthDependency(&graph, paths, index, index+1, "depends-on")
	}
	for index := 6; index < len(paths); index++ {
		addDepthDependency(&graph, paths, index, 3, "depends-on")
	}
	if findings := detectDependencyDepth(graph, "."); len(findings) != 0 {
		t.Fatalf("depth behind a highly reused central foundation must not be inherited by callers: %#v", findings)
	}
}

func TestDependencyDepthDefersCyclesToDependencyKnots(t *testing.T) {
	paths := []string{"a/a.go", "b/b.go", "c/c.go", "d/d.go", "e/e.go", "f/f.go", "g/g.go", "h/h.go"}
	graph := dependencyDepthTestGraph(paths)
	for index := 0; index < len(paths)-1; index++ {
		addDepthDependency(&graph, paths, index, index+1, "depends-on")
	}
	addDepthDependency(&graph, paths, 3, 2, "depends-on")
	if findings := detectDependencyDepth(graph, "."); len(findings) != 0 {
		t.Fatalf("cyclic members must be removed from depth analysis: %#v", findings)
	}
}

func TestDependencyDepthIgnoresReferenceOnlyChains(t *testing.T) {
	paths := []string{"a/a.go", "b/b.go", "c/c.go", "d/d.go", "e/e.go", "f/f.go", "g/g.go"}
	graph := dependencyDepthTestGraph(paths)
	for index := 0; index < len(paths)-1; index++ {
		addDepthDependency(&graph, paths, index, index+1, "references")
	}
	if findings := detectDependencyDepth(graph, "."); len(findings) != 0 {
		t.Fatalf("reference-only chains must not establish architectural depth: %#v", findings)
	}
}

func TestDependencyDepthIgnoresImportOnlyChains(t *testing.T) {
	paths := []string{"a/a.go", "b/b.go", "c/c.go", "d/d.go", "e/e.go", "f/f.go", "g/g.go"}
	graph := dependencyDepthTestGraph(paths)
	for index := 0; index < len(paths)-1; index++ {
		addDepthDependency(&graph, paths, index, index+1, "imports")
	}
	if findings := detectDependencyDepth(graph, "."); len(findings) != 0 {
		t.Fatalf("import-only chains must not establish architectural depth: %#v", findings)
	}
}
