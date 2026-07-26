package policy

import (
	"sort"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func SourcePrefixes(document Document) []string {
	return collectSourcePrefixes(document, true)
}

func RelationshipSourcePrefixes(document Document) []string {
	return collectSourcePrefixes(document, false)
}

func RequiresGraph(document Document) bool {
	for _, rule := range document.Rules {
		switch rule.Type {
		case RuleForbidDependency, RuleRequireOwnership, RuleForbidAreaCycles:
			return true
		}
	}
	return false
}

func RequiresRepository(document Document) bool {
	for _, rule := range document.Rules {
		switch rule.Type {
		case RuleForbidContent, RuleRequirePath, RuleForbidPath:
			return true
		}
	}
	return false
}

func collectSourcePrefixes(document Document, includeOwnership bool) []string {
	areas := areaIndex(document)
	seen := make(map[string]struct{})
	var prefixes []string
	add := func(values []string) {
		for _, prefix := range values {
			if _, exists := seen[prefix]; exists {
				continue
			}
			seen[prefix] = struct{}{}
			prefixes = append(prefixes, prefix)
		}
	}
	for _, rule := range document.Rules {
		switch rule.Type {
		case RuleForbidDependency:
			add(rule.FromPaths)
			for _, areaID := range rule.FromAreas {
				add(areas[areaID].Paths)
			}
		case RuleRequireOwnership:
			if includeOwnership {
				add(rule.ScopePaths)
			}
		case RuleForbidAreaCycles:
			for _, areaID := range rule.Areas {
				add(areas[areaID].Paths)
			}
		}
	}
	sort.Strings(prefixes)
	return prefixes
}

func Evaluate(document Document, graph arcana.Graph, repositoryRoots ...string) []Diagnostic {
	repositoryRoot := ""
	if len(repositoryRoots) > 0 {
		repositoryRoot = repositoryRoots[0]
	}
	areas := areaIndex(document)
	var repositoryFiles []repositoryEntry
	if RequiresRepository(document) {
		repositoryFiles = repositoryEntriesForDocument(repositoryRoot, document)
	}
	diagnostics := make([]Diagnostic, 0)
	for _, rule := range document.Rules {
		var diagnostic *Diagnostic
		switch rule.Type {
		case RuleForbidDependency:
			diagnostic = evaluateDependencyRule(rule, areas, graph)
		case RuleRequireOwnership:
			diagnostic = evaluateOwnershipRule(rule, areas, graph)
		case RuleForbidAreaCycles:
			diagnostic = evaluateAreaCycleRule(rule, areas, graph)
		case RuleForbidContent:
			diagnostic = evaluateContentRule(rule, repositoryRoot, repositoryFiles)
		case RuleRequirePath, RuleForbidPath:
			diagnostic = evaluatePathRule(rule, repositoryFiles)
		}
		if diagnostic != nil {
			diagnostics = append(diagnostics, *diagnostic)
		}
	}
	sort.Slice(diagnostics, func(i, j int) bool {
		if diagnostics[i].Severity != diagnostics[j].Severity {
			return diagnostics[i].Severity < diagnostics[j].Severity
		}
		return diagnostics[i].RuleID < diagnostics[j].RuleID
	})
	return diagnostics
}

func BuildSummary(document Document, graph arcana.Graph, diagnostics []Diagnostic) Summary {
	summary := Summary{
		RulesChecked:         len(document.Rules),
		SourceNodesScanned:   len(graph.Sources),
		RelationshipsScanned: graph.Relationships,
	}
	for _, diagnostic := range diagnostics {
		if diagnostic.Severity == "warning" {
			summary.Warnings++
		} else {
			summary.Errors++
		}
	}
	return summary
}

func CompareExpectation(expected []string, diagnostics []Diagnostic) *Expectation {
	if expected == nil {
		return nil
	}
	actualSet := make(map[string]struct{}, len(diagnostics))
	for _, diagnostic := range diagnostics {
		actualSet[diagnostic.RuleID] = struct{}{}
	}
	expectedSet := make(map[string]struct{}, len(expected))
	for _, id := range expected {
		expectedSet[id] = struct{}{}
	}
	actual := make([]string, 0, len(actualSet))
	missing := make([]string, 0)
	unexpected := make([]string, 0)
	for id := range actualSet {
		actual = append(actual, id)
		if _, exists := expectedSet[id]; !exists {
			unexpected = append(unexpected, id)
		}
	}
	for id := range expectedSet {
		if _, exists := actualSet[id]; !exists {
			missing = append(missing, id)
		}
	}
	sort.Strings(actual)
	sort.Strings(expected)
	sort.Strings(missing)
	sort.Strings(unexpected)
	return &Expectation{
		ExpectedIDs: expected,
		ActualIDs:   actual,
		MissingIDs:  missing,
		Unexpected:  unexpected,
		Matched:     len(missing) == 0 && len(unexpected) == 0,
	}
}

func areaIndex(document Document) map[string]Area {
	areas := make(map[string]Area, len(document.Areas))
	for _, area := range document.Areas {
		areas[area.ID] = area
	}
	return areas
}
