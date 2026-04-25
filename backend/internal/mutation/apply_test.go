// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package mutation

import (
	"context"
	"testing"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	"sstpa-tool/backend/internal/identity"
	"sstpa-tool/backend/internal/metadata"
	"sstpa-tool/backend/internal/schema"
	"sstpa-tool/backend/internal/telemetry"
	"sstpa-tool/backend/internal/testhelpers"
)

func TestApplyCreatesChangeNotificationInSameTransaction(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	if err := schema.Bootstrap(ctx, fixture.Driver, ""); err != nil {
		t.Fatal(err)
	}

	alice := metadata.Actor{Name: "Alice", Email: "alice@example.test"}
	bob := metadata.Actor{Name: "Bob", Email: "bob@example.test"}

	_, err := Apply(ctx, fixture.Driver, ApplyOptions{
		Actor:     alice,
		Now:       fixedTime(),
		CommitID:  "commit-create",
		VersionID: "test-schema",
	}, Plan{Operations: []Operation{{
		Kind:     OperationCreateNode,
		NodeType: identity.NodeTypeCapability,
		HID:      "CAP__0",
		UUID:     "00000000-0000-4000-8000-000000000001",
		Properties: map[string]any{
			"Name": "Capability",
		},
	}}})
	if err != nil {
		t.Fatal(err)
	}

	report, err := Apply(ctx, fixture.Driver, ApplyOptions{
		Actor:    bob,
		Now:      fixedTime().Add(time.Minute),
		CommitID: "commit-update",
	}, Plan{Operations: []Operation{{
		Kind: OperationUpdateNode,
		HID:  "CAP__0",
		Properties: map[string]any{
			"Name": "Changed by Bob",
		},
	}}})
	if err != nil {
		t.Fatal(err)
	}

	if report.MessagesGenerated != 1 {
		t.Fatalf("MessagesGenerated = %d, want 1", report.MessagesGenerated)
	}
	if len(report.RecipientsNotified) != 1 || report.RecipientsNotified[0] != alice.Email {
		t.Fatalf("RecipientsNotified = %#v, want Alice", report.RecipientsNotified)
	}

	session := fixture.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	values, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
MATCH (mailbox:Mailbox {MailboxID: $mailboxID})-[:HAS_MESSAGE]->(message:Message)
MATCH (cap {HID: "CAP__0"})
RETURN count(message) AS messages, cap.Name AS name, mailbox.UnreadCount AS unread
`, map[string]any{"mailboxID": alice.Email})
		if err != nil {
			return nil, err
		}
		return result.Single(ctx)
	})
	if err != nil {
		t.Fatal(err)
	}

	record := values.(*neo4j.Record)
	messageCount, _ := record.Get("messages")
	name, _ := record.Get("name")
	unread, _ := record.Get("unread")
	if messageCount != int64(1) || name != "Changed by Bob" || unread != int64(1) {
		t.Fatalf("unexpected db state: messages=%#v name=%#v unread=%#v", messageCount, name, unread)
	}
}

func TestApplyRollsBackWhenNotificationCannotBeCreated(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	session := fixture.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
CREATE (:Capability {
  HID: "CAP__0",
  uuid: "00000000-0000-4000-8000-000000000002",
  TypeName: "Capability",
  Name: "Original",
  Owner: "Alice",
  OwnerEmail: "",
  Creator: "Alice",
  CreatorEmail: "alice@example.test",
  Created: "2026-04-24T12:00:00Z",
  LastTouch: "2026-04-24T12:00:00Z",
  VersionID: "test-schema",
  ShortDescription: "Null",
  LongDescription: "Null"
})
`, nil)
		if err != nil {
			return nil, err
		}
		_, err = result.Consume(ctx)
		return nil, err
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = Apply(ctx, fixture.Driver, ApplyOptions{
		Actor:    metadata.Actor{Name: "Bob", Email: "bob@example.test"},
		Now:      fixedTime().Add(time.Minute),
		CommitID: "commit-rollback",
	}, Plan{Operations: []Operation{{
		Kind: OperationUpdateNode,
		HID:  "CAP__0",
		Properties: map[string]any{
			"Name": "Should Roll Back",
		},
	}}})
	if err == nil {
		t.Fatal("expected notification failure to roll back transaction")
	}

	value, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `MATCH (cap {HID: "CAP__0"}) RETURN cap.Name AS name`, nil)
		if err != nil {
			return nil, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		name, _ := record.Get("name")
		return name, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if value != "Original" {
		t.Fatalf("name = %#v, want rollback to preserve Original", value)
	}
}

