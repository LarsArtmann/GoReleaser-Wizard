# GoReleaser-Wizard: Research & Gap Analysis (2026-08-10)

Research performed against the current codebase (master, `358e420`) and GoReleaser
**v2.17.1** (official docs + source at `goreleaser/goreleaser`).

---

## 1. Project Health Check

### Verified Commands

| Check                                                                               | Result                                      | Notes                                                                            |
| ----------------------------------------------------------------------------------- | ------------------------------------------- | -------------------------------------------------------------------------------- |
| ~~`go build ./...`~~ done — point-in-time checks; still green 2026-11-07            | ~~✅ passes~~                               | ~~—~~                                                                            |
| ~~`go test ./...`~~ done — point-in-time checks; still green 2026-11-07             | ~~✅ 7 packages pass~~                      | ~~only 14 `_test.go` files / 83 `.go` files~~                                    |
| ~~`nix flake check`~~ done — point-in-time checks; still green 2026-11-07           | ~~✅ all 3 checks pass~~                    | ~~build / test / format~~                                                        |
| ~~`go-arch-lint check`~~ done — point-in-time checks; still green 2026-11-07        | ~~✅ OK~~                                   | ~~stale `internal/errors` refs from May 2026 are gone~~                          |
| ~~`golangci-lint run ./...`~~ done at `9703270`                                     | ~~⚠️ **52 issues**~~                         | ~~39 godoclint + 13 makezero~~                                                   |
| ~~Coverage~~ **Superseded — TODO_LIST.md (M19-M21).**                               | ~~⚠️~~                                       | ~~cmd 51.5%, domain 10%, validation 57%, **5 packages at 0%**~~                  |
| ~~`goreleaser check` (own `.goreleaser.yaml`)~~ done at `969259e`                   | ~~⚠️ valid, **uses deprecated properties**~~ | ~~`dockers`, `docker_manifests`, `brews`~~                                       |
| ~~`GOEXPERIMENT=jsonv2`~~ **Superseded — obsolete: jsonv2 default on go.mod 1.27.** | ~~✅ works~~                                | ~~`encoding/json/v2` is intentional; env must be set; go.mod targets go 1.26.5~~ |

### Repo Hygiene

- `.orig` files are still tracked despite `fa648db` ("chore(repo): remove leftover .orig backup files"):
  - `cmd/goreleaser-wizard/main.go.orig`
  - `cmd/goreleaser-wizard/validate_main.go.orig`
  - `internal/domain/errors.go.orig`
- `gopls` flags: unused `createTempFile` in `validate_test.go:59`, 5 `infertypeargs` hints in `internal/domain/ids.go`.

---

## 2. End-to-End Verification (clean temp project)

Ran `generate` in a scratch Go module:

### Found Failures

~~1. **`goreleaser check` fails without env vars.**~~ done at `e85bb4f`
`release:` section uses `{{.Env.GITHUB_OWNER}}` / `{{.Env.GITHUB_REPO}}`.
The wizard fills them from git detection at _generate_ time and the generated
GH Actions workflow sets them, but a local `goreleaser check` / `release`
without those env vars fails:
`template: failed to apply "{{.Env.GITHUB_REPO}}": map has no entry for key "GITHUB_REPO"`.

~~2. **Docker config generated, but no `Dockerfile`.** The wizard emits a~~ done at `e85bb4f`, `6f4f54d`
`dockers:` section whenever Docker support is enabled, but the
`DockerfileGenerator` is never wired into any workflow. A snapshot release
fails at the docker step:
`docker build failed: failed to copy dockerfile: ... lstat Dockerfile: no such file or directory`.

~~3. **Deprecated config emitted.** The generated `.goreleaser.yaml` contains~~ done at `e85bb4f`
`archives.format_overrides.format` (deprecated, use `formats`), `dockers` +
`docker_manifests` (deprecated, use `dockers_v2`), and `brews` (deprecated,
use `homebrew_casks`).

