// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package apihttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sstpa-tool/backend/internal/metadata"
)

const (
	defaultPageLimit = 50
	maxPageLimit     = 200
)

type api struct {
	version      string
	driver       neo4j.DriverWithContext
	databaseName string
	now          func() time.Time
}

type apiError struct {
	Error string `json:"error"`
}

type pageRequest struct {
	Page   int
	Limit  int
	Offset int
}

type listResponse[T any] struct {
	Items []T   `json:"items"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

func newAPI(options RouterOptions) api {
	version := options.Version
	if version == "" {
		version = "dev"
	}

	return api{
		version:      version,
		driver:       options.Driver,
		databaseName: options.DatabaseName,
		now:          func() time.Time { return time.Now().UTC() },
	}
}

func (api api) requireDriver(writer http.ResponseWriter) bool {
	if api.driver != nil {
		return true
	}

	writeError(writer, http.StatusServiceUnavailable, "graph persistence is not configured")
	return false
}

func parsePagination(request *http.Request) (pageRequest, error) {
	query := request.URL.Query()
	page := 1
	limit := defaultPageLimit

	if raw := query.Get("page"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 {
			return pageRequest{}, fmt.Errorf("page must be a positive integer")
		}
		page = parsed
	}

	if raw := query.Get("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 {
			return pageRequest{}, fmt.Errorf("limit must be a positive integer")
		}
		if parsed > maxPageLimit {
			return pageRequest{}, fmt.Errorf("limit must be less than or equal to %d", maxPageLimit)
		}
		limit = parsed
	}

	return pageRequest{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}, nil
}

func actorFromRequest(request *http.Request, fallback metadata.Actor) (metadata.Actor, error) {
	actor := fallback
	if actor.Name == "" {
		actor.Name = request.Header.Get("X-SSTPA-User")
	}
	if actor.Email == "" {
		actor.Email = request.Header.Get("X-SSTPA-User-Email")
	}
	if request.Header.Get("X-SSTPA-Admin") == "true" {
		actor.Admin = true
	}

	if actor.Name == "" || actor.Email == "" {
		return metadata.Actor{}, errors.New("actor name and email are required")
	}

	return actor, nil
}

func emailFromRequest(request *http.Request) (string, error) {
	email := request.Header.Get("X-SSTPA-User-Email")
	if email == "" {
		email = request.URL.Query().Get("userEmail")
	}
	if email == "" {
		return "", errors.New("user email is required")
	}

	return email, nil
}

func decodeJSON(request *http.Request, target any) error {
	defer request.Body.Close()
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}

func writeError(writer http.ResponseWriter, status int, message string) {
	writeJSON(writer, status, apiError{Error: message})
}

func sortDirection(value string) (string, error) {
	switch strings.ToLower(value) {
	case "", "asc":
		return "ASC", nil
	case "desc":
		return "DESC", nil
	default:
		return "", fmt.Errorf("sort direction must be asc or desc")
	}
}

func stringProperty(properties map[string]any, key string) string {
	value, _ := properties[key].(string)
	return value
}

func anySliceToStrings(value any) []string {
	switch typed := value.(type) {
	case []string:
		return append([]string(nil), typed...)
	case []any:
		items := make([]string, 0, len(typed))
		for _, item := range typed {
			if text, ok := item.(string); ok {
				items = append(items, text)
			}
		}
		return items
	default:
		return nil
	}
}
