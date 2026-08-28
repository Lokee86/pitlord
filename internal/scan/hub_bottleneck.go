package scan

import (
	"fmt"
	"math"
	"sort"

	"github.com/Lokee86/pitlord/internal/arcana"
)

const DetectorHubBottleneck = "hub-bottleneck"

const (
	minimumBottleneckPeers             = 8
	minimumBottleneckIncoming          = 8
	minimumBottleneckOutgoing          = 4
	bottleneckMedianMultiplier         = 3.0
	minimumBottleneckPercentile        = 95.0
	minimumBottleneckSourceRegions     = 3
	minimumBehavioralSourceRegions     = 6
	minimumBehavioralTargetRegions     = 4
	minimumBehavioralFlowBalance       = 0.4
	maximumBottleneckCandidateFraction = 0.125
	maximumBottleneckFindings          = 20
)

type hubBottleneckCandidate struct {
	unit                    dependencyUnit
	incoming                int
	outgoing                int
	behavioralIncoming      int
	behavioralOutgoing      int
	behavioralSourceRegions int
	behavioralTargetRegions int
	sourceRegions           int
	percentile              float64
}

var bottleneckBehaviorRelations = map[string]struct{}{
	"calls": {}, "reads": {}, "writes": {},
}

func detectHubBottlenecks(graph arcana.Graph, scopePath string) []Finding {
	units := dependencyUnits(graph)
	if len(units) < minimumBottleneckPeers {
		return nil
	}

	incomingDegrees := activeIncomingDegrees(units)
	if len(incomingDegrees) < minimumBottleneckPeers/2 {
		return nil
	}
	median := medianInt(incomingDegrees)
	trigger := max(minimumBottleneckIncoming, int(math.Ceil(float64(median)*bottleneckMedianMultiplier)))
	semantic := semanticDependencyRegions(graph, repositoryFilePaths(graph.Sources))
	behaviorByPath := dependencyUnitsByPath(dependencyUnitsForRelations(graph, bottleneckBehaviorRelations))

	candidates := make([]hubBottleneckCandidate, 0)
	for _, unit := range units {
		incoming := len(unit.incoming)
		outgoing := len(unit.outgoing)
		if incoming < trigger || outgoing < minimumBottleneckOutgoing {
			continue
		}
		rank := percentile(incomingDegrees, incoming)
		if rank < minimumBottleneckPercentile {
			continue
		}
		sourceRegions := incomingCrossRegionCount(unit, semantic)
		if sourceRegions < minimumBottleneckSourceRegions {
			continue
		}
		behavior := behaviorByPath[unit.path]
		behavioralIncoming := len(behavior.incoming)
		behavioralOutgoing := len(behavior.outgoing)
		behavioralSourceRegions := crossRegionCount(behavior.incoming, unit.path, semantic)
		behavioralTargetRegions := crossRegionCount(behavior.outgoing, unit.path, semantic)
		if behavioralSourceRegions < minimumBehavioralSourceRegions ||
			behavioralTargetRegions < minimumBehavioralTargetRegions ||
			behavioralIncoming == 0 ||
			float64(behavioralOutgoing)/float64(behavioralIncoming) < minimumBehavioralFlowBalance {
			continue
		}
		candidates = append(candidates, hubBottleneckCandidate{
			unit: unit, incoming: incoming, outgoing: outgoing,
			behavioralIncoming: behavioralIncoming, behavioralOutgoing: behavioralOutgoing,
			behavioralSourceRegions: behavioralSourceRegions,
			behavioralTargetRegions: behavioralTargetRegions,
			sourceRegions:           sourceRegions, percentile: rank,
		})
	}
	if float64(len(candidates))/float64(len(units)) > maximumBottleneckCandidateFraction {
		return nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].incoming != candidates[j].incoming {
			return candidates[i].incoming > candidates[j].incoming
		}
		if candidates[i].sourceRegions != candidates[j].sourceRegions {
			return candidates[i].sourceRegions > candidates[j].sourceRegions
		}
		if candidates[i].outgoing != candidates[j].outgoing {
			return candidates[i].outgoing > candidates[j].outgoing
		}
		return candidates[i].unit.path < candidates[j].unit.path
	})
	if len(candidates) > maximumBottleneckFindings {
		candidates = candidates[:maximumBottleneckFindings]
	}

	findings := make([]Finding, 0, len(candidates))
	for _, candidate := range candidates {
		findings = append(findings, candidate.finding(scopePath, trigger, median))
	}
	return findings
}

