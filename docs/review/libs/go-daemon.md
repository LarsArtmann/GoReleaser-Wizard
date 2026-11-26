# go-daemon — Library Adoption Audit

| Field            | Value                                                                 |
| ---------------- | --------------------------------------------------------------------- |
| Library          | `github.com/LarsArtmann/go-daemon` (note the capital module path)     |
| Local checkout   | `/home/lars/projects/go-daemon` @ `b6af9f1` (2026-11-18)              |
| Latest local tag | `v0.4.0`                                                              |
| Consumer         | GoReleaser-Wizard @ `a493950` (2026-11-22)                            |
| Used?            | **No** — `rg 'larsartmann/go-daemon'` (case-insensitive) → 0 hits     |
| Verdict          | NOT-ADOPTED → **STAY-OUT**                                            |
| Score            | n/a — 0 applicable capabilities (9 rows, all N/A)                     |
| Provenance       | Fleet audit Group C, 2026-11-22; citations QC-sampled by orchestrator |

## What the library offers

Daemon/server lifecycle machinery: a Unix-socket HTTP server with stale-socket cleanup and graceful drain, sd_notify integration (READY/STOPPING/watchdog), socket-dialing HTTP transports, content negotiation (JSON/CBOR), encoded response writers, SSE parsing, a flight recorder, and a liveness probe.

| Capability                 | What it does                                                    | Citation                                | Applicable here?                                                              |
| -------------------------- | --------------------------------------------------------------- | --------------------------------------- | ----------------------------------------------------------------------------- |
| Socket server lifecycle    | Stale-socket cleanup, 0600 perms, graceful drain, ordered hooks | go-daemon/socket.go:60, :81, :183, :213 | N/A — single-shot CLI with no listener; 0 `net/http` imports in non-test code |
| Any-network server variant | Same lifecycle over TCP / pre-built listener                    | go-daemon/socket.go:158, :150           | N/A                                                                           |
| sd_notify integration      | READY/STOPPING/watchdog keepalive; no-ops without systemd       | go-daemon/notify.go:59, :75, :46        | N/A — process exits in seconds under a user terminal                          |
| Socket HTTP transports     | Pooled keep-alive client dialing a daemon socket                | go-daemon/transport.go:25, :41          | N/A — consumer shells out to git via os/exec (internal/git/commands.go:57)    |
| Content negotiation        | JSON/CBOR resolution with q-value ranking                       | go-daemon/negotiation.go:73, :98        | N/A                                                                           |
| Encoded response writers   | Pooled-buffer JSON/encoded/error writes                         | go-daemon/write.go:29, :60, :66         | N/A — output is stdout via charm/log                                          |
| SSE parsing                | WHATWG-conformant framing                                       | go-daemon/sse.go:42                     | N/A                                                                           |
| Flight recorder            | Execution-trace recorder with slow-request trigger, metrics     | go-daemon/flightrecorder.go:198         | N/A                                                                           |
| Liveness probe             | `IsSocketAlive`                                                 | go-daemon/socket.go:418                 | N/A                                                                           |

Dependencies it would drag in: go-systemd/v22, go-codec, go-flightrecorder, go-sse/sseparse, prometheus/client_golang (go-daemon/go.mod:5-10).

## What GoReleaser-Wizard does today

No long-running mode: the entry is a Cobra root command with `init/generate/validate/version` subcommands (cmd/goreleaser-wizard/main.go:227); the process runs a Job/Workflow Execute+Rollback sequence and exits. `rg 'net/http'` over non-test `*.go` → 0 hits.

## Why STAY-OUT

- Every capability presupposes a resident process serving or dialing HTTP over a socket; the wizard exits after generation.
- Dependency weight (systemd, Prometheus, CBOR stack) is absurd for a config generator whose go.mod is deliberately tiny.
- ROADMAP.md "Non-goals" (:61-69) contains nothing daemon-shaped; only "Web UI" is named and it is rejected (:67).
- No TODO_LIST item would need a listener.

## RECONSIDER-WHEN

The wizard gains a resident mode — e.g. a file-watcher daemon auto-regenerating configs, or a socket API for editor/IDE integration. Revisit together with go-filewatcher and go-appkit (the same flip).

## How to verify

```bash
cd /home/lars/projects/GoReleaser-Wizard
rg -c -i 'larsartmann/go-daemon' .                # 0 hits
rg -c 'net/http' --glob '*.go' --glob '!*_test.go' # 0 hits
git -C /home/lars/projects/go-daemon tag --list 'v*' --sort=-v:refname | head -1   # v0.4.0
sed -n '5,10p' /home/lars/projects/go-daemon/go.mod # systemd/prometheus/sse dep weight
```
