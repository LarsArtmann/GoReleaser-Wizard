# GoReleaser-Wizard Features

> **Last Updated:** 2026-11-07 (docs-health audit; feature claims verified against code this sweep)\
> **Version:** Current Development State

## 🎯 Executive Summary

GoReleaser-Wizard is a CLI tool that generates production-ready GoReleaser v2.17 configurations plus release CI workflows (GitHub Actions, Forgejo, or both) for Go projects. Its correctness is machine-proven: an E2E test requires `goreleaser check` to exit 0 with zero deprecations on generated output, and golden files lock every generated artifact variant.

**Current State:** Core generation paths are fully functional and continuously verified. Package-manager publishing (Homebrew casks) and config-side signing are partial; Pro integration and config migration are planned.

---

## 🚀 Core Features

### 1. Interactive Configuration Wizard (TUI)

**State:** FULLY_FUNCTIONAL

- **Description:** Rich terminal wizard (`huh/v2`, Charm theme) that guides users through GoReleaser setup
- **Implementation:** `cmd/goreleaser-wizard/tui_wizard.go`
- **Features:**
  - Project information confirmation and modification ✅
  - Interactive platform and architecture selection ✅
  - CGO, Docker support level, and git provider configuration ✅
  - Advanced options (LDFlags, signing, SBOM, GitHub Actions) ✅
  - Final configuration review and confirmation ✅
  - Non-TTY detection with a helpful headless hint ✅

### 2. Project Auto-Detection

**State:** FULLY_FUNCTIONAL

- **Description:** Automatically detects project structure and settings
- **Features:**
  - Module name from go.mod ✅
  - Main package location (root, `cmd/`) ✅
  - Project type inference ✅
  - Owner/repo from git remote (flag-overridable) ✅

### 3. GoReleaser Configuration Generation

**State:** FULLY_FUNCTIONAL

- **Description:** Generates a modern `.goreleaser.yaml` (v2.17)
- **Features:**
  - Multi-platform builds with typed template data ✅
  - `dockers_v2` images + OCI annotations when Docker is enabled ✅
  - Checksums, GitHub-native changelog, snapshot config ✅
  - Literal owner/repo (checks clean locally, no env vars) ✅
  - Passes `goreleaser check` with zero deprecations (E2E-proven) ✅

### 4. CI Workflow Generation (GitHub / Forgejo / Dual)

**State:** FULLY_FUNCTIONAL

- **Description:** `--ci-provider {github,forgejo,dual,none}` selects the generated workflow set and release targets
- **Features:**
  - GitHub Actions release workflow (`ubuntu-latest`, buildx, conditional docker-login/cosign) ✅
  - Forgejo workflows (`.forgejo/workflows/ci.yml` + `release.yml`) with pinned BuildFlow flake rev and goreleaser v2.17.1 ✅
  - Forgejo releases: `gitea_urls` + git changelog; `GITHUB_TOKEN` is the Forgejo API token inside Forgejo Actions ✅
  - `dual`: both workflow sets; goreleaser still targets Forgejo ✅
  - Unset provider behaves byte-identically to legacy github output ✅

### 5. Dockerfile Generation

**State:** FULLY_FUNCTIONAL

- **Description:** goreleaser-prebuilt Dockerfile for `dockers_v2`
- **Features:**
  - `FROM scratch` (or `alpine` with CGO), `ARG TARGETPLATFORM`, copies prebuilt binaries ✅
  - Web API projects get `EXPOSE 8080` ✅
  - Wired into both full-wizard and config-only paths ✅
  - Rollback + `--force`/atomic preflight guards ✅

---

## 🔧 Configuration Features

### 6. Multi-Platform Build Support

**State:** FULLY_FUNCTIONAL

- **Platforms:** Linux, macOS, Windows, FreeBSD
- **Architectures:** amd64, arm64 (windows+arm64 rejected by domain validation by design)

### 7. Docker Integration

**State:** FULLY_FUNCTIONAL

- **Implementation:** `dockers_v2` multi-platform images; prebuilt-pattern Dockerfile (never rebuilds in-image)
- **Features:**
  - Build-only / publish-only / build-and-publish modes ✅
  - ghcr.io, Docker Hub, GitLab Registry, Quay.io, custom registries ✅
  - Verified end-to-end 2026-08-17: local snapshot release built and published amd64+arm64 images; 2026-11-07: own `.goreleaser.yaml` still passes `goreleaser check` (v2.17.1, zero deprecations)

### 8. Code Signing Support

**State:** PARTIALLY_FUNCTIONAL

- Cosign installer step in generated workflows (conditional) ✅
- Own-config cosign keyless signing with `--bundle` (CI-only, needs OIDC) ✅
- **Limitation:** generated `.goreleaser.yaml` contains no `signs:` section yet (ROADMAP)

### 9. Package Manager Integration

