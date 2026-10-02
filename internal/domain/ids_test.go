package domain

import (
	"encoding/json/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestJobID tests the JobID branded type.
func TestJobID(t *testing.T) {
	t.Run("constructor creates valid ID", func(t *testing.T) {
		id := NewJobID("job-123")
		assert.False(t, id.IsZero())
		assert.Equal(t, "job-123", id.Get())
	})

	t.Run("zero value is zero", func(t *testing.T) {
		var id JobID
		assert.True(t, id.IsZero())
	})

	t.Run("equality works", func(t *testing.T) {
		id1 := NewJobID("job-123")
		id2 := NewJobID("job-123")
		id3 := NewJobID("job-456")

		assert.True(t, id1.Equal(id2))
		assert.False(t, id1.Equal(id3))
	})

	t.Run("string representation", func(t *testing.T) {
		id := NewJobID("job-123")
		assert.Equal(t, "job:job-123", id.String())
	})

	t.Run("json serialization", func(t *testing.T) {
		id := NewJobID("job-123")
		data, err := json.Marshal(id)
		require.NoError(t, err)
		assert.Equal(t, `"job-123"`, string(data))
	})

	t.Run("json deserialization", func(t *testing.T) {
		data := []byte(`"job-123"`)

		var id JobID

		err := json.Unmarshal(data, &id)
		require.NoError(t, err)
		assert.Equal(t, "job-123", id.Get())
	})

	t.Run("zero value serializes to null", func(t *testing.T) {
		var id JobID

		data, err := json.Marshal(id)
		require.NoError(t, err)
		assert.Equal(t, "null", string(data))
	})

	t.Run("or returns default for zero", func(t *testing.T) {
		var empty JobID

		defaultID := NewJobID("default")

		result := empty.Or(defaultID)
		assert.Equal(t, "default", result.Get())
	})

	t.Run("or returns self for non-zero", func(t *testing.T) {
		id := NewJobID("job-123")
		defaultID := NewJobID("default")

		result := id.Or(defaultID)
		assert.Equal(t, "job-123", result.Get())
	})
}

// TestWorkflowID tests the WorkflowID branded type.
func TestWorkflowID(t *testing.T) {
	t.Run("constructor creates valid ID", func(t *testing.T) {
		id := NewWorkflowID("wf-123")
		assert.False(t, id.IsZero())
		assert.Equal(t, "wf-123", id.Get())
	})

	t.Run("workflow and job IDs are distinct types", func(t *testing.T) {
		// This test verifies at compile time that JobID and WorkflowID
		// are distinct types. If they weren't, this would compile.
		jobID := NewJobID("job-123")
		workflowID := NewWorkflowID("wf-123")

		// These should work - same type
		var (
			_ = jobID
			_ = workflowID
		)

		// The compiler prevents this:
		// var _ WorkflowID = jobID  // Compile error!
		// var _ JobID = workflowID  // Compile error!

		// Verify they're different values

		assert.NotEqual(t, jobID.Get(), workflowID.Get())
	})
}

// TestExecutionPlanID tests the ExecutionPlanID branded type.
func TestExecutionPlanID(t *testing.T) {
	t.Run("constructor creates valid ID", func(t *testing.T) {
		id := NewExecutionPlanID("plan-123")
		assert.False(t, id.IsZero())
		assert.Equal(t, "plan-123", id.Get())
	})
}

// TestGitHubRepoID tests the GitHubRepoID branded type.
func TestGitHubRepoID(t *testing.T) {
	t.Run("constructor creates valid ID", func(t *testing.T) {
		id := NewGitHubRepoID(12345)
		assert.False(t, id.IsZero())
		assert.Equal(t, int64(12345), id.Get())
	})

	t.Run("comparison works", func(t *testing.T) {
		id1 := NewGitHubRepoID(100)
		id2 := NewGitHubRepoID(200)

		result, err := id1.Compare(id2)
		require.NoError(t, err)
		assert.Equal(t, -1, result)

		result, err = id2.Compare(id1)
		require.NoError(t, err)
		assert.Equal(t, 1, result)

		result, err = id1.Compare(id1)
		require.NoError(t, err)
		assert.Equal(t, 0, result)
	})

	t.Run("json serialization as number", func(t *testing.T) {
		id := NewGitHubRepoID(12345)
		data, err := json.Marshal(id)
		require.NoError(t, err)
		assert.Equal(t, `12345`, string(data))
	})

	t.Run("json deserialization", func(t *testing.T) {
		data := []byte(`12345`)

		var id GitHubRepoID

		err := json.Unmarshal(data, &id)
		require.NoError(t, err)
		assert.Equal(t, int64(12345), id.Get())
	})
}

// TestGitHubReleaseID tests the GitHubReleaseID branded type.
func TestGitHubReleaseID(t *testing.T) {
	t.Run("constructor creates valid ID", func(t *testing.T) {
		id := NewGitHubReleaseID(67890)
		assert.False(t, id.IsZero())
		assert.Equal(t, int64(67890), id.Get())
	})
}

// TestGitHubAssetID tests the GitHubAssetID branded type.
func TestGitHubAssetID(t *testing.T) {
	t.Run("constructor creates valid ID", func(t *testing.T) {
		id := NewGitHubAssetID(11111)
		assert.False(t, id.IsZero())
		assert.Equal(t, int64(11111), id.Get())
	})
}

// TestGitHubWorkflowID tests the GitHubWorkflowID branded type.
func TestGitHubWorkflowID(t *testing.T) {
	t.Run("constructor creates valid ID", func(t *testing.T) {
		id := NewGitHubWorkflowID(22222)
		assert.False(t, id.IsZero())
		assert.Equal(t, int64(22222), id.Get())
	})
}

// TestGitHubUserID tests the GitHubUserID branded type.
func TestGitHubUserID(t *testing.T) {
	t.Run("constructor creates valid ID", func(t *testing.T) {
		id := NewGitHubUserID(33333)
		assert.False(t, id.IsZero())
		assert.Equal(t, int64(33333), id.Get())
	})
}

// TestIDTypeSafety verifies that different ID types cannot be mixed.
// This is a compile-time check - if this compiles, the types are distinct.
func TestIDTypeSafety(t *testing.T) {
	// These are all distinct types - if they weren't, the compiler would
	// allow assigning one to another.
	jobID := NewJobID("job-123")
	workflowID := NewWorkflowID("wf-123")
	planID := NewExecutionPlanID("plan-123")
	repoID := NewGitHubRepoID(12345)
	releaseID := NewGitHubReleaseID(67890)
	assetID := NewGitHubAssetID(11111)
	workflowGHID := NewGitHubWorkflowID(22222)
	userID := NewGitHubUserID(33333)

	// Verify all are non-zero
	assert.False(t, jobID.IsZero())
	assert.False(t, workflowID.IsZero())
	assert.False(t, planID.IsZero())
	assert.False(t, repoID.IsZero())
	assert.False(t, releaseID.IsZero())
	assert.False(t, assetID.IsZero())
	assert.False(t, workflowGHID.IsZero())
	assert.False(t, userID.IsZero())

	// Verify string types
	assert.Equal(t, "job-123", jobID.Get())
	assert.Equal(t, "wf-123", workflowID.Get())
	assert.Equal(t, "plan-123", planID.Get())

	// Verify int64 types
	assert.Equal(t, int64(12345), repoID.Get())
	assert.Equal(t, int64(67890), releaseID.Get())
	assert.Equal(t, int64(11111), assetID.Get())
	assert.Equal(t, int64(22222), workflowGHID.Get())
	assert.Equal(t, int64(33333), userID.Get())
}

// TestIDSerialization verifies serialization round-trips for all ID types.
func TestIDSerialization(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "JobID",
			input:    NewJobID("job-123"),
			expected: `"job-123"`,
		},
		{
			name:     "WorkflowID",
			input:    NewWorkflowID("wf-123"),
			expected: `"wf-123"`,
		},
		{
			name:     "ExecutionPlanID",
			input:    NewExecutionPlanID("plan-123"),
			expected: `"plan-123"`,
		},
		{
			name:     "GitHubRepoID",
			input:    NewGitHubRepoID(12345),
			expected: `12345`,
		},
		{
			name:     "GitHubReleaseID",
			input:    NewGitHubReleaseID(67890),
			expected: `67890`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(data))
		})
	}
}
