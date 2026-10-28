# Forgejo Migration Execution — Session Report (wizard T6–T10, dnsblockd pilot, SystemNix shadow)

**Date**: 2026-10-28 ~02:30
**Session goal**: execute the master plan autonomously (`docs/planning/2026-10-12_15-17_forgejo-primary-migration-master-plan.md`); privileged steps accumulate in `docs/planning/forgejo-migration-privileged-batch.md`.

## What shipped (all verified, all committed)

### 1. goreleaser-wizard: Forgejo CI provider (T6–T10)

- **Domain (T6)**: `CIProvider` enum (github/forgejo/dual/none) in `internal/domain/ci_provider.go`; unset behaves as github (`isGitHubLike`) so legacy configs are untouched; `ShouldGenerateForgejoWorkflows()` on `SafeProjectConfig`. Table tests in `ci_provider_test.go`.
- **Templates (T7/T8)**: `ForgejoCITemplate` + `ForgejoReleaseTemplate` (`[[ ]]` delims): BuildFlow full gate pinned to `github:LarsArtmann/BuildFlow/49ced3e0…`, goreleaser pinned to nixpkgs `3ec6363…#goreleaser` (v2.17.1), `runs-on: native`, release step uses `GITEA_TOKEN: ${{ secrets.GITHUB_TOKEN }}` (Forgejo's own token) and an idempotent release-shadow step that skips without `RELEASE_SHADOW_GITHUB_TOKEN`.
- **T9**: `GoReleaserTemplate` renders `gitea_urls` + `changelog.use: git` when `ReleasesToForgejo()`; goldens `goreleaser-forgejo.yaml` + `goreleaser-dual.yaml` pin the shapes; github output byte-identical to before (golden diff = one blank line).
- **Generator + wiring**: `generators/forgejo_workflows.go` mirrors the GitHubActionsGenerator API (both files, rollback, preview, GetWorkflowPaths); job + factory gating (`forgejo` → no GitHub job; `dual` → both; `none` → neither) + `generationTargets` preflight; `--ci-provider` flag in `addCommonFlags`, applied in `runGenerate` (which now `ApplyDefaults()` after overrides so job decisions see `ActionsOn`/`CIProvider`).
- **T10 verification**: `go test ./...` green (incl. new E2E `TestForgejoWorkflowsPassActionLint` with `-config-file` native-label actionlint), `golangci-lint run ./...` 0 issues, `go-arch-lint check` clean, goldens regenerated, live CLI smoke: forgejo/library path end-to-end + `goreleaser check` zero deprecations + actionlint clean; commit `9334dc3` (+ daemon absorb commits).

### 2. dnsblockd pilot content (T1 autonomous half)

- Generated via the wizard binary itself (`--ci-provider forgejo --project-type library`), pushed to GitHub master `411dce43`: `.forgejo/workflows/ci.yml`, `.forgejo/workflows/release.yml`, `.goreleaser.yaml` (gitea_urls → forgejo.home.lan, `use: git`). Verified: `goreleaser check` clean, actionlint clean.

### 3. SystemNix: `forgejo-release-shadow` + pilot staging (T2/T5)

- `forgejo-release-shadow` (writeShellApplication, PATH binary): copies Forgejo release notes+assets to the GitHub twin; empty `GITHUB_TOKEN` → exit 0; idempotent (≥ asset count = no-op); deficient release → loud refuse; missing tag/release → loud error. `FORGEJO_URL/OWNER/REPO`, `GITHUB_OWNER/REPO`, `FORGEJO_TAG` env contract matches the wizard's release template exactly.
- 5 fixture cases added to `forgejo-scripts-fixture` (skip / happy / already-mirrored / deficient / missing-release) — **check green** (`nix build .#checks.x86_64-linux.forgejo-scripts-fixture`).
- `canonicalRepos += "dnsblockd"` in `platforms/nixos/system/configuration.nix` (evaluated via `nix eval` on `nixosConfigurations.evo-x2`); runbook `docs/services/forgejo.md` documents the shadow + pilot; committed (`4ddc9a84`, `a1ca8384`), `nix fmt` clean.

### 4. Mirror-Actions experiment — ANSWERED (was recon risk)

Unauthenticated API sweep of the instance (2026-10-28): every visible mirror repo carries `has_actions: false` (incl. `lars/GoReleaser-Wizard`, `mirror_interval: 2h`). Conclusion: the 2h pull-sync of mirror repos does NOT trigger Actions runs; Actions arm at flip time (native + has_actions). Positive control = first push to dnsblockd after Phase B of the privileged batch. Private mirrors (dnsblockd, monitor365, wireguard-collector) remain unauth-invisible; authenticated confirmation is a Phase D one-liner.

## Pre-existing bugs found & fixed (root cause, each verified)

1. **Baseline red (3 tests)**: empty `CIProvider` broke `GeneratesGitHubWorkflows()` for all unset-provider configs → unset now github-like; E2E got a docker-presence skip mirroring its goreleaser skip.
2. **`main-push` invalid**: recommended default trigger had no `actionTriggerMetaMap` entry → `IsValid()` false on ApplyDefaults' own recommendation. Added the meta entry.
3. **Library arch recommendation**: `ProjectTypeLibrary.RecommendedArchitectures()` included 386, incompatible with its recommended darwin → every library config failed its own platform/arch validation. 386 removed with rationale.
4. **Docker-registry validation**: unconditional `IsValid()` rejected no-docker configs (empty registry). Now required only when `DockerSupport.IsEnabled()`.
5. **SystemNix flip fixture**: post-migrate stub body lacked `open_pr_counter` (gogenfilter-era verify arithmetic); jq `"null"` + `set -u` arithmetic died ("null: unbound variable") — fixture body updated (script was right, fixture drifted).
6. **`sed 's/{.*//'` on upload_url**: GNU sed brace pitfall → bash `${var%%\{*}` expansion.

## Known-left (not mine to chase today)

- `-update-golden` flag leaks into TestInitCommand's CLI invocation (update-mode-only failure, pre-existing).
- CLI docker dependency check consults pre-override DockerSupport (documented as gotcha 18; harmless in tests, needs a docker shim for host smoke tests).
- Forgejo SSH push auth still broken server-side until the Phase A deploy re-runs `forgejo-ssh-keys.service`.

## Next actions

1. User: run `docs/planning/forgejo-migration-privileged-batch.md` Phases A–E (deploy → flip pilot → shadow secret → live verify → optional pass handoff).
2. Agent (next session): T2 `forgejo-flip-batch@` design, T22 branch protection script, T16 runner capacity measurement after pilot load, GoReleaser-Wizard release (v2 tag ships `--ci-provider`), then wave planning per master plan §6.