**State:** PARTIALLY_FUNCTIONAL

- Typed Homebrew cask generator exists (`cmd/goreleaser-wizard/generators/homebrew.go`) but is unwired, blocked on personal tap defaults
- Own `.goreleaser.yaml` uses `homebrew_casks` (dogfood), quarantine hook + `Casks/` directory
- **Planned:** Snap, Scoop, AUR, winget, chocolatey (see ROADMAP)

---

## 🛠️ Validation and Testing

### 10. Configuration Validation

**State:** FULLY_FUNCTIONAL

- File existence, YAML structure, GoReleaser integration (when installed), dependency availability ✅
- Template validation for all generators; detailed error reporting with suggestions ✅
- **Planned (M22):** deprecated-key detection with before→after migration hints (TODO_LIST)

### 11. Project Structure Validation

**State:** FULLY_FUNCTIONAL

- Go module validation, main directory detection, alternative-structure warnings ✅

### 12. Automated Correctness Proof

**State:** FULLY_FUNCTIONAL

- E2E: generate into a temp module → `goreleaser check` exit 0, zero deprecations, no env vars (`cmd/goreleaser-wizard/e2e_test.go`; skips gracefully without the binary) ✅
- Golden files for goreleaser configs (github + forgejo + dual variants), workflows, Dockerfiles (`cmd/goreleaser-wizard/golden_test.go`, `-update-golden`) ✅
- CI: lint + race tests + fixture `goreleaser check` (`.github/workflows/ci.yml`) ✅
- actionlint E2E with self-hosted-runner config for the `native` label ✅

---

## 🏗️ Architecture and Infrastructure

### 13. Domain-Driven Architecture

**State:** FULLY_FUNCTIONAL

- Type-safe enums for every configuration axis; `SafeProjectConfig` with validation + defaults ✅
- Clean Architecture separation enforced by go-arch-lint (deep AST scanning) ✅
- Structured `DomainError` system (`internal/domain/errors.go`) ✅

### 14. Template System

**State:** FULLY_FUNCTIONAL

- Embedded templates in `cmd/goreleaser-wizard/templates/embedded.go` — single source of truth, rendered only by typed generators ✅
- Forgejo templates use `[[ ]]` delimiters; GitHub templates escape `{{ }}` correctly ✅
- Typed template data with owner/repo resolution (flag > git remote > placeholder warning) ✅

### 15. Workflow Engine

**State:** FULLY_FUNCTIONAL

- Job interface (Execute + Rollback), JobManager (sequential + parallel), JobFactory ✅
- All generation job types: project validation, config, workflows, Dockerfile, dependency check, preview ✅
- Atomic generation: preflight checks all target files before anything is written ✅

---

## 🚧 Planned Features

### 16. GoReleaser Pro Integration

**State:** PLANNED

- Custom publishers, advanced templating, nightly builds UI

### 17. Configuration Migration System

**State:** PLANNED

- Migrate existing `.goreleaser.yaml` to modern v2.17 (detect deprecated keys, propose diff, backup + rollback)

### 18. Advanced Customization

**State:** PLANNED

- Custom build hooks (opt-in `go test`/`go generate` before-hooks — needs a TypeSpec change), environment-specific configs, custom archive formats

---

## 🐛 Known Issues and Limitations

- **Large files:** several core files exceed 600 lines (`internal/types/validation.go` 898, `internal/validation/basic.go` 767, `cmd/goreleaser-wizard/jobs.go` 685) — splitting tracked in TODO_LIST
- **Package Managers:** Homebrew cask generation built but not wired into workflows (blocked on personal defaults)
- **Config-side signing:** no `signs:` section in generated configs (CI-side only)
- **Coverage:** domain layer ~10%, several packages at 0% — coverage lifts tracked in TODO_LIST (M19-M21)
- **Test coverage of tests:** dual-provider `goreleaser check` E2E missing (only github+forgejo goldens are check-proven) — TODO_LIST M16b
- **GoReleaser Dependency:** E2E check requires external GoReleaser installation (skips gracefully without it)

---

## 📊 Feature Status Summary

Computed from the rows above on 2026-11-07.

| Category      | Fully Functional | Partially Functional | Planned | Broken |
| ------------- | ---------------- | -------------------- | ------- | ------ |
| Core Features | 5                | 0                    | 0       | 0      |
| Configuration | 2                | 2                    | 0       | 0      |
| Validation    | 3                | 0                    | 0       | 0      |
| Architecture  | 3                | 0                    | 0       | 0      |
| Planned       | 0                | 0                    | 3       | 0      |
| **TOTAL**     | **13**           | **2**                | **3**   | **0**  |

---

Open work lives in [TODO_LIST.md](TODO_LIST.md); long-term direction in [ROADMAP.md](ROADMAP.md); change history in [CHANGELOG.md](CHANGELOG.md).
