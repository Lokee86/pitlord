package scan

import (
	"fmt"
	"math/bits"
	"sort"

	"github.com/Lokee86/pitlord/internal/arcana"
)

const DetectorImpactBlastRadius = "impact-blast-radius"

const (
	minimumImpactPeers             = 8
	minimumImpactReach             = 8
	minimumImpactFraction          = 0.25
	minimumImpactPercentile        = 95.0
	minimumImpactAmplification     = 2.0
	minimumImpactRegions           = 3
	maximumImpactCandidateFraction = 0.125
	maximumImpactFindings          = 20
)

type impactCandidate struct {
	unit          dependencyUnit
	reach         int
	regions       int
	fraction      float64
	percentile    float64
	amplification float64
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
		candidates = append(candidates, impactCandidate{
			unit: unit, reach: reach, regions: regions, fraction: fraction,
			percentile: candidatePercentile, amplification: amplification,
		})
	}
	if float64(len(candidates))/float64(len(units)) > maximumImpactCandidateFraction {
		return nil
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].reach != candidates[j].reach {
			return candidates[i].reach > candidates[j].reach
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

func transitiveDependents(units []dependencyUnit) [][]uint64 {
	wordCount := (len(units) + 63) / 64
	indexByPath := make(map[string]int, len(units))
	for index, unit := range units {
		indexByPath[unit.path] = index
	}
	reach := make([][]uint64, len(units))
	for index := range reach {
		reach[index] = make([]uint64, wordCount)
	}
	queue := make([]int, 0, len(units))
	queued := make([]bool, len(units))
	for target, unit := range units {
		for sourcePath := range unit.incoming {
			source := indexByPath[sourcePath]
			bitSetAdd(reach[target], source)
		}
		if bitSetCount(reach[target]) > 0 {
			queue = append(queue, target)
			queued[target] = true
		}
	}
	for len(queue) > 0 {
		source := queue[0]
		queue = queue[1:]
		queued[source] = false
		for targetPath := range units[source].outgoing {
			target := indexByPath[targetPath]
			changed := bitSetAdd(reach[target], source)
			changed = bitSetUnionExcept(reach[target], reach[source], target) || changed
			if changed && !queued[target] {
				queue = append(queue, target)
				queued[target] = true
			}
		}
	}
	return reach
}

func impactedRegionCount(dependents []uint64, units []dependencyUnit, semantic map[string]dependencyRegion) int {
	regions := make(map[string]struct{})
	for index, unit := range units {
		if bitSetHas(dependents, index) {
			regions[dependencyRegionForFile(unit.path, semantic).id] = struct{}{}
		}
	}
	return len(regions)
}

func bitSetAdd(set []uint64, index int) bool {
	word, mask := index/64, uint64(1)<<uint(index%64)
	before := set[word]
	set[word] |= mask
	return before != set[word]
}

func bitSetHas(set []uint64, index int) bool { return set[index/64]&(uint64(1)<<uint(index%64)) != 0 }

func bitSetUnionExcept(target, source []uint64, excluded int) bool {
	changed := false
	excludedWord := excluded / 64
	excludedMask := uint64(1) << uint(excluded%64)
	for index := range target {
		addition := source[index]
		if index == excludedWord {
			addition &^= excludedMask
		}
		before := target[index]
		target[index] |= addition
		changed = changed || before != target[index]
	}
	return changed
}

func bitSetCount(set []uint64) int {
	count := 0
	for _, word := range set {
		count += bits.OnesCount64(word)
	}
	return count
}

func (candidate impactCandidate) finding(scopePath string, unitCount int) Finding {
	severity := SeverityWarning
	if candidate.fraction >= 0.5 && candidate.amplification >= 4 && candidate.regions >= 6 {
		severity = SeverityHigh
	}
	return Finding{
		ID:          findingID(DetectorImpactBlastRadius, candidate.unit.path),
		Detector:    DetectorImpactBlastRadius,
		Disposition: DispositionAdvisory,
		Severity:    severity,
		Scope:       Scope{Kind: "file", ID: candidate.unit.path, Path: candidate.unit.path, Name: candidate.unit.path},
		Summary:     "File has an unusually broad transitive change blast radius",
		Rationale:   "Transitive dependents amplify the change surface well beyond direct fan-in, so modifications to this file can propagate through a large and architecturally broad portion of the repository.",
		Evidence: []Evidence{
			{Kind: "transitive-dependents", Message: fmt.Sprintf("%d transitive production-file dependents (%.0f%% of %d peers)", candidate.reach, candidate.fraction*100, unitCount-1)},
			{Kind: "direct-fan-in", Message: fmt.Sprintf("%d direct dependents; transitive amplification is %.1fx", len(candidate.unit.incoming), candidate.amplification)},
			{Kind: "impact-regions", Message: fmt.Sprintf("transitive dependents span %d architectural regions", candidate.regions)},
			{Kind: "impact-rank", Message: fmt.Sprintf("%.1fth percentile of active transitive impact across %s", candidate.percentile, scopePath)},
		},
		RequiredOutcome:   fmt.Sprintf("Reduce the transitive dependent set below %.0f%% of production peers, reduce transitive amplification below %.1fx direct fan-in, or narrow the affected architecture to fewer than %d regions.", minimumImpactFraction*100, minimumImpactAmplification, minimumImpactRegions),
		RecommendedAction: "Introduce a stable boundary, split volatile responsibility from the shared foundation, or redirect dependencies so changes propagate through fewer transitive consumers.",
	}
}
