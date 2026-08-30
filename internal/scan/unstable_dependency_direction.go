package scan

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

const DetectorUnstableDependencyDirection = "unstable-dependency-direction"

const (
	minimumDirectionRegionMembers   = 3
	minimumDirectionBoundarySources = 2
	minimumDirectionCouplingTotal   = 3
	minimumDirectionInstabilityGap  = 0.25
	maximumDirectionFindings        = 20
)

type unstableDependencyCandidate struct {
	source   *dependencyDirectionRegion
	target   *dependencyDirectionRegion
	boundary *dependencyDirectionBoundary
	gap      float64
}

func detectUnstableDependencyDirection(graph arcana.Graph, scopePath string) []Finding {
	files := repositoryFilePaths(graph.Sources)
	if len(files) < minimumDirectionRegionMembers*2 {
		return nil
	}
	semantic := semanticDependencyRegions(graph, files)
	regions, boundaries := dependencyDirectionTopology(
		dependencyUnitsForRelations(graph, dependencyKnotRelations), semantic,
	)
	candidates := unstableDependencyCandidates(regions, boundaries)
	if len(candidates) > maximumDirectionFindings {
		candidates = candidates[:maximumDirectionFindings]
	}
	findings := make([]Finding, 0, len(candidates))
	for _, candidate := range candidates {
		findings = append(findings, candidate.finding(scopePath))
	}
	return findings
}

func unstableDependencyCandidates(
	regions map[string]*dependencyDirectionRegion,
	boundaries map[string]*dependencyDirectionBoundary,
) []unstableDependencyCandidate {
	candidates := make([]unstableDependencyCandidate, 0)
	cycleMembership := dependencyDirectionCycleMembership(regions)
	for _, boundary := range boundaries {
		source := regions[boundary.sourceID]
		target := regions[boundary.targetID]
		if source == nil || target == nil ||
			len(source.members) < minimumDirectionRegionMembers ||
			len(target.members) < minimumDirectionRegionMembers ||
			len(boundary.sourceFiles) < minimumDirectionBoundarySources ||
			len(source.incoming)+len(source.outgoing) < minimumDirectionCouplingTotal ||
			len(target.incoming)+len(target.outgoing) < minimumDirectionCouplingTotal {
			continue
		}
		if len(source.incoming) <= len(source.outgoing) || len(target.outgoing) <= len(target.incoming) {
			continue
		}
		if sourceCycle := cycleMembership[boundary.sourceID]; sourceCycle != "" && sourceCycle == cycleMembership[boundary.targetID] {
			continue
		}
		gap := target.instability() - source.instability()
		if gap < minimumDirectionInstabilityGap {
			continue
		}
		candidates = append(candidates, unstableDependencyCandidate{
			source: source, target: target, boundary: boundary, gap: gap,
		})
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].gap != candidates[j].gap {
			return candidates[i].gap > candidates[j].gap
		}
		if candidates[i].boundary.edges != candidates[j].boundary.edges {
			return candidates[i].boundary.edges > candidates[j].boundary.edges
		}
		return candidates[i].source.region.id < candidates[j].source.region.id
	})
	return candidates
}

func (candidate unstableDependencyCandidate) finding(scopePath string) Finding {
	severity := SeverityWarning
	if candidate.gap >= 0.50 && len(candidate.source.incoming) >= 3 && len(candidate.target.outgoing) >= 3 {
		severity = SeverityHigh
	}
	sourceInstability := candidate.source.instability()
	targetInstability := candidate.target.instability()
	boundarySources := sortedStringSet(candidate.boundary.sourceFiles)
	boundaryTargets := sortedStringSet(candidate.boundary.targetFiles)
	return Finding{
		ID:          findingID(DetectorUnstableDependencyDirection, candidate.source.region.id+"\x00"+candidate.target.region.id),
		Detector:    DetectorUnstableDependencyDirection,
		Disposition: DispositionAdvisory,
		Severity:    severity,
		Scope: Scope{
			Kind: "architecture-dependency", ID: candidate.source.region.id + "->" + candidate.target.region.id,
			Path: boundarySources[0], Name: candidate.source.region.label + " -> " + candidate.target.region.label,
		},
		Summary:   fmt.Sprintf("Stable %s depends on less stable %s", candidate.source.region.label, candidate.target.region.label),
		Rationale: "An incoming-heavy region is comparatively stable because other regions depend on it. Making that region depend on an outgoing-heavy region points source dependency toward greater instability and makes the stable side sensitive to a more volatile boundary.",
		Evidence: []Evidence{
			{Kind: "dependency-direction", Message: fmt.Sprintf("%s -> %s through %d static file dependencies from %d source files to %d target files", candidate.source.region.label, candidate.target.region.label, candidate.boundary.edges, len(boundarySources), len(boundaryTargets))},
			{Kind: "source-stability", Message: fmt.Sprintf("%s has %d incoming regions and %d outgoing regions (instability %.2f)", candidate.source.region.label, len(candidate.source.incoming), len(candidate.source.outgoing), sourceInstability)},
			{Kind: "target-stability", Message: fmt.Sprintf("%s has %d incoming regions and %d outgoing regions (instability %.2f)", candidate.target.region.label, len(candidate.target.incoming), len(candidate.target.outgoing), targetInstability)},
			{Kind: "instability-gap", Message: fmt.Sprintf("dependency points across a %.2f instability increase; detector minimum is %.2f", candidate.gap, minimumDirectionInstabilityGap)},
			{Kind: "scope", Message: fmt.Sprintf("region coupling is measured within scan scope %s", scopePath)},
		},
		RequiredOutcome:   fmt.Sprintf("Remove the static %s-to-%s dependency, or change the region coupling so it no longer points from an incoming-heavy region to an outgoing-heavy region across an instability gap of %.2f or more.", candidate.source.region.label, candidate.target.region.label, minimumDirectionInstabilityGap),
		RecommendedAction: "Invert the dependency through an abstraction owned by the stable side, move the responsibility toward the stable boundary, or revise region ownership if the current direction is intentional.",
	}
}

func sortedStringSet(values map[string]struct{}) []string {
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func directionBoundaryKey(sourceID, targetID string) string {
	return strings.Join([]string{sourceID, targetID}, "\x00")
}
