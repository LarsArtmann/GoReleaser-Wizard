# Changelog

All notable changes to GoReleaser-Wizard are documented here. Newest first.
Format loosely follows [Keep a Changelog](https://keepachangelog.com); each entry cites its commit.

## [Unreleased]

### Hardening wave (2026-11-07)

Wave-6 of `docs/planning/2026-11-07_08-51_SUPERB-wave6-hardening-and-proof-plan.md`: 24 of 27 tasks executed; full status in `docs/status/2026-11-07_18-36_wave6-execution-24of27-tasks-status.md`.

**Fixed** (real bugs found by the hardening passes):

- Explicit `--config` files had no effect: loading early-returned before env vars, flags, and unmarshal were applied, and `GORELEASER_WIZARD_*` env vars mapped to dead keys because the prefix was never stripped (`260054f`)
- `ssh://git@github.com:22/owner/repo` parsed the port "22" as the repo owner; numeric remote-URL segments are now dropped (`bf2d3a1`)
- Five domain validators returned typed-nil `*DomainError`, making nil a non-nil error interface and silently invalidating valid configs; they return plain `error` now (`d38540d`)
- `extractBinaryFromPath` returned "main" instead of the binary directory name for `src/<bin>/main.go` layouts (`9c3fe88`)

**Added**:

- `validate` detects deprecated GoReleaser v2.17 keys (`dockers`, `docker_manifests`, `brews`, `archives.format_overrides.format`) with line numbers and before→after migration hints, empirically verified against goreleaser 2.17.1 (`b198b32`)
- Forgejo generations write `.forgejo/actionlint.yaml` (self-hosted `native` label) next to the workflows; the E2E gate actionlints the GENERATED config (`7bbf43c`)
- Provider-aware success tips: `--ci-provider forgejo/dual` points users at `.forgejo/workflows`, not `.github/workflows` (`66f1997`)
- Standing CI: golden-churn guard (`-update-golden` must leave `testdata/` clean), 5-provider generate matrix, flake go-pin ≥ go.mod guard, docs job running `cmd/docs-check` (`61d7efa`, `b160050`)
- `cmd/forgejo-preflight` (local flip-readiness gates) + `forgejo-flip-batch@.service` (one-sudo-per-wave privileged batch flipper) (`ae56dd7`)
- `scripts/docs-gates.sh`: dprint → docs-gates battery (advisory) → docs-check composition (`0308d5b`)
- Coverage lifts: domain 15.5%→48.7%, config 84.1%, git 0%→45.5%, validation registry extraction; parse-proof tests for all 7 embedded templates (`d38540d`, `260054f`, `bf2d3a1`, `4a153ab`, `b87636e`, `506129e`)

**Docs**:

- Root-sprawl disposition proposal: decision-ready inventory of the 19 root markdown files and 7 `github_issues/` drafts with archive/delete/harvest plan; execution owner-gated (`0f2b48f`)
- TODO_LIST synced with the wave-6 completions (17 completed rows removed); nix-flake-check-in-CI and `.docs-gates.json` seeding documented as deliberate non-goals in `ci.yml` and the wave-6 status report (`b706d1a`, `b198b32` era)

### Docs (2026-11-07)

- docs-health AUDIT: created `TODO_LIST.md` + `ROADMAP.md`, refreshed `FEATURES.md`/`README.md`/`AGENTS.md` to the post-Forgejo/Go-1.27 reality, annotated and archived the 2026-01…08 status corpus to `docs/status/archived/` (`2fcc8a2`, sweep commit)

### Added

- `.forgejo/actionlint.yaml`-ready workflows: generated Forgejo CI/release workflows pin the BuildFlow flake rev and goreleaser v2.17.1; `--ci-provider {github,forgejo,dual,none}` selects workflow + release targets (`9334dc3`, 2026-10-28)
- E2E proof harness: generate into a temp module → `goreleaser check` must exit 0 with zero deprecations; golden files lock every artifact variant; `-update-golden` regeneration flag (`8edf0b1`, 2026-08-17)
- Atomic generation: preflight checks all target files before writing anything; `--force` to overwrite (`6f4f54d`, 2026-08-17)
- Dockerfile generation whenever Docker is enabled (goreleaser-prebuilt `FROM scratch`/`alpine` pattern, `ARG TARGETPLATFORM`), wired into both full-wizard and config-only paths (`e85bb4f` + `6f4f54d`, 2026-08-16/17)
- `docs/goreleaser-guide.md`: GoReleaser v2.17 knowledge base (phases, `dockers_v2`, casks, template vars, snapshot vs release) (`969259e`, 2026-08-17)

### Changed

- Templates modernized to GoReleaser v2.17: `dockers_v2`, `archives.format_overrides.formats`, no deprecated keys; generated config checks clean with zero deprecations and no env vars (`e85bb4f`/`3282680`, 2026-08-16)
- Generators package is the single template source of truth; inline `jobs.go` template copies and `map[string]any` data prep deleted (`9f4bd71`, 2026-08-17)
- Own `.goreleaser.yaml` migrated to v2.17 (`dockers_v2`, `homebrew_casks`, cosign `--bundle`); snapshot release verified incl. multi-arch Docker (`969259e`, 2026-08-17)
- GitHub Actions release workflow: portable `ubuntu-latest` default, buildx setup, conditional Docker-login/cosign steps (`3282680`, 2026-08-16)
- Go toolchain: go.mod on 1.27 (`encoding/json/v2` is default, no `GOEXPERIMENT` needed); flake pins `go_1_27` via `buildGoModule.override`; CI no longer sets GOEXPERIMENT (2026-10-28 toolchain-alignment session, daemon commits `f9d9ee3`…`d0f7529`)
- GitHub Actions pinned to commit SHAs + dprint config (`3822656`, 2026-08-16)
- golangci-lint: 52 issues → 0 (slice allocation + godoc fixes `9703270`, 2026-08-10); config converged via auto-configure (`3242bca`, 2026-11-02)

### Fixed

- Test-suite flags no longer leak into the CLI harness (`-update-golden` etc.); goldens byte-identical (`27c037e`, 2026-11-04)
- Stale backups removed (3 `.orig` files), unused `createTempFile` deleted, broken `test-wizard/` scratch module removed, `infertypeargs` hints fixed (`fa648db` + `3cb91da`, 2026-07-28 / 2026-08-17)
- Broken `hasFile` helper restored (non-compiling `slices.Contains(dirs, nil)` rewrite) and filesystem I/O removed from the domain layer (`41aaf86`, 2026-06-30)
- arch-lint stale component refs cleaned after the errors→domain migration (`c6d0ec1`, 2026-05-28)
- Go template delimiter escaping + Docker registry double-dot rejection (`23889fb`/`29277e6` + `9bae006`, 2026-03-21)

### Migration history (condensed)

- 2026-10-17: external review (jdx/packslip lens) captured; 2026-10-12: Forgejo-primary fleet migration master plan (`e04ceb2`)
- 2026-09-29: nix real src/vendor hashes + sandbox-runnable tests in overlay (`27954ee`); 2026-09-04: Go 1.26.7 + Charm stack bump (`ab5f81e`)
- 2026-07-28: deduplication session — `getBoolFlag`/`getStringFlag` helpers replace 11 error-discarding flag reads (`22d656f`)
- 2026-05-25: error codes migrated `internal/errors` → `internal/domain` (`dde1081`); 2026-05-17: GoReleaser config v2 overhaul series (nfpm/Scoop/nixfmt/cosign/GHCR)
- 2026-04-30: `go-branded-id` migration (`cc8a753`); 2026-03-26/25: charm.land v2 migration incl. `huh/v2` TUI (`c99a279`, `52a5a96`)
- 2026-03-21: viper → koanf config migration (`50c0feb`); 2026-03-20: huh-based TUI wizard (`7e85f9e`) + GitHub Actions release workflow (`1b17127`)
- 2026-01-25: initial GoReleaser-Wizard implementation (`a10e28f`, `4fe717e`)

## Prehistory (2025, import order as originally logged)

- fix: critical GitHub Actions workflow issues (f08ea45)
- feat: Enhance verification system with comprehensive validation (e8c75de)
- feat: Fix and optimize GoReleaser configurations (017f382)
- feat: Complete license system dogfooding with EUPL-1.2 (93f0b89)
- feat: Add automated license generation system (79f93c8)
- feat: Add comprehensive README.md and fix readme-generator configuration (dc70cf0)
- feat: Create comprehensive GoReleaser template with free and pro versions (63b61df)
