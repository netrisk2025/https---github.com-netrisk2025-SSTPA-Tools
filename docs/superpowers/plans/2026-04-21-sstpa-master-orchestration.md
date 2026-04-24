# SSTPA Tools Master Orchestration Plan

> **For agentic workers:** This is a plan-of-plans. It defines the slice decomposition, dependency DAG, orchestration conventions, gap-fill workflow, and acceptance gates for the whole product. Per-slice plans live in `docs/superpowers/plans/slices/`. Slice execution uses `superpowers:subagent-driven-development` (recommended) or `superpowers:executing-plans`.

**Goal:** Deliver SSTPA Tools — a desktop-first systems security engineering application — to a working MVP on a single host, with a clear path to the future split-backend deployment. MVP runs on Linux, Windows, and macOS installers, operates air-gapped on Windows 11 Enterprise, and satisfies the Current-Version SHALL set in `SSTPA Tool SRS V56.md`.

**Architecture:** Tauri + React + TypeScript desktop shell (single window, add-ons in pop-ups) talking over HTTPS via Caddy to a Go/chi REST backend that drives a Neo4j Community Edition graph, with OTel/Prometheus/Tempo/Grafana telemetry. Reference frameworks (NIST 800-53r5, MITRE ATT&CK, MITRE EMB3D) flow through an offline `raw → staged → normalized → imported` pipeline into a read-only subgraph. All mutation is ACID and generates owner notifications in-transaction.

**Tech Stack:** Go 1.24+, `go-chi/chi/v5`, `neo4j/neo4j-go-driver/v5`, Cypher 25; React 19, TypeScript 5.8, Vite 6, Tailwind 3, Framer Motion 12, Zustand 5, TanStack Query 5, react-virtual, Cytoscape.js (Navigator + SoI popup only); Tauri 2; Docker Compose with Caddy 2, Neo4j 5.26-community, OTel Collector 0.122, Prometheus 3, Tempo 2.7, Grafana 11; Testcontainers-Go for integration tests; Playwright for end-to-end UI.

---

## 0. How to Read This Plan

This plan is designed to be executed by an AI orchestrator (Claude-class model) dispatching coding tasks to lower-capability models. Read this file first, then read `2026-04-21-agent-execution-conventions.md`. For any slice you intend to execute, read its per-slice plan. The orchestrator is responsible for:

1. Verifying the slice's prerequisites are complete (green `make verify` on prior slices).
2. If the slice is published as a task-level outline rather than a fully detailed plan, running the `superpowers:writing-plans` skill on the outline to expand it into bite-sized steps using the actual code signatures from completed prerequisite slices. This is a deliberate two-phase design — see §3.
3. Dispatching each task to a coder agent using the agent briefing template in the conventions doc.
4. Reviewing the task result against its acceptance tests before accepting the commit.
5. Updating `docs/verification/shall-register.md` and `docs/verification/verification-matrix.md` when requirements ship.

## 1. Non-Goals for MVP

- No remote backend, no multi-host deployment, no separate client/server packaging (deferred to v2).
- No authentication beyond the "Startup Software" user-selection stand-in required by SRS §2.1. Role-based access (SRS §2.2.10.9) stays placeholder.
- No live internet scraping of reference frameworks — all reference data is provided as offline artifacts per SRS §1.3.10.
- No Neo4j Enterprise features.
- No mobile clients.
- Only the reports explicitly listed in SRS §3.4.3.1–3.4.3.3 are in MVP.

## 2. Inventory: What Exists vs. What Must Be Built

### Present (functional scaffolding)
- Monorepo layout with npm workspaces (`apps/`, `addons/`, `packages/`) and a Go workspace (`backend/`, `tools/reference-pipeline/`).
- Go backend with `/healthz` and `/api/v1/health` chi routes, minimal config loader, one passing test (`backend/internal/http/health_test.go`).
- Tauri 2 + React 19 shell rendering a glass-styled scaffold page, three add-on stubs exporting manifests only (`addons/navigator`, `addons/requirements`, `addons/message-center`).
- Shared packages: `@sstpa/addon-sdk` (ToolModule/ToolManifest types), `@sstpa/domain` (RequirementRecord + hasVerificationGap), `@sstpa/ui` (glass panel className exports).
- Reference-pipeline CLI (`tools/reference-pipeline`) with a `NormalizedReferenceItem` record, JSON Schema, and two unit tests.
- Docker Compose stack wiring Caddy, backend, Neo4j, OTel Collector, Prometheus, Tempo, Grafana (no backend–Neo4j connection yet).
- Make targets: `bootstrap`, `backend-run`, `backend-test`, `reference-run`, `reference-test`, `frontend-dev`, `frontend-build`, `frontend-lint`, `frontend-typecheck`, `frontend-test`, `compose-config`, `verify`.
- Verification skeleton: `docs/verification/shall-register.md` (3 DRAFT entries), `docs/verification/verification-matrix.md` (same 3 entries, all `Planned`).

