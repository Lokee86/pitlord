package scan

import (
	"testing"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func TestCrossFileIntermediaryBypassFindsEstablishedNamedPeerRoute(t *testing.T) {
	graph := crossFileEstablishedPeerFixture()
	findings := detectCrossFileIntermediaryBypass(graph, ".")
	if len(findings) != 1 {
		t.Fatalf("expected one cross-file bypass, got %#v", findings)
	}
	finding := findings[0]
	if finding.Detector != DetectorCrossFileIntermediaryBypass || finding.Scope.Name != "MatchDecision" {
		t.Fatalf("unexpected finding: %+v", finding)
	}
	if len(finding.Evidence) != 3 || finding.Evidence[2].Kind != "peer-route" {
		t.Fatalf("expected peer-route evidence: %+v", finding)
	}
}

func TestCrossFileIntermediaryBypassKeepsAbandonedParallelWrappersQuiet(t *testing.T) {
	graph := crossFileNameAlignedFixture()
	if findings := detectCrossFileIntermediaryBypass(graph, "."); len(findings) != 0 {
		t.Fatalf("parallel abandoned wrapper has no defensible bypass direction: %#v", findings)
	}
}

func TestCrossFileIntermediaryBypassRejectsSameNamedOverloadTarget(t *testing.T) {
	graph := crossFileEstablishedPeerFixture()
	setSymbolName(&graph, 12, "matchDecisionLocked")
	if findings := detectCrossFileIntermediaryBypass(graph, "."); len(findings) != 0 {
		t.Fatalf("same-named overload chain should remain quiet: %#v", findings)
	}
}

func TestCrossFileIntermediaryBypassRequiresNarrowTargetOwnership(t *testing.T) {
	graph := crossFileEstablishedPeerFixture()
	otherFile := "services/game-server/internal/game/other.go"
	addFileNode(&graph, 20, otherFile)
	other := addSymbolNode(&graph, 21, "OtherDecision", otherFile)
	addSymbolCall(&graph, other, symbolNode(12, "evaluateMatchDecisionLocked", crossFileTargetPath()))
	if findings := detectCrossFileIntermediaryBypass(graph, "."); len(findings) != 0 {
		t.Fatalf("broad shared target should remain quiet: %#v", findings)
	}
}

func TestCrossFileIntermediaryBypassRequiresShortIntermediary(t *testing.T) {
	graph := crossFileEstablishedPeerFixture()
	setSymbolSpan(&graph, 11, 10, 30)
	if findings := detectCrossFileIntermediaryBypass(graph, "."); len(findings) != 0 {
		t.Fatalf("long behavioral method should not be inferred as a forwarding seam: %#v", findings)
	}
}

func TestCrossFileIntermediaryBypassRequiresSameRegion(t *testing.T) {
	graph := crossFileEstablishedPeerFixture()
	moveSymbolToFile(&graph, 10, "services/game-server/internal/other/control_match.go")
	addFileNode(&graph, 99, "services/game-server/internal/other/control_match.go")
	if findings := detectCrossFileIntermediaryBypass(graph, "."); len(findings) != 0 {
		t.Fatalf("cross-region direct call should remain quiet: %#v", findings)
	}
}

func TestCrossFileIntermediaryBypassLeavesSameFileCaseToLocalDetector(t *testing.T) {
	graph := crossFileEstablishedPeerFixture()
	moveSymbolToFile(&graph, 10, crossFileIntermediaryPath())
	if findings := detectCrossFileIntermediaryBypass(graph, "."); len(findings) != 0 {
		t.Fatalf("same-file bypass belongs to symbol-intermediary-bypass: %#v", findings)
	}
}

func TestGeneralizedScanIncludesCrossFileIntermediaryBypass(t *testing.T) {
	result := analyzeGraph(".", crossFileEstablishedPeerFixture())
	for _, finding := range result.Findings {
		if finding.Detector == DetectorCrossFileIntermediaryBypass {
			return
		}
	}
	t.Fatalf("generalized scan did not include %s: %#v", DetectorCrossFileIntermediaryBypass, result.Findings)
}

func crossFileNameAlignedFixture() arcana.Graph {
	graph := arcana.Graph{Outgoing: map[uint32][]arcana.Relationship{}}
	callerFile := "services/game-server/internal/devtools/streamruntime/simulation.go"
	intermediaryFile := "services/game-server/internal/devtools/streamruntime/runtime.go"
	targetFile := "services/game-server/internal/devtools/streamruntime/continuous_bullet_streams.go"
	addFileNode(&graph, 1, callerFile)
	addFileNode(&graph, 2, intermediaryFile)
	addFileNode(&graph, 3, targetFile)
	caller := addSymbolNode(&graph, 10, "StepContinuousBulletStreams", callerFile)
	intermediary := addSymbolNode(&graph, 11, "StepContinuousBulletStreams", intermediaryFile)
	target := addSymbolNode(&graph, 12, "Step", targetFile)
	setSymbolSpan(&graph, 11, 20, 22)
	addSymbolCall(&graph, intermediary, target)
	addSymbolCall(&graph, caller, target)
	return graph
}

func crossFileEstablishedPeerFixture() arcana.Graph {
	graph := arcana.Graph{Outgoing: map[uint32][]arcana.Relationship{}}
	callerFile := "services/game-server/internal/game/control_match.go"
	intermediaryFile := crossFileIntermediaryPath()
	targetFile := crossFileTargetPath()
	addFileNode(&graph, 1, callerFile)
	addFileNode(&graph, 2, intermediaryFile)
	addFileNode(&graph, 3, targetFile)
	caller := addSymbolNode(&graph, 10, "MatchDecision", callerFile)
	intermediary := addSymbolNode(&graph, 11, "matchDecisionLocked", intermediaryFile)
	target := addSymbolNode(&graph, 12, "evaluateMatchDecisionLocked", targetFile)
	peer := addSymbolNode(&graph, 13, "MatchDecision", intermediaryFile)
	otherPeer := addSymbolNode(&graph, 14, "IsGameOver", intermediaryFile)
	setSymbolSpan(&graph, 11, 42, 44)
	addSymbolCall(&graph, intermediary, target)
	addSymbolCall(&graph, peer, intermediary)
	addSymbolCall(&graph, otherPeer, intermediary)
	addSymbolCall(&graph, caller, target)
	return graph
}

func crossFileIntermediaryPath() string {
	return "services/game-server/internal/game/match.go"
}

func crossFileTargetPath() string {
	return "services/game-server/internal/game/match_mode_evaluation.go"
}

func addFileNode(graph *arcana.Graph, id uint32, file string) {
	graph.Sources = append(graph.Sources, arcana.Node{NodeID: id, Kind: "file", Path: file, Name: file})
}

func setSymbolName(graph *arcana.Graph, nodeID uint32, name string) {
	for index := range graph.Sources {
		if graph.Sources[index].NodeID == nodeID {
			graph.Sources[index].Name = name
			graph.Sources[index].Identity = graph.Sources[index].Path + "::" + name
		}
	}
	for sourceID, relationships := range graph.Outgoing {
		for index := range relationships {
			if relationships[index].Node.NodeID == nodeID {
				relationships[index].Node.Name = name
				relationships[index].Node.Identity = relationships[index].Node.Path + "::" + name
			}
		}
		graph.Outgoing[sourceID] = relationships
	}
}

func setSymbolSpan(graph *arcana.Graph, nodeID uint32, startLine, endLine int) {
	for index := range graph.Sources {
		if graph.Sources[index].NodeID == nodeID {
			graph.Sources[index].Span = &arcana.Span{Path: graph.Sources[index].Path, StartLine: startLine, EndLine: endLine}
		}
	}
}

func moveSymbolToFile(graph *arcana.Graph, nodeID uint32, file string) {
	for index := range graph.Sources {
		if graph.Sources[index].NodeID == nodeID {
			graph.Sources[index].Path = file
			graph.Sources[index].Identity = file + "::" + graph.Sources[index].Name
		}
	}
	for sourceID, relationships := range graph.Outgoing {
		for index := range relationships {
			if relationships[index].Node.NodeID == nodeID {
				relationships[index].Node.Path = file
				relationships[index].Node.Identity = file + "::" + relationships[index].Node.Name
			}
		}
		graph.Outgoing[sourceID] = relationships
	}
}
