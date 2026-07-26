package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Lokee86/pitlord/internal/arcana"
	"github.com/Lokee86/pitlord/internal/report"
	"github.com/Lokee86/pitlord/internal/snapshot"
)

func runInspect(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("inspect", flag.ContinueOnError)
	flags.SetOutput(stderr)
	repo := flags.String("repo", ".", "repository root containing .arcana/CURRENT")
	explicitSnapshot := flags.String("snapshot", "", "explicit Arcana snapshot directory")
	arcanaCommand := flags.String("arcana", "arcana", "Arcana executable")
	pathPrefix := flags.String("path-prefix", ".", "repository path prefix to inspect")
	relationList := flags.String("relations", "", "comma-separated normalized relationships")
	minCommunitySize := flags.Int("min-community-size", 2, "minimum nodes in a returned community")
	limit := flags.Int("limit", 20, "maximum communities to return")
	format := flags.String("format", "text", "output format: text or json")
	timeout := flags.Duration("timeout", 2*time.Minute, "maximum inspection duration")
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
	snapshotPath, err := snapshot.Resolve(*repo, *explicitSnapshot)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	summary, err := (arcana.Client{Command: *arcanaCommand}).InspectArchitecture(
		ctx,
		snapshotPath,
		*pathPrefix,
		splitCommaList(*relationList),
		*minCommunitySize,
		*limit,
	)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if *format == "json" {
		err = report.WriteArchitectureJSON(stdout, summary)
	} else {
		err = report.WriteArchitectureText(stdout, summary)
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	return 0
}

func splitCommaList(value string) []string {
	var result []string
	for _, current := range strings.Split(value, ",") {
		current = strings.TrimSpace(current)
		if current != "" {
			result = append(result, current)
		}
	}
	return result
}
