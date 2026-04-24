# Approved SHALL Register

| Req ID | Source Section | Requirement Summary | Status | Verification Intent | Notes |
| --- | --- | --- | --- | --- | --- |
| DRAFT-001 | 1.3.2 | Relationship names use `UPPERCASE_SNAKE_CASE` and singular labels | Approved | unit + contract | Implemented by backend graph label and relationship catalog tests |
| DRAFT-002 | 1.3.6 / 1.3.7 | Nodes carry identity and ownership metadata | Candidate | unit + integration | HID, uuid, and common metadata helpers are implemented; full mutation-layer enforcement remains open |
| DRAFT-003 | 2.2.10.8.1 | Non-owner changes generate mailbox notifications in the same transaction | Approved | integration | Implemented by backend mutation layer; REST exposure remains a later slice |
| SSTPA-COPY-001 | 5 | Source files include the required SSTPA Tool software copyright statement | Approved | automated check | Implemented by `make copyright-check` |
| SSTPA-VERIFY-001 | verification workflow | `make verify` is the standard automated verification gate after `make bootstrap` | Approved | automated command | Includes copyright, backend, reference pipeline, frontend lint/typecheck/test, and Docker Compose config |
| SSTPA-HID-001 | 1.3.6 | HID helpers format and parse `{TYPE}_{INDEX}_{SEQUENCE}` identifiers using SRS node type identifiers | Approved | unit | Implemented in `backend/internal/identity` |
| SSTPA-METADATA-001 | 1.3.7 / 1.3.7.1 | Common node metadata helper assigns identity, ownership, timestamps, schema version, and `"Null"` defaults | Approved | unit | Full graph-write enforcement remains part of mutation-layer work |
| SSTPA-SCHEMA-INDEX-001 | 1.3.6.2.1 | Backend can bootstrap SRS-required Neo4j indexes for HID, uuid, Name, and TypeName | Approved | unit + integration-ready | Bootstrap executes when `SSTPA_NEO4J_URI` is configured |
| SSTPA-VERIFY-NEO4J-001 | 2.2.5 / 2.2.10.8 | Backend has a reusable Testcontainers Neo4j fixture for graph integration tests | Approved | integration | Skips cleanly when Docker is unavailable |
| SSTPA-SBOM-001 | Compliance / future Help or Settings surface | Open source software integrated into SSTPA Tool is tracked in a maintained SBOM with name, version, and license | Approved | generated artifact + automated check | Current artifact is `docs/compliance/sbom.md`; future UI surfacing remains planned |
| SSTPA-MUTATION-001 | 2.2.10.8 | Backend mutations execute through a single transactional mutation layer | Approved | unit + integration | Current layer supports create node, update node, and create relationship operations |
| SSTPA-MUTATION-002 | 1.3.2.1 / 1.3.4.1 | Backend rejects duplicate logical relationships and prevents cycles for DAG-governed relationship types | Approved | unit + integration | Duplicate relationship covered by Neo4j integration test; DAG guard implemented with bounded traversal |
| SSTPA-MESSAGING-001 | 1.4.4 / 2.2.10.8.1 | Backend creates internal `CHANGE_NOTIFICATION` messages in owner mailboxes | Approved | unit + integration | Message model and mailbox append helper implemented in `backend/internal/messaging` |
| SSTPA-OWNERSHIP-001 | 1.3.7.1 / 2.2.10.8.2 | Ownership edits update Owner and OwnerEmail as a pair, preserve Creator fields for non-admin users, and enforce self-assumption for non-admin users | Approved | unit | Admin override support is present in the validator but broader role policy remains future security work |
| SSTPA-API-001 | 2.2.10.1 / 2.2.10.2 / 2.2.10.3 / 2.2.10.4 / 2.2.10.5 / 2.2.10.6 | Backend exposes JSON REST endpoints for node lookup, node type listing, hierarchy retrieval, graph search, relationship validation, and node context retrieval | Approved | unit + integration | Implemented under `/api/v1`; list/search endpoints enforce bounded pagination |
| SSTPA-API-MUTATION-001 | 2.2.10.8 / 2.2.10.8.1 | `POST /api/v1/mutations` executes graph writes through the transactional mutation layer and returns a `CommitReport` | Approved | integration | REST endpoint reuses `backend/internal/mutation` and preserves owner notification behavior |
| SSTPA-API-MESSAGE-001 | 2.2.10.10.1 / 3.4.8 | Backend exposes mailbox endpoints for message list, detail, create, reply, read, delete, and unread count | Approved | unit + integration | Current implementation uses internal mailbox records; future security work will refine authenticated access control |
| SSTPA-REF-API-001 | 1.5 / 2.2.10.10.2 / 2.2.10.10.3 / 2.2.10.10.4 / 2.2.10.10.5 | Backend exposes read-only reference framework retrieval/search endpoints and transactional `[:REFERENCES]` assignment endpoints | Approved | unit + integration | Imported reference item content remains read-only; assignment creation/removal only changes the relationship |
| SSTPA-OPENAPI-001 | 2.2.10 | Backend publishes an OpenAPI 3.1 contract for the implemented REST surface | Approved | static + unit | Contract lives in `docs/api/openapi.yaml` and is served from `/api/v1/openapi.yaml` |
| SSTPA-API-CLIENT-001 | 2.2.10 / public interfaces | TypeScript clients call the backend through a shared `@sstpa/api-client` package instead of ad hoc endpoint strings | Approved | unit | Package added under `packages/api-client` with Vitest coverage |

Status guidance:

- `Candidate` requirement still under review from the draft SRS
- `Approved` accepted for implementation and verification
- `Deferred` intentionally postponed
- `Superseded` replaced by a newer approved requirement
- `Rejected` invalid or removed from the active product scope
