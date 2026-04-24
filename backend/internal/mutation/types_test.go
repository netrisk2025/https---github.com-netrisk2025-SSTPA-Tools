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

func TestPlanValidateRejectsEmptyPlan(t *testing.T) {
	if err := (Plan{}).Validate(); err == nil {
		t.Fatal("expected empty plan to fail")
	}
}

func TestPlanValidateAcceptsKnownCreateNode(t *testing.T) {
	plan := Plan{Operations: []Operation{{
		Kind:     OperationCreateNode,
		NodeType: identity.NodeTypeCapability,
		HID:      "CAP__0",
	}}}

	if err := plan.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestPlanValidateRejectsInvalidRelationship(t *testing.T) {
	plan := Plan{Operations: []Operation{{
		Kind:             OperationCreateRelationship,
		RelationshipName: "HAS_SYSTEM",
		FromHID:          "REQ_1_1",
		FromType:         identity.NodeTypeRequirement,
		ToHID:            "SYS_1_0",
		ToType:           identity.NodeTypeSystem,
	}}}

	if err := plan.Validate(); err == nil {
		t.Fatal("expected invalid relationship endpoint types to fail")
	}
}
