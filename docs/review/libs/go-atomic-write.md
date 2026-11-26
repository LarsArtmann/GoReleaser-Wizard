# go-atomic-write — Library Adoption Audit

| Field            | Value                                                                                                 |
| ---------------- | ----------------------------------------------------------------------------------------------------- |
| Library          | `github.com/larsartmann/go-atomic-write`                                                              |
| Local checkout   | `/home/lars/projects/go-atomic-write` @ `4f10347` (2026-11-15)                                        |
| Latest local tag | `v0.6.2`                                                                                              |
| Consumer         | GoReleaser-Wizard @ `a493950` (2026-11-22)                                                            |
| Used?            | **No** — `rg 'larsartmann/go-atomic-write'` → 0 hits                                                  |
| Verdict          | NOT-ADOPTED → **CONSIDER** (the strongest adoption candidate of the 20-audit)                         |
| Score            | 0 of 6 applicable capabilities leveraged (0% — utilization band: low, expected for a non-adopted lib) |
| Provenance       | Fleet audit Group A, 2026-11-22; citations QC-sampled by orchestrator                                 |

## What the library offers

Crash-safe file writes: stage to a unique temp file → fsync → rename, plus optional concurrent-modification detection and skip-if-unchanged semantics. Two direct dependencies (xxhash, flock); `golang.org/x/sys` indirect — already in the consumer's graph.

| Capability                        | What it does                                                            | Citation                                             | Applicable here?                                                                                                                 |
| --------------------------------- | ----------------------------------------------------------------------- | ---------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| `Write(path, data)`               | Temp + fsync + rename; crash-safe single-file write                     | go-atomic-write/atomicwrite.go:17                    | **APPLICABLE** — replaces bare `os.WriteFile` in `generators.WriteFile`; a crash mid-write today can truncate `.goreleaser.yaml` |
| `WriteWithPerm(path, data, perm)` | Same, explicit perms (consumer mandates `0o644`)                        | go-atomic-write/atomicwrite.go:42                    | **APPLICABLE** — exact drop-in for the consumer's `WriteFile(path, data, perm)` signature                                        |
| `WriteVerified(path, data, fp)`   | + flock + fingerprint re-verify → `ErrConcurrentModification` on TOCTOU | go-atomic-write/atomicwrite.go:65; commit.go:13, :20 | **APPLICABLE** (borderline) — protects the backup-then-write sequence; two wizard processes in one repo are plausible in CI      |
| `WriteIfChanged(path, data)`      | No-op when bytes identical → no mtime bump, no watcher trigger          | go-atomic-write/atomicwrite.go:91                    | **APPLICABLE** — re-running `generate` (dogfood flow, AGENTS.md gotcha 19) currently always rewrites and bumps mtimes            |
| `Fingerprint` (xxhash64)          | `FingerprintFromBytes`/`FingerprintFile`/`Matches`                      | go-atomic-write/fingerprint.go:13, :26, :35          | **APPLICABLE** — input to `WriteVerified`                                                                                        |
| Durability details                | fsync temp before rename, `syncDir` after, parent `MkdirAll`            | go-atomic-write/rename_unix.go:11, :20; stage.go:16  | **APPLICABLE** — replaces ad-hoc `os.MkdirAll` in generators (github_actions.go:66, forgejo_workflows.go:55, homebrew.go:129)    |
| Streaming variants                | `WriteFunc*` callback writers for large content                         | go-atomic-write/atomicwrite.go:125, :151, :167       | N/A — generated artifacts are single small buffers                                                                               |

## What GoReleaser-Wizard does today

