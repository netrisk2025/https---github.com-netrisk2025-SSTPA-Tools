// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package schema

import (
	"context"
	"strings"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sstpa-tool/backend/internal/testhelpers"
)

func TestStatementsIncludeSRSIndexes(t *testing.T) {
	joined := strings.Join(Statements(), "\n")

	for _, want := range []string{
		"DROP INDEX node_hid_index IF EXISTS",
		"node_uuid_index",
		"node_name_index",
		"node_type_index",
		"user_email_index",
		"mailbox_id_index",
		"message_uuid_index",
		"node_hid_unique",
	} {
		if !strings.Contains(joined, want) {
			t.Fatalf("expected bootstrap statements to include %q", want)
		}
	}
}

func TestStatementsAreIdempotent(t *testing.T) {
	for _, statement := range Statements() {
		if strings.Contains(statement, "IF NOT EXISTS") || strings.Contains(statement, "IF EXISTS") {
			continue
		}
		t.Fatalf("statement must be idempotent: %s", statement)
	}
}

func TestStatementsDoNotUsePluralRegistryLabels(t *testing.T) {
	joined := strings.Join(Statements(), "\n")

	for _, forbidden := range []string{":" + "Users", ":" + "Admins"} {
		if strings.Contains(joined, forbidden) {
			t.Fatalf("bootstrap statements must not include plural label %q", forbidden)
		}
	}
}

func TestBootstrapMaterializesToolDataRegistries(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()

	if err := Bootstrap(ctx, fixture.Driver, ""); err != nil {
		t.Fatalf("first bootstrap: %v", err)
	}
	if err := Bootstrap(ctx, fixture.Driver, ""); err != nil {
		t.Fatalf("second bootstrap: %v", err)
	}

	session := fixture.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	for _, test := range []struct {
		label    string
		typeName string
		name     string
		hid      string
	}{
		{label: "SSTPA_Tool", typeName: "SSTPA_Tool", name: SSTPAToolName, hid: SSTPAToolHID},
		{label: "UserRegistry", typeName: "UserRegistry", name: UserRegistryName, hid: UserRegistryHID},
		{label: "AdminRegistry", typeName: "AdminRegistry", name: AdminRegistryName, hid: AdminRegistryHID},
	} {
		query := "MATCH (n:" + test.label + ":SSTPANode {HID: $hid, Name: $name}) RETURN count(n) AS c"
		result, err := session.Run(ctx, query, map[string]any{"hid": test.hid, "name": test.name})
		if err != nil {
			t.Fatalf("run %s: %v", test.label, err)
		}
		record, err := result.Single(ctx)
		if err != nil {
			t.Fatalf("single %s: %v", test.label, err)
		}
		countValue, hasCount := record.Get("c")
		if !hasCount {
			t.Fatalf("missing 'c' column for :%s", test.label)
		}
		count, ok := countValue.(int64)
		if !ok {
			t.Fatalf("expected int64 count for :%s, got %T", test.label, countValue)
		}
		if count != 1 {
			t.Fatalf("expected exactly one :%s node after two bootstraps, got %d", test.label, count)
		}

		query = "MATCH (n:" + test.label + ":SSTPANode {HID: $hid}) RETURN n.TypeName AS typeName, n.uuid AS uuid, n.Owner AS owner, n.OwnerEmail AS ownerEmail, n.Creator AS creator, n.CreatorEmail AS creatorEmail, n.Created AS created, n.LastTouch AS lastTouch LIMIT 1"
		result, err = session.Run(ctx, query, map[string]any{"hid": test.hid})
		if err != nil {
			t.Fatalf("run metadata %s: %v", test.label, err)
		}
		record, err = result.Single(ctx)
		if err != nil {
			t.Fatalf("single metadata %s: %v", test.label, err)
		}
		for _, field := range []string{"uuid", "owner", "ownerEmail", "creator", "creatorEmail", "created", "lastTouch"} {
			value, _ := record.Get(field)
			if value == "" || value == nil {
				t.Fatalf("expected :%s to carry common metadata field %s", test.label, field)
			}
		}
		typeName, _ := record.Get("typeName")
		if typeName != test.typeName {
			t.Fatalf("TypeName for :%s = %#v, want %q", test.label, typeName, test.typeName)
		}
	}

	for _, query := range []string{
		"MATCH (:SSTPA_Tool:SSTPANode {HID: $toolHID})-[:HAS_USER_REGISTRY]->(:UserRegistry:SSTPANode {HID: $registryHID}) RETURN count(*) AS c",
		"MATCH (:SSTPA_Tool:SSTPANode {HID: $toolHID})-[:HAS_ADMIN_REGISTRY]->(:AdminRegistry:SSTPANode {HID: $registryHID}) RETURN count(*) AS c",
	} {
		registryHID := UserRegistryHID
		if strings.Contains(query, "HAS_ADMIN_REGISTRY") {
			registryHID = AdminRegistryHID
		}
		result, err := session.Run(ctx, query, map[string]any{"toolHID": SSTPAToolHID, "registryHID": registryHID})
		if err != nil {
			t.Fatalf("run registry relationship: %v", err)
		}
		record, err := result.Single(ctx)
		if err != nil {
			t.Fatalf("single registry relationship: %v", err)
		}
		count, _ := record.Get("c")
		if count != int64(1) {
			t.Fatalf("expected one registry relationship, got %#v", count)
		}
	}
}

func TestBootstrapMigratesOldHIDIndex(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	session := fixture.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	result, err := session.Run(ctx, "CREATE INDEX node_hid_index IF NOT EXISTS FOR (n:SSTPANode) ON (n.HID)", nil)
	if err != nil {
		t.Fatalf("create old index: %v", err)
	}
	if _, err := result.Consume(ctx); err != nil {
		t.Fatalf("consume old index create: %v", err)
	}

	if err := Bootstrap(ctx, fixture.Driver, ""); err != nil {
		t.Fatalf("bootstrap after old index: %v", err)
	}
}
