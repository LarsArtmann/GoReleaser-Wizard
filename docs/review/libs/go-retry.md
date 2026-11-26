# go-retry — Library Adoption Audit

| Field            | Value                                                                 |
| ---------------- | --------------------------------------------------------------------- |
| Library          | `github.com/larsartmann/go-retry`                                     |
| Local checkout   | `/home/lars/projects/go-retry` @ `356fbd4` (2026-11-14)               |
| Latest local tag | `v0.9.0`                                                              |
| Consumer         | GoReleaser-Wizard @ `a493950` (2026-11-22)                            |
| Used?            | **No** — `rg 'larsartmann/go-retry'` → 0 hits                         |
| Verdict          | NOT-ADOPTED → **STAY-OUT**                                            |
| Score            | n/a — 0 applicable capabilities (7 rows, all N/A)                     |
| Provenance       | Fleet audit Group A, 2026-11-22; citations QC-sampled by orchestrator |

## What the library offers

A dependency-light retry loop with exponential backoff + jitter, context awareness, hooks, and generic value-returning variants.

| Capability                            | What it does                                                                     | Citation                                   | Applicable here?                                                                                   |
| ------------------------------------- | -------------------------------------------------------------------------------- | ------------------------------------------ | -------------------------------------------------------------------------------------------------- |
| `retry.Do(ctx, cfg, fn)`              | Retry loop with exponential backoff + jitter, ctx-aware                          | go-retry/retry.go:99                       | N/A — zero retry loops exist in the consumer                                                       |
| `DoWithValue[T]`                      | Generic retry returning the successful attempt's value                           | go-retry/retry.go:159                      | N/A — no fallible value-producing calls that could succeed on re-attempt                           |
| `Config` / `DefaultConfig()`          | MaxAttempts/InitialDelay/MaxDelay/Multiplier/IsRetryable                         | go-retry/config.go:20, :73                 | N/A                                                                                                |
| `FromPolicy(errorfamily.RetryPolicy)` | Bridge from error-family advisory policy                                         | go-retry/config.go:86                      | N/A — would drag in go-error-family (STAY-OUT); it is go-retry's only dependency (go-retry/go.mod) |
| Options: hooks & jitter               | `WithIsRetryable`, `WithDelayFunc`, `WithOnRetry`, `WithExhausted`, `WithJitter` | go-retry/options.go:25, :32, :38, :44, :74 | N/A                                                                                                |
| Sentinel errors                       | `ErrExhausted` (unwraps to last attempt), `ErrCanceled`, `ErrDeadlineExceeded`   | go-retry/retry.go:45, :55, :68             | N/A — no error-driven control flow branches on these                                               |
| Pure delay math                       | `Backoff(config, attempt)`, `ComputeDelay(...)`                                  | go-retry/retry.go:275, :290                | N/A                                                                                                |

## What GoReleaser-Wizard does today

No retry machinery at all. `rg 'time\.Sleep|attempt' --type go -g '!*_test.go'` returns only naming coincidences (`attemptFixes` in cmd/goreleaser-wizard/validation_fixes.go:21 — a one-shot "apply known fixes" pass, not a loop). The retry-adjacent artifacts are declarative and unconsumed: `DomainError.IsRetryable()` (internal/domain/errors.go:434, no callers) and a `Retryable bool` JSON field in cmd/goreleaser-wizard/types/template_data.go:248 (rendered output metadata). All real operations are single-shot: template render → preflight → `os.WriteFile` (cmd/goreleaser-wizard/generators/template_utils.go:24) → local `git` exec (internal/git/commands.go:57).

## Why STAY-OUT

- **No retryable operation exists**: the runtime profile is a short-lived single-shot process doing local file writes and one local git exec; retrying a failed local `os.WriteFile` returns the same ENOSPC/EPERM.
- **Fleet policy conflict**: the fleet resilience policy is failsafe-go; go-retry's own code positions it as the dependency-light alternative and defers loop policy to "retry libraries like failsafe-go" (go-retry/retry.go:22-24). If retries ever appear, the fleet standard wins.
- **Transitive coupling**: go-retry's only dependency is go-error-family v0.12.0 (go-retry/go.mod) — adopting it implicitly imports the taxonomy rejected in the go-error-family verdict via the default `IsRetryable` predicate.
- Well-engineered (fuzz/property tests, error preservation on `ErrExhausted`), but it solves a problem this binary does not have.

## RECONSIDER-WHEN

`generate` gains network calls (e.g. a Forgejo/GitHub API owner-repo resolution step or release-note fetching) AND the fleet has not standardized on failsafe-go for this class of binary — then `retry.Do` with `WithIsRetryable` mapped to `DomainError.IsRetryable`.

## How to verify

```bash
cd /home/lars/projects/GoReleaser-Wizard
rg -n 'time\.Sleep|attempt' --type go -g '!*_test.go'        # no retry loops; only attemptFixes (validation_fixes.go:21)
rg -c 'larsartmann/go-retry' .                                # 0 hits
sed -n '434,441p' internal/domain/errors.go                   # IsRetryable exists, zero callers
git -C /home/lars/projects/go-retry tag --list 'v*' --sort=-v:refname | head -1   # v0.9.0
head -12 /home/lars/projects/go-retry/go.mod                  # sole dep: go-error-family v0.12.0
```
