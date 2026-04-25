// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package systemcreation

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sstpa-tool/backend/internal/identity"
	"sstpa-tool/backend/internal/metadata"
	"sstpa-tool/backend/internal/mutation"
)

type FromElementOptions struct {
	DatabaseName string
	Actor        metadata.Actor
	Now          time.Time
	CommitID     string
	VersionID    string
}

type FromElementResult struct {
	mutation.CommitReport
	SystemHID         string   `json:"systemHid"`
	PurposeHID        string   `json:"purposeHid"`
	EnvironmentHID    string   `json:"environmentHid"`
	StateHID          string   `json:"stateHid"`
	SecurityHID       string   `json:"securityHid"`
	FunctionalFlowHID string   `json:"functionalFlowHid"`
	RequirementHIDs   []string `json:"requirementHids"`
	AssetHIDs         []string `json:"assetHids"`
}

type sourceNode struct {
	HID        string
	Name       string
	Properties map[string]any
}

type sourceSnapshot struct {
	Element      sourceNode
	Requirements []sourceNode
	Assets       []sourceNode
}

func CreateFromElement(ctx context.Context, driver neo4j.DriverWithContext, options FromElementOptions, elementHID string) (FromElementResult, error) {
	if elementHID == "" {
		return FromElementResult{}, fmt.Errorf("element HID is required")
	}
	if options.Actor.Name == "" || options.Actor.Email == "" {
		return FromElementResult{}, fmt.Errorf("actor name and email are required")
	}

	snapshot, err := loadSourceSnapshot(ctx, driver, options.DatabaseName, elementHID)
	if err != nil {
		return FromElementResult{}, err
	}

	plan, result, err := buildPlan(snapshot)
	if err != nil {
		return FromElementResult{}, err
	}

	report, err := mutation.Apply(ctx, driver, mutation.ApplyOptions{
		DatabaseName: options.DatabaseName,
		Actor:        options.Actor,
		Now:          options.Now,
		CommitID:     options.CommitID,
		VersionID:    options.VersionID,
	}, plan)
	if err != nil {
		return FromElementResult{}, err
	}

	result.CommitReport = report
	return result, nil
}

func loadSourceSnapshot(ctx context.Context, driver neo4j.DriverWithContext, databaseName string, elementHID string) (sourceSnapshot, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: databaseName})
	defer session.Close(ctx)

	value, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		element, err := loadElement(ctx, tx, elementHID)
		if err != nil {
			return nil, err
		}
		if err := rejectExistingChildSystem(ctx, tx, elementHID); err != nil {
			return nil, err
		}
		requirements, err := loadRelatedNodes(ctx, tx, elementHID, relatedRequirementsQuery)
		if err != nil {
			return nil, err
		}
		assets, err := loadRelatedNodes(ctx, tx, elementHID, relatedAssetsQuery)
		if err != nil {
			return nil, err
		}

		return sourceSnapshot{Element: element, Requirements: requirements, Assets: assets}, nil
	})
	if err != nil {
		return sourceSnapshot{}, err
	}

	snapshot, ok := value.(sourceSnapshot)
	if !ok {
		return sourceSnapshot{}, fmt.Errorf("unexpected source snapshot result")
	}

	return snapshot, nil
}

func loadElement(ctx context.Context, tx neo4j.ManagedTransaction, elementHID string) (sourceNode, error) {
	result, err := tx.Run(ctx, `
MATCH (element:Element:SSTPANode {HID: $elementHID})
RETURN element.HID AS hid, element.Name AS name, properties(element) AS properties
`, map[string]any{"elementHID": elementHID})
	if err != nil {
		return sourceNode{}, err
	}
	record, err := result.Single(ctx)
	if err != nil {
		return sourceNode{}, fmt.Errorf("element %s was not found", elementHID)
	}

	props, _ := record.Get("properties")
	properties, _ := props.(map[string]any)
	name, _ := record.Get("name")
	return sourceNode{HID: elementHID, Name: defaultName(stringValue(name), elementHID), Properties: properties}, nil
}

func rejectExistingChildSystem(ctx context.Context, tx neo4j.ManagedTransaction, elementHID string) error {
	result, err := tx.Run(ctx, `
MATCH (:Element:SSTPANode {HID: $elementHID})-[:PARENTS]->(system:System:SSTPANode)
RETURN count(system) AS count
`, map[string]any{"elementHID": elementHID})
	if err != nil {
		return err
	}
	record, err := result.Single(ctx)
	if err != nil {
		return err
	}
	value, _ := record.Get("count")
	count, _ := value.(int64)
	if count > 0 {
		return fmt.Errorf("element %s already parents a child System", elementHID)
	}

	return nil
}

