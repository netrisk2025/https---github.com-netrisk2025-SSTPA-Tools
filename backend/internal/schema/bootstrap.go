// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package schema

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sstpa-tool/backend/internal/identity"
	"sstpa-tool/backend/internal/metadata"
)

var bootstrapStatements = []string{
	"DROP INDEX node_hid_index IF EXISTS",
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
}

const (
	SSTPAToolHID     = "SST__1"
	UserRegistryHID  = "URG__1"
	AdminRegistryHID = "ARG__1"

	SSTPAToolName     = "SSTPA Tools Data"
	UserRegistryName  = "Users"
	AdminRegistryName = "Admins"
)

var bootstrapActor = metadata.Actor{Name: "SSTPA Tool Bootstrap", Email: "sstpa-tool@localhost", Admin: true}

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

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, ensureToolData(ctx, tx, time.Now().UTC())
	})
	return err
}

func ensureToolData(ctx context.Context, tx neo4j.ManagedTransaction, now time.Time) error {
	toolProps, err := commonProperties(identity.NodeTypeSSTPATool, SSTPAToolHID, SSTPAToolName, now)
	if err != nil {
		return err
	}
	userRegistryProps, err := commonProperties(identity.NodeTypeUserRegistry, UserRegistryHID, UserRegistryName, now)
	if err != nil {
		return err
	}
	adminRegistryProps, err := commonProperties(identity.NodeTypeAdminRegistry, AdminRegistryHID, AdminRegistryName, now)
	if err != nil {
		return err
	}

	queries := []struct {
		statement string
		params    map[string]any
	}{
		{
			statement: `
MERGE (tool:SSTPA_Tool:SSTPANode {HID: $hid})
ON CREATE SET tool = $props, tool.UserSequence = 0, tool.AdminSequence = 0
ON MATCH SET tool.UserSequence = coalesce(tool.UserSequence, 0),
             tool.AdminSequence = coalesce(tool.AdminSequence, 0)
RETURN tool.HID AS hid
`,
			params: map[string]any{"hid": SSTPAToolHID, "props": toolProps},
		},
		{
			statement: `
MERGE (registry:UserRegistry:SSTPANode {HID: $hid})
ON CREATE SET registry = $props
RETURN registry.HID AS hid
`,
			params: map[string]any{"hid": UserRegistryHID, "props": userRegistryProps},
		},
		{
			statement: `
MERGE (registry:AdminRegistry:SSTPANode {HID: $hid})
ON CREATE SET registry = $props
RETURN registry.HID AS hid
`,
			params: map[string]any{"hid": AdminRegistryHID, "props": adminRegistryProps},
		},
		{
			statement: `
MATCH (tool:SSTPA_Tool:SSTPANode {HID: $toolHID}), (registry:UserRegistry:SSTPANode {HID: $registryHID})
MERGE (tool)-[:HAS_USER_REGISTRY]->(registry)
RETURN registry.HID AS hid
`,
			params: map[string]any{"toolHID": SSTPAToolHID, "registryHID": UserRegistryHID},
		},
		{
			statement: `
MATCH (tool:SSTPA_Tool:SSTPANode {HID: $toolHID}), (registry:AdminRegistry:SSTPANode {HID: $registryHID})
MERGE (tool)-[:HAS_ADMIN_REGISTRY]->(registry)
RETURN registry.HID AS hid
`,
			params: map[string]any{"toolHID": SSTPAToolHID, "registryHID": AdminRegistryHID},
		},
	}

	for _, query := range queries {
		result, err := tx.Run(ctx, query.statement, query.params)
		if err != nil {
			return err
		}
		if _, err := result.Single(ctx); err != nil {
			return err
		}
	}

	return nil
}

func commonProperties(nodeType identity.NodeType, hid string, name string, now time.Time) (map[string]any, error) {
	common, err := metadata.NewCommon(metadata.NewCommonInput{
		NodeType:  nodeType,
		HID:       hid,
		UUID:      identity.NewUUID(),
		Actor:     bootstrapActor,
		Now:       now,
		VersionID: "",
	})
	if err != nil {
		return nil, fmt.Errorf("bootstrap %s: %w", nodeType, err)
	}

	props := common.Properties()
	props["Name"] = name
	return props, nil
}
