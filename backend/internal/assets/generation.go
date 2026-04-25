// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package assets

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sstpa-tool/backend/internal/graph"
	"sstpa-tool/backend/internal/identity"
	"sstpa-tool/backend/internal/metadata"
)

type GenerationInput struct {
	Actor     metadata.Actor
	Now       time.Time
	VersionID string
	AssetHIDs []string
}

type GenerationReport struct {
	CreatedNodes         []GeneratedNode
	CreatedRelationships []GeneratedRelationship
}

type GeneratedNode struct {
	HID      string
	NodeType identity.NodeType
}

type GeneratedRelationship struct {
	Name    string
	FromHID string
	ToHID   string
}

type namedFlag struct {
	Name     string
	Property string
}

var criticalityFlags = []namedFlag{
	{Name: "Safety", Property: "SafetyCritical"},
	{Name: "Mission", Property: "MissionCritical"},
	{Name: "Flight", Property: "FlightCritical"},
	{Name: "Security", Property: "SecurityCritical"},
}

var assuranceFlags = []namedFlag{
	{Name: "Confidentiality", Property: "Confidentiality"},
	{Name: "Availability", Property: "Availability"},
	{Name: "Authenticity", Property: "Authenticity"},
	{Name: "NonRepudiation", Property: "NonRepudiation"},
	{Name: "Durability", Property: "Durability"},
	{Name: "Privacy", Property: "Privacy"},
	{Name: "Trustworthy", Property: "Trustworthy"},
}

type assetRecord struct {
	HID        string
	Name       string
	Index      string
	Properties map[string]any
}

type environmentRecord struct {
	HID  string
	Name string
}

func EnsureLossGoalPairs(ctx context.Context, tx neo4j.ManagedTransaction, input GenerationInput) (GenerationReport, error) {
	if len(input.AssetHIDs) == 0 {
		return GenerationReport{}, nil
	}
	if input.Actor.Name == "" || input.Actor.Email == "" {
		return GenerationReport{}, fmt.Errorf("actor name and email are required")
	}
	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	report := GenerationReport{}
	sequenceCache := map[sequenceKey]int{}
	for _, assetHID := range uniqueStrings(input.AssetHIDs) {
		asset, ok, err := loadAsset(ctx, tx, assetHID)
		if err != nil {
			return GenerationReport{}, err
		}
		if !ok {
			return GenerationReport{}, fmt.Errorf("asset %s was not found", assetHID)
		}

		criticalities := selectedFlags(asset.Properties, criticalityFlags)
		assurances := selectedFlags(asset.Properties, assuranceFlags)
		if len(criticalities) == 0 || len(assurances) == 0 {
			continue
		}

		environments, err := loadEnvironmentsForIndex(ctx, tx, asset.Index)
		if err != nil {
			return GenerationReport{}, err
		}
		for _, environment := range environments {
			for _, criticality := range criticalities {
				for _, assurance := range assurances {
					created, err := ensurePair(ctx, tx, input, now, asset, environment, criticality, assurance, sequenceCache)
					if err != nil {
						return GenerationReport{}, err
					}
					report.CreatedNodes = append(report.CreatedNodes, created.CreatedNodes...)
					report.CreatedRelationships = append(report.CreatedRelationships, created.CreatedRelationships...)
				}
			}
		}
	}

	return report, nil
}

