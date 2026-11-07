# Wave-6 Remainder Completion — T25, T26, Full Gate + Red-Nix Rescue, Push — Status

**Session scope**: continuation session of 2026-11-07 (following `2026-11-07_18-36_wave6-execution-24of27-tasks-status.md`). Owner ordered: execute the whole remaining list (T25, T26, gate, push), then report.
**Method**: READ→UNDERSTAND→EXECUTE→VERIFY per task; every claim below cites a commit, command, or gate output from this session.
**End state**: tree clean, `master` synced with `origin/master` at `c7c32ed` (42 commits pushed, `650cfb2..c7c32ed`), remote CI GREEN on the pushed head (verified via `gh run list`, 55s run).

## Session commit ledger (this continuation)

| Commit    | What                                                                                                      |
| --------- | --------------------------------------------------------------------------------------------------------- |
| `b706d1a` | T25 docs: DOMAIN_LANGUAGE Commands/Events/Flip, ROADMAP review date, gotcha 11 slim                       |
| `0f2b48f` | T26: root-sprawl + github_issues disposition proposal + TODO_LIST proposal row closed                     |
| `d3a8a21` | daemon: carried the 16-row TODO_LIST sync deletion (1 of 5 dirty files)                                   |
| `285823f` | CHANGELOG hardening-wave curation + remaining TODO_LIST/docs dprint formatting                            |
| `ea957d3` | FIX: flake_guard sandbox skip + .gitignore root-anchor (daemon co-staged both files into this one commit) |
| `c7c32ed` | AGENTS.md gotchas: flake-source exclusion + gitignore shadow                                              |
| `0631f8e` | (crush-config repo, NOT pushed) 3 cross-project lessons appended to `references/lessons.md`               |

## a) FULLY DONE

1. **T25 lessons persistence** — 3 lessons written into crush-config `references/lessons.md` and committed (`0631f8e`): spec-fiction scaffolding repair (repair via real constructor + sweep sibling spec files), text-keyed checkbox ledgers (reword = create+delete to the tracker), auto-commit daemon write-race protocol (atomic RMW, grep-verify after every write, `git show --stat` before trusting). The write-race lesson was deliberately trimmed to complement — not duplicate — the pre-existing add/commit-window race lesson at line 412.
2. **T25 DOMAIN_LANGUAGE** — new `## Commands` section with the four real CLI commands (init/generate/validate/version, verified against `main.go`/`validate.go` `Use:` fields); two wave-6 Events rows (`DeprecationFinding`, `ForgejoPreflightGate` — both verified to exist in code); `Flip` glossary entry (`b706d1a`).
3. **T25 ROADMAP prune date** — "Next quarterly review: ~2027-02" added to the header blockquote (`b706d1a`).
4. **T25 gotcha 11 slim** — obsolete `GOEXPERIMENT=jsonv2` entry reduced to a one-line historical note, numbering preserved so gotcha 14b references stay stable (`b706d1a`).
5. **T26 disposition proposal** — `docs/planning/2026-11-07_18-53_root-sprawl-and-github-issues-disposition-proposal.md`: full inventory of all 19 root `.md` files + all 7 `github_issues/` files with lines/last-touch/classification; disposition per file (8 living docs keep, 11 archive to `docs/status/archived/` with date-prefixed names, `VERIFICATION.md` delete-pending-owner-glance, `github_issues/` archive as a directory with one live intent harvested); link/gate impact verified (zero inbound markdown links to any sprawl file); 3-phase execution plan (~1.5h); 4 owner decision points (`0f2b48f`).
6. **TODO_LIST wave-6 sync** — 17 completed rows removed (M17, M16b, M16c, M16, M14, M22, CI-matrix, golden-churn, M24, M26, hasFile, formatter re-gate, post-sweep proof, standing docs CI, docs-gates adoption, dated-claims, T26 proposal row itself). EVERY row was verified against the actual code before deletion (12 greps/file-checks; the initial M17 grep missed because the Must idiom is named `getBoolFlag`/`getStringFlag` — re-verified before proceeding). Deletions only — no rewording — so the check-rows text-key ledger stays coherent (`d3a8a21` + `285823f`).
7. **CHANGELOG curation** — new "Hardening wave (2026-11-07)" section: 4 real bug fixes and 7 added-capability bullets, all citing verified commit hashes (all 22 hashes checked with `git cat-file -t` before citation) (`285823f`).
8. **Full gate, executed honestly** — build OK; `go test ./...` 12/12 packages ok; `golangci-lint run ./...` 0 issues; `go-arch-lint check` OK; dprint byte-stable (2nd pass no-op); `cmd/docs-check .` all gates green (after formatting); actionlint clean on both workflows.
9. **Red-nix rescue** — `nix flake check` was RED: `TestFlakeGoVersionMeetsGoModRequirement` failed in the sandboxed `checks.test` with `open ../../flake.nix: no such file or directory` (nix build source excludes flake.nix). Root-caused via `nix log`, fixed with an honest `errors.Is(err, fs.ErrNotExist)` → `t.Skip` guard (guard still ACTIVE locally and in CI where the flake exists; proven: local test PASS + sandboxed check now passes), plus root-cause hygiene fix: the bare `goreleaser-wizard` .gitignore pattern had been shadowing the whole `cmd/goreleaser-wizard/` source directory since T21 (it made `git add` refuse tracked files under it) — now root-anchored to `/goreleaser-wizard` (`ea957d3`). Re-run: **`nix flake check` RC 0, "all checks passed"**.
10. **AGENTS.md memory capture** — gotcha 20 gained the flake-source-exclusion companion bullet; new gotcha 21 documents root-anchoring binary ignore patterns and the daemon co-staging trap (`c7c32ed`).
11. **Push** — 42 commits pushed plain (no history rewrite; the unanswered plain-vs-squash owner question resolved to the safe default), `master` fully synced, remote CI GREEN on the pushed head (verified via `gh run list`: `docs: capture flake-source-exclusion…` push run, success, 55s).
12. **Battery delta verification** — `docs-gates run .` re-executed AFTER my corpus edits (TODO_LIST deletions, new proposal doc): FAIL matrix 5/5, UNCHANGED from the T24 baseline — 30 unarchived reports (was 29; +1 is the 18-36 status report itself, expected self-accumulation), oldest unchanged at 508d. My deletions caused NO silent de-harvest. This closes the corpus-edit lesson loop this session opened.

