package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Lokee86/pitlord/internal/docguard"
)

func runDocs(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("docs", flag.ContinueOnError)
	flags.SetOutput(stderr)
	repository := flags.String("repo", ".", "repository root")
	snapshot := flags.String("snapshot", "", "explicit Arcana snapshot directory")
	docsRoot := flags.String("ddocs-root", "", "override the configured Demon Docs root")
	ddocs := flags.String("ddocs", "", "Demon Docs executable override")
	arcana := flags.String("arcana", "", "Arcana executable override")
	changedFrom := flags.String("changed-from", "", "Git revision used as the documentation-change baseline")
	format := flags.String("format", "text", "output format: text or json")
	timeout := flags.Duration("timeout", 2*time.Minute, "maximum guard duration")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if strings.TrimSpace(*changedFrom) == "" {
		fmt.Fprintln(stderr, "--changed-from is required")
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

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	dataset, err := docguard.LoadDataset(ctx, *repository, *docsRoot, *ddocs, *arcana)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	codeFiles, err := docguard.LoadCodeFiles(ctx, *repository, *snapshot, *arcana)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	changes, err := docguard.LoadChanges(ctx, *repository, *changedFrom)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	report := docguard.Evaluate(codeFiles, dataset, changes, *changedFrom)
	if err := writeDocsReport(stdout, report, *format); err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if len(report.Findings) > 0 {
		return 1
	}
	return 0
}

func writeDocsReport(writer io.Writer, report docguard.Report, format string) error {
	if format == "json" {
		encoder := json.NewEncoder(writer)
		encoder.SetIndent("", "  ")
		return encoder.Encode(report)
	}
	status := "PASS"
	if len(report.Findings) > 0 {
		status = "FAIL"
	}
	fmt.Fprintf(writer, "documentation guard: %s\n", status)
	fmt.Fprintf(writer, "code coverage: %d/%d mapped\n", report.MappedCodeFiles, report.CodeFiles)
	fmt.Fprintf(writer, "change check: %d code file(s), %d mapped document(s) changed\n", report.ChangedCode, report.ChangedDocs)
	for _, finding := range report.Findings {
		switch finding.Code {
		case "missing_codemap":
			fmt.Fprintf(writer, "- %s: no codemap ownership\n", finding.CodePath)
		case "documentation_not_changed":
			fmt.Fprintf(writer, "- %s: change requires one of [%s]\n", finding.CodePath, strings.Join(finding.Documents, ", "))
		}
	}
	return nil
}
