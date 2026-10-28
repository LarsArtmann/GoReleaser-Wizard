# Status Report: Go 1.27 Toolchain Alignment + Quality Gate Repair

**Date**: 2026-10-28 10:10
**Session scope**: Fix the red `buildflow` / `nix-build` pipeline (`go.mod requires go >= 1.27 (running go 1.26.10; GOTOOLCHAIN=local)`), then everything the green gate subsequently surfaced.
**Branch**: master, tree clean. Session commits (auto-commit daemon): `7c448c3`, `2df632b`, `f9d9ee3`, `f502d8d`, `804afa3`, `7819236`, `08a60d3`, `73dbecd`.

**End state (verified)**: `nix flake check` all green (format + build + test, incl. full E2E proof with `goreleaser check` + `actionlint` in the sandbox) · `go test ./...` 6/6 packages pass on Go 1.27.3 · `erraudit lint --type-aware` exit 0 · full `buildflow` run **exit 0** · binary smoke test (`nix run .# -- version`) works.

---

## a) FULLY DONE

| # | Item | Evidence | Scope |
|---|------|----------|-------|
| 1 | **Root cause diagnosed**: go.mod was deliberately bumped to `go 1.27` (commit `2e4dfb6`, 2026-10-25); the flake still used nixpkgs' default Go 1.26.10 with `GOTOOLCHAIN=local`. Decision: align flake UP, not go.mod down. | `git log -- go.mod`; error signature from paste | repo root |
| 2 | **Flake Go pinned to 1.27** via `pkgs.buildGoModule.override { go = pkgs.go_1_27; }`. Key discovery: a plain `go = pkgs.go_1_27;` attr *inside* the buildGoModule set is silently swallowed by mkDerivation — `go` is a callPackage-level arg in current nixpkgs `module.nix`. Verified by `nix eval .#...goreleaser-wizard.go.version` → `1.27.3`. | commit `f9d9ee3` (after failed attr attempt in `7c448c3`) | `flake.nix` |
| 3 | **devShells (default + ci) pinned to `go_1_27`** so local shells stop shadowing with nixpkgs' 1.26. | commit `7c448c3` | `flake.nix` |
| 4 | **vendorHash repaired** via `buildflow -s nix-hash-fix --fix` (go.sum had changed with go.mod). No manual hash pasting. | commit `2df632b`; `nix build` FOD passes | `flake.nix` |
| 5 | **Sandbox test failure eliminated**: current nixpkgs buildGoModule defaults `doCheck = true` on the package itself, so plain `nix build` ran tests without `git` on PATH. Package now sets `doCheck = false`; `checks.test` gained `nativeCheckInputs = [git goreleaser actionlint]`, so the E2E proof (git init → generate → `goreleaser check` → `actionlint`) runs fully instead of skipping. | commits `f502d8d` + doCheck edit in `f9d9ee3`..; `nix flake check` → "all checks passed" | `flake.nix` |
| 6 | **CI cleanup**: removed 6 stale `GOEXPERIMENT: jsonv2` env blocks from `.github/workflows/ci.yml` — json/v2 is the default in Go 1.27 (proven by the GOEXPERIMENT-free nix sandbox build). `go-version-file: go.mod` already yields 1.27. | commit `804afa3` (12 lines removed) | `.github/workflows/ci.yml` |
| 7 | **erraudit findings fixed**: 2× `errors.As` → `errors.AsType[*DomainError]` — `IsErrorCode` (also deleted a dead, duplicated AsType block below the As block) and `addValidationError`. Both unambiguous type-matches (read `.Code` / append typed error), no sentinel regressions. | commit `7819236`; `erraudit lint ./... --type-aware` exit 0 | `internal/domain/errors.go`, `internal/domain/validation.go` |
| 8 | **SECURITY.md rewritten**: now has the required `# Security Policy` header, `## Security Practices`, `## Supported Versions`, and zero unresolved template literals. Removed fiction: `just security-scan`, `./verify.sh`, `validate-strict.sh`, Pro-license vars — none exist in this repo. New content is verifiable-true (SHA-pinned actions, least-privilege permissions, hermetic nix build, secret handling incl. Forgejo `GITEA_TOKEN` note). | commits `08a60d3`, `73dbecd`; securitymd findings 4 → 0 | `SECURITY.md` |
| 9 | **AGENTS.md 414 → 99 lines** (linter max 220): kept all 20 gotchas (compressed), dropped generic style/naming/test boilerplate. Gotcha 11 rewritten (`GOEXPERIMENT` no longer needed — BuildFlow preflight warnings about it are now *correct*); new gotcha 20 documents the `.override` requirement, the `doCheck` default, and the drift symptom signature. | commit `73dbecd`; go-structure-linter finding cleared | `AGENTS.md` |
| 10 | **Full verification loop closed**: `go test ./...` re-run *after* the errors.As migration (TEST-EXIT=0, 6 ok); full `buildflow` EXIT=0 (was: 3 failed / exit 1 at session start); binary smoke-tested. | /tmp/gotest.log; buildflow run log | — |

