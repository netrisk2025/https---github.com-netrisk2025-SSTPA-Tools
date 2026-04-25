# Slice 14 - V58 Goal Keeper and Diagram Persistence (Phase B, Outline)

**Goal:** Implement the V58 Goal Keeper Tool on top of S13-generated Asset/Loss/Root Goal structures. The tool must support GSN node creation/editing, enforce a DAG rooted at one Root Goal, treat Solution nodes as terminal evidence references, and persist visual layout as structured JSON while keeping Neo4j authoritative.

**Architecture:** Extend canonical schema and mutation validation from S02-S04. Add backend services under `backend/internal/goals/**` for Goal Structure retrieval, GSN validation, Solution evidence validation, and GoalStructure JSON persistence. The frontend uses the scaffolded `addons/goal-keeper` package and the shared `DiagramPersistenceContract` from `@sstpa/addon-sdk`.

**Pre-reads:**
- `docs/srs/SSTPA Tool SRS V58.md` as intended version 0.5.8
- `docs/verification/SSTPA_SHALL_Requirements.md` requirements: `1.3.1.8-001`, `1.3.1.8-002`, `1.3.1.8-003`, `1.3.4.13-001`, `1.3.6-001`, `3.4.11.1-001`, `3.4.11.1-005`, `3.4.11.5-002`, `3.4.11.5-003`, `3.4.11.6-001`, `3.4.11.9-001`, `3.4.11.12-001`
- S13 outputs for Asset/Loss/Root Goal generation
- S09 AttackTreeJSON persistence patterns where available
- Conventions doc.

**Invariants:**
- Goal Structure is a DAG rooted at exactly one Root Goal.
- Solution nodes are terminal for `SUPPORTED_BY` and cannot have outgoing support edges.
- Solution evidence references may target Validation, Verification, or Loss nodes only.
- Neo4j nodes and relationships are authoritative for semantics.
- GoalStructure JSON is presentation state only and must be reconciled against the current graph on reopen.
- Stale diagram references are ignored gracefully and surfaced to the user.
- Semantic graph changes and diagram JSON persistence commit transactionally unless the user explicitly elects semantic-only save where supported by V58.

## Tasks

### S14-T01: Goal Keeper detailed design
Difficulty: high - Integration Checkpoint: yes - Files: `docs/architecture/v58-goal-keeper.md`.
**Agent Briefing:** Specify GSN node defaults, relationship allow-list usage, Root Goal identification, GSN ID assignment, diagram JSON schema, stale-reference behavior, and evidence validation. Cite every SHALL ID from this slice.

### S14-T02: Backend Goal Structure retrieval
Difficulty: medium - Files: `backend/internal/goals/query.go`, tests.
**Agent Briefing:** Load a Goal Structure from Asset, Loss, or Goal context. Return GSN nodes, relationships, evidence references, validation status, and persisted GoalStructure JSON where present. Use bounded traversal defaults.

### S14-T03: GSN mutation validation
Difficulty: high - Files: `backend/internal/goals/validate.go`, mutation validation hooks.
**Agent Briefing:** Enforce allowed GSN relationships, duplicate rejection, Root Goal rules, DAG acyclicity, reachability, terminal Solution rules, and same-SoI constraints except explicitly allowed evidence references.

### S14-T04: Solution evidence validation
Difficulty: medium - Files: extend `backend/internal/goals/**`.
**Agent Briefing:** Allow `(:Solution)-[:HAS_VALIDATION]->(:Validation)`, `(:Solution)-[:HAS_VERIFICATION]->(:Verification)`, and `(:Solution)-[:HAS_LOSS]->(:Loss)` only. Mark Solution nodes without evidence incomplete.

### S14-T05: Diagram JSON persistence and reopen reconciliation
Difficulty: high - Integration Checkpoint: yes - Files: `backend/internal/goals/diagram.go`, frontend diagram state helpers.
**Agent Briefing:** Persist node positions, viewport, zoom, collapsed state, filters, display options, and evidence panel state as structured JSON. On reopen, retrieve the authoritative graph, reconcile JSON references, ignore stale references, and report what was ignored.

### S14-T06: Goal Keeper canvas and panels
Difficulty: high - Files: `addons/goal-keeper/**`.
**Agent Briefing:** Implement a usable GSN canvas with create/edit/delete of Goal, Strategy, Context, Justification, Assumption, and Solution nodes. Provide property editing, validation display, evidence association UI, search, and path-to-root display.

### S14-T07: Export and report-ready views
Difficulty: medium - Files: extend `addons/goal-keeper/**`.
**Agent Briefing:** Export Goal Structure as JSON and report-ready image formats. Preserve visible GSN structure, evidence lists, and validation markers.

### S14-T08: Goal Keeper integration gate
Difficulty: medium - Integration Checkpoint: yes.
**Acceptance:** `make verify` green; backend tests cover cycle rejection, duplicate rejection, Solution terminal enforcement, Solution evidence validation, bounded traversal, and stale-reference reconciliation; Playwright covers open-from-Asset/Loss/Goal, edit, save, reopen, and evidence association.

**Verification against SRS V58 SHALL IDs:** Update `docs/verification/shall-register.md` and `docs/verification/verification-matrix.md` for the S14 IDs listed in `docs/superpowers/plans/srs-v58-slice-verification.md`.
