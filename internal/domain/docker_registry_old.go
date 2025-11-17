package domain

import (
	"regexp"
)

// DockerRegistry represents Docker registry types
// Generated from TypeSpec specification - DO NOT MODIFY MANUALLY
type DockerRegistry string

const (
	DockerRegistryDockerHub DockerRegistry = "docker.io"  // Docker Hub
	DockerRegistryGitHub    DockerRegistry = "ghcr.io"    // GitHub Container Registry
	DockerRegistryGitLab    DockerRegistry = "registry.gitlab.com" // GitLab Registry
	DockerRegistryQuay      DockerRegistry = "quay.io"    // Quay.io
	DockerRegistryCustom    DockerRegistry = "custom"      // Custom Registry
)

// DockerRegistry metadata - generated from TypeSpec invariants
type dockerRegistryMetadata struct {
	urlPattern           string
	supportsHTTPSOnly    bool
	requiresAuthentication bool
	defaultNamespace     string
}

var dockerRegistryMetadata = map[DockerRegistry]dockerRegistryMetadata{
	DockerRegistryDockerHub: {
		urlPattern:           "^[a-z0-9]([a-z0-9-]*[a-z0-9])?$",
		supportsHTTPSOnly:    false,
		requiresAuthentication: false,
		defaultNamespace:     "library",
	},
	DockerRegistryGitHub: {
		urlPattern:           "^ghcr\\.io/[a-z0-9-]+/[a-z0-9-]+$",
		supportsHTTPSOnly:    true,
		requiresAuthentication: true,
		defaultNamespace:     "",
	},
	DockerRegistryGitLab: {
		urlPattern:           "^registry\\.gitlab\\.com/[a-zA-Z0-9_-]+/[a-zA-Z0-9_-]+$",
		supportsHTTPSOnly:    true,
		requiresAuthentication: true,
		defaultNamespace:     "",
	},
	DockerRegistryQuay: {
		urlPattern:           "^quay\\.io/[a-z0-9-]+/[a-z0-9-]+$",
		supportsHTTPSOnly:    true,
		requiresAuthentication: true,
		defaultNamespace:     "",
	},
	DockerRegistryCustom: {
		urlPattern:           "", // User-defined
		supportsHTTPSOnly:    true,
		requiresAuthentication: true,
		defaultNamespace:     "",
	},
}

// IsValid returns true if DockerRegistry is valid
func (dr DockerRegistry) IsValid() bool {
	_, exists := dockerRegistryMetadata[dr]
	return exists
}

// String returns human-readable display name
func (dr DockerRegistry) String() string {
	switch dr {
	case DockerRegistryDockerHub:
		return "Docker Hub"
	case DockerRegistryGitHub:
		return "GitHub Container Registry"
	case DockerRegistryGitLab:
		return "GitLab Registry"
	case DockerRegistryQuay:
		return "Quay.io"
	case DockerRegistryCustom:
		return "Custom Registry"
	default:
		return string(dr)
	}
}

// URLPattern returns the URL validation pattern for this registry
func (dr DockerRegistry) URLPattern() string {
	if meta, exists := dockerRegistryMetadata[dr]; exists {
		return meta.urlPattern
	}
	return ""
}

// SupportsHTTPSOnly returns true if registry only supports HTTPS
func (dr DockerRegistry) SupportsHTTPSOnly() bool {
	if meta, exists := dockerRegistryMetadata[dr]; exists {
		return meta.supportsHTTPSOnly
	}
	return true
}

// RequiresAuthentication returns true if registry requires authentication
func (dr DockerRegistry) RequiresAuthentication() bool {
	if meta, exists := dockerRegistryMetadata[dr]; exists {
		return meta.requiresAuthentication
	}
	return true
}

// DefaultNamespace returns the default namespace for this registry
func (dr DockerRegistry) DefaultNamespace() string {
	if meta, exists := dockerRegistryMetadata[dr]; exists {
		return meta.defaultNamespace
	}
	return ""
}

// ValidateDockerRegistryURL validates a Docker registry URL
func ValidateDockerRegistryURL(registry DockerRegistry, url string) error {
	if !registry.IsValid() {
		return fmt.Errorf("invalid Docker registry: %s", registry)
	}
	
	if url == "" {
		return fmt.Errorf("Docker registry URL cannot be empty")
	}
	
	url = strings.TrimSpace(url)
	
	if registry == DockerRegistryDockerHub {
		// Docker Hub allows simple usernames or full registry URLs
		if !strings.Contains(url, "docker.io") && !strings.Contains(url, "/") {
			return fmt.Errorf("Docker Hub registry should include docker.io or be a valid username")
		}
	} else if registry == DockerRegistryGitHub {
		if !strings.Contains(url, "ghcr.io") {
			return fmt.Errorf("GitHub Container Registry should include ghcr.io")
		}
	}
	
	// Validate against registry pattern if available
	pattern := registry.URLPattern()
	if pattern != "" {
		if matched, err := regexp.MatchString(pattern, url); err != nil {
			return fmt.Errorf("invalid URL pattern for registry %s: %v", registry, err)
		} else if !matched {
			return fmt.Errorf("URL '%s' does not match expected pattern for registry %s", url, registry)
		}
	}
	
	return nil
}

// ValidateDockerImageName validates a Docker image name
func ValidateDockerImageName(name string) error {
	if len(name) == 0 {
		return nil // Empty is allowed, will default to project name
	}
	
	if len(name) > 255 {
		return fmt.Errorf("Docker image name must be 255 characters or less")
	}
	
	// Docker image name pattern: lowercase, numbers, dots, hyphens, underscores, forward slashes
	pattern := `^[a-z0-9][a-z0-9/_.-]*$`
	matched, err := regexp.MatchString(pattern, name)
	if err != nil {
		return fmt.Errorf("invalid image name pattern: %v", err)
	}
	if !matched {
		return fmt.Errorf("Docker image name '%s' must start with lowercase letter/number and contain only lowercase letters, numbers, dots, hyphens, underscores, and forward slashes", name)
	}
	
	return nil
}

// ValidateDockerRegistry validates a Docker registry
func ValidateDockerRegistry(registry DockerRegistry) error {
	if !registry.IsValid() {
		return fmt.Errorf("invalid Docker registry: %s", registry)
	}
	return nil
}