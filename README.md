# SSTPA Tool

SSTPA Tool is a desktop-first system security engineering environment built around a Neo4j-backed System of Interest (SoI) graph, a Go backend, and a Tauri/React shell.

This repository is scaffolded for two near-term needs:

- the SRS is still being reorganized and reduced to a cleaner approved-requirement set
- add-on tools and reference-data ingestion need stronger boundaries so they can evolve independently

## Repository Layout

- `apps/desktop-shell/` Tauri + React host application
- `addons/` independently developed add-on tool packages loaded by the shell
- `packages/` shared TypeScript contracts, domain types, and UI primitives
- `backend/` Go REST API and backend runtime
- `tools/reference-pipeline/` Go CLI for raw-to-normalized reference-data preparation
- `reference-data/` raw, staged, normalized, and manifest directories for imported frameworks
- `infra/docker/` local backend and telemetry stack
- `docs/` SRS copies, architecture notes, and verification planning
- `.opencode/` project-specific OpenCode agents and commands

## Working Model

- Treat `docs/srs/source/SSTPA Tool SRS V56.md` as the current upstream draft.
- Record approved, deferred, or superseded SHALL statements in `docs/verification/shall-register.md`.
- Map approved SHALL statements to verification assets in `docs/verification/verification-matrix.md`.
- Keep add-on tools thinly coupled to the shell through contracts in `packages/addon-sdk/`.
- Move external frameworks through `reference-data/raw -> reference-data/staged -> reference-data/normalized` before graph import.

## Key Commands

- `make bootstrap` install JavaScript workspace dependencies
- `make backend-run` run the Go API locally
- `make backend-test` run backend unit tests
- `make reference-run` exercise the reference-pipeline CLI
- `make reference-test` run reference-pipeline unit tests
- `make frontend-dev` run the desktop shell web workspace in Vite dev mode
- `make frontend-build` build the desktop shell web workspace
- `make frontend-lint` lint all TypeScript workspaces
- `make frontend-typecheck` typecheck all TypeScript workspaces
- `make frontend-test` run TypeScript workspace tests
- `make compose-config` validate the Docker Compose stack
- `make verify` run the current automated verification suite

## OpenCode Workflow

- Use `/verify-slice` to run targeted verification with the verification subagent.
- Use `/srs-audit` to review current code or plans against the draft SRS and approved SHALL set.
- Use `/traceability` to produce a requirement-to-code-to-test summary.
- Use `/reference-audit` when changing the framework ingestion pipeline or reference graph boundaries.
