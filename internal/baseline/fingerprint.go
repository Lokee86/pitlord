package baseline

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strconv"
	"strings"

	"github.com/Lokee86/pitlord/internal/arcana"
	"github.com/Lokee86/pitlord/internal/policy"
)

func Fingerprint(ruleID string, evidence policy.Evidence) string {
	areas := append([]string(nil), evidence.Areas...)
	sort.Strings(areas)
	parts := []string{
		ruleID,
		evidence.Issue,
		nodeIdentity(evidence.Source),
		evidence.Relation,
		targetIdentity(evidence.Target),
		strings.Join(areas, ","),
		evidence.SourceArea,
		evidence.TargetArea,
	}
	digest := sha256.Sum256([]byte(strings.Join(parts, "\x00")))
	return "sha256:" + hex.EncodeToString(digest[:])
}

func nodeIdentity(node arcana.Node) string {
	if node.Identity != "" {
		return "identity:" + node.Identity
	}
	if node.Key != "" {
		return "key:" + node.Key
	}
	parts := []string{node.Kind, node.Path, node.Name}
	if node.Span != nil {
		parts = append(parts,
			strconv.Itoa(node.Span.StartLine),
			strconv.Itoa(node.Span.StartColumn),
			strconv.Itoa(node.Span.EndLine),
			strconv.Itoa(node.Span.EndColumn),
		)
	}
	return "fallback:" + strings.Join(parts, "\x00")
}

func targetIdentity(node *arcana.Node) string {
	if node == nil {
		return ""
	}
	return nodeIdentity(*node)
}
