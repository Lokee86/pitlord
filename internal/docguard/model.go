package docguard

const DatasetSchemaVersion = 1

type Dataset struct {
	SchemaVersion int            `json:"schema_version"`
	Entries       []DatasetEntry `json:"entries"`
}

type DatasetEntry struct {
	Entry      Entry      `json:"entry"`
	Resolution Resolution `json:"resolution"`
}

type Entry struct {
	DocumentPath string `json:"document_path"`
	Target       string `json:"target"`
	Kind         string `json:"kind"`
}

type Resolution struct {
	Status       string        `json:"status"`
	ResolvedPath string        `json:"resolved_path,omitempty"`
	Matches      []TargetMatch `json:"matches,omitempty"`
	SemanticNode *SemanticNode `json:"semantic_node,omitempty"`
}

type SemanticNode struct {
	Path string `json:"path"`
}

type TargetMatch struct {
	Path  string `json:"path"`
	IsDir bool   `json:"is_directory"`
}

type Change struct {
	Status  string
	Path    string
	OldPath string
}

type Finding struct {
	Code      string   `json:"code"`
	CodePath  string   `json:"code_path"`
	Documents []string `json:"documents,omitempty"`
}

type Report struct {
	Schema          string    `json:"schema"`
	ChangedFrom     string    `json:"changed_from"`
	CodeFiles       int       `json:"code_files"`
	MappedCodeFiles int       `json:"mapped_code_files"`
	ChangedCode     int       `json:"changed_code"`
	ChangedDocs     int       `json:"changed_docs"`
	Findings        []Finding `json:"findings"`
}
