// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package apihttp

import (
	"net/http"
	"time"

	"sstpa-tool/backend/internal/identity"
	"sstpa-tool/backend/internal/metadata"
	"sstpa-tool/backend/internal/mutation"
)

type mutationRequest struct {
	Actor      actorRequest               `json:"actor"`
	CommitID   string                     `json:"commitId"`
	VersionID  string                     `json:"versionId"`
	Now        string                     `json:"now"`
	Operations []mutationOperationRequest `json:"operations"`
}

type actorRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Admin bool   `json:"admin"`
}

type mutationOperationRequest struct {
	Kind                   mutation.OperationKind `json:"kind"`
	NodeType               identity.NodeType      `json:"nodeType"`
	HID                    string                 `json:"hid"`
	UUID                   string                 `json:"uuid"`
	Properties             map[string]any         `json:"properties"`
	RelationshipName       string                 `json:"relationshipName"`
	FromHID                string                 `json:"fromHid"`
	FromType               identity.NodeType      `json:"fromType"`
	ToHID                  string                 `json:"toHid"`
	ToType                 identity.NodeType      `json:"toType"`
	RelationshipProperties map[string]any         `json:"relationshipProperties"`
}

func (api api) mutationsHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	var payload mutationRequest
	if err := decodeJSON(request, &payload); err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	actor, err := actorFromRequest(request, metadata.Actor{
		Name:  payload.Actor.Name,
		Email: payload.Actor.Email,
		Admin: payload.Actor.Admin,
	})
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	now := api.now()
	if payload.Now != "" {
		parsed, err := time.Parse(time.RFC3339, payload.Now)
		if err != nil {
			writeError(writer, http.StatusBadRequest, "now must be RFC3339")
			return
		}
		now = parsed
	}

	plan := mutation.Plan{Operations: make([]mutation.Operation, 0, len(payload.Operations))}
	for _, operation := range payload.Operations {
		plan.Operations = append(plan.Operations, mutation.Operation{
			Kind:                   operation.Kind,
			NodeType:               operation.NodeType,
			HID:                    operation.HID,
			UUID:                   operation.UUID,
			Properties:             operation.Properties,
			RelationshipName:       operation.RelationshipName,
			FromHID:                operation.FromHID,
			FromType:               operation.FromType,
			ToHID:                  operation.ToHID,
			ToType:                 operation.ToType,
			RelationshipProperties: operation.RelationshipProperties,
		})
	}

	report, err := mutation.Apply(request.Context(), api.driver, mutation.ApplyOptions{
		DatabaseName: api.databaseName,
		Actor:        actor,
		Now:          now,
		CommitID:     payload.CommitID,
		VersionID:    payload.VersionID,
	}, plan)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(writer, http.StatusCreated, report)
}
