# S04b — Users + Admins Onboarding & Backend Telemetry Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Close two S04 loose ends so slice S06 (Frontend Foundation) can proceed: (1) add Users and Admins onboarding endpoints that the Startup Software launcher needs, and (2) wire OpenTelemetry traces and Prometheus metrics through the chi router and around Neo4j transactions so the docker-compose observability stack (OTel Collector → Tempo, Prometheus, Grafana) actually receives data.

**Architecture:** `(:User:SSTPANode)` and `(:Admin:SSTPANode)` become first-class graph nodes linked from bootstrapped `(:Users)` and `(:Admins)` singleton containers via `HAS_USER` / `HAS_ADMIN` edges. A new `backend/internal/onboarding` package owns create/list/get flows parameterised on a `Kind` record so `:User` and `:Admin` share one implementation. Thin chi handlers under `backend/internal/http/{users,admins}.go` delegate to the package. Telemetry lands in a new `backend/internal/telemetry` package: a Prometheus registry with default HTTP instruments, an OTel tracer provider that exports OTLP/HTTP to the compose collector, and a chi middleware that combines `otelhttp.NewHandler` and Prometheus instrumentation using chi route patterns (not raw URLs) to bound cardinality. `backend/cmd/api/main.go` initialises telemetry, passes a registry and tracer into the router, and defers a tracer shutdown. The mutation layer opens an explicit span around each `session.ExecuteWrite` call and records commit metadata.

**Tech Stack:** Go 1.24, `github.com/go-chi/chi/v5`, `github.com/neo4j/neo4j-go-driver/v5`, `github.com/prometheus/client_golang/prometheus`, `go.opentelemetry.io/otel`, `go.opentelemetry.io/otel/sdk/trace`, `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp`, `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp`. Frontend: `@sstpa/api-client` TypeScript workspace extended with Vitest tests.

**SRS references:** §1.3.6, §1.4.1 Onboarding, §1.4.2 Admin Data, §1.4.3 User Data, §1.4.4 Messaging Data Model, §2.1 Startup Software, §2.2 Backend diagram (Prom/OTel instr.), §2.2.10 Backend API Requirements.

**Out of scope (MVP-aligned):** passwords/authentication (SRS §1.4.2 explicitly defers), role-based policy enforcement (§2.2.10.9 is placeholder), Grafana dashboards (S12), Neo4j bolt driver instrumentation beyond the mutation-layer span, remote backend support.

---

## File Structure

**Create**

- `backend/internal/onboarding/onboarding.go` — `Kind`, `Record`, `CreateInput`, `Create`, `List`, `GetByUUID`; owns the Cypher.
- `backend/internal/onboarding/onboarding_test.go` — unit tests + Testcontainers integration.
- `backend/internal/http/users.go` — chi handlers mapped to `onboarding.UserKind`.
- `backend/internal/http/admins.go` — chi handlers mapped to `onboarding.AdminKind`.
- `backend/internal/http/users_test.go` — integration tests covering both endpoint families.
- `backend/internal/telemetry/metrics.go` — Prometheus registry + default HTTP instruments.
- `backend/internal/telemetry/metrics_test.go` — unit tests verifying `/metrics` text output.
- `backend/internal/telemetry/tracer.go` — OTel tracer provider factory + `Shutdown`.
- `backend/internal/telemetry/tracer_test.go` — unit tests with the tracetest recorder.
- `backend/internal/telemetry/middleware.go` — chi middleware combining otelhttp + Prometheus.
- `backend/internal/telemetry/middleware_test.go` — unit tests asserting span creation + metrics increment.
- `backend/internal/config/config_test.go` — tests for new telemetry env knobs.

**Modify**

- `backend/internal/identity/types.go` — add `NodeTypeUser`, `NodeTypeAdmin`, HID prefixes `USR`, `ADM`, ordered entries.
- `backend/internal/identity/types_test.go` — bump expected count 27 → 29, rename test.
- `backend/internal/schema/bootstrap.go` — MERGE `:Users` and `:Admins` singletons; add `admin_email_index`; keep existing indexes.
- `backend/internal/schema/bootstrap_test.go` — assert new statements and singleton creation.
- `backend/internal/http/router.go` — install telemetry middleware, serve `/metrics`, register `/users` + `/admins` routes; accept `*prometheus.Registry` + `trace.Tracer` in `RouterOptions`.
- `backend/internal/http/api.go` — `api` struct carries `tracer` + `registry`; `newAPI` wires them.
- `backend/internal/http/openapi.go` — add `/users`, `/users/{uuid}`, `/admins`, `/admins/{uuid}` paths.
- `backend/internal/mutation/apply.go` — open a span around `session.ExecuteWrite`; add attributes.
- `backend/internal/config/config.go` — `OTLPEndpoint`, `ServiceName`, `MetricsEnabled`, `TracingEnabled` fields.
- `backend/cmd/api/main.go` — build telemetry, pass into router, defer tracer shutdown.
- `backend/go.mod` — promote existing OTel indirect deps to direct and add `prometheus/client_golang`, `otel/sdk`, `otlptracehttp`.
- `docs/api/openapi.yaml` — mirror the new endpoints (keep parity with the inline Go spec).
- `packages/api-client/src/index.ts` — `listUsers`, `createUser`, `getUser`, plus admin equivalents; new types.
- `packages/api-client/src/index.test.ts` — add two tests.
- `docs/verification/shall-register.md` — add `SSTPA-API-USER-001`, `SSTPA-API-ADMIN-001`, `SSTPA-USER-MODEL-001`, `SSTPA-TELEMETRY-HTTP-001`, `SSTPA-TELEMETRY-NEO4J-001`.
- `docs/verification/verification-matrix.md` — one row per new SHALL.

---

## Conventions the engineer must follow

- **Copyright banner.** Every new `.go`, `.ts`, `.tsx`, `.mjs`, `.cjs`, `.rs` source file begins with the five-line banner exactly as in existing files:

```
// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
```

`make copyright-check` fails otherwise.

- **Commit style.** `<type>(<area>): <subject>` (e.g. `feat(backend): register User and Admin node types`). Every commit body cites the SRS section and ends with the Co-Authored-By trailer. Examples in `git log`.
- **One task = one commit.** Don't batch multiple tasks.
- **Branch.** Work on `slice/04b-users-admins-telemetry` created from current `codex/sstpa-mvp-foundation` HEAD. Do not push to `main`.
- **Tests first.** For each behavior: write failing test → run (verify fail) → implement → run (verify pass) → commit.
- **Bootstrap on a fresh machine.** Before the final `make verify`: run `make bootstrap` (npm install) and `cd backend && go mod download` so Go modules and node_modules are populated.
- **Go module imports.** When adding a new dependency, run `cd backend && go mod tidy` after the import is in code so `go.sum` updates deterministically.
- **The `api` test package is named `apihttp`** (because `backend/internal/http` is the import path; `package apihttp` avoids colliding with the stdlib). Match it.
- **Integration tests that need Neo4j use `testhelpers.StartNeo4j(t)`.** That helper skips when Docker is unavailable; do not convert those tests to fail-open without Docker.

---

## Task 1 — Register `User` and `Admin` node types

**SRS:** §1.3.6, §1.4.2, §1.4.3.

**Files:**
- Modify: `backend/internal/identity/types.go`
- Modify: `backend/internal/identity/types_test.go`

- [ ] **Step 1: Update the count assertion to expect 29 types**

Replace the `TestAllTypesCountMatchesSRS` function in `backend/internal/identity/types_test.go` (lines 35-39) with:

```go
func TestAllTypesCountCoversCoreAndToolData(t *testing.T) {
	// 27 Core Data Model types (SRS §1.3.6.1) + 2 Tool Data types (User, Admin per §1.4).
	if got := len(AllTypes()); got != 29 {
		t.Fatalf("expected 29 node types (27 Core + 2 Tool Data), got %d", got)
	}
}
```

- [ ] **Step 2: Add a test that asserts the new type IDs**

Append to `backend/internal/identity/types_test.go`:

```go
func TestUserAndAdminTypeIDs(t *testing.T) {
	cases := map[NodeType]string{
		NodeTypeUser:  "USR",
		NodeTypeAdmin: "ADM",
	}

	for nodeType, want := range cases {
		got, ok := TypeID(nodeType)
		if !ok || got != want {
			t.Errorf("TypeID(%q) = (%q, %v), want (%q, true)", nodeType, got, ok, want)
		}
	}
}
```

- [ ] **Step 3: Run the tests to confirm they fail**

Run: `cd backend && go test ./internal/identity/... -run "TestAllTypesCountCoversCoreAndToolData|TestUserAndAdminTypeIDs" -v`
Expected: both FAIL — compile error on `NodeTypeUser`/`NodeTypeAdmin`.

- [ ] **Step 4: Add User and Admin constants plus catalog entries**

In `backend/internal/identity/types.go`, append `NodeTypeUser` and `NodeTypeAdmin` to the `const` block (after `NodeTypeAttack`):

```go
	NodeTypeUser  NodeType = "User"
	NodeTypeAdmin NodeType = "Admin"
```

Append to the `typeIDs` map:

```go
	NodeTypeUser:  "USR",
	NodeTypeAdmin: "ADM",
```

Append to `orderedTypes`:

```go
	NodeTypeUser,
	NodeTypeAdmin,
```

- [ ] **Step 5: Re-run the identity tests**

Run: `cd backend && go test ./internal/identity/...`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/identity/types.go backend/internal/identity/types_test.go
git commit -m "$(cat <<'EOF'
feat(backend): register User and Admin node types

Adds :User and :Admin to the identity catalog with HID prefixes USR and
ADM per SRS §1.3.6 + §1.4.2/§1.4.3. Unlocks HID allocation for the
Users and Admins onboarding flow referenced by §1.4.1 and §2.1.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 2 — Bootstrap `:Users` and `:Admins` singletons and their indexes

**SRS:** §1.4.2 "Tool Data will parent a single node called Admins …", §1.4.3 "Tool Data will parent a single node called Users …".

**Files:**
- Modify: `backend/internal/schema/bootstrap.go`
- Modify: `backend/internal/schema/bootstrap_test.go`