func TestApplyRejectsDuplicateRelationship(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	alice := metadata.Actor{Name: "Alice", Email: "alice@example.test"}

	_, err := Apply(ctx, fixture.Driver, ApplyOptions{Actor: alice, Now: fixedTime()}, Plan{Operations: []Operation{
		{Kind: OperationCreateNode, NodeType: identity.NodeTypeCapability, HID: "CAP__0"},
		{Kind: OperationCreateNode, NodeType: identity.NodeTypeSystem, HID: "SYS_1_0"},
		{
			Kind:             OperationCreateRelationship,
			RelationshipName: "HAS_SYSTEM",
			FromHID:          "CAP__0",
			FromType:         identity.NodeTypeCapability,
			ToHID:            "SYS_1_0",
			ToType:           identity.NodeTypeSystem,
		},
	}})
	if err != nil {
		t.Fatal(err)
	}

	_, err = Apply(ctx, fixture.Driver, ApplyOptions{Actor: alice, Now: fixedTime()}, Plan{Operations: []Operation{{
		Kind:             OperationCreateRelationship,
		RelationshipName: "HAS_SYSTEM",
		FromHID:          "CAP__0",
		FromType:         identity.NodeTypeCapability,
		ToHID:            "SYS_1_0",
		ToType:           identity.NodeTypeSystem,
	}}})
	if err == nil {
		t.Fatal("expected duplicate relationship to be rejected")
	}
}

func TestApplyCanonicalizesLegacyRelationshipAliasWhenAllowed(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	alice := metadata.Actor{Name: "Alice", Email: "alice@example.test"}

	_, err := Apply(ctx, fixture.Driver, ApplyOptions{
		Actor:                          alice,
		Now:                            fixedTime(),
		AllowLegacyRelationshipAliases: true,
	}, Plan{Operations: []Operation{
		{Kind: OperationCreateNode, NodeType: identity.NodeTypeAsset, HID: "AST_1_1"},
		{Kind: OperationCreateNode, NodeType: identity.NodeTypeRegime, HID: "REG_1_1"},
		{
			Kind:             OperationCreateRelationship,
			RelationshipName: "Has",
			FromHID:          "AST_1_1",
			FromType:         identity.NodeTypeAsset,
			ToHID:            "REG_1_1",
			ToType:           identity.NodeTypeRegime,
		},
	}})
	if err != nil {
		t.Fatal(err)
	}

	session := fixture.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	value, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
MATCH (:Asset {HID: "AST_1_1"})-[relationship:HAS_REGIME]->(:Regime {HID: "REG_1_1"})
RETURN count(relationship) AS count
`, nil)
		if err != nil {
			return nil, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		count, _ := record.Get("count")
		return count, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if value != int64(1) {
		t.Fatalf("HAS_REGIME count = %#v, want 1", value)
	}
}

func TestApplyRenamesLegacyRequirementBaronPropertyWhenAllowed(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	alice := metadata.Actor{Name: "Alice", Email: "alice@example.test"}

	_, err := Apply(ctx, fixture.Driver, ApplyOptions{
		Actor:                      alice,
		Now:                        fixedTime(),
		AllowLegacyPropertyAliases: true,
	}, Plan{Operations: []Operation{{
		Kind:     OperationCreateNode,
		NodeType: identity.NodeTypeRequirement,
		HID:      "REQ_1_1",
		Properties: map[string]any{
			"Baron": false,
		},
	}}})
	if err != nil {
		t.Fatal(err)
	}

	session := fixture.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	value, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
MATCH (requirement:Requirement {HID: "REQ_1_1"})
RETURN requirement.Barren AS barren, requirement.Baron AS baron
`, nil)
		if err != nil {
			return nil, err
		}
		return result.Single(ctx)
	})
	if err != nil {
		t.Fatal(err)
	}
	record := value.(*neo4j.Record)
	barren, _ := record.Get("barren")
	baron, hasBaron := record.Get("baron")
	if barren != false || (hasBaron && baron != nil) {
		t.Fatalf("legacy property normalization failed: Barren=%#v Baron=%#v", barren, baron)
	}
}

