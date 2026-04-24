# Slice 12 — Installer, Air-Gap Bundle, Telemetry Dashboards, Release (Phase B, Outline)

**Goal:** Produce distributable installers for Linux, Windows, and (pending user go/no-go per Master §10 Q6) macOS. Validate end-to-end on an air-gapped Windows 11 Enterprise environment. Ship Grafana dashboards wired to the OTel/Prom/Tempo stack. Close the project on an MVP release artifact.

**Architecture:** Tauri's bundler is the primary installer path. The installer bundle ships:
- Signed Tauri shell
- Bundled backend binary (static, `CGO_ENABLED=0`)
- Bundled Neo4j (Community Edition, offline zip; not embedded in Tauri bundle but delivered alongside)
- Bundled OTel Collector + Prometheus + Tempo + Grafana (Docker images pre-pulled to an offline tar archive; optional Grafana-less install flag)
- Bundled reference framework normalized artifacts (from S05) for first-run import
- Bundled example Capability fixture (from S11)
- Bundled help content (from S11)

A single-host MVP runs the stack via `docker load` + `docker compose up`. The installer orchestrates this per platform.

**Tech Stack:** Tauri 2 bundler, `nsis`/`wix` (Windows), `.deb`/`.rpm` (Linux), `.dmg` + codesign + notarytool (macOS pending), Grafana provisioning files, existing telemetry stack.

**Pre-reads:**
- SRS §2.1 Startup Software (shutdown behavior)
- SRS §2.3 Docker networks + §2.3.5 Backend Telemetry
- SRS §6 Constraints (Windows 11 Enterprise air-gapped; Ubuntu Studio 25.04 dev)
- SRS Section "Installer" (end of draft)
- Outputs from all prior slices
- `infra/docker/**`
- Conventions doc.

**Invariants:**
- No network access during install on air-gapped targets.
- Installer validates required resources before starting (SRS §6 Installer) — RAM, disk, CPU, Docker availability.
- Only Caddy is exposed on the edge network; Neo4j stays internal.
- Uninstall preserves user data unless user opts in to wipe (data = Neo4j volume + Grafana volumes).
- Every commit in this slice includes a CHANGELOG entry.

**Tasks:**

### S12-T01: Release engineering RFC
Difficulty: medium · Integration Checkpoint: yes (user approval) · Files: `docs/architecture/release-engineering-rfc.md`.
**Agent Briefing:** Specify the bundle layout per platform, Docker image preload method, resource-check UX, uninstall contract, update path for v1.1, signing/notarisation plan (or punt to v1.1). Pause for approval.

### S12-T02: Backend release build reproducibility
Difficulty: low · Files: `backend/Dockerfile` (already exists), add `backend/cmd/api/release.go` with build-info embedding; Makefile `backend-release`.
**Agent Briefing:** `go build -trimpath -ldflags "-X sstpa-tool/backend/internal/version.Version=$(VERSION) -X ...Commit=$(COMMIT)"`. Ensures version.Version (instead of version.Dev) renders in `/health`.

### S12-T03: Offline Docker image bundle
Difficulty: medium · Files: `tools/release/pack-images.sh` (or `.go`), `tools/release/image-manifest.txt`.
**Agent Briefing:** Pulls all compose images at pinned tags and `docker save`s them into a single tar for offline loading; manifest records SHA256 + image name.

### S12-T04: Resource-check library
Difficulty: medium · Files: `apps/desktop-shell/src-tauri/src/resourcecheck.rs`; Tauri command + UI step in installer.
**Agent Briefing:** Checks RAM ≥ 8 GiB, disk ≥ 40 GiB free, CPU ≥ 4 cores, Docker Engine responding, no stale SSTPA Neo4j container. Fails closed with a readable error before anything is written.

### S12-T05: Linux installer (.deb + .rpm + AppImage)
Difficulty: high · Integration Checkpoint: yes · Files: `tools/release/linux/**`; CI workflow.
**Agent Briefing:** Builds three artifacts. Post-install hook loads the Docker image tar and starts the compose stack. AppImage launches the Tauri shell.

### S12-T06: Windows installer (.msi via WiX)
Difficulty: high · Integration Checkpoint: yes · Files: `tools/release/windows/**`; CI workflow.
**Agent Briefing:** MSI bundles Tauri shell; post-install step configures Docker Desktop (or WSL2) detection and loads the image tar. Treats Docker absence as a user-recoverable error.

### S12-T07: macOS installer (.dmg, notarised) [conditional]
Difficulty: high · Integration Checkpoint: yes · Files: `tools/release/macos/**`; CI workflow.
**Agent Briefing:** Only executed if Master §10 Q6 resolves "macOS required". If deferred, file a stub doc and close this task with `deferred` status in Evidence.

### S12-T08: Reference framework data inclusion
Difficulty: medium · Files: `tools/release/pack-reference.sh`; installer integrates.
**Agent Briefing:** Pulls S05 normalized artifacts (under their license posture) into the installer bundle; first-run import triggered by Startup Software.

### S12-T09: Example Capability + Help content inclusion
Difficulty: low · Files: installer wiring.
**Agent Briefing:** First-run seeds example Capability (S11-T08) and Help content (S11-T04) after schema bootstrap.

### S12-T10: Grafana dashboards
Difficulty: medium · Files: `infra/docker/grafana/provisioning/**`; at minimum: HTTP latency histograms per route, Neo4j query duration, mutation TPS, notification emission counter, commit failure counter, container-level CPU/mem.
**Agent Briefing:** JSON dashboard definitions committed; Grafana provisioning loads them automatically.

### S12-T11: Air-gap validation procedure
Difficulty: medium · Integration Checkpoint: yes · Files: `docs/release/airgap-validation-procedure.md`.
**Agent Briefing:** Step-by-step manual test plan for running the installer on an air-gapped Windows 11 Enterprise VM with no internet. Records acceptable artefacts. Handed to a human for execution (we do not own the environment).

### S12-T12: CHANGELOG, version tag, MVP release notes
Difficulty: low · Files: `CHANGELOG.md`, `docs/release/v0.1.0-notes.md`; git tag `v0.1.0`.
**Agent Briefing:** Summarise every slice's shipped capability with SRS citations. Produce the release notes.

### S12-T13: Slice integration gate + final verification gate
Difficulty: low · Integration Checkpoint: yes.
**Acceptance:** every Req in shall-register either Approved, Deferred with sign-off, or Superseded. Verification matrix 100% populated for every Approved Req. Air-gap validation procedure executed (or explicitly deferred).

**Integration gate criteria:**
- Release artefacts build from clean on each platform's CI runner.
- Installer demo: offline install, first-run seed, launch shell, pick SoI, Commit a change, view Grafana dashboard showing the request.
- `make verify` green on main at release tag.
- Master orchestration plan closing note appended.

---

## Project Close-Out Checklist

After S12 Integration gate passes:

- [ ] Merge all slice PRs to `main`.
- [ ] Tag `v0.1.0`.
- [ ] Archive slice branches.
- [ ] Run the `srs-audit` OpenCode command once more and record gaps.
- [ ] File follow-up issues for every Deferred Req.
- [ ] Append "Project Close-Out" note to master orchestration plan.
