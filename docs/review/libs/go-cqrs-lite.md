# go-cqrs-lite — Library Adoption Audit

| Field | Value |
| --- | --- |
| Library | `github.com/larsartmann/go-cqrs-lite` — a 90+-module composable CQRS + event-sourcing family. Fleet-sanctioned entry points: `system/v4` `v4.19.0`, `metaengine/v4` `v4.26.0`; also catalog/v4 `v4.15.0`, cmd/cqrs-lint/v4 `v4.22.2`, event/v4/eventtest `v0.4.0`, metaengine/sqliteengine/v4 `v4.9.1` (root tag `v4.0.0`; many low-level v4 modules are legacy/migration-only per fleet ADR-0123). License: **Proprietary** (README.md:8). |
| Local checkout | `/home/lars/projects/go-cqrs-lite` @ `feec9bc74` (2026-11-22) |
| Consumer | GoReleaser-Wizard @ `a493950` (2026-11-22) |
| Used? | **No** — `rg 'larsartmann/go-cqrs-lite'` → 0 hits |
| Verdict | NOT-ADOPTED → **STAY-OUT** |
| Score | n/a — 0 applicable capabilities (9 rows, all N/A) |
| Provenance | Fleet audit Group E, 2026-11-22; citations QC-sampled by orchestrator |

## What the library offers

"CQRS and Event Sourcing for Go — without the framework tax": the `system.New(ctx, DomainConfig, DeploymentConfig)` composition root auto-wires event store/bus/projections/queries/snapshots/timers/health/lifecycle; `metaengine` routes each declared query to the cheapest of N interchangeable engines (memory, sqlite, turso, pebble, bbolt, badger, duckdb, postgres…); catalog renders typed message declarations to AsyncAPI/EventCatalog/D2.

| Capability | What it does | Citation | Applicable here? |
| --- | --- | --- | --- |
| system.New composition root | Builds the whole CQRS runtime from two config types | go-cqrs-lite/system/v4/constructor.go:26 | N/A — requires commands/events and a durable deployment; the wizard has neither |
| DomainConfig / DeploymentConfig split | Developer domain vs operator engines | go-cqrs-lite/system/v4/config_types.go:17, :147 | N/A — nothing to deploy; no engines to choose |
| System lifecycle (`Start`) | Boots projection processing / long-lived runtime | go-cqrs-lite/system/v4/system_lifecycle.go:47 | N/A — the wizard exits after generation |
| System.HealthCheck | Engine health across wired engines | go-cqrs-lite/system/v4/introspection.go:164 | N/A — no engines, no long-running process |
| metaengine.Plan cost-based planner | Routes queries to the cheapest engine | go-cqrs-lite/metaengine/v4/planner.go:49 | N/A — zero read models; config reads are direct YAML parses |
| QueryDecl[Q,R] + Store.Apply folds | Typed read-model declarations | go-cqrs-lite/metaengine/v4/query.go:12; go-cqrs-lite/metaengine/v4/store.go:422 | N/A — no events to fold; output is files, not accumulated state |
| catalog message declarations | Typed Event/Command entries → AsyncAPI/D2 exports | go-cqrs-lite/catalog/v4/message_config.go:234 | N/A — no messages; the wizard's "catalog" is static embedded templates |
| cqrs-lint | Domain-aware linter for go-cqrs-lite consumers | go-cqrs-lite/cmd/cqrs-lint/v4/main.go:169 | N/A — consumer doesn't use the library; golangci-lint + go-arch-lint already cover this repo |
| eventtest helpers | Golden-file/event-handler test helpers | go-cqrs-lite/event/v4/eventtest/handlers.go:11 | N/A — no events; testify + golden files already cover generated artifacts |

## What GoReleaser-Wizard does today (the tempting-but-wrong angle, answered)

The consumer's "Job/Workflow engine" is **not CQRS and would not benefit from event sourcing**:

- `Job` interface = ID/Name/Execute/Rollback only (cmd/goreleaser-wizard/job_manager.go:24-29). Jobs are created per-run by `JobFactory.CreateFullWizardJobs` — up to 7 (job_factory.go:19-60) or 4 in config-only mode (:63-77).
- Sequential execution; first failure aborts (job_manager.go:127-143). `Workflow.Execute` wraps with a timeout and on error calls `RollbackFailedJobs` (workflow.go:76-91).
- **Rollback = compensating `os.Remove`** of the failed job's files (jobs.go:125-136, :198-230, :361-369). `RollbackFailedJobs` iterates only jobs whose status is *Failed*, in reverse (job_manager.go:280-305) — actual atomicity comes from `GenerationPreflightJob` refusing before any write (jobs.go:439-503).
- **No event log, no replay, no persistence, no bus.** Total state: an in-memory `config.State` enum (jobs.go:88-118) and an in-memory `[]JobResult` (job_manager.go:60-68). `updateJobStatus` is a debug log with the comment "could be extended to update a UI or event system" (job_manager.go:229-232) — never extended. The "Events" in DOMAIN_LANGUAGE.md are ubiquitous-language prose, not runtime events.
- Caveat chased to ground: a stale file-level `// TODO: Add workflow persistence and recovery` exists at workflow.go:11 in a file self-flagged "SPLIT IMMEDIATELY" (:3). It is NOT open backlog: TODO_LIST.md contains no persistence/event-sourcing items.
- Reads vs writes: `validate` reads a file the *user* wrote, not state the wizard produced — no reads-aside-from-writes to separate.

## Why STAY-OUT

- **The Job engine is not CQRS**: no domain events are published anywhere (the bus/projection machinery would have zero input); rollback is compensating file deletion — an event store cannot replay "un-write a file" better than `os.Remove`; the process is single-shot so there is no durable timeline to source events from; re-runs are cheap, preflight-guarded, and idempotent by design.
- `system.New` presumes a long-lived runtime (Start/lifecycle/health); booting and tearing down a full CQRS composition root inside a seconds-long CLI run is pure overhead.
- The wizard's actual hard problems — template rendering, atomic multi-file writes, golden regression nets — are solved by the generators/templates layer and would be untouched.
- **License friction**: go-cqrs-lite is Proprietary; importing it into this otherwise-permissive-dep CLI changes redistribution properties for zero capability gained.
- The fleet skill itself scopes the library to CQRS/ES workloads — neither condition holds here.

## RECONSIDER-WHEN

The product model inverts: a roadmap item introducing **durable state that must survive the process** — a daemon/server mode, persistent generation history with audit/replay, or multi-project state with derived read models. Concretely: when `rg 'database|sql\.|sqlite'` in this repo returns real storage code (today the only 2 hits are CGO build-tag *names* in internal/domain/build_tag.go:97, :101) or the single-shot CLI model is removed from ROADMAP.md.

## How to verify

```bash
cd /home/lars/projects/GoReleaser-Wizard
rg -c 'larsartmann/go-cqrs-lite' .                                       # 0 hits
sed -n '24,29p' cmd/goreleaser-wizard/job_manager.go                     # Job interface: Execute+Rollback, no events
sed -n '280,305p' cmd/goreleaser-wizard/job_manager.go                   # rollback = reverse over Failed jobs only
sed -n '11p' cmd/goreleaser-wizard/workflow.go                           # stale file-level TODO (not open backlog)
rg -n 'database|sql\.|bbolt|badger|pebble|sqlite' --glob '!*.md' --glob '!testdata/**'  # build_tag.go:97,101 only
git -C /home/lars/projects/go-cqrs-lite tag --list 'system/v4/v*' --sort=-v:refname | head -1   # system/v4/v4.19.0
```
