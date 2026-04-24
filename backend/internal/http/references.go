// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package apihttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sstpa-tool/backend/internal/graph"
	"sstpa-tool/backend/internal/identity"
	"sstpa-tool/backend/internal/messaging"
	"sstpa-tool/backend/internal/metadata"
	"sstpa-tool/backend/internal/mutation"
)

type referenceFrameworkResponse struct {
	FrameworkName    string         `json:"frameworkName"`
	FrameworkVersion string         `json:"frameworkVersion"`
	Properties       map[string]any `json:"properties"`
}

type referenceItemResponse struct {
	UUID             string         `json:"uuid"`
	FrameworkName    string         `json:"frameworkName"`
	FrameworkVersion string         `json:"frameworkVersion"`
	ExternalID       string         `json:"externalId"`
	ExternalType     string         `json:"externalType"`
	Name             string         `json:"name"`
	ShortDescription string         `json:"shortDescription"`
	LongDescription  string         `json:"longDescription,omitempty"`
	SourceURI        string         `json:"sourceUri"`
	Properties       map[string]any `json:"properties"`
}

type referenceAssignmentResponse struct {
	SourceHID     string                `json:"sourceHid"`
	ReferenceItem referenceItemResponse `json:"referenceItem"`
}

type referenceAssignmentMutationResponse struct {
	Assignment   referenceAssignmentResponse `json:"assignment"`
	CommitReport mutation.CommitReport       `json:"commitReport"`
}

type validateReferenceAssignmentRequest struct {
	SourceHID        string `json:"sourceHid"`
	ReferenceUUID    string `json:"referenceUuid"`
	ExternalID       string `json:"externalId"`
	FrameworkName    string `json:"frameworkName"`
	FrameworkVersion string `json:"frameworkVersion"`
}

type referenceAssignmentRequest struct {
	Actor            actorRequest `json:"actor"`
	SourceHID        string       `json:"sourceHid"`
	ReferenceUUID    string       `json:"referenceUuid"`
	ExternalID       string       `json:"externalId"`
	FrameworkName    string       `json:"frameworkName"`
	FrameworkVersion string       `json:"frameworkVersion"`
	CommitID         string       `json:"commitId"`
}

func (api api) listReferenceFrameworksHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	page, err := parsePagination(request)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	session := api.driver.NewSession(request.Context(), neo4j.SessionConfig{DatabaseName: api.databaseName})
	defer session.Close(request.Context())

	value, err := session.ExecuteRead(request.Context(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(request.Context(), `
MATCH (framework:ReferenceFramework)
RETURN properties(framework) AS properties
ORDER BY framework.FrameworkName, framework.FrameworkVersion
SKIP $skip
LIMIT $limit
`, map[string]any{"skip": page.Offset, "limit": page.Limit})
		if err != nil {
			return nil, err
		}
		records, err := result.Collect(request.Context())
		if err != nil {
			return nil, err
		}
		frameworks := make([]referenceFrameworkResponse, 0, len(records))
		for _, record := range records {
			frameworks = append(frameworks, frameworkFromProperties(propertiesFromRecord(record, "properties")))
		}

		countResult, err := tx.Run(request.Context(), `
MATCH (framework:ReferenceFramework)
RETURN count(framework) AS total
`, nil)
		if err != nil {
			return nil, err
		}
		countRecord, err := countResult.Single(request.Context())
		if err != nil {
			return nil, err
		}
		totalValue, _ := countRecord.Get("total")
		total, _ := totalValue.(int64)

		return listResponse[referenceFrameworkResponse]{
			Items: frameworks,
			Page:  page.Page,
			Limit: page.Limit,
			Total: total,
		}, nil
	})
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(writer, http.StatusOK, value)
}

func (api api) listReferenceItemsHandler(writer http.ResponseWriter, request *http.Request) {
	api.searchReferenceItemsHandler(writer, request)
}