~~4. **Release hooks run the test suite.** `before.hooks` always contains~~ done at `e85bb4f`
`go mod tidy`, `go generate ./...`, `go test ./...`. Running tests as part of
every release build is usually unwanted and slows CI.

~~5. **GH Actions workflow is host-specific.** Generated `release.yml` uses~~ done at `3282680`
`runs-on: ["self-hosted", "linux", "x64"]`. It also hardcodes
`docker/login-action` + `DOCKER_USERNAME`/`DOCKER_PASSWORD` secrets and
`cosign-installer` regardless of whether those features are actually enabled.

### Working Pieces

- Non-interactive `generate` runs the full job workflow (validate → deps → config → actions).
- Project auto-detection (go.mod module name, `cmd/*/main.go` or root `main.go`) works.
- The generated config builds and archives binaries successfully (`goreleaser release --snapshot` up to the Docker failure).

---

## 3. Architecture Findings

### Dead Code: `generators/` and `templates/` packages are NOT wired in

- `go list -deps ./cmd/goreleaser-wizard` shows only `domain`, `git`, `types`, `config` — **not** `generators` or `templates`.
- No production file imports `cmd/goreleaser-wizard/generators` or uses `templates.GoReleaserTemplate`.
- The **live path** is:
  `init.go`/`generate.go` → `workflow.go` (`ExecuteWorkflow`) → `JobFactory` (`CreateFullWizardJobs`) → `jobs.go`
- `jobs.go` (839 lines) contains:
  - inline `goreleaserTemplateContent` / `githubActionsTemplateContent` (a **second, divergent copy** of the template in `templates/embedded.go`)
  - `prepareGoReleaserData()` / `prepareGitHubActionsData()` returning `map[string]any` — despite `cmd/goreleaser-wizard/types/template_data.go` already having **typed** `GoReleaserTemplateData` / `GitHubActionsTemplateData`
  - `// CRITICAL TODO: Replace map[string]any with strongly typed structs` acknowledged in code and again in status reports — unresolved.
- `templates/embedded.go` vs `jobs.go` template copy divergence (checked via diff): the two GoReleaser templates drifted; GitHub Actions / Dockerfile / Homebrew templates only exist in the dead package.

### Error-Discarding / Fragile Sites (non-test)

| File                                                                                    | Line        | Pattern                                                                                 |
| --------------------------------------------------------------------------------------- | ----------- | --------------------------------------------------------------------------------------- |
| ~~`cmd/goreleaser-wizard/flags.go`~~ **Superseded — TODO_LIST.md (M17).**               | ~~8, 16~~   | ~~`value, _ := cmd.Flags().GetBool/GetString(name)` — errors silently swallowed~~       |
| ~~`cmd/goreleaser-wizard/init.go`~~ **Superseded — TODO_LIST.md (M17).**                | ~~248~~     | ~~`matches, _ := filepath.Glob(...)` — filesystem errors hidden~~                       |
| ~~`cmd/goreleaser-wizard/validate_main.go`~~ **Superseded — TODO_LIST.md (M17).**       | ~~167~~     | ~~`if exists, _ := fileSystemRepo.DirExists(...)` — FS error treated as "dir missing"~~ |
| ~~`cmd/goreleaser-wizard/types/template_data.go`~~ **Superseded — TODO_LIST.md (M17).** | ~~236-261~~ | ~~`parseGitHubRemote()` swallows git errors, returns placeholder `"owner"`/`"repo"`~~   |
| ~~`cmd/goreleaser-wizard/jobs.go`~~ done — defensible by design                         | ~~490-491~~ | ~~`_ = os.Remove(...)` — best-effort cleanup (defensible)~~                             |
| ~~`internal/config/config.go`~~ done — defensible by design                             | ~~110~~     | ~~`_ = info` — vestigial assignment~~                                                   |
| ~~`cmd/goreleaser-wizard/validate.go`~~ done — defensible by design                     | ~~41~~      | ~~`os.Exit(exitCode)` — unconditional process exit~~                                    |
| ~~`cmd/goreleaser-wizard/main.go`~~ done — defensible by design                         | ~~162~~     | ~~`os.Exit(1)` inside `recoverFromPanic`~~                                              |

