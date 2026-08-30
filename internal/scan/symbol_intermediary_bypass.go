package scan

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

const DetectorSymbolIntermediaryBypass = "symbol-intermediary-bypass"

const (
	minimumSymbolWrapperPeers   = 3
	maximumSymbolBypassFindings = 20
)

type symbolIntermediaryBypassCandidate struct {
	caller        arcana.Node
	intermediary  arcana.Node
	target        arcana.Node
	sharedSupport arcana.Node
	peerWrappers  int
}

func detectSymbolIntermediaryBypass(graph arcana.Graph, scopePath string) []Finding {
	calls := buildSymbolCallGraph(graph)
	wrappers := symbolWrappers(calls)
	families := activeWrapperFamilies(calls, wrappers)
	candidates := symbolIntermediaryBypassCandidates(calls, wrappers, families)
	if len(candidates) > maximumSymbolBypassFindings {
		candidates = candidates[:maximumSymbolBypassFindings]
	}
	findings := make([]Finding, 0, len(candidates))
	for _, candidate := range candidates {
		findings = append(findings, candidate.finding(scopePath))
	}
	return findings
}

func symbolIntermediaryBypassCandidates(graph symbolCallGraph, wrappers []symbolWrapper, families map[string]symbolWrapperFamilyMetrics) []symbolIntermediaryBypassCandidate {
	candidates := make([]symbolIntermediaryBypassCandidate, 0)
	for _, wrapper := range wrappers {
		if len(graph.incoming[wrapper.intermediary.NodeID]) != 0 {
			continue
		}
		family := families[wrapperFamilyKey(wrapper)]
		if family.activeWrappers < minimumSymbolWrapperPeers {
			continue
		}
		sourceFile := normalizedUnitPath(wrapper.intermediary.Path)
		for callerID := range graph.incoming[wrapper.target.NodeID] {
			caller := graph.nodes[callerID]
			if caller.NodeID == wrapper.intermediary.NodeID || normalizedUnitPath(caller.Path) != sourceFile {
				continue
			}
			if _, stillUsesIntermediary := graph.outgoing[caller.NodeID][wrapper.intermediary.NodeID]; stillUsesIntermediary {
				continue
			}
			if isOneHopExternalWrapper(graph, caller) || callsIntermediarySibling(graph, caller, wrapper.intermediary) {
				continue
			}
			sharedSupport, ok := establishedWrapperSupport(graph, caller, family)
			if !ok {
				continue
			}
			candidates = append(candidates, symbolIntermediaryBypassCandidate{
				caller: caller, intermediary: wrapper.intermediary, target: wrapper.target, sharedSupport: sharedSupport, peerWrappers: family.activeWrappers,
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
		if left.intermediary.Name != right.intermediary.Name {
			return left.intermediary.Name < right.intermediary.Name
		}
		return left.target.Name < right.target.Name
	})
	return candidates
}

func callsIntermediarySibling(graph symbolCallGraph, caller, intermediary arcana.Node) bool {
	for targetID := range graph.outgoing[caller.NodeID] {
		target := graph.nodes[targetID]
		if target.NodeID == intermediary.NodeID {
			continue
		}
		if normalizedUnitPath(target.Path) == normalizedUnitPath(intermediary.Path) && target.Name == intermediary.Name {
			return true
		}
	}
	return false
}

func isOneHopExternalWrapper(graph symbolCallGraph, node arcana.Node) bool {
	targets := graph.outgoing[node.NodeID]
	if len(targets) != 1 {
		return false
	}
	for targetID := range targets {
		return normalizedUnitPath(graph.nodes[targetID].Path) != normalizedUnitPath(node.Path)
	}
	return false
}

func (candidate symbolIntermediaryBypassCandidate) finding(scopePath string) Finding {
	severity := SeverityWarning
	if candidate.peerWrappers >= 5 {
		severity = SeverityHigh
	}
	caller := symbolLabel(candidate.caller)
	intermediary := symbolLabel(candidate.intermediary)
	target := symbolLabel(candidate.target)
	support := symbolLabel(candidate.sharedSupport)
	return Finding{
		ID:          findingID(DetectorSymbolIntermediaryBypass, strings.Join([]string{candidate.caller.Identity, candidate.intermediary.Identity, candidate.target.Identity}, "\x00")),
		Detector:    DetectorSymbolIntermediaryBypass,
		Disposition: DispositionAdvisory,
		Severity:    severity,
		Scope:       Scope{Kind: "symbol-intermediary-bypass", Path: normalizedUnitPath(candidate.caller.Path), Name: candidate.caller.Name},
		Summary:     fmt.Sprintf("%s bypasses an established local intermediary", caller),
		Rationale: fmt.Sprintf(
			"%s calls %s directly while uncalled local wrapper %s still owns the same downstream call; %d peer wrapper paths preserve the boundary and share support step %s with the caller.",
			caller, target, intermediary, candidate.peerWrappers, support,
		),
		Evidence: []Evidence{
			{Kind: "direct-call", Message: fmt.Sprintf("%s calls %s directly", caller, target)},
			{Kind: "bypassed-intermediary", Message: fmt.Sprintf("%s is an uncalled one-hop wrapper for %s", intermediary, target)},
			{Kind: "peer-pattern", Message: fmt.Sprintf("%d active wrapper paths target the same downstream file and share support step %s with %s", candidate.peerWrappers, support, caller)},
		},
		RequiredOutcome:   "Remove the exceptional direct caller-to-target edge, or remove/refactor the abandoned intermediary so the file no longer contains an established wrapper layer around that target boundary.",
		RecommendedAction: fmt.Sprintf("Route %s through %s when that wrapper still owns the boundary; otherwise remove or redesign the obsolete intermediary explicitly.", caller, intermediary),
	}
}

func symbolLabel(node arcana.Node) string {
	if strings.TrimSpace(node.Name) != "" {
		return node.Name
	}
	if strings.TrimSpace(node.Identity) != "" {
		return node.Identity
	}
	return fmt.Sprintf("symbol %d", node.NodeID)
}
