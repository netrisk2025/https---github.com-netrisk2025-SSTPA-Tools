# Slice 06 — Frontend Foundation: Startup, Shell, Data Drawer (Phase B, Outline)

**Goal:** Make the desktop app real. Replace the scaffold page with the SRS-required layout: Branding Panel, Control Panel, SoI Panel, Main Panel. Implement the Startup Software, the Data Drawer as the only general edit surface, and commit/cancel/close flows that talk to the S04 backend. End-to-end goal: a user can launch the app, pick a `:Capability`, see a seeded Main Panel, open a Data Drawer for one node type, stage a change, Commit, and observe the change materialise.

**Architecture:** `apps/desktop-shell` gains `src/layout/*` (panel components), `src/startup/*` (launcher flow), `src/drawer/*` (Data Drawer), `src/api/*` (TanStack Query client wrapping the backend OpenAPI spec), `src/store/*` (Zustand stores for current user, current SoI, staged edits). The addon stubs stay out of the way — they're wired in S07+. Tauri's Rust side gains a small command set for launching/terminating the bundled backend during Startup.

**Tech Stack:** React 19, TypeScript 5.8, Vite 6, Tailwind 3, Framer Motion 12, Zustand 5, TanStack Query 5, react-virtual, Tauri 2 commands; Vitest for unit tests; Playwright for end-to-end.

**Pre-reads:**
- SRS §2.1 Startup Software
- SRS §3.1–3.3 GUI Overview, Style, Branding Panel
- SRS §3.4 Control Panel
- SRS §3.5 System of Interest Panel
- SRS §3.6 Main Panel (all subsections)
- SRS §3.7 Data Drawer (all subsections, especially 3.7.5.1 Commit Notification Behavior)
- SRS §6.2 UI Tech Stack
- Existing `apps/desktop-shell/**`, `packages/ui/**`, `packages/addon-sdk/**`, `packages/domain/**`
- Conventions doc.

**Invariants:**
- Single window; add-ons live in pop-ups only.
- Data Drawer is the only general edit surface.
- Only one Data Drawer open at a time.
- All commits go through the Backend; Frontend does not compute notification recipients.
- Dark liquid-glass aesthetic stays; no generic admin dashboards.
- Cytoscape.js is only for Navigator + SoI selection popup (S07).
- AG Grid is only for traceability matrix + search results + report tabulations (S07/S08).
- Empty string property values render as `"Null"`.

**Tasks:**

### S06-T01: Playwright harness + CI gate
Difficulty: medium · Files: `apps/desktop-shell/e2e/*.spec.ts`, `playwright.config.ts`; Makefile target `frontend-e2e`; add to `make verify`.
**Agent Briefing:** Configure Playwright to run against `vite dev` in CI mode. Library addition flagged; user approval required in the task's briefing.

### S06-T02: API client + types generated from OpenAPI
Difficulty: medium · Files: `packages/api-client/**` (new workspace); wire into `tsconfig.base.json` paths · Deps: S04-T14.
**Agent Briefing:** Generate TS types from `backend/api/openapi.yaml` via `openapi-typescript`. Expose a thin `fetch`-based client with request-ID propagation and TanStack Query integration (`queryKey` factory per route).

### S06-T03: Zustand stores for user, SoI, drawer
Difficulty: low · Files: `apps/desktop-shell/src/store/{user.ts,soi.ts,drawer.ts,index.ts,*.test.ts}`.
**Agent Briefing:** Stores: currentUser (hid, name, email, isAdmin), currentSoI (hid/uuid/name/shortDescription or null), drawer (openNodeHID|null, staged diff, isDirty). Unit tests cover staging/reset/commit transitions.

### S06-T04: Startup Software flow (S-launcher)
Difficulty: high · Integration Checkpoint: yes · Files: `apps/desktop-shell/src/startup/**`; Tauri commands in `src-tauri/src/**`.
**Agent Briefing:** Implements SRS §2.1: launch from icon/CLI, dialog with default theme + animation, theme picker, connect-or-start Backend on local machine, list users (via `GET /api/v1/users` — propose adding this endpoint if not already in S04; if missing, pause and escalate), user select/add, then launch the main shell. On shutdown, gracefully stop the backend and Neo4j. Playwright test covers happy path.

### S06-T05: Branding Panel
Difficulty: low · Files: `apps/desktop-shell/src/layout/BrandingPanel.tsx` + test.
**Agent Briefing:** Logo left, "SSTPA Tools" + version center; right: connection IP/port/status in Courier contrast color, user name, mail icon (stub wiring; real popup is S10), gear icon (Settings stub). Framer Motion subtle hover effects. Unread indicator pending real messaging in S10.

