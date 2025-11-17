package domain

// ActionTrigger represents GitHub Actions triggers
// Generated from TypeSpec specification - DO NOT MODIFY MANUALLY
type ActionTrigger string

const (
	ActionTriggerVersionTags ActionTrigger = "version-tags" // Version Tags (v*)
	ActionTriggerAllTags     ActionTrigger = "all-tags"     // All Tags (*)
	ActionTriggerManual      ActionTrigger = "manual"       // Manual Trigger
	ActionTriggerMain        ActionTrigger = "main"         // Push to Main
	ActionTriggerRelease     ActionTrigger = "release"      // Published Release
)

// ActionTrigger metadata - generated from TypeSpec invariants
type actionTriggerMetadata struct {
	githubPattern   string
	description     string
	recommendedFor  []ProjectType
}

var actionTriggerMetadata = map[ActionTrigger]actionTriggerMetadata{
	ActionTriggerVersionTags: {
		githubPattern: "push:\n  tags:\n    - 'v*'",
		description:  "Triggers on version tags like v1.0.0",
		recommendedFor: []ProjectType{ProjectTypeCLI, ProjectTypeWeb, ProjectTypeAPI},
	},
	ActionTriggerAllTags: {
		githubPattern: "push:\n  tags:\n    - '*'",
		description:  "Triggers on any tag",
		recommendedFor: []ProjectType{ProjectTypeLibrary},
	},
	ActionTriggerManual: {
		githubPattern: "workflow_dispatch:",
		description:  "Can be triggered manually from GitHub UI",
		recommendedFor: []ProjectType{ProjectTypeCLI, ProjectTypeWeb, ProjectTypeAPI, ProjectTypeDesktop},
	},
	ActionTriggerMain: {
		githubPattern: "push:\n  branches:\n    - main",
		description:  "Triggers on pushes to main branch",
		recommendedFor: []ProjectType{ProjectTypeWeb, ProjectTypeAPI},
	},
	ActionTriggerRelease: {
		githubPattern: "release:\n    types: [published]",
		description:  "Triggers when a GitHub release is published",
		recommendedFor: []ProjectType{ProjectTypeCLI, ProjectTypeWeb, ProjectTypeAPI, ProjectTypeDesktop},
	},
}

// IsValid returns true if ActionTrigger is valid
func (at ActionTrigger) IsValid() bool {
	_, exists := actionTriggerMetadata[at]
	return exists
}

// String returns human-readable display name
func (at ActionTrigger) String() string {
	switch at {
	case ActionTriggerVersionTags:
		return "Version Tags (v*)"
	case ActionTriggerAllTags:
		return "All Tags (*)"
	case ActionTriggerManual:
		return "Manual Trigger"
	case ActionTriggerMain:
		return "Push to Main"
	case ActionTriggerRelease:
		return "Published Release"
	default:
		return string(at)
	}
}

// GitHubPattern returns the GitHub Actions YAML pattern
func (at ActionTrigger) GitHubPattern() string {
	if meta, exists := actionTriggerMetadata[at]; exists {
		return meta.githubPattern
	}
	return ""
}

// Description returns the description of this trigger
func (at ActionTrigger) Description() string {
	if meta, exists := actionTriggerMetadata[at]; exists {
		return meta.description
	}
	return ""
}

// RecommendedFor returns the project types this trigger is recommended for
func (at ActionTrigger) RecommendedFor() []ProjectType {
	if meta, exists := actionTriggerMetadata[at]; exists {
		return meta.recommendedFor
	}
	return []ProjectType{}
}

// ValidateActionTriggers validates a slice of action triggers
func ValidateActionTriggers(triggers []ActionTrigger) error {
	if len(triggers) == 0 {
		return fmt.Errorf("at least one action trigger is required")
	}
	
	for _, trigger := range triggers {
		if !trigger.IsValid() {
			return fmt.Errorf("invalid action trigger: %s", trigger)
		}
	}
	
	return nil
}

// GetRecommendedTriggers returns recommended triggers for a project type
func GetRecommendedTriggers(projectType ProjectType) []ActionTrigger {
	recommended := []ActionTrigger{}
	for _, trigger := range GetAllActionTriggers() {
		meta := actionTriggerMetadata[trigger]
		for _, recommendedType := range meta.recommendedFor {
			if recommendedType == projectType {
				recommended = append(recommended, trigger)
				break
			}
		}
	}
	return recommended
}

// GetAllActionTriggers returns all available action triggers
func GetAllActionTriggers() []ActionTrigger {
	return []ActionTrigger{
		ActionTriggerVersionTags, ActionTriggerAllTags, ActionTriggerManual,
		ActionTriggerMain, ActionTriggerRelease,
	}
}

// GenerateGitHubActionsTriggersYAML generates YAML triggers for multiple triggers
func GenerateGitHubActionsTriggersYAML(triggers []ActionTrigger) string {
	if len(triggers) == 0 {
		return ""
	}
	
	if len(triggers) == 1 {
		return triggers[0].GitHubPattern()
	}
	
	// For multiple triggers, combine them
	patterns := []string{}
	for _, trigger := range triggers {
		patterns = append(patterns, trigger.GitHubPattern())
	}
	
	return strings.Join(patterns, "\n\n  ")
}