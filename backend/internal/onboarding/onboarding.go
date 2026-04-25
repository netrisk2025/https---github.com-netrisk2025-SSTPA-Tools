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
	"sstpa-tool/backend/internal/schema"
)

type Kind struct {
	NodeType         identity.NodeType
	NodeLabel        string
	RegistryNodeType identity.NodeType
	RegistryLabel    string
	RegistryHID      string
	Relationship     string
	SequenceProperty string
}

var UserKind = Kind{
	NodeType:         identity.NodeTypeUser,
	NodeLabel:        "User",
	RegistryNodeType: identity.NodeTypeUserRegistry,
	RegistryLabel:    "UserRegistry",
	RegistryHID:      schema.UserRegistryHID,
	Relationship:     "HAS_USER",
	SequenceProperty: "UserSequence",
}

var AdminKind = Kind{
	NodeType:         identity.NodeTypeAdmin,
	NodeLabel:        "Admin",
	RegistryNodeType: identity.NodeTypeAdminRegistry,
	RegistryLabel:    "AdminRegistry",
	RegistryHID:      schema.AdminRegistryHID,
	Relationship:     "HAS_ADMIN",
	SequenceProperty: "AdminSequence",
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

type BootstrapInput struct {
	InstallerName  string
	InstallerEmail string
	Now            time.Time
}

type BootstrapResult struct {
	Admin Record `json:"admin"`
	User  Record `json:"user"`
}

type ListResult struct {
	Items []Record `json:"items"`
	Page  int      `json:"page"`
	Limit int      `json:"limit"`
	Total int64    `json:"total"`
}

type Page struct {
	Page  int
	Limit int
}

var ErrAlreadyRegistered = errors.New("already registered")
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

	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: databaseName})
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return createWithTx(ctx, tx, kind, CreateInput{
			UserName:  input.UserName,
			UserEmail: input.UserEmail,
			Actor:     input.Actor,
			Now:       now,
		})
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

func BootstrapInstaller(ctx context.Context, driver neo4j.DriverWithContext, databaseName string, input BootstrapInput) (BootstrapResult, error) {
	if driver == nil {
		return BootstrapResult{}, errors.New("neo4j driver is required")
	}
	if input.InstallerName == "" || input.InstallerEmail == "" {
		return BootstrapResult{}, errors.New("installerName and installerEmail are required")
	}

	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}
	actor := metadata.Actor{Name: input.InstallerName, Email: input.InstallerEmail, Admin: true}

	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: databaseName})
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		userCount, err := registeredCount(ctx, tx, UserKind)
		if err != nil {
			return BootstrapResult{}, err
		}
		adminCount, err := registeredCount(ctx, tx, AdminKind)
		if err != nil {
			return BootstrapResult{}, err
		}
		if userCount > 0 || adminCount > 0 {
			return BootstrapResult{}, ErrAlreadyRegistered
		}

		admin, err := createWithTx(ctx, tx, AdminKind, CreateInput{
			UserName:  input.InstallerName,
			UserEmail: input.InstallerEmail,
			Actor:     actor,
			Now:       now,
		})
		if err != nil {
			return BootstrapResult{}, err
		}
		user, err := createWithTx(ctx, tx, UserKind, CreateInput{
			UserName:  input.InstallerName,
			UserEmail: input.InstallerEmail,
			Actor:     actor,
			Now:       now,
		})
		if err != nil {
			return BootstrapResult{}, err
		}

		return BootstrapResult{Admin: admin, User: user}, nil
	})
	if err != nil {
		return BootstrapResult{}, err
	}

	out, ok := result.(BootstrapResult)
	if !ok {
		return BootstrapResult{}, errors.New("unexpected bootstrap result")
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
	if page.Page < 1 {
		page.Page = 1
	}
	offset := (page.Page - 1) * page.Limit

	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: databaseName})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		rows, err := tx.Run(ctx, `
MATCH (registry:`+kind.RegistryLabel+`:SSTPANode {HID: $registryHID})-[:`+kind.Relationship+`]->(n:`+kind.NodeLabel+`:SSTPANode)
RETURN n.HID AS hid, n.uuid AS uuid, n.TypeName AS typeName,
       n.UserName AS userName, n.UserEmail AS userEmail,
       n.Created AS created, n.LastTouch AS lastTouch
ORDER BY n.HID
SKIP $skip LIMIT $limit
`, map[string]any{"registryHID": kind.RegistryHID, "skip": offset, "limit": page.Limit})
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
			"MATCH (registry:"+kind.RegistryLabel+":SSTPANode {HID: $registryHID})-[:"+kind.Relationship+"]->(n:"+kind.NodeLabel+":SSTPANode) RETURN count(n) AS c",
			map[string]any{"registryHID": kind.RegistryHID})
		if err != nil {
			return ListResult{}, err
		}

		return ListResult{Items: items, Page: page.Page, Limit: page.Limit, Total: total}, nil
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
MATCH (registry:`+kind.RegistryLabel+`:SSTPANode {HID: $registryHID})-[:`+kind.Relationship+`]->(n:`+kind.NodeLabel+`:SSTPANode {uuid: $uuid})
RETURN n.HID AS hid, n.uuid AS uuid, n.TypeName AS typeName,
       n.UserName AS userName, n.UserEmail AS userEmail,
       n.Created AS created, n.LastTouch AS lastTouch
