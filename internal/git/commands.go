package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

const defaultGitCommandTimeout = 30 * time.Second

// Command represents a git command with context.
//
//nolint:containedctx // Context is stored for builder pattern - commands execute with the context they were created with
type Command struct {
	ctx     context.Context
	dir     string
	env     []string
	timeout time.Duration
}

// NewCommand creates a new git command.
func NewCommand(ctx context.Context) *Command {
	return &Command{
		ctx:     ctx,
		timeout: defaultGitCommandTimeout,
	}
}

// WithDir sets the working directory for git commands.
func (c *Command) WithDir(dir string) *Command {
	c.dir = dir

	return c
}

// WithEnv sets environment variables for git commands.
func (c *Command) WithEnv(env []string) *Command {
	c.env = append(c.env, env...)

	return c
}

// WithTimeout sets the timeout for git commands.
func (c *Command) WithTimeout(timeout time.Duration) *Command {
	c.timeout = timeout

	return c
}

// execute executes a git command and returns the output.
func (c *Command) execute(args ...string) (string, error) {
	cmd := exec.CommandContext(c.ctx, "git", args...)
	if c.dir != "" {
		cmd.Dir = c.dir
	}

	if len(c.env) > 0 {
		cmd.Env = append(os.Environ(), c.env...)
	}

	output, err := cmd.Output()
	if err != nil {
		return "", domain.NewGitError(
			domain.ErrGitCommand,
			"Git command failed: git "+strings.Join(args, " "),
			err.Error(),
		).WithCause(err)
	}

	return strings.TrimSpace(string(output)), nil
}

// GetVersion returns the git version.
func (c *Command) GetVersion() (string, error) {
	return c.execute("version")
}

// GetDescribe gets the git describe output.
func (c *Command) GetDescribe() (string, error) {
	return c.execute("describe", "--tags", "--always", "--dirty")
}

// GetCommitHash gets the current commit hash.
func (c *Command) GetCommitHash() (string, error) {
	return c.execute("rev-parse", "HEAD")
}

// GetRemoteURL gets the URL for a remote.
func (c *Command) GetRemoteURL(remote string) (string, error) {
	return c.execute("remote", "get-url", remote)
}

func splitGitLines(output string, include func(string) bool) []string {
	var result []string

	for line := range strings.SplitSeq(output, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && include(line) {
			result = append(result, line)
		}
	}

	return result
}

// GetBranches gets all branches.
func (c *Command) GetBranches() ([]string, error) {
	output, err := c.execute("branch", "-a")
	if err != nil {
		return nil, err
	}

	return splitGitLines(output, func(branch string) bool {
		return !strings.HasPrefix(branch, "remotes/origin/HEAD")
	}), nil
}

// GetTags gets all tags.
func (c *Command) GetTags() ([]string, error) {
	output, err := c.execute("tag", "-l")
	if err != nil {
		return nil, err
	}

	return splitGitLines(output, func(string) bool { return true }), nil
}

// IsRepository checks if the current directory is a git repository.
func (c *Command) IsRepository() bool {
	_, err := c.execute("rev-parse", "--git-dir")

	return err == nil
}

// HasRemote checks if the repository has a remote.
func (c *Command) HasRemote(remote string) bool {
	_, err := c.execute("remote", "get-url", remote)

	return err == nil
}

// GetRepositoryInfo gets information about the repository.
func (c *Command) GetRepositoryInfo() (*RepositoryInfo, error) {
	if !c.IsRepository() {
		return nil, domain.NewGitError(
			domain.ErrGitNotFound,
			"Not a git repository",
			"Current directory is not a git repository",
		)
	}

	info := &RepositoryInfo{}

	// Get commit hash
	commitHash, err := c.GetCommitHash()
	if err == nil {
		info.CommitHash = commitHash
	}

	// Get description
	describe, err := c.GetDescribe()
	if err == nil {
		info.Description = describe
	}

	// Get remote info
	if c.HasRemote("origin") {
		remoteURL, err := c.GetRemoteURL("origin")
		if err == nil {
			info.RemoteURL = remoteURL
			info.Owner, info.Repo = ParseGitHubURL(remoteURL)
		}
	}

	// Get branches
	branches, err := c.GetBranches()
	if err == nil {
		info.Branches = branches
	}

	// Get tags
	tags, err := c.GetTags()
	if err == nil {
		info.Tags = tags
	}

	return info, nil
}

