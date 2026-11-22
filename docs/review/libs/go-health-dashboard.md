# go-health-dashboard — Library Adoption Audit

| Field | Value |
| --- | --- |
| Library | `github.com/larsartmann/go-health-dashboard` |
| Local checkout | `/home/lars/projects/go-health-dashboard` @ `19af9cf` (2026-11-20) |
| Latest local tag | `v0.13.3` |
| Consumer | GoReleaser-Wizard @ `a493950` (2026-11-22) |
| Used? | **No** — 0 hits for the go-health family |
| Verdict | NOT-ADOPTED → **STAY-OUT** |
| Score | n/a — 0 applicable capabilities (6 rows, all N/A) |
| Provenance | Fleet audit Group C, 2026-11-22; citations QC-sampled by orchestrator |

## What the library offers

A real-time HTML status dashboard riding on go-health probes: Datastar SSE updates, severity grouping, ~40 options (title, CSP nonce, rate limits, OIDC auth, webhooks, trends), one-call samber/do registration, a standalone `health-hub` federation binary, and Prometheus/Grafana deployment assets.

| Capability | What it does | Citation | Applicable here? |
| --- | --- | --- | --- |
| Real-time HTML dashboard | Datastar SSE status page mounted on a mux | go-health-dashboard/dashboard.go:143, :28 | N/A — needs a go-health Probe and an HTTP mux; consumer has neither |
| DI one-call registration | `dashboard.Register(injector, probe, …)` | go-health-dashboard/di.go:29 | N/A — no DI container |
| Route surface | `/health`, `/health/sse`, `/healthz`, trends, metrics | go-health-dashboard/routes.go:4, :30; README.md:102-114 | N/A |
| Options breadth | ~40 options: OIDC, webhooks, trends, rate limits | go-health-dashboard/options.go:179-597 (e.g. `WithMetrics` :187, `WithTrend` :198, `WithRateLimit` :283, `WithWebhookHMAC` :378) | N/A |
| Federation hub binary | Standalone `health-hub` aggregating remote probes | go-health-dashboard/cmd/health-hub/ | N/A — closest plausible use, still blocked (see verdict) |
| Prometheus/Grafana deploy | docker-compose + dashboards | go-health-dashboard/deploy/ | N/A |

Dependency tree: templ, go-datastar, go-health v0.10.0, go-sse, templ-components ×4 (go-health-dashboard/go.mod:6-18).

## What GoReleaser-Wizard does today

Nothing: no dashboard, no web surface. ROADMAP.md:67 — "**Web UI**: the CLI/TUI is the product; a web layer duplicates it without a consumer" (explicit non-goal). The TUI (cmd/goreleaser-wizard/tui_wizard.go, huh-based) is the entire interactive surface.

## Why STAY-OUT

- Requires adopting go-health + samber/do + an HTTP server first — a three-layer stack the roadmap explicitly rejects.
- Heavy deps (templ codegen, Datastar SDK, SSE) would roughly double the consumer's direct-dependency list for a page nobody requested.
- Even the fleet-observability angle (dashboard over generated projects' release health) belongs to the standalone `health-hub` binary deployed elsewhere — not to embedding this module in the wizard.
- No TODO_LIST/ROADMAP item mentions a status page.

## RECONSIDER-WHEN

The owner lifts the Web-UI non-goal AND wants a fleet status page — at that point run `health-hub` as a separate deployment; do not embed this module in the wizard.

## How to verify

```bash
cd /home/lars/projects/GoReleaser-Wizard
rg -c 'larsartmann/go-health' .                     # 0 hits
sed -n '67p' ROADMAP.md                              # Web UI non-goal
git -C /home/lars/projects/go-health-dashboard tag --list 'v*' --sort=-v:refname | head -1   # v0.13.3
```
