package scan

import (
	"sort"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

type symbolCallGraph struct {
	nodes    map[uint32]arcana.Node
	outgoing map[uint32]map[uint32]struct{}
	incoming map[uint32]map[uint32]struct{}
}

type symbolWrapper struct {
	intermediary arcana.Node
	target       arcana.Node
}

func buildSymbolCallGraph(graph arcana.Graph) symbolCallGraph {
	files := symbolBypassFilePaths(graph.Sources)
	result := symbolCallGraph{
		nodes:    make(map[uint32]arcana.Node),
		outgoing: make(map[uint32]map[uint32]struct{}),
		incoming: make(map[uint32]map[uint32]struct{}),
	}
	for _, node := range graph.Sources {
		if _, ok := files[normalizedUnitPath(node.Path)]; ok && isCallableSymbol(node) {
			result.nodes[node.NodeID] = node
		}
	}
	for sourceID, relationships := range graph.Outgoing {
		if _, ok := result.nodes[sourceID]; !ok {
			continue
		}
		for _, relationship := range relationships {
			if relationship.Relation != "calls" || !isCallableSymbol(relationship.Node) {
				continue
			}
			targetPath := normalizedUnitPath(relationship.Node.Path)
			if _, ok := files[targetPath]; !ok {
				continue
			}
			target := relationship.Node
			if sourceNode, exists := result.nodes[target.NodeID]; exists {
				target = sourceNode
			} else {
				result.nodes[target.NodeID] = target
			}
			if result.outgoing[sourceID] == nil {
				result.outgoing[sourceID] = make(map[uint32]struct{})
			}
			if result.incoming[target.NodeID] == nil {
				result.incoming[target.NodeID] = make(map[uint32]struct{})
			}
			result.outgoing[sourceID][target.NodeID] = struct{}{}
			result.incoming[target.NodeID][sourceID] = struct{}{}
		}
	}
	return result
}

func symbolWrappers(graph symbolCallGraph) []symbolWrapper {
	wrappers := make([]symbolWrapper, 0)
	for _, intermediary := range graph.nodes {
		targets := graph.outgoing[intermediary.NodeID]
		if len(targets) != 1 {
			continue
		}
		for targetID := range targets {
			target := graph.nodes[targetID]
			if strings.EqualFold(strings.TrimSpace(target.Kind), "constructor") || normalizedUnitPath(target.Path) == normalizedUnitPath(intermediary.Path) {
				continue
			}
			wrappers = append(wrappers, symbolWrapper{intermediary: intermediary, target: target})
		}
	}
	sort.Slice(wrappers, func(i, j int) bool {
		left, right := wrappers[i], wrappers[j]
		if left.intermediary.Path != right.intermediary.Path {
			return left.intermediary.Path < right.intermediary.Path
		}
		if left.intermediary.Name != right.intermediary.Name {
			return left.intermediary.Name < right.intermediary.Name
		}
		return left.intermediary.NodeID < right.intermediary.NodeID
	})
	return wrappers
}

func sameFileIncoming(graph symbolCallGraph, node arcana.Node) int {
	count := 0
	file := normalizedUnitPath(node.Path)
	for sourceID := range graph.incoming[node.NodeID] {
		if normalizedUnitPath(graph.nodes[sourceID].Path) == file {
			count++
		}
	}
	return count
}

type symbolWrapperFamilyMetrics struct {
	activeWrappers int
	support        map[uint32]int
}

func activeWrapperFamilies(graph symbolCallGraph, wrappers []symbolWrapper) map[string]symbolWrapperFamilyMetrics {
	families := make(map[string]symbolWrapperFamilyMetrics)
	for _, wrapper := range wrappers {
		callers := sameFileCallers(graph, wrapper.intermediary)
		if len(callers) == 0 {
			continue
		}
		key := wrapperFamilyKey(wrapper)
		metrics := families[key]
		if metrics.support == nil {
			metrics.support = make(map[uint32]int)
		}
		metrics.activeWrappers++
		seenSupport := make(map[uint32]struct{})
		for _, caller := range callers {
			for targetID := range graph.outgoing[caller.NodeID] {
				target := graph.nodes[targetID]
				if targetID == wrapper.intermediary.NodeID || normalizedUnitPath(target.Path) == normalizedUnitPath(wrapper.target.Path) {
					continue
				}
				seenSupport[targetID] = struct{}{}
			}
		}
		for targetID := range seenSupport {
			metrics.support[targetID]++
		}
		families[key] = metrics
	}
	return families
}

func sameFileCallers(graph symbolCallGraph, node arcana.Node) []arcana.Node {
	callers := make([]arcana.Node, 0)
	file := normalizedUnitPath(node.Path)
	for sourceID := range graph.incoming[node.NodeID] {
		source := graph.nodes[sourceID]
		if normalizedUnitPath(source.Path) == file {
			callers = append(callers, source)
		}
	}
	return callers
}

func establishedWrapperSupport(graph symbolCallGraph, caller arcana.Node, metrics symbolWrapperFamilyMetrics) (arcana.Node, bool) {
	matches := make([]arcana.Node, 0)
	for targetID := range graph.outgoing[caller.NodeID] {
		if metrics.support[targetID] >= minimumSymbolWrapperPeers {
			matches = append(matches, graph.nodes[targetID])
		}
	}
	if len(matches) == 0 {
		return arcana.Node{}, false
	}
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Path != matches[j].Path {
			return matches[i].Path < matches[j].Path
		}
		if matches[i].Name != matches[j].Name {
			return matches[i].Name < matches[j].Name
		}
		return matches[i].NodeID < matches[j].NodeID
	})
	return matches[0], true
}

func wrapperFamilyKey(wrapper symbolWrapper) string {
	return normalizedUnitPath(wrapper.intermediary.Path) + "\x00" + normalizedUnitPath(wrapper.target.Path)
}

func isCallableSymbol(node arcana.Node) bool {
	switch strings.ToLower(strings.TrimSpace(node.Kind)) {
	case "function", "method", "constructor", "procedure", "symbol":
		return true
	default:
		return false
	}
}
