package domain

import (
	"fmt"
	"strings"
)

// SafeProjectConfig represents the single source of truth for project configuration
// Generated from TypeSpec specification - DO NOT MODIFY MANUALLY
type SafeProjectConfig struct {
	// Basic Information
	ProjectName        string      `json:"project_name" yaml:"project_name"`
	ProjectDescription string      `json:"project_description,omitempty" yaml:"project_description,omitempty"`
	ProjectType        ProjectType `json:"project_type" yaml:"project_type"`
	BinaryName         string      `json:"binary_name" yaml:"binary_name"`
	MainPath           string      `json:"main_path" yaml:"main_path"`

	// Build Configuration
	Platforms     []Platform     `json:"platforms" yaml:"platforms"`
	Architectures []Architecture `json:"architectures" yaml:"architectures"`
	CGOEnabled    bool           `json:"cgo_enabled" yaml:"cgo_enabled"`
	BuildTags     []BuildTag     `json:"build_tags,omitempty" yaml:"build_tags,omitempty"`
	LDFlags       bool           `json:"ldflags" yaml:"ldflags"`

	// Release Configuration
	GitProvider    GitProvider    `json:"git_provider" yaml:"git_provider"`
	DockerEnabled  bool           `json:"docker_enabled" yaml:"docker_enabled"`
	DockerRegistry DockerRegistry `json:"docker_registry" yaml:"docker_registry"`
	DockerImage    string         `json:"docker_image,omitempty" yaml:"docker_image,omitempty"`
	Signing        bool           `json:"signing" yaml:"signing"`
	Homebrew       bool           `json:"homebrew" yaml:"homebrew"`
	Snap           bool           `json:"snap" yaml:"snap"`
	SBOM           bool           `json:"sbom" yaml:"sbom"`

	// CI/CD Configuration
	GenerateActions bool            `json:"generate_actions" yaml:"generate_actions"`
	ActionsOn       []ActionTrigger `json:"actions_on" yaml:"actions_on"`

	// Advanced Features
	ProVersion bool        `json:"pro_version" yaml:"pro_version"`
	State      ConfigState `json:"state" yaml:"state"`
}

// NewSafeProjectConfig creates a new configuration with safe defaults
func NewSafeProjectConfig() *SafeProjectConfig {
	config := &SafeProjectConfig{
		// Set initial state
		State: ConfigStateDraft,
		
		// Apply smart defaults
		ProjectType:   ProjectTypeCLI,
		GitProvider:   GitProviderGitHub,
		DockerRegistry: DockerRegistryDockerHub,
		LDFlags:       true,
		SBOM:          true,
		GenerateActions: true,
		ActionsOn:      []ActionTrigger{ActionTriggerVersionTags},
	}
	
	// Apply project type specific defaults
	config.ApplyDefaults()
	
	return config
}

// ApplyDefaults applies intelligent defaults based on project type and context
func (spc *SafeProjectConfig) ApplyDefaults() {
	// Apply CGO defaults based on project type
	if spc.ProjectType.DefaultCGOEnabled() {
		spc.CGOEnabled = spc.ProjectType.DefaultCGOEnabled()
	}
	
	// Set recommended platforms if not specified
	if len(spc.Platforms) == 0 {
		spc.Platforms = spc.ProjectType.RecommendedPlatforms()
	}
	
	// Set recommended architectures if not specified
	if len(spc.Architectures) == 0 {
		spc.Architectures = []Architecture{ArchitectureAMD64, ArchitectureARM64}
	}
	
	// Set default Docker registry if Docker is enabled
	if spc.DockerEnabled && spc.DockerRegistry == "" {
		spc.DockerRegistry = DockerRegistryDockerHub
	}
	
	// Set default binary name if not specified
	if spc.BinaryName == "" {
		spc.BinaryName = spc.ProjectType.DefaultBinaryName()
	}
	
	// Set default main path if not specified
	if spc.MainPath == "" && spc.ProjectType.RequiresMainPath() {
		spc.MainPath = "./cmd/" + spc.ProjectType.DefaultBinaryName()
	}
	
	// Set default image name if not specified
	if spc.DockerImage == "" {
		spc.DockerImage = strings.ToLower(spc.ProjectName)
	}
}

