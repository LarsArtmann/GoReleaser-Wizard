# Wave 1 Execution — Brutal Status Report

**Date:** 2026-08-17 00:04
**Session:** Wave 1 (M1–M4 / F1–F18) + opportunistic F27/F30/F32 of the SUPERB Wizard Modernization plan
**Commits this session:** `e85bb4f` (Wave 1 core), `3282680` (CI template fixes) — **local only, NOT pushed**
**Plan:** `docs/planning/2026-08-16_18-45_SUPERB-WIZARD-MODERNIZATION.md`

---

## a) FULLY DONE (verified this session)

| #                                              | Item                                                                                                                                                                                                                                                   | Proof                                    |
| ---------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------- |
| ~~F1~~ done at `e85bb4f`, `3282680`            | ~~Template key/data-source inventory (jobs.go, embedded.go, template_data.go, dockerfile.go, git/commands.go)~~                                                                                                                                        | ~~read in full~~                         |
| ~~F2~~ done at `e85bb4f`, `3282680`            | ~~`archives.format_overrides.format` → `formats: [zip]`~~                                                                                                                                                                                              | ~~`goreleaser check`~~                   |
| ~~F3~~ done at `e85bb4f`, `3282680`            | ~~`dockers`+`build_flag_templates` → `dockers_v2` (id, images, tags w/ `IsNightly` guard, platforms, OCI annotations)~~                                                                                                                                | ~~check exit 0, zero deprecations~~      |
| ~~F6~~ done at `e85bb4f`, `3282680`            | ~~Deprecated `release.mode: append` dropped~~                                                                                                                                                                                                          | ~~check output~~                         |
| ~~F7–F8~~ done at `e85bb4f`, `3282680`         | ~~Rendered output passes `goreleaser check` **exit 0, zero deprecations, NO env vars**~~                                                                                                                                                               | ~~manual E2E, twice~~                    |
| ~~F9–F11~~ done at `e85bb4f`, `3282680`        | ~~Git-remote owner/repo detection rendered as literals; `--github-owner`/`--github-repo` flags on `init`+`generate`~~                                                                                                                                  | ~~flag rendered `owner: "LarsArtmann"`~~ |
| ~~F13–F17~~ done at `e85bb4f`, `3282680`       | ~~`DockerfileGenerationJob` wired into `CreateFullWizardJobs` when Docker enabled; prebuilt pattern Dockerfile (`FROM scratch`, `ARG TARGETPLATFORM`, `COPY $TARGETPLATFORM/<bin>`; alpine fallback when CGO); rollback implemented; `--force` guard~~ | ~~job runs, file exists~~                |
| ~~F18~~ done at `e85bb4f`, `3282680`           | ~~Full `goreleaser release --snapshot --clean` succeeds **through Docker image builds** (amd64 + arm64, digests emitted)~~                                                                                                                             | ~~run log, 14s, clean~~                  |
| ~~F27~~ done at `e85bb4f`, `3282680`           | ~~GH Actions `runs-on: ubuntu-latest`~~                                                                                                                                                                                                                | ~~rendered YAML~~                        |
| ~~F30~~ done at `e85bb4f`, `3282680`           | ~~`docker/setup-buildx-action` step when Docker enabled~~                                                                                                                                                                                              | ~~rendered YAML~~                        |
| ~~F32~~ done at `e85bb4f`, `3282680`           | ~~`before.hooks` = `go mod tidy` only~~                                                                                                                                                                                                                | ~~golden eyeball~~                       |
| ~~Gates~~ done at `e85bb4f`, `3282680`         | ~~`go build`, `go test ./...`, `golangci-lint` 0 issues, `go-arch-lint` OK, BuildFlow pre-commit passed (warnings pre-existing)~~                                                                                                                      | ~~all green~~                            |
| ~~Tests updated~~ done at `e85bb4f`, `3282680` | ~~`generate_test.go`, `generate_extended_test.go` expectations moved to `dockers_v2`~~                                                                                                                                                                 | ~~suite green~~                          |

---

## b) PARTIALLY DONE

