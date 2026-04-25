// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package graph

import (
	"regexp"
	"testing"

	"sstpa-tool/backend/internal/identity"
)

func TestCatalogRelationshipNamesUseUppercaseSnakeCase(t *testing.T) {
	pattern := regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

	for _, relationship := range Catalog() {
		if !pattern.MatchString(relationship.Name) {
			t.Fatalf("relationship %q does not use UPPERCASE_SNAKE_CASE", relationship.Name)
		}
	}
}

func TestCatalogHasExpectedCoverage(t *testing.T) {
	if got := len(Catalog()); got < 95 {
		t.Fatalf("expected at least 95 relationship entries, got %d", got)
	}
}

func TestCatalogSpotChecks(t *testing.T) {
	cases := []Relationship{
		{Name: "HAS_USER_REGISTRY", From: identity.NodeTypeSSTPATool, To: identity.NodeTypeUserRegistry},
		{Name: "HAS_ADMIN_REGISTRY", From: identity.NodeTypeSSTPATool, To: identity.NodeTypeAdminRegistry},
		{Name: "HAS_USER", From: identity.NodeTypeUserRegistry, To: identity.NodeTypeUser},
		{Name: "HAS_ADMIN", From: identity.NodeTypeAdminRegistry, To: identity.NodeTypeAdmin},
		{Name: "HAS_SYSTEM", From: identity.NodeTypeCapability, To: identity.NodeTypeSystem},
		{Name: "HAS_CONNECTION", From: identity.NodeTypeSystem, To: identity.NodeTypeConnection},
		{Name: "HAS_REQUIREMENT", From: identity.NodeTypePurpose, To: identity.NodeTypeRequirement},
		{Name: "TRANSITIONS_TO", From: identity.NodeTypeState, To: identity.NodeTypeState},
		{Name: "HAS_FUNCTIONAL_FLOW", From: identity.NodeTypeSystem, To: identity.NodeTypeFunctionalFlow},
		{Name: "ALLOCATED_TO", From: identity.NodeTypeInterface, To: identity.NodeTypeElement},
		{Name: "CONTAINS", From: identity.NodeTypeInterface, To: identity.NodeTypeAsset},
		{Name: "HAS_REGIME", From: identity.NodeTypeAsset, To: identity.NodeTypeRegime},
		{Name: "HAS_LOSS", From: identity.NodeTypeAsset, To: identity.NodeTypeLoss},
		{Name: "HAS_GOAL", From: identity.NodeTypeAsset, To: identity.NodeTypeGoal},
		{Name: "DERIVED_FROM", From: identity.NodeTypeAsset, To: identity.NodeTypeAsset},
		{Name: "MITIGATES", From: identity.NodeTypeControl, To: identity.NodeTypeHazard},
		{Name: "SATISFIES", From: identity.NodeTypeCountermeasure, To: identity.NodeTypeControl},
		{Name: "PARENTS", From: identity.NodeTypeRequirement, To: identity.NodeTypeRequirement},
		{Name: "HAS_HAZARD", From: identity.NodeTypeEnvironment, To: identity.NodeTypeHazard},
		{Name: "SUPPORTED_BY", From: identity.NodeTypeGoal, To: identity.NodeTypeStrategy},
		{Name: "HAS_VERIFICATION", From: identity.NodeTypeSolution, To: identity.NodeTypeVerification},
	}

	for _, test := range cases {
		if !AllowedRelationship(test.Name, test.From, test.To) {
			t.Fatalf("expected catalog to allow %s from %s to %s", test.Name, test.From, test.To)
		}
	}
}

func TestRecursiveRelationshipsDeclareGovernance(t *testing.T) {
	for _, relationship := range Catalog() {
		if relationship.From == relationship.To && relationship.Recursion == RecursionNone {
			t.Fatalf("recursive relationship %s must declare recursion governance", relationship.Name)
		}
	}

	if got := recursionFor("TRANSITIONS_TO", identity.NodeTypeState, identity.NodeTypeState); got != RecursionCyclicByDesign {
		t.Fatalf("TRANSITIONS_TO recursion = %q, want %q", got, RecursionCyclicByDesign)
	}

	if got := recursionFor("PARENTS", identity.NodeTypeRequirement, identity.NodeTypeRequirement); got != RecursionDAG {
		t.Fatalf("Requirement PARENTS recursion = %q, want %q", got, RecursionDAG)
	}

	if got := recursionFor("SUPPORTED_BY", identity.NodeTypeGoal, identity.NodeTypeStrategy); got != RecursionDAG {
		t.Fatalf("Goal SUPPORTED_BY Strategy recursion = %q, want %q", got, RecursionDAG)
	}
}

