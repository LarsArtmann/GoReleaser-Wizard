# Session Status — Plan Re-Verification + M8/M9/M3/M13 Executed and Pushed

**Created**: 2026-11-04 07:55 CET
**Author**: Crush session (status-report skill)
**Session window**: 2026-11-04 ~06:50–07:55 CET (single session, two phases: plan review → execution)
**Input**: user directive "Is this plan superb!???" + re-pasted pareto-planning instruction template (steps 1–8)
**Scope of this report**: everything below is from THIS session's run. Where I cite a command, I ran it this session and the exit code is recorded here.

## Scope: what this report verified vs. what it did NOT verify

**Verified this session (command → artifact):**

- Plan arithmetic — `awk -F'|'` recomputation of both plan tables (medium 28 tasks/1765 min, fine 134/1191 min).
- depgraph CLI behavior — `--help` capture + executed runs (small-dir rc=0/175KB; full run rc=0/4,559,298 B `/tmp/fleet.json`; first full run failed silently, see §d).
- Fleet counts — `jq` aggregations over `/tmp/fleet.json` (241 dirs → 110 own projects / 85 Go / 300 library modules / 137 mr-tagged / 977 third-party).
- M3 fixture suite — `nix build .#checks.x86_64-linux.forgejo-scripts-fixture` → **rc=0** (rc-checked directly, after one false-green, see §d).
- M13 fix — `go test ./cmd/goreleaser-wizard/` green, `go test … -update-golden` green, `git status testdata/` empty (goldens byte-identical).
- Wizard gates — `go build ./...`, `go test ./...`, `golangci-lint run` (0 issues), `go-arch-lint check` (no warnings), run against the working tree including the foreign go.mod/go.sum diff.
- Pushes — git push outputs: wizard `f674b44..de2dcb1`, `de2dcb1..90ea06e`, `90ea06e..27c037e`; SystemNix `d34e239d7..ab09e3e71`. Both trees clean at session end.

**NOT verified (and why):**

- Live Forgejo state (is dnsblockd flipped? is Actions armed? protection applied?) — **no FORGEJO_TOKEN agent-side** (root-gated by design, batch doc Phase E handoff pending). The M3 harness has NOT touched the live instance.
- dnsblockd CI/release behavior on the real runner (M5/M6) — owner-gated.
- Full `nix flake check` — only the scoped fixture check was built; host evals not re-run this session.
- `mr list` fleet count (~265 claim) — `mr list` in the wizard cwd returned 3 rows (wrong cwd/config scope); I substituted depgraph's own mrconfig enrichment (137 tagged) instead of chasing mr semantics.
- depgraph git enrichment (`commitCount`, `commitsSinceTag`) — all zeros/`hasGit:false` across nodes; cause unknown, flagged in §e, not investigated (out of session scope).
- The 4 repos where `git log` failed (stoptube, papdashboard, autocart, buildflow) — recorded as preflight blockers, root cause (safe.directory/worktree?) not diagnosed.

---

## a) FULLY DONE

