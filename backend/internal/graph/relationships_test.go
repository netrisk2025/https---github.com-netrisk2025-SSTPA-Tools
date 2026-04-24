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
	if got := len(Catalog()); got < 70 {
		t.Fatalf("expected at least 70 relationship entries, got %d", got)
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
		{Name: "HAS_LOSS", From: identity.NodeTypeAsset, To: identity.NodeTypeLoss},
		{Name: "MITIGATES", From: identity.NodeTypeControl, To: identity.NodeTypeHazard},
		{Name: "SATISFIES", From: identity.NodeTypeCountermeasure, To: identity.NodeTypeControl},
		{Name: "PARENTS", From: identity.NodeTypeRequirement, To: identity.NodeTypeRequirement},
		{Name: "HAS_HAZARD", From: identity.NodeTypeEnvironment, To: identity.NodeTypeHazard},
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
}

func recursionFor(name string, from identity.NodeType, to identity.NodeType) RecursionKind {
	for _, relationship := range Catalog() {
		if relationship.Name == name && relationship.From == from && relationship.To == to {
			return relationship.Recursion
		}
	}

	return RecursionNone
}
