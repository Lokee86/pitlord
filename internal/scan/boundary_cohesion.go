package scan

import (
	"fmt"
	"sort"

	"github.com/Lokee86/pitlord/internal/arcana"
)

const DetectorBoundaryCohesion = "boundary-cohesion"

const maximumBoundaryCohesionFindings = 20

var boundaryCohesionRelations = map[string]struct{}{
	"references": {}, "imports": {}, "implements": {}, "extends": {},
	"uses-trait": {}, "overrides": {}, "annotates": {}, "includes": {},
	"depends-on": {}, "converts-to": {},
}

type boundaryCohesionRegion struct {
	region       dependencyRegion
	members      []string
	internal     int
	outgoing     int
	incoming     int
	outgoingRate float64
	internalRate float64
}

func detectBoundaryCohesion(graph arcana.Graph, scopePath string) []Finding {
	filePaths := repositoryFilePaths(graph.Sources)
	units := dependencyUnitsForRelations(graph, boundaryCohesionRelations)
	regions := boundaryCohesionRegions(units, semanticDependencyRegions(graph, filePaths))
	if len(regions) < 4 {
		return nil
	}

	outgoingRates := make([]float64, 0, len(regions))
	internalRates := make([]float64, 0, len(regions))
	for _, region := range regions {
		outgoingRates = append(outgoingRates, region.outgoingRate)
		internalRates = append(internalRates, region.internalRate)
	}
	upperBoundaryFence := upperOutlierFence(outgoingRates)
	lowerInternalFence := lowerOutlierFence(internalRates)

	candidates := make([]boundaryCohesionRegion, 0)
	for _, region := range regions {
		if len(region.members) < 2 || region.outgoing == 0 {
			continue
		}
		if region.outgoingRate <= upperBoundaryFence || region.internalRate >= lowerInternalFence {
			continue
		}
		if region.incoming > region.outgoing && region.internal > 0 {
			continue
		}
		candidates = append(candidates, region)
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].outgoingRate != candidates[j].outgoingRate {
			return candidates[i].outgoingRate > candidates[j].outgoingRate
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
		findings = append(findings, candidate.finding(scopePath, upperBoundaryFence, lowerInternalFence))
	}
	return findings
}

func boundaryCohesionRegions(units []dependencyUnit, semantic map[string]dependencyRegion) []boundaryCohesionRegion {
	byPath := make(map[string]dependencyUnit, len(units))
	byRegion := make(map[string]*boundaryCohesionRegion)
	for _, unit := range units {
		byPath[unit.path] = unit
		region := dependencyRegionForFile(unit.path, semantic)
		entry := byRegion[region.id]
		if entry == nil {
			entry = &boundaryCohesionRegion{region: region}
			byRegion[region.id] = entry
		}
		entry.members = append(entry.members, unit.path)
	}
	for _, unit := range units {
		sourceRegion := dependencyRegionForFile(unit.path, semantic)
		source := byRegion[sourceRegion.id]
		for targetPath := range unit.outgoing {
			target, ok := byPath[targetPath]
			if !ok {
				continue
			}
			targetRegion := dependencyRegionForFile(target.path, semantic)
			if targetRegion.id == sourceRegion.id {
				source.internal++
				continue
			}
			source.outgoing++
			byRegion[targetRegion.id].incoming++
		}
	}
	result := make([]boundaryCohesionRegion, 0, len(byRegion))
	for _, region := range byRegion {
		sort.Strings(region.members)
		region.outgoingRate = float64(region.outgoing) / float64(max(1, region.internal+region.outgoing))
		region.internalRate = float64(region.internal) / float64(max(1, len(region.members)))
		result = append(result, *region)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].region.id < result[j].region.id })
	return result
}

func upperOutlierFence(values []float64) float64 {
	q1, q3 := quartiles(values)
	return q3 + 1.5*(q3-q1)
}

func lowerOutlierFence(values []float64) float64 {
	q1, q3 := quartiles(values)
	return q1 - 1.5*(q3-q1)
}

func quartiles(values []float64) (float64, float64) {
	ordered := append([]float64(nil), values...)
	sort.Float64s(ordered)
	return linearQuantile(ordered, 0.25), linearQuantile(ordered, 0.75)
}

func linearQuantile(values []float64, quantile float64) float64 {
	if len(values) == 0 {
		return 0
	}
	position := quantile * float64(len(values)-1)
	lower := int(position)
	upper := min(lower+1, len(values)-1)
	fraction := position - float64(lower)
	return values[lower] + (values[upper]-values[lower])*fraction
}

func (region boundaryCohesionRegion) finding(scopePath string, boundaryFence, internalFence float64) Finding {
	severity := SeverityWarning
	if region.outgoingRate >= 0.9 && region.internal == 0 {
		severity = SeverityHigh
	}
	return Finding{
		ID:          findingID(DetectorBoundaryCohesion, region.region.id),
		Detector:    DetectorBoundaryCohesion,
		Disposition: DispositionAdvisory,
		Severity:    severity,
		Scope:       Scope{Kind: "architecture-region", ID: region.region.id, Path: region.members[0], Name: region.region.label},
		Summary:     "Architectural region has weak internal cohesion and excessive outward coupling",
		Rationale:   "A region whose source relationships point outward far more than peer regions while its files have unusually little internal support may not represent a coherent ownership boundary.",
		Evidence: []Evidence{
			{Kind: "region-members", Message: fmt.Sprintf("%d production files in %s", len(region.members), region.region.label)},
			{Kind: "internal-cohesion", Message: fmt.Sprintf("%d internal relationships (%.2f per file); peer lower outlier fence is %.2f", region.internal, region.internalRate, internalFence)},
			{Kind: "outgoing-boundary", Message: fmt.Sprintf("%d outgoing cross-region relationships; %.0f%% of owned outgoing relationships cross the boundary (peer upper outlier fence %.0f%%)", region.outgoing, region.outgoingRate*100, boundaryFence*100)},
			{Kind: "incoming-boundary", Message: fmt.Sprintf("%d incoming cross-region relationships", region.incoming)},
			{Kind: "scope", Message: fmt.Sprintf("peer comparison is within scan scope %s", scopePath)},
		},
		RequiredOutcome:   "Restore a coherent region boundary by increasing local ownership of related behavior or reducing direct outward source dependencies until the region is no longer a peer outlier.",
		RecommendedAction: "Move behavior to the region that owns it, introduce a narrower boundary abstraction, or split the region if its files do not share a coherent responsibility.",
	}
}
