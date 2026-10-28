package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/generators"
	"github.com/LarsArtmann/GoReleaser-Wizard/cmd/goreleaser-wizard/types"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

// update allows regenerating golden files: go test ./... -update-golden.
var updateGolden = flag.Bool("update-golden", false, "update golden files instead of comparing")

// goldenConfig builds a deterministic config: fixed GitHub target so goldens
// never depend on the developer's git remote.
func goldenConfig(
	t *testing.T,
	docker domain.DockerSupport,
	signing domain.SigningLevel,
	projectType domain.ProjectType,
) *ProjectConfig {
	t.Helper()

	types.SetGitHubOwnerOverride("golden-owner")
	types.SetGitHubRepoOverride("golden-repo")
	t.Cleanup(func() {
		types.SetGitHubOwnerOverride("")
		types.SetGitHubRepoOverride("")
	})

	cfg := baseTestProjectConfig()
	cfg.ProjectName = "golden-app"
	cfg.BinaryName = "golden-app"
	cfg.MainPath = "./cmd/golden-app"
	cfg.ProjectType = projectType
	cfg.Platforms = []domain.Platform{"linux", "darwin"}
	cfg.Architectures = []domain.Architecture{"amd64", "arm64"}
	cfg.DockerSupport = docker
	cfg.DockerRegistry = domain.DockerRegistryGitHub
	cfg.SigningLevel = signing

	return &cfg
}

// assertGolden compares generated output against the golden file, writing the
// golden file when -update-golden is set.
func assertGolden(t *testing.T, goldenName, output string) {
	t.Helper()

	goldenPath := filepath.Join("testdata", "golden", goldenName)

	if *updateGolden {
		if err := os.MkdirAll(filepath.Dir(goldenPath), workflowDirPermission); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(goldenPath, []byte(output), 0o644); err != nil {
			t.Fatal(err)
		}

		return
	}

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("golden file missing (run go test ./cmd/goreleaser-wizard -update-golden): %v", err)
	}

	if string(want) != output {
		t.Errorf("generated output differs from golden %s\n--- golden ---\n%s\n--- generated ---\n%s",
			goldenName, want, output)
	}
}

// discardGeneratorLogger satisfies generators.Logger without output.
type discardGeneratorLogger struct{}

func (discardGeneratorLogger) Debug(string, ...any) {}
func (discardGeneratorLogger) Info(string, ...any)  {}
func (discardGeneratorLogger) Warn(string, ...any)  {}
func (discardGeneratorLogger) Error(string, ...any) {}

func TestGoldenGoReleaserConfig(t *testing.T) {
	tests := []struct {
		name   string
		golden string
		docker domain.DockerSupport
	}{
		{
			name:   "no_docker",
			golden: "goreleaser-nodocker.yaml",
			docker: domain.DockerSupportNone,
		},
		{
			name:   "docker_enabled",
			golden: "goreleaser-docker.yaml",
			docker: domain.DockerSupportBoth,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := goldenConfig(t, tt.docker, domain.SigningLevelBasic, domain.ProjectTypeCLI)

			generator := generators.NewGoReleaserGenerator(config, discardGeneratorLogger{})

			output, err := generator.GeneratePreview(context.Background())
			if err != nil {
				t.Fatalf("GeneratePreview failed: %v", err)
			}

			assertGolden(t, tt.golden, output)
		})
	}
}

func TestGoldenGitHubActionsWorkflow(t *testing.T) {
	combinations := []struct {
		name    string
		golden  string
		docker  domain.DockerSupport
		signing domain.SigningLevel
	}{
		{
			name:    "docker_and_signing",
			golden:  "release-docker-signing.yml",
			docker:  domain.DockerSupportBoth,
			signing: domain.SigningLevelBasic,
		},
		{
			name:    "docker_only",
			golden:  "release-docker.yml",
			docker:  domain.DockerSupportBoth,
			signing: domain.SigningLevelNone,
		},
		{
			name:    "signing_only",
			golden:  "release-signing.yml",
			docker:  domain.DockerSupportNone,
			signing: domain.SigningLevelBasic,
		},
		{
			name:    "neither",
			golden:  "release-plain.yml",
			docker:  domain.DockerSupportNone,
			signing: domain.SigningLevelNone,
		},
	}

	for _, tt := range combinations {
		t.Run(tt.name, func(t *testing.T) {
			config := goldenConfig(t, tt.docker, tt.signing, domain.ProjectTypeCLI)

			generator := generators.NewGitHubActionsGenerator(config, discardGeneratorLogger{})

			output, err := generator.GeneratePreview(context.Background())
			if err != nil {
				t.Fatalf("GeneratePreview failed: %v", err)
			}

			assertGolden(t, tt.golden, output)
		})
	}
}

