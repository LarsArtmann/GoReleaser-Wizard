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

// checkRelativeLinks walks every .md file under root (skipping archived
// history, which is closed) and verifies each relative link target exists.
func checkRelativeLinks(root string) []string {
	var findings []string

	linkPattern := regexp.MustCompile(`\[[^\]]*\]\(([^)]+)\)`)

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			// Closed history is not re-gated; node_modules/vendor never exist here.
			if d.Name() == "archived" || d.Name() == "node_modules" || d.Name() == "vendor" || d.Name() == ".git" {
				return filepath.SkipAll
			}

			return nil
		}

		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, readErr := os.ReadFile(path)
		if readErr != nil {
			findings = append(findings, path+": unreadable: "+readErr.Error())

			return nil
		}

		for _, match := range linkPattern.FindAllStringSubmatch(string(content), -1) {
			target := match[1]

			if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") ||
				strings.HasPrefix(target, "#") || strings.HasPrefix(target, "mailto:") ||
				target == "" {
				continue
			}

			// Strip anchors before resolving.
			resolved := target
			if idx := strings.Index(resolved, "#"); idx != -1 {
				resolved = resolved[:idx]
			}

			if resolved == "" {
				continue
			}

			base := filepath.Dir(path)
			full := filepath.Join(base, resolved)

			if _, statErr := os.Stat(full); statErr != nil {
				findings = append(findings, path+" -> "+target)
			}
		}

		return nil
	})
	if err != nil {
		findings = append(findings, "walk failed: "+err.Error())
	}

	return findings
}

// checkFeaturesCounts verifies FEATURES.md's counts table against its own
// status rows: the counts must be derivable, never hand-rounded.
func checkFeaturesCounts(featuresPath string) []string {
	content, err := os.ReadFile(featuresPath)
	if err != nil {
		return []string{featuresPath + ": unreadable: " + err.Error()}
	}

	var findings []string

	lines := strings.Split(string(content), "\n")

	// Find the counts table row(s): | <label> | <n> | style entries in a
	// section whose rows declare statuses.
	statusCounts := map[string]int{
		"FULLY_FUNCTIONAL":       0,
		"PARTIALLY_FUNCTIONAL":   0,
		"BROKEN":                 0,
		"PLANNED":                0,
		"DEMAND_GATED":           0,
		"ON_HOLD":                0,
		"OUT_OF_SCOPE":           0,
		"DISABLED":               0,
	}

	for _, line := range lines {
		for status := range statusCounts {
			if strings.Contains(line, "| "+status+" |") ||
				strings.Contains(line, "| "+status+"  ") {
				statusCounts[status]++
			}
		}
	}

	// Match declared totals: "13 FULLY_FUNCTIONAL" style phrases.
	declared := regexp.MustCompile(`(\d+)\s+([A-Z_]+)`)

	declaredTotal := map[string]int{}

	for _, line := range lines {
		for _, match := range declared.FindAllStringSubmatch(line, -1) {
			count := 0

			fmt.Sscanf(match[1], "%d", &count)

			status := match[2]
			if _, tracked := statusCounts[status]; tracked {
				declaredTotal[status] = count
			}
		}
	}

	for status, declared := range declaredTotal {
		if statusCounts[status] == 0 {
			continue
		}

		if statusCounts[status] != declared {
			findings = append(findings, fmt.Sprintf(
				"%s declares %d %s rows but the table holds %d",
				featuresPath, declared, status, statusCounts[status]))
		}
	}

	return findings
}
