package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

// setupTempGoModule creates a minimal Go module with a main package and a git
// repository, matching the layout the wizard is designed to run inside.
func setupTempGoModule(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()

	goMod := "module demo.example/e2e\n\ngo 1.25\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}

	mainContent := "package main\n\nfunc main() {}\n"
	mainDir := filepath.Join(tmpDir, "cmd", "demo")

	if err := os.MkdirAll(mainDir, workflowDirPermission); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(mainDir, "main.go"), []byte(mainContent), 0o644); err != nil {
		t.Fatal(err)
	}

	if out, err := exec.Command("git", "-C", tmpDir, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init failed: %v: %s", err, out)
	}

	return tmpDir
}

// generateIntoModule runs the full wizard workflow programmatically in the
// given directory and fails the test if generation errors.
func generateIntoModule(t *testing.T, tmpDir string, config *ProjectConfig) {
	t.Helper()

	t.Chdir(tmpDir)

	workflow, err := NewWorkflowBuilder(quietLogger()).BuildWorkflow(WorkflowTypeFullWizard, config, false)
	if err != nil {
		t.Fatalf("failed to build workflow: %v", err)
	}

	if err := workflow.Execute(context.Background()); err != nil {
		t.Fatalf("generation workflow failed: %v", err)
	}
}

// TestGeneratedConfigPassesGoReleaserCheck is the North Star proof: generate
// into a fresh module and require `goreleaser check` to exit zero with no
// deprecation warnings and no environment variables set.
func TestGeneratedConfigPassesGoReleaserCheck(t *testing.T) {
	goreleaserPath, err := exec.LookPath("goreleaser")
	if err != nil {
		t.Skip("goreleaser binary not found in PATH; skipping E2E proof")
	}

	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("docker binary not found in PATH; the docker-enabled E2E config requires it")
	}

	tmpDir := setupTempGoModule(t)

	config := baseTestProjectConfig()
	config.ProjectName = "e2e-app"
	config.BinaryName = "e2e-app"
	config.MainPath = "./cmd/demo"
	config.Platforms = []domain.Platform{"linux", "darwin"}
	config.Architectures = []domain.Architecture{"amd64", "arm64"}
	config.DockerSupport = domain.DockerSupportBoth
	config.DockerRegistry = domain.DockerRegistryGitHub

	generateIntoModule(t, tmpDir, &config)

	AssertFileExists(t, ".goreleaser.yaml", "full wizard must generate .goreleaser.yaml")
	AssertFileExists(t, "Dockerfile", "docker-enabled project must get a Dockerfile")
	AssertFileExists(t, releaseWorkflowTargetPath, "full wizard must generate release.yml")

	// Run goreleaser check with a clean environment: the generated config must
	// not depend on GITHUB_OWNER/GITHUB_REPO being set.
	cmd := exec.Command(goreleaserPath, "check")
	cmd.Dir = tmpDir
	cmd.Env = []string{"HOME=" + tmpDir, "PATH=" + os.Getenv("PATH")}

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("goreleaser check failed (exit %v):\n%s", err, output)
	}

	if strings.Contains(strings.ToUpper(string(output)), "DEPRECATED") {
		t.Errorf("goreleaser check reported deprecations:\n%s", output)
	}
}

// TestDualProviderGeneratedConfigPassesGoReleaserCheck extends the North Star
// proof to --ci-provider dual: both workflow sets generate and the Forgejo-
// targeted goreleaser config (gitea_urls release block) must still pass
// `goreleaser check` with zero deprecations and a clean environment.
func TestDualProviderGeneratedConfigPassesGoReleaserCheck(t *testing.T) {
	goreleaserPath, err := exec.LookPath("goreleaser")
	if err != nil {
		t.Skip("goreleaser binary not found in PATH; skipping E2E proof")
	}

	tmpDir := setupTempGoModule(t)

	// The Forgejo-targeted release block makes goreleaser check resolve the
	// origin remote (real Forgejo repos always have one); a non-resolvable
	// URL is enough, check never contacts the host.
	if out, err := exec.Command("git", "-C", tmpDir, "remote", "add", "origin",
		"https://forgejo.example.com/owner/demo.git").CombinedOutput(); err != nil {
		t.Fatalf("git remote add failed: %v: %s", err, out)
	}

	config := baseTestProjectConfig()
	config.ProjectName = "e2e-app"
	config.BinaryName = "e2e-app"
	config.MainPath = "./cmd/demo"
	config.CIProvider = domain.CIProviderDual

	generateIntoModule(t, tmpDir, &config)

	AssertFileExists(t, ".goreleaser.yaml", "dual provider must generate .goreleaser.yaml")
	AssertFileExists(t, releaseWorkflowTargetPath, "dual provider must generate the GitHub release workflow")
	AssertFileExists(t, forgejoCIWorkflowPath, "dual provider must generate the Forgejo CI workflow")
	AssertFileExists(t, forgejoReleaseWorkflowPath, "dual provider must generate the Forgejo release workflow")

	cmd := exec.Command(goreleaserPath, "check")
	cmd.Dir = tmpDir
	cmd.Env = []string{"HOME=" + tmpDir, "PATH=" + os.Getenv("PATH")}

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("goreleaser check failed for dual-provider config (exit %v):\n%s", err, output)
	}

	if strings.Contains(strings.ToUpper(string(output)), "DEPRECATED") {
		t.Errorf("goreleaser check reported deprecations for dual-provider config:\n%s", output)
	}
}