func activeIncomingDegrees(units []dependencyUnit) []int {
	values := make([]int, 0, len(units))
	for _, unit := range units {
		if degree := len(unit.incoming); degree > 0 {
			values = append(values, degree)
		}
	}
	sort.Ints(values)
	return values
}

func incomingCrossRegionCount(unit dependencyUnit, semantic map[string]dependencyRegion) int {
	return crossRegionCount(unit.incoming, unit.path, semantic)
}

func dependencyUnitsByPath(units []dependencyUnit) map[string]dependencyUnit {
	result := make(map[string]dependencyUnit, len(units))
	for _, unit := range units {
		result[unit.path] = unit
	}
	return result
}

func crossRegionCount(paths map[string]struct{}, ownerPath string, semantic map[string]dependencyRegion) int {
	ownerRegion := dependencyRegionForFile(ownerPath, semantic)
	regions := make(map[string]struct{})
	for candidatePath := range paths {
		region := dependencyRegionForFile(candidatePath, semantic)
		if region.id != ownerRegion.id {
			regions[region.id] = struct{}{}
		}
	}
	return len(regions)
}

func (candidate hubBottleneckCandidate) finding(scopePath string, trigger, median int) Finding {
	severity := SeverityWarning
	if candidate.behavioralIncoming >= 24 && candidate.behavioralOutgoing >= 16 &&
		candidate.behavioralSourceRegions >= 8 && candidate.behavioralTargetRegions >= 6 {
		severity = SeverityHigh
	}
	return Finding{
		ID:          findingID(DetectorHubBottleneck, candidate.unit.path),
		Detector:    DetectorHubBottleneck,
		Disposition: DispositionAdvisory,
		Severity:    severity,
		Scope:       Scope{Kind: "file", ID: candidate.unit.path, Path: candidate.unit.path, Name: candidate.unit.path},
		Summary:     "File is a many-to-many behavioral coordination bottleneck",
		Rationale:   "A production file with extreme direct fan-in that behaviorally receives work from many architectural regions and dispatches behavior across many others forms a central coordination waist rather than merely a shared type or data contract.",
		Evidence: []Evidence{
			{Kind: "incoming-degree", Message: fmt.Sprintf("%d direct production-file dependents; peer median is %d and trigger is %d", candidate.incoming, median, trigger)},
			{Kind: "incoming-percentile", Message: fmt.Sprintf("direct fan-in ranks at the %.1fth percentile among active production peers", candidate.percentile)},
			{Kind: "source-regions", Message: fmt.Sprintf("dependents span %d architectural regions", candidate.sourceRegions)},
			{Kind: "outgoing-degree", Message: fmt.Sprintf("file also directly depends on %d production files", candidate.outgoing)},
			{Kind: "behavioral-flow", Message: fmt.Sprintf("%d behavioral dependents from %d regions; %d behavioral dependencies across %d regions", candidate.behavioralIncoming, candidate.behavioralSourceRegions, candidate.behavioralOutgoing, candidate.behavioralTargetRegions)},
			{Kind: "scope", Message: fmt.Sprintf("peer comparison is within scan scope %s", scopePath)},
		},
		RequiredOutcome:   fmt.Sprintf("Reduce direct fan-in below %d, behavioral source breadth below %d regions, behavioral target breadth below %d regions, or move coordination responsibility out of this shared file.", trigger, minimumBehavioralSourceRegions, minimumBehavioralTargetRegions),
		RecommendedAction: "Split coordination from the shared abstraction, introduce narrower stable interfaces, or move dependent-specific behavior toward the regions that own it.",
	}
}
