package main

import "time"

const (
	defaultTimeout    = 30 * time.Minute
	fullWizardTimeout = 10 * time.Minute
	configOnlyTimeout = 5 * time.Minute
	validationTimeout = 2 * time.Minute
	migrationTimeout  = 15 * time.Minute
	updateTimeout     = 10 * time.Minute
)

// Paths of the artifacts the wizard generates. Workflow paths use forward
// slashes: git repositories and GitHub Actions require them on every platform.
const (
	goreleaserConfigFilename   = ".goreleaser.yaml"
	dockerfileFilename         = "Dockerfile"
	releaseWorkflowTargetPath  = ".github/workflows/release.yml"
	forgejoCIWorkflowPath      = ".forgejo/workflows/ci.yml"
	forgejoReleaseWorkflowPath = ".forgejo/workflows/release.yml"
)