| Item                                                                                                | Done                                             | Missing                                                                                                                                                                                                                                  |
| --------------------------------------------------------------------------------------------------- | ------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| ~~M2 owner/repo~~ done at `6f4f54d`                                                                 | ~~detection, flags, env-tolerant check~~         | ~~**no unit tests** for overrides/detection; silent `owner:"owner"` placeholder when no git remote and no flags (check passes, release would target a fake repo — no warning printed)~~                                                  |
| ~~M3 Dockerfile job~~ done at `6f4f54d`                                                             | ~~full-wizard path wired~~                       | ~~**NOT wired into `CreateConfigOnlyJobs`** → `generate --config-only` + Docker still emits config referencing a Dockerfile that is never generated — the exact defect class Wave 1 was supposed to kill, surviving in a narrower path~~ |
| ~~M7 GH Actions~~ **Superseded — F33 → ROADMAP.md (TypeSpec-blocked).**                             | ~~runner, buildx, login+cosign conditionals~~    | ~~F31 (render 4 docker×signing combos, YAML-parse each) **not done**; F33 (opt-in test/generate hooks) **not done** — I removed hooks outright instead of making them configurable~~                                                     |
| ~~M4 Dockerfile template~~ **Won't implement — distroless variant dropped; scratch+alpine chosen.** | ~~scratch + TARGETPLATFORM + binary name wired~~ | ~~distroless variant skipped; WebAPI EXPOSE 8080 hardcoded; no CA-cert handling for HTTPS CLIs from scratch~~                                                                                                                            |
| ~~F12 check-without-env~~ done at `8edf0b1`                                                         | ~~manually verified~~                            | ~~**not automated** (manual /tmp runs, since deleted nothing — /tmp/gw-e2e, /tmp/gw-cfg.yaml still around)~~                                                                                                                             |
| ~~M1 template modernization~~ done at `9f4bd71`                                                     | ~~live inline template modern~~                  | ~~dead copies (`templates/embedded.go`, `templates/*.tmpl`, `generators/`) untouched — Wave 2 target~~                                                                                                                                   |

---

## c) NOT STARTED (per plan)

- **Wave 2:** M5 (wire `generators` pkg into jobs), M6 (delete inline templates + `map[string]any` prep), F19–F26, F31, F33
- **Wave 3:** M9 automated E2E (`generate` → `goreleaser check` in test), M10 golden files, M11 CI workflow (`ci.yml` + `goreleaser check` step + `-race`)
- **Wave 4:** M12 hygiene — 3 `.orig` files still tracked, unused `createTempFile`, `infertypeargs` hints
- **Wave 5:** M13 GoReleaser guide (Goal A), M14 README, M15 AGENTS.md (still `just`-based, no jsonv2/cache gotchas; BuildFlow flags it 31 days stale), M16 FEATURES.md
- **Wave 6:** M17–M27 (error-discarding fixes, personal defaults, coverage lifts, `validate` deprecation detection, own `.goreleaser.yaml` migration — still fails check on `dockers`/`brews`, flake cache robustness, TODO_LIST/ROADMAP, v1.0.0 release, post-release verify)

---

## d) TOTALLY FUCKED UP (session honesty log)

