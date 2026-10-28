package domain

// ProjectType represents different types of Go projects
// This enum replaces string-based project types for type safety.
type ProjectType string

const (
	ProjectTypeCLI          ProjectType = "cli"
	ProjectTypeLibrary      ProjectType = "library"
	ProjectTypeWebAPI       ProjectType = "webapi"
	ProjectTypeGRPCService  ProjectType = "grpc-service"
	ProjectTypeMicroservice ProjectType = "microservice"
	ProjectTypeDesktop      ProjectType = "desktop"
	ProjectTypeMobile       ProjectType = "mobile"
	ProjectTypePlugin       ProjectType = "plugin"
	ProjectTypeDaemon       ProjectType = "daemon"
	ProjectTypeTool         ProjectType = "tool"
)

// IsValid returns true if ProjectType is valid.
func (pt ProjectType) IsValid() bool {
	switch pt {
	case ProjectTypeCLI, ProjectTypeLibrary, ProjectTypeWebAPI,
		ProjectTypeGRPCService, ProjectTypeMicroservice, ProjectTypeDesktop,
		ProjectTypeMobile, ProjectTypePlugin, ProjectTypeDaemon, ProjectTypeTool:
		return true
	default:
		return false
	}
}

// String returns human-readable display name.
func (pt ProjectType) String() string {
	switch pt {
	case ProjectTypeCLI:
		return "CLI Application"
	case ProjectTypeLibrary:
		return "Library"
	case ProjectTypeWebAPI:
		return "Web API"
	case ProjectTypeGRPCService:
		return "gRPC Service"
	case ProjectTypeMicroservice:
		return "Microservice"
	case ProjectTypeDesktop:
		return "Desktop Application"
	case ProjectTypeMobile:
		return "Mobile Application"
	case ProjectTypePlugin:
		return "Plugin"
	case ProjectTypeDaemon:
		return "Daemon/Service"
	case ProjectTypeTool:
		return "Command Line Tool"
	default:
		return unknownValue
	}
}

// IsWebRelated returns true for web-related project types.
func (pt ProjectType) IsWebRelated() bool {
	return pt == ProjectTypeWebAPI || pt == ProjectTypeGRPCService || pt == ProjectTypeMicroservice
}

// IsDesktopRelated returns true for desktop-related project types.
func (pt ProjectType) IsDesktopRelated() bool {
	return pt == ProjectTypeDesktop || pt == ProjectTypeDaemon || pt == ProjectTypeTool
}

// IsMobileRelated returns true for mobile-related project types.
func (pt ProjectType) IsMobileRelated() bool {
	return pt == ProjectTypeMobile
}

// IsLibrary returns true for library project type.
func (pt ProjectType) IsLibrary() bool {
	return pt == ProjectTypeLibrary
}

// IsService returns true for service-oriented project types.
func (pt ProjectType) IsService() bool {
	return pt == ProjectTypeWebAPI || pt == ProjectTypeGRPCService ||
		pt == ProjectTypeMicroservice ||
		pt == ProjectTypeDaemon
}

// DefaultCGOEnabled returns true if project type typically requires CGO.
func (pt ProjectType) DefaultCGOEnabled() bool {
	switch pt {
	case ProjectTypeDesktop, ProjectTypeWebAPI, ProjectTypeGRPCService:
		return true
	default:
		return false
	}
}

// DockerSupported returns true if project type supports Docker.
func (pt ProjectType) DockerSupported() bool {
	return !pt.IsMobileRelated() && !pt.IsLibrary()
}

// RecommendedPlatforms returns recommended platforms for project type.
func (pt ProjectType) RecommendedPlatforms() []Platform {
	switch pt {
	case ProjectTypeCLI, ProjectTypeTool:
		// Note: Windows is not included by default because it doesn't support ARM64
		// Add Windows explicitly if needed for AMD64/386 only builds
		return []Platform{PlatformLinux, PlatformDarwin}
	case ProjectTypeDesktop, ProjectTypeDaemon:
		return []Platform{PlatformLinux, PlatformDarwin}
	case ProjectTypeWebAPI, ProjectTypeGRPCService, ProjectTypeMicroservice:
		return []Platform{PlatformLinux, PlatformDarwin}
	case ProjectTypeLibrary:
		// Note: Windows is not included by default because it doesn't support ARM64
		return []Platform{PlatformLinux, PlatformDarwin}
	case ProjectTypeMobile:
		return []Platform{} // Mobile projects handled differently
	case ProjectTypePlugin:
		// Note: Windows is not included by default because it doesn't support ARM64
		return []Platform{PlatformLinux, PlatformDarwin}
	default:
		return []Platform{PlatformLinux}
	}
}