All generator writes funnel through one choke point: `generators.WriteFile` (cmd/goreleaser-wizard/generators/template_utils.go:23-30) = `os.WriteFile` at :24 + DomainError wrap; called by goreleaser.go:68, forgejo_workflows.go:99, dockerfile.go:79, homebrew.go:139. Two stragglers bypass it: validation_fixes.go:39 and simple_filesystem.go:29. Backups exist in exactly ONE generator: `GoReleaserGenerator.createBackup` (goreleaser.go:81-94) does `os.Rename(file, file+".backup")` at :85 — leaving a window where the target is absent between rename and the subsequent write, and only for `.goreleaser.yaml`. Multi-file atomicity is handled one level up: `GenerationPreflightJob` (jobs.go:439-503) refuses before any write; per-job `Rollback` deletes generated files on failure.

## Why CONSIDER

- **Real gap, right scope**: a crash mid-`os.WriteFile` (template_utils.go:24) can leave a truncated `.goreleaser.yaml`, and `createBackup`'s rename-away (goreleaser.go:85) leaves the target missing if the process dies before the write — temp+fsync+rename closes both windows per file.
- **It complements, not duplicates, the atomic-generation property**: DOMAIN_LANGUAGE.md's "all target files checked before any is written; a refusal leaves zero partial state" is enforced by preflight+rollback at the job level; this library adds crash-atomicity at the byte level _within_ each write, which the current machinery does not cover.
- **Does not replace the `.backup` user-facing feature** (goreleaser.go:83) — that stays; only the write mechanics change.
- **Cost is genuinely small**: 2 direct deps, one choke-point function + 2 stragglers; golden tests unaffected (content byte-identical).
- **Against**: single-shot CLI with tiny files makes the crash window microscopic; `WriteIfChanged` is nice-to-have; not urgent — hence CONSIDER, not ADOPT-NOW.

## Before/after

Before — `cmd/goreleaser-wizard/generators/template_utils.go:23-30`:

```go
func WriteFile(path string, data []byte, perm os.FileMode) error {
	err := os.WriteFile(path, data, perm)
	if err != nil {
		return WrapFileError(err, "failed to write file")
	}
	return nil
}
```

After:

```go
import atomicwrite "github.com/larsartmann/go-atomic-write"

func WriteFile(path string, data []byte, perm os.FileMode) error {
	if err := atomicwrite.WriteWithPerm(path, data, perm); err != nil {
		return WrapFileError(err, "failed to write file")
	}
	return nil
}
```

Optional hardening in `goreleaser.go Generate`: fingerprint the existing file before `createBackup` and use `atomicwrite.WriteVerified` so a concurrent writer yields `ErrConcurrentModification` instead of a silent clobber. Also migrate validation_fixes.go:39 to `generators.WriteFile`.

**Effort**: ~1-2 hours. **Tests**: none break (golden/E2E assert content); add a unit test that a failed write leaves no `*.tmp` residue and the original intact, a perms test (0o644 preserved), and — if `WriteVerified` is adopted — a TOCTOU test asserting `errors.Is(err, atomicwrite.ErrConcurrentModification)`.

## RECONSIDER-WHEN (escalation triggers to ADOPT-NOW)

1. `generate` starts writing files >1 MB or many files where partial-state risk grows.
2. Concurrent-wizard runs in CI become a supported scenario (fingerprint verification becomes load-bearing).
3. Rerun-idempotence ("no mtime bump when unchanged") becomes a documented feature.

## How to verify

```bash
cd /home/lars/projects/GoReleaser-Wizard
rg -n 'os\.WriteFile' --type go -g '!*_test.go'   # template_utils.go:24 choke point + stragglers validation_fixes.go:39, simple_filesystem.go:29
sed -n '81,94p' cmd/goreleaser-wizard/generators/goreleaser.go   # createBackup rename-away window
rg -c 'larsartmann/go-atomic-write' .             # 0 hits
git -C /home/lars/projects/go-atomic-write tag --list 'v*' --sort=-v:refname | head -1   # v0.6.2
cat /home/lars/projects/go-atomic-write/go.mod    # deps: cespare/xxhash v2.3.0, gofrs/flock v0.13.1, x/sys indirect
```
