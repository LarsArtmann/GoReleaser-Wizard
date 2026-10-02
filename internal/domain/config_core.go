package domain

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	"github.com/go-faster/yaml"
)

// SafeProjectConfig represents single source of truth for project configuration
// Generated from TypeSpec specification - DO NOT MODIFY MANUALLY.
type SafeProjectConfig struct {
	// Basic Information
	ProjectName        string      `json:"project_name"                  yaml:"project_name"`
	ProjectDescription string      `json:"project_description,omitempty" yaml:"project_description,omitempty"`
	ProjectType        ProjectType `json:"project_type"                  yaml:"project_type"`
	BinaryName         string      `json:"binary_name"                   yaml:"binary_name"`
	MainPath           string      `json:"main_path"                     yaml:"main_path"`

	// Build Configuration
	Platforms     []Platform     `json:"platforms"            yaml:"platforms"`
	Architectures []Architecture `json:"architectures"        yaml:"architectures"`
	CGOStatus     CGOStatus      `json:"cgo_status"           yaml:"cgo_status"`
	BuildTags     []BuildTag     `json:"build_tags,omitempty" yaml:"build_tags,omitempty"`
	LDFlags       bool           `json:"ldflags"              yaml:"ldflags"`

	// Release Configuration
	GitProvider    GitProvider    `json:"git_provider"           yaml:"git_provider"`
	DockerSupport  DockerSupport  `json:"docker_support"         yaml:"docker_support"`
	DockerRegistry DockerRegistry `json:"docker_registry"        yaml:"docker_registry"`
	DockerImage    string         `json:"docker_image,omitempty" yaml:"docker_image,omitempty"`
	SigningLevel   SigningLevel   `json:"signing_level"          yaml:"signing_level"`
	Homebrew       bool           `json:"homebrew"               yaml:"homebrew"`
	Snap           bool           `json:"snap"                   yaml:"snap"`
	SBOM           bool           `json:"sbom"                   yaml:"sbom"`

	// CI/CD Configuration
	ActionLevel ActionLevel     `json:"action_level" yaml:"action_level"`
	ActionsOn   []ActionTrigger `json:"actions_on"   yaml:"actions_on"`

	// Advanced Features
	FeatureLevel FeatureLevel `json:"feature_level" yaml:"feature_level"`

	// State Management
	State ConfigState `json:"state" yaml:"state"`
}

// NewSafeProjectConfig creates a new safe configuration with smart defaults.
func NewSafeProjectConfig() *SafeProjectConfig {
	return &SafeProjectConfig{
		// Smart defaults based on project analysis
		ProjectType:    GetRecommendedProjectType(),
		Platforms:      GetRecommendedPlatforms(),
		Architectures:  GetRecommendedArchitectures(),
		GitProvider:    GetRecommendedGitProvider(),
		DockerRegistry: GetRecommendedDockerRegistry(),
		CGOStatus:      CGOStatusDisabled,
		ActionLevel:    ActionLevelBasic,
		SigningLevel:   SigningLevelNone,
		FeatureLevel:   FeatureLevelBasic,
		State:          ConfigStateDraft,
		LDFlags:        true,
		Homebrew:       false,
		Snap:           false,
		SBOM:           false,
	}
}

// Validate validates the configuration.
func (spc *SafeProjectConfig) Validate() error {
	if err := spc.validateBasicFields(); err != nil {
		return err
	}

	if err := spc.validateEnums(); err != nil {
		return err
	}

	if err := spc.validateCollections(); err != nil {
		return err
	}

	// Validate business rules
	return spc.validateBusinessRules()
}

// validateBasicFields validates the scalar config fields.
func (spc *SafeProjectConfig) validateBasicFields() error {
	for _, err := range []error{
		ValidateProjectName(spc.ProjectName),
		ValidateBinaryName(spc.BinaryName),
		ValidateMainPath(spc.MainPath),
	} {
		if err != nil {
			return err
		}
	}

	if spc.ProjectDescription != "" {
		return ValidateProjectDescription(spc.ProjectDescription)
	}

	return nil
}

