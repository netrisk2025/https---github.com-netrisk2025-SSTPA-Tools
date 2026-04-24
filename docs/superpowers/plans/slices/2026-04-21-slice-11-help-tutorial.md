# Slice 11 — Help, Tutorial, Example Capability (Phase B, Outline + Gap-Fill)

**Goal:** Close three user-flagged gaps: Help Data Model (SRS §1.6 is a stub), Tutorial (not in SRS), and Example Capability (not in SRS). Deliver a runnable in-app help system, a guided tutorial triggered from the Help menu, and a seeded example Capability the user can explore immediately after installation.

**Architecture:** Help content lives as `(:HelpItem)` nodes in the graph (authored graph-first per gap-fill T01), seeded from Markdown bundled with the installer. Tutorial is a scripted overlay on top of existing tools; no forks of the tool code. Example Capability is a normalised graph fixture loaded by the installer (or by `make seed-example`) into the running Neo4j instance.

**Tech Stack:** Existing React/Tailwind/Framer; Markdown renderer (`react-markdown` + `remark-gfm` — propose in T01); new backend package `backend/internal/help`.

**Pre-reads:**
- SRS §1.6 Help Data Model (current stub)
- SRS §3.3 Branding Panel (where Help icon sits; confirm in T01)
- Master §10 Q3, Q4, Q5 (open questions)
- Outputs from all prior slices
- Conventions doc.

**Invariants:**
- Help content is read-only to users; Admin may edit via a dedicated Admin tool surface (out of scope for MVP unless user asks).
- Tutorial does not mutate any real SoI; it runs against a disposable sandbox SoI created for the tutorial session and deleted on exit.
- Example Capability ships as idempotent seed; running the seeder twice does not duplicate.
- Tutorial is skippable at any step.

**Tasks:**

### S11-T01 (gap-fill): Author `gap-fill-help-data-model.md`
Difficulty: medium · Integration Checkpoint: yes (user approval) · Files: `docs/architecture/gap-fill-help-data-model.md`.
**Agent Briefing:** Propose:
1. `(:HelpItem)` nodes with properties `{HelpID, Title, Body (Markdown), Category, RelatedNodeTypes[], Anchor (GUI field/input id), Version}`.
2. `(:HelpTutorial)` nodes composed of ordered `(:HelpStep)` children.
3. `(:HelpGlossary)` for SSTPA terminology definitions.
4. Relationships: `(:HelpItem)-[:EXPLAINS]->(:HelpAnchor)`, `(:HelpTutorial)-[:HAS_STEP]->(:HelpStep)`, `(:HelpGlossary)-[:HAS_DEFINITION]->(:HelpDefinition)`.
5. Content sourcing: author Markdown in `docs/help/<category>/<id>.md` bundled with installer; `refimport help` loads into Neo4j idempotently.
6. Backend read-only endpoints under `/api/v1/help/*`.
Pause for user sign-off.

### S11-T02 (gap-fill): Author `gap-fill-tutorial.md`
Difficulty: medium · Integration Checkpoint: yes (user approval) · Files: `docs/architecture/gap-fill-tutorial.md`.
**Agent Briefing:** Propose the tutorial UX (overlay cards, "Next"/"Skip"/"Done"), sandbox-SoI lifecycle, which steps exercise which tools, how progress persists per user. Target length per Master §10 Q5. Pause for sign-off.

### S11-T03 (gap-fill): Author `gap-fill-example-capability.md`
Difficulty: medium · Integration Checkpoint: yes (user approval) · Files: `docs/architecture/gap-fill-example-capability.md`.
**Agent Briefing:** Propose an example system (size/domain per Master §10 Q4). Enumerate every node with HID/Name and every relationship. Include at least one example each of: Capability, System, Environment, Interface, Function, Element, Purpose, State, ControlStructure, Asset, Security, Connection, Requirement, Countermeasure, Control, Attack, Hazard, Loss. Pause for sign-off.

### S11-T04: Help content seed + backend
Difficulty: medium · Files: `backend/internal/help/**`; route `GET /api/v1/help/items`, `GET /api/v1/help/items/{id}`, `GET /api/v1/help/by-anchor/{anchor}`; seed loader CLI `tools/reference-pipeline/cmd/refimport` gains `help` subcommand.
**Agent Briefing:** Load authored Markdown into `(:HelpItem)` nodes. Idempotent. Content ships in `docs/help/**`.

### S11-T05: Help panel + field-level tooltips in Data Drawer
Difficulty: medium · Files: `apps/desktop-shell/src/help/**`; modify Data Drawer to fetch anchor-keyed help.
**Agent Briefing:** Hover/`?` icon on each field shows short help; expandable to full help item. Reads from `/api/v1/help/by-anchor/*` with TanStack Query caching.

### S11-T06: Help menu in Branding Panel
Difficulty: low · Files: modify `apps/desktop-shell/src/layout/BrandingPanel.tsx`.
**Agent Briefing:** Help icon in the Branding Panel opens a Help pop-up; listing top-level Help categories, search, glossary.

### S11-T07: Tutorial runtime
Difficulty: high · Integration Checkpoint: yes · Files: `apps/desktop-shell/src/tutorial/**`.
**Agent Briefing:** Implements the approved tutorial UX. Sandbox SoI is created at start, deleted at exit. Steps are Markdown-rendered inside overlay cards. Selectors target real shell elements with ARIA/data-attribute hooks.

### S11-T08: Example Capability seeder
Difficulty: medium · Files: `tools/example-capability/cmd/seed/main.go`, `tools/example-capability/fixtures/**.cypher`; Makefile `seed-example` target.
**Agent Briefing:** Idempotent Cypher seeder matching the approved fixture. Integration test via Testcontainers: seed twice, assert node counts equal single-seed counts.

### S11-T09: Tutorial hooks into Example Capability
Difficulty: medium · Files: extend tutorial module.
**Agent Briefing:** Tutorial uses the Example Capability so users see real data during the walkthrough.

### S11-T10: Bundled Markdown help content (initial set)
Difficulty: medium · Files: `docs/help/**/*.md`.
**Agent Briefing:** Author help pages for at least: SoI selection, Data Drawer, Commit flow, Navigator, Requirements Tool, Reference Catalog, Loss Tool, Reports, Message Center, Glossary (20+ SSTPA terms). Low-capability agent ok for drafting; orchestrator reviews for accuracy against SRS.

### S11-T11: Slice integration gate
Difficulty: low · Integration Checkpoint: yes.
**Acceptance:** Gap-fill SHALLs (Help Data Model, Tutorial, Example Capability) promoted Approved. Reqs for Help menu, field tooltips, tutorial runtime, example seeder → Approved.

**Integration gate criteria:**
- Playwright demo: fresh install → seed example → open Help menu → run Tutorial to completion → verify no residual sandbox SoI in Neo4j.
- `make verify` green.
