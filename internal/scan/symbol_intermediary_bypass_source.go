package scan

import (
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func symbolBypassFilePaths(nodes []arcana.Node) map[string]struct{} {
	paths := make(map[string]struct{})
	for _, node := range nodes {
		if node.Kind != "file" {
			continue
		}
		file := normalizedUnitPath(node.Path)
		if file == "" || strings.HasPrefix(file, "@") || excludedSymbolBypassPath(file) {
			continue
		}
		paths[file] = struct{}{}
	}
	return paths
}

func excludedSymbolBypassPath(file string) bool {
	normalized := strings.ToLower(normalizedUnitPath(file))
	segments := strings.Split(normalized, "/")
	for index := 0; index+1 < len(segments); index++ {
		if segments[index] == "addons" && segments[index+1] == "gut" {
			return true
		}
	}
	for _, segment := range segments {
		if segment == "devtools" || segment == "tooling" {
			continue
		}
		if nonProductionSegment(segment) {
			return true
		}
	}
	return strings.Contains(normalized, "/snippets/docs/") || strings.HasPrefix(normalized, "snippets/docs/")
}