### S06-T06: Control Panel (icon row + Shutdown)
Difficulty: low · Files: `apps/desktop-shell/src/layout/ControlPanel.tsx` + test.
**Agent Briefing:** Icons for Navigator, Requirements, State, Flow, Loss, Reports, Reference Catalog, Message Center, Shutdown (red power icon right-aligned). Clicking any tool without real wiring opens the SRS-mandated "Under Construction" alert. Real wiring happens in later slices.

### S06-T07: SoI Panel
Difficulty: low · Files: `apps/desktop-shell/src/layout/SoIPanel.tsx` + test.
**Agent Briefing:** Displays HID, Name, ShortDescription when a current SoI exists; "Select a System of Interest" placeholder otherwise. Not editable (SRS §3.5). Reads from the soi Zustand store.

### S06-T08: Main Panel skeleton with one Node Type Section
Difficulty: medium · Files: `apps/desktop-shell/src/layout/MainPanel.tsx`, `src/layout/nodeTypeSection/**`, uses `react-virtual` · Deps: S06-T02.
**Agent Briefing:** Implement the SRS §3.6 hierarchical card interface for ONE Node Type Section (Requirement) end-to-end: fetch `(:Requirement)` nodes for the current SoI, render as collapsible cards with HID/Name/ShortDescription. Progressive disclosure stubbed for other types. Virtualised list.

### S06-T09: Data Drawer shell + Zustand plumbing
Difficulty: medium · Files: `apps/desktop-shell/src/drawer/DataDrawer.tsx`, `src/drawer/sections/*.tsx`, `src/drawer/hooks/*.ts` · Deps: S06-T03.
**Agent Briefing:** Slide-in from right (Framer Motion). Header (Node Type, Name, HID; Commit/Cancel/Close buttons). Property groups collapsible. All empty displayed as `"Null"`. Only one drawer open at a time. Relationship groups below type-specific properties — display only in this task; add/remove/associate arrives in later tasks.

### S06-T10: Data Drawer staged editing
Difficulty: medium · Files: extend `src/drawer/**` and `src/store/drawer.ts`.
**Agent Briefing:** Edits update the drawer's `staged` object, not the fetched node. Cancel reverts. Commit button disabled unless `isDirty`. Close-while-dirty triggers the SRS-required confirmation.

### S06-T11: Data Drawer commit + notification summary
Difficulty: high · Integration Checkpoint: yes · Files: extend `src/drawer/**` · Deps: S04-T09.
**Agent Briefing:** Commit sends the staged delta to the backend via the S04 mutation endpoint. On success, consume `CommitReport` and render the dialog summary per SRS §3.7.5.1: nodes changed, relationships changed, messages generated, recipients notified. On failure (including notification failure) display "overall commit failure" and do not clear the drawer. Covered by Playwright end-to-end test.

### S06-T12: Relationship Groups — add/remove/associate flows
Difficulty: high · Integration Checkpoint: yes · Files: extend `src/drawer/**`.
**Agent Briefing:** Implement SRS §3.7.4: "Add" / "Associate" / "Remove" actions per relationship group, including the Create-related commit dialog and the orphaned-node warning on removal.

### S06-T13: Out-of-SoI edit guard
Difficulty: medium · Files: `src/drawer/OutOfSoIAlert.tsx` + logic.
**Agent Briefing:** SRS §3.7.6 — attempts to edit a node outside the current SoI produce the "Navigate to: {HID} to edit" alert with copy-HID icon.

### S06-T14: Consistency + accessibility sweep
Difficulty: low · Files: `apps/desktop-shell/src/**`; add Storybook-free visual test via Playwright screenshot on key layouts.
**Agent Briefing:** Ensure glass aesthetic consistent; add keyboard navigation basics; document the style primitives in `packages/ui`.

### S06-T15: Slice integration gate
Difficulty: low · Integration Checkpoint: yes.
**Acceptance:** new Candidates (Branding, Control, SoI Panel, Main Panel, Data Drawer, Commit Flow, Out-of-SoI Guard) promoted Approved with Playwright evidence.

**Integration gate criteria:**
- Playwright end-to-end demo: launch → pick user → pick SoI → edit Requirement → Commit → message visible in owner's mailbox (backend-level check).
- `make verify` green (including `frontend-e2e`).
- No stray Cytoscape or AG Grid usage (those libraries land in S07).
