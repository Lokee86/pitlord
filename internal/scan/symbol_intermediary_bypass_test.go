package scan

import (
	"fmt"
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestSymbolIntermediaryBypassFindsAbandonedWrapperInsideEstablishedFamily(t *testing.T) {
	graph := symbolBypassFixture(true, 3)
	findings := detectSymbolIntermediaryBypass(graph, ".")
	if len(findings) != 1 {
		t.Fatalf("expected one symbol intermediary bypass, got %#v", findings)
	}
	finding := findings[0]
	if finding.Detector != DetectorSymbolIntermediaryBypass || finding.Scope.Path != "services/game-server/internal/devtools/player_counters.go" || finding.Scope.Name != "handleDebugAddScore" {
		t.Fatalf("unexpected finding: %+v", finding)
	}
	if len(finding.Evidence) != 3 || finding.RequiredOutcome == "" || finding.RecommendedAction == "" {
		t.Fatalf("finding must be mechanically actionable: %+v", finding)
	}
}

func TestSymbolIntermediaryBypassKeepsIntactWrapperQuiet(t *testing.T) {
	graph := symbolBypassFixture(false, 3)
	if findings := detectSymbolIntermediaryBypass(graph, "."); len(findings) != 0 {
		t.Fatalf("intact wrapper route should remain quiet: %#v", findings)
	}
}

func TestSymbolIntermediaryBypassRequiresEstablishedPeerPattern(t *testing.T) {
	graph := symbolBypassFixture(true, 2)
	if findings := detectSymbolIntermediaryBypass(graph, "."); len(findings) != 0 {
		t.Fatalf("two peer wrappers do not establish an intermediary convention: %#v", findings)
	}
}

func TestSymbolIntermediaryBypassKeepsStillUsedIntermediaryQuiet(t *testing.T) {
	graph := symbolBypassFixture(true, 3)
	caller := addSymbolNode(&graph, 90, "otherCaller", sourceFixturePath())
	addSymbolCall(&graph, caller, symbolNode(12, "addDebugScoreForPlayer", sourceFixturePath()))
	if findings := detectSymbolIntermediaryBypass(graph, "."); len(findings) != 0 {
		t.Fatalf("an intermediary still used elsewhere is not abandoned: %#v", findings)
	}
}

func TestSymbolIntermediaryBypassRequiresSharedPeerPipeline(t *testing.T) {
	graph := symbolBypassFixture(true, 3)
	for _, callerID := range []uint32{40, 43, 46} {
		removeSymbolCall(&graph, callerID, 30)
	}
	if findings := detectSymbolIntermediaryBypass(graph, "."); len(findings) != 0 {
		t.Fatalf("unrelated wrappers without a shared support step must remain quiet: %#v", findings)
	}
}

func TestSymbolIntermediaryBypassKeepsIntermediaryOverloadUseQuiet(t *testing.T) {
	graph := symbolBypassFixture(true, 3)
	sibling := arcana.Node{NodeID: 91, Identity: sourceFixturePath() + "::addDebugScoreForPlayer#overload", Kind: "function", Path: sourceFixturePath(), Name: "addDebugScoreForPlayer"}
	graph.Sources = append(graph.Sources, sibling)
	addSymbolCall(&graph, sibling, symbolNode(20, "AddPlayerScore", targetFixturePath()))
	addSymbolCall(&graph, symbolNode(10, "handleDebugAddScore", sourceFixturePath()), sibling)
	if findings := detectSymbolIntermediaryBypass(graph, "."); len(findings) != 0 {
		t.Fatalf("calling another overload in the intermediary family is not a bypass: %#v", findings)
	}
}

func TestSymbolIntermediaryBypassDoesNotTreatConstructionAsBehavioralBypass(t *testing.T) {
	graph := symbolBypassFixture(true, 3)
	setSymbolKind(&graph, 20, "constructor")
	if findings := detectSymbolIntermediaryBypass(graph, "."); len(findings) != 0 {
		t.Fatalf("direct construction belongs to creation-boundary analysis, not intermediary bypass: %#v", findings)
	}
}

func TestSymbolIntermediaryBypassIncludesDevtoolsButExcludesTests(t *testing.T) {
	if excludedSymbolBypassPath(sourceFixturePath()) {
		t.Fatal("runtime devtools code must remain eligible for symbol-level bypass analysis")
	}
	if !excludedSymbolBypassPath("services/game-server/internal/devtools/player_counters_test.go") {
		t.Fatal("test source must remain excluded")
	}
	if !excludedSymbolBypassPath("client/addons/gut/test.gd") {
		t.Fatal("vendored GUT test framework must remain excluded")
	}
}

func TestGeneralizedScanIncludesSymbolIntermediaryBypass(t *testing.T) {
	result := analyzeGraph(".", symbolBypassFixture(true, 3))
	for _, finding := range result.Findings {
		if finding.Detector == DetectorSymbolIntermediaryBypass {
			return
		}
	}
	t.Fatalf("generalized scan did not include %s: %#v", DetectorSymbolIntermediaryBypass, result.Findings)
}

func symbolBypassFixture(bypass bool, peers int) arcana.Graph {
	graph := arcana.Graph{Outgoing: map[uint32][]arcana.Relationship{}}
	graph.Sources = append(graph.Sources,
		arcana.Node{NodeID: 1, Kind: "file", Path: sourceFixturePath(), Name: "player_counters.go"},
		arcana.Node{NodeID: 2, Kind: "file", Path: targetFixturePath(), Name: "control_counters.go"},
	)

	handler := addSymbolNode(&graph, 10, "handleDebugAddScore", sourceFixturePath())
	intermediary := addSymbolNode(&graph, 12, "addDebugScoreForPlayer", sourceFixturePath())
	target := addSymbolNode(&graph, 20, "AddPlayerScore", targetFixturePath())
	resolve := addSymbolNode(&graph, 30, "resolveCommandTargetPlayerIDs", sourceFixturePath())
	addSymbolCall(&graph, intermediary, target)
	addSymbolCall(&graph, handler, resolve)
	if bypass {
		addSymbolCall(&graph, handler, target)
	} else {
		addSymbolCall(&graph, handler, intermediary)
	}

	for index := 0; index < peers; index++ {
		h := addSymbolNode(&graph, uint32(40+index*3), fmt.Sprintf("handler%d", index), sourceFixturePath())
		w := addSymbolNode(&graph, uint32(41+index*3), fmt.Sprintf("wrapper%d", index), sourceFixturePath())
		t := addSymbolNode(&graph, uint32(42+index*3), fmt.Sprintf("Target%d", index), targetFixturePath())
		addSymbolCall(&graph, h, resolve)
		addSymbolCall(&graph, h, w)
		addSymbolCall(&graph, w, t)
	}
	return graph
}

func addSymbolNode(graph *arcana.Graph, id uint32, name, file string) arcana.Node {
	node := symbolNode(id, name, file)
	graph.Sources = append(graph.Sources, node)
	return node
}

func symbolNode(id uint32, name, file string) arcana.Node {
	return arcana.Node{NodeID: id, Identity: file + "::" + name, Kind: "function", Path: file, Name: name}
}

func addSymbolCall(graph *arcana.Graph, source, target arcana.Node) {
	graph.Outgoing[source.NodeID] = append(graph.Outgoing[source.NodeID], arcana.Relationship{Relation: "calls", Node: target})
}

func removeSymbolCall(graph *arcana.Graph, sourceID, targetID uint32) {
	relationships := graph.Outgoing[sourceID]
	kept := relationships[:0]
	for _, relationship := range relationships {
		if relationship.Relation == "calls" && relationship.Node.NodeID == targetID {
			continue
		}
		kept = append(kept, relationship)
	}
	graph.Outgoing[sourceID] = kept
}

func setSymbolKind(graph *arcana.Graph, nodeID uint32, kind string) {
	for index := range graph.Sources {
		if graph.Sources[index].NodeID == nodeID {
			graph.Sources[index].Kind = kind
		}
	}
	for sourceID, relationships := range graph.Outgoing {
		for index := range relationships {
			if relationships[index].Node.NodeID == nodeID {
				relationships[index].Node.Kind = kind
			}
		}
		graph.Outgoing[sourceID] = relationships
	}
}

func sourceFixturePath() string {
	return "services/game-server/internal/devtools/player_counters.go"
}

func targetFixturePath() string {
	return "services/game-server/internal/game/control_counters.go"
}
