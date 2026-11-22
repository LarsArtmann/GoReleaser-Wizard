# go-datastar — Library Adoption Audit

| Field | Value |
| --- | --- |
| Library | `github.com/larsartmann/go-datastar` (+ submodules `broadcast`, `datastartest`, `static`) |
| Local checkout | `/home/lars/projects/go-datastar` @ `24a808f` (2026-11-21) |
| Latest local tags | root `v0.7.0` · broadcast/datastartest `v0.7.0` · static `v0.8.0` |
| Consumer | GoReleaser-Wizard @ `a493950` (2026-11-22) |
| Used? | **No** — `rg 'larsartmann/go-datastar'` → 0 hits |
| Verdict | NOT-ADOPTED → **STAY-OUT** |
| Score | n/a — 0 applicable capabilities (8 rows, all N/A) |
| Provenance | Fleet audit Group B, 2026-11-22; citations QC-sampled by orchestrator |

## What the library offers

A Go SDK for the Datastar hypermedia architecture: every patch (elements/signals/redirect/script) is a value that serializes to an SSE event, plus a broadcaster hub with replay, an embedded pinned Datastar JS runtime served via `http.Handler`, and test helpers.

| Capability | What it does | Citation | Applicable here? |
| --- | --- | --- | --- |
| Patch-as-value model | `Patch` interface → every patch yields an `sse.Event` | go-datastar/patch.go:25 | N/A — patches target a browser DOM the consumer doesn't serve |
| Elements patch | `NewElementsPatch(html, opts...)` DOM-merge fragments | go-datastar/elements.go:97 | N/A — no HTML fragments are pushed anywhere |
| Signals patch | `NewSignalsPatch(v, opts...)` state updates as JSON signals | go-datastar/signals.go:61 | N/A — reactive state is client-side in Datastar apps only |
| Redirect / script patches | `NewRedirectPatch(url)` etc. as SSE events | go-datastar/script_convenience.go:16 | N/A — a CLI redirects nothing; it exits with a status code |
| Serve Datastar JS client | `ScriptHandler()` serving the embedded runtime with ETag+Cache-Control | go-datastar/script_handler.go:37 | N/A — requires an HTTP server and a browser |
| Broadcast hub | Pre-wired broadcaster w/ replay, heartbeats | go-datastar/broadcast/broadcaster.go:57 | N/A — fan-out to N connected browsers; N=0 forever in a CLI |
| Embedded JS runtime | `static.Bytes()` pinned Datastar client | go-datastar/static/static.go:66 | N/A — nothing serves it, no browser consumes it |
| Test helpers | `datastartest.Diff(want, got)` | go-datastar/datastartest/diff.go:12 | N/A — tests Datastar-emitting handlers |

## What GoReleaser-Wizard does today

No web UI: the entire UX is a terminal TUI (cmd/goreleaser-wizard/tui_wizard.go, huh/bubbletea) plus generated static files written with `os.WriteFile`. `rg 'net/http'` → 0 hits. Generated *templates* reference web infra for other projects, but the wizard itself never renders or serves HTML. ROADMAP.md:67 lists **"Web UI" under Non-goals**.

## Why STAY-OUT

- go-datastar is a *server-side hypermedia* library: every capability terminates in either an SSE event to a connected browser or an `http.Handler` serving JS; the consumer has neither a server nor a browser client.
- It transitively requires go-sse and an HTTP server — i.e. the entire server stack the consumer's roadmap explicitly disavows (ROADMAP.md:67).
- Heaviest of the Group-B four (JS runtime embedded via `static`); it would bloat the wizard binary for zero reachable code.
- The closest plausible hook — "config preview screen before generation" (ROADMAP Theme 4) — is specced as a TUI screen; huh/lipgloss already own that space.
- No TODO_LIST item mentions Datastar, hypermedia, or a dashboard.

## RECONSIDER-WHEN

Only if the Web-UI non-goal is formally revoked and a browser dashboard is planned — at which point go-datastar + go-sse + httputil would be evaluated *together* as a server stack. No intermediate trigger exists for this library in isolation.

## How to verify

```bash
cd /home/lars/projects/GoReleaser-Wizard
rg -c 'larsartmann/go-datastar|net/http' --glob '*.go'   # 0 hits
sed -n '67p' ROADMAP.md                                    # Web UI non-goal
git -C /home/lars/projects/go-datastar tag --list 'v*' --sort=-v:refname | head -1   # v0.7.0
```
