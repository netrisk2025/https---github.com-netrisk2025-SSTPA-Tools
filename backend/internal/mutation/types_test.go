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

func TestPlanValidateRejectsCrossSoIRelationshipWithoutException(t *testing.T) {
	plan := Plan{Operations: []Operation{{
		Kind:             OperationCreateRelationship,
		RelationshipName: "FLOWS_TO_FUNCTION",
		FromHID:          "FUN_1_1",
		FromType:         identity.NodeTypeFunction,
		ToHID:            "FUN_2_1",
		ToType:           identity.NodeTypeFunction,
	}}}

	if err := plan.Validate(); err == nil {
		t.Fatal("expected cross-SoI relationship without exception to fail")
	}
}

func TestPlanValidateAcceptsRecordedCrossSoIException(t *testing.T) {
	plan := Plan{Operations: []Operation{{
		Kind:             OperationCreateRelationship,
		RelationshipName: "FLOWS_TO_FUNCTION",
		FromHID:          "FUN_1_1",
		FromType:         identity.NodeTypeFunction,
		ToHID:            "FUN_2_1",
		ToType:           identity.NodeTypeFunction,
		RelationshipProperties: map[string]any{
			"CrossSoIJustification": "analysis exception",
		},
	}}}

	if err := plan.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestPlanValidateRejectsLegacyRelationshipAliasForNewWrites(t *testing.T) {
	plan := Plan{Operations: []Operation{{
		Kind:             OperationCreateRelationship,
		RelationshipName: "Has",
		FromHID:          "AST_1_1",
		FromType:         identity.NodeTypeAsset,
		ToHID:            "REG_1_1",
		ToType:           identity.NodeTypeRegime,
	}}}

	if err := plan.Validate(); err == nil {
		t.Fatal("expected legacy relationship alias to fail for new writes")
	}
	if err := plan.ValidateWithOptions(ValidationOptions{AllowLegacyRelationshipAliases: true}); err != nil {
		t.Fatal(err)
	}
}

func TestPlanValidateRejectsLegacyBaronPropertyForNewWrites(t *testing.T) {
	plan := Plan{Operations: []Operation{{
		Kind:     OperationCreateNode,
		NodeType: identity.NodeTypeRequirement,
		HID:      "REQ_1_1",
		Properties: map[string]any{
			"Baron": true,
		},
	}}}

	if err := plan.Validate(); err == nil {
		t.Fatal("expected legacy Baron property to fail for new writes")
	}
	if err := plan.ValidateWithOptions(ValidationOptions{AllowLegacyPropertyAliases: true}); err != nil {
		t.Fatal(err)
	}
}

func TestPlanValidateSolutionEvidenceRules(t *testing.T) {
	valid := Plan{Operations: []Operation{{
		Kind:             OperationCreateRelationship,
		RelationshipName: "HAS_VERIFICATION",
		FromHID:          "SOL_1_1",
		FromType:         identity.NodeTypeSolution,
		ToHID:            "VER_1_1",
		ToType:           identity.NodeTypeVerification,
	}}}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}

	invalid := Plan{Operations: []Operation{{
		Kind:             OperationCreateRelationship,
		RelationshipName: "HAS_ASSET",
		FromHID:          "SOL_1_1",
		FromType:         identity.NodeTypeSolution,
		ToHID:            "AST_1_1",
		ToType:           identity.NodeTypeAsset,
	}}}
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected Solution to Asset evidence relationship to fail")
	}
}
