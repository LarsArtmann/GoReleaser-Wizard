# go-etag — Library Adoption Audit

| Field             | Value                                                                                      |
| ----------------- | ------------------------------------------------------------------------------------------ |
| Library           | `github.com/larsartmann/go-etag` (+ submodules `server`, `client`, `entitytag`, `metrics`) |
| Local checkout    | `/home/lars/projects/go-etag` @ `39f9f37` (2026-11-18)                                     |
| Latest local tags | root `v0.7.1` (deprecated alias shim) · server/client/entitytag/metrics `v0.6.2`           |
| Consumer          | GoReleaser-Wizard @ `a493950` (2026-11-22)                                                 |
| Used?             | **No** — `rg 'larsartmann/go-etag'` → 0 hits; `rg 'http\.Client'` → **0 hits**             |
| Verdict           | NOT-ADOPTED → **STAY-OUT**                                                                 |
| Score             | n/a — 0 applicable capabilities (7 rows, all N/A)                                          |
| Provenance        | Fleet audit Group B, 2026-11-22; citations QC-sampled by orchestrator                      |

## What the library offers

HTTP conditional-request handling on both sides: server ETag middleware (compute + answer `If-None-Match` with 304), a client conditional-GET cache `RoundTripper` with RFC 9111 freshening, an RFC 7232 entity-tag domain type, and hit/miss metrics.

| Capability                   | What it does                                                               | Citation                                           | Applicable here?                                                   |
| ---------------------------- | -------------------------------------------------------------------------- | -------------------------------------------------- | ------------------------------------------------------------------ |
| Server ETag middleware       | `etag.New(cfg)` buffers body, computes ETag, 304s                          | go-etag/server/config.go:157                       | N/A — middleware for an HTTP server the consumer doesn't have      |
| Client conditional-GET cache | `etagclient.NewTransport(next, opts)` replays validators, rebuilds 304→200 | go-etag/client/transport.go:40                     | N/A — consumer makes zero HTTP requests; it execs binaries instead |
| RFC 9111 freshening          | 304 freshens cached headers; unsafe-method invalidation                    | go-etag/client/rebuild.go:13; client/options.go:72 | N/A — no cacheable HTTP exchanges exist                            |
| Entity-tag domain type       | `NewETag`/`ParseETag`/`MatchesIfNoneMatch` (weak/strong comparison)        | go-etag/entitytag/entity_tag.go:52, :111, :176     | N/A — no input source for validator headers                        |
| Observability counters       | `metrics.Attach(cfg)` hit/miss counters                                    | go-etag/metrics/metrics.go:38                      | N/A — measures a server that doesn't exist                         |
| Bounded-memory eviction      | FIFO `MaxEntries` + `MaxBodyBytes` caps                                    | go-etag/client/rebuild.go:107                      | N/A — long-lived-process optimization; consumer is single-shot     |
| Root alias shim              | Root module re-exports server package; slated for removal in v1.0.0        | go-etag README (deprecated alias shim)             | N/A — nothing to alias                                             |

## What GoReleaser-Wizard does today

No HTTP client or server: `rg 'http\.Client'` → 0 hits; `rg 'net/http'` → 0 hits. All external interaction is `exec.CommandContext` (internal/git/commands.go:57). File "caching" concerns are local-file backups (`createBackup()` per AGENTS.md gotcha 4), not HTTP cache validators.

## Why STAY-OUT

- Both sides of go-etag (server middleware, client `RoundTripper`) attach to `net/http` objects the consumer never constructs.
- The wizard talks to git/docker/goreleaser as subprocesses precisely to avoid network plumbing; adding an HTTP transport layer would be an architectural reversal, not an increment.
- go-etag's value (304 savings on repeated GETs) presumes a long-lived or high-volume HTTP caller; a one-shot CLI cannot amortize the dependency.
- No roadmap theme involves fetching remote resources.
- Note for any hypothetical future adoption: the root module is a deprecated shim — target `go-etag/server` or `go-etag/client` submodules directly.

## RECONSIDER-WHEN

The consumer gains an HTTP client surface that re-fetches the same resource repeatedly — e.g. a `wizard validate --remote` mode polling a Forgejo/GitHub API, or version-check GETs to the module proxy. Until an `http.Client` exists in the codebase, no.

## How to verify

```bash
cd /home/lars/projects/GoReleaser-Wizard
rg -c 'http\.Client|net/http' --glob '*.go'   # 0 hits
rg -c 'larsartmann/go-etag' .                  # 0 hits
git -C /home/lars/projects/go-etag tag --list 'v*' --sort=-v:refname | head -1   # v0.7.1
git -C /home/lars/projects/go-etag tag --list 'server/v*' --sort=-v:refname | head -1   # server/v0.6.2
```
