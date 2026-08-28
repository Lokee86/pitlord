package scan

import (
	"fmt"
	"sort"

	"github.com/Lokee86/pitlord/internal/arcana"
)

const DetectorBoundaryCohesion = "boundary-cohesion"

const (
	minimumBoundaryCohesionMembers  = 6
	minimumBoundaryCohesionOutgoing = 6
	minimumBoundaryCohesionRegions  = 3
	minimumBoundaryRegionSpread     = 0.25
	minimumBoundaryReachPercentile  = 90.0
	maximumWeakInternalSupportRate  = 0.5
	maximumBoundaryCohesionFindings = 20
)

var boundaryCohesionRelations = map[string]struct{}{
	"references": {}, "imports": {}, "implements": {}, "extends": {},
	"uses-trait": {}, "overrides": {}, "annotates": {}, "includes": {},
	"depends-on": {}, "converts-to": {},
}

var cohesionSupportRelations = map[string]struct{}{
	"references": {}, "imports": {}, "calls": {}, "implements": {},
	"extends": {}, "uses-trait": {}, "overrides": {}, "reads": {},
	"writes": {}, "annotates": {}, "includes": {}, "depends-on": {},
	"converts-to": {},
}

func detectBoundaryCohesion(graph arcana.Graph, scopePath string) []Finding {
	filePaths := repositoryFilePaths(graph.Sources)
	semantic := semanticDependencyRegions(graph, filePaths)
	regions := boundaryCohesionRegions(
		dependencyUnitsForRelations(graph, boundaryCohesionRelations),
		dependencyUnitsForRelations(graph, cohesionSupportRelations),
		semantic,
	)
	if len(regions) < 4 {
		return nil
	}

	eligible := make([]boundaryCohesionRegion, 0, len(regions))
	reachValues := make([]int, 0, len(regions))
	internalRates := make([]float64, 0, len(regions))
	for _, region := range regions {
		if !semanticBoundaryRegion(region.region) || len(region.members) < minimumBoundaryCohesionMembers {
			continue
		}
		eligible = append(eligible, region)
		reachValues = append(reachValues, region.targetRegions)
		internalRates = append(internalRates, region.internalRate)
	}
	if len(eligible) < 4 {
		return nil
	}
	sort.Ints(reachValues)
	sort.Float64s(internalRates)
	internalThreshold := linearQuantile(internalRates, 0.25)

	candidates := make([]boundaryCohesionRegion, 0)
	for _, region := range eligible {
		if region.outgoing < minimumBoundaryCohesionOutgoing ||
			region.targetRegions < minimumBoundaryCohesionRegions ||
			region.regionSpread < minimumBoundaryRegionSpread ||
			percentile(reachValues, region.targetRegions) < minimumBoundaryReachPercentile ||
			region.internalRate > internalThreshold ||
			region.internalRate > maximumWeakInternalSupportRate ||
			region.incoming >= region.outgoing {
			continue
		}
		candidates = append(candidates, region)
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].targetRegions != candidates[j].targetRegions {
			return candidates[i].targetRegions > candidates[j].targetRegions
		}
		if candidates[i].regionSpread != candidates[j].regionSpread {
			return candidates[i].regionSpread > candidates[j].regionSpread
		}
		if candidates[i].internalRate != candidates[j].internalRate {
			return candidates[i].internalRate < candidates[j].internalRate
		}
		return candidates[i].region.id < candidates[j].region.id
	})
	if len(candidates) > maximumBoundaryCohesionFindings {
		candidates = candidates[:maximumBoundaryCohesionFindings]
	}

	findings := make([]Finding, 0, len(candidates))
	for _, candidate := range candidates {
		findings = append(findings, candidate.finding(scopePath, internalThreshold, percentile(reachValues, candidate.targetRegions)))
	}
	return findings
}

func (region boundaryCohesionRegion) finding(scopePath string, internalThreshold, reachPercentile float64) Finding {
	severity := SeverityWarning
	if region.targetRegions >= 6 && region.regionSpread >= 0.5 && region.internalSupport == 0 {
		severity = SeverityHigh
	}
	return Finding{
		ID:          findingID(DetectorBoundaryCohesion, region.region.id),
		Detector:    DetectorBoundaryCohesion,
		Disposition: DispositionAdvisory,
		Severity:    severity,
		Scope:       Scope{Kind: "architecture-region", ID: region.region.id, Path: region.members[0], Name: region.region.label},
		Summary:     "Architectural region has weak internal cohesion and broad boundary spread",
		Rationale:   "A semantic region that reaches across unusually many independent architectural regions while its files have unusually little internal collaboration may not represent a coherent ownership boundary.",
		Evidence: []Evidence{
			{Kind: "region-members", Message: fmt.Sprintf("%d production files in %s", len(region.members), region.region.label)},
			{Kind: "internal-cohesion", Message: fmt.Sprintf("%d internal support relationships (%.2f per file); peer lower quartile is %.2f", region.internalSupport, region.internalRate, internalThreshold)},
			{Kind: "boundary-spread", Message: fmt.Sprintf("%d independent target regions across %d outgoing cross-region dependencies (%.2f target regions per member; %.1fth reach percentile)", region.targetRegions, region.outgoing, region.regionSpread, reachPercentile)},
			{Kind: "outgoing-boundary", Message: fmt.Sprintf("%.0f%% of static owned outgoing relationships cross the region boundary", region.outgoingRate*100)},
			{Kind: "incoming-boundary", Message: fmt.Sprintf("%d incoming cross-region relationships", region.incoming)},
			{Kind: "scope", Message: fmt.Sprintf("peer comparison is within scan scope %s", scopePath)},
		},
		RequiredOutcome:   fmt.Sprintf("Reduce direct ownership spread below %d independent target regions or restore enough internal collaboration that the region is no longer in the weakest peer quartile.", minimumBoundaryCohesionRegions),
		RecommendedAction: "Move behavior to the region that owns it, introduce a narrower boundary abstraction, or split the region if its files do not share a coherent responsibility.",
	}
}
