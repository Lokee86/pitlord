package scan

import "sort"

func stronglyConnectedDependencyComponents(units []dependencyUnit) [][]string {
	adjacency := make(map[string][]string, len(units))
	reverse := make(map[string][]string, len(units))
	paths := make([]string, 0, len(units))
	for _, unit := range units {
		paths = append(paths, unit.path)
		adjacency[unit.path] = sortedSetValues(unit.outgoing)
		reverse[unit.path] = sortedSetValues(unit.incoming)
	}
	sort.Strings(paths)

	visited := make(map[string]struct{}, len(units))
	order := make([]string, 0, len(units))
	for _, start := range paths {
		if _, ok := visited[start]; ok {
			continue
		}
		order = append(order, dependencyDFSPostorder(start, adjacency, visited)...)
	}

	visited = make(map[string]struct{}, len(units))
	components := make([][]string, 0)
	for index := len(order) - 1; index >= 0; index-- {
		start := order[index]
		if _, ok := visited[start]; ok {
			continue
		}
		component := dependencyDFSCollect(start, reverse, visited)
		sort.Strings(component)
		components = append(components, component)
	}
	sort.Slice(components, func(i, j int) bool {
		return components[i][0] < components[j][0]
	})
	return components
}

type dependencyDFSFrame struct {
	node      string
	neighbors []string
	next      int
}

func dependencyDFSPostorder(start string, adjacency map[string][]string, visited map[string]struct{}) []string {
	visited[start] = struct{}{}
	stack := []dependencyDFSFrame{{node: start, neighbors: adjacency[start]}}
	order := make([]string, 0)
	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		if top.next < len(top.neighbors) {
			neighbor := top.neighbors[top.next]
			top.next++
			if _, ok := visited[neighbor]; ok {
				continue
			}
			visited[neighbor] = struct{}{}
			stack = append(stack, dependencyDFSFrame{node: neighbor, neighbors: adjacency[neighbor]})
			continue
		}
		order = append(order, top.node)
		stack = stack[:len(stack)-1]
	}
	return order
}

func dependencyDFSCollect(start string, adjacency map[string][]string, visited map[string]struct{}) []string {
	visited[start] = struct{}{}
	stack := []string{start}
	component := make([]string, 0)
	for len(stack) > 0 {
		last := len(stack) - 1
		node := stack[last]
		stack = stack[:last]
		component = append(component, node)
		neighbors := adjacency[node]
		for index := len(neighbors) - 1; index >= 0; index-- {
			neighbor := neighbors[index]
			if _, ok := visited[neighbor]; ok {
				continue
			}
			visited[neighbor] = struct{}{}
			stack = append(stack, neighbor)
		}
	}
	return component
}

func sortedSetValues(values map[string]struct{}) []string {
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
