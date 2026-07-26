package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/Lokee86/pitlord/internal/baseline"
	"github.com/Lokee86/pitlord/internal/policy"
	"github.com/Lokee86/pitlord/internal/report"
)

func runCheck(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("check", flag.ContinueOnError)
	flags.SetOutput(stderr)
	options := bindEvaluationOptions(flags)
	baselinePath := flags.String("baseline", "", "evidence baseline to suppress")
	format := flags.String("format", "text", "output format: text, json, or sarif")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *format != "text" && *format != "json" && *format != "sarif" {
		fmt.Fprintf(stderr, "unsupported format %q\n", *format)
		return 2
	}

	result, err := options.run()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if *baselinePath != "" {
		document, err := baseline.Load(*baselinePath)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		filtered, suppressed := baseline.Apply(document, result.Report.Diagnostics)
		result.Report.Diagnostics = filtered
		result.Report.Summary = policy.BuildSummary(result.Document, result.Graph, filtered)
		result.Report.Summary.SuppressedEvidence = suppressed
	}

	switch *format {
	case "json":
		err = report.WriteJSON(stdout, result.Report)
	case "sarif":
		err = report.WriteSARIF(stdout, result.Report, version)
	default:
		err = report.WriteText(stdout, result.Report)
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if len(result.Report.Diagnostics) > 0 ||
		(result.Report.Expectation != nil && !result.Report.Expectation.Matched) {
		return 1
	}
	return 0
}
