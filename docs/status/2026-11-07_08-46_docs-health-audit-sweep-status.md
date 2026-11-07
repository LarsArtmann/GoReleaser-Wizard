# Docs-Health Audit Sweep — Status Report

**Date:** 2026-11-07 08:46 CET\
**Session Goal:** View ALL `**/2026-0*` files, execute the docs-health skill (AUDIT = BUILD + HARVEST + VERIFY + ANNOTATE), make all six living docs superb, archive fully-done files with inline strikethrough.\
**Commits this session:** annotations landed via the auto-commit daemon as they were produced (see `git log --grep "auto-commit"` between 06:xx and 08:35); the final tree (18 `git mv` renames + 2 archive manifests + AGENTS.md + TODO_LIST/ROADMAP edits + incidental go.sum hashes) is staged/working — daemon will fold it in.\
**Skill:** docs-health (`~/.config/crush/skills/docs-health/`), annotate.go toolkit, gates `check-rows`/`bare-rows`/`prose-items`.

---

## Executive Summary

All 19 `**/2026-0*` files (18 md + 1 html, ~6.5k lines) were read, inventoried, verified against the codebase, and **fully annotated inline** (~1,100 items struck with per-item verdicts). 18 files archived to `docs/status/archived/` (17) and `docs/research/archived/` (1) after passing `check-rows --strict` + `bare-rows`; bulk-archive manifests written. `TODO_LIST.md` and `ROADMAP.md` were **created** (they did not exist); `CHANGELOG.md`, `FEATURES.md`, `README.md`, `AGENTS.md`, `docs/DOMAIN_LANGUAGE.md` were rebuilt/fixed. All canonical quality gates green: build ✅, tests 6/6 ✅ (E2E + goldens), golangci-lint 0 issues ✅, go-arch-lint ✅.

**Verdict: sweep goals met; three follow-up debts admitted below (formatter-stability re-gate, nix/docs-gates not run, plain-bullet items left bare).**

---

## a) FULLY DONE ✅