func (api api) searchReferenceItemsHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	page, err := parsePagination(request)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	result, err := api.searchReferenceItems(request.Context(), referenceSearchInput{
		Query:            request.URL.Query().Get("q"),
		FrameworkName:    request.URL.Query().Get("frameworkName"),
		FrameworkVersion: request.URL.Query().Get("frameworkVersion"),
		ExternalType:     request.URL.Query().Get("externalType"),
		Page:             page,
	})
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(writer, http.StatusOK, result)
}

func (api api) getReferenceItemByExternalIDHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	item, err := api.findReferenceItem(request.Context(), referenceLookup{
		ExternalID:       chi.URLParam(request, "externalID"),
		FrameworkName:    request.URL.Query().Get("frameworkName"),
		FrameworkVersion: request.URL.Query().Get("frameworkVersion"),
	})
	if err != nil {
		handleNeo4jReadError(writer, err, "reference item not found")
		return
	}

	writeJSON(writer, http.StatusOK, item)
}

func (api api) getReferenceItemByUUIDHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	item, err := api.findReferenceItem(request.Context(), referenceLookup{UUID: chi.URLParam(request, "uuid")})
	if err != nil {
		handleNeo4jReadError(writer, err, "reference item not found")
		return
	}

	writeJSON(writer, http.StatusOK, item)
}

func (api api) relatedReferenceItemsHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	uuid := chi.URLParam(request, "uuid")
	page, err := parsePagination(request)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	session := api.driver.NewSession(request.Context(), neo4j.SessionConfig{DatabaseName: api.databaseName})
	defer session.Close(request.Context())
	value, err := session.ExecuteRead(request.Context(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(request.Context(), `
MATCH (item:ReferenceItem {uuid: $uuid})-[relationship:HAS_CHILD|RELATED_TO]-(related:ReferenceItem)
RETURN properties(related) AS properties
ORDER BY related.ExternalID
SKIP $skip
LIMIT $limit
`, map[string]any{"uuid": uuid, "skip": page.Offset, "limit": page.Limit})
		if err != nil {
			return nil, err
		}
		records, err := result.Collect(request.Context())
		if err != nil {
			return nil, err
		}
		items := make([]referenceItemResponse, 0, len(records))
		for _, record := range records {
			items = append(items, referenceItemFromProperties(propertiesFromRecord(record, "properties")))
		}

		countResult, err := tx.Run(request.Context(), `
MATCH (item:ReferenceItem {uuid: $uuid})-[relationship:HAS_CHILD|RELATED_TO]-(related:ReferenceItem)
RETURN count(related) AS total
`, map[string]any{"uuid": uuid})
		if err != nil {
			return nil, err
		}
		countRecord, err := countResult.Single(request.Context())
		if err != nil {
			return nil, err
		}
		totalValue, _ := countRecord.Get("total")
		total, _ := totalValue.(int64)

		return listResponse[referenceItemResponse]{Items: items, Page: page.Page, Limit: page.Limit, Total: total}, nil
	})
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(writer, http.StatusOK, value)
}

func (api api) validateReferenceAssignmentHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	var payload validateReferenceAssignmentRequest
	if err := decodeJSON(request, &payload); err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	result, err := api.validateReferenceAssignment(request.Context(), referenceAssignmentLookup{
		SourceHID: payload.SourceHID,
		Lookup: referenceLookup{
			UUID:             payload.ReferenceUUID,
			ExternalID:       payload.ExternalID,
			FrameworkName:    payload.FrameworkName,
			FrameworkVersion: payload.FrameworkVersion,
		},
	})
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(writer, http.StatusOK, result)
}

func (api api) listReferenceAssignmentsHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	sourceHID := chi.URLParam(request, "sourceHID")
	page, err := parsePagination(request)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	session := api.driver.NewSession(request.Context(), neo4j.SessionConfig{DatabaseName: api.databaseName})
	defer session.Close(request.Context())
	value, err := session.ExecuteRead(request.Context(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(request.Context(), `
MATCH (source:SSTPANode {HID: $sourceHID})-[:REFERENCES]->(item:ReferenceItem)
RETURN properties(item) AS properties
ORDER BY item.FrameworkName, item.ExternalID
SKIP $skip
LIMIT $limit
`, map[string]any{"sourceHID": sourceHID, "skip": page.Offset, "limit": page.Limit})
		if err != nil {
			return nil, err
		}
		records, err := result.Collect(request.Context())
		if err != nil {
			return nil, err
		}
		items := make([]referenceAssignmentResponse, 0, len(records))
		for _, record := range records {
			items = append(items, referenceAssignmentResponse{
				SourceHID:     sourceHID,
				ReferenceItem: referenceItemFromProperties(propertiesFromRecord(record, "properties")),
			})
		}

		countResult, err := tx.Run(request.Context(), `
MATCH (source:SSTPANode {HID: $sourceHID})-[:REFERENCES]->(item:ReferenceItem)
RETURN count(item) AS total
`, map[string]any{"sourceHID": sourceHID})
		if err != nil {
			return nil, err
		}
		countRecord, err := countResult.Single(request.Context())
		if err != nil {
			return nil, err
		}
		totalValue, _ := countRecord.Get("total")
		total, _ := totalValue.(int64)

		return listResponse[referenceAssignmentResponse]{Items: items, Page: page.Page, Limit: page.Limit, Total: total}, nil
	})
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(writer, http.StatusOK, value)
}

func (api api) createReferenceAssignmentHandler(writer http.ResponseWriter, request *http.Request) {
	api.referenceAssignmentMutationHandler(writer, request, false)
}

func (api api) deleteReferenceAssignmentHandler(writer http.ResponseWriter, request *http.Request) {
	api.referenceAssignmentMutationHandler(writer, request, true)
}

func (api api) referenceAssignmentMutationHandler(writer http.ResponseWriter, request *http.Request, deleteRelationship bool) {
	if !api.requireDriver(writer) {
		return
	}

	var payload referenceAssignmentRequest
	if err := decodeJSON(request, &payload); err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	actor, err := actorFromRequest(request, metadata.Actor{Name: payload.Actor.Name, Email: payload.Actor.Email, Admin: payload.Actor.Admin})
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	if payload.SourceHID == "" {
		writeError(writer, http.StatusBadRequest, "sourceHid is required")
		return
	}

	response, err := api.mutateReferenceAssignment(request.Context(), referenceAssignmentMutationInput{
		SourceHID:           payload.SourceHID,
		Lookup:              referenceLookup{UUID: payload.ReferenceUUID, ExternalID: payload.ExternalID, FrameworkName: payload.FrameworkName, FrameworkVersion: payload.FrameworkVersion},
		Actor:               actor,
		CommitID:            payload.CommitID,
		DeleteRelationship:  deleteRelationship,
		Now:                 api.now(),
		ChangeTypeSummary:   referenceAssignmentChangeSummary(deleteRelationship),
		RelationshipChanged: "REFERENCES",
	})
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	if deleteRelationship {
		writeJSON(writer, http.StatusOK, response)
		return
	}
	writeJSON(writer, http.StatusCreated, response)
}

type referenceSearchInput struct {
	Query            string
	FrameworkName    string
	FrameworkVersion string
	ExternalType     string
	Page             pageRequest
}

type referenceLookup struct {
	UUID             string
	ExternalID       string
	FrameworkName    string
	FrameworkVersion string
}

type referenceAssignmentLookup struct {
	SourceHID string
	Lookup    referenceLookup
}

type referenceAssignmentMutationInput struct {
	SourceHID           string
	Lookup              referenceLookup
	Actor               metadata.Actor
	CommitID            string
	DeleteRelationship  bool
	Now                 time.Time
	ChangeTypeSummary   string
	RelationshipChanged string
}

