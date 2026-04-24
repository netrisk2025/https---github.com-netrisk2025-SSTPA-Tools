# Repo Layout

This repository uses a mixed-language monorepo because the project has four different change rates:

- `backend/` for the transactional graph API and telemetry runtime
- `apps/desktop-shell/` for the Tauri host application
- `addons/` for user-facing tools that should remain independently developable
- `tools/reference-pipeline/` plus `reference-data/` for offline framework ingestion and normalization

The shell owns composition, navigation, and shared UX conventions. Add-on tools own their own local behavior and communicate through stable contracts from `packages/addon-sdk/`.

The reference pipeline is intentionally separate from the live backend so framework ingestion can be tested, versioned, and audited before any graph import occurs.
