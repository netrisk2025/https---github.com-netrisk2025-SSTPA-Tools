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

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sstpa-tool/backend/internal/schema"
	"sstpa-tool/backend/internal/testhelpers"
)

func TestRelationshipValidationWithoutGraph(t *testing.T) {
	body := `{"relationshipName":"HAS_SYSTEM","fromType":"Capability","toType":"System"}`
	request := httptest.NewRequest(http.MethodPost, "/api/v1/validate/relationship", strings.NewReader(body))
	recorder := httptest.NewRecorder()

	NewRouter("test").ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}

	var response validationResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if !response.Valid {
		t.Fatalf("expected valid relationship, got %#v", response)
	}
}

func TestListEndpointRejectsLimitAboveMaximum(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/nodes?limit=201", nil)
	recorder := httptest.NewRecorder()

	NewRouterWithOptions(RouterOptions{Version: "test", Driver: fakeConfiguredDriver{}}).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestOpenAPIEndpointPublishesVersion31Contract(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/openapi.yaml", nil)
	recorder := httptest.NewRecorder()

	NewRouter("test").ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), "openapi: 3.1.0") {
		t.Fatalf("expected OpenAPI 3.1 contract, got %s", recorder.Body.String())
	}
}

func TestMutationAndReadRoutes(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	if err := schema.Bootstrap(ctx, fixture.Driver, ""); err != nil {
		t.Fatal(err)
	}

	router := NewRouterWithOptions(RouterOptions{Version: "test", Driver: fixture.Driver})

	createCapability := `{
  "actor": {"name": "Alice", "email": "alice@example.test"},
  "commitId": "commit-api-create",
  "versionId": "test-schema",
  "operations": [{
    "kind": "create_node",
    "nodeType": "Capability",
    "hid": "CAP__0",
    "uuid": "00000000-0000-4000-8000-000000000201",
    "properties": {"Name": "Capability Root", "ShortDescription": "Customer intent"}
  }]
}`
	recorder := performJSON(router, http.MethodPost, "/api/v1/mutations", createCapability)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("create capability status = %d body=%s", recorder.Code, recorder.Body.String())
	}

	recorder = performJSON(router, http.MethodGet, "/api/v1/nodes/CAP__0", "")
	if recorder.Code != http.StatusOK {
		t.Fatalf("get node status = %d body=%s", recorder.Code, recorder.Body.String())
	}
	var node nodeResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &node); err != nil {
		t.Fatal(err)
	}
	if node.HID != "CAP__0" || node.Properties["Name"] != "Capability Root" {
		t.Fatalf("unexpected node response: %#v", node)
	}

	recorder = performJSON(router, http.MethodGet, "/api/v1/search?q=Customer&limit=10", "")
	if recorder.Code != http.StatusOK {
		t.Fatalf("search status = %d body=%s", recorder.Code, recorder.Body.String())
	}
	var search listResponse[nodeSummary]
	if err := json.Unmarshal(recorder.Body.Bytes(), &search); err != nil {
		t.Fatal(err)
	}
	if search.Total != 1 || len(search.Items) != 1 {
		t.Fatalf("unexpected search response: %#v", search)
	}

	updateCapability := `{
  "actor": {"name": "Bob", "email": "bob@example.test"},
  "commitId": "commit-api-update",
  "operations": [{
    "kind": "update_node",
    "hid": "CAP__0",
    "properties": {"Name": "Changed Root"}
  }]
}`
	recorder = performJSON(router, http.MethodPost, "/api/v1/mutations", updateCapability)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("update capability status = %d body=%s", recorder.Code, recorder.Body.String())
	}
	var report struct {
		MessagesGenerated  int      `json:"messagesGenerated"`
		RecipientsNotified []string `json:"recipientsNotified"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.MessagesGenerated != 1 || len(report.RecipientsNotified) != 1 || report.RecipientsNotified[0] != "alice@example.test" {
		t.Fatalf("unexpected commit report: %#v", report)
	}

	recorder = performJSON(router, http.MethodGet, "/api/v1/messages?userEmail=alice@example.test", "")
	if recorder.Code != http.StatusOK {
		t.Fatalf("messages status = %d body=%s", recorder.Code, recorder.Body.String())
	}
	var messages listResponse[messageResponse]
	if err := json.Unmarshal(recorder.Body.Bytes(), &messages); err != nil {
		t.Fatal(err)
	}
	if messages.Total != 1 || messages.Items[0].MessageType != "CHANGE_NOTIFICATION" {
		t.Fatalf("unexpected messages response: %#v", messages)
	}
}

func TestReferenceAssignmentRoutes(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	router := NewRouterWithOptions(RouterOptions{Version: "test", Driver: fixture.Driver})
	seedReferenceAssignmentGraph(t, fixture.Driver)

	recorder := performJSON(router, http.MethodGet, "/api/v1/reference/search?q=Technique", "")
	if recorder.Code != http.StatusOK {
		t.Fatalf("reference search status = %d body=%s", recorder.Code, recorder.Body.String())
	}
	var referenceSearch listResponse[referenceItemResponse]
	if err := json.Unmarshal(recorder.Body.Bytes(), &referenceSearch); err != nil {
		t.Fatal(err)
	}
	if referenceSearch.Total != 1 || referenceSearch.Items[0].ExternalID != "T1000" {
		t.Fatalf("unexpected reference search response: %#v", referenceSearch)
	}

	validateBody := `{"sourceHid":"HAZ_1_1","referenceUuid":"00000000-0000-4000-8000-000000000301"}`
	recorder = performJSON(router, http.MethodPost, "/api/v1/reference/validate-assignment", validateBody)
	if recorder.Code != http.StatusOK {
		t.Fatalf("validate assignment status = %d body=%s", recorder.Code, recorder.Body.String())
	}
	var validation validationResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &validation); err != nil {
		t.Fatal(err)
	}
	if !validation.Valid {
		t.Fatalf("expected valid assignment, got %#v", validation)
	}

	createBody := `{
  "actor": {"name": "Alice", "email": "alice@example.test"},
  "sourceHid": "HAZ_1_1",
  "referenceUuid": "00000000-0000-4000-8000-000000000301",
  "commitId": "commit-reference-create"
}`
	recorder = performJSON(router, http.MethodPost, "/api/v1/references/assignments", createBody)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("create assignment status = %d body=%s", recorder.Code, recorder.Body.String())
	}

	recorder = performJSON(router, http.MethodGet, "/api/v1/references/assignments/HAZ_1_1", "")
	if recorder.Code != http.StatusOK {
		t.Fatalf("list assignments status = %d body=%s", recorder.Code, recorder.Body.String())
	}
	var assignments listResponse[referenceAssignmentResponse]
	if err := json.Unmarshal(recorder.Body.Bytes(), &assignments); err != nil {
		t.Fatal(err)
	}
	if assignments.Total != 1 || assignments.Items[0].ReferenceItem.ExternalID != "T1000" {
		t.Fatalf("unexpected assignments response: %#v", assignments)
	}

	recorder = performJSON(router, http.MethodDelete, "/api/v1/references/assignments", createBody)
	if recorder.Code != http.StatusOK {
		t.Fatalf("delete assignment status = %d body=%s", recorder.Code, recorder.Body.String())
	}
}

func performJSON(handler http.Handler, method string, target string, body string) *httptest.ResponseRecorder {
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
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

func seedReferenceAssignmentGraph(t *testing.T, driver neo4j.DriverWithContext) {
	t.Helper()
	ctx := context.Background()
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
CREATE (:Hazard:SSTPANode {
  HID: "HAZ_1_1",
  uuid: "00000000-0000-4000-8000-000000000300",
  TypeName: "Hazard",
  Owner: "Alice",
  OwnerEmail: "alice@example.test",
  Creator: "Alice",
  CreatorEmail: "alice@example.test",
  Created: "2026-04-24T12:00:00Z",
  LastTouch: "2026-04-24T12:00:00Z",
  VersionID: "test-schema",
  Name: "Hazard"
})
CREATE (:ReferenceFramework {
  FrameworkName: "MITRE ATT&CK",
  FrameworkVersion: "v15"
})
CREATE (:ReferenceItem:AttackTechnique {
  uuid: "00000000-0000-4000-8000-000000000301",
  FrameworkName: "MITRE ATT&CK",
  FrameworkVersion: "v15",
  ExternalID: "T1000",
  ExternalType: "Technique",
  Name: "Example Technique",
  ShortDescription: "Technique summary",
  LongDescription: "Technique detail",
  SourceURI: "https://example.test/T1000"
})
`, nil)
		if err != nil {
			return nil, err
		}
		_, err = result.Consume(ctx)
		return nil, err
	})
	if err != nil {
		t.Fatal(err)
	}
}

type fakeConfiguredDriver struct {
	neo4j.DriverWithContext
}
