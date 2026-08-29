package scan

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

const DetectorDependencyDepth = "dependency-depth"

const (
	minimumDependencyDepthHops        = 4
	minimumDependencyDepthRegions     = 4
	minimumDependencyDepthTransitions = 3
	minimumDominantDepthGap           = 2
	maximumDependencyDepthFindings    = 20
)

var dependencyDepthRelations = map[string]struct{}{
	"calls": {}, "implements": {}, "extends": {}, "uses-trait": {},
	"overrides": {}, "includes": {}, "depends-on": {}, "converts-to": {},
}

type dependencyDepthCandidate struct {
	paths          []string
	regions        []string
	transitions    int
	singleSteps    int
	dominatedSteps int
	minimumGap     int
}

func detectDependencyDepth(graph arcana.Graph, scopePath string) []Finding {
	files := repositoryFilePaths(graph.Sources)
	return detectDependencyDepthFromUnits(
		dependencyUnitsForRelations(graph, dependencyDepthRelations),
		scopePath,
		semanticDependencyRegions(graph, files),
	)
}

func detectDependencyDepthFromUnits(units []dependencyUnit, scopePath string, semantic map[string]dependencyRegion) []Finding {
	if len(units) < minimumDependencyDepthHops+1 {
		return nil
	}
	cyclic := cyclicDependencyPaths(units)
	byPath := dependencyUnitsByPath(units)
	profiles := dependencyDepthProfiles(units, cyclic)
	candidates := make([]dependencyDepthCandidate, 0)
	for _, unit := range units {
		profile := profiles[unit.path]
		if _, blocked := cyclic[unit.path]; blocked || !profile.dominant || hasDominantPredecessor(unit.path, byPath, profiles, cyclic) {
			continue
		}
		walk := followDominantDependencyChain(unit.path, byPath, profiles, cyclic)
		if len(walk.paths)-1 < minimumDependencyDepthHops {
			continue
		}
		regions, transitions := dependencyChainRegions(walk.paths, semantic)
		if len(regions) < minimumDependencyDepthRegions || transitions < minimumDependencyDepthTransitions {
			continue
		}
		candidates = append(candidates, dependencyDepthCandidate{
			paths: walk.paths, regions: regions, transitions: transitions,
			singleSteps: walk.singleSteps, dominatedSteps: walk.dominatedSteps, minimumGap: walk.minimumGap,
		})
	}
	sort.Slice(candidates, func(i, j int) bool {
		if len(candidates[i].paths) != len(candidates[j].paths) {
			return len(candidates[i].paths) > len(candidates[j].paths)
		}
		if len(candidates[i].regions) != len(candidates[j].regions) {
			return len(candidates[i].regions) > len(candidates[j].regions)
		}
		return candidates[i].paths[0] < candidates[j].paths[0]
	})
	if len(candidates) > maximumDependencyDepthFindings {
		candidates = candidates[:maximumDependencyDepthFindings]
	}
	findings := make([]Finding, 0, len(candidates))
	for _, candidate := range candidates {
		findings = append(findings, candidate.finding(scopePath))
	}
	return findings
}

func dependencyChainRegions(chain []string, semantic map[string]dependencyRegion) ([]string, int) {
	seen := make(map[string]struct{})
	regions := make([]string, 0)
	transitions := 0
	previous := ""
	for _, file := range chain {
		region := dependencyRegionForFile(file, semantic)
		if previous != "" && region.id != previous {
			transitions++
		}
		previous = region.id
		if _, exists := seen[region.id]; exists {
			continue
		}
		seen[region.id] = struct{}{}
		regions = append(regions, region.label)
	}
	return regions, transitions
}

func (candidate dependencyDepthCandidate) finding(scopePath string) Finding {
	hops := len(candidate.paths) - 1
	severity := SeverityWarning
	if hops >= 9 && len(candidate.regions) >= 6 {
		severity = SeverityHigh
	}
	shape := fmt.Sprintf("%d single-dependency steps and %d depth-dominant steps", candidate.singleSteps, candidate.dominatedSteps)
	if candidate.dominatedSteps > 0 {
		shape += fmt.Sprintf("; dominant branches stay at least %d hops deeper than alternatives", candidate.minimumGap)
	}
	return Finding{
		ID:          findingID(DetectorDependencyDepth, strings.Join(candidate.paths, "\x00")),
		Detector:    DetectorDependencyDepth,
		Disposition: DispositionAdvisory,
		Severity:    severity,
		Scope:       Scope{Kind: "dependency-chain", Path: candidate.paths[0], Name: fmt.Sprintf("%d-hop dependency-depth corridor", hops)},
		Summary:     fmt.Sprintf("%d dependency hops form a narrow corridor across %d architectural regions", hops, len(candidate.regions)),
		Rationale:   "At each hop the chain either has one acyclic behavioral dependency or one branch that remains materially deeper than every alternative. Import-only edges do not establish depth, while cycles, shared convergence points, and heavily reused central foundations terminate the walk so incidental or shared infrastructure depth is not inherited by callers.",
		Evidence: []Evidence{
			{Kind: "dependency-chain", Message: summarizedPaths(candidate.paths, 10)},
			{Kind: "chain-regions", Message: fmt.Sprintf("%d regions with %d region transitions: %s", len(candidate.regions), candidate.transitions, summarizedPaths(candidate.regions, 8))},
			{Kind: "chain-shape", Message: shape},
		},
		RequiredOutcome:   fmt.Sprintf("Reduce or deliberately stabilize the %d-hop dependency corridor rooted at %s so routine changes do not require traversing successive intermediary seams.", hops, candidate.paths[0]),
		RecommendedAction: "Collapse pass-through layers, move responsibility to a stable owner, or introduce a deliberate boundary that shortens the dominant dependency path without increasing cyclic coupling.",
	}
}
