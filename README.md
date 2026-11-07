# GoReleaser Wizard

**Stop copy-pasting GoReleaser configs. Get a production-ready setup with one command.**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-1.27+-blue.svg)](https://golang.org)
[![GoReleaser](https://img.shields.io/badge/powered%20by-GoReleaser-blue.svg)](https://goreleaser.com)

An interactive CLI tool that generates production-ready GoReleaser configurations for Go projects. It auto-detects your project structure, guides you through options with a rich TUI, and outputs `.goreleaser.yaml`, release CI workflows (GitHub Actions, Forgejo, or both via `--ci-provider`), and a Dockerfile -- all verified against GoReleaser v2.17 with zero deprecations.

## Quick Start

```bash
go install github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard@latest
goreleaser-wizard init
```

That's it. Your `.goreleaser.yaml` passes `goreleaser check` with zero deprecations, no environment variables required -- GitHub owner/repo are detected from your git remote (override with `--github-owner` / `--github-repo`).

## Demo

```
$ goreleaser-wizard init

  GoReleaser Wizard
  Let's configure your GoReleaser setup!

? Project Name › my-awesome-cli
? Binary Name › my-awesome-cli
? Main Package Path › ./cmd/my-awesome-cli
? Project Type › CLI Application
? Platforms › ✓ Linux  ✓ macOS  ✓ Windows  ✓ FreeBSD
? Architectures › ✓ amd64  ✓ arm64
? CGO Configuration › Disabled (recommended)
? Docker Support › Build and publish
? Git Provider › GitHub
? Include LDFlags? › Yes
? Enable Code Signing? › Yes
? Generate SBOM? › Yes
? Generate Homebrew Formula? › Yes
? Generate GitHub Actions? › Yes

Created  .goreleaser.yaml
Created  .github/workflows/release.yml

GoReleaser configuration initialized successfully!
```

## Installation

### Go Install

```bash
go install github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard@latest
```

### From Source

```bash
git clone https://github.com/LarsArtmann/GoReleaser-Wizard.git
cd GoReleaser-Wizard
go build -o goreleaser-wizard ./cmd/goreleaser-wizard
```

### Pre-built Binary

Download the latest release from the [releases page](https://github.com/LarsArtmann/GoReleaser-Wizard/releases).

### Homebrew

```bash
brew tap LarsArtmann/tap
brew install goreleaser-wizard
```

### Docker

```bash
docker pull ghcr.io/larsartmann/goreleaser-wizard:latest
```

## Commands

### `init` -- Interactive Wizard

Guides you through creating a complete configuration with a rich terminal UI:

```bash
goreleaser-wizard init
```

| Flag             | Default       | Description                                                                                      |
| ---------------- | ------------- | ------------------------------------------------------------------------------------------------ |
| `--force`        | `false`       | Overwrite existing configuration                                                                 |
| `--interactive`  | `true`        | Run the TUI wizard (set `false` for headless)                                                    |
| `--project-name` | auto-detected | Override project name                                                                            |
| `--binary-name`  | auto-detected | Override binary name                                                                             |
| `--main-path`    | auto-detected | Override main package path                                                                       |
| `--project-type` | auto-detected | Override project type                                                                            |
| `--github-owner` | git remote    | Override GitHub repository owner                                                                 |
| `--github-repo`  | git remote    | Override GitHub repository name                                                                  |
| `--ci-provider`  | `github`      | `github`, `forgejo`, `dual`, or `none` -- selects the generated workflow set and release targets |

When stdout is not a terminal, `init` prints a helpful message with the non-interactive flag.

### `generate` -- Non-Interactive Generation

For CI pipelines and scripting:

```bash
goreleaser-wizard generate \
  --project-name my-project \
  --binary-name my-app \
  --project-type cli \
  --config-only=false
```

| Flag             | Default       | Description                                                                                      |
| ---------------- | ------------- | ------------------------------------------------------------------------------------------------ |
| `--force`        | `false`       | Overwrite existing configuration                                                                 |
| `--config-only`  | `false`       | Generate only `.goreleaser.yaml` (no GitHub Actions)                                             |
| `--project-name` | auto-detected | Override project name                                                                            |
| `--binary-name`  | auto-detected | Override binary name                                                                             |
| `--main-path`    | auto-detected | Override main package path                                                                       |
| `--project-type` | auto-detected | Override project type                                                                            |
| `--github-owner` | git remote    | Override GitHub repository owner                                                                 |
| `--github-repo`  | git remote    | Override GitHub repository name                                                                  |
| `--ci-provider`  | `github`      | `github`, `forgejo`, `dual`, or `none` -- selects the generated workflow set and release targets |

Forgejo example:

```bash
goreleaser-wizard generate --ci-provider forgejo \
  --project-name my-project --binary-name my-app --project-type cli
```

This emits `.forgejo/workflows/ci.yml` + `release.yml` and a `.goreleaser.yaml` with `gitea_urls` (releases target your Forgejo instance; inside Forgejo Actions `secrets.GITHUB_TOKEN` is the Forgejo API token). Workflow pins (BuildFlow flake rev, goreleaser v2.17.1) are deliberate -- see `cmd/goreleaser-wizard/types/template_data.go`.

### `validate` -- Check Configuration

Validates your existing GoReleaser configuration:

```bash
goreleaser-wizard validate
goreleaser-wizard validate --verbose
goreleaser-wizard validate --fix
goreleaser-wizard validate --project-only
```

| Flag             | Default | Description                                |
| ---------------- | ------- | ------------------------------------------ |
| `--verbose`      | `false` | Show detailed validation output            |
| `--fix`          | `false` | Attempt to fix common issues automatically |
| `--project-only` | `false` | Validate project structure only            |

### `version` -- Version Info

```bash
goreleaser-wizard version
```

## What It Generates

### `.goreleaser.yaml`

Modern GoReleaser v2.17 configuration: multi-platform builds, archives with Windows zip overrides, checksums, GitHub-native changelog, literal owner/repo (checks clean locally without env vars), and `dockers_v2` multi-platform images with OCI annotations when Docker is enabled.

### `.github/workflows/release.yml` and/or `.forgejo/workflows/`

Tag-triggered release workflow on `ubuntu-latest` (GitHub) or your Forgejo Actions runner (default label `native`): checkout with full history, Go from `go.mod`, Docker buildx + registry login when Docker is enabled, cosign installer when signing is enabled, then `goreleaser release --clean`. `--ci-provider` picks the family: `github` (default, byte-identical to legacy output), `forgejo` (`.forgejo/workflows/ci.yml` + `release.yml`, goreleaser targets Forgejo via `gitea_urls`), `dual` (both sets), `none` (no workflows).

### `Dockerfile`

Prebuilt-pattern Dockerfile for `dockers_v2`: `FROM scratch` (or `alpine` when CGO is enabled), `ARG TARGETPLATFORM`, copies the binaries GoReleaser already built -- never rebuilds inside the image. Web API projects get `EXPOSE 8080`.

Generation is atomic: if any target file already exists the wizard refuses up front unless `--force` is set, so you never end up half-generated.

## Supported Project Types

The wizard adapts configuration to your project:

| Type                    | Description                                    |
| ----------------------- | ---------------------------------------------- |
| **CLI Application**     | Single binary with version info via ldflags    |
| **Web API**             | Includes Docker and health check configuration |
| **Library**             | Focuses on the CLI component if present        |
| **gRPC Service**        | Docker with multi-platform images              |
| **Microservice**        | Docker publishing, minimal cross-compilation   |
| **Desktop Application** | Platform-specific builds (Linux, macOS)        |
| **Daemon/Service**      | System service with nFPM packages              |
| **Command Line Tool**   | Lightweight single-binary distribution         |

## Platforms and Architectures

| Platform | Architectures |
| -------- | ------------- |
| Linux    | amd64, arm64  |
| macOS    | amd64, arm64  |
| Windows  | amd64, arm64  |
| FreeBSD  | amd64         |

## Docker Support

Four levels of Docker integration:

| Level             | Behavior                                              |
| ----------------- | ----------------------------------------------------- |
| None              | No Docker configuration                               |
| Build only        | Generates Dockerfile, builds images locally           |
| Publish only      | Pushes pre-built images to registry                   |
| Build and Publish | Full pipeline: build + push with multi-arch manifests |

Supported registries: GitHub Container Registry (ghcr.io), Docker Hub, GitLab Registry, Quay.io, and custom registries.

## Git Providers

| Provider    | Generated CI/CD                                     |
| ----------- | --------------------------------------------------- |
| GitHub      | GitHub Actions release workflow                     |
| Forgejo     | Forgejo Actions workflows (`--ci-provider forgejo`) |
| Gitea       | Manual (config only)                                |
| GitLab      | Config only (CI templates not generated)            |
| Bitbucket   | Config only (CI templates not generated)            |
| Self-hosted | Manual                                              |

## Architecture

GoReleaser Wizard follows Domain-Driven Design with Clean Architecture:

```
internal/domain/            Pure business logic, zero external dependencies
internal/validation/        Field-level and business rule validation
internal/git/               Git operations, URL parsing, version helpers
internal/config/            Configuration management

cmd/goreleaser-wizard/      CLI layer (Cobra), TUI (huh), workflow + jobs
cmd/goreleaser-wizard/generators/   Typed generators (single template render path)
cmd/goreleaser-wizard/templates/    Embedded templates -- single source of truth
cmd/goreleaser-wizard/types/        Typed template data + GitHub target resolution
```

The domain layer uses type-safe enums for all configuration options (platforms, architectures, project types, Docker support levels, signing levels, etc.) and enforces validation through a centralized error system with structured diagnostics.

Correctness is pinned by tests: an E2E test generates into a temp module and requires `goreleaser check` to exit clean, and golden files lock every generated artifact variant (`go test ./cmd/goreleaser-wizard -update-golden` to regenerate).

## Development

```bash
nix develop                                   # dev shell
go build ./...                               # go.mod requires Go 1.27 (encoding/json/v2 is the default)
go test ./...                                # includes E2E + golden tests
golangci-lint run ./...
go-arch-lint check
nix flake check
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines and [docs/goreleaser-guide.md](docs/goreleaser-guide.md) to learn GoReleaser v2.17 itself.

## License

[MIT](LICENSE) -- Lars Artmann