// TestNoneProviderGeneratesNoWorkflows proves --ci-provider none writes the
// goreleaser config and nothing else: no .github, no .forgejo artifacts.
func TestNoneProviderGeneratesNoWorkflows(t *testing.T) {
	tmpDir := setupTempGoModule(t)

	config := baseTestProjectConfig()
	config.ProjectName = "e2e-app"
	config.BinaryName = "e2e-app"
	config.MainPath = "./cmd/demo"
	config.CIProvider = domain.CIProviderNone

	generateIntoModule(t, tmpDir, &config)

	AssertFileExists(t, ".goreleaser.yaml", "none provider must still generate the goreleaser config")

	for _, forbidden := range []string{
		releaseWorkflowTargetPath,
		forgejoCIWorkflowPath,
		forgejoReleaseWorkflowPath,
		forgejoActionlintConfigPath,
	} {
		if _, err := os.Stat(filepath.Join(tmpDir, forbidden)); !os.IsNotExist(err) {
			t.Errorf("none provider must not generate %s", forbidden)
		}
	}
}

// TestForgejoWorkflowsPassActionLint generates Forgejo workflows for a
// forgejo-provider project and requires actionlint to accept them. The
// self-hosted runner label must be declared in an actionlint config, as on
// the real fleet.
func TestForgejoWorkflowsPassActionLint(t *testing.T) {
	actionlintPath, err := exec.LookPath("actionlint")
	if err != nil {
		t.Skip("actionlint binary not found in PATH; skipping Forgejo workflow lint")
	}

	tmpDir := setupTempGoModule(t)

	config := baseTestProjectConfig()
	config.ProjectName = "e2e-app"
	config.BinaryName = "e2e-app"
	config.MainPath = "./cmd/demo"
	config.CIProvider = domain.CIProviderForgejo

	generateIntoModule(t, tmpDir, &config)

	ciPath := filepath.Join(tmpDir, forgejoCIWorkflowPath)
	releasePath := filepath.Join(tmpDir, forgejoReleaseWorkflowPath)

	AssertFileExists(t, ciPath, "forgejo provider must generate the CI workflow")
	AssertFileExists(t, releasePath, "forgejo provider must generate the release workflow")

	// The generated actionlint config declares the self-hosted runner label,
	// so the workflows lint out of the box (no hand-written config needed).
	actionlintConfigPath := filepath.Join(tmpDir, forgejoActionlintConfigPath)
	AssertFileExists(t, actionlintConfigPath, "forgejo provider must generate the actionlint config")

	for _, path := range []string{ciPath, releasePath} {
		cmd := exec.Command(actionlintPath, "-config-file", actionlintConfigPath, path)
		cmd.Dir = tmpDir

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("actionlint rejected %s:\n%s", path, output)
		}
	}
}

// TestGenerateRefusesToOverwriteWithoutForce proves the atomicity contract:
// an existing target artifact fails the whole workflow before anything new
// is written.
func TestGenerateRefusesToOverwriteWithoutForce(t *testing.T) {
	tmpDir := setupTempGoModule(t)

	if err := os.WriteFile(filepath.Join(tmpDir, ".goreleaser.yaml"), []byte("version: 2\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Chdir(tmpDir)

	config := baseTestProjectConfig()
	config.ProjectName = "e2e-app"
	config.BinaryName = "e2e-app"
	config.MainPath = "./cmd/demo"

	workflow, err := NewWorkflowBuilder(quietLogger()).BuildWorkflow(WorkflowTypeFullWizard, &config, false)
	if err != nil {
		t.Fatalf("failed to build workflow: %v", err)
	}

	if err := workflow.Execute(context.Background()); err == nil {
		t.Fatal("workflow must fail when a target artifact exists without --force")
	}

	if _, err := os.Stat(releaseWorkflowTargetPath); !os.IsNotExist(err) {
		t.Error("no artifacts may be written when the preflight check fails")
	}
}
