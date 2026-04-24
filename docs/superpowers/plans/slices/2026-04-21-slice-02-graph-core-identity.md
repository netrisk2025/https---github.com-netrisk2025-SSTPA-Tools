# Slice 02 — Graph Core: Schema, Identity, Ownership (Phase A, Detailed)

> **For agentic workers:** REQUIRED SUB-SKILL: `superpowers:subagent-driven-development` or `superpowers:executing-plans`. Read `docs/superpowers/plans/2026-04-21-agent-execution-conventions.md` first. One task = one commit.

**Goal:** Stand up the Neo4j-backed core: driver wiring, node-type catalog, HID + uuid generation, common property group (identity + ownership + description), and the constraint/index bootstrap. No REST endpoints beyond `/health`, no mutation layer, no notifications. Those land in S03.

**Architecture:** Pure-Go library in `backend/internal/{identity,graph,metadata,schema}`. An opt-in Neo4j driver wired into `cmd/api` at startup that runs constraint/index migrations and then idles. All new code is library-first so S03 can assemble it without touching the HTTP layer.

**Tech Stack:** Go 1.24, `github.com/neo4j/neo4j-go-driver/v5`, `github.com/google/uuid`. Cypher 25.

**Pre-reads (mandatory before any task):**
- `SSTPA Tool SRS V56.md` §1.3 (Core Data Model), §1.3.6 (Identity), §1.3.7 (Common Property Groups), §2.2.4–2.2.5 (Backend Software, Backend Database)
- `backend/cmd/api/main.go`
- `backend/internal/config/config.go`
- `backend/internal/http/router.go`
- `backend/internal/testhelpers/neo4j.go` (from S01-T02)
- `CLAUDE.md`
- `docs/superpowers/plans/2026-04-21-agent-execution-conventions.md`

**Invariants:**
- Node labels are singular (`(:System)`, not `(:Systems)`).
- Relationship names are `UPPERCASE_SNAKE_CASE`.
- No reverse relationships unless SRS explicitly requires one.
- Property values that would otherwise be empty are the literal string `"Null"`.
- HID format is `{TYPE}_{INDEX}_{SEQUENCE}` exactly.
- `Capability` is the only type whose `Index` is empty.
- `System` nodes always carry sequence `0`.
- Every node created through this slice carries the full identity + ownership metadata.
- `make verify` stays green after every task.

---

### Task S02-T01: Node type catalog

**Task ID:** S02-T01
**Depends On:** none (S01 must be Accepted)
**Difficulty:** low
**Integration Checkpoint:** no
**Governing SRS sections:** §1.3.1 Core Data Model Nodes, §1.3.6.1 Node Type Identifier
**Req IDs touched:** `SSTPA-GRAPH-RULES-001`, proposes `SSTPA-NODETYPE-CATALOG-001`

**Files:**
- Create: `backend/internal/identity/types.go`
- Create: `backend/internal/identity/types_test.go`

**Agent Briefing:**
Build a Go source of truth for the 25 SRS Core Data Model node types and their three-letter identifiers. Export:
- `type NodeType string` with typed constants `NodeTypeCapability`, `NodeTypeSandbox`, ... for all 25 types (see SRS §1.3.6.1 for the canonical list and abbreviations).
- `TypeID(nt NodeType) (string, bool)` returning the canonical ID (e.g., `"SYS"`) and `ok=false` for unknown types.
- `AllTypes() []NodeType` returning the 25 types in SRS listing order.
- `IsValidTypeID(id string) bool`.

The list from SRS §1.3.6.1 is authoritative: Capability CAP, Sandbox SB, System SYS, Environment ENV, Connection CNN, Interface INT, Function FUN, Element EL, Purpose PUR, State ST, ControlStructure CS, Asset AST, Security SEC, Constraint CONSTR, Requirement REQ, Validation VAL, Control CTRL, Countermeasure CM, Verification VER, ControlAlgorithm CAL, ProcessModel PM, ControlAction ACT, Feedback FB, ControlledProcess CP, Hazard HAZ, Loss LOS, Attack ATK. Count: 27 IDs → 27 types (cross-check with SRS before shipping).