func TestGoldenForgejoWorkflows(t *testing.T) {
	config := goldenConfig(t, domain.DockerSupportNone, domain.SigningLevelNone, domain.ProjectTypeCLI)
	config.CIProvider = domain.CIProviderForgejo

	generator := generators.NewForgejoWorkflowsGenerator(config, discardGeneratorLogger{})

	if err := generator.ValidateTemplate(); err != nil {
		t.Fatalf("template validation failed: %v", err)
	}

	output, err := generator.GeneratePreview(context.Background())
	if err != nil {
		t.Fatalf("GeneratePreview failed: %v", err)
	}

	assertGolden(t, "forgejo-workflows.yml", output)
}

// TestGoldenGoReleaserForgejo pins the forgejo-provider GoReleaser config:
// gitea_urls targeting the Forgejo instance, changelog sourced from git,
// and no github release block.
func TestGoldenGoReleaserForgejo(t *testing.T) {
	config := goldenConfig(t, domain.DockerSupportNone, domain.SigningLevelNone, domain.ProjectTypeCLI)
	config.CIProvider = domain.CIProviderForgejo

	generator := generators.NewGoReleaserGenerator(config, discardGeneratorLogger{})

	output, err := generator.GeneratePreview(context.Background())
	if err != nil {
		t.Fatalf("GeneratePreview failed: %v", err)
	}

	assertGolden(t, "goreleaser-forgejo.yaml", output)
}

// TestGoldenDualProvider verifies the transition state: both workflow sets
// render, while the GoReleaser config targets Forgejo (the GitHub copy is
// maintained by the release shadow, not goreleaser).
func TestGoldenDualProvider(t *testing.T) {
	config := goldenConfig(t, domain.DockerSupportNone, domain.SigningLevelNone, domain.ProjectTypeCLI)
	config.CIProvider = domain.CIProviderDual

	if !config.CIProvider.GeneratesGitHubWorkflows() || !config.CIProvider.GeneratesForgejoWorkflows() {
		t.Fatal("dual provider must generate both workflow sets")
	}

	goreleaserOutput, err := generators.NewGoReleaserGenerator(config, discardGeneratorLogger{}).
		GeneratePreview(context.Background())
	if err != nil {
		t.Fatalf("goreleaser GeneratePreview failed: %v", err)
	}

	assertGolden(t, "goreleaser-dual.yaml", goreleaserOutput)

	forgejoOutput, err := generators.NewForgejoWorkflowsGenerator(config, discardGeneratorLogger{}).
		GeneratePreview(context.Background())
	if err != nil {
		t.Fatalf("forgejo GeneratePreview failed: %v", err)
	}

	if forgejoOutput == "" {
		t.Error("dual provider must render Forgejo workflows")
	}
}

func TestGoldenDockerfile(t *testing.T) {
	tests := []struct {
		name        string
		golden      string
		cgo         domain.CGOStatus
		projectType domain.ProjectType
	}{
		{
			name:        "cli_static_scratch",
			golden:      "Dockerfile-cli-scratch",
			cgo:         domain.CGOStatusDisabled,
			projectType: domain.ProjectTypeCLI,
		},
		{
			name:        "cgo_alpine",
			golden:      "Dockerfile-cgo-alpine",
			cgo:         domain.CGOStatusEnabled,
			projectType: domain.ProjectTypeCLI,
		},
		{
			name:        "webapi_exposes_port",
			golden:      "Dockerfile-webapi",
			cgo:         domain.CGOStatusDisabled,
			projectType: domain.ProjectTypeWebAPI,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := goldenConfig(t, domain.DockerSupportBoth, domain.SigningLevelNone, tt.projectType)
			config.CGOStatus = tt.cgo

			generator := generators.NewDockerfileGenerator(config, discardGeneratorLogger{})

			output, err := generator.GeneratePreview(context.Background())
			if err != nil {
				t.Fatalf("GeneratePreview failed: %v", err)
			}

			assertGolden(t, tt.golden, output)
		})
	}
}
