package scan

import "github.com/Lokee86/pitlord/internal/arcana"

// dependency-pressure fixtures use generic references for some reciprocal edges.
// Knot calibration keeps the same topology shape but projects those reciprocal
// edges onto a static dependency relation so the cycle detector is tested
// against the relation family it actually owns.
func directionalizeCalibrationReferences(graph arcana.Graph) arcana.Graph {
	for sourceID, relationships := range graph.Outgoing {
		updated := append([]arcana.Relationship(nil), relationships...)
		for index := range updated {
			if updated[index].Relation == "references" || updated[index].Relation == "calls" {
				updated[index].Relation = "imports"
			}
		}
		graph.Outgoing[sourceID] = updated
	}
	return graph
}
