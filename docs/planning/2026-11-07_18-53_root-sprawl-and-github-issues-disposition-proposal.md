# Root Sprawl & github_issues/ Disposition Proposal

**Status**: PROPOSAL — execution is owner-gated (T26 of `2026-11-07_08-51_SUPERB-wave6-hardening-and-proof-plan.md`, source item §f27/41).
**Date**: 2026-11-07. **Inventory method**: `ls *.md` + `git log -1 --format=%ad` per file, same day.
**Scope**: the 19 root-level `.md` files and the 7 files in `github_issues/`. Execution decisions need the owner because root moves break muscle memory, external bookmarks, and any uncommitted references.

## Why this exists

The 2026-06-30 report (§06 item 6) set the canon: root keeps only living docs (README, FEATURES, TODO_LIST, ROADMAP, CHANGELOG, AGENTS, SECURITY, CONTRIBUTING). That canon is now enforced by `cmd/docs-check` for links, but eight historical files still sit at root next to the living docs, and `github_issues/` holds five issue drafts from the 2025-12/2026-09 era whose fate nobody has decided. This document is the decision-ready inventory: what each file IS, what should happen to it, and what execution costs.

## Inventory — root `.md` files (19)

| File | Lines | Last touched | Classification | Proposed disposition |
| --- | --- | --- | --- | --- |
| README.md | 268 | 2026-11-07 | Living | **KEEP** (sales page, current) |
| CHANGELOG.md | 56 | 2026-11-07 | Living | **KEEP** |
| FEATURES.md | 218 | 2026-11-07 | Living | **KEEP** (counts gate-checked) |
| TODO_LIST.md | 72 | 2026-11-07 | Living | **KEEP** (check-rows ledger) |
| ROADMAP.md | 71 | 2026-11-07 | Living | **KEEP** (review date set ~2027-02) |
| AGENTS.md | 98 | 2026-11-07 | Living | **KEEP** (session context) |
| CONTRIBUTING.md | 445 | 2026-01-24 | Living | **KEEP** (canon root set; content freshness is a separate sweep question, not a move) |
| SECURITY.md | 61 | 2026-10-28 | Living | **KEEP** (GitHub reads it from root; must not move) |
| BDD_TESTS_REVIEW.md | 378 | 2026-03-28 | Historical review snapshot | **ARCHIVE** → `docs/status/archived/` with pre-canon routing banner |
| COMPREHENSIVE_REFACTORING_TODO.md | 378 | 2025-12-15 | Historical todo snapshot | **ARCHIVE** → `docs/status/archived/` (superseded by TODO_LIST.md harvest) |
| CRITICAL_ARCHITECTURAL_ANALYSIS.md | 342 | 2025-12-15 | Historical analysis | **ARCHIVE** → `docs/status/archived/` (superseded by go-arch-lint enforcement + current Clean Architecture) |
| EXECUTION_PLAN.md | 351 | 2026-09-04 | Historical plan | **ARCHIVE** → `docs/status/archived/` (pre-dates the docs/planning canon; its successor plans already live in `docs/planning/`) |
| IMPLEMENTATION_SUCCESS_REPORT.md | 273 | 2026-09-04 | Historical status report | **ARCHIVE** → `docs/status/archived/` |
| PROJECT_SPLIT_EXECUTIVE_REPORT.md | 86 | 2026-09-04 | Historical status report | **ARCHIVE** → `docs/status/archived/` |
| CONSUMER_PERSPECTIVE.md | 135 | 2026-09-04 | Historical perspective/ideas | **ARCHIVE** → `docs/status/archived/` (live ideas already harvested into ROADMAP themes; verify during execution) |
| MIGRATION_TO_NIX_FLAKES_PROPOSAL.md | 639 | 2026-09-04 | Executed proposal | **ARCHIVE** → `docs/status/archived/` (the migration it proposes is DONE — repo has shipped on flakes since) |
| PARTS.md | 535 | 2026-02-27 | Historical component analysis | **ARCHIVE** → `docs/status/archived/` |
| dedup-acceptance.md | 274 | 2026-07-29 | Historical acceptance log | **ARCHIVE** → `docs/status/archived/` |
| VERIFICATION.md | 63 | 2026-01-24 | Likely tool scaffold junk (title: "Created directories and files") | **DELETE** after owner glance — 60 seconds to confirm it carries nothing |

## Inventory — `github_issues/` (7 files)

