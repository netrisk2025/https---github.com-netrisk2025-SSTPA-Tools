# Reference Pipeline

This module owns the offline ingest path from delivered framework payloads to a normalized intermediate representation.

Short-term goals:

- keep framework parsing independent from backend graph import
- make raw, staged, and normalized payload transitions explicit
- give the project a place to verify framework data before it touches Neo4j

The first scaffold exposes a small CLI and a validation model for normalized records.

The intended normalized graph artifact shape is documented in
`docs/architecture/reference-data-graph-model.md`.

Current staging commands:

- `go run ./cmd/refstage --print-layout`
- `go run ./cmd/refstage manifest validate <path>`
- `go run ./cmd/refstage stage nist --catalog <path> --license <path> --out-dir <dir> --manifest <path>`

The first implemented extractor stages the NIST SP 800-53 OSCAL catalog into:

- `metadata.json`
- `collections.ndjson`
- `items.ndjson`
- `edge-candidates.ndjson`
- `citations.ndjson`

and writes a provenance manifest for the staged release.
