package scan

import (
	"fmt"
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestBoundaryBypassFindsMinorityDirectRouteAroundEstablishedGateway(t *testing.T) {
	files, semantic, edges := bypassFixture()
	candidates := boundaryBypassCandidates(edges, files, semantic)
	if len(candidates) != 1 {
		t.Fatalf("expected one bypass candidate, got %#v", candidates)
	}
	candidate := candidates[0]
	if candidate.source != "ui/bypass.go" || candidate.gateway != "service/gateway.go" || candidate.target != "storage/store.go" {
		t.Fatalf("unexpected bypass candidate: %#v", candidate)
	}
	if candidate.gatewayPeers != 3 || candidate.bypassers != 1 {
		t.Fatalf("unexpected seam strength: %#v", candidate)
	}
}

func TestBoundaryBypassKeepsCommonDirectAccessQuiet(t *testing.T) {
	files, semantic, edges := bypassFixture()
	addBypassEdge(&edges, "ui/peer1.go", "storage/store.go", "calls")
	if candidates := boundaryBypassCandidates(edges, files, semantic); len(candidates) != 0 {
		t.Fatalf("direct access used by a broad share of peers is not a minority bypass: %#v", candidates)
	}
}

func TestBoundaryBypassKeepsSharedFoundationQuiet(t *testing.T) {
	files, semantic, edges := bypassFixture()
	files["worker/job.go"] = struct{}{}
	semantic["worker/job.go"] = dependencyRegion{id: "namespace:worker", label: "namespace worker"}
	addBypassEdge(&edges, "worker/job.go", "storage/store.go", "reads")
	if candidates := boundaryBypassCandidates(edges, files, semantic); len(candidates) != 0 {
		t.Fatalf("widely shared downstream foundations should not be treated as gateway-owned internals: %#v", candidates)
	}
}

func TestBoundaryBypassRequiresGatewayConcentrationTowardTargetRegion(t *testing.T) {
	files, semantic, edges := bypassFixture()
	for index := 1; index <= 4; index++ {
		file := fmt.Sprintf("misc/side%d.go", index)
		files[file] = struct{}{}
		semantic[file] = dependencyRegion{id: "namespace:misc", label: "namespace misc"}
		addBypassEdge(&edges, "service/gateway.go", file, "calls")
	}
	if candidates := boundaryBypassCandidates(edges, files, semantic); len(candidates) != 0 {
		t.Fatalf("a broadly mixed dependency should not be inferred as a downstream gateway: %#v", candidates)
	}
}

func TestBoundaryBypassRequiresEstablishedGatewayPeerUse(t *testing.T) {
	files, semantic, edges := bypassFixture()
	delete(edges.outgoing, "ui/peer2.go")
	delete(edges.outgoing, "ui/peer3.go")
	delete(edges.incoming["service/gateway.go"], "ui/peer2.go")
	delete(edges.incoming["service/gateway.go"], "ui/peer3.go")
	if candidates := boundaryBypassCandidates(edges, files, semantic); len(candidates) != 0 {
		t.Fatalf("one upstream user does not establish an architectural intermediary: %#v", candidates)
	}
}

func TestBoundaryBypassIgnoresImportOnlyShortcut(t *testing.T) {
	graph := dependencyTestGraph(6)
	setDependencyTestPath(&graph, 1, "ui/peer1.go")
	setDependencyTestPath(&graph, 2, "ui/peer2.go")
	setDependencyTestPath(&graph, 3, "ui/peer3.go")
	setDependencyTestPath(&graph, 4, "service/gateway.go")
	setDependencyTestPath(&graph, 5, "storage/store.go")
	setDependencyTestPath(&graph, 6, "storage/other.go")
	addDependency(&graph, 1, 4, "calls")
	addDependency(&graph, 2, 4, "calls")
	addDependency(&graph, 3, 4, "calls")
	addDependency(&graph, 4, 5, "calls")
	addDependency(&graph, 4, 6, "calls")
	addDependency(&graph, 1, 5, "imports")
	if findings := detectBoundaryBypass(graph, "."); len(findings) != 0 {
		t.Fatalf("import-only reach must not establish a behavioral bypass: %#v", findings)
	}
}

func TestBoundaryBypassFindingIsActionable(t *testing.T) {
	files, semantic, edges := bypassFixture()
	candidate := boundaryBypassCandidates(edges, files, semantic)[0]
	finding := candidate.finding(".")
	if finding.Detector != DetectorBoundaryBypass || finding.Scope.Path != "ui/bypass.go" || len(finding.Evidence) != 4 {
		t.Fatalf("unexpected finding: %+v", finding)
	}
	if finding.RequiredOutcome == "" || finding.RecommendedAction == "" {
		t.Fatalf("finding must define a repair target: %+v", finding)
	}
}

func bypassFixture() (map[string]struct{}, map[string]dependencyRegion, bypassGraph) {
	paths := []string{"ui/peer1.go", "ui/peer2.go", "ui/peer3.go", "ui/bypass.go", "service/gateway.go", "storage/store.go", "storage/other.go"}
	files := make(map[string]struct{}, len(paths))
	semantic := make(map[string]dependencyRegion, len(paths))
	for _, file := range paths {
		files[file] = struct{}{}
	}
	for _, file := range paths[:4] {
		semantic[file] = dependencyRegion{id: "namespace:ui", label: "namespace ui"}
	}
	semantic["service/gateway.go"] = dependencyRegion{id: "namespace:service", label: "namespace service"}
	semantic["storage/store.go"] = dependencyRegion{id: "namespace:storage", label: "namespace storage"}
	semantic["storage/other.go"] = dependencyRegion{id: "namespace:storage", label: "namespace storage"}
	edges := bypassGraph{outgoing: map[string]map[string]map[string]struct{}{}, incoming: map[string]map[string]struct{}{}}
	for _, peer := range paths[:3] {
		addBypassEdge(&edges, peer, "service/gateway.go", "calls")
	}
	addBypassEdge(&edges, "service/gateway.go", "storage/store.go", "calls")
	addBypassEdge(&edges, "service/gateway.go", "storage/other.go", "calls")
	addBypassEdge(&edges, "ui/bypass.go", "storage/store.go", "calls")
	return files, semantic, edges
}

func addBypassEdge(graph *bypassGraph, source, target, relation string) {
	if graph.outgoing[source] == nil {
		graph.outgoing[source] = map[string]map[string]struct{}{}
	}
	if graph.outgoing[source][target] == nil {
		graph.outgoing[source][target] = map[string]struct{}{}
	}
	graph.outgoing[source][target][relation] = struct{}{}
	if graph.incoming[target] == nil {
		graph.incoming[target] = map[string]struct{}{}
	}
	graph.incoming[target][source] = struct{}{}
}

func setDependencyTestPath(graph *arcana.Graph, fileID uint32, file string) {
	for index := range graph.Sources {
		if graph.Sources[index].NodeID == fileID || graph.Sources[index].NodeID == dependencySymbolID(fileID) {
			graph.Sources[index].Path = file
		}
	}
}