- [ ] **Step 1: Add test that asserts the new bootstrap statements**

Append to `backend/internal/schema/bootstrap_test.go`:

```go
func TestStatementsIncludeUserAndAdminContainers(t *testing.T) {
	joined := strings.Join(Statements(), "\n")

	for _, want := range []string{
		"MERGE (:Users)",
		"MERGE (:Admins)",
		"user_uuid_index",
		"admin_email_index",
		"admin_uuid_index",
	} {
		if !strings.Contains(joined, want) {
			t.Fatalf("expected bootstrap statements to include %q", want)
		}
	}
}
```

- [ ] **Step 2: Run the test and confirm it fails**

Run: `cd backend && go test ./internal/schema/... -run TestStatementsIncludeUserAndAdminContainers`
Expected: FAIL — substrings not found.

- [ ] **Step 3: Extend the bootstrap statements**

Replace the `bootstrapStatements` slice in `backend/internal/schema/bootstrap.go` with:

```go
var bootstrapStatements = []string{
	"CREATE INDEX node_hid_index IF NOT EXISTS FOR (n:SSTPANode) ON (n.HID)",
	"CREATE INDEX node_uuid_index IF NOT EXISTS FOR (n:SSTPANode) ON (n.uuid)",
	"CREATE INDEX node_name_index IF NOT EXISTS FOR (n:SSTPANode) ON (n.Name)",
	"CREATE INDEX node_type_index IF NOT EXISTS FOR (n:SSTPANode) ON (n.TypeName)",
	"CREATE INDEX user_email_index IF NOT EXISTS FOR (u:User) ON (u.UserEmail)",
	"CREATE INDEX user_uuid_index IF NOT EXISTS FOR (u:User) ON (u.uuid)",
	"CREATE INDEX admin_email_index IF NOT EXISTS FOR (a:Admin) ON (a.UserEmail)",
	"CREATE INDEX admin_uuid_index IF NOT EXISTS FOR (a:Admin) ON (a.uuid)",
	"CREATE INDEX mailbox_id_index IF NOT EXISTS FOR (m:Mailbox) ON (m.MailboxID)",
	"CREATE INDEX message_uuid_index IF NOT EXISTS FOR (m:Message) ON (m.uuid)",
	"MERGE (:Users)",
	"MERGE (:Admins)",
}
```

(The existing `TestStatementsAreIdempotent` tolerates `MERGE` because `MERGE` itself is idempotent — add `MERGE` to the idempotency allowlist in Step 4.)

- [ ] **Step 4: Loosen the idempotency test to permit MERGE**

Replace `TestStatementsAreIdempotent` in `backend/internal/schema/bootstrap_test.go` with:

```go
func TestStatementsAreIdempotent(t *testing.T) {
	for _, statement := range Statements() {
		if strings.Contains(statement, "IF NOT EXISTS") || strings.HasPrefix(statement, "MERGE ") {
			continue
		}
		t.Fatalf("statement must be idempotent: %s", statement)
	}
}
```

- [ ] **Step 5: Run the schema tests**

Run: `cd backend && go test ./internal/schema/...`
Expected: PASS.

- [ ] **Step 6: Add an integration smoke test that the singletons materialise**

Append to `backend/internal/schema/bootstrap_test.go`:

```go
func TestBootstrapMaterializesContainers(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()

	if err := Bootstrap(ctx, fixture.Driver, ""); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}

	session := fixture.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	for _, label := range []string{"Users", "Admins"} {
		query := "MATCH (n:" + label + ") RETURN count(n) AS c"
		result, err := session.Run(ctx, query, nil)
		if err != nil {
			t.Fatalf("run %s: %v", label, err)
		}
		record, err := result.Single(ctx)
		if err != nil {
			t.Fatalf("single %s: %v", label, err)
		}
		count, _ := record.Get("c")
		if count.(int64) != 1 {
			t.Fatalf("expected exactly one :%s node, got %v", label, count)
		}
	}
}
```

Add imports to the test file header:

```go
import (
	"context"
	"strings"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sstpa-tool/backend/internal/testhelpers"
)
```

- [ ] **Step 7: Run the full schema test suite (requires Docker)**

Run: `cd backend && go test ./internal/schema/... -v`
Expected: PASS. If Docker is not available the integration test skips automatically.

- [ ] **Step 8: Commit**

```bash
git add backend/internal/schema/bootstrap.go backend/internal/schema/bootstrap_test.go
git commit -m "$(cat <<'EOF'
feat(backend): bootstrap Users and Admins singleton containers

Adds idempotent MERGE of (:Users) and (:Admins) plus indexes on
User.uuid, Admin.UserEmail, Admin.uuid so onboarding lookups and
HID allocation stay O(1). Implements SRS §1.4.2 and §1.4.3.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 3 — Create the generic `onboarding` package

**SRS:** §1.4.1 Onboarding, §1.4.2, §1.4.3, §1.4.4 `:User` shape.

**Files:**
- Create: `backend/internal/onboarding/onboarding.go`
- Create: `backend/internal/onboarding/onboarding_test.go`

- [ ] **Step 1: Write a failing unit test for `Kind` defaults**

Create `backend/internal/onboarding/onboarding_test.go`:

```go
// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package onboarding

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sstpa-tool/backend/internal/identity"
	"sstpa-tool/backend/internal/metadata"
	"sstpa-tool/backend/internal/schema"
	"sstpa-tool/backend/internal/testhelpers"
)

func TestUserKindShape(t *testing.T) {
	if UserKind.NodeType != identity.NodeTypeUser {
		t.Fatalf("UserKind.NodeType = %q, want %q", UserKind.NodeType, identity.NodeTypeUser)
	}
	if UserKind.ContainerLabel != "Users" {
		t.Fatalf("UserKind.ContainerLabel = %q, want Users", UserKind.ContainerLabel)
	}
	if UserKind.Relationship != "HAS_USER" {
		t.Fatalf("UserKind.Relationship = %q, want HAS_USER", UserKind.Relationship)
	}
}

func TestAdminKindShape(t *testing.T) {
	if AdminKind.NodeType != identity.NodeTypeAdmin {
		t.Fatalf("AdminKind.NodeType = %q, want %q", AdminKind.NodeType, identity.NodeTypeAdmin)
	}
	if AdminKind.ContainerLabel != "Admins" {
		t.Fatalf("AdminKind.ContainerLabel = %q, want Admins", AdminKind.ContainerLabel)
	}
	if AdminKind.Relationship != "HAS_ADMIN" {
		t.Fatalf("AdminKind.Relationship = %q, want HAS_ADMIN", AdminKind.Relationship)
	}
}
```

- [ ] **Step 2: Run the tests to confirm they fail (compile error)**

Run: `cd backend && go test ./internal/onboarding/...`
Expected: FAIL — package does not exist.

- [ ] **Step 3: Implement the package**

Create `backend/internal/onboarding/onboarding.go`:

```go
// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package onboarding

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sstpa-tool/backend/internal/identity"
	"sstpa-tool/backend/internal/metadata"
)

type Kind struct {
	NodeType       identity.NodeType
	NodeLabel      string
	ContainerLabel string
	Relationship   string
}

var UserKind = Kind{
	NodeType:       identity.NodeTypeUser,
	NodeLabel:      "User",
	ContainerLabel: "Users",
	Relationship:   "HAS_USER",
}

var AdminKind = Kind{
	NodeType:       identity.NodeTypeAdmin,
	NodeLabel:      "Admin",
	ContainerLabel: "Admins",
	Relationship:   "HAS_ADMIN",
}

type Record struct {
	HID       string `json:"hid"`
	UUID      string `json:"uuid"`
	TypeName  string `json:"typeName"`
	UserName  string `json:"userName"`
	UserEmail string `json:"userEmail"`
	Created   string `json:"created"`
	LastTouch string `json:"lastTouch"`
}

type CreateInput struct {
	UserName  string
	UserEmail string
	Actor     metadata.Actor
	Now       time.Time
}

type ListResult struct {
	Items []Record `json:"items"`
	Page  int      `json:"page"`
	Limit int      `json:"limit"`
	Total int64    `json:"total"`
}

type Page struct {
	Page   int
	Limit  int
	Offset int
}

var ErrAlreadyRegistered = errors.New("user already registered")
var ErrNotFound = errors.New("not found")

func Create(ctx context.Context, driver neo4j.DriverWithContext, databaseName string, kind Kind, input CreateInput) (Record, error) {
	if driver == nil {
		return Record{}, errors.New("neo4j driver is required")
	}
	if input.UserName == "" || input.UserEmail == "" {
		return Record{}, errors.New("UserName and UserEmail are required")
	}
	if input.Actor.Name == "" || input.Actor.Email == "" {
		return Record{}, errors.New("actor name and email are required")
	}

	typeID, ok := identity.TypeID(kind.NodeType)
	if !ok {
		return Record{}, fmt.Errorf("unknown node type %q", kind.NodeType)
	}

	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: databaseName})
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		existing, err := scalarInt(ctx, tx,
			"MATCH (n:"+kind.NodeLabel+":SSTPANode {UserEmail: $email}) RETURN count(n) AS c",
			map[string]any{"email": input.UserEmail})
		if err != nil {
			return Record{}, err
		}
		if existing > 0 {
			return Record{}, ErrAlreadyRegistered
		}

		nextSeq, err := scalarInt(ctx, tx,
			"MATCH (n:"+kind.NodeLabel+":SSTPANode) RETURN coalesce(max(n.HIDSequence), 0) + 1 AS c",
			nil)
		if err != nil {
			return Record{}, err
		}

		hid, err := identity.FormatHID(typeID, "", int(nextSeq))
		if err != nil {
			return Record{}, err
		}

		uuid := identity.NewUUID()
		common, err := metadata.NewCommon(metadata.NewCommonInput{
			NodeType:  kind.NodeType,
			HID:       hid,
			UUID:      uuid,
			Actor:     input.Actor,
			Now:       now,
			VersionID: "",
		})
		if err != nil {
			return Record{}, err
		}

		props := common.Properties()
		props["Name"] = input.UserName
		props["UserName"] = input.UserName
		props["UserEmail"] = input.UserEmail
		props["UserHash"] = input.UserEmail
		props["HIDSequence"] = nextSeq

		createCypher := fmt.Sprintf(`
MATCH (container:%s)
CREATE (n:%s:SSTPANode)
SET n = $props
MERGE (container)-[:%s]->(n)
RETURN n.HID AS hid, n.uuid AS uuid, n.TypeName AS typeName,
       n.UserName AS userName, n.UserEmail AS userEmail,
       n.Created AS created, n.LastTouch AS lastTouch
`, kind.ContainerLabel, kind.NodeLabel, kind.Relationship)

		row, err := tx.Run(ctx, createCypher, map[string]any{"props": props})
		if err != nil {
			return Record{}, err
		}
		record, err := row.Single(ctx)
		if err != nil {
			return Record{}, err
		}
		return recordFromRow(record), nil
	})
	if err != nil {
		return Record{}, err
	}

	out, ok := result.(Record)
	if !ok {
		return Record{}, errors.New("unexpected onboarding result")
	}
	return out, nil
}

