package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/Lokee86/pitlord/internal/baseline"
	"github.com/Lokee86/pitlord/internal/mutation"
	"github.com/Lokee86/pitlord/internal/policy"
)

func runValidate(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("validate", flag.ContinueOnError)
	flags.SetOutput(stderr)
	policyPath := flags.String("policy", "", "Pitlord JSON policy")
	baselinePath := flags.String("baseline", "", "Pitlord evidence baseline")
	manifestPath := flags.String("homunculus-manifest", "", "Homunculus specimen manifest convertible to Pitlord policy")
	mutationPath := flags.String("homunculus-mutation", "", "Homunculus architecture-mutation.json")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	provided := 0
	for _, path := range []string{*policyPath, *baselinePath, *manifestPath, *mutationPath} {
		if path != "" {
			provided++
		}
	}
	if provided != 1 {
		fmt.Fprintln(stderr, "exactly one of --policy, --baseline, --homunculus-manifest, or --homunculus-mutation is required")
		return 2
	}

	switch {
	case *policyPath != "":
		document, err := policy.Load(*policyPath)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		fmt.Fprintf(
			stdout,
			"Valid Pitlord policy v%d: %d area(s), %d rule(s).\n",
			document.Version,
			len(document.Areas),
			len(document.Rules),
		)
	case *baselinePath != "":
		document, err := baseline.Load(*baselinePath)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		fmt.Fprintf(
			stdout,
			"Valid Pitlord baseline v%d: %d finding(s).\n",
			document.Version,
			len(document.Entries),
		)
	case *manifestPath != "":
		document, expected, err := policy.FromHomunculus(*manifestPath)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		fmt.Fprintf(
			stdout,
			"Valid Homunculus architecture contract: %d area(s), %d Pitlord rule(s), %d expected diagnostic(s).\n",
			len(document.Areas),
			len(document.Rules),
			len(expected),
		)
	case *mutationPath != "":
		manifest, err := mutation.Load(*mutationPath)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		fmt.Fprintf(
			stdout,
			"Valid Homunculus mutation contract v%d: %d added and %d removed relationship expectation(s).\n",
			manifest.Version,
			len(manifest.ExpectedAddedRelationships),
			len(manifest.ExpectedRemovedRelationships),
		)
	}
	return 0
}
