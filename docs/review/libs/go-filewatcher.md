# go-filewatcher — Library Adoption Audit

| Field            | Value                                                                                                                  |
| ---------------- | ---------------------------------------------------------------------------------------------------------------------- |
| Library          | `github.com/larsartmann/go-filewatcher/v2`                                                                             |
| Local checkout   | `/home/lars/projects/go-filewatcher` @ `26459ba` (2026-11-17)                                                          |
| Latest local tag | `v2.5.0`                                                                                                               |
| Consumer         | GoReleaser-Wizard @ `a493950` (2026-11-22)                                                                             |
| Used?            | **No** — `rg 'larsartmann/go-filewatcher'` → 0 hits; `fsnotify` appears in go.mod only as an indirect dep of bubbletea |
| Verdict          | NOT-ADOPTED → **STAY-OUT**                                                                                             |
| Score            | n/a — 0 applicable capabilities (7 rows, all N/A)                                                                      |
| Provenance       | Fleet audit Group A, 2026-11-22; citations QC-sampled by orchestrator                                                  |

## What the library offers

A daemon-grade recursive file watcher: fsnotify-based `Watcher` with rich events, ~20 filters, 19 middleware, ~40 options, and operational resilience (self-heal, circuit breaker, health, metrics, OTel).

| Capability                                                    | What it does                                                                                   | Citation                                                                     | Applicable here?                                                                |
| ------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ------------------------------------------------------------------------------- |
| `New(paths, opts...)` + `Watch(ctx)`                          | Recursive watching returning `<-chan Event`                                                    | go-filewatcher/watcher.go:237, :369; :62                                     | N/A — single-shot generate-then-exit; no event loop exists or is wanted         |
| `WatchOnce(ctx)`                                              | Block for a single event                                                                       | go-filewatcher/watcher.go:462                                                | N/A — nothing to wait for                                                       |
| Rich `Event`                                                  | Path/Op/Timestamp/IsDir/Size/ModTime/optional SHA-256                                          | go-filewatcher/event.go:99, :14                                              | N/A                                                                             |
| 20 filters + combinators                                      | `FilterExtensions`, `FilterAnd/Or/Not`, gitignore-aware                                        | go-filewatcher/filter.go:24, :289, :303, :316; go-filewatcher/options.go:349 | N/A                                                                             |
| 19 middleware                                                 | Logging, recovery, rate limiting                                                               | go-filewatcher/middleware.go:39, :46, :66                                    | N/A                                                                             |
| ~40 options                                                   | Debounce (options.go:24), per-path debounce (:36), polling fallback (:232), max-watches (:380) | go-filewatcher/options.go                                                    | N/A                                                                             |
| Resilience: self-heal, circuit breaker, health, metrics, OTel | ENOSPC recovery, NFS polling, Prometheus/OTel hooks                                            | watcher_selfheal.go, health.go, metrics.go, otel.go                          | N/A — long-running-daemon machinery; dead weight in a CLI that exits in seconds |

## What GoReleaser-Wizard does today

Nothing. `rg 'watch|Watch|fsnotify' --type go` over the consumer returns no file-watching code; there is no dev-server, no `watch` subcommand, no regeneration loop. The closest concept — "rerun generation" — is manual (`go build` + rerun, AGENTS.md gotcha 14). fsnotify v1.10.1 is in go.mod only as `// indirect` (pulled by charm.land/bubbletea).

## Why STAY-OUT

- **No use case**: a wizard that generates files and exits has nothing to observe; watching would require inventing a daemon mode first.
- **Heaviest dependency tree of the five Group-A libs** (gogenfilter, go-gitignore, x/time, doublestar, logr, pprof indirect) — a poor trade for code that would never execute.
- Its headline features (debounce, self-heal, circuit breaker, OTel) address exactly the operational profile (long-running, flaky filesystems, NFS) this consumer lacks.
- Compatibility note (non-blocker): module requires Go ≥1.26.7 and uses `encoding/json/v2` (event.go:5) — fine on the consumer's Go 1.27.
- ROADMAP.md / TODO_LIST.md contain no watch/dev-mode item.

## RECONSIDER-WHEN

The wizard gains an interactive/dev mode — e.g. `goreleaser-wizard dev` re-rendering `.goreleaser.yaml` when the user edits their config or templates — i.e. any feature requiring a long-lived process observing files. At that point this is the right library.

## How to verify

```bash
cd /home/lars/projects/GoReleaser-Wizard
rg -ln 'fsnotify|Watch\(' --type go -g '!*_test.go'    # no watching code
grep -n fsnotify go.mod                                  # v1.10.1 // indirect only
rg -c 'larsartmann/go-filewatcher' .                    # 0 hits
git -C /home/lars/projects/go-filewatcher tag --list 'v*' --sort=-v:refname | head -1   # v2.5.0
```
