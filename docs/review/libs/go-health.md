# go-health — Library Adoption Audit

| Field | Value |
| --- | --- |
| Library | `github.com/larsartmann/go-health` |
| Local checkout | `/home/lars/projects/go-health` @ `3ca6644` (2026-11-20) |
| Latest local tag | `v0.10.0` |
| Consumer | GoReleaser-Wizard @ `a493950` (2026-11-22) |
| Used? | **No** — `rg 'larsartmann/go-health'` → 0 hits |
| Verdict | NOT-ADOPTED → **STAY-OUT** |
| Score | n/a — 0 applicable capabilities (9 rows, all N/A) |
| Provenance | Fleet audit Group C, 2026-11-22; citations QC-sampled by orchestrator |

## What the library offers

Kubelet-semantics health probing for services: three probes (liveness/readiness/startup) from one `Probe`, critical-service classification, `/healthz`-family JSON routes, a version endpoint, a response envelope with per-check detail, federation/aggregate merging, and shutdown-aware draining.

| Capability | What it does | Citation | Applicable here? |
| --- | --- | --- | --- |
| Three-probe SDK | Liveness/readiness/startup from one `Probe` with kubelet semantics | go-health/probe.go:161 (`New(injector, …)`), :67 | N/A — requires a samber/do DI injector and an HTTP mux; consumer has neither |
| Critical-service classification | Non-critical failures surface without 503 | go-health/options.go:48 | N/A — no services to check |
| Route registration | `/healthz` `/readyz` `/startupz` JSON handlers | go-health/handlers.go:37, :14 | N/A — no HTTP server |
| Version endpoint | `VersionHandler` answering JSON | go-health/version.go:56 | N/A — consumer's version is `fmt.Printf` at cmd/goreleaser-wizard/main.go:256-263 |
| Response envelope | pass/warn/fail status, uptime, per-check detail, `InstanceID` | go-health/types.go:151, :148, :11 | N/A |
| Health-recorder decoupling | Build probe from a plain function | go-health/options.go:127; probe.go:22 | N/A |
| Checks batteries | Named `CheckFunc` map with durations | go-health/checks.go:15, :44 | N/A |
| Federation/aggregate | Merge remote probes | go-health/federation/, aggregate/ | N/A |
| Shutdown awareness | One-way drain with grace period | go-health/options.go:86 | N/A — process lifetime is seconds |

Dependencies: samber/do v2, go-branded-id (go-health/go.mod:5-8).

## What GoReleaser-Wizard does today

No health/status concept. The `version` command prints name/version/date/commit via `fmt.Printf` (cmd/goreleaser-wizard/main.go:256-263); `fang.WithVersion(version)` at main.go:141 handles `--version`. Error reporting is typed domain errors + charm/log + Cobra exit codes.

## Why STAY-OUT

- The library's value is kubelet-facing HTTP probes for long-lived services with real dependencies (DB, redis); the wizard is a single-shot generator with zero runtime dependencies.
- Adoption would require introducing samber/do DI and an HTTP server first — architecture the project has never had and its roadmap explicitly refuses (Web UI non-goal, ROADMAP.md:67).
- Its `VersionHandler` — the only superficially relevant piece — would replace three `fmt.Printf` lines with an HTTP handler: a regression.

## RECONSIDER-WHEN

The wizard becomes a long-running service (daemon/REST API) with checkable dependencies — e.g. a `wizard serve` validating registry/Forgejo connectivity continuously. Revisit together with go-daemon.

## How to verify

```bash
cd /home/lars/projects/GoReleaser-Wizard
rg -c 'larsartmann/go-health' .                        # 0 hits
sed -n '256,263p' cmd/goreleaser-wizard/main.go        # plain Printf version command
git -C /home/lars/projects/go-health tag --list 'v*' --sort=-v:refname | head -1   # v0.10.0
```
