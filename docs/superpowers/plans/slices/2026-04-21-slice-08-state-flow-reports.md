# Slice 08 — State Tool, Flow Tool, Reports (Phase B, Outline)

**Goal:** Ship the three next analytical surfaces: the State Tool (SRS §3.4.5), the Flow Tool (SRS §3.4.6 with Functional + STPA Control modes), and the Reports dropdown (SRS §3.4.3.1–3.4.3.3) with the three required reports. The Loss Tool is deliberately excluded — it has a gap-fill and owns its own slice (S09).

**Architecture:** Two addon packages added to the monorepo: `addons/state-tool`, `addons/flow-tool`. Reports live inside the shell (`apps/desktop-shell/src/reports/**`) because they're accessed via a dropdown on the Control Panel, not a pop-up tool. Reports render server-side where possible (backend emits text/Markdown) and client-side rasterises to PDF/Word via a minimal adapter. State Tool and Flow Tool reuse the Requirements Tool's diagram engine chosen in S07-T01 where appropriate, else a local overlay.

**Tech Stack:** React + diagram engine from S07; AG Grid for report tabulations; `pandoc` invoked from backend for Markdown→PDF/docx or a pure-Go equivalent (propose in T01 RFC).

**Pre-reads:**
- SRS §3.4.3 Reports Dropdown Menu and §3.4.3.1–3.4.3.3
- SRS §3.4.5 State Tool (all subsections)
- SRS §3.4.6 Flow Tool (all subsections, both modes)
- SRS §1.3.4 Secondary Relationships (State, Flow semantics)
- `docs/verification/SSTPA_SHALL_Requirements.md` requirements for this slice: `1.3.4.3-001`, `1.3.4.5-001`, `1.3.4.12-001`, `3.4.3.1-001`, `3.4.3.2-001`, `3.4.3.3-001`, `3.4.5.1-001`, `3.4.6.1-001`, `3.4.6.2-002`
- S04 + S07 outputs
- Conventions doc.

**Invariants:**
- State transitions remain a relationship, not a node (SRS §1.3.4 `[:TRANSITIONS_TO]`).
- Flow Tool respects cross-SoI Function-flow prohibition.
- STPA Control Flow mode enforces role validity (Interface can't be a ProcessModel, etc.).
- Reports are read-only artifacts; generation does not mutate the graph except for the `Orphan`/`Barren` update allowed by SRS §3.4.3.3. Legacy imports may normalize `Baron` to `Barren`, but new writes use `Barren`.
- Report formats: text, markdown, MS Word, PDF (all four).

**Tasks:**

### S08-T01: Reporting pipeline RFC
Difficulty: medium · Integration Checkpoint: yes (user approval) · Files: `docs/architecture/reporting-pipeline-rfc.md`.
**Agent Briefing:** Decide between (a) backend emits Markdown + Pandoc for conversions, (b) pure-Go path with `goldmark` + `gofpdf` + `unioffice`. Recommend one; flag any library additions.

### S08-T02: System Description Report backend
Difficulty: medium · Files: `backend/internal/reports/description/**`; route `POST /api/v1/reports/system-description/{soiHID}`.
**Agent Briefing:** Per SRS §3.4.3.1, emits report in text/markdown/docx/pdf. Integration-tested with seeded SoI.

### S08-T03: System Specification Report backend
Difficulty: medium · Files: `backend/internal/reports/specification/**`; route added.
**Agent Briefing:** Per SRS §3.4.3.2.

### S08-T04: Requirement-Traceability Gap Analysis Report
Difficulty: medium · Files: `backend/internal/reports/traceability/**`; route added.
**Agent Briefing:** Per SRS §3.4.3.3, updates `Orphan` and `Barren` properties on `(:Requirement)` as part of report generation — this is the one allowed mutation. ACID.

### S08-T05: Reports dropdown UI
Difficulty: low · Files: `apps/desktop-shell/src/reports/**`; Control Panel wiring.
**Agent Briefing:** Dropdown listing three reports; on selection, prompt for format (text/markdown/docx/pdf), call backend, save file via Tauri.

### S08-T06: State Tool scaffold (pop-up + Observation mode)
Difficulty: medium · Files: `addons/state-tool/**`.
**Agent Briefing:** SRS §3.4.5 Observation mode: read-only diagram of the current SoI's `(:State)` graph; SysML 2 state-transition conventions. Visually distinguishes TransitionKind (functional / countermeasure / both).

### S08-T07: State Tool — Analytical mode
Difficulty: high · Files: extend state-tool.
**Agent Briefing:** Overlays related `(:Hazard)`, `(:Countermeasure)`, `(:Requirement)` nodes and relationships; filters; context-highlighting per SRS §3.4.5.

### S08-T08: State Tool — Edit mode with transaction validation
Difficulty: high · Integration Checkpoint: yes · Files: extend state-tool + backend `backend/internal/mutation/transitions.go`.
**Agent Briefing:** Create/delete `[:TRANSITIONS_TO]`; edit TransitionKind/Trigger/GuardCondition/Rationale/Priority/ResidualRiskNote; require `RequiredByCountermeasureHID/UUID` when TransitionKind in {COUNTERMEASURE_REQUIRED, BOTH}. All through mutation layer.

### S08-T09: State Tool export (PNG/SVG/PDF)
Difficulty: medium · Files: extend state-tool.

### S08-T10: Flow Tool scaffold + Functional Flow mode
Difficulty: high · Files: `addons/flow-tool/**`.
**Agent Briefing:** SRS §3.4.6 Functional Flow Mode: renders `[:FLOWS_TO_FUNCTION]` / `[:FLOWS_TO_INTERFACE]` within the current SoI; cycles allowed but flagged; no cross-SoI expansion.

### S08-T11: Flow Tool — STPA Control Flow mode with role validation
Difficulty: high · Integration Checkpoint: yes · Files: extend flow-tool; backend `backend/internal/mutation/control_flow.go`.
**Agent Briefing:** Casts Functions/Interfaces into ControlAlgorithm / ControlledProcess / ProcessModel / ControlAction / Feedback roles per SRS §3.4.6. Enforces: `(:Interface)` cannot be assigned to `(:ProcessModel)`; one-implements-at-most-one-role rule for both `:Function` and `:Interface`.

### S08-T12: Flow Tool export
Difficulty: low.

### S08-T13: Slice integration gate
Difficulty: low · Integration Checkpoint: yes.
**Acceptance:** new Reqs for three reports, State Tool modes, Flow Tool modes → Approved.

**Integration gate criteria:**
- Playwright demo: run each report, open State Tool, open Flow Tool in both modes, edit a transition, Commit.
- Every generated report validates as Markdown/PDF/docx in unit tests.
- `make verify` green.
