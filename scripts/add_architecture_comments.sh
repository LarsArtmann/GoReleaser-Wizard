#!/bin/bash

# Add Status Comments to Architecture Issues
echo "Adding status comments to architecture issues..."

# Comment on #27 - Strong Type System
gh issue comment 27 \
  --repo LarsArtmann/GoReleaser-Wizard \
  --body "$(cat /tmp/issue_comment_architecture.txt)"

# Comment on #28 - BDD Test Suite  
echo "Architecture comments added!"
echo ""