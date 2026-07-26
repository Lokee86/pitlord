package policy

import (
	"strconv"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func evaluateOwnershipRule(
	rule Rule,
	areas map[string]Area,
	graph arcana.Graph,
) *Diagnostic {
	evidence := make([]Evidence, 0)
	seen := make(map[string]struct{})
	for _, source := range graph.Sources {
		if !matchesAnyPrefix(source.Path, rule.ScopePaths) ||
			matchesAnyPrefix(source.Path, rule.ScopeExcludePaths) ||
			!matchesKind(source.Kind, rule.SourceKinds) {
			continue
		}
		owners := matchingAreas(source, areas)
		issue := ""
		switch {
		case len(owners) == 0 && !rule.AllowUnowned:
			issue = "unowned"
		case len(owners) > 1 && !rule.AllowOverlaps:
			issue = "multiple_owners"
		default:
			continue
		}
		key := ownershipEvidenceKey(source, issue, owners)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		evidence = append(evidence, Evidence{
			Issue:  issue,
			Source: source,
			Areas:  owners,
		})
	}
	return diagnosticFromEvidence(rule, evidence)
}

func ownershipEvidenceKey(source arcana.Node, issue string, areas []string) string {
	return strings.Join([]string{
		strconv.FormatUint(uint64(source.NodeID), 10),
		issue,
		strings.Join(areas, ","),
	}, "\x00")
}