**Step 1: Failing test**

Create `backend/internal/identity/types_test.go`:
```go
package identity

import "testing"

func TestTypeIDKnown(t *testing.T) {
    cases := map[NodeType]string{
        NodeTypeCapability: "CAP",
        NodeTypeSystem:     "SYS",
        NodeTypeHazard:     "HAZ",
        NodeTypeLoss:       "LOS",
        NodeTypeAttack:     "ATK",
    }
    for nt, want := range cases {
        got, ok := TypeID(nt)
        if !ok || got != want {
            t.Errorf("TypeID(%q) = (%q,%v), want (%q,true)", nt, got, ok, want)
        }
    }
}

func TestTypeIDUnknown(t *testing.T) {
    if _, ok := TypeID(NodeType("Nope")); ok {
        t.Fatal("expected ok=false for unknown type")
    }
}

func TestAllTypesCountMatchesSRS(t *testing.T) {
    if got := len(AllTypes()); got < 25 {
        t.Fatalf("expected >=25 node types per SRS §1.3.6.1, got %d", got)
    }
}

func TestIsValidTypeID(t *testing.T) {
    if !IsValidTypeID("SYS") {
        t.Fatal("SYS must be valid")
    }
    if IsValidTypeID("NOPE") {
        t.Fatal("NOPE must not be valid")
    }
}
```

**Step 2: Run — expect fail**
```
cd backend && go test ./internal/identity/...
```
Expected: undefined symbols.

**Step 3: Implement**

Create `backend/internal/identity/types.go` with all 27 constants, the `TypeID` map, `AllTypes()`, and `IsValidTypeID()`. Use `NodeType` as a string alias. Declare IDs in the order of SRS §1.3.6.1.

**Step 4: Run — expect pass**
```
cd backend && go test ./internal/identity/... -v
```
Expected: all 4 tests pass.

**Step 5: Commit**
```
git add backend/internal/identity
git commit -m "$(cat <<'EOF'
feat(identity): introduce node-type catalog per SRS §1.3.6.1

Adds the 27-entry NodeType enum and TypeID lookup shared by the HID
generator and future schema bootstrapper. Pure library; no DB.

Refs: SRS §1.3.1, §1.3.6.1. Req: SSTPA-GRAPH-RULES-001.
Proposes: SSTPA-NODETYPE-CATALOG-001. Task: S02-T01.
Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

**Acceptance tests:**
- `cd backend && go test ./internal/identity/...` is green.
- `go vet ./...` clean.

**Out-of-scope:** No DB writes. No label strings beyond the catalog (label mapping = next task).

**Evidence:** _________

---

### Task S02-T02: Label + relationship catalog

**Task ID:** S02-T02
**Depends On:** S02-T01
**Difficulty:** medium
**Integration Checkpoint:** no
**Governing SRS sections:** §1.3.2 General Modeling Rules, §1.3.3–1.3.4 Relationships
**Req IDs touched:** `SSTPA-GRAPH-RULES-001`, proposes `SSTPA-RELATIONSHIP-CATALOG-001`

**Files:**
- Create: `backend/internal/graph/labels.go`
- Create: `backend/internal/graph/labels_test.go`
- Create: `backend/internal/graph/relationships.go`
- Create: `backend/internal/graph/relationships_test.go`

**Agent Briefing:**
Two files of pure catalog: (a) a function `LabelFor(NodeType) string` that returns the singular Neo4j label for each NodeType (exactly matching SRS §1.3.1 spellings — e.g., `"Capability"`, `"System"`, `"ControlStructure"`). (b) a `Relationship` record with `Name string`, `From NodeType`, `To NodeType`, and `DAG bool` (see SRS §1.3.2.1) plus a `Catalog() []Relationship` seeded with every primary and secondary relationship listed in SRS §1.3.3 and §1.3.4. Names must be `UPPERCASE_SNAKE_CASE`.

The tests MUST (1) assert every NodeType maps to a label, (2) assert every catalog entry's Name matches `^[A-Z][A-Z_]+$`, (3) assert at least 20 relationships are catalogued (cross-check SRS §1.3.3/§1.3.4 count), (4) include a spot check: `HAS_SYSTEM`, `HAS_CONNECTION`, `HAS_REQUIREMENT`, `TRANSITIONS_TO`, `HAS_LOSS`, `REFERENCES`, `MITIGATES`, `SATISFIES`, `PARENTS`, `HAS_HAZARD` are all present with correct endpoints.

**Step 1–4: TDD loop** — test, fail, implement, pass.

**Step 5: Commit**
```
git commit -m "$(cat <<'EOF'
feat(graph): add label + relationship catalog