~~1. **Fake verification #1:** first E2E used `go run github.com/.../cmd/goreleaser-wizard` which resolved the **published module from the proxy**, not my edited code. I initially read its (old) output as my result. Caught only because the diff looked unchanged. Lesson: build local binary first, always.~~ done — session honesty log; resolved by 6f4f54d/9f4bd71/8edf0b1/3cb91da
~~2. **Silent multiedit partial failure:** 2 of 3 edits applied; the critical template edit failed on trailing-whitespace lines and I didn't notice until E2E. Then my retry put the closing backtick on the wrong line and failed again. 3 wasted round trips from sloppy whitespace copying.~~ done — session honesty log; resolved by 6f4f54d/9f4bd71/8edf0b1/3cb91da
~~3. **First-write logic bug:** `dockerPlatforms()` v1 returned early and never appended `linux/arm64`. Caught on re-read, but I shipped a bug in 10 lines of fresh logic.~~ done — session honesty log; resolved by 6f4f54d/9f4bd71/8edf0b1/3cb91da
~~4. **Created a syntax error** (comment merged onto `func parseGitHubRemote()` line) via careless edit; build caught it.~~ done — session honesty log; resolved by 6f4f54d/9f4bd71/8edf0b1/3cb91da
~~5. **Split-brain increased, deliberately:** added a **third** Dockerfile template (new inline const in jobs.go) instead of wiring the existing dead `DockerfileGenerator`. Justifiable only because M5/M6 consolidates everything next wave; if Wave 2 stalls, the codebase is worse than before.~~ done — session honesty log; resolved by 6f4f54d/9f4bd71/8edf0b1/3cb91da
~~6. **Zero new automated tests.** Every new symbol (`dockerPlatforms`, owner/repo overrides, `DockerfileGenerationJob`, rollback) is verified by hand only. The plan's own words: "correctness is proven, not hoped." I hoped.~~ done — session honesty log; resolved by 6f4f54d/9f4bd71/8edf0b1/3cb91da
~~7. **Doc drift created:** plan file still says `Status: PLANNED — awaiting execution`. I didn't update it after executing a wave.~~ done — session honesty log; resolved by 6f4f54d/9f4bd71/8edf0b1/3cb91da
~~8. **Ignored BuildFlow signals:** `go mod tidy` needed (root module) — not run; `GOEXPERIMENT=jsonv2 possibly redundant` — not verified whether the build still needs it.~~ done — session honesty log; resolved by 6f4f54d/9f4bd71/8edf0b1/3cb91da
~~9. **Job-order partial state:** Dockerfile job runs LAST; without `--force`, a pre-existing Dockerfile fails the workflow after `.goreleaser.yaml` + `release.yml` are already written. Did not verify JobManager's rollback-on-later-failure semantics.~~ done — session honesty log; resolved by 6f4f54d/9f4bd71/8edf0b1/3cb91da
~~10. **TUI path unverified:** `tui_wizard.go` presumably flows through JobFactory, but I never ran the interactive wizard to confirm the Dockerfile job appears/behaves there.~~ done — session honesty log; resolved by 6f4f54d/9f4bd71/8edf0b1/3cb91da

---

## e) WHAT WE SHOULD IMPROVE

