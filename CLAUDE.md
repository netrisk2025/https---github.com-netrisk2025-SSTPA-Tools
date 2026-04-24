# SSTPA Tool

## Source Of Truth
- `SSTPA Tool SRS V56.md` is the governing specification for this project.
- Treat every `shall` as mandatory and every `should` as mandatory unless the user explicitly approves a deferral.
- For any non-trivial task, identify the governing SRS sections before editing and mention them in the final summary, commit notes, or PR summary.
- If the SRS is ambiguous, internally inconsistent, or appears to conflict with prior code, ask one focused question instead of guessing.
- Sections described as placeholder or future-version work require explicit confirmation before implementation.

## Spec Status
- `docs/srs/source/SSTPA Tool SRS V56.md` is the current upstream draft.
- The SRS is being reorganized. Treat it as authoritative background, but route active implementation planning through the approved requirement set in `docs/verification/shall-register.md` when entries exist.
- Do not invent new mandatory requirements. If the draft SRS, approved SHALL register, and code disagree, stop and ask one focused question unless a repo document explicitly states it supersedes the draft.
- Placeholder or future-version sections still require explicit confirmation before implementation.

## Current Architecture Direction
- Keep the repository split into independent slices: Go backend, Tauri shell, add-on tool packages, and reference-data tooling.
- Add-on tools must remain independently developable. Prefer contracts in `packages/addon-sdk/` and avoid burying tool-specific behavior directly inside the shell.
- Reference frameworks must flow through an intermediate data pipeline before graph import: raw source -> staged extraction -> normalized records -> backend import.
- Preserve the required stack from the draft unless the user explicitly approves a change: Go + `chi` + Neo4j + Caddy + OTel/Prometheus/Tempo/Grafana on the backend, Tauri + React + TypeScript + Vite + Tailwind + Framer Motion + Zustand + TanStack Query + `react-virtual` + Cytoscape.js on the frontend.

## Delivery Strategy
- Build in thin, end-to-end slices that leave the project runnable and testable.
- Prefer this implementation order:
  1. repo scaffolding and developer tooling
  2. core Neo4j data model, HID/uuid generation, ownership metadata, and relationship constraints
  3. Go backend REST API, transactional mutation layer, and telemetry wiring
  4. Tauri/React shell, SoI selection, and Data Drawer commit flow
  5. external reference framework import, search, and assignment
  6. reports, MBSE tools, and installer work
- Keep early work aligned to the current-version single-machine workflow before remote or enterprise extensions.

## Core Graph Rules
- Use singular Neo4j labels and `UPPERCASE_SNAKE_CASE` relationship names.
- Do not create reverse relationships without a concrete semantic or performance need.
- Represent missing property values as `"Null"` unless a newer approved requirement says otherwise.
- Prevent duplicate logical relationships.
- Recursive edges and traversals must stay bounded and explicitly governed as DAG or cyclic-by-design.
- All list-returning APIs must support pagination and maximum limits.

## Identity, Ownership, And Transactions
- Every graph node needs `HID`, `uuid`, `TypeName`, `Owner`, `OwnerEmail`, `Creator`, `CreatorEmail`, `Created`, `LastTouch`, and `VersionID`.
- HID format stays `{TYPE}_{INDEX}_{SEQUENCE}`.
- System nodes use sequence `0`; non-system sequence numbers are unique per node type inside an SoI.
- Mutations are ACID and must create required owner notifications in the same transaction when non-owners change data they do not own.
- Reference framework data is read-only to users and must stay separate from SoI graph data.

## Repo Commands
- Install JavaScript dependencies: `make bootstrap`
- Run backend API: `make backend-run`
- Test backend: `make backend-test`
- Exercise reference CLI: `make reference-run`
- Test reference pipeline: `make reference-test`
- Run desktop shell dev server: `make frontend-dev`
- Build TypeScript workspaces: `make frontend-build`
- Lint TypeScript workspaces: `make frontend-lint`
- Typecheck TypeScript workspaces: `make frontend-typecheck`
- Run TypeScript workspace tests: `make frontend-test`
- Validate Docker Compose: `make compose-config`
- Run the automated verification set: `make verify`

