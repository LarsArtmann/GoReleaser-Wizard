# TODO List

> Short-term, actionable, bounded work items, verified against the actual code on 2026-11-07.
> Long-term vision lives in [ROADMAP.md](ROADMAP.md). Completed items are removed and logged in [CHANGELOG.md](CHANGELOG.md).
> Primary sources: the Wave-6 remainder of `docs/planning/2026-08-16_18-45_SUPERB-WIZARD-MODERNIZATION.md` (M17-M27), the Forgejo cutover plan (`docs/planning/2026-11-03_17-05_SUPERB-forgejo-fleet-cutover-pareto-plan.md`), and the 2026-10/11 status reports.

## Status legend

<!-- check-rows:ledger -->

| Status           | Meaning                                                 |
| ---------------- | ------------------------------------------------------- |
| 🔴 `TODO`        | Not started. Needs doing.                               |
| 🟡 `IN_PROGRESS` | Actively being worked on.                               |
| 🔵 `BLOCKED`     | Cannot proceed, external dependency or decision needed. |

## High Impact

<!-- check-rows:ledger -->

| Task | Status | Impact | Effort | Evidence |
| ---- | ------ | ------ | ------ | -------- |

## Medium Impact

<!-- check-rows:ledger -->

| Task                                                                                                                                                                                                              | Status    | Impact  | Effort | Evidence                                                                                                                      |
| ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------- | ------- | ------ | ----------------------------------------------------------------------------------------------------------------------------- |
| M19-M21 coverage lifts: `internal/git`, `internal/config`, `internal/utils`, `generators`, `templates`, `types` (0-5% → solid)                                                                                    | 🔴 `TODO` | Med     | 6h     | `go test -cover ./...`; goldens cover `generators` partially (plan M19/M20, 2026-08-17 01-15 c)                               |
| M21: domain coverage 10% → 60%+ (`Validate()`, `ApplyDefaults()`, enum `IsValid` table tests)                                                                                                                     | 🔴 `TODO` | Med     | 3h     | `internal/domain/` (plan M21 / F74-F76; open since the 2026-06-30 report)                                                     |
| Split oversized files (>600 lines): `internal/types/validation.go` (898), `internal/validation/basic.go` (767), `cmd/goreleaser-wizard/jobs.go` (685), `internal/validation/business_rules.go` (650)              | 🔴 `TODO` | Med     | 4h     | `wc -l` 2026-11-07; open since the 2026-03-25 report, revisit scope before starting                                           |
| M12: wave-1 preflight generator for the 7 Phase-F repos + preflight checklists (starred/pending-verdict gates)                                                                                                    | 🔴 `TODO` | Med     | 2h     | builds on the M3 `forgejo-post-flip-verify` harness, fixture-green but never run live (2026-11-04 f.17/f.42)                  |
| Fix `git log` failures on stoptube/papdashboard/autocart/buildflow (safe.directory/worktree) — unblocks wave-3 preflights                                                                                         | 🔴 `TODO` | Med     | 30min  | (repo: machine/fleet) 2026-11-04 f.18                                                                                         |
| Consolidate root-level doc sprawl: 17 root `.md` files (2025-era: `BDD_TESTS_REVIEW.md`, `COMPREHENSIVE_REFACTORING_TODO.md`, `EXECUTION_PLAN.md`, …) → archive to `docs/` or delete; root keeps only living docs | 🔴 `TODO` | Low-Med | 1h     | `ls *.md` (2026-06-30 report §06 item 6; root keeps README/FEATURES/TODO_LIST/ROADMAP/CHANGELOG/AGENTS/SECURITY/CONTRIBUTING) |

## Blocked (external decision or owner action needed)

<!-- check-rows:ledger -->

| Task                                                                                                                                       | Status       | Blocker                                                                                                               |
| ------------------------------------------------------------------------------------------------------------------------------------------ | ------------ | --------------------------------------------------------------------------------------------------------------------- |
| M18: personal defaults — `--forgejo-owner`/`--runner` flags, tap repo pattern, commit author, homebrew casks wiring                        | 🔵 `BLOCKED` | Needs Lars's answers: tap pattern, author email, runner default (2026-08-17 00-04 §e)                                 |
| M28: wizard v2 release prep (CHANGELOG cut, version) — tag only on approval                                                                | 🔵 `BLOCKED` | Owner approval for the tag (2026-11-04 f.26)                                                                          |
| Forgejo cutover owner-gated phases: A (deploy), B (dnsblockd flip), C (secret), E (token handoff); then M5/M6 live proofs + M7 retro       | 🔵 `BLOCKED` | No FORGEJO_TOKEN agent-side by design; see `docs/planning/forgejo-migration-privileged-batch.md` (2026-11-04 line 54) |
| T20: decommission `.github/workflows` per repo                                                                                             | 🔵 `BLOCKED` | Post-flip, after the 2-week green-shadow window (2026-10-28 00-28 line 115)                                           |
| M10: authenticated mirror sweep (`has_actions` on private mirrors)                                                                         | 🔵 `BLOCKED` | Needs the Forgejo API token (2026-11-04 line 55)                                                                      |
| Answer-gated: M17-fleet docker-missing CLI behavior (g3), branch-protection strictness (g1), `RELEASE_SHADOW_GITHUB_TOKEN` PAT policy (g2) | 🔵 `BLOCKED` | Questions for Lars in `docs/status/2026-10-28_02-35_forgejo-migration-full-status-a-g.md` lines 146-148               |

## Sweep debts (2026-11-07 docs-health audit follow-ups)

<!-- check-rows:ledger -->

| Task                                                                                                                                                                                                                                                    | Status    | Impact | Effort | Evidence                                                             |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------- | ------ | ------ | -------------------------------------------------------------------- |
| 2025 pre-canon corpus sweep: annotate (pre-canon files get routing banners, not fake verdicts) + archive the 29 unarchived 2025 status reports so docs-gates cadence/routing go green (`docs-gates run .` is wired advisory in `scripts/docs-gates.sh`) | 🔴 `TODO` | Med    | 3h     | `docs-gates run .` 5-gate FAIL matrix 2026-11-07; oldest report 508d |

---

<!-- Harvested 2026-11-07 via docs-health AUDIT: HARVEST from docs/status/2026-08-17_*, docs/status/2026-10-*, docs/status/2026-11-04, docs/planning/2026-08-16 (Wave 6), verified against code. Superseded/done items from the 2026-01..08 corpus were annotated in place and archived to docs/status/archived/. Sweep debts added from docs/planning/2026-11-07_08-51_SUPERB-wave6-hardening-and-proof-plan.md. -->