func ensurePair(
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	input GenerationInput,
	now time.Time,
	asset assetRecord,
	environment environmentRecord,
	criticality namedFlag,
	assurance namedFlag,
	sequenceCache map[sequenceKey]int,
) (GenerationReport, error) {
	report := GenerationReport{}
	loss, createdLoss, err := ensureLoss(ctx, tx, input, now, asset, environment, criticality, assurance, sequenceCache)
	if err != nil {
		return GenerationReport{}, err
	}
	if createdLoss {
		report.CreatedNodes = append(report.CreatedNodes, GeneratedNode{HID: loss.HID, NodeType: identity.NodeTypeLoss})
	}

	if created, err := ensureRelationship(ctx, tx, "HAS_LOSS", identity.NodeTypeAsset, asset.HID, identity.NodeTypeLoss, loss.HID); err != nil {
		return GenerationReport{}, err
	} else if created {
		report.CreatedRelationships = append(report.CreatedRelationships, GeneratedRelationship{Name: "HAS_LOSS", FromHID: asset.HID, ToHID: loss.HID})
	}
	if created, err := ensureRelationship(ctx, tx, "HAS_ENVIRONMENT", identity.NodeTypeLoss, loss.HID, identity.NodeTypeEnvironment, environment.HID); err != nil {
		return GenerationReport{}, err
	} else if created {
		report.CreatedRelationships = append(report.CreatedRelationships, GeneratedRelationship{Name: "HAS_ENVIRONMENT", FromHID: loss.HID, ToHID: environment.HID})
	}

	goal, createdGoal, err := ensureRootGoal(ctx, tx, input, now, asset, environment, loss, criticality, assurance, sequenceCache)
	if err != nil {
		return GenerationReport{}, err
	}
	if createdGoal {
		report.CreatedNodes = append(report.CreatedNodes, GeneratedNode{HID: goal.HID, NodeType: identity.NodeTypeGoal})
	}
	if created, err := ensureRelationship(ctx, tx, "HAS_GOAL", identity.NodeTypeAsset, asset.HID, identity.NodeTypeGoal, goal.HID); err != nil {
		return GenerationReport{}, err
	} else if created {
		report.CreatedRelationships = append(report.CreatedRelationships, GeneratedRelationship{Name: "HAS_GOAL", FromHID: asset.HID, ToHID: goal.HID})
	}

	if err := setPairBacklinks(ctx, tx, loss.HID, goal.HID); err != nil {
		return GenerationReport{}, err
	}

	return report, nil
}

type generatedLoss struct {
	HID  string
	Name string
}

type generatedGoal struct {
	HID string
}

func ensureLoss(
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	input GenerationInput,
	now time.Time,
	asset assetRecord,
	environment environmentRecord,
	criticality namedFlag,
	assurance namedFlag,
	sequenceCache map[sequenceKey]int,
) (generatedLoss, bool, error) {
	key := lossKey(asset.HID, environment.HID, criticality.Name, assurance.Name)
	result, err := tx.Run(ctx, `
MATCH (loss:Loss:SSTPANode {LossKey: $lossKey})
RETURN loss.HID AS hid, loss.Name AS name
`, map[string]any{"assetHID": asset.HID, "lossKey": key})
	if err != nil {
		return generatedLoss{}, false, err
	}
	record, err := result.Single(ctx)
	if err == nil {
		hid, _ := record.Get("hid")
		name, _ := record.Get("name")
		return generatedLoss{HID: stringValue(hid), Name: stringValue(name)}, false, nil
	}
	if !isNoRecord(err) {
		return generatedLoss{}, false, err
	}

	hid, err := nextHID(ctx, tx, identity.NodeTypeLoss, asset.Index, sequenceCache)
	if err != nil {
		return generatedLoss{}, false, err
	}
	name := fmt.Sprintf("Loss: %s %s/%s in %s", asset.Name, criticality.Name, assurance.Name, environment.Name)
	props, err := commonProperties(identity.NodeTypeLoss, hid, name, input.Actor, now, input.VersionID)
	if err != nil {
		return generatedLoss{}, false, err
	}
	props["LossKey"] = key
	props["SourceAssetHID"] = asset.HID
	props["EnvironmentHID"] = environment.HID
	props["Criticality"] = criticality.Name
	props["Assurance"] = assurance.Name
	props["AttackTreeFormat"] = "SSTPA-ATF-1.0"
	props["AttackTreeVersion"] = int64(0)
	props["AttackTreeCreated"] = now.UTC().Format(time.RFC3339)
	props["AttackTreeLastModified"] = now.UTC().Format(time.RFC3339)
	props["AttackTreeCreatedBy"] = input.Actor.Name
	props["AttackTreeCreatedByEmail"] = input.Actor.Email
	props["AttackTreeStatus"] = "AUTO_GENERATED"
	props["AttackTreeJSON"] = metadata.NullValue
	for _, flag := range criticalityFlags {
		props[flag.Property] = flag.Property == criticality.Property
	}
	for _, flag := range assuranceFlags {
		props[flag.Property] = flag.Property == assurance.Property
	}

	if err := createGeneratedNode(ctx, tx, identity.NodeTypeLoss, props); err != nil {
		return generatedLoss{}, false, err
	}

	return generatedLoss{HID: hid, Name: name}, true, nil
}

