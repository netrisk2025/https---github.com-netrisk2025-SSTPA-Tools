# Slice 01 — Scaffold Hardening & Verify Pipeline (Phase A, Detailed)

> **For agentic workers:** REQUIRED SUB-SKILL: `superpowers:subagent-driven-development` or `superpowers:executing-plans`. Steps use `- [ ]` checkboxes. One task = one commit. Read `docs/superpowers/plans/2026-04-21-agent-execution-conventions.md` before starting.

**Goal:** Make `make verify` a real acceptance gate. Add CI, license headers, a reference-data `.gitignore` skeleton, a Testcontainers baseline for backend integration tests, and documentation that every later slice depends on. No product features are added.

**Architecture:** Pure tooling / scaffolding work. No new runtime libraries. Only developer-dependency additions.

**Tech Stack:** GitHub Actions, Go `testing`, Testcontainers-Go, existing Vitest + TypeScript toolchain, existing Docker Compose stack.

**Pre-reads (mandatory before any task):**
- `Makefile`
- `package.json`
- `go.work`
- `backend/Dockerfile`
- `infra/docker/compose.yaml`
- `docs/verification/README.md`
- `docs/verification/shall-register.md`
- `docs/verification/verification-matrix.md`
- `CLAUDE.md` (project rules)
- `docs/superpowers/plans/2026-04-21-agent-execution-conventions.md`

**Invariants:**
- No new runtime dependencies in this slice. Dev-dependencies only.
- Do not touch `apps/`, `addons/`, `packages/` source code beyond adding license banners.
- Keep commits small; one task = one commit.
- Do not disable or skip existing tests.

---

### Task S01-T01: Add copyright header generator + apply to all source files

**Task ID:** S01-T01
**Depends On:** none
**Difficulty:** low
**Integration Checkpoint:** no
**Governing SRS sections:** §5 SSTPA Tool Component Copyright
**Req IDs touched:** proposes `SSTPA-COPY-001`

**Files:**
- Create: `tools/devtools/copyright/cmd/apply/main.go`
- Create: `tools/devtools/copyright/go.mod`
- Create: `tools/devtools/copyright/internal/banner/banner.go`
- Create: `tools/devtools/copyright/internal/banner/banner_test.go`
- Modify: `go.work` — add `./tools/devtools/copyright`
- Modify: `Makefile` — add `copyright-apply` and `copyright-check` targets
- Modify: every `.go`, `.ts`, `.tsx` file under `backend/`, `apps/`, `addons/`, `packages/`, `tools/` (to prepend banner)

**Agent Briefing:**
SRS §5 requires every SSTPA Tool source file to carry a copyright banner. Build a Go CLI that idempotently prepends the banner to `.go`, `.ts`, `.tsx` files under the listed top-level directories. The banner text is the literal block from SRS §5 (reproduced in Step 3 below). Use `//` for Go and TS. The CLI accepts `--check` (exits nonzero if any file is missing the banner) and `--apply` (idempotently inserts it after the first-line shebang if present, else at position 0). Write unit tests for `banner.HasBanner(content string) bool` and `banner.Prepend(content string) string`. Then run `copyright-apply` over the repo. Do not add the banner to `node_modules/`, `dist/`, `.git/`, `reference-data/`, or generated files.

**Step 1: Write failing tests for banner package**

