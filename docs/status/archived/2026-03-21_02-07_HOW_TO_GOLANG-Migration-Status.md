# Status Report: HOW_TO_GOLANG.md Migration

**Date**: 2026-03-21_02:07
**Author**: Lars Artmann
**Audience**: Development team
**Version**: 1.0

---

## Executive Summary

Successfully migrated from `gopkg.in/yaml.v3` to `github.com/go-faster/yaml` and replaced `github.com/spf13/viper` with `github.com/knadh/koanf/v2` configuration management. **Build currently FAILING** - blocked by koanf env provider API mismatch. **All tests currently fail** due to compilation errors. **22 files exceed 350-line limit from HOW_TO_GOLANG.md** **Disk space critical** (91% used, 22GB free at 96% used).

**Blockers**:\*\*

- koanf env provider v1 API incompatible with koanf v2
- Build failures
- Low disk space (22GB free)
- All tests failing
- 22 files exceed 350-line limit (not addressed)
- Cobra replacement not evaluated (documentation mentions fang, readme mentions viper)
- Test files have placeholder implementations
- Documentation out of sync with code changes
  - `internal/domain/config_core.go` - yaml import changed
- `internal/config/config.go` - New koanf config manager (180 lines)
  - `cmd/goreleaser-wizard/main.go` - Replaced viper with config.Manager
  - `cmd/goreleaser-wizard/init_test.go` - Updated to use config.Manager
  - `cmd/goreleaser-wizard/validate_test.go` - Updated to use config.Manager
  - `go.mod` - Removed viper, added koanf dependencies

  - `.golangci.yml` - Linter configuration updated (243 lines)

  - `docs/status/2026-03-20_22-19_TUI-Implementation-Complete.md` - Updated with migration progress
  - `internal/types/validation.go` - 857 lines - Core validation types, duplicate code
  - `cmd/goreleaser-wizard/jobs.go` - 833 lines - Job orchestration, very complex
  - `internal/domain/validation.go` - 659 lines - Domain validation logic, duplicate with `internal/validation/`
  - `internal/validation/business_rules.go` - 626 lines - Business rule validation, duplicate with `internal/domain/`
  - `internal/validation/basic.go` - 617 lines - Basic field validation, duplicate with `internal/validation/`
  - `cmd/goreleaser-wizard/jobs/implementations.go` - 573 lines - Job implementations, complex workflow
  - `internal/domain/interfaces.go` - 490 lines - Domain interfaces ( bloated)
  - `cmd/goreleaser-wizard/workflow.go` - 467 lines - Workflow orchestration ( complex)
  - `internal/domain/config_core.go` - 427 lines - Core configuration types ( bloated)
  - `cmd/goreleaser-wizard/jobs/factory.go` - 408 lines - Job factory pattern ( complex)
  - `internal/domain/config_defaults.go` - 381 lines - Default values and recommendations ( bloated)
  - `internal/domain/errors.go` - 368 lines - Domain error types ( bloated)
  - `cmd/goreleaser-wizard/job_manager.go` - 367 lines - Job manager with workflow execution
  - `cmd/goreleaser-wizard/tui_wizard.go` - 367 lines - TUI wizard implementation
  - `cmd/goreleaser-wizard/jobs/types.go` - 364 lines - Job types and enums ( bloated)
  - `internal/domain/enums.go` - 317 lines - Domain enumerations ( bloated)
  - `cmd/goreleaser-wizard/interactive.go` - 331 lines - Unused interactive prompt methods
  - `cmd/goreleaser-wizard/generators/constants.go` - 23 lines - Unused constants

  - `internal/utils/recommendations.go` - 195 lines - Utility functions ( bloated)
  - `cmd/goreleaser-wizard/validate_main.go` - 250 lines - Main validation logic ( bloated)
  - `internal/domain/ids.go` - 225 lines - ID types ( bloated)
  - `internal/domain/architecture.go` - 275 lines - Architecture configuration ( bloated)
  - `internal/domain/ids_test.go` - 200 lines - ID type tests ( bloated)
  - `internal/domain/architecture_test.go` - 275 lines - Architecture tests ( bloated)
  - `internal/domain/enums_test.go` - 317 lines - Enum tests ( bloated)
  - `cmd/goreleaser-wizard/generate.go` - 300 lines - Generate command logic ( bloated)
  - `internal/validation/template_escaping.go` - 156 lines - Template escaping utilities ( bloated)
  - `internal/validation/form_validator.go` - 132 lines - Form validation utilities ( bloated)
  - `cmd/goreleaser-wizard/generators/goreleaser.go` - 115 lines - GoReleaser generator ( bloated)
  - `cmd/goreleaser-wizard/generators/github_actions.go` - 107 lines - GitHub Actions generator ( bloated)
  - `cmd/goreleaser-wizard/generators/dockerfile.go` - 102 lines - Dockerfile generator ( bloated)
  - `cmd/goreleaser-wizard/generators/homebrew.go` - 95 lines - Homebrew generator ( bloated)
  - `cmd/goreleaser-wizard/generators/validate.go` - 84 lines - Validation generator ( bloated)
  - `internal/git/commands.go` - 76 lines - Git commands ( bloated)
  - `cmd/goreleaser-wizard/types/template_data.go` - 69 lines - Template data structures ( bloated)
  - `cmd/goreleaser-wizard/validate_display.go` - 63 lines - Validation display ( bloated)
  - `internal/validation/validate.go` - 60 lines - Validation entry point ( bloated)
  - `cmd/goreleaser-wizard/generators/sbom.go` - 55 lines - SBOM generator ( bloated)
  - `internal/validation/validation_utils.go` - 54 lines - Validation utilities ( bloated)
  - `cmd/goreleaser-wizard/validate_workflow.go` - 52 lines - Validation workflow ( bloated)
  - `cmd/goreleaser-wizard/generators/homebrew_formula.go` - 47 lines - Homebrew formula template ( bloated)
  - **Test Coverage**: 0% (unknown, cannot measure)
  - **Dependencies**:
  - `gopkg.in/yaml.v3` → `github.com/go-faster/yaml` (DONE)
  - `github.com/spf13/viper` → REMOVED from go.mod (in progress)
  - `github.com/knadh/koanf/v2` → ADDED (configuration management)
  - `github.com/knadh/koanf/parsers/yaml` → ADDED (YAML parsing)
  - `github.com/knadh/koanf/providers/confmap` → ADDED (default values)
  - `github.com/knadh/koanf/providers/env` → ADDED (environment variables) - NOTE: v1 API incompatible with koanf v2
  - `github.com/knadh/koanf/providers/file` → ADDED (file provider)
  - `github.com/knadh/koanf/providers/posflag` → ADDED (flag binding)

  - `github.com/stretchr/testify` → RETAINED (testing framework)
  - `github.com/charmbracelet/huh` → RETAINED (interactive TUI)
  - `github.com/charmbracelet/lipgloss` → RETAINED (terminal styling)
  - `github.com/charmbracelet/log` → RETAINED (logging)
  - `github.com/spf13/cobra` → RETAINED (CLI framework)
  - `github.com/spf13/pflag` → RETAINED (flag parsing)
  - `github.com/larsartmann/go-composable-business-types` → RETAINED (local replace)

