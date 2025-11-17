package domain

import (
	"fmt"
	"strings"
)

// SafeProjectConfig represents single source of truth for project configuration
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
	CGOStatus     CGOStatus     `json:"cgo_status" yaml:"cgo_status"`
	BuildTags     []BuildTag     `json:"build_tags,omitempty" yaml:"build_tags,omitempty"`
	LDFlags       bool           `json:"ldflags" yaml:"ldflags"`

	// Release Configuration
	GitProvider    GitProvider    `json:"git_provider" yaml:"git_provider"`
	DockerSupport  DockerSupport  `json:"docker_support" yaml:"docker_support"`
	DockerRegistry DockerRegistry `json:"docker_registry" yaml:"docker_registry"`
	DockerImage    string         `json:"docker_image,omitempty" yaml:"docker_image,omitempty"`
	SigningLevel   SigningLevel   `json:"signing_level" yaml:"signing_level"`
	Homebrew       bool           `json:"homebrew" yaml:"homebrew"`
	Snap           bool           `json:"snap" yaml:"snap"`
	SBOM           bool           `json:"sbom" yaml:"sbom"`

	// CI/CD Configuration
	ActionLevel   ActionLevel    `json:"action_level" yaml:"action_level"`
	ActionsOn     []ActionTrigger `json:"actions_on" yaml:"actions_on"`

	// Advanced Features
	FeatureLevel FeatureLevel `json:"feature_level" yaml:"feature_level"`

	// State Management
	State ConfigState `json:"state" yaml:"state"`
}

// NewSafeProjectConfig creates a new safe configuration with smart defaults
func NewSafeProjectConfig() *SafeProjectConfig {
	return &SafeProjectConfig{
		// Smart defaults based on project analysis
		ProjectType:      GetRecommendedProjectType(),
		Platforms:        GetRecommendedPlatforms(),
		Architectures:    GetRecommendedArchitectures(),
		GitProvider:      GetRecommendedGitProvider(),
		DockerRegistry:   GetRecommendedDockerRegistry(),
		CGOStatus:        CGOStatusDisabled,
		ActionLevel:      ActionLevelBasic,
		SigningLevel:     SigningLevelNone,
		FeatureLevel:     FeatureLevelBasic,
		State:            ConfigStateDraft,
		LDFlags:          true,
		Homebrew:         false,
		Snap:             false,
		SBOM:             false,
	}
}

// ApplyDefaults applies smart defaults based on project type and context
func (spc *SafeProjectConfig) ApplyDefaults() {
	// Apply project type-specific defaults
	if spc.CGOStatus == CGOStatusDisabled && spc.ProjectType.DefaultCGOEnabled() {
		spc.CGOStatus = CGOStatusEnabled
	}

	if spc.Platforms == nil || len(spc.Platforms) == 0 {
		spc.Platforms = spc.ProjectType.RecommendedPlatforms()
	}

	if spc.Architectures == nil || len(spc.Architectures) == 0 {
		spc.Architectures = GetRecommendedArchitectures()
	}

	if spc.GitProvider == "" {
		spc.GitProvider = GetRecommendedGitProvider()
	}

	// Apply Docker support defaults
	if spc.DockerSupport == DockerSupportNone && spc.ProjectType.DockerSupported() {
		spc.DockerSupport = DockerSupportBuild
	}

	if spc.DockerRegistry == "" && spc.DockerSupport.IsEnabled() {
		spc.DockerRegistry = spc.GitProvider.DefaultRegistry()
	}

	// Apply action level defaults
	if spc.ActionLevel == ActionLevelNone && spc.GitProvider.ActionsSupported() {
		spc.ActionLevel = ActionLevelBasic
		if spc.ActionsOn == nil || len(spc.ActionsOn) == 0 {
			spc.ActionsOn = GetRecommendedTriggers(spc.ProjectType)
		}
	}

	// Apply feature level defaults based on project type
	if spc.FeatureLevel == FeatureLevelBasic {
		spc.FeatureLevel = GetRecommendedFeatureLevel(spc.ProjectType)
	}

	// Apply signing level defaults
	if spc.SigningLevel == SigningLevelNone {
		spc.SigningLevel = GetRecommendedSigningLevel(spc.ProjectType)
	}

	// Apply defaults for other fields
	if spc.BinaryName == "" && spc.ProjectName != "" {
		spc.BinaryName = spc.ProjectName
	}

	if spc.DockerImage == "" && spc.ProjectName != "" && spc.DockerSupport.IsEnabled() {
		spc.DockerImage = strings.ToLower(spc.ProjectName)
	}
}

