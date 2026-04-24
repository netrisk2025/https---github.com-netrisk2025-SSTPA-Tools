# Slice 09 — Loss Tool & Loss Data Model (Phase B, Outline + Gap-Fill)

**Goal:** Close the second SRS gap the user explicitly called out — the Loss Tool is mentioned but not fully specified. Author the Loss data-model gap-fill (AttackTreeJSON schema, Loss auto-generation rules, DAG constraints), then implement the Loss Tool: auto-generating Loss nodes on Asset commit, rendering the Attack Tree DAG with SAND/XOR semantics, allowing the user to extend/modify it, terminate branches in Attack (residual vulnerability) or a new Derived Asset (which spawns more Loss nodes).

**Architecture:** New addon `addons/loss-tool`. Backend gains `backend/internal/loss/**` for auto-generation and AttackTreeJSON persistence (the DAG is stored as a property on the `(:Loss)` node per SRS §3.4.7). Loss auto-generation is a side effect of Asset commits, executed in the same ACID transaction by the mutation layer. The Attack Tree is a client-side editor on top of the serialized JSON document; mutations to the graph (new Attacks, new Countermeasures, new Derived Assets) route through the standard mutation layer.

**Tech Stack:** Diagram engine chosen in S07-T01 (reactflow or custom SVG); Ajv-compatible JSON Schema validator; existing mutation layer.

**Pre-reads:**
- SRS §1.3.1 `(:Loss)` definition and Attack Tree intent
- SRS §1.3.3/1.3.4 (Loss relationships: `[:HAS_ENVIRONMENT]`, `[:HAS_ELEMENT]`, `[:HAS_STATE]`, `[:HAS_ATTACK]`, `[:HAS_COUNTERMEASURE]`)
- SRS §1.3.8.24 Loss per-type properties
- SRS §3.4.7 The Loss Tool
- Any reference literature the user supplies on attack-tree analysis conventions (Schneier/ADTool/ATSyRA) — T01 blocks until user provides guidance per Master §10 Q2
- Core Data Model §1.3 for `(:Asset)`, `(:Environment)`, Criticality, Assurance
- Conventions doc.

**Invariants:**
- Loss auto-generation happens server-side in the same transaction as the Asset commit.
- One `(:Loss)` per (Asset, Criticality, Assurance) pair; allocation to `(:Environment)` is user-driven.
- Each `(:Loss)` has exactly one Criticality, one Assurance, and one Environment (SRS §1.3.4 constraint).
- AttackTreeJSON is the canonical Attack Tree document; editing the DAG edits this property and any node/relationship mutations it implies.
- Attack Tree branches terminate in `:Attack` (residual vulnerability) or a new `:Asset` (derived, spawns a new Loss).
- SAND (Sequential AND) and XOR operators apply only to intra-tree logic; they are not Neo4j relationships.
- Reference Tool assignment of ATT&CK Techniques to `(:Attack)` nodes remains the canonical source; Loss Tool does not invent attack content.

**Tasks:**

### S09-T01 (gap-fill): Author `gap-fill-loss-data-model.md`
Difficulty: high · Integration Checkpoint: yes (user approval) · Files: `docs/architecture/gap-fill-loss-data-model.md`.
**Agent Briefing:** Deliver a complete Loss Data Model specification as proposal SHALLs. Cover:
1. `AttackTreeJSON` JSON Schema (root = Loss, interior nodes = SAND/XOR operators, leaves = Attack HID or derived-Asset HID). Each tree node carries: uuid, nodeType (operator|attack|asset|loss), operator (sand|xor|null), children (array of nodeId), reference to graph HID where applicable, displayMetadata (label, note).
2. `(:Loss)` property list beyond the common group: `AttackTreeJSON: String`, `AttackTreeVersion: Int`, `LossCriticality: enum`, `LossAssurance: enum`, `Accepted: Bool`, `ResidualRiskNote: String`.
3. Auto-generation rules: on commit of `(:Asset)` N, compute the set of (Criticality, Assurance) pairs currently true on N; for each pair, ensure exactly one `(:Loss)` node exists with `[:HAS_LOSS]` from N, creating missing ones; do not delete extras (user-owned decision). Handle Derived Assets spawned by the tree: when a branch terminates in a new `(:Asset)`, spawn its Losses on Asset commit.
4. DAG rules: AttackTreeJSON must be a DAG; sharing subtrees between branches is allowed but cycles are not.
5. Backend validation endpoint: `POST /api/v1/validate/attack-tree` returning `{valid, reasons[]}`.
6. Loss Tool invocation rules per SRS §3.4.7 (open from Control Panel; if Data Drawer open for Loss with valid AttackTreeJSON, display that; if no valid JSON, create one; else list Loss nodes to pick from).
7. Conflict resolution when two users edit the same Loss (optimistic with AttackTreeVersion).
Pause for user sign-off. If user has prior-art guidance per Master §10 Q2, incorporate it.