func (api api) searchReferenceItems(ctx context.Context, input referenceSearchInput) (listResponse[referenceItemResponse], error) {
	session := api.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: api.databaseName})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		params := map[string]any{
			"query":            input.Query,
			"lower":            strings.ToLower(input.Query),
			"frameworkName":    input.FrameworkName,
			"frameworkVersion": input.FrameworkVersion,
			"externalType":     input.ExternalType,
			"skip":             input.Page.Offset,
			"limit":            input.Page.Limit,
		}
		where := `
($frameworkName = "" OR item.FrameworkName = $frameworkName)
AND ($frameworkVersion = "" OR item.FrameworkVersion = $frameworkVersion)
AND ($externalType = "" OR item.ExternalType = $externalType)
AND (
  $query = ""
  OR item.ExternalID = $query
  OR toLower(coalesce(item.Name, "")) CONTAINS $lower
  OR toLower(coalesce(item.ShortDescription, "")) CONTAINS $lower
)
`
		records, err := tx.Run(ctx, `
MATCH (item:ReferenceItem)
WHERE `+where+`
RETURN properties(item) AS properties
ORDER BY item.FrameworkName, item.ExternalID
SKIP $skip
LIMIT $limit
`, params)
		if err != nil {
			return nil, err
		}
		collected, err := records.Collect(ctx)
		if err != nil {
			return nil, err
		}
		items := make([]referenceItemResponse, 0, len(collected))
		for _, record := range collected {
			items = append(items, referenceItemFromProperties(propertiesFromRecord(record, "properties")))
		}

		countResult, err := tx.Run(ctx, `
MATCH (item:ReferenceItem)
WHERE `+where+`
RETURN count(item) AS total
`, params)
		if err != nil {
			return nil, err
		}
		countRecord, err := countResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		totalValue, _ := countRecord.Get("total")
		total, _ := totalValue.(int64)
		return listResponse[referenceItemResponse]{Items: items, Page: input.Page.Page, Limit: input.Page.Limit, Total: total}, nil
	})
	if err != nil {
		return listResponse[referenceItemResponse]{}, err
	}
	response, _ := result.(listResponse[referenceItemResponse])
	return response, nil
}

func (api api) findReferenceItem(ctx context.Context, lookup referenceLookup) (referenceItemResponse, error) {
	session := api.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: api.databaseName})
	defer session.Close(ctx)

	value, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return readReferenceItem(ctx, tx, lookup)
	})
	if err != nil {
		return referenceItemResponse{}, err
	}
	item, _ := value.(referenceItemResponse)
	return item, nil
}

func (api api) validateReferenceAssignment(ctx context.Context, lookup referenceAssignmentLookup) (validationResponse, error) {
	session := api.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: api.databaseName})
	defer session.Close(ctx)

	value, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		source, item, duplicate, err := readReferenceAssignmentState(ctx, tx, lookup)
		if err != nil {
			return validationResponse{Valid: false, Reason: err.Error()}, nil
		}
		return validateReferenceAssignmentState(source, item, duplicate), nil
	})
	if err != nil {
		return validationResponse{}, err
	}
	response, _ := value.(validationResponse)
	return response, nil
}

