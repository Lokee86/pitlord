package main

import (
	"fmt"
	"io"
	"os"
)

var version = "0.1.2"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		writeUsage(stderr)
		return 2
	}
	switch args[0] {
	case "check":
		return runCheck(args[1:], stdout, stderr)
	case "baseline":
		return runBaseline(args[1:], stdout, stderr)
	case "validate":
		return runValidate(args[1:], stdout, stderr)
	case "schema":
		return runSchema(args[1:], stdout, stderr)
	case "verify-mutation":
		return runVerifyMutation(args[1:], stdout, stderr)
	case "diff":
		return runDiff(args[1:], stdout, stderr)
	case "analyze":
		return runAnalyze(args[1:], stdout, stderr)
	case "scan":
		return runScan(args[1:], stdout, stderr)
	case "inspect":
		return runInspect(args[1:], stdout, stderr)
	case "generate":
		return runGenerate(args[1:], stdout, stderr)
	case "version", "--version", "-version":
		fmt.Fprintln(stdout, version)
		return 0
	case "help", "--help", "-h":
		writeUsage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		writeUsage(stderr)
		return 2
	}
}

func writeUsage(writer io.Writer) {
	fmt.Fprintln(writer, "Pitlord enforces repository architecture policy over source files and Arcana snapshots.")
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "Usage:")
	fmt.Fprintln(writer, "  pitlord check --repo <root> --policy <pitlord.json> [--baseline <baseline.json>]")
	fmt.Fprintln(writer, "  pitlord check --repo <root> --homunculus-manifest <manifest.json>")
	fmt.Fprintln(writer, "  pitlord baseline --repo <root> --policy <pitlord.json> [--output <baseline.json>]")
	fmt.Fprintln(writer, "  pitlord validate --policy <pitlord.json>")
	fmt.Fprintln(writer, "  pitlord schema --kind policy|baseline")
	fmt.Fprintln(writer, "  pitlord verify-mutation --manifest <architecture-mutation.json> --baseline-snapshot <path>")
	fmt.Fprintln(writer, "  pitlord diff --before-snapshot <path> --repo <after-root> --policy <pitlord.json>")
	fmt.Fprintln(writer, "  pitlord analyze --repo <root> --policy <pitlord.json> [--format text|json|dot]")
	fmt.Fprintln(writer, "  pitlord scan --repo <root> [--path-prefix src] [--format text|json]")
	fmt.Fprintln(writer, "  pitlord inspect --repo <root> [--path-prefix src]")
	fmt.Fprintln(writer, "  pitlord generate --homunculus-manifest <manifest.json> [--output pitlord.json]")
	fmt.Fprintln(writer, "  pitlord version")
}
