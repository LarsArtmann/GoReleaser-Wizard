package types

import (
	"context"
	"os/exec"
	"slices"
	"strings"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/LarsArtmann/GoReleaser-Wizard/internal/git"
)

// GoReleaserTemplateData is the strongly typed template data for the GoReleaser
// configuration template. Every field maps to a template key; version, tag and
// commit fields are intentionally absent because GoReleaser evaluates those
// itself at release time.
type GoReleaserTemplateData struct {
	// Project information
	ProjectName string
	BinaryName  string
	MainPath    string

	// Build configuration
	CGOEnabled         string
	Platforms          []string
	Architectures      []string
	BuildTags          []string
	IgnoreCombinations []PlatformIgnored

	// Release configuration
	DockerEnabled   bool
	DockerRegistry  string
	DockerImage     string
	DockerPlatforms []string

	// GitHub repository, detected from the git remote at generate time so the
	// generated config passes `goreleaser check` without environment variables.
	GitHubOwner string
	GitHubRepo  string

	// Forgejo release target. ForgejoRelease switches the release block and
	// changelog source from github.com to the self-hosted Forgejo instance
	// (gitea_urls); the public GitHub copy is maintained by the release shadow.
	ForgejoRelease    bool
	GiteaAPIURL       string
	GiteaDownloadURL  string
	ChangelogUse      string
}

// GitHubActionsTemplateData is the strongly typed template data for the
// GitHub Actions workflow template.
type GitHubActionsTemplateData struct {
	ProjectName    string
	Triggers       []string
	DockerEnabled  bool
	SigningEnabled bool
	DockerRegistry string
	DockerImage    string
}

// Forgejo workflow generation pins. These are immutable flake refs verified
// during the Forgejo-primary migration (2026-10-28); bump via a deliberate
// upgrade pass, never a floating ref (reproducibility + runner bandwidth).
const (
	// BuildFlowNixPin is the pinned BuildFlow flake used as the CI gate.
	BuildFlowNixPin = "github:LarsArtmann/BuildFlow/49ced3e0039a5f6e88028697bd779ec2f7f28826"
	// GoReleaserNixPin is the pinned nixpkgs revision providing goreleaser
	// v2.17, the exact GoReleaser target of the generated configs.
	GoReleaserNixPin = "github:NixOS/nixpkgs/3ec636362e1816bc2f0888fa0fa2c69f1bc9e504#goreleaser"
	// ForgejoRunnerLabel is the act_runner label of the self-hosted runner.
	ForgejoRunnerLabel = "native"
)

// ForgejoWorkflowsTemplateData is the strongly typed template data for the
// Forgejo CI and release workflow templates (both rendered with [[ ]]
// delimiters, so ${{ secrets.* }} literals need no escaping).
type ForgejoWorkflowsTemplateData struct {
	ProjectName string
	RunnerLabel string
	BuildFlowPin string
	GoReleaserPin string

	// Forgejo instance and repository coordinates for the release shadow.
	ForgejoURL   string
	ForgejoOwner string
	GitHubOwner  string
	GitHubRepo   string
}

// NewForgejoWorkflowsTemplateData creates typed template data from
// SafeProjectConfig.
func NewForgejoWorkflowsTemplateData(config *domain.SafeProjectConfig) *ForgejoWorkflowsTemplateData {
	return &ForgejoWorkflowsTemplateData{
		ProjectName:   config.ProjectName,
		RunnerLabel:   ForgejoRunnerLabel,
		BuildFlowPin:  BuildFlowNixPin,
		GoReleaserPin: GoReleaserNixPin,
		ForgejoURL:    domain.ForgejoInstanceURL,
		ForgejoOwner:  domain.ForgejoPrimaryUser,
		GitHubOwner:   GetGitHubOwner(),
		GitHubRepo:    GetGitHubRepo(),
	}
}

// PlatformIgnored represents platform/architecture combinations to ignore.
type PlatformIgnored struct {
	GoOS   string
	GoArch string
}

// goOSDarwin is the Go OS used in the default ignore combination.
const goOSDarwin = "darwin"

// dockerPlatforms returns the linux platforms for dockers_v2 based on the
// configured architectures. amd64 is always included; arm64 only when the
// project actually builds it, so images are not built for dead targets.
func dockerPlatforms(config *domain.SafeProjectConfig) []string {
	platforms := []string{"linux/amd64"}
	if slices.Contains(config.Architectures, domain.ArchitectureARM64) {
		platforms = append(platforms, "linux/arm64")
	}

	return platforms
}

