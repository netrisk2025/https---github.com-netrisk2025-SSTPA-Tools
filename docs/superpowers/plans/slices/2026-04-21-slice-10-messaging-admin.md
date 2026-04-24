# Slice 10 — Message Center, Reference Catalog, Admin & Onboarding (Phase B, Outline)

**Goal:** Ship Message Center as a real mailbox (SRS §3.4.8), the Reference Catalog Tool for external-framework assignment (SRS §3.4.3.5 / §1.5), and the Admin + User onboarding flow (SRS §1.4.1–1.4.3).

**Architecture:** Three work areas:
1. **Message Center** — promote `addons/message-center` from stub to full implementation; uses the `/api/v1/messages/*` routes from S04-T10 and owner-notification messages from S03.
2. **Reference Catalog Tool** — new addon `addons/reference-catalog`, pop-up, Assignment + Research modes per SRS §3.4.3.5. Reads via `/api/v1/reference/*` (S04-T11) and writes `[:REFERENCES]` via `/api/v1/reference/assignments` (S04-T12).
3. **Admin + Onboarding** — new tool surface + startup integration. `(:User)` / `(:Admin)` / `(:Users)` / `(:Admins)` node management; login screens (user, new-user, admin, new-admin) per SRS §1.4.1.

**Tech Stack:** Existing React/TanStack; AG Grid allowed for the mailbox list and reference search results (already approved in S07); no new libraries expected.

**Pre-reads:**
- SRS §1.4.1 Onboarding
- SRS §1.4.2 Admin Data
- SRS §1.4.3 User Data
- SRS §1.4.4 Messaging Data Model
- SRS §3.3.1 Message Center Add-on Tool
- SRS §3.4.3.4 / §3.4.3.5 Reference Catalog Tool (section reads as §3.4.3.5 in the draft; confirm heading during T01)
- SRS §3.4.8 Message Center
- S03 messaging model, S04 messaging + reference endpoints
- Conventions doc.

**Invariants:**
- Message Center is read-only to the node graph except for read/delete/reply operations against messages the user owns.
- Reference Catalog does not edit imported reference items; only creates/removes `[:REFERENCES]`.
- Onboarding does not bypass the authentication stand-in — no password auth is introduced here.
- Admin role cannot own data (SRS §1.4.2); Admin actions on others' nodes still generate change-notifications.
- Installer seeds a single initial `(:Admin)` and `(:User)` from Installer inputs (SRS §1.4.1) — implementation of the installer itself is S12; S10 only defines the shapes.

**Tasks:**

### S10-T01: Docs audit — reconcile Reference Catalog heading and §1.4 user model
Difficulty: low · Integration Checkpoint: yes · Files: `docs/architecture/user-model-reconciliation.md`.
**Agent Briefing:** Reconcile SRS inconsistencies flagged during S10 pre-read (section number drift; `(:Users)` vs `(:User)` parent container phrasing; Admin relationship to ownership rules). Produce proposals as SHALL candidates. Pause for user approval before T02.

### S10-T02: Backend `:Users` + `:Admins` parent containers + mutation support
Difficulty: medium · Files: `backend/internal/mutation/onboarding.go`, `onboarding_test.go`; modify `backend/internal/schema/bootstrap.go` for new labels.
**Agent Briefing:** Idempotent `EnsureUsersContainer(ctx, tx)` and `EnsureAdminsContainer(ctx, tx)`. Used by the Startup Software and Installer.

### S10-T03: Backend onboarding endpoints
Difficulty: medium · Files: `backend/internal/http/onboarding/**`; routes `GET /api/v1/users`, `POST /api/v1/users`, `POST /api/v1/admins` (admin-only).
**Agent Briefing:** List registered users for startup; create new user from onboarding screen; create new admin with admin-login stand-in. Contract-tested.

### S10-T04: Startup login flow
Difficulty: high · Integration Checkpoint: yes · Files: extend `apps/desktop-shell/src/startup/**` from S06-T04.
**Agent Briefing:** Four screens per SRS §1.4.1 Onboarding: User Login, New User, Admin Login, New Admin. On success, establishes Zustand currentUser with isAdmin flag. Playwright test covers all four paths.

### S10-T05: Admin tool surface
Difficulty: medium · Files: `addons/admin-tool/**`.
**Agent Briefing:** Pop-up tool (admin-only visibility on Control Panel). Lists users/admins, allows ownership reassignment (of ANY node), triggers owner-notification messages per §2.2.10.8.2.

### S10-T06: Message Center mailbox list view
Difficulty: medium · Files: rewrite `addons/message-center/src/**`.
**Agent Briefing:** AG Grid columns: Subject, DateTime, HID, Sender, Message Type, Read/Unread. Sort + reverse, row selection, keyboard navigation, unread highlighting. Reads from `/api/v1/messages`.

### S10-T07: Message Center detail view (Reply, Delete, Close)
Difficulty: medium · Files: extend message-center.
**Agent Briefing:** SRS §3.4.8.4 + 3.4.8.5 Direct Messaging; Reply uses `POST /api/v1/messages/{id}/reply`. Delete is soft (SRS §3.4.8.7).

### S10-T08: Message Center change-notification integration
Difficulty: medium · Integration Checkpoint: yes · Files: extend message-center; connect to S03 output format.
**Agent Briefing:** Automatic `CHANGE_NOTIFICATION` messages land correctly; detail view shows all affected HIDs; clicking an HID jumps to its Data Drawer (uses S06 drawer open API).

### S10-T09: Reference Catalog Tool — Research Mode
Difficulty: high · Files: `addons/reference-catalog/**`.
**Agent Briefing:** Pop-up with framework navigation, search, filter per SRS §3.4.3.5. Read-only inspector for selected reference item.

### S10-T10: Reference Catalog Tool — Assignment Mode
Difficulty: high · Integration Checkpoint: yes · Files: extend reference-catalog.
**Agent Briefing:** Auto-initializes to Assignment mode when Data Drawer open for a valid node type. Validates assignment via backend; on success returns the reference item to the calling Data Drawer. Only allowed node-type/item-type pairs are enabled (SRS §1.5.5).

### S10-T11: Reference Catalog — Data Drawer relationship group
Difficulty: medium · Files: extend `apps/desktop-shell/src/drawer/**`.
**Agent Briefing:** Data Drawer displays assigned `[:REFERENCES]` in a dedicated group with unlink/inspect actions; node-local external-reference properties editable per SRS §1.3.8 prelude text.

### S10-T12: Slice integration gate
Difficulty: low · Integration Checkpoint: yes.
**Acceptance:** Reqs for onboarding, admin tool, message center real, reference catalog (both modes), node-local external-reference properties → Approved.

**Integration gate criteria:**
- Playwright demos: new user onboarding; admin login; change owner → notification arrives; open Message Center; reply; assign an ATT&CK Technique to a `(:Hazard)` via Reference Catalog.
- `make verify` green.