Create `tools/devtools/copyright/internal/banner/banner_test.go`:
```go
package banner

import "testing"

func TestHasBannerDetectsExistingBanner(t *testing.T) {
    src := Prepend("package foo\n")
    if !HasBanner(src) {
        t.Fatal("expected HasBanner to be true after Prepend")
    }
}

func TestHasBannerDetectsMissingBanner(t *testing.T) {
    if HasBanner("package foo\n") {
        t.Fatal("expected HasBanner to be false on bare source")
    }
}

func TestPrependIsIdempotent(t *testing.T) {
    once := Prepend("package foo\n")
    twice := Prepend(once)
    if once != twice {
        t.Fatalf("expected Prepend to be idempotent\nonce=%q\ntwice=%q", once, twice)
    }
}

func TestPrependPreservesShebang(t *testing.T) {
    src := "#!/usr/bin/env bash\necho hi\n"
    got := Prepend(src)
    if got[:len("#!/usr/bin/env bash")] != "#!/usr/bin/env bash" {
        t.Fatal("shebang must remain first line")
    }
    if !HasBanner(got) {
        t.Fatal("banner must be present after shebang")
    }
}
```

**Step 2: Run tests — expect compile failure**

Run:
```
cd tools/devtools/copyright && go test ./internal/banner/...
```
Expected: `no Go files in ... /banner` or `undefined: Prepend` / `undefined: HasBanner`.

**Step 3: Implement banner package**

Create `tools/devtools/copyright/internal/banner/banner.go`:
```go
// Package banner inserts the SSTPA Tools copyright header into source files.
package banner

import "strings"

const Marker = "SSTPA Tools proprietary intellectual property"

// Text is the literal banner mandated by SRS §5.
const Text = `// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source
// code are proprietary intellectual property of Nicholas Triska.
// Unauthorized reproduction, modification, or distribution is strictly
// prohibited. Licensed copies may be used under specific contractual terms
// provided by the author.
`

func HasBanner(content string) bool {
    return strings.Contains(content, Marker)
}

func Prepend(content string) string {
    if HasBanner(content) {
        return content
    }
    if strings.HasPrefix(content, "#!") {
        idx := strings.Index(content, "\n")
        if idx < 0 {
            return content + "\n" + Text
        }
        return content[:idx+1] + Text + content[idx+1:]
    }
    return Text + content
}
```

**Step 4: Run tests — expect PASS**

```
cd tools/devtools/copyright && go test ./internal/banner/...
```
Expected: all 4 subtests pass.

**Step 5: Implement CLI**

Create `tools/devtools/copyright/go.mod`:
```
module sstpa-tool/devtools/copyright

go 1.24
```

Create `tools/devtools/copyright/cmd/apply/main.go`:
```go
package main

import (
    "flag"
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "strings"

    "sstpa-tool/devtools/copyright/internal/banner"
)

var roots = []string{"backend", "apps", "addons", "packages", "tools"}
var skipDirs = map[string]struct{}{
    "node_modules": {}, "dist": {}, ".git": {}, "target": {},
    "reference-data": {}, "coverage": {}, ".vite": {},
}
var exts = map[string]struct{}{".go": {}, ".ts": {}, ".tsx": {}}

func main() {
    check := flag.Bool("check", false, "exit nonzero if any file lacks the banner")
    apply := flag.Bool("apply", false, "prepend the banner to files that lack it")
    flag.Parse()
    if *check == *apply {
        fmt.Fprintln(os.Stderr, "exactly one of --check or --apply is required")
        os.Exit(2)
    }

    var missing []string
    for _, root := range roots {
        _ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
            if err != nil {
                return nil
            }
            if d.IsDir() {
                if _, skip := skipDirs[d.Name()]; skip {
                    return filepath.SkipDir
                }
                return nil
            }
            if _, ok := exts[filepath.Ext(path)]; !ok {
                return nil
            }
            if strings.Contains(path, "node_modules") {
                return nil
            }
            content, err := os.ReadFile(path)
            if err != nil {
                return nil
            }
            s := string(content)
            if banner.HasBanner(s) {
                return nil
            }
            if *check {
                missing = append(missing, path)
                return nil
            }
            updated := banner.Prepend(s)
            return os.WriteFile(path, []byte(updated), 0o644)
        })
    }
    if *check && len(missing) > 0 {
        fmt.Fprintln(os.Stderr, "files missing copyright banner:")
        for _, p := range missing {
            fmt.Fprintln(os.Stderr, "  ", p)
        }
        os.Exit(1)
    }
}
```

