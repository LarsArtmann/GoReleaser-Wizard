# Security Policy

Security policy for GoReleaser-Wizard, an interactive CLI that generates GoReleaser
configurations and CI release workflows. Because its output ships into other
projects' release pipelines, correctness and supply-chain hygiene of this repo
directly affect every generated config's consumers.

## Reporting a Vulnerability

- Use GitHub's private security advisory for this repository
  (Security tab -> Report a vulnerability). Do not open a public issue for
  anything that looks exploitable.
- Include the affected version/commit, a minimal reproduction, and the impact.
- You can expect an initial response within 7 days.

## Supported Versions

Only the latest released version and the current `master` branch receive
security fixes.

| Version            | Supported |
| ------------------ | --------- |
| latest release tag | yes       |
| master (HEAD)      | yes       |
| older tags         | no        |

## Security Practices

### Source and build integrity

- All GitHub Actions are pinned to full commit SHAs (`actions/checkout@<sha> # vX.Y.Z`),
  never to moving tags.
- Workflows run with least-privilege `permissions` (CI jobs: `contents: read`).
- The Nix flake builds the binary hermetically: pinned nixpkgs revision, pinned
  Go toolchain (`go_1_27`), pinned `vendorHash` over a vendored dependency set,
  no network access during the build.

### Dependency hygiene

- `govulncheck` is available in the dev shell; run it before releases.
- Dependencies are pinned via `go.mod`/`go.sum` and verified by the Nix
  `vendorHash`; hash drift fails the build instead of silently upgrading.
- The BuildFlow pipeline can run secret scanning (`gitleaks`) and spell/secret
  hygiene checks on demand (`buildflow -s gitleaks`).

### Secrets handling

- No secrets are stored in the repository.
- CI/release tokens (`GITHUB_TOKEN`, Forgejo `GITEA_TOKEN`, release-shadow
  `RELEASE_SHADOW_GITHUB_TOKEN`) are provided exclusively through the CI
  secret store and consumed via environment variables in generated workflows.
- Generated GoReleaser configs reference secrets only through CI-native secret
  interpolation inside the target CI system (GitHub Actions `${{ secrets.* }}`,
  Forgejo Actions `${{ secrets.GITHUB_TOKEN }}`).

### Generated output guarantees

- Generated workflows pin tool versions (BuildFlow flake rev, goreleaser v2.17.1)
  instead of floating tags, so consumers get reproducible release tooling.
- Generated configs are validated with `goreleaser check` and `actionlint` in
  this repo's E2E tests before they can ship.