// validateEnums validates every enum-typed field.
func (spc *SafeProjectConfig) validateEnums() error {
	enumChecks := []struct {
		valid   bool
		errType ErrorCode
		message string
		value   string
		field   string
	}{
		{spc.ProjectType.IsValid(), ErrInvalidConfig, "Invalid project type", string(spc.ProjectType), "project_type"},
		{
			spc.GitProvider.IsValid(),
			ErrInvalidGitProvider,
			"Invalid Git provider",
			string(spc.GitProvider),
			"git_provider",
		},
		{spc.CGOStatus.IsValid(), ErrInvalidConfig, "Invalid CGO status", string(spc.CGOStatus), "cgo_status"},
		{
			spc.DockerSupport.IsValid(),
			ErrInvalidConfig,
			"Invalid Docker support",
			string(spc.DockerSupport),
			"docker_support",
		},
		{
			spc.DockerRegistry.IsValid(),
			ErrInvalidDockerRegistry,
			"Invalid Docker registry",
			string(spc.DockerRegistry),
			"docker_registry",
		},
		{
			spc.SigningLevel.IsValid(),
			ErrInvalidConfig,
			"Invalid signing level",
			string(spc.SigningLevel),
			"signing_level",
		},
		{spc.ActionLevel.IsValid(), ErrInvalidConfig, "Invalid action level", string(spc.ActionLevel), "action_level"},
		{
			spc.FeatureLevel.IsValid(),
			ErrInvalidConfig,
			"Invalid feature level",
			string(spc.FeatureLevel),
			"feature_level",
		},
		{spc.State.IsValid(), ErrInvalidConfigState, "Invalid config state", string(spc.State), "state"},
	}

	for _, check := range enumChecks {
		if !check.valid {
			return NewValidationError(
				check.errType,
				check.message,
				check.value,
			).WithField(check.field)
		}
	}

	return nil
}

// validateCollections validates the slice-typed fields.
func (spc *SafeProjectConfig) validateCollections() error {
	for _, platform := range spc.Platforms {
		if !platform.IsValid() {
			return NewValidationError(
				ErrInvalidPlatform,
				"Invalid platform",
				string(platform),
			).WithField("platforms")
		}
	}

	for _, arch := range spc.Architectures {
		if !arch.IsValid() {
			return NewValidationError(
				ErrInvalidArchitecture,
				"Invalid architecture",
				string(arch),
			).WithField("architectures")
		}
	}

	for _, tag := range spc.BuildTags {
		if !tag.IsValid() {
			return NewValidationError(
				ErrInvalidBuildTag,
				"Invalid build tag",
				tag.String(),
			).WithField("build_tags")
		}
	}

	for _, trigger := range spc.ActionsOn {
		if !trigger.IsValid() {
			return NewValidationError(
				ErrInvalidActionTrigger,
				"Invalid action trigger",
				string(trigger),
			).WithField("actions_on")
		}
	}

	return nil
}

// validateBusinessRules validates business logic rules.
func (spc *SafeProjectConfig) validateBusinessRules() error {
	// Validate platform/architecture compatibility
	err := spc.validatePlatformArchCompatibility()
	if err != nil {
		return err
	}

	// Validate Docker configuration
	err = spc.validateDockerConfiguration()
	if err != nil {
		return err
	}

	// Validate Actions configuration
	err = spc.validateActionsConfiguration()
	if err != nil {
		return err
	}

	return nil
}

// validatePlatformArchCompatibility validates platform/architecture compatibility.
func (spc *SafeProjectConfig) validatePlatformArchCompatibility() error {
	for _, platform := range spc.Platforms {
		for _, arch := range spc.Architectures {
			if !arch.IsCompatibleWith(platform) {
				return NewValidationError(
					ErrInvalidProjectName,
					"Platform/architecture combination not supported",
					fmt.Sprintf("%s + %s", platform.String(), arch.String()),
				).WithField("platform_arch_compatibility")
			}
		}
	}

	return nil
}

// validateDockerConfiguration validates Docker configuration.
func (spc *SafeProjectConfig) validateDockerConfiguration() error {
	if spc.DockerSupport.IsEnabled() {
		// Check if Docker is supported by project type
		if !spc.ProjectType.DockerSupported() {
			return NewValidationError(
				ErrInvalidProjectName,
				"Docker not supported for project type",
				spc.ProjectType.String(),
			).WithField("docker_support")
		}

		// Validate Docker image name
		if spc.DockerSupport.IsDeployEnabled() && spc.GetDockerImageName() == "" {
			return NewValidationError(
				ErrInvalidProjectName,
				"Docker image name is required for deployment",
				"Docker image name cannot be empty when deployment is enabled",
			).WithField("docker_image")
		}
	}

	return nil
}

