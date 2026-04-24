// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package metadata

import (
	"fmt"
	"time"

	"sstpa-tool/backend/internal/identity"
)

const NullValue = "Null"

type Actor struct {
	Name  string
	Email string
	Admin bool
}

type Common struct {
	Name             string
	HID              string
	UUID             string
	TypeName         string
	Owner            string
	OwnerEmail       string
	Creator          string
	CreatorEmail     string
	Created          time.Time
	LastTouch        time.Time
	VersionID        string
	ShortDescription string
	LongDescription  string
}

type NewCommonInput struct {
	NodeType  identity.NodeType
	HID       string
	UUID      string
	Actor     Actor
	Now       time.Time
	VersionID string
}

func NewCommon(input NewCommonInput) (Common, error) {
	if _, ok := identity.TypeID(input.NodeType); !ok {
		return Common{}, fmt.Errorf("unknown node type %q", input.NodeType)
	}

	if input.HID == "" {
		return Common{}, fmt.Errorf("HID is required")
	}

	if input.UUID == "" {
		return Common{}, fmt.Errorf("uuid is required")
	}

	if input.Actor.Name == "" || input.Actor.Email == "" {
		return Common{}, fmt.Errorf("actor name and email are required")
	}

	if input.Now.IsZero() {
		return Common{}, fmt.Errorf("timestamp is required")
	}

	return Common{
		Name:             "New",
		HID:              input.HID,
		UUID:             input.UUID,
		TypeName:         string(input.NodeType),
		Owner:            input.Actor.Name,
		OwnerEmail:       input.Actor.Email,
		Creator:          input.Actor.Name,
		CreatorEmail:     input.Actor.Email,
		Created:          input.Now.UTC(),
		LastTouch:        input.Now.UTC(),
		VersionID:        defaultString(input.VersionID),
		ShortDescription: NullValue,
		LongDescription:  NullValue,
	}, nil
}

func (common Common) Properties() map[string]any {
	return map[string]any{
		"Name":             nullIfEmpty(common.Name),
		"HID":              common.HID,
		"uuid":             common.UUID,
		"TypeName":         common.TypeName,
		"Owner":            common.Owner,
		"OwnerEmail":       common.OwnerEmail,
		"Creator":          common.Creator,
		"CreatorEmail":     common.CreatorEmail,
		"Created":          common.Created.Format(time.RFC3339),
		"LastTouch":        common.LastTouch.Format(time.RFC3339),
		"VersionID":        nullIfEmpty(common.VersionID),
		"ShortDescription": nullIfEmpty(common.ShortDescription),
		"LongDescription":  nullIfEmpty(common.LongDescription),
	}
}

func defaultString(value string) string {
	if value == "" {
		return NullValue
	}

	return value
}

func nullIfEmpty(value string) string {
	if value == "" {
		return NullValue
	}

	return value
}
