// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package mutation

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"sstpa-tool/backend/internal/graph"
	"sstpa-tool/backend/internal/identity"
	"sstpa-tool/backend/internal/messaging"
	"sstpa-tool/backend/internal/metadata"
)

var mutationTracer trace.Tracer

// SetTracer installs a tracer used to open an sstpa.mutation.apply span around
// every Apply call's write transaction. Passing nil disables tracing.
func SetTracer(tracer trace.Tracer) { mutationTracer = tracer }

func Apply(ctx context.Context, driver neo4j.DriverWithContext, options ApplyOptions, plan Plan) (CommitReport, error) {
	if err := plan.Validate(); err != nil {
		return CommitReport{}, err
	}

	if options.Actor.Name == "" || options.Actor.Email == "" {
		return CommitReport{}, fmt.Errorf("actor name and email are required")
	}

	now := options.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}
	commitID := options.CommitID
	if commitID == "" {
		commitID = identity.NewUUID()
	}

	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: options.DatabaseName})
	defer session.Close(ctx)

	ctxForWrite := ctx
	var span trace.Span
	if mutationTracer != nil {
		ctxForWrite, span = mutationTracer.Start(ctx, "sstpa.mutation.apply")
		defer span.End()
	}

	result, err := session.ExecuteWrite(ctxForWrite, func(tx neo4j.ManagedTransaction) (any, error) {
		before, err := readSnapshot(ctx, tx, preExistingHIDs(plan))
		if err != nil {
			return CommitReport{}, err
		}

		for _, operation := range plan.Operations {
			if err := applyOperation(ctx, tx, options, now, operation); err != nil {
				return CommitReport{}, err
			}
		}

		allHIDs := allOperationHIDs(plan)
		after, err := readSnapshot(ctx, tx, allHIDs)
		if err != nil {
			return CommitReport{}, err
		}

		affected := ComputeAffected(plan, before, after)
		report := CommitReport{
			CommitID:             commitID,
			NodesChanged:         affectedHIDs(affected),
			RelationshipsChanged: relationshipChanges(plan),
		}

		recipients, err := notifyAffectedOwners(ctx, tx, options.Actor, commitID, now, affected, before, after, report.RelationshipsChanged)
		if err != nil {
			return CommitReport{}, err
		}

		report.RecipientsNotified = recipients
		report.MessagesGenerated = len(recipients)
		return report, nil
	})

	if span != nil {
		span.SetAttributes(
			attribute.String("sstpa.commit_id", commitID),
			attribute.Int("sstpa.operations_count", len(plan.Operations)),
			attribute.String("sstpa.actor_email", options.Actor.Email),
		)
		if err != nil {
			span.RecordError(err)
		} else if report, ok := result.(CommitReport); ok {
			span.SetAttributes(
				attribute.Int("sstpa.messages_generated", report.MessagesGenerated),
				attribute.Int("sstpa.nodes_changed", len(report.NodesChanged)),
			)
		}
	}

	if err != nil {
		return CommitReport{}, err
	}

	report, ok := result.(CommitReport)
	if !ok {
		return CommitReport{}, fmt.Errorf("unexpected commit report result")
	}

	return report, nil
}

func applyOperation(ctx context.Context, tx neo4j.ManagedTransaction, options ApplyOptions, now time.Time, operation Operation) error {
	switch operation.Kind {
	case OperationCreateNode:
		return createNode(ctx, tx, options, now, operation)
	case OperationUpdateNode:
		return updateNode(ctx, tx, options.Actor, now, operation)
	case OperationCreateRelationship:
		return createRelationship(ctx, tx, operation)
	default:
		return fmt.Errorf("unsupported operation kind %q", operation.Kind)
	}
}

func createNode(ctx context.Context, tx neo4j.ManagedTransaction, options ApplyOptions, now time.Time, operation Operation) error {
	label, ok := graph.LabelFor(operation.NodeType)
	if !ok {
		return fmt.Errorf("unknown node type %q", operation.NodeType)
	}

	uuid := operation.UUID
	if uuid == "" {
		uuid = identity.NewUUID()
	}

	common, err := metadata.NewCommon(metadata.NewCommonInput{
		NodeType:  operation.NodeType,
		HID:       operation.HID,
		UUID:      uuid,
		Actor:     options.Actor,
		Now:       now,
		VersionID: options.VersionID,
	})
	if err != nil {
		return err
	}

	props := common.Properties()
	for key, value := range operation.Properties {
		if isCreationProtectedProperty(key) {
			return fmt.Errorf("property %s is assigned by the mutation layer", key)
		}
		props[key] = normalizePropertyValue(value)
	}

	result, err := tx.Run(ctx, fmt.Sprintf("CREATE (n:%s:SSTPANode) SET n = $props RETURN n.HID AS HID", label), map[string]any{"props": props})
	if err != nil {
		return err
	}

	_, err = result.Single(ctx)
	return err
}