| # | Work                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 | Evidence / commits                                                                                                |
| - | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------- |
| 1 | **Brutal re-verification of the 11-03 cutover plan**: 28 claims checked against git, `--help`, upstream Gitea source, live gates. **8 defects found and fixed in the plan doc** (understated totals 26.5h→29.4h; graph↔table ID drift with a v2-release task existing only in the graph → restored as M28/F133-F134; 8 of 50 §f items unmapped → folded into M7/M16/M19/M25/M26+M28; M15 rescoped — ForgejoOwner/RunnerLabel fields already ship in `types/template_data.go:79-99`; depgraph commands corrected against `--help`; F29 scoped — dnsblockd is PRIVATE per gh API; fine-table 4-col header over 3-cell rows; stale working-tree errata replaced by a §6 consistency guard). **1 reviewer false alarm tested and cleared**: actions-secret API takes plaintext (`CreateOrUpdateSecretOption.Data`, go-gitea/gitea source), server seals — M21's simple-PUT design holds. | plan commit `de2dcb1` (pushed); review report `docs/reviews/2026-11-04_07-01_brutal-self-review.html`             |
| 2 | **M8 — depgraph fleet inventory** (F36–F38, F40-modified): built `/tmp/depgraph` (workspace-mode fix: `GOFLAGS=`); full run over `~/projects` with `--output`; counts extracted via jq. Measured truth replaces the "~265 repos" estimate: **110 own projects (85 Go), 300 own library modules, 977 third-party, 137 mr-tagged**. Errata line added to the plan.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     | `90ea06e` (pushed); `/tmp/fleet.json` (transient)                                                                 |
| 3 | **M9 — wave assignment** (F41–F45): Phase F written into `docs/planning/forgejo-migration-privileged-batch.md` — measured-counts table, activity-signal disproof, 4 waves with flip order and rationale (wave 1: 7 Go CLIs incl. wizard dogfood; wave 2: supply chain with BuildFlow LAST; wave 3: 64 remaining Go apps; wave 4: Rust/JS/CV/private), 4 preflight blockers flagged.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  | `90ea06e` (pushed)                                                                                                |
| 4 | **M3 — post-flip verification harness** (F9/F10/F13/F14): `forgejo-post-flip-verify <repo> [--patch-actions]` as a writeShellApplication in `_forgejo-scripts.nix` — 5 read-only checks (native repo, Actions armed, ci.yml+release.yml present, push mirror attached, default branch protected), reusing the API shapes of the shipped flip/branch-protect/attach-mirrors scripts; `--patch-actions` is the only mutation. Wired into forgejo.nix (inherit list + `environment.systemPackages`). **7 fixture scenarios green** (happy, pull-mirror refusal, actions-off→PATCH with jq payload assert, workflows missing, no mirror, no protection).                                                                                                                                                                                                                                 | SystemNix `4a947ce46` (script+wiring, daemon-absorbed, content verified at HEAD) + `ab09e3e71` (fixtures), pushed |
| 5 | **M13 — `-update-golden` flag-leak fix** (F62–F65): root cause = TestInitCommand's cobra `Execute()` parsing `os.Args[1:]`; fixed with `cmd.SetArgs([]string{})` (init_test.go:107). Proven: suite green both with and without the flag, `testdata/` byte-identical after an update run. AGENTS.md gotcha 14b flipped from "known failure" to "FIXED".                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               | `27c037e` (pushed)                                                                                                |
| 6 | **Wizard gates re-run green on the foreign go.mod/go.sum diff** (indirect-dep patch bumps that appeared after the plan commit; judged on merits, kept, validated).                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   | build OK / tests ok / lint 0 issues / arch-lint clean                                                             |

## b) PARTIALLY DONE

| # | Work                                 | Done                                                          | Missing                                                                                             | Blocker                                                  |
| - | ------------------------------------ | ------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- | -------------------------------------------------------- |
| 1 | M3 harness **run** (F11/F12)         | script + fixtures                                             | never executed against live Forgejo/dnsblockd                                                       | no `FORGEJO_TOKEN` agent-side (Phase E handoff)          |
| 2 | M8 Mermaid export (F39)              | JSON inventory complete                                       | `depgraph generate --format mermaid` not run                                                        | none — deferred as low-value while JSON holds the graph  |
| 3 | F40 inventory↔mr cross-check         | substituted with depgraph's mrconfig enrichment (137 tagged)  | `mr list` full-fleet count not obtained (3 rows in repo cwd — wrong scope)                          | mr semantics/cwd; not chased per "no unrelated research" |
| 4 | Fleet activity data for wave scoring | waves assigned via coupling + blast radius + domain knowledge | `commitCount`/`commitsSinceTag` empty (`hasGit:false` everywhere) — git enrichment did not populate | cause unknown; next-session investigation item (§f-31)   |

## c) NOT STARTED

- **Owner-gated critical path**: M1 Phase A deploy, M2 Phase B flip, M4 Phase C secret, M5 live CI proof, M6 live release proof (tag v0.9.5 — verified correct next tag), M7 pilot retro (also gated on flip), Phase E token handoff.
- **Token-blocked**: M10 authenticated mirror sweep (`has_actions` on private mirrors).
- **Answer-gated**: M17 docker-missing behavior (needs g3), policy halves of M19 (g1) and M21 (g2).
- **Unblocked, queued**: M11 `forgejo-flip-batch@` (SystemNix unit), M12 wave-1 preflight generator (now unblocked by Phase F), M14 provider-aware CLI output, M15 `--forgejo-owner`/`--runner-label` flags, M16 actionlint.yaml template + dual-provider goreleaser-check E2E + pin-field unit tests, M18 README/FEATURES/DOMAIN_LANGUAGE, M19 rollout verification, M20 mirror×protection proof, M21 secret loop script, M22 Gatus drift detector, M23 capacity, M24 attic templates, M25 lessons/runbook consolidation, M26 cleanup sweep (incl. OrEmpty deletion — zero callers verified), M27 self-host CI + tap ADR, M28 wizard v2 release.

## d) TOTALLY FUCKED UP (own failures this session, honestly)

