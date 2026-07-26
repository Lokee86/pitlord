package mutation

import (
	"encoding/json"
	"fmt"
	"io"
)

func WriteJSON(writer io.Writer, result Result) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(result)
}

func WriteText(writer io.Writer, result Result) error {
	matched := 0
	for _, check := range result.Checks {
		status := "FAIL"
		if check.Matched {
			status = "PASS"
			matched++
		}
		if _, err := fmt.Fprintf(
			writer,
			"%s %s %s --%s--> %s (baseline=%t mutated=%t)\n",
			status,
			check.Expectation,
			check.Relationship.Source,
			check.Relationship.Relation,
			check.Relationship.Target,
			check.BaselinePresent,
			check.MutatedPresent,
		); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(
		writer,
		"Pitlord mutation verification %s: %d/%d expectation(s) matched.\n",
		result.MutationID,
		matched,
		len(result.Checks),
	)
	return err
}
