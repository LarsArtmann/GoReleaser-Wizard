package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	yaml "go.yaml.in/yaml/v3"
)

const deprecationsDocURL = "https://goreleaser.com/deprecations"

// DeprecationFinding is one deprecated key occurrence with its source location.
type DeprecationFinding struct {
	Key         string
	Line        int
	Replacement string
}

// deprecationDocURL returns the upstream doc anchor for a deprecated key.
func deprecationDocURL(key string) string {
	return deprecationsDocURL + "#" + strings.NewReplacer(".", "", " ", "").Replace(key)
}

// detectDeprecatedGoreleaserKeys parses goreleaser config content as YAML and
// walks the shapes GoReleaser v2.17 deprecated. The list is verified against
// goreleaser v2.17.1 `check` output (2026-11-07): dockers and docker_manifests
// (both replaced by dockers_v2), brews (replaced by homebrew_casks), and
// archives[].format_overrides[].format (string, replaced by the formats list).
// Findings carry the source line so users can jump straight to the key.
func detectDeprecatedGoreleaserKeys(content []byte) ([]DeprecationFinding, error) {
	var root yaml.Node

	if err := yaml.Unmarshal(content, &root); err != nil {
		return nil, domain.NewConfigError(
			domain.ErrInvalidFileFormat,
			"Failed to parse GoReleaser configuration as YAML",
			err.Error(),
		).WithCause(err)
	}

	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return nil, nil
	}

	configNode := root.Content[0]
	if configNode.Kind != yaml.MappingNode {
		return nil, nil
	}

	var findings []DeprecationFinding

	topLevelDeprecated := []struct {
		key         string
		replacement string
	}{
		{"dockers", "dockers_v2"},
		{"docker_manifests", "dockers_v2"},
		{"brews", "homebrew_casks"},
	}

	for _, dep := range topLevelDeprecated {
		if line, ok := mappingKeyLine(configNode, dep.key); ok {
			findings = append(findings, DeprecationFinding{
				Key:         dep.key,
				Line:        line,
				Replacement: dep.replacement,
			})
		}
	}

	findings = append(findings, detectArchiveFormatDeprecations(configNode)...)

	return findings, nil
}

// detectArchiveFormatDeprecations walks archives[*].format_overrides[*].format,
// the string form goreleaser replaced with the formats list.
func detectArchiveFormatDeprecations(configNode *yaml.Node) []DeprecationFinding {
	archivesNode := mappingValue(configNode, "archives")
	if archivesNode == nil || archivesNode.Kind != yaml.SequenceNode {
		return nil
	}

	var findings []DeprecationFinding

	for _, archive := range archivesNode.Content {
		if archive.Kind != yaml.MappingNode {
			continue
		}

		overrides := mappingValue(archive, "format_overrides")
		if overrides == nil || overrides.Kind != yaml.SequenceNode {
			continue
		}

		for _, override := range overrides.Content {
			if override.Kind != yaml.MappingNode {
				continue
			}

			if line, ok := mappingKeyLine(override, "format"); ok {
				findings = append(findings, DeprecationFinding{
					Key:         "archives.format_overrides.format",
					Line:        line,
					Replacement: "archives.format_overrides.formats",
				})
			}
		}
	}

	return findings
}

// detectDeprecatedKeysInFile reads and checks the goreleaser config at path.
func detectDeprecatedKeysInFile(path string) ([]DeprecationFinding, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, domain.NewSystemError(
			domain.ErrFileReadFailed,
			"Failed to read GoReleaser configuration",
			err.Error(),
			err,
		).WithContext(path)
	}

	return detectDeprecatedGoreleaserKeys(content)
}

// mappingKeyLine reports the line of a key inside a YAML mapping node (the
// key's own line, so users land on the deprecated key, not its first value).
func mappingKeyLine(mappingNode *yaml.Node, key string) (int, bool) {
	for i := 0; i+1 < len(mappingNode.Content); i += 2 {
		keyNode := mappingNode.Content[i]
		if keyNode.Value == key {
			return keyNode.Line, true
		}
	}

	return 0, false
}

// mappingValue returns the value node for a key inside a YAML mapping node.
func mappingValue(mappingNode *yaml.Node, key string) *yaml.Node {
	for i := 0; i+1 < len(mappingNode.Content); i += 2 {
		if mappingNode.Content[i].Value == key {
			return mappingNode.Content[i+1]
		}
	}

	return nil
}

// validateDeprecations appends deprecation findings to the validation results
// as structured warnings with before→after migration hints.
func validateDeprecations(results *ValidationResults, path string) {
	findings, err := detectDeprecatedKeysInFile(path)
	if err != nil {
		if domainErr, ok := errors.AsType[*domain.DomainError](err); ok {
			results.Warnings = append(results.Warnings, domainErr)
		}

		return
	}

	for _, finding := range findings {
		results.Warnings = append(results.Warnings,
			NewValidationError(
				"DEPRECATED_KEY",
				fmt.Sprintf(
					"Deprecated key %q at %s:%d",
					finding.Key, path, finding.Line,
				),
				fmt.Sprintf(
					"GoReleaser v2.17 deprecated %q; migrate to %q before the next major release removes it",
					finding.Key, finding.Replacement,
				),
			).
				WithContext(fmt.Sprintf("%s:%d", path, finding.Line)).
				WithSuggestion(fmt.Sprintf(
					"Replace %q with %q (before→after hints in docs/goreleaser-guide.md; full list: %s)",
					finding.Key, finding.Replacement, deprecationDocURL(finding.Key),
				)))

		results.Recommendations = append(results.Recommendations,
			fmt.Sprintf("Migrate %s:%d from %q to %q", path, finding.Line, finding.Key, finding.Replacement))
	}
}
