package main

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"

	"github.com/Lokee86/pitlord/internal/baseline"
)

func runBaseline(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("baseline", flag.ContinueOnError)
	flags.SetOutput(stderr)
	options := bindEvaluationOptions(flags)
	output := flags.String("output", "pitlord.baseline.json", "baseline output path")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	result, err := options.run()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	path, err := filepath.Abs(*output)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	document := baseline.Build(result.Report.Diagnostics)
	if err := baseline.Write(path, document); err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	fmt.Fprintf(stdout, "Wrote %d baseline finding(s) to %s\n", len(document.Entries), path)
	return 0
}
