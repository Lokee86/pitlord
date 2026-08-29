package scan

import (
	"fmt"
	"math"
)

const (
	minimumDirectImpactIncoming            = 6
	minimumDirectImpactFraction            = 0.03
	minimumDirectImpactPercentile          = 95.0
	minimumDirectImpactRegions             = 2
	minimumDirectImpactEvidence            = 2
	minimumDirectImpactStrongFraction      = 0.05
	minimumDirectImpactExceptionalFraction = 0.20
	minimumDirectImpactBroadRegions        = 10
	minimumDirectPercentilePeers           = 32
)

var impactDirectEvidenceRelations = map[string]struct{}{
	"calls": {}, "implements": {}, "extends": {}, "uses-trait": {},
	"overrides": {}, "reads": {}, "writes": {}, "includes": {},
	"depends-on": {}, "converts-to": {},
}

type directImpactProfile struct {
	incoming           int
	trigger            int
	fraction           float64
	percentile         float64
	regions            int
	meaningfulIncoming int
}

func directImpactExposure(
	unit dependencyUnit,
	peerCount int,
	incomingDegrees []int,
	semantic map[string]dependencyRegion,
	meaningful dependencyUnit,
) (directImpactProfile, bool) {
	incoming := len(unit.incoming)
	trigger := max(minimumDirectImpactIncoming, int(math.Ceil(float64(peerCount)*minimumDirectImpactFraction)))
	if incoming < trigger || len(incomingDegrees) == 0 {
		return directImpactProfile{}, false
	}
	fraction := float64(incoming) / float64(max(1, peerCount))
	rank := percentile(incomingDegrees, incoming)
	regions := incomingCrossRegionCount(unit, semantic)
	if peerCount >= minimumDirectPercentilePeers && rank < minimumDirectImpactPercentile &&
		fraction < minimumDirectImpactExceptionalFraction && regions < minimumDirectImpactBroadRegions {
		return directImpactProfile{}, false
	}
	if regions < minimumDirectImpactRegions {
		return directImpactProfile{}, false
	}
	meaningfulIncoming := len(meaningful.incoming)
	if meaningfulIncoming < minimumDirectImpactEvidence {
		return directImpactProfile{}, false
	}
	if fraction < minimumDirectImpactStrongFraction && regions < minimumDirectImpactBroadRegions {
		return directImpactProfile{}, false
	}
	return directImpactProfile{
		incoming: incoming, trigger: trigger, fraction: fraction,
		percentile: rank, regions: regions, meaningfulIncoming: meaningfulIncoming,
	}, true
}

func (candidate impactCandidate) directFinding(scopePath string, unitCount int) Finding {
	profile := candidate.direct
	severity := SeverityWarning
	if profile.fraction >= 0.25 && profile.regions >= 6 && profile.meaningfulIncoming >= 16 {
		severity = SeverityHigh
	}
	return Finding{
		ID:          findingID(DetectorImpactBlastRadius, candidate.unit.path),
		Detector:    DetectorImpactBlastRadius,
		Disposition: DispositionAdvisory,
		Severity:    severity,
		Scope:       Scope{Kind: "file", ID: candidate.unit.path, Path: candidate.unit.path, Name: candidate.unit.path},
		Summary:     "File has unusually broad direct change exposure",
		Rationale:   "An unusually large and architecturally broad set of production files depends directly on this file, so even a cohesive or intentionally shared contract has a large immediate change surface.",
		Evidence: []Evidence{
			{Kind: "direct-fan-in", Message: fmt.Sprintf("%d direct production-file dependents (%.0f%% of %d peers); trigger is %d", profile.incoming, profile.fraction*100, unitCount-1, profile.trigger)},
			{Kind: "incoming-percentile", Message: fmt.Sprintf("direct fan-in ranks at the %.1fth percentile among active production peers", profile.percentile)},
			{Kind: "source-regions", Message: fmt.Sprintf("direct dependents span %d architectural regions", profile.regions)},
			{Kind: "direct-evidence", Message: fmt.Sprintf("%d direct dependents are supported by behavioral, inheritance, include, or explicit dependency relations", profile.meaningfulIncoming)},
			{Kind: "scope", Message: fmt.Sprintf("peer comparison is within scan scope %s", scopePath)},
		},
		RequiredOutcome:   "Keep changes to this shared surface intentionally compatible, or reduce its direct dependent breadth if the exposure is accidental.",
		RecommendedAction: "Treat this file as a high-impact contract: preserve compatibility deliberately, isolate volatile implementation behind stable boundaries, and split the surface only when its broad dependency role is accidental.",
	}
}
