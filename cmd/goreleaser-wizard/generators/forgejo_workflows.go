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

// Forgejo workflow file names under .forgejo/workflows.
const (
	forgejoWorkflowDir      = ".forgejo/workflows"
	forgejoCIWorkflowName   = "ci.yml"
	forgejoReleaseName      = "release.yml"
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
		name     string
		template string
	}{
		{forgejoCIWorkflowName, templates.ForgejoCITemplate},
		{forgejoReleaseName, templates.ForgejoReleaseTemplate},
	}

	if err := os.MkdirAll(forgejoWorkflowDir, directoryPermission); err != nil {
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

		if err := g.renderAndWrite(file.name, file.template); err != nil {
			return err
		}
	}

	g.logger.Info("Forgejo Actions workflows generated successfully", "dir", forgejoWorkflowDir)

	return nil
}

func (g *ForgejoWorkflowsGenerator) renderAndWrite(name, templateText string) error {
	tmpl, err := template.New(name).Delims("[[", "]]").Parse(templateText)
	if err != nil {
		return domain.NewConfigError(
			domain.ErrTemplateParsing,
			"Failed to parse Forgejo workflow template: "+name,
			err.Error(),
		).WithCause(err)
	}

	var output bytes.Buffer
	if err := tmpl.Execute(&output, g.templateData); err != nil {
		return domain.NewConfigError(
			domain.ErrTemplateRendering,
			"Failed to execute Forgejo workflow template: "+name,
			err.Error(),
		).WithCause(err)
	}

	workflowPath := filepath.Join(forgejoWorkflowDir, name)
	if err := WriteFile(workflowPath, output.Bytes(), filePermission); err != nil {
		return domain.NewFileError(
			domain.ErrFileOperation,
			"Failed to write Forgejo workflow: "+name,
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
		filepath.Join(forgejoWorkflowDir, forgejoCIWorkflowName),
		filepath.Join(forgejoWorkflowDir, forgejoReleaseName),
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
		forgejoWorkflowDir,
		".forgejo",
	})

	return nil
}

// UpdateConfig updates the generator configuration.
func (g *ForgejoWorkflowsGenerator) UpdateConfig(config *domain.SafeProjectConfig) {
	g.templateData = types.NewForgejoWorkflowsTemplateData(config)
}

// GetWorkflowPaths returns the paths of both generated workflows.
func (g *ForgejoWorkflowsGenerator) GetWorkflowPaths() []string {
	return []string{
		filepath.Join(forgejoWorkflowDir, forgejoCIWorkflowName),
		filepath.Join(forgejoWorkflowDir, forgejoReleaseName),
	}
}
