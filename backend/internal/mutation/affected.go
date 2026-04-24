// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package mutation

import "sort"

type NodeSnapshot struct {
	HID        string
	Owner      string
	OwnerEmail string
	Properties map[string]any
}

type RelationshipSnapshot struct {
	Name    string
	FromHID string
	ToHID   string
}

type GraphSnapshot struct {
	Nodes         map[string]NodeSnapshot
	Relationships []RelationshipSnapshot
}

type AffectedNode struct {
	HID        string
	Owner      string
	OwnerEmail string
	Reasons    []string
}

func ComputeAffected(plan Plan, before GraphSnapshot, after GraphSnapshot) []AffectedNode {
	reasons := map[string]map[string]struct{}{}

	addReason := func(hid string, reason string) {
		if hid == "" {
			return
		}
		if _, ok := reasons[hid]; !ok {
			reasons[hid] = map[string]struct{}{}
		}
		reasons[hid][reason] = struct{}{}
	}

	for _, operation := range plan.Operations {
		switch operation.Kind {
		case OperationCreateNode:
			addReason(operation.HID, "node_created")
		case OperationUpdateNode:
			addReason(operation.HID, "node_updated")
		case OperationCreateRelationship:
			addReason(operation.FromHID, "relationship_created")
			addReason(operation.ToHID, "relationship_created")
		}
	}

	affected := make([]AffectedNode, 0, len(reasons))
	for hid, reasonSet := range reasons {
		snapshot := after.Nodes[hid]
		if snapshot.HID == "" {
			snapshot = before.Nodes[hid]
		}

		reasonList := make([]string, 0, len(reasonSet))
		for reason := range reasonSet {
			reasonList = append(reasonList, reason)
		}
		sort.Strings(reasonList)

		affected = append(affected, AffectedNode{
			HID:        hid,
			Owner:      snapshot.Owner,
			OwnerEmail: snapshot.OwnerEmail,
			Reasons:    reasonList,
		})
	}

	sort.Slice(affected, func(i, j int) bool {
		return affected[i].HID < affected[j].HID
	})

	return affected
}
