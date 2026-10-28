# GoReleaser-Wizard Agent Guide

**Project**: GoReleaser-Wizard - Interactive GoReleaser configuration wizard
**Repository**: https://github.com/LarsArtmann/GoReleaser-Wizard

---

## Overview

GoReleaser-Wizard is an interactive CLI tool that generates production-ready GoReleaser configurations. It follows Domain-Driven Design (DDD) principles with Clean Architecture, using a workflow-based job execution system for configuration generation and validation.

**Key Technologies**:

- Go 1.26.5 (requires `GOEXPERIMENT=jsonv2` to build — see Gotchas)
- Cobra CLI framework
- Viper for configuration
- Charm libraries (lipgloss, log) for terminal UI
- golangci-lint for code quality
- go-arch-lint for architecture enforcement
- GoReleaser v2.17 as the target of generated configs

---

## Essential Commands

All Go commands need the environment prefixes below while `/mnt/buildcache` is broken.

### Build / Test / Lint

```bash
export GOCACHE=/tmp/gw-cache GOMODCACHE=/tmp/gw-modcache GOLANGCI_LINT_CACHE=/tmp/gw-lint-cache GOEXPERIMENT=jsonv2

go build ./...
go test ./...                            # includes E2E + golden tests
go test ./cmd/goreleaser-wizard -update-golden   # regenerate golden files
golangci-lint run ./...
go-arch-lint check
```

### Verify the product (generated configs)

```bash
go build -o /tmp/gw-bin ./cmd/goreleaser-wizard
# in a fresh Go module with git init:
/tmp/gw-bin generate --config <cfg.yaml> --github-owner <owner> --github-repo <repo>
goreleaser check                        # must exit 0 with zero deprecations
```

The E2E test `TestGeneratedConfigPassesGoReleaserCheck` automates exactly this and skips gracefully when `goreleaser` is not in PATH.

### Project-Specific Commands

- `nix flake check`, `nix build`, `nix run .#test`, `nix run .#lint` (never Makefile/justfile)
- BuildFlow runs as pre-commit hook: ~60 tools, auto-fixes files — check `git status` after every commit

---

## Project Architecture

### Directory Structure

### Architectural Layers

The project implements Clean Architecture with DDD principles:

1. **Domain Layer** (`internal/domain/`)
   - Pure business logic with zero external dependencies
   - Core types: `SafeProjectConfig`, enumerations, validation
   - All business rules and domain constraints
   - Interfaces for repository and service abstractions

2. **Application Layer** (`cmd/goreleaser-wizard/`)
   - Workflow orchestration and job execution
   - CLI commands (Cobra-based)
   - Template generation
   - File system operations

3. **Infrastructure Layer** (`internal/git/`, `internal/utils/`)
   - External system integrations
   - Git command wrappers
   - Utility functions

4. **Configuration Layer** (`cmd/goreleaser-wizard/templates/`)
   - Embedded Go templates — THE single source of truth for all generated files
   - Rendered exclusively by the typed generators in `cmd/goreleaser-wizard/generators/`
   - `cmd/goreleaser-wizard/types/` holds the typed template data + GitHub owner/repo resolution (flag overrides > git remote detection > placeholders)
   - Go templates for file generation
   - Template data structures
   - Template escaping utilities

### Architecture Enforcement

The project uses **go-arch-lint** to enforce architectural boundaries:

- **Deep scanning enabled**: AST-based method call analysis (v1.14.0+)
- **Domain isolation**: Domain layer cannot import application/infrastructure
- **Test file exclusions**: `*_test.go` files excluded from architecture validation
- **Vendor support**: Can validate with or without vendor dependencies

**Key rules** (from `.go-arch-lint.yml`):

- Domain components can only depend on other domain components
- Application layer can depend on domain but not infrastructure
- Infrastructure implements domain interfaces
- Test files have no architectural constraints

---

## Code Patterns and Conventions

### Domain-Driven Design Patterns

**1. Safe Configuration Types**

```go
// Single source of truth for configuration
type SafeProjectConfig struct {
    ProjectName        string      `json:"project_name" yaml:"project_name"`
    ProjectType        ProjectType `json:"project_type" yaml:"project_type"`
    // ... with validation methods
}

func (spc *SafeProjectConfig) Validate() error { }
func (spc *SafeProjectConfig) ApplyDefaults() { }
```

**2. Typed Enumerations**

