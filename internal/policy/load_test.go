package policy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadMergesNestedRelativeIncludesOnce(t *testing.T) {
	root := t.TempDir()
	writePolicyFile(t, root, "shared/content.json", `{
  "version": 1,
  "rules": [
    {"id": "content", "type": "forbid_content", "literal": "bad", "include_paths": ["**/*.go"]}
  ]
}`)
	writePolicyFile(t, root, "areas/graph.json", `{
  "version": 1,
  "areas": [
    {"id": "api", "paths": ["api"]},
    {"id": "storage", "paths": ["storage"]}
  ],
  "rules": [
    {"id": "dependency", "type": "forbid_dependency", "from_areas": ["api"], "to_areas": ["storage"]}
  ]
}`)
	writePolicyFile(t, root, "nested.json", `{
  "version": 1,
  "includes": ["shared/content.json"],
  "rules": [
    {"id": "required", "type": "require_path", "path": "README.md"}
  ]
}`)
	writePolicyFile(t, root, "policy.json", `{
  "version": 1,
  "includes": ["nested.json", "shared/content.json", "areas/graph.json"],
  "rules": [
    {"id": "forbidden", "type": "forbid_path", "path": "tmp/**"}
  ]
}`)

	document, err := Load(filepath.Join(root, "policy.json"))
	if err != nil {
		t.Fatalf("load policy: %v", err)
	}
	if len(document.Includes) != 0 {
		t.Fatalf("loaded policy was not flattened: %#v", document.Includes)
	}
	if len(document.Areas) != 2 || len(document.Rules) != 4 {
		t.Fatalf("unexpected merged policy: areas=%d rules=%d", len(document.Areas), len(document.Rules))
	}
	ids := []string{document.Rules[0].ID, document.Rules[1].ID, document.Rules[2].ID, document.Rules[3].ID}
	joined := strings.Join(ids, ",")
	if joined != "dependency,content,required,forbidden" && joined != "content,required,dependency,forbidden" {
		t.Fatalf("unexpected deterministic include order: %s", joined)
	}
}

func TestLoadRejectsIncludeCycle(t *testing.T) {
	root := t.TempDir()
	writePolicyFile(t, root, "a.json", `{"version":1,"includes":["b.json"],"rules":[{"id":"a","type":"require_path","path":"a"}]}`)
	writePolicyFile(t, root, "b.json", `{"version":1,"includes":["a.json"],"rules":[{"id":"b","type":"require_path","path":"b"}]}`)
	_, err := Load(filepath.Join(root, "a.json"))
	if err == nil || !strings.Contains(err.Error(), "policy include cycle") {
		t.Fatalf("expected include cycle error, got %v", err)
	}
}

func TestLoadRejectsDuplicateRuleAcrossIncludes(t *testing.T) {
	root := t.TempDir()
	writePolicyFile(t, root, "a.json", `{"version":1,"rules":[{"id":"same","type":"require_path","path":"a"}]}`)
	writePolicyFile(t, root, "b.json", `{"version":1,"rules":[{"id":"same","type":"require_path","path":"b"}]}`)
	writePolicyFile(t, root, "policy.json", `{"version":1,"includes":["a.json","b.json"],"rules":[{"id":"root","type":"require_path","path":"root"}]}`)
	_, err := Load(filepath.Join(root, "policy.json"))
	if err == nil || !strings.Contains(err.Error(), `duplicate rule id "same"`) {
		t.Fatalf("expected duplicate rule error, got %v", err)
	}
}

func writePolicyFile(t *testing.T, root, relative, contents string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create policy directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write policy: %v", err)
	}
}