func updateNode(ctx context.Context, tx neo4j.ManagedTransaction, actor metadata.Actor, now time.Time, operation Operation) error {
	props, err := validatedUpdateProperties(actor, now, operation.Properties)
	if err != nil {
		return err
	}

	result, err := tx.Run(ctx, `
MATCH (n {HID: $hid})
SET n += $props
RETURN n.HID AS HID
`, map[string]any{"hid": operation.HID, "props": props})
	if err != nil {
		return err
	}

	_, err = result.Single(ctx)
	return err
}

func createRelationship(ctx context.Context, tx neo4j.ManagedTransaction, operation Operation) error {
	relationship, ok := graph.LookupRelationship(operation.RelationshipName, operation.FromType, operation.ToType)
	if !ok {
		return fmt.Errorf("relationship %s from %s to %s is not allowed", operation.RelationshipName, operation.FromType, operation.ToType)
	}

	if relationship.Recursion == graph.RecursionDAG {
		if err := rejectCycle(ctx, tx, operation); err != nil {
			return err
		}
	}

	existsQuery := fmt.Sprintf(`
MATCH (from {HID: $fromHID})-[r:%s]->(to {HID: $toHID})
RETURN count(r) AS count
`, relationship.Name)
	existing, err := scalarInt(ctx, tx, existsQuery, map[string]any{"fromHID": operation.FromHID, "toHID": operation.ToHID})
	if err != nil {
		return err
	}
	if existing > 0 {
		return fmt.Errorf("duplicate relationship %s from %s to %s", relationship.Name, operation.FromHID, operation.ToHID)
	}

	props := map[string]any{}
	for key, value := range operation.RelationshipProperties {
		props[key] = normalizePropertyValue(value)
	}

	createQuery := fmt.Sprintf(`
MATCH (from {HID: $fromHID}), (to {HID: $toHID})
CREATE (from)-[r:%s]->(to)
SET r = $props
RETURN type(r) AS type
`, relationship.Name)
	result, err := tx.Run(ctx, createQuery, map[string]any{
		"fromHID": operation.FromHID,
		"toHID":   operation.ToHID,
		"props":   props,
	})
	if err != nil {
		return err
	}

	_, err = result.Single(ctx)
	return err
}

func rejectCycle(ctx context.Context, tx neo4j.ManagedTransaction, operation Operation) error {
	query := fmt.Sprintf(`
MATCH (from {HID: $fromHID}), (to {HID: $toHID})
OPTIONAL MATCH path = (to)-[:%s*1..50]->(from)
RETURN count(path) AS count
`, operation.RelationshipName)
	count, err := scalarInt(ctx, tx, query, map[string]any{"fromHID": operation.FromHID, "toHID": operation.ToHID})
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("relationship %s would introduce a cycle", operation.RelationshipName)
	}

	return nil
}

