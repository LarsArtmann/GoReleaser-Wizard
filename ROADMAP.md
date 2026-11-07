<!-- ROADMAP.md template, copy into the project root and fill in.
     This is LONG-TERM DIRECTION and raw ideas.
     No bounded tasks, no status indicators, no completed milestones.
     See docs-health/SKILL.md. -->

# Roadmap

> Long-term direction and raw ideas. Items here are NOT actionable tasks.
> When an idea is refined into bounded work, it moves to [TODO_LIST.md](TODO_LIST.md).
> Rebuilt 2026-11-07 from the docs-health audit of `docs/status/` (2026-01 … 2026-11) and `docs/research/2026-08-10_GORELASER_17_RESEARCH.md`.

## Themes

### 1. One command from any Go repo to a releasable project

The north star stays: `goreleaser-wizard generate` should be the first thing a new LarsArtmann repo runs, and everything it emits should be releasable without hand-editing.

Raw ideas:

- Wire the full `homebrew_casks` path (generator exists in `cmd/goreleaser-wizard/generators/homebrew.go`, unwired pending personal tap defaults)
- Config-side signing: emit a `signs:` section (today only the CI-side cosign installer is generated)
- More GoReleaser v2.17 surface: `winget`, `chocolatey`, `nsis`, `aur`, `krew`, NUR, `verifiable_builds`, root-level `retry` tuning
- Parameterized everything: tap repo, NUR repo, commit author, runner label as first-class config (feeds off TODO_LIST M18)
- Self-hosted runner revival (`--runner` flag) alongside the `ubuntu-latest` default
- Snap/Scoop support returned as opt-in sections (the old unwired generators were removed in the v2.17 modernization; re-add on demand, not speculatively)

### 2. Validation that understands GoReleaser, not just YAML

`validate` should be the tool a user runs BEFORE a release, not after a failure.

Raw ideas:

- Deprecated-key detection with exact before→after migration snippets (bounded work, see TODO_LIST M22)
- Split "config generation correctness" from "release-readiness" validation: a blocking pre-release environment check (`gpg`/`cosign`/registry auth present) as a separate command — the signing-tool question from the 2026-06-30 report, still awaiting a product decision
- Ingest a user's existing `.goreleaser.yaml` and propose a modernization diff (migration system: version compat, backup, rollback)

### 3. CI-provider families beyond GitHub + Forgejo

`--ci-provider {github,forgejo,dual,none}` proved the template-family pattern.

Raw ideas:

- GitLab CI and Bitbucket Pipelines workflow families (the domain already has the git-provider enums)
- Gitea self-hosted flavor distinct from Forgejo (shared `gitea_urls`, different Actions assumptions)
- Release-shadow as a generated optional workflow (today it is a SystemNix fleet script)

### 4. Extensibility

Raw ideas:

- Plugin system for custom generators/templates (custom template repositories; the generators package is already the seam)
- Multi-project/monorepo support: one wizard run, per-module configs
- AI-assisted config suggestions inside the TUI (suggest sensible defaults from repo inspection)
- TUI polish: accessible mode, a config preview screen before generation, progress indicators, i18n

### 5. Mechanical guarantees

Raw ideas:

- Enforce domain purity mechanically (go-arch-lint catches imports, not `exec.LookPath` behavior; a hermetic test harness or custom linter rule — the 2026-06-30 lesson)
- Property-based/fuzz tests for validation rules
- BDD (Ginkgo) coverage for the interactive wizard flow (skill exists, unused)

## Non-goals

Things we are deliberately NOT pursuing and why:

- **Multi-VCS rewrite of the domain model**: git-provider enums cover the real fleet (GitHub/Forgejo); abstracting further is speculative.
- **Pin `anchore/sbom-action/download-sbom`-style floating subactions to SHA where upstream 404s**: known-unpinnable upstream (BuildFlow warning accepted, 2026-08-17 report).
- **Web UI**: the CLI/TUI is the product; a web layer duplicates it without a consumer.
- **Chaos engineering tests**: no production infra to speak of; effort wildly outweighs impact for a config generator.
- **Re-adding `brews` (formula) generation**: superseded by `homebrew_casks` direction since the v2.17 modernization; old formulas in the personal tap get cleaned, not regenerated.

---
