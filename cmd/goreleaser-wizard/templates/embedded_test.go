package templates

import (
	"testing"
	"text/template"
)

// TestAllEmbeddedTemplatesParse proves every template in the single source of
// truth is syntactically valid with its documented delimiters. Unwired
// templates rot silently without this: the homebrew template has no
// generator yet, and a syntax break would only surface when someone wires it.
// The GitHub Actions and Forgejo templates use [[ ]] delimiters so
// ${{ secrets.* }} expression literals need no escaping (gotcha 3).
func TestAllEmbeddedTemplatesParse(t *testing.T) {
	tests := []struct {
		name     string
		template string
		left     string
		right    string
	}{
		{name: "goreleaser", template: GoReleaserTemplate, left: "{{", right: "}}"},
		{name: "github_actions", template: GitHubActionsTemplate, left: "[[", right: "]]"},
		{name: "forgejo_ci", template: ForgejoCITemplate, left: "[[", right: "]]"},
		{name: "forgejo_release", template: ForgejoReleaseTemplate, left: "[[", right: "]]"},
		{name: "forgejo_actionlint", template: ForgejoActionlintTemplate, left: "[[", right: "]]"},
		{name: "dockerfile", template: DockerfileTemplate, left: "{{", right: "}}"},
		{name: "homebrew", template: HomebrewTemplate, left: "{{", right: "}}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := template.New(tt.name).Delims(tt.left, tt.right).Parse(tt.template); err != nil {
				t.Errorf("template %s does not parse with delims %s...%s: %v", tt.name, tt.left, tt.right, err)
			}
		})
	}
}
