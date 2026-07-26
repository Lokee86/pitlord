package baseline

const Version = 1

type Document struct {
	Version int     `json:"version"`
	Entries []Entry `json:"entries"`
}

type Entry struct {
	Fingerprint string   `json:"fingerprint"`
	RuleID      string   `json:"rule_id"`
	Issue       string   `json:"issue"`
	Source      string   `json:"source"`
	Relation    string   `json:"relation,omitempty"`
	Target      string   `json:"target,omitempty"`
	Areas       []string `json:"areas,omitempty"`
	SourceArea  string   `json:"source_area,omitempty"`
	TargetArea  string   `json:"target_area,omitempty"`
}
