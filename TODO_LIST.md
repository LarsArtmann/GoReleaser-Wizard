# TODO List

> Short-term, actionable, bounded work items, verified against the actual code on 2026-11-07.
> Long-term vision lives in [ROADMAP.md](ROADMAP.md). Completed items are removed and logged in [CHANGELOG.md](CHANGELOG.md).
> Primary sources: the Wave-6 remainder of `docs/planning/2026-08-16_18-45_SUPERB-WIZARD-MODERNIZATION.md` (M17-M27), the Forgejo cutover plan (`docs/planning/2026-11-03_17-05_SUPERB-forgejo-fleet-cutover-pareto-plan.md`), and the 2026-10/11 status reports.

## Status legend

| Status           | Meaning                                                     |
| ---------------- | ----------------------------------------------------------- |
| 🔴 `TODO`        | Not started. Needs doing.                                   |
| 🟡 `IN_PROGRESS` | Actively being worked on.                                   |
| 🔵 `BLOCKED`     | Cannot proceed, external dependency or decision needed.     |

## High Impact

| Task                                                                                     | Status     | Impact | Effort | Evidence                                                                                     |
| ---------------------------------------------------------------------------------------- | ---------- | ------ | ------ | -------------------------------------------------------------------------------------------- |
| M17: fix error-discarding sites (`flags.go` helpers, `init.go` glob, `validate_main.go` DirExists, `parseGitHubRemote` placeholders) | 🔴 `TODO` | High | 1h | `cmd/goreleaser-wizard/flags.go`, `cmd/goreleaser-wizard/init.go:248`, `cmd/goreleaser-wizard/validate_main.go:167`, `cmd/goreleaser-wizard/types/template_data.go:236` (plan M17 / F61-F64) |
| M16b: dual-provider `goreleaser check` E2E — only github+forgejo goldens are check-proven today | 🔴 `TODO` | High | 1h | `cmd/goreleaser-wizard/golden_test.go`, goldens `forgejo-workflows.yml`/`goreleaser-dual.yaml` (2026-11-04 report f.13) |
| M16c: unit tests for `NewForgejoWorkflowsTemplateData` pin fields (BuildFlow rev, goreleaser v2.17.1 must never float) | 🔴 `TODO` | High | 1h | `cmd/goreleaser-wizard/types/template_data.go`, `cmd/goreleaser-wizard/generators/forgejo_workflows.go` (2026-11-04 report f.14) |
| M16: generate `.forgejo/actionlint.yaml` (self-hosted `native` label) + golden            | 🔴 `TODO` | Med   | 1h     | E2E needs the `<(...)` config hack from AGENTS.md gotcha 17 (2026-11-04 report f.12)          |
| M14: provider-aware CLI success tips (`forgejo`/`dual` → point at `.forgejo/workflows`, not `.github/workflows`) | 🔴 `TODO` | Med | 30min | `cmd/goreleaser-wizard/generate.go` success tip block (2026-10-28 02-35 report c, 2026-11-04 f.10) |
| M22: `validate` detects deprecated GoReleaser keys (`dockers`, `brews`, `format:`) + prints before→after migration hints | 🔴 `TODO` | High | 2h | plan M22 / F77-F78; `cmd/goreleaser-wizard/validate_main.go` has no deprecation pass |
| CI/test-matrix gaps: `generate` × {github, forgejo, dual, none} matrix, `nix flake check` in CI, `go_1_27 ≥ go.mod` guard test | 🔴 `TODO` | Med | 2h | `.github/workflows/ci.yml` covers lint+test only (2026-10-28 10-10 report f.11/f.15/f.28/f.30) |
| Golden-churn guard: CI asserts `testdata/` stays clean after `-update-golden` runs         | 🔴 `TODO` | Med   | 30min  | regression class fixed by `27c037e` needs a tripwire (2026-11-04 report f.33)                 |

## Medium Impact

