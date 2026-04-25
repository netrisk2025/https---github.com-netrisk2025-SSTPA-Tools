# SRS V58 Migration Notes

SRS V58 is the source of truth for current development and is treated as intended SSTPA Tool SRS version 0.5.8.

## Owner Review Items

- V58 document header says version 0.5.7, while the migration baseline is intended version 0.5.8.
- `(:MasterRegime)` is required by the Asset Manager, but it is not listed in Section 1.3.3 canonical node labels.
- `(:FunctionalFlow)` is listed in Section 1.3.3 canonical node labels, but Section 1.3.8.1 does not assign it a type identifier.
- `(:DiagramView)` is referenced as a preferred future persistence pattern, but it is not fully canonicalized as a node label or relationship model.
- The reference model mixes language about cloning reference properties with the rule that GUI mutation of imported reference data is limited to creating or removing `[:REFERENCES]`.
- Some section references are stale after renumbering, including Requirement property headings and references to the older section layout.

## Implementation Assumptions

- `SSTPA Tool SRS V58.md` supersedes archived SRS files. Archived files are historical context only.
- The active consolidated SHALL extraction is `docs/verification/SSTPA_SHALL_Requirements.md`.
- `(:MasterRegime)` is handled as a tool-data/template node until the SRS canonical node-label list is updated.
- `(:FunctionalFlow)` uses the HID type identifier `FF` until Section 1.3.8.1 assigns an official identifier.
- `(:DiagramView)` is not added to the core graph model in this migration pass. Existing JSON properties remain the MVP persistence mechanism.
- New graph writes use canonical SRS 0.5.8 labels and relationship names. Legacy aliases are accepted only at import or migration boundaries.