// cgoEnabledValue renders the CGO_ENABLED build environment value.
func cgoEnabledValue(status domain.CGOStatus) string {
	if status.IsEnabled() {
		return "1"
	}

	return "0"
}

// NewGoReleaserTemplateData creates typed template data from SafeProjectConfig.
func NewGoReleaserTemplateData(config *domain.SafeProjectConfig) *GoReleaserTemplateData {
	data := &GoReleaserTemplateData{
		ProjectName:     config.ProjectName,
		BinaryName:      config.BinaryName,
		MainPath:        config.MainPath,
		CGOEnabled:      cgoEnabledValue(config.CGOStatus),
		DockerPlatforms: dockerPlatforms(config),
		GitHubOwner:     GetGitHubOwner(),
		GitHubRepo:      GetGitHubRepo(),
		IgnoreCombinations: []PlatformIgnored{
			{GoOS: goOSDarwin, GoArch: "386"},
			{GoOS: "windows", GoArch: string(domain.ArchitectureARM64)},
		},
	}

	if len(config.Platforms) > 0 {
		platforms := make([]string, 0, len(config.Platforms))
		for _, platform := range config.Platforms {
			platforms = append(platforms, string(platform))
		}

		data.Platforms = platforms
	}

	if len(config.Architectures) > 0 {
		architectures := make([]string, 0, len(config.Architectures))
		for _, arch := range config.Architectures {
			architectures = append(architectures, string(arch))
		}

		data.Architectures = architectures
	}

	if len(config.BuildTags) > 0 {
		tags := make([]string, 0, len(config.BuildTags))
		for _, tag := range config.BuildTags {
			tags = append(tags, tag.String())
		}

		data.BuildTags = tags
	}

	if config.DockerSupport.IsEnabled() {
		data.DockerEnabled = true
		data.DockerRegistry = config.DockerRegistry.GetURL()
		data.DockerImage = config.GetDockerImageName()
	}

	data.ChangelogUse = "github"
	if config.CIProvider.ReleasesToForgejo() {
		data.ForgejoRelease = true
		data.GiteaAPIURL = domain.ForgejoAPIURL
		data.GiteaDownloadURL = domain.ForgejoInstanceURL
		data.ChangelogUse = "git"
	}

	return data
}

// NewGitHubActionsTemplateData creates typed template data from SafeProjectConfig.
func NewGitHubActionsTemplateData(config *domain.SafeProjectConfig) *GitHubActionsTemplateData {
	data := &GitHubActionsTemplateData{
		ProjectName:    config.ProjectName,
		DockerEnabled:  config.DockerSupport.IsEnabled(),
		SigningEnabled: config.SigningLevel.IsEnabled(),
	}

	if len(config.ActionsOn) > 0 {
		triggers := make([]string, 0, len(config.ActionsOn))
		for _, trigger := range config.ActionsOn {
			triggers = append(triggers, trigger.GitHubPattern())
		}

		data.Triggers = triggers
	}

	if config.DockerSupport.IsEnabled() {
		data.DockerRegistry = config.DockerRegistry.GetURL()
		data.DockerImage = config.GetDockerImageName()
	}

	return data
}

// JobExecutionResult represents the result of a job execution.
type JobExecutionResult struct {
	JobID    string      `json:"jobId"`
	JobName  string      `json:"jobName"`
	Status   JobStatus   `json:"status"`
	Duration string      `json:"duration"`
	Error    *JobError   `json:"error,omitempty"`
	Output   string      `json:"output,omitempty"`
	Metadata JobMetadata `json:"metadata"`
}

// JobStatus represents the status of a job.
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCancelled JobStatus = "cancelled"
)

// JobError represents a job execution error.
type JobError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Details   string `json:"details,omitempty"`
	Context   string `json:"context,omitempty"`
	Retryable bool   `json:"retryable"`
}

// JobMetadata represents job execution metadata.
type JobMetadata struct {
	StartedAt   string            `json:"startedAt"`
	CompletedAt string            `json:"completedAt,omitempty"`
	Retries     int               `json:"retries"`
	Tags        map[string]string `json:"tags,omitempty"`
}