## b) PARTIALLY DONE

1. **T26 execution** — the proposal is decision-ready, but phases 1–3 (moves, banners, harvest, proof, ~1.5h) are untouched BY DESIGN: the plan marks execution owner-gated ("execution awaits owner call").
2. **LSP-vs-CLI lint divergence** — the LSP reports 23 warnings in files I touched (`gosec G703` path-traversal in `cmd/docs-check`/`cmd/forgejo-preflight`, `gocognit 28>25` on `checkRelativeLinks`, `unused: func waveList`, `gci`, `wsl_v5`, `golines`) while `golangci-lint run ./...` says 0 issues. I trusted the CLI (per the verify-tool-output lesson) but did NOT triage WHY the two configs disagree — the LSP is running different rules or different severity mapping. Unresolved.
3. **crush-config deliverable** — the lessons commit exists locally; the repo was NOT pushed and its own gates (buildflow/flake checks) were NOT run. The lesson content is in; the delivery layer is not.
4. **CHANGELOG non-goal notes** — the hardening wave is curated, but the two deliberate non-goals (nix-in-CI, `.docs-gates.json` seeding) are documented in `ci.yml` comments and the 18-36 report, only summarized in CHANGELOG. Acceptable, but a reader of CHANGELOG alone doesn't get the reasoning.
5. **Coverage remainders** (carried, untouched this session): `internal/types` 5.2%, `forgejo-preflight` 32.4%, `validation` 57.2%, domain at 48.7% vs the TODO row's 60%+ target.

## c) NOT STARTED