**Step 6: Wire into go.work and Makefile**

Edit `go.work` — under `use (` add `    ./tools/devtools/copyright` keeping alphabetical order.

Edit `Makefile` — add before `verify:`:
```
copyright-check:
	cd tools/devtools/copyright && go run ./cmd/apply --check

copyright-apply:
	cd tools/devtools/copyright && go run ./cmd/apply --apply
```

And update `verify:` to depend on `copyright-check`:
```
verify: copyright-check backend-test reference-test frontend-typecheck frontend-test
```

Add `copyright-check copyright-apply` to the `.PHONY:` list.

**Step 7: Apply banner across the repo**

Run:
```
make copyright-apply
make copyright-check
```
Expected: `copyright-apply` modifies existing `.go/.ts/.tsx` files. `copyright-check` exits 0.

**Step 8: Verify**

Run:
```
make backend-test
make reference-test
make frontend-typecheck
make frontend-test
```
Expected: all pass (banner insertion does not break any existing test).

**Step 9: Commit**

```
git add tools/devtools/copyright Makefile go.work backend apps addons packages tools
git commit -m "$(cat <<'EOF'
chore(docs): add copyright banner per SRS §5

Introduces a dev-only CLI (tools/devtools/copyright) that idempotently
applies the SSTPA Tools copyright banner to .go/.ts/.tsx sources. Wires
make copyright-check into make verify so drift fails CI.

Refs: SRS §5. Proposes: SSTPA-COPY-001. Task: S01-T01.
Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

**Acceptance tests:**
- `make copyright-check` exits 0.
- `make verify` exits 0.
- `git grep -L 'SSTPA Tools proprietary intellectual property' -- 'backend/**/*.go' 'apps/**/*.ts' 'apps/**/*.tsx' 'addons/**/*.ts' 'packages/**/*.ts' 'tools/**/*.go'` outputs nothing.

**Out-of-scope:**
- Do not add a banner to `README.md`, `CLAUDE.md`, or any `.md` / `.yaml` / `.toml` file.
- Do not modify the reference-data tree.

**Evidence:** (fill after completion)
- Commit SHA: _________
- `make verify` timestamp: _________
- SHALL register update: add `SSTPA-COPY-001` as Candidate.

---

### Task S01-T02: Add Testcontainers-Go Neo4j fixture for future integration tests

**Task ID:** S01-T02
**Depends On:** S01-T01
**Difficulty:** medium
**Integration Checkpoint:** no
**Governing SRS sections:** §2.2.5 Backend Database, §2.2.10.8 Transaction Requirements
**Req IDs touched:** proposes `SSTPA-VERIFY-NEO4J-001`

**Files:**
- Create: `backend/internal/testhelpers/neo4j.go`
- Create: `backend/internal/testhelpers/neo4j_test.go`
- Modify: `backend/go.mod` (new dev dep: `github.com/testcontainers/testcontainers-go`, `github.com/testcontainers/testcontainers-go/modules/neo4j`, `github.com/neo4j/neo4j-go-driver/v5`)

**Agent Briefing:**
Provide a reusable Go test helper that spins up a Neo4j Community container, returns a configured driver, and cleans up on test end. Later slices depend on this. The helper MUST skip (not fail) when Docker is unavailable, using `t.Skip("docker unavailable: " + err.Error())`, so `make backend-test` remains runnable in environments without Docker. Write a smoke test that starts the container, runs `RETURN 1 AS n`, and asserts `n == 1`. This is the only task in S01 that adds runtime-adjacent libraries — they are test-only and pull into the backend module. Ask for user approval if `CLAUDE.md` rules require it; the orchestrator has a standing approval for test infrastructure additions in S01.

**Step 1: Write the helper contract test**

Create `backend/internal/testhelpers/neo4j_test.go`:
```go
package testhelpers