### Absent (must be built)
- Startup Software / launcher.
- Neo4j driver wiring; any Cypher; HID + uuid generation; ownership metadata; constraint & index bootstrap.
- Any mutation layer, affected-node diffing, change-notification-in-transaction logic.
- Every API endpoint beyond `/health` (node retrieval, hierarchy, search, context, validation, messages, reference framework CRUD of [:REFERENCES]).
- Every per-type property group for the 25 Core Data Model node types.
- All Core Data Model relationships except as listed in docs.
- Reference framework importers for NIST 800-53r5, MITRE ATT&CK, MITRE EMB3D, and the offline intermediate format that feeds them.
- Branding Panel, Control Panel, SoI Panel, Main Panel, Data Drawer — all GUI surfaces beyond the scaffold placeholder.
- Navigator Tool (Cytoscape), Requirements Tool (AG Grid + SysML 2.0 diagram), State Tool, Flow Tool, Loss Tool, Reference Catalog Tool, Message Center real implementation, Reports dropdown and its three reports.
- Help content, tutorial, example Capability.
- Installer for any platform.
- Telemetry content (dashboards, trace exports from backend, metric names).

## 3. The Two-Phase Plan Model

Because Claude-class orchestration cannot productively pre-write thousands of lines of implementation code before earlier architectural choices settle, each slice below exists in one of two states:

- **Phase A — Detailed Plan.** Bite-sized, TDD-style, code-verbatim steps ready for a low-capability coder agent. Slices 01 and 02 ship this way as templates.
- **Phase B — Outline Plan.** A structured task breakdown: scope, file inventory, per-task agent briefing, acceptance tests, dependencies, and difficulty — but not yet verbatim code. Before execution, the orchestrator runs `superpowers:writing-plans` on the outline to produce a Detailed Plan using real code signatures from completed slices.

This is **not a placeholder pattern**. The outline is a fully specified decomposition; what's deferred is only the final keystrokes. The "no placeholders" rule from `superpowers:writing-plans` applies *within a Detailed Plan being executed*, not to this plan-of-plans.

**Promotion rule:** A Phase-B outline is promoted to a Phase-A Detailed Plan only after all its prerequisite slices are `completed` in the task tracker and their `make verify` is green on the current main branch.

## 4. Slice Catalog (Dependency DAG)

All slice plans live under `docs/superpowers/plans/slices/`. Prerequisites must reach "Accepted" (see §7) before a dependent slice begins.

| ID | Title | Phase | Prereqs | One-liner |
|----|-------|-------|---------|-----------|
| S01 | Scaffold Hardening & Verify Pipeline | A (detailed) | none | CI, license headers, testcontainers baseline, `make verify` truly green, copyright banner, reference-data README |
| S02 | Graph Core: Schema, Identity, Ownership | A (detailed) | S01 | Neo4j driver, label catalog, HID + uuid generation, common property group, constraint/index bootstrap |
| S03 | Transactional Mutation Layer & Change Notifications | B (outline) | S02 | Generic mutation DSL, affected-node diff, same-transaction owner notifications, rollback contract, Messaging data model nodes |
| S04 | Backend REST API (Core + Messaging + Reference endpoints) | B (outline) | S03 | chi routes for retrieval, hierarchy, search, validation, context, pagination, messages, `/api/reference`, OTel + Prom middleware |
| S05 | Reference Data Pipeline & Intermediate Format (gap-fill) | B (outline) | S02 | Designs the intermediate format (SRS gap), builds `raw → staged → normalized`, Neo4j importer for NIST 800-53r5, ATT&CK, EMB3D; read-only enforcement |
| S06 | Frontend Foundation: Startup, Shell, Data Drawer | B (outline) | S04 | Startup Software launcher, Tauri lifecycle, Branding & Control Panels, SoI Panel placeholder, Main Panel with one Node Type Section, Data Drawer end-to-end for one node type |
| S07 | Navigator, SoI Selection & Requirements Tool | B (outline) | S06 | Cytoscape Navigator with hierarchy/search/clone modes, SoI Panel, Requirements Tool with SysML 2.0 diagram + traceability matrix |
| S08 | State Tool, Flow Tool, Reports Dropdown | B (outline) | S07 | State Tool, Flow Tool (Functional + STPA Control modes), System Description / System Specification / Requirement Traceability Gap reports |
| S09 | Loss Tool & Loss Data Model (gap-fill) | B (outline) | S07 | Loss data model gap-fill (AttackTreeJSON schema, Loss auto-generation rules), Loss Tool with Attack Tree DAG + SAND/XOR |
| S10 | Message Center, Reference Catalog Tool, Admin & Onboarding | B (outline) | S06 | Message Center real implementation, Reference Catalog Tool, Admin/User onboarding per SRS §1.4.1–1.4.3 |
| S11 | Help, Tutorial, Example Capability (gap-fill) | B (outline) | S08, S10 | Help Data Model gap-fill, Help content scaffolding, runnable Tutorial from Help menu, bundled example Capability |
| S12 | Installer, Air-Gap Bundle, Telemetry Dashboards, Release | B (outline) | S11 | Linux/Windows/macOS installers, air-gapped distribution layout, resource-check dialog, Grafana dashboards, release pipeline |

