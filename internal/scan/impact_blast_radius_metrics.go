package scan

import (
	"math"
	"math/bits"
)

func transitiveDependents(units []dependencyUnit) [][]uint64 {
	wordCount := (len(units) + 63) / 64
	indexByPath := make(map[string]int, len(units))
	for index, unit := range units {
		indexByPath[unit.path] = index
	}
	reach := make([][]uint64, len(units))
	for index := range reach {
		reach[index] = make([]uint64, wordCount)
	}
	queue := make([]int, 0, len(units))
	queued := make([]bool, len(units))
	for target, unit := range units {
		for sourcePath := range unit.incoming {
			source := indexByPath[sourcePath]
			bitSetAdd(reach[target], source)
		}
		if bitSetCount(reach[target]) > 0 {
			queue = append(queue, target)
			queued[target] = true
		}
	}
	for len(queue) > 0 {
		source := queue[0]
		queue = queue[1:]
		queued[source] = false
		for targetPath := range units[source].outgoing {
			target := indexByPath[targetPath]
			changed := bitSetAdd(reach[target], source)
			changed = bitSetUnionExcept(reach[target], reach[source], target) || changed
			if changed && !queued[target] {
				queue = append(queue, target)
				queued[target] = true
			}
		}
	}
	return reach
}

func independentImpactBranches(
	candidate int,
	unit dependencyUnit,
	transitive [][]uint64,
	units []dependencyUnit,
	indexByPath map[string]int,
) (int, float64) {
	if len(unit.incoming) == 0 {
		return 0, 1
	}
	branches := make([][]uint64, 0, len(unit.incoming))
	membership := make([]int, len(units))
	largest := 0
	for sourcePath := range unit.incoming {
		source := indexByPath[sourcePath]
		branch := append([]uint64(nil), transitive[source]...)
		bitSetAdd(branch, source)
		branch[candidate/64] &^= uint64(1) << uint(candidate%64)
		branches = append(branches, branch)
		if size := bitSetCount(branch); size > largest {
			largest = size
		}
		for index := range units {
			if bitSetHas(branch, index) {
				membership[index]++
			}
		}
	}
	reach := max(1, bitSetCount(transitive[candidate]))
	threshold := max(2, int(math.Ceil(float64(reach)*minimumUniqueImpactBranchFraction)))
	independent := 0
	for _, branch := range branches {
		unique := 0
		for index := range units {
			if bitSetHas(branch, index) && membership[index] == 1 {
				unique++
			}
		}
		if unique >= threshold {
			independent++
		}
	}
	return independent, float64(largest) / float64(reach)
}

func impactedRegionCount(dependents []uint64, units []dependencyUnit, semantic map[string]dependencyRegion) int {
	regions := make(map[string]struct{})
	for index, unit := range units {
		if bitSetHas(dependents, index) {
			regions[dependencyRegionForFile(unit.path, semantic).id] = struct{}{}
		}
	}
	return len(regions)
}

func bitSetAdd(set []uint64, index int) bool {
	word, mask := index/64, uint64(1)<<uint(index%64)
	before := set[word]
	set[word] |= mask
	return before != set[word]
}

func bitSetHas(set []uint64, index int) bool {
	return set[index/64]&(uint64(1)<<uint(index%64)) != 0
}

func bitSetUnionExcept(target, source []uint64, excluded int) bool {
	changed := false
	excludedWord := excluded / 64
	excludedMask := uint64(1) << uint(excluded%64)
	for index := range target {
		addition := source[index]
		if index == excludedWord {
			addition &^= excludedMask
		}
		before := target[index]
		target[index] |= addition
		changed = changed || before != target[index]
	}
	return changed
}

func bitSetCount(set []uint64) int {
	count := 0
	for _, word := range set {
		count += bits.OnesCount64(word)
	}
	return count
}
