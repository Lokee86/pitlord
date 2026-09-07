package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/Lokee86/pitlord/internal/arcana"
	"github.com/Lokee86/pitlord/internal/calibration"
	"github.com/Lokee86/pitlord/internal/report"
	"github.com/Lokee86/pitlord/internal/scan"
	"github.com/Lokee86/pitlord/internal/snapshot"
)

func runCalibrate(args []string, stdout, stderr io.Writer) int {
	return runCalibrateWithDependencies(args, stdout, stderr, nil, calibration.GitRevisionResolver{})
}

func runCalibrateWithDependencies(
	args []string,
	stdout, stderr io.Writer,
	loader scan.GraphLoader,
	revisions calibration.RevisionResolver,
) int {
	flags := flag.NewFlagSet("calibrate", flag.ContinueOnError)
	flags.SetOutput(stderr)
	repo := flags.String("repo", ".", "pinned calibration corpus root")
	referenceFile := flags.String("reference", "", "calibration reference JSON")
	explicitSnapshot := flags.String("snapshot", "", "explicit Arcana snapshot directory")
	arcanaCommand := flags.String("arcana", "", "Arcana executable override")
	format := flags.String("format", "text", "output format: text or json")
	timeout := flags.Duration("timeout", 4*time.Minute, "maximum scan duration")
	skipRevision := flags.Bool("skip-revision-check", false, "do not verify corpus Git HEAD")
	failOnMismatch := flags.Bool("fail-on-mismatch", false, "exit 1 when labelled expectations mismatch")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *referenceFile == "" {
		fmt.Fprintln(stderr, "--reference is required")
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

	reference, err := calibration.LoadReference(*referenceFile)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if !*skipRevision {
		if revisions == nil {
			fmt.Fprintln(stderr, "revision resolver is required")
			return 2
		}
		if err := calibration.VerifyRevision(revisions, *repo, reference.SourceRevision); err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		if reference.WorktreeDiffSHA256 != "" {
			diffResolver, ok := revisions.(calibration.WorktreeDiffResolver)
			if !ok {
				fmt.Fprintln(stderr, "worktree diff resolver is required")
				return 2
			}
			if err := calibration.VerifyWorktreeDiff(diffResolver, *repo, reference.WorktreeDiffSHA256); err != nil {
				fmt.Fprintln(stderr, err)
				return 2
			}
		}
	}

	analyzerID := reference.Detector
	if reference.Analyzer != "" {
		analyzerID = reference.Analyzer
	}
	analyzer, err := scan.BuiltInAnalyzer(analyzerID)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	analyzers := []scan.Analyzer{analyzer}
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
	scanResult, err := (scan.Engine{Loader: loader, Analyzers: analyzers}).Run(ctx, scan.Input{
		RepositoryRoot: *repo,
		SnapshotPath:   resolvedSnapshot,
		PathPrefix:     reference.PathPrefix,
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	result := calibration.Evaluate(reference, scanResult)
	if *format == "json" {
		err = report.WriteCalibrationJSON(stdout, result)
	} else {
		err = report.WriteCalibrationText(stdout, result)
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if *failOnMismatch && result.HasMismatch() {
		return 1
	}
	return 0
}
