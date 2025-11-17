package domain

import (
	"context"
	"fmt"
	"strings"
)

// CGOStatus represents CGO compilation status with compile-time safety
// Replaces bool CGOEnabled for better type safety and semantic clarity
type CGOStatus string

const (
	// CGOStatusDisabled disables CGO compilation completely
	CGOStatusDisabled CGOStatus = "disabled"
	// CGOStatusEnabled enables CGO compilation when available
	CGOStatusEnabled CGOStatus = "enabled"
	// CGOStatusRequired requires CGO compilation and will fail if not available
	CGOStatusRequired CGOStatus = "required"
)

// IsValid returns true if CGOStatus is valid
func (cs CGOStatus) IsValid() bool {
	switch cs {
	case CGOStatusDisabled, CGOStatusEnabled, CGOStatusRequired:
		return true
	default:
		return false
	}
}

// String returns human-readable display name
func (cs CGOStatus) String() string {
	switch cs {
	case CGOStatusDisabled:
		return "Disabled"
	case CGOStatusEnabled:
		return "Enabled"
	case CGOStatusRequired:
		return "Required"
	default:
		return "Unknown"
	}
}

// IsEnabled returns true if CGO is enabled (enabled or required)
func (cs CGOStatus) IsEnabled() bool {
	return cs == CGOStatusEnabled || cs == CGOStatusRequired
}

// IsRequired returns true if CGO is required
func (cs CGOStatus) IsRequired() bool {
	return cs == CGOStatusRequired
}

// ToBool converts to legacy boolean for compatibility
func (cs CGOStatus) ToBool() bool {
	return cs.IsEnabled()
}

// ValidateCGOStatus validates a CGO status
func ValidateCGOStatus(status CGOStatus) error {
	if !status.IsValid() {
		return NewValidationError(
			ErrInvalidCharacters,
			"Invalid CGO status",
			fmt.Sprintf("'%s' is not a valid CGO status", status),
		)
	}
	return nil
}

// DockerSupport represents Docker support level with compile-time safety
// Replaces bool DockerEnabled for better type safety and semantic clarity
type DockerSupport string

const (
	// DockerSupportNone disables Docker completely
	DockerSupportNone DockerSupport = "none"
	// DockerSupportBuild enables Docker image building only
	DockerSupportBuild DockerSupport = "build"
	// DockerSupportPublish enables Docker image publishing only
	DockerSupportPublish DockerSupport = "publish"
	// DockerSupportBoth enables both building and publishing
	DockerSupportBoth DockerSupport = "both"
)

// IsValid returns true if DockerSupport is valid
func (ds DockerSupport) IsValid() bool {
	switch ds {
	case DockerSupportNone, DockerSupportBuild, DockerSupportPublish, DockerSupportBoth:
		return true
	default:
		return false
	}
}

// String returns human-readable display name
func (ds DockerSupport) String() string {
	switch ds {
	case DockerSupportNone:
		return "None"
	case DockerSupportBuild:
		return "Build Only"
	case DockerSupportPublish:
		return "Publish Only"
	case DockerSupportBoth:
		return "Build and Publish"
	default:
		return "Unknown"
	}
}

// IsEnabled returns true if Docker support is enabled (build, publish, or both)
func (ds DockerSupport) IsEnabled() bool {
	return ds != DockerSupportNone
}

// ShouldBuild returns true if Docker images should be built
func (ds DockerSupport) ShouldBuild() bool {
	return ds == DockerSupportBuild || ds == DockerSupportBoth
}

// ShouldPublish returns true if Docker images should be published
func (ds DockerSupport) ShouldPublish() bool {
	return ds == DockerSupportPublish || ds == DockerSupportBoth
}

// ToBool converts to legacy boolean for compatibility
func (ds DockerSupport) ToBool() bool {
	return ds.IsEnabled()
}

// ValidateDockerSupport validates a Docker support level
func ValidateDockerSupport(support DockerSupport) error {
	if !support.IsValid() {
		return NewValidationError(
			ErrInvalidCharacters,
			"Invalid Docker support level",
			fmt.Sprintf("'%s' is not a valid Docker support level", support),
		)
	}
	return nil
}

// SigningLevel represents code signing level with compile-time safety
// Replaces bool Signing for better type safety and semantic clarity
type SigningLevel string