1. **False green on the fixture suite.** I ran `nix build … -L | tail -6 && echo FIXTURE_GREEN` — the pipe masked the build's exit code, so I asserted green while the build was FAILING (missing PATCH stub). Caught one step later with `rc=$?` re-run. A gate asserted through a pipe is not a gate.
2. **First full depgraph run died silently**: exit 0, zero-byte output, logs ending at `progress=196/198`. I had encoded the stdout-redirect invocation into the background job instead of the `--output` flag that the tool's own `--help` documents. One lost cycle; `--output /tmp/fleet.json` worked immediately. (The plan itself had already been corrected for a different depgraph invocation error — this session still stepped into the next one.)
3. **Artifact leak**: depgraph wrote a 4.5 MB `output.json` into the wizard repo root during the inventory runs. Trashed; but tooling that writes default outputs into cwd should have been run from a scratch dir.
4. **Two fixture-authoring mistakes shipped in my first draft**: (a) asserted `"--patch-actions"` through `need_has`, whose `grep -qF` parses leading-dash strings as options; (b) never stubbed the `.PATCH` response, so the default 404 silently failed the patch path. Both caught by the harness, both fixed, both lessons now encoded in the fixture file.

## e) WHAT WE SHOULD IMPROVE

1. **Never pipe verification gates.** Assert `$?`/`PIPESTATUS[0]` directly; a piped build "passes" whenever `tail` exits 0. This session produced one false green this way.
2. **Fixture stub matrix before scenarios**: enumerate every HTTP call a script can make (method × URL × outcome) and stub all of them first; assert request payloads with `jq`, never grep; never assert strings starting with `-`.
3. **`--help` before encoding**: any CLI invocation written into a plan or a script gets verified against the tool's own help output first (depgraph cost two invocations across two sessions before this stuck).
4. **Daemon interplay**: the auto-commit daemon absorbed the M3 script+wiring mid-feature with a terse message. Flow that works: verify content at HEAD (`git grep`), then commit the remaining piece with the full rationale referencing the daemon commit. Don't fight it; do document it.
5. **depgraph git enrichment** (`hasGit:false`, all-zero commit counts) is worth one investigation — a working `commitsSinceTag` would give wave scoring a real activity signal even under the commit daemon (tag-age survives daemon commits).
6. **`safe.directory`/worktree quirks** on stoptube/papdashboard/autocart/buildflow block `git log` and will block batch-flip preflights — small fix, real payoff.
7. **Status of tool-generated files in repos**: run repo-mutating tools from a scratch cwd (the `output.json` leak), or add it to `.gitignore`.

## f) Up to 50 things we should get done next (ordered by impact; owner/agent tagged)

