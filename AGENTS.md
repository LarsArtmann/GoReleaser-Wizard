# GoReleaser-Wizard Agent Guide

**Project**: GoReleaser-Wizard — interactive GoReleaser configuration wizard
**Repository**: https://github.com/LarsArtmann/GoReleaser-Wizard

Only non-obvious, session-critical context lives here. Full history in git; details in README.md / FEATURES.md / docs/.

## Overview

Interactive CLI that generates production-ready GoReleaser configurations with CI workflows (github/forgejo/dual targets). Clean Architecture + DDD:

1. **Domain** (`internal/domain/`): pure business logic, zero external deps. Core: `SafeProjectConfig`, typed enums, `DomainError` system (`internal/errors/domain_errors.go`)
2. **Application** (`cmd/goreleaser-wizard/`): Cobra CLI, workflow/job execution system, typed generators (`generators/`)
3. **Infrastructure** (`internal/git/`, `internal/utils/`): git wrappers, utilities
4. **Templates** (`cmd/goreleaser-wizard/templates/`, embedded in `templates/embedded.go`): THE single source of truth for generated files; rendered only by generators. `types/template_data.go` holds typed template data + owner/repo resolution (flag > git remote > placeholder)

**Stack**: Go 1.27, Cobra, Viper, Charm (lipgloss/log), golangci-lint, go-arch-lint, GoReleaser v2.17 as the generation target.

## Essential Commands

Export these prefixes for every Go command while `/mnt/buildcache` is broken (gotcha 12):

```bash
export GOCACHE=/tmp/gw-cache GOMODCACHE=/tmp/gw-modcache GOLANGCI_LINT_CACHE=/tmp/gw-lint-cache

go build ./...
go test ./...                            # includes E2E + golden tests
go test ./cmd/goreleaser-wizard -update-golden   # regenerate goldens (see gotcha 14b)
golangci-lint run ./...
go-arch-lint check
```

Verify the product (generated configs) in a fresh Go module with git init:

```bash
go build -o /tmp/gw-bin ./cmd/goreleaser-wizard
/tmp/gw-bin generate --config <cfg.yaml> --github-owner <owner> --github-repo <repo>
goreleaser check                        # must exit 0 with zero deprecations
```

The E2E test `TestGeneratedConfigPassesGoReleaserCheck` automates this; skips gracefully when `goreleaser` is not in PATH.

- `nix flake check`, `nix build`, `nix run .#test`, `nix run .#lint` — never Makefile/justfile
- BuildFlow pre-commit hook: ~60 tools, auto-fixes files — check `git status` after every commit

## Architecture Enforcement (go-arch-lint)

- Domain cannot import application/infrastructure; deep AST scanning enabled
- `*_test.go` excluded from architecture validation — don't rely on this for production code
- Cycles are caught by arch-lint and compilation; break them with DI

## Important Gotchas

**1. Domain layer purity** — domain MUST NOT import application/infrastructure; test it in isolation.

**2. Always use typed domain errors** — never `errors.New`/`fmt.Errorf` directly; use the `errors.New*Error()` constructors. Propagate `context.Context` through async ops.

**3. Template escaping** — Go template delimiters need `{{"{{"}}` escapes; see `internal/validation/template_escaping.go`. Forgejo templates use `[[ ]]` delims so `${{ secrets.* }}` literals need no escaping.

**4. `os.WriteFile` with explicit perms** (`0o644`); always create backups before overwriting existing files (generators' `createBackup()`).

**5. Nix first** — no justfiles, no Makefile.

**11. Go 1.27: `GOEXPERIMENT=jsonv2` no longer needed (2026-10-28)** — `internal/types` imports `encoding/json/v2`; since go.mod bumped to 1.27, jsonv2 is the default. CI no longer sets the flag; a user-level `go env -w GOEXPERIMENT=jsonv2` is harmless.

**12. `/mnt/buildcache` may be broken** — symptom: `failed to initialize build cache at /mnt/buildcache/...`. Fix: the GOCACHE/GOMODCACHE/GOLANGCI_LINT_CACHE exports above.

**13. BuildFlow pre-commit can hang on telemetry** — PostHog TLS timeouts block the commit after the pipeline went green. Recovery: kill buildflow/git-commit, verify log shows `0 failed`, commit with `--no-verify` and note why.

**14. E2E testing gotcha** — `go run github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard` resolves the PUBLISHED module from the proxy, not local code. Always `go build -o /tmp/gw-bin` first.

**14b. Golden regeneration flag leaks** — passing `-update-golden` to `go test` makes `TestInitCommand` fail (CLI harness inherits the flag). Tests are green without it.

**15. Domain validation rejects windows+arm64** — test configs must not combine them (`architecture_compat.go`).

**16. `--ci-provider {github,forgejo,dual,none}` selects workflow + release targets** (2026-10-28)

- Default github = byte-identical legacy output; unset CIProvider behaves as github everywhere
- Forgejo templates: `ForgejoCITemplate`/`ForgejoReleaseTemplate` in `templates/embedded.go`; generator `generators/forgejo_workflows.go` writes both files
- Forgejo releases: GoReleaser config renders `gitea_urls` + `changelog.use: git`; inside Forgejo Actions `secrets.GITHUB_TOKEN` IS the Forgejo API token (`GITEA_TOKEN`)
- Release shadow: `forgejo-release-shadow` (SystemNix) copies assets to GitHub; skips when `RELEASE_SHADOW_GITHUB_TOKEN` is absent
- Pins live in `types/template_data.go`: BuildFlow flake rev + nixpkgs goreleaser v2.17.1 — bump deliberately, never float
- Goldens: `forgejo-workflows.yml`, `goreleaser-forgejo.yaml`, `goreleaser-dual.yaml` (dual = both workflow sets, goreleaser still targets Forgejo)

**17. actionlint + self-hosted runner labels** — generated workflows use `runs-on: native`; actionlint rejects it without a config:
`actionlint -config-file <(printf 'self-hosted-runner:\n  labels:\n    - native\n') ...` (the E2E test does exactly this)

**18. CLI defaults are applied twice (quirk, not a bug to chase)** — `runGenerate` applies `ApplyDefaults()` after flag overrides; `ConfigGenerationJob` re-applies harmlessly. The docker dependency check still consults pre-override DockerSupport on the CLI path — pass a `docker` shim on hosts without docker when smoke-testing.

**19. Dogfood proof** — dnsblockd's `.forgejo/workflows` + `.goreleaser.yaml` were generated by this binary (`--ci-provider forgejo --project-type library`) and validated with goreleaser check + actionlint before push.

**20. Flake Go pin: use `buildGoModule.override`, not a `go` attr (2026-10-28)**

- go.mod requires 1.27; nixpkgs default `go` trails (1.26) and buildGoModule sets `GOTOOLCHAIN=local` → `go: go.mod requires go >= 1.27`
- A `go = pkgs.go_1_27;` attr INSIDE the buildGoModule set is silently swallowed by mkDerivation (callPackage-level arg); working form: `pkgs.buildGoModule.override { go = pkgs.go_1_27; } { ... }`
- Current nixpkgs buildGoModule defaults `doCheck = true` on the package: plain `nix build .#goreleaser-wizard` RUNS tests and fails in the sandbox (`git` not on PATH). Flake sets `doCheck = false` on the package; tests live in `checks.test` with `git`/`goreleaser`/`actionlint` as `nativeCheckInputs` (E2E proof runs fully there instead of skipping)
- Symptom when these drift: `nix flake check` fails on `packages...` (not `checks.test`) with `git init failed: exec: "git": executable file not found in $PATH`

> This file was auto-trimmed. Full history in git.
