# Slice 04 — Backend REST API (Core + Messaging + Reference endpoints) (Phase B, Outline)

**Goal:** Expose the REST API required by SRS §2.2.10 and §2.2.10.10. Every route is HTTPS-terminated at Caddy, returns JSON, uses the S03 mutation layer for writes, and enforces pagination and validation rules. Telemetry middleware (OTel + Prometheus) is wired here.

**Architecture:** Existing `backend/internal/http` package expands into sub-packages: `http/nodes`, `http/hierarchy`, `http/search`, `http/context`, `http/validation`, `http/messages`, `http/reference`. Handlers are thin — they validate input, call services in `backend/internal/services/*`, and marshal responses. Middleware chain: request ID → OTel trace → Prom counter/histogram → user identification (stub) → routing. Pagination uses a `page_token` cursor derived from the result's last `(HID, uuid)` tuple (stable, collision-free).

**Tech Stack:** `chi/v5`, `go.opentelemetry.io/otel`, `github.com/prometheus/client_golang`, existing JSON stdlib.

**Pre-reads:**
- SRS §2.2.10 General Requirements through §2.2.10.10.6 Performance Requirements
- `backend/internal/mutation/*` (S03 output)
- `backend/internal/http/router.go`, `health.go`, `health_test.go`
- `infra/docker/caddy/Caddyfile`, `infra/docker/compose.yaml`
- `infra/docker/otel-collector/config.yaml`, `prometheus.yml`
- Conventions doc.

**Invariants:**
- Every list endpoint paginates (`?cursor=...&limit=...`, max 200).
- Every write uses `mutation.Apply`; no direct driver writes from handlers.
- Every response body is JSON; error body is `{"error": {"code":"...","message":"..."}}`.
- Validation endpoints return `{"valid": bool, "reasons": []string}`.
- Never exposes Neo4j; Bolt port stays on the backend network.
- OTel span around every handler; Prom histogram per route.

**Tasks:**

### S04-T01: HTTP middleware chain (request ID, OTel, Prom, recover)
Difficulty: medium · Files: `backend/internal/http/middleware/*.go` · Deps: none · Tests: handler round-trip.

### S04-T02: Error and pagination helpers
Difficulty: low · Files: `backend/internal/http/httpx/{error.go,pagination.go,*_test.go}` · Deps: S04-T01 · Tests: encode/decode cursor round-trip; malformed cursor returns 400.

### S04-T03: User identification stub
Difficulty: low · Files: `backend/internal/http/middleware/user.go` · Deps: S04-T01 · Tests: requests with `X-SSTPA-User: ada@example.com` populate a `user.Current(ctx)`; missing header → 401. This is the placeholder per SRS §2.2.10.9.

### S04-T04: `GET /api/v1/nodes/by-hid/{hid}` and `/by-uuid/{uuid}` and `/by-type/{type}`
Difficulty: medium · Files: `backend/internal/http/nodes/*.go` and service `backend/internal/services/nodes/*.go` · Deps: S04-T02 · Tests: 200 with full properties; 404 when not found; `by-type` paginates.

### S04-T05: `GET /api/v1/hierarchy` + `/api/v1/hierarchy/{capabilityHID}`
Difficulty: medium · Files: `backend/internal/http/hierarchy/*.go` · Deps: S04-T04 · Tests: returns compact hierarchy (HID, Name, TypeName, parent HID) with payload size under a threshold for a seeded hierarchy of 100 Systems.

### S04-T06: `GET /api/v1/search`
Difficulty: medium · Files: `backend/internal/http/search/*.go` · Deps: S04-T04 · Tests: HID/uuid exact faster than text; Name/ShortDescription partial; type filtering; response includes containing SoI.

### S04-T07: `GET /api/v1/nodes/{hid}/context`
Difficulty: low · Files: `backend/internal/http/context/*.go` · Deps: S04-T04 · Tests: returns containing System, path through hierarchy, parent relationships.

### S04-T08: `POST /api/v1/validate/relationship`
Difficulty: medium · Files: `backend/internal/http/validation/*.go` · Deps: S04-T02 · Tests: confirms allowed node types per `graph.Catalog()`, enforces relationship rules, rejects duplicates; returns `{valid, reasons[]}` per SRS §2.2.10.5.

### S04-T09: Mutation write endpoints
Difficulty: high · Integration Checkpoint: yes · Files: `backend/internal/http/nodes/create.go`, `update.go`, `delete.go`, `relationships.go` · Deps: S04-T08 · Tests: create node; update properties; delete; create relationship. Every write returns the `CommitReport` from S03 including notification recipients.

### S04-T10: Messaging endpoints per SRS §2.2.10.10.1 messaging
Difficulty: medium · Files: `backend/internal/http/messages/*.go` · Deps: S04-T02, S03 messaging · Tests: list (sort + asc/desc), detail, send, reply, mark read, delete, unread-count.

### S04-T11: Reference framework endpoints per SRS §2.2.10.10.{2-6}
Difficulty: medium · Files: `backend/internal/http/reference/*.go` · Deps: S04-T04 · Tests: list frameworks, list versions, get by ExternalID/uuid, search with framework/version/type filters, hierarchy, related items. NO imports here — importing happens via the CLI in S05.

### S04-T12: `[:REFERENCES]` assignment endpoints per SRS §2.2.10.10.5
Difficulty: medium · Files: `backend/internal/http/reference/assignments.go` · Deps: S04-T09, S04-T11 · Tests: create/delete/list `[:REFERENCES]` with mutation layer; validate against allowed assignments table (SRS §1.5.5).

### S04-T13: Telemetry wiring end-to-end
Difficulty: medium · Integration Checkpoint: yes · Files: `backend/internal/telemetry/*.go`; modify `cmd/api/main.go`; modify Prometheus scrape target, Caddy route (no change expected) · Deps: S04-T01 · Tests: spin up backend + OTel + Tempo + Prometheus via compose test; assert `sstpa_http_requests_total` incremented and a trace reached Tempo.

### S04-T14: OpenAPI 3.1 spec committed
Difficulty: medium · Files: `backend/api/openapi.yaml`; `backend/internal/http/spec_test.go` (snapshot verifying route list matches the spec) · Deps: all routes · Tests: spec file exists, routes match registered routes.

### S04-T15: Slice integration gate
Difficulty: low · Integration Checkpoint: yes · Files: verification docs + Evidence Summary. Move `SSTPA-PAGINATION-001` to Approved; add/approve Reqs for node retrieval, hierarchy, search, validation, messaging, reference endpoints.

**Integration gate criteria:**
- Contract tests cover every route.
- `make verify` green.
- OpenAPI spec parses.
- Compose-up demo: 200 on `/api/v1/hierarchy` with seeded data.