func scalarInt(ctx context.Context, tx neo4j.ManagedTransaction, query string, params map[string]any) (int64, error) {
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

func validatedUpdateProperties(actor metadata.Actor, now time.Time, input map[string]any) (map[string]any, error) {
	props := map[string]any{}
	for key, value := range input {
		switch key {
		case "HID", "uuid", "TypeName", "Created", "VersionID":
			return nil, fmt.Errorf("property %s is fixed", key)
		case "Creator", "CreatorEmail":
			if !actor.Admin {
				return nil, fmt.Errorf("property %s is immutable except for Admin", key)
			}
		}
		props[key] = normalizePropertyValue(value)
	}

	owner, hasOwner := props["Owner"].(string)
	ownerEmail, hasOwnerEmail := props["OwnerEmail"].(string)
	if hasOwner || hasOwnerEmail {
		if !hasOwner || !hasOwnerEmail {
			return nil, fmt.Errorf("Owner and OwnerEmail must be changed as a pair")
		}
		if !actor.Admin && (owner != actor.Name || ownerEmail != actor.Email) {
			return nil, fmt.Errorf("non-admin users may only assume ownership for themselves")
		}
	}

	props["LastTouch"] = now.UTC().Format(time.RFC3339)
	return props, nil
}

func isCreationProtectedProperty(key string) bool {
	switch key {
	case "HID", "uuid", "TypeName", "Owner", "OwnerEmail", "Creator", "CreatorEmail", "Created", "LastTouch", "VersionID":
		return true
	default:
		return false
	}
}

func normalizePropertyValue(value any) any {
	if text, ok := value.(string); ok && text == "" {
		return metadata.NullValue
	}

	return value
}

func readSnapshot(ctx context.Context, tx neo4j.ManagedTransaction, hids []string) (GraphSnapshot, error) {
	snapshot := GraphSnapshot{Nodes: map[string]NodeSnapshot{}}
	if len(hids) == 0 {
		return snapshot, nil
	}

	result, err := tx.Run(ctx, `
MATCH (n)
WHERE n.HID IN $hids
RETURN n.HID AS HID, n.Owner AS Owner, n.OwnerEmail AS OwnerEmail, properties(n) AS Properties
`, map[string]any{"hids": uniqueStrings(hids)})
	if err != nil {
		return snapshot, err
	}

	records, err := result.Collect(ctx)
	if err != nil {
		return snapshot, err
	}

	for _, record := range records {
		hid, _ := record.Get("HID")
		owner, _ := record.Get("Owner")
		ownerEmail, _ := record.Get("OwnerEmail")
		props, _ := record.Get("Properties")
		hidText, _ := hid.(string)
		if hidText == "" {
			continue
		}

		properties, _ := props.(map[string]any)
		snapshot.Nodes[hidText] = NodeSnapshot{
			HID:        hidText,
			Owner:      stringValue(owner),
			OwnerEmail: stringValue(ownerEmail),
			Properties: properties,
		}
	}

	return snapshot, nil
}

func notifyAffectedOwners(
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	actor metadata.Actor,
	commitID string,
	now time.Time,
	affected []AffectedNode,
	before GraphSnapshot,
	after GraphSnapshot,
	relationshipTypes []string,
) ([]string, error) {
	type recipientChange struct {
		Actor metadata.Actor
		HIDs  []string
	}

	recipients := map[string]recipientChange{}
	addRecipient := func(owner string, ownerEmail string, hid string) error {
		if owner == "" || ownerEmail == "" {
			return fmt.Errorf("affected node %s is missing Owner/OwnerEmail", hid)
		}
		if owner == actor.Name && ownerEmail == actor.Email {
			return nil
		}

		change := recipients[ownerEmail]
		change.Actor = metadata.Actor{Name: owner, Email: ownerEmail}
		change.HIDs = append(change.HIDs, hid)
		recipients[ownerEmail] = change
		return nil
	}

	for _, node := range affected {
		beforeNode := before.Nodes[node.HID]
		afterNode := after.Nodes[node.HID]
		if beforeNode.HID != "" && afterNode.HID != "" && beforeNode.OwnerEmail != afterNode.OwnerEmail {
			if err := addRecipient(beforeNode.Owner, beforeNode.OwnerEmail, node.HID); err != nil {
				return nil, err
			}
		}

		owner := node.Owner
		ownerEmail := node.OwnerEmail
		if owner == "" && afterNode.HID != "" {
			owner = afterNode.Owner
			ownerEmail = afterNode.OwnerEmail
		}
		if err := addRecipient(owner, ownerEmail, node.HID); err != nil {
			return nil, err
		}
	}

	emails := make([]string, 0, len(recipients))
	for email, change := range recipients {
		hids := uniqueStrings(change.HIDs)
		sort.Strings(hids)
		_, err := messaging.AppendChangeNotification(ctx, tx, messaging.ChangeNotification{
			CommitID:                 commitID,
			Subject:                  "SSTPA data changed",
			Body:                     fmt.Sprintf("%s changed SSTPA data affecting %s.", actor.Name, strings.Join(hids, ", ")),
			SentAt:                   now,
			Sender:                   actor,
			Recipient:                change.Actor,
			RelatedNodeHIDs:          hids,
			RelatedRelationshipTypes: relationshipTypes,
			ChangeTypeSummary:        "mutation commit",
		})
		if err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	sort.Strings(emails)
	return emails, nil
}

func preExistingHIDs(plan Plan) []string {
	var hids []string
	for _, operation := range plan.Operations {
		switch operation.Kind {
		case OperationUpdateNode:
			hids = append(hids, operation.HID)
		case OperationCreateRelationship:
			hids = append(hids, operation.FromHID, operation.ToHID)
		}
	}
	return uniqueStrings(hids)
}

func allOperationHIDs(plan Plan) []string {
	var hids []string
	for _, operation := range plan.Operations {
		switch operation.Kind {
		case OperationCreateNode:
			hids = append(hids, operation.HID)
		case OperationUpdateNode:
			hids = append(hids, operation.HID)
		case OperationCreateRelationship:
			hids = append(hids, operation.FromHID, operation.ToHID)
		}
	}
	return uniqueStrings(hids)
}

func affectedHIDs(affected []AffectedNode) []string {
	hids := make([]string, 0, len(affected))
	for _, node := range affected {
		hids = append(hids, node.HID)
	}
	sort.Strings(hids)
	return hids
}

func relationshipChanges(plan Plan) []string {
	var names []string
	for _, operation := range plan.Operations {
		if operation.Kind == OperationCreateRelationship {
			names = append(names, operation.RelationshipName)
		}
	}
	return uniqueStrings(names)
}

func uniqueStrings(values []string) []string {
	seen := map[string]struct{}{}
	var unique []string
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
	sort.Strings(unique)
	return unique
}

func stringValue(value any) string {
	text, _ := value.(string)
	return text
}