Seeds the authoritative mapping NodeType→Neo4j label and the primary +
secondary relationship catalog per SRS §1.3.3/§1.3.4 with naming rules
enforced by regex tests.

Refs: SRS §1.3.2, §1.3.3, §1.3.4. Req: SSTPA-GRAPH-RULES-001.
Proposes: SSTPA-RELATIONSHIP-CATALOG-001. Task: S02-T02.
Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

**Acceptance tests:**
- `cd backend && go test ./internal/graph/... -v` is green.
- No relationship name fails the regex `^[A-Z][A-Z_]+$`.

**Out-of-scope:** Constraint definitions beyond naming. DAG enforcement is S03.

**Evidence:** _________

---

### Task S02-T03: HID formatter + sequence/index validators

**Task ID:** S02-T03
**Depends On:** S02-T01
**Difficulty:** low
**Integration Checkpoint:** no
**Governing SRS sections:** §1.3.6, §1.3.6.2, §1.3.6.3
**Req IDs touched:** `SSTPA-HID-001`

**Files:**
- Create: `backend/internal/identity/hid.go`
- Create: `backend/internal/identity/hid_test.go`

**Agent Briefing:**
Pure formatting. Export `FormatHID(typeID, index string, sequence int) (string, error)` and `ParseHID(hid string) (typeID, index string, sequence int, err error)`. Rules from SRS §1.3.6:
- Shape `{TYPE}_{INDEX}_{SEQUENCE}`. Capability has empty INDEX so its HID looks like `CAP__0`.
- `System` nodes always carry sequence `0` (enforcement belongs to the mutation layer; the formatter permits any sequence ≥ 0).
- Sequence is a non-negative integer.
- INDEX is a dotted non-negative-integer string or empty: `^(\d+(\.\d+)*)?$`.
- Unknown TYPE IDs rejected (cross-checked against `IsValidTypeID` from S02-T01).

**Step 1: Failing table-driven test**

```go
package identity

import "testing"

func TestFormatHID(t *testing.T) {
    tests := []struct {
        name     string
        typeID   string
        index    string
        sequence int
        want     string
        wantErr  bool
    }{
        {"capability", "CAP", "", 0, "CAP__0", false},
        {"root system", "SYS", "1", 0, "SYS_1_0", false},
        {"nested element", "EL", "1.2.3", 4, "EL_1.2.3_4", false},
        {"negative sequence", "SYS", "1", -1, "", true},
        {"unknown type", "XYZ", "1", 0, "", true},
        {"malformed index", "SYS", "1..2", 0, "", true},
    }
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got, err := FormatHID(tc.typeID, tc.index, tc.sequence)
            if (err != nil) != tc.wantErr {
                t.Fatalf("err=%v wantErr=%v", err, tc.wantErr)
            }
            if got != tc.want {
                t.Fatalf("got %q want %q", got, tc.want)
            }
        })
    }
}

func TestParseHIDRoundTrip(t *testing.T) {
    in := "EL_1.2.3_4"
    typeID, index, seq, err := ParseHID(in)
    if err != nil {
        t.Fatal(err)
    }
    got, err := FormatHID(typeID, index, seq)
    if err != nil || got != in {
        t.Fatalf("round-trip failed: got %q err %v", got, err)
    }
}

func TestParseHIDRejectsMalformed(t *testing.T) {
    for _, bad := range []string{"", "NOPE", "SYS_1_", "SYS__", "SYS_1_abc"} {
        if _, _, _, err := ParseHID(bad); err == nil {
            t.Errorf("expected error for %q", bad)
        }
    }
}
```