## Verification Workflow
- First reduce or clarify the active SHALL set in `docs/verification/shall-register.md` when the user is refining scope.
- For any implemented requirement slice, update or reference `docs/verification/verification-matrix.md`.
- Prefer focused automated tests over broad placeholders.
- Use OpenCode subagents for targeted review and verification:
  - `/verify-slice`
  - `/srs-audit`
  - `/traceability`
  - `/reference-audit`

## Delivery Strategy
- Build thin, end-to-end slices that leave the repo runnable and testable.
- Prefer this order unless the user redirects it:
  1. repo scaffolding and tooling
  2. core graph identity and ownership rules
  3. transactional backend mutations and validation
  4. shell plus SoI selection and Data Drawer flow
  5. reference pipeline and reference graph import
  6. reports, MBSE tools, and installer work

## Change Management
- Keep commits small and traceable to a single requirement slice.
- Mention the governing SRS sections or approved requirement IDs in completion notes, commit messages, and PR summaries.
- Avoid combining unrelated backend, frontend, and reporting work unless required for a working vertical slice.


## Required Stack
- Backend: Go, `chi`, Neo4j Community Edition, Cypher 25, Caddy, OpenTelemetry Collector, Prometheus, Tempo, Grafana, Docker Compose. See SRS section 2.2.
- Frontend: Tauri, React, TypeScript, Vite, Tailwind CSS, Framer Motion, Zustand, TanStack Query, `react-virtual`, Cytoscape.js for the Navigator only. See SRS section 6.2.
- Target deployment is air-gapped Windows 11 Enterprise. Ubuntu Studio 25.04 is also a supported development platform. See SRS section 6.
- Do not replace the named stack or add large alternate frameworks without explicit approval.

## Core Graph Invariants
- Use singular node labels and `UPPERCASE_SNAKE_CASE` relationship names in Neo4j syntax. See SRS section 1.3.2.
- Do not add reverse relationships unless there is a clear semantic or performance requirement.
- Represent missing property values as `"Null"` unless the SRS defines a different fixed default.
- Prevent duplicate logical relationships.
- Recursive relationships must declare DAG vs cyclic-by-design behavior, and recursive traversals must always be bounded.
- All list-returning endpoints must support pagination and maximum result limits.
- Default cross-system modeling goes through `(:Connection)`.

## Identity, Ownership, And Transactions
- Every node must have `HID`, `uuid`, `TypeName`, `Owner`, `OwnerEmail`, `Creator`, `CreatorEmail`, `Created`, `LastTouch`, and `VersionID`. See SRS sections 1.3.6 and 1.3.7.
- HID format is `{TYPE}_{INDEX}_{SEQUENCE}` using the SRS node type identifiers.
- System nodes always use sequence `0`. Non-system sequence numbers are unique per node type within an SoI.
- Creation and mutation must update ownership and touch metadata correctly.
- When a user changes data or relationships they do not own, the same ACID transaction must also create the required internal mailbox notification. If notification creation fails, the mutation fails. See SRS sections 1.3.7.1 and 2.2.10.8.1.

## Backend Rules
- All mutations are transactional and ACID compliant. Validate relationships before commit. See SRS section 2.2.10.
- Neo4j must never be publicly exposed; only the reverse proxy is internet-facing.
- Imported framework data is read-only to users, stored separately from SoI graphs, and must not be scraped live from the internet. See SRS section 1.3.10 and 2.2.10.10.
- Prefer integration and contract tests for HID generation, clone behavior, relationship validation, pagination, and rollback rules.

## Frontend Rules
- The GUI is a desktop app with one main window. Add-on tools open in pop-ups. See SRS section 3.
- The Data Drawer is the only general edit surface.
- Editing outside the current SoI is blocked except for explicit clone flows allowed by the SRS.
- Preserve the dark liquid-glass visual language; do not drift into a generic admin-dashboard UI.
- Use Cytoscape.js only for the Navigator and SoI graph, not for the main editor surface.

## Change Management
- Keep commits small and traceable to one coherent SRS slice.
- Explain which SRS sections were affected and why in completion notes and PR summaries.
- Do not combine unrelated backend, frontend, and reporting work unless it is required for a working vertical slice.
- After the repo scaffold exists, update this file with exact build, lint, test, and run commands for each workspace.
