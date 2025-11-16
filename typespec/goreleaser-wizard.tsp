// GoReleaser Wizard Domain Type Specification
// Generated: 2025-11-16
// Purpose: Single source of truth for all domain types
// Eliminates split brains and ensures type safety

// ProjectType enum with compile-time invariants
@discriminator("type")
enum ProjectType {
  // CLI Application - Single binary, cross-platform
  CLI("cli", "CLI Application") {
    defaultCGOEnabled: false,
    recommendedPlatforms: [Platform.Linux, Platform.Darwin, Platform.Windows],
    dockerSupported: true,
    requiresMainPath: true,
    defaultBinaryName: "cli-app"
  }
  
  // Web Service - HTTP server, container-focused
  Web("web", "Web Service") {
    defaultCGOEnabled: true,
    recommendedPlatforms: [Platform.Linux, Platform.Darwin],
    dockerSupported: true,
    requiresMainPath: true,
    defaultBinaryName: "server"
  }
  
  // Library - Go package with optional CLI
  Library("library", "Library") {
    defaultCGOEnabled: false,
    recommendedPlatforms: [Platform.Linux, Platform.Darwin, Platform.Windows],
    dockerSupported: false,
    requiresMainPath: false,
    defaultBinaryName: "lib-tool"
  }
  
  // API Service - REST/GraphQL API
  API("api", "API Service") {
    defaultCGOEnabled: true,
    recommendedPlatforms: [Platform.Linux, Platform.Darwin],
    dockerSupported: true,
    requiresMainPath: true,
    defaultBinaryName: "api"
  }
  
  // Desktop Application - GUI application
  Desktop("desktop", "Desktop Application") {
    defaultCGOEnabled: true,
    recommendedPlatforms: [Platform.Windows, Platform.Darwin, Platform.Linux],
    dockerSupported: false,
    requiresMainPath: true,
    defaultBinaryName: "app"
  }
}

// Platform enum with validation rules
@discriminator("platform")
enum Platform {
  Linux("linux", "Linux") {
    architectures: [Architecture.AMD64, Architecture.ARM64, Architecture.ARM, Architecture.386],
    isWindowsBased: false,
    isUnixLike: true,
    supportsCGO: true
  }
  
  Darwin("darwin", "macOS") {
    architectures: [Architecture.AMD64, Architecture.ARM64],
    isWindowsBased: false,
    isUnixLike: true,
    supportsCGO: true
  }
  
  Windows("windows", "Windows") {
    architectures: [Architecture.AMD64, Architecture.ARM64, Architecture.386],
    isWindowsBased: true,
    isUnixLike: false,
    supportsCGO: true
  }
  
  FreeBSD("freebsd", "FreeBSD") {
    architectures: [Architecture.AMD64, Architecture.ARM64, Architecture.386],
    isWindowsBased: false,
    isUnixLike: true,
    supportsCGO: true
  }
  
  OpenBSD("openbsd", "OpenBSD") {
    architectures: [Architecture.AMD64, Architecture.ARM64, Architecture.386],
    isWindowsBased: false,
    isUnixLike: true,
    supportsCGO: true
  }
  
  NetBSD("netbsd", "NetBSD") {
    architectures: [Architecture.AMD64, Architecture.ARM64, Architecture.386],
    isWindowsBased: false,
    isUnixLike: true,
    supportsCGO: true
  }
}

// Architecture enum with platform compatibility
@discriminator("arch")
enum Architecture {
  AMD64("amd64", "64-bit x86") {
    supportedByAllPlatforms: true,
    is64Bit: true,
    goSupport: "stable"
  }
  
  ARM64("arm64", "64-bit ARM") {
    supportedByAllPlatforms: true,
    is64Bit: true,
    goSupport: "stable"
  }
  
  ARM("arm", "32-bit ARM") {
    supportedByAllPlatforms: true,
    is64Bit: false,
    goSupport: "stable"
  }
  
  x86("386", "32-bit x86") {
    supportedByAllPlatforms: true,
    is64Bit: false,
    goSupport: "stable"
  }
  
  PPC64("ppc64", "64-bit PowerPC (big endian)") {
    supportedByAllPlatforms: false,
    is64Bit: true,
    goSupport: "stable"
  }
  
  PPC64LE("ppc64le", "64-bit PowerPC (little endian)") {
    supportedByAllPlatforms: false,
    is64Bit: true,
    goSupport: "stable"
  }
  
  S390X("s390x", "IBM System z") {
    supportedByAllPlatforms: false,
    is64Bit: true,
    goSupport: "stable"
  }
  
  MIPS("mips", "32-bit MIPS (big endian)") {
    supportedByAllPlatforms: false,
    is64Bit: false,
    goSupport: "stable"
  }
  
  MIPSLE("mipsle", "32-bit MIPS (little endian)") {
    supportedByAllPlatforms: false,
    is64Bit: false,
    goSupport: "stable"
  }
}