// ValidateInvariants enforces compile-time invariants defined in TypeSpec
func (spc *SafeProjectConfig) ValidateInvariants() error {
	// Basic Information Validation
	if err := ValidateProjectName(spc.ProjectName); err != nil {
		return fmt.Errorf("project name validation failed: %w", err)
	}
	
	if err := ValidateBinaryName(spc.BinaryName); err != nil {
		return fmt.Errorf("binary name validation failed: %w", err)
	}
	
	if err := ValidateProjectDescription(spc.ProjectDescription); err != nil {
		return fmt.Errorf("project description validation failed: %w", err)
	}
	
	if err := ValidateMainPath(spc.MainPath); err != nil {
		return fmt.Errorf("main path validation failed: %w", err)
	}
	
	if err := ValidateGitProvider(spc.GitProvider); err != nil {
		return fmt.Errorf("git provider validation failed: %w", err)
	}
	
	// Type Validation
	if !spc.ProjectType.IsValid() {
		return fmt.Errorf("invalid project type: %s", spc.ProjectType)
	}
	
	if err := ValidatePlatforms(spc.Platforms); err != nil {
		return fmt.Errorf("platforms validation failed: %w", err)
	}
	
	if err := ValidateArchitectures(spc.Architectures); err != nil {
		return fmt.Errorf("architectures validation failed: %w", err)
	}
	
	if err := ValidatePlatformArchCompatibility(spc.Platforms, spc.Architectures); err != nil {
		return fmt.Errorf("platform-architecture compatibility failed: %w", err)
	}
	
	if err := ValidateActionTriggers(spc.ActionsOn); err != nil {
		return fmt.Errorf("action triggers validation failed: %w", err)
	}
	
	// Registry Validation
	if err := ValidateDockerRegistry(spc.DockerRegistry); err != nil {
		return fmt.Errorf("docker registry validation failed: %w", err)
	}
	
	if err := ValidateDockerImageName(spc.DockerImage); err != nil {
		return fmt.Errorf("docker image name validation failed: %w", err)
	}
	
	if err := ValidateBuildTags(spc.BuildTags); err != nil {
		return fmt.Errorf("build tags validation failed: %w", err)
	}
	
	if err := ValidateConfigState(spc.State); err != nil {
		return fmt.Errorf("config state validation failed: %w", err)
	}
	
	// Domain-Specific Invariants
	if err := spc.validateDockerSupportInvariant(); err != nil {
		return err
	}
	
	if err := spc.validateMainPathRequirementInvariant(); err != nil {
		return err
	}
	
	if err := spc.validateStateTransitionInvariant(); err != nil {
		return err
	}
	
	return nil
}

// validateDockerSupportInvariant ensures Docker is only enabled for supported project types
func (spc *SafeProjectConfig) validateDockerSupportInvariant() error {
	if spc.DockerEnabled && !spc.ProjectType.DockerSupported() {
		return fmt.Errorf("Docker is not supported for %s project type", spc.ProjectType)
	}
	return nil
}

// validateMainPathRequirementInvariant ensures main path is provided when required
func (spc *SafeProjectConfig) validateMainPathRequirementInvariant() error {
	if spc.ProjectType.RequiresMainPath() && spc.MainPath == "" {
		return fmt.Errorf("main path is required for %s project type", spc.ProjectType)
	}
	return nil
}

// validateStateTransitionInvariant ensures state transitions are valid
func (spc *SafeProjectConfig) validateStateTransitionInvariant() error {
	if !spc.State.AllowsValidation() {
		return fmt.Errorf("configuration in state '%s' cannot be validated", spc.State)
	}
	
	if !spc.State.AllowsGeneration() && spc.GenerateActions {
		return fmt.Errorf("configuration in state '%s' cannot generate actions", spc.State)
	}
	
	return nil
}