const (
	// SigningLevelNone disables code signing completely
	SigningLevelNone SigningLevel = "none"
	// SigningLevelBasic enables basic code signing
	SigningLevelBasic SigningLevel = "basic"
	// SigningLevelAdvanced enables advanced code signing with additional verification
	SigningLevelAdvanced SigningLevel = "advanced"
	// SigningLevelEnterprise enables enterprise-level code signing with full compliance
	SigningLevelEnterprise SigningLevel = "enterprise"
)

// IsValid returns true if SigningLevel is valid
func (sl SigningLevel) IsValid() bool {
	switch sl {
	case SigningLevelNone, SigningLevelBasic, SigningLevelAdvanced, SigningLevelEnterprise:
		return true
	default:
		return false
	}
}

// String returns human-readable display name
func (sl SigningLevel) String() string {
	switch sl {
	case SigningLevelNone:
		return "None"
	case SigningLevelBasic:
		return "Basic"
	case SigningLevelAdvanced:
		return "Advanced"
	case SigningLevelEnterprise:
		return "Enterprise"
	default:
		return "Unknown"
	}
}

// IsEnabled returns true if signing is enabled (basic, advanced, or enterprise)
func (sl SigningLevel) IsEnabled() bool {
	return sl != SigningLevelNone
}

// ToBool converts to legacy boolean for compatibility
func (sl SigningLevel) ToBool() bool {
	return sl.IsEnabled()
}

// ValidateSigningLevel validates a signing level
func ValidateSigningLevel(level SigningLevel) error {
	if !level.IsValid() {
		return NewValidationError(
			ErrInvalidCharacters,
			"Invalid signing level",
			fmt.Sprintf("'%s' is not a valid signing level", level),
		)
	}
	return nil
}

// ActionLevel represents GitHub Actions generation level with compile-time safety
// Replaces bool GenerateActions for better type safety and semantic clarity
type ActionLevel string

const (
	// ActionLevelNone disables GitHub Actions generation completely
	ActionLevelNone ActionLevel = "none"
	// ActionLevelBasic generates basic GitHub Actions workflow
	ActionLevelBasic ActionLevel = "basic"
	// ActionLevelAdvanced generates advanced GitHub Actions with caching and optimization
	ActionLevelAdvanced ActionLevel = "advanced"
)

// IsValid returns true if ActionLevel is valid
func (al ActionLevel) IsValid() bool {
	switch al {
	case ActionLevelNone, ActionLevelBasic, ActionLevelAdvanced:
		return true
	default:
		return false
	}
}

// String returns human-readable display name
func (al ActionLevel) String() string {
	switch al {
	case ActionLevelNone:
		return "None"
	case ActionLevelBasic:
		return "Basic"
	case ActionLevelAdvanced:
		return "Advanced"
	default:
		return "Unknown"
	}
}

// IsEnabled returns true if actions should be generated
func (al ActionLevel) IsEnabled() bool {
	return al != ActionLevelNone
}

// ToBool converts to legacy boolean for compatibility
func (al ActionLevel) ToBool() bool {
	return al.IsEnabled()
}

// ValidateActionLevel validates an action level
func ValidateActionLevel(level ActionLevel) error {
	if !level.IsValid() {
		return NewValidationError(
			ErrInvalidCharacters,
			"Invalid action level",
			fmt.Sprintf("'%s' is not a valid action level", level),
		)
	}
	return nil
}

// FeatureLevel represents feature tier with compile-time safety
// Replaces bool ProVersion for better type safety and semantic clarity
type FeatureLevel string

const (
	// FeatureLevelBasic includes basic features for standard projects
	FeatureLevelBasic FeatureLevel = "basic"
	// FeatureLevelProfessional includes professional features for larger projects
	FeatureLevelProfessional FeatureLevel = "professional"
	// FeatureLevelEnterprise includes enterprise features with full compliance
	FeatureLevelEnterprise FeatureLevel = "enterprise"
)

// IsValid returns true if FeatureLevel is valid
func (fl FeatureLevel) IsValid() bool {
	switch fl {
	case FeatureLevelBasic, FeatureLevelProfessional, FeatureLevelEnterprise:
		return true
	default:
		return false
	}
}