func TestApplyDefaultsRequiredRelationshipProperties(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	alice := metadata.Actor{Name: "Alice", Email: "alice@example.test"}

	_, err := Apply(ctx, fixture.Driver, ApplyOptions{Actor: alice, Now: fixedTime()}, Plan{Operations: []Operation{
		{Kind: OperationCreateNode, NodeType: identity.NodeTypeState, HID: "ST_1_1"},
		{Kind: OperationCreateNode, NodeType: identity.NodeTypeState, HID: "ST_1_2"},
		{
			Kind:             OperationCreateRelationship,
			RelationshipName: "TRANSITIONS_TO",
			FromHID:          "ST_1_1",
			FromType:         identity.NodeTypeState,
			ToHID:            "ST_1_2",
			ToType:           identity.NodeTypeState,
		},
	}})
	if err != nil {
		t.Fatal(err)
	}

	session := fixture.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	value, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
MATCH (:State {HID: "ST_1_1"})-[relationship:TRANSITIONS_TO]->(:State {HID: "ST_1_2"})
RETURN relationship.TransitionKind AS transitionKind
`, nil)
		if err != nil {
			return nil, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		transitionKind, _ := record.Get("transitionKind")
		return transitionKind, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if value != "FUNCTIONAL" {
		t.Fatalf("TransitionKind = %#v, want FUNCTIONAL", value)
	}
}

func TestApplyRejectsDerivedAssetFromNonPrimaryTarget(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	alice := metadata.Actor{Name: "Alice", Email: "alice@example.test"}

	_, err := Apply(ctx, fixture.Driver, ApplyOptions{Actor: alice, Now: fixedTime()}, Plan{Operations: []Operation{
		{Kind: OperationCreateNode, NodeType: identity.NodeTypeAsset, HID: "AST_1_1", Properties: map[string]any{"IsPrimary": true}},
		{Kind: OperationCreateNode, NodeType: identity.NodeTypeAsset, HID: "AST_1_2", Properties: map[string]any{"IsPrimary": false}},
		{Kind: OperationCreateNode, NodeType: identity.NodeTypeAsset, HID: "AST_1_3", Properties: map[string]any{"IsPrimary": false}},
		{
			Kind:             OperationCreateRelationship,
			RelationshipName: "DERIVED_FROM",
			FromHID:          "AST_1_2",
			FromType:         identity.NodeTypeAsset,
			ToHID:            "AST_1_1",
			ToType:           identity.NodeTypeAsset,
		},
	}})
	if err != nil {
		t.Fatal(err)
	}

	_, err = Apply(ctx, fixture.Driver, ApplyOptions{Actor: alice, Now: fixedTime()}, Plan{Operations: []Operation{{
		Kind:             OperationCreateRelationship,
		RelationshipName: "DERIVED_FROM",
		FromHID:          "AST_1_3",
		FromType:         identity.NodeTypeAsset,
		ToHID:            "AST_1_2",
		ToType:           identity.NodeTypeAsset,
	}}})
	if err == nil {
		t.Fatal("expected DERIVED_FROM to reject a non-primary target")
	}
}

func TestApplyAutoGeneratesLossAndRootGoalForAssetCriticalityAssuranceEnvironment(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	alice := metadata.Actor{Name: "Alice", Email: "alice@example.test"}

	report, err := Apply(ctx, fixture.Driver, ApplyOptions{Actor: alice, Now: fixedTime(), VersionID: "v58-test"}, Plan{Operations: []Operation{
		{
			Kind:     OperationCreateNode,
			NodeType: identity.NodeTypeEnvironment,
			HID:      "ENV_1_1",
			Properties: map[string]any{
				"Name": "Operational",
			},
		},
		{
			Kind:     OperationCreateNode,
			NodeType: identity.NodeTypeAsset,
			HID:      "AST_1_1",
			Properties: map[string]any{
				"Name":             "Mission Data",
				"AssetType":        "PRIMARY",
				"IsPrimary":        true,
				"SafetyCritical":   true,
				"Confidentiality":  true,
				"Availability":     true,
				"NonRepudiation":   false,
				"MissionCritical":  false,
				"FlightCritical":   false,
				"SecurityCritical": false,
			},
		},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if !containsString(report.NodesChanged, "LOS_1_1") || !containsString(report.NodesChanged, "G_1_1") {
		t.Fatalf("generated Loss/Goal HIDs not reported in NodesChanged: %#v", report.NodesChanged)
	}

	assertLossGoalCounts(t, ctx, fixture.Driver, 2, 2)

	_, err = Apply(ctx, fixture.Driver, ApplyOptions{Actor: alice, Now: fixedTime().Add(time.Minute), VersionID: "v58-test"}, Plan{Operations: []Operation{{
		Kind: OperationUpdateNode,
		HID:  "AST_1_1",
		Properties: map[string]any{
			"ShortDescription": "still the same logical asset",
		},
	}}})
	if err != nil {
		t.Fatal(err)
	}

	assertLossGoalCounts(t, ctx, fixture.Driver, 2, 2)
}

func TestApplyRejectsGoalStructureCycle(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	alice := metadata.Actor{Name: "Alice", Email: "alice@example.test"}

	_, err := Apply(ctx, fixture.Driver, ApplyOptions{Actor: alice, Now: fixedTime()}, Plan{Operations: []Operation{
		{Kind: OperationCreateNode, NodeType: identity.NodeTypeGoal, HID: "G_1_1"},
		{Kind: OperationCreateNode, NodeType: identity.NodeTypeStrategy, HID: "SGY_1_1"},
		{
			Kind:             OperationCreateRelationship,
			RelationshipName: "SUPPORTED_BY",
			FromHID:          "G_1_1",
			FromType:         identity.NodeTypeGoal,
			ToHID:            "SGY_1_1",
			ToType:           identity.NodeTypeStrategy,
		},
	}})
	if err != nil {
		t.Fatal(err)
	}

	_, err = Apply(ctx, fixture.Driver, ApplyOptions{Actor: alice, Now: fixedTime()}, Plan{Operations: []Operation{{
		Kind:             OperationCreateRelationship,
		RelationshipName: "SUPPORTED_BY",
		FromHID:          "SGY_1_1",
		FromType:         identity.NodeTypeStrategy,
		ToHID:            "G_1_1",
		ToType:           identity.NodeTypeGoal,
	}}})
	if err == nil {
		t.Fatal("expected Goal Structure cycle to be rejected")
	}
}

func fixedTime() time.Time {
	return time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
}

func assertLossGoalCounts(t *testing.T, ctx context.Context, driver neo4j.DriverWithContext, wantLosses int64, wantGoals int64) {
	t.Helper()
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	value, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
MATCH (asset:Asset {HID: "AST_1_1"})
OPTIONAL MATCH (asset)-[:HAS_LOSS]->(loss:Loss)-[:HAS_ENVIRONMENT]->(:Environment {HID: "ENV_1_1"})
WITH asset, collect(DISTINCT loss) AS losses
OPTIONAL MATCH (asset)-[:HAS_GOAL]->(goal:Goal)
WHERE goal.RootForLossHID IN [loss IN losses | loss.HID]
RETURN size(losses) AS losses, count(DISTINCT goal) AS goals
`, nil)
		if err != nil {
			return nil, err
		}
		return result.Single(ctx)
	})
	if err != nil {
		t.Fatal(err)
	}

	record := value.(*neo4j.Record)
	losses, _ := record.Get("losses")
	goals, _ := record.Get("goals")
	if losses != wantLosses || goals != wantGoals {
		t.Fatalf("generated counts = losses:%#v goals:%#v, want losses:%d goals:%d", losses, goals, wantLosses, wantGoals)
	}
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func TestApplyRecordsTraceSpan(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	if err := schema.Bootstrap(ctx, fixture.Driver, ""); err != nil {
		t.Fatal(err)
	}

	recorder := tracetest.NewSpanRecorder()
	provider := telemetry.NewTestTracerProvider(recorder)
	SetTracer(provider.Tracer("mutation-test"))
	t.Cleanup(func() { SetTracer(nil) })

	actor := metadata.Actor{Name: "Alice", Email: "alice@example.test"}
	plan := Plan{Operations: []Operation{{
		Kind:       OperationCreateNode,
		NodeType:   identity.NodeTypeCapability,
		HID:        "CAP__1",
		UUID:       "00000000-0000-4000-8000-000000000900",
		Properties: map[string]any{"Name": "Root"},
	}}}

	if _, err := Apply(ctx, fixture.Driver, ApplyOptions{Actor: actor, VersionID: "v1"}, plan); err != nil {
		t.Fatalf("apply: %v", err)
	}

	spans := recorder.Ended()
	if len(spans) == 0 {
		t.Fatalf("expected a mutation span, got none")
	}
	found := false
	for _, span := range spans {
		if span.Name() == "sstpa.mutation.apply" {
			found = true
			attrs := span.Attributes()
			hasCommit := false
			for _, a := range attrs {
				if string(a.Key) == "sstpa.commit_id" && a.Value.AsString() != "" {
					hasCommit = true
				}
			}
			if !hasCommit {
				t.Fatalf("span missing sstpa.commit_id attribute")
			}
		}
	}
	if !found {
		t.Fatalf("sstpa.mutation.apply span not found; names: %v", spanNames(spans))
	}
}

func spanNames(spans []sdktrace.ReadOnlySpan) []string {
	names := make([]string, 0, len(spans))
	for _, span := range spans {
		names = append(names, span.Name())
	}
	return names
}