// Clone creates a deep copy of the configuration
func (spc *SafeProjectConfig) Clone() *SafeProjectConfig {
	clone := *spc
	
	// Deep copy slices
	if len(spc.Platforms) > 0 {
		clone.Platforms = make([]Platform, len(spc.Platforms))
		copy(clone.Platforms, spc.Platforms)
	}
	
	if len(spc.Architectures) > 0 {
		clone.Architectures = make([]Architecture, len(spc.Architectures))
		copy(clone.Architectures, spc.Architectures)
	}
	
	if len(spc.BuildTags) > 0 {
		clone.BuildTags = make([]BuildTag, len(spc.BuildTags))
		copy(clone.BuildTags, spc.BuildTags)
	}
	
	if len(spc.ActionsOn) > 0 {
		clone.ActionsOn = make([]ActionTrigger, len(spc.ActionsOn))
		copy(clone.ActionsOn, spc.ActionsOn)
	}
	
	return &clone
}

// Equal returns true if two configurations are equal (deep comparison)
func (spc *SafeProjectConfig) Equal(other *SafeProjectConfig) bool {
	if spc.ProjectName != other.ProjectName ||
		spc.ProjectDescription != other.ProjectDescription ||
		spc.ProjectType != other.ProjectType ||
		spc.BinaryName != other.BinaryName ||
		spc.MainPath != other.MainPath ||
		spc.GitProvider != other.GitProvider ||
		spc.DockerEnabled != other.DockerEnabled ||
		spc.DockerRegistry != other.DockerRegistry ||
		spc.DockerImage != other.DockerImage ||
		spc.Signing != other.Signing ||
		spc.Homebrew != other.Homebrew ||
		spc.Snap != other.Snap ||
		spc.SBOM != other.SBOM ||
		spc.GenerateActions != other.GenerateActions ||
		spc.ProVersion != other.ProVersion ||
		spc.State != other.State ||
		spc.CGOEnabled != other.CGOEnabled ||
		spc.LDFlags != other.LDFlags {
		return false
	}
	
	// Compare slices
	if !equalStringSlices(spc.Platforms, other.Platforms) ||
		!equalStringSlices(spc.Architectures, other.Architectures) ||
		!equalBuildTags(spc.BuildTags, other.BuildTags) ||
		!equalStringSlices(spc.ActionsOn, other.ActionsOn) {
		return false
	}
	
	return true
}

// Helper function to compare string slices
func equalStringSlices[T ~string](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}

// Helper function to compare BuildTag slices
func equalBuildTags(a, b []BuildTag) bool {
	if len(a) != len(b) {
		return false
	}
	for i, tag := range a {
		if b[i].Name != tag.Name || b[i].Description != tag.Description {
			return false
		}
	}
	return true
}

// Summary returns a human-readable summary of the configuration
func (spc *SafeProjectConfig) Summary() string {
	parts := []string{
		fmt.Sprintf("Project: %s (%s)", spc.ProjectName, spc.ProjectType),
		fmt.Sprintf("Binary: %s", spc.BinaryName),
		fmt.Sprintf("Platforms: %s", spc.joinPlatforms()),
		fmt.Sprintf("Architectures: %s", spc.joinArchitectures()),
		fmt.Sprintf("Provider: %s", spc.GitProvider),
	}
	
	if spc.DockerEnabled {
		parts = append(parts, fmt.Sprintf("Docker: %s", spc.DockerRegistry))
	}
	
	return strings.Join(parts, "\n")
}

// joinPlatforms creates a comma-separated list of platforms
func (spc *SafeProjectConfig) joinPlatforms() string {
	platforms := make([]string, len(spc.Platforms))
	for i, p := range spc.Platforms {
		platforms[i] = string(p)
	}
	return strings.Join(platforms, ", ")
}

// joinArchitectures creates a comma-separated list of architectures
func (spc *SafeProjectConfig) joinArchitectures() string {
	architectures := make([]string, len(spc.Architectures))
	for i, a := range spc.Architectures {
		architectures[i] = string(a)
	}
	return strings.Join(architectures, ", ")
}