package main

import (
	"fmt"
	"strings"
)

// TODO: This should be generated from TypeSpec!
// Eliminate split brains by making this the single source of truth.

// ProjectType represents the type of project being configured
type ProjectType string

const (
	ProjectTypeCLI     ProjectType = "cli" // Renamed for type safety
	ProjectTypeWeb     ProjectType = "web"
	ProjectTypeLibrary ProjectType = "library"
	ProjectTypeAPI     ProjectType = "api"
	ProjectTypeDesktop ProjectType = "desktop"
)

func (pt ProjectType) IsValid() bool {
	switch pt {
	case ProjectTypeCLI, ProjectTypeWeb, ProjectTypeLibrary,
		ProjectTypeAPI, ProjectTypeDesktop:
		return true
	default:
		return false
	}
}

func (pt ProjectType) String() string {
	// TODO: Should be generated from TypeSpec with human-readable names
	switch pt {
	case ProjectTypeCLI:
		return "CLI Application"
	case ProjectTypeWeb:
		return "Web Service"
	case ProjectTypeLibrary:
		return "Library"
	case ProjectTypeAPI:
		return "API Service"
	case ProjectTypeDesktop:
		return "Desktop Application"
	default:
		return string(pt)
	}
}

// TODO: Generate from TypeSpec with invariants!
// Docker should be invalid for Library types, etc.
func (pt ProjectType) DefaultCGOEnabled() bool {
	switch pt {
	case ProjectTypeCLI, ProjectTypeLibrary:
		return false
	case ProjectTypeWeb, ProjectTypeAPI, ProjectTypeDesktop:
		return true
	default:
		return false
	}
}

// TODO: Should be exhaustive pattern matching from TypeSpec
func (pt ProjectType) RecommendedPlatforms() []Platform {
	switch pt {
	case ProjectTypeCLI:
		return []Platform{PlatformLinux, PlatformDarwin, PlatformWindows}
	case ProjectTypeWeb, ProjectTypeAPI:
		return []Platform{PlatformLinux, PlatformDarwin}
	case ProjectTypeLibrary:
		return []Platform{PlatformLinux, PlatformDarwin, PlatformWindows}
	case ProjectTypeDesktop:
		return []Platform{PlatformWindows, PlatformDarwin, PlatformLinux}
	default:
		return []Platform{PlatformLinux, PlatformDarwin, PlatformWindows}
	}
}

// GitProvider represents the git hosting provider
type GitProvider string

const (
	GitProviderGitHub     GitProvider = "github"
	GitProviderGitLab     GitProvider = "gitlab"
	GitProviderBitbucket  GitProvider = "bitbucket"
	GitProviderGitea      GitProvider = "gitea"
	GitProviderSelfHosted GitProvider = "self-hosted"
)

func (gp GitProvider) IsValid() bool {
	switch gp {
	case GitProviderGitHub, GitProviderGitLab, GitProviderBitbucket,
		GitProviderGitea, GitProviderSelfHosted:
		return true
	default:
		return false
	}
}

func (gp GitProvider) String() string {
	switch gp {
	case GitProviderGitHub:
		return "GitHub"
	case GitProviderGitLab:
		return "GitLab"
	case GitProviderBitbucket:
		return "Bitbucket"
	case GitProviderGitea:
		return "Gitea"
	case GitProviderSelfHosted:
		return "Self-hosted"
	default:
		return string(gp)
	}
}

// Platform represents supported target platforms
type Platform string

const (
	PlatformLinux   Platform = "linux"
	PlatformDarwin  Platform = "darwin"
	PlatformWindows Platform = "windows"
	PlatformFreeBSD Platform = "freebsd"
	PlatformOpenBSD Platform = "openbsd"
	PlatformNetBSD  Platform = "netbsd"
)

func (p Platform) IsValid() bool {
	switch p {
	case PlatformLinux, PlatformDarwin, PlatformWindows,
		PlatformFreeBSD, PlatformOpenBSD, PlatformNetBSD:
		return true
	default:
		return false
	}
}

// Architecture represents supported CPU architectures
type Architecture string

