#!/bin/bash

# Create Milestones for Project Organization
echo "Creating project milestones..."

# v0.1.0 - Critical Build Recovery
gh milestone create \
  --repo LarsArtmann/GoReleaser-Wizard \
  --title "v0.1.0 - Critical Build Recovery" \
  --description "Foundation restoration - fix build system, complete domain layer, basic CLI functionality. Critical path that enables all subsequent development work." \
  --due-on "2025-11-24"

# v0.1.1 - Testing & Validation  
gh milestone create \
  --repo LarsArtmann/GoReleaser-Wizard \
  --title "v0.1.1 - Testing & Validation" \
  --description "Comprehensive testing infrastructure - unit tests, integration tests, BDD scenarios, security validation." \
  --due-on "2025-12-08"

# v0.1.2 - Feature Foundation
gh milestone create \
  --repo LarsArtmann/GoReleaser-Wizard \
  --title "v0.1.2 - Feature Foundation" \
  --description "Core feature implementation - migrate command, advanced CLI features, multi-binary support, configuration templates." \
  --due-on "2025-12-22"

# v0.2.0 - Production Readiness
gh milestone create \
  --repo LarsArtmann/GoReleaser-Wizard \
  --title "v0.2.0 - Production Readiness" \
  --description "Production deployment - distribution packages, CI/CD pipeline, GoReleaser integration, plugin architecture." \
  --due-on "2026-01-12"

# v0.3.0 - Professional Polish
gh milestone create \
  --repo LarsArtmann/GoReleaser-Wizard \
  --title "v0.3.0 - Professional Polish" \
  --description "Professional completion - web UI, REST API, animated demos, enterprise features, community building." \
  --due-on "2026-02-02"

echo "Milestones created!"
echo ""