# Wave-6 Hardening Execution Status — 24/27 Plan Tasks Done

**Date:** 2026-11-07 18:36
**Session:** execution of `docs/planning/2026-11-07_08-51_SUPERB-wave6-hardening-and-proof-plan.md` (T1-T27)
**Baseline:** plan committed `07ec54c`, tree clean at session start
**End state:** tree clean, 35 commits ahead of origin (16 substantive + daemon chores), all gates green

---

## Executive summary

Executed 24 of 27 plan tasks to completion with every gate green: `go build`, `go test ./...` (11 packages, includes E2E `goreleaser check` + goldens), `golangci-lint run ./...` 0 issues, `go-arch-lint check` clean, `nix flake check` all checks passed, `dprint` byte-stable, docs-check (links + counts) green. The execution found and fixed **nine real bugs** the plan only suspected, plus two more classes nobody suspected (typed-nil landmine, dead URL-injection validator). T25-T27 remain: T25 not started, T26 not started, T27 done (it was started out of order because its evidence fed the CHANGELOG).

Nothing is pushed. Everything below is local on `master`.

---

## a) FULLY DONE (19 tasks + extras)

| Task                     | What landed                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    | Proof                                      |
| ------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------ |
| T1 formatter re-gate     | dprint scope fixed (goldens + tool-managed `.config/` excluded — the YAML plugin was rewriting goldens and breaking byte-exact tests), byte-stability proven with two consecutive passes, post-format gates re-run over the 17-file archived corpus (check-rows 17/17, bare-rows clean, prose-items clean)                                                                                                                                                                                                                                                                                                                                     | `2057baa`; gates in commit msg             |
| T2 post-sweep proof      | `nix flake check` PASS (after vendorHash fix — nix-hash-fix failed, 88% step-failure rate per its own telemetry, so the `got:` hash was pasted with recorded justification), `go mod tidy` clean, daemon end-state verified (renames + manifests committed), sweep hash pinned into CHANGELOG                                                                                                                                                                                                                                                                                                                                                  | flake check "all checks passed"; `c784f57` |
| T3 M17a flag+glob errors | `getBoolFlag`/`getStringFlag` → Must idiom (loud panic via recover middleware), glob error propagates as `INVALID_PROJECT_STRUCTURE` domain error; **exposed 2 latent test bugs** (simplified init command missing production flags, validate tests missing `--project-only`) — both now wire through shared `configure*Command` helpers; **extractBinaryFromPath bug fixed on sight** (returned "main" instead of the directory name for `src/<bin>/main.go`)                                                                                                                                                                                 | `9c3fe88`, new `flags_test.go`             |
| T4 M17b stat+remote      | alternatives-loop `DirExists` error distinction (warning + continue instead of silent skip); `parseGitHubRemote` returns honest empty (placeholders injected by readers) — **fixes unparseable-remote silent empty render** (old code cached "" and `HasPlaceholderGitHubTarget` couldn't flag it); placeholder-warning tests with deterministic override                                                                                                                                                                                                                                                                                      | `f4dfe62`                                  |
| T5 M16b dual E2E         | dual-provider `goreleaser check` E2E — **caught a real gap**: gitea-scm release blocks need an origin remote; fake non-resolvable URL suffices (check never contacts host), network-free; skip-graceful proven both ways with compiled binary                                                                                                                                                                                                                                                                                                                                                                                                  | `9caea8b`                                  |
| T6 M16c pin tests        | exact-value pin guards (BuildFlow rev, nixpkgs goreleaser rev, runner label), well-formedness (full-length revision, not branch/short-hash), constructor wiring + both-templates-render pins; red-then-green tamper proof executed                                                                                                                                                                                                                                                                                                                                                                                                             | `1ea7343`                                  |
| T7 M16 actionlint.yaml   | wizard now generates `.forgejo/actionlint.yaml` (template from the typed runner-label constant, generation target, rollback, golden); **retired gotcha 17's `<(...)` hack** — the actionlint E2E lints workflows with the GENERATED config; AGENTS gotcha 17 updated                                                                                                                                                                                                                                                                                                                                                                           | `7bbf43c`                                  |
| T8 M14 provider tips     | `forgejo`/`dual` runs no longer point users at `.github/workflows` (actively wrong guidance); one helper derives tips from domain workflow-generation semantics; init gained the provider tip; none/config-only contribute no tip; unset-provider default pinned                                                                                                                                                                                                                                                                                                                                                                               | `66f1997`, `tips_test.go`                  |
| T9 golden-churn guard    | CI regenerates goldens then asserts `testdata/` byte-identical; **red/green proven locally** with a template tweak (generator drift); hand-edited goldens stay covered by the comparison test (regen self-heals them)                                                                                                                                                                                                                                                                                                                                                                                                                          | `61d7efa`                                  |
| T10 CI matrix            | goreleaser-check job became a 5-entry matrix: github/forgejo/dual/none + flake-go-pin guard; **new `TestNoneProviderGeneratesNoWorkflows`**; **new `TestFlakeGoVersionMeetsGoModRequirement`** (gotcha-20 drift guard); nix-in-CI documented as deliberately absent (SHA-pin posture vs unverifiable installer); both workflows actionlint-clean                                                                                                                                                                                                                                                                                               | `61d7efa`, `flake_guard_test.go`           |
| T11+T12 M22 validate     | deprecation detection in `validate` (line-aware YAML walk, file:line findings, before→after hints, works WITHOUT the goreleaser binary); **list verified empirically against goreleaser v2.17.1** — dockers/docker_manifests/brews/format_overrides.format in; release.mode + per-docker retry dropped as empirically wrong; **caught an off-by-one** (report the key node's line, not the value's); README documented; smoke-proven end-to-end (brews @ line 9 detected)                                                                                                                                                                      | `b198b32`                                  |
| T13 M26 cleanup          | init header unified through `printCommandHeader`; `CIProvider.OrEmpty` deleted (zero callers); wsl/nlreturn triage verdict recorded (every LSP warning was a real finding, fixed at site, no suppressions)                                                                                                                                                                                                                                                                                                                                                                                                                                     | `76d7f3c`                                  |
| T14 hasFile tests        | hasFile table (root hit/miss/subdir hit) + GetRecommendedProjectType fixture table (all 5 branches)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            | daemon `200fc16`                           |
| T15 M24 flake fallback   | devShell shellHook falls back to /tmp caches when `/mnt/buildcache` is unwritable; both branches proven (fallback fires on dead path and respects pre-set vars; dormant on live mount); `nix flake check` green with hook                                                                                                                                                                                                                                                                                                                                                                                                                      | daemon `4c76ae6`; flake check PASS         |
| T16 git coverage         | **real bug: `ssh://git@github.com:22/owner/repo` parsed port "22" as owner** — fixed with numeric-segment drop; ParseGitHubURL table (8 shapes), version helpers pinned to ACTUAL behavior (incl. incpatch v-prefix normalization the template depends on), Command wrappers against real checkout + temp dir; coverage 0% → 45.5%                                                                                                                                                                                                                                                                                                             | `bf2d3a1`                                  |
| T17 config precedence    | **three real bugs**: explicit `--config` file returned before env/flags/unmarshal (files had NO effect on the struct); env TransformFunc never stripped the prefix (env vars were dead config — mapped to an island key `goreleaser-wizard-debug`); `Set()` returned `fmt.Errorf("...: %w", nil)` (non-nil error on success). Full precedence chain tests (defaults < file < env < flags) + empty-env-means-unset semantic                                                                                                                                                                                                                     | `260054f`                                  |
| T18 templates/types      | ALL 7 embedded templates parse-proofed with documented delims (**the wrong-delim parse attempt immediately caught** that github-actions uses `[[ ]]` for `${{ secrets.* }}` literals); HomebrewTemplate render-proof (both service variants) despite being unwired; NewGitHubActionsTemplateData mapping test                                                                                                                                                                                                                                                                                                                                  | `4a153ab`                                  |
| T19 M21 domain           | coverage 15.5% → **48.7%** (past the 40% gate): 16-path Validate error table, Clone deep-independence, ApplyDefaults, JSON/YAML round-trips, security validators, ValidationUseCase walk against a REAL temp-module filesystem; **typed-nil landmine fixed** (5 validators returned `*DomainError` — nil was a non-nil interface, appending five `<nil>` entries and invalidating VALID configs); addValidationError now records (not drops) non-domain errors; basic-field failures now typed with field; FromJSON/FromYAML `%w(nil)`-on-success bug fixed; **containsURLInjection was dead code — never wired into ValidateMainPath**; wired | `d38540d`                                  |
| T20 basic.go split       | docker block (151 lines, 4 functions) → `basic_registry.go`; build-FIRST guardrail followed; coverage unchanged 57.2%                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          | `b87636e`                                  |
| T21 jobs.go split        | JobFactory → `job_factory.go` (82 lines), rollback/context plumbing → `job_helpers.go` (34); jobs.go 685 → 582; plan's "git-utility helpers" honestly noted as no longer existing in that file                                                                                                                                                                                                                                                                                                                                                                                                                                                 | daemon `506129e`                           |
| T22 fleet skeletons      | `scripts/forgejo/forgejo-flip-batch@.service` (one sudo per wave, sequential, stop-on-first-red) + `cmd/forgejo-preflight` (5 local flip gates with remediation hints, exit 0/1/2, fixture scenarios both shapes, wave-file parser); **smoke-proven: correctly judges the wizard's own repo not flip-ready**                                                                                                                                                                                                                                                                                                                                   | `ae56dd7`                                  |
| T23 standing docs CI     | `cmd/docs-check` (relative-link resolution + FEATURES counts re-derived from rows, archived history excluded) + `docs` job in ci.yml; **markdown formatting deliberately NOT CI-gated** — curl\|sh dprint installer would breach the SHA-pin posture (documented decision); repo gates green                                                                                                                                                                                                                                                                                                                                                   | `b160050`                                  |
| T24 docs-gates adoption  | battery runs (`nix run github:Artmann-Technologies/docs-gates -- run .`) — **5 gates red, findings are REAL**: 29 unarchived 2025 pre-canon reports (oldest 508d); threshold-guessing rejected (keys didn't apply; no config = honest fleet defaults); `scripts/docs-gates.sh` composition wired ADVISORY (gaming thresholds green would be the lie the tooling prevents); remediation routed as TODO_LIST sweep row; 4 completed TODO rows deleted per canon; TODO_LIST tables marked `check-rows:ledger`                                                                                                                                     | `0308d5b`, daemon `0ec2cbe`                |
| T27 dated claims         | own `.goreleaser.yaml` passes `goreleaser check` (v2.17.1, zero deprecations) captured 2026-11-07; FEATURES Docker claim re-dated with both dates; living-docs scan found only legitimate historical dates                                                                                                                                                                                                                                                                                                                                                                                                                                     | daemon `be03670`                           |

**Extras not in the plan:** stray 15MB build binary the daemon committed at repo root — removed, gitignored (`3821692`); `CIProvider.GeneratesAnyWorkflows()` added to the domain method family.

## b) PARTIALLY DONE (0)

None. Every started task reached done-or-routed state. (T25 was reconnoitered — target lines located — but no edits landed, so it counts as NOT STARTED.)

## c) NOT STARTED (3)

| Task                     | What it is                                                                                                                                                                                                | Effort |
| ------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ |
| T25 housekeeping         | 3 lessons → crush-config `references/lessons.md` (cross-repo commit); ROADMAP quarterly-prune date (~2027-02); DOMAIN_LANGUAGE Commands/Events example rows; AGENTS gotcha 11 slim to "obsolete, history" | 30min  |
| T26 disposition proposal | inventory + classification of 17 root `.md` files + `github_issues/` (25 files); execution awaits owner answer (g3)                                                                                       | 30min  |
| Push                     | 35 commits local-only. Previous session pushed on request; this session has no explicit push instruction                                                                                                  | 1min   |

## d) TOTALLY FUCKED UP (session's own errors, all recovered)

1. **Lost-update races with the daemon (3×).** The auto-commit daemon's formatter raced my file edits — notably `addValidationError` in validation.go silently reverted once (edit tool reported success, file showed old content, test kept failing). Cost: ~15 minutes of ghost-debugging a fix that "was applied". Mitigation adopted mid-session: atomic python read-modify-write for contested files, verify-with-grep after every write. This is a systemic hazard for any session that edits while the daemon commits.
2. **Committed a 15MB build binary.** A stray `go build` output at repo root got swept into a daemon auto-commit. Caught and fixed (`3821692`, gitignored), but the binary sits in history — repo size permanently carries ~15MB until a filter-rewrite (not worth it).
3. **Test-fixture assertion bugs (5×)** in my own new tests: wrong actual-behavior assumptions (IncPatchVersion v-prefix, DefaultCGOEnabled(mobile)=false, reserved names "gorm"/"example" not reserved, ConvertToDockerRegistry fallback to DockerHub, empty docker-image-name allowed). Each was fixed by probing ACTUAL behavior first, then pinning with why-comments — the right pattern, but I wrote assumptions before probes five times in a row.
4. **Churn-guard red/green demo tested the wrong failure class first** (hand-edited golden — regen self-heals it, so the guard "failed to catch"); had to re-derive the guard's true target (generator output drift) and re-prove. The demo is correct now, but the first attempt shipped a false "GUARD_FAILED" conclusion into my own notes.
5. **Validate line-number off-by-one shipped before the test caught it** — `mappingKeyLine` returned the value node's line; the fixture test caught it pre-commit. Good gate, bad first implementation.
6. **dprint damaged goldens on first run** (T1): I ran `dprint fmt` repo-wide before realizing the YAML plugin rewrites generated goldens (breaking byte-exact tests) and `.config/metadata.yaml` is daemon-managed. Restored + excluded + verified stable — but the damage class was predictable and I should have scoped the formatter before running it.
7. **`bash scripts/docs-gates.sh`-style quoting slip** — one malformed `git commit -m "$(...)"` line aborted a command batch (no damage, but sloppy).

## e) WHAT WE SHOULD IMPROVE

1. **The daemon needs an edit-race protocol.** Every lost update cost 5-15 min. Options: a lock-file convention, a "verify after write" habit (now manual), or disabling the daemon's formatter during interactive sessions.
2. **Coverage gates should be standing, not aspirational.** domain went 15.5%→48.7% only because a plan task said so; `internal/types` still sits at 5.2% and generators at 0%. A `go test -cover` threshold in CI would make regression visible.
3. **The battery-red 2025 corpus is now the single biggest docs debt.** 29 pre-canon reports make 4/5 docs-gates red. The annotate+archive sweep (~3h) would let the advisory battery become blocking.
4. **Empirical verification beats plan guesses — every time.** Two plan-claimed deprecations (`release.mode`, per-docker retry) were wrong; two guard designs (churn guard, dual-E2E remote) needed re-derivation against real behavior. Plans should mark unverified claims as such.
5. **Providers beyond github are still placeholder-prone.** `ParseGitHubURL` fallbacks to literal "owner"/"repo"; GitLab/Forgejo remotes fall into the same bucket. The forgejo templates carry `ForgejoOwner` from a domain constant, not detection.
6. **`internal/types` (5.2%) and TypeSpec-generated configs are untested.** The branded-ID layer is the least-tested code in the repo.
7. **CI cannot run the nix proof or the formatter.** Both documented as deliberate, but a self-hosted runner with nix would close both gaps honestly.
8. **The goreleaser-check E2E depends on PATH goreleaser version.** A pinned-version matrix entry would catch version-drift deprecations earlier than fleet rollouts do.
9. **Wave-file parsing lives in two places** (Go `waveList` + the systemd unit's shell expectations). The batch wrapper's Go port should consolidate before wave-1 flip.
10. **Session-level: sub-agent throttling and probe-first testing worked well** — keep both as defaults for future hardening sessions.

## f) UP TO 50 NEXT THINGS (ordered by impact)

**Immediate (this session's loose ends):**

1. Push the 35 commits (awaiting owner word).
2. T25: commit 3 lessons to crush-config `references/lessons.md` (single-specfile repair trap, checkbox text-keys, daemon edit-race protocol).
3. T25: ROADMAP quarterly-prune date line.
4. T25: DOMAIN_LANGUAGE Commands/Events example rows with real wizard commands.
5. T25: AGENTS gotcha 11 → "obsolete, kept for history".
6. T26: root-sprawl inventory + classification table (17 root .md files).
7. T26: `github_issues/` (25 files) classification.
8. T26: disposition proposal doc for owner decision.
9. CHANGELOG entries for the whole wave-6 execution session (the daemon commits carry fragments; a curated entry is missing).
10. Wave-6 plan file: mark T1-T24 done with commit hashes (the plan file still says PLANNED).

**Bug-class closures surfaced this session:**
11. Sweep ALL error constructors for the `%w(nil)`-on-success pattern (FromJSON/FromYAML/Set fixed; audit the rest of the codebase for `fmt.Errorf` wrapping possibly-nil errors).
12. Sweep ALL validators for dead security helpers (containsURLInjection pattern: helper exists, never called).
13. Audit remaining `*DomainError` concrete-return signatures for the typed-nil trap (validation.go fixed; errors.go constructors, interfaces.go use cases remain).
14. `git.ParseGitHubURL`: add Forgejo/GitLab remote detection or return typed "unsupported-host" instead of owner/repo placeholders.
15. `internal/types` coverage lift (5.2% — the branded-ID layer).
16. `cmd/goreleaser-wizard/generators` coverage (0% — goldens exercise it indirectly only).
17. `cmd/forgejo-preflight` coverage (32.4% — main() path untested).

**Structural hardening:**
18. Coverage threshold in CI (e.g. fail if any package drops below its recorded floor).
19. Self-hosted Forgejo runner for CI: enables real `nix flake check` + dprint check without breaking the SHA-pin posture.
20. goreleaser-version-pinned E2E matrix entry (catch version-drift deprecations).
21. Consolidate wave-file parsing into the preflight tool; the batch wrapper becomes a thin systemd call.
22. `forgejo-preflight`: add the Forgejo API gates (repo exists, mirror attached, branch protection) behind an opt-in `--live` flag (skeleton stage today).
23. Golden for the full dual-provider workflow SET (ci.yml + release.yml + actionlint.yaml together as one fixture dir).
24. Snapshot-mode golden (`goreleaser release --snapshot --clean` dry proof for a generated config in CI).
25. Actionlint the GENERATED release.yml for the github provider too (only forgejo workflows are linted today).
26. TestInitCommand-style production-wiring test: assert every command's flag registration matches what its Run function reads (the divergence class T3 fixed is still possible for future commands).

**Docs:**
27. 2025 pre-canon corpus sweep (annotate-rout + archive 29 reports) → docs-gates battery goes green → flip `scripts/docs-gates.sh` from advisory to blocking.
28. Root-sprawl execution (post owner decision from T26).
29. FEATURES.md: add forgejo-preflight + docs-check as features (new user-facing tools shipped this session).
30. README: document `cmd/forgejo-preflight` + `cmd/docs-check` in a "Companion tools" section.
31. AGENTS: new gotcha — daemon edit-race protocol (verify-after-write, python atomic writes).
32. AGENTS: new gotcha — vendorHash let-binding shape defeats nix-hash-fix (paste `got:` hash manually).
33. docs/goreleaser-guide.md: record the empirical deprecation table (release.mode NOT deprecated in v2.17.1, per-docker retry is a type error) so future validate work doesn't re-guess.
34. DOMAIN_LANGUAGE: CIProvider/flip/wave/preflight terms from this session.

**Fleet migration readiness:**
35. forgejo-flip-batch Go port (retire the shell expectations in the unit file).
36. Wave-1 repo list generation from the 110-repo fleet inventory.
37. `forgejo-post-flip-verify` integration test against a fixture Forgejo instance.
38. Release-shadow smoke test (RELEASE_SHADOW_GITHUB_TOKEN path) in a sandbox.
39. M12 preflight: add goreleaser check + actionlint as LIVE gates to forgejo-preflight (run the binaries, not just file existence).

**Quality:**
40. `golangci-lint` erraudit pass over the new code (AsType migrations landed; Is-sentinel sites need review per the decision tree).
41. Add `-race` to the local default test invocation (CI has it; local doesn't).
42. Split `validate_deprecations.go` reporting from detection (detect is pure, report touches results — easier testing).
43. Extract `generationTargets` + preflight into the same package section as GenerationPreflightJob (cohesion).
44. `tips.go` output capture test (printGenerateSuccess currently only has the tips table tested, not the rendered banner).
45. Consider `os.Root`-scoped walks in docs-check (gosec G122 suggested; low risk, recorded as nolint today).

**Housekeeping:**
46. Prune `docs/status/` 2025 files list into the disposition proposal (overlaps T26).
47. `scripts/` README: document the new scripts (docs-gates.sh, forgejo/).
48. Milestone labels for the wave-6 work (the repo has create_milestones.sh; nothing assigned this session).
49. Verify the 11 test packages' coverage floors are recorded somewhere durable (AGENTS or a coverage file) so the CI threshold (#18) has a reference.
50. Re-run `nix flake check` after any flake.nix touch as a standing session ritual (caught the vendorHash break this session).

## g) QUESTIONS (cannot answer myself)

1. **Push now?** 35 commits (16 substantive) sit local-only on `master`. Previous session pushed on explicit request; this one has no instruction. If yes: plain push, or squash the 19 `chore: auto-commit` daemon fragments first (history rewrite on unpushed commits is safe but changes the audit trail).
2. **Should the 2025 pre-canon corpus (29 status reports, oldest 508 days) get the annotate+archive sweep (~3h) so the docs-gates battery can go from advisory to BLOCKING — or are those files delete/archive-as-is candidates under the T26 root-sprawl disposition?** The answer determines whether I schedule the sweep or fold it into the disposition.
3. **For the Forgejo fleet flip (BLOCKED tier): is `docs/planning/forgejo-migration-privileged-batch.md` still the current owner-run script, or should the new `forgejo-flip-batch@` skeleton replace its Phase-C loop?** The two overlap; I won't touch the owner-run file without knowing which is canonical.

---

## Gate snapshot at report time

```
go build ./...          exit 0
go test ./...           11/11 packages ok (E2E + goldens included)
  coverage: docs-check 70.9%, config 84.1%, types(app) 75.0%, utils 62.3%,
            validation 57.2%, goreleaser-wizard 54.8%, domain 48.7%,
            git 45.5%, forgejo-preflight 32.4%, internal/types 5.2%
golangci-lint ./...     0 issues
go-arch-lint check      OK
nix flake check         all checks passed (aarch64 omitted: incompatible systems)
dprint (2nd pass)       no-op (byte-stable)
docs-check .            all gates green
actionlint (both wf)    clean
docs-gates battery      5 gates RED (2025 pre-canon corpus; advisory, routed)
pushed?                 NO — 35 commits ahead of origin/master
```
