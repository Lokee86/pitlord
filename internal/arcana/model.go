package arcana

type Span struct {
	Path        string `json:"path"`
	StartLine   int    `json:"start_line"`
	StartColumn int    `json:"start_column"`
	EndLine     int    `json:"end_line"`
	EndColumn   int    `json:"end_column"`
}

type Node struct {
	NodeID   uint32 `json:"node_id"`
	Key      string `json:"key"`
	Identity string `json:"identity"`
	Kind     string `json:"kind"`
	Path     string `json:"path"`
	Name     string `json:"name"`
	Span     *Span  `json:"span,omitempty"`
}

type Relationship struct {
	Relation string `json:"relation"`
	Node     Node   `json:"node"`
}

type Graph struct {
	Sources       []Node
	Outgoing      map[uint32][]Relationship
	Relationships int
}
