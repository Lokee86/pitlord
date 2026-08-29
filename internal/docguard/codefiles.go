package docguard

import (
	"context"
	"fmt"
	"sort"

	"github.com/Lokee86/pitlord/internal/arcana"
	"github.com/Lokee86/pitlord/internal/snapshot"
)

func LoadCodeFiles(ctx context.Context, repository, snapshotPath, arcanaCommand string) ([]string, error) {
	resolvedSnapshot, err := snapshot.Resolve(repository, snapshotPath)
	if err != nil {
		return nil, err
	}
	resolvedCommand, err := arcana.ResolveCommand(repository, arcanaCommand)
	if err != nil {
		return nil, err
	}
	graph, err := (arcana.Client{Command: resolvedCommand}).LoadGraphWithOptions(
		ctx,
		resolvedSnapshot,
		arcana.LoadOptions{SourcePrefixes: []string{"."}},
	)
	if err != nil {
		return nil, fmt.Errorf("load Arcana code files: %w", err)
	}
	seen := make(map[string]struct{})
	for _, node := range graph.Sources {
		file := normalizePath(node.Path)
		if node.Kind == "file" && file != "" {
			seen[file] = struct{}{}
		}
	}
	files := make([]string, 0, len(seen))
	for file := range seen {
		files = append(files, file)
	}
	sort.Strings(files)
	return files, nil
}