- **Blockers**:
  - **koanf env provider v1 API mismatch** (CRITICAL)
    - `internal/config/config.go:13` uses v1 API: `env.Provider(prefix, delimiter, callback)`
    - koanf v2 expects: `env.Provider(config *env.Config{...})` with TransformFunc
    - Fix: Update to koanf v2 env provider OR use v1 API
    - **Build failures** (CRITICAL)
    - All tests currently fail due to undefined `viper` references
    - Low disk space (91% used, 22GB free) 96% used) causing slow downloads and cache corruption
    - Need to fix env provider API before tests can run
  - **Test failures** (CRITICAL)
    - `cmd/goreleaser-wizard/init_test.go` - viper.Reset() undefined
    - `cmd/goreleaser-wizard/validate_test.go` - viper.Reset() undefined
    - Need to fix env provider API and update tests
  - **Files > 350 lines** (NOT STARTED)
  - 22 files exceed 350-line limit from HOW_TO_GOLANG.md
    - Largest: 857 lines (`internal/types/validation.go`)
    - Smallest over limit: 360 lines (`cmd/goreleaser-wizard/jobs/types.go`)
  - Splitting required to maintain single responsibility
    - Consider incremental approach to avoid breaking existing functionality
  - **Cobra replacement** (Not Evaluated)
  - Documentation mentions `charmbracelet/fang` as alternative
  - README.md still mentions viper as dependency
  - HOW_TO_GOLANG.md doesn't mention fang specifically
  - Research needed to determine if fang is mature enough for production use
  - **Documentation Updates** (Partial)
  - AGENTS.md still references viper in examples (needs update)
  - README.md still mentions viper (needs update)
  - **Test Coverage** (Unknown)
  - All tests currently failing due to compilation errors
  - Cannot run test suite to verify changes
  - **Technical Debt** (Significant)
  - Large test files (200+ lines) need maintenance
  - Duplicate validation logic across multiple files
  - Deep module coupling between domain and application layers
  - Some functions with 5+ parameters
  - **Disk Space Critical** (LOW)
  - 91% disk usage (22GB free)
  - Risk of Running out of disk space during builds
  - Workaround: Clean caches before builds (temporary solution)
  - **CI Pipeline Status** (Unknown)
  - Cannot run `just ci` until build succeeds
  - **Pre-commit Hooks** (Unknown)
  - May fail due to build errors
  - **File Size Analysis**:
    | File                                                   | Lines | Category                 |
    | ------------------------------------------------------ | ----- | ------------------------ |
