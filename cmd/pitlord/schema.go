package main

import (
	"flag"
	"fmt"
	"io"

	pitlordschema "github.com/Lokee86/pitlord/internal/schema"
)

func runSchema(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("schema", flag.ContinueOnError)
	flags.SetOutput(stderr)
	kind := flags.String("kind", "policy", "schema kind: policy or baseline")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	data, err := pitlordschema.Read(*kind)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if _, err := stdout.Write(data); err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if len(data) == 0 || data[len(data)-1] != '\n' {
		fmt.Fprintln(stdout)
	}
	return 0
}