// WorkflowState represents the state of a workflow.
type WorkflowState struct {
	WorkflowID   string                `json:"workflowId"`
	WorkflowName string                `json:"workflowName"`
	Status       WorkflowStatus        `json:"status"`
	CurrentStep  string                `json:"currentStep"`
	TotalSteps   int                   `json:"totalSteps"`
	Progress     float64               `json:"progress"`
	StartedAt    string                `json:"startedAt"`
	Results      []*JobExecutionResult `json:"results"`
	Metadata     WorkflowMetadata      `json:"metadata"`
}

// WorkflowStatus represents the status of a workflow.
type WorkflowStatus string

const (
	WorkflowStatusPending   WorkflowStatus = "pending"
	WorkflowStatusRunning   WorkflowStatus = "running"
	WorkflowStatusCompleted WorkflowStatus = "completed"
	WorkflowStatusFailed    WorkflowStatus = "failed"
	WorkflowStatusCancelled WorkflowStatus = "cancelled"
)

// WorkflowMetadata represents workflow metadata.
type WorkflowMetadata struct {
	CreatedBy   string            `json:"createdBy"`
	Timeout     string            `json:"timeout"`
	Parallel    bool              `json:"parallel"`
	Environment string            `json:"environment"`
	Tags        map[string]string `json:"tags,omitempty"`
}

// PlaceholderGitHubOwner and PlaceholderGitHubRepo are the fallback values used
// when the GitHub remote cannot be detected. They keep `goreleaser check`
// green without environment variables, but a release would target a
// nonexistent repository, so callers must warn when they leak into generated
// output.
const (
	PlaceholderGitHubOwner = "owner"
	PlaceholderGitHubRepo  = "repo"
)

// githubRemote is a resolved git origin remote.
type githubRemote struct {
	owner string
	repo  string
}

//nolint:gochecknoglobals // The origin remote cannot change during one run; cached after first lookup.
var cachedGitHubRemote *githubRemote

// parseGitHubRemote tries to get GitHub owner and repo from git remote.
// The result is cached for the lifetime of the process.
func parseGitHubRemote() (string, string) {
	if cachedGitHubRemote != nil {
		return cachedGitHubRemote.owner, cachedGitHubRemote.repo
	}

	ctx := context.Background()

	cmd := exec.CommandContext(ctx, "git", "remote", "get-url", "origin")
	if cmd != nil {
		output, err := cmd.Output()
		if err == nil {
			remote := strings.TrimSpace(string(output))
			owner, repo := git.ParseGitHubURL(remote)
			cachedGitHubRemote = &githubRemote{owner: owner, repo: repo}

			return owner, repo
		}
	}

	return PlaceholderGitHubOwner, PlaceholderGitHubRepo
}

// HasPlaceholderGitHubTarget reports whether GitHub owner or repository
// resolution fell back to placeholders: no git remote was detected and no
// explicit override was set via flags.
func HasPlaceholderGitHubTarget() bool {
	return GetGitHubOwner() == PlaceholderGitHubOwner || GetGitHubRepo() == PlaceholderGitHubRepo
}

var (
	//nolint:gochecknoglobals // Set once from CLI flags at startup, before any reader runs.
	githubOwnerOverride string
	//nolint:gochecknoglobals // Set once from CLI flags at startup, before any reader runs.
	githubRepoOverride string
)

// SetGitHubOwnerOverride sets an explicit GitHub owner, taking precedence
// over git remote detection. Used by the --github-owner flag.
func SetGitHubOwnerOverride(owner string) {
	githubOwnerOverride = owner
}

// SetGitHubRepoOverride sets an explicit GitHub repository, taking precedence
// over git remote detection. Used by the --github-repo flag.
func SetGitHubRepoOverride(repo string) {
	githubRepoOverride = repo
}

// GetGitHubOwner tries to get GitHub owner from git remote.
func GetGitHubOwner() string {
	if githubOwnerOverride != "" {
		return githubOwnerOverride
	}

	owner, _ := parseGitHubRemote()

	return owner
}

// GetGitHubRepo tries to get GitHub repo from git remote.
func GetGitHubRepo() string {
	if githubRepoOverride != "" {
		return githubRepoOverride
	}

	_, repo := parseGitHubRemote()

	return repo
}
