// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package onboarding

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sstpa-tool/backend/internal/identity"
	"sstpa-tool/backend/internal/metadata"
)

type Kind struct {
	NodeType       identity.NodeType
	NodeLabel      string
	ContainerLabel string
	Relationship   string
}

var UserKind = Kind{
	NodeType:       identity.NodeTypeUser,
	NodeLabel:      "User",
	ContainerLabel: "Users",
	Relationship:   "HAS_USER",
}

var AdminKind = Kind{
	NodeType:       identity.NodeTypeAdmin,
	NodeLabel:      "Admin",
	ContainerLabel: "Admins",
	Relationship:   "HAS_ADMIN",
}

type Record struct {
	HID       string `json:"hid"`
	UUID      string `json:"uuid"`
	TypeName  string `json:"typeName"`
	UserName  string `json:"userName"`
	UserEmail string `json:"userEmail"`
	Created   string `json:"created"`
	LastTouch string `json:"lastTouch"`
}

type CreateInput struct {
	UserName  string
	UserEmail string
	Actor     metadata.Actor
	Now       time.Time
}

type ListResult struct {
	Items []Record `json:"items"`
	Page  int      `json:"page"`
	Limit int      `json:"limit"`
	Total int64    `json:"total"`
}

type Page struct {
	Page   int
	Limit  int
	Offset int
}

var ErrAlreadyRegistered = errors.New("user already registered")
var ErrNotFound = errors.New("not found")

func Create(ctx context.Context, driver neo4j.DriverWithContext, databaseName string, kind Kind, input CreateInput) (Record, error) {
	if driver == nil {
		return Record{}, errors.New("neo4j driver is required")
	}
	if input.UserName == "" || input.UserEmail == "" {
		return Record{}, errors.New("UserName and UserEmail are required")
	}
	if input.Actor.Name == "" || input.Actor.Email == "" {
		return Record{}, errors.New("actor name and email are required")
	}

	typeID, ok := identity.TypeID(kind.NodeType)
	if !ok {
		return Record{}, fmt.Errorf("unknown node type %q", kind.NodeType)
	}

	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: databaseName})
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		existing, err := scalarInt(ctx, tx,
			"MATCH (n:"+kind.NodeLabel+":SSTPANode {UserEmail: $email}) RETURN count(n) AS c",
			map[string]any{"email": input.UserEmail})
		if err != nil {
			return Record{}, err
		}
		if existing > 0 {
			return Record{}, ErrAlreadyRegistered
		}

		nextSeq, err := scalarInt(ctx, tx,
			"MATCH (n:"+kind.NodeLabel+":SSTPANode) RETURN coalesce(max(n.HIDSequence), 0) + 1 AS c",
			nil)
		if err != nil {
			return Record{}, err
		}

		hid, err := identity.FormatHID(typeID, "", int(nextSeq))
		if err != nil {
			return Record{}, err
		}

		uuid := identity.NewUUID()
		common, err := metadata.NewCommon(metadata.NewCommonInput{
			NodeType:  kind.NodeType,
			HID:       hid,
			UUID:      uuid,
			Actor:     input.Actor,
			Now:       now,
			VersionID: "",
		})
		if err != nil {
			return Record{}, err
		}

		props := common.Properties()
		props["Name"] = input.UserName
		props["UserName"] = input.UserName
		props["UserEmail"] = input.UserEmail
		props["UserHash"] = input.UserEmail
		props["HIDSequence"] = nextSeq

		createCypher := fmt.Sprintf(`
MATCH (container:%s)
CREATE (n:%s:SSTPANode)
SET n = $props
MERGE (container)-[:%s]->(n)
RETURN n.HID AS hid, n.uuid AS uuid, n.TypeName AS typeName,
       n.UserName AS userName, n.UserEmail AS userEmail,
       n.Created AS created, n.LastTouch AS lastTouch
`, kind.ContainerLabel, kind.NodeLabel, kind.Relationship)

		row, err := tx.Run(ctx, createCypher, map[string]any{"props": props})
		if err != nil {
			return Record{}, err
		}
		record, err := row.Single(ctx)
		if err != nil {
			return Record{}, err
		}
		return recordFromRow(record), nil
	})
	if err != nil {
		return Record{}, err
	}

	out, ok := result.(Record)
	if !ok {
		return Record{}, errors.New("unexpected onboarding result")
	}
	return out, nil
}

