package main

import (
	"strings"
	"testing"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

func TestProviderWorkflowTips(t *testing.T) {
	tests := []struct {
		name       string
		provider   domain.CIProvider
		configOnly bool
		wantSubs   []string
		wantAbsent []string
	}{
		{
			name:       "github_full_run_points_at_release_workflow",
			provider:   domain.CIProviderGitHub,
			configOnly: false,
			wantSubs:   []string{".github/workflows/release.yml"},
		},
		{
			name:       "github_config_only_has_no_workflow_tip",
			provider:   domain.CIProviderGitHub,
			configOnly: true,
			wantAbsent: []string{"workflows"},
		},
		{
			name:       "forgejo_points_at_forgejo_workflows",
			provider:   domain.CIProviderForgejo,
			configOnly: false,
			wantSubs:   []string{".forgejo/workflows/ci.yml", ".forgejo/workflows/release.yml"},
			wantAbsent: []string{".github/workflows"},
		},
		{
			name:       "forgejo_config_only_has_no_workflow_tip",
			provider:   domain.CIProviderForgejo,
			configOnly: true,
			wantAbsent: []string{"workflows"},
		},
		{
			name:       "dual_points_at_both_workflow_sets",
			provider:   domain.CIProviderDual,
			configOnly: false,
			wantSubs:   []string{".github/workflows/release.yml", ".forgejo/workflows/"},
		},
		{
			name:       "none_generates_no_tip",
			provider:   domain.CIProviderNone,
			configOnly: false,
			wantAbsent: []string{"workflows"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tips := providerWorkflowTips(tt.provider, tt.configOnly)

			joined := strings.Join(tips, "\n")
			for _, want := range tt.wantSubs {
				if !strings.Contains(joined, want) {
					t.Errorf("tips %q do not contain %q", joined, want)
				}
			}

			for _, absent := range tt.wantAbsent {
				if strings.Contains(joined, absent) {
					t.Errorf("tips %q must not mention %q", joined, absent)
				}
			}
		})
	}
}

// TestProviderWorkflowTipsUnsetDefaultsToGitHub pins the backward-compatible
// default: an unset provider behaves as github (gotcha 16).
func TestProviderWorkflowTipsUnsetDefaultsToGitHub(t *testing.T) {
	unset := domain.CIProvider("")
	got := providerWorkflowTips(unset, false)

	if len(got) != 1 || !strings.Contains(got[0], ".github/workflows/release.yml") {
		t.Errorf("unset provider tips = %v, want the github release workflow tip", got)
	}
}

func TestGeneratesAnyWorkflows(t *testing.T) {
	tests := []struct {
		provider domain.CIProvider
		want     bool
	}{
		{domain.CIProviderGitHub, true},
		{domain.CIProviderForgejo, true},
		{domain.CIProviderDual, true},
		{domain.CIProviderNone, false},
		{domain.CIProvider(""), true},
	}

	for _, tt := range tests {
		if got := tt.provider.GeneratesAnyWorkflows(); got != tt.want {
			t.Errorf("CIProvider(%q).GeneratesAnyWorkflows() = %v, want %v", tt.provider, got, tt.want)
		}
	}
}
