package arcana

import (
	"context"
	"encoding/json"
	"fmt"
)

func loadNodesForPrefix(
	ctx context.Context,
	run protocolRun,
	command string,
	snapshot string,
	pathPrefix string,
	prefixIndex int,
) ([]Node, error) {
	var nodes []Node
	offset := 0
	expectedCount := -1
	for page := 0; ; page++ {
		current := request{
			ID:         fmt.Sprintf("list-%d-%d", prefixIndex, page),
			Op:         "list_nodes",
			PathPrefix: pathPrefix,
			Offset:     offset,
			Limit:      maxResultLimit,
		}
		responses, err := run(ctx, command, snapshot, []request{current})
		if err != nil {
			return nil, err
		}
		var result nodeListResult
		if err := json.Unmarshal(responses[current.ID].Result, &result); err != nil {
			return nil, fmt.Errorf("decode Arcana node list %q: %w", current.ID, err)
		}
		if result.Offset != offset {
			return nil, fmt.Errorf(
				"Arcana node list for %q returned offset %d for requested offset %d",
				pathPrefix,
				result.Offset,
				offset,
			)
		}
		if result.Returned != len(result.Nodes) {
			return nil, fmt.Errorf(
				"Arcana node list for %q reported %d nodes but returned %d",
				pathPrefix,
				result.Returned,
				len(result.Nodes),
			)
		}
		if expectedCount < 0 {
			expectedCount = result.Count
		} else if result.Count != expectedCount {
			return nil, fmt.Errorf(
				"Arcana node list for %q changed count from %d to %d during pagination",
				pathPrefix,
				expectedCount,
				result.Count,
			)
		}
		if result.Count < 0 || result.Offset+result.Returned > result.Count {
			return nil, fmt.Errorf("Arcana node list for %q returned invalid page bounds", pathPrefix)
		}
		nodes = append(nodes, result.Nodes...)
		if !result.Truncated {
			if result.NextOffset != nil {
				return nil, fmt.Errorf("Arcana node list for %q returned a terminal next_offset", pathPrefix)
			}
			if result.Offset+result.Returned != result.Count {
				return nil, fmt.Errorf(
					"Arcana node list for %q ended at %d of %d nodes",
					pathPrefix,
					result.Offset+result.Returned,
					result.Count,
				)
			}
			return nodes, nil
		}
		if result.NextOffset == nil {
			return nil, fmt.Errorf(
				"Arcana node list for %q was truncated without next_offset; upgrade Arcana",
				pathPrefix,
			)
		}
		if *result.NextOffset <= offset || *result.NextOffset > result.Count {
			return nil, fmt.Errorf(
				"Arcana node list for %q returned invalid next_offset %d",
				pathPrefix,
				*result.NextOffset,
			)
		}
		offset = *result.NextOffset
	}
}