LIMIT 1
`, map[string]any{"registryHID": kind.RegistryHID, "uuid": uuid})
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

func IsRegisteredAdmin(ctx context.Context, driver neo4j.DriverWithContext, databaseName string, actor metadata.Actor) (bool, error) {
	if driver == nil {
		return false, errors.New("neo4j driver is required")
	}
	if actor.Email == "" {
		return false, nil
	}

	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: databaseName})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		count, err := scalarInt(ctx, tx,
			"MATCH (registry:AdminRegistry:SSTPANode {HID: $registryHID})-[:HAS_ADMIN]->(n:Admin:SSTPANode {UserEmail: $email}) RETURN count(n) AS c",
			map[string]any{"registryHID": AdminKind.RegistryHID, "email": actor.Email})
		if err != nil {
			return false, err
		}
		return count > 0, nil
	})
	if err != nil {
		return false, err
	}
	out, ok := result.(bool)
	if !ok {
		return false, errors.New("unexpected admin lookup result")
	}
	return out, nil
}

func createWithTx(ctx context.Context, tx neo4j.ManagedTransaction, kind Kind, input CreateInput) (Record, error) {
	typeID, ok := identity.TypeID(kind.NodeType)
	if !ok {
		return Record{}, fmt.Errorf("unknown node type %q", kind.NodeType)
	}
	if _, ok := identity.TypeID(kind.RegistryNodeType); !ok {
		return Record{}, fmt.Errorf("unknown registry node type %q", kind.RegistryNodeType)
	}

	existing, err := registeredEmailCount(ctx, tx, kind, input.UserEmail)
	if err != nil {
		return Record{}, err
	}
	if existing > 0 {
		return Record{}, ErrAlreadyRegistered
	}

	nextSeq, err := allocateSequence(ctx, tx, kind)
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
		Now:       input.Now,
		VersionID: "",
	})
	if err != nil {
		return Record{}, err
	}

	props := common.Properties()
	props["Name"] = input.UserName
	props["UserName"] = input.UserName
	props["UserEmail"] = input.UserEmail
	// TODO(sstpa-auth): compute salted hash once the auth layer lands;
	// SRS §1.4.4 permits email as an interim equivalent identifier.
	props["UserHash"] = input.UserEmail
	props["HIDSequence"] = nextSeq

	createCypher := fmt.Sprintf(`
MATCH (registry:%s:SSTPANode {HID: $registryHID})
CREATE (n:%s:SSTPANode)
SET n = $props
MERGE (registry)-[:%s]->(n)
RETURN n.HID AS hid, n.uuid AS uuid, n.TypeName AS typeName,
       n.UserName AS userName, n.UserEmail AS userEmail,
       n.Created AS created, n.LastTouch AS lastTouch
`, kind.RegistryLabel, kind.NodeLabel, kind.Relationship)

	row, err := tx.Run(ctx, createCypher, map[string]any{"registryHID": kind.RegistryHID, "props": props})
	if err != nil {
		return Record{}, err
	}
	record, err := row.Single(ctx)
	if err != nil {
		return Record{}, err
	}
	return recordFromRow(record), nil
}

func allocateSequence(ctx context.Context, tx neo4j.ManagedTransaction, kind Kind) (int64, error) {
	query := fmt.Sprintf(`
MATCH (tool:SSTPA_Tool:SSTPANode {HID: $toolHID})
SET tool.%s = coalesce(tool.%s, 0) + 1
RETURN tool.%s AS c
`, kind.SequenceProperty, kind.SequenceProperty, kind.SequenceProperty)
	return scalarInt(ctx, tx, query, map[string]any{"toolHID": schema.SSTPAToolHID})
}

func registeredEmailCount(ctx context.Context, tx neo4j.ManagedTransaction, kind Kind, email string) (int64, error) {
	return scalarInt(ctx, tx,
		"MATCH (registry:"+kind.RegistryLabel+":SSTPANode {HID: $registryHID})-[:"+kind.Relationship+"]->(n:"+kind.NodeLabel+":SSTPANode {UserEmail: $email}) RETURN count(n) AS c",
		map[string]any{"registryHID": kind.RegistryHID, "email": email})
}

func registeredCount(ctx context.Context, tx neo4j.ManagedTransaction, kind Kind) (int64, error) {
	return scalarInt(ctx, tx,
		"MATCH (registry:"+kind.RegistryLabel+":SSTPANode {HID: $registryHID})-[:"+kind.Relationship+"]->(n:"+kind.NodeLabel+":SSTPANode) RETURN count(n) AS c",
		map[string]any{"registryHID": kind.RegistryHID})
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
