#!/bin/bash

# Create Critical Missing Issue
echo "Creating critical build recovery issue..."

gh issue create \
  --repo LarsArtmann/GoReleaser-Wizard \
  --title "🚨 CRITICAL: EMERGENCY BUILD SYSTEM RECOVERY" \
  --label "bug,critical,urgent" \
  --milestone "v0.1.0" \
  --body "$(cat /tmp/issue_critical_build_recovery.txt)"

echo "Critical issue created!"
echo ""