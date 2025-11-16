package domain

import (
	"fmt"
	"strings"
)

// ProjectType represents the type of project being configured
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
type projectTypeMetadata struct {
	defaultCGOEnabled     bool
	recommendedPlatforms  []Platform
	dockerSupported      bool
	requiresMainPath     bool
	defaultBinaryName    string
}

var projectTypeMetadata = map[ProjectType]projectTypeMetadata{
	ProjectTypeCLI: {
		defaultCGOEnabled:    false,
		recommendedPlatforms: []Platform{PlatformLinux, PlatformDarwin, PlatformWindows},
		dockerSupported:     true,
		requiresMainPath:    true,
		defaultBinaryName:   "cli-app",
	},
	ProjectTypeWeb: {
		defaultCGOEnabled:    true,
		recommendedPlatforms: []Platform{PlatformLinux, PlatformDarwin},
		dockerSupported:     true,
		requiresMainPath:    true,
		defaultBinaryName:   "server",
	},
	ProjectTypeLibrary: {
		defaultCGOEnabled:    false,
		recommendedPlatforms: []Platform{PlatformLinux, PlatformDarwin, PlatformWindows},
		dockerSupported:     false,
		requiresMainPath:    false,
		defaultBinaryName:   "lib-tool",
	},
	ProjectTypeAPI: {
		defaultCGOEnabled:    true,
		recommendedPlatforms: []Platform{PlatformLinux, PlatformDarwin},
		dockerSupported:     true,
		requiresMainPath:    true,
		defaultBinaryName:   "api",
	},
	ProjectTypeDesktop: {
		defaultCGOEnabled:    true,
		recommendedPlatforms: []Platform{PlatformWindows, PlatformDarwin, PlatformLinux},
		dockerSupported:     false,
		requiresMainPath:    true,
		defaultBinaryName:   "app",
	},
}

// IsValid returns true if the ProjectType is valid
func (pt ProjectType) IsValid() bool {
	_, exists := projectTypeMetadata[pt]
	return exists
}

// String returns the human-readable display name
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
	if meta, exists := projectTypeMetadata[pt]; exists {
		return meta.defaultCGOEnabled
	}
	return false
}

// RecommendedPlatforms returns the recommended platforms for this project type
func (pt ProjectType) RecommendedPlatforms() []Platform {
	if meta, exists := projectTypeMetadata[pt]; exists {
		return meta.recommendedPlatforms
	}
	return []Platform{PlatformLinux, PlatformDarwin, PlatformWindows}
}

// DockerSupported returns true if Docker is supported for this project type
func (pt ProjectType) DockerSupported() bool {
	if meta, exists := projectTypeMetadata[pt]; exists {
		return meta.dockerSupported
	}
	return false
}

// RequiresMainPath returns true if main path is required for this project type
func (pt ProjectType) RequiresMainPath() bool {
	if meta, exists := projectTypeMetadata[pt]; exists {
		return meta.requiresMainPath
	}
	return true
}

// DefaultBinaryName returns the default binary name for this project type
func (pt ProjectType) DefaultBinaryName() string {
	if meta, exists := projectTypeMetadata[pt]; exists {
		return meta.defaultBinaryName
	}
	return "app"
}

// ValidateProjectName validates a project name according to TypeSpec rules
func ValidateProjectName(name string) error {
	if len(name) < 1 || len(name) > 63 {
		return fmt.Errorf("project name must be between 1 and 63 characters")
	}
	
	if !isAlphaNumeric(name[0]) {
		return fmt.Errorf("project name must start with a letter")
	}
	
	for _, r := range name {
		if !isAlphaNumeric(r) && r != '_' && r != '-' {
			return fmt.Errorf("project name can only contain letters, numbers, hyphens, and underscores")
		}
	}
	
	return nil
}

// ValidateBinaryName validates a binary name according to TypeSpec rules
func ValidateBinaryName(name string) error {
	if len(name) < 1 || len(name) > 63 {
		return fmt.Errorf("binary name must be between 1 and 63 characters")
	}
	
	if !isAlphaNumeric(name[0]) {
		return fmt.Errorf("binary name must start with a letter")
	}
	
	for _, r := range name {
		if !isAlphaNumeric(r) && r != '_' && r != '-' {
			return fmt.Errorf("binary name can only contain letters, numbers, hyphens, and underscores")
		}
	}
	
	// Check for reserved Windows names
	reservedNames := []string{"con", "prn", "aux", "nul", "com1", "com2", "com3", "com4", "com5", "com6", "com7", "com8", "com9", "lpt1", "lpt2", "lpt3", "lpt4", "lpt5", "lpt6", "lpt7", "lpt8", "lpt9"}
	lowerName := strings.ToLower(name)
	for _, reserved := range reservedNames {
		if lowerName == reserved {
			return fmt.Errorf("binary name '%s' is reserved", name)
		}
	}
	
	return nil
}

// ValidateMainPath validates a main path according to TypeSpec rules
func ValidateMainPath(path string) error {
	if len(path) == 0 || len(path) > 255 {
		return fmt.Errorf("main path must be between 1 and 255 characters")
	}
	
	for _, r := range path {
		if !isAlphaNumeric(r) && r != '/' && r != '_' && r != '.' && r != '-' {
			return fmt.Errorf("main path can only contain letters, numbers, slashes, underscores, dots, and hyphens")
		}
	}
	
	// Prevent path traversal attempts
	if strings.Contains(path, "..") {
		return fmt.Errorf("main path cannot contain parent directory references")
	}
	
	return nil
}

// ValidateProjectDescription validates a project description according to TypeSpec rules
func ValidateProjectDescription(desc string) error {
	if len(desc) > 255 {
		return fmt.Errorf("project description must be 255 characters or less")
	}
	
	// Check for control characters
	for _, r := range desc {
		if r < 32 || r == 127 {
			return fmt.Errorf("project description cannot contain control characters")
		}
	}
	
	return nil
}

// Helper function to check if a rune is alphanumeric
func isAlphaNumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}