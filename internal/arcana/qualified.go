package arcana

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

func (client Client) ResolveQualified(
	ctx context.Context,
	snapshot string,
	qualifiedName string,
) ([]Node, error) {
	path, name, err := splitQualifiedName(qualifiedName)
	if err != nil {
		return nil, err
	}
	run := client.Run
	if run == nil {
		run = runProtocol
	}
	command := client.Command
	if command == "" {
		command = "arcana"
	}
	current := request{
		ID:    "resolve-qualified",
		Op:    "resolve_symbol",
		Name:  name,
		Path:  path,
		Limit: maxResultLimit,
	}
	responses, err := run(ctx, command, snapshot, []request{current})
	if err != nil {
		return nil, err
	}
	var result nodeListResult
	if err := json.Unmarshal(responses[current.ID].Result, &result); err != nil {
		return nil, fmt.Errorf("decode Arcana qualified-symbol result: %w", err)
	}
	if result.Truncated || result.Returned != result.Count || result.Returned != len(result.Nodes) {
		return nil, fmt.Errorf(
			"Arcana qualified-symbol result for %q was incomplete: %d of %d nodes",
			qualifiedName,
			result.Returned,
			result.Count,
		)
	}
	return result.Nodes, nil
}

func (client Client) HasQualifiedRelationship(
	ctx context.Context,
	snapshot string,
	sourceQualifiedName string,
	relation string,
	targetQualifiedName string,
) (bool, error) {
	sources, err := client.ResolveQualified(ctx, snapshot, sourceQualifiedName)
	if err != nil {
		return false, err
	}
	if len(sources) == 0 {
		return false, nil
	}
	run := client.Run
	if run == nil {
		run = runProtocol
	}
	command := client.Command
	if command == "" {
		command = "arcana"
	}
	requests := make([]request, 0, len(sources))
	for _, source := range sources {
		requests = append(requests, request{
			ID:        fmt.Sprintf("qualified-neighbors-%d", source.NodeID),
			Op:        "neighbors",
			NodeID:    source.NodeID,
			Direction: "outgoing",
			Relation:  relation,
		})
	}
	responses, err := run(ctx, command, snapshot, requests)
	if err != nil {
		return false, err
	}
	for _, current := range requests {
		var result neighborResult
		if err := json.Unmarshal(responses[current.ID].Result, &result); err != nil {
			return false, fmt.Errorf("decode Arcana qualified neighbors %q: %w", current.ID, err)
		}
		if result.Count != len(result.Relationships) {
			return false, fmt.Errorf(
				"Arcana qualified neighbors %q reported %d relationships but returned %d",
				current.ID,
				result.Count,
				len(result.Relationships),
			)
		}
		for _, relationship := range result.Relationships {
			if qualifiedNodeName(relationship.Node) == targetQualifiedName {
				return true, nil
			}
		}
	}
	return false, nil
}

func splitQualifiedName(value string) (string, string, error) {
	value = strings.TrimSpace(value)
	separator := strings.LastIndex(value, "::")
	if separator <= 0 || separator+2 >= len(value) {
		return "", "", fmt.Errorf("qualified symbol %q must use path::name", value)
	}
	path := strings.TrimSpace(value[:separator])
	name := strings.TrimSpace(value[separator+2:])
	if path == "" || name == "" {
		return "", "", fmt.Errorf("qualified symbol %q must use path::name", value)
	}
	return path, name, nil
}

func qualifiedNodeName(node Node) string {
	return node.Path + "::" + node.Name
}
