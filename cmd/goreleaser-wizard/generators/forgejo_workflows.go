package generators

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"text/template"

	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/templates"
	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/types"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

// Forgejo workflow artifact paths, slash-separated (git and Forgejo Actions
// require forward slashes on every platform).
const (
	forgejoCIWorkflowRelPath = ".forgejo/workflows/ci.yml"
	forgejoReleaseRelPath    = ".forgejo/workflows/release.yml"
	forgejoActionlintRelPath = ".forgejo/actionlint.yaml"
)

// ForgejoWorkflowsGenerator generates the Forgejo Actions CI and release
// workflows. Forgejo Actions is GitHub-Actions-compatible, so the same [[ ]]
// delimiter strategy applies.
type ForgejoWorkflowsGenerator struct {
	templateData *types.ForgejoWorkflowsTemplateData
	logger       Logger
}

// NewForgejoWorkflowsGenerator creates a new Forgejo workflows generator.
func NewForgejoWorkflowsGenerator(
	config *domain.SafeProjectConfig,
	logger Logger,
) *ForgejoWorkflowsGenerator {
	return &ForgejoWorkflowsGenerator{
		templateData: types.NewForgejoWorkflowsTemplateData(config),
		logger:       logger,
	}
}

// Generate writes both Forgejo workflow files.
func (g *ForgejoWorkflowsGenerator) Generate(ctx context.Context) error {
	g.logger.Info("Generating Forgejo Actions workflows")

	files := []struct {
		path     string
		template string
	}{
		{forgejoCIWorkflowRelPath, templates.ForgejoCITemplate},
		{forgejoReleaseRelPath, templates.ForgejoReleaseTemplate},
		{forgejoActionlintRelPath, templates.ForgejoActionlintTemplate},
	}

	if err := os.MkdirAll(".forgejo/workflows", directoryPermission); err != nil {
		return domain.NewFileError(
			domain.ErrFileOperation,
			"Failed to create Forgejo workflow directory",
			err.Error(),
		).WithCause(err)
	}

	for _, file := range files {
		if err := CheckContext(ctx); err != nil {
			return err
		}

		if err := g.renderAndWriteAt(filepath.FromSlash(file.path), file.template); err != nil {
			return err
		}
	}

	g.logger.Info("Forgejo Actions workflows generated successfully", "dir", ".forgejo")

	return nil
}

// renderAndWriteAt renders templateText and writes it to the given relative
// path. Failed parses and writes surface as typed domain errors.
func (g *ForgejoWorkflowsGenerator) renderAndWriteAt(relPath, templateText string) error {
	tmpl, err := template.New(relPath).Delims("[[", "]]").Parse(templateText)
	if err != nil {
		return domain.NewConfigError(
			domain.ErrTemplateParsing,
			"Failed to parse Forgejo workflow template: "+relPath,
			err.Error(),
		).WithCause(err)
	}

	var output bytes.Buffer
	if err := tmpl.Execute(&output, g.templateData); err != nil {
		return domain.NewConfigError(
			domain.ErrTemplateRendering,
			"Failed to execute Forgejo workflow template: "+relPath,
			err.Error(),
		).WithCause(err)
	}

	if err := WriteFile(relPath, output.Bytes(), filePermission); err != nil {
		return domain.NewFileError(
			domain.ErrFileOperation,
			"Failed to write Forgejo workflow artifact: "+relPath,
			err.Error(),
		).WithCause(err)
	}

	return nil
}

// ValidateTemplate validates both Forgejo workflow templates.
func (g *ForgejoWorkflowsGenerator) ValidateTemplate() error {
	if err := ValidateTemplateWithDelims("forgejo-ci", "[[", "]]", templates.ForgejoCITemplate); err != nil {
		return err
	}

	return ValidateTemplateWithDelims("forgejo-release", "[[", "]]", templates.ForgejoReleaseTemplate)
}

// GeneratePreview generates a preview of both workflows without writing files.
func (g *ForgejoWorkflowsGenerator) GeneratePreview(ctx context.Context) (string, error) {
	ciPreview, err := GeneratePreviewWithDelims(
		ctx,
		g.logger,
		"forgejo-ci",
		"[[", "]]",
		templates.ForgejoCITemplate,
		"Generating Forgejo CI workflow preview",
		g.templateData,
	)
	if err != nil {
		return "", err
	}

	releasePreview, err := GeneratePreviewWithDelims(
		ctx,
		g.logger,
		"forgejo-release",
		"[[", "]]",
		templates.ForgejoReleaseTemplate,
		"Generating Forgejo release workflow preview",
		g.templateData,
	)
	if err != nil {
		return "", err
	}

	return ciPreview + "\n---\n" + releasePreview, nil
}

// Rollback removes the generated Forgejo workflows.
func (g *ForgejoWorkflowsGenerator) Rollback(ctx context.Context) error {
	g.logger.Info("Rolling back Forgejo workflow generation")

	if err := CheckContext(ctx); err != nil {
		return err
	}

	paths := []string{
		forgejoCIWorkflowRelPath,
		forgejoReleaseRelPath,
		forgejoActionlintRelPath,
	}
	for _, workflowPath := range paths {
		err := removeGeneratedFile(
			g.logger,
			workflowPath,
			"Failed to remove generated workflow",
			"Removed generated workflow",
		)
		if err != nil {
			return err
		}
	}

	removeEmptyDirectories(g.logger, []string{
		".forgejo/workflows",
		".forgejo",
	})

	return nil
}

// UpdateConfig updates the generator configuration.
func (g *ForgejoWorkflowsGenerator) UpdateConfig(config *domain.SafeProjectConfig) {
	g.templateData = types.NewForgejoWorkflowsTemplateData(config)
}

// GetWorkflowPaths returns the paths of all generated Forgejo artifacts.
func (g *ForgejoWorkflowsGenerator) GetWorkflowPaths() []string {
	return []string{
		forgejoCIWorkflowRelPath,
		forgejoReleaseRelPath,
		forgejoActionlintRelPath,
	}
}
