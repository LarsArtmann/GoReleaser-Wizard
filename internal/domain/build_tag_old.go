package domain

import (
	"fmt"
	"strings"
)

// BuildTag represents build tags for conditional compilation
type BuildTag struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
}

// IsValid returns true if BuildTag has valid name
func (bt BuildTag) IsValid() bool {
	if len(bt.Name) == 0 {
		return false
	}
	
	// Build tag pattern: alphanumerics, hyphens, underscores
	if !isValidBuildTagName(bt.Name) {
		return false
	}
	
	return true
}

// String returns string representation of BuildTag
func (bt BuildTag) String() string {
	return bt.Name
}

// isValidBuildTagName validates build tag name format
func isValidBuildTagName(name string) bool {
	if len(name) == 0 {
		return false
	}
	
	for _, char := range name {
		if !(char == '_' || char == '-' || (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return false
		}
	}
	
	return true
}

// ValidateBuildTag validates a single build tag
func ValidateBuildTag(tag BuildTag) error {
	if !tag.IsValid() {
		return fmt.Errorf("invalid build tag: %s", tag.Name)
	}
	return nil
}

// ValidateBuildTags validates a slice of build tags
func ValidateBuildTags(tags []BuildTag) error {
	if len(tags) > 50 {
		return fmt.Errorf("too many build tags (max 50)")
	}

	for _, tag := range tags {
		if err := ValidateBuildTag(tag); err != nil {
			return err
		}
	}

	// Check for duplicates
	seen := make(map[string]bool)
	for _, tag := range tags {
		if seen[tag.Name] {
			return fmt.Errorf("duplicate build tag: %s", tag.Name)
		}
		seen[tag.Name] = true
	}

	return nil
}

// GetCommonBuildTags returns commonly used build tags
func GetCommonBuildTags() []BuildTag {
	return []BuildTag{
		{
			Name:        "netgo",
			Description: "Use netgo for networking",
		},
		{
			Name:        "osusergo",
			Description: "Use osusergo for user lookup",
		},
		{
			Name:        "sqlite_omit_load_extension",
			Description: "Omit SQLite load extension",
		},
		{
			Name:        "sqlite_unlock_notify",
			Description: "Enable SQLite unlock notify",
		},
		{
			Name:        "inotify",
			Description: "Enable inotify support",
		},
		{
			Name:        "kqueue",
			Description: "Enable kqueue support",
		},
	}
}

// FilterBuildTagsByPlatform filters build tags by platform compatibility
func FilterBuildTagsByPlatform(tags []BuildTag, platform Platform) []BuildTag {
	// Platform-specific build tag filtering
	platformTags := make(map[Platform][]string){
		PlatformLinux:   {"inotify"},
		PlatformDarwin:  {"kqueue"},
		PlatformWindows:  {},
		PlatformFreeBSD: {"kqueue"},
		PlatformOpenBSD: {"kqueue"},
		PlatformNetBSD:  {"kqueue"},
	}

	filtered := []BuildTag{}
	compatibleTags := platformTags[platform]
	
	for _, tag := range tags {
		// Keep tags that are platform-specific and compatible, or general tags
		isPlatformSpecific := false
		for _, pt := range GetAllPlatforms() {
			if contains(platformTags[pt], tag.Name) {
				isPlatformSpecific = true
				if contains(compatibleTags, tag.Name) {
					filtered = append(filtered, tag)
				}
				break
			}
		}
		
		// Keep general tags
		if !isPlatformSpecific {
			filtered = append(filtered, tag)
		}
	}
	
	return filtered
}

// CreateBuildTag creates a new build tag
func CreateBuildTag(name, description string) BuildTag {
	return BuildTag{
		Name:        name,
		Description: description,
	}
}

// contains checks if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}