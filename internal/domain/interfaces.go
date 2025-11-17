package domain

import (
	"context"
	"io"
	"os"
)

// Repository interfaces for external dependency management
// Follows Clean Architecture principles - domain defines contracts

// FileSystemRepository handles all file system operations
type FileSystemRepository interface {
	// File operations
	ReadFile(ctx context.Context, path string) ([]byte, error)
	WriteFile(ctx context.Context, path string, data []byte, perm os.FileMode) error
	CreateFile(ctx context.Context, path string) (io.WriteCloser, error)
	DeleteFile(ctx context.Context, path string) error
	FileExists(ctx context.Context, path string) (bool, error)
	
	// Directory operations
	CreateDir(ctx context.Context, path string, perm os.FileMode) error
	CreateDirAll(ctx context.Context, path string, perm os.FileMode) error
	DirExists(ctx context.Context, path string) (bool, error)
	ReadDir(ctx context.Context, path string) ([]os.DirEntry, error)
	
	// Permission and metadata operations
	GetFileInfo(ctx context.Context, path string) (os.FileInfo, error)
	CheckPermissions(ctx context.Context, path string) (bool, error)
	
	// Utility operations
	AbsPath(path string) (string, error)
	RelPath(base, target string) (string, error)
	CleanPath(path string) string
	JoinPath(elem ...string) string
	TempDir(dir, pattern string) (string, error)
}

// TemplateRepository handles template rendering and management
type TemplateRepository interface {
	// Template operations
	LoadTemplate(ctx context.Context, name string) (string, error)
	RenderTemplate(ctx context.Context, templateContent string, data interface{}) (string, error)
	ValidateTemplate(ctx context.Context, content string) error
	
	// Template discovery and management
	ListTemplates(ctx context.Context) ([]string, error)
	TemplateExists(ctx context.Context, name string) (bool, error)
	GetTemplatePath(ctx context.Context, name string) (string, error)
	
	// Template categories
	GetTemplatesForType(ctx context.Context, projectType ProjectType) ([]string, error)
	GetTemplateMetadata(ctx context.Context, name string) (TemplateMetadata, error)
}

// GoReleaserRepository handles GoReleaser integration
type GoReleaserRepository interface {
	// Configuration management
	ValidateConfig(ctx context.Context, config *SafeProjectConfig) error
	CheckConfig(ctx context.Context, configPath string) error
	GenerateConfig(ctx context.Context, config *SafeProjectConfig) (string, error)
	
	// Build operations
	BuildSnapshot(ctx context.Context, config *SafeProjectConfig) error
	BuildRelease(ctx context.Context, config *SafeProjectConfig, version string) error
	
	// Release operations
	ReleaseDryRun(ctx context.Context, config *SafeProjectConfig) error
	Release(ctx context.Context, config *SafeProjectConfig, version string) error
	
	// Utility operations
	GetSupportedPlatforms() ([]Platform, error)
	GetSupportedArchitectures() ([]Architecture, error)
	GetVersion() (string, error)
	CheckInstallation() (bool, error)
}

// GitHubRepository handles GitHub API operations
type GitHubRepository interface {
	// Repository operations
	GetRepo(ctx context.Context, owner, name string) (*GitHubRepo, error)
	CreateRepo(ctx context.Context, repo *GitHubRepoRequest) (*GitHubRepo, error)
	UpdateRepo(ctx context.Context, owner, name string, updates *GitHubRepoUpdate) (*GitHubRepo, error)
	
	// Branch and tag operations
	GetBranches(ctx context.Context, owner, name string) ([]*GitHubBranch, error)
	GetTags(ctx context.Context, owner, name string) ([]*GitHubTag, error)
	CreateTag(ctx context.Context, owner, name string, tag *GitHubTagRequest) (*GitHubTag, error)
	
	// Release operations
	GetReleases(ctx context.Context, owner, name string) ([]*GitHubRelease, error)
	CreateRelease(ctx context.Context, owner, name string, release *GitHubReleaseRequest) (*GitHubRelease, error)
	UploadReleaseAsset(ctx context.Context, owner, name, releaseID string, asset *GitHubAsset) (*GitHubAsset, error)
	
	// Workflow operations
	GetWorkflows(ctx context.Context, owner, name string) ([]*GitHubWorkflow, error)
	TriggerWorkflow(ctx context.Context, owner, name, workflowID string, inputs map[string]interface{}) error
	
	// Authentication and permissions
	ValidateToken(ctx context.Context, token string) error
	GetUser(ctx context.Context, token string) (*GitHubUser, error)
	CheckPermissions(ctx context.Context, owner, name string, token string) (*GitHubPermissions, error)
}

