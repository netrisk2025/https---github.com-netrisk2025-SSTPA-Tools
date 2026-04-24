# Reference Data Workspace

This directory stores framework data before it becomes graph data.

- `raw/` immutable source deliveries
- `staged/` extracted or re-packed source material ready for normalization
- `normalized/` framework-independent intermediate records
- `manifests/` version manifests, checksums, and import metadata

The goal is to make parsing, normalization, validation, and Neo4j import independently testable.

The target graph model for the normalized outputs is documented in
`docs/architecture/reference-data-graph-model.md`.