| File | Lines | Last touched | Content state | Proposed disposition |
| --- | --- | --- | --- | --- |
| #001_Generic_Configuration_System_Design_Question.md | 204 | 2026-09-04 | Question ANSWERED: the typed `SafeProjectConfig` + TypeSpec pipeline shipped | **ARCHIVE** directory-wide; substance already landed |
| #002_Fix_All_Compilation_Errors.md | 247 | 2026-09-04 | Obsolete: build has been green through wave-6 hardening | (same) |
| #003_Complete_Test_Suite_Reorganization.md | 380 | 2026-09-04 | Largely landed: 11 test packages, in-package tests, E2E + golden harness | (same; skim for any unharvested intent during execution) |
| #004_Complete_Documentation_Generation.md | 458 | 2026-09-04 | Landed differently: docs-health battery + docs-check instead of generated docs | (same) |
| #005_Performance_Optimization_and_Benchmarking.md | 394 | 2026-09-04 | NOT DONE: repo has no benchmarks; the only live intent | **HARVEST** one ROADMAP idea row ("benchmark the generate hot path") before archiving |
| INDEX.md | 337 | 2026-09-04 | Draft index, states "Last Updated 2025-06-17", claims statuses that no longer match code | (archive with directory) |
| README.md | 51 | 2025-12-15 | Directory explainer | (archive with directory) |

**Directory proposal**: move all 7 to `docs/status/archived/github-issues-drafts/` in one `git mv`, with a one-paragraph routing banner at the top of `INDEX.md` ("drafts from 2025-12…2026-09; outcomes recorded here; none were filed to GitHub"). Alternatively `docs/research/github-issues-drafts/` if the owner prefers keeping them off the status corpus — the status location puts them in docs-gates cadence scope, which is the honest place for documents that make dated claims.

## Link and gate impact (verified 2026-11-07)

- **Inbound references**: the only living-doc reference to any sprawl file is one TODO_LIST evidence cell (text mention, not a markdown link). No README/FEATURES/docs link breaks under the proposed moves.
- **`cmd/docs-check`**: relative-link gate — moving files does not break it as long as no doc links TO the moved paths (verified none do). Links FROM moved files into the repo keep working (archived files resolve `../..` the same as existing archived corpus).
- **docs-gates battery**: moving the 8 root files into `docs/status/archived/` classifies them as historical corpus — the same re-classification the 2026-11-07 sweep applied to the 2026-01…08 reports. This SHRINKS the living corpus (good for the matcher's distinctiveness bars) but each moved file needs the pre-canon routing banner per the sweep canon, so execution should ride the same annotate→archive→format→verify loop as the pending 2025-corpus sweep row in TODO_LIST.
- **dprint**: archived/ is formatter-covered (T1 excluded only testdata and .config) — run `dprint fmt` after the moves, before the verify pass.

## Execution phases (if approved)

1. **Phase 1 — mechanical moves** (~40min): `git mv` the 11 archive-tier root files into `docs/status/archived/` with date-prefixed names matching the corpus convention (`YYYY-MM-DD_<current-name>.md`), prepend the pre-canon routing banner to each (2-line "archived from root, never re-verified" per sweep canon). Delete `VERIFICATION.md` (after owner glance).
2. **Phase 2 — github_issues/** (~30min): harvest the one live intent (performance/benchmarks) into ROADMAP; `git mv` the directory to `docs/status/archived/github-issues-drafts/`; add the routing banner to INDEX.md.
3. **Phase 3 — proof** (~20min): `dprint fmt` → `go run ./cmd/docs-check .` → `nix run github:Artmann-Technologies/docs-gates -- run .` (advisory matrix must not regress beyond the known 5-gate 2025-corpus debt) → `go test ./...` (golden/docs tests unaffected, but run them anyway).

Total: ~1.5h. No code changes; zero production risk; all moves are `git mv` (history preserved).

## Owner decision points

1. **VERIFICATION.md** — delete outright, or archive for one cycle first? (Evidence says tool scaffold junk, but the call is yours.)
2. **github_issues/ destination** — `docs/status/archived/` (in cadence scope, recommended) or `docs/research/` (out of cadence scope)?
3. **CONTRIBUTING.md content freshness** — the file is 2026-01-24 era and predates the Forgejo/dual-provider reality; keep location, but want a content refresh row in TODO_LIST? (Not scheduled here — location is in scope, content is a separate decision.)
4. **Execution timing** — bundle with the pending "2025 pre-canon corpus sweep" TODO row (same annotate→archive→verify loop, one session) or execute standalone?
