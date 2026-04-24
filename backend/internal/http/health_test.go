package apihttp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	recorder := httptest.NewRecorder()

	NewRouter("test").ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response healthResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("expected valid JSON response: %v", err)
	}

	if response.Status != "ok" {
		t.Fatalf("expected ok status, got %q", response.Status)
	}
}

func TestHealthEndpointRejectsUnsupportedMethods(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	recorder := httptest.NewRecorder()

	NewRouter("test").ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, recorder.Code)
	}
}
