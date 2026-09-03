package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/Lokee86/pitlord/internal/arcana"
	"github.com/Lokee86/pitlord/internal/report"
	"github.com/Lokee86/pitlord/internal/scan"
	"github.com/Lokee86/pitlord/internal/snapshot"
)

func runScan(args []string, stdout, stderr io.Writer) int {
	return runScanWithLoader(args, stdout, stderr, nil)
}

func runScanWithLoader(args []string, stdout, stderr io.Writer, loader scan.GraphLoader) int {
	flags := flag.NewFlagSet("scan", flag.ContinueOnError)
	flags.SetOutput(stderr)
	repo := flags.String("repo", ".", "repository root containing .arcana/CURRENT")
	explicitSnapshot := flags.String("snapshot", "", "explicit Arcana snapshot directory")
	arcanaCommand := flags.String("arcana", "", "Arcana executable override")
	pathPrefix := flags.String("path-prefix", ".", "repository-relative path prefix to scan")
	analyzerSpec := flags.String("analyzers", "architecture", "built-in analyzers: architecture, semantic, clippy, or a comma-separated combination")
	format := flags.String("format", "text", "output format: text or json")
	timeout := flags.Duration("timeout", 4*time.Minute, "maximum scan duration")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *format != "text" && *format != "json" {
		fmt.Fprintf(stderr, "unsupported format %q\n", *format)
		return 2
	}
	if *timeout <= 0 {
		fmt.Fprintln(stderr, "--timeout must be positive")
		return 2
	}
	analyzers, err := scan.BuiltInAnalyzers(*analyzerSpec)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	resolvedSnapshot := ""
	if scan.AnalyzersRequireGraph(analyzers) {
		resolvedSnapshot, err = snapshot.Resolve(*repo, *explicitSnapshot)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		if loader == nil {
			resolvedCommand, resolveErr := arcana.ResolveCommand(*repo, *arcanaCommand)
			if resolveErr != nil {
				fmt.Fprintln(stderr, resolveErr)
				return 2
			}
			loader = arcana.Client{Command: resolvedCommand}
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	result, err := (scan.Engine{Loader: loader, Analyzers: analyzers}).Run(ctx, scan.Input{
		RepositoryRoot: *repo,
		SnapshotPath:   resolvedSnapshot,
		PathPrefix:     *pathPrefix,
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if *format == "json" {
		err = report.WriteScanJSON(stdout, result)
	} else {
		err = report.WriteScanText(stdout, result)
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	return 0
}