func ensureRootGoal(
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	input GenerationInput,
	now time.Time,
	asset assetRecord,
	environment environmentRecord,
	loss generatedLoss,
	criticality namedFlag,
	assurance namedFlag,
	sequenceCache map[sequenceKey]int,
) (generatedGoal, bool, error) {
	result, err := tx.Run(ctx, `
MATCH (goal:Goal:SSTPANode {RootForLossHID: $lossHID})
RETURN goal.HID AS hid
`, map[string]any{"assetHID": asset.HID, "lossHID": loss.HID})
	if err != nil {
		return generatedGoal{}, false, err
	}
	record, err := result.Single(ctx)
	if err == nil {
		hid, _ := record.Get("hid")
		return generatedGoal{HID: stringValue(hid)}, false, nil
	}
	if !isNoRecord(err) {
		return generatedGoal{}, false, err
	}

	hid, err := nextHID(ctx, tx, identity.NodeTypeGoal, asset.Index, sequenceCache)
	if err != nil {
		return generatedGoal{}, false, err
	}
	props, err := commonProperties(identity.NodeTypeGoal, hid, fmt.Sprintf("Root Goal: %s", loss.Name), input.Actor, now, input.VersionID)
	if err != nil {
		return generatedGoal{}, false, err
	}
	props["IsRootGoal"] = true
	props["AssetHID"] = asset.HID
	props["LossHID"] = loss.HID
	props["RootForLossHID"] = loss.HID
	props["EnvironmentHID"] = environment.HID
	props["Criticality"] = criticality.Name
	props["Assurance"] = assurance.Name
	props["GSNID"] = "G-ROOT"
	props["GoalStructure"] = metadata.NullValue
	props["GoalStatement"] = fmt.Sprintf(
		"The evidence supports certification that %s maintains %s for %s in %s such that %s is acceptably mitigated.",
		asset.Name,
		assurance.Name,
		criticality.Name,
		environment.Name,
		loss.Name,
	)

	if err := createGeneratedNode(ctx, tx, identity.NodeTypeGoal, props); err != nil {
		return generatedGoal{}, false, err
	}

	return generatedGoal{HID: hid}, true, nil
}

func loadAsset(ctx context.Context, tx neo4j.ManagedTransaction, hid string) (assetRecord, bool, error) {
	_, index, _, err := identity.ParseHID(hid)
	if err != nil {
		return assetRecord{}, false, err
	}

	result, err := tx.Run(ctx, `
MATCH (asset:Asset:SSTPANode {HID: $hid})
RETURN asset.HID AS hid, asset.Name AS name, properties(asset) AS properties
`, map[string]any{"hid": hid})
	if err != nil {
		return assetRecord{}, false, err
	}
	record, err := result.Single(ctx)
	if err != nil {
		if isNoRecord(err) {
			return assetRecord{}, false, nil
		}
		return assetRecord{}, false, err
	}

	props, _ := record.Get("properties")
	properties, _ := props.(map[string]any)
	name, _ := record.Get("name")
	return assetRecord{HID: hid, Name: defaultName(stringValue(name), hid), Index: index, Properties: properties}, true, nil
}

