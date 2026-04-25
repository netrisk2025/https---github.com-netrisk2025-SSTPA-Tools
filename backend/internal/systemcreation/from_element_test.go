// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package systemcreation

import (
	"context"
	"testing"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sstpa-tool/backend/internal/identity"
	"sstpa-tool/backend/internal/metadata"
	"sstpa-tool/backend/internal/mutation"
	"sstpa-tool/backend/internal/testhelpers"
)

func TestCreateFromElementBuildsChildSoIAndCopiesRequirementsAssets(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	actor := metadata.Actor{Name: "Alice", Email: "alice@example.test"}
	seedParentElement(t, ctx, fixture.Driver, actor)

	result, err := CreateFromElement(ctx, fixture.Driver, FromElementOptions{
		Actor:     actor,
		Now:       fixedTime(),
		VersionID: "v58-test",
	}, "EL_1_4")
	if err != nil {
		t.Fatal(err)
	}

	if result.SystemHID != "SYS_1.4_0" ||
		result.PurposeHID != "PUR_1.4_1" ||
		result.EnvironmentHID != "ENV_1.4_1" ||
		result.StateHID != "ST_1.4_1" ||
		result.SecurityHID != "SEC_1.4_1" ||
		result.FunctionalFlowHID != "FF_1.4_1" {
		t.Fatalf("unexpected child scaffold HIDs: %#v", result)
	}
	if len(result.RequirementHIDs) != 3 {
		t.Fatalf("copied requirements = %#v, want 3", result.RequirementHIDs)
	}
	if len(result.AssetHIDs) != 2 {
		t.Fatalf("copied assets = %#v, want 2", result.AssetHIDs)
	}

	session := fixture.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	value, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
MATCH (:Element {HID: "EL_1_4"})-[:PARENTS]->(:System {HID: "SYS_1.4_0"})
MATCH (:System {HID: "SYS_1.4_0"})-[:REALIZES]->(:Purpose {HID: "PUR_1.4_1"})
MATCH (:System {HID: "SYS_1.4_0"})-[:ACTS_IN]->(:Environment {HID: "ENV_1.4_1"})
MATCH (:System {HID: "SYS_1.4_0"})-[:EXHIBITS]->(:State {HID: "ST_1.4_1"})
MATCH (:System {HID: "SYS_1.4_0"})-[:HAS_SECURITY]->(:Security {HID: "SEC_1.4_1"})
MATCH (:System {HID: "SYS_1.4_0"})-[:HAS_FUNCTIONAL_FLOW]->(:FunctionalFlow {HID: "FF_1.4_1"})
OPTIONAL MATCH (:Purpose {HID: "PUR_1.4_1"})-[:HAS_REQUIREMENT]->(requirement:Requirement)
WITH count(DISTINCT requirement) AS requirements
OPTIONAL MATCH (:System {HID: "SYS_1.4_0"})-[:HAS_ASSET]->(asset:Asset)
WITH requirements, count(DISTINCT asset) AS assets
OPTIONAL MATCH (:Asset)-[:HAS_LOSS]->(loss:Loss)-[:HAS_ENVIRONMENT]->(:Environment {HID: "ENV_1.4_1"})
WHERE loss.SourceAssetHID IN ["AST_1.4_1", "AST_1.4_2"]
WITH requirements, assets, count(DISTINCT loss) AS losses
OPTIONAL MATCH (:Asset)-[:HAS_LOSS]->(lossForGoal:Loss)
WHERE lossForGoal.SourceAssetHID IN ["AST_1.4_1", "AST_1.4_2"]
WITH requirements, assets, losses, collect(DISTINCT lossForGoal.HID) AS lossHIDs
OPTIONAL MATCH (:Asset)-[:HAS_GOAL]->(goal:Goal)
WHERE goal.RootForLossHID IN lossHIDs
RETURN requirements, assets, losses, count(DISTINCT goal) AS goals
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
	requirements, _ := record.Get("requirements")
	assets, _ := record.Get("assets")
	losses, _ := record.Get("losses")
	goals, _ := record.Get("goals")
	if requirements != int64(3) || assets != int64(2) || losses != int64(2) || goals != int64(2) {
		t.Fatalf("child SoI counts requirements=%#v assets=%#v losses=%#v goals=%#v", requirements, assets, losses, goals)
	}
}

func TestCreateFromElementRejectsSecondChildSystem(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	actor := metadata.Actor{Name: "Alice", Email: "alice@example.test"}
	seedParentElement(t, ctx, fixture.Driver, actor)

	_, err := CreateFromElement(ctx, fixture.Driver, FromElementOptions{Actor: actor, Now: fixedTime()}, "EL_1_4")
	if err != nil {
		t.Fatal(err)
	}

	_, err = CreateFromElement(ctx, fixture.Driver, FromElementOptions{Actor: actor, Now: fixedTime()}, "EL_1_4")
	if err == nil {
		t.Fatal("expected second child System creation to be rejected")
	}
}

