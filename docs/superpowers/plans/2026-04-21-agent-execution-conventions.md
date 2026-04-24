# SSTPA Tools Agent Execution Conventions

> **Purpose:** Define the shared format that every slice plan in `slices/` assumes. An orchestrator must read this file once per session before dispatching tasks.

## 1. Task Anatomy

Every slice task follows this schema, whether published as a Phase-A detailed plan or a Phase-B outline.

```markdown
### Task <SID>-T<TID>: <short imperative subject>

**Task ID:** <SID>-T<TID>
**Depends On:** <comma-separated list of Task IDs, or `none`>
**Difficulty:** low | medium | high
**Integration Checkpoint:** yes | no
**Governing SRS sections:** <list of section numbers or `n/a - gap-fill`>
**Req IDs touched:** <shall-register IDs or `proposes: NEW-REQ-XYZ`>

**Files:**
- Create: `relative/path/to/new_file.ext`
- Modify: `relative/path/to/existing_file.ext[:start-end]`
- Test: `relative/path/to/test_file.ext`

**Agent Briefing:**
<3–10 sentences that a fresh coder agent with no prior context can execute from. State the goal, the invariants, any surprising constraints (e.g., "do not touch any file outside the Files list"), and the exact acceptance test the agent will self-verify with.>

**Steps:** (Phase A only — bite-sized 2–5 minute items with literal code)
- [ ] Step 1: Write failing test `<test name>`
- [ ] Step 2: Run `<command>`, expect FAIL with `<error>`
- [ ] Step 3: Implement `<function/type>`
- [ ] Step 4: Run `<command>`, expect PASS
- [ ] Step 5: Commit `<commit subject>`

**Acceptance tests:**
- `<shell command>` → `<expected output or pass condition>`
- `<shell command>` → `<expected output or pass condition>`

**Out-of-scope:**
<anything a naive agent might think belongs here but doesn't — e.g., "do not wire this into the REST layer; that is S04-T07">

**Evidence (filled after execution):**
- Commit: `<sha>`
- `make verify`: <pass timestamp>
- SHALL register update: <row ID + status change>
```

## 2. Commit Convention

```
<type>(<area>): <short imperative subject ≤ 72 chars>

<optional body — why, not what. Reference SRS sections.>

Refs: SRS §<sec1>, §<sec2>. Req: <ID or proposed ID>. Task: S<NN>-T<TT>.
Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
```

- `type` ∈ `feat | fix | test | chore | docs | refactor | build`.
- `area` ∈ `backend | shell | addon-<name> | reference | infra | docs | verify | identity | graph | api | ui | installer`.
- One task = one commit whenever possible. If a task legitimately requires multiple commits (e.g., test-then-impl), each commit cites the same Task ID and the final commit in the sequence closes the task.

## 3. Branch Convention

- `slice/<NN>-<short-name>` — one branch per slice, taken from `main` after prior slice is Accepted.
- Tasks land as sequential commits on the slice branch.
- Slice branch merges to `main` via PR after the Integration gate passes.

## 4. Difficulty Tiers (orchestrator routing)

| Tier | Typical work | Model guidance |
|------|--------------|----------------|
| `low` | Mechanical: add a struct field, add a route, add a Tailwind class, write a pure-function unit test. No cross-file coupling. | Haiku-class; minimal review. |
| `medium` | Wire two components: handler + service, component + store, serde struct + JSON Schema. Touches ≤ 3 files. Requires reading at least one existing file to get signatures right. | Sonnet-class; single-file spot review. |
| `high` | Integration: mutation transactions, cross-SoI validation, full-tool wiring, diagram rendering with layout. Touches ≥ 4 files or requires architectural judgment. | Opus-class; mandatory human or senior-model review. |

## 5. Agent Briefing Template (to paste into a coder subagent)

```
You are implementing Task <SID>-T<TID> of the SSTPA Tools project.
Working directory: /home/netrisk/Projects/sstpa-tool
Plan file: docs/superpowers/plans/slices/<slice-file>.md
Conventions: docs/superpowers/plans/2026-04-21-agent-execution-conventions.md

READ ONLY THESE PLAN SECTIONS:
- The task block for <SID>-T<TID>
- The "Pre-reads" list at the top of the slice plan
- The "Invariants" section of the slice plan

FILES YOU MAY CREATE OR MODIFY:
<paste Files: block>

YOUR DELIVERABLE:
<paste Agent Briefing>

WHEN YOU FINISH:
1. Run each acceptance-test command and paste the output into your response.
2. Create a single commit following the convention in §2 of the conventions file.
3. Report the commit SHA back.

DO NOT:
- Modify any file not listed above.
- Add dependencies not already in go.mod / package.json without an explicit instruction in your briefing.
- "Fix" code outside your scope that looks wrong — report it instead.
- Merge or push. Leave the branch as-is.
```

