package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/Lokee86/pitlord/internal/checker"
	"github.com/Lokee86/pitlord/internal/policy"
	"github.com/Lokee86/pitlord/internal/snapshot"
	"github.com/Lokee86/pitlord/internal/snapshotdiff"
)

type diffEvaluation struct {
	result checker.Result
	err    error
}

func runDiff(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("diff", flag.ContinueOnError)
	flags.SetOutput(stderr)
	beforeSnapshot := flags.String("before-snapshot", "", "baseline Arcana snapshot directory")
	afterRepo := flags.String("repo", ".", "after repository root containing .arcana/CURRENT")
	afterSnapshot := flags.String("after-snapshot", "", "explicit after Arcana snapshot directory")
	policyPath := flags.String("policy", "", "Pitlord JSON policy")
	arcanaCommand := flags.String("arcana", "arcana", "Arcana executable")
	format := flags.String("format", "text", "output format: text, json, or sarif")
	timeout := flags.Duration("timeout", 4*time.Minute, "maximum comparison duration")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *beforeSnapshot == "" || *policyPath == "" {
		fmt.Fprintln(stderr, "--before-snapshot and --policy are required")
		return 2
	}
	if *format != "text" && *format != "json" && *format != "sarif" {
		fmt.Fprintf(stderr, "unsupported format %q\n", *format)
		return 2
	}
	if *timeout <= 0 {
		fmt.Fprintln(stderr, "--timeout must be positive")
		return 2
	}

	document, err := policy.Load(*policyPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	resolvedBefore, err := snapshot.Resolve("", *beforeSnapshot)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	resolvedAfter, err := snapshot.Resolve(*afterRepo, *afterSnapshot)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	before, after, err := evaluateDiffSnapshots(
		ctx,
		document,
		*policyPath,
		*afterRepo,
		resolvedBefore,
		resolvedAfter,
		*arcanaCommand,
	)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	result := snapshotdiff.Compare(
		*policyPath,
		resolvedBefore,
		resolvedAfter,
		before.Report.Diagnostics,
		after.Report.Diagnostics,
	)

	switch *format {
	case "json":
		err = snapshotdiff.WriteJSON(stdout, result)
	case "sarif":
		err = snapshotdiff.WriteIntroducedSARIF(stdout, result, version)
	default:
		err = snapshotdiff.WriteText(stdout, result)
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if len(result.Introduced) > 0 {
		return 1
	}
	return 0
}

func evaluateDiffSnapshots(
	ctx context.Context,
	document policy.Document,
	policySource string,
	repositoryRoot string,
	beforeSnapshot string,
	afterSnapshot string,
	arcanaCommand string,
) (checker.Result, checker.Result, error) {
	beforeChannel := make(chan diffEvaluation, 1)
	afterChannel := make(chan diffEvaluation, 1)
	go func() {
		result, err := checker.Evaluate(ctx, document, policySource, repositoryRoot, beforeSnapshot, arcanaCommand)
		beforeChannel <- diffEvaluation{result: result, err: err}
	}()
	go func() {
		result, err := checker.Evaluate(ctx, document, policySource, repositoryRoot, afterSnapshot, arcanaCommand)
		afterChannel <- diffEvaluation{result: result, err: err}
	}()
	before := <-beforeChannel
	after := <-afterChannel
	if before.err != nil {
		return checker.Result{}, checker.Result{}, fmt.Errorf("evaluate before snapshot: %w", before.err)
	}
	if after.err != nil {
		return checker.Result{}, checker.Result{}, fmt.Errorf("evaluate after snapshot: %w", after.err)
	}
	return before.result, after.result, nil
}
