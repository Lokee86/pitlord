package scan

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

type Input struct {
	RepositoryRoot string
	SnapshotPath   string
	PathPrefix     string
}

type GraphLoader interface {
	LoadGraphWithOptions(context.Context, string, arcana.LoadOptions) (arcana.Graph, error)
}

type Engine struct {
	Loader    GraphLoader
	Analyzers []Analyzer
}

func (engine Engine) Run(ctx context.Context, input Input) (Result, error) {
	pathPrefix, err := normalizePathPrefix(input.PathPrefix)
	if err != nil {
		return Result{}, err
	}
	analyzers := engine.Analyzers
	if analyzers == nil {
		analyzers = defaultAnalyzers()
	}
	analyzerInput := AnalyzerContext{
		RepositoryRoot: input.RepositoryRoot,
		SnapshotPath:   input.SnapshotPath,
		PathPrefix:     pathPrefix,
	}
	if analyzersRequireGraph(analyzers) {
		if strings.TrimSpace(input.SnapshotPath) == "" {
			return Result{}, fmt.Errorf("Arcana snapshot is required")
		}
		if engine.Loader == nil {
			return Result{}, fmt.Errorf("Arcana graph loader is required")
		}
		graph, loadErr := engine.Loader.LoadGraphWithOptions(ctx, input.SnapshotPath, arcana.LoadOptions{
			SourcePrefixes:   []string{pathPrefix},
			OutgoingPrefixes: []string{pathPrefix},
		})
		if loadErr != nil {
			return Result{}, loadErr
		}
		analyzerInput.Graph = graph
	}
	return analyzeWithAnalyzers(ctx, analyzerInput, analyzers)
}

func analyzeGraph(pathPrefix string, graph arcana.Graph) Result {
	result, err := analyzeWithAnalyzers(context.Background(), AnalyzerContext{
		PathPrefix: pathPrefix,
		Graph:      graph,
	}, defaultAnalyzers())
	if err != nil {
		panic(err)
	}
	return result
}

func analyzeWithAnalyzers(ctx context.Context, input AnalyzerContext, analyzers []Analyzer) (Result, error) {
	findings, err := runAnalyzers(ctx, input, analyzers)
	if err != nil {
		return Result{}, err
	}
	return finalize(Result{
		Schema: Schema,
		Scope: Scope{
			Kind: "repository",
			Path: input.PathPrefix,
		},
		Findings: findings,
	}), nil
}

func normalizePathPrefix(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" || value == "." {
		return ".", nil
	}
	if filepath.IsAbs(value) {
		return "", fmt.Errorf("scan path prefix must be repository-relative")
	}
	value = filepath.ToSlash(value)
	cleaned := path.Clean(value)
	if cleaned == "." {
		return ".", nil
	}
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return "", fmt.Errorf("scan path prefix must not escape the repository")
	}
	return strings.TrimPrefix(cleaned, "./"), nil
}