func loadEnvironmentsForIndex(ctx context.Context, tx neo4j.ManagedTransaction, index string) ([]environmentRecord, error) {
	prefix, err := hidPrefix(identity.NodeTypeEnvironment, index)
	if err != nil {
		return nil, err
	}

	result, err := tx.Run(ctx, `
MATCH (environment:Environment:SSTPANode)
WHERE environment.HID STARTS WITH $prefix
RETURN environment.HID AS hid, environment.Name AS name
ORDER BY environment.HID
`, map[string]any{"prefix": prefix})
	if err != nil {
		return nil, err
	}
	records, err := result.Collect(ctx)
	if err != nil {
		return nil, err
	}

	environments := make([]environmentRecord, 0, len(records))
	for _, record := range records {
		hid, _ := record.Get("hid")
		name, _ := record.Get("name")
		hidText := stringValue(hid)
		environments = append(environments, environmentRecord{
			HID:  hidText,
			Name: defaultName(stringValue(name), hidText),
		})
	}

	return environments, nil
}

func selectedFlags(properties map[string]any, flags []namedFlag) []namedFlag {
	selected := []namedFlag{}
	for _, flag := range flags {
		if truthy(properties[flag.Property]) {
			selected = append(selected, flag)
		}
	}

	return selected
}

func commonProperties(nodeType identity.NodeType, hid string, name string, actor metadata.Actor, now time.Time, versionID string) (map[string]any, error) {
	common, err := metadata.NewCommon(metadata.NewCommonInput{
		NodeType:  nodeType,
		HID:       hid,
		UUID:      identity.NewUUID(),
		Actor:     actor,
		Now:       now,
		VersionID: versionID,
	})
	if err != nil {
		return nil, err
	}

	props := common.Properties()
	props["Name"] = name
	return props, nil
}

func createGeneratedNode(ctx context.Context, tx neo4j.ManagedTransaction, nodeType identity.NodeType, props map[string]any) error {
	label, ok := graph.LabelFor(nodeType)
	if !ok {
		return fmt.Errorf("unknown node type %q", nodeType)
	}

	result, err := tx.Run(ctx, fmt.Sprintf(`
CREATE (node:%s:SSTPANode)
SET node = $props
RETURN node.HID AS hid
`, label), map[string]any{"props": props})
	if err != nil {
		return err
	}
	_, err = result.Single(ctx)
	return err
}

func ensureRelationship(
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	name string,
	fromType identity.NodeType,
	fromHID string,
	toType identity.NodeType,
	toHID string,
) (bool, error) {
	relationship, ok := graph.LookupRelationship(name, fromType, toType)
	if !ok {
		return false, fmt.Errorf("relationship %s from %s to %s is not allowed", name, fromType, toType)
	}
	props := graph.DefaultRelationshipProperties(relationship)
	if err := graph.ValidateRelationshipProperties(relationship, props); err != nil {
		return false, err
	}
	if err := graph.ValidateSoIBoundary(relationship, fromHID, toHID, props); err != nil {
		return false, err
	}

	fromLabel, _ := graph.LabelFor(fromType)
	toLabel, _ := graph.LabelFor(toType)
	existsQuery := fmt.Sprintf(`
MATCH (from:%s:SSTPANode {HID: $fromHID})-[relationship:%s]->(to:%s:SSTPANode {HID: $toHID})
RETURN count(relationship) AS count
`, fromLabel, name, toLabel)
	count, err := countQuery(ctx, tx, existsQuery, map[string]any{"fromHID": fromHID, "toHID": toHID})
	if err != nil {
		return false, err
	}
	if count > 0 {
		return false, nil
	}

	createQuery := fmt.Sprintf(`
MATCH (from:%s:SSTPANode {HID: $fromHID}), (to:%s:SSTPANode {HID: $toHID})
CREATE (from)-[relationship:%s]->(to)
SET relationship = $props
RETURN type(relationship) AS relationshipType
`, fromLabel, toLabel, name)
	result, err := tx.Run(ctx, createQuery, map[string]any{"fromHID": fromHID, "toHID": toHID, "props": props})
	if err != nil {
		return false, err
	}
	_, err = result.Single(ctx)
	return true, err
}

