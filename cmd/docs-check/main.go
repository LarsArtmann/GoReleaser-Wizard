// Command docs-check runs the standing documentation gates over the repo:
// every relative markdown link resolves, and FEATURES.md's honest-counts
// table totals match its rows. Hand-maintained counts rot; this gate
// re-derives them on every CI run.
//
// Usage:
//
//	docs-check [repo-root]
//
// Exit 0 = all gates green, 1 = findings, 2 = usage error.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "usage: docs-check [repo-root]")

		os.Exit(2)
	}

	exitCode := 0

	for _, finding := range checkRelativeLinks(root) {
		fmt.Fprintf(os.Stderr, "BROKEN LINK: %s\n", finding)

		exitCode = 1
	}

	for _, finding := range checkFeaturesCounts(filepath.Join(root, "FEATURES.md")) {
		fmt.Fprintf(os.Stderr, "COUNTS MISMATCH: %s\n", finding)

		exitCode = 1
	}

	if exitCode == 0 {
		fmt.Fprintln(os.Stdout, "docs-check: all gates green")
	}

	os.Exit(exitCode)
}

// markdownLinkPattern captures group 1 = the link target.
var markdownLinkPattern = regexp.MustCompile(`\[[^\]]*\]\(([^)]+)\)`)

// checkRelativeLinks walks every .md file under root (skipping archived
// history, which is closed) and verifies each relative link target exists.
//
//nolint:gosec // G703: read-only local docs walk over the repo being gated
func checkRelativeLinks(root string) []string {
	var findings []string

	walkErr := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if skip, reason := skipPath(entry); skip {
			if reason == skipReasonDirectory {
				return filepath.SkipAll
			}

			return nil
		}

		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		findings = append(findings, checkFileLinks(path)...)

		return nil
	})
	if walkErr != nil {
		findings = append(findings, "walk failed: "+walkErr.Error())
	}

	return findings
}

type skipReason int

const (
	skipReasonNone skipReason = iota
	skipReasonDirectory
)

// skipPath reports whether a walked entry is excluded from gating.
func skipPath(entry os.DirEntry) (bool, skipReason) {
	switch entry.Name() {
	case "archived", "node_modules", "vendor", ".git":
		if entry.IsDir() {
			return true, skipReasonDirectory
		}

		return true, skipReasonNone
	default:
		return false, skipReasonNone
	}
}

// checkFileLinks verifies relative links in one markdown file.
func checkFileLinks(path string) []string {
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		return []string{path + ": unreadable: " + readErr.Error()}
	}

	var findings []string

	for _, match := range markdownLinkPattern.FindAllStringSubmatch(string(content), -1) {
		resolved := relativeTarget(match[1])
		if resolved == "" {
			continue
		}

		full := filepath.Join(filepath.Dir(path), resolved)

		//nolint:gosec // G703: path derives from the walked repo tree
		if _, statErr := os.Stat(full); statErr != nil {
			findings = append(findings, path+" -> "+match[1])
		}
	}

	return findings
}

// relativeTarget strips schemes, anchors and pure-hash links from a markdown
// link target, returning the filesystem path to verify ("" = nothing to check).
func relativeTarget(target string) string {
	if isExternal(target) {
		return ""
	}

	resolved := target
	if idx := strings.Index(resolved, "#"); idx != -1 {
		resolved = resolved[:idx]
	}

	return resolved
}

func isExternal(target string) bool {
	for _, prefix := range []string{"http://", "https://", "#", "mailto:", ""} {
		if strings.HasPrefix(target, prefix) && prefix != "" {
			return true
		}
	}

	return target == ""
}

// statusVocabulary is the FEATURES.md status set (docs-health canon).
//
//nolint:gochecknoglobals // fixed vocabulary table, never mutated
var statusVocabulary = []string{
	"FULLY_FUNCTIONAL",
	"PARTIALLY_FUNCTIONAL",
	"BROKEN",
	"PLANNED",
	"DEMAND_GATED",
	"ON_HOLD",
	"OUT_OF_SCOPE",
	"DISABLED",
}

// declaredCountsPattern matches "13 FULLY_FUNCTIONAL" style totals.
var declaredCountsPattern = regexp.MustCompile(`(\d+)\s+([A-Z_]+)`)

// checkFeaturesCounts verifies FEATURES.md's counts table against its own
// status rows: the counts must be derivable, never hand-rounded.
//
//nolint:gosec // G703: featuresPath is repo-root-joined inside main
func checkFeaturesCounts(featuresPath string) []string {
	content, err := os.ReadFile(featuresPath)
	if err != nil {
		return []string{featuresPath + ": unreadable: " + err.Error()}
	}

	rowCounts := countStatusRows(string(content))
	declaredTotals := declaredTotals(string(content), rowCounts)

	var findings []string

	for status, declared := range declaredTotals {
		if rowCounts[status] == 0 {
			continue
		}

		if rowCounts[status] != declared {
			findings = append(findings, fmt.Sprintf(
				"%s declares %d %s rows but the table holds %d",
				featuresPath, declared, status, rowCounts[status]))
		}
	}

	return findings
}

// countStatusRows counts table rows carrying each status literal.
func countStatusRows(content string) map[string]int {
	counts := make(map[string]int, len(statusVocabulary))

	for line := range strings.SplitSeq(content, "\n") {
		for _, status := range statusVocabulary {
			if strings.Contains(line, "| "+status+" |") {
				counts[status]++
			}
		}
	}

	return counts
}

// declaredTotals extracts declared totals for statuses present in the rows.
func declaredTotals(content string, rowCounts map[string]int) map[string]int {
	totals := make(map[string]int, len(rowCounts))

	for line := range strings.SplitSeq(content, "\n") {
		for _, match := range declaredCountsPattern.FindAllStringSubmatch(line, -1) {
			count, parseErr := strconv.Atoi(match[1])
			if parseErr != nil {
				continue
			}

			status := match[2]
			if rowCounts[status] > 0 {
				totals[status] = count
			}
		}
	}

	return totals
}
