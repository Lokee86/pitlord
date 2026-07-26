package arcana

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

func (client Client) LoadGraphWithOptions(
	ctx context.Context,
	snapshot string,
	options LoadOptions,
) (Graph, error) {
	if snapshot == "" {
		return Graph{}, fmt.Errorf("Arcana snapshot is required")
	}
	if len(options.SourcePrefixes) == 0 {
		return Graph{}, fmt.Errorf("at least one source path prefix is required")
	}
	run := client.Run
	if run == nil {
		run = runProtocol
	}
	command := client.Command
	if command == "" {
		command = "arcana"
	}

	sources, err := loadSourceNodes(
		ctx,
		run,
		command,
		snapshot,
		options.SourcePrefixes,
	)
	if err != nil {
		return Graph{}, err
	}
	graph := Graph{
		Sources:  sources,
		Outgoing: make(map[uint32][]Relationship, len(sources)),
	}
	outgoingSources := selectOutgoingSources(sources, options.OutgoingPrefixes)
	if len(outgoingSources) == 0 {
		return graph, nil
	}
	if err := loadOutgoingRelationships(
		ctx,
		run,
		command,
		snapshot,
		outgoingSources,
		&graph,
	); err != nil {
		return Graph{}, err
	}
	return graph, nil
}

func loadSourceNodes(
	ctx context.Context,
	run protocolRun,
	command string,
	snapshot string,
	pathPrefixes []string,
) ([]Node, error) {
	byID := make(map[uint32]Node)
	for index, prefix := range pathPrefixes {
		queryPrefix := prefix
		if queryPrefix == "." {
			queryPrefix = ""
		}
		nodes, err := loadNodesForPrefix(
			ctx,
			run,
			command,
			snapshot,
			queryPrefix,
			index,
		)
		if err != nil {
			return nil, err
		}
		for _, node := range nodes {
			byID[node.NodeID] = node
		}
	}

	sources := make([]Node, 0, len(byID))
	for _, node := range byID {
		sources = append(sources, node)
	}
	sort.Slice(sources, func(i, j int) bool {
		if sources[i].Path != sources[j].Path {
			return sources[i].Path < sources[j].Path
		}
		if sources[i].Name != sources[j].Name {
			return sources[i].Name < sources[j].Name
		}
		return sources[i].NodeID < sources[j].NodeID
	})
	return sources, nil
}

func selectOutgoingSources(sources []Node, prefixes []string) []Node {
	if len(prefixes) == 0 {
		return nil
	}
	selected := make([]Node, 0)
	for _, source := range sources {
		if matchesPathPrefixes(source.Path, prefixes) {
			selected = append(selected, source)
		}
	}
	return selected
}

func loadOutgoingRelationships(
	ctx context.Context,
	run protocolRun,
	command string,
	snapshot string,
	sources []Node,
	graph *Graph,
) error {
	requests := make([]request, 0, len(sources))
	for _, node := range sources {
		requests = append(requests, request{
			ID:        fmt.Sprintf("neighbors-%d", node.NodeID),
			Op:        "neighbors",
			NodeID:    node.NodeID,
			Direction: "outgoing",
		})
	}
	responses, err := run(ctx, command, snapshot, requests)
	if err != nil {
		return err
	}
	for _, current := range requests {
		var result neighborResult
		if err := json.Unmarshal(responses[current.ID].Result, &result); err != nil {
			return fmt.Errorf("decode Arcana neighbors %q: %w", current.ID, err)
		}
		if result.Count != len(result.Relationships) {
			return fmt.Errorf(
				"Arcana neighbors %q reported %d relationships but returned %d",
				current.ID,
				result.Count,
				len(result.Relationships),
			)
		}
		graph.Outgoing[result.Node.NodeID] = result.Relationships
		graph.Relationships += len(result.Relationships)
	}
	return nil
}

func matchesPathPrefixes(path string, prefixes []string) bool {
	path = normalizeGraphPath(path)
	for _, prefix := range prefixes {
		prefix = normalizeGraphPath(prefix)
		if prefix == "." || path == prefix || strings.HasPrefix(path, prefix+"/") {
			return true
		}
	}
	return false
}

func normalizeGraphPath(path string) string {
	path = strings.ReplaceAll(strings.TrimSpace(path), "\\", "/")
	if path == "." || path == "./" || path == "/" {
		return "."
	}
	path = strings.TrimPrefix(path, "./")
	return strings.Trim(path, "/")
}