~~1. **Prove, don't eyeball:** every template change should land WITH a test that renders + validates (YAML-parse at minimum, `goreleaser check` exec when available). Wave 3's E2E harness should be pulled forward — it retro-covers Wave 1.~~ done at `8edf0b1`
~~2. **Kill the split brain NOW, not later:** M5/M6 is the highest-leverage next step; two full waves of drift have already accumulated (dead generators/templates/embedded + repo-root `templates/`).~~ done at `9f4bd71`
~~3. **Pre-flight existence checks:** check ALL target files (`.goreleaser.yaml`, `Dockerfile`, `release.yml`) for existence/force BEFORE generating anything → fail atomically, no partial state.~~ done at `6f4f54d`
~~4. **No silent placeholders:** if git detection fails and no flags given, print a warning and/or require `--github-owner/--github-repo`. A config that checks green but releases to `owner/repo` is a trap.~~ done at `6f4f54d`
~~5. **Flags through types, not package globals:** the `githubOwnerOverride` globals (nolint'ed) are startup-only mutable state; thread overrides through config/TemplateData when M6 lands.~~ done at `9f4bd71`
~~6. **Clean unused template data:** after `dockers_v2`, wizard-side `Version/Tag/Major/Date/FullCommit` data keys and `Os`/`Arch` placeholders are dead — remove during M6.~~ done at `9f4bd71`
~~7. **Trust-but-verify tooling:** BuildFlow's preflight warnings (tidy, GOEXPERIMENT redundancy) are cheap to resolve — stop leaving them red.~~ done — BuildFlow signals resolved by later runs

---

## f) NEXT — up to 50 concrete tasks (rough priority order)

**Close Wave 1 gaps (quick wins)**

~~1. Wire `DockerfileGenerationJob` into `CreateConfigOnlyJobs` (+ test)~~ done at `6f4f54d`
~~2. Pre-flight existence check for all three artifacts before any generation~~ done at `6f4f54d`
~~3. Unit tests: `dockerPlatforms` (arm64 on/off), owner/repo override precedence, `DockerfileGenerationJob` execute/rollback/force~~ done at `6f4f54d`
~~4. Warn (or error) on placeholder owner/repo in generated config~~ done at `6f4f54d`
~~5. Update plan file status + wave-1 checkboxes to reflect reality~~ done — plan header updated to EXECUTING
~~6. Run `go mod tidy`; verify whether `GOEXPERIMENT=jsonv2` is still required; document truth~~ done at `6f4f54d`
~~7. F31: render release.yml for all 4 docker×signing combos; YAML-parse each (test)~~ done at `6f4f54d`
~~8. Manual TUI run to confirm Dockerfile job in interactive path~~ done — structural verification (2026-08-17 01-15 report Gap 7)

**Wave 2 — single source of truth (M5/M6)**
~~9. F19: `ConfigGenerationJob` → `generators.NewGoReleaserGenerator`~~ done at `9f4bd71`
~~10. F20: `GitHubActionsGenerationJob` → `GitHubActionsGenerator`~~ done at `9f4bd71`
~~11. F21: JobFactory constructs generators with `LoggerAdapter`~~ done at `9f4bd71`
~~12. F22: full suite green after wiring~~ done at `9f4bd71`
~~13. F23: delete `goreleaserTemplateContent`/`githubActionsTemplateContent` (build FIRST per AGENTS rule)~~ done at `9f4bd71`
~~14. F24: delete `prepareGoReleaserData`/`prepareGitHubActionsData`~~ done at `9f4bd71`
~~15. F25: delete dead helpers (`getVersion`, `getCommitHash`, `runGitCommand`, `addDockerConfig`, `mapToStrings` if unused)~~ done at `9f4bd71`
~~16. F26: grep evidence — no stray `map[string]any` prep in cmd/~~ done at `9f4bd71`
~~17. Port my Wave-1 inline `dockerfileTemplate` into `generators.DockerfileGenerator`; delete the dead Dockerfile template from `templates/embedded.go` + repo-root `templates/` dir (or wire them — decide once)~~ done at `9f4bd71`
~~18. Modernize `templates/embedded.go` GoReleaserTemplate to v2.17 during the port (it still has old `dockers`)~~ done at `9f4bd71`
~~19. Replace override globals with typed TemplateData fields~~ done at `9f4bd71`
~~20. F33: config option for opt-in `go test`/`go generate` before-hooks~~ **Superseded — ROADMAP.md (TypeSpec-blocked).**

**Wave 3 — automated proof**
~~21. F34: E2E scaffolding — temp Go module in `t.TempDir()`~~ done at `8edf0b1`
~~22. F35: programmatic generation into temp module~~ done at `8edf0b1`
~~23. F36: exec `goreleaser check` (skip gracefully if binary missing); assert exit 0 + zero deprecations~~ done at `8edf0b1`
~~24. F37: wire E2E into `go test ./...`~~ done at `8edf0b1`
~~25. F38: golden-file layout + `-update` flag helper~~ done at `8edf0b1`
~~26. F39–F42: golden tests — `.goreleaser.yaml` (2 variants), `release.yml` (matrix), `Dockerfile` (scratch/alpine)~~ done at `8edf0b1`
~~27. F43: author `.github/workflows/ci.yml` (lint+test+build, `-race`)~~ done at `8edf0b1`
~~28. F44: CI step `goreleaser check` on own config + generated fixture~~ done at `8edf0b1`
~~29. F45: keep own `release.yml` in sync after M23~~ done at `8edf0b1`

**Wave 4 — hygiene**
~~30. F46: `git rm` 3 `.orig` files; `.gitignore *.orig`~~ done at `3cb91da`
~~31. F47: remove unused `createTempFile` (validate_test.go)~~ done at `3cb91da`
~~32. F48: fix 5 `infertypeargs` hints (internal/domain/ids.go)~~ done at `3cb91da`
~~33. Delete or repair `test-wizard/` module (invalid module path, BuildFlow error)~~ done at `3cb91da`

**Wave 5 — knowledge & docs (Goal A)**
~~34. F49–F53: `docs/goreleaser-guide.md` (phases, dockers_v2 deep-dive, casks migration, template vars, snapshot vs release)~~ done at `969259e`
~~35. F54–F56: README truth-up (What It Generates, quick start re-run, flags table incl. new flags)~~ done at `969259e`
~~36. F57–F58: AGENTS.md real commands + jsonv2/cache gotchas~~ done at `969259e`
~~37. M16: FEATURES.md refresh~~ done at `969259e`
~~38. Update AGENTS.md with Wave-1 architectural facts (new job, flags)~~ done at `969259e`

**Wave 6 — hardening & release**
~~39. M17: fix error-discarding sites (flags.go, init.go:248, validate_main.go:167, parseGitHubRemote)~~ **Superseded — TODO_LIST.md (M17).**
~~40. M18: parameterize tap repo / NUR / commit author; `--runner` option~~ **Superseded — TODO_LIST.md (M18, BLOCKED).**
~~41. M19: internal/git tests (0% → solid)~~ **Superseded — TODO_LIST.md (M19-M21).**
~~42. M20: generators/templates/types/config tests~~ **Superseded — TODO_LIST.md (M19-M21).**
~~43. M21: domain coverage 10% → 60%+~~ **Superseded — TODO_LIST.md (M19-M21).**
~~44. M22: `validate` detects deprecated keys + migration hints~~ **Superseded — TODO_LIST.md (M22).**
~~45. M23: migrate own `.goreleaser.yaml` to v2.17 + snapshot verify (currently fails check!)~~ done at `969259e`
~~46. M24: flake.nix cache fallback when `/mnt/buildcache` dead~~ **Superseded — TODO_LIST.md (M24).**
~~47. M25: refresh TODO_LIST.md / ROADMAP.md~~ done — docs-health sweep 2026-11-07: TODO_LIST.md + ROADMAP.md created
~~48. Pin `anchore/sbom-action/download-syft@v0` subaction to SHA (BuildFlow unfixable warning)~~ **Won't implement — upstream 404; known-unpinnable.**
~~49. M26: CHANGELOG cut → tag v1.0.0 → release → verify artifacts/docker/cask~~ **Superseded — TODO_LIST.md (M28 release prep, tag gated).**
~~50. M27: post-release verification (tap cask, `brew install`, `docker pull`)~~ **Superseded — TODO_LIST.md (M28 release prep, tag gated).**

---

## g) Questions I cannot answer myself

~~1. **Runner default:** I replaced hardcoded `runs-on: [self-hosted, linux, x64]` with `ubuntu-latest` (per plan default). Was the self-hosted runner intentional for your infrastructure? Should M18's `--runner` flag revive self-hosted as an option?~~ done — ubuntu-latest kept (3282680); commits pushed
~~2. **Push now or hold?** Commits `e85bb4f` + `3282680` are local only. Push immediately, or hold until Wave 2 (split-brain removal) lands so master never carries the third-Dockerfile-template state?~~ done — ubuntu-latest kept (3282680); commits pushed
~~3. **Personal defaults (M18):** for generated Homebrew tap / NUR / commit author — what values? (e.g. tap repo pattern `LarsArtmann/homebrew-<project>`? commit author `LarsArtmann <...?>`) I refuse to guess emails.~~ **Superseded — TODO_LIST.md (M18, BLOCKED).**

---

**Bottom line:** North Star for Wave 1 is met and manually proven (`check` exit 0 / zero deprecations / no env; snapshot release incl. Docker). The honest costs: zero new automated tests, one deliberately-added split brain, and a config-only path still shipping the Dockerfile defect. Wave 2 must land next or the debt compounds.
