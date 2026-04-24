// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package identity

import "testing"

func TestTypeIDKnown(t *testing.T) {
	cases := map[NodeType]string{
		NodeTypeCapability:        "CAP",
		NodeTypeSystem:            "SYS",
		NodeTypeControlStructure:  "CS",
		NodeTypeControlledProcess: "CP",
		NodeTypeHazard:            "HAZ",
		NodeTypeLoss:              "LOS",
		NodeTypeAttack:            "ATK",
	}

	for nodeType, want := range cases {
		got, ok := TypeID(nodeType)
		if !ok || got != want {
			t.Errorf("TypeID(%q) = (%q, %v), want (%q, true)", nodeType, got, ok, want)
		}
	}
}

func TestTypeIDUnknown(t *testing.T) {
	if _, ok := TypeID(NodeType("Unknown")); ok {
		t.Fatal("expected ok=false for unknown type")
	}
}

func TestAllTypesCountCoversCoreAndToolData(t *testing.T) {
	// 27 Core Data Model types (SRS §1.3.6.1) + 2 Tool Data types (User, Admin per §1.4).
	if got := len(AllTypes()); got != 29 {
		t.Fatalf("expected 29 node types (27 Core + 2 Tool Data), got %d", got)
	}
}

func TestIsValidTypeID(t *testing.T) {
	if !IsValidTypeID("SYS") {
		t.Fatal("SYS must be valid")
	}

	if IsValidTypeID("NOPE") {
		t.Fatal("NOPE must not be valid")
	}
}

func TestUserAndAdminTypeIDs(t *testing.T) {
	cases := map[NodeType]string{
		NodeTypeUser:  "USR",
		NodeTypeAdmin: "ADM",
	}

	for nodeType, want := range cases {
		got, ok := TypeID(nodeType)
		if !ok || got != want {
			t.Errorf("TypeID(%q) = (%q, %v), want (%q, true)", nodeType, got, ok, want)
		}
	}
}