func List(ctx context.Context, driver neo4j.DriverWithContext, databaseName string, kind Kind, page Page) (ListResult, error) {
	if driver == nil {
		return ListResult{}, errors.New("neo4j driver is required")
	}
	if page.Limit <= 0 {
		page.Limit = 50
	}

	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: databaseName})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		rows, err := tx.Run(ctx, `
MATCH (n:`+kind.NodeLabel+`:SSTPANode)
RETURN n.HID AS hid, n.uuid AS uuid, n.TypeName AS typeName,
       n.UserName AS userName, n.UserEmail AS userEmail,
       n.Created AS created, n.LastTouch AS lastTouch
ORDER BY n.HID
SKIP $skip LIMIT $limit
`, map[string]any{"skip": page.Offset, "limit": page.Limit})
		if err != nil {
			return ListResult{}, err
		}
		collected, err := rows.Collect(ctx)
		if err != nil {
			return ListResult{}, err
		}

		items := make([]Record, 0, len(collected))
		for _, record := range collected {
			items = append(items, recordFromRow(record))
		}

		total, err := scalarInt(ctx, tx,
			"MATCH (n:"+kind.NodeLabel+":SSTPANode) RETURN count(n) AS c",
			nil)
		if err != nil {
			return ListResult{}, err
		}

		pageNumber := page.Page
		if pageNumber < 1 {
			pageNumber = 1
		}
		return ListResult{Items: items, Page: pageNumber, Limit: page.Limit, Total: total}, nil
	})
	if err != nil {
		return ListResult{}, err
	}

	out, ok := result.(ListResult)
	if !ok {
		return ListResult{}, errors.New("unexpected list result")
	}
	return out, nil
}

func GetByUUID(ctx context.Context, driver neo4j.DriverWithContext, databaseName string, kind Kind, uuid string) (Record, error) {
	if driver == nil {
		return Record{}, errors.New("neo4j driver is required")
	}
	if uuid == "" {
		return Record{}, errors.New("uuid is required")
	}

	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: databaseName})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		row, err := tx.Run(ctx, `
MATCH (n:`+kind.NodeLabel+`:SSTPANode {uuid: $uuid})
RETURN n.HID AS hid, n.uuid AS uuid, n.TypeName AS typeName,
       n.UserName AS userName, n.UserEmail AS userEmail,
       n.Created AS created, n.LastTouch AS lastTouch
LIMIT 1
`, map[string]any{"uuid": uuid})
		if err != nil {
			return Record{}, err
		}
		record, err := row.Single(ctx)
		if err != nil {
			return Record{}, ErrNotFound
		}
		return recordFromRow(record), nil
	})
	if err != nil {
		return Record{}, err
	}
	out, ok := result.(Record)
	if !ok {
		return Record{}, errors.New("unexpected get result")
	}
	return out, nil
}

func recordFromRow(record *neo4j.Record) Record {
	get := func(key string) string {
		value, _ := record.Get(key)
		text, _ := value.(string)
		return text
	}
	return Record{
		HID:       get("hid"),
		UUID:      get("uuid"),
		TypeName:  get("typeName"),
		UserName:  get("userName"),
		UserEmail: get("userEmail"),
		Created:   get("created"),
		LastTouch: get("lastTouch"),
	}
}

func scalarInt(ctx context.Context, tx neo4j.ManagedTransaction, query string, params map[string]any) (int64, error) {
	result, err := tx.Run(ctx, query, params)
	if err != nil {
		return 0, err
	}
	record, err := result.Single(ctx)
	if err != nil {
		return 0, err
	}
	value, ok := record.Get("c")
	if !ok {
		return 0, errors.New("c scalar not returned")
	}
	count, ok := value.(int64)
	if !ok {
		return 0, fmt.Errorf("c has unexpected type %T", value)
	}
	return count, nil
}
```

- [ ] **Step 4: Run the package unit tests**

Run: `cd backend && go test ./internal/onboarding/... -run "Kind"`
Expected: PASS on `TestUserKindShape` and `TestAdminKindShape`.

- [ ] **Step 5: Add an integration test covering Create, List, GetByUUID, and duplicate rejection**

Append to `backend/internal/onboarding/onboarding_test.go`:

```go
func TestCreateListGetUser(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	if err := schema.Bootstrap(ctx, fixture.Driver, ""); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}

	installer := metadata.Actor{Name: "Installer", Email: "installer@example.test", Admin: true}
	now := time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)

	record, err := Create(ctx, fixture.Driver, "", UserKind, CreateInput{
		UserName:  "Alice Analyst",
		UserEmail: "alice@example.test",
		Actor:     installer,
		Now:       now,
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if !strings.HasPrefix(record.HID, "USR__") {
		t.Fatalf("expected HID to start with USR__, got %q", record.HID)
	}
	if record.TypeName != "User" || record.UserEmail != "alice@example.test" {
		t.Fatalf("unexpected record: %#v", record)
	}

	if _, err := Create(ctx, fixture.Driver, "", UserKind, CreateInput{
		UserName:  "Alice Again",
		UserEmail: "alice@example.test",
		Actor:     installer,
		Now:       now,
	}); err == nil {
		t.Fatal("expected duplicate registration to error")
	}

	list, err := List(ctx, fixture.Driver, "", UserKind, Page{Page: 1, Limit: 10, Offset: 0})
	if err != nil {
		t.Fatalf("list users: %v", err)
	}
	if list.Total != 1 || len(list.Items) != 1 {
		t.Fatalf("unexpected list: %#v", list)
	}

	fetched, err := GetByUUID(ctx, fixture.Driver, "", UserKind, record.UUID)
	if err != nil {
		t.Fatalf("get user: %v", err)
	}
	if fetched.HID != record.HID {
		t.Fatalf("fetched HID = %q, want %q", fetched.HID, record.HID)
	}
}

func TestCreateListGetAdmin(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	if err := schema.Bootstrap(ctx, fixture.Driver, ""); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}

	installer := metadata.Actor{Name: "Installer", Email: "installer@example.test", Admin: true}
	now := time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)

	record, err := Create(ctx, fixture.Driver, "", AdminKind, CreateInput{
		UserName:  "Root Admin",
		UserEmail: "root@example.test",
		Actor:     installer,
		Now:       now,
	})
	if err != nil {
		t.Fatalf("create admin: %v", err)
	}
	if !strings.HasPrefix(record.HID, "ADM__") {
		t.Fatalf("expected HID to start with ADM__, got %q", record.HID)
	}

	list, err := List(ctx, fixture.Driver, "", AdminKind, Page{Page: 1, Limit: 10, Offset: 0})
	if err != nil {
		t.Fatalf("list admins: %v", err)
	}
	if list.Total != 1 {
		t.Fatalf("unexpected admin list: %#v", list)
	}
}
```

- [ ] **Step 6: Run the integration tests (Docker required)**

Run: `cd backend && go test ./internal/onboarding/... -v`
Expected: PASS on both integration tests (or SKIP if Docker is unavailable).

- [ ] **Step 7: Commit**

```bash
git add backend/internal/onboarding
git commit -m "$(cat <<'EOF'
feat(backend): add onboarding package for User and Admin registration

Provides Create/List/GetByUUID parameterised on a Kind so :User and
:Admin share one transactional flow. Nodes carry SSTPA common metadata
and are linked from the bootstrapped :Users / :Admins containers via
HAS_USER / HAS_ADMIN. Implements SRS §1.4.1–§1.4.4.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 4 — Expose `/api/v1/users` REST endpoints

**SRS:** §2.1 Startup Software (list + add User), §2.2.10.1 General Requirements (JSON, transactional).

**Files:**
- Create: `backend/internal/http/users.go`
- Modify: `backend/internal/http/router.go`

- [ ] **Step 1: Register the new routes in the router**

In `backend/internal/http/router.go`, inside the `/api/v1` `Route` block (after the `Delete("/references/assignments", ...)` line 60 block), append:

```go
		group.Get("/users", api.listOnboardingHandler(onboarding.UserKind))
		group.Post("/users", api.createOnboardingHandler(onboarding.UserKind))
		group.Get("/users/{uuid}", api.getOnboardingHandler(onboarding.UserKind))
		group.Get("/admins", api.listOnboardingHandler(onboarding.AdminKind))
		group.Post("/admins", api.createOnboardingHandler(onboarding.AdminKind))
		group.Get("/admins/{uuid}", api.getOnboardingHandler(onboarding.AdminKind))
```

Add `"sstpa-tool/backend/internal/onboarding"` to the import block at the top.

- [ ] **Step 2: Create the handler file**

Create `backend/internal/http/users.go`:

```go
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
			Page:   page.Page,
			Limit:  page.Limit,
			Offset: page.Offset,
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
```

Add the `metadata` import to the top of the file:

```go
	"sstpa-tool/backend/internal/metadata"
```

- [ ] **Step 3: Build the backend to confirm the router compiles**

Run: `cd backend && go build ./...`
Expected: builds cleanly.

- [ ] **Step 4: Commit**

