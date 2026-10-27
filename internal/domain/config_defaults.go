package domain

import (
	"fmt"
	"slices"
	"strings"
)

// ApplyDefaults applies smart defaults based on project type and context.
func (spc *SafeProjectConfig) ApplyDefaults() {
	// Set default feature level first (required by other defaults)
	if spc.FeatureLevel == "" {
		spc.FeatureLevel = FeatureLevelStandard
	}

	// Apply project type-specific defaults
	spc.applyProjectTypeDefaults()

	// Apply feature level defaults
	spc.applyFeatureLevelDefaults()

	// Apply general defaults
	spc.applyGeneralDefaults()
}

// applyProjectTypeDefaults applies defaults based on project type.
func (spc *SafeProjectConfig) applyProjectTypeDefaults() {
	// Set CGO status based on project type
	if spc.CGOStatus == "" {
		if spc.ProjectType.DefaultCGOEnabled() {
			spc.CGOStatus = CGOStatusEnabled
		} else {
			spc.CGOStatus = CGOStatusDisabled
		}
	}

	// Set platforms based on project type
	if len(spc.Platforms) == 0 {
		spc.Platforms = spc.ProjectType.RecommendedPlatforms()
	}

	// Set architectures based on project type
	if len(spc.Architectures) == 0 {
		spc.Architectures = spc.ProjectType.RecommendedArchitectures()
	}

	// Set Docker support based on project type
	if spc.DockerSupport == "" {
		if spc.ProjectType.DockerSupported() {
			spc.DockerSupport = DockerSupportBuild
		} else {
			spc.DockerSupport = DockerSupportNone
		}
	}

	// Set CI/CD requirements based on project type
	if spc.ActionLevel == "" && spc.ProjectType.RequiresCI() {
		spc.ActionLevel = ActionLevelBasic
	}
}

// applyFeatureLevelDefaults applies defaults based on feature level.
func (spc *SafeProjectConfig) applyFeatureLevelDefaults() {
	// Set Docker support based on feature level
	if spc.DockerSupport == "" {
		spc.DockerSupport = spc.FeatureLevel.GetRecommendedDockerSupport()
	}

	// Set signing level based on feature level
	if spc.SigningLevel == "" {
		spc.SigningLevel = spc.FeatureLevel.GetRecommendedSigningLevel()
	}

	// Set action level based on feature level
	if spc.ActionLevel == "" {
		spc.ActionLevel = spc.FeatureLevel.GetRecommendedActionLevel()
	}

	// Set feature-specific defaults
	if spc.FeatureLevel.IncludesBasic() {
		spc.LDFlags = true
	}

	if spc.FeatureLevel.IncludesAdvanced() {
		spc.Homebrew = true
		spc.Snap = true
	}

	if spc.FeatureLevel.IncludesEnterprise() {
		spc.SBOM = true
	}
}

// applyGeneralDefaults applies general defaults.
func (spc *SafeProjectConfig) applyGeneralDefaults() {
	// Set Git provider if not specified
	if spc.GitProvider == "" {
		spc.GitProvider = GetRecommendedGitProvider()
	}

	// Set Docker registry if not specified
	if spc.DockerRegistry == "" && spc.DockerSupport.IsEnabled() {
		spc.DockerRegistry = spc.GitProvider.DefaultRegistry()
	}

	// Set default state
	if spc.State == "" {
		spc.State = ConfigStateDraft
	}

	// Set default main path
	if spc.MainPath == "" {
		spc.MainPath = "."
	}

	// Set default binary name if not specified
	if spc.BinaryName == "" && spc.ProjectName != "" {
		spc.BinaryName = spc.ProjectName
	}

	// Set default action triggers if actions are enabled
	if spc.ActionLevel.IsEnabled() && len(spc.ActionsOn) == 0 {
		spc.ActionsOn = spc.ActionLevel.GetRecommendedTriggers()
	}

	// Set default CI provider if not specified
	if spc.CIProvider == "" {
		spc.CIProvider = GetRecommendedCIProvider()
	}

	// Set default build tags based on CGO status
	if spc.CGOStatus.IsDisabled() &&
		!slices.Contains(spc.BuildTags, CreateBuildTag("pure", "Pure Go compilation")) {
		spc.BuildTags = append(spc.BuildTags, CreateBuildTag("pure", "Pure Go compilation"))
	}
}