const relatedRequirementsQuery = `
MATCH (element:Element:SSTPANode {HID: $elementHID})
OPTIONAL MATCH (element)-[:HAS_REQUIREMENT]->(direct:Requirement:SSTPANode)
WITH element, collect(direct) AS nodes
OPTIONAL MATCH (function:Function:SSTPANode)-[:ALLOCATED_TO]->(element)
OPTIONAL MATCH (function)-[:HAS_REQUIREMENT]->(functionRequirement:Requirement:SSTPANode)
WITH element, nodes + collect(functionRequirement) AS nodes
OPTIONAL MATCH (interface:Interface:SSTPANode)-[:ALLOCATED_TO]->(element)
OPTIONAL MATCH (interface)-[:HAS_REQUIREMENT]->(interfaceRequirement:Requirement:SSTPANode)
WITH nodes + collect(interfaceRequirement) AS nodes
UNWIND nodes AS node
WITH DISTINCT node
WHERE node IS NOT NULL
RETURN node.HID AS hid, node.Name AS name, properties(node) AS properties
ORDER BY hid
`

const relatedAssetsQuery = `
MATCH (element:Element:SSTPANode {HID: $elementHID})
OPTIONAL MATCH (element)-[:CONTAINS]->(direct:Asset:SSTPANode)
WITH element, collect(direct) AS nodes
OPTIONAL MATCH (function:Function:SSTPANode)-[:ALLOCATED_TO]->(element)
OPTIONAL MATCH (function)-[:CONTAINS]->(functionAsset:Asset:SSTPANode)
WITH element, nodes + collect(functionAsset) AS nodes
OPTIONAL MATCH (interface:Interface:SSTPANode)-[:ALLOCATED_TO]->(element)
OPTIONAL MATCH (interface)-[:CONTAINS]->(interfaceAsset:Asset:SSTPANode)
WITH nodes + collect(interfaceAsset) AS nodes
UNWIND nodes AS node
WITH DISTINCT node
WHERE node IS NOT NULL
RETURN node.HID AS hid, node.Name AS name, properties(node) AS properties
ORDER BY hid
`

func loadRelatedNodes(ctx context.Context, tx neo4j.ManagedTransaction, elementHID string, query string) ([]sourceNode, error) {
	result, err := tx.Run(ctx, query, map[string]any{"elementHID": elementHID})
	if err != nil {
		return nil, err
	}
	records, err := result.Collect(ctx)
	if err != nil {
		return nil, err
	}

	nodes := make([]sourceNode, 0, len(records))
	for _, record := range records {
		hid, _ := record.Get("hid")
		name, _ := record.Get("name")
		props, _ := record.Get("properties")
		properties, _ := props.(map[string]any)
		hidText := stringValue(hid)
		nodes = append(nodes, sourceNode{
			HID:        hidText,
			Name:       defaultName(stringValue(name), hidText),
			Properties: properties,
		})
	}

	return nodes, nil
}