```go
type ProjectType string

const (
    ProjectTypeCLI           ProjectType = "cli"
    ProjectTypeWebService    ProjectType = "web-service"
    ProjectTypeLibrary       ProjectType = "library"
)

func (pt ProjectType) IsValid() bool { }
```

**3. Smart Defaults**

```go
func NewSafeProjectConfig() *SafeProjectConfig {
    return &SafeProjectConfig{
        ProjectType:    GetRecommendedProjectType(),
        Platforms:      GetRecommendedPlatforms(),
        CGOStatus:      CGOStatusDisabled,
        // ... intelligent defaults based on project analysis
    }
}
```

### Error Handling

**Centralized Error System** (`internal/errors/domain_errors.go`)

All errors use typed error codes for structured error handling:

**Error Usage Pattern**:

**Recovery from Panics**:

### Workflow and Job System

The project uses a sophisticated workflow-based execution system:

**Workflow Pattern**:

**Job Factory Pattern**:

**Job Implementation**:

### Validation Patterns

**Multi-Level Validation**:

### Template Generation Pattern

### Testing Patterns

**Table-Driven Tests**:

**Test Setup Pattern**:

### CLI Pattern (Cobra)

---

## Code Style and Conventions

### Naming Conventions

**Packages**: Lowercase, single word when possible

- `domain`, `validation`, `errors`, `generators`, `jobs`

**Types**: PascalCase

- `SafeProjectConfig`, `Workflow`, `JobFactory`, `DomainError`

**Interfaces**: PascalCase, usually end with "er"

- `Logger`, `Job`, `FileSystem`

**Functions**: PascalCase (exported), camelCase (unexported)

- `ValidateProjectName()`, `prepareTemplateData()`, `createBackup()`

**Constants**: PascalCase

- `ErrValidationFailed`, `ProjectTypeCLI`, `JobExecutionStatusRunning`

**Variables**: camelCase

- `projectName`, `templateData`, `jobManager`

**Acronyms**: Capitalized in PascalCase, lowercase in camelCase

- `HTTPResponse` (not HttpResponse)
- `httpClient` (not HTTPClient)
- `CGOStatus`, `SBOM`, `LDFlags` (well-known acronyms all caps)

### File Organization

**Keep files under 300 lines** - Split immediately when exceeded

Example file splitting pattern (from workflow.go):

```
workflow_core.go         - Core Workflow struct and basic methods
workflow_execution.go    - Execution logic and state management
workflow_types.go        - Workflow-related types and enums
workflow_templates.go    - Workflow templates and factories
workflow_validation.go    - Workflow validation logic
```

**Related code together** - Group by functionality, not by type

### Import Order

### Function Patterns

**Early returns preferred** over nested conditionals:

**Small, focused functions** - Prefer <30 lines

**Explicit over implicit** - Clear function signatures

---

## Testing Guidelines

### Test Structure

- **Unit tests**: `*_test.go` in same package
- **Integration tests**: `integration_test.go`
- **Performance tests**: `*_test.go` with `Benchmark*` functions

### Coverage Requirements

- **Minimum 80% coverage** enforced by `just coverage`
- Test both success and error paths
- Test edge cases and boundary conditions

### Test Utilities

**Use testify/assert**:

```go
import "github.com/stretchr/testify/assert"

assert.Equal(t, expected, actual)
assert.NoError(t, err)
assert.True(t, condition)
assert.Error(t, err)
assert.Contains(t, str, substr)
```

**Setup/Teardown**:

---

## Important Gotchas

### Architecture Constraints

**1. Domain Layer Purity**

- Domain layer MUST NOT import application or infrastructure
- Domain cannot use external dependencies (except through interfaces)
- Test domain layer in isolation - no external services

**2. Test File Exclusions**

- `*_test.go` files are excluded from architecture validation
- This allows testing patterns like package imports
- Don't rely on this for production code

**3. Circular Dependencies**

- go-arch-lint will catch circular dependencies
- Import cycles prevent compilation
- Use dependency injection to break cycles

### Error Handling

**4. Always Use Domain Errors**

- Never use `errors.New()` or `fmt.Errorf()` directly
- Always create errors with `errors.New*Error()`
- Provides structured error context and recovery suggestions

**5. Context Propagation**

- All async operations should accept `context.Context`
- Check `ctx.Err()` for cancellation
- Pass context through to downstream operations

### Template System

**6. Template Escaping**

- Go templates use `{{` and `}}` delimiters
- Must escape literal braces with `{{"{{"}}` and `{{"}}"}}`
- Use `internal/validation/template_escaping.go` utilities

