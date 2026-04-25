// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package apihttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sstpa-tool/backend/internal/onboarding"
	"sstpa-tool/backend/internal/schema"
	"sstpa-tool/backend/internal/testhelpers"
)

func TestUsersEndpointCreateListGet(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	if err := schema.Bootstrap(ctx, fixture.Driver, ""); err != nil {
		t.Fatal(err)
	}

	router := NewRouterWithOptions(RouterOptions{Version: "test", Driver: fixture.Driver})

	body := `{"userName": "Alice Analyst", "userEmail": "alice@example.test"}`
	recorder := performJSON(router, http.MethodPost, "/api/v1/users", body)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("create user status = %d body=%s", recorder.Code, recorder.Body.String())
	}
	var created onboarding.Record
	if err := json.Unmarshal(recorder.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(created.HID, "USR__") {
		t.Fatalf("unexpected HID %q", created.HID)
	}

	recorder = performJSON(router, http.MethodPost, "/api/v1/users", body)
	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected 409 on duplicate, got %d", recorder.Code)
	}

	recorder = performJSON(router, http.MethodGet, "/api/v1/users", "")
	if recorder.Code != http.StatusOK {
		t.Fatalf("list users status = %d", recorder.Code)
	}
	var list onboarding.ListResult
	if err := json.Unmarshal(recorder.Body.Bytes(), &list); err != nil {
		t.Fatal(err)
	}
	if list.Total != 1 || len(list.Items) != 1 {
		t.Fatalf("unexpected list: %#v", list)
	}

	recorder = performJSON(router, http.MethodGet, "/api/v1/users/"+created.UUID, "")
	if recorder.Code != http.StatusOK {
		t.Fatalf("get user status = %d", recorder.Code)
	}
	var fetched onboarding.Record
	if err := json.Unmarshal(recorder.Body.Bytes(), &fetched); err != nil {
		t.Fatal(err)
	}
	if fetched.HID != created.HID {
		t.Fatalf("fetched HID = %q, want %q", fetched.HID, created.HID)
	}
}

func TestAdminsEndpointCreateListGet(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	if err := schema.Bootstrap(ctx, fixture.Driver, ""); err != nil {
		t.Fatal(err)
	}

	router := NewRouterWithOptions(RouterOptions{Version: "test", Driver: fixture.Driver})

	bootstrapBody := `{"installerName": "Installer", "installerEmail": "installer@example.test"}`
	recorder := performJSON(router, http.MethodPost, "/api/v1/onboarding/bootstrap", bootstrapBody)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("bootstrap status = %d body=%s", recorder.Code, recorder.Body.String())
	}

	body := `{"userName": "Root Admin", "userEmail": "root@example.test"}`
	recorder = performJSONWithHeaders(router, http.MethodPost, "/api/v1/admins", body, map[string]string{
		"X-SSTPA-User":       "Installer",
		"X-SSTPA-User-Email": "installer@example.test",
		"X-SSTPA-Admin":      "true",
	})
	if recorder.Code != http.StatusCreated {
		t.Fatalf("create admin status = %d body=%s", recorder.Code, recorder.Body.String())
	}
	var created onboarding.Record
	if err := json.Unmarshal(recorder.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(created.HID, "ADM__") {
		t.Fatalf("unexpected HID %q", created.HID)
	}

	recorder = performJSON(router, http.MethodGet, "/api/v1/admins/"+created.UUID, "")
	if recorder.Code != http.StatusOK {
		t.Fatalf("get admin status = %d", recorder.Code)
	}
}

func TestBootstrapInstallerEndpointCreatesFirstAdminAndUser(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	if err := schema.Bootstrap(ctx, fixture.Driver, ""); err != nil {
		t.Fatal(err)
	}

	router := NewRouterWithOptions(RouterOptions{Version: "test", Driver: fixture.Driver})
	body := `{"installerName": "Installer", "installerEmail": "installer@example.test"}`
	recorder := performJSON(router, http.MethodPost, "/api/v1/onboarding/bootstrap", body)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("bootstrap status = %d body=%s", recorder.Code, recorder.Body.String())
	}
	var created onboarding.BootstrapResult
	if err := json.Unmarshal(recorder.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(created.Admin.HID, "ADM__") || !strings.HasPrefix(created.User.HID, "USR__") {
		t.Fatalf("unexpected bootstrap result: %#v", created)
	}

	recorder = performJSON(router, http.MethodPost, "/api/v1/onboarding/bootstrap", body)
	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected duplicate bootstrap 409, got %d", recorder.Code)
	}
}

func TestAdminsEndpointRejectsUnauthorizedCreate(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	if err := schema.Bootstrap(ctx, fixture.Driver, ""); err != nil {
		t.Fatal(err)
	}

	router := NewRouterWithOptions(RouterOptions{Version: "test", Driver: fixture.Driver})
	body := `{"userName": "Root Admin", "userEmail": "root@example.test"}`
	recorder := performJSON(router, http.MethodPost, "/api/v1/admins", body)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for missing admin actor, got %d", recorder.Code)
	}
}

func TestUsersEndpointRejectsMissingFields(t *testing.T) {
	router := NewRouterWithOptions(RouterOptions{Version: "test", Driver: fakeConfiguredDriver{}})

	recorder := performJSON(router, http.MethodPost, "/api/v1/users", `{"userName": "Alice"}`)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing email, got %d", recorder.Code)
	}
}

func performJSONWithHeaders(handler http.Handler, method string, target string, body string, headers map[string]string) *httptest.ResponseRecorder {
	var reader *bytes.Reader
	if body == "" {
		reader = bytes.NewReader(nil)
	} else {
		reader = bytes.NewReader([]byte(body))
	}
	request := httptest.NewRequest(method, target, reader)
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}