func buildPlan(snapshot sourceSnapshot) (mutation.Plan, FromElementResult, error) {
	_, parentIndex, parentSequence, err := identity.ParseHID(snapshot.Element.HID)
	if err != nil {
		return mutation.Plan{}, FromElementResult{}, err
	}
	childIndex := childSystemIndex(parentIndex, parentSequence)

	systemHID, err := formatHID(identity.NodeTypeSystem, childIndex, 0)
	if err != nil {
		return mutation.Plan{}, FromElementResult{}, err
	}
	purposeHID, err := formatHID(identity.NodeTypePurpose, childIndex, 1)
	if err != nil {
		return mutation.Plan{}, FromElementResult{}, err
	}
	environmentHID, err := formatHID(identity.NodeTypeEnvironment, childIndex, 1)
	if err != nil {
		return mutation.Plan{}, FromElementResult{}, err
	}
	stateHID, err := formatHID(identity.NodeTypeState, childIndex, 1)
	if err != nil {
		return mutation.Plan{}, FromElementResult{}, err
	}
	securityHID, err := formatHID(identity.NodeTypeSecurity, childIndex, 1)
	if err != nil {
		return mutation.Plan{}, FromElementResult{}, err
	}
	functionalFlowHID, err := formatHID(identity.NodeTypeFunctionalFlow, childIndex, 1)
	if err != nil {
		return mutation.Plan{}, FromElementResult{}, err
	}

	result := FromElementResult{
		SystemHID:         systemHID,
		PurposeHID:        purposeHID,
		EnvironmentHID:    environmentHID,
		StateHID:          stateHID,
		SecurityHID:       securityHID,
		FunctionalFlowHID: functionalFlowHID,
	}
	operations := []mutation.Operation{
		createNode(identity.NodeTypeSystem, systemHID, map[string]any{"Name": "Child System of " + snapshot.Element.Name}),
		createNode(identity.NodeTypePurpose, purposeHID, map[string]any{"Name": "Purpose for " + snapshot.Element.Name}),
		createNode(identity.NodeTypeEnvironment, environmentHID, map[string]any{"Name": "Default Environment"}),
		createNode(identity.NodeTypeState, stateHID, map[string]any{"Name": "Default State"}),
		createNode(identity.NodeTypeSecurity, securityHID, map[string]any{"Name": "Security for " + snapshot.Element.Name}),
		createNode(identity.NodeTypeFunctionalFlow, functionalFlowHID, map[string]any{"Name": "Functional Flow for " + snapshot.Element.Name}),
		createRelationship("PARENTS", identity.NodeTypeElement, snapshot.Element.HID, identity.NodeTypeSystem, systemHID),
		createRelationship("REALIZES", identity.NodeTypeSystem, systemHID, identity.NodeTypePurpose, purposeHID),
		createRelationship("ACTS_IN", identity.NodeTypeSystem, systemHID, identity.NodeTypeEnvironment, environmentHID),
		createRelationship("EXHIBITS", identity.NodeTypeSystem, systemHID, identity.NodeTypeState, stateHID),
		createRelationship("HAS_SECURITY", identity.NodeTypeSystem, systemHID, identity.NodeTypeSecurity, securityHID),
		createRelationship("HAS_FUNCTIONAL_FLOW", identity.NodeTypeSystem, systemHID, identity.NodeTypeFunctionalFlow, functionalFlowHID),
	}

	for index, requirement := range snapshot.Requirements {
		hid, err := formatHID(identity.NodeTypeRequirement, childIndex, index+1)
		if err != nil {
			return mutation.Plan{}, FromElementResult{}, err
		}
		props := copiedProperties(requirement.Properties)
		props["Name"] = requirement.Name
		operations = append(operations,
			createNode(identity.NodeTypeRequirement, hid, props),
			createRelationship("HAS_REQUIREMENT", identity.NodeTypePurpose, purposeHID, identity.NodeTypeRequirement, hid),
		)
		result.RequirementHIDs = append(result.RequirementHIDs, hid)
	}

	for index, asset := range snapshot.Assets {
		hid, err := formatHID(identity.NodeTypeAsset, childIndex, index+1)
		if err != nil {
			return mutation.Plan{}, FromElementResult{}, err
		}
		props := copiedProperties(asset.Properties)
		props["Name"] = asset.Name
		operations = append(operations,
			createNode(identity.NodeTypeAsset, hid, props),
			createRelationship("HAS_ASSET", identity.NodeTypeSystem, systemHID, identity.NodeTypeAsset, hid),
		)
		result.AssetHIDs = append(result.AssetHIDs, hid)
	}

	return mutation.Plan{Operations: operations}, result, nil
}

func createNode(nodeType identity.NodeType, hid string, properties map[string]any) mutation.Operation {
	return mutation.Operation{
		Kind:       mutation.OperationCreateNode,
		NodeType:   nodeType,
		HID:        hid,
		Properties: properties,
	}
}

func createRelationship(name string, fromType identity.NodeType, fromHID string, toType identity.NodeType, toHID string) mutation.Operation {
	return mutation.Operation{
		Kind:             mutation.OperationCreateRelationship,
		RelationshipName: name,
		FromType:         fromType,
		FromHID:          fromHID,
		ToType:           toType,
		ToHID:            toHID,
	}
}

func childSystemIndex(parentIndex string, parentSequence int) string {
	sequence := strconv.Itoa(parentSequence)
	if parentIndex == "" {
		return sequence
	}

	return parentIndex + "." + sequence
}

func formatHID(nodeType identity.NodeType, index string, sequence int) (string, error) {
	typeID, ok := identity.TypeID(nodeType)
	if !ok {
		return "", fmt.Errorf("unknown node type %q", nodeType)
	}

	return identity.FormatHID(typeID, index, sequence)
}

func copiedProperties(source map[string]any) map[string]any {
	protected := map[string]struct{}{
		"HID": {}, "uuid": {}, "TypeName": {}, "Owner": {}, "OwnerEmail": {},
		"Creator": {}, "CreatorEmail": {}, "Created": {}, "LastTouch": {}, "VersionID": {},
	}
	target := map[string]any{}
	keys := make([]string, 0, len(source))
	for key := range source {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if _, ok := protected[key]; ok {
			continue
		}
		target[key] = source[key]
	}

	return target
}

func stringValue(value any) string {
	text, _ := value.(string)
	return text
}

func defaultName(name string, fallback string) string {
	if name == "" || name == metadata.NullValue {
		return fallback
	}

	return name
}
