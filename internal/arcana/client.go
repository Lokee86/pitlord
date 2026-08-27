package arcana

import (
	"context"
	"fmt"
)

const maxResultLimit = 10_000

type Client struct {
	Command string
	Run     protocolRun
}

type LoadOptions struct {
	SourcePrefixes   []string
	OutgoingPrefixes []string
}

type nodeListResult struct {
	Count      int    `json:"count"`
	Offset     int    `json:"offset"`
	Returned   int    `json:"returned"`
	Truncated  bool   `json:"truncated"`
	NextOffset *int   `json:"next_offset"`
	Nodes      []Node `json:"nodes"`
}

type neighborResult struct {
	Node          Node           `json:"node"`
	Count         int            `json:"count"`
	Returned      int            `json:"returned"`
	Truncated     bool           `json:"truncated"`
	Relationships []Relationship `json:"relationships"`
}

func validateCompleteNeighbors(label string, result neighborResult) error {
	if result.Returned != len(result.Relationships) {
		return fmt.Errorf(
			"%s reported %d returned relationships but contained %d",
			label,
			result.Returned,
			len(result.Relationships),
		)
	}
	if result.Truncated {
		return fmt.Errorf(
			"%s was truncated at %d relationships; Arcana's protocol limit is %d",
			label,
			result.Returned,
			maxResultLimit,
		)
	}
	if result.Count != result.Returned {
		return fmt.Errorf(
			"%s reported %d relationships but returned %d",
			label,
			result.Count,
			result.Returned,
		)
	}
	return nil
}

func (client Client) LoadGraph(
	ctx context.Context,
	snapshot string,
	pathPrefixes []string,
) (Graph, error) {
	return client.LoadGraphWithOptions(ctx, snapshot, LoadOptions{
		SourcePrefixes:   pathPrefixes,
		OutgoingPrefixes: pathPrefixes,
	})
}
