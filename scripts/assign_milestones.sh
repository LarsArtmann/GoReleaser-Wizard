#!/bin/bash

# Assign Issues to Milestones
echo "Assigning issues to milestones..."

# Critical Build Recovery
echo "Assigning critical build recovery issue (once created)..."
# This would be done after issue creation

# v0.1.0 - Critical Build Recovery & Foundation
gh issue edit 27 --repo LarsArtmann/GoReleaser-Wizard --milestone "v0.1.0"  # Strong Type System
gh issue edit 28 --repo LarsArtmann/GoReleaser-Wizard --milestone "v0.1.0"  # BDD Test Suite
gh issue edit 29 --repo LarsArtmann/GoReleaser-Wizard --milestone "v0.1.0"  # Configuration State Machine
gh issue edit 30 --repo LarsArtmann/GoReleaser-Wizard --milestone "v0.1.0"  # Documentation Update

# v0.1.2 - Feature Foundation  
gh issue edit 7 --repo LarsArtmann/GoReleaser-Wizard --milestone "v0.1.2"   # Advanced CLI Features
gh issue edit 19 --repo LarsArtmann/GoReleaser-Wizard --milestone "v0.1.2"  # Migrate Command
gh issue edit 20 --repo LarsArtmann/GoReleaser-Wizard --milestone "v0.1.2"  # Multi-binary Support
gh issue edit 18 --repo LarsArtmann/GoReleaser-Wizard --milestone "v0.1.2"  # Test on Popular Projects

# v0.2.0 - Production Readiness
gh issue edit 22 --repo LarsArtmann/GoReleaser-Wizard --milestone "v0.2.0"  # GitHub Actions CI/CD
gh issue edit 23 --repo LarsArtmann/GoReleaser-Wizard --milestone "v0.2.0"  # Distribute via Homebrew
gh issue edit 24 --repo LarsArtmann/GoReleaser-Wizard --milestone "v0.2.0"  # GoReleaser Integration

# v0.3.0 - Professional Polish
gh issue edit 21 --repo LarsArtmann/GoReleaser-Wizard --milestone "v0.3.0"  # Animated GIF Demo
gh issue edit 6 --repo LarsArtmann/GoReleaser-Wizard --milestone "v0.3.0"   # Final Documentation

echo "Issues assigned to milestones!"
echo ""