Dependency DAG:

```
S01 → S02 ──┬─> S03 ─> S04 ─> S06 ─┬─> S07 ─┬─> S08 ──┐
            │                      │        ├─> S09 ──┤
            └─> S05 ────────────────┘        └──┬─────┤
                                 S06 ─> S10 ────┘     │
                                 S08 + S10 ─> S11 ────┤
                                              S11 ─> S12
```

## 5. Gap-Fill Workflow (For SRS Holes Called Out by the User)

The user flagged these gaps: Loss Tool, reference data intermediate format, Help Data Model, Tutorial, Example Capability. Each gets a dedicated Gap-Fill Spec authored *before* its implementation tasks begin. The spec is authored during the first task of the owning slice:

| Gap | Slice | Gap-Fill Spec file | Task |
|-----|-------|---------------------|------|
| Reference intermediate data format & ingestion methodology | S05 | `docs/architecture/gap-fill-reference-intermediate-format.md` | S05-T01 |
| Loss data model (AttackTreeJSON schema, Loss auto-generation rules, DAG constraints) | S09 | `docs/architecture/gap-fill-loss-data-model.md` | S09-T01 |
| Help Data Model (node types, properties, content sourcing) | S11 | `docs/architecture/gap-fill-help-data-model.md` | S11-T01 |
| Tutorial structure (runnable steps, state resets, bundled content) | S11 | `docs/architecture/gap-fill-tutorial.md` | S11-T02 |
| Example Capability bundle (content, HIDs, expected analytics) | S11 | `docs/architecture/gap-fill-example-capability.md` | S11-T03 |

Each Gap-Fill Spec MUST:
1. Propose SHALL statements for addition to `docs/verification/shall-register.md` with `Candidate` status.
2. Call out any contradictions with the SRS draft and route them through a user question (per `CLAUDE.md`: one focused question when SRS is ambiguous).
3. Reach user acceptance before the slice proceeds past its first task. The orchestrator MUST pause for human approval at the Gap-Fill Spec task.

If, during later slices, additional gaps surface (e.g., an unspecified report format), the orchestrator authors a new Gap-Fill Spec in the same directory and adds it to this table.

## 5a. Coverage Strategy for SRS §1.3.8 Type-Specific Property Groups

SRS §1.3.8 defines per-type property groups for all 25 Core Data Model node types plus the messaging types. No slice in this plan owns §1.3.8 as a whole because the type-specific property group for each node type should ship alongside the first tool that requires that type to be editable end-to-end. The ownership map is:

| Node types | Slice that brings property groups online |
|---|---|
| Capability, System, Requirement, Asset, State | S03 (used by mutation + integration smoke) |
| Environment, Interface, Function, Element, Purpose, ControlStructure, Security, Constraint, Connection, Validation, Verification | S06–S07 (Main Panel + Requirements Tool) |
| ControlAlgorithm, ProcessModel, ControlAction, Feedback, ControlledProcess | S08 (Flow Tool STPA Control mode) |
| Control, Countermeasure, Hazard, Attack, Loss | S08–S09 (State + Flow overlays; Loss Tool) |
| Sandbox | S07 (Navigator clone flows reach the Sandbox in Clone Mode) |
| User, Mailbox, Message, Users, Admins, Admin | S03 + S10 |
| ReferenceFramework, ReferenceItem, NistControl, AttackTechnique, etc. | S05 (read-only, pipeline-owned) |
| HelpItem, HelpTutorial, HelpStep, HelpGlossary, HelpDefinition, HelpAnchor | S11 |