const (
	ArchAMD64   Architecture = "amd64"
	ArchARM64   Architecture = "arm64"
	Arch386     Architecture = "386"
	ArchARM     Architecture = "arm"
	ArchPPC64   Architecture = "ppc64"
	ArchPPC64LE Architecture = "ppc64le"
	ArchS390X   Architecture = "s390x"
	ArchMIPS    Architecture = "mips"
	ArchMIPSLE  Architecture = "mipsle"
)

func (a Architecture) IsValid() bool {
	switch a {
	case ArchAMD64, ArchARM64, Arch386, ArchARM,
		ArchPPC64, ArchPPC64LE, ArchS390X, ArchMIPS, ArchMIPSLE:
		return true
	default:
		return false
	}
}

// ActionTrigger represents GitHub Actions triggers
type ActionTrigger string

const (
	ActionTriggerVersionTags ActionTrigger = "version-tags" // v*, v1.0.0, etc.
	ActionTriggerAllTags     ActionTrigger = "all-tags"     // *
	ActionTriggerManual      ActionTrigger = "manual"       // workflow_dispatch
	ActionTriggerMain        ActionTrigger = "main"         // push to main
	ActionTriggerRelease     ActionTrigger = "release"      // published release
)

func (at ActionTrigger) IsValid() bool {
	switch at {
	case ActionTriggerVersionTags, ActionTriggerAllTags, ActionTriggerManual,
		ActionTriggerMain, ActionTriggerRelease:
		return true
	default:
		return false
	}
}

func (at ActionTrigger) ToGitHubActions() []string {
	// TODO: Should be generated from TypeSpec
	switch at {
	case ActionTriggerVersionTags:
		return []string{"push:\n  tags:\n    - 'v*'"}
	case ActionTriggerAllTags:
		return []string{"push:\n  tags:\n    - '*'"}
	case ActionTriggerManual:
		return []string{"workflow_dispatch:"}
	case ActionTriggerMain:
		return []string{"push:\n  branches:\n    - main"}
	case ActionTriggerRelease:
		return []string{"release:\n    types: [published]"}
	default:
		return []string{}
	}
}

// DockerRegistry represents Docker registry types
type DockerRegistry string

const (
	DockerRegistryDockerHub DockerRegistry = "docker.io"
	DockerRegistryGitHub    DockerRegistry = "ghcr.io"
	DockerRegistryGitLab    DockerRegistry = "registry.gitlab.com"
	DockerRegistryQuay      DockerRegistry = "quay.io"
	DockerRegistryCustom    DockerRegistry = "custom"
)

func (dr DockerRegistry) IsValid() bool {
	switch dr {
	case DockerRegistryDockerHub, DockerRegistryGitHub, DockerRegistryGitLab,
		DockerRegistryQuay, DockerRegistryCustom:
		return true
	default:
		return false
	}
}

// TODO: Generate from TypeSpec with validation!
func (dr DockerRegistry) ValidateRegistryURL(url string) error {
	if url == "" {
		return fmt.Errorf("docker registry URL cannot be empty")
	}

	url = strings.TrimSpace(url)

	// TODO: Should be type-safe validation from TypeSpec
	switch dr {
	case DockerRegistryDockerHub:
		if !strings.Contains(url, "docker.io") && !strings.Contains(url, "/") {
			return fmt.Errorf("Docker Hub registry should include docker.io or be username")
		}
	case DockerRegistryGitHub:
		if !strings.Contains(url, "ghcr.io") {
			return fmt.Errorf("GitHub Container Registry should include ghcr.io")
		}
		// TODO: Add validation for all registry types
	}

	return nil
}

// ConfigState represents the validation state of configuration
type ConfigState string

const (
	ConfigStateDraft      ConfigState = "draft" // TODO: Should be generated
	ConfigStateValid      ConfigState = "valid"
	ConfigStateInvalid    ConfigState = "invalid"
	ConfigStateProcessing ConfigState = "processing"
)

func (cs ConfigState) IsValid() bool {
	switch cs {
	case ConfigStateDraft, ConfigStateValid, ConfigStateInvalid, ConfigStateProcessing:
		return true
	default:
		return false
	}
}

// TODO: CRITICAL: Single source of truth configuration
// This eliminates split brains and provides type safety