**Step 2: Run — expect fail**
**Step 3: Implement**

Key implementation details:
- Use `regexp.MustCompile(`^(\d+(\.\d+)*)?$`)` for index.
- Sequence parsed with `strconv.Atoi` and checked `>= 0`.
- `ParseHID` splits on underscore: `strings.SplitN(hid, "_", 3)` expecting exactly 3 parts where part 1 is type, part 2 is index (possibly empty), part 3 is sequence.
- Capability round-trip: empty index accepted; `CAP__0` tokenises to `["CAP", "", "0"]`.

**Step 4: Pass** — green.

**Step 5: Commit**
```
git commit -m "$(cat <<'EOF'
feat(identity): HID formatter and parser per SRS §1.3.6

Pure functions for {TYPE}_{INDEX}_{SEQUENCE} with validation. Capability
shape (empty INDEX) round-trips correctly. Sequencer and index
generator live in separate tasks once the mutation layer exists.

Refs: SRS §1.3.6, §1.3.6.2, §1.3.6.3. Req: SSTPA-HID-001.
Task: S02-T03.
Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

**Acceptance tests:**
- `cd backend && go test ./internal/identity/...` is green.

**Out-of-scope:** Sequence allocation against a live graph. That is a backend concern introduced in S03.

**Evidence:** _________

---

### Task S02-T04: UUID wrapper

**Task ID:** S02-T04
**Depends On:** S02-T01
**Difficulty:** low
**Integration Checkpoint:** no
**Governing SRS sections:** §1.3.6 Identity Model
**Req IDs touched:** proposes `SSTPA-UUID-001`

**Files:**
- Create: `backend/internal/identity/uuid.go`
- Create: `backend/internal/identity/uuid_test.go`
- Modify: `backend/go.mod` (add `github.com/google/uuid`)

**Agent Briefing:**
Thin wrapper so we can swap implementations later. Export `NewUUID() string` returning a UUIDv4 canonical string. SRS says `uuid: apoc.create.uuid()` but APOC isn't available in Community Edition by default; generate UUIDs in Go. Tests: format matches `^[0-9a-f-]{36}$`, two successive calls differ.

TDD loop then commit. This task introduces the `github.com/google/uuid` dependency — flag it in the commit body as a library addition per `CLAUDE.md`.

**Step 5: Commit**
```
git commit -m "$(cat <<'EOF'
feat(identity): UUID wrapper per SRS §1.3.6

Introduces github.com/google/uuid as the canonical uuid source, wrapped
to keep callers swappable. APOC is not assumed in Community Edition.

Library addition approved per S01 standing approval for identity/graph
primitives.

Refs: SRS §1.3.6. Proposes: SSTPA-UUID-001. Task: S02-T04.
Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

**Evidence:** _________

---

### Task S02-T05: Common metadata struct (identity + ownership + description)

**Task ID:** S02-T05
**Depends On:** S02-T01, S02-T03, S02-T04
**Difficulty:** medium
**Integration Checkpoint:** no
**Governing SRS sections:** §1.3.7, §1.3.7.1
**Req IDs touched:** `SSTPA-IDENTITY-001`

**Files:**
- Create: `backend/internal/metadata/common.go`
- Create: `backend/internal/metadata/common_test.go`

**Agent Briefing:**
Encapsulate the common property group from SRS §1.3.7 as a Go value:
```go
type Common struct {
    Name             string
    HID              string
    UUID             string
    TypeName         string
    Owner            string
    OwnerEmail       string
    Creator          string
    CreatorEmail     string
    Created          time.Time
    LastTouch        time.Time
    VersionID        string
    ShortDescription string
    LongDescription  string
}
```
Provide:
- `NewCommon(params NewCommonParams) Common` — fills HID via `identity.FormatHID`, UUID via `identity.NewUUID`, sets Creator/Owner/CreatorEmail/OwnerEmail to the current user, Created=LastTouch=now, VersionID from a compile-time constant `SchemaVersion` declared in the same package (initial value `"0.1.0"`), ShortDescription/LongDescription default `"Null"`.
- `TouchOnWrite(c *Common, now time.Time)` — updates LastTouch; never touches Creator.
- `TransferOwnership(c *Common, newOwner, newEmail string, now time.Time)` — changes Owner/OwnerEmail/LastTouch. Leaves Creator/CreatorEmail unchanged. Per SRS §1.3.7.1, this is ALWAYS considered a node modification for notification purposes; S03 is responsible for emitting the notification — this function only updates the struct.
- `ToCypherProperties(c Common) map[string]any` — returns the flat property map suitable for a `CREATE/MERGE` `SET n = $props` Cypher, substituting empty strings with `"Null"`.