// ValidateInvariants enforces domain invariants and returns any violations
func (spc *SafeProjectConfig) ValidateInvariants() error {
	// Basic validation
	if err := ValidateProjectName(spc.ProjectName); err != nil {
		return err
	}

	if err := ValidateBinaryName(spc.BinaryName); err != nil {
		return err
	}

	if err := ValidateMainPath(spc.MainPath); err != nil {
		return err
	}

	if err := ValidateProjectDescription(spc.ProjectDescription); err != nil {
		return err
	}

	// Type validation
	if !spc.ProjectType.IsValid() {
		return fmt.Errorf("invalid project type: %s", spc.ProjectType)
	}

	if err := ValidatePlatforms(spc.Platforms); err != nil {
		return err
	}

	if err := ValidateArchitectures(spc.Architectures); err != nil {
		return err
	}

	if err := ValidateGitProvider(spc.GitProvider); err != nil {
		return err
	}

	if err := ValidateDockerRegistry(spc.DockerRegistry); err != nil {
		return err
	}

	if err := ValidateActionTriggers(spc.ActionsOn); err != nil {
		return err
	}

	// CGO status validation
	if err := ValidateCGOStatus(spc.CGOStatus); err != nil {
		return err
	}

	// Docker support validation
	if err := ValidateDockerSupport(spc.DockerSupport); err != nil {
		return err
	}

	// Signing level validation
	if err := ValidateSigningLevel(spc.SigningLevel); err != nil {
		return err
	}

	// Action level validation
	if err := ValidateActionLevel(spc.ActionLevel); err != nil {
		return err
	}

	// Feature level validation
	if err := ValidateFeatureLevel(spc.FeatureLevel); err != nil {
		return err
	}

	// Config state validation
	if err := ValidateConfigState(spc.State); err != nil {
		return err
	}

	// Cross-field invariants
	if spc.DockerSupport.IsEnabled() && !spc.ProjectType.DockerSupported() {
		return fmt.Errorf("docker support enabled but project type %s does not support docker", spc.ProjectType)
	}

	if spc.CGOStatus.IsEnabled() && spc.CGOStatus.IsRequired() {
		hasCGOSupport := false
		for _, platform := range spc.Platforms {
			if platform.SupportsCGO() {
				hasCGOSupport = true
				break
			}
		}
		if !hasCGOSupport {
			return fmt.Errorf("cgo required but no selected platforms support cgo")
		}
	}

	// Platform-architecture compatibility
	return ValidatePlatformArchCompatibility(spc.Platforms, spc.Architectures)
}

// Clone creates a deep copy of the configuration
func (spc *SafeProjectConfig) Clone() *SafeProjectConfig {
	clone := *spc
	
	// Deep copy slices
	if spc.Platforms != nil {
		clone.Platforms = make([]Platform, len(spc.Platforms))
		copy(clone.Platforms, spc.Platforms)
	}
	
	if spc.Architectures != nil {
		clone.Architectures = make([]Architecture, len(spc.Architectures))
		copy(clone.Architectures, spc.Architectures)
	}
	
	if spc.BuildTags != nil {
		clone.BuildTags = make([]BuildTag, len(spc.BuildTags))
		copy(clone.BuildTags, spc.BuildTags)
	}
	
	if spc.ActionsOn != nil {
		clone.ActionsOn = make([]ActionTrigger, len(spc.ActionsOn))
		copy(clone.ActionsOn, spc.ActionsOn)
	}
	
	return &clone
}

// Equals returns true if two configurations are equivalent
func (spc *SafeProjectConfig) Equals(other *SafeProjectConfig) bool {
	if spc == nil || other == nil {
		return spc == other
	}
	
	return spc.ProjectName == other.ProjectName &&
		spc.ProjectDescription == other.ProjectDescription &&
		spc.ProjectType == other.ProjectType &&
		spc.BinaryName == other.BinaryName &&
		spc.MainPath == other.MainPath &&
		spc.LDFlags == other.LDFlags &&
		spc.GitProvider == other.GitProvider &&
		spc.DockerRegistry == other.DockerRegistry &&
		spc.DockerImage == other.DockerImage &&
		spc.Homebrew == other.Homebrew &&
		spc.Snap == other.Snap &&
		spc.SBOM == other.SBOM &&
		spc.State == other.State
}

// HasChanged returns true if any critical field has changed
func (spc *SafeProjectConfig) HasChanged(other *SafeProjectConfig) bool {
	return !spc.Equals(other)
}

// IsReadyForGeneration returns true if configuration is ready for file generation
func (spc *SafeProjectConfig) IsReadyForGeneration() bool {
	return spc.State.AllowsGeneration() &&
		spc.ProjectName != "" &&
		spc.ProjectType.IsValid() &&
		spc.BinaryName != "" &&
		len(spc.Platforms) > 0 &&
		len(spc.Architectures) > 0
}

