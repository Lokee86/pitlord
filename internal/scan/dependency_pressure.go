package scan

import (
	"fmt"
	"math"
	"path"
	"sort"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

const DetectorDependencyPressure = "dependency-pressure"

const (
	minimumDependencyPeers      = 8
	minimumHubDegree            = 4
	peerMedianMultiplier        = 3.0
	minimumBoundaryRegions      = 3
	minimumBoundarySpread       = 0.35
	minimumBoundaryPercentile   = 95.0
	minimumCentralFanIn         = 20
	centralFanInOutRatio        = 1.0
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
	boundaryAware := distinctUnitDirectories(units) > 1
	candidates := make([]dependencyCandidate, 0)
	for _, unit := range units {
		degree := len(unit.outgoing)
		if degree < trigger || isConventionalCompositionSeam(unit.path) || isHighlyReusedCentralHub(unit) {
			continue
		}
		regions := outgoingBoundaryRegions(unit)
		spread := float64(regions) / float64(max(1, degree))
		candidatePercentile := percentile(degrees, degree)
		if boundaryAware && (regions < minimumBoundaryRegions || spread < minimumBoundarySpread || candidatePercentile < minimumBoundaryPercentile) {
			continue
		}
		candidates = append(candidates, dependencyCandidate{
			unit:          unit,
			degree:        degree,
			median:        median,
			trigger:       trigger,
			percentile:    candidatePercentile,
			fraction:      float64(degree) / float64(max(1, len(units)-1)),
			boundaryAware: boundaryAware,
			regions:       regions,
			spread:        spread,
		})
	}
	if float64(len(candidates))/float64(len(units)) > maximumHubCandidateFraction {
		return nil
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].regions != candidates[j].regions {
			return candidates[i].regions > candidates[j].regions
		}
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
	unit          dependencyUnit
	degree        int
	median        int
	trigger       int
	percentile    float64
	fraction      float64
	boundaryAware bool
	regions       int
	spread        float64
}

func (candidate dependencyCandidate) finding(scopePath string, unitCount int) Finding {
	severity := SeverityWarning
	if candidate.boundaryAware {
		if !isCompatibilityPath(candidate.unit.path) && candidate.regions >= 6 && candidate.spread >= 0.5 && float64(candidate.degree)/float64(candidate.median) >= 6 {
			severity = SeverityHigh
		}
	} else if candidate.fraction >= 0.25 || float64(candidate.degree)/float64(candidate.median) >= 6 {
		severity = SeverityHigh
	}
	boundaryMessage := fmt.Sprintf("%d distinct target regions across %d outgoing file dependencies (%.0f%% boundary spread)", candidate.regions, candidate.degree, candidate.spread*100)
	if !candidate.boundaryAware {
		boundaryMessage = "single-directory scope; boundary diversity unavailable, using raw topology fallback"
	}
	return Finding{
		ID:          findingID(DetectorDependencyPressure, candidate.unit.path),
		Detector:    DetectorDependencyPressure,
		Disposition: DispositionAdvisory,
		Severity:    severity,
		Scope:       Scope{Kind: "file", Path: candidate.unit.path},
		Summary:     "File has excessive outgoing dependency pressure",
		Rationale:   "Broad outgoing dependencies distributed across architectural regions increase change surface and make the code harder for humans and agents to modify locally.",
		Evidence: []Evidence{
			{Kind: "fan-out-pressure", Message: fmt.Sprintf("%d outgoing file dependencies; active-peer median is %d (%.1fx)", candidate.degree, candidate.median, float64(candidate.degree)/float64(candidate.median))},
			{Kind: "boundary-spread", Message: boundaryMessage},
			{Kind: "fan-out", Message: fmt.Sprintf("%d outgoing file dependencies", len(candidate.unit.outgoing))},
			{Kind: "fan-in", Message: fmt.Sprintf("%d incoming file dependencies within scan scope", len(candidate.unit.incoming))},
			{Kind: "rank", Message: fmt.Sprintf("%.1fth percentile of active outgoing dependency surfaces across %d production files in %s", candidate.percentile, unitCount, scopePath)},
		},
		RequiredOutcome:   fmt.Sprintf("Narrow the outgoing dependency surface below %d direct file dependencies or concentrate it within fewer cohesive architectural regions; current value is %d across %d regions.", candidate.trigger, candidate.degree, candidate.regions),
		RecommendedAction: "Keep orchestration at intentional seams, but move unrelated policy or behavior behind cohesive owners so direct dependencies cross fewer architectural regions.",
	}
}

func distinctUnitDirectories(units []dependencyUnit) int {
	directories := make(map[string]struct{})
	for _, unit := range units {
		directories[path.Dir(unit.path)] = struct{}{}
	}
	return len(directories)
}

func outgoingBoundaryRegions(unit dependencyUnit) int {
	sourceDirectory := path.Dir(unit.path)
	regions := make(map[string]struct{})
	for target := range unit.outgoing {
		targetDirectory := path.Dir(target)
		if targetDirectory == sourceDirectory {
			continue
		}
		regions[targetDirectory] = struct{}{}
	}
	return len(regions)
}

func isConventionalCompositionSeam(filePath string) bool {
	base := path.Base(filePath)
	name := strings.TrimSuffix(base, path.Ext(base))
	lowerName := strings.ToLower(name)
	if lowerName == "main" || lowerName == "app" || lowerName == "application" || lowerName == "app_entry" || lowerName == "entrypoint" || lowerName == "startup" || lowerName == "bootstrap" {
		return true
	}
	return lowerName == "composition" ||
		strings.HasSuffix(lowerName, "_composition") ||
		strings.HasSuffix(lowerName, "-composition") ||
		strings.HasSuffix(name, "Composition") ||
		lowerName == "composer" ||
		strings.HasSuffix(lowerName, "_composer") ||
		strings.HasSuffix(lowerName, "-composer") ||
		strings.HasSuffix(name, "Composer") ||
		strings.HasSuffix(lowerName, "factory") ||
		strings.HasSuffix(lowerName, "controller") ||
		strings.HasSuffix(lowerName, "invoker")
}

func isCompatibilityPath(filePath string) bool {
	for _, segment := range strings.Split(strings.ToLower(normalizedUnitPath(filePath)), "/") {
		switch segment {
		case "compat", "compatibility", "legacy", "deprecated":
			return true
		}
	}
	return false
}

func isHighlyReusedCentralHub(unit dependencyUnit) bool {
	outgoing := len(unit.outgoing)
	incoming := len(unit.incoming)
	if outgoing == 0 || incoming < minimumCentralFanIn {
		return false
	}
	return float64(incoming)/float64(outgoing) >= centralFanInOutRatio
}
