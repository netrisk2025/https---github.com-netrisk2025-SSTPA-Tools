# Slice 05 — Reference Data Pipeline & Intermediate Format (Phase B, Outline + Gap-Fill)

**Goal:** Close the SRS gap the user explicitly called out: "intermediate data format or methodology for downloading, relating and installing reference data". Design the intermediate format, build the offline `raw → staged → normalized → imported` pipeline, and ingest NIST SP 800-53r5, MITRE ATT&CK, and MITRE EMB3D into the Neo4j read-only reference graph.

**Architecture:** Expand `tools/reference-pipeline/` with per-framework modules (`internal/nist`, `internal/attack`, `internal/emb3d`) each implementing a shared `Normalizer` interface. A new `cmd/refimport` CLI applies the normalized JSON to Neo4j via Cypher, using the same constraints/indexes bootstrapped in S02. Provenance manifests (`reference-data/manifests/<framework>/<version>.yaml`) are committed to git; payloads stay untracked. License compliance is enforced by the manifest schema — no redistributable payloads live in the repo.

**Tech Stack:** Go, `encoding/json`, `encoding/xml`, `stix2` library only if ATT&CK STIX ingestion needs it (propose in T01 RFC before adding), `neo4j-go-driver/v5`.

**Pre-reads:**
- SRS §1.3.10, §1.5.1–1.5.6 (Reference Data Model)
- SRS §2.2.10.10.1 Framework Import Requirements
- `tools/reference-pipeline/**`
- `reference-data/README.md` (from S01-T05)
- `backend/internal/schema/bootstrap.go` (S02 output)
- Conventions doc.

**Invariants:**
- Imported reference nodes are read-only to the GUI (only `[:REFERENCES]` mutations allowed there).
- `ExternalID` is unique within a framework version.
- No live internet access from any pipeline tool — all sources are file-system inputs.
- Import is idempotent; re-running does not create duplicates.
- Every imported item carries `FrameworkName, FrameworkVersion, ExternalID, ExternalType, Name, ShortDescription, LongDescription, SourceURI, Imported, LastUpdated, RawData`.
- License posture is documented per framework in its manifest.

**Tasks:**

### S05-T01 (gap-fill): Author `gap-fill-reference-intermediate-format.md`
Difficulty: medium · Integration Checkpoint: yes (user approval) · Files: `docs/architecture/gap-fill-reference-intermediate-format.md`
**Agent Briefing:** Specify, as proposal SHALLs with proposed Req IDs:
1. The intermediate format (likely a JSON array of `NormalizedReferenceItem` records, one file per framework+version, with a sibling `relationships.json` describing `[:HAS_CHILD]` and `[:RELATED_TO]` edges).
2. The `raw → staged → normalized` methodology per framework (what format the pipeline accepts, what transformations happen at each stage, where framework-specific labels like `:NistControl` are attached).
3. The manifest schema (framework name, version, sha256 of each raw artifact, license string, import timestamp, staged/normalized file list).
4. License posture for NIST 800-53r5 (public domain), MITRE ATT&CK (Apache 2.0), MITRE EMB3D (check — flag if unclear).
5. Idempotent import Cypher: `MERGE` on `(FrameworkName, FrameworkVersion, ExternalID)`.
6. Versioning: how older imported versions coexist with newer ones in the same graph.
Pause for user sign-off before T02.

### S05-T02: Extend `NormalizedReferenceItem` with relationship fields + migrate schema
Difficulty: low · Files: `tools/reference-pipeline/internal/manifest/*.go`, `schemas/normalized-reference-item.schema.json` · Deps: S05-T01.
**Agent Briefing:** Add fields per the approved gap-fill: `ParentExternalID *string`, `RelatedExternalIDs []string`, `License string`, `RawData any`. Update JSON Schema. Update unit tests.

### S05-T03: Provenance manifest schema + validator
Difficulty: medium · Files: `tools/reference-pipeline/internal/provenance/*.go`, `schemas/provenance-manifest.schema.json` · Deps: S05-T01.
**Agent Briefing:** YAML-based manifest with sha256 of every raw input and every produced normalized file. CLI subcommand `refstage manifest validate <path>`.

### S05-T04: NIST SP 800-53r5 normalizer
Difficulty: medium · Files: `tools/reference-pipeline/internal/nist/*.go` · Deps: S05-T02, S05-T03.
**Agent Briefing:** Consumes NIST's JSON distribution (stored in `reference-data/raw/nist-800-53r5/5.1.1/`). Emits NormalizedReferenceItems with `ExternalType = "Control"` or `"ControlEnhancement"`, establishes `ParentExternalID` for enhancements. Unit tests with a tiny fixture (one control + one enhancement).

### S05-T05: MITRE ATT&CK normalizer
Difficulty: high · Integration Checkpoint: yes · Files: `tools/reference-pipeline/internal/attack/*.go` · Deps: S05-T02, S05-T03.
**Agent Briefing:** STIX 2.1 input (Enterprise ATT&CK bundle). Emits Tactic, Technique, Sub-technique, Mitigation items with inter-item RelatedExternalIDs (uses `x_mitre_` properties). Verify a few known records round-trip.

### S05-T06: MITRE EMB3D normalizer
Difficulty: high · Files: `tools/reference-pipeline/internal/emb3d/*.go` · Deps: S05-T02, S05-T03.
**Agent Briefing:** EMB3D source format TBD by T01 gap-fill; if still unclear, BLOCK and escalate rather than guess.

### S05-T07: `refimport` CLI — idempotent import into Neo4j
Difficulty: high · Integration Checkpoint: yes · Files: `tools/reference-pipeline/cmd/refimport/main.go`, `internal/importer/*.go` · Deps: S05-T04..T06, S02 schema bootstrap.
**Agent Briefing:** CLI accepts `--framework <name> --version <ver> --dir <path>` and imports all normalized JSON + relationships. Integration test (Testcontainers): import the NIST fixture, re-import, assert only one copy of each control exists. Uses `MERGE (f:ReferenceFramework {Name:$n, Version:$v}) ... MERGE (i:ReferenceItem:NistControl {FrameworkName:$n, FrameworkVersion:$v, ExternalID:$eid}) SET i += $props MERGE (f)-[:HAS_ITEM]->(i)`.

### S05-T08: Backend read-only enforcement
Difficulty: medium · Files: Modify `backend/internal/mutation/validate.go` · Deps: S03 mutation layer.
**Agent Briefing:** Any mutation plan that would set properties on a node carrying `:ReferenceItem` or `:ReferenceFramework` — other than creating/removing `[:REFERENCES]` — is rejected. Integration test via Testcontainers.

### S05-T09: Add `make reference-import` target
Difficulty: low · Files: `Makefile` · Deps: S05-T07.
**Agent Briefing:** `reference-import FRAMEWORK=<name> VERSION=<ver>` wrapping the CLI.

### S05-T10: Update `reference-data/README.md` with runbook
Difficulty: low · Files: `reference-data/README.md`; docs pointer in top-level README.

### S05-T11: Slice integration gate
Difficulty: low · Integration Checkpoint: yes · Files: verification docs + Evidence Summary.
**Acceptance:** `SSTPA-REFDATA-LAYOUT-001` → Approved confirmed; `SSTPA-REFDATA-READONLY-001` → Approved; new Candidates from gap-fill promoted to Approved.

**Integration gate criteria:**
- Pipeline runs end-to-end for NIST fixture.
- Backend rejects read-only violations.
- Manifests validate.
- `make verify` green.
