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
	byPath := make(map[string]*dependencyUnit)
	for _, source := range graph.Sources {
		if !isDependencyUnitNode(source) {
			continue
		}
		sourcePath := normalizedUnitPath(source.Path)
		if sourcePath == "" {
			continue
		}
		unit := ensureDependencyUnit(byPath, sourcePath)
		for _, relationship := range graph.Outgoing[source.NodeID] {
			if _, selected := dependencyRelations[relationship.Relation]; !selected || !isDependencyUnitNode(relationship.Node) {
				continue
			}
			targetPath := normalizedUnitPath(relationship.Node.Path)
			if targetPath == "" || targetPath == sourcePath {
				continue
			}
			unit.outgoing[targetPath] = struct{}{}
		}
	}
	for _, source := range graph.Sources {
		if !isDependencyUnitNode(source) {
			continue
		}
		sourcePath := normalizedUnitPath(source.Path)
		for _, relationship := range graph.Outgoing[source.NodeID] {
			if _, selected := dependencyRelations[relationship.Relation]; !selected || !isDependencyUnitNode(relationship.Node) {
				continue
			}
			targetPath := normalizedUnitPath(relationship.Node.Path)
			if target, exists := byPath[targetPath]; exists && targetPath != sourcePath {
				target.incoming[sourcePath] = struct{}{}
			}
		}
	}
	units := make([]dependencyUnit, 0, len(byPath))
	for _, unit := range byPath {
		units = append(units, *unit)
	}
	sort.Slice(units, func(i, j int) bool { return units[i].path < units[j].path })
	return units
}

func isDependencyUnitNode(node arcana.Node) bool {
	return node.Kind != "repository" && node.Kind != "directory"
}

func ensureDependencyUnit(units map[string]*dependencyUnit, path string) *dependencyUnit {
	if unit, exists := units[path]; exists {
		return unit
	}
	unit := &dependencyUnit{path: path, incoming: map[string]struct{}{}, outgoing: map[string]struct{}{}}
	units[path] = unit
	return unit
}

func combinedDegree(unit dependencyUnit) int {
	neighbors := make(map[string]struct{}, len(unit.incoming)+len(unit.outgoing))
	for path := range unit.incoming {
		neighbors[path] = struct{}{}
	}
	for path := range unit.outgoing {
		neighbors[path] = struct{}{}
	}
	return len(neighbors)
}

func activeDegrees(units []dependencyUnit) []int {
	values := make([]int, 0, len(units))
	for _, unit := range units {
		if degree := combinedDegree(unit); degree > 0 {
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