// GetDockerImageName returns Docker image name based on configuration.
func (spc *SafeProjectConfig) GetDockerImageName() string {
	if spc.DockerImage != "" {
		return spc.DockerImage
	}

	// Generate default image name
	if spc.ProjectName != "" {
		// Convert to lowercase and replace spaces with hyphens
		imageName := strings.ToLower(spc.ProjectName)
		imageName = strings.ReplaceAll(imageName, " ", "-")

		return imageName
	}

	return "app"
}

// GetDockerRegistryURL returns full Docker registry URL.
func (spc *SafeProjectConfig) GetDockerRegistryURL() string {
	baseURL := spc.DockerRegistry.GetURL()

	// For GitHub, append username
	if spc.DockerRegistry == DockerRegistryGitHub && spc.GetDockerOwner() != "" {
		return fmt.Sprintf("%s/%s", baseURL, spc.GetDockerOwner())
	}

	// For Azure, handle registry differently
	if spc.DockerRegistry == DockerRegistryCustom {
		// Azure registries end with .azurecr.io
		if strings.Contains(baseURL, "azurecr.io") {
			return baseURL
		}
	}

	return baseURL
}

// GetDockerOwner returns Docker image owner.
func (spc *SafeProjectConfig) GetDockerOwner() string {
	// Try to get from git remote
	if owner := GetGitHubOwner(spc.ProjectName); owner != "owner" {
		return owner
	}

	// Fallback to project name
	if spc.ProjectName != "" {
		return strings.ToLower(spc.ProjectName)
	}

	return "user"
}

// GetAzureContainerRegistry returns Azure container registry name.
func (spc *SafeProjectConfig) GetAzureContainerRegistry() string {
	// This would typically come from configuration or environment
	// For now, return a default
	return "myregistry"
}

// ShouldGenerateActionsFiles returns true if actions files should be generated.
func (spc *SafeProjectConfig) ShouldGenerateActionsFiles() bool {
	return spc.ActionLevel.IsEnabled() &&
		spc.CIProvider.GeneratesGitHubWorkflows() &&
		spc.GitProvider.ActionsSupported() &&
		len(spc.ActionsOn) > 0
}

// ShouldGenerateForgejoWorkflows returns true if Forgejo CI workflow files
// should be generated. Forgejo Actions do not depend on GitProvider (the
// Forgejo instance is the CI host regardless of where the public mirror
// lives).
func (spc *SafeProjectConfig) ShouldGenerateForgejoWorkflows() bool {
	return spc.ActionLevel.IsEnabled() &&
		spc.CIProvider.GeneratesForgejoWorkflows() &&
		len(spc.ActionsOn) > 0
}

// ShouldGenerateDockerFiles returns true if Docker files should be generated.
func (spc *SafeProjectConfig) ShouldGenerateDockerFiles() bool {
	return spc.DockerSupport.IsBuildEnabled() &&
		spc.ProjectType.DockerSupported()
}

// GetDockerEnabled returns true if Docker is fully enabled.
func (spc *SafeProjectConfig) GetDockerEnabled() bool {
	return spc.DockerSupport.IsEnabled()
}

// GetSigningEnabled returns true if signing is enabled.
func (spc *SafeProjectConfig) GetSigningEnabled() bool {
	return spc.SigningLevel.IsEnabled()
}

// GetCGOEnabled returns true if CGO is enabled.
func (spc *SafeProjectConfig) GetCGOEnabled() bool {
	return spc.CGOStatus.IsEnabled()
}

// GetPlatformCount returns number of platforms.
func (spc *SafeProjectConfig) GetPlatformCount() int {
	return len(spc.Platforms)
}

// GetArchitectureCount returns number of architectures.
func (spc *SafeProjectConfig) GetArchitectureCount() int {
	return len(spc.Architectures)
}

// GetBuildTagCount returns number of build tags.
func (spc *SafeProjectConfig) GetBuildTagCount() int {
	return len(spc.BuildTags)
}

// GetActionTriggerCount returns number of action triggers.
func (spc *SafeProjectConfig) GetActionTriggerCount() int {
	return len(spc.ActionsOn)
}

// IsDesktopApplication returns true for desktop applications.
func (spc *SafeProjectConfig) IsDesktopApplication() bool {
	return spc.ProjectType.IsDesktopRelated()
}

