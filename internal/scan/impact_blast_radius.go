package scan

import (
	"fmt"
	"sort"

	"github.com/Lokee86/pitlord/internal/arcana"
)

const DetectorImpactBlastRadius = "impact-blast-radius"

const (
	minimumImpactPeers                = 8
	minimumImpactReach                = 8
	minimumImpactFraction             = 0.25
	minimumImpactPercentile           = 95.0
	minimumImpactAmplification        = 2.0
	minimumImpactRegions              = 3
	minimumIndependentImpactBranches  = 3
	minimumUniqueImpactBranchFraction = 0.05
	maximumDominantImpactBranchShare  = 0.80
	maximumImpactCandidateFraction    = 0.125
	maximumImpactFindings             = 20
)

type impactCandidate struct {
	unit                dependencyUnit
	reach               int
	regions             int
	fraction            float64
	percentile          float64
	amplification       float64
	independentBranches int
	dominantBranchShare float64
}

func detectImpactBlastRadius(graph arcana.Graph, scopePath string) []Finding {
	units := dependencyUnits(graph)
	if len(units) < minimumImpactPeers {
		return nil
	}
	transitive := transitiveDependents(units)
	active := make([]int, 0, len(units))
	for _, dependents := range transitive {
		if reach := bitSetCount(dependents); reach > 0 {
			active = append(active, reach)
		}
	}
	if len(active) < minimumImpactPeers/2 {
		return nil
	}
	sort.Ints(active)
	semantic := semanticDependencyRegions(graph, repositoryFilePaths(graph.Sources))
	indexByPath := make(map[string]int, len(units))
	for index, unit := range units {
		indexByPath[unit.path] = index
	}
	candidates := make([]impactCandidate, 0)
	for index, unit := range units {
		reach := bitSetCount(transitive[index])
		fraction := float64(reach) / float64(max(1, len(units)-1))
		direct := len(unit.incoming)
		amplification := float64(reach) / float64(max(1, direct))
		candidatePercentile := percentile(active, reach)
		regions := impactedRegionCount(transitive[index], units, semantic)
		if reach < minimumImpactReach || fraction < minimumImpactFraction ||
			candidatePercentile < minimumImpactPercentile || amplification < minimumImpactAmplification ||
			regions < minimumImpactRegions {
			continue
		}
		branches, dominant := independentImpactBranches(index, unit, transitive, units, indexByPath)
		if branches < minimumIndependentImpactBranches || dominant > maximumDominantImpactBranchShare {
			continue
		}
		candidates = append(candidates, impactCandidate{
			unit: unit, reach: reach, regions: regions, fraction: fraction,
			percentile: candidatePercentile, amplification: amplification,
			independentBranches: branches, dominantBranchShare: dominant,
		})
	}
	if float64(len(candidates))/float64(len(units)) > maximumImpactCandidateFraction {
		return nil
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].reach != candidates[j].reach {
			return candidates[i].reach > candidates[j].reach
		}
		if candidates[i].independentBranches != candidates[j].independentBranches {
			return candidates[i].independentBranches > candidates[j].independentBranches
		}
		if candidates[i].regions != candidates[j].regions {
			return candidates[i].regions > candidates[j].regions
		}
		return candidates[i].unit.path < candidates[j].unit.path
	})
	if len(candidates) > maximumImpactFindings {
		candidates = candidates[:maximumImpactFindings]
	}
	findings := make([]Finding, 0, len(candidates))
	for _, candidate := range candidates {
		findings = append(findings, candidate.finding(scopePath, len(units)))
	}
	return findings
}

func (candidate impactCandidate) finding(scopePath string, unitCount int) Finding {
	severity := SeverityWarning
	if candidate.fraction >= 0.5 && candidate.amplification >= 4 &&
		candidate.regions >= 6 && candidate.independentBranches >= 4 {
		severity = SeverityHigh
	}
	return Finding{
		ID:          findingID(DetectorImpactBlastRadius, candidate.unit.path),
		Detector:    DetectorImpactBlastRadius,
		Disposition: DispositionAdvisory,
		Severity:    severity,
		Scope:       Scope{Kind: "file", ID: candidate.unit.path, Path: candidate.unit.path, Name: candidate.unit.path},
		Summary:     "File has an unusually broad independent transitive change blast radius",
		Rationale:   "Several independently expanding dependent branches amplify changes through a large and architecturally broad portion of the repository, rather than inheriting one gateway's downstream closure.",
		Evidence: []Evidence{
			{Kind: "transitive-dependents", Message: fmt.Sprintf("%d transitive production-file dependents (%.0f%% of %d peers)", candidate.reach, candidate.fraction*100, unitCount-1)},
			{Kind: "direct-fan-in", Message: fmt.Sprintf("%d direct dependents; transitive amplification is %.1fx", len(candidate.unit.incoming), candidate.amplification)},
			{Kind: "impact-branches", Message: fmt.Sprintf("%d substantial independent first-hop branches; largest branch contributes %.0f%% of transitive reach", candidate.independentBranches, candidate.dominantBranchShare*100)},
			{Kind: "impact-regions", Message: fmt.Sprintf("transitive dependents span %d architectural regions", candidate.regions)},
			{Kind: "impact-rank", Message: fmt.Sprintf("%.1fth percentile of active transitive impact across %s", candidate.percentile, scopePath)},
		},
		RequiredOutcome:   fmt.Sprintf("Reduce transitive reach below %.0f%% of production peers, reduce substantial independent dependent branches below %d, or concentrate at least %.0f%% of downstream reach behind one stable boundary.", minimumImpactFraction*100, minimumIndependentImpactBranches, maximumDominantImpactBranchShare*100),
		RecommendedAction: "Introduce stable subsystem boundaries or redirect dependent branches so one implementation change cannot propagate independently through several large downstream regions.",
	}
}
