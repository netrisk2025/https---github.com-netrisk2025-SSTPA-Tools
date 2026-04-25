// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package apihttp

import (
	"net/http"
	"time"

	"sstpa-tool/backend/internal/metadata"
	"sstpa-tool/backend/internal/systemcreation"
)

type systemFromElementRequest struct {
	Actor      actorRequest `json:"actor"`
	ElementHID string       `json:"elementHid"`
	CommitID   string       `json:"commitId"`
	VersionID  string       `json:"versionId"`
	Now        string       `json:"now"`
}

func (api api) createSystemFromElementHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	var payload systemFromElementRequest
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

	result, err := systemcreation.CreateFromElement(request.Context(), api.driver, systemcreation.FromElementOptions{
		DatabaseName: api.databaseName,
		Actor:        actor,
		Now:          now,
		CommitID:     payload.CommitID,
		VersionID:    payload.VersionID,
	}, payload.ElementHID)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(writer, http.StatusCreated, result)
}
