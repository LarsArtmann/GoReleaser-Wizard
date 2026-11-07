package types

import (
	"strings"
	"testing"
	"text/template"

	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/templates"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

// TestForgejoPinsAreExactValues is the tamper guard for the Forgejo workflow
// pins: the whole Forgejo story rests on reproducible CI gates, so the pins
// must never float. A bump must land here FIRST (deliberate, reviewed), then
// flow through the templates and goldens. If this test fails, someone changed
// a pin without documenting the upgrade.
func TestForgejoPinsAreExactValues(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{
			name: "buildflow_flake_pin",
			got:  BuildFlowNixPin,
			want: "github:LarsArtmann/BuildFlow/49ced3e0039a5f6e88028697bd779ec2f7f28826",
		},
		{
			name: "goreleaser_nixpkgs_pin",
			got:  GoReleaserNixPin,
			want: "github:NixOS/nixpkgs/3ec636362e1816bc2f0888fa0fa2c69f1bc9e504#goreleaser",
		},
		{
			name: "forgejo_runner_label",
			got:  ForgejoRunnerLabel,
			want: "native",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("pin %s = %q, want %q (bumps are deliberate: update this test, the templates, and the goldens together)", tt.name, tt.got, tt.want)
			}
		})
	}
}

// TestForgejoPinsAreWellFormed rejects malformed pins a hand-edit could
// introduce: a flake ref needs an exact full-length revision, not a branch or
// short hash.
func TestForgejoPinsAreWellFormed(t *testing.T) {
	tests := []struct {
		name       string
		pin        string
		wantScheme string
		revLength  int
	}{
		{
			name:       "buildflow",
			pin:        BuildFlowNixPin,
			wantScheme: "github:",
			revLength:  40,
		},
		{
			name:       "goreleaser",
			pin:        GoReleaserNixPin,
			wantScheme: "github:",
			revLength:  40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.HasPrefix(tt.pin, tt.wantScheme) {
				t.Errorf("pin %q must start with %q", tt.pin, tt.wantScheme)
			}

			ref := strings.TrimPrefix(tt.pin, tt.wantScheme)
			rev := ref[strings.LastIndex(ref, "/")+1:]
			if attr := strings.Index(rev, "#"); attr != -1 {
				rev = rev[:attr]
			}

			if len(rev) != tt.revLength {
				t.Errorf("pin revision %q is %d chars, want full %d-char git revision (never a branch or short hash)", rev, len(rev), tt.revLength)
			}
		})
	}
}

// TestForgejoTemplateDataWiresPins proves the template data constructor
// forwards the constants into the rendered fields: a renamed constant or a
// rewired constructor must not silently drop a pin from the output.
func TestForgejoTemplateDataWiresPins(t *testing.T) {
	config := baseForgejoTestConfig()

	data := NewForgejoWorkflowsTemplateData(config)

	if data.BuildFlowPin != BuildFlowNixPin {
		t.Errorf("BuildFlowPin = %q, want the pinned constant %q", data.BuildFlowPin, BuildFlowNixPin)
	}

	if data.GoReleaserPin != GoReleaserNixPin {
		t.Errorf("GoReleaserPin = %q, want the pinned constant %q", data.GoReleaserPin, GoReleaserNixPin)
	}

	if data.RunnerLabel != ForgejoRunnerLabel {
		t.Errorf("RunnerLabel = %q, want the pinned constant %q", data.RunnerLabel, ForgejoRunnerLabel)
	}
}

// TestForgejoTemplatesUsePins proves both Forgejo workflow templates render
// every pin: a template edit that drops a pin field turns the generated CI
// gate into a floating ref without any other test noticing.
func TestForgejoTemplatesUsePins(t *testing.T) {
	config := baseForgejoTestConfig()

	data := NewForgejoWorkflowsTemplateData(config)

	tests := []struct {
		name      string
		template  string
		mustOccur []string
	}{
		{
			name:      "ci_workflow",
			template:  templates.ForgejoCITemplate,
			mustOccur: []string{data.BuildFlowPin, data.RunnerLabel},
		},
		{
			name:      "release_workflow",
			template:  templates.ForgejoReleaseTemplate,
			mustOccur: []string{data.BuildFlowPin, data.GoReleaserPin, data.RunnerLabel, data.ForgejoURL},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rendered := renderForgejoTemplate(t, tt.template, data)

			for _, want := range tt.mustOccur {
				if !strings.Contains(rendered, want) {
					t.Errorf("rendered %s does not contain %q; a pin or field went missing", tt.name, want)
				}
			}
		})
	}
}

func baseForgejoTestConfig() *domain.SafeProjectConfig {
	config := &domain.SafeProjectConfig{}
	config.ProjectName = "pin-test-app"
	config.BinaryName = "pin-test-app"
	config.MainPath = "."
	config.ProjectType = domain.ProjectTypeCLI
	config.CIProvider = domain.CIProviderForgejo

	return config
}

func renderForgejoTemplate(t *testing.T, templateText string, data *ForgejoWorkflowsTemplateData) string {
	t.Helper()

	tmpl, err := template.New("test").Delims("[[", "]]").Parse(templateText)
	if err != nil {
		t.Fatalf("template parse failed: %v", err)
	}

	var output strings.Builder

	if err := tmpl.Execute(&output, data); err != nil {
		t.Fatalf("template execute failed: %v", err)
	}

	return output.String()
}