| Task                                                                                  | Status     | Impact | Effort | Evidence                                                                                |
| ------------------------------------------------------------------------------------- | ---------- | ------ | ------ | --------------------------------------------------------------------------------------- |
| M19-M21 coverage lifts: `internal/git`, `internal/config`, `internal/utils`, `generators`, `templates`, `types` (0-5% → solid) | 🔴 `TODO` | Med | 6h | `go test -cover ./...`; goldens cover `generators` partially (plan M19/M20, 2026-08-17 01-15 c) |
| M21: domain coverage 10% → 60%+ (`Validate()`, `ApplyDefaults()`, enum `IsValid` table tests) | 🔴 `TODO` | Med | 3h | `internal/domain/` (plan M21 / F74-F76; open since the 2026-06-30 report)                |
| M24: flake.nix shellHook fallback to local cache dirs when `/mnt/buildcache` is dead   | 🔴 `TODO`  | Med    | 30min  | `flake.nix` has no fallback; workaround is manual exports (AGENTS.md gotcha 12, plan M24/F83-F84) |
| M26 sweep: `init.go` should use `printCommandHeader()` like the other commands; remove `OrEmpty()` (zero callers verified twice); review stale `nlreturn`/`wsl` LSP warnings | 🔴 `TODO` | Low-Med | 1h | `cmd/goreleaser-wizard/init.go:25`, `cmd/goreleaser-wizard/types/` (2026-10-28 02-35 lines 102, 136-137; 2026-11-04 f.19) |
| Add `hasFile` unit tests (the function whose 2026-06 breakage killed the build has zero direct tests) | 🔴 `TODO` | Low-Med | 30min | `internal/utils/recommendations.go` (2026-06-30 report §06)                              |
| Split oversized files (>600 lines): `internal/types/validation.go` (898), `internal/validation/basic.go` (767), `cmd/goreleaser-wizard/jobs.go` (685), `internal/validation/business_rules.go` (650) | 🔴 `TODO` | Med | 4h | `wc -l` 2026-11-07; open since the 2026-03-25 report, revisit scope before starting      |
| M12: wave-1 preflight generator for the 7 Phase-F repos + preflight checklists (starred/pending-verdict gates) | 🔴 `TODO` | Med | 2h | builds on the M3 `forgejo-post-flip-verify` harness, fixture-green but never run live (2026-11-04 f.17/f.42) |
| Fix `git log` failures on stoptube/papdashboard/autocart/buildflow (safe.directory/worktree) — unblocks wave-3 preflights | 🔴 `TODO` | Med | 30min | (repo: machine/fleet) 2026-11-04 f.18                                                    |
| Consolidate root-level doc sprawl: 17 root `.md` files (2025-era: `BDD_TESTS_REVIEW.md`, `COMPREHENSIVE_REFACTORING_TODO.md`, `EXECUTION_PLAN.md`, …) → archive to `docs/` or delete; root keeps only living docs | 🔴 `TODO` | Low-Med | 1h | `ls *.md` (2026-06-30 report §06 item 6; root keeps README/FEATURES/TODO_LIST/ROADMAP/CHANGELOG/AGENTS/SECURITY/CONTRIBUTING) |

## Blocked (external decision or owner action needed)

| Task                                                                                     | Status     | Blocker                                                                                     |
| ---------------------------------------------------------------------------------------- | ---------- | ------------------------------------------------------------------------------------------- |
| M18: personal defaults — `--forgejo-owner`/`--runner` flags, tap repo pattern, commit author, homebrew casks wiring | 🔵 `BLOCKED` | Needs Lars's answers: tap pattern, author email, runner default (2026-08-17 00-04 §e)       |
| M28: wizard v2 release prep (CHANGELOG cut, version) — tag only on approval               | 🔵 `BLOCKED` | Owner approval for the tag (2026-11-04 f.26)                                                |
| Forgejo cutover owner-gated phases: A (deploy), B (dnsblockd flip), C (secret), E (token handoff); then M5/M6 live proofs + M7 retro | 🔵 `BLOCKED` | No FORGEJO_TOKEN agent-side by design; see `docs/planning/forgejo-migration-privileged-batch.md` (2026-11-04 line 54) |
| T20: decommission `.github/workflows` per repo                                            | 🔵 `BLOCKED` | Post-flip, after the 2-week green-shadow window (2026-10-28 00-28 line 115)                 |
| M10: authenticated mirror sweep (`has_actions` on private mirrors)                        | 🔵 `BLOCKED` | Needs the Forgejo API token (2026-11-04 line 55)                                            |
| Answer-gated: M17-fleet docker-missing CLI behavior (g3), branch-protection strictness (g1), `RELEASE_SHADOW_GITHUB_TOKEN` PAT policy (g2) | 🔵 `BLOCKED` | Questions for Lars in `docs/status/2026-10-28_02-35_forgejo-migration-full-status-a-g.md` lines 146-148 |

---

<!-- Harvested 2026-11-07 via docs-health AUDIT: HARVEST from docs/status/2026-08-17_*, docs/status/2026-10-*, docs/status/2026-11-04, docs/planning/2026-08-16 (Wave 6), verified against code. Superseded/done items from the 2026-01..08 corpus were annotated in place and archived to docs/status/archived/. -->