func seedParentElement(t *testing.T, ctx context.Context, driver neo4j.DriverWithContext, actor metadata.Actor) {
	t.Helper()
	_, err := mutation.Apply(ctx, driver, mutation.ApplyOptions{Actor: actor, Now: fixedTime(), VersionID: "v58-test"}, mutation.Plan{Operations: []mutation.Operation{
		{Kind: mutation.OperationCreateNode, NodeType: identity.NodeTypeElement, HID: "EL_1_4", Properties: map[string]any{"Name": "Guidance Computer"}},
		{Kind: mutation.OperationCreateNode, NodeType: identity.NodeTypeFunction, HID: "FUN_1_1", Properties: map[string]any{"Name": "Compute Guidance"}},
		{Kind: mutation.OperationCreateNode, NodeType: identity.NodeTypeInterface, HID: "INT_1_1", Properties: map[string]any{"Name": "Guidance Bus"}},
		{Kind: mutation.OperationCreateNode, NodeType: identity.NodeTypeRequirement, HID: "REQ_1_1", Properties: map[string]any{"Name": "Element requirement", "RStatement": "Element shall preserve command integrity."}},
		{Kind: mutation.OperationCreateNode, NodeType: identity.NodeTypeRequirement, HID: "REQ_1_2", Properties: map[string]any{"Name": "Function requirement", "RStatement": "Function shall validate commands."}},
		{Kind: mutation.OperationCreateNode, NodeType: identity.NodeTypeRequirement, HID: "REQ_1_3", Properties: map[string]any{"Name": "Interface requirement", "RStatement": "Interface shall authenticate traffic."}},
		{Kind: mutation.OperationCreateNode, NodeType: identity.NodeTypeAsset, HID: "AST_1_1", Properties: assetProperties("Mission Data")},
		{Kind: mutation.OperationCreateNode, NodeType: identity.NodeTypeAsset, HID: "AST_1_2", Properties: assetProperties("Command Key")},
		{Kind: mutation.OperationCreateRelationship, RelationshipName: "ALLOCATED_TO", FromType: identity.NodeTypeFunction, FromHID: "FUN_1_1", ToType: identity.NodeTypeElement, ToHID: "EL_1_4"},
		{Kind: mutation.OperationCreateRelationship, RelationshipName: "ALLOCATED_TO", FromType: identity.NodeTypeInterface, FromHID: "INT_1_1", ToType: identity.NodeTypeElement, ToHID: "EL_1_4"},
		{Kind: mutation.OperationCreateRelationship, RelationshipName: "HAS_REQUIREMENT", FromType: identity.NodeTypeElement, FromHID: "EL_1_4", ToType: identity.NodeTypeRequirement, ToHID: "REQ_1_1"},
		{Kind: mutation.OperationCreateRelationship, RelationshipName: "HAS_REQUIREMENT", FromType: identity.NodeTypeFunction, FromHID: "FUN_1_1", ToType: identity.NodeTypeRequirement, ToHID: "REQ_1_2"},
		{Kind: mutation.OperationCreateRelationship, RelationshipName: "HAS_REQUIREMENT", FromType: identity.NodeTypeInterface, FromHID: "INT_1_1", ToType: identity.NodeTypeRequirement, ToHID: "REQ_1_3"},
		{Kind: mutation.OperationCreateRelationship, RelationshipName: "CONTAINS", FromType: identity.NodeTypeElement, FromHID: "EL_1_4", ToType: identity.NodeTypeAsset, ToHID: "AST_1_1"},
		{Kind: mutation.OperationCreateRelationship, RelationshipName: "CONTAINS", FromType: identity.NodeTypeInterface, FromHID: "INT_1_1", ToType: identity.NodeTypeAsset, ToHID: "AST_1_2"},
	}})
	if err != nil {
		t.Fatal(err)
	}
}

func assetProperties(name string) map[string]any {
	return map[string]any{
		"Name":             name,
		"AssetType":        "PRIMARY",
		"IsPrimary":        true,
		"SafetyCritical":   true,
		"MissionCritical":  false,
		"FlightCritical":   false,
		"SecurityCritical": false,
		"Confidentiality":  true,
		"Availability":     false,
		"Authenticity":     false,
		"NonRepudiation":   false,
		"Durability":       false,
		"Privacy":          false,
		"Trustworthy":      false,
	}
}

func fixedTime() time.Time {
	return time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
}
