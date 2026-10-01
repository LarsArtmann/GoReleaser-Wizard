// Package domain provides branded, strongly-typed identifiers for domain entities.
//
// This file implements compile-time type safety for entity IDs using the
// go-branded-id package. This prevents accidental mixing
// of different entity IDs (e.g., using a JobID where a WorkflowID is expected).
//
// # Usage
//
// Use the type aliases and constructor functions provided:
//
//	jobID := domain.NewJobID("job-123")
//	workflowID := domain.NewWorkflowID("wf-456")
//
//	// Compile error - cannot use JobID where WorkflowID is expected:
//	func ProcessWorkflow(id WorkflowID) { ... }
//	ProcessWorkflow(jobID) // Error: type mismatch
//
// # Serialization
//
// All ID types support JSON, SQL, and text serialization. They serialize as
// their underlying primitive type (string or int64) and deserialize back to
// the branded type.
//
//	jobID := NewJobID("job-123")
//	data, _ := json.Marshal(jobID) // "\"job-123\""
//
//	var restored JobID
//	json.Unmarshal(data, &restored)
package domain

import (
	id "github.com/larsartmann/go-branded-id"
)

// Brand types for compile-time distinctness.
// These are phantom types - they have no instances but create distinct ID types.

type (
	jobBrand           struct{} // Brand for JobID
	workflowBrand      struct{} // Brand for WorkflowID
	executionPlanBrand struct{} // Brand for ExecutionPlanID
	configBrand        struct{} // Brand for ConfigID
	idBrand            struct{} // Brand for IdID
	iDIDBrand          struct{} // Brand for IDID
	keyBrand           struct{} // Brand for KeyID
	aggregateBrand     struct{} // Brand for AggregateID
	dockerImageBrand   struct{} // Brand for DockerImageID
)

func (jobBrand) Name() string { return "job" }

func (workflowBrand) Name() string { return "workflow" }

func (executionPlanBrand) Name() string { return "executionPlan" }

func (configBrand) Name() string { return "config" }

func (idBrand) Name() string { return "id" }

func (iDIDBrand) Name() string { return "iD" }

func (keyBrand) Name() string { return "key" }

func (aggregateBrand) Name() string { return "aggregate" }

func (dockerImageBrand) Name() string { return "dockerImage" }

// GitHub entity brands.
type (
	gitHubRepoBrand     struct{} // Brand for GitHubRepoID
	gitHubReleaseBrand  struct{} // Brand for GitHubReleaseID
	gitHubAssetBrand    struct{} // Brand for GitHubAssetID
	gitHubWorkflowBrand struct{} // Brand for GitHubWorkflowID
	gitHubUserBrand     struct{} // Brand for GitHubUserID
)

func (gitHubRepoBrand) Name() string { return "gitHubRepo" }

func (gitHubReleaseBrand) Name() string { return "gitHubRelease" }

func (gitHubAssetBrand) Name() string { return "gitHubAsset" }

func (gitHubWorkflowBrand) Name() string { return "gitHubWorkflow" }

func (gitHubUserBrand) Name() string { return "gitHubUser" }

