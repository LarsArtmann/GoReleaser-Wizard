# go-idempotency — Library Adoption Audit

| Field             | Value                                                                 |
| ----------------- | --------------------------------------------------------------------- |
| Library           | `github.com/larsartmann/go-idempotency` (+ `idemtest`)                |
| Local checkout    | `/home/lars/projects/go-idempotency` @ `220c218` (2026-11-20)         |
| Latest local tags | root `v0.6.0` · idemtest `v0.1.1`                                     |
| Consumer          | GoReleaser-Wizard @ `a493950` (2026-11-22)                            |
| Used?             | **No** — `rg 'larsartmann/go-idempotency'` → 0 hits                   |
| Verdict           | NOT-ADOPTED → **STAY-OUT**                                            |
| Score             | n/a — 0 applicable capabilities (8 rows, all N/A)                     |
| Provenance        | Fleet audit Group A, 2026-11-22; citations QC-sampled by orchestrator |

## What the library offers

An idempotency-key SDK for services: a `Store` contract (you implement the backend), atomic `CheckAndRecord` dedup, an HTTP middleware, a contract test suite, and test drills.

| Capability                       | What it does                                                    | Citation                                                       | Applicable here?                                                                       |
| -------------------------------- | --------------------------------------------------------------- | -------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| `Store` interface                | `Seen`/`Record`/`CheckAndRecord` with atomicity + TTL contracts | go-idempotency/store.go:77                                     | N/A — no at-least-once delivery, no commands, no racing invocations                    |
| `CheckAndRecord` atomic dedup    | Concurrent same-key callers: exactly one wins                   | go-idempotency/store.go:103 (interface), :271 (MemoryStore)    | N/A — the "run twice" hazard is solved in-process by `GenerationPreflightJob`          |
| `ErrDuplicate` / `ErrInvalidTTL` | Conflict-/Rejection-classified sentinels                        | go-idempotency/store.go:21, :41                                | N/A — duplicate-file refusal is a plain error at cmd/goreleaser-wizard/jobs.go:486-490 |
| `Forever` TTL sentinel           | Permanent "only once ever" claims                               | go-idempotency/store.go:63                                     | N/A — generation is deliberately re-runnable (`--force`, backups)                      |
| `MemoryStore` (dev/test only)    | In-map store with sweeper                                       | go-idempotency/store.go:157, :177                              | N/A — production needs a Redis/SQL backend; consumer has no store of any kind          |
| HTTP middleware                  | Key-extraction + dedup around `http.Handler`                    | go-idempotency/middleware/http.go:30                           | N/A — no HTTP server                                                                   |
| Contract test suite              | Conformance suite for custom Stores                             | go-idempotency/contract/contract.go                            | N/A — nothing to conform                                                               |
| `idemtest` drills                | `RequireExactlyOnce`, `RequireOneWinner`, `PinnedClock`         | go-idempotency/idemtest/drill.go:36, :82; idemtest/clock.go:22 | N/A                                                                                    |

## What GoReleaser-Wizard does today

The nearest analog is **idempotency-by-refusal**: `GenerationPreflightJob` (cmd/goreleaser-wizard/jobs.go:439-503) stats every target file up front and refuses — "refusing to overwrite existing file(s): … (use --force to overwrite)" at jobs.go:487 — unless `--force`; per-job guards repeat this (jobs.go:98-104 for `.goreleaser.yaml`, jobs.go:346 for Dockerfile). This delivers the DOMAIN_LANGUAGE.md "Atomic generation" property without any store: one process, one pass, existence checks before any write. There is no keyed-command dispatch, no ack-loss scenario, no cross-process concurrency.

## Why STAY-OUT

- **Category mismatch**: the library deduplicates retried network commands under at-least-once delivery; the consumer is a single-shot local CLI with no delivery semantics, no server, no DB to back a `Store`.
- **Adoption would weaken the domain property**: a TTL-windowed "already generated" claim could silently skip a legitimate regeneration after a crash, whereas preflight refusal is deterministic and inspectable.
- The SDK philosophy (implement your own backend) means the real cost is a Redis/SQL store + lifecycle the consumer has no infrastructure for.
- Transitive dependency on go-error-family v0.12.0 (go-idempotency/go.mod) — inherits that STAY-OUT verdict.

## RECONSIDER-WHEN

The wizard becomes a long-running service exposing an API where clients retry `generate` requests — then `Store` + `middleware.HTTP` + the contract suite become the right tools with a shared Redis backend.

## How to verify

```bash
cd /home/lars/projects/GoReleaser-Wizard
rg -c 'larsartmann/go-idempotency' .                          # 0 hits
sed -n '439,503p' cmd/goreleaser-wizard/jobs.go               # preflight refusal; message at :487
git -C /home/lars/projects/go-idempotency tag --list 'v*' --sort=-v:refname | head -1   # v0.6.0
rg 'go-error-family' /home/lars/projects/go-idempotency/go.mod # transitive dep confirmed
```