func TestDefaultRelationshipPropertiesSatisfyRegistry(t *testing.T) {
	cases := []Relationship{
		mustRelationship(t, "TRANSITIONS_TO", identity.NodeTypeState, identity.NodeTypeState),
		mustRelationship(t, "FLOWS_TO_FUNCTION", identity.NodeTypeFunction, identity.NodeTypeFunction),
		mustRelationship(t, "PARTICIPATES_IN", identity.NodeTypeInterface, identity.NodeTypeConnection),
	}

	for _, relationship := range cases {
		props := DefaultRelationshipProperties(relationship)
		if err := ValidateRelationshipProperties(relationship, props); err != nil {
			t.Fatalf("%s default properties failed validation: %v", relationship.Name, err)
		}
	}
}

func TestSoIBoundaryRules(t *testing.T) {
	flow := mustRelationship(t, "FLOWS_TO_FUNCTION", identity.NodeTypeFunction, identity.NodeTypeFunction)
	if err := ValidateSoIBoundary(flow, "FUN_1_1", "FUN_1_2", DefaultRelationshipProperties(flow)); err != nil {
		t.Fatalf("same-SoI flow should be allowed: %v", err)
	}
	if err := ValidateSoIBoundary(flow, "FUN_1_1", "FUN_2_1", DefaultRelationshipProperties(flow)); err == nil {
		t.Fatal("expected cross-SoI flow without justification to be rejected")
	}
	props := DefaultRelationshipProperties(flow)
	props[CrossSoIJustificationProperty] = "recorded analytical exception"
	if err := ValidateSoIBoundary(flow, "FUN_1_1", "FUN_2_1", props); err != nil {
		t.Fatalf("cross-SoI flow with justification should be allowed: %v", err)
	}

	participates := mustRelationship(t, "PARTICIPATES_IN", identity.NodeTypeInterface, identity.NodeTypeConnection)
	if err := ValidateSoIBoundary(participates, "INT_1_1", "CNN_2_1", DefaultRelationshipProperties(participates)); err != nil {
		t.Fatalf("PARTICIPATES_IN may cross SoI boundaries: %v", err)
	}
}

func TestLegacyRelationshipAliases(t *testing.T) {
	if AllowedRelationship("HAS", identity.NodeTypeAsset, identity.NodeTypeRegime) {
		t.Fatal("new writes must not treat generic HAS as canonical")
	}

	canonical, ok := LegacyRelationshipAlias("Has", identity.NodeTypeAsset, identity.NodeTypeRegime)
	if !ok || canonical != "HAS_REGIME" {
		t.Fatalf("legacy Asset-Has-Regime alias = (%q, %v), want HAS_REGIME", canonical, ok)
	}

	canonical, ok = LegacyRelationshipAlias("HAS", identity.NodeTypeAsset, identity.NodeTypeGoal)
	if !ok || canonical != "HAS_GOAL" {
		t.Fatalf("legacy Asset-HAS-Goal alias = (%q, %v), want HAS_GOAL", canonical, ok)
	}
}

func recursionFor(name string, from identity.NodeType, to identity.NodeType) RecursionKind {
	for _, relationship := range Catalog() {
		if relationship.Name == name && relationship.From == from && relationship.To == to {
			return relationship.Recursion
		}
	}

	return RecursionNone
}

func mustRelationship(t *testing.T, name string, from identity.NodeType, to identity.NodeType) Relationship {
	t.Helper()
	relationship, ok := LookupRelationship(name, from, to)
	if !ok {
		t.Fatalf("missing relationship %s from %s to %s", name, from, to)
	}

	return relationship
}
