# Slice 03 — Transactional Mutation Layer & Change Notifications (Phase B, Outline)

> **Promotion:** This outline is promoted to Phase A by running `superpowers:writing-plans` with access to the completed S02 code (package paths, function signatures, type names). Fill each task's Steps with verbatim TDD code using the actual S02 symbols.

**Goal:** Introduce the ACID mutation layer. All writes to the SSTPA graph flow through this layer, it computes affected nodes, and it emits owner-change-notification messages in the same transaction. This slice also materialises the Messaging data model nodes (`:User`, `:Mailbox`, `:Message`) per SRS §1.4.4.

**Architecture:** Go package `backend/internal/mutation` exposing a small DSL: `Plan` (one or more node/relationship operations) → `Apply(ctx, driver, user, plan)` returning `CommitReport` (affected HIDs, messages emitted, rollback reason if applicable). Internally: single `ExecuteWrite` session per commit; compute affected-node set by diffing pre/post node and relationship state; synthesise change-notification messages; fail the whole transaction if any notification cannot be written. Sequence and index allocation use Neo4j counters (or `apoc.atomic` equivalents) scoped per label-per-SoI — if APOC isn't available, use `MERGE … ON MATCH SET n.counter = n.counter + 1 RETURN n.counter` pattern on a `(:SequenceCounter {label,soi})` sidecar node.

**Tech Stack:** Go, `neo4j-go-driver/v5` v5 sessions with explicit `ExecuteWrite`, Testcontainers-Go for integration tests.

**Pre-reads:**
- SRS §1.3.7.1 Data Ownership Rules
- SRS §1.3.8.1–1.3.8.24 (type-specific property groups; scope constrained to Capability, System, Requirement, Asset, State for S03 — rest carried over when touched)
- SRS §1.4.4 Messaging Data Model
- SRS §2.2.10.8 Transaction Requirements
- SRS §2.2.10.8.1 Ownership and Change Notification Requirements
- SRS §2.2.10.8.2 Ownership Change Rules
- `backend/internal/identity/`, `backend/internal/metadata/`, `backend/internal/graph/`, `backend/internal/schema/`, `backend/internal/neo4jx/` (S02 outputs)
- `docs/superpowers/plans/2026-04-21-agent-execution-conventions.md`

**Invariants:**
- Every write is single-transaction. No multi-statement mutation that can half-commit.
- If notification generation fails, the whole commit rolls back.
- `Creator` and `CreatorEmail` are immutable except when the committing user is Admin.
- `Owner` updates always pair Owner + OwnerEmail.
- A single commit may emit multiple notifications, but duplicate targets are deduplicated (one message per (Owner, Commit) tuple).
- Messaging data model nodes obey the same identity invariants as Core Data Model nodes.
- Recursive relationships declared DAG in the catalog cannot form cycles through this layer.
- List-returning queries used by this layer are bounded/paginated.

**Tasks:**

### Task S03-T01: Author a transactional-commit RFC
**Task ID:** S03-T01 · **Depends On:** none · **Difficulty:** medium · **Integration Checkpoint:** yes (user review)
**Governing SRS:** §2.2.10.8, §2.2.10.8.1 · **Req IDs:** `SSTPA-OWNERSHIP-TX-001`
**Files:** Create `docs/architecture/mutation-layer-rfc.md`
**Agent Briefing:** Draft a ≤400-line RFC that specifies the Commit DSL (Plan/Operation/Result types), how affected nodes are computed (pre-read snapshot → compute diff), how notifications are built (one per distinct Owner not equal to current user), the failure/rollback contract, the concurrency model (one writer session per commit, optimistic retry on transient errors), and the sequence/index allocation strategy (APOC vs. counter node — pick one and justify against the CLAUDE.md "minimum complexity" rule). End with three worked examples: create `(:Capability)`, create `(:System)` under a `(:Capability)`, update `Owner` of a `(:Requirement)`. Pause for user approval before T02.
**Acceptance:** RFC committed; user-approved before proceeding.

### Task S03-T02: Define `mutation.Plan` and `mutation.Operation` types
**Task ID:** S03-T02 · **Depends On:** S03-T01 · **Difficulty:** medium
**Files:** Create `backend/internal/mutation/types.go`, `types_test.go`
**Agent Briefing:** Implement the Plan/Operation types exactly as specified by the approved RFC. Write table-driven tests for plan validation (e.g., a plan with no ops is invalid; an op referring to an unknown NodeType is invalid).

### Task S03-T03: Affected-node diff utility
**Task ID:** S03-T03 · **Depends On:** S03-T02 · **Difficulty:** medium
**Files:** Create `backend/internal/mutation/affected.go`, `affected_test.go`
**Agent Briefing:** Pure function `ComputeAffected(plan Plan, before, after GraphSnapshot) []AffectedNode` returning the union of nodes whose properties changed, both endpoints of any added/removed relationship. Unit-tested with fixture snapshots (no DB).