**7. Template Functions**

- Custom template functions registered in `template.FuncMap`
- Example: `"incpatch": git.IncPatchVersion`
- Functions must be exported and error-safe

### File Operations

**8. Use `os.WriteFile` with Explicit Permissions**

```go
// Good
os.WriteFile("config.yaml", data, 0o644)

// Avoid
ioutil.WriteFile("config.yaml", data, 0o644)  // Deprecated
```

**9. Backup Before Overwrite**

- Always create backups before overwriting existing files
- Pattern used in generators: `createBackup()` method

### Build and Dependencies

**10. Nix First, Never Justfile**

- Prefer `nix` commands (`nix flake check`, `nix run .#test`); never create justfiles

### Environment Gotchas (verified 2026-08-17)

**11. `GOEXPERIMENT=jsonv2` is REQUIRED to build**

- `internal/types` imports `encoding/json/v2` and `encoding/json/jsontext`
- Build fails without it; BuildFlow's preflight claims it is redundant — that static check is WRONG (it misses `internal/types`); trust the build
- CI (`.github/workflows/ci.yml`) sets it on every Go step

**12. `/mnt/buildcache` may be broken**

- Symptom: `failed to initialize build cache at /mnt/buildcache/...`
- Fix: export `GOCACHE=/tmp/gw-cache GOMODCACHE=/tmp/gw-modcache GOLANGCI_LINT_CACHE=/tmp/gw-lint-cache`

**13. BuildFlow pre-commit hook can hang on telemetry**

- PostHog TLS timeouts (network outage) block the commit AFTER the pipeline finished green
- Recovery: kill the `buildflow`/`git commit` processes, verify the pipeline log shows `0 failed`, then commit with `--no-verify` and note why in the message

**14. E2E testing gotcha**

- `go run github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard` resolves the PUBLISHED module from the proxy, not local code — always `go build -o /tmp/gw-bin ./cmd/goreleaser-wizard` first

**15. Domain validation rejects windows+arm64**

- Test configs must not combine them (`architecture_compat.go`); template ignore lists are belt-and-braces, not a license

### Forgejo CI provider (added 2026-10-28)

**16. `--ci-provider {github,forgejo,dual,none}` selects workflow + release targets**

- Default github = byte-identical legacy output; unset CIProvider behaves as github everywhere (`isGitHubLike`), legacy configs need no migration
- Forgejo templates: `ForgejoCITemplate`/`ForgejoReleaseTemplate` in `templates/embedded.go`, `[[ ]]` delims so `${{ secrets.* }}` literals need no escaping; generator `generators/forgejo_workflows.go` writes both files
- Forgejo releases: GoReleaser config renders `gitea_urls` + `changelog.use: git`; inside Forgejo Actions `secrets.GITHUB_TOKEN` IS the Forgejo API token (`GITEA_TOKEN`)
- Release shadow: `forgejo-release-shadow` (SystemNix) copies assets to GitHub; skips when `RELEASE_SHADOW_GITHUB_TOKEN` is absent
- Pins live in `types/template_data.go`: BuildFlow flake rev + nixpkgs goreleaser v2.17.1 — bump deliberately, never float
- Goldens: `forgejo-workflows.yml`, `goreleaser-forgejo.yaml`, `goreleaser-dual.yaml` (dual = both workflow sets, goreleaser still targets Forgejo); regenerate with `go test ./cmd/goreleaser-wizard -update-golden` (NOTE: passing that flag makes TestInitCommand fail — the CLI harness inherits the flag; tests are green without it)

**17. actionlint + self-hosted runner labels**

- Generated workflows use `runs-on: native`; actionlint rejects it without a config: `actionlint -config-file <(printf 'self-hosted-runner:\n  labels:\n    - native\n') ...` (the E2E test does exactly this)

**18. CLI defaults are applied twice (quirk, not a bug to chase)**

- `runGenerate` applies `ApplyDefaults()` after flag overrides (job-list decisions need defaulted fields); `ConfigGenerationJob` re-applies harmlessly. The docker dependency check still consults pre-override DockerSupport on the CLI path — pass a `docker` shim on hosts without docker when smoke-testing

**19. Dogfood proof**

- dnsblockd's `.forgejo/workflows` + `.goreleaser.yaml` were generated by this binary (`--ci-provider forgejo --project-type library`) and validated with goreleaser check + actionlint before push

> This file was auto-trimmed. Full history in git.
