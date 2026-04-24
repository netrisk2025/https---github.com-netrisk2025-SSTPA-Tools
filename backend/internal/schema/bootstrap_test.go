// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package schema

import (
	"strings"
	"testing"
)

func TestStatementsIncludeSRSIndexes(t *testing.T) {
	joined := strings.Join(Statements(), "\n")

	for _, want := range []string{
		"node_hid_index",
		"node_uuid_index",
		"node_name_index",
		"node_type_index",
		"user_email_index",
		"mailbox_id_index",
		"message_uuid_index",
	} {
		if !strings.Contains(joined, want) {
			t.Fatalf("expected bootstrap statements to include %q", want)
		}
	}
}

func TestStatementsAreIdempotent(t *testing.T) {
	for _, statement := range Statements() {
		if !strings.Contains(statement, "IF NOT EXISTS") {
			t.Fatalf("statement must be idempotent: %s", statement)
		}
	}
}
