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
		NodeTypeFunctionalFlow:    "FF",
		NodeTypeRegime:            "REG",
		NodeTypeControlledProcess: "CP",
		NodeTypeHazard:            "HAZ",
		NodeTypeLoss:              "LOS",
		NodeTypeAttack:            "ATK",
		NodeTypeGoal:              "G",
		NodeTypeSolution:          "SOL",
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
	// 35 V58 canonical labels + 6 tool/template data types.
	if got := len(AllTypes()); got != 41 {
		t.Fatalf("expected 41 node types (35 canonical + 6 tool/template), got %d", got)
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

func TestToolDataTypeIDs(t *testing.T) {
	cases := map[NodeType]string{
		NodeTypeSSTPATool:     "SST",
		NodeTypeUserRegistry:  "URG",
		NodeTypeAdminRegistry: "ARG",
		NodeTypeUser:          "USR",
		NodeTypeAdmin:         "ADM",
		NodeTypeMasterRegime:  "MRG",
	}

	for nodeType, want := range cases {
		got, ok := TypeID(nodeType)
		if !ok || got != want {
			t.Errorf("TypeID(%q) = (%q, %v), want (%q, true)", nodeType, got, ok, want)
		}
	}
}
