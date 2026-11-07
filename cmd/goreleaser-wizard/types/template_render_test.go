package types

import (
	"strings"
	"testing"
	"text/template"

	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/templates"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

// homebrewRenderData is the data shape HomebrewTemplate expects. The template
// has no generator wiring yet (PLANNED in FEATURES); this render-only proof
// keeps it healthy until then.
type homebrewRenderData struct {
	FormulaName        string
	ProjectDescription string
	Homepage           string
	ArchiveURL         string
	Checksum           string
	License            string
	LDFlags            string
	MainPath           string
	BinaryName         string
	TestOutput         string
	Service            bool
}

func TestHomebrewTemplateRendersFormula(t *testing.T) {
	data := homebrewRenderData{
		FormulaName:        "mytool",
		ProjectDescription: "A tool",
		Homepage:           "https://example.com",
		ArchiveURL:         "https://github.com/octo/mytool/archive/v1.0.0.tar.gz",
		Checksum:           "abc123",
		License:            "MIT",
		LDFlags:            "-s -w",
		MainPath:           "./cmd/mytool",
		BinaryName:         "mytool",
		TestOutput:         "mytool version",
		Service:            false,
	}

	tmpl, err := template.New("homebrew").Parse(templates.HomebrewTemplate)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	var output strings.Builder

	if err := tmpl.Execute(&output, data); err != nil {
		t.Fatalf("execute: %v", err)
	}

	rendered := output.String()

	for _, want := range []string{
		"class mytool < Formula",
		`desc "A tool"`,
		`url "` + data.ArchiveURL + `"`,
		`ldflags: "-s -w"`,
		`"#{bin}/mytool --version"`,
	} {
		if !strings.Contains(rendered, want) {
			t.Errorf("rendered formula missing %q", want)
		}
	}

	if strings.Contains(rendered, "service do") {
		t.Error("service block rendered despite Service=false")
	}
}

func TestHomebrewTemplateOptionalServiceBlock(t *testing.T) {
	data := homebrewRenderData{FormulaName: "svc", BinaryName: "svc", Service: true}

	tmpl, err := template.New("homebrew").Parse(templates.HomebrewTemplate)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	var output strings.Builder

	if err := tmpl.Execute(&output, data); err != nil {
		t.Fatalf("execute: %v", err)
	}

	if !strings.Contains(output.String(), "service do") {
		t.Error("Service=true must render the service block")
	}
}

func TestNewGitHubActionsTemplateDataMapping(t *testing.T) {
	config := &domain.SafeProjectConfig{}
	config.ProjectName = "actions-app"
	config.DockerSupport = domain.DockerSupportNone
	config.SigningLevel = domain.SigningLevelNone
	config.ActionsOn = []domain.ActionTrigger{domain.ActionTriggerVersionTags}

	data := NewGitHubActionsTemplateData(config)

	if data.ProjectName != "actions-app" {
		t.Errorf("ProjectName = %q, want actions-app", data.ProjectName)
	}

	if data.DockerEnabled {
		t.Error("DockerEnabled = true for DockerSupportNone")
	}

	if data.SigningEnabled {
		t.Error("SigningEnabled = true for SigningLevelNone")
	}

	if len(data.Triggers) == 0 {
		t.Fatal("version-tag trigger must produce at least one GitHub pattern")
	}

	for _, trigger := range data.Triggers {
		if !strings.Contains(trigger, "tags") && trigger != "v*" && trigger != "v**" {
			t.Errorf("trigger %q does not look like a version-tag pattern", trigger)
		}
	}
}