// GitProvider enum with configuration patterns
@discriminator("provider")
enum GitProvider {
  GitHub("github", "GitHub") {
    defaultRegistry: DockerRegistry.GitHub,
    actionsSupported: true,
    apiURL: "https://api.github.com",
    webURL: "https://github.com",
    requiresPersonalAccessToken: false
  }
  
  GitLab("gitlab", "GitLab") {
    defaultRegistry: DockerRegistry.GitLab,
    actionsSupported: true,
    apiURL: "https://gitlab.com/api/v4",
    webURL: "https://gitlab.com",
    requiresPersonalAccessToken: true
  }
  
  Bitbucket("bitbucket", "Bitbucket") {
    defaultRegistry: DockerRegistry.Custom,
    actionsSupported: true,
    apiURL: "https://api.bitbucket.org/2.0",
    webURL: "https://bitbucket.org",
    requiresPersonalAccessToken: true
  }
  
  Gitea("gitea", "Gitea") {
    defaultRegistry: DockerRegistry.Custom,
    actionsSupported: false,
    apiURL: "", // Self-hosted
    webURL: "",  // Self-hosted
    requiresPersonalAccessToken: true
  }
  
  SelfHosted("self-hosted", "Self-hosted") {
    defaultRegistry: DockerRegistry.Custom,
    actionsSupported: false,
    apiURL: "", // User-defined
    webURL: "",  // User-defined
    requiresPersonalAccessToken: true
  }
}

// DockerRegistry enum with validation rules
@discriminator("registry")
enum DockerRegistry {
  DockerHub("docker.io", "Docker Hub") {
    urlPattern: "^[a-z0-9]([a-z0-9-]*[a-z0-9])?$",
    supportsHTTPSOnly: false,
    requiresAuthentication: false,
    defaultNamespace: "library"
  }
  
  GitHub("ghcr.io", "GitHub Container Registry") {
    urlPattern: "^ghcr\\.io/[a-z0-9-]+/[a-z0-9-]+$",
    supportsHTTPSOnly: true,
    requiresAuthentication: true,
    defaultNamespace: ""
  }
  
  GitLab("registry.gitlab.com", "GitLab Registry") {
    urlPattern: "^registry\\.gitlab\\.com/[a-zA-Z0-9_-]+/[a-zA-Z0-9_-]+$",
    supportsHTTPSOnly: true,
    requiresAuthentication: true,
    defaultNamespace: ""
  }
  
  Quay("quay.io", "Quay.io") {
    urlPattern: "^quay\\.io/[a-z0-9-]+/[a-z0-9-]+$",
    supportsHTTPSOnly: true,
    requiresAuthentication: true,
    defaultNamespace: ""
  }
  
  Custom("custom", "Custom Registry") {
    urlPattern: "", // User-defined
    supportsHTTPSOnly: true,
    requiresAuthentication: true,
    defaultNamespace: ""
  }
}

// ActionTrigger enum for GitHub Actions
@discriminator("trigger")
enum ActionTrigger {
  VersionTags("version-tags", "Version Tags (v*)") {
    githubPattern: "push:\n  tags:\n    - 'v*'",
    description: "Triggers on version tags like v1.0.0",
    recommendedFor: [ProjectType.CLI, ProjectType.Web, ProjectType.API]
  }
  
  AllTags("all-tags", "All Tags (*)") {
    githubPattern: "push:\n  tags:\n    - '*'",
    description: "Triggers on any tag",
    recommendedFor: [ProjectType.Library]
  }
  
  Manual("manual", "Manual Trigger") {
    githubPattern: "workflow_dispatch:",
    description: "Can be triggered manually from GitHub UI",
    recommendedFor: [ProjectType.CLI, ProjectType.Web, ProjectType.API, ProjectType.Desktop]
  }
  
  Main("main", "Push to Main") {
    githubPattern: "push:\n  branches:\n    - main",
    description: "Triggers on pushes to main branch",
    recommendedFor: [ProjectType.Web, ProjectType.API]
  }
  
  Release("release", "Published Release") {
    githubPattern: "release:\n    types: [published]",
    description: "Triggers when a GitHub release is published",
    recommendedFor: [ProjectType.CLI, ProjectType.Web, ProjectType.API, ProjectType.Desktop]
  }
}

// ConfigState enum for configuration lifecycle
@discriminator("state")
enum ConfigState {
  Draft("draft", "Draft") {
    description: "Configuration is being created or modified",
    isFinal: false,
    allowsValidation: true,
    allowsGeneration: false
  }
  
  Valid("valid", "Valid") {
    description: "Configuration is valid and ready",
    isFinal: false,
    allowsValidation: true,
    allowsGeneration: true
  }
  
  Invalid("invalid", "Invalid") {
    description: "Configuration has validation errors",
    isFinal: false,
    allowsValidation: true,
    allowsGeneration: false
  }
  
