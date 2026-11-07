package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckRelativeLinksFindsBrokenLinks(t *testing.T) {
	repo := t.TempDir()

	if err := os.WriteFile(filepath.Join(repo, "README.md"),
		[]byte("[good](FEATURES.md) [broken](MISSING.md) [external](https://example.com) [anchor](#section)\n"),
		0o644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(repo, "FEATURES.md"), []byte("# Features\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	findings := checkRelativeLinks(repo)

	if len(findings) != 1 {
		t.Fatalf("findings = %v, want exactly the broken link", findings)
	}

	if !strings.Contains(findings[0], "MISSING.md") {
		t.Errorf("finding %q does not name the broken target", findings[0])
	}
}

func TestCheckRelativeLinksSkipsArchived(t *testing.T) {
	repo := t.TempDir()

	if err := os.MkdirAll(filepath.Join(repo, "docs", "status", "archived"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(repo, "docs", "status", "archived", "old.md"),
		[]byte("[broken](GONE.md)\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if findings := checkRelativeLinks(repo); len(findings) != 0 {
		t.Errorf("archived history must not be gated; findings: %v", findings)
	}
}

func TestCheckFeaturesCountsDetectsMismatch(t *testing.T) {
	repo := t.TempDir()

	features := filepath.Join(repo, "FEATURES.md")

	content := "# Features\n\n## Status\n\n" +
		"| Feature | Status | Evidence |\n" +
		"| ------- | ------ | -------- |\n" +
		"| a | FULLY_FUNCTIONAL | x.go:1 |\n" +
		"| b | FULLY_FUNCTIONAL | y.go:2 |\n" +
		"| c | PLANNED | - |\n" +
		"\nCounts: 3 FULLY_FUNCTIONAL and 1 PLANNED.\n"

	if err := os.WriteFile(features, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	findings := checkFeaturesCounts(features)

	if len(findings) != 1 {
		t.Fatalf("findings = %v, want exactly the FULLY_FUNCTIONAL mismatch", findings)
	}

	if !strings.Contains(findings[0], "FULLY_FUNCTIONAL") {
		t.Errorf("finding %q does not name the mismatched status", findings[0])
	}
}

func TestCheckFeaturesCountsConsistentCountsPass(t *testing.T) {
	repo := t.TempDir()

	features := filepath.Join(repo, "FEATURES.md")

	content := "# Features\n\n" +
		"| a | FULLY_FUNCTIONAL | x |\n" +
		"| b | FULLY_FUNCTIONAL | y |\n" +
		"| c | PLANNED | - |\n" +
		"\nTotals: 2 FULLY_FUNCTIONAL, 1 PLANNED.\n"

	if err := os.WriteFile(features, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	if findings := checkFeaturesCounts(features); len(findings) != 0 {
		t.Errorf("consistent counts flagged: %v", findings)
	}
}
