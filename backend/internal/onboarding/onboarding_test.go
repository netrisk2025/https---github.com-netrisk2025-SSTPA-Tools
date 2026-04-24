// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package onboarding

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sstpa-tool/backend/internal/identity"
	"sstpa-tool/backend/internal/metadata"
	"sstpa-tool/backend/internal/schema"
	"sstpa-tool/backend/internal/testhelpers"
)

func TestUserKindShape(t *testing.T) {
	if UserKind.NodeType != identity.NodeTypeUser {
		t.Fatalf("UserKind.NodeType = %q, want %q", UserKind.NodeType, identity.NodeTypeUser)
	}
	if UserKind.RegistryLabel != "UserRegistry" {
		t.Fatalf("UserKind.RegistryLabel = %q, want UserRegistry", UserKind.RegistryLabel)
	}
	if UserKind.Relationship != "HAS_USER" {
		t.Fatalf("UserKind.Relationship = %q, want HAS_USER", UserKind.Relationship)
	}
}

func TestAdminKindShape(t *testing.T) {
	if AdminKind.NodeType != identity.NodeTypeAdmin {
		t.Fatalf("AdminKind.NodeType = %q, want %q", AdminKind.NodeType, identity.NodeTypeAdmin)
	}
	if AdminKind.RegistryLabel != "AdminRegistry" {
		t.Fatalf("AdminKind.RegistryLabel = %q, want AdminRegistry", AdminKind.RegistryLabel)
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
	assertRegistryRelationship(t, fixture.Driver, "UserRegistry", "HAS_USER", schema.UserRegistryHID, record.HID)

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
	assertRegistryRelationship(t, fixture.Driver, "AdminRegistry", "HAS_ADMIN", schema.AdminRegistryHID, record.HID)

	list, err := List(ctx, fixture.Driver, "", AdminKind, Page{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("list admins: %v", err)
	}
	if list.Total != 1 {
		t.Fatalf("unexpected admin list: %#v", list)
	}
}

func TestBootstrapInstallerCreatesFirstAdminAndUserAtomically(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	if err := schema.Bootstrap(ctx, fixture.Driver, ""); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}

	result, err := BootstrapInstaller(ctx, fixture.Driver, "", BootstrapInput{
		InstallerName:  "Installer",
		InstallerEmail: "installer@example.test",
		Now:            time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("bootstrap installer: %v", err)
	}
	if !strings.HasPrefix(result.Admin.HID, "ADM__") || !strings.HasPrefix(result.User.HID, "USR__") {
		t.Fatalf("unexpected bootstrap records: %#v", result)
	}

	if _, err := BootstrapInstaller(ctx, fixture.Driver, "", BootstrapInput{
		InstallerName:  "Installer",
		InstallerEmail: "installer@example.test",
	}); err == nil {
		t.Fatal("expected duplicate installer bootstrap to fail")
	}

	registered, err := IsRegisteredAdmin(ctx, fixture.Driver, "", metadata.Actor{Name: "Installer", Email: "installer@example.test", Admin: true})
	if err != nil {
		t.Fatalf("registered admin lookup: %v", err)
	}
	if !registered {
		t.Fatal("expected bootstrapped installer to be a registered admin")
	}
}

func TestConcurrentUserCreatesAllocateUniqueHIDs(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	if err := schema.Bootstrap(ctx, fixture.Driver, ""); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}

	actor := metadata.Actor{Name: "Installer", Email: "installer@example.test", Admin: true}
	const count = 5
	var wg sync.WaitGroup
	hids := make(chan string, count)
	errs := make(chan error, count)
	for i := 0; i < count; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			record, err := Create(ctx, fixture.Driver, "", UserKind, CreateInput{
				UserName:  "User",
				UserEmail: "user" + string(rune('a'+i)) + "@example.test",
				Actor:     actor,
				Now:       time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC),
			})
			if err != nil {
				errs <- err
				return
			}
			hids <- record.HID
		}()
	}
	wg.Wait()
	close(hids)
	close(errs)

	for err := range errs {
		t.Fatalf("concurrent create failed: %v", err)
	}

	seen := map[string]struct{}{}
	for hid := range hids {
		if _, exists := seen[hid]; exists {
			t.Fatalf("duplicate HID allocated: %s", hid)
		}
		seen[hid] = struct{}{}
	}
	if len(seen) != count {
		t.Fatalf("created %d users, want %d", len(seen), count)
	}
}

func assertRegistryRelationship(t *testing.T, driver neo4j.DriverWithContext, registryLabel string, relationship string, registryHID string, childHID string) {
	t.Helper()
	ctx := context.Background()
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := "MATCH (registry:" + registryLabel + ":SSTPANode {HID: $registryHID})-[:" + relationship + "]->(child:SSTPANode {HID: $childHID}) RETURN count(child) AS c"
	result, err := session.Run(ctx, query, map[string]any{"registryHID": registryHID, "childHID": childHID})
	if err != nil {
		t.Fatalf("query registry relationship: %v", err)
	}
	record, err := result.Single(ctx)
	if err != nil {
		t.Fatalf("read registry relationship: %v", err)
	}
	count, _ := record.Get("c")
	if count != int64(1) {
		t.Fatalf("expected one %s relationship from %s to %s, got %#v", relationship, registryHID, childHID, count)
	}
}
