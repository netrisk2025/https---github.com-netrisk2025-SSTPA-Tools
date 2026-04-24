// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package mutation

import (
	"testing"

	"sstpa-tool/backend/internal/identity"
)

func TestComputeAffectedIncludesChangedNodeAndRelationshipEndpoints(t *testing.T) {
	plan := Plan{Operations: []Operation{
		{Kind: OperationUpdateNode, HID: "REQ_1_1", Properties: map[string]any{"Name": "Updated"}},
		{
			Kind:             OperationCreateRelationship,
			RelationshipName: "HAS_REQUIREMENT",
			FromHID:          "PUR_1_1",
			FromType:         identity.NodeTypePurpose,
			ToHID:            "REQ_1_1",
			ToType:           identity.NodeTypeRequirement,
		},
	}}

	after := GraphSnapshot{Nodes: map[string]NodeSnapshot{
		"REQ_1_1": {HID: "REQ_1_1", Owner: "Alice", OwnerEmail: "alice@example.test"},
		"PUR_1_1": {HID: "PUR_1_1", Owner: "Bob", OwnerEmail: "bob@example.test"},
	}}

	got := ComputeAffected(plan, GraphSnapshot{}, after)
	if len(got) != 2 {
		t.Fatalf("expected 2 affected nodes, got %#v", got)
	}

	if got[0].HID != "PUR_1_1" || got[1].HID != "REQ_1_1" {
		t.Fatalf("affected nodes sorted by HID, got %#v", got)
	}
}
