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

func TestLabelForAllNodeTypes(t *testing.T) {
	for _, nodeType := range identity.AllTypes() {
		label, ok := LabelFor(nodeType)
		if !ok {
			t.Fatalf("expected label for %q", nodeType)
		}

		if label != string(nodeType) {
			t.Fatalf("expected canonical singular label %q, got %q", nodeType, label)
		}
	}
}

func TestLabelForUnknownNodeType(t *testing.T) {
	if _, ok := LabelFor(identity.NodeType("Unknown")); ok {
		t.Fatal("expected unknown node type to be rejected")
	}
}
