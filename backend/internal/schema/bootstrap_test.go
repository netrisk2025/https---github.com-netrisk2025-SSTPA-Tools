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
		if strings.Contains(statement, "IF NOT EXISTS") || strings.HasPrefix(statement, "MERGE ") {
			continue
		}
		t.Fatalf("statement must be idempotent: %s", statement)
	}
}

func TestStatementsIncludeUserAndAdminContainers(t *testing.T) {
	joined := strings.Join(Statements(), "\n")

	for _, want := range []string{
		"MERGE (:Users)",
		"MERGE (:Admins)",
		"user_uuid_index",
		"admin_email_index",
		"admin_uuid_index",
	} {
		if !strings.Contains(joined, want) {
			t.Fatalf("expected bootstrap statements to include %q", want)
		}
	}
}

func TestBootstrapMaterializesContainers(t *testing.T) {
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

	for _, label := range []string{"Users", "Admins"} {
		query := "MATCH (n:" + label + ") RETURN count(n) AS c"
		result, err := session.Run(ctx, query, nil)
		if err != nil {
			t.Fatalf("run %s: %v", label, err)
		}
		record, err := result.Single(ctx)
		if err != nil {
			t.Fatalf("single %s: %v", label, err)
		}
		countValue, hasCount := record.Get("c")
		if !hasCount {
			t.Fatalf("missing 'c' column for :%s", label)
		}
		count, ok := countValue.(int64)
		if !ok {
			t.Fatalf("expected int64 count for :%s, got %T", label, countValue)
		}
		if count != 1 {
			t.Fatalf("expected exactly one :%s node after two bootstraps, got %d", label, count)
		}
	}
}