Each of those slices MUST include a task that (a) codifies the type's property group as Go + Cypher, (b) wires it into the Data Drawer's property-group registry, and (c) adds integration tests covering defaults, `"Null"` substitution, and editability rules per SRS §1.3.7 / §1.3.8. This is not a placeholder — it's a distributed coverage plan, and each slice's task list above already includes the relevant type(s).

If a future slice introduces a new node type not listed here, that slice owns its property group.

## 6. Cross-Cutting Conventions (Summary)

Full details live in `2026-04-21-agent-execution-conventions.md`. Highlights:

- **Task ID format:** `S{slice:02d}-T{task:02d}` (e.g., `S02-T05`).
- **Commit convention:** `<type>(<area>): <subject>` with body referencing SRS section(s) and Req ID(s). Types: `feat`, `fix`, `test`, `chore`, `docs`, `refactor`.
- **Branch convention:** `slice/<NN>-<short-name>` per slice; each task lands as one commit.
- **Difficulty tiers:** `low` (mechanical, suitable for Haiku-class), `medium` (requires judgement, Sonnet-class), `high` (integration, review-heavy, Opus-class or senior human).
- **Integration checkpoints:** every slice ends with a checkpoint where a higher-capability agent or human runs the full slice `make verify`, reviews the diff, and updates the SHALL register.
- **Test strategy:** unit tests for pure logic, contract tests for API routes, Testcontainers-Go Neo4j integration tests for mutation/notification slices, Playwright for UI end-to-end, Vitest for frontend unit tests.
- **Out-of-scope flag:** any slice task that discovers behavior not listed in the slice plan MUST stop and emit a `BLOCKED:` comment on the task rather than invent scope.

## 7. Acceptance Gates

Each slice has four gates, in order. A slice is `Accepted` only when all four pass.

1. **Task gate** — every task in the slice is `completed` in the task tracker, with a green `make verify` run recorded in the slice plan's "Evidence" section.
2. **Verification gate** — the slice's SHALL set in `docs/verification/shall-register.md` has been updated (Candidate → Approved where appropriate, or Deferred with user sign-off) and `docs/verification/verification-matrix.md` has at least one automated verification row per Approved SHALL.
3. **SRS-section gate** — the slice's completion note in the plan file lists the governing SRS sections (per `CLAUDE.md` Source Of Truth rule) and flags any open questions.
4. **Integration gate** — a higher-capability agent or human performs a structural review. For slices that touch backend mutation or notification behavior this is mandatory, not optional.

## 8. Risk Register

| ID | Risk | Mitigation |
|----|------|------------|
| R1 | Cypher 25 requires Neo4j 5.x features not in Community Edition | Pin to Neo4j 5.26-community and validate Cypher 25 syntax coverage during S02-T03. If a feature is missing, file a Gap-Fill Spec. |
| R2 | Tauri 2 macOS signing is non-trivial | Target Linux + Windows installers in S12 first; macOS signing captured as S12-T10 with explicit "may be deferred" flag. |
| R3 | AttackTreeJSON schema unspecified by SRS | S09-T01 authors gap-fill and blocks until user signs off. Do not allow Loss Tool implementation to proceed without it. |
| R4 | Reference framework license agreements forbid redistribution | S05-T01 gap-fill must enumerate license posture for each framework and produce a provenance manifest rather than committing raw data into the repo. |
| R5 | Air-gapped Windows 11 Enterprise environment unavailable for validation | Treat S12 air-gap acceptance as manual with a written validation procedure; orchestrator produces the procedure, not the run. |
| R6 | Tool surface (Navigator/State/Flow/Loss) requires diagram libraries the SRS didn't pick | SRS §6.2 only names Cytoscape for Navigator + SoI popup. For SysML 2 diagrams (Requirement/State/Flow/Loss) the orchestrator MUST produce a short "diagram engine choice" RFC as the first task of the relevant slice and route through user approval before introducing a new library (per `CLAUDE.md` "minimum complexity" rule). |
| R7 | `CLAUDE.md` forbids adding libraries without approval | Every slice whose first task involves a new runtime library MUST include a one-paragraph justification and explicit user approval before Task 2. |
| R8 | `make verify` is currently a smoke-level check | S01 hardens it into the real acceptance gate used by all later slices. Until S01 is Accepted, no later slice may begin. |

