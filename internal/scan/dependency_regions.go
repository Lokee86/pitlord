package scan

import (
	"path"
	"sort"

	"github.com/Lokee86/pitlord/internal/arcana"
)

type dependencyRegion struct {
	id    string
	label string
}

func semanticDependencyRegions(graph arcana.Graph, filePaths map[string]struct{}) map[string]dependencyRegion {
	regions := make(map[string]dependencyRegion)
	assignContainerRegions(graph, filePaths, regions, "namespace", false)
	assignContainerRegions(graph, filePaths, regions, "module", true)
	return regions
}

func assignContainerRegions(
	graph arcana.Graph,
	filePaths map[string]struct{},
	regions map[string]dependencyRegion,
	kind string,
	requireMultiFile bool,
) {
	containers := append([]arcana.Node(nil), graph.Sources...)
	sort.Slice(containers, func(i, j int) bool {
		if containers[i].Kind != containers[j].Kind {
			return containers[i].Kind < containers[j].Kind
		}
		return containers[i].Key < containers[j].Key
	})
	for _, container := range containers {
		if container.Kind != kind {
			continue
		}
		members := containerProductionFiles(graph, container, filePaths)
		if requireMultiFile && len(members) < 2 {
			continue
		}
		region := dependencyRegion{
			id:    kind + ":" + container.Key,
			label: semanticRegionLabel(container, members),
		}
		for _, member := range members {
			if _, exists := regions[member]; !exists {
				regions[member] = region
			}
		}
	}
}

func containerProductionFiles(graph arcana.Graph, container arcana.Node, filePaths map[string]struct{}) []string {
	members := make(map[string]struct{})
	for _, relationship := range graph.Outgoing[container.NodeID] {
		if relationship.Relation != "contains" && relationship.Relation != "defines" {
			continue
		}
		member := normalizedUnitPath(relationship.Node.Path)
		if _, ok := filePaths[member]; ok {
			members[member] = struct{}{}
			continue
		}
		if relationship.Node.Span == nil {
			continue
		}
		member = normalizedUnitPath(relationship.Node.Span.Path)
		if _, ok := filePaths[member]; ok {
			members[member] = struct{}{}
		}
	}
	result := make([]string, 0, len(members))
	for member := range members {
		result = append(result, member)
	}
	sort.Strings(result)
	return result
}

func semanticRegionLabel(container arcana.Node, members []string) string {
	if container.Kind == "namespace" && container.Name != "" {
		return "namespace " + container.Name
	}
	if container.Kind == "module" && container.Path != "" {
		return "module " + normalizedUnitPath(container.Path)
	}
	if len(members) > 0 {
		return path.Dir(members[0])
	}
	return container.Kind + " " + container.Key
}

func dependencyRegionForFile(file string, semantic map[string]dependencyRegion) dependencyRegion {
	if region, ok := semantic[file]; ok {
		return region
	}
	directory := path.Dir(file)
	return dependencyRegion{id: "directory:" + directory, label: directory}
}