Tests cover each path including that `"Null"` substitution, that TouchOnWrite preserves Creator, and that TransferOwnership does not modify Creator.

**Step 5: Commit**
```
git commit -m "$(cat <<'EOF'
feat(metadata): common property group struct per SRS §1.3.7

Introduces metadata.Common covering identity + ownership + description
with factory, touch, and ownership-transfer helpers. Empty strings
serialise to \"Null\" per SRS §1.3.2. Notification emission belongs to
S03; this task stays library-only.

Refs: SRS §1.3.7, §1.3.7.1. Req: SSTPA-IDENTITY-001. Task: S02-T05.
Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

**Evidence:** _________

---

### Task S02-T06: Schema bootstrapper — constraints + indexes

**Task ID:** S02-T06
**Depends On:** S02-T02, S02-T05
**Difficulty:** medium
**Integration Checkpoint:** no
**Governing SRS sections:** §1.3.6.2.1 Index Strategy, §2.2.5 Backend Database
**Req IDs touched:** proposes `SSTPA-SCHEMA-BOOTSTRAP-001`

**Files:**
- Create: `backend/internal/schema/bootstrap.go`
- Create: `backend/internal/schema/bootstrap_test.go`
- Create: `backend/internal/schema/statements.go` (list of Cypher statements as `[]string`)

**Agent Briefing:**
Provide `schema.Bootstrap(ctx, driver)` that executes, in dependency order and idempotently, all constraints and indexes listed in SRS §1.3.6.2.1:
```cypher
CREATE INDEX node_hid_index IF NOT EXISTS FOR (n) ON (n.HID);
CREATE INDEX node_uuid_index IF NOT EXISTS FOR (n) ON (n.uuid);
CREATE INDEX node_name_index IF NOT EXISTS FOR (n) ON (n.Name);
CREATE INDEX node_type_index IF NOT EXISTS FOR (n) ON (n.TypeName);
```
Also create per-label uniqueness constraints on `uuid` and `HID`:
```cypher
CREATE CONSTRAINT <label>_uuid_unique IF NOT EXISTS FOR (n:<label>) REQUIRE n.uuid IS UNIQUE;
CREATE CONSTRAINT <label>_hid_unique  IF NOT EXISTS FOR (n:<label>) REQUIRE n.HID  IS UNIQUE;
```
Generated for every label in `graph.LabelFor` for every `identity.AllTypes()` type plus `Mailbox`, `Message`, `User` (messaging data model — added later in S03 but the labels constant lives here so the constraint set covers them on first bootstrap).

Integration test uses `testhelpers.StartNeo4j(ctx, t)` (skip when Docker missing). Asserts that running Bootstrap twice in a row does not error (idempotency). Asserts that after bootstrap `SHOW CONSTRAINTS` returns ≥ 2 × |AllTypes()| rows. Unit tests cover `statements.All()` containing the expected strings.

**Step 5: Commit**
```
git commit -m "$(cat <<'EOF'
feat(schema): Neo4j constraint + index bootstrap per SRS §1.3.6.2.1

Idempotent schema.Bootstrap creates HID/uuid/Name/TypeName indexes and
per-label uniqueness constraints for HID and uuid across all Core Data
Model labels plus messaging labels. Covered by unit + Testcontainers
integration tests.

Refs: SRS §1.3.6.2.1, §2.2.5. Proposes: SSTPA-SCHEMA-BOOTSTRAP-001.
Task: S02-T06.
Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

