package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

func WriteArchitectureJSON(writer io.Writer, summary arcana.ArchitectureSummary) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(summary)
}

func WriteArchitectureText(writer io.Writer, summary arcana.ArchitectureSummary) error {
	if _, err := fmt.Fprintf(
		writer,
		"Architecture: %d nodes, %d internal edges, %d boundary edges, %d communities (%d returned).\n",
		summary.NodeCount,
		summary.InternalEdgeCount,
		summary.BoundaryEdgeCount,
		summary.CommunityCount,
		summary.Returned,
	); err != nil {
		return err
	}
	if len(summary.Relations) > 0 {
		if _, err := fmt.Fprintf(writer, "Relations: %s\n", strings.Join(summary.Relations, ", ")); err != nil {
			return err
		}
	}
	for index, community := range summary.Communities {
		if _, err := fmt.Fprintf(
			writer,
			"\nCommunity %d: %d nodes, %d edges\n",
			index+1,
			community.NodeCount,
			community.EdgeCount,
		); err != nil {
			return err
		}
		if len(community.Paths) > 0 {
			paths := make([]string, 0, len(community.Paths))
			for _, path := range community.Paths {
				paths = append(paths, fmt.Sprintf("%s (%d)", path.Path, path.NodeCount))
			}
			if _, err := fmt.Fprintf(writer, "  Paths: %s\n", strings.Join(paths, ", ")); err != nil {
				return err
			}
		}
		if len(community.RepresentativeNodes) > 0 {
			nodes := make([]string, 0, len(community.RepresentativeNodes))
			for _, node := range community.RepresentativeNodes {
				nodes = append(nodes, fmt.Sprintf("%s [%s]", displayNode(node), node.Path))
			}
			if _, err := fmt.Fprintf(writer, "  Representatives: %s\n", strings.Join(nodes, ", ")); err != nil {
				return err
			}
		}
		if counts := formatCounts(community.RelationCounts); counts != "" {
			if _, err := fmt.Fprintf(writer, "  Internal: %s\n", counts); err != nil {
				return err
			}
		}
		if counts := formatCounts(community.IncomingBoundary); counts != "" {
			if _, err := fmt.Fprintf(writer, "  Incoming boundary: %s\n", counts); err != nil {
				return err
			}
		}
		if counts := formatCounts(community.OutgoingBoundary); counts != "" {
			if _, err := fmt.Fprintf(writer, "  Outgoing boundary: %s\n", counts); err != nil {
				return err
			}
		}
	}
	if summary.Truncated {
		_, err := fmt.Fprintln(writer, "\nResult truncated; increase --limit to inspect more communities.")
		return err
	}
	return nil
}

func formatCounts(counts map[string]int) string {
	if len(counts) == 0 {
		return ""
	}
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%d", key, counts[key]))
	}
	return strings.Join(parts, ", ")
}
