// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package apihttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sstpa-tool/backend/internal/graph"
	"sstpa-tool/backend/internal/identity"
)

type validateRelationshipRequest struct {
	RelationshipName string            `json:"relationshipName"`
	FromType         identity.NodeType `json:"fromType"`
	ToType           identity.NodeType `json:"toType"`
	FromHID          string            `json:"fromHid"`
	ToHID            string            `json:"toHid"`
}

type validationResponse struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason"`
}

func (api api) validateRelationshipHandler(writer http.ResponseWriter, request *http.Request) {
	var payload validateRelationshipRequest
	if err := decodeJSON(request, &payload); err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	response := validateRelationshipCatalog(payload)
	if !response.Valid || payload.FromHID == "" || payload.ToHID == "" {
		writeJSON(writer, http.StatusOK, response)
		return
	}

	if !api.requireDriver(writer) {
		return
	}

	duplicate, err := api.relationshipExists(request.Context(), payload)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	if duplicate {
		writeJSON(writer, http.StatusOK, validationResponse{Valid: false, Reason: "duplicate relationship already exists"})
		return
	}

	writeJSON(writer, http.StatusOK, response)
}

func validateRelationshipCatalog(payload validateRelationshipRequest) validationResponse {
	if _, ok := identity.TypeID(payload.FromType); !ok {
		return validationResponse{Valid: false, Reason: fmt.Sprintf("unknown fromType %q", payload.FromType)}
	}
	if _, ok := identity.TypeID(payload.ToType); !ok {
		return validationResponse{Valid: false, Reason: fmt.Sprintf("unknown toType %q", payload.ToType)}
	}
	relationship, ok := graph.LookupRelationship(payload.RelationshipName, payload.FromType, payload.ToType)
	if !ok {
		return validationResponse{Valid: false, Reason: fmt.Sprintf("relationship %s from %s to %s is not allowed", payload.RelationshipName, payload.FromType, payload.ToType)}
	}
	if payload.FromHID != "" && payload.ToHID != "" {
		if err := graph.ValidateSoIBoundary(relationship, payload.FromHID, payload.ToHID, graph.DefaultRelationshipProperties(relationship)); err != nil {
			return validationResponse{Valid: false, Reason: err.Error()}
		}
	}

	return validationResponse{Valid: true, Reason: "allowed by relationship catalog"}
}

func (api api) relationshipExists(ctx context.Context, payload validateRelationshipRequest) (bool, error) {
	session := api.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: api.databaseName})
	defer session.Close(ctx)

	value, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := fmt.Sprintf(`
MATCH (from:SSTPANode {HID: $fromHID}), (to:SSTPANode {HID: $toHID})
OPTIONAL MATCH (from)-[r:%s]->(to)
RETURN count(r) AS count
`, payload.RelationshipName)
		result, err := tx.Run(ctx, query, map[string]any{"fromHID": payload.FromHID, "toHID": payload.ToHID})
		if err != nil {
			return false, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return false, err
		}
		countValue, _ := record.Get("count")
		count, _ := countValue.(int64)
		return count > 0, nil
	})
	if err != nil {
		return false, err
	}

	exists, _ := value.(bool)
	return exists, nil
}