// Type aliases for compile-time safety.
// Using type aliases allows the IDs to be used interchangeably with id.ID
// for generic operations while maintaining distinct type identity.
type (
	// JobID is a strongly-typed job identifier.
	// Use NewJobID() to create instances.
	JobID = id.ID[jobBrand, string]

	// WorkflowID is a strongly-typed workflow execution identifier.
	// Use NewWorkflowID() to create instances.
	WorkflowID = id.ID[workflowBrand, string]

	// ExecutionPlanID is a strongly-typed execution plan identifier.
	// Use NewExecutionPlanID() to create instances.
	ExecutionPlanID = id.ID[executionPlanBrand, string]

	// GitHubRepoID is a strongly-typed GitHub repository identifier.
	// Use NewGitHubRepoID() to create instances.
	GitHubRepoID = id.ID[gitHubRepoBrand, int64]

	// GitHubReleaseID is a strongly-typed GitHub release identifier.
	// Use NewGitHubReleaseID() to create instances.
	GitHubReleaseID = id.ID[gitHubReleaseBrand, int64]

	// GitHubAssetID is a strongly-typed GitHub asset identifier.
	// Use NewGitHubAssetID() to create instances.
	GitHubAssetID = id.ID[gitHubAssetBrand, int64]

	// GitHubWorkflowID is a strongly-typed GitHub workflow identifier.
	// Use NewGitHubWorkflowID() to create instances.
	GitHubWorkflowID = id.ID[gitHubWorkflowBrand, int64]

	// GitHubUserID is a strongly-typed GitHub user identifier.
	// Use NewGitHubUserID() to create instances.
	GitHubUserID = id.ID[gitHubUserBrand, int64]

	// ConfigID is a strongly-typed configuration identifier.
	// Use NewConfigID() to create instances.
	ConfigID = id.ID[configBrand, string]

	// IdID is a strongly-typed generic identifier.
	// Use NewIdID() to create instances.
	IdID = id.ID[idBrand, string]

	// IDID is a strongly-typed identifier for generic ID fields.
	// Use NewIDID() to create instances.
	IDID = id.ID[iDIDBrand, string]

	// KeyID is a strongly-typed signing key identifier.
	// Use NewKeyID() to create instances.
	KeyID = id.ID[keyBrand, string]

	// AggregateID is a strongly-typed aggregate identifier for events.
	// Use NewAggregateID() to create instances.
	AggregateID = id.ID[aggregateBrand, string]

	// DockerImageID is a strongly-typed Docker image identifier.
	// Use NewDockerImageID() to create instances.
	DockerImageID = id.ID[dockerImageBrand, string]
)

// Constructor functions for better ergonomics and type inference.

// NewJobID creates a new JobID from a string value.
// Returns a zero JobID if the input is empty.
//
// Example:
//
//	jobID := domain.NewJobID("config-generation")
//	fmt.Println(jobID.Get()) // "config-generation"
func NewJobID(value string) JobID { return id.NewID[jobBrand](value) }

// NewWorkflowID creates a new WorkflowID from a string value.
// Returns a zero WorkflowID if the input is empty.
//
// Example:
//
//	wfID := domain.NewWorkflowID("init-abc123")
//	fmt.Println(wfID.Get()) // "init-abc123"
func NewWorkflowID(value string) WorkflowID { return id.NewID[workflowBrand](value) }

// NewExecutionPlanID creates a new ExecutionPlanID from a string value.
// Returns a zero ExecutionPlanID if the input is empty.
//
// Example:
//
//	planID := domain.NewExecutionPlanID("plan-xyz789")
//	fmt.Println(planID.Get()) // "plan-xyz789"
func NewExecutionPlanID(value string) ExecutionPlanID {
	return id.NewID[executionPlanBrand](value)
}

// NewGitHubRepoID creates a new GitHubRepoID from an int64 value.
// Returns a zero GitHubRepoID if the input is 0.
//
// Example:
//
//	repoID := domain.NewGitHubRepoID(12345)
//	fmt.Println(repoID.Get()) // 12345
func NewGitHubRepoID(value int64) GitHubRepoID { return id.NewID[gitHubRepoBrand](value) }

// NewGitHubReleaseID creates a new GitHubReleaseID from an int64 value.
// Returns a zero GitHubReleaseID if the input is 0.
//
// Example:
//
//	releaseID := domain.NewGitHubReleaseID(67890)
//	fmt.Println(releaseID.Get()) // 67890
func NewGitHubReleaseID(value int64) GitHubReleaseID {
	return id.NewID[gitHubReleaseBrand](value)
}

// NewGitHubAssetID creates a new GitHubAssetID from an int64 value.
// Returns a zero GitHubAssetID if the input is 0.
//
// Example:
//
//	assetID := domain.NewGitHubAssetID(11111)
//	fmt.Println(assetID.Get()) // 11111
func NewGitHubAssetID(value int64) GitHubAssetID {
	return id.NewID[gitHubAssetBrand](value)
}

// NewGitHubWorkflowID creates a new GitHubWorkflowID from an int64 value.
// Returns a zero GitHubWorkflowID if the input is 0.
//
// Example:
//
//	workflowID := domain.NewGitHubWorkflowID(22222)
//	fmt.Println(workflowID.Get()) // 22222
func NewGitHubWorkflowID(value int64) GitHubWorkflowID {
	return id.NewID[gitHubWorkflowBrand](value)
}

