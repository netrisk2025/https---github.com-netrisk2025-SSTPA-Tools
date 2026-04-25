# Slice 07 — Navigator Tool, SoI Selection & Requirements Tool (Phase B, Outline)

**Goal:** Ship the two tools that unlock real use of the app: the Navigator (SoI selection, hierarchy search, clone flows, connection selection) and the Requirements Tool (hierarchical requirement editing, SysML 2.0-style diagram, traceability matrix).

**Architecture:** Promote the two existing addon stubs (`addons/navigator`, `addons/requirements`) to full implementations. Each mounts in a Tauri pop-up child window. Navigator uses Cytoscape.js + react-cytoscapejs with `cytoscape-fcose` layout. Requirements Tool uses a custom SysML-2-styled React canvas (not Cytoscape — SRS restricts Cytoscape to Navigator + SoI popup) plus AG Grid for the traceability matrix.

**Tech Stack:** Cytoscape.js, react-cytoscapejs, cytoscape-fcose; AG Grid Community; dagre (optional, for Requirements diagram hierarchical layout — propose in T01 RFC if chosen); existing React/Tailwind/Framer/Zustand/TanStack.

**Pre-reads:**
- SRS §3.4.1 Navigator Tool (all subsections, esp. 3.4.1.x Clone Modes and 3.4.1 constraints)
- SRS §3.4.2 Requirements Tool (all subsections)
- SRS §3.4.3.3 Requirement Traceability Gap Analysis (feeds the tool's `Orphan`/`Barren` flags; legacy `Baron` is normalized only at import/migration boundaries)
- SRS §6.2 UI Tech Stack (Cytoscape scope; AG Grid scope)
- S06 outputs
- Conventions doc.

**Invariants:**
- Navigator stays a pop-up — does not edit nodes except via clone operations.
- Clone operations are transactional via the mutation layer (SRS §3.4.1 clone requirements).
- Cloned nodes receive fresh HID/uuid and metadata per SRS clone rules.
- Requirements Tool editing flows still route through Data Drawer commits where SRS demands (it's a specialised view, not a new edit surface).
- Cytoscape is not used in the Requirements Tool.
- AG Grid is only for the traceability matrix in this slice.

**Tasks:**

### S07-T01: Diagram-engine choice RFC for Requirements Tool
Difficulty: medium · Integration Checkpoint: yes (user approval) · Files: `docs/architecture/requirements-diagram-rfc.md`.
**Agent Briefing:** Per SRS §6.1 "minimum complexity" + §6.2 Cytoscape scope, produce a short RFC comparing options for the SysML 2.0-conformant Requirement diagram: (a) custom SVG with dagre layout, (b) reactflow, (c) other. Recommend one and list the library addition cost. Pause for user sign-off.

### S07-T02: Cytoscape mount + Tauri child-window plumbing
Difficulty: medium · Files: `apps/desktop-shell/src/popup/**`, `src-tauri/src/popup.rs` · Deps: S06 shell.
**Agent Briefing:** Helper `openToolPopup(toolId, props)` creates a Tauri WebView window and mounts the addon's React root. Lifecycle: close button; never touches the main window's SoI unless the tool returns a `{type: 'select-soi', ...}` message.

### S07-T03: Navigator graph view (Explore Mode)
Difficulty: high · Integration Checkpoint: yes · Files: `addons/navigator/src/**`.
**Agent Briefing:** Cytoscape.js with `fcose` layout; default shows `:Capability` + `:System` nodes, roll-down reveals other primaries; current SoI visually distinct; smooth zoom/pan/animated re-centering (SRS §3.4.1).

### S07-T04: Navigator search interface
Difficulty: medium · Files: extend `addons/navigator/src/**` · Deps: S04-T06.
**Agent Briefing:** Search by HID/uuid exact + Name/ShortDescription partial; exact matches fast-path per SRS §3.4.1. Results update the graph view with highlighting.

### S07-T05: Navigator SoI Selection Mode
Difficulty: medium · Files: extend `addons/navigator/src/**`.
**Agent Briefing:** "Select Mode": on confirmation, selected node becomes the current SoI, main window's panels refresh. Cancel returns to Explore. Emits message to the main window via Tauri IPC.

### S07-T06: Navigator Clone-Target Selection Mode
Difficulty: medium · Files: extend `addons/navigator/src/**`.
**Agent Briefing:** Select-without-SoI-change mode. Invalid parent nodes visually muted; valid targets selectable. Emits the selected target HID back to the caller (Data Drawer or calling tool).

### S07-T07: Clone Mode (Node Only)
Difficulty: high · Integration Checkpoint: yes · Files: `addons/navigator/src/**`; backend support in `backend/internal/mutation/clone.go` · Deps: S04-T09.
**Agent Briefing:** Implements SRS §3.4.1 Clone Node flow: clone a single node with no retained relationships; attach to user-selected parent via a valid Core Data Model relationship; new HID/uuid, fresh metadata. Transactional.

### S07-T08: Clone Mode (Node + Requirements)
Difficulty: high · Files: extend clone.
**Agent Briefing:** Clones selected node + `[:HAS_REQUIREMENT]` children; cloned Requirement.Orphan=true; attached to parent SoI's `(:Purpose)` via `[:HAS_REQUIREMENT]`. Per SRS §3.4.1 Clone Node + Requirements.

### S07-T09: Connection selection / participation
Difficulty: medium · Files: extend `addons/navigator/src/**`; backend assist in `backend/internal/mutation/interfaces.go` if needed.
**Agent Briefing:** Allow the user to connect an SoI Interface to a Connection owned by another System. Enforces SRS §1.3.3.1 constraint that Connection ownership does not imply all participating Interfaces belong to owner.

### S07-T10: Requirements Tool — hierarchy view
Difficulty: high · Files: `addons/requirements/src/**`; diagram engine chosen in T01.
**Agent Briefing:** SRS §3.4.2 Hierarchy View: focused requirement plus heritage/lineage to user-selected depth; SysML 2 conventions; distinct node shapes/colors; user-toggleable display properties.

### S07-T11: Requirements Tool — allocation view
Difficulty: medium · Files: extend `addons/requirements/src/**`.
**Agent Briefing:** Displays the valid node with `[:HAS_REQUIREMENT]` relationships and all allocated requirements; zoom/pan/filter.

### S07-T12: Requirements Tool — create/edit/delete
Difficulty: high · Integration Checkpoint: yes · Files: extend tool + backend mutation · Deps: S04-T09.
**Agent Briefing:** Creation of `(:Requirement)` nodes with required fields; edit via Data Drawer; delete with orphan detection per SRS §3.4.2. All backend interactions transactional ACID.

### S07-T13: Requirements Traceability Matrix (AG Grid)
Difficulty: medium · Files: `addons/requirements/src/matrix/**`; library addition approval for `ag-grid-community` + `ag-grid-react`.
**Agent Briefing:** Traceability matrix columns: Requirement HID/Name/Status, Allocated-To HIDs, Verifications, Validations, Orphan/Barren flags per SRS §3.4.3.3. Exportable per tool subsection.

### S07-T14: Requirements Tool — SysML 2 export (PNG/SVG/PDF)
Difficulty: medium · Files: extend tool.
**Agent Briefing:** Export diagrams per SRS §3.4.2 export requirements.

### S07-T15: Slice integration gate
Difficulty: low · Integration Checkpoint: yes.
**Acceptance:** Navigator (all modes), Clone flows, Requirements Tool (views, edit, matrix, export) → Approved. Library additions (Cytoscape ecosystem, AG Grid, diagram engine) recorded with user approvals.

**Integration gate criteria:**
- Playwright demo: open Navigator, change SoI, clone a Requirement, open Requirements Tool, edit a Requirement, export traceability matrix.
- `make verify` green.