// DockerRepository handles Docker registry operations
type DockerRepository interface {
	// Registry operations
	ValidateRegistry(ctx context.Context, registry DockerRegistry, url string) error
	GetRegistryURL(ctx context.Context, registry DockerRegistry) (string, error)
	
	// Image operations
	BuildImage(ctx context.Context, dockerfile string, tag string, config *DockerBuildConfig) error
	PushImage(ctx context.Context, image, registry DockerRegistry, credentials *DockerCredentials) error
	PullImage(ctx context.Context, image string) error
	
	// Image metadata
	GetImageInfo(ctx context.Context, image string) (*DockerImage, error)
	TagImage(ctx context.Context, source, target string) error
	
	// Authentication and security
	Login(ctx context.Context, registry DockerRegistry, credentials *DockerCredentials) error
	Logout(ctx context.Context, registry DockerRegistry) error
	ValidateCredentials(ctx context.Context, registry DockerRegistry, credentials *DockerCredentials) error
}

// Logger interface for dependency injection
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	
	// Context logging
	DebugContext(ctx context.Context, msg string, args ...interface{})
	InfoContext(ctx context.Context, msg string, args ...interface{})
	WarnContext(ctx context.Context, msg string, args ...interface{})
	ErrorContext(ctx context.Context, msg string, args ...interface{})
	
	// Structured logging
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithError(err error) Logger
}

// Validator interface for configuration validation
type Validator interface {
	Validate(ctx context.Context, config *SafeProjectConfig) error
	ValidateField(ctx context.Context, field, value string) error
	GetValidationRules(ctx context.Context) (*ValidationRules, error)
}

// Use case interfaces following Clean Architecture

// ConfigUseCase handles configuration management
type ConfigUseCase interface {
	CreateConfig(ctx context.Context, projectType ProjectType) (*SafeProjectConfig, error)
	LoadConfig(ctx context.Context, path string) (*SafeProjectConfig, error)
	SaveConfig(ctx context.Context, config *SafeProjectConfig, path string) error
	ValidateConfig(ctx context.Context, config *SafeProjectConfig) (*ValidationResult, error)
	UpdateConfig(ctx context.Context, config *SafeProjectConfig, updates *ConfigUpdate) (*SafeProjectConfig, error)
}

// GenerationUseCase handles template and configuration generation
type GenerationUseCase interface {
	GenerateGoReleaserConfig(ctx context.Context, config *SafeProjectConfig) (string, error)
	GenerateGitHubActions(ctx context.Context, config *SafeProjectConfig) (string, error)
	GenerateDockerfile(ctx context.Context, config *SafeProjectConfig) (string, error)
	GenerateAll(ctx context.Context, config *SafeProjectConfig, outputPath string) error
}

// ProjectUseCase handles project-level operations
type ProjectUseCase interface {
	InitializeProject(ctx context.Context, config *SafeProjectConfig) error
	ValidateProject(ctx context.Context, projectPath string) (*ProjectValidationResult, error)
	GetProjectInfo(ctx context.Context, projectPath string) (*ProjectInfo, error)
	UpgradeProject(ctx context.Context, projectPath string, config *SafeProjectConfig) error
}

// Template metadata and configuration types

type TemplateMetadata struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	ProjectType ProjectType        `json:"project_type"`
	Category    TemplateCategory    `json:"category"`
	Author      string             `json:"author"`
	Version     string             `json:"version"`
	License     string             `json:"license"`
	Tags        []string           `json:"tags"`
	Variables   []TemplateVariable `json:"variables"`
}