1. Sprawl execution phases 1–3 (archive 11 root files, banners, github_issues move, proof) — owner-gated.
2. `VERIFICATION.md` deletion (owner glance pending; evidence says tool scaffold junk).
3. `github_issues/` destination decision (status/archived vs research) — owner call in the proposal.
4. CONTRIBUTING.md content freshness (location kept; content is 2026-01-24 era, predates Forgejo/dual reality).
5. Benchmarks → ROADMAP harvest (phase-2 sub-item of sprawl execution).
6. 2025 pre-canon corpus sweep (29→30 unarchived reports, oldest 508d, ~3h; battery's 5 red gates exist because of it).
7. `.docs-gates.json` baseline decision (deliberately absent since T24; threshold keys were rejected on evidence, needs a real decision not a guess).
8. Oversized-file splits, round 2: `internal/types/validation.go` (898 lines), `internal/validation/business_rules.go` (650).
9. Coverage lifts round 2 (see b.5).
10. Entire Blocked tier: M18 personal defaults, M28 v2 tag, cutover phases A/B/C/E, M5/M6 live proofs + M7 retro, T20 decommission, M10 mirror sweep, g1/g2/g3 answers.
11. M12 preflight generator live run; fleet `git log` fixes on stoptube/papdashboard/autocart/buildflow.
12. ROADMAP v2.17 surface ideas (winget/chocolatey/nsis/aur/krew), `homebrew_casks` wiring, config-side `signs:` section, ingest-and-modernize-diff feature, blocking pre-release env check, runner revival, snap/scoop opt-in, parameterized flags (tap repo, NUR repo, commit author, runner label).
13. Performance benchmarks (the one live intent inside `github_issues/#005` — harvested only at sprawl execution).

## d) TOTALLY FUCKED UP (honest)

1. **My exit-code probe lied to me.** First nix gate run used `nix flake check 2>&1 | tail; echo "NIX EXIT: $?"` — the `$?` was tail's exit code. It printed `NIX EXIT: 0` while nix was FAILING. Pipeline masking — the exact failure mode enshrined in the fleet lessons — and I committed it myself mid-session. Caught only because I read the output instead of trusting the probe. If I had trusted it, a red sandbox state ships behind a green-looking number.
2. **I propagated an inherited stale-green claim.** The 18-36 report (and my session summary) said "nix flake check all checks passed at last verify" — false for the post-T10 tree, because `flake_guard_test.go` landed after the last real nix verify and reads a file the sandbox doesn't have. I initially repeated "all gates green at last verify" as context before the run proved otherwise. Lesson now encoded: inherited gate claims are per-commit re-verified, never quoted.
3. **I violated my own freshly-written lesson.** The text-keyed-ledger lesson (written this session into lessons.md) says: re-run the matcher/gate after corpus edits. I edited TODO_LIST (17 rows) + added a proposal doc and did NOT re-run the battery until the user's status request forced the verification sweep. Delta turned out clean (30 vs 29, +1 explained), but the compliance was luck-shaped, not process-shaped.
4. **T21's gitignore footgun sat in the tree for a full session** — the bare `goreleaser-wizard` pattern shadowed the source directory. It surfaced only when it blocked MY commit; nothing tests .gitignore semantics, so it could have bitten any future `git add` under `cmd/goreleaser-wizard/`.
5. **Daemon co-staging reshaped my commits twice** — `d3a8a21` (took 1 of 5 dirty files mid-gate) and `ea957d3` (carried BOTH my gitignore fix and the flake-guard fix under one message). No data lost, but commit attribution was mutated under me; both were caught by reading `git show --stat`, which is now gotcha 21.

## e) WHAT WE SHOULD IMPROVE

1. **Exit-probe discipline**: never `cmd | tail; echo $?`. Capture `rc=$?` before the pipe, `set -o pipefail`, or redirect-to-file-then-tail (the pattern this session ended up using). This should be a fleet lesson, not a memory.
2. **Corpus-edit ritual**: any TODO_LIST/FEATURES/ROADMAP edit triggers an immediate battery re-run — make it part of the commit loop, not the post-hoc sweep.
3. **Inherited-claim verification**: when resuming, re-run each gate once before quoting the prior report's verdicts; the summary's "all green" cost this session a red rescue.
4. **TODO_LIST canon cadence**: completed rows should be deleted in the SAME session that completes them (T24 deleted 4 rows while ~17 stale ones remained) — sync should ride the execution, not lag it to the next session.
5. **Daemon-fragment policy before push**: decide plain-vs-squash BEFORE accumulating 19 `chore:` fragments; next time squash locally pre-push if the owner wants clean history (this push went plain; history is now public and immutable-by-policy).
6. **LSP/CLI lint unification**: one triage — either the CLI lint adopts the LSP's extra rules (gosec, gocognit, wsl as warnings) or the LSP config aligns down. Currently two tools give contradictory verdicts on the same tree.
7. **gitignore hygiene check**: add a cheap guard (or habit) that binary-name ignore patterns are root-anchored — nothing catches this class today.
8. **crush-config delivery**: lessons commits should be pushed + gates run in the same session, same as gw work.
9. **CHANGELOG curation per wave-commit** instead of one batch at session end (batching cost a re-verification pass over 22 hashes).
10. **The battery's routing gate text** ("open actionable rows outside §f sections") fired on the new proposal doc's owner questions — either proposals need a `routing:skip` marker or routing should learn planning-doc semantics; worth one upstream look.

## f) NEXT THINGS (up to 50, grouped, sorted by leverage)

**Immediate / this week**

1. Owner decisions on the sprawl proposal (§g.1 below) → execute phases 1–3 (~1.5h).
2. 2025 pre-canon corpus sweep: annotate + archive 30 reports (~3h) → battery 5 gates to green, then decide `.docs-gates.json` baseline for real.
3. Merge or reject the open dependabot PR (pflag 1.0.10→1.1.0, CI green on the PR).
4. Triage LSP-vs-CLI lint divergence to one source of truth.
5. Push crush-config + run its gates.
6. Watch next CI runs for the 5-provider generate matrix + golden-churn guard in anger (first real push-through-CI of the new jobs went green already, but the matrix jobs deserve one deliberate review of their logs).
7. Add the exit-probe discipline lesson to crush-config lessons.md (this session's d.1).

**Coverage (round 2)**
8. `internal/types` 5.2% → meaningful (validation.go paths).
9. `forgejo-preflight` 32.4% → gate-table coverage.
10. domain 48.7% → 60%+ (close the M21 row honestly).
11. `validation` 57.2% → registry/business-rules paths.
12. `cmd/docs-check` + `cmd/forgejo-preflight` integration fixtures beyond unit.

**Splits (round 2)**
13. `internal/types/validation.go` (898 lines).
14. `internal/validation/business_rules.go` (650 lines).

**Feature surface (ROADMAP-fed)**
15. `homebrew_casks` generator wiring (pending M18 defaults).
16. Config-side `signs:` section emission.
17. v2.17 surface: winget/chocolatey/nsis/aur/krew/NUR generators.
18. Ingest existing `.goreleaser.yaml` → modernization diff.
19. Blocking pre-release environment check as separate command.
20. `--runner` self-hosted revival alongside `ubuntu-latest`.
21. Snap/Scoop opt-in re-add.
22. Parameterized flags: tap repo, NUR repo, commit author, runner label (M18-adjacent).

**Blocked tier (needs owner)**
23. M18 personal defaults (tap pattern, author email, runner default).
24. M28 wizard v2 release: CHANGELOG cut + tag on approval.
25. Cutover phase A (deploy).
26. Cutover phase B (dnsblockd flip).
27. Cutover phase C (secret).
28. Cutover phase E (token handoff).
29. M5/M6 live proofs post-cutover.
30. M7 cutover retro.
31. T20 decommission `.github/workflows` per repo (post-shadow-window).
32. M10 authenticated mirror sweep.
33. g1/g2/g3 answers (docker-missing CLI behavior, branch protection strictness, release-shadow PAT policy).

**Fleet tooling**
34. M12 wave-1 preflight generator live run on the 7 Phase-F repos.
35. Fix `git log` failures on stoptube/papdashboard/autocart/buildflow.

**Docs/tooling hygiene**
36. `.docs-gates.json` decision (post-sweep).
37. CONTRIBUTING.md refresh to Forgejo/dual reality (owner: content vs location decision).
38. Benchmarks implementation (post-harvest of #005 intent).
39. go-arch-lint `deepScan` switch evaluation (currently documented as available).
40. `nix flake check --all-systems` consideration (aarch64-darwin/aarch64-linux omitted today).
41. E2E workflow timeout (2m0s) audit against the dual-provider path duration.
42. Routing-gate semantics for planning docs (proposal-doc `routing:skip` or upstream rule).
43. CHANGELOG: promote non-goal notes into the Unreleased prose.
44. AGENTS.md: gotcha 17/20 cross-check after CI matrix soak.
45. Decide whether `scripts/forgejo/forgejo-flip-batch@.service` canon supersedes `docs/planning/forgejo-migration-privileged-batch.md` (the unanswered flip-script question, folded here from the 18-36 report).

## g) QUESTIONS FOR THE OWNER (cannot self-answer)

1. **Sprawl execution**: Approve phases 1–3 as proposed — archive the 11 root files to `docs/status/archived/`, move `github_issues/` (which destination: `docs/status/archived/github-issues-drafts/` or `docs/research/github-issues-drafts/`?), and DELETE `VERIFICATION.md` outright vs archive-one-cycle-first?
2. **Corpus sweep timing**: Run the 2025 pre-canon sweep (30 reports, ~3h) standalone NOW, or bundle it with sprawl execution into one annotate→archive→format→verify session (shared ritual, one session, one battery green)?
3. **Lint duality**: Is the LSP's stricter rule set (gosec G703/G122, gocognit ≤25, wsl, golines — currently non-blocking warnings) fleet-intentional policy the CLI lint should adopt, or local misalignment the LSP config should drop? (Determines whether `cmd/docs-check`'s walk and `forgejo-preflight`'s parser get refactored for gocognit/gosec or left as-is.)

---

_Report ends. Waiting for owner instructions; nothing executes until answered._