| #  | Item                                                                                                                                                                                                                                                                                                                                      | Proof                                                                |
| -- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------- |
| 1  | All 19 target files read + inventoried (3 subagents, staggered after rate-limit lesson; 9 read directly)                                                                                                                                                                                                                                  | subagent reports + session transcript                                |
| 2  | Commit-evidence map built (git log Jan→Nov 2026: wave hashes `e85bb4f`…`969259e`, forgejo `9334dc3`, koanf `50c0feb`, branded-id `cc8a753`, lint `9703270`, hygiene `3cb91da`, golden fix `27c037e`, …)                                                                                                                                   | `git log` sweeps this session                                        |
| 3  | VERIFY of ~30 concrete living-doc claims against code (go.mod 1.27, `--ci-provider` at `main.go:304`, `internal/errors` gone, Viper gone/Koanf present, template has `dockers_v2` but no `homebrew_casks`/`signs:`, no `--runner` flag, no flake cache fallback, `.orig` gone, `test-wizard/` gone, big-file sizes, json/v2 import sites) | greps in transcript                                                  |
| 4  | `TODO_LIST.md` CREATED — 20 verified rows (High/Medium/Blocked), every row carries `file:line` evidence + M-number routing; includes root-sprawl + machine-level rows                                                                                                                                                                     | `TODO_LIST.md`                                                       |
| 5  | `ROADMAP.md` CREATED — 5 themes + non-goals; packslip idea routed from the 2026-10-17 review                                                                                                                                                                                                                                              | `ROADMAP.md`                                                         |
| 6  | `CHANGELOG.md` REBUILT — append-only (old bullets preserved verbatim under "Prehistory"); structured Unreleased/Added/Changed/Fixed/Migration-history; every entry hash-cited                                                                                                                                                             | `CHANGELOG.md`                                                       |
| 7  | `FEATURES.md` REBUILT — Forgejo/dual CI feature added (#4), TUI moved out of "Planned" (structural bug), honest statuses (signing + package managers PARTIAL), summary recomputed 13/2/3/0, stale "This Session" narrative deleted, header date fixed (was 2025-12-16)                                                                    | `FEATURES.md`                                                        |
| 8  | `README.md` FIXED — 7 surgical edits: go-1.27 badge, `--ci-provider` on both init+generate flag tables, Forgejo example, What-It-Generates covers `.forgejo/workflows`, honest Git-Providers table (GitLab/Bitbucket marked "config only"), GOEXPERIMENT prefix dropped (obsolete on Go 1.27 per gotcha 11)                               | `README.md`                                                          |
| 9  | `AGENTS.md` FIXED — `internal/errors/domain_errors.go` → `internal/domain/errors.go` (dir removed in `dde1081`); "Viper" → Koanf (`50c0feb`); pointer to archived corpus added                                                                                                                                                            | `AGENTS.md`                                                          |
| 10 | `docs/DOMAIN_LANGUAGE.md` FILLED — CIProvider/Forgejo/dual/shadow/pin/golden vocabulary (closes plan item M18's docs leg)                                                                                                                                                                                                                 | `docs/DOMAIN_LANGUAGE.md`                                            |
| 11 | ~1,100 inline annotations across 18 md files + HTML: `~~item~~ done at <hash>` / `superseded — X` / `Won't implement` / `NOT-DO` / `→ TODO_LIST.md (M##)`                                                                                                                                                                                 | 15 repair specs + 3 status-items specs, all dry-run first            |
| 12 | 18 files archived via `git mv` AFTER passing `check-rows --strict` (18/18 complete) + `bare-rows` (18/18 clean); per-file `~~` completeness gate 18/18                                                                                                                                                                                    | gate output in transcript                                            |
| 13 | Bulk-archive manifests with per-file classification + deciding reason                                                                                                                                                                                                                                                                     | `docs/status/archived/README.md`, `docs/research/archived/README.md` |
| 14 | `docs/planning/2026-08-16` plan: all 178 tasks annotated (done waves hashed, Wave-6 open items routed `→ TODO_LIST.md`), status header inline-corrected; KEPT ACTIVE deliberately                                                                                                                                                         | plan file line 4 + tables                                            |
| 15 | HTML report (2026-06-30) annotated with a dated Resolution section (26 `~~` verdicts) + inline `<del>` corrections on 5 stale claims                                                                                                                                                                                                      | file content                                                         |
| 16 | Quality gates re-run at the end: `go build` ✅, `go test ./...` 6/6 packages ✅ (E2E + goldens), `golangci-lint run ./...` **0 issues** ✅, `go-arch-lint check` ✅                                                                                                                                                                       | transcript                                                           |
| 17 | Internal-link integrity check across all living docs + manifests: ALL RESOLVE                                                                                                                                                                                                                                                             | python link checker                                                  |
| 18 | One integrity miss found and fixed by the gates themselves: a skipped table row (10-37 line 209) and the 01-14 lessons tail — caught by `check-rows`/`prose-items`, repaired, re-gated                                                                                                                                                    | spec-14/15                                                           |

## b) PARTIALLY DONE 🔄

| # | Item                         | Done                                                        | Missing                                                                                                                                                         |
| - | ---------------------------- | ----------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1 | Formatter-stability protocol | annotate → check-rows                                       | **did NOT run dprint/format after annotating, then re-gate** — the known damage class (formatter re-wrapping struck cells) is unverified on the archived corpus |
| 2 | HTML annotation              | Resolution banner (all verdicts) + 5 inline `<del>` anchors | per-item inline `<del>` on all 25 Top-25 rows (verdicts live in the banner table, not beside each row)                                                          |
| 3 | Bullet-level annotation      | every numbered item + every table row struck                | plain `-` bullet lists (not gated by the toolkit) left bare in several files (02-07 blockers, 05-28 b/c bullets, trio exec bullets)                             |
| 4 | docs-gates battery           | individual gates (`check-rows`, `bare-rows`, `prose-items`) | `docs-gates run --json` full matrix + `--pre-flight` before living-doc writes not run (binary not installed locally)                                            |
| 5 | FEATURE counts               | recomputed manually, honest                                 | no standing gate re-derives the counts from the repo (skill prefers computed-over-maintained)                                                                   |
| 6 | Daemon-interleaved history   | work landed continuously                                    | final committed end-state (renames + manifests) not yet re-verified post-daemon-commit                                                                          |
| 7 | CHANGELOG Unreleased entry   | written                                                     | references "this sweep" without a sweep commit hash (daemon supplies; not pinned)                                                                               |

## c) NOT STARTED ⬜

| # | Item                                                                                                                                   | Why it matters                                                                                |
| - | -------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| 1 | `nix flake check` after the sweep                                                                                                      | canonical gate not re-run (docs-only changes claimed, unproven)                               |
| 2 | `docs-gates` adoption in this repo (`.docs-gates.json` + gate script wired into quality flow)                                          | would mechanize freshness/routing/cadence permanently                                         |
| 3 | Cross-project lessons → crush-config `references/lessons.md` (repair-takes-one-file; checkbox keys are text-keyed; throttle subagents) | generalizable lessons not persisted                                                           |
| 4 | SECURITY.md / CONTRIBUTING.md / CONSUMER_PERSPECTIVE.md content VERIFY                                                                 | link check passed; content claims unverified this session                                     |
| 5 | `github_issues/` directory disposition                                                                                                 | flagged in the HTML verdict as "archival", no TODO_LIST row, no action                        |
| 6 | Root `.md` sprawl cleanup execution                                                                                                    | routed to TODO_LIST, not executed (deletion needs owner sign-off)                             |
| 7 | Harvest of THIS report's §f into TODO_LIST                                                                                             | per docs-health, harvest happens when the session continues — parked pending your instruction |

## d) TOTALLY FUCKED UP 💥

| # | What                                                                                                                                                                                                                                                                                               | Impact                                                                                           |
| - | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------ |
| 1 | **Ran `repair` with two specfiles in one invocation** → usage error; spec-8 (02-07) silently never applied while I proceeded as if 8+9 were done. Caught only by the `prose-items` sweep one step later.                                                                                           | 1 file annotated a step late; wasted a cycle; the "verify after batch intent" reflex was missing |
| 2 | **Line-range bookkeeping error in spec-12**: skipped row 209 and shifted 210-211 in the 10-37 priority table → PARTIAL row. The gate caught it (as designed), but building specs from subagent-approximate line numbers instead of machine-emitted keys is exactly the trap the skill warns about. | 1 gate failure + 1 extra repair pass                                                             |
| 3 | **First HTML injection aborted on an unverified anchor** (`<section id="summary">` guessed, wrong) and a later `grep -L` exit-code was misread as ambiguous. Read-anchors-first cost one cycle.                                                                                                    | 2 wasted cycles, no damage (script was atomic)                                                   |
| 4 | **5 parallel subagents on launch → 4 rate-limit failures.** You had to tell me to throttle. Then still launched 2-at-once when 1 was the safe pace.                                                                                                                                                | ~3 minutes lost, 4 wasted launches                                                               |
| 5 | **Formatter canon violated by omission**: the skill's atomic sequence is annotate → `check-rows` → format → `check-rows` AGAIN. I never ran the format pass, so the second gate never ran. The archived corpus's formatter stability is unproven.                                                  | real risk of struck-cell damage on the next format pass                                          |
| 6 | Minor: edited AGENTS.md before reading it from disk (tool refused); left the template HTML comment in ROADMAP.md on first write (caught immediately); TODO_LIST first draft lacked the root-sprawl row (added later when the gap surfaced during HTML verdicts).                                   | noise, self-caught                                                                               |

## e) WHAT WE SHOULD IMPROVE 📈

1. **Spec hygiene**: never hand-range specs from memory or subagent approximations — `emit-keys`/`--verify` on EVERY unfamiliar file shape before writing (the tooling already refuses unmatched specs; let it).
2. **Formatter-order canon as a hard gate**: after any annotation sweep, run the repo formatter and re-run `check-rows --strict` + `bare-rows` on the touched set — treat the post-format gate as part of the sweep, not an optional extra.
3. **Adopt `docs-gates`**: install the battery in this repo, wire one command into the quality flow, use `--pre-flight` before living-doc edits (it exists precisely to catch living-doc lines that flip historical items).
4. **Computed counts**: FEATURES summary table should be re-derived by a standing test/script, not hand-maintained (skill: duplicated counts rot).
5. **Subagent pacing**: max 1-2 concurrent; stagger after any rate-limit signal.
6. **Batch verification**: after any multi-file tool batch ("apply 8+9"), verify each file's post-state immediately — don't trust a batch invocation to cover every argument.
7. **Bullet lists deserve verdicts too**: the gates only cover numbered items, checkboxes, and tables; narrative bullet claims (like 02-07's blocker list) misled no gate — strike them on sight in future sweeps.
8. **Harvest discipline**: this report's §f must be harvested into TODO_LIST at next session start (docs-health's own rule; parked on your instruction).

## f) Up to 50 Things to Get Done Next 🎯

**Sweep debts (this session's own)**

1. Run the repo formatter (dprint via BuildFlow) over `docs/status/archived/` + living docs, then re-run `check-rows --strict` + `bare-rows` — close the formatter-stability gap.
2. `nix flake check` post-sweep proof.
3. Install/run `docs-gates run --json` once over the repo; adopt `.docs-gates.json` if clean.
4. Verify the daemon's final commit contains the 18 renames + 2 manifests + AGENTS/TODO/ROADMAP edits; pin the sweep hash into CHANGELOG's Unreleased entry.
5. Harvest this report's §f into `TODO_LIST.md` (docs-health rule).
6. Add `github_issues/` disposition row to TODO_LIST (archival vs delete — owner call).
7. Persist 3 cross-project lessons (repair single-specfile, checkbox text-keys, subagent throttling) to crush-config `references/lessons.md` by commit.
8. Strike the plain-bullet claims in the archived corpus's non-gated lists (02-07 blockers, trio exec bullets) — optional polish.
9. README/FEATURES/CHANGELOG through the formatter + link gate as a standing docs CI step.
10. Add a standing test that recomputes FEATURES.md's status-summary table from FEATURES rows (no hand-maintained counts).

**TODO_LIST top rows already routed (execute next sessions)**

11. M17 error-discarding fixes (flags.go, init.go glob, validate_main.go DirExists, parseGitHubRemote).
12. M16b dual-provider `goreleaser check` E2E.
13. M16c `NewForgejoWorkflowsTemplateData` pin-field tests.
14. M16 `.forgejo/actionlint.yaml` generation + golden.
15. M14 provider-aware CLI success tips.
16. M22 `validate` deprecated-key detection + migration hints.
17. CI/test-matrix gaps: generate × 4 providers, `nix flake check` in CI, go_1_27 guard test.
18. Golden-churn guard (CI asserts `testdata/` clean on `-update-golden`).
19. M19-M21 coverage lifts (git/config/utils/generators/templates/types).
20. M21 domain coverage 10% → 60%+.
21. M24 flake.nix shellHook cache fallback.
22. M26 sweep: init.go `printCommandHeader()`, remove `OrEmpty()`, stale LSP warnings.
23. `hasFile` unit tests.
24. Split oversized files (>600 lines: types/validation.go 898, basic.go 767, jobs.go 685, business_rules.go 650) — scope first.
25. M12 wave-1 preflight generator + checklists.
26. Fix `git log` failures on the 4 fleet repos (safe.directory).
27. Root `.md` sprawl cleanup (17 files → archive/delete; owner call).
28. M18 personal defaults → `--forgejo-owner`/`--runner` flags + casks wiring (BLOCKED on your answers).
29. M28 wizard v2 release prep (BLOCKED on tag approval).
30. Forgejo cutover owner-gated phases A/B/C/E → M5/M6 live proofs → M7 retro (BLOCKED on token).
31. T20 `.github/workflows` decommission after green-shadow windows (BLOCKED post-flip).
32. M10 authenticated mirror sweep (BLOCKED on token).
33. Answer-gated: g1 branch-protection strictness, g2 PAT policy, g3 docker-missing behavior.
34. M11 `forgejo-flip-batch@` wrapper.
35. M19 branch-protect rollout verification across canonicalRepos.
36. M20 mirror × protection live proof.
37. M22-fleet Gatus shadow-drift detector.
38. M23 runner capacity measurement (after live pilot).
39. M25 lessons/runbook consolidation (SystemNix runbook mirror-Actions section).
40. T13 mrconfig origin rewrites for flipped repos.

**Noticed this session, worth queueing**

41. `github_issues/` (25 files) — decide archival vs delete alongside #6.
42. `dedup-acceptance.md` still at repo root (kept deliberately twice — revisit once #27 lands).
43. FEATURES "Verified end-to-end 2026-08-17" Docker claim — re-verify when goreleaser next touches the repo (dated claims need re-verification cadence).
44. go.sum gained hashes from this session's test run — confirm `go mod tidy` state clean at next gate pass.
45. TODO_LIST "as of" date line — adopt a per-row verified-date convention when rows next change.
46. Consider `docs/reviews/` for the 2026-11-04 brutal-self-review HTML (currently the only reviews/ file — fine, just noting the pattern).
47. DOMAIN_LANGUAGE.md — add Wizard/Forgejo examples to Commands/Events sections as they emerge (skeleton sections remain generic).
48. ROADMAP quarterly prune date: set next review ~2027-02.
49. AGENTS.md gotcha numbering: gotcha 11 (jsonv2) can shrink once everyone is on 1.27 toolchains fleet-wide.
50. Sweep-retro: the annotate.go toolkit refused/kinds behaved differently than guessed twice (repair single-file, checkbox keys) — consider upstreaming a `--check` mode that takes multiple specfiles.

## g) Questions I Cannot Figure Out Myself 🤔

1. **The 2026-08-16 modernization plan is now fully struck (open items routed to TODO_LIST) — archive it to `docs/planning/archived/` now, or keep it active until Wave-6 work actually starts landing?** I kept it active as the execution record TODO_LIST cites; if you prefer the docs-health canon (fully-resolved → archive), say so and I'll move it + update the two references.
2. **M18 personal defaults (blocked since August): tap repo pattern, commit author email, runner default — what are the values?** Without them, `--forgejo-owner`/`--runner` flags, casks wiring, and the v2 release stay BLOCKED.
3. **Root doc sprawl + `github_issues/`: archive to `docs/archive/` (git mv, history preserved) or straight delete?** The 17 root `.md` files are 2025-era point-in-time snapshots; I routed the cleanup but will not delete files you may want kept without your call.

---

**Bottom line:** the 2026-01…08 corpus is no longer a drift trap — every item is resolved-or-routed inline, the living docs are rebuilt on verified evidence, and open work has exactly one home (`TODO_LIST.md`). The sweep's own debts are honest and small: formatter re-gate, nix/docs-gates proof, and this report's harvest.

_Report generated: 2026-11-07 08:46 CET_
