#!/usr/bin/env bash
# docs-gates.sh — standing docs-gate composition for GoReleaser-Wizard.
#
# Order (docs-health canon): formatter FIRST (a cooperating rewriter must not
# run after verdict gates without a re-gate), then the docs-gates battery,
# then the repo's own Go docs gate.
#
# Adoption posture (2026-11-07): the battery currently reports 5 red gates —
# the 2025 pre-canon status corpus (29 unarchived reports, oldest 508 days)
# predates the annotate canon and needs an owner-scale annotate+archive
# sweep before cadence/routing go green. This script is therefore ADVISORY
# until that sweep lands: it prints the matrix, never fails the build.
# Remediation is routed in TODO_LIST.md / the root-sprawl disposition
# proposal (docs/planning/).
set -uo pipefail

cd "$(git rev-parse --show-toplevel)" || exit 1

echo "== 1/3 dprint format check (markdown + config files) =="
if command -v dprint >/dev/null 2>&1; then
  dprint check || echo "[advisory] formatting drift; run: dprint fmt"
else
  echo "[advisory] dprint not installed; skipping format check"
fi

echo
echo "== 2/3 docs-gates battery (freshness, audit-all, annotate, routing, cadence) =="
if command -v docs-gates >/dev/null 2>&1; then
  docs-gates run .
elif command -v nix >/dev/null 2>&1; then
  nix run github:Artmann-Technologies/docs-gates -- run . || echo "[advisory] battery red; see matrix above"
else
  echo "[advisory] neither docs-gates nor nix available; skipping battery"
fi

echo
echo "== 3/3 repo Go docs gate (links + FEATURES counts) =="
go run ./cmd/docs-check . || echo "[blocking] docs-check failed"

echo
echo "docs-gates.sh done (battery is advisory until the 2025-corpus sweep; see file header)"
