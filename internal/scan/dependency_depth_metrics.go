package scan

import "sort"

type dependencyDepthProfile struct {
	next        string
	depth       int
	secondDepth int
	outgoing    int
	dominant    bool
}

type dependencyDepthWalk struct {
	paths          []string
	singleSteps    int
	dominatedSteps int
	minimumGap     int
}

type dependencyDepthBranch struct {
	path  string
	depth int
}

func cyclicDependencyPaths(units []dependencyUnit) map[string]struct{} {
	cyclic := make(map[string]struct{})
	for _, component := range stronglyConnectedDependencyComponents(units) {
		if len(component) < 2 {
			continue
		}
		for _, member := range component {
			cyclic[member] = struct{}{}
		}
	}
	return cyclic
}

func dependencyDepthProfiles(units []dependencyUnit, cyclic map[string]struct{}) map[string]dependencyDepthProfile {
	byPath := dependencyUnitsByPath(units)
	outgoing := make(map[string][]string, len(units))
	indegree := make(map[string]int, len(units))
	for _, unit := range units {
		if _, blocked := cyclic[unit.path]; blocked {
			continue
		}
		indegree[unit.path] = 0
	}
	for _, unit := range units {
		if _, blocked := cyclic[unit.path]; blocked {
			continue
		}
		if isHighlyReusedCentralHub(unit) {
			continue
		}
		for target := range unit.outgoing {
			if _, blocked := cyclic[target]; blocked {
				continue
			}
			if _, exists := byPath[target]; !exists {
				continue
			}
			outgoing[unit.path] = append(outgoing[unit.path], target)
			indegree[target]++
		}
		sort.Strings(outgoing[unit.path])
	}
	queue := make([]string, 0)
	for _, unit := range units {
		if _, blocked := cyclic[unit.path]; blocked {
			continue
		}
		if indegree[unit.path] == 0 {
			queue = append(queue, unit.path)
		}
	}
	order := make([]string, 0, len(indegree))
	for index := 0; index < len(queue); index++ {
		current := queue[index]
		order = append(order, current)
		for _, target := range outgoing[current] {
			indegree[target]--
			if indegree[target] == 0 {
				queue = append(queue, target)
			}
		}
	}
	profiles := make(map[string]dependencyDepthProfile, len(order))
	for index := len(order) - 1; index >= 0; index-- {
		path := order[index]
		branches := make([]dependencyDepthBranch, 0, len(outgoing[path]))
		for _, target := range outgoing[path] {
			branches = append(branches, dependencyDepthBranch{path: target, depth: profiles[target].depth + 1})
		}
		sort.Slice(branches, func(i, j int) bool {
			if branches[i].depth != branches[j].depth {
				return branches[i].depth > branches[j].depth
			}
			return branches[i].path < branches[j].path
		})
		profile := dependencyDepthProfile{outgoing: len(branches)}
		if len(branches) > 0 {
			profile.next = branches[0].path
			profile.depth = branches[0].depth
		}
		if len(branches) > 1 {
			profile.secondDepth = branches[1].depth
		}
		profile.dominant = len(branches) == 1 || len(branches) > 1 && profile.depth-profile.secondDepth >= minimumDominantDepthGap
		profiles[path] = profile
	}
	return profiles
}

func eligibleDepthIncoming(path string, units map[string]dependencyUnit, cyclic map[string]struct{}) int {
	unit, ok := units[path]
	if !ok {
		return 0
	}
	count := 0
	for predecessor := range unit.incoming {
		if _, blocked := cyclic[predecessor]; blocked {
			continue
		}
		if _, exists := units[predecessor]; exists {
			count++
		}
	}
	return count
}

func hasDominantPredecessor(path string, units map[string]dependencyUnit, profiles map[string]dependencyDepthProfile, cyclic map[string]struct{}) bool {
	if eligibleDepthIncoming(path, units, cyclic) > 1 {
		return false
	}
	unit, ok := units[path]
	if !ok {
		return false
	}
	for predecessor := range unit.incoming {
		if _, blocked := cyclic[predecessor]; blocked {
			continue
		}
		profile := profiles[predecessor]
		if profile.dominant && profile.next == path {
			return true
		}
	}
	return false
}

func followDominantDependencyChain(start string, units map[string]dependencyUnit, profiles map[string]dependencyDepthProfile, cyclic map[string]struct{}) dependencyDepthWalk {
	walk := dependencyDepthWalk{minimumGap: int(^uint(0) >> 1)}
	seen := make(map[string]struct{})
	current := start
	for current != "" {
		if _, duplicate := seen[current]; duplicate {
			break
		}
		seen[current] = struct{}{}
		walk.paths = append(walk.paths, current)
		profile := profiles[current]
		if !profile.dominant || profile.next == "" {
			break
		}
		if profile.outgoing == 1 {
			walk.singleSteps++
		} else {
			walk.dominatedSteps++
			gap := profile.depth - profile.secondDepth
			if gap < walk.minimumGap {
				walk.minimumGap = gap
			}
		}
		if eligibleDepthIncoming(profile.next, units, cyclic) > 1 {
			walk.paths = append(walk.paths, profile.next)
			break
		}
		current = profile.next
	}
	if walk.dominatedSteps == 0 {
		walk.minimumGap = 0
	}
	return walk
}
