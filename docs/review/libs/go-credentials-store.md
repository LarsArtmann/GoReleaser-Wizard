# go-credentials-store — Library Adoption Audit

| Field            | Value                                                                 |
| ---------------- | --------------------------------------------------------------------- |
| Library          | `github.com/larsartmann/go-credentials-store`                         |
| Local checkout   | `/home/lars/projects/go-credentials-store` @ `16ab83f` (2026-11-20)   |
| Latest local tag | `v0.4.3`                                                              |
| Consumer         | GoReleaser-Wizard @ `a493950` (2026-11-22)                            |
| Used?            | **No** — `rg 'larsartmann/go-credentials-store'` → 0 hits             |
| Verdict          | NOT-ADOPTED → **STAY-OUT**                                            |
| Score            | n/a — 0 applicable capabilities (6 rows, all N/A)                     |
| Provenance       | Fleet audit Group C, 2026-11-22; citations QC-sampled by orchestrator |

## What the library offers

An encrypted credential store SDK: XChaCha20-Poly1305 sealing with KeyID rotation, an event-sourced lifecycle (created/rotated/probed/used/revealed), a `Resolver` API, master-key management with file-mode guards, health hooks, and an admin UI shim mounted behind host auth.

| Capability            | What it does                                    | Citation                                                 | Applicable here?                                         |
| --------------------- | ----------------------------------------------- | -------------------------------------------------------- | -------------------------------------------------------- |
| Encrypted store SDK   | `Build(ctx, Options)` composition root          | go-credentials-store/container.go:124, :28               | N/A — consumer never stores or resolves real credentials |
| Resolve API           | `Resolver.Resolve(ctx, platform, kind)`         | go-credentials-store/resolve.go:35, :22                  | N/A                                                      |
| Master-key management | Env/file key loading, rotation, file-mode guard | go-credentials-store/keys.go:22, :78, :120               | N/A — no keys to manage                                  |
| Health check hook     | Per-subsystem health map                        | go-credentials-store/health.go:24                        | N/A                                                      |
| Admin UI shim         | HTML UI + JSON API + SSE + audit trail          | go-credentials-store/adminui/adminui.go; README.md:71-77 | N/A                                                      |
| Lifecycle audit       | Event-sourced trail                             | go-credentials-store/README.md:4-6                       | N/A                                                      |

Dependency weight is the fleet's heaviest: 12 `go-cqrs-lite/*` modules, templ-components ×4, modernc.org/sqlite, prometheus, samber/do, cobra, **and go-daemon v0.4.0** (go-credentials-store/go.mod:5-33).

## What GoReleaser-Wizard does today — every token/secret hit, classified

`rg -i 'token|credential|secret' --glob '*.go' --glob '!*_test.go'` → 40+ hits, all in four classes; **none store or handle real secret material**:

1. **CI template literals** (writes _references_, never values) — cmd/goreleaser-wizard/templates/embedded.go: `${{ secrets.DOCKER_USERNAME/PASSWORD }}` (:152-153), `${{ secrets.GITHUB_TOKEN }}` (:164, :238), `${{ secrets.RELEASE_SHADOW_GITHUB_TOKEN }}` (:244), `id-token: write` (:128). The wizard binary never sees a value.
2. **Unimplemented domain interface declarations** — internal/domain/interfaces.go:116-118 (`ValidateToken`/`GetUser`/`CheckPermissions`), :141/:143 (`Login`/`ValidateCredentials`), :444-448 (`DockerCredentials` struct). No non-test implementation exists.
3. **Provider metadata** — internal/domain/git_provider.go:127-131 `RequiresPersonalAccessToken()` (enum flag); enums_actions.go:114-119 (permission scopes emitted into workflows).
4. **Escape-analysis data** — internal/validation/template_escaping.go:324 (scans for `"$GITHUB_"`, `secrets.` to decide template escaping).

Plus golden testdata with the same `${{ secrets.* }}` literals, and TODO_LIST.md:45,47 mentioning operator-held Forgejo API tokens for _external_ scripts (blocked on the owner, not wizard state).

## Why STAY-OUT

- **There is no credential to store**: the product contract is generating configs that _reference_ `secrets.GITHUB_TOKEN` etc., resolved by CI at runtime.
- Pulling it in transitively adds the go-cqrs-lite stack, SQLite, an event store, and go-daemon itself — catastrophic dependency bloat for a zero-requirement.
- Security posture would _worsen_: introducing a master-key file and encrypted store into a tool that today touches no secret material creates a new attack surface (its threat model assumes a server deployment).
- No roadmap/TODO item plans login, keychain, or token persistence.

## RECONSIDER-WHEN

The wizard starts performing releases itself (a hypothetical `wizard release` command) and must hold real PATs/registry credentials across runs — only then does encrypted-at-rest storage with audit become a genuine requirement.

## How to verify

```bash
cd /home/lars/projects/GoReleaser-Wizard
rg -n -i 'token|credential|secret' --glob '*.go' --glob '!*_test.go' | head -30   # all four classes above
rg -c 'larsartmann/go-credentials-store' .                                       # 0 hits
git -C /home/lars/projects/go-credentials-store tag --list 'v*' --sort=-v:refname | head -1   # v0.4.3
```