**Acceptance tests:**
- Unit: `go test ./internal/schema/...` green.
- Integration (with Docker): same command passes including a real bootstrap.
- `make verify` green.

**Evidence:** _________

---

### Task S02-T07: Driver wiring in `cmd/api` with bootstrap-on-start

**Task ID:** S02-T07
**Depends On:** S02-T06
**Difficulty:** medium
**Integration Checkpoint:** no
**Governing SRS sections:** §2.2.4 Backend Software, §2.2.5 Backend Database, §6.1 minimum complexity
**Req IDs touched:** proposes `SSTPA-BACKEND-BOOT-001`

**Files:**
- Modify: `backend/cmd/api/main.go`
- Modify: `backend/internal/config/config.go` (add `Neo4jURI`, `Neo4jUser`, `Neo4jPassword`, `BootstrapSchema` with env fallbacks)
- Create: `backend/internal/neo4jx/driver.go` (constructor + health ping)
- Create: `backend/internal/neo4jx/driver_test.go`

**Agent Briefing:**
The backend currently starts the HTTP server only. Wire Neo4j so that:
1. `cmd/api` reads Neo4j config via `config.Load()` (new fields: `Neo4jURI` default `bolt://neo4j:7687`, `Neo4jUser` default `neo4j`, `Neo4jPassword` required in non-dev, `BootstrapSchema` default `true`).
2. On startup the process constructs a driver via `neo4jx.Open(cfg)`, verifies connectivity with `driver.VerifyConnectivity(ctx)`, and if `cfg.BootstrapSchema` is true runs `schema.Bootstrap(ctx, driver)`. Any failure aborts startup (log + exit 1).
3. The driver is shared with future HTTP handlers via `apihttp.NewRouter(version, driver)` (extend the signature; the health handler keeps the current behavior but the router now owns the driver for future slices).
4. On `SIGINT/SIGTERM` the process calls `driver.Close(ctx)` before exiting so Bolt sessions drain (SRS §2.1 "preserving stored data").

Tests: `neo4jx_test.go` uses the Testcontainers helper to check that `Open` then `VerifyConnectivity` succeeds.

**Step 5: Commit**
```
git commit -m "$(cat <<'EOF'
feat(backend): wire Neo4j driver + schema bootstrap into cmd/api

On startup, cmd/api now constructs a Neo4j 5.26 bolt driver, verifies
connectivity, runs schema.Bootstrap, and shares the driver with the
chi router. SIGTERM triggers ordered driver close before HTTP shutdown
per SRS §2.1.

Refs: SRS §2.1, §2.2.4, §2.2.5. Proposes: SSTPA-BACKEND-BOOT-001.
Task: S02-T07.
Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

**Acceptance tests:**
- `cd backend && go test ./... -v` green (integration tests skip without Docker).
- With the compose stack running, `curl -s localhost:8080/api/v1/health` returns `{"status":"ok",...}` and backend logs show "schema bootstrap complete".

**Out-of-scope:**
- No new HTTP routes here. `/api/v1/nodes/...` lives in S04.

**Evidence:** _________

---

### Task S02-T08: Integration smoke — create Capability with full metadata, read back

**Task ID:** S02-T08
**Depends On:** S02-T05, S02-T06, S02-T07
**Difficulty:** medium
**Integration Checkpoint:** no
**Governing SRS sections:** §1.3.6, §1.3.7, §1.3.7.1
**Req IDs touched:** `SSTPA-IDENTITY-001`, `SSTPA-HID-001`

**Files:**
- Create: `backend/internal/graph/nodetest/capability_test.go`

**Agent Briefing:**
End-to-end integration test: start Neo4j via Testcontainers, bootstrap schema, insert one `(:Capability)` with a full `metadata.Common` populated, then read it back by HID and assert every metadata field round-trips and TypeName equals `"Capability"`. Uses a tiny inline helper `createCapability(ctx, driver, common) error` that writes via `MERGE (c:Capability {HID: $HID}) SET c = $props`. This helper stays in the test file (not production) because the production mutation layer arrives in S03.

Test MUST skip cleanly when Docker is unavailable.

**Step 5: Commit**
```
git commit -m "$(cat <<'EOF'
test(graph): smoke create+read (:Capability) with full common metadata

