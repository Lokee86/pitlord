package scan

import (
	"context"
	"fmt"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
)

type AnalyzerContext struct {
	RepositoryRoot string
	SnapshotPath   string
	PathPrefix     string
	Graph          arcana.Graph
}

type AnalyzerMetadata struct {
	ID            string
	Category      string
	Language      string
	RequiresGraph bool
}

type Analyzer interface {
	Metadata() AnalyzerMetadata
	Analyze(context.Context, AnalyzerContext) ([]Finding, error)
}

type graphAnalyzer struct {
	metadata AnalyzerMetadata
	detect   func(arcana.Graph, string) []Finding
}

func (analyzer graphAnalyzer) Metadata() AnalyzerMetadata {
	return analyzer.metadata
}

func (analyzer graphAnalyzer) Analyze(_ context.Context, input AnalyzerContext) ([]Finding, error) {
	return analyzer.detect(input.Graph, input.PathPrefix), nil
}

func defaultAnalyzers() []Analyzer {
	return []Analyzer{
		architectureAnalyzer(DetectorDependencyPressure, detectDependencyPressure),
		architectureAnalyzer(DetectorHubBottleneck, detectHubBottlenecks),
		architectureAnalyzer(DetectorDependencyKnots, detectDependencyKnots),
		architectureAnalyzer(DetectorDependencyDepth, detectDependencyDepth),
		architectureAnalyzer(DetectorUnstableDependencyDirection, detectUnstableDependencyDirection),
		architectureAnalyzer(DetectorSymbolIntermediaryBypass, detectSymbolIntermediaryBypass),
		architectureAnalyzer(DetectorCrossFileIntermediaryBypass, detectCrossFileIntermediaryBypass),
		architectureAnalyzer(DetectorBoundaryBypass, detectBoundaryBypass),
		architectureAnalyzer(DetectorBoundaryCohesion, detectBoundaryCohesion),
		architectureAnalyzer(DetectorImpactBlastRadius, detectImpactBlastRadius),
	}
}

func architectureAnalyzer(id string, detect func(arcana.Graph, string) []Finding) Analyzer {
	return graphAnalyzer{
		metadata: AnalyzerMetadata{ID: id, Category: "architecture", RequiresGraph: true},
		detect:   detect,
	}
}

func analyzersRequireGraph(analyzers []Analyzer) bool {
	for _, analyzer := range analyzers {
		if analyzer != nil && analyzer.Metadata().RequiresGraph {
			return true
		}
	}
	return false
}

func runAnalyzers(ctx context.Context, input AnalyzerContext, analyzers []Analyzer) ([]Finding, error) {
	findings := make([]Finding, 0)
	seen := make(map[string]struct{}, len(analyzers))
	for _, analyzer := range analyzers {
		if analyzer == nil {
			return nil, fmt.Errorf("scan analyzer is required")
		}
		metadata := analyzer.Metadata()
		metadata.ID = strings.TrimSpace(metadata.ID)
		if metadata.ID == "" {
			return nil, fmt.Errorf("scan analyzer ID is required")
		}
		if _, exists := seen[metadata.ID]; exists {
			return nil, fmt.Errorf("duplicate scan analyzer %q", metadata.ID)
		}
		seen[metadata.ID] = struct{}{}
		analyzerFindings, err := analyzer.Analyze(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("scan analyzer %q: %w", metadata.ID, err)
		}
		for index := range analyzerFindings {
			normalizeAnalyzerFinding(&analyzerFindings[index], metadata)
		}
		findings = append(findings, analyzerFindings...)
	}
	return findings, nil
}

func normalizeAnalyzerFinding(finding *Finding, metadata AnalyzerMetadata) {
	if finding.Analyzer == "" {
		finding.Analyzer = metadata.ID
	}
	if finding.Category == "" {
		finding.Category = metadata.Category
	}
	if finding.Language == "" {
		finding.Language = metadata.Language
	}
	if finding.RuleID == "" {
		finding.RuleID = finding.Detector
		if finding.RuleID == "" {
			finding.RuleID = metadata.ID
		}
	}
	if finding.Detector == "" {
		finding.Detector = metadata.ID
	}
}
