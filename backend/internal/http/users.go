// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package apihttp

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"sstpa-tool/backend/internal/metadata"
	"sstpa-tool/backend/internal/onboarding"
)

type createOnboardingRequest struct {
	UserName  string `json:"userName"`
	UserEmail string `json:"userEmail"`
}

func (api api) listOnboardingHandler(kind onboarding.Kind) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if !api.requireDriver(writer) {
			return
		}

		page, err := parsePagination(request)
		if err != nil {
			writeError(writer, http.StatusBadRequest, err.Error())
			return
		}

		result, err := onboarding.List(request.Context(), api.driver, api.databaseName, kind, onboarding.Page{
			Page:  page.Page,
			Limit: page.Limit,
		})
		if err != nil {
			writeError(writer, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(writer, http.StatusOK, result)
	}
}

func (api api) getOnboardingHandler(kind onboarding.Kind) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if !api.requireDriver(writer) {
			return
		}

		uuid := chi.URLParam(request, "uuid")
		record, err := onboarding.GetByUUID(request.Context(), api.driver, api.databaseName, kind, uuid)
		if err != nil {
			if errors.Is(err, onboarding.ErrNotFound) {
				writeError(writer, http.StatusNotFound, "not found")
				return
			}
			writeError(writer, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(writer, http.StatusOK, record)
	}
}

func (api api) createOnboardingHandler(kind onboarding.Kind) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if !api.requireDriver(writer) {
			return
		}

		var payload createOnboardingRequest
		if err := decodeJSON(request, &payload); err != nil {
			writeError(writer, http.StatusBadRequest, err.Error())
			return
		}
		if payload.UserName == "" || payload.UserEmail == "" {
			writeError(writer, http.StatusBadRequest, "userName and userEmail are required")
			return
		}

		// TODO(sstpa-auth): the self-declaration fallback (Actor == payload values) is an
		// installer affordance per SRS §1.4.2: the first Admin + User have no pre-existing
		// identity to attribute. Once the auth layer lands, require X-SSTPA-User headers
		// for POST /admins and tighten the fallback for POST /users to first-run only.
		actor, err := actorFromRequest(request, metadata.Actor{})
		if err != nil {
			actor = metadata.Actor{Name: payload.UserName, Email: payload.UserEmail}
		}

		record, err := onboarding.Create(request.Context(), api.driver, api.databaseName, kind, onboarding.CreateInput{
			UserName:  payload.UserName,
			UserEmail: payload.UserEmail,
			Actor:     actor,
			Now:       api.now(),
		})
		if err != nil {
			if errors.Is(err, onboarding.ErrAlreadyRegistered) {
				writeError(writer, http.StatusConflict, err.Error())
				return
			}
			writeError(writer, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(writer, http.StatusCreated, record)
	}
}