## 6. Pre-reads Section

Every slice plan opens with a **Pre-reads** list. An orchestrator dispatching a task MUST include these files in the coder agent's initial read set and forbid reads beyond that list + the task's own Files set. This bounds the agent's context and makes results reproducible.

## 7. Invariants Section

Every slice plan has an **Invariants** section that enumerates properties the slice must preserve. Common invariants:

- Singular node labels, `UPPERCASE_SNAKE_CASE` relationships (SRS §1.3.2).
- Never create reverse relationships without concrete justification.
- Missing property values are the literal string `"Null"`.
- All list-returning endpoints paginate.
- Neo4j is never exposed on the edge network.
- Imported reference nodes are read-only from the GUI.
- Every mutation is ACID; change-notifications ship in-transaction.
- Add-on tools live in pop-ups; the Data Drawer is the only general edit surface.

Every coder agent briefing restates the relevant invariants.

## 8. Gap-Fill Specs

When a slice owns a gap-fill (per Master §5), its first task authors a Gap-Fill Spec in `docs/architecture/gap-fill-<topic>.md`. The spec is:

1. Short (≤ 400 lines preferred).
2. Proposal-style: each requirement is written as a SHALL statement with a proposed Req ID.
3. Followed by a PAUSE for human approval before the slice's Task 2 begins.
4. Added to `docs/verification/shall-register.md` as Candidate entries on acceptance.

## 9. Verification Artifacts (post-task updates)

Whenever a task implements or verifies a requirement, it MUST (as part of its own commit or the slice's closing commit) update:

1. `docs/verification/shall-register.md` — set the Req row's `Status` to `Approved` and its `Verification Intent` to the specific unit/contract/integration path used.
2. `docs/verification/verification-matrix.md` — add or update the row with the test command, test file location, and Automated column.

An Integration gate FAILS if any task in the slice implemented a requirement without updating these files.

## 10. Evidence Capture

At the end of each task, the orchestrator fills the task's **Evidence** block in the slice plan file. This keeps the plan itself the single source of truth for what's been delivered, and makes the whole effort auditable later. An Integration gate also produces an "Evidence Summary" at the bottom of the slice plan:

```
## Evidence Summary
Slice: S<NN> — <title>
Branch: slice/<NN>-<name>
PR: <url or "n/a">
verify at slice close: <timestamp>
Tasks: <all T## → commit SHA>
SHALL register deltas: <list>
Open questions logged for next slice: <list>
```

## 11. When Blocked

If a coder agent cannot complete a task within its Files set without violating a rule:

1. The agent writes a `BLOCKED: <reason>` comment to the task's Evidence block and stops.
2. The orchestrator decides: (a) amend the task (stay in-slice), (b) open a new task (stay in-slice), (c) escalate to a human (stop the slice).
3. The orchestrator never silently expands scope.

## 12. Verify Cadence

- `make verify` runs as the final step of every task.
- `make verify` must be green before a task is marked `completed`.
- A slice is not Accepted until `make verify` is green on the slice branch after the final task commits.

## 13. SRS Section Citation

Every PR description, commit body, and plan completion note includes the governing SRS sections (e.g., `SRS §2.2.10.8.1`). Per `CLAUDE.md`, missing citation is a review-blocker.

## 14. Library-Addition Rule

Per `CLAUDE.md` / SRS §6.1 "minimum complexity": any task that introduces a new runtime dependency MUST do so in a dedicated task whose agent briefing explicitly says "introduces dependency X; user approved in <link or 'pending'>". The orchestrator pauses for user approval before executing such a task.

## 15. Naming Consistency Across Slices

Shared identifiers (Go types, TS types, route paths, Cypher procedure names) are declared in the slice that first introduces them and re-used verbatim by downstream slices. The orchestrator, before promoting a Phase-B outline to Phase A, scans the completed upstream slices and fills the outline's code blocks with actual names. This is the whole reason the two-phase model exists.

## 16. Pre-Commit Hook Behavior

If a pre-commit hook blocks a commit, the coder agent must NOT use `--no-verify`. It must:
1. Read the hook's failure message.
2. Fix the underlying issue.
3. Re-run `git add` and `git commit`.
4. Record the fix in the task's Evidence block.

## 17. Ending Notes

These conventions are live. If a slice finds a convention doesn't fit, the orchestrator updates this file AND notes the change in the slice Evidence Summary. Future slices pick up the new convention.