import (
    "context"
    "testing"
)

func TestNeo4jContainerSmoke(t *testing.T) {
    ctx := context.Background()
    fixture, err := StartNeo4j(ctx, t)
    if err != nil {
        t.Skipf("docker unavailable: %v", err)
    }
    defer fixture.Close(ctx)

    session := fixture.Driver.NewSession(ctx, neo4jSessionConfig())
    defer session.Close(ctx)

    result, err := session.Run(ctx, "RETURN 1 AS n", nil)
    if err != nil {
        t.Fatalf("cypher RETURN 1 failed: %v", err)
    }
    rec, err := result.Single(ctx)
    if err != nil {
        t.Fatalf("expected single record: %v", err)
    }
    got, ok := rec.Get("n")
    if !ok || got.(int64) != 1 {
        t.Fatalf("expected n=1, got %#v ok=%v", got, ok)
    }
}
```

**Step 2: Run test — expect compile failure**

```
cd backend && go test ./internal/testhelpers/...
```
Expected: `undefined: StartNeo4j` and `undefined: neo4jSessionConfig`.

**Step 3: Add dependencies**

```
cd backend && go get github.com/testcontainers/testcontainers-go@latest github.com/testcontainers/testcontainers-go/modules/neo4j@latest github.com/neo4j/neo4j-go-driver/v5@latest
go mod tidy
```

**Step 4: Implement helper**

Create `backend/internal/testhelpers/neo4j.go`:
```go
// Package testhelpers provides shared fixtures for backend tests.
package testhelpers

import (
    "context"
    "testing"

    "github.com/neo4j/neo4j-go-driver/v5/neo4j"
    "github.com/testcontainers/testcontainers-go"
    tcneo4j "github.com/testcontainers/testcontainers-go/modules/neo4j"
)

type Neo4jFixture struct {
    Driver    neo4j.DriverWithContext
    Container testcontainers.Container
}

func StartNeo4j(ctx context.Context, t *testing.T) (*Neo4jFixture, error) {
    t.Helper()
    container, err := tcneo4j.Run(ctx,
        "neo4j:5.26-community",
        tcneo4j.WithAdminPassword("test-password"),
    )
    if err != nil {
        return nil, err
    }
    uri, err := container.BoltUrl(ctx)
    if err != nil {
        _ = container.Terminate(ctx)
        return nil, err
    }
    driver, err := neo4j.NewDriverWithContext(uri,
        neo4j.BasicAuth("neo4j", "test-password", ""))
    if err != nil {
        _ = container.Terminate(ctx)
        return nil, err
    }
    return &Neo4jFixture{Driver: driver, Container: container}, nil
}

func (f *Neo4jFixture) Close(ctx context.Context) {
    _ = f.Driver.Close(ctx)
    _ = f.Container.Terminate(ctx)
}

func neo4jSessionConfig() neo4j.SessionConfig {
    return neo4j.SessionConfig{DatabaseName: "neo4j"}
}
```

**Step 5: Run test — expect PASS or SKIP**

```
cd backend && go test ./internal/testhelpers/... -v
```
Expected: PASS on hosts with Docker; SKIP (not FAIL) otherwise.

**Step 6: Confirm make verify still passes**

```
make verify
```
Expected: PASS.

**Step 7: Commit**

```
git add backend/internal/testhelpers backend/go.mod backend/go.sum
git commit -m "$(cat <<'EOF'
test(backend): add Testcontainers Neo4j fixture

Provides StartNeo4j() fixture that spins up a Neo4j 5.26-community
container for integration tests. Skips cleanly when Docker is unavailable
so make backend-test stays runnable in CI-lite environments.

Refs: SRS §2.2.5, §2.2.10.8. Proposes: SSTPA-VERIFY-NEO4J-001.
Task: S01-T02.
Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

