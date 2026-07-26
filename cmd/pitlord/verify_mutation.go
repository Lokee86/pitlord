package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/Lokee86/pitlord/internal/mutation"
	"github.com/Lokee86/pitlord/internal/snapshot"
)

func runVerifyMutation(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("verify-mutation", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestPath := flags.String("manifest", "", "Homunculus architecture-mutation.json")
	baselineSnapshot := flags.String("baseline-snapshot", "", "baseline Arcana snapshot directory")
	mutatedRepo := flags.String("repo", ".", "mutated repository root containing .arcana/CURRENT")
	mutatedSnapshot := flags.String("mutated-snapshot", "", "explicit mutated Arcana snapshot directory")
	arcanaCommand := flags.String("arcana", "arcana", "Arcana executable")
	format := flags.String("format", "text", "output format: text or json")
	timeout := flags.Duration("timeout", 2*time.Minute, "maximum verification duration")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *manifestPath == "" || *baselineSnapshot == "" {
		fmt.Fprintln(stderr, "--manifest and --baseline-snapshot are required")
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
	resolvedBaseline, err := snapshot.Resolve("", *baselineSnapshot)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	resolvedMutated, err := snapshot.Resolve(*mutatedRepo, *mutatedSnapshot)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	result, err := mutation.Verify(ctx, mutation.VerifyRequest{
		ManifestPath:     *manifestPath,
		BaselineSnapshot: resolvedBaseline,
		MutatedSnapshot:  resolvedMutated,
		ArcanaCommand:    *arcanaCommand,
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if *format == "json" {
		err = mutation.WriteJSON(stdout, result)
	} else {
		err = mutation.WriteText(stdout, result)
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if !result.Matched {
		return 1
	}
	return 0
}
