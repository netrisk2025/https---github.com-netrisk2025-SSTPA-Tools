# Verification Workflow

This project intends to shrink the draft SRS down to an approved set of valid SHALL statements and then create verification for every approved SHALL.

Use the following artifacts together:

- `shall-register.md` records which SHALL statements are currently approved, deferred, superseded, or rejected.
- `verification-matrix.md` maps approved SHALL statements to automated tests, contract checks, integration tests, or manual verification.

Verification policy:

- No non-trivial implementation slice is complete without a stated verification path.
- Use focused unit, contract, or integration tests where they provide the best signal.
- Prefer subagent-driven verification for routine code reviews and test execution.
- If a requirement is not yet approved for implementation, capture that as a gap instead of guessing.