```bash
git add backend/internal/http/users.go backend/internal/http/router.go
git commit -m "$(cat <<'EOF'
feat(backend): expose /api/v1/users and /api/v1/admins endpoints

Registers list/create/get handlers for User and Admin onboarding. Bodies
are JSON; creation returns the created record with HID/uuid. Implements
SRS §2.1 Startup Software user-list requirement and §1.4.1 onboarding
screens.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 5 — Add HTTP integration tests for `/users` and `/admins`

**SRS:** §2.1, §2.2.10.1.

**Files:**
- Create: `backend/internal/http/users_test.go`

- [ ] **Step 1: Write integration tests that exercise create, duplicate, list, and get**

Create `backend/internal/http/users_test.go`:

```go
// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package apihttp

import (
	"context"
	"encoding/json"
	"net/http"
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

	body := `{"userName": "Root Admin", "userEmail": "root@example.test"}`
	recorder := performJSON(router, http.MethodPost, "/api/v1/admins", body)
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

func TestUsersEndpointRejectsMissingFields(t *testing.T) {
	router := NewRouterWithOptions(RouterOptions{Version: "test", Driver: fakeConfiguredDriver{}})

	recorder := performJSON(router, http.MethodPost, "/api/v1/users", `{"userName": "Alice"}`)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing email, got %d", recorder.Code)
	}
}
```

- [ ] **Step 2: Run the HTTP integration tests (Docker required)**

Run: `cd backend && go test ./internal/http/... -run "Users|Admins" -v`
Expected: PASS (or SKIP the Testcontainers tests if Docker is unavailable; the "RejectsMissingFields" variant runs without Docker and must PASS).

- [ ] **Step 3: Commit**

```bash
git add backend/internal/http/users_test.go
git commit -m "$(cat <<'EOF'
test(backend): integration tests for /users and /admins endpoints

Covers create + duplicate + list + get-by-uuid for both onboarding
kinds plus an offline validation test for missing payload fields.
Exercises SRS §2.1 and §1.4.1 end-to-end through the chi router.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 6 — Publish the new endpoints in the OpenAPI surface

**SRS:** §2.2.10.1 Backend SHALL expose JSON REST API.

**Files:**
- Modify: `backend/internal/http/openapi.go`
- Modify: `docs/api/openapi.yaml`

- [ ] **Step 1: Write a failing test asserting the new paths are in the spec**

Append to `backend/internal/http/api_test.go` (inside the existing package):

```go
func TestOpenAPIIncludesUserAndAdminPaths(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/openapi.yaml", nil)
	recorder := httptest.NewRecorder()

	NewRouter("test").ServeHTTP(recorder, request)

	for _, want := range []string{"/users:", "/users/{uuid}:", "/admins:", "/admins/{uuid}:"} {
		if !strings.Contains(recorder.Body.String(), want) {
			t.Fatalf("expected OpenAPI spec to declare %q", want)
		}
	}
}
```

- [ ] **Step 2: Run the test and confirm it fails**

Run: `cd backend && go test ./internal/http/... -run TestOpenAPIIncludesUserAndAdminPaths`
Expected: FAIL.

- [ ] **Step 3: Extend the inline OpenAPI constant**

In `backend/internal/http/openapi.go`, insert the following block inside the `openAPISpec` raw string, immediately before the closing `components:` line:

```
  /users:
    get:
      parameters:
        - name: page
          in: query
          schema: { type: integer, minimum: 1 }
        - name: limit
          in: query
          schema: { type: integer, minimum: 1, maximum: 200 }
      responses:
        "200": { description: Paginated registered users. }
    post:
      responses:
        "201": { description: Newly registered user. }
        "409": { description: User email already registered. }
  /users/{uuid}:
    get:
      parameters:
        - name: uuid
          in: path
          required: true
          schema: { type: string }
      responses:
        "200": { description: Registered user. }
        "404": { description: User not found. }
  /admins:
    get:
      parameters:
        - name: page
          in: query
          schema: { type: integer, minimum: 1 }
        - name: limit
          in: query
          schema: { type: integer, minimum: 1, maximum: 200 }
      responses:
        "200": { description: Paginated registered admins. }
    post:
      responses:
        "201": { description: Newly registered admin. }
        "409": { description: Admin email already registered. }
  /admins/{uuid}:
    get:
      parameters:
        - name: uuid
          in: path
          required: true
          schema: { type: string }
      responses:
        "200": { description: Registered admin. }
        "404": { description: Admin not found. }
```

- [ ] **Step 4: Run the test to confirm it passes**

Run: `cd backend && go test ./internal/http/... -run TestOpenAPIIncludesUserAndAdminPaths`
Expected: PASS.

- [ ] **Step 5: Mirror the same paths in `docs/api/openapi.yaml`**

Open `docs/api/openapi.yaml`. Before the `components:` line, insert the same block as Step 3 (copy-paste verbatim).

- [ ] **Step 6: Commit**

```bash
git add backend/internal/http/openapi.go backend/internal/http/api_test.go docs/api/openapi.yaml
git commit -m "$(cat <<'EOF'
docs(api): publish /users and /admins in OpenAPI 3.1 contract

Mirrors the new REST surface in both the inline Go spec and the
committed docs/api/openapi.yaml so consumers have a single source of
truth per SRS §2.2.10.1.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 7 — Extend the `@sstpa/api-client` TypeScript package

**SRS:** §2.2.10.1.

**Files:**
- Modify: `packages/api-client/src/index.ts`
- Modify: `packages/api-client/src/index.test.ts`

- [ ] **Step 1: Add a failing test for the client's `createUser` and `listUsers` methods**

Append to `packages/api-client/src/index.test.ts`, inside the existing `describe("SSTPAClient", ...)` block (before the closing `})`):

```ts
  it("creates a user and lists registered users", async () => {
    const requests: Request[] = []
    const client = new SSTPAClient({
      baseUrl: "http://localhost:8080/api/v1/",
      fetchImpl: async (input, init) => {
        requests.push(new Request(input, init))
        if (init?.method === "POST") {
          return Response.json(
            {
              hid: "USR__1",
              uuid: "u-1",
              typeName: "User",
              userName: "Alice",
              userEmail: "alice@example.test",
              created: "2026-04-24T12:00:00Z",
              lastTouch: "2026-04-24T12:00:00Z",
            },
            { status: 201 },
          )
        }
        return Response.json({
          items: [
            {
              hid: "USR__1",
              uuid: "u-1",
              typeName: "User",
              userName: "Alice",
              userEmail: "alice@example.test",
              created: "2026-04-24T12:00:00Z",
              lastTouch: "2026-04-24T12:00:00Z",
            },
          ],
          page: 1,
          limit: 50,
          total: 1,
        })
      },
    })

    const created = await client.createUser({ userName: "Alice", userEmail: "alice@example.test" })
    expect(created.hid).toBe("USR__1")
    expect(requests[0].method).toBe("POST")
    expect(requests[0].url).toBe("http://localhost:8080/api/v1/users")
    await expect(requests[0].json()).resolves.toEqual({
      userName: "Alice",
      userEmail: "alice@example.test",
    })

    const list = await client.listUsers()
    expect(list.total).toBe(1)
    expect(list.items[0].userEmail).toBe("alice@example.test")
  })

  it("creates an admin via createAdmin", async () => {
    const requests: Request[] = []
    const client = new SSTPAClient({
      baseUrl: "http://localhost:8080/api/v1/",
      fetchImpl: async (input, init) => {
        requests.push(new Request(input, init))
        return Response.json(
          {
            hid: "ADM__1",
            uuid: "a-1",
            typeName: "Admin",
            userName: "Root",
            userEmail: "root@example.test",
            created: "2026-04-24T12:00:00Z",
            lastTouch: "2026-04-24T12:00:00Z",
          },
          { status: 201 },
        )
      },
    })

    const created = await client.createAdmin({ userName: "Root", userEmail: "root@example.test" })
    expect(created.hid).toBe("ADM__1")
    expect(requests[0].url).toBe("http://localhost:8080/api/v1/admins")
  })
```

- [ ] **Step 2: Run the tests to confirm they fail**

Run: `cd packages/api-client && npx vitest run`
Expected: FAIL — `client.createUser is not a function` (and `createAdmin`).

- [ ] **Step 3: Add the new types and methods**

In `packages/api-client/src/index.ts`, add the following interfaces after `interface ReferenceAssignmentMutationResponse` (before `class APIError`):

```ts
export interface OnboardingRecord {
  hid: string
  uuid: string
  typeName: string
  userName: string
  userEmail: string
  created: string
  lastTouch: string
}

export interface CreateOnboardingRequest {
  userName: string
  userEmail: string
}
```

Add the following methods inside the `SSTPAClient` class, after `deleteReferenceAssignment`:

```ts
  listUsers(params: { page?: number; limit?: number } = {}) {
    return this.request<ListResponse<OnboardingRecord>>(`/users${queryString(params)}`)
  }

  getUser(uuid: string) {
    return this.request<OnboardingRecord>(`/users/${encodeURIComponent(uuid)}`)
  }

  createUser(payload: CreateOnboardingRequest) {
    return this.request<OnboardingRecord>("/users", { method: "POST", body: payload })
  }

  listAdmins(params: { page?: number; limit?: number } = {}) {
    return this.request<ListResponse<OnboardingRecord>>(`/admins${queryString(params)}`)
  }

  getAdmin(uuid: string) {
    return this.request<OnboardingRecord>(`/admins/${encodeURIComponent(uuid)}`)
  }

  createAdmin(payload: CreateOnboardingRequest) {
    return this.request<OnboardingRecord>("/admins", { method: "POST", body: payload })
  }
```

- [ ] **Step 4: Run the vitest suite**

Run: `cd packages/api-client && npx vitest run`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add packages/api-client/src/index.ts packages/api-client/src/index.test.ts
git commit -m "$(cat <<'EOF'
feat(api-client): add user and admin onboarding methods

Exposes listUsers / getUser / createUser and listAdmins / getAdmin /
createAdmin on SSTPAClient, returning the shared OnboardingRecord shape.
Matches the Go handler contract added alongside and unblocks the S06
Startup Software launcher.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 8 — Add telemetry config knobs

**SRS:** §2.2 Backend diagram (Prom/OTel instr.), §2.2.7 Performance.

**Files:**
- Modify: `backend/internal/config/config.go`
- Create: `backend/internal/config/config_test.go`

- [ ] **Step 1: Write a failing test for the new config fields**

Create `backend/internal/config/config_test.go`:

```go
// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package config

