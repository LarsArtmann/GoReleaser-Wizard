package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCIProviderIsValid(t *testing.T) {
	tests := []struct {
		name     string
		provider CIProvider
		want     bool
	}{
		{"github", CIProviderGitHub, true},
		{"forgejo", CIProviderForgejo, true},
		{"dual", CIProviderDual, true},
		{"none", CIProviderNone, true},
		{"empty is invalid per enum, but validated separately as unset", "", false},
		{"unknown", "gitlab", false},
		{"typo", "forgejo ", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.provider.IsValid())
		})
	}
}

func TestCIProviderWorkflowGeneration(t *testing.T) {
	tests := []struct {
		name       string
		provider   CIProvider
		github     bool
		forgejo    bool
		forgejoRel bool
	}{
		{"github generates github only", CIProviderGitHub, true, false, false},
		{"unset behaves as github (back-compat)", "", true, false, false},
		{"forgejo generates forgejo only", CIProviderForgejo, false, true, true},
		{"dual generates both, releases to forgejo", CIProviderDual, true, true, true},
		{"none generates nothing", CIProviderNone, false, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.github, tt.provider.GeneratesGitHubWorkflows())
			assert.Equal(t, tt.forgejo, tt.provider.GeneratesForgejoWorkflows())
			assert.Equal(t, tt.forgejoRel, tt.provider.ReleasesToForgejo())
		})
	}
}

func TestValidateCIProvider(t *testing.T) {
	assert.NoError(t, ValidateCIProvider(CIProviderGitHub))
	assert.NoError(t, ValidateCIProvider(CIProviderForgejo))
	assert.NoError(t, ValidateCIProvider(CIProviderDual))
	assert.NoError(t, ValidateCIProvider(CIProviderNone))
	assert.Error(t, ValidateCIProvider(CIProvider("gitlab")))
}

func TestGetAllCIProviders(t *testing.T) {
	all := GetAllCIProviders()
	assert.Len(t, all, 4)
	for _, provider := range all {
		assert.True(t, provider.IsValid())
	}
}

func TestConvertToCIProvider(t *testing.T) {
	tests := []struct {
		input string
		want  CIProvider
	}{
		{"github", CIProviderGitHub},
		{"GitHub", CIProviderGitHub},
		{" forgejo ", CIProviderForgejo},
		{"dual", CIProviderDual},
		{"none", CIProviderNone},
		{"", CIProviderGitHub},
		{"unknown", CIProviderGitHub},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, ConvertToCIProvider(tt.input), "input %q", tt.input)
	}
}

func TestShouldGenerateForgejoWorkflows(t *testing.T) {
	tests := []struct {
		name        string
		mutate      func(*SafeProjectConfig)
		wantActions bool
		wantForgejo bool
	}{
		{
			name:        "default config: github actions only",
			mutate:      func(c *SafeProjectConfig) {},
			wantActions: true,
			wantForgejo: false,
		},
		{
			name: "forgejo provider: forgejo only",
			mutate: func(c *SafeProjectConfig) {
				c.CIProvider = CIProviderForgejo
			},
			wantActions: false,
			wantForgejo: true,
		},
		{
			name: "dual provider: both",
			mutate: func(c *SafeProjectConfig) {
				c.CIProvider = CIProviderDual
			},
			wantActions: true,
			wantForgejo: true,
		},
		{
			name: "none provider: neither",
			mutate: func(c *SafeProjectConfig) {
				c.CIProvider = CIProviderNone
			},
			wantActions: false,
			wantForgejo: false,
		},
		{
			name: "actions disabled: no forgejo either",
			mutate: func(c *SafeProjectConfig) {
				c.CIProvider = CIProviderForgejo
				c.ActionLevel = ActionLevelNone
			},
			wantActions: false,
			wantForgejo: false,
		},
		{
			name: "library project: no github actions files, forgejo yes (fleet CI covers libraries)",
			mutate: func(c *SafeProjectConfig) {
				c.CIProvider = CIProviderForgejo
				c.ProjectType = ProjectTypeLibrary
			},
			wantActions: false,
			wantForgejo: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewSafeProjectConfig()
			config.ApplyDefaults()
			tt.mutate(config)
			assert.Equal(t, tt.wantActions, config.ShouldGenerateActionsFiles())
			assert.Equal(t, tt.wantForgejo, config.ShouldGenerateForgejoWorkflows())
		})
	}
}

func TestApplyDefaultsFillsCIProvider(t *testing.T) {
	config := NewSafeProjectConfig()
	config.ApplyDefaults()
	assert.Equal(t, GetRecommendedCIProvider(), config.CIProvider)
}