// GetDockerImageName returns the full Docker image name
func (spc *SafeProjectConfig) GetDockerImageName() string {
	if spc.DockerImage != "" {
		return spc.DockerImage
	}
	return strings.ToLower(spc.ProjectName)
}

// ShouldGenerateDockerFiles returns true if Docker files should be generated
func (spc *SafeProjectConfig) ShouldGenerateDockerFiles() bool {
	return spc.DockerSupport.ShouldBuild() &&
		spc.ProjectType.DockerSupported() &&
		spc.DockerRegistry != ""
}

// ShouldGenerateActionsFiles returns true if Actions files should be generated
func (spc *SafeProjectConfig) ShouldGenerateActionsFiles() bool {
	return spc.ActionLevel.IsEnabled() &&
		spc.GitProvider.ActionsSupported() &&
		len(spc.ActionsOn) > 0
}

// ShouldSignReleases returns true if releases should be signed
func (spc *SafeProjectConfig) ShouldSignReleases() bool {
	return spc.SigningLevel.IsEnabled() &&
		spc.SigningLevel != SigningLevelNone
}

// IsProFeatures returns true if pro features are enabled
func (spc *SafeProjectConfig) IsProFeatures() bool {
	return spc.FeatureLevel.IsPro()
}

// TransitionToState safely transitions configuration to a new state
func (spc *SafeProjectConfig) TransitionToState(newState ConfigState) error {
	// Validate state transition
	if !spc.State.AllowsTransitionTo(newState) {
		return fmt.Errorf("invalid state transition from %s to %s", spc.State, newState)
	}
	
	spc.State = newState
	return nil
}

// AllowsTransitionTo checks if state transition is allowed
func (cs ConfigState) AllowsTransitionTo(newState ConfigState) bool {
	// Define allowed state transitions
	allowedTransitions := map[ConfigState][]ConfigState{
		ConfigStateDraft:      {ConfigStateValid, ConfigStateInvalid},
		ConfigStateValid:      {ConfigStateInvalid, ConfigStateProcessing},
		ConfigStateInvalid:    {ConfigStateValid, ConfigStateProcessing},
		ConfigStateProcessing: {ConfigStateValid, ConfigStateInvalid, ConfigStateGenerated},
		ConfigStateGenerated:  {}, // Final state - no transitions allowed
	}
	
	for _, allowed := range allowedTransitions[cs] {
		if allowed == newState {
			return true
		}
	}
	
	return false
}

// Legacy compatibility methods - DEPRECATED but kept for migration

// GetCGOEnabled returns legacy boolean for CGO
func (spc *SafeProjectConfig) GetCGOEnabled() bool {
	return spc.CGOStatus.ToBool()
}

// SetCGOEnabled sets CGO status from legacy boolean
func (spc *SafeProjectConfig) SetCGOEnabled(enabled bool) {
	if enabled {
		spc.CGOStatus = CGOStatusEnabled
	} else {
		spc.CGOStatus = CGOStatusDisabled
	}
}

// GetDockerEnabled returns legacy boolean for Docker
func (spc *SafeProjectConfig) GetDockerEnabled() bool {
	return spc.DockerSupport.ToBool()
}

// SetDockerEnabled sets Docker support from legacy boolean
func (spc *SafeProjectConfig) SetDockerEnabled(enabled bool) {
	spc.DockerSupport = DockerSupportFromBool(enabled)
}

// GetSigning returns legacy boolean for Signing
func (spc *SafeProjectConfig) GetSigning() bool {
	return spc.SigningLevel.ToBool()
}

// SetSigning sets signing level from legacy boolean
func (spc *SafeProjectConfig) SetSigning(enabled bool) {
	spc.SigningLevel = SigningLevelFromBool(enabled)
}

// GetGenerateActions returns legacy boolean for Actions
func (spc *SafeProjectConfig) GetGenerateActions() bool {
	return spc.ActionLevel.ToBool()
}

// SetGenerateActions sets action level from legacy boolean
func (spc *SafeProjectConfig) SetGenerateActions(enabled bool) {
	spc.ActionLevel = ActionLevelFromBool(enabled)
}

// GetProVersion returns legacy boolean for Pro features
func (spc *SafeProjectConfig) GetProVersion() bool {
	return spc.FeatureLevel.ToBool()
}

// SetProVersion sets feature level from legacy boolean
func (spc *SafeProjectConfig) SetProVersion(enabled bool) {
	spc.FeatureLevel = FeatureLevelFromBool(enabled)
}