import (
	"os"
	"testing"
)

func TestLoadDefaultTelemetrySettings(t *testing.T) {
	os.Unsetenv("SSTPA_OTLP_ENDPOINT")
	os.Unsetenv("SSTPA_SERVICE_NAME")
	os.Unsetenv("SSTPA_TELEMETRY_METRICS")
	os.Unsetenv("SSTPA_TELEMETRY_TRACING")

	cfg := Load()
	if cfg.OTLPEndpoint != "http://otel-collector:4318" {
		t.Fatalf("OTLPEndpoint default = %q", cfg.OTLPEndpoint)
	}
	if cfg.ServiceName != "sstpa-backend" {
		t.Fatalf("ServiceName default = %q", cfg.ServiceName)
	}
	if !cfg.MetricsEnabled {
		t.Fatalf("MetricsEnabled default must be true")
	}
	if !cfg.TracingEnabled {
		t.Fatalf("TracingEnabled default must be true")
	}
}

func TestLoadOverridesTelemetrySettings(t *testing.T) {
	t.Setenv("SSTPA_OTLP_ENDPOINT", "http://localhost:4318")
	t.Setenv("SSTPA_SERVICE_NAME", "sstpa-test")
	t.Setenv("SSTPA_TELEMETRY_METRICS", "false")
	t.Setenv("SSTPA_TELEMETRY_TRACING", "false")

	cfg := Load()
	if cfg.OTLPEndpoint != "http://localhost:4318" {
		t.Fatalf("OTLPEndpoint override = %q", cfg.OTLPEndpoint)
	}
	if cfg.ServiceName != "sstpa-test" {
		t.Fatalf("ServiceName override = %q", cfg.ServiceName)
	}
	if cfg.MetricsEnabled {
		t.Fatalf("MetricsEnabled override must be false")
	}
	if cfg.TracingEnabled {
		t.Fatalf("TracingEnabled override must be false")
	}
}
```

- [ ] **Step 2: Run and confirm failure**

Run: `cd backend && go test ./internal/config/...`
Expected: FAIL (compile error — fields missing).

- [ ] **Step 3: Add the new fields to `Config` and `Load`**

Replace the contents of `backend/internal/config/config.go` with:

```go
// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Address           string
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	Neo4jURI          string
	Neo4jUser         string
	Neo4jPassword     string
	Neo4jDatabase     string
	Neo4jTimeout      time.Duration
	OTLPEndpoint      string
	ServiceName       string
	MetricsEnabled    bool
	TracingEnabled    bool
}

func Load() Config {
	return Config{
		Address:           stringFromEnv("SSTPA_API_ADDR", ":8080"),
		ReadHeaderTimeout: durationFromEnv("SSTPA_API_READ_HEADER_TIMEOUT", 5*time.Second),
		WriteTimeout:      durationFromEnv("SSTPA_API_WRITE_TIMEOUT", 15*time.Second),
		Neo4jURI:          stringFromEnv("SSTPA_NEO4J_URI", ""),
		Neo4jUser:         stringFromEnv("SSTPA_NEO4J_USER", "neo4j"),
		Neo4jPassword:     stringFromEnv("SSTPA_NEO4J_PASSWORD", ""),
		Neo4jDatabase:     stringFromEnv("SSTPA_NEO4J_DATABASE", "neo4j"),
		Neo4jTimeout:      durationFromEnv("SSTPA_NEO4J_TIMEOUT", 10*time.Second),
		OTLPEndpoint:      stringFromEnv("SSTPA_OTLP_ENDPOINT", "http://otel-collector:4318"),
		ServiceName:       stringFromEnv("SSTPA_SERVICE_NAME", "sstpa-backend"),
		MetricsEnabled:    boolFromEnv("SSTPA_TELEMETRY_METRICS", true),
		TracingEnabled:    boolFromEnv("SSTPA_TELEMETRY_TRACING", true),
	}
}

func stringFromEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func durationFromEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func boolFromEnv(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
```

- [ ] **Step 4: Run config tests**

Run: `cd backend && go test ./internal/config/...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add backend/internal/config
git commit -m "$(cat <<'EOF'
feat(config): add telemetry configuration knobs

Introduces OTLPEndpoint, ServiceName, MetricsEnabled, TracingEnabled
env overrides (SSTPA_OTLP_ENDPOINT, SSTPA_SERVICE_NAME,
SSTPA_TELEMETRY_METRICS, SSTPA_TELEMETRY_TRACING). Defaults target
the docker-compose OTel Collector per SRS §2.2.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 9 — Create the `telemetry` package: Prometheus metrics

**SRS:** §2.2 (Prom instr.), §2.2.10.7 (performance indicators).

**Files:**
- Create: `backend/internal/telemetry/metrics.go`
- Create: `backend/internal/telemetry/metrics_test.go`
- Modify: `backend/go.mod` (promote / add Prometheus deps)

- [ ] **Step 1: Add Prometheus client dependency**

Run: `cd backend && go get github.com/prometheus/client_golang@latest`
Expected: `go.mod` gains a direct `github.com/prometheus/client_golang` entry.

- [ ] **Step 2: Write a failing test asserting the metrics handler emits the expected metric names**

Create `backend/internal/telemetry/metrics_test.go`:

```go
// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package telemetry

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMetricsHandlerExposesSSTPAInstruments(t *testing.T) {
	metrics := NewMetrics()
	metrics.RecordHTTPRequest("GET", "/api/v1/nodes", 200, 0.012)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/metrics", nil)
	metrics.Handler().ServeHTTP(recorder, request)

	body := recorder.Body.String()
	for _, want := range []string{
		"sstpa_http_requests_total",
		"sstpa_http_request_duration_seconds_bucket",
		"go_goroutines",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected metrics output to include %q, got:\n%s", want, body)
		}
	}
}
```

- [ ] **Step 3: Run the test and confirm failure**

Run: `cd backend && go test ./internal/telemetry/...`
Expected: FAIL — package does not exist.

- [ ] **Step 4: Implement `metrics.go`**

Create `backend/internal/telemetry/metrics.go`:

```go
// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package telemetry

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	registry   *prometheus.Registry
	requests   *prometheus.CounterVec
	durations  *prometheus.HistogramVec
	inflight   prometheus.Gauge
}

func NewMetrics() *Metrics {
	registry := prometheus.NewRegistry()
	requests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "sstpa",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total SSTPA API HTTP requests labelled by method, route, and status.",
		},
		[]string{"method", "route", "status"},
	)
	durations := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "sstpa",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "SSTPA API HTTP request duration histogram.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "route", "status"},
	)
	inflight := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "sstpa",
		Subsystem: "http",
		Name:      "inflight_requests",
		Help:      "In-flight SSTPA API HTTP requests.",
	})

	registry.MustRegister(requests, durations, inflight)
	registry.MustRegister(collectors.NewGoCollector())
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	return &Metrics{registry: registry, requests: requests, durations: durations, inflight: inflight}
}

func (m *Metrics) Registry() *prometheus.Registry { return m.registry }

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{Registry: m.registry})
}

func (m *Metrics) RecordHTTPRequest(method string, route string, status int, durationSeconds float64) {
	labels := prometheus.Labels{
		"method": method,
		"route":  route,
		"status": strconv.Itoa(status),
	}
	m.requests.With(labels).Inc()
	m.durations.With(labels).Observe(durationSeconds)
}

func (m *Metrics) InflightBegin() { m.inflight.Inc() }
func (m *Metrics) InflightEnd()   { m.inflight.Dec() }
```

- [ ] **Step 5: Run the metrics test**

Run: `cd backend && go test ./internal/telemetry/... -run TestMetricsHandlerExposesSSTPAInstruments -v`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/telemetry/metrics.go backend/internal/telemetry/metrics_test.go backend/go.mod backend/go.sum
git commit -m "$(cat <<'EOF'
feat(telemetry): add Prometheus metrics registry with SSTPA HTTP instruments

Adds backend/internal/telemetry.Metrics wrapping a dedicated
prometheus.Registry with sstpa_http_requests_total,
sstpa_http_request_duration_seconds, sstpa_http_inflight_requests,
plus Go runtime and process collectors. Implements the Prom side of
SRS §2.2 Backend diagram.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 10 — Create the `telemetry` package: OTel tracer provider

**SRS:** §2.2 (OTel instr.).

**Files:**
- Create: `backend/internal/telemetry/tracer.go`
- Create: `backend/internal/telemetry/tracer_test.go`
- Modify: `backend/go.mod` (add sdk + otlptracehttp)

- [ ] **Step 1: Add the OTel SDK and OTLP/HTTP exporter dependencies**

Run: `cd backend && go get go.opentelemetry.io/otel/sdk@v1.41.0 go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp@v1.41.0`
Expected: `go.mod` now has direct entries for these.

- [ ] **Step 2: Write a failing unit test with a tracetest recorder**

Create `backend/internal/telemetry/tracer_test.go`:

```go
// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package telemetry

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestTracerProviderRecordsSpans(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := NewTestTracerProvider(recorder)
	tracer := provider.Tracer("test")

	_, span := tracer.Start(context.Background(), "unit.span")
	span.End()

	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	if spans[0].Name() != "unit.span" {
		t.Fatalf("unexpected span name: %q", spans[0].Name())
	}
}

func TestShutdownIsIdempotent(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := NewTestTracerProvider(recorder)

	if err := provider.Shutdown(context.Background()); err != nil {
		t.Fatalf("first shutdown: %v", err)
	}
	if err := provider.Shutdown(context.Background()); err != nil {
		t.Fatalf("second shutdown: %v", err)
	}
}
```

- [ ] **Step 3: Implement `tracer.go`**

Create `backend/internal/telemetry/tracer.go`:

```go
// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package telemetry

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"sync"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

type TracerProvider struct {
	provider *sdktrace.TracerProvider
	mu       sync.Mutex
	closed   bool
}

type TracerOptions struct {
	Enabled      bool
	OTLPEndpoint string
	ServiceName  string
}

func NewTracerProvider(ctx context.Context, options TracerOptions) (*TracerProvider, error) {
	if !options.Enabled {
		provider := sdktrace.NewTracerProvider()
		return &TracerProvider{provider: provider}, nil
	}

	endpoint, path, secure, err := parseOTLPEndpoint(options.OTLPEndpoint)
	if err != nil {
		return nil, err
	}

	exporterOpts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(endpoint),
	}
	if path != "" {
		exporterOpts = append(exporterOpts, otlptracehttp.WithURLPath(path))
	}
	if !secure {
		exporterOpts = append(exporterOpts, otlptracehttp.WithInsecure())
	}

	exporter, err := otlptracehttp.New(ctx, exporterOpts...)
	if err != nil {
		return nil, err
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(buildResource(options.ServiceName)),
	)
	return &TracerProvider{provider: provider}, nil
}

func NewTestTracerProvider(recorder *tracetest.SpanRecorder) *TracerProvider {
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	return &TracerProvider{provider: provider}
}

func (p *TracerProvider) Tracer(name string) trace.Tracer {
	return p.provider.Tracer(name)
}

func (p *TracerProvider) Shutdown(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil
	}
	p.closed = true
	return p.provider.Shutdown(ctx)
}

func parseOTLPEndpoint(raw string) (string, string, bool, error) {
	if raw == "" {
		return "", "", false, errors.New("OTLP endpoint is empty")
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return "", "", false, err
	}
	if parsed.Host == "" {
		return "", "", false, errors.New("OTLP endpoint missing host")
	}

	secure := strings.EqualFold(parsed.Scheme, "https")
	path := parsed.Path
	return parsed.Host, path, secure, nil
}
```

- [ ] **Step 4: Add the resource builder**

Append to the same file:

```go
func buildResource(serviceName string) *sdkresource.Resource {
	if serviceName == "" {
		serviceName = "sstpa-backend"
	}
	return sdkresource.NewSchemaless(
		semconv.ServiceNameKey.String(serviceName),
	)
}
```

Add the following imports to `tracer.go`:

```go
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
```

Run: `cd backend && go mod tidy`

- [ ] **Step 5: Run tracer tests**

Run: `cd backend && go test ./internal/telemetry/... -run "Tracer|Shutdown" -v`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/telemetry/tracer.go backend/internal/telemetry/tracer_test.go backend/go.mod backend/go.sum
git commit -m "$(cat <<'EOF'
feat(telemetry): OTel tracer provider with OTLP/HTTP exporter

Adds TracerProvider wrapping the OTel SDK with a batch span processor
exporting to the docker-compose OTel Collector over OTLP/HTTP. Ships
with a TestTracerProvider for unit tests and an idempotent Shutdown.
Implements the OTel side of SRS §2.2.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 11 — Create the telemetry chi middleware

**SRS:** §2.2 (instrumented routes), §2.2.10.7.

**Files:**
- Create: `backend/internal/telemetry/middleware.go`
- Create: `backend/internal/telemetry/middleware_test.go`

- [ ] **Step 1: Write a failing test asserting the middleware creates a span and records a metric**

Create `backend/internal/telemetry/middleware_test.go`:

```go
// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package telemetry

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestMiddlewareRecordsSpanAndMetric(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := NewTestTracerProvider(recorder)
	metrics := NewMetrics()

	router := chi.NewRouter()
	router.Use(Middleware(MiddlewareOptions{Tracer: provider.Tracer("test"), Metrics: metrics}))
	router.Get("/api/v1/ping/{id}", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ping/42", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("unexpected status %d", res.Code)
	}

	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	if !strings.Contains(spans[0].Name(), "/api/v1/ping/{id}") {
		t.Fatalf("span name should contain the route pattern, got %q", spans[0].Name())
	}

	metricsRecorder := httptest.NewRecorder()
	metricsReq := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metrics.Handler().ServeHTTP(metricsRecorder, metricsReq)
	body := metricsRecorder.Body.String()
	if !strings.Contains(body, `route="/api/v1/ping/{id}"`) {
		t.Fatalf("expected metrics to use chi route pattern, got:\n%s", body)
	}
}
```

- [ ] **Step 2: Run and confirm failure**

Run: `cd backend && go test ./internal/telemetry/... -run TestMiddlewareRecordsSpanAndMetric`
Expected: FAIL — `Middleware` undefined.

- [ ] **Step 3: Implement the middleware**

Create `backend/internal/telemetry/middleware.go`:

```go
// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package telemetry

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type MiddlewareOptions struct {
	Tracer  trace.Tracer
	Metrics *Metrics
}

func Middleware(options MiddlewareOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			start := time.Now()
			spanWriter := &statusWriter{ResponseWriter: writer, status: http.StatusOK}

			var (
				ctx  = request.Context()
				span trace.Span
			)
			if options.Tracer != nil {
				ctx, span = options.Tracer.Start(request.Context(), request.Method+" "+request.URL.Path)
				defer span.End()
				request = request.WithContext(ctx)
			}

			if options.Metrics != nil {
				options.Metrics.InflightBegin()
			}

			next.ServeHTTP(spanWriter, request)

			route := chiRoutePattern(request)
			duration := time.Since(start).Seconds()

			if span != nil {
				span.SetAttributes(
					attribute.String("http.method", request.Method),
					attribute.String("http.route", route),
					attribute.Int("http.status_code", spanWriter.status),
					attribute.String("http.status_text", strconv.Itoa(spanWriter.status)),
				)
				if route != "" {
					span.SetName(request.Method + " " + route)
				}
			}

			if options.Metrics != nil {
				options.Metrics.InflightEnd()
				target := route
				if target == "" {
					target = request.URL.Path
				}
				options.Metrics.RecordHTTPRequest(request.Method, target, spanWriter.status, duration)
			}
		})
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (s *statusWriter) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func chiRoutePattern(request *http.Request) string {
	if ctx := chi.RouteContext(request.Context()); ctx != nil {
		if pattern := ctx.RoutePattern(); pattern != "" {
			return pattern
		}
	}
	return request.URL.Path
}
```

- [ ] **Step 4: Run the middleware test**

Run: `cd backend && go test ./internal/telemetry/... -run TestMiddlewareRecordsSpanAndMetric -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add backend/internal/telemetry/middleware.go backend/internal/telemetry/middleware_test.go
git commit -m "$(cat <<'EOF'
feat(telemetry): chi middleware combining OTel spans and Prometheus metrics

Middleware captures chi route patterns (bounded cardinality), starts a
span named "METHOD /api/v1/route/{param}", records status code and
duration, and increments the shared Prometheus registry. Implements
the instrumentation contract expected by SRS §2.2 and the docker-
compose scrape config.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 12 — Wire telemetry into the HTTP router

**SRS:** §2.2 Backend diagram, §2.2.10.1.

**Files:**
- Modify: `backend/internal/http/router.go`
- Modify: `backend/internal/http/api.go`

- [ ] **Step 1: Extend `RouterOptions` and `api` with telemetry fields**

In `backend/internal/http/router.go`, replace the `RouterOptions` struct and `NewRouterWithOptions` signature block with:

```go
type RouterOptions struct {
	Version      string
	Driver       neo4j.DriverWithContext
	DatabaseName string
	Tracer       trace.Tracer
	Metrics      *telemetry.Metrics
}
```

Add imports:

```go
	"go.opentelemetry.io/otel/trace"

	"sstpa-tool/backend/internal/telemetry"
```

- [ ] **Step 2: Plug in the middleware and `/metrics` route**

Still in `backend/internal/http/router.go`, replace the body of `NewRouterWithOptions` with:

```go
func NewRouterWithOptions(options RouterOptions) http.Handler {
	api := newAPI(options)
	router := chi.NewRouter()

	if options.Metrics != nil || options.Tracer != nil {
		router.Use(telemetry.Middleware(telemetry.MiddlewareOptions{
			Tracer:  options.Tracer,
			Metrics: options.Metrics,
		}))
	}

	router.Get("/healthz", healthHandler(api.version))
	if options.Metrics != nil {
		router.Handle("/metrics", options.Metrics.Handler())
	}

	router.Route("/api/v1", func(group chi.Router) {
		group.Get("/health", healthHandler(api.version))
		group.Get("/openapi.yaml", api.openapiHandler)

		group.Get("/nodes", api.listNodesHandler)
		group.Get("/nodes/uuid/{uuid}", api.getNodeByUUIDHandler)
		group.Get("/nodes/{hid}", api.getNodeByHIDHandler)
		group.Get("/nodes/{hid}/context", api.nodeContextHandler)
		group.Get("/hierarchy", api.hierarchyHandler)
		group.Get("/search", api.searchHandler)
		group.Post("/validate/relationship", api.validateRelationshipHandler)
		group.Post("/mutations", api.mutationsHandler)

		group.Get("/messages/unread-count", api.unreadMessageCountHandler)
		group.Get("/messages", api.listMessagesHandler)
		group.Post("/messages", api.createMessageHandler)
		group.Get("/messages/{messageId}", api.getMessageHandler)
		group.Post("/messages/{messageId}/reply", api.replyMessageHandler)
		group.Post("/messages/{messageId}/read", api.markMessageReadHandler)
		group.Delete("/messages/{messageId}", api.deleteMessageHandler)

		group.Get("/reference/frameworks", api.listReferenceFrameworksHandler)
		group.Get("/reference/items", api.listReferenceItemsHandler)
		group.Get("/reference/items/uuid/{uuid}", api.getReferenceItemByUUIDHandler)
		group.Get("/reference/items/{externalID}", api.getReferenceItemByExternalIDHandler)
		group.Get("/reference/items/{uuid}/related", api.relatedReferenceItemsHandler)
		group.Get("/reference/search", api.searchReferenceItemsHandler)
		group.Post("/reference/validate-assignment", api.validateReferenceAssignmentHandler)

		group.Get("/references/assignments/{sourceHID}", api.listReferenceAssignmentsHandler)
		group.Post("/references/assignments", api.createReferenceAssignmentHandler)
		group.Delete("/references/assignments", api.deleteReferenceAssignmentHandler)

		group.Get("/users", api.listOnboardingHandler(onboarding.UserKind))
		group.Post("/users", api.createOnboardingHandler(onboarding.UserKind))
		group.Get("/users/{uuid}", api.getOnboardingHandler(onboarding.UserKind))
		group.Get("/admins", api.listOnboardingHandler(onboarding.AdminKind))
		group.Post("/admins", api.createOnboardingHandler(onboarding.AdminKind))
		group.Get("/admins/{uuid}", api.getOnboardingHandler(onboarding.AdminKind))
	})

	return router
}
```

- [ ] **Step 3: Add tracer/metrics to `api` struct**

In `backend/internal/http/api.go`, replace the `api` struct and `newAPI` with:

```go
type api struct {
	version      string
	driver       neo4j.DriverWithContext
	databaseName string
	now          func() time.Time
	tracer       trace.Tracer
	metrics      *telemetry.Metrics
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
		tracer:       options.Tracer,
		metrics:      options.Metrics,
	}
}
```

Add the imports to `api.go`:

```go
	"go.opentelemetry.io/otel/trace"

	"sstpa-tool/backend/internal/telemetry"
