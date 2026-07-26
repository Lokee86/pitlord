package policy

import (
	"sort"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func matchesEndpoint(
	node arcana.Node,
	areaIDs []string,
	paths []string,
	excludePaths []string,
	kinds []string,
	areas map[string]Area,
) bool {
	if !matchesKind(node.Kind, kinds) {
		return false
	}
	if len(paths) > 0 && matchesAnyPrefix(node.Path, paths) && !matchesAnyPrefix(node.Path, excludePaths) {
		return true
	}
	for _, areaID := range areaIDs {
		if matchesArea(node, areas[areaID]) {
			return true
		}
	}
	return false
}

func matchingAreas(node arcana.Node, areas map[string]Area) []string {
	matches := make([]string, 0)
	for id, area := range areas {
		if matchesArea(node, area) {
			matches = append(matches, id)
		}
	}
	sort.Strings(matches)
	return matches
}

func matchesArea(node arcana.Node, area Area) bool {
	return matchesAnyPrefix(node.Path, area.Paths) &&
		!matchesAnyPrefix(node.Path, area.ExcludePaths) &&
		matchesKind(node.Kind, area.Kinds)
}

func matchesAnyPrefix(path string, prefixes []string) bool {
	if len(prefixes) == 0 {
		return false
	}
	path = normalizePath(path)
	for _, prefix := range prefixes {
		if prefix == "." || path == prefix || strings.HasPrefix(path, prefix+"/") {
			return true
		}
	}
	return false
}

func matchesKind(kind string, kinds []string) bool {
	if len(kinds) == 0 {
		return true
	}
	kind = strings.ToLower(strings.TrimSpace(kind))
	for _, candidate := range kinds {
		if kind == candidate {
			return true
		}
	}
	return false
}
