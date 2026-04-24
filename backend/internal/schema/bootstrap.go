// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package schema

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var bootstrapStatements = []string{
	"CREATE INDEX node_uuid_index IF NOT EXISTS FOR (n:SSTPANode) ON (n.uuid)",
	"CREATE INDEX node_name_index IF NOT EXISTS FOR (n:SSTPANode) ON (n.Name)",
	"CREATE INDEX node_type_index IF NOT EXISTS FOR (n:SSTPANode) ON (n.TypeName)",
	"CREATE INDEX user_email_index IF NOT EXISTS FOR (u:User) ON (u.UserEmail)",
	"CREATE INDEX user_uuid_index IF NOT EXISTS FOR (u:User) ON (u.uuid)",
	"CREATE INDEX admin_email_index IF NOT EXISTS FOR (a:Admin) ON (a.UserEmail)",
	"CREATE INDEX admin_uuid_index IF NOT EXISTS FOR (a:Admin) ON (a.uuid)",
	"CREATE INDEX mailbox_id_index IF NOT EXISTS FOR (m:Mailbox) ON (m.MailboxID)",
	"CREATE INDEX message_uuid_index IF NOT EXISTS FOR (m:Message) ON (m.uuid)",
	"CREATE CONSTRAINT node_hid_unique IF NOT EXISTS FOR (n:SSTPANode) REQUIRE n.HID IS UNIQUE",
	"MERGE (:Users)",
	"MERGE (:Admins)",
}

func Statements() []string {
	return append([]string(nil), bootstrapStatements...)
}

func Bootstrap(ctx context.Context, driver neo4j.DriverWithContext, databaseName string) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: databaseName})
	defer session.Close(ctx)

	for _, statement := range bootstrapStatements {
		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			result, err := tx.Run(ctx, statement, nil)
			if err != nil {
				return nil, err
			}

			_, err = result.Consume(ctx)
			return nil, err
		})
		if err != nil {
			return err
		}
	}

	return nil
}