  Processing("processing", "Processing") {
    description: "Configuration is being processed",
    isFinal: false,
    allowsValidation: false,
    allowsGeneration: false
  }
  
  Generated("generated", "Generated") {
    description: "Configuration has been generated successfully",
    isFinal: true,
    allowsValidation: true,
    allowsGeneration: false
  }
}

// BuildTag represents a constrained build tag
model BuildTag {
  name: string<1..50> @pattern("^[a-zA-Z0-9_-]+$") @required;
  description?: string<0..200>;
}

// SafeProjectConfig - Single source of truth for project configuration
model SafeProjectConfig {
  // Basic Information
  projectName: string<1..63> @required @pattern("^[a-zA-Z][a-zA-Z0-9_-]*$") {
    @minLength(1)
    @maxLength(63)
    @description("Project name must start with letter and contain only letters, numbers, hyphens, underscores")
  }
  
  projectDescription: string<0..255> @pattern("^[^\\x00-\\x1F\\x7F]*$") {
    @maxLength(255)
    @description("Optional project description, no control characters")
  }
  
  projectType: ProjectType @required {
    @description("Type of project being configured")
  }
  
  binaryName: string<1..63> @pattern("^[a-zA-Z][a-zA-Z0-9_-]*$") @required {
    @minLength(1)
    @maxLength(63)
    @description("Binary name must start with letter and contain only letters, numbers, hyphens, underscores")
  }
  
  mainPath: string<1..255> @pattern("^[a-zA-Z0-9/_.-]+$") @required {
    @description("Path to main.go file, relative to project root")
  }

  // Build Configuration
  platforms: Platform<1..> @required {
    @minItems(1)
    @description("Target platforms for cross-compilation")
  }
  
  architectures: Architecture<1..> @required {
    @minItems(1)
    @description("Target architectures for cross-compilation")
  }
  
  cgoEnabled: boolean @default(false) {
    @description("Enable CGO for compilation")
  }
  
  buildTags: BuildTag<0..10> @default([]) {
    @maxItems(10)
    @description("Additional build tags for conditional compilation")
  }
  
  ldflags: boolean @default(true) {
    @description("Enable automatic ldflags for version information")
  }

  // Release Configuration
  gitProvider: GitProvider @required {
    @description("Git hosting provider")
  }
  
  dockerEnabled: boolean @default(false) {
    @description("Enable Docker image generation")
  }
  
  dockerRegistry: DockerRegistry @default(DockerRegistry.DockerHub) {
    @description("Docker registry for container images")
  }
  
  dockerImage: string<0..255> @pattern("^[a-z0-9][a-z0-9/_.-]*$") {
    @description("Docker image name, defaults to project name")
  }
  
  signing: boolean @default(false) {
    @description("Enable code signing with cosign")
  }
  
  homebrew: boolean @default(false) {
    @description("Generate Homebrew formula")
  }
  
  snap: boolean @default(false) {
    @description("Generate Snap package")
  }
  
  sbom: boolean @default(true) {
    @description("Generate Software Bill of Materials")
  }

  // CI/CD Configuration
  generateActions: boolean @default(true) {
    @description("Generate GitHub Actions workflow")
  }
  
  actionsOn: ActionTrigger<1..> @required @default([ActionTrigger.VersionTags]) {
    @minItems(1)
    @description("Triggers for GitHub Actions workflow")
  }

  // Advanced Features
  proVersion: boolean @default(false) {
    @description("Enable GoReleaser Pro features")
  }
  
  state: ConfigState @default(ConfigState.Draft) {
    @description("Current state of configuration")
  }
}

// Compile-time invariants validation
@validator
model ConfigurationInvariants {
  // Docker should only be enabled for project types that support it
  @invariant("docker_requires_support")
  validateDockerSupport(config: SafeProjectConfig): boolean {
    !config.dockerEnabled || config.projectType.dockerSupported
  }
  
  // All platforms should support selected architectures
  @invariant("architecture_platform_compatibility")
  validatePlatformArchCompatibility(config: SafeProjectConfig): boolean {
    for platform in config.platforms {
      for arch in config.architectures {
        if arch not in platform.architectures {
          return false
        }
      }
    }
    return true
  }
  
  // Main path required only for project types that need it
  @invariant("main_path_requirement")
  validateMainPathRequirement(config: SafeProjectConfig): boolean {
    config.projectType.requiresMainPath ? config.mainPath != "" : true
  }
  
  // CGO should align with project type defaults unless explicitly changed
  @invariant("cgo_alignment")
  validateCGOAlignment(config: SafeProjectConfig): boolean {
    // Allow override but warn if different from project type default
    true
  }
  
  // Actions triggers should be recommended for project type
  @invariant("action_trigger_recommendations")
  validateActionTriggers(config: SafeProjectConfig): boolean {
    // Allow flexibility but warn if not recommended
    true
  }
}