// NewGitHubUserID creates a new GitHubUserID from an int64 value.
// Returns a zero GitHubUserID if the input is 0.
//
// Example:
//
//	userID := domain.NewGitHubUserID(33333)
//	fmt.Println(userID.Get()) // 33333
func NewGitHubUserID(value int64) GitHubUserID {
	return id.NewID[gitHubUserBrand](value)
}

// NewConfigID creates a new ConfigID from a string value.
// Returns a zero ConfigID if the input is empty.
func NewConfigID(value string) ConfigID { return id.NewID[configBrand](value) }

// NewIdID creates a new IdID from a string value.
// Returns a zero IdID if the input is empty.
func NewIdID(value string) IdID { return id.NewID[idBrand](value) }

// NewKeyID creates a new KeyID from a string value.
// Returns a zero KeyID if the input is empty.
func NewKeyID(value string) KeyID { return id.NewID[keyBrand](value) }

// NewAggregateID creates a new AggregateID from a string value.
// Returns a zero AggregateID if the input is empty.
func NewAggregateID(value string) AggregateID { return id.NewID[aggregateBrand](value) }

// NewDockerImageID creates a new DockerImageID from a string value.
// Returns a zero DockerImageID if the input is empty.
func NewDockerImageID(value string) DockerImageID { return id.NewID[dockerImageBrand](value) }

// Compile-time interface assertions to ensure ID types implement expected interfaces.
// These will fail at compile time if the ID types don't satisfy the interfaces.

// String-based IDs implement fmt.Stringer via the embedded id.ID type.
var (
	_ interface{ Get() string } = JobID{}
	_ interface{ Get() string } = WorkflowID{}
	_ interface{ Get() string } = ExecutionPlanID{}
	_ interface{ Get() string } = ConfigID{}
	_ interface{ Get() string } = IdID{}
	_ interface{ Get() string } = KeyID{}
	_ interface{ Get() string } = AggregateID{}
	_ interface{ Get() string } = DockerImageID{}

	_ interface{ IsZero() bool }     = JobID{}
	_ interface{ Equal(JobID) bool } = JobID{}
	_ interface{ Or(JobID) JobID }   = JobID{}

	_ interface{ IsZero() bool }             = WorkflowID{}
	_ interface{ Equal(WorkflowID) bool }    = WorkflowID{}
	_ interface{ Or(WorkflowID) WorkflowID } = WorkflowID{}

	_ interface{ IsZero() bool }               = ExecutionPlanID{}
	_ interface{ Equal(ExecutionPlanID) bool } = ExecutionPlanID{}
	_ interface {
		Or(ExecutionPlanID) ExecutionPlanID
	} = ExecutionPlanID{}

	_ interface{ IsZero() bool }        = ConfigID{}
	_ interface{ Equal(ConfigID) bool } = ConfigID{}

	_ interface{ IsZero() bool }    = IdID{}
	_ interface{ Equal(IdID) bool } = IdID{}

	_ interface{ IsZero() bool }     = KeyID{}
	_ interface{ Equal(KeyID) bool } = KeyID{}

	_ interface{ IsZero() bool }           = AggregateID{}
	_ interface{ Equal(AggregateID) bool } = AggregateID{}

	_ interface{ IsZero() bool }             = DockerImageID{}
	_ interface{ Equal(DockerImageID) bool } = DockerImageID{}
)

// Int64-based IDs implement comparable operations.
var (
	_ interface{ Get() int64 } = GitHubRepoID{}
	_ interface{ Get() int64 } = GitHubReleaseID{}
	_ interface{ Get() int64 } = GitHubAssetID{}
	_ interface{ Get() int64 } = GitHubWorkflowID{}
	_ interface{ Get() int64 } = GitHubUserID{}

	_ interface{ IsZero() bool }            = GitHubRepoID{}
	_ interface{ Equal(GitHubRepoID) bool } = GitHubRepoID{}

	_ interface{ IsZero() bool }               = GitHubReleaseID{}
	_ interface{ Equal(GitHubReleaseID) bool } = GitHubReleaseID{}
)