// validateActionsConfiguration validates Actions configuration.
func (spc *SafeProjectConfig) validateActionsConfiguration() error {
	if spc.ActionLevel.IsEnabled() {
		// Check if Git provider supports actions
		if !spc.GitProvider.ActionsSupported() {
			return NewValidationError(
				ErrInvalidProjectName,
				"Git provider doesn't support actions",
				spc.GitProvider.String(),
			).WithField("git_provider")
		}

		// Validate action triggers
		if len(spc.ActionsOn) == 0 {
			return NewValidationError(
				ErrInvalidProjectName,
				"Action triggers are required when actions are enabled",
				"At least one action trigger must be specified",
			).WithField("actions_on")
		}
	}

	return nil
}

// Clone creates a deep copy of the configuration.
func (spc *SafeProjectConfig) Clone() *SafeProjectConfig {
	clone := *spc

	// Deep copy slices
	if spc.Platforms != nil {
		clone.Platforms = make([]Platform, 0, len(spc.Platforms))
		clone.Platforms = append(clone.Platforms, spc.Platforms...)
	}

	if spc.Architectures != nil {
		clone.Architectures = make([]Architecture, 0, len(spc.Architectures))
		clone.Architectures = append(clone.Architectures, spc.Architectures...)
	}

	if spc.BuildTags != nil {
		clone.BuildTags = make([]BuildTag, 0, len(spc.BuildTags))
		clone.BuildTags = append(clone.BuildTags, spc.BuildTags...)
	}

	if spc.ActionsOn != nil {
		clone.ActionsOn = make([]ActionTrigger, 0, len(spc.ActionsOn))
		clone.ActionsOn = append(clone.ActionsOn, spc.ActionsOn...)
	}

	return &clone
}

// serializeToFormat serializes the config using the provided marshal function.
func (spc *SafeProjectConfig) serializeToFormat(
	marshalFunc func(any) ([]byte, error),
	formatName string,
) (string, error) {
	data, err := marshalFunc(spc)
	if err != nil {
		return "", NewValidationError(
			ErrInvalidProjectName,
			"Failed to serialize configuration to "+formatName,
			err.Error(),
		)
	}

	return string(data), nil
}

// ToJSON converts configuration to JSON string.
func (spc *SafeProjectConfig) ToJSON() (string, error) {
	return spc.serializeToFormat(func(v any) ([]byte, error) {
		return json.Marshal(v, jsontext.WithIndentPrefix(""), jsontext.WithIndent("  "))
	}, "JSON")
}

// ToYAML converts configuration to YAML string.
func (spc *SafeProjectConfig) ToYAML() (string, error) {
	return spc.serializeToFormat(yaml.Marshal, "YAML")
}

// FromJSON loads configuration from JSON string.
func (spc *SafeProjectConfig) FromJSON(jsonStr string) error {
	return fmt.Errorf("failed to unmarshal JSON: %w", json.Unmarshal([]byte(jsonStr), spc))
}

// FromYAML loads configuration from YAML string.
func (spc *SafeProjectConfig) FromYAML(yamlStr string) error {
	return fmt.Errorf("failed to unmarshal YAML: %w", yaml.Unmarshal([]byte(yamlStr), spc))
}

// GetSummary returns a summary of the configuration.
func (spc *SafeProjectConfig) GetSummary() ConfigSummary {
	return ConfigSummary{
		ProjectName:    spc.ProjectName,
		ProjectType:    spc.ProjectType.String(),
		BinaryName:     spc.BinaryName,
		PlatformCount:  len(spc.Platforms),
		ArchCount:      len(spc.Architectures),
		CGOEnabled:     spc.CGOStatus.IsEnabled(),
		DockerEnabled:  spc.DockerSupport.IsEnabled(),
		SigningEnabled: spc.SigningLevel.IsEnabled(),
		ActionsEnabled: spc.ActionLevel.IsEnabled(),
		State:          spc.State.String(),
		FeatureLevel:   spc.FeatureLevel.String(),
	}
}

// ConfigSummary represents a summary of configuration.
type ConfigSummary struct {
	ProjectName    string `json:"project_name"`
	ProjectType    string `json:"project_type"`
	BinaryName     string `json:"binary_name"`
	PlatformCount  int    `json:"platform_count"`
	ArchCount      int    `json:"arch_count"`
	CGOEnabled     bool   `json:"cgo_enabled"`
	DockerEnabled  bool   `json:"docker_enabled"`
	SigningEnabled bool   `json:"signing_enabled"`
	ActionsEnabled bool   `json:"actions_enabled"`
	State          string `json:"state"`
	FeatureLevel   string `json:"feature_level"`
}
