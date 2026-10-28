# Forgejo Migration — Privileged Batch (user-run)

**Created**: 2026-10-28
**Owner**: lars (every command here needs root / FORGEJO_TOKEN / deploy rights that agent sessions deliberately lack)
**Master plan**: `docs/planning/2026-10-12_15-17_forgejo-primary-migration-master-plan.md`
**Session report**: `docs/status/2026-10-28_00-28_forgejo-migration-execution-session-recon-and-t6.md`
**Status report for this batch**: see "Execution state" at the bottom — everything NOT in this file is already done and verified.

Run the phases in order. Each phase lists its verification; stop on the first red.

---

## Phase A — Deploy SystemNix (carries: dnsblockd canonicalRepos entry, `forgejo-release-shadow` + `forgejo-branch-protect` binaries, fixture fixes)

```bash
cd ~/projects/SystemNix
nix flake check                       # fixture check forgejo-scripts-fixture must pass (verified green 2026-10-28)
sudo nixos-rebuild switch --flake .#evo-x2
```

This deploy also re-runs `forgejo-ssh-keys.service`, which re-pushes the account-side SSH key
(`lars@evo-x2` — the key push had staled server-side; `ssh git@forgejo.home.lan` failing with
`Permission denied (publickey)` was the symptom). Verify afterwards:

```bash
ssh -T git@forgejo.home.lan                # expect a forgejo greeting, not publickey denial
which forgejo-release-shadow               # /run/current-system/sw/bin/forgejo-release-shadow
```

## Phase B — Flip the pilot (dnsblockd)

```bash
# dry-run first: read the plan it prints
sudo systemctl start forgejo-flip-check@dnsblockd.service
journalctl -u forgejo-flip-check@dnsblockd --no-pager -n 50

# execute (deletes the disposable mirror, full-imports issues/PRs/releases/wiki/LFS,
# verifies, attaches the 8h+sync_on_commit GitHub push mirror)
sudo systemctl start forgejo-flip@dnsblockd.service
journalctl -u forgejo-flip@dnsblockd --no-pager -n 80
```

Verify (unauthenticated + authenticated):

```bash
curl -s https://forgejo.home.lan/api/v1/repos/lars/dnsblockd | jq '{mirror, default_branch, has_actions}'
# expect: mirror=false, default_branch=master
# if has_actions=false, enable it:
#   source /var/lib/forgejo/.admin-token.env && curl -s -X PATCH \
#     -H "Authorization: token $FORGEJO_TOKEN" -H 'Content-Type: application/json' \
#     -d '{"has_actions":true}' https://forgejo.home.lan/api/v1/repos/lars/dnsblockd
```

Push-mirror attach is idempotent via phase 5 of `forgejo-github-sync` (dnsblockd is in
`canonicalRepos`); either wait for the next hourly sync run or start it now:

```bash
sudo systemctl start forgejo-github-sync.service
journalctl -u forgejo-github-sync --no-pager -n 30 | grep dnsblockd   # expect "push mirror attached"
```

## Phase C — Release-shadow secret (one command per canonical repo; dnsblockd today)

The release workflow reads `secrets.RELEASE_SHADOW_GITHUB_TOKEN` (a repo-level Forgejo Actions
secret holding a GitHub PAT with `repo` scope) and `secrets.GITHUB_TOKEN` (Forgejo-provided
already). As root:

```bash
GH_PAT='<github PAT with repo scope>'   # or: pass show github/release-shadow-pat
source /var/lib/forgejo/.admin-token.env
curl -s -X PUT \
  -H "Authorization: token $FORGEJO_TOKEN" -H 'Content-Type: application/json' \
  -d "{\"data\":\"$GH_PAT\"}" \
  https://forgejo.home.lan/api/v1/repos/lars/dnsblockd/actions/secrets/RELEASE_SHADOW_GITHUB_TOKEN
# expect HTTP 201/204. List to confirm:
curl -s -H "Authorization: token $FORGEJO_TOKEN" \
  https://forgejo.home.lan/api/v1/repos/lars/dnsblockd/actions/secrets | jq .
```

Without this secret the release still succeeds end-to-end on Forgejo; only the GitHub asset
copy is skipped (the workflow step exits 0 by design).

## Phase D — Live verification of CI + release + shadow

```bash
cd ~/projects/dnsblockd
git push origin master                       # mirror-sync pushes flow both ways after flip
# CI run appears at https://forgejo.home.lan/lars/dnsblockd/actions (run-on: native runner)

# release path (optional, cuts a real release):
#   bump version, tag, push the TAG TO FORGEJO directly now that it is primary:
git remote add forgejo git@forgejo.home.lan:lars/dnsblockd.git 2>/dev/null || true
#   git tag v0.9.5 && git push forgejo v0.9.5
# then watch the release workflow: BuildFlow gate -> goreleaser (GITEA_TOKEN) -> shadow step
```

Mirror-Actions observability answer (was UNVERIFIED in recon): public mirrors on the instance
all show `has_actions:false` (checked unauthenticated 2026-10-28, incl. `lars/GoReleaser-Wizard`),
so the 2h pull-sync of mirror repos does NOT trigger Actions runs — the flip (native +
`has_actions:true`) is what arms workflows. After Phase B, the first push to dnsblockd is the
positive control.

## Phase E — Optional token handoff for future agent sessions

```bash
pass insert forgejo/lars-api-token    # paste the admin/user Forgejo API token
```

This unblocks future non-interactive verification (runner status, Actions runs, secrets listing)
without root. Until then, agent sessions stay read-only on Forgejo (recon finding: all tokens
root/hermes-gated by design).

---

## Execution state (2026-10-28, what is ALREADY DONE)

| Piece | State |
| --- | --- |
| Wizard `--ci-provider {github,forgejo,dual,none}` (T6–T10) | DONE, committed `9334dc3` (+daemon commits); build/tests/lint/arch-lint green; goldens pinned incl. forgejo+dual; goreleaser 2.17.1 check zero-deprecation; actionlint E2E green |
| dnsblockd `.forgejo/workflows/{ci,release}.yml` + `.goreleaser.yaml` | DONE, pushed to GitHub master (`411dce43`) |
| SystemNix `forgejo-release-shadow` + fixture tests | DONE, committed (`4ddc9a84`/`a1ca8384`), fixture check green |
| SystemNix `canonicalRepos += dnsblockd` + runbook | DONE, committed |
| Phases A–E above | WAITING ON USER |

## Later waves (unchanged from master plan §6)

After the pilot verifies: waves of remaining ~260 repos via `forgejo-flip-batch` (design pending,
T2), runner capacity check (T16: single `native` runner,
capacity 2, 16G watermark), GoReleaser-Wizard release (v2 with `--ci-provider`), then per-repo
wizard regeneration sweeps (`goreleaser-wizard generate --ci-provider forgejo`).