Testcontainers-backed integration test that exercises Bootstrap +
Common metadata serialisation + HID uniqueness against a real Neo4j.
Production mutation layer arrives in S03.

Refs: SRS §1.3.6, §1.3.7, §1.3.7.1. Reqs: SSTPA-IDENTITY-001,
SSTPA-HID-001. Task: S02-T08.
Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

**Acceptance tests:**
- `cd backend && go test ./internal/graph/nodetest/... -v` green (skipped without Docker).

**Evidence:** _________

---

### Task S02-T09: Slice integration gate — verification updates

**Task ID:** S02-T09
**Depends On:** S02-T08
**Difficulty:** low
**Integration Checkpoint:** yes (slice close)
**Governing SRS sections:** §1.3, §1.3.6, §1.3.7, §2.2.5
**Req IDs touched:** `SSTPA-GRAPH-RULES-001`, `SSTPA-HID-001`, `SSTPA-IDENTITY-001`, `SSTPA-NODETYPE-CATALOG-001`, `SSTPA-RELATIONSHIP-CATALOG-001`, `SSTPA-UUID-001`, `SSTPA-SCHEMA-BOOTSTRAP-001`, `SSTPA-BACKEND-BOOT-001`

**Files:**
- Modify: `docs/verification/shall-register.md`
- Modify: `docs/verification/verification-matrix.md`
- Modify: `docs/superpowers/plans/slices/2026-04-21-slice-02-graph-core-identity.md` (fill Evidence Summary)

**Agent Briefing:**
Promote the Req IDs listed above from Candidate to Approved in the shall-register. For each, add or update the matching row in the verification-matrix with the precise test command and file location. Fill the Evidence Summary block at the bottom of this slice plan with commit SHAs from tasks T01–T08.

Then run `make verify` and capture the full output in the Evidence block. The slice is Accepted when all three gate conditions (task, verification, SRS-section) hold.

**Step 5: Commit**
```
git commit -m "$(cat <<'EOF'
docs(verify): promote S02 requirements to Approved

S02 ships: node-type catalog, relationship catalog, HID formatter, UUID
wrapper, common metadata, schema bootstrapper, backend driver wiring,
integration smoke. All listed SHALLs move Candidate→Approved with
test-command traceability.

Refs: SRS §1.3.2, §1.3.6, §1.3.7, §2.2.5, §2.1.
Reqs: SSTPA-GRAPH-RULES-001, SSTPA-HID-001, SSTPA-IDENTITY-001,
SSTPA-NODETYPE-CATALOG-001, SSTPA-RELATIONSHIP-CATALOG-001,
SSTPA-UUID-001, SSTPA-SCHEMA-BOOTSTRAP-001, SSTPA-BACKEND-BOOT-001.
Task: S02-T09.
Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

**Integration gate criteria:**
- All 8 prior tasks have filled Evidence blocks.
- `make verify` is green.
- `git log --oneline slice/02-graph-core-identity ^main | wc -l` ≥ 9.
- Open PR titled `slice/02: graph core identity`.

**Evidence Summary** (fill at close):
```
Slice: S02 — Graph Core: Schema, Identity, Ownership
Branch: slice/02-graph-core-identity
PR: __________________________________________
verify at slice close: ________________________
Tasks: T01..T09 commit SHAs ____________________
SHALL register deltas:
  SSTPA-NODETYPE-CATALOG-001 → Approved
  SSTPA-GRAPH-RULES-001 → Approved
  SSTPA-RELATIONSHIP-CATALOG-001 → Approved
  SSTPA-HID-001 → Approved
  SSTPA-UUID-001 → Approved
  SSTPA-IDENTITY-001 → Approved
  SSTPA-SCHEMA-BOOTSTRAP-001 → Approved
  SSTPA-BACKEND-BOOT-001 → Approved
Open questions logged for S03:
  - Confirm Cypher 25 feature coverage for transactional mutation use cases.
  - Confirm sequencer/index generator approach under concurrent SoI creation.
```
