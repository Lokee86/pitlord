package scan

import (
	"sort"
	"strings"
)

type boundaryCohesionRegion struct {
	region          dependencyRegion
	members         []string
	internal        int
	internalSupport int
	outgoing        int
	incoming        int
	targetRegions   int
	outgoingRate    float64
	internalRate    float64
	regionSpread    float64
}

func boundaryCohesionRegions(
	boundaryUnits []dependencyUnit,
	supportUnits []dependencyUnit,
	semantic map[string]dependencyRegion,
) []boundaryCohesionRegion {
	byPath := make(map[string]dependencyUnit, len(boundaryUnits))
	byRegion := make(map[string]*boundaryCohesionRegion)
	for _, unit := range boundaryUnits {
		byPath[unit.path] = unit
		region := dependencyRegionForFile(unit.path, semantic)
		entry := byRegion[region.id]
		if entry == nil {
			entry = &boundaryCohesionRegion{region: region}
			byRegion[region.id] = entry
		}
		entry.members = append(entry.members, unit.path)
	}

	targets := make(map[string]map[string]struct{}, len(byRegion))
	for _, unit := range boundaryUnits {
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
			if targets[sourceRegion.id] == nil {
				targets[sourceRegion.id] = make(map[string]struct{})
			}
			targets[sourceRegion.id][targetRegion.id] = struct{}{}
		}
	}

	for _, unit := range supportUnits {
		sourceRegion := dependencyRegionForFile(unit.path, semantic)
		source := byRegion[sourceRegion.id]
		if source == nil {
			continue
		}
		for targetPath := range unit.outgoing {
			targetRegion := dependencyRegionForFile(targetPath, semantic)
			if targetRegion.id == sourceRegion.id {
				source.internalSupport++
			}
		}
	}

	result := make([]boundaryCohesionRegion, 0, len(byRegion))
	for id, region := range byRegion {
		sort.Strings(region.members)
		region.targetRegions = len(targets[id])
		region.outgoingRate = float64(region.outgoing) / float64(max(1, region.internal+region.outgoing))
		region.internalRate = float64(region.internalSupport) / float64(max(1, len(region.members)))
		region.regionSpread = float64(region.targetRegions) / float64(max(1, len(region.members)))
		result = append(result, *region)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].region.id < result[j].region.id })
	return result
}

func semanticBoundaryRegion(region dependencyRegion) bool {
	return strings.HasPrefix(region.id, "namespace:") || strings.HasPrefix(region.id, "module:")
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
