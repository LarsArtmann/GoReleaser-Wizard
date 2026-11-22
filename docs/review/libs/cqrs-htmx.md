# cqrs-htmx — Library Adoption Audit

| Field | Value |
| --- | --- |
| Library | `github.com/larsartmann/cqrs-htmx/v4` (+ submodules: usermgmt `v4.21.0`, adminui `v4.18.0`, setup `v4.18.0`, loginpage `v4.19.0`, auditlog `v4.13.0`, health `v4.15.0`, datastar `v4.13.0`, systemadapter `v4.14.0`, huma `v4.1.0`, usermgmt/oauth2 `v4.13.1`, usermgmt/totp `v4.12.1`, usermgmt/webauthn `v4.13.0`) |
| Local checkout | `/home/lars/projects/cqrs-htmx` @ `d4e14b73` (2026-11-22) |
| Latest local tag | root `v4.19.0` |
| Current state (README) | Current major is **v4**; the tree is already SSE-only — WebSocket transport was removed in v5 per ADR 0046 (cqrs-htmx/README.md:571, :440); v5 also deletes deprecated re-export aliases (cqrs-htmx/README.md:391, :573) |
| Consumer | GoReleaser-Wizard @ `a493950` (2026-11-22) |
| Used? | **No** — `rg 'larsartmann/cqrs-htmx'` → 0 hits |
| Verdict | NOT-ADOPTED → **STAY-OUT** |
| Score | n/a — 0 applicable capabilities (7 rows, all N/A) |
| Provenance | Fleet audit Group E, 2026-11-22; citations QC-sampled by orchestrator |

## What the library offers

A CQRS/HTMX web-application layer over go-cqrs-lite: a typed command/query HTTP pipeline (`App.Command`/`App.Query`), an event-sourced `usermgmt.Service` (register/authenticate/logout, roles/tenants/bots, WebAuthn/OAuth2/TOTP, Casbin authz), a ready-made templ+HTMX admin panel, a login page, an SSE broadcaster, and a `setup` one-call bundle composing all of it.

| Capability | What it does | Citation | Applicable here? |
| --- | --- | --- | --- |
| App builder | `Config{Commands,Queries}` → typed `http.HandlerFunc` pipeline | cqrs-htmx/app.go:148 (`New`), :195 (`MustNew`) | N/A — consumer has zero `net/http` |
| App.Command / App.Query endpoints | Declarative POST→command→201 and query endpoints | cqrs-htmx/app.go:243, :270 | N/A — no commands, no queries, no HTTP surface |
| usermgmt.Service | Event-sourced register/authenticate/logout, roles/tenants | cqrs-htmx/usermgmt/service_factory.go:85, cqrs-htmx/usermgmt/service_register.go:68, cqrs-htmx/usermgmt/service_login.go:19 | N/A — no users, no auth, no sessions; the tool runs as the local user's CLI |
| adminui panel | Templ+HTMX admin dashboard, role-gated | cqrs-htmx/adminui/handler.go:33, authz.go:30 | N/A — Web UI is an explicit non-goal (ROADMAP.md:67); nothing to administer |
| setup one-call bundle | Composes event infra + usermgmt + panels, mounts routes | cqrs-htmx/setup/setup.go:30, cqrs-htmx/setup/mount.go:35 | N/A — runtime HTTP wiring; **generates no files** (see verdict) |
| SSE broadcaster | Fan-out live events to browsers | cqrs-htmx/sse_broadcaster.go:54 | N/A — no browser, no long-lived process |
| Transport posture | v4 current; SSE-only tree; v5 codemod planned | cqrs-htmx/README.md:571 | N/A — informational; confirms the library targets server-rendered web apps |

## What GoReleaser-Wizard does today

Single-shot Cobra CLI: no DB (only CGO build-tag name strings at internal/domain/build_tag.go:97, :101), no `net/http` (0 hits), no web UI, no users/auth. The adjacent "engine" is the Job/Workflow list (≤7 in-memory generation jobs, compensating-delete rollback, preflight-gated atomicity — job_manager.go:280-305, jobs.go:439-503). ROADMAP.md:67 states verbatim: "**Web UI**: the CLI/TUI is the product; a web layer duplicates it without a consumer". TODO_LIST.md contains no persistence, event-sourcing, projection, user-management, or admin-UI plans.

## Why STAY-OUT

- **The setup-codegen temptation is real but wrong**: `setup.New` (setup/setup.go:30) is a *runtime composition root* — it builds event stores, a usermgmt.Service, three UI panels and `Mount`s them on an `http.ServeMux` (setup/mount.go:35). It generates zero files; it requires exactly the three things this consumer lacks (net/http, event sourcing, users). The wizard's codegen is template rendering — `go:embed` + text/template + golden tests — already solved without this library.
- Every consumer-relevant submodule bottoms out in a web-server model: App/handlers (app.go:243, :270), admin UI (adminui/handler.go:33), auth (usermgmt/service_factory.go:85) — all presuppose a long-lived HTTP process.
- Pulling cqrs-htmx transitively drags go-cqrs-lite (and its proprietary license) — the go-cqrs-lite STAY-OUT applies doubly here.
- No roadmap item asks for anything this library provides; adopting would be dependency weight with zero callable code paths.

## RECONSIDER-WHEN

Flip only if ROADMAP.md:67's Web UI non-goal is formally reversed **and** the product grows a hosted/service form with accounts — a web-delivered wizard needing user management, sessions, and an admin panel. Either condition alone is insufficient (a web UI without accounts needs no usermgmt; accounts without a web UI have no mount point). Until the non-goal line is edited out of ROADMAP.md, never.

## How to verify

```bash
cd /home/lars/projects/GoReleaser-Wizard
rg -c 'larsartmann/cqrs-htmx' .                     # 0 hits
sed -n '67p' ROADMAP.md                              # Web UI non-goal (verbatim)
rg -c 'net/http' --glob '*.go' --glob '!*_test.go'   # 0 hits
git -C /home/lars/projects/cqrs-htmx tag --list 'v*' --sort=-v:refname | head -1   # v4.19.0
sed -n '571p' /home/lars/projects/cqrs-htmx/README.md # v4 current, SSE-only tree
```
