// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package graph

import (
	"testing"

	"sstpa-tool/backend/internal/identity"
)

func TestReferenceAssignmentAllowed(t *testing.T) {
	if !ReferenceAssignmentAllowed(identity.NodeTypeHazard, "MITRE ATT&CK", "Technique") {
		t.Fatal("expected Hazard to allow MITRE ATT&CK Technique references")
	}
}

func TestReferenceAssignmentAllowedRejectsInvalidMapping(t *testing.T) {
	if ReferenceAssignmentAllowed(identity.NodeTypeRequirement, "MITRE ATT&CK", "Technique") {
		t.Fatal("expected Requirement to reject external references")
	}
}
