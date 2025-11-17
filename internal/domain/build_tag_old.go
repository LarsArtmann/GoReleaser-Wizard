package domain

import (
	"fmt"
)

// BuildTag represents a constrained build tag
// Generated from TypeSpec specification - DO NOT MODIFY MANUALLY
type BuildTag struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

// ValidateBuildTag validates a build tag according to TypeSpec rules
func ValidateBuildTag(tag BuildTag) error {
	if len(tag.Name) < 1 || len(tag.Name) > 50 {
		return fmt.Errorf("build tag name must be between 1 and 50 characters")
	}
	
	// Build tag pattern: letters, numbers, hyphens, underscores only
	pattern := `^[a-zA-Z0-9_-]+$`
	if matched, err := regexp.MatchString(pattern, tag.Name); err != nil {
		return fmt.Errorf("invalid build tag pattern: %v", err)
	} else if !matched {
		return fmt.Errorf("build tag name '%s' can only contain letters, numbers, hyphens, and underscores", tag.Name)
	}
	
	if len(tag.Description) > 200 {
		return fmt.Errorf("build tag description must be 200 characters or less")
	}
	
	return nil
}

// ValidateBuildTags validates a slice of build tags
func ValidateBuildTags(tags []BuildTag) error {
	if len(tags) > 10 {
		return fmt.Errorf("maximum 10 build tags allowed")
	}
	
	for _, tag := range tags {
		if err := ValidateBuildTag(tag); err != nil {
			return fmt.Errorf("build tag '%s': %w", tag.Name, err)
		}
	}
	
	// Check for duplicate tag names
	tagNames := make(map[string]bool)
	for _, tag := range tags {
		if tagNames[tag.Name] {
			return fmt.Errorf("duplicate build tag name: %s", tag.Name)
		}
		tagNames[tag.Name] = true
	}
	
	return nil
}

// CreateBuildTag creates a new build tag with validation
func CreateBuildTag(name, description string) (BuildTag, error) {
	tag := BuildTag{
		Name:        name,
		Description: description,
	}
	
	if err := ValidateBuildTag(tag); err != nil {
		return BuildTag{}, err
	}
	
	return tag, nil
}

// GetCommonBuildTags returns commonly used build tags
func GetCommonBuildTags() []BuildTag {
	return []BuildTag{
		{Name: "netgo", Description: "Enable networking support for Android"},
		{Name: "sqlite", Description: "Enable SQLite support"},
		{Name: "sqlite_omit_load_extension", Description: "Omit SQLite load extension support"},
		{Name: "timetzdata", Description: "Embed timezone database"},
		{Name: "osusergo", Description: "Enable osusergo support"},
		{Name: "static_build", Description: "Enable static linking"},
	}
}

// FilterBuildTagsByPlatform returns platform-specific build tags
func FilterBuildTagsByPlatform(tags []BuildTag, platforms []Platform) []BuildTag {
	filtered := []BuildTag{}
	
	for _, tag := range tags {
		// Add platform-specific filtering logic here
		// For now, return all tags
		filtered = append(filtered, tag)
	}
	
	return filtered
}