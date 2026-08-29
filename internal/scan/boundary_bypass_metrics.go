package scan

import (
	"path"
	"sort"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

type bypassGraph struct {
	outgoing map[string]map[string]map[string]struct{}
	incoming map[string]map[string]struct{}
}

func boundaryBypassGraph(graph arcana.Graph, files map[string]struct{}) bypassGraph {
	result := bypassGraph{
		outgoing: make(map[string]map[string]map[string]struct{}),
		incoming: make(map[string]map[string]struct{}),
	}
	for _, source := range graph.Sources {
		sourcePath := normalizedUnitPath(source.Path)
		if _, ok := files[sourcePath]; !ok {
			continue
		}
		for _, relationship := range graph.Outgoing[source.NodeID] {
			if _, ok := boundaryBypassRelations[relationship.Relation]; !ok {
				continue
			}
			targetPath := normalizedUnitPath(relationship.Node.Path)
			if targetPath == sourcePath {
				continue
			}
			if _, ok := files[targetPath]; !ok {
				continue
			}
			if result.outgoing[sourcePath] == nil {
				result.outgoing[sourcePath] = make(map[string]map[string]struct{})
			}
			if result.outgoing[sourcePath][targetPath] == nil {
				result.outgoing[sourcePath][targetPath] = make(map[string]struct{})
			}
			result.outgoing[sourcePath][targetPath][relationship.Relation] = struct{}{}
			if result.incoming[targetPath] == nil {
				result.incoming[targetPath] = make(map[string]struct{})
			}
			result.incoming[targetPath][sourcePath] = struct{}{}
		}
	}
	return result
}

func gatewayPeerGroups(gateway string, incoming map[string]map[string]struct{}, semantic map[string]dependencyRegion) map[string][]string {
	groups := make(map[string][]string)
	for source := range incoming[gateway] {
		region := dependencyRegionForFile(source, semantic)
		groups[region.id] = append(groups[region.id], source)
	}
	for region := range groups {
		sort.Strings(groups[region])
	}
	return groups
}

func directRegionSources(target, regionID string, files map[string]struct{}, edges bypassGraph, semantic map[string]dependencyRegion) []string {
	result := make([]string, 0)
	for source := range edges.incoming[target] {
		if _, ok := files[source]; !ok || dependencyRegionForFile(source, semantic).id != regionID {
			continue
		}
		result = append(result, source)
	}
	sort.Strings(result)
	return result
}

func incomingRegionCount(target string, incoming map[string]map[string]struct{}, semantic map[string]dependencyRegion) int {
	regions := make(map[string]struct{})
	for source := range incoming[target] {
		regions[dependencyRegionForFile(source, semantic).id] = struct{}{}
	}
	return len(regions)
}

func gatewayTargetRegionConcentration(gateway, regionID string, edges bypassGraph, semantic map[string]dependencyRegion) (int, float64) {
	total := len(edges.outgoing[gateway])
	if total == 0 {
		return 0, 0
	}
	matching := 0
	for target := range edges.outgoing[gateway] {
		if dependencyRegionForFile(target, semantic).id == regionID {
			matching++
		}
	}
	return matching, float64(matching) / float64(total)
}

func nestedFileOwnership(source, target string) bool {
	sourceDir := path.Dir(source)
	targetDir := path.Dir(target)
	if sourceDir == "." || targetDir == "." || sourceDir == targetDir {
		return false
	}
	return strings.HasPrefix(sourceDir, targetDir+"/") || strings.HasPrefix(targetDir, sourceDir+"/")
}

func sortedRelationSet(values map[string]struct{}) []string {
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