type TemplateCategory string

const (
	TemplateCategoryConfig        TemplateCategory = "config"
	TemplateCategoryWorkflow     TemplateCategory = "workflow"
	TemplateCategoryDockerfile  TemplateCategory = "dockerfile"
	TemplateCategoryTemplate     TemplateCategory = "template"
	TemplateCategoryExample      TemplateCategory = "example"
)

type TemplateVariable struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Description  string      `json:"description"`
	Required     bool        `json:"required"`
	DefaultValue interface{} `json:"default_value,omitempty"`
	Options      []string    `json:"options,omitempty"`
}

type ValidationRules struct {
	ProjectName        FieldRules `json:"project_name"`
	BinaryName         FieldRules `json:"binary_name"`
	MainPath           FieldRules `json:"main_path"`
	ProjectDescription  FieldRules `json:"project_description"`
	DockerRegistry     FieldRules `json:"docker_registry"`
	DockerImage        FieldRules `json:"docker_image"`
}

type FieldRules struct {
	Required     bool     `json:"required"`
	MinLength    int      `json:"min_length"`
	MaxLength    int      `json:"max_length"`
	Pattern      string   `json:"pattern"`
	Options      []string `json:"options"`
	Description  string   `json:"description"`
}

type ValidationResult struct {
	IsValid  bool             `json:"is_valid"`
	Errors   []*DomainError    `json:"errors"`
	Warnings []*DomainError    `json:"warnings"`
	Rules    *ValidationRules  `json:"rules"`
}

type ConfigUpdate struct {
	ProjectName        *string        `json:"project_name,omitempty"`
	ProjectDescription *string        `json:"project_description,omitempty"`
	ProjectType       *ProjectType   `json:"project_type,omitempty"`
	BinaryName        *string        `json:"binary_name,omitempty"`
	MainPath          *string        `json:"main_path,omitempty"`
	Platforms         []Platform      `json:"platforms,omitempty"`
	Architectures     []Architecture `json:"architectures,omitempty"`
	CGOEnabled        *bool          `json:"cgo_enabled,omitempty"`
	BuildTags         []BuildTag     `json:"build_tags,omitempty"`
	LDFlags           *bool          `json:"ldflags,omitempty"`
	GitProvider       *GitProvider   `json:"git_provider,omitempty"`
	DockerEnabled     *bool          `json:"docker_enabled,omitempty"`
	DockerRegistry    *DockerRegistry `json:"docker_registry,omitempty"`
	DockerImage       *string        `json:"docker_image,omitempty"`
	Signing           *bool          `json:"signing,omitempty"`
	Homebrew          *bool          `json:"homebrew,omitempty"`
	Snap              *bool          `json:"snap,omitempty"`
	SBOM              *bool          `json:"sbom,omitempty"`
	GenerateActions   *bool          `json:"generate_actions,omitempty"`
	ActionsOn         []ActionTrigger `json:"actions_on,omitempty"`
	ProVersion        *bool          `json:"pro_version,omitempty"`
	State             *ConfigState   `json:"state,omitempty"`
}

// Project information and validation types
type ProjectInfo struct {
	Name         string           `json:"name"`
	Path         string           `json:"path"`
	HasGoMod     bool             `json:"has_go_mod"`
	HasMainFile  bool             `json:"has_main_file"`
	MainFilePath  string           `json:"main_file_path"`
	ProjectType  ProjectType      `json:"project_type"`
	BinaryName   string           `json:"binary_name"`
	Buildable    bool             `json:"buildable"`
	Dependencies  []string         `json:"dependencies"`
	GoVersion    string           `json:"go_version"`
	Modules      []string         `json:"modules"`
}

type ProjectValidationResult struct {
	IsValid        bool                     `json:"is_valid"`
	Info           *ProjectInfo              `json:"info,omitempty"`
	Issues         []*DomainError           `json:"issues,omitempty"`
	Recommendations []string                  `json:"recommendations,omitempty"`
	Warnings       []*DomainError           `json:"warnings,omitempty"`
}

