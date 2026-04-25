// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package graph

import "sstpa-tool/backend/internal/identity"

type RecursionKind string

const (
	RecursionNone           RecursionKind = "none"
	RecursionDAG            RecursionKind = "dag"
	RecursionCyclicByDesign RecursionKind = "cyclic-by-design"
)

type Relationship struct {
	Name      string
	From      identity.NodeType
	To        identity.NodeType
	Recursion RecursionKind
}

func Catalog() []Relationship {
	return append([]Relationship(nil), relationshipCatalog...)
}

func AllowedRelationship(name string, from identity.NodeType, to identity.NodeType) bool {
	_, ok := LookupRelationship(name, from, to)
	return ok
}

func LookupRelationship(name string, from identity.NodeType, to identity.NodeType) (Relationship, bool) {
	for _, relationship := range relationshipCatalog {
		if relationship.Name == name && relationship.From == from && relationship.To == to {
			return relationship, true
		}
	}

	return Relationship{}, false
}

var relationshipCatalog = []Relationship{
	{Name: "HAS_USER_REGISTRY", From: identity.NodeTypeSSTPATool, To: identity.NodeTypeUserRegistry},
	{Name: "HAS_ADMIN_REGISTRY", From: identity.NodeTypeSSTPATool, To: identity.NodeTypeAdminRegistry},
	{Name: "HAS_USER", From: identity.NodeTypeUserRegistry, To: identity.NodeTypeUser},
	{Name: "HAS_ADMIN", From: identity.NodeTypeAdminRegistry, To: identity.NodeTypeAdmin},

	{Name: "HAS_SYSTEM", From: identity.NodeTypeCapability, To: identity.NodeTypeSystem},
	{Name: "HAS_REQUIREMENT", From: identity.NodeTypeCapability, To: identity.NodeTypeRequirement},
	{Name: "HAS_SYSTEM", From: identity.NodeTypeSandbox, To: identity.NodeTypeSystem},

	{Name: "ACTS_IN", From: identity.NodeTypeSystem, To: identity.NodeTypeEnvironment},
	{Name: "HAS_CONNECTION", From: identity.NodeTypeSystem, To: identity.NodeTypeConnection},
	{Name: "HAS_INTERFACE", From: identity.NodeTypeSystem, To: identity.NodeTypeInterface},
	{Name: "HAS_FUNCTION", From: identity.NodeTypeSystem, To: identity.NodeTypeFunction},
	{Name: "HAS_ELEMENT", From: identity.NodeTypeSystem, To: identity.NodeTypeElement},
	{Name: "REALIZES", From: identity.NodeTypeSystem, To: identity.NodeTypePurpose},
	{Name: "EXHIBITS", From: identity.NodeTypeSystem, To: identity.NodeTypeState},
	{Name: "EXECUTES", From: identity.NodeTypeSystem, To: identity.NodeTypeControlStructure},
	{Name: "HAS_ASSET", From: identity.NodeTypeSystem, To: identity.NodeTypeAsset},
	{Name: "HAS_SECURITY", From: identity.NodeTypeSystem, To: identity.NodeTypeSecurity},

	{Name: "HAS_HAZARD", From: identity.NodeTypeEnvironment, To: identity.NodeTypeHazard},
	{Name: "HAS_REQUIREMENT", From: identity.NodeTypeConnection, To: identity.NodeTypeRequirement},

	{Name: "HAS_REQUIREMENT", From: identity.NodeTypeInterface, To: identity.NodeTypeRequirement},
	{Name: "IMPLEMENTS", From: identity.NodeTypeInterface, To: identity.NodeTypeControlAlgorithm},
	{Name: "IMPLEMENTS", From: identity.NodeTypeInterface, To: identity.NodeTypeControlledProcess},
	{Name: "PARTICIPATES_IN", From: identity.NodeTypeInterface, To: identity.NodeTypeConnection},
	{Name: "CONNECTS", From: identity.NodeTypeInterface, To: identity.NodeTypeFunction},

	{Name: "HAS_REQUIREMENT", From: identity.NodeTypeFunction, To: identity.NodeTypeRequirement},
	{Name: "IMPLEMENTS", From: identity.NodeTypeFunction, To: identity.NodeTypeControlAlgorithm},
	{Name: "IMPLEMENTS", From: identity.NodeTypeFunction, To: identity.NodeTypeControlledProcess},
	{Name: "IMPLEMENTS", From: identity.NodeTypeFunction, To: identity.NodeTypeProcessModel},
	{Name: "FLOWS_TO_FUNCTION", From: identity.NodeTypeFunction, To: identity.NodeTypeFunction, Recursion: RecursionCyclicByDesign},
	{Name: "FLOWS_TO_INTERFACE", From: identity.NodeTypeFunction, To: identity.NodeTypeInterface},

	{Name: "HAS_REQUIREMENT", From: identity.NodeTypeElement, To: identity.NodeTypeRequirement},
	{Name: "CONTAINS", From: identity.NodeTypeElement, To: identity.NodeTypeAsset},
	{Name: "PARENTS", From: identity.NodeTypeElement, To: identity.NodeTypeSystem, Recursion: RecursionDAG},

	{Name: "HAS_CONSTRAINT", From: identity.NodeTypePurpose, To: identity.NodeTypeConstraint},
	{Name: "HAS_REQUIREMENT", From: identity.NodeTypePurpose, To: identity.NodeTypeRequirement},
	{Name: "HAS_VALIDATION", From: identity.NodeTypePurpose, To: identity.NodeTypeValidation},

	{Name: "TRANSITIONS_TO", From: identity.NodeTypeState, To: identity.NodeTypeState, Recursion: RecursionCyclicByDesign},
	{Name: "HAS_HAZARD", From: identity.NodeTypeState, To: identity.NodeTypeHazard},
	{Name: "CONTAINS", From: identity.NodeTypeState, To: identity.NodeTypeAsset},

	{Name: "HAS_CONTROL_ALGORITHM", From: identity.NodeTypeControlStructure, To: identity.NodeTypeControlAlgorithm},
	{Name: "HAS_PROCESS_MODEL", From: identity.NodeTypeControlStructure, To: identity.NodeTypeProcessModel},
	{Name: "HAS_CONTROLLED_PROCESS", From: identity.NodeTypeControlStructure, To: identity.NodeTypeControlledProcess},
	{Name: "HAS_CONTROL_ACTION", From: identity.NodeTypeControlStructure, To: identity.NodeTypeControlAction},
	{Name: "HAS_FEEDBACK", From: identity.NodeTypeControlStructure, To: identity.NodeTypeFeedback},

	{Name: "HAS_LOSS", From: identity.NodeTypeAsset, To: identity.NodeTypeLoss},

	{Name: "HAS_CONTROL", From: identity.NodeTypeSecurity, To: identity.NodeTypeControl},
	{Name: "HAS_COUNTERMEASURE", From: identity.NodeTypeSecurity, To: identity.NodeTypeCountermeasure},

	{Name: "HAS_REQUIREMENT", From: identity.NodeTypeConstraint, To: identity.NodeTypeRequirement},
	{Name: "VIOLATES", From: identity.NodeTypeHazard, To: identity.NodeTypeConstraint},
	{Name: "THREATENS", From: identity.NodeTypeHazard, To: identity.NodeTypeAsset},
	{Name: "USES_ATTACK", From: identity.NodeTypeHazard, To: identity.NodeTypeAttack},

	{Name: "ENFORCES", From: identity.NodeTypeControl, To: identity.NodeTypeConstraint},
	{Name: "MITIGATES", From: identity.NodeTypeControl, To: identity.NodeTypeHazard},

	{Name: "SATISFIES", From: identity.NodeTypeCountermeasure, To: identity.NodeTypeControl},
	{Name: "HAS_REQUIREMENT", From: identity.NodeTypeCountermeasure, To: identity.NodeTypeRequirement},
	{Name: "APPLIES_TO_FUNCTION", From: identity.NodeTypeCountermeasure, To: identity.NodeTypeFunction},
	{Name: "APPLIES_TO_INTERFACE", From: identity.NodeTypeCountermeasure, To: identity.NodeTypeInterface},
	{Name: "APPLIES_TO_ELEMENT", From: identity.NodeTypeCountermeasure, To: identity.NodeTypeElement},
	{Name: "APPLIES_TO_STATE", From: identity.NodeTypeCountermeasure, To: identity.NodeTypeState},
	{Name: "APPLIES_TO_FEEDBACK", From: identity.NodeTypeCountermeasure, To: identity.NodeTypeFeedback},
	{Name: "BLOCKS", From: identity.NodeTypeCountermeasure, To: identity.NodeTypeAttack},

	{Name: "PARENTS", From: identity.NodeTypeRequirement, To: identity.NodeTypeRequirement, Recursion: RecursionDAG},
	{Name: "VERIFIED_BY", From: identity.NodeTypeRequirement, To: identity.NodeTypeVerification},

	{Name: "GENERATES", From: identity.NodeTypeControlAlgorithm, To: identity.NodeTypeControlAction},
	{Name: "CAUSES", From: identity.NodeTypeControlAction, To: identity.NodeTypeHazard},
	{Name: "COMMANDS", From: identity.NodeTypeControlAction, To: identity.NodeTypeControlledProcess},
	{Name: "PRODUCES", From: identity.NodeTypeControlledProcess, To: identity.NodeTypeFeedback},
	{Name: "INFORMS", From: identity.NodeTypeFeedback, To: identity.NodeTypeProcessModel},
	{Name: "TUNES", From: identity.NodeTypeProcessModel, To: identity.NodeTypeControlAlgorithm},

	{Name: "HAS_ENVIRONMENT", From: identity.NodeTypeLoss, To: identity.NodeTypeEnvironment},
	{Name: "HAS_ELEMENT", From: identity.NodeTypeLoss, To: identity.NodeTypeElement},
	{Name: "HAS_STATE", From: identity.NodeTypeLoss, To: identity.NodeTypeState},
	{Name: "HAS_ATTACK", From: identity.NodeTypeLoss, To: identity.NodeTypeAttack},
	{Name: "HAS_COUNTERMEASURE", From: identity.NodeTypeLoss, To: identity.NodeTypeCountermeasure},

	{Name: "DEFEATS", From: identity.NodeTypeAttack, To: identity.NodeTypeCountermeasure},
	{Name: "EXPLOITS", From: identity.NodeTypeAttack, To: identity.NodeTypeElement},
}