### Task S03-T04: Sequence/index allocator
**Task ID:** S03-T04 · **Depends On:** S03-T02 · **Difficulty:** high · **Integration Checkpoint:** yes
**Files:** Create `backend/internal/mutation/sequence.go`, `sequence_test.go`; Modify `backend/internal/schema/bootstrap.go` to add the `(:SequenceCounter)` label if counter-node strategy is chosen.
**Agent Briefing:** Implement per-(label, SoI) sequence generation atomic under Neo4j semantics. Test with concurrent goroutines against a real Testcontainers Neo4j: 100 concurrent requests for the next `:System` sequence in the same SoI must yield 1..100 without duplicates or gaps.

### Task S03-T05: Messaging data model nodes + Cypher fragments
**Task ID:** S03-T05 · **Depends On:** S03-T02 · **Difficulty:** medium
**Files:** Create `backend/internal/messaging/model.go`, `model_test.go`, `messaging/cypher.go`
**Agent Briefing:** Model `(:User)`, `(:Mailbox)`, `(:Message)` per SRS §1.4.4 with exact property names. Provide `EnsureMailbox(ctx, tx, userHID)` and `AppendMessage(ctx, tx, mailboxHID, msg)` helpers usable inside an existing Neo4j transaction (take `neo4j.ManagedTransaction`). Constraints: `:User.UserHash` unique; `:Mailbox.MailboxID` unique; `:Message.uuid` unique.

### Task S03-T06: `mutation.Apply` happy path
**Task ID:** S03-T06 · **Depends On:** S03-T03, S03-T04, S03-T05 · **Difficulty:** high · **Integration Checkpoint:** yes
**Files:** Create `backend/internal/mutation/apply.go`, `apply_test.go`
**Agent Briefing:** Implement `Apply(ctx, driver, user, plan) (CommitReport, error)`. Pseudocode:
1. Begin ExecuteWrite session.
2. Read the "before" snapshot of every node mentioned in the plan.
3. Execute the plan's node/relationship operations within the transaction.
4. Read the "after" snapshot.
5. Compute affected nodes.
6. For each affected node whose Owner ≠ current user, append a CHANGE_NOTIFICATION message into the Owner's mailbox in the SAME transaction.
7. Commit. Return CommitReport with HIDs, message count, recipient list.
Integration test via Testcontainers: create `(:Capability)` as user A, update its Name as user B, verify a message lands in user A's mailbox with the expected payload.

### Task S03-T07: Rollback on notification failure
**Task ID:** S03-T07 · **Depends On:** S03-T06 · **Difficulty:** high · **Integration Checkpoint:** yes
**Files:** Modify `backend/internal/mutation/apply.go`; Create `backend/internal/mutation/rollback_test.go`
**Agent Briefing:** Inject a fault (e.g., missing mailbox for an Owner) and assert the transaction rolls back atomically: the pre-existing `(:Capability)` must retain its unchanged Name, and no Message is created. Covers SRS §2.2.10.8.1 "If required messages fail, the entire transaction SHALL roll back."

### Task S03-T08: Ownership transfer rules
**Task ID:** S03-T08 · **Depends On:** S03-T06 · **Difficulty:** medium
**Files:** Modify `backend/internal/mutation/apply.go`; Create `backend/internal/mutation/ownership_test.go`
**Agent Briefing:** Implement and test SRS §2.2.10.8.2 — ownership change updates Owner, OwnerEmail, LastTouch; does not modify Creator/CreatorEmail; if performed by a user other than the prior Owner, a notification to the prior Owner is generated; the current user may assume ownership only (unless Admin). Admin override tested separately.

### Task S03-T09: Plan validation — DAG + no-duplicate relationships
**Task ID:** S03-T09 · **Depends On:** S03-T02 · **Difficulty:** medium
**Files:** Modify `backend/internal/mutation/types.go`; Create `backend/internal/mutation/validate_test.go`
**Agent Briefing:** Reject plans that would create duplicate logical relationships (same endpoints + type without distinguishing properties) or that would introduce a cycle in a DAG-declared relationship. Validation runs in-transaction before commit so a concurrent write can't slip through. Uses `graph.Catalog()` for DAG flags.

### Task S03-T10: Commit report & test harness for S04
**Task ID:** S03-T10 · **Depends On:** S03-T06 · **Difficulty:** low
**Files:** Modify `backend/internal/mutation/types.go` (CommitReport struct); Create `backend/internal/mutation/harness/harness.go` (test helper `CommitAs(t, driver, user, plan)` used by S04 contract tests).

### Task S03-T11: Slice integration gate
**Task ID:** S03-T11 · **Depends On:** S03-T07, S03-T08, S03-T09, S03-T10 · **Difficulty:** low · **Integration Checkpoint:** yes
**Files:** Modify `docs/verification/shall-register.md`, `docs/verification/verification-matrix.md`, this slice plan (Evidence Summary).
**Acceptance:** `SSTPA-OWNERSHIP-TX-001` → Approved. New Candidates for sequencer, messaging model, rollback behavior, plan validation.

**Integration gate criteria:**
- Every task Evidence block filled.
- `make verify` green on branch.
- SRS citations present in every commit body.
- PR opened: `slice/03: mutation notifications`.

**Evidence Summary (fill at close):** see S02 template.