```

- [ ] **Step 4: Write a test that the router serves `/metrics` when Metrics are provided**

Append to `backend/internal/http/api_test.go`:

```go
func TestRouterServesMetricsEndpoint(t *testing.T) {
	metrics := telemetry.NewMetrics()
	router := NewRouterWithOptions(RouterOptions{Version: "test", Metrics: metrics})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), "sstpa_http_inflight_requests") {
		t.Fatalf("expected sstpa metrics in body, got:\n%s", recorder.Body.String())
	}
}
```

Add the import `"sstpa-tool/backend/internal/telemetry"` to `api_test.go`.

- [ ] **Step 5: Run the test**

Run: `cd backend && go test ./internal/http/... -run TestRouterServesMetricsEndpoint -v`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/http/router.go backend/internal/http/api.go backend/internal/http/api_test.go
git commit -m "$(cat <<'EOF'
feat(http): install telemetry middleware and /metrics endpoint

Router now accepts trace.Tracer and *telemetry.Metrics options,
installs the telemetry middleware, and serves /metrics when a
registry is provided. Implements the instrumented surface SRS §2.2
assumes for the docker-compose Prom scrape and OTel export.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 13 — Wrap Neo4j mutations with an explicit span

**SRS:** §2.2.10.8 Transaction Requirements, §2.2 OTel instr.

**Files:**
- Modify: `backend/internal/mutation/apply.go`

- [ ] **Step 1: Write a failing test that records a mutation span**

Append to `backend/internal/mutation/apply_test.go` (the test file exists alongside `apply.go`):

```go
func TestApplyRecordsTraceSpan(t *testing.T) {
	fixture := testhelpers.StartNeo4j(t)
	ctx := context.Background()
	if err := schema.Bootstrap(ctx, fixture.Driver, ""); err != nil {
		t.Fatal(err)
	}

	recorder := tracetest.NewSpanRecorder()
	provider := telemetry.NewTestTracerProvider(recorder)
	SetTracer(provider.Tracer("mutation-test"))
	t.Cleanup(func() { SetTracer(nil) })

	actor := metadata.Actor{Name: "Alice", Email: "alice@example.test"}
	plan := Plan{Operations: []Operation{{
		Kind:     OperationCreateNode,
		NodeType: identity.NodeTypeCapability,
		HID:      "CAP__1",
		UUID:     "00000000-0000-4000-8000-000000000900",
		Properties: map[string]any{"Name": "Root"},
	}}}

	if _, err := Apply(ctx, fixture.Driver, ApplyOptions{Actor: actor, VersionID: "v1"}, plan); err != nil {
		t.Fatalf("apply: %v", err)
	}

	spans := recorder.Ended()
	if len(spans) == 0 {
		t.Fatalf("expected a mutation span, got none")
	}
	found := false
	for _, span := range spans {
		if span.Name() == "sstpa.mutation.apply" {
			found = true
			attrs := span.Attributes()
			hasCommit := false
			for _, a := range attrs {
				if string(a.Key) == "sstpa.commit_id" && a.Value.AsString() != "" {
					hasCommit = true
				}
			}
			if !hasCommit {
				t.Fatalf("span missing sstpa.commit_id attribute")
			}
		}
	}
	if !found {
		t.Fatalf("sstpa.mutation.apply span not found; names: %v", spanNames(spans))
	}
}

func spanNames(spans []sdktrace.ReadOnlySpan) []string {
	names := make([]string, 0, len(spans))
	for _, span := range spans {
		names = append(names, span.Name())
	}
	return names
}
```

Add the imports at the top of `apply_test.go`:

```go
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"sstpa-tool/backend/internal/telemetry"
```

- [ ] **Step 2: Run and confirm failure**

Run: `cd backend && go test ./internal/mutation/... -run TestApplyRecordsTraceSpan`
Expected: FAIL — `SetTracer` undefined.

- [ ] **Step 3: Add a package-level tracer hook to the mutation package**

Append to the top of `backend/internal/mutation/apply.go` (after imports):

```go
var mutationTracer trace.Tracer

func SetTracer(tracer trace.Tracer) { mutationTracer = tracer }
```

Add `"go.opentelemetry.io/otel/trace"` and `"go.opentelemetry.io/otel/attribute"` to the imports.

- [ ] **Step 4: Wrap the write transaction in a span**

Inside `Apply`, replace the `session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) { ... })` call with a wrapped version:

```go
	ctxForWrite := ctx
	var span trace.Span
	if mutationTracer != nil {
		ctxForWrite, span = mutationTracer.Start(ctx, "sstpa.mutation.apply")
		defer span.End()
	}

	result, err := session.ExecuteWrite(ctxForWrite, func(tx neo4j.ManagedTransaction) (any, error) {
		before, err := readSnapshot(ctx, tx, preExistingHIDs(plan))
		if err != nil {
			return CommitReport{}, err
		}

		for _, operation := range plan.Operations {
			if err := applyOperation(ctx, tx, options, now, operation); err != nil {
				return CommitReport{}, err
			}
		}

		allHIDs := allOperationHIDs(plan)
		after, err := readSnapshot(ctx, tx, allHIDs)
		if err != nil {
			return CommitReport{}, err
		}

		affected := ComputeAffected(plan, before, after)
		report := CommitReport{
			CommitID:             commitID,
			NodesChanged:         affectedHIDs(affected),
			RelationshipsChanged: relationshipChanges(plan),
		}

		recipients, err := notifyAffectedOwners(ctx, tx, options.Actor, commitID, now, affected, before, after, report.RelationshipsChanged)
		if err != nil {
			return CommitReport{}, err
		}

		report.RecipientsNotified = recipients
		report.MessagesGenerated = len(recipients)
		return report, nil
	})

	if span != nil {
		span.SetAttributes(
			attribute.String("sstpa.commit_id", commitID),
			attribute.Int("sstpa.operations_count", len(plan.Operations)),
			attribute.String("sstpa.actor_email", options.Actor.Email),
		)
		if err != nil {
			span.RecordError(err)
		} else if report, ok := result.(CommitReport); ok {
			span.SetAttributes(
				attribute.Int("sstpa.messages_generated", report.MessagesGenerated),
				attribute.Int("sstpa.nodes_changed", len(report.NodesChanged)),
			)
		}
	}

	if err != nil {
		return CommitReport{}, err
	}
```

(Remove the preceding `result, err := session.ExecuteWrite(...)` and the standalone `if err != nil { return CommitReport{}, err }` that existed before this block — the rewrite above subsumes them.)

- [ ] **Step 5: Run the mutation tests**

Run: `cd backend && go test ./internal/mutation/... -v`
Expected: PASS (including the new `TestApplyRecordsTraceSpan`).

- [ ] **Step 6: Commit**

```bash
git add backend/internal/mutation/apply.go backend/internal/mutation/apply_test.go
git commit -m "$(cat <<'EOF'
feat(mutation): wrap Apply in an OTel span with commit attributes

Adds a package-level mutationTracer + SetTracer hook and wraps the
ExecuteWrite call in a sstpa.mutation.apply span carrying commit_id,
operations_count, actor_email, messages_generated, and nodes_changed.
Implements the Neo4j-tracing scope of the S04b telemetry slice.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 14 — Initialise telemetry in `cmd/api/main.go`

**SRS:** §2.2 Backend diagram (service boundary).

**Files:**
- Modify: `backend/cmd/api/main.go`

- [ ] **Step 1: Replace `main.go` to build telemetry and pass it to the router**

Replace the contents of `backend/cmd/api/main.go` with:

```go
// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"sstpa-tool/backend/internal/config"
	apihttp "sstpa-tool/backend/internal/http"
	"sstpa-tool/backend/internal/mutation"
	"sstpa-tool/backend/internal/neo4jx"
	"sstpa-tool/backend/internal/schema"
	"sstpa-tool/backend/internal/telemetry"
	"sstpa-tool/backend/internal/version"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	driver, err := neo4jx.Open(ctx, neo4jx.Config{
		URI:      cfg.Neo4jURI,
		User:     cfg.Neo4jUser,
		Password: cfg.Neo4jPassword,
		Database: cfg.Neo4jDatabase,
		Timeout:  cfg.Neo4jTimeout,
	})
	if err != nil {
		slog.Error("neo4j connection failed", "error", err)
		os.Exit(1)
	}
	if driver != nil {
		defer driver.Close(ctx)
		if err := schema.Bootstrap(ctx, driver, cfg.Neo4jDatabase); err != nil {
			slog.Error("neo4j schema bootstrap failed", "error", err)
			os.Exit(1)
		}
		slog.Info("neo4j schema ready", "database", cfg.Neo4jDatabase)
	} else {
		slog.Info("neo4j disabled; set SSTPA_NEO4J_URI to enable graph persistence")
	}

	var metrics *telemetry.Metrics
	if cfg.MetricsEnabled {
		metrics = telemetry.NewMetrics()
	}

	tracerProvider, err := telemetry.NewTracerProvider(ctx, telemetry.TracerOptions{
		Enabled:      cfg.TracingEnabled,
		OTLPEndpoint: cfg.OTLPEndpoint,
		ServiceName:  cfg.ServiceName,
	})
	if err != nil {
		slog.Error("telemetry bootstrap failed", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			slog.Error("telemetry shutdown failed", "error", err)
		}
	}()

	tracer := tracerProvider.Tracer("sstpa.backend")
	mutation.SetTracer(tracer)

	server := &http.Server{
		Addr: cfg.Address,
		Handler: apihttp.NewRouterWithOptions(apihttp.RouterOptions{
			Version:      version.Dev,
			Driver:       driver,
			DatabaseName: cfg.Neo4jDatabase,
			Tracer:       tracer,
			Metrics:      metrics,
		}),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
	}

	slog.Info("starting sstpa api", "addr", server.Addr)

	err = server.ListenAndServe()
	if err == nil || errors.Is(err, http.ErrServerClosed) {
		return
	}

	slog.Error("sstpa api exited", "error", err)
	os.Exit(1)
}
```