func (api api) mutateReferenceAssignment(ctx context.Context, input referenceAssignmentMutationInput) (referenceAssignmentMutationResponse, error) {
	commitID := input.CommitID
	if commitID == "" {
		commitID = identity.NewUUID()
	}

	session := api.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: api.databaseName})
	defer session.Close(ctx)

	value, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		source, item, duplicate, err := readReferenceAssignmentState(ctx, tx, referenceAssignmentLookup{SourceHID: input.SourceHID, Lookup: input.Lookup})
		if err != nil {
			return nil, err
		}
		validation := validateReferenceAssignmentState(source, item, duplicate)
		if !validation.Valid && !(input.DeleteRelationship && duplicate) {
			return nil, errors.New(validation.Reason)
		}
		if input.DeleteRelationship && !duplicate {
			return nil, errors.New("reference assignment does not exist")
		}
		if !input.DeleteRelationship && duplicate {
			return nil, errors.New("duplicate reference assignment already exists")
		}

		query := `
MATCH (source:SSTPANode {HID: $sourceHID}), (item:ReferenceItem {uuid: $referenceUUID})
CREATE (source)-[:REFERENCES {
  Created: $now,
  Creator: $actorName,
  CreatorEmail: $actorEmail,
  CommitID: $commitID
}]->(item)
RETURN properties(item) AS properties
`
		if input.DeleteRelationship {
			query = `
MATCH (source:SSTPANode {HID: $sourceHID})-[relationship:REFERENCES]->(item:ReferenceItem {uuid: $referenceUUID})
DELETE relationship
RETURN properties(item) AS properties
`
		}

		result, err := tx.Run(ctx, query, map[string]any{
			"sourceHID":     input.SourceHID,
			"referenceUUID": item.UUID,
			"now":           input.Now.UTC().Format(time.RFC3339),
			"actorName":     input.Actor.Name,
			"actorEmail":    input.Actor.Email,
			"commitID":      commitID,
		})
		if err != nil {
			return nil, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		changedItem := referenceItemFromProperties(propertiesFromRecord(record, "properties"))

		recipients := []string{}
		if source.Owner == "" || source.OwnerEmail == "" {
			return nil, fmt.Errorf("source node %s is missing Owner/OwnerEmail", source.HID)
		}
		if source.Owner != input.Actor.Name || source.OwnerEmail != input.Actor.Email {
			_, err := messaging.AppendChangeNotification(ctx, tx, messaging.ChangeNotification{
				CommitID:                 commitID,
				Subject:                  "SSTPA reference assignment changed",
				Body:                     fmt.Sprintf("%s changed reference assignments for %s.", input.Actor.Name, source.HID),
				SentAt:                   input.Now,
				Sender:                   input.Actor,
				Recipient:                metadata.Actor{Name: source.Owner, Email: source.OwnerEmail},
				RelatedNodeHIDs:          []string{source.HID},
				RelatedRelationshipTypes: []string{input.RelationshipChanged},
				ChangeTypeSummary:        input.ChangeTypeSummary,
			})
			if err != nil {
				return nil, err
			}
			recipients = append(recipients, source.OwnerEmail)
		}

		return referenceAssignmentMutationResponse{
			Assignment: referenceAssignmentResponse{
				SourceHID:     source.HID,
				ReferenceItem: changedItem,
			},
			CommitReport: mutation.CommitReport{
				CommitID:             commitID,
				NodesChanged:         []string{source.HID},
				RelationshipsChanged: []string{input.RelationshipChanged},
				MessagesGenerated:    len(recipients),
				RecipientsNotified:   recipients,
			},
		}, nil
	})
	if err != nil {
		return referenceAssignmentMutationResponse{}, err
	}
	response, _ := value.(referenceAssignmentMutationResponse)
	return response, nil
}

type referenceAssignmentSource struct {
	HID        string
	TypeName   identity.NodeType
	Owner      string
	OwnerEmail string
}

func readReferenceAssignmentState(ctx context.Context, tx neo4j.ManagedTransaction, lookup referenceAssignmentLookup) (referenceAssignmentSource, referenceItemResponse, bool, error) {
	if lookup.SourceHID == "" {
		return referenceAssignmentSource{}, referenceItemResponse{}, false, errors.New("sourceHid is required")
	}

	sourceResult, err := tx.Run(ctx, `
MATCH (source:SSTPANode {HID: $sourceHID})
RETURN source.HID AS hid, source.TypeName AS typeName, source.Owner AS owner, source.OwnerEmail AS ownerEmail
LIMIT 1
`, map[string]any{"sourceHID": lookup.SourceHID})
	if err != nil {
		return referenceAssignmentSource{}, referenceItemResponse{}, false, err
	}
	sourceRecord, err := sourceResult.Single(ctx)
	if err != nil {
		return referenceAssignmentSource{}, referenceItemResponse{}, false, errors.New("source SSTPA node does not exist")
	}
	source := referenceAssignmentSource{
		HID:        stringFromRecord(sourceRecord, "hid"),
		TypeName:   identity.NodeType(stringFromRecord(sourceRecord, "typeName")),
		Owner:      stringFromRecord(sourceRecord, "owner"),
		OwnerEmail: stringFromRecord(sourceRecord, "ownerEmail"),
	}

	item, err := readReferenceItem(ctx, tx, lookup.Lookup)
	if err != nil {
		return referenceAssignmentSource{}, referenceItemResponse{}, false, errors.New("reference item does not exist")
	}

	duplicateResult, err := tx.Run(ctx, `
MATCH (source:SSTPANode {HID: $sourceHID}), (item:ReferenceItem {uuid: $referenceUUID})
OPTIONAL MATCH (source)-[relationship:REFERENCES]->(item)
RETURN count(relationship) AS count
`, map[string]any{"sourceHID": source.HID, "referenceUUID": item.UUID})
	if err != nil {
		return referenceAssignmentSource{}, referenceItemResponse{}, false, err
	}
	duplicateRecord, err := duplicateResult.Single(ctx)
	if err != nil {
		return referenceAssignmentSource{}, referenceItemResponse{}, false, err
	}
	countValue, _ := duplicateRecord.Get("count")
	count, _ := countValue.(int64)

	return source, item, count > 0, nil
}