// RecommendedArchitectures returns recommended architectures for project type.
func (pt ProjectType) RecommendedArchitectures() []Architecture {
	switch pt {
	case ProjectTypeCLI, ProjectTypeTool:
		return []Architecture{ArchitectureAMD64, ArchitectureARM64}
	case ProjectTypeDesktop:
		return []Architecture{ArchitectureAMD64}
	case ProjectTypeWebAPI, ProjectTypeGRPCService, ProjectTypeMicroservice:
		return []Architecture{ArchitectureAMD64, ArchitectureARM64}
	case ProjectTypeLibrary:
		// 386 excluded: it is incompatible with the recommended darwin
		// platform, so recommending it would produce configs that fail
		// their own platform/architecture validation.
		return []Architecture{ArchitectureAMD64, ArchitectureARM64}
	case ProjectTypeMobile:
		return []Architecture{ArchitectureARM64}
	default:
		return []Architecture{ArchitectureAMD64}
	}
}

// RequiresGit returns true if project type requires Git.
func (pt ProjectType) RequiresGit() bool {
	return pt != ProjectTypeTool && pt != ProjectTypePlugin
}

// RequiresCI returns true if project type benefits from CI/CD.
func (pt ProjectType) RequiresCI() bool {
	return !pt.IsLibrary() && pt != ProjectTypePlugin
}

// FeatureLevel represents feature maturity levels.
type FeatureLevel string

const (
	FeatureLevelNone       FeatureLevel = "none"
	FeatureLevelBasic      FeatureLevel = "basic"
	FeatureLevelStandard   FeatureLevel = "standard"
	FeatureLevelAdvanced   FeatureLevel = "advanced"
	FeatureLevelEnterprise FeatureLevel = "enterprise"
)

// IsValid returns true if FeatureLevel is valid.
func (fl FeatureLevel) IsValid() bool {
	switch fl {
	case FeatureLevelNone, FeatureLevelBasic, FeatureLevelStandard,
		FeatureLevelAdvanced, FeatureLevelEnterprise:
		return true
	default:
		return false
	}
}

// String returns human-readable display name.
func (fl FeatureLevel) String() string {
	switch fl {
	case FeatureLevelNone:
		return "None"
	case FeatureLevelBasic:
		return "Basic"
	case FeatureLevelStandard:
		return "Standard"
	case FeatureLevelAdvanced:
		return "Advanced"
	case FeatureLevelEnterprise:
		return "Enterprise"
	default:
		return unknownValue
	}
}

// IsEnabled returns true if feature level is enabled.
func (fl FeatureLevel) IsEnabled() bool {
	return fl != FeatureLevelNone
}

// IncludesBasic returns true if feature level includes basic features.
func (fl FeatureLevel) IncludesBasic() bool {
	return fl == FeatureLevelBasic || fl == FeatureLevelStandard || fl == FeatureLevelAdvanced ||
		fl == FeatureLevelEnterprise
}

// IncludesAdvanced returns true if feature level includes advanced features.
func (fl FeatureLevel) IncludesAdvanced() bool {
	return fl == FeatureLevelAdvanced || fl == FeatureLevelEnterprise
}

// IncludesEnterprise returns true if feature level includes enterprise features.
func (fl FeatureLevel) IncludesEnterprise() bool {
	return fl == FeatureLevelEnterprise
}

// GetRecommendedDockerSupport returns recommended Docker support for feature level.
func (fl FeatureLevel) GetRecommendedDockerSupport() DockerSupport {
	switch fl {
	case FeatureLevelNone:
		return DockerSupportNone
	case FeatureLevelBasic:
		return DockerSupportBuild
	case FeatureLevelStandard:
		return DockerSupportBuild
	case FeatureLevelAdvanced:
		return DockerSupportBoth
	case FeatureLevelEnterprise:
		return DockerSupportBoth
	default:
		return DockerSupportNone
	}
}

// GetRecommendedSigningLevel returns recommended signing level for feature level.
func (fl FeatureLevel) GetRecommendedSigningLevel() SigningLevel {
	switch fl {
	case FeatureLevelNone:
		return SigningLevelNone
	case FeatureLevelBasic:
		return SigningLevelBasic
	case FeatureLevelStandard:
		return SigningLevelBasic
	case FeatureLevelAdvanced:
		return SigningLevelAdvanced
	case FeatureLevelEnterprise:
		return SigningLevelEnterprise
	default:
		return SigningLevelNone
	}
}

// GetRecommendedActionLevel returns recommended GitHub Actions level for feature level.
func (fl FeatureLevel) GetRecommendedActionLevel() ActionLevel {
	switch fl {
	case FeatureLevelNone:
		return ActionLevelNone
	case FeatureLevelBasic:
		return ActionLevelBasic
	case FeatureLevelStandard:
		return ActionLevelStandard
	case FeatureLevelAdvanced:
		return ActionLevelAdvanced
	case FeatureLevelEnterprise:
		return ActionLevelEnterprise
	default:
		return ActionLevelNone
	}
}
