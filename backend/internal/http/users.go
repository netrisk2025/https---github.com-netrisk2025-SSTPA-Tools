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

type bootstrapInstallerRequest struct {
	InstallerName  string `json:"installerName"`
	InstallerEmail string `json:"installerEmail"`
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

		actor, ok := api.actorForOnboardingCreate(writer, request, kind, payload)
		if !ok {
			return
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

func (api api) bootstrapInstallerHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	var payload bootstrapInstallerRequest
	if err := decodeJSON(request, &payload); err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	if payload.InstallerName == "" || payload.InstallerEmail == "" {
		writeError(writer, http.StatusBadRequest, "installerName and installerEmail are required")
		return
	}

	result, err := onboarding.BootstrapInstaller(request.Context(), api.driver, api.databaseName, onboarding.BootstrapInput{
		InstallerName:  payload.InstallerName,
		InstallerEmail: payload.InstallerEmail,
		Now:            api.now(),
	})
	if err != nil {
		if errors.Is(err, onboarding.ErrAlreadyRegistered) {
			writeError(writer, http.StatusConflict, err.Error())
			return
		}
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(writer, http.StatusCreated, result)
}

func (api api) actorForOnboardingCreate(writer http.ResponseWriter, request *http.Request, kind onboarding.Kind, payload createOnboardingRequest) (metadata.Actor, bool) {
	if kind == onboarding.UserKind {
		actor, err := actorFromRequest(request, metadata.Actor{})
		if err != nil {
			return metadata.Actor{Name: payload.UserName, Email: payload.UserEmail}, true
		}
		return actor, true
	}

	actor, err := actorFromRequest(request, metadata.Actor{})
	if err != nil || !actor.Admin {
		writeError(writer, http.StatusForbidden, "registered admin actor is required")
		return metadata.Actor{}, false
	}

	registered, err := onboarding.IsRegisteredAdmin(request.Context(), api.driver, api.databaseName, actor)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return metadata.Actor{}, false
	}
	if !registered {
		writeError(writer, http.StatusForbidden, "registered admin actor is required")
		return metadata.Actor{}, false
	}

	return actor, true
}