// TODO: Should be generated from TypeSpec!
// All fields should be typed, no string[] anywhere!
type SafeProjectConfig struct {
	// Basic Information
	ProjectName        string      `json:"project_name" yaml:"project_name" validate:"required,min=1,max=63"`
	ProjectDescription string      `json:"project_description" yaml:"project_description" validate:"max=255"`
	ProjectType        ProjectType `json:"project_type" yaml:"project_type" validate:"required"`
	BinaryName         string      `json:"binary_name" yaml:"binary_name" validate:"required,min=1,max=63"`
	MainPath           string      `json:"main_path" yaml:"main_path" validate:"required"`

	// Build Configuration
	Platforms     []Platform     `json:"platforms" yaml:"platforms" validate:"required,min=1"`
	Architectures []Architecture `json:"architectures" yaml:"architectures" validate:"required,min=1"`
	CGOEnabled    bool           `json:"cgo_enabled" yaml:"cgo_enabled"`
	BuildTags     []string       `json:"build_tags" yaml:"build_tags"`
	LDFlags       bool           `json:"ldflags" yaml:"ldflags"`

	// Release Configuration
	GitProvider    GitProvider    `json:"git_provider" yaml:"git_provider" validate:"required"`
	DockerEnabled  bool           `json:"docker_enabled" yaml:"docker_enabled"`
	DockerRegistry DockerRegistry `json:"docker_registry" yaml:"docker_registry" validate:"required_if=DockerEnabled"`
	DockerImage    string         `json:"docker_image" yaml:"docker_image"` // TODO: Should be typed
	Signing        bool           `json:"signing" yaml:"signing"`
	Homebrew       bool           `json:"homebrew" yaml:"homebrew"`
	Snap           bool           `json:"snap" yaml:"snap"`
	SBOM           bool           `json:"sbom" yaml:"sbom"`

	// CI/CD Configuration
	GenerateActions bool            `json:"generate_actions" yaml:"generate_actions"`
	ActionsOn       []ActionTrigger `json:"actions_on" yaml:"actions_on" validate:"required_if=GenerateActions"`

	// Advanced Features
	ProVersion bool        `json:"pro_version" yaml:"pro_version"`
	State      ConfigState `json:"state" yaml:"state"`
}

// TODO: Should be generated from TypeSpec with invariants!
func (spc *SafeProjectConfig) ValidateInvariants() error {
	// TODO: Should be compile-time generated validation

	// Invariant: Required fields must not be empty
	if spc.ProjectName == "" {
		return fmt.Errorf("project name is required")
	}

	if spc.BinaryName == "" {
		return fmt.Errorf("binary name is required")
	}

	if spc.MainPath == "" {
		return fmt.Errorf("main path is required")
	}

	// Invariant: Project type must be valid
	if !spc.ProjectType.IsValid() {
		return fmt.Errorf("invalid project type: %s", spc.ProjectType)
	}

	// Invariant: Platforms must not be empty
	if len(spc.Platforms) == 0 {
		return fmt.Errorf("at least one platform is required")
	}

	// Validate each platform
	for _, platform := range spc.Platforms {
		if !platform.IsValid() {
			return fmt.Errorf("invalid platform: %s", platform)
		}
	}

	// Invariant: Architectures must not be empty
	if len(spc.Architectures) == 0 {
		return fmt.Errorf("at least one architecture is required")
	}

	// Validate each architecture
	for _, arch := range spc.Architectures {
		if !arch.IsValid() {
			return fmt.Errorf("invalid architecture: %s", arch)
		}
	}

	// Invariant: Git provider must be valid
	if !spc.GitProvider.IsValid() {
		return fmt.Errorf("invalid git provider: %s", spc.GitProvider)
	}

	// Invariant: Docker enabled requires valid registry
	if spc.DockerEnabled && !spc.DockerRegistry.IsValid() {
		return fmt.Errorf("Docker enabled requires valid registry")
	}

	// Invariant: Actions enabled requires triggers
	if spc.GenerateActions && len(spc.ActionsOn) == 0 {
		return fmt.Errorf("GitHub Actions enabled requires at least one trigger")
	}

	// Validate action triggers
	for _, trigger := range spc.ActionsOn {
		if !trigger.IsValid() {
			return fmt.Errorf("invalid action trigger: %s", trigger)
		}
	}

	// TODO: Add more invariants for all type combinations

	return nil
}

