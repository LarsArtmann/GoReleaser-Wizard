# cmdguard — Library Adoption Audit

| Field | Value |
| --- | --- |
| Library | `github.com/larsartmann/cmdguard/v5` (+ submodules: lint `v0.4.0`, prompts/glamour/spinner/telemetry `v0.3.0`, flightrecorder `v0.2.0`, schema `v0.1.0`) |
| Local checkout | `/home/lars/projects/cmdguard` @ `0bf6678` (2026-11-21) |
| Latest local tag | root `v5.1.0` |
| Consumer | GoReleaser-Wizard @ `a493950` (2026-11-22) |
| Used? | **No** — `rg 'larsartmann/cmdguard'` → 0 hits |
| Verdict | NOT-ADOPTED → **CONSIDER** (with a cheaper local fix for the sharpest gap; near-term sub-adoption: `cmdguard/lint`) |
| Score | 0 of 8 applicable capabilities leveraged (0% — utilization band: low, expected for a non-adopted lib) |
| Provenance | Fleet audit Group D, 2026-11-22; citations QC-sampled by orchestrator |

## What the library offers

A typed Cobra wrapper: the root command is generated from a config struct `T` via struct tags (`NewCLI[T]`), each subcommand gets its own typed flag struct `F` (`NewCommand[T,F]`), with flag scoping (persistent/local/per-command), flag→env→config-file precedence (koanf internally), typed middleware, samber/do DI scope, an exit-code contract (`ExecuteAndExit`, `ExitCode(err)`), production batteries (strict config, signals, auto `doctor`/`version`), a `describe` machine contract, and a builder API. The package lives at `pkg/cmdguard/v5/` inside the root module.

| Capability | What it does | Citation | Applicable here? |
| --- | --- | --- | --- |
| Typed root CLI via type param | `NewCLI[T]` turns a config struct into root flags | cmdguard/pkg/cmdguard/v5/cli.go:163 | **APPLICABLE** — consumer hand-registers ~15 flags as strings (main.go:218-224, :296-305) and reads them back via name lookups `getStringFlag`/`getBoolFlag` (generate.go:33-51) |
| Per-command flag structs | `NewCommand[T,F]` per-command typed flags | cmdguard/pkg/cmdguard/v5/command.go:284 | **APPLICABLE** — init/generate/validate each own flags, duplicated via `addCommonFlags` (main.go:296-305) + per-command blocks (validate.go:40-44) |
| Flag scoping | persistent/local/per-command via tags | cmdguard/pkg/cmdguard/v5/cli.go:292 | **APPLICABLE** — consumer's only persistent flags are config/debug/no-color (main.go:218-221); cmdguard formalizes exactly this split |
| Config precedence + files | flag → env → config-file default (koanf internally) | cmdguard/pkg/cmdguard/v5/flags_parse.go:46-71; cmdguard/pkg/cmdguard/v5/koanf_loader.go:12-15; cmdguard/pkg/cmdguard/v5/cli_options.go:287 | **PARTIAL** — consumer already has this via its own koanf Manager (internal/config/config.go:51-117); adopting means *replacing* a working layer |
| Exit-code contract | `ExecuteAndExit` + `ExitCode(err)` + `NewExitError`; error printed exactly once | cmdguard/pkg/cmdguard/v5/cli.go:854; cmdguard/pkg/cmdguard/v5/cli_errors_json.go:73; cmdguard/pkg/cmdguard/v5/errors.go:246 | **APPLICABLE** — fixes a real gap (see below) |
| Production batteries | strict validation, signals, auto `doctor`/`version` | cmdguard/pkg/cmdguard/v5/cli_batteries.go:63-64 | **PARTIAL** — version/strict yes; `doctor` is meaningless for a file generator |
| Builder API | Error-collecting command registration | cmdguard/pkg/cmdguard/v5/cli_builder.go:35, :117 | **APPLICABLE** — marginal at 4 commands (main.go:214-231) |
| `describe` machine contract | JSON CLI contract for agents/tests | cmdguard/pkg/cmdguard/v5/describe.go:122, :138 | N/A — surface pinned by golden/E2E tests; no agent-driving requirement |
| prompts submodule (huh/v2) | Single-field `prompt:"..."` tag fallbacks | cmdguard/prompts/prompts.go:26-27 | N/A — consumer's wizard is a multi-group huh form with cross-field state (tui_wizard.go:132, :153-358); a capability *downgrade* |
| lint submodule | CG001-CG008 usage linter | cmdguard/lint/rules.go:17-24 | **APPLICABLE standalone** — CG001/CG006 would flag real defects in the consumer today |
| Middleware | Typed `Middleware[T]` | cmdguard/pkg/cmdguard/v5/middleware.go:26 | N/A — zero middleware concept exists (`rg 'Middleware'` non-test → 0 hits) |
| DI Scope | samber/do-backed lazy services | cmdguard/pkg/cmdguard/v5/scope.go:21 | N/A — no services to inject; single-shot binary |
| glamour/spinner/telemetry/flightrecorder/schema | Markdown help, spinner, OTel, trace recorder, JSON Schema | glamour/glamour.go:30; spinner/spinner.go:59; telemetry/telemetry.go:34; flightrecorder/middleware.go:114; schema/schema.go:125 | N/A — no long-running ops, no telemetry consumers |