## b) PARTIALLY DONE

| # | Item | Works now | Open | Effort |
|---|------|-----------|------|--------|
| 1 | **Non-x86_64 flake checks** | x86_64-linux: fully green | `nix flake check` omits aarch64-linux + aarch64-darwin (host is x86_64). Not verified anywhere this session; no CI coverage known. | M |
| 2 | **`nix/package.nix` overlay** | Untouched, still evaluable | Pinned to v0.1.0 with `fetchFromGitHub` + old `vendorHash` (`sha256-CTXzZ...`) — almost certainly stale/broken vs the current module set. Deliberately out of session scope; whether it's even consumed is unknown. | S/M |
| 3 | **`nix-hash-fix` step health** | It DID repair the vendorHash this session (first success in a long while) | Track record 80/82 failures (98%); preflight still warns "vendorHash stale" every run. Root cause of the failure rate not investigated — possible BuildFlow bug or dead-cache environment artifact. | M |
| 4 | **`/mnt/buildcache` broken mount** | Worked around via `/tmp` cache exports + buildflow env-guard rewrites (every run, 4 warnings) | The mount itself is still dead; every Go spawn pays cold-cache + probe latency (preflight-doctor +277% regression). Environment-level fix out of repo scope. | ? |
| 5 | **go-error-modernization skill accuracy** | Used it successfully for the migration judgment | Its documented flags are stale: `erraudit fix` no longer accepts `--type-aware` (binary `61a744c` rejects it; used by other subcommands). Skill not updated; per buildflow's close-the-loop rule this is an unclosed state (fix in skill repo or file issue). | S |
| 6 | **Commit hygiene** | All changes committed (daemon), tree clean | 8 commits all named "chore: auto-commit (heuristic)" — the session's work is not reviewable as logical units (flake fix / CI cleanup / error migration / docs are interleaved). Daemon behavior, not restructured this session. | S (if wanted) |

## c) NOT STARTED