// TODO: Should be generated from TypeSpec!
func (spc *SafeProjectConfig) ApplyDefaults() {
	// TODO: Should be intelligent defaults from TypeSpec
	if spc.ProjectType == "" {
		spc.ProjectType = ProjectTypeCLI
	}

	// Apply CGO defaults based on project type
	if spc.ProjectType.DefaultCGOEnabled() && !spc.CGOEnabled {
		spc.CGOEnabled = spc.ProjectType.DefaultCGOEnabled()
	}

	if len(spc.Platforms) == 0 {
		spc.Platforms = spc.ProjectType.RecommendedPlatforms()
	}

	if len(spc.Architectures) == 0 {
		spc.Architectures = []Architecture{ArchAMD64, ArchARM64}
	}

	if spc.GitProvider == "" {
		spc.GitProvider = GitProviderGitHub
	}

	if spc.DockerEnabled && spc.DockerRegistry == "" {
		spc.DockerRegistry = DockerRegistryDockerHub
	}

	if spc.State == "" {
		spc.State = ConfigStateDraft
	}
}

// TODO: Should be generated from TypeSpec!
// Converts legacy ProjectConfig to SafeProjectConfig
func (spc *SafeProjectConfig) FromLegacy(legacy ProjectConfig) error {
	spc.ProjectName = legacy.ProjectName
	spc.ProjectDescription = legacy.ProjectDescription

	// TODO: This conversion is problematic - eliminates type safety
	switch legacy.ProjectType {
	case "CLI Application":
		spc.ProjectType = ProjectTypeCLI
	case "Web Service":
		spc.ProjectType = ProjectTypeWeb
	case "Library":
		spc.ProjectType = ProjectTypeLibrary
	case "API Service":
		spc.ProjectType = ProjectTypeAPI
	case "Desktop Application":
		spc.ProjectType = ProjectTypeDesktop
	default:
		return fmt.Errorf("invalid project type: %s", legacy.ProjectType)
	}

	spc.BinaryName = legacy.BinaryName
	spc.MainPath = legacy.MainPath

	// TODO: This is split brain! Should be type-safe!
	spc.Platforms = make([]Platform, len(legacy.Platforms))
	for i, p := range legacy.Platforms {
		spc.Platforms[i] = Platform(p)
	}

	spc.Architectures = make([]Architecture, len(legacy.Architectures))
	for i, a := range legacy.Architectures {
		spc.Architectures[i] = Architecture(a)
	}

	spc.CGOEnabled = legacy.CGOEnabled
	spc.BuildTags = legacy.BuildTags
	spc.LDFlags = legacy.LDFlags

	switch legacy.GitProvider {
	case "GitHub":
		spc.GitProvider = GitProviderGitHub
	case "GitLab":
		spc.GitProvider = GitProviderGitLab
	case "Bitbucket":
		spc.GitProvider = GitProviderBitbucket
	case "Gitea":
		spc.GitProvider = GitProviderGitea
	case "Self-hosted":
		spc.GitProvider = GitProviderSelfHosted
	default:
		return fmt.Errorf("invalid git provider: %s", legacy.GitProvider)
	}

	spc.DockerEnabled = legacy.DockerEnabled
	spc.DockerImage = legacy.DockerRegistry   // TODO: Fix mapping
	spc.DockerRegistry = DockerRegistryGitHub // Default

	// TODO: More split brain cleanup needed
	spc.GenerateActions = legacy.GenerateActions

	spc.ActionsOn = make([]ActionTrigger, len(legacy.ActionsOn))
	for i, a := range legacy.ActionsOn {
		// TODO: This mapping is error-prone
		if strings.Contains(a, "version tags") {
			spc.ActionsOn[i] = ActionTriggerVersionTags
		} else if strings.Contains(a, "all tags") {
			spc.ActionsOn[i] = ActionTriggerAllTags
		} else if strings.Contains(a, "manual") {
			spc.ActionsOn[i] = ActionTriggerManual
		} else {
			spc.ActionsOn[i] = ActionTriggerManual // Default
		}
	}

	spc.ProVersion = legacy.ProVersion
	spc.State = ConfigStateDraft

	return nil
}
