package scan

import (
	"fmt"
	"math"
	"sort"

	"github.com/Lokee86/pitlord/internal/arcana"
)

const DetectorDependencyPressure = "dependency-pressure"

const (
	minimumDependencyPeers      = 8
	minimumHubDegree            = 4
	peerMedianMultiplier        = 3.0
	maximumHubCandidateFraction = 0.125
	maximumHubFindings          = 20
)

var dependencyRelations = map[string]struct{}{
	"references": {}, "imports": {}, "calls": {}, "implements": {},
	"extends": {}, "uses-trait": {}, "overrides": {}, "reads": {},
	"writes": {}, "annotates": {}, "includes": {}, "depends-on": {},
	"converts-to": {},
}

func detectDependencyPressure(graph arcana.Graph, scopePath string) []Finding {
	units := dependencyUnits(graph)
	if len(units) < minimumDependencyPeers {
		return nil
	}
	degrees := activeOutgoingDegrees(units)
	if len(degrees) < minimumDependencyPeers/2 {
		return nil
	}
	median := medianInt(degrees)
	if median <= 0 {
		return nil
	}
	trigger := max(minimumHubDegree, int(math.Ceil(float64(median)*peerMedianMultiplier)))
	candidates := make([]dependencyCandidate, 0)
	for _, unit := range units {
		degree := len(unit.outgoing)
		if degree < trigger {
			continue
		}
		candidates = append(candidates, dependencyCandidate{
			unit:       unit,
			degree:     degree,
			median:     median,
			trigger:    trigger,
			percentile: percentile(degrees, degree),
			fraction:   float64(degree) / float64(max(1, len(units)-1)),
		})
	}
	if float64(len(candidates))/float64(len(units)) > maximumHubCandidateFraction {
		return nil
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].degree != candidates[j].degree {
			return candidates[i].degree > candidates[j].degree
		}
		return candidates[i].unit.path < candidates[j].unit.path
	})
	if len(candidates) > maximumHubFindings {
		candidates = candidates[:maximumHubFindings]
	}
	findings := make([]Finding, 0, len(candidates))
	for _, candidate := range candidates {
		findings = append(findings, candidate.finding(scopePath, len(units)))
	}
	return findings
}

type dependencyCandidate struct {
	unit       dependencyUnit
	degree     int
	median     int
	trigger    int
	percentile float64
	fraction   float64
}

func (candidate dependencyCandidate) finding(scopePath string, unitCount int) Finding {
	severity := SeverityWarning
	if candidate.fraction >= 0.25 || float64(candidate.degree)/float64(candidate.median) >= 6 {
		severity = SeverityHigh
	}
	return Finding{
		ID:          findingID(DetectorDependencyPressure, candidate.unit.path),
		Detector:    DetectorDependencyPressure,
		Disposition: DispositionAdvisory,
		Severity:    severity,
		Scope:       Scope{Kind: "file", Path: candidate.unit.path},
		Summary:     "File has excessive outgoing dependency pressure",
		Rationale:   "Broad outgoing dependency surfaces force a change to span many architectural concepts and make the code harder for humans and agents to modify locally.",
		Evidence: []Evidence{
			{Kind: "fan-out-pressure", Message: fmt.Sprintf("%d outgoing file dependencies; active-peer median is %d (%.1fx)", candidate.degree, candidate.median, float64(candidate.degree)/float64(candidate.median))},
			{Kind: "fan-out", Message: fmt.Sprintf("%d outgoing file dependencies", len(candidate.unit.outgoing))},
			{Kind: "fan-in", Message: fmt.Sprintf("%d incoming file dependencies within scan scope", len(candidate.unit.incoming))},
			{Kind: "rank", Message: fmt.Sprintf("%.1fth percentile of active outgoing dependency surfaces across %d files in %s", candidate.percentile, unitCount, scopePath)},
		},
		RequiredOutcome:   fmt.Sprintf("Reduce unique outgoing file dependencies below %d; current value is %d.", candidate.trigger, candidate.degree),
		RecommendedAction: "Split orchestration or unrelated responsibilities, or introduce a smaller cohesive boundary that narrows direct dependencies.",
	}
}