**Acceptance tests:**
- `cd backend && go test ./internal/testhelpers/... -v` passes on a host with Docker.
- `cd backend && go test ./internal/testhelpers/... -v` skips cleanly without failure on a host without Docker.
- `make verify` passes.

**Out-of-scope:**
- Do not introduce any Cypher beyond `RETURN 1`. Schema belongs to S02.
- Do not wire Neo4j into the running backend (`cmd/api`). Also S02.

**Evidence:** _________

---

### Task S01-T03: Add GitHub Actions CI running `make verify`

**Task ID:** S01-T03
**Depends On:** S01-T02
**Difficulty:** low
**Integration Checkpoint:** no
**Governing SRS sections:** §6.1 minimum complexity
**Req IDs touched:** proposes `SSTPA-CI-001`

**Files:**
- Create: `.github/workflows/verify.yml`

**Agent Briefing:**
Add a GitHub Actions workflow that runs `make verify` on `ubuntu-latest` for every push and pull request. Pin action versions. Use `actions/setup-go@v5` with `go-version: 1.24`, `actions/setup-node@v4` with `node-version: 22`, and `docker/setup-buildx-action@v3`. The workflow MUST call `make bootstrap` before `make verify`. Do not split the verify target; run it whole.

**Step 1: Create workflow**

Create `.github/workflows/verify.yml`:
```yaml
name: verify
on:
  push:
    branches: [main]
  pull_request:
  workflow_dispatch:

jobs:
  verify:
    name: make verify
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
          cache: true
          cache-dependency-path: |
            backend/go.sum
            tools/reference-pipeline/go.sum
            tools/devtools/copyright/go.sum
      - uses: actions/setup-node@v4
        with:
          node-version: '22'
          cache: 'npm'
      - uses: docker/setup-buildx-action@v3
      - name: bootstrap
        run: make bootstrap
      - name: verify
        run: make verify
```

**Step 2: Local dry-run**

Run:
```
make bootstrap
make verify
```
Expected: PASS (asserts the workflow would succeed).

**Step 3: Commit**

```
git add .github/workflows/verify.yml
git commit -m "$(cat <<'EOF'
chore(build): add GitHub Actions verify workflow

Runs make bootstrap + make verify on ubuntu-latest for every push and PR
so scaffold hardening gates all downstream slices.

Refs: SRS §6.1. Proposes: SSTPA-CI-001. Task: S01-T03.
Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

**Acceptance tests:**
- `make verify` passes locally.
- Workflow file parses: `yq '.jobs.verify.runs-on' .github/workflows/verify.yml` returns `ubuntu-latest` (or a visual check if `yq` unavailable).

**Out-of-scope:**
- No release, lint, or deploy jobs here. Those land in S12.

**Evidence:** _________

---

### Task S01-T04: Tighten `make verify` to include lint and compose-config

**Task ID:** S01-T04
**Depends On:** S01-T03
**Difficulty:** low
**Integration Checkpoint:** no
**Governing SRS sections:** §6.1 minimum complexity, §2.3 Docker networks
**Req IDs touched:** `SSTPA-CI-001`

**Files:**
- Modify: `Makefile`

**Agent Briefing:**
The current `make verify` target runs `backend-test reference-test frontend-typecheck frontend-test`. It misses lint and compose validation. Add `copyright-check`, `frontend-lint`, and `compose-config` to the verify target. Ensure the compose-config step works in CI by checking only configuration (not launching services), which the existing target already does.

**Step 1: Edit Makefile**

Modify the `verify:` target to read:
```
verify: copyright-check frontend-lint backend-test reference-test frontend-typecheck frontend-test compose-config
```

**Step 2: Run full verify**

```
make verify
```
Expected: every sub-step passes.

**Step 3: Commit**

```
git add Makefile
git commit -m "$(cat <<'EOF'
chore(build): widen make verify to include lint and compose-config

