package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/Lokee86/pitlord/internal/arcana"
	"github.com/Lokee86/pitlord/internal/policy"
	"github.com/Lokee86/pitlord/internal/report"
	"github.com/Lokee86/pitlord/internal/snapshot"
)

func runAnalyze(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("analyze", flag.ContinueOnError)
	flags.SetOutput(stderr)
	repo := flags.String("repo", ".", "repository root containing .arcana/CURRENT")
	explicitSnapshot := flags.String("snapshot", "", "explicit Arcana snapshot directory")
	policyPath := flags.String("policy", "", "Pitlord JSON policy with declared areas")
	arcanaCommand := flags.String("arcana", "arcana", "Arcana executable")
	scopeList := flags.String("scope", ".", "comma-separated ownership coverage paths")
	scopeExcludeList := flags.String("scope-exclude", "", "comma-separated excluded coverage paths")
	ownershipKindList := flags.String("ownership-kinds", "file", "comma-separated node kinds checked for ownership")
	relationList := flags.String("relations", "", "comma-separated normalized architecture relationships")
	exampleLimit := flags.Int("example-limit", 20, "maximum unowned and overlap examples")
	format := flags.String("format", "text", "output format: text, json, or dot")
	timeout := flags.Duration("timeout", 4*time.Minute, "maximum analysis duration")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *policyPath == "" {
		fmt.Fprintln(stderr, "--policy is required")
		return 2
	}
	if *format != "text" && *format != "json" && *format != "dot" {
		fmt.Fprintf(stderr, "unsupported format %q\n", *format)
		return 2
	}
	if *timeout <= 0 || *exampleLimit <= 0 {
		fmt.Fprintln(stderr, "--timeout and --example-limit must be positive")
		return 2
	}
	document, err := policy.Load(*policyPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if len(document.Areas) == 0 {
		fmt.Fprintln(stderr, "policy must declare at least one area for analysis")
		return 2
	}
	resolvedSnapshot, err := snapshot.Resolve(*repo, *explicitSnapshot)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	scopePaths := splitCommaList(*scopeList)
	areaPaths := policy.AreaPaths(document)
	sourcePrefixes := uniqueStrings(append(append([]string(nil), scopePaths...), areaPaths...))
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	graph, err := (arcana.Client{Command: *arcanaCommand}).LoadGraphWithOptions(
		ctx,
		resolvedSnapshot,
		arcana.LoadOptions{
			SourcePrefixes:   sourcePrefixes,
			OutgoingPrefixes: areaPaths,
		},
	)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	analysis := policy.AnalyzeAreas(document, graph, policy.AreaAnalysisOptions{
		ScopePaths:        scopePaths,
		ScopeExcludePaths: splitCommaList(*scopeExcludeList),
		OwnershipKinds:    splitCommaList(*ownershipKindList),
		Relations:         splitCommaList(*relationList),
		ExampleLimit:      *exampleLimit,
	})
	switch *format {
	case "json":
		err = report.WriteAreaAnalysisJSON(stdout, analysis)
	case "dot":
		err = report.WriteAreaAnalysisDOT(stdout, analysis)
	default:
		err = report.WriteAreaAnalysisText(stdout, analysis)
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	return 0
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