## What GoReleaser-Wizard does today

- **Registration**: manual Cobra — persistent flags at main.go:218-224, per-command string flags at main.go:296-330 and validate.go:40-44; values fetched by name (`getStringFlag`, generate.go:33-51). Handlers are `Run:` funcs returning void (validate.go:46, generate.go:12, init.go:17) — no `RunE`.
- **Error rendering**: `fang.Execute(ctx, rootCmd, …)` at main.go:138, **then** `displayError(err)` re-renders the same error with lipgloss at main.go:145 within the `if err != nil` block at :144-146 — the double-print pattern cmdguard's CG006 rule exists to catch.
- **Exit codes (the real defect)**: because handlers swallow errors via `displayError(err); return` (generate.go:26-30, displayError at :27), a failed `generate`/`init` run **exits 0**; only `validate` sets a real exit code (`os.Exit(exitCode)` at validate.go:48, codes 1/2 from `GetExitCode()`, validation_results.go:18-28) and panics exit 1 (main.go:163).
- **Prompts**: huh/v2 directly (charm.land/huh/v2 v2.0.3) in a 435-line wizard (tui_wizard.go) — richer than cmdguard's prompts submodule.
- **Dependency overlap**: cmdguard pins the same cobra v1.10.2 (cmdguard/go.mod:86) and koanf/v2 v2.4.1 (:78) and prompts pins the same huh/v2 v2.0.3 — but cmdguard uses **`charm.land/fang/v2 v2.0.1`** (:73) while the consumer uses **`github.com/charmbracelet/fang v1.0.0`**: different module paths, so adoption ships two fang majors unless the consumer also migrates fang v1→v2.

## Why CONSIDER

- **The exit-code gap is real and valuable to close**: failed `generate`/`init` runs exit 0 today — a CI-scripting bug for a tool whose `generate` path exists *for* scripting (DOMAIN_LANGUAGE.md).
- **But the cheaper 90% fix is local**: switch handlers to `RunE`, propagate the error, exit on it — no new dependency, no test-harness rewrite. Do that regardless of cmdguard.
- **Dependency duplication**: fang v1 (consumer) vs `charm.land/fang/v2` (cmdguard) would coexist in one binary.
- The flagship UX (huh wizard) is better than cmdguard's prompts submodule; the config layer is already solved with the same koanf versions — cmdguard would replace, not enhance.
- **Best near-term move**: adopt **`cmdguard/lint` standalone** into the lint pipeline — zero runtime coupling, directly targets the two defects above (CG001 execute-bypass, CG006 error re-print).

## Migration sketch (if adopted)

Mechanical: define root `T` + three `F` structs (~15 flags), convert `Run:` handlers to `(ctx, *T, F) error`, swap `fang.Execute` for `ExecuteAndExit`, delete the manual koanf Manager (or keep it and use only flag features — overlapping config layers). **Effort: ~2-4 focused days.** Risks: golden/E2E tests drive cobra commands and `performValidation` directly (validate.go:52-53 exists to avoid `os.Exit` in tests) and need rewiring; the fang v1/v2 duplication. The Job/Workflow engine is orthogonal and survives untouched.

## RECONSIDER-WHEN (flip to full ADOPT)

The command layer changes shape: subcommands added, flag count grows past ~20, a `describe`-style agent contract or `--format json` machine output is requested, or the fang v1→v2 migration happens anyway — then migrate in one pass (`ExecuteAndExit` + deleting the manual koanf Manager).

## How to verify

```bash
cd /home/lars/projects/GoReleaser-Wizard
sed -n '24,31p' cmd/goreleaser-wizard/generate.go    # Run: handler swallows error via displayError (exit 0)
sed -n '138,148p' cmd/goreleaser-wizard/main.go      # fang.Execute THEN displayError — double print
sed -n '48p' cmd/goreleaser-wizard/validate.go       # the only real os.Exit(exitCode)
rg -c 'larsartmann/cmdguard' .                        # 0 hits
grep -n 'fang' go.mod                                  # github.com/charmbracelet/fang v1.0.0
grep -n 'fang\|cobra\|koanf/v2 ' /home/lars/projects/cmdguard/go.mod | head   # fang/v2 v2.0.1, cobra v1.10.2, koanf v2.4.1
git -C /home/lars/projects/cmdguard tag --list 'v*' --sort=-v:refname | head -1   # v5.1.0
```