func List(ctx context.Context, driver neo4j.DriverWithContext, databaseName string, kind Kind, page Page) (ListResult, error) {
	if driver == nil {
		return ListResult{}, errors.New("neo4j driver is required")
	}
	if page.Limit <= 0 {
		page.Limit = 50
	}

	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: databaseName})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		rows, err := tx.Run(ctx, `
MATCH (n:`+kind.NodeLabel+`:SSTPANode)
RETURN n.HID AS hid, n.uuid AS uuid, n.TypeName AS typeName,
       n.UserName AS userName, n.UserEmail AS userEmail,
       n.Created AS created, n.LastTouch AS lastTouch
ORDER BY n.HID
SKIP $skip LIMIT $limit
`, map[string]any{"skip": page.Offset, "limit": page.Limit})
		if err != nil {
			return ListResult{}, err
		}
		collected, err := rows.Collect(ctx)
		if err != nil {
			return ListResult{}, err
		}

		items := make([]Record, 0, len(collected))
		for _, record := range collected {
			items = append(items, recordFromRow(record))
		}

		total, err := scalarInt(ctx, tx,
			"MATCH (n:"+kind.NodeLabel+":SSTPANode) RETURN count(n) AS c",
			nil)
		if err != nil {
			return ListResult{}, err
		}

		pageNumber := page.Page
		if pageNumber < 1 {
			pageNumber = 1
		}
		return ListResult{Items: items, Page: pageNumber, Limit: page.Limit, Total: total}, nil
	})
	if err != nil {
		return ListResult{}, err
	}

	out, ok := result.(ListResult)
	if !ok {
		return ListResult{}, errors.New("unexpected list result")
	}
	return out, nil
}

func GetByUUID(ctx context.Context, driver neo4j.DriverWithContext, databaseName string, kind Kind, uuid string) (Record, error) {
	if driver == nil {
		return Record{}, errors.New("neo4j driver is required")
	}
	if uuid == "" {
		return Record{}, errors.New("uuid is required")
	}

	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: databaseName})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		row, err := tx.Run(ctx, `
MATCH (n:`+kind.NodeLabel+`:SSTPANode {uuid: $uuid})
RETURN n.HID AS hid, n.uuid AS uuid, n.TypeName AS typeName,
       n.UserName AS userName, n.UserEmail AS userEmail,
       n.Created AS created, n.LastTouch AS lastTouch
LIMIT 1
`, map[string]any{"uuid": uuid})
		if err != nil {
			return Record{}, err
		}
		record, err := row.Single(ctx)
		if err != nil {
			return Record{}, ErrNotFound
		}
		return recordFromRow(record), nil
	})
	if err != nil {
		return Record{}, err
	}
	out, ok := result.(Record)
	if !ok {
		return Record{}, errors.New("unexpected get result")
	}
	return out, nil
}

func recordFromRow(record *neo4j.Record) Record {
	get := func(key string) string {
		value, _ := record.Get(key)
		text, _ := value.(string)
		return text
	}
	return Record{
		HID:       get("hid"),
		UUID:      get("uuid"),
		TypeName:  get("typeName"),
		UserName:  get("userName"),
		UserEmail: get("userEmail"),
		Created:   get("created"),
		LastTouch: get("lastTouch"),
	}
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
	value, ok := record.Get("c")
	if !ok {
		return 0, errors.New("c scalar not returned")
	}
	count, ok := value.(int64)
	if !ok {
		return 0, fmt.Errorf("c has unexpected type %T", value)
	}
	return count, nil
}
