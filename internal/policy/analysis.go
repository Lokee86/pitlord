package policy

import (
	"sort"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

type dependencyAccumulator struct {
	count          int
	relationCounts map[string]int
	representative Evidence
	hasEvidence    bool
}

func AreaPaths(document Document) []string {
	seen := make(map[string]struct{})
	paths := make([]string, 0)
	for _, area := range document.Areas {
		for _, path := range area.Paths {
			if _, exists := seen[path]; exists {
				continue
			}
			seen[path] = struct{}{}
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	return paths
}

func AnalyzeAreas(
	document Document,
	graph arcana.Graph,
	options AreaAnalysisOptions,
) AreaAnalysis {
	options.ScopePaths = normalizePrefixes(options.ScopePaths)
	if len(options.ScopePaths) == 0 {
		options.ScopePaths = []string{"."}
	}
	options.ScopeExcludePaths = normalizePrefixes(options.ScopeExcludePaths)
	options.OwnershipKinds = normalizeKinds(options.OwnershipKinds)
	if len(options.OwnershipKinds) == 0 {
		options.OwnershipKinds = []string{"file"}
	}
	options.Relations = normalizeRelations(options.Relations)
	if len(options.Relations) == 0 {
		options.Relations = append([]string(nil), defaultRelations...)
	}
	selectedRelations := make(map[string]struct{}, len(options.Relations))
	for _, relation := range options.Relations {
		selectedRelations[relation] = struct{}{}
	}
	if options.ExampleLimit <= 0 {
		options.ExampleLimit = 20
	}

	areas := areaIndex(document)
	statistics := make(map[string]*AreaStatistics, len(document.Areas))
	for _, area := range document.Areas {
		statistics[area.ID] = &AreaStatistics{
			ID:         area.ID,
			KindCounts: make(map[string]int),
		}
	}
	ownership := OwnershipStatistics{}
	dependencies := make(map[string]*dependencyAccumulator)
	selectedRelationshipCount := 0

	for _, source := range graph.Sources {
		owners := matchingAreas(source, areas)
		for _, owner := range owners {
			current := statistics[owner]
			current.NodeCount++
			current.KindCounts[source.Kind]++
		}
		if matchesAnyPrefix(source.Path, options.ScopePaths) &&
			!matchesAnyPrefix(source.Path, options.ScopeExcludePaths) &&
			matchesKind(source.Kind, options.OwnershipKinds) {
			ownership.CheckedNodes++
			switch len(owners) {
			case 0:
				ownership.UnownedNodes++
				if len(ownership.UnownedExamples) < options.ExampleLimit {
					ownership.UnownedExamples = append(ownership.UnownedExamples, source)
				}
			case 1:
				ownership.OwnedNodes++
			default:
				ownership.OverlappingNodes++
				if len(ownership.OverlapExamples) < options.ExampleLimit {
					ownership.OverlapExamples = append(ownership.OverlapExamples, OwnershipOverlap{
						Node:  source,
						Areas: append([]string(nil), owners...),
					})
				}
			}
		}
		if len(owners) == 0 {
			continue
		}
		for _, relationship := range graph.Outgoing[source.NodeID] {
			if _, selected := selectedRelations[relationship.Relation]; !selected {
				continue
			}
			selectedRelationshipCount++
			targetOwners := matchingAreas(relationship.Node, areas)
			if len(targetOwners) == 0 {
				targetOwners = []string{ExternalArea}
			}
			for _, sourceArea := range owners {
				for _, targetArea := range targetOwners {
					if sourceArea == targetArea {
						statistics[sourceArea].InternalRelationships++
						continue
					}
					key := sourceArea + "\x00" + targetArea
					accumulator := dependencies[key]
					if accumulator == nil {
						accumulator = &dependencyAccumulator{relationCounts: make(map[string]int)}
						dependencies[key] = accumulator
					}
					accumulator.count++
					accumulator.relationCounts[relationship.Relation]++
					target := relationship.Node
					candidate := Evidence{
						Issue:      "area_dependency",
						Source:     source,
						Relation:   relationship.Relation,
						Target:     &target,
						SourceArea: sourceArea,
						TargetArea: targetArea,
					}
					if !accumulator.hasEvidence || areaAnalysisEvidenceLess(candidate, accumulator.representative) {
						accumulator.representative = candidate
						accumulator.hasEvidence = true
					}
				}
			}
		}
	}

	areaStatistics := make([]AreaStatistics, 0, len(document.Areas))
	for _, area := range document.Areas {
		areaStatistics = append(areaStatistics, *statistics[area.ID])
	}
	sort.Slice(areaStatistics, func(i, j int) bool {
		return areaStatistics[i].ID < areaStatistics[j].ID
	})

	areaDependencies := make([]AreaDependency, 0, len(dependencies))
	cycleEdges := make([]areaEdge, 0)
	for key, accumulator := range dependencies {
		parts := strings.SplitN(key, "\x00", 2)
		dependency := AreaDependency{
			SourceArea:     parts[0],
			TargetArea:     parts[1],
			EdgeCount:      accumulator.count,
			RelationCounts: accumulator.relationCounts,
			Representative: accumulator.representative,
		}
		areaDependencies = append(areaDependencies, dependency)
		if dependency.TargetArea != ExternalArea {
			cycleEdges = append(cycleEdges, areaEdge{
				from:     dependency.SourceArea,
				to:       dependency.TargetArea,
				evidence: dependency.Representative,
			})
		}
	}
	sort.Slice(areaDependencies, func(i, j int) bool {
		if areaDependencies[i].SourceArea != areaDependencies[j].SourceArea {
			return areaDependencies[i].SourceArea < areaDependencies[j].SourceArea
		}
		return areaDependencies[i].TargetArea < areaDependencies[j].TargetArea
	})

	areaIDs := make([]string, 0, len(document.Areas))
	for _, area := range document.Areas {
		areaIDs = append(areaIDs, area.ID)
	}
	cycles := cyclicAreaComponents(areaIDs, cycleEdges)
	return AreaAnalysis{
		Schema:       "pitlord.area-analysis.v1",
		ScopePaths:   options.ScopePaths,
		Areas:        areaStatistics,
		Ownership:    ownership,
		Dependencies: areaDependencies,
		Cycles:       cycles,
		Summary: AreaAnalysisSummary{
			SourceNodesScanned:    len(graph.Sources),
			RelationshipsScanned:  graph.Relationships,
			RelationshipsSelected: selectedRelationshipCount,
			AreaCount:             len(areaStatistics),
			DependencyCount:       len(areaDependencies),
			CycleCount:            len(cycles),
		},
	}
}

func areaAnalysisEvidenceLess(left, right Evidence) bool {
	return areaEvidenceSortKey(left) < areaEvidenceSortKey(right)
}