// External service data types

type GitHubRepo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	Fork        bool   `json:"fork"`
	HTMLURL     string `json:"html_url"`
	CloneURL    string `json:"clone_url"`
	DefaultBranch string `json:"default_branch"`
	Language    string `json:"language"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type GitHubRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	AutoInit    bool   `json:"auto_init"`
	GitignoreTemplate string `json:"gitignore_template"`
	LicenseTemplate string `json:"license_template"`
}

type GitHubRepoUpdate struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Private     *bool  `json:"private,omitempty"`
}

type GitHubBranch struct {
	Name   string `json:"name"`
	SHA    string `json:"commit.sha"`
	URL    string `json:"url"`
	Protected bool `json:"protected"`
}

type GitHubTag struct {
	Name      string `json:"name"`
	SHA       string `json:"commit.sha"`
	URL       string `json:"url"`
	ZipballURL string `json:"zipball_url"`
	TarballURL string `json:"tarball_url"`
}

type GitHubTagRequest struct {
	Tag         string `json:"tag"`
	Message     string `json:"message"`
	ObjectType  string `json:"object"`
	ObjectSHA   string `json:"sha"`
	Draft       bool   `json:"draft"`
	Prerelease  bool   `json:"prerelease"`
}

type GitHubRelease struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	TagName     string          `json:"tag_name"`
	Draft       bool            `json:"draft"`
	Prerelease  bool            `json:"prerelease"`
	PublishedAt string          `json:"published_at"`
	Assets      []*GitHubAsset   `json:"assets"`
	HTMLURL     string          `json:"html_url"`
	URL         string          `json:"url"`
}

type GitHubReleaseRequest struct {
	Name            string              `json:"name"`
	TagName         string              `json:"tag_name"`
	Body            string              `json:"body"`
	Draft           bool                `json:"draft"`
	Prerelease      bool                `json:"prerelease"`
	GenerateReleaseNotes bool             `json:"generate_release_notes"`
	TargetCommitish  string              `json:"target_commitish,omitempty"`
}

type GitHubAsset struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	ContentType   string `json:"content_type"`
	Size          int64  `json:"size"`
	DownloadCount int    `json:"download_count"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	BrowserURL    string `json:"browser_download_url"`
	URL           string `json:"url"`
}

type GitHubWorkflow struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	State     string `json:"state"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	URL       string `json:"html_url"`
	BadgeURL  string `json:"badge_url"`
}

type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
}

type GitHubPermissions struct {
	Admin    bool `json:"admin"`
	Push     bool `json:"push"`
	Pull     bool `json:"pull"`
	Maintain  bool `json:"maintain"`
	Triage   bool `json:"triage"`
}

// Docker configuration and build types
type DockerBuildConfig struct {
	Context       string            `json:"context"`
	Dockerfile    string            `json:"dockerfile"`
	BuildArgs     map[string]string `json:"build_args"`
	Labels        map[string]string `json:"labels"`
	Tags          []string          `json:"tags"`
	NoCache       bool              `json:"no_cache"`
	Target        string            `json:"target"`
	Platform      string            `json:"platform"`
}

type DockerCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Registry string `json:"registry"`
}

type DockerImage struct {
	ID           string            `json:"id"`
	RepoTags     []string          `json:"repo_tags"`
	Size         int64             `json:"size"`
	Created      int64             `json:"created"`
	Labels       map[string]string `json:"labels"`
	Architecture string            `json:"architecture"`
	OS           string            `json:"os"`
}

// Utility interfaces

// ProgressTracker for long-running operations
type ProgressTracker interface {
	Start(total int64, message string)
	Update(current int64, message string)
	Finish(message string)
	WithError(err error)
	WithLogger(logger Logger) ProgressTracker
}

// Configuration store interface
type ConfigurationStore interface {
	Save(ctx context.Context, key string, value interface{}) error
	Load(ctx context.Context, key string, value interface{}) (bool, error)
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, pattern string) (map[string]interface{}, error)
}