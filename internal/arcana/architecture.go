package arcana

import (
	"context"
	"encoding/json"
	"fmt"
)

type ArchitectureSummary struct {
	PathPrefix        string                  `json:"path_prefix,omitempty"`
	Relations         []string                `json:"relations"`
	NodeCount         int                     `json:"node_count"`
	InternalEdgeCount int                     `json:"internal_edge_count"`
	BoundaryEdgeCount int                     `json:"boundary_edge_count"`
	ComponentCount    int                     `json:"component_count"`
	CommunityCount    int                     `json:"community_count"`
	Returned          int                     `json:"returned"`
	Truncated         bool                    `json:"truncated"`
	MinCommunitySize  int                     `json:"min_community_size"`
	ExcludedNodeCount int                     `json:"excluded_node_count"`
	Communities       []ArchitectureCommunity `json:"communities"`
}

type ArchitectureCommunity struct {
	CommunityID         uint32          `json:"community_id"`
	NodeCount           int             `json:"node_count"`
	EdgeCount           int             `json:"edge_count"`
	RelationCounts      map[string]int  `json:"relation_counts"`
	IncomingBoundary    map[string]int  `json:"incoming_boundary_counts"`
	OutgoingBoundary    map[string]int  `json:"outgoing_boundary_counts"`
	KindCounts          map[string]int  `json:"kind_counts"`
	Paths               []CommunityPath `json:"paths"`
	RepresentativeNodes []Node          `json:"representative_nodes"`
}

type CommunityPath struct {
	Path      string `json:"path"`
	NodeCount int    `json:"node_count"`
}

func (client Client) InspectArchitecture(
	ctx context.Context,
	snapshot string,
	pathPrefix string,
	relations []string,
	minCommunitySize int,
	limit int,
) (ArchitectureSummary, error) {
	if snapshot == "" {
		return ArchitectureSummary{}, fmt.Errorf("Arcana snapshot is required")
	}
	run := client.Run
	if run == nil {
		run = runProtocol
	}
	command := client.Command
	if command == "" {
		command = "arcana"
	}
	if pathPrefix == "." {
		pathPrefix = ""
	}
	if minCommunitySize <= 0 {
		minCommunitySize = 2
	}
	if limit <= 0 {
		limit = 20
	}
	current := request{
		ID:               "architecture",
		Op:               "architecture_summary",
		PathPrefix:       pathPrefix,
		Relations:        relations,
		MinCommunitySize: minCommunitySize,
		Limit:            limit,
	}
	responses, err := run(ctx, command, snapshot, []request{current})
	if err != nil {
		return ArchitectureSummary{}, err
	}
	var summary ArchitectureSummary
	if err := json.Unmarshal(responses[current.ID].Result, &summary); err != nil {
		return ArchitectureSummary{}, fmt.Errorf("decode Arcana architecture summary: %w", err)
	}
	if summary.Returned != len(summary.Communities) {
		return ArchitectureSummary{}, fmt.Errorf(
			"Arcana architecture summary reported %d communities but returned %d",
			summary.Returned,
			len(summary.Communities),
		)
	}
	return summary, nil
}
