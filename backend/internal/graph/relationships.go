// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package graph

import (
	"fmt"

	"sstpa-tool/backend/internal/identity"
	"sstpa-tool/backend/internal/metadata"
)

type RecursionKind string

const (
	RecursionNone           RecursionKind = "none"
	RecursionDAG            RecursionKind = "dag"
	RecursionCyclicByDesign RecursionKind = "cyclic-by-design"
)

type SoIBoundaryRule string

const (
	SoIUnrestricted           SoIBoundaryRule = "unrestricted"
	SoISameIndex              SoIBoundaryRule = "same-index"
	SoISameIndexWithException SoIBoundaryRule = "same-index-with-exception"
	SoIMayCross               SoIBoundaryRule = "may-cross"

	DefaultRecursiveTraversalMaxDepth = 50
	CrossSoIJustificationProperty     = "CrossSoIJustification"
)

type RelationshipProperty struct {
	Name          string
	Required      bool
	Default       any
	AllowedValues []string
}

type Relationship struct {
	Name              string
	From              identity.NodeType
	To                identity.NodeType
	Recursion         RecursionKind
	SoI               SoIBoundaryRule
	Properties        []RelationshipProperty
	DefaultMaxDepth   int
	AllowMultiplicity bool
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

func LookupRelationshipWithLegacyAliases(name string, from identity.NodeType, to identity.NodeType, allowLegacy bool) (Relationship, string, bool) {
	if relationship, ok := LookupRelationship(name, from, to); ok {
		return relationship, name, true
	}

	if !allowLegacy {
		return Relationship{}, "", false
	}

	canonical, ok := LegacyRelationshipAlias(name, from, to)
	if !ok {
		return Relationship{}, "", false
	}

	relationship, ok := LookupRelationship(canonical, from, to)
	return relationship, canonical, ok
}

func LegacyRelationshipAlias(name string, from identity.NodeType, to identity.NodeType) (string, bool) {
	if name != "HAS" && name != "Has" {
		return "", false
	}

	candidates := []string{}
	for _, relationship := range relationshipCatalog {
		if relationship.From == from && relationship.To == to {
			candidates = append(candidates, relationship.Name)
		}
	}
	if len(candidates) != 1 {
		return "", false
	}

	return candidates[0], true
}

func DefaultRelationshipProperties(relationship Relationship) map[string]any {
	props := map[string]any{}
	for _, property := range relationship.Properties {
		if property.Required {
			props[property.Name] = property.Default
		}
	}

	return props
}

func TraversalMaxDepth(relationship Relationship) int {
	if relationship.DefaultMaxDepth > 0 {
		return relationship.DefaultMaxDepth
	}
	if relationship.Recursion == RecursionNone {
		return 0
	}

	return DefaultRecursiveTraversalMaxDepth
}

func ValidateRelationshipProperties(relationship Relationship, props map[string]any) error {
	for _, property := range relationship.Properties {
		value, ok := props[property.Name]
		if property.Required && (!ok || value == "") {
			return fmt.Errorf("relationship %s requires property %s", relationship.Name, property.Name)
		}
		if len(property.AllowedValues) == 0 || !ok {
			continue
		}
		text, ok := value.(string)
		if !ok {
			return fmt.Errorf("relationship %s property %s must be a string enum", relationship.Name, property.Name)
		}
		if !stringIn(text, property.AllowedValues) {
			return fmt.Errorf("relationship %s property %s value %q is not allowed", relationship.Name, property.Name, text)
		}
	}

	return nil
}

func ValidateSoIBoundary(relationship Relationship, fromHID string, toHID string, props map[string]any) error {
	switch relationship.SoI {
	case SoIUnrestricted, SoIMayCross:
		return nil
	}

	_, fromIndex, _, err := identity.ParseHID(fromHID)
	if err != nil {
		return fmt.Errorf("invalid from HID %s: %w", fromHID, err)
	}
	_, toIndex, _, err := identity.ParseHID(toHID)
	if err != nil {
		return fmt.Errorf("invalid to HID %s: %w", toHID, err)
	}
	if fromIndex == toIndex {
		return nil
	}

	if relationship.SoI == SoISameIndexWithException {
		if justification, ok := props[CrossSoIJustificationProperty].(string); ok && justification != "" && justification != metadata.NullValue {
			return nil
		}
	}

	return fmt.Errorf("relationship %s cannot cross SoI boundary from index %q to %q", relationship.Name, fromIndex, toIndex)
}

type relationshipOption func(*Relationship)

func rel(name string, from identity.NodeType, to identity.NodeType, options ...relationshipOption) Relationship {
	relationship := Relationship{
		Name: name,
		From: from,
		To:   to,
		SoI:  SoISameIndex,
	}
	for _, option := range options {
		option(&relationship)
	}
	if relationship.Recursion != RecursionNone && relationship.DefaultMaxDepth == 0 {
		relationship.DefaultMaxDepth = DefaultRecursiveTraversalMaxDepth
	}

	return relationship
}

func unrestricted() relationshipOption {
	return func(relationship *Relationship) { relationship.SoI = SoIUnrestricted }
}

func mayCross() relationshipOption {
	return func(relationship *Relationship) { relationship.SoI = SoIMayCross }
}

func sameIndexWithException() relationshipOption {
	return func(relationship *Relationship) { relationship.SoI = SoISameIndexWithException }
}

func recursion(kind RecursionKind) relationshipOption {
	return func(relationship *Relationship) { relationship.Recursion = kind }
}

func properties(properties ...RelationshipProperty) relationshipOption {
	return func(relationship *Relationship) {
		relationship.Properties = append(relationship.Properties, properties...)
	}
}

func requiredEnum(name string, defaultValue string, allowed ...string) RelationshipProperty {
	return RelationshipProperty{Name: name, Required: true, Default: defaultValue, AllowedValues: allowed}
}

func requiredDefault(name string, defaultValue any) RelationshipProperty {
	return RelationshipProperty{Name: name, Required: true, Default: defaultValue}
}

func stringIn(value string, candidates []string) bool {
	for _, candidate := range candidates {
		if candidate == value {
			return true
		}
	}

	return false
}

var transitionProperties = []RelationshipProperty{
	requiredEnum("TransitionKind", "FUNCTIONAL", "FUNCTIONAL", "COUNTERMEASURE_REQUIRED", "BOTH"),
}

var flowProperties = []RelationshipProperty{
	requiredEnum("RelationshipNature", "LOGICAL", "PHYSICAL", "LOGICAL", "BOTH"),
	requiredDefault("PhysicalType", metadata.NullValue),
	requiredDefault("LogicalLayer", metadata.NullValue),
	requiredDefault("Protocol", metadata.NullValue),
	requiredEnum("FlowDirectionality", "Unidirectional", "Unidirectional", "Bidirectional", "Multicast"),
	requiredDefault("TimingClass", metadata.NullValue),
	requiredDefault("SecurityClass", metadata.NullValue),
}

var relationshipCatalog = []Relationship{
	rel("HAS_USER_REGISTRY", identity.NodeTypeSSTPATool, identity.NodeTypeUserRegistry, unrestricted()),
	rel("HAS_ADMIN_REGISTRY", identity.NodeTypeSSTPATool, identity.NodeTypeAdminRegistry, unrestricted()),
	rel("HAS_USER", identity.NodeTypeUserRegistry, identity.NodeTypeUser, unrestricted()),
	rel("HAS_ADMIN", identity.NodeTypeAdminRegistry, identity.NodeTypeAdmin, unrestricted()),
	rel("HAS_MASTER_REGIME", identity.NodeTypeSSTPATool, identity.NodeTypeMasterRegime, unrestricted()),

	rel("HAS_SYSTEM", identity.NodeTypeCapability, identity.NodeTypeSystem, unrestricted()),
	rel("HAS_REQUIREMENT", identity.NodeTypeCapability, identity.NodeTypeRequirement, unrestricted()),
	rel("HAS_SYSTEM", identity.NodeTypeSandbox, identity.NodeTypeSystem, unrestricted()),

	rel("ACTS_IN", identity.NodeTypeSystem, identity.NodeTypeEnvironment),
	rel("HAS_CONNECTION", identity.NodeTypeSystem, identity.NodeTypeConnection),
	rel("HAS_INTERFACE", identity.NodeTypeSystem, identity.NodeTypeInterface),
	rel("HAS_FUNCTION", identity.NodeTypeSystem, identity.NodeTypeFunction),
	rel("HAS_ELEMENT", identity.NodeTypeSystem, identity.NodeTypeElement),
	rel("REALIZES", identity.NodeTypeSystem, identity.NodeTypePurpose),
	rel("EXHIBITS", identity.NodeTypeSystem, identity.NodeTypeState),
	rel("EXECUTES", identity.NodeTypeSystem, identity.NodeTypeControlStructure),
	rel("HAS_ASSET", identity.NodeTypeSystem, identity.NodeTypeAsset),
	rel("HAS_SECURITY", identity.NodeTypeSystem, identity.NodeTypeSecurity),
	rel("HAS_FUNCTIONAL_FLOW", identity.NodeTypeSystem, identity.NodeTypeFunctionalFlow),

	rel("HAS_HAZARD", identity.NodeTypeEnvironment, identity.NodeTypeHazard),
	rel("HAS_REQUIREMENT", identity.NodeTypeConnection, identity.NodeTypeRequirement),

	rel("HAS_REQUIREMENT", identity.NodeTypeInterface, identity.NodeTypeRequirement),
	rel("IMPLEMENTS", identity.NodeTypeInterface, identity.NodeTypeControlAlgorithm),
	rel("IMPLEMENTS", identity.NodeTypeInterface, identity.NodeTypeControlledProcess),
	rel("ALLOCATED_TO", identity.NodeTypeInterface, identity.NodeTypeElement),
	rel("CONTAINS", identity.NodeTypeInterface, identity.NodeTypeAsset),
	rel("PARTICIPATES_IN", identity.NodeTypeInterface, identity.NodeTypeConnection, mayCross(), properties(flowProperties...)),
	rel("CONNECTS", identity.NodeTypeInterface, identity.NodeTypeFunction, properties(flowProperties...)),

	rel("HAS_REQUIREMENT", identity.NodeTypeFunction, identity.NodeTypeRequirement),
	rel("IMPLEMENTS", identity.NodeTypeFunction, identity.NodeTypeControlAlgorithm),
	rel("IMPLEMENTS", identity.NodeTypeFunction, identity.NodeTypeControlledProcess),
	rel("IMPLEMENTS", identity.NodeTypeFunction, identity.NodeTypeProcessModel),
	rel("FLOWS_TO_FUNCTION", identity.NodeTypeFunction, identity.NodeTypeFunction, recursion(RecursionCyclicByDesign), sameIndexWithException(), properties(flowProperties...)),
	rel("FLOWS_TO_INTERFACE", identity.NodeTypeFunction, identity.NodeTypeInterface, sameIndexWithException(), properties(flowProperties...)),
	rel("ALLOCATED_TO", identity.NodeTypeFunction, identity.NodeTypeElement),
	rel("CONTAINS", identity.NodeTypeFunction, identity.NodeTypeAsset),

	rel("HAS_REQUIREMENT", identity.NodeTypeElement, identity.NodeTypeRequirement),
	rel("CONTAINS", identity.NodeTypeElement, identity.NodeTypeAsset),
	rel("PARENTS", identity.NodeTypeElement, identity.NodeTypeSystem, recursion(RecursionDAG), mayCross()),

	rel("HAS_CONSTRAINT", identity.NodeTypePurpose, identity.NodeTypeConstraint),
	rel("HAS_REQUIREMENT", identity.NodeTypePurpose, identity.NodeTypeRequirement),
	rel("HAS_VALIDATION", identity.NodeTypePurpose, identity.NodeTypeValidation),

	rel("TRANSITIONS_TO", identity.NodeTypeState, identity.NodeTypeState, recursion(RecursionCyclicByDesign), sameIndexWithException(), properties(transitionProperties...)),
	rel("HAS_HAZARD", identity.NodeTypeState, identity.NodeTypeHazard),
	rel("CONTAINS", identity.NodeTypeState, identity.NodeTypeAsset),

	rel("HAS_CONTROL_ALGORITHM", identity.NodeTypeControlStructure, identity.NodeTypeControlAlgorithm),
	rel("HAS_PROCESS_MODEL", identity.NodeTypeControlStructure, identity.NodeTypeProcessModel),
	rel("HAS_CONTROLLED_PROCESS", identity.NodeTypeControlStructure, identity.NodeTypeControlledProcess),
	rel("HAS_CONTROL_ACTION", identity.NodeTypeControlStructure, identity.NodeTypeControlAction),
	rel("HAS_FEEDBACK", identity.NodeTypeControlStructure, identity.NodeTypeFeedback),

	rel("CONTAINS", identity.NodeTypeFunctionalFlow, identity.NodeTypeFunction),
	rel("CONTAINS", identity.NodeTypeFunctionalFlow, identity.NodeTypeInterface),
	rel("CONTAINS", identity.NodeTypeFunctionalFlow, identity.NodeTypeConnection),
	rel("CONTAINS", identity.NodeTypeFunctionalFlow, identity.NodeTypeElement),
	rel("CONTAINS", identity.NodeTypeFunctionalFlow, identity.NodeTypeAsset),

	rel("HAS_REGIME", identity.NodeTypeAsset, identity.NodeTypeRegime),
	rel("HAS_LOSS", identity.NodeTypeAsset, identity.NodeTypeLoss),
	rel("HAS_GOAL", identity.NodeTypeAsset, identity.NodeTypeGoal),
	rel("DERIVED_FROM", identity.NodeTypeAsset, identity.NodeTypeAsset, recursion(RecursionDAG)),

	rel("HAS_CONTROL", identity.NodeTypeSecurity, identity.NodeTypeControl),
	rel("HAS_COUNTERMEASURE", identity.NodeTypeSecurity, identity.NodeTypeCountermeasure),

	rel("HAS_REQUIREMENT", identity.NodeTypeConstraint, identity.NodeTypeRequirement),
	rel("VIOLATES", identity.NodeTypeHazard, identity.NodeTypeConstraint),
	rel("THREATENS", identity.NodeTypeHazard, identity.NodeTypeAsset),
	rel("USES_ATTACK", identity.NodeTypeHazard, identity.NodeTypeAttack),

	rel("ENFORCES", identity.NodeTypeControl, identity.NodeTypeConstraint),
	rel("MITIGATES", identity.NodeTypeControl, identity.NodeTypeHazard),

	rel("SATISFIES", identity.NodeTypeCountermeasure, identity.NodeTypeControl),
	rel("HAS_REQUIREMENT", identity.NodeTypeCountermeasure, identity.NodeTypeRequirement),
	rel("APPLIES_TO_FUNCTION", identity.NodeTypeCountermeasure, identity.NodeTypeFunction, sameIndexWithException()),
	rel("APPLIES_TO_INTERFACE", identity.NodeTypeCountermeasure, identity.NodeTypeInterface, sameIndexWithException()),
	rel("APPLIES_TO_ELEMENT", identity.NodeTypeCountermeasure, identity.NodeTypeElement, sameIndexWithException()),
	rel("APPLIES_TO_STATE", identity.NodeTypeCountermeasure, identity.NodeTypeState, sameIndexWithException()),
	rel("APPLIES_TO_FEEDBACK", identity.NodeTypeCountermeasure, identity.NodeTypeFeedback),
	rel("BLOCKS", identity.NodeTypeCountermeasure, identity.NodeTypeAttack),

	rel("PARENTS", identity.NodeTypeRequirement, identity.NodeTypeRequirement, recursion(RecursionDAG), mayCross()),
	rel("VERIFIED_BY", identity.NodeTypeRequirement, identity.NodeTypeVerification),

	rel("GENERATES", identity.NodeTypeControlAlgorithm, identity.NodeTypeControlAction, recursion(RecursionCyclicByDesign)),
	rel("CAUSES", identity.NodeTypeControlAction, identity.NodeTypeHazard),
	rel("COMMANDS", identity.NodeTypeControlAction, identity.NodeTypeControlledProcess, recursion(RecursionCyclicByDesign)),
	rel("PRODUCES", identity.NodeTypeControlledProcess, identity.NodeTypeFeedback, recursion(RecursionCyclicByDesign)),
	rel("INFORMS", identity.NodeTypeFeedback, identity.NodeTypeProcessModel, recursion(RecursionCyclicByDesign)),
	rel("TUNES", identity.NodeTypeProcessModel, identity.NodeTypeControlAlgorithm, recursion(RecursionCyclicByDesign)),

	rel("HAS_ENVIRONMENT", identity.NodeTypeLoss, identity.NodeTypeEnvironment),
	rel("HAS_ELEMENT", identity.NodeTypeLoss, identity.NodeTypeElement),
	rel("HAS_STATE", identity.NodeTypeLoss, identity.NodeTypeState),
	rel("HAS_ATTACK", identity.NodeTypeLoss, identity.NodeTypeAttack),
	rel("HAS_COUNTERMEASURE", identity.NodeTypeLoss, identity.NodeTypeCountermeasure),

	rel("DEFEATS", identity.NodeTypeAttack, identity.NodeTypeCountermeasure),
	rel("EXPLOITS", identity.NodeTypeAttack, identity.NodeTypeElement),

	rel("SUPPORTED_BY", identity.NodeTypeGoal, identity.NodeTypeGoal, recursion(RecursionDAG)),
	rel("SUPPORTED_BY", identity.NodeTypeGoal, identity.NodeTypeStrategy, recursion(RecursionDAG)),
	rel("SUPPORTED_BY", identity.NodeTypeGoal, identity.NodeTypeSolution, recursion(RecursionDAG)),
	rel("IN_CONTEXT_OF", identity.NodeTypeGoal, identity.NodeTypeContext),
	rel("IN_CONTEXT_OF", identity.NodeTypeGoal, identity.NodeTypeJustification),
	rel("IN_CONTEXT_OF", identity.NodeTypeGoal, identity.NodeTypeAssumption),
	rel("SUPPORTED_BY", identity.NodeTypeStrategy, identity.NodeTypeGoal, recursion(RecursionDAG)),
	rel("SUPPORTED_BY", identity.NodeTypeStrategy, identity.NodeTypeSolution, recursion(RecursionDAG)),
	rel("IN_CONTEXT_OF", identity.NodeTypeStrategy, identity.NodeTypeContext),
	rel("IN_CONTEXT_OF", identity.NodeTypeStrategy, identity.NodeTypeJustification),
	rel("IN_CONTEXT_OF", identity.NodeTypeStrategy, identity.NodeTypeAssumption),
	rel("HAS_ENVIRONMENT", identity.NodeTypeContext, identity.NodeTypeEnvironment),
	rel("HAS_VALIDATION", identity.NodeTypeSolution, identity.NodeTypeValidation),
	rel("HAS_VERIFICATION", identity.NodeTypeSolution, identity.NodeTypeVerification),
	rel("HAS_LOSS", identity.NodeTypeSolution, identity.NodeTypeLoss),
}