### Other Structural Notes

- `ProjectConfig = domain.SafeProjectConfig` alias in `main.go:17`; `config` (koanf) `Manager` is wired only for debug/no-color flags.
- Domain layer is rich (type-safe enums for platforms, project types, CGO, Docker support, git providers, signing, actions, features) and validates well — the architecture foundation is good.
- `integration_test.go` and `architecture_test.go` exist but are not coverage-heavy for the real generation paths.
- Test-wizard fixture exists at `test-wizard/` (used for deduplication testing, not part of the suite).

---

## 4. GoReleaser v2.17.1 — Current Feature Surface (Official Docs)

### Key Deprecations (from goreleaser.com/deprecations)

| Deprecated                                                                                   | Since                 | Replacement                                                   |
| -------------------------------------------------------------------------------------------- | --------------------- | ------------------------------------------------------------- |
| ~~`dockers` + `docker_manifests`~~ done — upstream doc facts; dated 2026-08-10 in header     | ~~v2.12 / v2.16~~     | ~~**`dockers_v2`** (fast-tracked to become `dockers` in v3)~~ |
| ~~`brews`~~ done — upstream doc facts; dated 2026-08-10 in header                            | ~~v2.10 soft, v2.16~~ | ~~**`homebrew_casks`**~~                                      |
| ~~`archives.format_overrides.format`~~ done — upstream doc facts; dated 2026-08-10 in header | ~~v2.6~~              | ~~`archives.format_overrides.formats` (list)~~                |
| ~~per-`dockers` `retry`~~ done — upstream doc facts; dated 2026-08-10 in header              | ~~v2.15.3~~           | ~~root-level `retry`~~                                        |

### `dockers_v2` Essentials (verified against docs source)

```yaml
dockers_v2:
  - id: myimg # default: project name
    dockerfile: Dockerfile # path from project root; templates allowed
    ids: [mybuild] # filter binaries/packages to copy
    images: # image names (templates allowed)
      - "user/repo"
      - "ghcr.io/user/repo"
    tags: # tag templates
      - "v{{ .Version }}"
      - "{{ if not .IsNightly }}latest{{ end }}"
    platforms: # default: [linux/amd64 linux/arm64]
      - linux/amd64
      - linux/arm64
    extra_files: [] # files to add to the build context
    labels: {} # image labels
    annotations: {} # OCI annotations ({{.Date}}, {{.FullCommit}}, etc.)
    sbom: "{{ not .IsNightly }}"
    build_args: {}
    flags: []
    retry: { attempts: 10, delay: 10s, max_delay: 5m }
```

Behavior notes:

- **Build and push are one step**: `docker buildx build --push`; runs in the **publish phase** (NOT in `goreleaser build`).
- `--snapshot` builds per-platform images instead (`repo:1.2.4-amd64`, `repo:1.2.4-arm64`) so the Dockerfile can be verified locally.
- The build context is a _temporary directory_ containing prebuilt binaries/packages — **do not build binaries inside the Dockerfile**. Recommended minimal Dockerfile:

  ```dockerfile
  FROM scratch
  ARG TARGETPLATFORM
  ENTRYPOINT ["/usr/bin/myprogram"]
  COPY $TARGETPLATFORM/myprogram /usr/bin/
  ```

- Requires a multi-platform buildx builder:
  ```sh
  docker buildx create --name=goreleaser --use
  docker run --privileged --rm tonistiigi/binfmt --install all
  ```

### Modern sections the wizard could generate

- `homebrew_casks` (tap casks, `directory: Casks`)
- `winget`, `chocolatey`, `msi`, `nsis` (Windows), `aur` (Arch), `krew` (kubectl), `nix`/NUR, `scoops`
- `builds.verifiable_builds` (reproducible builds)
- `nightlies`, `snapshot`, `release.mode`, `source` archives, `sboms`, `signs` (cosign keyless), `attestations`
- `retry` (root level), `docker_digests`
- `before` hooks, `publishers` (custom), `verify` (post-release verification, v2.17)