### S09-T02: Backend auto-generation hook on Asset commit
Difficulty: high · Integration Checkpoint: yes · Files: `backend/internal/loss/generate.go`, `generate_test.go`; modify `backend/internal/mutation/apply.go` to invoke the hook.
**Agent Briefing:** Integration-tested via Testcontainers. Commit an Asset with (Flight-Critical, Availability) and (Cyber-Critical, Confidentiality) flags; assert two `(:Loss)` nodes appear with `[:HAS_LOSS]`. Re-commit with the same flags — no new Losses appear. Remove a flag — existing Loss retained, no new one.

### S09-T03: AttackTreeJSON schema + validator
Difficulty: medium · Files: `backend/internal/loss/attacktree/{schema.go,validate.go,schema.json,*_test.go}`.
**Agent Briefing:** Load the gap-fill's JSON Schema. Validator enforces DAG. Unit tests cover valid roots, cycle rejection, unknown nodeType rejection.

### S09-T04: Backend validation endpoint
Difficulty: low · Files: `backend/internal/http/validation/attacktree.go`; register with router.
**Agent Briefing:** Wraps S09-T03 validator.

### S09-T05: Loss Tool scaffold (pop-up) + invocation rules
Difficulty: medium · Files: `addons/loss-tool/**`.
**Agent Briefing:** Pop-up mounted via S07-T02 plumbing; implements SRS §3.4.7 invocation: read Data Drawer state; if Loss is open and AttackTreeJSON present, render it; else create new; else list Losses in the SoI.

### S09-T06: Attack Tree canvas with SAND/XOR operators
Difficulty: high · Integration Checkpoint: yes · Files: extend loss-tool.
**Agent Briefing:** Canvas renders the DAG with operator nodes distinct from leaf nodes; user can add/remove/connect; local edits staged; Commit persists new AttackTreeJSON via mutation layer (AttackTreeVersion increments atomically).

### S09-T07: Branch termination — Attack (residual vulnerability)
Difficulty: medium · Files: extend loss-tool; backend mutation support.
**Agent Briefing:** Selecting a branch to terminate in `:Attack` either links to an existing Attack node in the SoI or creates a new one (via Data Drawer commit confirmation pattern from S06-T11).

### S09-T08: Branch termination — Derived Asset
Difficulty: high · Files: extend loss-tool; ties into S09-T02 auto-generation.
**Agent Briefing:** Terminating in a new Asset creates the `(:Asset)` with Derived flag (per SRS §1.3.1), triggers Loss auto-generation for the new Asset (cascades correctly).

### S09-T09: Residual vulnerability acceptance flow
Difficulty: medium · Files: extend loss-tool.
**Agent Briefing:** UI to mark residual vulnerabilities Accepted with a ResidualRiskNote; SRS §3.4.7 acceptability loop.

### S09-T10: Attack Tree export
Difficulty: low · Files: extend loss-tool.
**Agent Briefing:** Export AttackTreeJSON and a rendered PNG/SVG/PDF.

### S09-T11: Slice integration gate
Difficulty: low · Integration Checkpoint: yes.
**Acceptance:** new Reqs (AttackTreeJSON schema, Loss auto-generation, SAND/XOR, derived Asset spawning, acceptance flow) → Approved.

**Integration gate criteria:**
- Playwright demo: commit an Asset, open Loss Tool, build a small Attack Tree with SAND/XOR, terminate branches in Attack + Derived Asset, see the Derived Asset's own Losses appear, Commit.
- Validator rejects cyclic trees.
- `make verify` green.
