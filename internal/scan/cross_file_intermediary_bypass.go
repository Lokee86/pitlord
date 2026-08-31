package scan

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

const DetectorCrossFileIntermediaryBypass = "cross-file-intermediary-bypass"

const (
	minimumEstablishedIntermediaryCallers = 2
	maximumIntermediarySpanLines          = 5
	maximumTargetCallerFiles              = 2
	maximumCrossFileBypassFindings        = 20
)

type crossFileIntermediaryBypassCandidate struct {
	caller       arcana.Node
	intermediary arcana.Node
	target       arcana.Node
	peerCaller   arcana.Node
	incoming     int
}

func detectCrossFileIntermediaryBypass(graph arcana.Graph, scopePath string) []Finding {
	calls := buildSymbolCallGraph(graph)
	files := symbolBypassFilePaths(graph.Sources)
	regions := semanticDependencyRegions(graph, files)
	candidates := crossFileIntermediaryBypassCandidates(calls, symbolWrappers(calls), regions)
	if len(candidates) > maximumCrossFileBypassFindings {
		candidates = candidates[:maximumCrossFileBypassFindings]
	}
	findings := make([]Finding, 0, len(candidates))
	for _, candidate := range candidates {
		findings = append(findings, candidate.finding())
	}
	return findings
}

func crossFileIntermediaryBypassCandidates(
	graph symbolCallGraph,
	wrappers []symbolWrapper,
	regions map[string]dependencyRegion,
) []crossFileIntermediaryBypassCandidate {
	candidates := make([]crossFileIntermediaryBypassCandidate, 0)
	for _, wrapper := range wrappers {
		intermediaryFile := normalizedUnitPath(wrapper.intermediary.Path)
		targetFile := normalizedUnitPath(wrapper.target.Path)
		if !sameDependencyRegion(intermediaryFile, targetFile, regions) ||
			sameSymbolName(wrapper.intermediary, wrapper.target) ||
			!shortIntermediary(wrapper.intermediary) ||
			targetCallerFileCount(graph, wrapper.target) > maximumTargetCallerFiles {
			continue
		}
		incoming := len(graph.incoming[wrapper.intermediary.NodeID])
		if incoming < minimumEstablishedIntermediaryCallers {
			continue
		}
		for callerID := range graph.incoming[wrapper.target.NodeID] {
			caller := graph.nodes[callerID]
			callerFile := normalizedUnitPath(caller.Path)
			if caller.NodeID == wrapper.intermediary.NodeID || callerFile == intermediaryFile ||
				!sameDependencyRegion(callerFile, intermediaryFile, regions) {
				continue
			}
			if _, stillUsesIntermediary := graph.outgoing[caller.NodeID][wrapper.intermediary.NodeID]; stillUsesIntermediary {
				continue
			}
			peerCaller, ok := namedSameFileIntermediaryPeer(graph, caller, wrapper.intermediary)
			if !ok {
				continue
			}
			candidates = append(candidates, crossFileIntermediaryBypassCandidate{
				caller: caller, intermediary: wrapper.intermediary, target: wrapper.target,
				peerCaller: peerCaller, incoming: incoming,
			})
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		left, right := candidates[i], candidates[j]
		if left.caller.Path != right.caller.Path {
			return left.caller.Path < right.caller.Path
		}
		if left.caller.Name != right.caller.Name {
			return left.caller.Name < right.caller.Name
		}
		if left.intermediary.Path != right.intermediary.Path {
			return left.intermediary.Path < right.intermediary.Path
		}
		return left.intermediary.Name < right.intermediary.Name
	})
	return candidates
}

func sameDependencyRegion(left, right string, regions map[string]dependencyRegion) bool {
	return dependencyRegionForFile(left, regions).id == dependencyRegionForFile(right, regions).id
}

func sameSymbolName(left, right arcana.Node) bool {
	return strings.EqualFold(strings.TrimSpace(left.Name), strings.TrimSpace(right.Name)) && strings.TrimSpace(left.Name) != ""
}

func shortIntermediary(node arcana.Node) bool {
	if node.Span == nil || node.Span.StartLine <= 0 || node.Span.EndLine < node.Span.StartLine {
		return false
	}
	return node.Span.EndLine-node.Span.StartLine+1 <= maximumIntermediarySpanLines
}

func targetCallerFileCount(graph symbolCallGraph, target arcana.Node) int {
	files := make(map[string]struct{})
	for sourceID := range graph.incoming[target.NodeID] {
		files[normalizedUnitPath(graph.nodes[sourceID].Path)] = struct{}{}
	}
	return len(files)
}

func namedSameFileIntermediaryPeer(graph symbolCallGraph, caller, intermediary arcana.Node) (arcana.Node, bool) {
	peers := make([]arcana.Node, 0)
	intermediaryFile := normalizedUnitPath(intermediary.Path)
	for sourceID := range graph.incoming[intermediary.NodeID] {
		peer := graph.nodes[sourceID]
		if normalizedUnitPath(peer.Path) == intermediaryFile && sameSymbolName(peer, caller) {
			peers = append(peers, peer)
		}
	}
	if len(peers) == 0 {
		return arcana.Node{}, false
	}
	sort.Slice(peers, func(i, j int) bool {
		if peers[i].Name != peers[j].Name {
			return peers[i].Name < peers[j].Name
		}
		return peers[i].NodeID < peers[j].NodeID
	})
	return peers[0], true
}

func (candidate crossFileIntermediaryBypassCandidate) finding() Finding {
	caller := symbolLabel(candidate.caller)
	intermediary := symbolLabel(candidate.intermediary)
	target := symbolLabel(candidate.target)
	peer := symbolLabel(candidate.peerCaller)
	return Finding{
		ID:          findingID(DetectorCrossFileIntermediaryBypass, strings.Join([]string{candidate.caller.Identity, candidate.intermediary.Identity, candidate.target.Identity}, "\x00")),
		Detector:    DetectorCrossFileIntermediaryBypass,
		Disposition: DispositionAdvisory,
		Severity:    SeverityWarning,
		Scope:       Scope{Kind: DetectorCrossFileIntermediaryBypass, Path: normalizedUnitPath(candidate.caller.Path), Name: candidate.caller.Name},
		Summary:     fmt.Sprintf("%s bypasses a cross-file intermediary", caller),
		Rationale: fmt.Sprintf(
			"%s calls %s directly while established same-region intermediary %s still owns that narrow downstream call; same-named peer %s routes through the intermediary alongside %d total callers.",
			caller, target, intermediary, peer, candidate.incoming,
		),
		Evidence: []Evidence{
			{Kind: "direct-call", Message: fmt.Sprintf("%s calls %s directly across files", caller, target)},
			{Kind: "intermediary", Message: fmt.Sprintf("%s is a short one-hop intermediary for %s", intermediary, target)},
			{Kind: "peer-route", Message: fmt.Sprintf("same-file peer %s still routes through %s; intermediary has %d callers", peer, intermediary, candidate.incoming)},
		},
		RequiredOutcome:   "Remove the exceptional direct caller-to-target edge, or remove/refactor the intermediary so the same-region call path no longer presents it as the owning seam.",
		RecommendedAction: fmt.Sprintf("Route %s through %s when that intermediary still owns the call boundary; otherwise redesign or remove the intermediary explicitly.", caller, intermediary),
	}
}