- [ ] **Step 2: Build to confirm the wiring compiles**

Run: `cd backend && go build ./...`
Expected: builds cleanly.

- [ ] **Step 3: Commit**

```bash
git add backend/cmd/api/main.go
git commit -m "$(cat <<'EOF'
feat(api): bootstrap telemetry provider and wire mutation tracer

cmd/api now builds the Prometheus registry (when enabled) and the OTel
tracer provider, registers the mutation tracer hook, passes both to the
router, and defers a graceful tracer shutdown. Satisfies the service
boundary expected by SRS §2.2.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 15 — Update the verification documents

**SRS:** CLAUDE.md verification workflow.

**Files:**
- Modify: `docs/verification/shall-register.md`
- Modify: `docs/verification/verification-matrix.md`

- [ ] **Step 1: Append five rows to the SHALL register**

In `docs/verification/shall-register.md`, inside the table body (before the `Status guidance:` line), append these rows:

```
| SSTPA-USER-MODEL-001 | 1.4.2 / 1.4.3 | Users and Admins are bootstrapped as :Users / :Admins singleton containers with :User:SSTPANode and :Admin:SSTPANode children carrying identity metadata | Approved | unit + integration | Implemented via `backend/internal/onboarding` and schema bootstrap |
| SSTPA-API-USER-001 | 2.1 / 2.2.10 | Backend exposes `/api/v1/users` list/create/get endpoints returning SSTPA identity metadata | Approved | integration | Covered by `backend/internal/http/users_test.go` |
| SSTPA-API-ADMIN-001 | 1.4.2 / 2.1 | Backend exposes `/api/v1/admins` list/create/get endpoints returning SSTPA identity metadata | Approved | integration | Covered by `backend/internal/http/users_test.go` |
| SSTPA-TELEMETRY-HTTP-001 | 2.2 / 2.2.10.7 | Backend serves `/metrics` in Prometheus exposition format and emits OTel HTTP spans using chi route patterns | Approved | unit + integration | Implemented in `backend/internal/telemetry`; `/metrics` served from the chi router |
| SSTPA-TELEMETRY-NEO4J-001 | 2.2 / 2.2.10.8 | Mutation transactions emit an OTel span `sstpa.mutation.apply` with commit_id, operations_count, actor_email, messages_generated attributes | Approved | integration | Implemented in `backend/internal/mutation` via `SetTracer` |
```

- [ ] **Step 2: Append five rows to the verification matrix**

In `docs/verification/verification-matrix.md`, inside the table body, append:

```
| SSTPA-USER-MODEL-001 | Users/Admins container + typed nodes | unit + integration | Yes | `make backend-test` | `backend/internal/schema`, `backend/internal/onboarding` | Bootstrap creates the singletons; onboarding package enforces HID allocation and identity metadata |
| SSTPA-API-USER-001 | `/api/v1/users` REST endpoints | integration | Yes | `make backend-test` | `backend/internal/http` | Create, duplicate-rejection, list, get-by-uuid covered |
| SSTPA-API-ADMIN-001 | `/api/v1/admins` REST endpoints | integration | Yes | `make backend-test` | `backend/internal/http` | Same shape as SSTPA-API-USER-001 |
| SSTPA-TELEMETRY-HTTP-001 | `/metrics` exposition + HTTP spans | unit + integration | Yes | `make backend-test` | `backend/internal/telemetry`, `backend/internal/http` | Route-pattern cardinality enforced by chi middleware |
| SSTPA-TELEMETRY-NEO4J-001 | Mutation-layer trace span | integration | Yes | `make backend-test` | `backend/internal/mutation` | Test asserts `sstpa.mutation.apply` span with required attributes |
```

- [ ] **Step 3: Commit**

```bash
git add docs/verification/shall-register.md docs/verification/verification-matrix.md
git commit -m "$(cat <<'EOF'
docs(verification): approve S04b SHALLs for onboarding + telemetry

Adds SSTPA-USER-MODEL-001, SSTPA-API-USER-001, SSTPA-API-ADMIN-001,
SSTPA-TELEMETRY-HTTP-001, SSTPA-TELEMETRY-NEO4J-001 with matching
verification matrix rows. Closes the SRS §2.1 and §2.2 Prom/OTel
loose ends from S04.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 16 — Bootstrap dependencies and run `make verify`

**SRS:** CLAUDE.md Verification Workflow.

**Files:** none (build + verification only).

- [ ] **Step 1: Install npm workspaces**

Run from repo root: `make bootstrap`
Expected: `npm install` completes. If it prints vulnerability warnings, that is informational — not a failure.

- [ ] **Step 2: Populate Go modules**

Run: `cd backend && go mod download && cd ../tools/reference-pipeline && go mod download && cd ../devtools/copyright && go mod download`
Expected: no errors.

- [ ] **Step 3: Run the copyright check**

Run from repo root: `make copyright-check`
Expected: PASS. If it fails on new files, the banner at the top of each new `.go` / `.ts` file is missing or out of format.

- [ ] **Step 4: Run the backend tests**

Run from repo root: `make backend-test`
Expected: PASS (Testcontainers tests skip cleanly when Docker is unavailable; they must pass when Docker is running locally).

- [ ] **Step 5: Run the frontend tests**

Run from repo root: `make frontend-test`
Expected: PASS (Vitest covers the api-client additions).

- [ ] **Step 6: Run the full `make verify`**

Run from repo root: `make verify`
Expected: PASS. The gate runs copyright-check → sbom-check → backend-test → reference-test → frontend-lint → frontend-typecheck → frontend-test → compose-config in that order.

- [ ] **Step 7: If `sbom-check` flags new Go packages, regenerate the SBOM**

If `make sbom-check` fails because it detected new Go dependencies (Prometheus client, OTel sdk, otlptracehttp), run `make sbom-generate` to refresh `docs/compliance/sbom.md`, inspect the diff to confirm the new entries are reasonable, then commit:

```bash
git add docs/compliance/sbom.md
git commit -m "$(cat <<'EOF'
chore(sbom): refresh SBOM for telemetry dependencies

Adds prometheus/client_golang, go.opentelemetry.io/otel/sdk, and
otlptracehttp entries introduced by the S04b telemetry slice. Keeps
SSTPA-SBOM-001 green.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

- [ ] **Step 8: Capture evidence of a green verify run**

Append the following line to the end of this plan file (still inside the repo), listing the commit SHA at `git log -1 --format=%H`:

```
**Verify Evidence:** `make verify` green at commit <SHA> on `slice/04b-users-admins-telemetry`.
```

Then commit:

```bash
git add docs/superpowers/plans/2026-04-24-s04b-users-admins-telemetry.md
git commit -m "$(cat <<'EOF'
docs(plans): record S04b make-verify evidence

Captures the green `make verify` run that closes the Users + Admins
onboarding and telemetry loose ends from S04.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Acceptance criteria

- `GET /api/v1/users`, `POST /api/v1/users`, `GET /api/v1/users/{uuid}`, `GET /api/v1/admins`, `POST /api/v1/admins`, `GET /api/v1/admins/{uuid}` return the shape above and are covered by Testcontainers integration tests.
- `MERGE (:Users)` and `MERGE (:Admins)` singletons are present after `schema.Bootstrap`; reregistration does not duplicate them.
- `POST /api/v1/users` with a duplicate email returns HTTP 409.
- `GET /metrics` returns Prometheus exposition format including `sstpa_http_requests_total` and `sstpa_http_request_duration_seconds_bucket`, and metric labels use the chi route pattern (not the raw URL path).
- A mutation against the backend emits a span named `sstpa.mutation.apply` carrying `sstpa.commit_id`, `sstpa.operations_count`, `sstpa.actor_email` attributes; integration test enforces this.
- `make verify` runs green end-to-end.
- `docs/verification/shall-register.md` and `docs/verification/verification-matrix.md` carry the five new `Approved` rows.

## Out-of-scope / deferred

- Password-based authentication or user sessions (§1.4.2 explicit defer).
- Role-based policy enforcement beyond the existing mutation-layer Admin flag (§2.2.10.9 placeholder).
- Grafana dashboards and alerting rules — owned by S12.
- Bolt driver instrumentation beyond the mutation-layer span — future work when Grafana dashboards need it.
- Replacing the MERGE-on-email behavior inside `backend/internal/messaging/model.go`. That path remains a messaging-only upsert; registered Users keep the `:SSTPANode` label and container edge so GET /users filters on that.

---

## Verify Evidence

`make verify` green at commit `c15618e29fad5380144fefe67f2a3df180fce330` on branch `slice/04b-users-admins-telemetry` — all targets (copyright-check, sbom-check, backend-test, reference-test, frontend-lint, frontend-typecheck, frontend-test, compose-config) exited successfully. Testcontainers Neo4j fixtures exercised via Docker 29.2.1; 12 backend packages, 3 reference-pipeline packages, and 2 TypeScript workspaces all green.
