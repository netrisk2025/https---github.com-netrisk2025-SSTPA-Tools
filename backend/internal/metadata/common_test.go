// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package metadata

import (
	"testing"
	"time"

	"sstpa-tool/backend/internal/identity"
)

func TestNewCommonAssignsOwnershipAndDefaults(t *testing.T) {
	now := time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)

	common, err := NewCommon(NewCommonInput{
		NodeType:  identity.NodeTypeSystem,
		HID:       "SYS_1_0",
		UUID:      "123e4567-e89b-12d3-a456-426614174000",
		Actor:     Actor{Name: "A. User", Email: "user@example.test"},
		Now:       now,
		VersionID: "schema-1",
	})
	if err != nil {
		t.Fatal(err)
	}

	if common.Owner != "A. User" || common.Creator != "A. User" {
		t.Fatalf("owner/creator not assigned from current user: %#v", common)
	}

	if common.Name != "New" {
		t.Fatalf("Name default = %q, want New", common.Name)
	}

	if common.ShortDescription != NullValue || common.LongDescription != NullValue {
		t.Fatalf("description defaults = %q/%q, want Null", common.ShortDescription, common.LongDescription)
	}
}

func TestPropertiesUseSRSPropertyNames(t *testing.T) {
	now := time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
	common, err := NewCommon(NewCommonInput{
		NodeType: identity.NodeTypeRequirement,
		HID:      "REQ_1_1",
		UUID:     "123e4567-e89b-12d3-a456-426614174000",
		Actor:    Actor{Name: "A. User", Email: "user@example.test"},
		Now:      now,
	})
	if err != nil {
		t.Fatal(err)
	}

	props := common.Properties()
	for _, key := range []string{
		"Name",
		"HID",
		"uuid",
		"TypeName",
		"Owner",
		"OwnerEmail",
		"Creator",
		"CreatorEmail",
		"Created",
		"LastTouch",
		"VersionID",
		"ShortDescription",
		"LongDescription",
	} {
		if _, ok := props[key]; !ok {
			t.Fatalf("missing common property %q", key)
		}
	}

	if props["VersionID"] != NullValue {
		t.Fatalf("empty VersionID = %q, want Null", props["VersionID"])
	}
}

func TestNewCommonRejectsIncompleteIdentity(t *testing.T) {
	_, err := NewCommon(NewCommonInput{
		NodeType: identity.NodeTypeSystem,
		Actor:    Actor{Name: "A. User", Email: "user@example.test"},
		Now:      time.Now(),
	})
	if err == nil {
		t.Fatal("expected missing HID/uuid to fail")
	}
}
