# Slice 13 - V58 Asset Manager, System Creation, Loss/Goal Generation (Phase B, Outline)

**Goal:** Bring the canonical SRS V58 Asset Manager and System creation behavior online before downstream Loss Tool and Goal Keeper feature work. This slice owns Element-to-System creation defaults, copying eligible Requirements/Assets into the child SoI, PRIMARY/DERIVED Assets, Regime cloning from MasterRegime, inherited Derived Asset Criticality display, and idempotent Loss/Root Goal generation.

**Architecture:** Refactor existing backend schema/mutation code instead of replacing it. Add focused services under `backend/internal/systemcreation/**` and `backend/internal/assets/**`, then expose them through S04-style REST endpoints. The frontend uses the scaffolded `addons/asset-manager` package and the shared add-on SDK context contracts. Loss Tool and Goal Keeper consume outputs from this slice.

**Pre-reads:**
- `docs/srs/SSTPA Tool SRS V58.md` as intended version 0.5.8
- `docs/srs/srs-v58-migration-notes.md` for known editorial assumptions
- `docs/verification/SSTPA_SHALL_Requirements.md` requirements: `1.3.4.1-001`, `1.3.4.2-001`, `1.3.4.11-001`, `1.3.7-001`, `3.4.7.1-001`, `3.4.7.6-001`, `3.4.7.7-001`, `3.4.7.8-005`, `3.4.7.10-001`, `3.4.7.14-001`, `3.4.7.19-001`
- S02-S04 outputs for schema registry, mutation validation, and REST conventions
- S06 output for shell/Data Drawer integration
- Conventions doc.

**Invariants:**
- HID Index is the canonical SoI membership rule.
- Creating a System from an Element creates default Purpose, Environment, State, Security, and FunctionalFlow nodes.
- Copied Requirements and Assets receive recomputed HID and uuid values in the child SoI.
- All new writes use canonical V58 relationship names; legacy aliases are accepted only at import/migration boundaries.
- `DERIVED_FROM` target must be a PRIMARY Asset.
- `MasterRegime` is treated as a tool-data/template node pending owner review because V58 requires it operationally but omits it from canonical labels.
- Loss/Root Goal generation is idempotent by Asset, Criticality, Assurance, and Environment combination.

## Tasks

### S13-T01: Asset/System Creation detailed design
Difficulty: high - Integration Checkpoint: yes - Files: `docs/architecture/v58-asset-system-creation.md`.
**Agent Briefing:** Turn the V58 assumptions into an implementation design. Confirm exact backend service boundaries, Root Goal naming/defaults, MasterRegime storage, idempotency keys, and the test fixture shape. Cite every SHALL ID from this slice.

### S13-T02: Backend System-from-Element service
Difficulty: high - Files: `backend/internal/systemcreation/**`; tests under same package.
**Agent Briefing:** Implement creation of child System from Element with default Purpose, Environment, State, Security, and FunctionalFlow nodes. Recompute HID/uuid. Copy eligible Requirements and Assets from the parent Element, allocated Function, and allocated Interface into the child SoI. Enforce the Element parent DAG rule.

### S13-T03: Asset service and canonical Asset relationships
Difficulty: high - Files: `backend/internal/assets/**`; mutation registry extensions.
**Agent Briefing:** Implement PRIMARY/DERIVED Asset creation, `DERIVED_FROM` validation, Asset association to Element/Function/Interface/State/Environment, and `HAS_REGIME`/`HAS_GOAL` relationships. Use registry constants only.

### S13-T04: Regime template and clone workflow
Difficulty: medium - Files: extend `backend/internal/assets/**`.
**Agent Briefing:** Add MasterRegime support as template/tool-data. Clone selected MasterRegime into an Asset-specific Regime, generate new HID/uuid, preserve eligible properties, and link through `(:Asset)-[:HAS_REGIME]->(:Regime)`.

### S13-T05: Idempotent Loss and Root Goal generation
Difficulty: high - Integration Checkpoint: yes - Files: `backend/internal/assets/loss_goal_generation.go`, tests.
**Agent Briefing:** On Asset commit, ensure missing Loss and Root Goal pairs exist for each Asset/Criticality/Assurance/Environment combination. Never create duplicate logical Loss/Goal pairs. Generate Root Goals associated with the Asset/Loss case for S14.

### S13-T06: REST endpoints for Asset Manager and System creation
Difficulty: medium - Files: `backend/internal/http/assets/**`, router wiring.
**Agent Briefing:** Add endpoints for listing Asset table rows, creating/updating Assets, creating child System from Element, cloning Regime templates, validating staged Asset changes, and regenerating missing Loss/Goal pairs.

### S13-T07: Asset Manager table-first UI
Difficulty: high - Files: `addons/asset-manager/**`.
**Agent Briefing:** Build table-first UI with sorting/filtering/column visibility/multi-select/search. Use progressive disclosure for detailed edit panels. Display PRIMARY/DERIVED type, associated Regimes, associations, Loss count, Goal status, and validation status.

### S13-T08: Derived Asset and inherited Criticality UI
Difficulty: medium - Files: extend `addons/asset-manager/**`.
**Agent Briefing:** Support Derived Asset creation from one or more PRIMARY Assets. Display inherited Criticality distinctly from directly selected Criticality. Prevent commit when the DERIVED target rule is violated.

### S13-T09: Asset Manager integration gate
Difficulty: medium - Integration Checkpoint: yes.
**Acceptance:** `make verify` green; backend tests cover System-from-Element, Asset creation, DERIVED_FROM target validation, Regime cloning, idempotent Loss/Goal generation; Playwright covers the table-first Asset Manager workflow.

**Verification against SRS V58 SHALL IDs:** Update `docs/verification/shall-register.md` and `docs/verification/verification-matrix.md` for the S13 IDs listed in `docs/superpowers/plans/srs-v58-slice-verification.md`.
