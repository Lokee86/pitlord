package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/Lokee86/pitlord/internal/report"
	"github.com/Lokee86/pitlord/internal/scan"
	"github.com/Lokee86/pitlord/internal/snapshot"
)

func runScan(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("scan", flag.ContinueOnError)
	flags.SetOutput(stderr)
	repo := flags.String("repo", ".", "repository root containing .arcana/CURRENT")
	explicitSnapshot := flags.String("snapshot", "", "explicit Arcana snapshot directory")
	pathPrefix := flags.String("path-prefix", ".", "repository-relative path prefix to scan")
	format := flags.String("format", "text", "output format: text or json")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *format != "text" && *format != "json" {
		fmt.Fprintf(stderr, "unsupported format %q\n", *format)
		return 2
	}
	resolvedSnapshot, err := snapshot.Resolve(*repo, *explicitSnapshot)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	result, err := scan.Run(scan.Input{
		SnapshotPath: resolvedSnapshot,
		PathPrefix:   *pathPrefix,
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