// String returns human-readable display name
func (fl FeatureLevel) String() string {
	switch fl {
	case FeatureLevelBasic:
		return "Basic"
	case FeatureLevelProfessional:
		return "Professional"
	case FeatureLevelEnterprise:
		return "Enterprise"
	default:
		return "Unknown"
	}
}

// IsPro returns true if professional or enterprise features are enabled
func (fl FeatureLevel) IsPro() bool {
	return fl == FeatureLevelProfessional || fl == FeatureLevelEnterprise
}

// ToBool converts to legacy boolean for compatibility
func (fl FeatureLevel) ToBool() bool {
	return fl.IsPro()
}

// ValidateFeatureLevel validates a feature level
func ValidateFeatureLevel(level FeatureLevel) error {
	if !level.IsValid() {
		return NewValidationError(
			ErrInvalidCharacters,
			"Invalid feature level",
			fmt.Sprintf("'%s' is not a valid feature level", level),
		)
	}
	return nil
}

// Enum migration utilities for backward compatibility

// CGOStatusFromBool converts legacy boolean to CGOStatus
func CGOStatusFromBool(enabled bool) CGOStatus {
	if enabled {
		return CGOStatusEnabled
	}
	return CGOStatusDisabled
}

// DockerSupportFromBool converts legacy boolean to DockerSupport
func DockerSupportFromBool(enabled bool) DockerSupport {
	if enabled {
		return DockerSupportBoth
	}
	return DockerSupportNone
}

// SigningLevelFromBool converts legacy boolean to SigningLevel
func SigningLevelFromBool(enabled bool) SigningLevel {
	if enabled {
		return SigningLevelBasic
	}
	return SigningLevelNone
}

// ActionLevelFromBool converts legacy boolean to ActionLevel
func ActionLevelFromBool(enabled bool) ActionLevel {
	if enabled {
		return ActionLevelBasic
	}
	return ActionLevelNone
}

// FeatureLevelFromBool converts legacy boolean to FeatureLevel
func FeatureLevelFromBool(enabled bool) FeatureLevel {
	if enabled {
		return FeatureLevelProfessional
	}
	return FeatureLevelBasic
}

// Smart conversion functions based on project type context

// GetDefaultCGOStatus returns smart CGO status based on project type
func GetDefaultCGOStatus(projectType ProjectType) CGOStatus {
	if projectType.DefaultCGOEnabled() {
		return CGOStatusEnabled
	}
	return CGOStatusDisabled
}

// GetDefaultDockerSupport returns smart Docker support based on project type
func GetDefaultDockerSupport(projectType ProjectType) DockerSupport {
	if projectType.DockerSupported() {
		return DockerSupportBuild
	}
	return DockerSupportNone
}

// GetRecommendedActionLevel returns recommended action level based on project type
func GetRecommendedActionLevel(projectType ProjectType) ActionLevel {
	switch projectType {
	case ProjectTypeCLI, ProjectTypeAPI:
		return ActionLevelAdvanced
	case ProjectTypeWeb:
		return ActionLevelBasic
	case ProjectTypeLibrary:
		return ActionLevelBasic
	case ProjectTypeDesktop:
		return ActionLevelAdvanced
	default:
		return ActionLevelBasic
	}
}

// GetRecommendedSigningLevel returns recommended signing level based on project type
func GetRecommendedSigningLevel(projectType ProjectType) SigningLevel {
	switch projectType {
	case ProjectTypeCLI:
		return SigningLevelBasic
	case ProjectTypeWeb, ProjectTypeAPI:
		return SigningLevelAdvanced
	case ProjectTypeDesktop:
		return SigningLevelEnterprise
	case ProjectTypeLibrary:
		return SigningLevelNone
	default:
		return SigningLevelBasic
	}
}

// GetRecommendedFeatureLevel returns recommended feature level based on project type
func GetRecommendedFeatureLevel(projectType ProjectType) FeatureLevel {
	switch projectType {
	case ProjectTypeAPI, ProjectTypeWeb:
		return FeatureLevelProfessional
	case ProjectTypeDesktop:
		return FeatureLevelEnterprise
	case ProjectTypeCLI, ProjectTypeLibrary:
		return FeatureLevelBasic
	default:
		return FeatureLevelBasic
	}
}