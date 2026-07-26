package main

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"

	"github.com/Lokee86/pitlord/internal/policy"
)

func runGenerate(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("generate", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestPath := flags.String("homunculus-manifest", "", "Homunculus manifest")
	outputPath := flags.String("output", "pitlord.json", "generated policy path")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *manifestPath == "" {
		fmt.Fprintln(stderr, "--homunculus-manifest is required")
		return 2
	}
	document, _, err := policy.FromHomunculus(*manifestPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	path, err := filepath.Abs(*outputPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if err := policy.Write(path, document); err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	fmt.Fprintln(stdout, path)
	return 0
}
