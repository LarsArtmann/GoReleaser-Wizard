// Command forgejo-preflight generates the flip-preflight checklist for a
// repository scheduled into a Forgejo migration wave (M12 wave-1 skeleton).
//
// It inspects a local checkout and prints PASS/FAIL lines for every gate the
// fleet flip requires, so the owner-run batch (forgejo-flip-batch@) only
// reaches repos whose preflight is green. Skeleton stage: checks run locally,
// no Forgejo API calls yet.
//
// Usage:
//
//	forgejo-preflight /path/to/repo
//
// Exit 0 = all gates green, 1 = at least one gate failed, 2 = usage error.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type check struct {
	name string
	ok   bool
	hint string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: forgejo-preflight <repo-dir>")

		os.Exit(2)
	}

	repo := filepath.Clean(os.Args[1])

	// The user-supplied repo path IS the tool's contract (a local admin
	// diagnostic); gosec's taint tracking cannot know that. Cleaned above.
	info, err := os.Stat(repo) //nolint:gosec // G703: repo path is the explicit CLI input, filepath.Clean applied
	if err != nil || !info.IsDir() {
		fmt.Fprintf(os.Stderr, "not a directory: %s\n", repo)

		os.Exit(2)
	}

	checks := runPreflight(repo)

	failed := 0

	for _, chk := range checks {
		status := "PASS"

		if !chk.ok {
			status = "FAIL"

			failed++
		}

		fmt.Fprintf(os.Stdout, "%s  %s", status, chk.name)

		if chk.hint != "" {
			fmt.Fprintf(os.Stdout, "  (%s)", chk.hint)
		}

		fmt.Fprintln(os.Stdout)
	}

	if failed > 0 {
		fmt.Fprintf(os.Stderr, "\n%d gate(s) failed; repo is not flip-ready\n", failed)

		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "\n%s is flip-ready\n", filepath.Base(repo))
}

// runPreflight evaluates the local gates a repo must pass before its flip.
func runPreflight(repo string) []check {
	checks := []check{
		{
			name: ".goreleaser.yaml exists",
			ok:   fileExists(filepath.Join(repo, ".goreleaser.yaml")),
			hint: "run goreleaser-wizard generate",
		},
		{
			name: ".forgejo/workflows/ci.yml exists",
			ok:   fileExists(filepath.Join(repo, ".forgejo", "workflows", "ci.yml")),
			hint: "run goreleaser-wizard generate --ci-provider forgejo",
		},
		{
			name: ".forgejo/workflows/release.yml exists",
			ok:   fileExists(filepath.Join(repo, ".forgejo", "workflows", "release.yml")),
			hint: "run goreleaser-wizard generate --ci-provider forgejo",
		},
		{
			name: ".forgejo/actionlint.yaml exists",
			ok:   fileExists(filepath.Join(repo, ".forgejo", "actionlint.yaml")),
			hint: "regenerate with a wizard version >= 2026-11-07",
		},
		{
			name: "no stale GitHub workflows scheduled on tags",
			ok:   !fileExists(filepath.Join(repo, ".github", "workflows", "release.yml")),
			hint: "the legacy release.yml must be deleted pre-flip or both systems fire on a tag",
		},
	}

	return checks
}

func fileExists(path string) bool {
	// Paths are built from the (cleaned) repo root + fixed constants, not
	// user strings; gosec cannot track that across the call boundary.
	_, err := os.Stat(path) //nolint:gosec // G703: paths derive from the cleaned repo root

	return err == nil
}

// waveList reads a wave file: one repo name per line, blank lines and
// # comments ignored. The batch wrapper's future Go port parses waves with
// this shape.
func waveList(content string) []string {
	var repos []string

	for line := range strings.SplitSeq(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		repos = append(repos, line)
	}

	return repos
}