func setPairBacklinks(ctx context.Context, tx neo4j.ManagedTransaction, lossHID string, goalHID string) error {
	result, err := tx.Run(ctx, `
MATCH (loss:Loss:SSTPANode {HID: $lossHID}), (goal:Goal:SSTPANode {HID: $goalHID})
SET loss.RootGoalHID = goal.HID,
    goal.LossHID = loss.HID
RETURN loss.HID AS lossHID
`, map[string]any{"lossHID": lossHID, "goalHID": goalHID})
	if err != nil {
		return err
	}
	_, err = result.Single(ctx)
	return err
}

type sequenceKey struct {
	NodeType identity.NodeType
	Index    string
}

func nextHID(ctx context.Context, tx neo4j.ManagedTransaction, nodeType identity.NodeType, index string, cache map[sequenceKey]int) (string, error) {
	key := sequenceKey{NodeType: nodeType, Index: index}
	next, ok := cache[key]
	if !ok {
		maxSequence, err := maxSequence(ctx, tx, nodeType, index)
		if err != nil {
			return "", err
		}
		next = maxSequence + 1
	}
	cache[key] = next + 1

	typeID, _ := identity.TypeID(nodeType)
	return identity.FormatHID(typeID, index, next)
}

func maxSequence(ctx context.Context, tx neo4j.ManagedTransaction, nodeType identity.NodeType, index string) (int, error) {
	prefix, err := hidPrefix(nodeType, index)
	if err != nil {
		return 0, err
	}
	label, _ := graph.LabelFor(nodeType)
	result, err := tx.Run(ctx, fmt.Sprintf(`
MATCH (node:%s:SSTPANode)
WHERE node.HID STARTS WITH $prefix
RETURN node.HID AS hid
`, label), map[string]any{"prefix": prefix})
	if err != nil {
		return 0, err
	}
	records, err := result.Collect(ctx)
	if err != nil {
		return 0, err
	}

	maximum := 0
	for _, record := range records {
		value, _ := record.Get("hid")
		_, hidIndex, sequence, err := identity.ParseHID(stringValue(value))
		if err != nil || hidIndex != index {
			continue
		}
		if sequence > maximum {
			maximum = sequence
		}
	}

	return maximum, nil
}

func hidPrefix(nodeType identity.NodeType, index string) (string, error) {
	typeID, ok := identity.TypeID(nodeType)
	if !ok {
		return "", fmt.Errorf("unknown node type %q", nodeType)
	}

	return fmt.Sprintf("%s_%s_", typeID, index), nil
}

func countQuery(ctx context.Context, tx neo4j.ManagedTransaction, query string, params map[string]any) (int64, error) {
	result, err := tx.Run(ctx, query, params)
	if err != nil {
		return 0, err
	}
	record, err := result.Single(ctx)
	if err != nil {
		return 0, err
	}
	value, ok := record.Get("count")
	if !ok {
		return 0, fmt.Errorf("count not returned")
	}
	count, ok := value.(int64)
	if !ok {
		return 0, fmt.Errorf("count has unexpected type %T", value)
	}
	return count, nil
}

func lossKey(assetHID string, environmentHID string, criticality string, assurance string) string {
	return strings.Join([]string{assetHID, environmentHID, criticality, assurance}, "|")
}

func uniqueStrings(values []string) []string {
	seen := map[string]struct{}{}
	unique := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		unique = append(unique, value)
	}
	return unique
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

func truthy(value any) bool {
	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		return strings.EqualFold(typed, "true")
	default:
		return false
	}
}

func isNoRecord(err error) bool {
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "no records") || strings.Contains(message, "no more records")
}