1. **Owner**: Phase A deploy (M1) — canonicalRepos, shadow/protect binaries, ssh-keys fix; F1–F4.
2. **Owner**: Phase B flip pilot (M2) — `forgejo-flip-check@dnsblockd` → `forgejo-flip@dnsblockd` → sync phase 5; F5–F8.
3. **Agent**: run `forgejo-post-flip-verify dnsblockd` immediately after the flip — the harness now exists and is fixture-green (M3 F11).
4. **Owner**: Phase C secret `RELEASE_SHADOW_GITHUB_TOKEN` on dnsblockd (M4/F16–F17).
5. **Agent+owner**: live CI proof — push to dnsblockd, watch `native` runner run BuildFlow gate (M5/F18–F23).
6. **Agent+owner**: live release proof — tag v0.9.5, watch BuildFlow → goreleaser (`GITEA_TOKEN`) → shadow → GitHub assets (M6/F24–F30).
7. **Agent**: M7 pilot retro — annotate master plan §6, batch doc Phase F evidence links (F31–F35).
8. **Owner**: Phase E token handoff (pass insert forgejo/lars-api-token) — unblocks agent-side verification.
9. **Agent**: M10 authenticated mirror sweep (`has_actions` on dnsblockd-pre-flip/monitor365/wireguard-collector mirrors; F46–F48).
10. **Agent**: M14 provider-aware CLI success tips (forgejo/dual → `.forgejo/workflows`; F66–F69).
11. **Agent**: M15 `--forgejo-owner` + `--runner-label` flags + precedence tests (template fields already exist — flags only; F70–F75).
12. **Agent**: M16 wizard-generated `.forgejo/actionlint.yaml` (native label) + golden (F76–F80).
13. **Agent**: M16b dual-provider `goreleaser check` E2E (§f-16) — currently only forgejo+github goldens are check-proven.
14. **Agent**: M16c `NewForgejoWorkflowsTemplateData` pin-field unit tests (§f-17).
15. **Agent**: M18 docs — README CI-providers section, FEATURES Forgejo DONE, DOMAIN_LANGUAGE vocabulary (F85–F90).
16. **Agent**: M11 `forgejo-flip-batch@` design + implement + `--dry-run` twin + fixtures (F49–F56).
17. **Agent**: M12 wave-1 preflight generator for the 7 Phase-F wave-1 repos (F57–F61).
18. **Agent**: fix `git log` failures on stoptube/papdashboard/autocart/buildflow (safe.directory/worktree) — unblocks wave-3 preflights.
19. **Agent**: M26 `OrEmpty()` deletion (zero callers verified twice now; F122).
20. **Agent**: M26 verdict on SystemNix `scripts/io-writer-attribution.sh` (F123).
21. **Agent**: M26 judge dnsblockd's daemon-touched `docs/status/README.md` (F124).
22. **Agent**: M26 stale `nlreturn`/`wsl_v5` LSP warnings vs real golangci-lint authority (§f-47).
23. **Agent**: M25 lessons into crush-config `references/lessons.md`: fixture-assert-parse-don't-grep; pipefail swallowing command-substitution rc; **never pipe verification gates** (this session's addition); stub-matrix-before-scenarios.
24. **Agent**: M25 SystemNix runbook — mirror-Actions answer + positive control (F119).
25. **Agent**: F120 goreleaser `use: gitea` vs `git` verification against v2.17 docs — switch if PR attribution wins.
26. **Agent**: M28 F133 — wizard v2 CHANGELOG + version prep (go-release checklist); tag only on approval.
27. **Agent**: M19 branch-protect rollout verification across the 6 canonicalRepos after first phase-6 run (+ OnFailure paging check; F91–F94).
28. **Agent+owner**: M20 push-mirror × branch-protection interplay proof (F95–F98).
29. **Agent**: M21 secret-provisioning loop script (list API idempotency + PUT; F99–F102; owner runs it).
30. **Agent**: M22 shadow drift detector — Gatus check Forgejo tag == GitHub twin tag (F103–F107).
31. **Agent**: depgraph `hasGit`/commitsSinceTag enrichment investigation — restore a real activity signal (§b-4).
32. **Agent**: F39 leftover — Mermaid export of own-project coupling from the existing JSON.
33. **Agent**: golden-churn guard — CI asserts `testdata/` stays clean on `-update-golden` runs (hardens M13's proof).
34. **Agent**: M27/F128–F129 — wizard self-hosts its CI on forgejo (dogfood the templates on the pilot infra).
35. **Agent**: M27/F130 — Homebrew tap implications ADR (tap stays GitHub; release-name scheme).
36. **Agent**: M23 capacity measurement once M5/M6 produce wall-time data (F108–F112).
37. **Agent**: M24 attic cache steps in forgejo templates, guarded by env (F113–F116; monitor365 `.forgejo/workflows/nix-cache.yml` is the pattern).
38. **Owner**: g1 — branch-protection strictness policy (tighten vs keep push-open) → folds into M19/M93.
39. **Owner**: g2 — release-shadow PAT policy (reuse sync PAT vs dedicated secret) → folds into M4/M21.
40. **Owner**: g3 — docker-missing CLI behavior (fail/warn/skip) → folds into M17.
41. **Agent**: M17 implement the docker decision + host-without-docker test (F81–F84).
42. **Agent**: M12b preflight checklist generation for wave 1 incl. starred/pending-verdict gates (F58–F59).
43. **Agent**: capture runner wall-times from M5/M6 harness runs → feeds M23 (F108/F30).
44. **Agent**: dnsblockd v0.9.5 release-notes draft (prep for M6 tag).
45. **Agent**: SystemNix runbook post-flip checklist section pointing at `forgejo-post-flip-verify` (runbook proper, beyond the batch doc).
46. **Agent**: `mr list` cwd/config semantics — why 3 rows in the wizard repo; document the correct invocation for F40 closure.
47. **Agent**: §f-37 optional Actions-dashboard tile (low; only if cheap).
48. **Agent**: §f-31 non-wizard workflow shapes (Rust collectors) — build AFTER wave-3 counts confirm demand.
49. **Agent**: F118 commit harvested lessons in crush-config repo.
50. **Agent**: F132 final full gates sweep (build+test+lint+arch+`nix flake check` both repos) once P2 batch lands.

## g) Questions I cannot figure out myself

1. **Did you run privileged Phases A–B since 10-28 — and if not, when do you want to run them?** Everything agent-side is staged for an instant post-flip verify (`forgejo-post-flip-verify dnsblockd`), but I cannot see the live instance without the Phase E token, so I don't know whether the pilot is already flipped.
2. **g2 — release-shadow PAT policy**: reuse the existing repo-scoped sync PAT for the `RELEASE_SHADOW_GITHUB_TOKEN` secret, or provision a dedicated PAT (and if dedicated: one per repo, or one shared shadow PAT)? This decides M4's actual command and the M21 loop design.
3. **g3 — docker-missing CLI behavior**: when the `docker` dependency check can't find docker on the host, should `generate` fail hard, warn-and-continue, or silently skip the docker job? This decides M17 and retires the `/tmp/fakebin/docker` shim gotcha.
