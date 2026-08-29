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
	SnapshotPath string
	PathPrefix   string
}

type GraphLoader interface {
	LoadGraphWithOptions(context.Context, string, arcana.LoadOptions) (arcana.Graph, error)
}

type Engine struct {
	Loader GraphLoader
}

func (engine Engine) Run(ctx context.Context, input Input) (Result, error) {
	if strings.TrimSpace(input.SnapshotPath) == "" {
		return Result{}, fmt.Errorf("Arcana snapshot is required")
	}
	pathPrefix, err := normalizePathPrefix(input.PathPrefix)
	if err != nil {
		return Result{}, err
	}
	if engine.Loader == nil {
		return Result{}, fmt.Errorf("Arcana graph loader is required")
	}
	graph, err := engine.Loader.LoadGraphWithOptions(ctx, input.SnapshotPath, arcana.LoadOptions{
		SourcePrefixes:   []string{pathPrefix},
		OutgoingPrefixes: []string{pathPrefix},
	})
	if err != nil {
		return Result{}, err
	}
	return analyzeGraph(pathPrefix, graph), nil
}

func analyzeGraph(pathPrefix string, graph arcana.Graph) Result {
	findings := detectDependencyPressure(graph, pathPrefix)
	findings = append(findings, detectHubBottlenecks(graph, pathPrefix)...)
	findings = append(findings, detectDependencyKnots(graph, pathPrefix)...)
	findings = append(findings, detectDependencyDepth(graph, pathPrefix)...)
	findings = append(findings, detectBoundaryCohesion(graph, pathPrefix)...)
	findings = append(findings, detectImpactBlastRadius(graph, pathPrefix)...)
	return finalize(Result{
		Schema: Schema,
		Scope: Scope{
			Kind: "repository",
			Path: pathPrefix,
		},
		Findings: findings,
	})
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
