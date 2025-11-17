package domain

import (
	"fmt"
)

// ProjectType represents type of project being configured
// Generated from TypeSpec specification - DO NOT MODIFY MANUALLY
type ProjectType string

const (
	ProjectTypeCLI     ProjectType = "cli"     // CLI Application - Single binary, cross-platform
	ProjectTypeWeb     ProjectType = "web"     // Web Service - HTTP server, container-focused
	ProjectTypeLibrary ProjectType = "library" // Library - Go package with optional CLI
	ProjectTypeAPI     ProjectType = "api"     // API Service - REST/GraphQL API
	ProjectTypeDesktop ProjectType = "desktop" // Desktop Application - GUI application
)

// ProjectType metadata - generated from TypeSpec invariants
type projectTypeMeta struct {
	defaultCGOEnabled     bool
	recommendedPlatforms  []Platform
	dockerSupported      bool
	requiresMainPath     bool
	defaultBinaryName    string
}

var projectTypeMetaMap = map[ProjectType]projectTypeMeta{
	ProjectTypeCLI: {
		defaultCGOEnabled:    false,
		recommendedPlatforms: []Platform{PlatformLinux, PlatformDarwin, PlatformWindows},
		dockerSupported:     true,
		requiresMainPath:    true,
		defaultBinaryName:   "cli-app",
	},
	ProjectTypeWeb: {
		defaultCGOEnabled:    true,
		recommendedPlatforms: []Platform{PlatformLinux},
		dockerSupported:     true,
		requiresMainPath:    true,
		defaultBinaryName:   "web-service",
	},
	ProjectTypeLibrary: {
		defaultCGOEnabled:    false,
		recommendedPlatforms: []Platform{PlatformLinux, PlatformDarwin, PlatformWindows},
		dockerSupported:     false,
		requiresMainPath:    false,
		defaultBinaryName:   "library",
	},
	ProjectTypeAPI: {
		defaultCGOEnabled:    true,
		recommendedPlatforms: []Platform{PlatformLinux},
		dockerSupported:     true,
		requiresMainPath:    true,
		defaultBinaryName:   "api-server",
	},
	ProjectTypeDesktop: {
		defaultCGOEnabled:    true,
		recommendedPlatforms: []Platform{PlatformDarwin, PlatformWindows},
		dockerSupported:     false,
		requiresMainPath:    true,
		defaultBinaryName:   "desktop-app",
	},
}

// IsValid returns true if ProjectType is valid
func (pt ProjectType) IsValid() bool {
	_, exists := projectTypeMetaMap[pt]
	return exists
}

// String returns human-readable display name
func (pt ProjectType) String() string {
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

// DefaultCGOEnabled returns the default CGO setting for this project type
func (pt ProjectType) DefaultCGOEnabled() bool {
	if meta, exists := projectTypeMetaMap[pt]; exists {
		return meta.defaultCGOEnabled
	}
	return false
}

// RecommendedPlatforms returns the recommended platforms for this project type
func (pt ProjectType) RecommendedPlatforms() []Platform {
	if meta, exists := projectTypeMetaMap[pt]; exists {
		return meta.recommendedPlatforms
	}
	return []Platform{PlatformLinux, PlatformDarwin, PlatformWindows}
}

// DockerSupported returns true if Docker is supported for this project type
func (pt ProjectType) DockerSupported() bool {
	if meta, exists := projectTypeMetaMap[pt]; exists {
		return meta.dockerSupported
	}
	return false
}

// RequiresMainPath returns true if main path is required for this project type
func (pt ProjectType) RequiresMainPath() bool {
	if meta, exists := projectTypeMetaMap[pt]; exists {
		return meta.requiresMainPath
	}
	return true
}

// DefaultBinaryName returns the default binary name for this project type
func (pt ProjectType) DefaultBinaryName() string {
	if meta, exists := projectTypeMetaMap[pt]; exists {
		return meta.defaultBinaryName
	}
	return "app"
}

// ValidateProjectType validates a project type
func ValidateProjectType(pt ProjectType) error {
	if !pt.IsValid() {
		return fmt.Errorf("invalid project type: %s", pt)
	}
	return nil
}

// GetAllProjectTypes returns all available project types
func GetAllProjectTypes() []ProjectType {
	return []ProjectType{
		ProjectTypeCLI, ProjectTypeWeb, ProjectTypeLibrary,
		ProjectTypeAPI, ProjectTypeDesktop,
	}
}

// GetRecommendedProjectType returns the recommended project type (CLI)
func GetRecommendedProjectType() ProjectType {
	return ProjectTypeCLI
}