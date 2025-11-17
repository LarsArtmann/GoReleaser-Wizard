#!/bin/bash

# Complete GitHub Organization Script
echo "=== GITHUB ISSUES ORGANIZATION EXECUTION ==="
echo ""

# Step 1: Create critical missing issue
echo "STEP 1: Creating critical build recovery issue..."
./create_critical_issue.sh

# Step 2: Add status comments to architecture issues
echo "STEP 2: Adding status comments to architecture issues..."
./add_architecture_comments.sh

# Step 3: Create all milestones
echo "STEP 3: Creating project milestones..."
./create_milestones.sh

# Step 4: Assign issues to milestones
echo "STEP 4: Assigning issues to milestones..."
./assign_milestones.sh

echo ""
echo "=== GITHUB ORGANIZATION COMPLETE ==="
echo "Summary:"
echo "- ✅ Critical build recovery issue created"
echo "- ✅ Status comments added to architecture issues" 
echo "- ✅ 5 milestones created (v0.1.0 through v0.3.0)"
echo "- ✅ All issues assigned to appropriate milestones"
echo ""
echo "Project is now properly organized with clear execution path!"
echo "Critical path: v0.1.0 → Build Recovery → Foundation"
echo "Next priority: Focus on v0.1.0 issues only!"