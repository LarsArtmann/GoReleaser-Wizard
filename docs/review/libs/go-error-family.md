# go-error-family — Library Adoption Audit

| Field | Value |
| --- | --- |
| Library | `github.com/larsartmann/go-error-family` (+ submodules `agent`, `bridge`, `diagnose`, `diagnose/git`, `diagnose/postgres`) |
| Local checkout | `/home/lars/projects/go-error-family` @ `898004e` (2026-11-15) |
| Latest local tags | root `v0.12.0` · agent `v0.2.5` · bridge `v0.3.4` · diagnose `v0.2.6` · diagnose/git `v0.5.4` · diagnose/postgres `v0.5.4` |
| Consumer | GoReleaser-Wizard @ `a493950` (2026-11-22) |
| Used? | **No** — `rg 'larsartmann/go-error-family'` → 0 hits; go.mod lists no such module |
| Verdict | NOT-ADOPTED → **STAY-OUT** |
| Score | n/a — 0 applicable capabilities (11 rows, all N/A with reasons below) |
| Provenance | Fleet audit Group A, 2026-11-22; citations QC-sampled by orchestrator |

## What the library offers

Error-family is a zero-dependency behavioral error taxonomy: every error is classified into one of six families (Rejection, Conflict, Transient, Corruption, Infrastructure, Orchestration), and the family drives retry decisions, CLI exit codes, HTTP status mapping, and user-facing tone.

| Capability | What it does | Citation | Applicable here? |
| --- | --- | --- | --- |
| Family classification | 6 behavioral families driving retry/exit-code/tone decisions | go-error-family/family.go:29-30 (`const` block, `Rejection = iota`) | N/A — consumer already has a richer, project-specific taxonomy: 90+ `ErrorCode` constants in `internal/domain/errors.go:26-119` |
| `Classify(err)` / `IsRetryable(err)` / `ExitCode(err)` | Universal classification incl. `errors.Join` worst-of | go-error-family/classify.go:61, :67, :78 | N/A — nothing branches on retryability; `DomainError.IsRetryable()` (internal/domain/errors.go:434) has zero callers |
| `HandleError(err) int` CLI boundary | What/Why/Fix/WayOut output + derived exit code | go-error-family/handle.go:121 | N/A — boundary is `cmd/goreleaser-wizard/main.go:134-148` (`fang.Execute` + `displayError`); fang owns CLI error presentation |
| Reference `Error` type | `[family:code] message: cause` | go-error-family/error.go:22, :37 | N/A — direct duplicate of the consumer's `DomainError` (internal/domain/errors.go:10: Code/Message/Details/Context/Field/Suggestion/Cause) |
| Constructors + `WrapOnce` | `NewRejection`/`WrapTransientf`/…, idempotent double-wrap guard | go-error-family/constructors.go:57, :162 | N/A — AGENTS.md gotcha 2 mandates the consumer's own `New*Error()` constructors; a parallel family would split the codebase |
| `Registry` + `Clone()` | Injectable sentinel/classifier registry | go-error-family/registry.go:31, :260 | N/A — tests use testify + `IsErrorCode` (internal/domain/errors.go:170) |
| Stdlib defaults | Pre-registered context/sql/os classifications | go-error-family/stdlib.go:30 | N/A — os/exec git errors are wrapped via `NewGitError` (internal/domain/errors.go:244) |
| `RetryPolicy()` advisory | Attempts + backoff bounds derived from family | go-error-family/retry.go:18, :26 | N/A — no retry loop exists anywhere in the consumer |
| `errorfamilytest` assertions | `AssertFamily`, `AssertExitCode`, … | go-error-family/errorfamilytest/errorfamilytest.go:31 | N/A — testify + `IsErrorCode` already cover error tests |
| Diagnose rules (fs/git/pg) | Deterministic "why did this happen" checks emitting fix commands | go-error-family/diagnose/rules_filesystem.go:17 | N/A — closest analog `GetRecoverySuggestion()` (internal/domain/errors.go:372) is deliberately hand-written per domain code |
| HTTP middleware / JSON view | Family → status mapping | go-error-family/http.go:17, :54 | N/A — consumer has no HTTP server (`rg 'net/http'` → 0 hits) |

## What GoReleaser-Wizard does today

- A full typed error system: `internal/domain/errors.go` — `DomainError` (:10), `Error()`/`Unwrap()` (:122/:131), fluent `WithContext/Field/Suggestion/Cause` (:136-167), `IsErrorCode` (:170), 9 `New*Error` constructors (:179-250), `GetSeverity` (:405), `IsRecoverable` (:427), `IsRetryable` (:434).
- CLI boundary: `cmd/goreleaser-wizard/main.go:134-148` (`Execute()` → `fang.Execute` → `displayError`), panic recovery → `os.Exit(1)` at main.go:163.

## Why STAY-OUT

- **Head-on duplication**: the consumer deliberately built a DomainError system with domain-specific codes, suggestions, and severity — error-family's generic `Family`/`Error` re-implements the same shape without the domain knowledge.
- **Fleet policy conflict**: the fleet's errors policy is cockroachdb/errors (`how-to-golang` required-libraries); adopting a first-party classification lib would introduce a third competing taxonomy with no consumer for the extra classifications.
- **Nothing consumes the behavioral outputs**: no caller of `IsRetryable` semantics; exit codes are a hardcoded `os.Exit(1)`; error presentation is owned by fang/charm.
- **Domain purity**: `internal/domain` is zero-external-deps by arch-lint rule; importing error-family into the domain layer (where errors live) would break the architecture contract.
- The library is high quality (zero deps, clean retracts) — the verdict is fit, not merit.
- Partial adoption (`errorfamilytest` or `ExitCode` at the boundary only) adds a module dependency for trivially inlinable functionality.

## RECONSIDER-WHEN

The wizard's error surface gains behavioral consumers — differentiated exit codes per error class (EX_TEMPFAIL vs EX_USAGE), retry decisions fed into a retry loop, or a fleet-wide standardization on error-family replacing cockroachdb/errors — and the DomainError system is scheduled for retirement rather than coexistence.

## How to verify

```bash
cd /home/lars/projects/GoReleaser-Wizard
rg -c 'larsartmann/go-error-family' . ; rg 'larsartmann/' go.mod        # → 0 hits; only go-branded-id v0.7.1
sed -n '10p;170p;244p;372p;434p' internal/domain/errors.go              # DomainError / IsErrorCode / NewGitError / GetRecoverySuggestion / IsRetryable
sed -n '134,148p;163p' cmd/goreleaser-wizard/main.go                    # fang boundary + os.Exit(1)
git -C /home/lars/projects/go-error-family tag --list 'v*' --sort=-v:refname | head -1   # v0.12.0
```
