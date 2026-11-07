package main

import (
	"fmt"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

// providerWorkflowTips returns the "commit the generated workflows" tip
// lines matching what the provider actually generates. Config-only runs
// generate no workflows for any provider, and CIProviderNone generates no
// workflows at all, so both contribute no tip: pointing users at files that
// do not exist is worse than no tip.
func providerWorkflowTips(provider domain.CIProvider, configOnly bool) []string {
	if configOnly || !provider.GeneratesAnyWorkflows() {
		return nil
	}

	switch provider {
	case domain.CIProviderForgejo:
		return []string{
			"  • Commit .forgejo/workflows/ci.yml and .forgejo/workflows/release.yml to enable Forgejo Actions CI/CD",
		}
	case domain.CIProviderDual:
		return []string{
			"  • Commit .github/workflows/release.yml and .forgejo/workflows/ (ci.yml, release.yml) to enable both CI systems",
		}
	case domain.CIProviderGitHub:
		return []string{
			"  • Commit .github/workflows/release.yml to enable CI/CD",
		}
	case domain.CIProviderNone:
		return nil
	default:
		return nil
	}
}

// printGenerateSuccess prints the success banner and provider-aware next
// steps for a finished generation run.
func printGenerateSuccess(provider domain.CIProvider, configOnly bool) {
	fmt.Println()

	if configOnly {
		fmt.Println(successStyle.Render("📄 GoReleaser configuration generated successfully!"))
	} else {
		fmt.Println(
			successStyle.Render("🎉 GoReleaser configuration and workflow generated successfully!"),
		)
	}

	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  • Review generated files")
	fmt.Println("  • Run 'goreleaser-wizard validate' to check configuration")

	for _, tip := range providerWorkflowTips(provider, configOnly) {
		fmt.Println(tip)
	}

	fmt.Println("  • Test configuration with 'goreleaser check'")
}