func validateReferenceAssignmentState(source referenceAssignmentSource, item referenceItemResponse, duplicate bool) validationResponse {
	if _, ok := identity.TypeID(source.TypeName); !ok {
		return validationResponse{Valid: false, Reason: "source SSTPA node has an unknown TypeName"}
	}
	if !graph.ReferenceAssignmentAllowed(source.TypeName, item.FrameworkName, item.ExternalType) {
		return validationResponse{Valid: false, Reason: fmt.Sprintf("%s cannot reference %s %s", source.TypeName, item.FrameworkName, item.ExternalType)}
	}
	if duplicate {
		return validationResponse{Valid: false, Reason: "duplicate reference assignment already exists"}
	}

	return validationResponse{Valid: true, Reason: "reference assignment is allowed"}
}

func readReferenceItem(ctx context.Context, tx neo4j.ManagedTransaction, lookup referenceLookup) (referenceItemResponse, error) {
	if lookup.UUID == "" && lookup.ExternalID == "" {
		return referenceItemResponse{}, errors.New("referenceUuid or externalId is required")
	}

	where := "item.ExternalID = $externalID"
	if lookup.UUID != "" {
		where = "item.uuid = $uuid"
	}
	result, err := tx.Run(ctx, `
MATCH (item:ReferenceItem)
WHERE `+where+`
  AND ($frameworkName = "" OR item.FrameworkName = $frameworkName)
  AND ($frameworkVersion = "" OR item.FrameworkVersion = $frameworkVersion)
RETURN properties(item) AS properties
ORDER BY item.FrameworkName, item.FrameworkVersion
LIMIT 1
`, map[string]any{
		"uuid":             lookup.UUID,
		"externalID":       lookup.ExternalID,
		"frameworkName":    lookup.FrameworkName,
		"frameworkVersion": lookup.FrameworkVersion,
	})
	if err != nil {
		return referenceItemResponse{}, err
	}
	record, err := result.Single(ctx)
	if err != nil {
		return referenceItemResponse{}, err
	}
	return referenceItemFromProperties(propertiesFromRecord(record, "properties")), nil
}

func frameworkFromProperties(properties map[string]any) referenceFrameworkResponse {
	return referenceFrameworkResponse{
		FrameworkName:    stringProperty(properties, "FrameworkName"),
		FrameworkVersion: stringProperty(properties, "FrameworkVersion"),
		Properties:       properties,
	}
}

func referenceItemFromProperties(properties map[string]any) referenceItemResponse {
	return referenceItemResponse{
		UUID:             stringProperty(properties, "uuid"),
		FrameworkName:    stringProperty(properties, "FrameworkName"),
		FrameworkVersion: stringProperty(properties, "FrameworkVersion"),
		ExternalID:       stringProperty(properties, "ExternalID"),
		ExternalType:     stringProperty(properties, "ExternalType"),
		Name:             stringProperty(properties, "Name"),
		ShortDescription: stringProperty(properties, "ShortDescription"),
		LongDescription:  stringProperty(properties, "LongDescription"),
		SourceURI:        stringProperty(properties, "SourceURI"),
		Properties:       properties,
	}
}

func stringFromRecord(record *neo4j.Record, key string) string {
	value, _ := record.Get(key)
	return stringValue(value)
}

func referenceAssignmentChangeSummary(deleteRelationship bool) string {
	if deleteRelationship {
		return "reference assignment removed"
	}

	return "reference assignment created"
}
