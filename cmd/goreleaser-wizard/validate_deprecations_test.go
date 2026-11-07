package main

import (
	"os"
	"strings"
	"testing"
)

const deprecatedConfigFixture = `version: 2
project_name: legacy-app
builds:
  - id: legacy-app
    main: .
    binary: legacy-app
    goos: [linux]
    goarch: [amd64]
dockers:
  - goos: linux
    goarch: amd64
docker_manifests:
  - name_template: legacy-app:latest
brews:
  - name: legacy-app
    repository:
      owner: octocat
      name: homebrew-tap
archives:
  - id: legacy-app
    format_overrides:
      - goos: windows
        format: zip
`

const modernConfigFixture = `version: 2
project_name: modern-app
builds:
  - id: modern-app
    main: .
    binary: modern-app
    goos: [linux]
    goarch: [amd64]
dockers_v2:
  - images:
      - modern-app
archives:
  - id: modern-app
    formats: [tar.gz]
    format_overrides:
      - goos: windows
        formats: [zip]
`

func TestDetectDeprecatedGoreleaserKeys(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		wantKeys      []string
		wantMinLine   map[string]int
		wantNotInside []string
	}{
		{
			name:     "deprecated_config_lists_all_four_keys",
			content:  deprecatedConfigFixture,
			wantKeys: []string{"dockers", "docker_manifests", "brews", "archives.format_overrides.format"},
		},
		{
			name:     "modern_config_is_clean",
			content:  modernConfigFixture,
			wantKeys: nil,
		},
		{
			name:          "nested_format_is_found_not_the_archive_format_list",
			content:       modernConfigFixture,
			wantNotInside: []string{"archives.format_overrides.format"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings, err := detectDeprecatedGoreleaserKeys([]byte(tt.content))
			if err != nil {
				t.Fatalf("detectDeprecatedGoreleaserKeys() = %v, want nil", err)
			}

			gotKeys := map[string]DeprecationFinding{}
			for _, finding := range findings {
				gotKeys[finding.Key] = finding
			}

			for _, want := range tt.wantKeys {
				finding, ok := gotKeys[want]
				if !ok {
					t.Errorf("finding for key %q missing; got %v", want, findings)
				}

				if finding.Line < 1 {
					t.Errorf("finding for %q has line %d, want >= 1", want, finding.Line)
				}

				if finding.Replacement == "" {
					t.Errorf("finding for %q has empty replacement", want)
				}
			}

			for _, absent := range tt.wantNotInside {
				if _, ok := gotKeys[absent]; ok {
					t.Errorf("unexpected finding for %q", absent)
				}
			}
		})
	}
}

func TestDetectDeprecatedKeyLinesPointAtSource(t *testing.T) {
	findings, err := detectDeprecatedGoreleaserKeys([]byte(deprecatedConfigFixture))
	if err != nil {
		t.Fatalf("detect = %v", err)
	}

	lines := map[string]int{}
	for _, finding := range findings {
		lines[finding.Key] = finding.Line
	}

	sourceLines := map[string]int{}

	for i, line := range strings.Split(deprecatedConfigFixture, "\n") {
		trimmed := strings.TrimSpace(line)
		for key := range lines {
			// Only top-level keys have a whole-line match; the nested
			// archives.format_overrides.format is line-asserted in the table.
			if trimmed == key+":" && !strings.Contains(key, ".") {
				sourceLines[key] = i + 1
			}
		}
	}

	for key, want := range sourceLines {
		if lines[key] != want {
			t.Errorf("key %q reported at line %d, want %d", key, lines[key], want)
		}
	}
}

func TestDetectDeprecatedKeysMalformedYAMLIsTypedError(t *testing.T) {
	_, err := detectDeprecatedGoreleaserKeys([]byte("not: [valid: yaml"))
	if err == nil {
		t.Fatal("malformed YAML must return an error")
	}

	if !strings.Contains(err.Error(), "YAML") {
		t.Errorf("error %q does not mention YAML", err.Error())
	}
}

func TestDetectDeprecatedKeysEmptyDocumentIsClean(t *testing.T) {
	findings, err := detectDeprecatedGoreleaserKeys([]byte(""))
	if err != nil {
		t.Fatalf("empty document returned error: %v", err)
	}

	if len(findings) != 0 {
		t.Errorf("empty document yielded findings: %v", findings)
	}
}

func TestValidateDeprecationsPopulatesResults(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/.goreleaser.yaml"

	if err := os.WriteFile(path, []byte(deprecatedConfigFixture), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	results := &ValidationResults{}
	validateDeprecations(results, path)

	if len(results.Warnings) == 0 {
		t.Fatal("deprecated fixture produced no warnings")
	}

	var builder strings.Builder

	for _, warning := range results.Warnings {
		builder.WriteString(warning.Message)
		builder.WriteString(" ")
		builder.WriteString(warning.Suggestion)

		if warning.Code != "DEPRECATED_KEY" {
			t.Errorf("warning code = %q, want DEPRECATED_KEY", warning.Code)
		}
	}

	joinedWarnings := builder.String()

	for _, want := range []string{
		`"dockers"`,
		"docs/goreleaser-guide.md",
	} {
		if !strings.Contains(joinedWarnings, want) {
			t.Errorf("warnings do not mention %q; got: %s", want, joinedWarnings)
		}
	}

	if len(results.Recommendations) == 0 {
		t.Error("deprecated fixture produced no recommendations")
	}
}

func TestValidateDeprecationsCleanConfigStaysSilent(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/.goreleaser.yaml"

	if err := os.WriteFile(path, []byte(modernConfigFixture), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	results := &ValidationResults{}
	validateDeprecations(results, path)

	if len(results.Warnings) != 0 {
		t.Errorf("modern config produced warnings: %v", results.Warnings)
	}
}