| ~~`internal/types/validation.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~857~~ | ~~Validation Types~~ |
| ~~`cmd/goreleaser-wizard/jobs.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~833~~ | ~~Job Orchestration~~ |
| ~~`internal/domain/validation.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~659~~ | ~~Domain Validation~~ |
| ~~`internal/validation/business_rules.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~626~~ | ~~Business Rules~~ |
| ~~`internal/validation/basic.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~617~~ | ~~Basic Validation~~ |
| ~~`cmd/goreleaser-wizard/jobs/implementations.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~573~~ | ~~Job Implementations~~ |
| ~~`internal/domain/interfaces.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~490~~ | ~~Domain Interfaces~~ |
| ~~`cmd/goreleaser-wizard/workflow.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~467~~ | ~~Workflow~~ |
| ~~`internal/domain/config_core.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~427~~ | ~~Core Configuration~~ |
| ~~`cmd/goreleaser-wizard/jobs/factory.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~408~~ | ~~Job Factory~~ |
| ~~`internal/domain/config_defaults.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~381~~ | ~~Configuration Defaults~~ |
| ~~`internal/domain/errors.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~368~~ | ~~Domain Errors~~ |
| ~~`cmd/goreleaser-wizard/tui_wizard.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~367~~ | ~~TUI Wizard~~ |
| ~~`cmd/goreleaser-wizard/jobs/types.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~364~~ | ~~Job Types~~ |
| ~~`cmd/goreleaser-wizard/job_manager.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~367~~ | ~~Job Manager~~ |
| ~~`internal/domain/enums.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~317~~ | ~~Domain Enums~~ |
| ~~`cmd/goreleaser-wizard/interactive.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~331~~ | ~~Interactive Prompts~~ |
| ~~`cmd/goreleaser-wizard/generate.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~300~~ | ~~Generate Command~~ |
| ~~`internal/domain/architecture.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~275~~ | ~~Architecture Config~~ |
| ~~`cmd/goreleaser-wizard/validate_main.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~250~~ | ~~Validation Main~~ |
| ~~`internal/domain/ids.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~225~~ | ~~ID Types~~ |
| ~~`internal/utils/recommendations.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~195~~ | ~~Recommendations~~ |
| ~~`internal/domain/ids_test.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~200~~ | ~~ID Tests~~ |
| ~~`internal/validation/template_escaping.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~156~~ | ~~Template Escaping~~ |
| ~~`internal/validation/form_validator.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~132~~ | ~~Form Validator~~ |
| ~~`cmd/goreleaser-wizard/generators/goreleaser.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~115~~ | ~~GoReleaser Generator~~ |
| ~~`cmd/goreleaser-wizard/generators/github_actions.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~107~~ | ~~GitHub Actions Generator~~ |
| ~~`cmd/goreleaser-wizard/generators/dockerfile.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~102~~ | ~~Dockerfile Generator~~ |
| ~~`cmd/goreleaser-wizard/generators/homebrew.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~95~~ | ~~Homebrew Generator~~ |
| ~~`internal/domain/architecture_test.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~275~~ | ~~Architecture Tests~~ |
| ~~`internal/domain/enums_test.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~317~~ | ~~Enum Tests~~ |
| ~~`cmd/goreleaser-wizard/generators/validate.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~84~~ | ~~Validate Generator~~ |
| ~~`internal/git/commands.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~76~~ | ~~Git Commands~~ |
| ~~`cmd/goreleaser-wizard/types/template_data.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~69~~ | ~~Template Data~~ |
| ~~`cmd/goreleaser-wizard/validate_display.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~63~~ | ~~Validation Display~~ |
| ~~`internal/validation/validate.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~60~~ | ~~Validation Entry~~ |
| ~~`cmd/goreleaser-wizard/generators/sbom.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~55~~ | ~~SBOM Generator~~ |
| ~~`internal/validation/validation_utils.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~54~~ | ~~Validation Utils~~ |
| ~~`cmd/goreleaser-wizard/validate_workflow.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~52~~ | ~~Validation Workflow~~ |
| ~~`cmd/goreleaser-wizard/generators/homebrew_formula.go`~~ **Superseded — point-in-time audit; splitting superseded — TODO_LIST.md (split-oversized-files).** | ~~47~~ | ~~Homebrew Formula~~ |

  ## Priority Recommendations (Next Steps)

  | Priority | Task | Status | Effort | Impact |
  | -------- | ----------------------------------------------------------------- | ------- | --------- | --------------------------------- | --------------------- |
| ~~**1**~~ done at `50c0feb`, ` a7680ac` | ~~Fix koanf env provider API~~ | ~~BLOCKED~~ | ~~15 min~~ | ~~HIGH - Blocks all builds/tests~~ |
| ~~**2**~~ done at `50c0feb`, ` a7680ac` | ~~Verify build succeeds~~ | ~~PENDING~~ | ~~5 min~~ | ~~HIGH - Confirms migration success~~ |
| ~~**3**~~ done at `50c0feb`, ` a7680ac` | ~~Run test suite~~ | ~~PENDING~~ | ~~10 min~~ | ~~HIGH - Ensures no regressions~~ |
| ~~**4**~~ **Superseded — TODO_LIST.md (split-oversized-files).** | ~~Split `internal/types/validation.go` (857 lines)~~ | ~~PENDING~~ | ~~2 hours~~ | ~~HIGH - Largest violation~~ |
| ~~**5**~~ **Superseded — TODO_LIST.md (split-oversized-files).** | ~~Split `cmd/goreleaser-wizard/jobs.go` (833 lines)~~ | ~~PENDING~~ | ~~2 hours~~ | ~~HIGH - Complex orchestration~~ |
| ~~**6**~~ **Superseded — TODO_LIST.md (split-oversized-files).** | ~~Split `internal/domain/validation.go` (659 lines)~~ | ~~PENDING~~ | ~~1.5 hours~~ | ~~HIGH - Core domain logic~~ |
| ~~**7**~~ **Superseded — TODO_LIST.md (split-oversized-files).** | ~~Split `internal/validation/business_rules.go` (626 lines)~~ | ~~PENDING~~ | ~~1.5 hours~~ | ~~MEDIUM - Business validation~~ |
| ~~**8**~~ **Superseded — TODO_LIST.md (split-oversized-files).** | ~~Split `internal/validation/basic.go` (617 lines)~~ | ~~PENDING~~ | ~~1.5 hours~~ | ~~MEDIUM - Basic validation~~ |
| ~~**9**~~ **Won't implement — cobra kept; fang never adopted.** | ~~Evaluate Cobra → fang migration~~ | ~~PENDING~~ | ~~2 hours~~ | ~~MEDIUM - Modern CLI framework~~ |
| ~~**10**~~ done at `50c0feb`, ` ed01089` | ~~Update AGENTS.md viper examples~~ | ~~PENDING~~ | ~~30 min~~ | ~~LOW - Documentation consistency~~ |
| ~~**11**~~ done at `50c0feb`, ` ed01089` | ~~Update README.md dependencies~~ | ~~PENDING~~ | ~~15 min~~ | ~~LOW - User-facing docs~~ |
| ~~**12**~~ **Superseded — TODO_LIST.md (M20).** | ~~Add integration tests for config manager~~ | ~~PENDING~~ | ~~1 hour~~ | ~~MEDIUM - Test coverage~~ |
| ~~**13**~~ done at `a5e790d` | ~~Consolidate duplicate validation logic~~ | ~~PENDING~~ | ~~3 hours~~ | ~~MEDIUM - Code quality~~ |
| ~~**14**~~ done at `7e85f9e` | ~~Remove unused code in interactive.go~~ | ~~PENDING~~ | ~~30 min~~ | ~~LOW - Code cleanup~~ |
| ~~**15**~~ done at `9703270` | ~~Address golangci-lint warnings~~ | ~~PENDING~~ | ~~1 hour~~ | ~~MEDIUM - Code quality~~ |
| ~~**16**~~ done at `fb51e45` | ~~Fix local go-composable-business-types dependency~~ | ~~PENDING~~ | ~~30 min~~ | ~~MEDIUM - Dependency management~~ |
| ~~**17**~~ done — green in gates | ~~Review architecture violations~~ | ~~PENDING~~ | ~~2 hours~~ | ~~MEDIUM - Architecture compliance~~ |
| ~~**18**~~ done — BuildFlow active | ~~Add pre-commit hook verification~~ | ~~PENDING~~ | ~~30 min~~ | ~~LOW - CI reliability~~ |
| ~~**19**~~ **Superseded — environment-era.** | ~~Free up disk space (currently 91%)~~ | ~~PENDING~~ | ~~1 hour~~ | ~~HIGH - Blocks builds~~ |
| ~~**20**~~ **Superseded — TODO_LIST.md (split-oversized-files).** | ~~Split `cmd/goreleaser-wizard/jobs/implementations.go` (573 lines)~~ | ~~PENDING~~ | ~~1.5 hours~~ | ~~MEDIUM~~ | ~~Job implementations~~ |
| ~~**21**~~ **Superseded — TODO_LIST.md (split-oversized-files).** | ~~Split `internal/domain/interfaces.go` (490 lines)~~ | ~~PENDING~~ | ~~1 hour~~ | ~~MEDIUM~~ | ~~Interface definitions~~ |
| ~~**22**~~ **Superseded — TODO_LIST.md (split-oversized-files).** | ~~Split `cmd/goreleaser-wizard/workflow.go` (467 lines)~~ | ~~PENDING~~ | ~~1 hour~~ | ~~MEDIUM~~ | ~~Workflow logic~~ |
| ~~**23**~~ **Superseded — TODO_LIST.md (split-oversized-files).** | ~~Split `internal/domain/config_core.go` (427 lines)~~ | ~~PENDING~~ | ~~1 hour~~ | ~~MEDIUM~~ | ~~Core configuration~~ |
| ~~**24**~~ **Superseded — TODO_LIST.md (split-oversized-files).** | ~~Split `cmd/goreleaser-wizard/jobs/factory.go` (408 lines)~~ | ~~PENDING~~ | ~~1 hour~~ | ~~MEDIUM~~ | ~~Factory pattern~~ |
| ~~**25**~~ **NOT-DO — superseded by completed migration.** | ~~Create migration checklist document~~ | ~~PENDING~~ | ~~30 min~~ | ~~LOW~~ | ~~Process tracking~~ |

  ## Metrics Summary

  | Metric                    | Value    | Target  | Status         |
  | ------------------------- | -------- | ------- | -------------- |
| ~~Banned libraries replaced~~ **Superseded — superseded same day (build/test green); disk environment-era.** | ~~2/3~~ | ~~3~~ | ~~67% Complete~~ |
| ~~Files under 350 lines~~ **Superseded — superseded same day (build/test green); disk environment-era.** | ~~~200/222~~ | ~~222~~ | ~~~90% Compliant~~ |
| ~~Files over 350 lines~~ **Superseded — superseded same day (build/test green); disk environment-era.** | ~~22~~ | ~~0~~ | ~~NEEDS WORK~~ |
| ~~Build status~~ **Superseded — superseded same day (build/test green); disk environment-era.** | ~~FAILING~~ | ~~PASSING~~ | ~~BLOCKED~~ |
| ~~Test status~~ **Superseded — superseded same day (build/test green); disk environment-era.** | ~~FAILING~~ | ~~PASSING~~ | ~~BLOCKED~~ |
| ~~Disk usage~~ **Superseded — superseded same day (build/test green); disk environment-era.** | ~~91%~~ | ~~<85%~~ | ~~CRITICAL~~ |

  ## Open Questions
~~  1. Should we use koanf v1 or v2 env provider? (v1 uses callback, v2 uses Config struct)~~ done at `50c0feb (env provider v2 in go.mod)`
~~  2. What is the priority order for file splitting? (By size, by complexity, by domain importance?)~~ **Superseded — TODO_LIST.md (split-oversized-files).**
~~  3. Should we address architecture violations before or after file splitting?~~ done — architecture green in gates
~~  4. Is the local `go-composable-business-types` dependency intentional or a mistake?~~ done at `fb51e45`
~~  5. Should we add integration tests before or after fixing the current blocker?~~ done at `8edf0b1`

~~  6. What is the minimum Go version we should support? (Currently requires 1.26.1)~~ done — go.mod on 1.27 (2026-11-07)
~~  7. Should we create a separate migration branch or continue on master?~~ done — master
~~  8. How should we handle the 22 oversized files? (All at once, incrementally, or by priority?)~~ **Superseded — TODO_LIST.md (split-oversized-files).**
~~  9. Should we document the splitting strategy before starting?~~ **NOT-DO — moot: migration completed.**
~~  10. What is the rollback plan if file splitting breaks functionality?~~ **NOT-DO — moot: migration completed.**

  ## Top #1 Question (Cannot Answer Myself)

  **Should we use koanf v1 env provider (with callback API) or upgrade to koanf v2 env provider (with Config struct)?**
  - v1 API: `env.Provider(prefix, delimiter, callback)`
  - v2 API: `env.Provider(&env.Config{Prefix: "...", TransformFunc: ...})`
  - Current go.mod has v1 (`github.com/knadh/koanf/providers/env v1.1.0`)
  - Current config.go imports v1 path (`github.com/knadh/koanf/providers/env`)
  - But the callback signature is v1 style, while the implementation attempted v2 style
  - **Decision needed**: Use v1 API correctly OR upgrade to v2 providers

  ## Detailed File Changes

  ### internal/config/config.go (180 lines - NEW)

  ````go
  // Package config provides configuration management using koanf.
          package config

          import (
              "fmt"
              "os"
              "path/filepath"
              "strings"
              "sync"

              "github.com/knadh/koanf/parsers/yaml"
              "github.com/knadh/koanf/providers/confmap"
              "github.com/knadh/koanf/providers/env"  // v1 import
              "github.com/knadh/koanf/providers/file"
              "github.com/knadh/koanf/providers/posflag"
              "github.com/knadh/koanf/v2"
              "github.com/spf13/pflag"
          )

          // Config holds application configuration.
          type Config struct {
              Debug   bool `koanf:"debug"`
              NoColor bool `koanf:"no-color"`
          }

          // Manager manages configuration loading and access.
          type Manager struct {
              k     *koanf.Koanf
              mu    sync.RWMutex
              cfg   Config
              flags *pflag.FlagSet
          }

          // NewManager, Load, RegisterFlags, loadFile, Get, IsDebug,
          // NoColors, Set, Reset, GetRaw, GetManager, SetGlobalManager
          // ... (full implementation)
          ```

  ### cmd/goreleaser-wizard/main.go (Changes)
  ```diff
          - import "github.com/spf13/viper"
          + import "github.com/LarsArtmann/GoReleaser-Wizard/internal/config"

          - viper.GetBool("debug")
          + config.GetManager().IsDebug()

          - viper.BindPFlag(...)
          + config.GetManager().RegisterFlags(rootCmd.PersistentFlags())

          - func initConfig() {
          -     viper.SetConfigFile(cfgFile)
          -     viper.AddConfigPath(home)
          -     viper.SetConfigName(".goreleaser-wizard")
          -     viper.SetConfigType("yaml")
          -     viper.SetEnvPrefix("GORELEASER_WIZARD")
          -     viper.AutomaticEnv()
          -     viper.ReadInConfig()
          +     config.GetManager().Load(cfgFile)
          +     // Simplified error handling
          }
  ````

  ### Test File Changes

            ```diff
            // init_test.go
            - viper.Reset()
            - viper.Set("debug", false)
            + config.GetManager().Reset()
            + config.GetManager().Set("debug", false)

            + viper.Get(key)
            + config.GetManager().GetRaw(key)

            // validate_test.go
            - viper.Reset()
            - viper.Set("debug", false)
            + config.GetManager().Reset()
            + config.GetManager().Set("debug", false)
            ```

  ## Commit Plan

~~      1. **Fix env provider API** (Priority 1)~~ done at `50c0feb`, ` a7680ac`
         - Update internal/config/config.go to use correct v1 API
         - Run `go build ./...` to verify
~~         2. **Verify and commit** (Priority 2)~~ done at `50c0feb`, ` a7680ac`
         - Run `go test ./...` to ensure tests pass
         - Commit with message: "refactor(config): fix koanf env provider API"
~~         3. **Split oversized files** (Priority 3)~~ **Superseded — TODO_LIST.md (split-oversized-files).**
         - Start with largest file: internal/types/validation.go (857 lines)
         - Create splitting plan document
         - Execute splits incrementally
         - Run tests after each split
~~         4. **Update documentation** (Priority 4)~~ done at `50c0feb`
         - Remove viper references from AGENTS.md
         - Update README.md dependencies section
         - Add migration notes to changelog
~~      5. **Final verification** (Priority 5)~~ done — gates green (2026-11-07)
         - Run full CI pipeline locally
         - Verify all changes work correctly
         - Create PR or merge to master