verify now runs: copyright-check, frontend-lint, backend-test,
reference-test, frontend-typecheck, frontend-test, compose-config.
Ensures downstream slices detect lint drift and compose breakage early.

Refs: SRS §6.1, §2.3. Req: SSTPA-CI-001. Task: S01-T04.
Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

**Acceptance tests:**
- `grep '^verify:' Makefile` shows all seven sub-targets.
- `make verify` is green.

**Evidence:** _________

---

### Task S01-T05: Add reference-data tree README + .gitkeep stubs

**Task ID:** S01-T05
**Depends On:** S01-T04
**Difficulty:** low
**Integration Checkpoint:** no
**Governing SRS sections:** §1.3.10 External Reference Framework Data Model
**Req IDs touched:** proposes `SSTPA-REFDATA-LAYOUT-001`

**Files:**
- Create: `reference-data/README.md`
- Create: `reference-data/raw/.gitkeep`
- Create: `reference-data/staged/.gitkeep`
- Create: `reference-data/normalized/.gitkeep`
- Create: `reference-data/manifests/.gitkeep`

**Agent Briefing:**
`.gitignore` already excludes reference-data payloads but the directory structure is implicit. Add a README that names the four stages (raw / staged / normalized / manifests), describes what each holds, and references the forthcoming Gap-Fill Spec authored in S05. Add `.gitkeep` files so the directories exist on fresh clones.

**Step 1: Create README**

Create `reference-data/README.md`:
```markdown
# Reference Data Pipeline Workspace

External reference frameworks (NIST SP 800-53r5, MITRE ATT&CK, MITRE
EMB3D) flow through four stages before reaching the Neo4j graph as
read-only `(:ReferenceFramework)` + `(:ReferenceItem)` nodes.

| Stage | Purpose | Committed? |
| --- | --- | --- |
| `raw/` | Untouched vendor artifacts (XML, JSON, STIX, CSV). | No (gitignored). |
| `staged/` | First-pass extraction to a normalized intermediate format (see Gap-Fill Spec S05-T01). | No. |
| `normalized/` | Final `NormalizedReferenceItem` JSON ready for backend import. | No. |
| `manifests/` | Provenance manifests describing what was ingested and when. | Yes, one per framework+version. |

Each framework version lives under `reference-data/<stage>/<framework>/<version>/`.

See `docs/architecture/gap-fill-reference-intermediate-format.md`
(authored in slice S05) for the authoritative schema.

SRS reference: §1.3.10, §1.5.
```

**Step 2: Create .gitkeep files**

Each directory gets an empty `.gitkeep`. Adjust `.gitignore` if needed so these `.gitkeep` files are tracked (existing rules already whitelist `!reference-data/raw/.gitignore` patterns; replicate for `.gitkeep`). Update `.gitignore`:
```
reference-data/raw/*
!reference-data/raw/.gitkeep
reference-data/staged/*
!reference-data/staged/.gitkeep
reference-data/normalized/*
!reference-data/normalized/.gitkeep
```
(Preserve existing `.gitignore` exceptions.)

**Step 3: Verify untracked state**

```
git status reference-data
```
Expected: README.md and four .gitkeep files tracked; no payload files.

**Step 4: Commit**

