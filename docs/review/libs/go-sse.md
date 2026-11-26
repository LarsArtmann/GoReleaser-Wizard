# go-sse — Library Adoption Audit

| Field             | Value                                                                     |
| ----------------- | ------------------------------------------------------------------------- |
| Library           | `github.com/larsartmann/go-sse` (+ submodules `sseparse`, `ssetest`)      |
| Local checkout    | `/home/lars/projects/go-sse` @ `8acf04c` (2026-11-21)                     |
| Latest local tags | root `v0.7.0` · sseparse `v0.2.1` · ssetest `v0.5.0`                      |
| Consumer          | GoReleaser-Wizard @ `a493950` (2026-11-22)                                |
| Used?             | **No** — `rg 'larsartmann/go-sse'` → 0 hits; `rg 'net/http'` → **0 hits** |
| Verdict           | NOT-ADOPTED → **STAY-OUT**                                                |
| Score             | n/a — 0 applicable capabilities (7 rows, all N/A)                         |
| Provenance        | Fleet audit Group B, 2026-11-22; citations QC-sampled by orchestrator     |

## What the library offers

Server-Sent Events plumbing: wire types and an allocation-minimized serializer, connection lifecycle (headers/heartbeat/Last-Event-ID), generic fan-out broadcasting, reconnection replay, client-side stream parsing, and a handler test harness.

| Capability               | What it does                                          | Citation                      | Applicable here?                                              |
| ------------------------ | ----------------------------------------------------- | ----------------------------- | ------------------------------------------------------------- |
| SSE wire types           | `Event` (name/data/id/retry) + branded `EventID`      | go-sse/event.go:91, :29       | N/A — values exist only for serialization into an SSE stream  |
| Serializer               | `WriteEvent(w, evt)` — multi-line/CRLF-safe           | go-sse/event.go:167           | N/A — consumer writes files/stdout, never `text/event-stream` |
| Connection lifecycle     | `NewStream(w, r)` — headers, heartbeat, Last-Event-ID | go-sse/stream.go:83           | N/A — requires `http.ResponseWriter`                          |
| Fan-out broadcasting     | Generic `NewBroadcaster[T]`                           | go-sse/broadcaster.go:77      | N/A — no subscribers, no long-lived process                   |
| Reconnection replay      | `EventStore` interface + replay                       | go-sse/replay.go:9            | N/A — no reconnecting clients                                 |
| sseparse: client parsing | `NewStreamReader(r)` parses SSE from any reader       | go-sse/sseparse/reader.go:197 | N/A — consumer never consumes an SSE feed                     |
| ssetest: handler harness | `ReadEvents(r)` + request options                     | go-sse/ssetest/compat.go:35   | N/A — only useful for SSE-emitting handlers                   |

## What GoReleaser-Wizard does today

No streaming, no realtime, no long-lived process: `rg 'net/http'` → 0 hits; no `Flush()`/`text/event-stream` anywhere; the process model is run-then-exit (Cobra command tree). "Streaming" output in the wizard means printing progress to the terminal. ROADMAP "TUI polish" items (preview screen, progress indicators) are bubbletea/huh concerns, incompatible with SSE by design.

## Why STAY-OUT

- Every public entrypoint (`NewStream`, `Broadcaster`, `WriteEvent` to a stream) presupposes a server holding connections open; heartbeat/reconnect/replay are meaningless in a process that exits in seconds.
- The library's own README scope ("every Go project that serves SSE…") excludes this consumer by definition.
- ROADMAP non-goal on Web UI + zero TODO items mentioning streaming/realtime/events = no demand signal anywhere in the planning corpus.
- No partial value: even sseparse (a pure reader) needs an SSE producer somewhere in the pipeline; there is none.
- Adopting would add transitive deps for zero callable code paths.

## RECONSIDER-WHEN

A `goreleaser-wizard serve`/watch mode keeps a connection open to push live release progress — which first requires the ROADMAP Web-UI non-goal to be revoked and a server dependency to land. Concretely: when the `net/http` hit count in this repo goes above zero _for a listening server_.

## How to verify

```bash
cd /home/lars/projects/GoReleaser-Wizard
rg -c 'net/http|text/event-stream' --glob '*.go'   # 0 hits
rg -c 'larsartmann/go-sse' .                        # 0 hits
git -C /home/lars/projects/go-sse tag --list 'v*' --sort=-v:refname | head -1   # v0.7.0
```
