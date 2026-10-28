package domain

import (
	"slices"
	"strings"
)

// ActionTrigger represents GitHub Actions triggers
// Generated from TypeSpec specification - DO NOT MODIFY MANUALLY.
type ActionTrigger string

const (
	ActionTriggerVersionTags ActionTrigger = "version-tags" // Version Tags (v*)
	ActionTriggerAllTags     ActionTrigger = "all-tags"     // All Tags (*)
	ActionTriggerManual      ActionTrigger = "manual"       // Manual Trigger
	ActionTriggerMain        ActionTrigger = "main"         // Push to Main
	ActionTriggerRelease     ActionTrigger = "release"      // Published Release
	ActionTriggerMainPush    ActionTrigger = "main-push"    // Push to Main (explicit)
	ActionTriggerPullRequest ActionTrigger = "pull-request" // Pull Request
	ActionTriggerSchedule    ActionTrigger = "schedule"     // Scheduled Trigger
)

// ActionTrigger metadata - generated from TypeSpec invariants.
type actionTriggerMeta struct {
	githubPattern  string
	description    string
	recommendedFor []ProjectType
}

var actionTriggerMetaMap = map[ActionTrigger]actionTriggerMeta{
	ActionTriggerVersionTags: {
		githubPattern:  "push:\n  tags:\n    - 'v*'",
		description:    "Triggers on version tags like v1.0.0",
		recommendedFor: []ProjectType{ProjectTypeCLI, ProjectTypeWebAPI, ProjectTypeGRPCService},
	},
	ActionTriggerAllTags: {
		githubPattern:  "push:\n  tags:\n    - '*'",
		description:    "Triggers on any tag",
		recommendedFor: []ProjectType{ProjectTypeLibrary},
	},
	ActionTriggerManual: {
		githubPattern: "workflow_dispatch:",
		description:   "Can be triggered manually from GitHub UI",
		recommendedFor: []ProjectType{
			ProjectTypeCLI,
			ProjectTypeWebAPI,
			ProjectTypeGRPCService,
			ProjectTypeDesktop,
		},
	},
	ActionTriggerMain: {
		githubPattern:  "push:\n  branches:\n    - main",
		description:    "Triggers on pushes to main branch",
		recommendedFor: []ProjectType{ProjectTypeWebAPI, ProjectTypeGRPCService},
	},
	ActionTriggerRelease: {
		githubPattern: "release:\n    types: [published]",
		description:   "Triggers when a GitHub release is published",
		recommendedFor: []ProjectType{
			ProjectTypeCLI,
			ProjectTypeWebAPI,
			ProjectTypeGRPCService,
			ProjectTypeDesktop,
		},
	},
	ActionTriggerMainPush: {
		githubPattern:  "push:\n  branches:\n    - main",
		description:    "Triggers on pushes to the main branch (explicit trigger name)",
		recommendedFor: []ProjectType{ProjectTypeCLI, ProjectTypeWebAPI, ProjectTypeGRPCService},
	},
}

// IsValid returns true if ActionTrigger is valid.
func (at ActionTrigger) IsValid() bool {
	_, exists := actionTriggerMetaMap[at]

	return exists
}

// String returns human-readable display name.
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
	case ActionTriggerMainPush:
		return "Push to Main (explicit)"
	case ActionTriggerPullRequest:
		return "Pull Request"
	case ActionTriggerSchedule:
		return "Scheduled Trigger"
	default:
		return string(at)
	}
}

// GitHubPattern returns the GitHub Actions YAML pattern.
func (at ActionTrigger) GitHubPattern() string {
	if meta, exists := actionTriggerMetaMap[at]; exists {
		return meta.githubPattern
	}

	return ""
}

// Description returns the description of this trigger.
func (at ActionTrigger) Description() string {
	if meta, exists := actionTriggerMetaMap[at]; exists {
		return meta.description
	}

	return ""
}

// RecommendedFor returns the project types this trigger is recommended for.
func (at ActionTrigger) RecommendedFor() []ProjectType {
	if meta, exists := actionTriggerMetaMap[at]; exists {
		return meta.recommendedFor
	}

	return []ProjectType{}
}

// ValidateActionTriggers validates a slice of action triggers.
func ValidateActionTriggers(triggers []ActionTrigger) error {
	return validateEnumSlice(triggers, "action trigger", ActionTrigger.IsValid)
}

// GetRecommendedTriggers returns recommended triggers for a project type.
func GetRecommendedTriggers(projectType ProjectType) []ActionTrigger {
	recommended := []ActionTrigger{}

	for _, trigger := range GetAllActionTriggers() {
		meta := actionTriggerMetaMap[trigger]
		if slices.Contains(meta.recommendedFor, projectType) {
			recommended = append(recommended, trigger)
		}
	}

	return recommended
}

// GetAllActionTriggers returns all available action triggers.
func GetAllActionTriggers() []ActionTrigger {
	return []ActionTrigger{
		ActionTriggerVersionTags, ActionTriggerAllTags, ActionTriggerManual,
		ActionTriggerMain, ActionTriggerRelease,
	}
}

// GenerateGitHubActionsTriggersYAML generates YAML triggers for multiple triggers.
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
