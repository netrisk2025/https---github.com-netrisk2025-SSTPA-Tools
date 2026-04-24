// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package onboarding

import (
	"context"
	"strings"
	"testing"
	"time"

	"sstpa-tool/backend/internal/identity"
	"sstpa-tool/backend/internal/metadata"
	"sstpa-tool/backend/internal/schema"
	"sstpa-tool/backend/internal/testhelpers"
)

func TestUserKindShape(t *testing.T) {
	if UserKind.NodeType != identity.NodeTypeUser {
		t.Fatalf("UserKind.NodeType = %q, want %q", UserKind.NodeType, identity.NodeTypeUser)
	}
	if UserKind.ContainerLabel != "Users" {
		t.Fatalf("UserKind.ContainerLabel = %q, want Users", UserKind.ContainerLabel)
	}
	if UserKind.Relationship != "HAS_USER" {
		t.Fatalf("UserKind.Relationship = %q, want HAS_USER", UserKind.Relationship)
	}
}

func TestAdminKindShape(t *testing.T) {
	if AdminKind.NodeType != identity.NodeTypeAdmin {
		t.Fatalf("AdminKind.NodeType = %q, want %q", AdminKind.NodeType, identity.NodeTypeAdmin)
	}
	if AdminKind.ContainerLabel != "Admins" {
		t.Fatalf("AdminKind.ContainerLabel = %q, want Admins", AdminKind.ContainerLabel)
	}
	if AdminKind.Relationship != "HAS_ADMIN" {
		t.Fatalf("AdminKind.Relationship = %q, want HAS_ADMIN", AdminKind.Relationship)
	}
}

func TestCreateListGetUser(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	if err := schema.Bootstrap(ctx, fixture.Driver, ""); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}

	installer := metadata.Actor{Name: "Installer", Email: "installer@example.test", Admin: true}
	now := time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)

	record, err := Create(ctx, fixture.Driver, "", UserKind, CreateInput{
		UserName:  "Alice Analyst",
		UserEmail: "alice@example.test",
		Actor:     installer,
		Now:       now,
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if !strings.HasPrefix(record.HID, "USR__") {
		t.Fatalf("expected HID to start with USR__, got %q", record.HID)
	}
	if record.TypeName != "User" || record.UserEmail != "alice@example.test" {
		t.Fatalf("unexpected record: %#v", record)
	}

	if _, err := Create(ctx, fixture.Driver, "", UserKind, CreateInput{
		UserName:  "Alice Again",
		UserEmail: "alice@example.test",
		Actor:     installer,
		Now:       now,
	}); err == nil {
		t.Fatal("expected duplicate registration to error")
	}

	list, err := List(ctx, fixture.Driver, "", UserKind, Page{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("list users: %v", err)
	}
	if list.Total != 1 || len(list.Items) != 1 {
		t.Fatalf("unexpected list: %#v", list)
	}

	fetched, err := GetByUUID(ctx, fixture.Driver, "", UserKind, record.UUID)
	if err != nil {
		t.Fatalf("get user: %v", err)
	}
	if fetched.HID != record.HID {
		t.Fatalf("fetched HID = %q, want %q", fetched.HID, record.HID)
	}
}

func TestCreateListGetAdmin(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	if err := schema.Bootstrap(ctx, fixture.Driver, ""); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}

	installer := metadata.Actor{Name: "Installer", Email: "installer@example.test", Admin: true}
	now := time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)

	record, err := Create(ctx, fixture.Driver, "", AdminKind, CreateInput{
		UserName:  "Root Admin",
		UserEmail: "root@example.test",
		Actor:     installer,
		Now:       now,
	})
	if err != nil {
		t.Fatalf("create admin: %v", err)
	}
	if !strings.HasPrefix(record.HID, "ADM__") {
		t.Fatalf("expected HID to start with ADM__, got %q", record.HID)
	}

	list, err := List(ctx, fixture.Driver, "", AdminKind, Page{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("list admins: %v", err)
	}
	if list.Total != 1 {
		t.Fatalf("unexpected admin list: %#v", list)
	}
}
