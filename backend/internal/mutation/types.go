// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package mutation

import (
	"fmt"
	"time"

	"sstpa-tool/backend/internal/graph"
	"sstpa-tool/backend/internal/identity"
	"sstpa-tool/backend/internal/metadata"
)

type OperationKind string

const (
	OperationCreateNode         OperationKind = "create_node"
	OperationUpdateNode         OperationKind = "update_node"
	OperationCreateRelationship OperationKind = "create_relationship"
)

type Plan struct {
	Operations []Operation
}

type Operation struct {
	Kind OperationKind

	NodeType   identity.NodeType
	HID        string
	UUID       string
	Properties map[string]any

	RelationshipName       string
	FromHID                string
	FromType               identity.NodeType
	ToHID                  string
	ToType                 identity.NodeType
	RelationshipProperties map[string]any
}

type CommitReport struct {
	CommitID             string   `json:"commitId"`
	NodesChanged         []string `json:"nodesChanged"`
	RelationshipsChanged []string `json:"relationshipsChanged"`
	MessagesGenerated    int      `json:"messagesGenerated"`
	RecipientsNotified   []string `json:"recipientsNotified"`
}

type ApplyOptions struct {
	DatabaseName                   string
	Actor                          metadata.Actor
	Now                            time.Time
	CommitID                       string
	VersionID                      string
	AllowLegacyRelationshipAliases bool
	AllowLegacyPropertyAliases     bool
}

func (plan Plan) Validate() error {
	return plan.ValidateWithOptions(ValidationOptions{})
}

type ValidationOptions struct {
	AllowLegacyRelationshipAliases bool
	AllowLegacyPropertyAliases     bool
}

func (plan Plan) ValidateWithOptions(options ValidationOptions) error {
	if len(plan.Operations) == 0 {
		return fmt.Errorf("mutation plan must include at least one operation")
	}

	for index, operation := range plan.Operations {
		if err := operation.ValidateWithOptions(options); err != nil {
			return fmt.Errorf("operation %d: %w", index, err)
		}
	}

	return nil
}

func (operation Operation) Validate() error {
	return operation.ValidateWithOptions(ValidationOptions{})
}

func (operation Operation) ValidateWithOptions(options ValidationOptions) error {
	switch operation.Kind {
	case OperationCreateNode:
		if _, ok := identity.TypeID(operation.NodeType); !ok {
			return fmt.Errorf("unknown node type %q", operation.NodeType)
		}
		if operation.HID == "" {
			return fmt.Errorf("HID is required for node creation")
		}
		if err := validateHIDMatchesNodeType(operation.HID, operation.NodeType); err != nil {
			return err
		}
		if err := validatePropertyAliases(operation.Properties, options); err != nil {
			return err
		}
	case OperationUpdateNode:
		if operation.HID == "" {
			return fmt.Errorf("HID is required for node update")
		}
		if _, _, _, err := identity.ParseHID(operation.HID); err != nil {
			return err
		}
		if len(operation.Properties) == 0 {
			return fmt.Errorf("properties are required for node update")
		}
		if err := validatePropertyAliases(operation.Properties, options); err != nil {
			return err
		}
	case OperationCreateRelationship:
		if operation.FromHID == "" || operation.ToHID == "" {
			return fmt.Errorf("from and to HIDs are required for relationship creation")
		}
		if err := validateHIDMatchesNodeType(operation.FromHID, operation.FromType); err != nil {
			return fmt.Errorf("from HID/type mismatch: %w", err)
		}
		if err := validateHIDMatchesNodeType(operation.ToHID, operation.ToType); err != nil {
			return fmt.Errorf("to HID/type mismatch: %w", err)
		}
		relationship, _, ok := graph.LookupRelationshipWithLegacyAliases(operation.RelationshipName, operation.FromType, operation.ToType, options.AllowLegacyRelationshipAliases)
		if !ok {
			return fmt.Errorf("relationship %s from %s to %s is not allowed", operation.RelationshipName, operation.FromType, operation.ToType)
		}
		props := graph.DefaultRelationshipProperties(relationship)
		for key, value := range operation.RelationshipProperties {
			props[key] = value
		}
		if err := graph.ValidateSoIBoundary(relationship, operation.FromHID, operation.ToHID, props); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown operation kind %q", operation.Kind)
	}

	return nil
}

func validateHIDMatchesNodeType(hid string, nodeType identity.NodeType) error {
	typeID, _, _, err := identity.ParseHID(hid)
	if err != nil {
		return err
	}

	want, ok := identity.TypeID(nodeType)
	if !ok {
		return fmt.Errorf("unknown node type %q", nodeType)
	}
	if typeID != want {
		return fmt.Errorf("HID %s has type id %s, want %s for %s", hid, typeID, want, nodeType)
	}

	return nil
}

func validatePropertyAliases(properties map[string]any, options ValidationOptions) error {
	if _, ok := properties["Baron"]; ok && !options.AllowLegacyPropertyAliases {
		return fmt.Errorf("legacy Requirement property Baron is not allowed in new writes; use Barren")
	}

	return nil
}