## 9. Interaction Model With Existing Artifacts

- **`CLAUDE.md`** — treated as project-level authoritative guidance. Any conflict with a slice plan is resolved in favor of `CLAUDE.md`; the slice plan is amended, not overridden.
- **`SSTPA Tool SRS V56.md`** — the current draft. Every slice's completion note must cite the SRS sections it implements.
- **`docs/verification/shall-register.md`** — the approved SHALL set. Today it holds 3 `Candidate` entries (DRAFT-001/002/003); S01 will expand it materially.
- **`docs/verification/verification-matrix.md`** — the matrix that closes the loop. Each Approved SHALL must have a row here with an automated or explicitly-documented manual verification.
- **`.opencode/commands/`** — the existing `/verify-slice`, `/srs-audit`, `/traceability`, `/reference-audit` commands remain the recommended tools for the Verification gate; this plan does not replace them.

## 10. Open Questions for the User (to resolve before S05 / S09 / S11 start)

1. **Reference framework artifact delivery** — the SRS says frameworks will be "manually inserted into the database prior to delivery" but does not specify the intermediate artifact. Should S05 target (a) a versioned BLOB in `reference-data/normalized/<framework>/<version>/*.json` committed to the repo, (b) a signed tarball stored externally and fetched at install time, or (c) both with a provenance manifest?
2. **Loss Tool AttackTreeJSON** — is there a prior art schema from structured attack-tree analysis (e.g., Schneier's original, ADTool, ATSyRA) the user wants S09-T01 to align with, or should we design fresh?
3. **Help content sourcing** — will help text be authored in Markdown bundled with the installer, fetched from the graph as `:HelpItem` nodes, or both? (S11-T01 currently assumes graph-first with Markdown seed content.)
4. **Example Capability scope** — rough guidance on system size (small aircraft subsystem? IoT sensor? generic? 5 Systems × 10 Requirements vs. 50×500)?
5. **Tutorial depth** — 5-minute guided overlay, 30-minute hands-on walkthrough, or both?
6. **macOS target** — is macOS a hard requirement for MVP or can S12 ship Linux + Windows first?

## 11. File Map Produced by This Planning Effort

```
docs/superpowers/plans/
  2026-04-21-sstpa-master-orchestration.md          ← this file
  2026-04-21-agent-execution-conventions.md         ← orchestration conventions
  slices/
    2026-04-21-slice-01-scaffold-hardening.md       ← Phase A (detailed)
    2026-04-21-slice-02-graph-core-identity.md      ← Phase A (detailed)
    2026-04-21-slice-03-mutation-notifications.md   ← Phase B (outline)
    2026-04-21-slice-04-rest-api.md                 ← Phase B (outline)
    2026-04-21-slice-05-reference-pipeline.md       ← Phase B (outline + gap-fill)
    2026-04-21-slice-06-frontend-foundation.md      ← Phase B (outline)
    2026-04-21-slice-07-navigator-requirements.md   ← Phase B (outline)
    2026-04-21-slice-08-state-flow-reports.md       ← Phase B (outline)
    2026-04-21-slice-09-loss-tool.md                ← Phase B (outline + gap-fill)
    2026-04-21-slice-10-messaging-admin.md          ← Phase B (outline)
    2026-04-21-slice-11-help-tutorial.md            ← Phase B (outline + gap-fill)
    2026-04-21-slice-12-installer-telemetry.md      ← Phase B (outline)
```

## 12. Execution Entry Point

Begin with **Slice 01**. Do not skip it. S01 establishes the verification rails every later slice depends on.

Approach choice for this plan set:

- **Subagent-Driven (recommended)** — orchestrator dispatches a fresh coder subagent per task, reviews between tasks, and promotes outlines to detailed plans only when prerequisites are Accepted. Use `superpowers:subagent-driven-development`.
- **Inline Execution** — execute Slice 01 inline, then decide whether to stay inline or switch to subagent-driven for remaining slices. Use `superpowers:executing-plans`.

Both approaches are valid; subagent-driven is preferred because the orchestration conventions in this plan are explicitly designed to isolate each task behind a self-contained briefing.