// IsWebApplication returns true for web applications.
func (spc *SafeProjectConfig) IsWebApplication() bool {
	return spc.ProjectType.IsWebRelated()
}

// IsMobileApplication returns true for mobile applications.
func (spc *SafeProjectConfig) IsMobileApplication() bool {
	return spc.ProjectType.IsMobileRelated()
}

// IsLibrary returns true for libraries.
func (spc *SafeProjectConfig) IsLibrary() bool {
	return spc.ProjectType.IsLibrary()
}

// IsService returns true for services.
func (spc *SafeProjectConfig) IsService() bool {
	return spc.ProjectType.IsService()
}

// IsProductionReady returns true if configuration is production-ready.
func (spc *SafeProjectConfig) IsProductionReady() bool {
	return spc.FeatureLevel.IncludesAdvanced() &&
		spc.SigningLevel.IsEnabled() &&
		spc.DockerSupport.IsDeployEnabled() &&
		spc.ActionLevel.IsProductionReady()
}

// IsMinimal returns true if configuration is minimal.
func (spc *SafeProjectConfig) IsMinimal() bool {
	return spc.FeatureLevel == FeatureLevelNone ||
		spc.FeatureLevel == FeatureLevelBasic
}

// RequiresCrossCompilation returns true if cross-compilation is required.
func (spc *SafeProjectConfig) RequiresCrossCompilation() bool {
	if len(spc.Platforms) <= 1 || len(spc.Architectures) <= 1 {
		return false
	}

	// Check if we need to compile for different platforms/archs
	for _, platform := range spc.Platforms {
		for _, arch := range spc.Architectures {
			if !isNativePlatform(platform, arch) {
				return true
			}
		}
	}

	return false
}

// isNativePlatform returns true if platform/arch is native to current system.
func isNativePlatform(platform Platform, arch Architecture) bool {
	// This is a simplified check - in reality would need to detect current platform
	return (platform == PlatformLinux && arch == ArchitectureAMD64) ||
		(platform == PlatformDarwin && arch == ArchitectureAMD64) ||
		(platform == PlatformWindows && arch == ArchitectureAMD64)
}

// GetBuildComplexity returns build complexity score.
func (spc *SafeProjectConfig) GetBuildComplexity() int {
	complexity := 0

	// Base complexity from platforms and architectures
	complexity += len(spc.Platforms) * len(spc.Architectures)

	// Add complexity from CGO
	if spc.CGOStatus.IsEnabled() {
		complexity += complexityCGOEnabled
	}

	// Add complexity from build tags
	complexity += len(spc.BuildTags) * 2

	// Add complexity from Docker
	if spc.DockerSupport.IsEnabled() {
		complexity += complexityDockerEnabled
	}

	// Add complexity from signing
	if spc.SigningLevel.IsEnabled() {
		complexity += len(spc.SigningLevel.GetRequiredTools()) * complexitySigningTool
	}

	return complexity
}

// GetEstimatedBuildTime returns estimated build time in seconds.
func (spc *SafeProjectConfig) GetEstimatedBuildTime() int {
	baseTime := baseBuildTimeSeconds // Base build time in seconds

	// Multiply by number of platform/arch combinations
	combos := len(spc.Platforms) * len(spc.Architectures)
	buildTime := baseTime * combos

	// Add time for CGO
	if spc.CGOStatus.IsEnabled() {
		buildTime += cgoBuildTimeAdditional * combos
	}

	// Add time for Docker
	if spc.DockerSupport.IsEnabled() {
		buildTime += dockerBuildTimeAdditional // Additional time for Docker builds
	}

	// Add time for signing
	if spc.SigningLevel.IsEnabled() {
		buildTime += signingBuildTimeAdditional * combos
	}

	return buildTime
}

// GetDependencyCount returns estimated number of external dependencies.
func (spc *SafeProjectConfig) GetDependencyCount() int {
	deps := baseDependencies // Base Go dependencies

	// Add dependencies for features
	if spc.DockerSupport.IsEnabled() {
		deps += dockerDependency // Docker
	}

	if spc.SigningLevel.RequiresCosign() {
		deps += 1 // Cosign
	}

	if spc.SigningLevel.RequiresKeyManagement() {
		deps += 1 // GPG
	}

	if spc.Homebrew {
		deps += 1 // Homebrew tools
	}

	if spc.Snap {
		deps += 1 // Snap tools
	}

	return deps
}
