package policy

import (
	"strconv"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func evaluateDependencyRule(
	rule Rule,
	areas map[string]Area,
	graph arcana.Graph,
) *Diagnostic {
	relations := make(map[string]struct{}, len(rule.Relations))
	for _, relation := range rule.Relations {
		relations[relation] = struct{}{}
	}
	evidence := make([]Evidence, 0)
	seen := make(map[string]struct{})
	for _, source := range graph.Sources {
		if !matchesEndpoint(
			source,
			rule.FromAreas,
			rule.FromPaths,
			rule.FromExcludePaths,
			rule.SourceKinds,
			areas,
		) {
			continue
		}
		for _, relationship := range graph.Outgoing[source.NodeID] {
			if _, allowed := relations[relationship.Relation]; !allowed {
				continue
			}
			if !matchesEndpoint(
				relationship.Node,
				rule.ToAreas,
				rule.ToPaths,
				rule.ToExcludePaths,
				rule.TargetKinds,
				areas,
			) {
				continue
			}
			key := dependencyEvidenceKey(source, relationship)
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			target := relationship.Node
			evidence = append(evidence, Evidence{
				Issue:    "forbidden_dependency",
				Source:   source,
				Relation: relationship.Relation,
				Target:   &target,
			})
		}
	}
	return diagnosticFromEvidence(rule, evidence)
}

func dependencyEvidenceKey(source arcana.Node, relationship arcana.Relationship) string {
	return strings.Join([]string{
		strconv.FormatUint(uint64(source.NodeID), 10),
		relationship.Relation,
		strconv.FormatUint(uint64(relationship.Node.NodeID), 10),
	}, "\x00")
}
