package scan

import "sort"

type dependencyDirectionRegion struct {
	region   dependencyRegion
	members  []string
	incoming map[string]struct{}
	outgoing map[string]struct{}
}

type dependencyDirectionBoundary struct {
	sourceID    string
	targetID    string
	edges       int
	sourceFiles map[string]struct{}
	targetFiles map[string]struct{}
}

func (region dependencyDirectionRegion) instability() float64 {
	total := len(region.incoming) + len(region.outgoing)
	if total == 0 {
		return 0
	}
	return float64(len(region.outgoing)) / float64(total)
}

func dependencyDirectionTopology(
	units []dependencyUnit,
	semantic map[string]dependencyRegion,
) (map[string]*dependencyDirectionRegion, map[string]*dependencyDirectionBoundary) {
	regions := make(map[string]*dependencyDirectionRegion)
	boundaries := make(map[string]*dependencyDirectionBoundary)
	for _, unit := range units {
		region := dependencyRegionForFile(unit.path, semantic)
		current := ensureDependencyDirectionRegion(regions, region)
		current.members = append(current.members, unit.path)
	}
	for _, unit := range units {
		sourceRegion := dependencyRegionForFile(unit.path, semantic)
		for targetPath := range unit.outgoing {
			targetRegion := dependencyRegionForFile(targetPath, semantic)
			if sourceRegion.id == targetRegion.id {
				continue
			}
			source := ensureDependencyDirectionRegion(regions, sourceRegion)
			target := ensureDependencyDirectionRegion(regions, targetRegion)
			source.outgoing[targetRegion.id] = struct{}{}
			target.incoming[sourceRegion.id] = struct{}{}

			key := directionBoundaryKey(sourceRegion.id, targetRegion.id)
			boundary := boundaries[key]
			if boundary == nil {
				boundary = &dependencyDirectionBoundary{
					sourceID: sourceRegion.id, targetID: targetRegion.id,
					sourceFiles: make(map[string]struct{}), targetFiles: make(map[string]struct{}),
				}
				boundaries[key] = boundary
			}
			boundary.edges++
			boundary.sourceFiles[unit.path] = struct{}{}
			boundary.targetFiles[targetPath] = struct{}{}
		}
	}
	for _, region := range regions {
		sort.Strings(region.members)
	}
	return regions, boundaries
}

func ensureDependencyDirectionRegion(
	regions map[string]*dependencyDirectionRegion,
	region dependencyRegion,
) *dependencyDirectionRegion {
	if current := regions[region.id]; current != nil {
		return current
	}
	current := &dependencyDirectionRegion{
		region: region, incoming: make(map[string]struct{}), outgoing: make(map[string]struct{}),
	}
	regions[region.id] = current
	return current
}

func dependencyDirectionCycleMembership(regions map[string]*dependencyDirectionRegion) map[string]string {
	units := make([]dependencyUnit, 0, len(regions))
	for id, region := range regions {
		units = append(units, dependencyUnit{path: id, incoming: region.incoming, outgoing: region.outgoing})
	}
	membership := make(map[string]string)
	for _, component := range stronglyConnectedDependencyComponents(units) {
		if len(component) < 2 {
			continue
		}
		componentID := component[0]
		for _, regionID := range component {
			membership[regionID] = componentID
		}
	}
	return membership
}