---

## 5. Recommended Plan (Priority Order)

### A. Fix broken output first — correctness

~~1. **Generate a `Dockerfile` whenever Docker is enabled** (wire `DockerfileGenerator` into the job factory; use a `scratch` + `COPY $TARGETPLATFORM/...` layout).~~ done at `6f4f54d`
~~2. **Make `release` env-tolerant**: keep `{{.Env.GITHUB_OWNER}}`/`{{.Env.GITHUB_REPO}}` but detect defaults from git, or surface a clear error at generate/check time (fail fast with instructions).~~ done at `e85bb4f`
~~3. **Remove `go test ./...` / `go generate ./...` from `before.hooks`** by default (offer as opt-in).~~ done at `e85bb4f`, `3282680`
~~4. **Fix GH Actions template**: use `ubuntu-latest` by default (or keep self-hosted as an option), only emit Docker login / cosign steps when the respective features are enabled.~~ done at `e85bb4f`, `3282680`

### B. Consolidate dead code — maintainability

~~5. **Delete or wire `cmd/goreleaser-wizard/generators/` and `templates/`**. Recommended: wire the typed generators in (they already exist: GoReleaser, GitHubActions, Dockerfile, Homebrew) and route `jobs.go` through them; delete `jobs.go`'s inline templates and `map[string]any` data prep.~~ done at `9f4bd71`
~~6. **Remove the 3 `.orig` files** (git rm) — they were intended to be deleted in `fa648db`.~~ done at `3cb91da`
~~7. Fix error-discarding sites in `flags.go`, `init.go:248`, `validate_main.go:167`, `template_data.go` (`parseGitHubRemote`).~~ **Superseded — TODO_LIST.md (M17).**

### C. Modernize to GoReleaser 2.17

~~8. **Migrate templates** to `dockers_v2` (buildx), `homebrew_casks`, `archives.formats`.~~ **Superseded — dockers_v2+formats done (e85bb4f); casks unwired → TODO_LIST.md (M18).**
~~9. Add wizard **knowledge guide**: a `docs/` + README section explaining modern GoReleaser concepts (dockers_v2, casks, phases, snapshot) — the "understand goreleaser" deliverable.~~ done at `969259e`, `8edf0b1`
~~10. Add a `--verify` step that runs `goreleaser check` on the generated output and reports deprecations/warnings.~~ done at `969259e`, `8edf0b1`

### D. Quality gates

~~11. Fix the 52 lint issues (39 godoclint comment start-symbol fixes, 13 makezero slice preallocation) — mostly mechanical.~~ done at `9703270`
~~12. Lift coverage: golden-file tests for generated configs, tests for `generators`, `templates`, `types`, `internal/config`, `internal/git`.~~ **Superseded — TODO_LIST.md (M19-M21); goldens done 8edf0b1.**
~~13. CI: run `goreleaser check` on the wizard's own config and on generated fixtures so deprecations cannot silently slip in again.~~ done at `8edf0b1`

---

## 6. Open Questions for the Owner

~~1. **Docker**: is per-arch `dockers` behavior (build in build-phase, local images) relied upon, or is `dockers_v2` (publish-phase, manifest-only) acceptable? This changes the template shape and CI workflow.~~ done at `e85bb4f`, `3282680`
~~2. **CI runners**: keep `self-hosted` in generated workflows, or default to `ubuntu-latest`?~~ done at `e85bb4f`, `3282680`
~~3. **Homebrew**: migrate from `brews` (formulas) to `homebrew_casks`, or keep formulas for legacy users?~~ **Superseded — TODO_LIST.md (M18, BLOCKED).**
~~4. **Dead packages**: full delete, or effort to wire the typed `generators` package in as the single source of truth?~~ done at `9f4bd71`
~~5. **`encoding/json/v2`**: remains intentional (requires `GOEXPERIMENT=jsonv2`)? Should this be documented in AGENTS.md + flake devShell?~~ **Superseded — obsolete: jsonv2 default on Go 1.27.**