// RepositoryInfo contains information about a git repository.
type RepositoryInfo struct {
	CommitHash  string   `json:"commit_hash"`
	Description string   `json:"description"`
	RemoteURL   string   `json:"remote_url,omitempty"`
	Owner       string   `json:"owner,omitempty"`
	Repo        string   `json:"repo,omitempty"`
	Branches    []string `json:"branches,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// ParseGitHubURL parses a GitHub URL and extracts owner and repo.
// Supports formats: git@github.com:owner/repo.git, https://github.com/owner/repo.git,
// and ssh://git@github.com:port/owner/repo.git (custom-port hosts must not
// leak the port into the owner). Anything unparseable falls back to the
// documented "owner"/"repo" placeholders.
func ParseGitHubURL(url string) (owner, repo string) {
	if strings.Contains(url, "github.com") {
		parts := strings.Split(url, "github.com")
		if len(parts) > 1 {
			repoPath := strings.TrimPrefix(parts[1], ":")
			repoPath = strings.TrimPrefix(repoPath, "/")
			repoPath = strings.TrimSuffix(repoPath, ".git")

			// ssh://git@github.com:22/owner/repo leaves "22/owner/repo"
			// after the colon/slash trims; drop the port segment so it is
			// never mistaken for the owner.
			if slash := strings.Index(repoPath, "/"); slash != -1 && isAllDigits(repoPath[:slash]) {
				repoPath = repoPath[slash+1:]
			}

			pathParts := strings.Split(repoPath, "/")
			if len(pathParts) > 1 && pathParts[0] != "" && pathParts[1] != "" {
				return pathParts[0], pathParts[1]
			}
		}
	}

	return "owner", "repo"
}

// isAllDigits reports whether s consists solely of ASCII digits (non-empty).
func isAllDigits(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

// GetVersionInfo gets version information from git.
func GetVersionInfo(ctx context.Context) (*VersionInfo, error) {
	cmd := NewCommand(ctx)

	info := &VersionInfo{
		Timestamp: time.Now(),
	}

	// Get version from git describe
	version, err := cmd.GetDescribe()
	if err == nil {
		info.Version = version
	}

	// Get commit hash
	commitHash, err := cmd.GetCommitHash()
	if err == nil {
		info.CommitHash = commitHash
	}

	// Get repository info
	repoInfo, err := cmd.GetRepositoryInfo()
	if err == nil {
		info.Owner = repoInfo.Owner
		info.Repo = repoInfo.Repo
	}

	return info, nil
}

// VersionInfo contains version information from git.
type VersionInfo struct {
	Version    string    `json:"version"`
	CommitHash string    `json:"commit_hash"`
	Owner      string    `json:"owner,omitempty"`
	Repo       string    `json:"repo,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// GetMajorVersion extracts major version from version string.
func GetMajorVersion(version string) string {
	version = stripVersionPrefix(version)

	parts := strings.Split(version, ".")
	if len(parts) > 0 {
		return parts[0]
	}

	return "0"
}

// GetCurrentDate returns current date in YYYY-MM-DD format.
func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

func repositoryInfoPart(selector func(*RepositoryInfo) string, fallback string) string {
	info, err := NewCommand(context.Background()).GetRepositoryInfo()
	if err != nil {
		return fallback
	}

	return selector(info)
}

// GetGitHubOwner tries to get GitHub owner from git remote.
func GetGitHubOwner() string {
	return repositoryInfoPart(func(info *RepositoryInfo) string { return info.Owner }, "owner")
}

// GetGitHubRepo tries to get GitHub repo from git remote.
func GetGitHubRepo() string {
	return repositoryInfoPart(func(info *RepositoryInfo) string { return info.Repo }, "repo")
}

// IncPatchVersion increments the patch version for snapshots.
// It removes the 'v' prefix if present, splits the version string,
// and returns the version with the patch number incremented and "-next" suffix.
func IncPatchVersion(v string) string {
	v = stripVersionPrefix(v)

	parts := strings.Split(v, ".")
	if len(parts) == 3 {
		patch := 0
		if len(parts) > 2 {
			n, err := fmt.Sscanf(parts[2], "%d", &patch)
			if err == nil && n == 1 {
				return fmt.Sprintf("v%s.%s.%d-next", parts[0], parts[1], patch+1)
			}
		}
	}

	return v + "-next"
}

// stripVersionPrefix removes the leading "v" from a semantic-version-like string.
// It centralises the prefix handling used by GetMajorVersion and IncPatchVersion so
// the two callers stay in lockstep if the prefix convention ever changes.
func stripVersionPrefix(version string) string {
	return strings.TrimPrefix(version, "v")
}

// GetGitHubRepoDescription fetches the repository description from GitHub using gh CLI.
// Returns an empty string if gh is not available or the description cannot be fetched.
func GetGitHubRepoDescription() string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"gh",
		"repo",
		"view",
		"--json",
		"description",
		"-q",
		".description",
	)

	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}
