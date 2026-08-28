package scan

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

const DetectorDependencyKnots = "dependency-knots"

const (
	minimumSingleRegionKnotDensity = 0.06
	maximumDependencyKnotFindings  = 20
)

var dependencyKnotRelations = map[string]struct{}{
	"imports": {}, "implements": {}, "extends": {},
	"uses-trait": {}, "overrides": {}, "includes": {}, "depends-on": {},
	"converts-to": {},
}

type dependencyKnot struct {
	members       []string
	internalEdges int
	regions       []string
	density       float64
}

func detectDependencyKnots(graph arcana.Graph, scopePath string) []Finding {
	filePaths := repositoryFilePaths(graph.Sources)
	return detectDependencyKnotsFromUnitsWithRegions(
		dependencyUnitsForRelations(graph, dependencyKnotRelations),
		scopePath,
		semanticDependencyRegions(graph, filePaths),
	)
}

func detectDependencyKnotsFromUnits(units []dependencyUnit, scopePath string) []Finding {
	return detectDependencyKnotsFromUnitsWithRegions(units, scopePath, nil)
}

func detectDependencyKnotsFromUnitsWithRegions(
	units []dependencyUnit,
	scopePath string,
	semanticRegions map[string]dependencyRegion,
) []Finding {
	if len(units) < 2 {
		return nil
	}
	boundaryAware := distinctKnotRegions(units, semanticRegions) > 1
	components := stronglyConnectedDependencyComponents(units)
	knots := make([]dependencyKnot, 0, len(components))
	for _, members := range components {
		if len(members) < 2 {
			continue
		}
		regions := dependencyComponentRegions(members, semanticRegions)
		internalEdges := dependencyComponentInternalEdges(members, units)
		density := float64(internalEdges) / float64(len(members)*(len(members)-1))
		if boundaryAware {
			if len(regions) < 2 {
				continue
			}
		} else if density < minimumSingleRegionKnotDensity {
			continue
		}
		if len(regions) == 2 && conventionalImplementationSeam(members) {
			continue
		}
		knots = append(knots, dependencyKnot{
			members:       members,
			internalEdges: internalEdges,
			regions:       regions,
			density:       density,
		})
	}
	sort.Slice(knots, func(i, j int) bool {
		if len(knots[i].regions) != len(knots[j].regions) {
			return len(knots[i].regions) > len(knots[j].regions)
		}
		if len(knots[i].members) != len(knots[j].members) {
			return len(knots[i].members) > len(knots[j].members)
		}
		if knots[i].density != knots[j].density {
			return knots[i].density > knots[j].density
		}
		return knots[i].members[0] < knots[j].members[0]
	})
	if len(knots) > maximumDependencyKnotFindings {
		knots = knots[:maximumDependencyKnotFindings]
	}
	findings := make([]Finding, 0, len(knots))
	for _, knot := range knots {
		findings = append(findings, knot.finding(scopePath))
	}
	return findings
}

func dependencyComponentRegions(members []string, semantic map[string]dependencyRegion) []string {
	regions := make(map[string]string)
	for _, member := range members {
		region := dependencyRegionForFile(member, semantic)
		regions[region.id] = region.label
	}
	labels := make([]string, 0, len(regions))
	for _, label := range regions {
		labels = append(labels, label)
	}
	sort.Strings(labels)
	return labels
}

func distinctKnotRegions(units []dependencyUnit, semantic map[string]dependencyRegion) int {
	regions := make(map[string]struct{})
	for _, unit := range units {
		regions[dependencyRegionForFile(unit.path, semantic).id] = struct{}{}
	}
	return len(regions)
}

func conventionalImplementationSeam(members []string) bool {
	if len(members) != 2 {
		return false
	}
	first := path.Dir(members[0])
	second := path.Dir(members[1])
	return implementationChildDirectory(first, second) || implementationChildDirectory(second, first)
}

func implementationChildDirectory(parent, child string) bool {
	if !strings.HasPrefix(child, parent+"/") {
		return false
	}
	relative := strings.TrimPrefix(child, parent+"/")
	firstSegment := strings.SplitN(relative, "/", 2)[0]
	switch firstSegment {
	case "impl", "internal", "implementation":
		return true
	default:
		return false
	}
}

func dependencyComponentInternalEdges(members []string, units []dependencyUnit) int {
	memberSet := make(map[string]struct{}, len(members))
	for _, member := range members {
		memberSet[member] = struct{}{}
	}
	count := 0
	for _, unit := range units {
		if _, ok := memberSet[unit.path]; !ok {
			continue
		}
		for target := range unit.outgoing {
			if _, ok := memberSet[target]; ok {
				count++
			}
		}
	}
	return count
}

func (knot dependencyKnot) finding(scopePath string) Finding {
	severity := SeverityWarning
	if len(knot.regions) >= 3 || len(knot.members) >= 8 {
		severity = SeverityHigh
	}
	memberSummary := summarizedPaths(knot.members, 8)
	regionSummary := summarizedPaths(knot.regions, 6)
	componentKey := strings.Join(knot.members, "\x00")
	return Finding{
		ID:          findingID(DetectorDependencyKnots, componentKey),
		Detector:    DetectorDependencyKnots,
		Disposition: DispositionAdvisory,
		Severity:    severity,
		Scope: Scope{
			Kind: "dependency-component",
			Path: knot.members[0],
			Name: fmt.Sprintf("%d-file dependency knot", len(knot.members)),
		},
		Summary:   fmt.Sprintf("%d files form a cyclic dependency knot", len(knot.members)),
		Rationale: "Every file in the component can reach every other through dependency edges, so changes can circulate through the component instead of following one directional boundary.",
		Evidence: []Evidence{
			{Kind: "cycle-members", Message: fmt.Sprintf("%d mutually reachable production files: %s", len(knot.members), memberSummary)},
			{Kind: "cycle-regions", Message: fmt.Sprintf("component spans %d architectural regions: %s", len(knot.regions), regionSummary)},
			{Kind: "cycle-density", Message: fmt.Sprintf("%d internal directed file dependencies (%.1f%% of %d possible)", knot.internalEdges, knot.density*100, len(knot.members)*(len(knot.members)-1))},
		},
		RequiredOutcome:   fmt.Sprintf("Break the strongly connected dependency component rooted at %s so its %d files are no longer mutually reachable through dependency edges.", knot.members[0], len(knot.members)),
		RecommendedAction: "Reverse, remove, or interpose a dependency boundary so the component has a directional dependency flow while preserving cohesive ownership.",
	}
}

func summarizedPaths(values []string, limit int) string {
	if len(values) <= limit {
		return strings.Join(values, ", ")
	}
	return fmt.Sprintf("%s, … (+%d more)", strings.Join(values[:limit], ", "), len(values)-limit)
}
