package schema

import (
	"embed"
	"fmt"
)

//go:embed files/*.json
var files embed.FS

func Read(kind string) ([]byte, error) {
	var path string
	switch kind {
	case "policy":
		path = "files/policy-v1.schema.json"
	case "baseline":
		path = "files/baseline-v1.schema.json"
	default:
		return nil, fmt.Errorf("unsupported schema kind %q", kind)
	}
	data, err := files.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read embedded %s schema: %w", kind, err)
	}
	return data, nil
}