| # | Item | Why not started | Priority |
|---|------|-----------------|----------|
| 1 | **HARVEST section (f) into TODO_LIST.md / ROADMAP.md** | Report only just written; skill says harvest when session continues | High |
| 2 | **CI green verification on GitHub** | Daemon commits locally; whether/how push happens was not established this session | High |
| 3 | **govulncheck in CI** | Noticed gap (devShell has it, CI doesn't); not in scope | Medium |
| 4 | **release.yml review under Go 1.27** | Only confirmed `go-version-file: go.mod`; full release path untested | Medium |
| 5 | **README/FEATURES staleness sweep** (may still claim Go 1.26.5/GOEXPERIMENT) | AGENTS.md was the session's doc scope; broader docs-health VERIFY not started | Medium |
| 6 | **`.buildflow.yml` for this repo** | Repo has none (defaults + pre-commit hook only); skip_steps/env tuning not started | Low |

## d) TOTALLY FUCKED UP

Nothing left broken — final state is fully green — but radical honesty about this session's own failures:

1. **False-positive verification claim (process bug, mine).** Early toolchain test: `go build ... | head -5 && echo "BUILD OK without GOEXPERIMENT"` — the `&&` fired on a *failed* build because piping to `head` masked the exit code. For one step I "verified" the opposite of reality. Caught one command later, but the pattern (piped verification without `PIPESTATUS`/pipefail) is exactly how false claims get recorded.
2. **Wandering diagnosis on the overrideAttrs mystery.** It took ~6 tool cycles (drv hash chasing, derivation env redaction, an arm64-store-go detour, a GOROOT theory) to find the real cause (buildGoModule's `doCheck = args.doCheck or (!buildTestBinaries)` default = true). The signal was there immediately: the *failing drv was `packages...`, not `checks...`* — reading the failing attr name carefully would have halved the hunt. Severity: wasted cycles only; no harm shipped.
3. **Shell env scoping mistake.** `PATH=go1.27.3:$PATH go build ... && erraudit ...` — the PATH prefix applied only to `go build`, so erraudit still ran with 1.26.10 and failed; one wasted round-trip. Same class as (1): env not exported.
4. **Raced the auto-commit daemon blind.** Mid-session drv hashes flip-flopped (dirty ↔ clean evals) because the daemon was committing my edits while I built. I noticed late; checking `git log`/`git status` at the *first* drv-hash inconsistency would have explained it instantly (gotcha 13/daemon was documented; I didn't apply it).
5. **Pre-existing but now-visible debt the session did not create but also did not close**: `nix-build-verify` 12/12 failures (100%), art-dupl warnings in `architecture_test.go` (27-token ×3 clones), timing regressions across the board. These gate nothing today (warnings) but are acknowledged-not-addressed.

## e) WHAT WE SHOULD IMPROVE

1. **Never pipe verification commands without explicit exit-code capture.** Use `cmd > log 2>&1; echo EXIT=$?` or `PIPESTATUS`. (This session produced one false "BUILD OK".) → candidate lesson for `references/lessons.md`.
2. **Read the failing attr name before theorizing.** `packages.x86_64-linux...` failing (not `checks.test`) pointed straight at the package's own `doCheck`. Diagnostic-first, hypothesis-second.
3. **Check git state early in daemon-managed repos.** Any "same file, different drv hash" anomaly → `git log -3` first.
4. **Eval builder signatures before attribute-level overrides.** `nix eval .#pkg.go.version` (or reading `module.nix`) answers "did my override land?" in one command — cheaper than a build cycle.
5. **Re-verify skill-documented CLI flags against the installed binary** (verify-external-claims): the erraudit skill's `fix --type-aware` no longer matches the binary; skills drift silently.
6. **Docs-rot is real**: SECURITY.md had drifted into documenting tools that don't exist in this repo (justfile, verify.sh). Periodic docs-health VERIFY on security/release docs, not just AGENTS.md.
7. **Know the budget before appending to budgeted files**: AGENTS.md is linter-enforced ≤220 lines; I added ~10 lines before learning the constraint.
8. **Test immediately after semantic code edits** even when a fast gate is running: buildflow fast-mode skips tests, so the errors.As migration sat unverified until the end.

## f) Top 50 next tasks (brainstorm — HARVEST fuel, not commitments)

Ranked roughly by impact. Impact / Effort (S <30min, M 30min–2h, L >2h) / Category.

| # | Task | Impact | Effort | Category |
|---|------|--------|--------|----------|
| 1 | HARVEST this report's (f) items into TODO_LIST.md / ROADMAP.md | High | S | Documentation |
| 2 | Verify GitHub CI is green after the GOEXPERIMENT removal push | High | S | Quality |
| 3 | Investigate `nix-build-verify` 100% failure rate (12/12); fix or exclude with rationale | High | M | Bug |
| 4 | Investigate `nix-hash-fix` 98% failure rate; file upstream BuildFlow issue if tool bug | High | M | Bug |
| 5 | Update go-error-modernization skill: `erraudit fix` dropped `--type-aware`; document real flags | Medium | S | Documentation |
| 6 | Decide fate of `nix/package.nix` overlay (v0.1.0, stale hashes): modernize, test, or delete | Medium | S | Cleanup |
| 7 | Add govulncheck step to CI | High | S | Quality |
| 8 | Add aarch64-linux/aarch64-darwin check coverage (CI or remote builder) | Medium | M | Quality |
| 9 | Review + smoke-test release.yml under Go 1.27 end-to-end | High | M | Quality |
| 10 | docs-health VERIFY sweep: README/FEATURES/CHANGELOG for stale Go 1.26.5/GOEXPERIMENT claims | Medium | S | Documentation |
| 11 | Add `nix flake check` to CI so nix regressions fail on GitHub, not just locally | High | S | Quality |
| 12 | Fix or remove `/mnt/buildcache` references from global env; standardize on explicit caches | High | M | Cleanup |
| 13 | Add erraudit lint (`--type-aware`) as a CI lint step so the gate isn't local-only | Medium | S | Quality |
| 14 | Add AGENTS.md ≤220-line size check to CI (currently local go-structure-linter only) | Low | S | Quality |
| 15 | Guard test: assert nixpkgs `go_1_27` version ≥ go.mod `go` directive (catches the next bump early) | High | S | Quality |
| 16 | Guard test: generated workflows' goreleaser pin stays in sync with `types/template_data.go` | Medium | S | Quality |
| 17 | Run + clear on-demand scans: `buildflow -s gitleaks`, `-s codespell`, `-s markdown-lint` | Medium | S | Quality |
| 18 | Fix art-dupl clones in `cmd/goreleaser-wizard/architecture_test.go` (27 tokens ×3) | Low | S | Cleanup |
| 19 | Decide `toolchain go1.27.x` line in go.mod vs relying on nix/flake pinning (standing policy) | Medium | S | Feature |
| 20 | Measure `checks.test` duration after adding goreleaser+actionlint; split e2e vs unit checks if slow | Low | S | Quality |
| 21 | Confirm the in-sandbox `goreleaser check` E2E needs no network (long-term sandbox-proofing) | Medium | S | Quality |
| 22 | Pin flake input nixpkgs to a rev instead of branch head (update deliberately via `nix flake update`) | Medium | S | Cleanup |
| 23 | Add renovate/dependabot for flake.lock + GitHub Action SHAs | Medium | M | Feature |
| 24 | Rebuild erraudit binary (embeds Go 1.26.10 GOROOT → fails package loads on 1.27 modules without PATH tricks) | Medium | S | Bug |
| 25 | Document the erraudit PATH/GOROOT workaround in project AGENTS.md gotchas | Low | S | Documentation |
| 26 | Create minimal `.buildflow.yml` with skip_steps + rationale (e.g. if nix-hash-fix stays red) | Low | S | Cleanup |
| 27 | Re-baseline buildflow timings (preflight-doctor +277%, golangci-lint +146% — cold-cache noise?) | Low | S | Cleanup |
| 28 | Test matrix: `generate` × {github, forgejo, dual, none} × project types, validated by goreleaser check + actionlint | High | M | Quality |
| 29 | Dogfood: regenerate dnsblockd configs with current binary, diff against shipped ones | Medium | S | Quality |
| 30 | Verify goldens cover all 4 ci-provider outputs in CI (`TestGolden` scope check) | Medium | S | Quality |
| 31 | Cross-reference/comment flake.nix vendorHash ↔ nix/package.nix vendorHash (split-brain risk) | Low | S | Documentation |
| 32 | Improve auto-commit daemon messages (batch or semantic) — BuildFlow/meta task | Low | M | Feature |
| 33 | Check `internal/types` comments/docs for leftover GOEXPERIMENT instructions | Low | S | Cleanup |
| 34 | Verify goreleaser release path builds module-proxy artifacts unaffected by nix doCheck change | Low | S | Quality |
| 35 | Add security.txt or PGP contact to SECURITY.md for verifiable reports | Low | S | Documentation |
| 36 | Decide whether devShell should export GOTOOLCHAIN explicitly (host shell pins `local` globally) | Medium | S | Feature |
| 37 | Remove stale `/tmp/gw-go127-test` scratch dir (should be gone; verify) | Low | S | Cleanup |
| 38 | Explore caching the go-modules FOD across flake input updates (check-time) | Low | M | Quality |
| 39 | Consider `buildTestBinaries`/check-phase parallelism tuning for checks.test | Low | M | Quality |
| 40 | Sweep `//nolint` directives for staleness after the errors.As migration (`nolint-audit`, module-scoped) | Low | S | Cleanup |
| 41 | Add CHANGELOG entry for the Go 1.27 + CI-behavior change (consumers of generated configs unaffected, but CI templates changed) | Medium | S | Documentation |
| 42 | Confirm Forgejo goldens byte-stable after toolchain bump (regenerate + diff) | Medium | S | Quality |
| 43 | Add flake-check gate for treefmt covering SECURITY.md/AGENTS.md (prettier/markdownlint opt-in) | Low | S | Quality |
| 44 | Investigate flake "incompatible systems omitted" warning — encode expected systems explicitly | Low | S | Cleanup |
| 45 | Check whether any consumer docs reference `GOEXPERIMENT` setup steps for this repo | Low | S | Documentation |
| 46 | Evaluate Go 1.27 stdlib wins for the codebase (e.g. jsonv2 adoption beyond internal/types) | Low | M | Feature |
| 47 | Wire `nix run .#test` and `nix run .#lint` apps if referenced anywhere (AGENTS.md mentions them) | Low | S | Cleanup |
| 48 | Schedule `nix flake update` + `buildflow update` cadence (monthly?) | Low | S | Cleanup |
| 49 | Confirm flake comment about nixpkgs 26.11 dropping x86_64-darwin is still accurate | Low | S | Documentation |
| 50 | Split-brain sweep: any other place pinning a Go version (CI, docs, goreleaser config, devshell) besides go.mod + flake | Medium | S | Cleanup |

## g) Questions I cannot answer myself

1. **`/mnt/buildcache`: fix it or amputate it?** Every run pays 4 env-guard rewrites + cold-cache probe failures (preflight-doctor +277%). Is this mount coming back (then I leave configs alone), or is it dead permanently (then I strip its references from global/project env and standardize on explicit cache paths)? I can't see the system-level mount plan.
2. **`nix/package.nix` overlay: who consumes it?** It's pinned to v0.1.0 with a GitHub `fetchFromGitHub` source and a stale vendorHash. Is `flake.overlays.default` used by another flake/consumer (then it needs modernizing + a test), or is it dead weight (then I delete it)? Grepping this repo alone can't answer cross-repo consumption.
3. **Commit history policy on master:** the daemon produced 8 heuristic "chore: auto-commit" commits for this session's work. Before the next push — should I restructure these into semantic commits (build: / ci: / fix: / docs:), or is raw auto-commit history acceptable on master? This determines whether the Go 1.27 migration is reviewable as a unit.

---

*Per the status-report skill: section (f) is the primary input for `docs-health` HARVEST — it belongs in TODO_LIST.md/ROADMAP.md, not entombed here. HARVEST not yet run this session (awaiting instructions). Report left uncommitted deliberately (harness rule: no commits without explicit request); the auto-commit daemon will pick it up.*
