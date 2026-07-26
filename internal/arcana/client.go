package arcana

import "context"

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
	Relationships []Relationship `json:"relationships"`
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
