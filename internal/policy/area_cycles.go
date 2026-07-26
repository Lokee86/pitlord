package policy

import (
	"sort"
	"strconv"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

type areaEdge struct {
	from     string
	to       string
	evidence Evidence
}

func evaluateAreaCycleRule(
	rule Rule,
	areas map[string]Area,
	graph arcana.Graph,
) *Diagnostic {
	selected := make(map[string]Area, len(rule.Areas))
	for _, areaID := range rule.Areas {
		selected[areaID] = areas[areaID]
	}
	relations := make(map[string]struct{}, len(rule.Relations))
	for _, relation := range rule.Relations {
		relations[relation] = struct{}{}
	}

	edges := collectAreaEdges(rule, selected, relations, graph)
	components := cyclicAreaComponents(rule.Areas, edges)
	if len(components) == 0 {
		return nil
	}

	evidence := make([]Evidence, 0)
	for _, component := range components {
		members := make(map[string]struct{}, len(component))
		for _, areaID := range component {
			members[areaID] = struct{}{}
		}
		for _, edge := range edges {
			if _, exists := members[edge.from]; !exists {
				continue
			}
			if _, exists := members[edge.to]; !exists {
				continue
			}
			current := edge.evidence
			current.Areas = append([]string(nil), component...)
			evidence = append(evidence, current)
		}
	}
	return diagnosticFromEvidence(rule, evidence)
}

func collectAreaEdges(
	rule Rule,
	areas map[string]Area,
	relations map[string]struct{},
	graph arcana.Graph,
) []areaEdge {
	best := make(map[string]areaEdge)
	for _, source := range graph.Sources {
		if !matchesKind(source.Kind, rule.SourceKinds) {
			continue
		}
		sourceAreas := matchingSelectedAreas(source, areas)
		if len(sourceAreas) == 0 {
			continue
		}
		for _, relationship := range graph.Outgoing[source.NodeID] {
			if _, allowed := relations[relationship.Relation]; !allowed {
				continue
			}
			if !matchesKind(relationship.Node.Kind, rule.TargetKinds) {
				continue
			}
			targetAreas := matchingSelectedAreas(relationship.Node, areas)
			for _, sourceArea := range sourceAreas {
				for _, targetArea := range targetAreas {
					if sourceArea == targetArea {
						continue
					}
					target := relationship.Node
					candidate := areaEdge{
						from: sourceArea,
						to:   targetArea,
						evidence: Evidence{
							Issue:      "area_cycle",
							Source:     source,
							Relation:   relationship.Relation,
							Target:     &target,
							SourceArea: sourceArea,
							TargetArea: targetArea,
						},
					}
					key := sourceArea + "\x00" + targetArea
					current, exists := best[key]
					if !exists || areaEdgeLess(candidate, current) {
						best[key] = candidate
					}
				}
			}
		}
	}
	edges := make([]areaEdge, 0, len(best))
	for _, edge := range best {
		edges = append(edges, edge)
	}
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].from != edges[j].from {
			return edges[i].from < edges[j].from
		}
		return edges[i].to < edges[j].to
	})
	return edges
}

func matchingSelectedAreas(node arcana.Node, areas map[string]Area) []string {
	matches := make([]string, 0)
	for areaID, area := range areas {
		if matchesArea(node, area) {
			matches = append(matches, areaID)
		}
	}
	sort.Strings(matches)
	return matches
}

func areaEdgeLess(left, right areaEdge) bool {
	return areaEvidenceSortKey(left.evidence) < areaEvidenceSortKey(right.evidence)
}

func areaEvidenceSortKey(evidence Evidence) string {
	line := 0
	column := 0
	if evidence.Source.Span != nil {
		line = evidence.Source.Span.StartLine
		column = evidence.Source.Span.StartColumn
	}
	targetPath := ""
	targetName := ""
	if evidence.Target != nil {
		targetPath = evidence.Target.Path
		targetName = evidence.Target.Name
	}
	return strings.Join([]string{
		evidence.Source.Path,
		strconv.Itoa(line),
		strconv.Itoa(column),
		evidence.Source.Name,
		evidence.Relation,
		targetPath,
		targetName,
	}, "\x00")
}

func cyclicAreaComponents(areaIDs []string, edges []areaEdge) [][]string {
	adjacency := make(map[string][]string, len(areaIDs))
	for _, areaID := range areaIDs {
		adjacency[areaID] = nil
	}
	for _, edge := range edges {
		adjacency[edge.from] = append(adjacency[edge.from], edge.to)
	}
	for areaID := range adjacency {
		adjacency[areaID] = uniqueSorted(adjacency[areaID])
	}

	index := 0
	indices := make(map[string]int, len(areaIDs))
	lowlinks := make(map[string]int, len(areaIDs))
	onStack := make(map[string]bool, len(areaIDs))
	stack := make([]string, 0, len(areaIDs))
	components := make([][]string, 0)
	var visit func(string)
	visit = func(areaID string) {
		indices[areaID] = index
		lowlinks[areaID] = index
		index++
		stack = append(stack, areaID)
		onStack[areaID] = true

		for _, neighbor := range adjacency[areaID] {
			neighborIndex, visited := indices[neighbor]
			if !visited {
				visit(neighbor)
				if lowlinks[neighbor] < lowlinks[areaID] {
					lowlinks[areaID] = lowlinks[neighbor]
				}
			} else if onStack[neighbor] && neighborIndex < lowlinks[areaID] {
				lowlinks[areaID] = neighborIndex
			}
		}

		if lowlinks[areaID] != indices[areaID] {
			return
		}
		component := make([]string, 0)
		for {
			last := len(stack) - 1
			member := stack[last]
			stack = stack[:last]
			onStack[member] = false
			component = append(component, member)
			if member == areaID {
				break
			}
		}
		if len(component) > 1 {
			sort.Strings(component)
			components = append(components, component)
		}
	}

	sortedAreas := append([]string(nil), areaIDs...)
	sort.Strings(sortedAreas)
	for _, areaID := range sortedAreas {
		if _, visited := indices[areaID]; !visited {
			visit(areaID)
		}
	}
	sort.Slice(components, func(i, j int) bool {
		return strings.Join(components[i], "\x00") < strings.Join(components[j], "\x00")
	})
	return components
}

func uniqueSorted(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	sort.Strings(values)
	result := values[:0]
	for _, value := range values {
		if len(result) == 0 || result[len(result)-1] != value {
			result = append(result, value)
		}
	}
	return result
}
