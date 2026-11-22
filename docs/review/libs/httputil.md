# httputil (+ server_timing, endpoints) — Library Adoption Audit

| Field | Value |
| --- | --- |
| Library | `github.com/larsartmann/httputil` (+ submodules `server_timing`, `endpoints`) |
| Local checkout | `/home/lars/projects/httputil` @ `34dbf05` (2026-11-21) |
| Latest local tags | root `v1.9.0` · server_timing `v1.0.3` · endpoints `v0.1.0` |
| Consumer | GoReleaser-Wizard @ `a493950` (2026-11-22) |
| Used? | **No** — `rg 'larsartmann/httputil'` → 0 hits; `rg 'net/http'` over consumer `.go` files → **0 hits** |
| Verdict | NOT-ADOPTED → **STAY-OUT** |
| Score | n/a — 0 applicable capabilities (12 rows, all N/A) |
| Provenance | Fleet audit Group B, 2026-11-22; citations QC-sampled by orchestrator |

## What the library offers

A production HTTP middleware toolkit: chaining, panic recovery, CORS, security headers/CSP nonce, request IDs, structured logging, timeout/body-size/rate-limit/CSRF/language middlewares, compression/decompression, and a configurable server with graceful shutdown — plus a W3C `Server-Timing` instrumentation submodule and a self-describing `/endpoints` registry submodule.

| Capability | What it does | Citation | Applicable here? |
| --- | --- | --- | --- |
| Middleware chaining | `Chain(handler, middlewares...)` | httputil/recorder.go:158 | N/A — requires an `http.Handler`; consumer has none |
| Panic recovery | `Recovery(logger)` → 500s | httputil/recovery.go:33 | N/A — CLI panics surface via Cobra/fang |
| CORS | Configurable cross-origin policy | httputil/cors.go:113 | N/A — meaningless for a local CLI |
| Security headers / CSP nonce | `SecurityHeaders(cfg)` | httputil/security.go:95 | N/A — applies to HTTP responses, not stdout |
| Request ID | Per-request IDs in context | httputil/requestid.go:69 | N/A — single-shot process, no multiplexing |
| Structured request logging | `Logging(logger)` | httputil/logging.go:12 | N/A — consumer logs via charm.land/log directly |
| Timeout / body-size / rate-limit / CSRF / language | Handler-safety middlewares | httputil/timeout.go:12, maxbodysize.go:62, ratelimit_keyed.go:223, csrf.go:586, language.go:445 | N/A — no inbound request lifecycles |
| Compression / decompression | Response/request bodies with bomb protection | httputil/compression.go:269, decompression.go:107 | N/A — no HTTP bodies in either direction |
| Server lifecycle | `NewServer(cfg, handler)` w/ graceful shutdown | httputil/server.go:197 | N/A — the binary never listens on a socket |
| server_timing submodule | W3C `Server-Timing` header instrumentation | httputil/server_timing/server_timing.go:55, :369 | N/A — emits an HTTP response header; no responses |
| endpoints submodule | Self-describing `GET /endpoints` route listing | httputil/endpoints/endpoints.go:93, :111 | N/A — enumerates `http.ServeMux` routes; consumer has no mux |

## What GoReleaser-Wizard does today

The space does not exist: `rg 'net/http'` → 0 hits across consumer `.go` files; go.mod pulls no HTTP libs. The only external-tool interaction is process spawning: `internal/git/commands.go:57` (`exec.CommandContext(c.ctx, "git", args...)`) and :359. Interactive UX is terminal-side via huh (cmd/goreleaser-wizard/tui_wizard.go). ROADMAP.md:67 lists **"Web UI" under Non-goals** ("the CLI/TUI is the product; a web layer duplicates it without a consumer").

## Why STAY-OUT

- **Zero `net/http` surface**: importing httputil would drag a multi-dependency HTTP middleware stack into a binary that never opens a socket — dead weight in the module graph and the nix vendorHash.
- **The only conceivable bridge is explicitly rejected** (ROADMAP.md:67); no roadmap theme (config surface, validation, CI families, extensibility, mechanical guarantees) involves serving HTTP.
- The wizard's "server-ish" roadmap items (self-hosted runner label, release-shadow workflow) are generated YAML *text*, not runtime behavior of this binary.
- Even the self-contained submodules (server_timing, endpoints) presuppose a live `http.ResponseWriter`/`ServeMux`; no partial adoption is worth taking.
- Adoption would violate the project's minimal-deps clean-architecture narrative (domain purity, AGENTS.md).

## RECONSIDER-WHEN

The wizard ever *serves* anything — a `goreleaser-wizard serve` preview/dashboard mode or an HTTP API for CI integration. First flip trigger: any roadmap change deleting the "Web UI — non-goal" entry at ROADMAP.md:67.

## How to verify

```bash
cd /home/lars/projects/GoReleaser-Wizard
rg -c 'net/http' --glob '*.go'          # 0 hits
rg -c 'larsartmann/httputil' .          # 0 hits
sed -n '67p' ROADMAP.md                 # "Web UI" non-goal
sed -n '57p' internal/git/commands.go   # exec.CommandContext — subprocess, not HTTP
git -C /home/lars/projects/httputil tag --list 'v*' --sort=-v:refname | head -1   # v1.9.0
```
