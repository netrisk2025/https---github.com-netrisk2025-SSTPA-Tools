// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package mutation

import (
	"testing"

	"sstpa-tool/backend/internal/metadata"
)

func TestValidatedUpdatePropertiesRequiresOwnerPair(t *testing.T) {
	_, err := validatedUpdateProperties(metadata.Actor{Name: "Alice", Email: "alice@example.test"}, fixedTime(), map[string]any{
		"Owner": "Alice",
	})
	if err == nil {
		t.Fatal("expected owner without owner email to fail")
	}
}

func TestValidatedUpdatePropertiesAllowsUserToAssumeOwnOwnership(t *testing.T) {
	props, err := validatedUpdateProperties(metadata.Actor{Name: "Alice", Email: "alice@example.test"}, fixedTime(), map[string]any{
		"Owner":      "Alice",
		"OwnerEmail": "alice@example.test",
	})
	if err != nil {
		t.Fatal(err)
	}

	if props["Owner"] != "Alice" || props["OwnerEmail"] != "alice@example.test" {
		t.Fatalf("unexpected ownership props: %#v", props)
	}
}

func TestValidatedUpdatePropertiesRejectsNonAdminAssigningOtherOwner(t *testing.T) {
	_, err := validatedUpdateProperties(metadata.Actor{Name: "Alice", Email: "alice@example.test"}, fixedTime(), map[string]any{
		"Owner":      "Bob",
		"OwnerEmail": "bob@example.test",
	})
	if err == nil {
		t.Fatal("expected non-admin assignment to another owner to fail")
	}
}

func TestValidatedUpdatePropertiesRejectsCreatorChangeForNonAdmin(t *testing.T) {
	_, err := validatedUpdateProperties(metadata.Actor{Name: "Alice", Email: "alice@example.test"}, fixedTime(), map[string]any{
		"Creator": "Bob",
	})
	if err == nil {
		t.Fatal("expected non-admin creator edit to fail")
	}
}