```
git add reference-data .gitignore
git commit -m "$(cat <<'EOF'
docs(reference): add reference-data pipeline README and directory stubs

Documents raw/staged/normalized/manifests staging per SRS §1.3.10 and
§1.5 and pins empty directories via .gitkeep so clones have the layout.

Refs: SRS §1.3.10, §1.5. Proposes: SSTPA-REFDATA-LAYOUT-001.
Task: S01-T05.
Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

**Acceptance tests:**
- `ls reference-data` shows README.md plus raw/ staged/ normalized/ manifests/ each containing `.gitkeep`.
- `make verify` is green.

**Evidence:** _________

---

### Task S01-T06: Expand SHALL register with Slice-01 candidates + bootstrap Core Data Model IDs

**Task ID:** S01-T06
**Depends On:** S01-T05
**Difficulty:** low
**Integration Checkpoint:** yes (integration gate for S01)
**Governing SRS sections:** §1 through §6 (structural)
**Req IDs touched:** `SSTPA-COPY-001`, `SSTPA-CI-001`, `SSTPA-VERIFY-NEO4J-001`, `SSTPA-REFDATA-LAYOUT-001`

**Files:**
- Modify: `docs/verification/shall-register.md`
- Modify: `docs/verification/verification-matrix.md`

**Agent Briefing:**
Add one row for each of this slice's four new Candidate requirements. Also add placeholder Candidate rows for the Core Data Model SHALLs that S02 will implement (to avoid re-editing in S02): graph rule conformance, HID format, identity metadata fields, ownership mutation transaction rules. Do not mark any row Approved yet — acceptance happens when the implementing slice closes.

**Step 1: Update shall-register**

Append to `docs/verification/shall-register.md` (keep header intact):
```
| SSTPA-COPY-001 | §5 | Every source file carries the mandated copyright banner | Candidate | tooling check | verified by make copyright-check |
| SSTPA-CI-001 | §6.1 | make verify runs on every push and PR | Candidate | tooling check | GH Actions verify.yml |
| SSTPA-VERIFY-NEO4J-001 | §2.2.5 | Backend integration tests run against a real Neo4j 5.26-community instance | Candidate | integration | Testcontainers-Go fixture |
| SSTPA-REFDATA-LAYOUT-001 | §1.3.10, §1.5 | Reference data moves through raw/staged/normalized/manifests stages | Candidate | docs + tooling | enforced by pipeline CLI |
| SSTPA-GRAPH-RULES-001 | §1.3.2 | Node labels singular; relationships UPPERCASE_SNAKE_CASE; no gratuitous reverse edges | Candidate | unit + contract | implemented in S02 |
| SSTPA-HID-001 | §1.3.6 | HID format is {TYPE}_{INDEX}_{SEQUENCE} | Candidate | unit | implemented in S02 |
| SSTPA-IDENTITY-001 | §1.3.7 | Every node carries HID, uuid, TypeName, Owner, OwnerEmail, Creator, CreatorEmail, Created, LastTouch, VersionID | Candidate | integration | implemented in S02 |
| SSTPA-OWNERSHIP-TX-001 | §2.2.10.8.1 | Non-owner changes generate change-notification messages in the same ACID transaction | Candidate | integration | implemented in S03 |
| SSTPA-PAGINATION-001 | §1.3.2.1 | All list-returning endpoints support pagination and maximum result limits | Candidate | contract | implemented in S04 |
| SSTPA-REFDATA-READONLY-001 | §1.5.6 | Imported reference framework nodes are read-only from the GUI; only [:REFERENCES] mutations allowed | Candidate | integration | implemented in S05 |
```

**Step 2: Update verification-matrix**

Append matching rows to `docs/verification/verification-matrix.md`:
```
| SSTPA-COPY-001 | Source files carry SRS §5 banner | tooling | Yes | make copyright-check | tools/devtools/copyright | Added in S01-T01 |
| SSTPA-CI-001 | make verify runs in CI | tooling | Yes | GitHub Actions verify.yml | .github/workflows/verify.yml | Added in S01-T03 |
| SSTPA-VERIFY-NEO4J-001 | Neo4j integration fixture | integration | Yes | cd backend && go test ./internal/testhelpers/... | backend/internal/testhelpers | Added in S01-T02 |
| SSTPA-REFDATA-LAYOUT-001 | Reference-data staging layout | docs + tooling | Planned | `ls reference-data` | reference-data/README.md | Enforcement in S05 |
| SSTPA-GRAPH-RULES-001 | Graph naming rules | unit + contract | Planned | cd backend && go test ./internal/graph/... | backend/internal/graph | Implementation in S02 |
| SSTPA-HID-001 | HID format | unit | Planned | cd backend && go test ./internal/identity/... | backend/internal/identity | Implementation in S02 |
| SSTPA-IDENTITY-001 | Identity + ownership metadata on every node | integration | Planned | cd backend && go test ./internal/graph/nodetest/... | backend/internal/graph/nodetest | Implementation in S02 |
| SSTPA-OWNERSHIP-TX-001 | Same-transaction change notifications | integration | Planned | cd backend && go test ./internal/mutation/... | backend/internal/mutation | Implementation in S03 |
| SSTPA-PAGINATION-001 | Pagination on list endpoints | contract | Planned | cd backend && go test ./internal/http/... | backend/internal/http | Implementation in S04 |
| SSTPA-REFDATA-READONLY-001 | Reference graph read-only | integration | Planned | cd backend && go test ./internal/reference/... | backend/internal/reference | Implementation in S05 |
```

For the four `SSTPA-*-001` rows from this slice (COPY, CI, VERIFY-NEO4J, REFDATA-LAYOUT), also update the matching shall-register row's `Status` from `Candidate` to `Approved` since they have shipped.

**Step 3: Commit**

```
git add docs/verification
git commit -m "$(cat <<'EOF'
docs(verify): seed shall-register and verification-matrix for S01 + downstream

