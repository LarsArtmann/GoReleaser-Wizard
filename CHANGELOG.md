# Changelog

All notable changes to GoReleaser-Wizard are documented here. Newest first.
Format loosely follows [Keep a Changelog](https://keepachangelog.com); each entry cites its commit.

## [Unreleased]

### Docs (2026-11-07)

- docs-health AUDIT: created `TODO_LIST.md` + `ROADMAP.md`, refreshed `FEATURES.md`/`README.md`/`AGENTS.md` to the post-Forgejo/Go-1.27 reality, annotated and archived the 2026-01…08 status corpus to `docs/status/archived/` (this sweep)

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
