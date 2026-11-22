# go-appkit — Library Adoption Audit

| Field | Value |
| --- | --- |
| Library | `github.com/larsartmann/go-appkit` (+ submodules: cqrs `v0.9.0`, otel `v0.4.0`, health `v0.3.2`, security `v0.3.0`, httpx `v0.2.0`, systemd `v0.2.0`, realtime `v0.2.0`, config `v0.1.1`, preset `v0.1.1`, docs `v0.3.3`, errorpages `v0.1.3`, flightrecorder `v0.2.0`, flightrecorderhealth `v0.1.7`) |
| Local checkout | `/home/lars/projects/go-appkit` @ `10a17ab` (2026-11-21) |
| Latest local tag | root `v0.13.0` |
| Consumer | GoReleaser-Wizard @ `a493950` (2026-11-22) |
| Used? | **No** — `rg 'larsartmann/go-appkit'` → 0 hits |
| Verdict | NOT-ADOPTED → **STAY-OUT** |
| Score | n/a — 0 applicable capabilities (10 rows, all N/A; `config` borderline) |
| Provenance | Fleet audit Group D, 2026-11-22; citations QC-sampled by orchestrator |

## What the library offers

The fleet service foundation: `NewService(ServiceConfig)` HTTP service lifecycle (listener, mux, middleware, graceful shutdown), `preset.New` one-call service base (health + error pages + version from env switches), `config.Load[T]` env-first loader, health/dashboard mounts, a CQRS command bus, OTel wiring, systemd sd_notify, realtime dashboards, HTML error pages, HTTP helpers, and docs serving.

| Capability | What it does | Citation | Applicable here? |
| --- | --- | --- | --- |
| HTTP service lifecycle | Listener, mux, middleware, graceful shutdown | go-appkit/service.go:51; config.go:51 | N/A — consumer runs single-shot and exits |
| `preset` one-call base | Boots service + health + error pages + version from env | go-appkit/preset/preset.go:33; preset/README.md:5-9 | N/A — presupposes a long-running process |
| `config.Load[T]` env-first loader | Defaults + optional file + env + flags onto a flat struct | go-appkit/config/config.go:88 | N/A (borderline) — the one submodule touching the consumer's space, but the consumer's koanf Manager already does file+env+posflag (internal/config/config.go:51-117); appkit/config's own design note *declined* koanf for flat env-first service structs — swapping is a lateral move with churn |
| `health` probes/dashboard | `/health`-family mounts | go-appkit/health/mount.go:91 | N/A — process lives milliseconds |
| `cqrs` command bus | RegisterCommand/RegisterQuery/Dispatch | go-appkit/cqrs/commands.go:55, :65, :82 | N/A — the consumer's "commands" are 4 cobra verbs; the Job engine (job_manager.go:24-31) already covers Execute+Rollback |
| `otel` observability | Metrics/tracing setup | go-appkit/otel/ (client.go, views.go) | N/A |
| `systemd` integration | sd_notify hooks | go-appkit/systemd/hooks.go:57 | N/A — no daemon, no unit file |
| realtime / errorpages / httpx / security / docs | Dashboards, HTML errors, HTTP helpers, headers, docs serving | go-appkit/ dirs; preset/README.md | N/A — every one mounts on an HTTP mux the consumer doesn't have |
| flightrecorder (appkit variant) | Runtime trace with **HTTP** trigger | go-appkit/flightrecorder/ (handler.go, middleware.go) | N/A — trace-dump-on-request needs a server |
| Fleet positioning | "A NEW fleet *service* starts on go-appkit" | go-appkit/README.md:5-6 | N/A — this consumer is not a service |

## What GoReleaser-Wizard does today

No service infrastructure whatsoever — by profile and by architecture ruling (domain purity per AGENTS.md/go-arch-lint). The nearest neighbors in *role*: config loading (koanf Manager, internal/config/config.go:51-117) and lifecycle (Job/Workflow engine with rollback, job_manager.go:24-31, workflow.go:81) — both already purpose-built CLI-side.

## Why STAY-OUT

- **The entire family solves problems a single-shot CLI does not have** — there is no service to boot, probe, drain, or notify.
- The one overlapping submodule (`config`) is a *competing* design to the working koanf setup, positioned as the anti-koanf choice for flat service envs — adopting here would be churn, not an upgrade.
- `preset`'s env-switch model presumes a daemon lifecycle; nothing maps onto "parse flags → render templates → exit".
- Zero dependency overlap with the consumer's graph (charm + cobra + koanf), so staying out costs nothing.
- If stricter config errors are ever wanted, the cheaper local move is typed validation on the existing Manager.

## RECONSIDER-WHEN

The wizard ever grows a resident process — a watch/daemon mode or hosted wizard service (which would also flip the templ-components and go-daemon rulings). At that point `preset.New` becomes the front door; `config` would still be evaluated against the incumbent koanf Manager on its own merits.

## How to verify

```bash
cd /home/lars/projects/GoReleaser-Wizard
rg -c 'larsartmann/go-appkit' .                     # 0 hits
sed -n '5,6p' /home/lars/projects/go-appkit/README.md # "fleet service" framing
git -C /home/lars/projects/go-appkit tag --list 'v*' --sort=-v:refname | head -1   # v0.13.0
```