Registers the four S01 requirements as Approved and seeds Candidate rows
for the Core Data Model, HID, identity, ownership-transaction,
pagination, and reference-read-only requirements that downstream slices
will implement. Keeps the matrix as the single traceability surface.

Refs: SRS §5, §1.3.2, §1.3.6, §1.3.7, §1.3.10, §1.5.6, §2.2.10.8.1.
Reqs: SSTPA-COPY-001, SSTPA-CI-001, SSTPA-VERIFY-NEO4J-001,
SSTPA-REFDATA-LAYOUT-001 (Approved); SSTPA-GRAPH-RULES-001,
SSTPA-HID-001, SSTPA-IDENTITY-001, SSTPA-OWNERSHIP-TX-001,
SSTPA-PAGINATION-001, SSTPA-REFDATA-READONLY-001 (Candidate).
Task: S01-T06.
Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

**Acceptance tests:**
- `grep -c 'SSTPA-' docs/verification/shall-register.md` >= 10.
- `grep -c 'SSTPA-' docs/verification/verification-matrix.md` >= 10.
- `make verify` is green.

**Integration gate:**
- Confirm all six S01 tasks have Evidence blocks filled.
- Confirm the slice branch has six commits (one per task).
- Open a PR to main titled `slice/01: scaffold hardening`.

**Evidence:** _________

---

## Slice 01 Evidence Summary (fill on close)

```
Slice: S01 — Scaffold Hardening & Verify Pipeline
Branch: slice/01-scaffold-hardening
PR: ____________________________________
verify at slice close: ________________
Tasks:
  S01-T01: _________
  S01-T02: _________
  S01-T03: _________
  S01-T04: _________
  S01-T05: _________
  S01-T06: _________
SHALL register deltas:
  SSTPA-COPY-001 → Approved
  SSTPA-CI-001 → Approved
  SSTPA-VERIFY-NEO4J-001 → Approved
  SSTPA-REFDATA-LAYOUT-001 → Approved
  SSTPA-GRAPH-RULES-001 → Candidate (seeded)
  SSTPA-HID-001 → Candidate (seeded)
  SSTPA-IDENTITY-001 → Candidate (seeded)
  SSTPA-OWNERSHIP-TX-001 → Candidate (seeded)
  SSTPA-PAGINATION-001 → Candidate (seeded)
  SSTPA-REFDATA-READONLY-001 → Candidate (seeded)
Open questions logged for S02:
  - Confirm Neo4j 5.26-community supports Cypher 25 feature set we need.
```
