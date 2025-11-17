package domain

// GitProvider represents git hosting providers
// Generated from TypeSpec specification - DO NOT MODIFY MANUALLY
type GitProvider string

const (
	GitProviderGitHub     GitProvider = "github"     // GitHub
	GitProviderGitLab     GitProvider = "gitlab"     // GitLab
	GitProviderBitbucket  GitProvider = "bitbucket"  // Bitbucket
	GitProviderGitea      GitProvider = "gitea"      // Gitea
	GitProviderSelfHosted GitProvider = "self-hosted" // Self-hosted
)

// GitProvider metadata - generated from TypeSpec invariants
type gitProviderMetadata struct {
	defaultRegistry           DockerRegistry
	actionsSupported         bool
	apiURL                  string
	webURL                  string
	requiresPersonalAccessToken bool
}

var gitProviderMetadata = map[GitProvider]gitProviderMetadata{
	GitProviderGitHub: {
		defaultRegistry:           DockerRegistryGitHub,
		actionsSupported:         true,
		apiURL:                  "https://api.github.com",
		webURL:                  "https://github.com",
		requiresPersonalAccessToken: false,
	},
	GitProviderGitLab: {
		defaultRegistry:           DockerRegistryGitLab,
		actionsSupported:         true,
		apiURL:                  "https://gitlab.com/api/v4",
		webURL:                  "https://gitlab.com",
		requiresPersonalAccessToken: true,
	},
	GitProviderBitbucket: {
		defaultRegistry:           DockerRegistryCustom,
		actionsSupported:         true,
		apiURL:                  "https://api.bitbucket.org/2.0",
		webURL:                  "https://bitbucket.org",
		requiresPersonalAccessToken: true,
	},
	GitProviderGitea: {
		defaultRegistry:           DockerRegistryCustom,
		actionsSupported:         false,
		apiURL:                  "", // Self-hosted
		webURL:                  "", // Self-hosted
		requiresPersonalAccessToken: true,
	},
	GitProviderSelfHosted: {
		defaultRegistry:           DockerRegistryCustom,
		actionsSupported:         false,
		apiURL:                  "", // User-defined
		webURL:                  "", // User-defined
		requiresPersonalAccessToken: true,
	},
}

// IsValid returns true if GitProvider is valid
func (gp GitProvider) IsValid() bool {
	_, exists := gitProviderMetadata[gp]
	return exists
}

// String returns human-readable display name
func (gp GitProvider) String() string {
	switch gp {
	case GitProviderGitHub:
		return "GitHub"
	case GitProviderGitLab:
		return "GitLab"
	case GitProviderBitbucket:
		return "Bitbucket"
	case GitProviderGitea:
		return "Gitea"
	case GitProviderSelfHosted:
		return "Self-hosted"
	default:
		return string(gp)
	}
}

// DefaultRegistry returns the default Docker registry for this provider
func (gp GitProvider) DefaultRegistry() DockerRegistry {
	if meta, exists := gitProviderMetadata[gp]; exists {
		return meta.defaultRegistry
	}
	return DockerRegistryCustom
}

// ActionsSupported returns true if GitHub Actions are supported
func (gp GitProvider) ActionsSupported() bool {
	if meta, exists := gitProviderMetadata[gp]; exists {
		return meta.actionsSupported
	}
	return false
}

// APIURL returns the API URL for this provider
func (gp GitProvider) APIURL() string {
	if meta, exists := gitProviderMetadata[gp]; exists {
		return meta.apiURL
	}
	return ""
}

// WebURL returns the web URL for this provider
func (gp GitProvider) WebURL() string {
	if meta, exists := gitProviderMetadata[gp]; exists {
		return meta.webURL
	}
	return ""
}

// RequiresPersonalAccessToken returns true if provider requires personal access token
func (gp GitProvider) RequiresPersonalAccessToken() bool {
	if meta, exists := gitProviderMetadata[gp]; exists {
		return meta.requiresPersonalAccessToken
	}
	return true
}

// ValidateGitProvider validates a git provider
func ValidateGitProvider(provider GitProvider) error {
	if !provider.IsValid() {
		return fmt.Errorf("invalid git provider: %s", provider)
	}
	return nil
}