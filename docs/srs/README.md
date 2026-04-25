# SRS Workflow

`SSTPA Tool SRS V58.md` is the authoritative SRS for current development and is treated as intended version 0.5.8 even though its document header still says 0.5.7.

Prior SRS baselines are archived in `archive/` and may be used only as historical context. Do not use archived SRS files to override V58 behavior.

Use the following workflow:

1. Start from `SSTPA Tool SRS V58.md`.
2. Use `../verification/SSTPA_SHALL_Requirements.md` as the consolidated V58 SHALL reference when testing functionality.
3. Keep `../verification/shall-register.md` as the curated implementation/verification tracking register.
4. Map accepted SHALL statements to automated or manual verification in `../verification/verification-matrix.md`.
5. Document V58 editorial assumptions in `srs-v58-migration-notes.md` instead of blocking implementation on obvious non-behavioral issues.
