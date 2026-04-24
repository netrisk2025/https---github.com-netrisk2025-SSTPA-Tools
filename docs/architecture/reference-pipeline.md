# Reference Pipeline

Reference frameworks must not be scraped live during normal product operation. They move through an offline pipeline:

1. `reference-data/raw/` stores the original source payloads delivered to the project.
2. `reference-data/staged/` stores extracted or decompressed payloads in a tool-friendly layout.
3. `reference-data/normalized/` stores a stable intermediate representation shared across frameworks.
   This should be treated as two logical passes:
   - pass 1: per-framework internal graph normalization
   - pass 2: cross-framework link resolution
4. Backend import code turns normalized records into a separate read-only graph structure inside Neo4j.

Initial target frameworks:

- MITRE ATT&CK
- MITRE EMB3D
- NIST SP 800-53

The intermediate format exists so parsing, normalization, validation, and graph import can be developed and verified independently.

The concrete graph model and normalized artifact shape are defined in
`docs/architecture/reference-data-graph-model.md`. That document is the
authoritative design for reference item keying, cross-framework links,
assignability, and normalized `items/relationships/citations` outputs.
