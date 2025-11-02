package main

import (
	"fmt"
	"strings"
)

// ProjectType represents the type of project being configured
type ProjectType string

const (
	ProjectTypeCLI      ProjectType = "CLI Application"
	ProjectTypeWeb      ProjectType = "Web Service"
	ProjectTypeLibrary  ProjectType = "Library"
	ProjectTypeAPI      ProjectType = "API Service"
	ProjectTypeDesktop  ProjectType = "Desktop Application"
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

func (pt ProjectType) Validate() error {
	if !pt.IsValid() {
		return fmt.Errorf("invalid project type: %s (must be one of: CLI Application, Web Service, Library, API Service, Desktop Application)", pt)
	}
	return nil
}

// GitProvider represents the git hosting provider
type GitProvider string

const (
	GitProviderGitHub GitProvider = "GitHub"
	GitProviderGitLab GitProvider = "GitLab"
	GitProviderBitbucket GitProvider = "Bitbucket"
	GitProviderGitea GitProvider = "Gitea"
	GitProviderSelfHosted GitProvider = "Self-hosted"
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

func (gp GitProvider) Validate() error {
	if !gp.IsValid() {
		return fmt.Errorf("invalid git provider: %s (must be one of: GitHub, GitLab, Bitbucket, Gitea, Self-hosted)", gp)
	}
	return nil
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

func (p Platform) Validate() error {
	if !p.IsValid() {
		return fmt.Errorf("invalid platform: %s", p)
	}
	return nil
}

// Architecture represents supported CPU architectures
type Architecture string

const (
	ArchAMD64 Architecture = "amd64"
	ArchARM64 Architecture = "arm64"
	Arch386   Architecture = "386"
	ArchARM   Architecture = "arm"
	ArchPPC64  Architecture = "ppc64"
	ArchPPC64LE Architecture = "ppc64le"
	ArchS390X  Architecture = "s390x"
	ArchMIPS   Architecture = "mips"
	ArchMIPSLE Architecture = "mipsle"
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

func (a Architecture) Validate() error {
	if !a.IsValid() {
		return fmt.Errorf("invalid architecture: %s", a)
	}
	return nil
}

// DockerRegistryType represents the type of Docker registry
type DockerRegistryType string

const (
	DockerRegistryDockerHub DockerRegistryType = "docker.io"
	DockerRegistryGitHub DockerRegistryType = "ghcr.io"
	DockerRegistryGitLab DockerRegistryType = "registry.gitlab.com"
	DockerRegistryQuay DockerRegistryType = "quay.io"
	DockerRegistryCustom DockerRegistryType = "custom"
)

func (drt DockerRegistryType) IsValid() bool {
	switch drt {
	case DockerRegistryDockerHub, DockerRegistryGitHub, DockerRegistryGitLab,
	     DockerRegistryQuay, DockerRegistryCustom:
		return true
	default:
		return false
	}
}

// ValidateDockerRegistry validates and normalizes Docker registry URLs
func ValidateDockerRegistry(registry string) (string, error) {
	if registry == "" {
		return "", fmt.Errorf("docker registry cannot be empty")
	}

	// Normalize common registry patterns
	registry = strings.TrimSpace(registry)
	
	// Add docker.io if just username
	if !strings.Contains(registry, ".") && !strings.Contains(registry, "/") {
		registry = "docker.io/" + registry
	}

	// Validate common registries
	parts := strings.Split(registry, "/")
	switch parts[0] {
	case "docker.io", "ghcr.io", "registry.gitlab.com", "quay.io":
		return registry, nil
	default:
		// Allow custom registries but warn
		return registry, nil
	}
}

// EnhancedProjectConfig represents the configuration with type safety
type EnhancedProjectConfig struct {
	// Basic Info
	ProjectName        string       `validate:"required,min=1,max=63"`
	ProjectDescription string       `validate:"max=255"`
	ProjectType        ProjectType  `validate:"required"`
	BinaryName         string       `validate:"required,min=1,max=63"`
	MainPath           string       `validate:"required"`

	// Build Options
	Platforms     []Platform    `validate:"required,min=1"`
	Architectures []Architecture `validate:"required,min=1"`
	CGOEnabled    bool
	BuildTags     []string
	LDFlags       bool

	// Release Options
	GitProvider    GitProvider `validate:"required"`
	DockerEnabled  bool
	DockerRegistry string       `validate:"required_if=DockerEnabled"`
	Signing        bool
	Homebrew       bool
	Snap           bool
	SBOM           bool

	// GitHub Actions
	GenerateActions bool
	ActionsOn       []string `validate:"required_if=GenerateActions"`

	// Advanced
	ProVersion bool
}

// ToProjectConfig converts EnhancedProjectConfig to legacy ProjectConfig
func (epc *EnhancedProjectConfig) ToProjectConfig() ProjectConfig {
	// Convert enums to strings
	platforms := make([]string, len(epc.Platforms))
	for i, p := range epc.Platforms {
		platforms[i] = string(p)
	}

	architectures := make([]string, len(epc.Architectures))
	for i, a := range epc.Architectures {
		architectures[i] = string(a)
	}

	return ProjectConfig{
		ProjectName:        epc.ProjectName,
		ProjectDescription: epc.ProjectDescription,
		ProjectType:        string(epc.ProjectType),
		BinaryName:         epc.BinaryName,
		MainPath:           epc.MainPath,
		Platforms:          platforms,
		Architectures:      architectures,
		CGOEnabled:         epc.CGOEnabled,
		BuildTags:          epc.BuildTags,
		LDFlags:            epc.LDFlags,
		GitProvider:         string(epc.GitProvider),
		DockerEnabled:       epc.DockerEnabled,
		DockerRegistry:      epc.DockerRegistry,
		Signing:            epc.Signing,
		Homebrew:           epc.Homebrew,
		Snap:               epc.Snap,
		SBOM:               epc.SBOM,
		GenerateActions:     epc.GenerateActions,
		ActionsOn:          epc.ActionsOn,
		ProVersion:         epc.ProVersion,
	}
}

// FromProjectConfig creates EnhancedProjectConfig from legacy ProjectConfig
func FromProjectConfig(pc ProjectConfig) (*EnhancedProjectConfig, error) {
	epc := &EnhancedProjectConfig{
		ProjectName:        pc.ProjectName,
		ProjectDescription: pc.ProjectDescription,
		BinaryName:         pc.BinaryName,
		MainPath:           pc.MainPath,
		CGOEnabled:         pc.CGOEnabled,
		BuildTags:          pc.BuildTags,
		LDFlags:            pc.LDFlags,
		DockerRegistry:      pc.DockerRegistry,
		Signing:            pc.Signing,
		Homebrew:           pc.Homebrew,
		Snap:               pc.Snap,
		SBOM:               pc.SBOM,
		GenerateActions:     pc.GenerateActions,
		ActionsOn:          pc.ActionsOn,
		ProVersion:         pc.ProVersion,
	}

	// Parse enums from strings
	epc.ProjectType = ProjectType(pc.ProjectType)
	if err := epc.ProjectType.Validate(); err != nil {
		return nil, fmt.Errorf("invalid project type: %w", err)
	}

	epc.GitProvider = GitProvider(pc.GitProvider)
	if err := epc.GitProvider.Validate(); err != nil {
		return nil, fmt.Errorf("invalid git provider: %w", err)
	}

	// Parse platforms
	epc.Platforms = make([]Platform, len(pc.Platforms))
	for i, p := range pc.Platforms {
		platform := Platform(p)
		if err := platform.Validate(); err != nil {
			return nil, fmt.Errorf("invalid platform %d: %w", i, err)
		}
		epc.Platforms[i] = platform
	}

	// Parse architectures
	epc.Architectures = make([]Architecture, len(pc.Architectures))
	for i, a := range pc.Architectures {
		arch := Architecture(a)
		if err := arch.Validate(); err != nil {
			return nil, fmt.Errorf("invalid architecture %d: %w", i, err)
		}
		epc.Architectures[i] = arch
	}

	return epc, nil
}

// Validate performs comprehensive validation of the configuration
func (epc *EnhancedProjectConfig) Validate() error {
	// Basic validation
	if epc.ProjectName == "" {
		return fmt.Errorf("project name is required")
	}
	
	if len(epc.ProjectName) > 63 {
		return fmt.Errorf("project name must be 63 characters or less")
	}

	if epc.BinaryName == "" {
		return fmt.Errorf("binary name is required")
	}

	if epc.MainPath == "" {
		return fmt.Errorf("main path is required")
	}

	// Validate enums
	if err := epc.ProjectType.Validate(); err != nil {
		return fmt.Errorf("project type validation: %w", err)
	}

	if err := epc.GitProvider.Validate(); err != nil {
		return fmt.Errorf("git provider validation: %w", err)
	}

	// Validate slices
	if len(epc.Platforms) == 0 {
		return fmt.Errorf("at least one platform is required")
	}

	for i, platform := range epc.Platforms {
		if err := platform.Validate(); err != nil {
			return fmt.Errorf("platform %d validation: %w", i, err)
		}
	}

	if len(epc.Architectures) == 0 {
		return fmt.Errorf("at least one architecture is required")
	}

	for i, arch := range epc.Architectures {
		if err := arch.Validate(); err != nil {
			return fmt.Errorf("architecture %d validation: %w", i, err)
		}
	}

	// Conditional validation
	if epc.DockerEnabled && epc.DockerRegistry == "" {
		return fmt.Errorf("docker registry is required when docker is enabled")
	}

	if epc.GenerateActions && len(epc.ActionsOn) == 0 {
		return fmt.Errorf("actions triggers are required when github actions is enabled")
	}

	return nil
}

// ApplyDefaults applies intelligent defaults based on project type
func (epc *EnhancedProjectConfig) ApplyDefaults() {
	// Set default project type if not specified
	if epc.ProjectType == "" {
		epc.ProjectType = ProjectTypeCLI
	}

	// Set default CGO based on project type
	if !epc.CGOEnabled && !epc.LDFlags {
		epc.CGOEnabled = epc.ProjectType.DefaultCGOEnabled()
	}

	// Set default platforms if not specified
	if len(epc.Platforms) == 0 {
		epc.Platforms = []Platform{PlatformLinux, PlatformDarwin, PlatformWindows}
	}

	// Set default architectures if not specified
	if len(epc.Architectures) == 0 {
		epc.Architectures = []Architecture{ArchAMD64, ArchARM64}
	}

	// Set default git provider if not specified
	if epc.GitProvider == "" {
		epc.GitProvider = GitProviderGitHub
	}

	// Set default docker registry if docker is enabled but registry not specified
	if epc.DockerEnabled && epc.DockerRegistry == "" {
		epc.DockerRegistry = "docker.io"
	}
}

// GetRecommendedPlatforms returns recommended platforms for the project type
func (pt ProjectType) GetRecommendedPlatforms() []Platform {
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

// GetRecommendedArchitectures returns recommended architectures for the project type
func (pt ProjectType) GetRecommendedArchitectures() []Architecture {
	switch pt {
	case ProjectTypeCLI:
		return []Architecture{ArchAMD64, ArchARM64}
	case ProjectTypeWeb, ProjectTypeAPI:
		return []Architecture{ArchAMD64, ArchARM64}
	case ProjectTypeLibrary:
		return []Architecture{ArchAMD64, ArchARM64, Arch386}
	case ProjectTypeDesktop:
		return []Architecture{ArchAMD64, ArchARM64, Arch386}
	default:
		return []Architecture{ArchAMD64, ArchARM64}
	}
}