package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/Lokee86/pitlord/internal/policy"
)

func WriteAreaAnalysisJSON(writer io.Writer, analysis policy.AreaAnalysis) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(analysis)
}

func WriteAreaAnalysisText(writer io.Writer, analysis policy.AreaAnalysis) error {
	if _, err := fmt.Fprintf(
		writer,
		"Pitlord area analysis: %d areas, %d cross-area/external dependencies, %d cycle(s).\n",
		analysis.Summary.AreaCount,
		analysis.Summary.DependencyCount,
		analysis.Summary.CycleCount,
	); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(
		writer,
		"Scanned %d nodes and %d relationships; %d relationships matched the selected architecture relations.\n",
		analysis.Summary.SourceNodesScanned,
		analysis.Summary.RelationshipsScanned,
		analysis.Summary.RelationshipsSelected,
	); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(
		writer,
		"Ownership: %d checked, %d singly owned, %d unowned, %d overlapping.\n",
		analysis.Ownership.CheckedNodes,
		analysis.Ownership.OwnedNodes,
		analysis.Ownership.UnownedNodes,
		analysis.Ownership.OverlappingNodes,
	); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(writer, "\nAreas:"); err != nil {
		return err
	}
	for _, area := range analysis.Areas {
		if _, err := fmt.Fprintf(
			writer,
			"  %s: %d nodes, %d internal relationships",
			area.ID,
			area.NodeCount,
			area.InternalRelationships,
		); err != nil {
			return err
		}
		if kinds := formatCounts(area.KindCounts); kinds != "" {
			if _, err := fmt.Fprintf(writer, " (%s)", kinds); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(writer); err != nil {
			return err
		}
	}

	if len(analysis.Dependencies) > 0 {
		if _, err := fmt.Fprintln(writer, "\nDependencies:"); err != nil {
			return err
		}
		for _, dependency := range analysis.Dependencies {
			if _, err := fmt.Fprintf(
				writer,
				"  %s -> %s: %d edges (%s)\n",
				dependency.SourceArea,
				dependency.TargetArea,
				dependency.EdgeCount,
				formatCounts(dependency.RelationCounts),
			); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(
				writer,
				"    representative: %s\n",
				evidenceText(dependency.Representative),
			); err != nil {
				return err
			}
		}
	}

	if len(analysis.Cycles) > 0 {
		if _, err := fmt.Fprintln(writer, "\nCycles:"); err != nil {
			return err
		}
		for _, cycle := range analysis.Cycles {
			if _, err := fmt.Fprintf(writer, "  %s\n", strings.Join(cycle, " <-> ")); err != nil {
				return err
			}
		}
	}

	if len(analysis.Ownership.UnownedExamples) > 0 {
		if _, err := fmt.Fprintln(writer, "\nUnowned examples:"); err != nil {
			return err
		}
		for _, node := range analysis.Ownership.UnownedExamples {
			if _, err := fmt.Fprintf(writer, "  %s %s\n", location(node), displayNode(node)); err != nil {
				return err
			}
		}
	}
	if len(analysis.Ownership.OverlapExamples) > 0 {
		if _, err := fmt.Fprintln(writer, "\nOverlap examples:"); err != nil {
			return err
		}
		for _, overlap := range analysis.Ownership.OverlapExamples {
			if _, err := fmt.Fprintf(
				writer,
				"  %s %s: %s\n",
				location(overlap.Node),
				displayNode(overlap.Node),
				strings.Join(overlap.Areas, ", "),
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func WriteAreaAnalysisDOT(writer io.Writer, analysis policy.AreaAnalysis) error {
	if _, err := fmt.Fprintln(writer, "digraph pitlord_architecture {"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(writer, "  rankdir=LR;"); err != nil {
		return err
	}
	for _, area := range analysis.Areas {
		label := area.ID + "\\n" + strconv.Itoa(area.NodeCount) + " nodes"
		if _, err := fmt.Fprintf(
			writer,
			"  %s [label=\"%s\"];\n",
			dotID(area.ID),
			dotEscape(label),
		); err != nil {
			return err
		}
	}
	hasExternal := false
	for _, dependency := range analysis.Dependencies {
		if dependency.TargetArea == policy.ExternalArea {
			hasExternal = true
			break
		}
	}
	if hasExternal {
		if _, err := fmt.Fprintln(writer, "  external [label=\"external\", shape=box, style=dashed];"); err != nil {
			return err
		}
	}
	cycleEdges := cycleEdgeSet(analysis.Cycles)
	for _, dependency := range analysis.Dependencies {
		targetID := dotID(dependency.TargetArea)
		if dependency.TargetArea == policy.ExternalArea {
			targetID = "external"
		}
		attributes := []string{
			"label=\"" + dotEscape(formatCounts(dependency.RelationCounts)) + "\"",
			"penwidth=" + strconv.Itoa(dotPenWidth(dependency.EdgeCount)),
		}
		if _, cyclic := cycleEdges[dependency.SourceArea+"\x00"+dependency.TargetArea]; cyclic {
			attributes = append(attributes, "color=red")
		}
		if _, err := fmt.Fprintf(
			writer,
			"  %s -> %s [%s];\n",
			dotID(dependency.SourceArea),
			targetID,
			strings.Join(attributes, ", "),
		); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(writer, "}")
	return err
}

func cycleEdgeSet(cycles [][]string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, cycle := range cycles {
		members := make(map[string]struct{}, len(cycle))
		for _, area := range cycle {
			members[area] = struct{}{}
		}
		for source := range members {
			for target := range members {
				if source != target {
					result[source+"\x00"+target] = struct{}{}
				}
			}
		}
	}
	return result
}

func dotID(value string) string {
	return strconv.QuoteToASCII(value)
}

func dotEscape(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	value = strings.ReplaceAll(value, "\n", "\\n")
	return value
}

func dotPenWidth(edges int) int {
	switch {
	case edges >= 100:
		return 5
	case edges >= 25:
		return 4
	case edges >= 5:
		return 3
	case edges >= 2:
		return 2
	default:
		return 1
	}
}

func sortedCountKeys(counts map[string]int) []string {
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
