package scan

import (
	"crypto/sha256"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

type dependencyUnit struct {
	path     string
	incoming map[string]struct{}
	outgoing map[string]struct{}
}

func dependencyUnits(graph arcana.Graph) []dependencyUnit {
	filePaths := repositoryFilePaths(graph.Sources)
	if len(filePaths) == 0 {
		return nil
	}
	byPath := make(map[string]*dependencyUnit, len(filePaths))
	for filePath := range filePaths {
		ensureDependencyUnit(byPath, filePath)
	}
	for _, source := range graph.Sources {
		sourcePath := normalizedUnitPath(source.Path)
		if _, selected := filePaths[sourcePath]; !selected {
			continue
		}
		unit := byPath[sourcePath]
		for _, relationship := range graph.Outgoing[source.NodeID] {
			if _, selected := dependencyRelations[relationship.Relation]; !selected {
				continue
			}
			targetPath := normalizedUnitPath(relationship.Node.Path)
			if targetPath == sourcePath {
				continue
			}
			if _, selected := filePaths[targetPath]; !selected {
				continue
			}
			unit.outgoing[targetPath] = struct{}{}
			byPath[targetPath].incoming[sourcePath] = struct{}{}
		}
	}
	units := make([]dependencyUnit, 0, len(byPath))
	for _, unit := range byPath {
		units = append(units, *unit)
	}
	sort.Slice(units, func(i, j int) bool { return units[i].path < units[j].path })
	return units
}

func repositoryFilePaths(nodes []arcana.Node) map[string]struct{} {
	paths := make(map[string]struct{})
	for _, node := range nodes {
		if node.Kind != "file" {
			continue
		}
		path := normalizedUnitPath(node.Path)
		if path == "" || strings.HasPrefix(path, "@") || classifySourceRole(path) != sourceRoleProduction {
			continue
		}
		paths[path] = struct{}{}
	}
	return paths
}

func ensureDependencyUnit(units map[string]*dependencyUnit, path string) *dependencyUnit {
	if unit, exists := units[path]; exists {
		return unit
	}
	unit := &dependencyUnit{path: path, incoming: map[string]struct{}{}, outgoing: map[string]struct{}{}}
	units[path] = unit
	return unit
}

func activeOutgoingDegrees(units []dependencyUnit) []int {
	values := make([]int, 0, len(units))
	for _, unit := range units {
		if degree := len(unit.outgoing); degree > 0 {
			values = append(values, degree)
		}
	}
	sort.Ints(values)
	return values
}

func medianInt(values []int) int {
	if len(values) == 0 {
		return 0
	}
	if len(values)%2 == 1 {
		return values[len(values)/2]
	}
	return int(math.Ceil(float64(values[len(values)/2-1]+values[len(values)/2]) / 2))
}

func percentile(values []int, value int) float64 {
	count := sort.Search(len(values), func(index int) bool { return values[index] > value })
	return 100 * float64(count) / float64(len(values))
}

func normalizedUnitPath(value string) string {
	value = strings.Trim(strings.ReplaceAll(strings.TrimSpace(value), "\\", "/"), "/")
	return strings.TrimPrefix(value, "./")
}

func findingID(detector, scope string) string {
	digest := sha256.Sum256([]byte(detector + "\x00" + scope))
	return fmt.Sprintf("%s:%x", detector, digest[:8])
}
