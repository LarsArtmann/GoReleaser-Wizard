package domain

import "strings"

// CIProvider selects which CI system the generated workflows target.
// The fleet's primary CI is the self-hosted Forgejo instance; GitHub stays
// as the public shadow (module paths, Homebrew taps, public clones).
type CIProvider string

const (
	// CIProviderGitHub generates only .github/workflows (legacy behavior).
	CIProviderGitHub CIProvider = "github"
	// CIProviderForgejo generates only .forgejo/workflows.
	CIProviderForgejo CIProvider = "forgejo"
	// CIProviderDual generates both workflow sets (transition state).
	CIProviderDual CIProvider = "dual"
	// CIProviderNone generates no CI workflows at all.
	CIProviderNone CIProvider = "none"
)

// IsValid returns true if CIProvider is a known provider.
func (p CIProvider) IsValid() bool {
	switch p {
	case CIProviderGitHub, CIProviderForgejo, CIProviderDual, CIProviderNone:
		return true
	default:
		return false
	}
}

// String returns human-readable display name.
func (p CIProvider) String() string {
	switch p {
	case CIProviderGitHub:
		return "GitHub"
	case CIProviderForgejo:
		return "Forgejo"
	case CIProviderDual:
		return "Dual (Forgejo + GitHub)"
	case CIProviderNone:
		return "None"
	default:
		return string(p)
	}
}

// GeneratesGitHubWorkflows reports whether .github/workflows are generated.
// An unset provider (empty string, pre-migration configs) defaults to GitHub.
func (p CIProvider) GeneratesGitHubWorkflows() bool {
	return p.isGitHubLike()
}

// GeneratesForgejoWorkflows reports whether .forgejo/workflows are generated.
func (p CIProvider) GeneratesForgejoWorkflows() bool {
	return p == CIProviderForgejo || p == CIProviderDual
}

// isGitHubLike treats the unset provider as GitHub for backward compatibility.
func (p CIProvider) isGitHubLike() bool {
	return p == CIProviderGitHub || p == CIProviderDual || p == ""
}

// ReleasesToForgejo reports whether GoReleaser releases target the Forgejo
// instance (gitea_urls) instead of github.com. Dual providers also release
// to Forgejo: the public GitHub copy is maintained by the release shadow.
func (p CIProvider) ReleasesToForgejo() bool {
	return p == CIProviderForgejo || p == CIProviderDual
}

// OrEmpty normalizes an unset provider to GitHub.
func (p CIProvider) OrEmpty() CIProvider {
	if p == "" {
		return CIProviderGitHub
	}
	return p
}

// ValidateCIProvider validates a CI provider value.
func ValidateCIProvider(provider CIProvider) error {
	return validateEnum("CI provider", string(provider), provider.IsValid())
}

// GetAllCIProviders returns all available CI providers.
func GetAllCIProviders() []CIProvider {
	return []CIProvider{
		CIProviderGitHub, CIProviderForgejo, CIProviderDual, CIProviderNone,
	}
}

// GetRecommendedCIProvider returns the recommended CI provider (GitHub for
// backward compatibility; flip to forgejo when a repo migrates).
func GetRecommendedCIProvider() CIProvider {
	return CIProviderGitHub
}

// ConvertToCIProvider converts a string to CIProvider, defaulting to GitHub.
func ConvertToCIProvider(value string) CIProvider {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "github":
		return CIProviderGitHub
	case "forgejo":
		return CIProviderForgejo
	case "dual":
		return CIProviderDual
	case "none":
		return CIProviderNone
	default:
		return CIProviderGitHub
	}
}
