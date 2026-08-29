package scan

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

const DetectorBoundaryBypass = "boundary-bypass"

const (
	minimumBypassGatewayPeers     = 3
	maximumBypassShare            = 0.34
	maximumBypassTargetRegions    = 2
	minimumGatewayTargetFiles     = 2
	minimumGatewayTargetRegionPct = 0.40
	maximumBypassFindings         = 20
)

var boundaryBypassRelations = map[string]struct{}{
	"calls": {}, "reads": {}, "writes": {}, "depends-on": {}, "includes": {}, "converts-to": {},
}

type boundaryBypassCandidate struct {
	source        string
	gateway       string
	target        string
	sourceRegion  dependencyRegion
	gatewayRegion dependencyRegion
	targetRegion  dependencyRegion
	gatewayPeers  int
	bypassers     int
	relations     []string
}

func detectBoundaryBypass(graph arcana.Graph, scopePath string) []Finding {
	files := repositoryFilePaths(graph.Sources)
	if len(files) < minimumBypassGatewayPeers+2 {
		return nil
	}
	semantic := semanticDependencyRegions(graph, files)
	candidates := boundaryBypassCandidates(boundaryBypassGraph(graph, files), files, semantic)
	if len(candidates) > maximumBypassFindings {
		candidates = candidates[:maximumBypassFindings]
	}
	findings := make([]Finding, 0, len(candidates))
	for _, candidate := range candidates {
		findings = append(findings, candidate.finding(scopePath))
	}
	return findings
}

func boundaryBypassCandidates(edges bypassGraph, files map[string]struct{}, semantic map[string]dependencyRegion) []boundaryBypassCandidate {
	best := make(map[string]boundaryBypassCandidate)
	for gateway, downstream := range edges.outgoing {
		gatewayRegion := dependencyRegionForFile(gateway, semantic)
		for target := range downstream {
			targetRegion := dependencyRegionForFile(target, semantic)
			targetFiles, targetShare := gatewayTargetRegionConcentration(gateway, targetRegion.id, edges, semantic)
			if gatewayRegion.id == targetRegion.id ||
				incomingRegionCount(target, edges.incoming, semantic) > maximumBypassTargetRegions ||
				targetFiles < minimumGatewayTargetFiles || targetShare < minimumGatewayTargetRegionPct {
				continue
			}
			for sourceRegionID, peers := range gatewayPeerGroups(gateway, edges.incoming, semantic) {
				if len(peers) < minimumBypassGatewayPeers || sourceRegionID == gatewayRegion.id || sourceRegionID == targetRegion.id {
					continue
				}
				bypassers := directRegionSources(target, sourceRegionID, files, edges, semantic)
				if len(bypassers) == 0 || float64(len(bypassers))/float64(len(peers)) > maximumBypassShare {
					continue
				}
				for _, source := range bypassers {
					if nestedFileOwnership(source, target) {
						continue
					}
					candidate := boundaryBypassCandidate{
						source: source, gateway: gateway, target: target,
						sourceRegion: dependencyRegionForFile(source, semantic), gatewayRegion: gatewayRegion, targetRegion: targetRegion,
						gatewayPeers: len(peers), bypassers: len(bypassers), relations: sortedRelationSet(edges.outgoing[source][target]),
					}
					key := source + "\x00" + target
					current, exists := best[key]
					if !exists || candidate.gatewayPeers > current.gatewayPeers || (candidate.gatewayPeers == current.gatewayPeers && candidate.gateway < current.gateway) {
						best[key] = candidate
					}
				}
			}
		}
	}
	result := make([]boundaryBypassCandidate, 0, len(best))
	for _, candidate := range best {
		result = append(result, candidate)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].gatewayPeers != result[j].gatewayPeers {
			return result[i].gatewayPeers > result[j].gatewayPeers
		}
		if result[i].source != result[j].source {
			return result[i].source < result[j].source
		}
		return result[i].target < result[j].target
	})
	return result
}

func (candidate boundaryBypassCandidate) finding(scopePath string) Finding {
	severity := SeverityWarning
	if candidate.gatewayPeers >= 6 && candidate.bypassers == 1 {
		severity = SeverityHigh
	}
	return Finding{
		ID:       findingID(DetectorBoundaryBypass, strings.Join([]string{candidate.source, candidate.gateway, candidate.target}, "\x00")),
		Detector: DetectorBoundaryBypass, Disposition: DispositionAdvisory, Severity: severity,
		Scope:     Scope{Kind: "boundary-bypass", Path: candidate.source, Name: "direct dependency bypass"},
		Summary:   fmt.Sprintf("Direct dependency bypasses an established %s intermediary", candidate.gatewayRegion.label),
		Rationale: fmt.Sprintf("%d peer files from %s use %s before reaching %s, while only %d files in that upstream region reach the downstream target directly. The exceptional direct edge is behavioral rather than import/reference-only evidence.", candidate.gatewayPeers, candidate.sourceRegion.label, candidate.gateway, candidate.targetRegion.label, candidate.bypassers),
		Evidence: []Evidence{
			{Kind: "bypass-edge", Message: fmt.Sprintf("%s -> %s via %s", candidate.source, candidate.target, strings.Join(candidate.relations, ", "))},
			{Kind: "established-intermediary", Message: fmt.Sprintf("%d peer files in %s depend on %s", candidate.gatewayPeers, candidate.sourceRegion.label, candidate.gateway)},
			{Kind: "downstream-edge", Message: fmt.Sprintf("%s -> %s", candidate.gateway, candidate.target)},
			{Kind: "boundary-path", Message: fmt.Sprintf("%s -> %s -> %s", candidate.sourceRegion.label, candidate.gatewayRegion.label, candidate.targetRegion.label)},
		},
		RequiredOutcome:   fmt.Sprintf("Remove the exceptional direct dependency from %s to %s, or make direct %s-to-%s access the deliberate architectural norm rather than a minority bypass.", candidate.source, candidate.target, candidate.sourceRegion.label, candidate.targetRegion.label),
		RecommendedAction: "Route the operation through the established intermediary, move the responsibility to the owning boundary, or intentionally collapse the intermediary if it no longer represents a real architectural seam.",
	}
}
