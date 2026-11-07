package git

import (
	"os"
	"strings"
	"testing"
)

func TestParseGitHubURL(t *testing.T) {
	tests := []struct {
		name  string
		url   string
		owner string
		repo  string
	}{
		{
			name:  "https_url",
			url:   "https://github.com/octocat/hello.git",
			owner: "octocat",
			repo:  "hello",
		},
		{
			name:  "https_url_without_git_suffix",
			url:   "https://github.com/octocat/hello",
			owner: "octocat",
			repo:  "hello",
		},
		{
			name:  "ssh_url",
			url:   "git@github.com:octocat/hello.git",
			owner: "octocat",
			repo:  "hello",
		},
		{
			name:  "ssh_url_with_port",
			url:   "ssh://git@github.com:22/octocat/hello.git",
			owner: "octocat",
			repo:  "hello",
		},
		{
			name:  "nested_repo_path_takes_first_two_segments",
			url:   "https://github.com/octocat/hello/extra",
			owner: "octocat",
			repo:  "hello",
		},
		{
			name:  "non_github_url_falls_back_to_placeholders",
			url:   "https://gitlab.com/octocat/hello.git",
			owner: "owner",
			repo:  "repo",
		},
		{
			name:  "empty_url_falls_back_to_placeholders",
			url:   "",
			owner: "owner",
			repo:  "repo",
		},
		{
			name:  "github_host_without_path_falls_back",
			url:   "https://github.com/",
			owner: "owner",
			repo:  "repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo := ParseGitHubURL(tt.url)
			if owner != tt.owner || repo != tt.repo {
				t.Errorf("ParseGitHubURL(%q) = (%q, %q), want (%q, %q)", tt.url, owner, repo, tt.owner, tt.repo)
			}
		})
	}
}

func TestGetMajorVersion(t *testing.T) {
	tests := []struct {
		version string
		want    string
	}{
		{"v1.2.3", "1"},
		{"1.2.3", "1"},
		{"v10.0.0", "10"},
		{"", ""},
		{"v", ""},
	}

	for _, tt := range tests {
		t.Run("version_"+tt.version, func(t *testing.T) {
			if got := GetMajorVersion(tt.version); got != tt.want {
				t.Errorf("GetMajorVersion(%q) = %q, want %q", tt.version, got, tt.want)
			}
		})
	}
}

func TestIncPatchVersion(t *testing.T) {
	tests := []struct {
		version string
		want    string
	}{
		{"v1.2.3", "v1.2.4-next"},
		// A missing v prefix is ADDED by the increment path; pinned because
		// the goreleaser template's incpatch output shape depends on it.
		{"1.2.3", "v1.2.4-next"},
		{"v0.0.0", "v0.0.1-next"},
		// Two segments never reach the increment branch: shape kept, no v added.
		{"v1.2", "1.2-next"},
		// Sscanf parses the numeric prefix and drops a prerelease suffix.
		{"v1.2.3-rc1", "v1.2.4-next"},
	}

	for _, tt := range tests {
		t.Run("version_"+tt.version, func(t *testing.T) {
			if got := IncPatchVersion(tt.version); got != tt.want {
				t.Errorf("IncPatchVersion(%q) = %q, want %q", tt.version, got, tt.want)
			}
		})
	}
}

func TestStripVersionPrefix(t *testing.T) {
	tests := []struct {
		version string
		want    string
	}{
		{"v1.2.3", "1.2.3"},
		{"1.2.3", "1.2.3"},
		{"vv1", "v1"},
	}

	for _, tt := range tests {
		t.Run("version_"+tt.version, func(t *testing.T) {
			if got := stripVersionPrefix(tt.version); got != tt.want {
				t.Errorf("stripVersionPrefix(%q) = %q, want %q", tt.version, got, tt.want)
			}
		})
	}
}

func TestGetCurrentDate(t *testing.T) {
	got := GetCurrentDate()
	if len(got) != 10 || got[4] != '-' || got[7] != '-' {
		t.Errorf("GetCurrentDate() = %q, want YYYY-MM-DD", got)
	}
}

// TestCommandInsideRealRepo exercises the Command wrappers against the wizard's
// own checkout, which is a git repository with an origin remote in CI and
// local clones alike.
func TestCommandInsideRealRepo(t *testing.T) {
	if _, err := os.Stat("../../.git"); err != nil {
		t.Skip("not running inside the wizard checkout; real-repo cases not verifiable")
	}

	cmd := NewCommand(t.Context())

	if !cmd.IsRepository() {
		t.Fatal("IsRepository() = false inside the wizard checkout")
	}

	if !cmd.HasRemote("origin") {
		t.Fatal("HasRemote(origin) = false inside the wizard checkout")
	}

	url, err := cmd.GetRemoteURL("origin")
	if err != nil {
		t.Fatalf("GetRemoteURL(origin) = %v", err)
	}

	if url == "" {
		t.Error("GetRemoteURL(origin) = empty")
	}

	version, err := cmd.GetVersion()
	if err != nil {
		t.Fatalf("GetVersion() = %v", err)
	}

	if !strings.Contains(version, ".") {
		t.Errorf("GetVersion() = %q, want a dotted go version", version)
	}

	branches, err := cmd.GetBranches()
	if err != nil {
		t.Fatalf("GetBranches() = %v", err)
	}

	if len(branches) == 0 {
		t.Error("GetBranches() = empty inside a real repository")
	}
}

// TestCommandOutsideRepository pins the graceful non-repo behavior.
func TestCommandOutsideRepository(t *testing.T) {
	t.Chdir(t.TempDir())

	cmd := NewCommand(t.Context())

	if cmd.IsRepository() {
		t.Error("IsRepository() = true in an empty temp dir")
	}

	if cmd.HasRemote("origin") {
		t.Error("HasRemote(origin) = true in an empty temp dir")
	}

	if _, err := cmd.GetRemoteURL("origin"); err == nil {
		t.Error("GetRemoteURL(origin) = nil error outside a repository, want error")
	}
}

func TestSplitGitLines(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		include func(string) bool
		want    []string
	}{
		{
			name:    "splits_and_trims",
			output:  "main\nfeature/x\n",
			include: func(string) bool { return true },
			want:    []string{"main", "feature/x"},
		},
		{
			name:    "filter_applies_to_trimmed_line",
			output:  "keep\nskip\n",
			include: func(line string) bool { return line != "skip" },
			want:    []string{"keep"},
		},
		{
			name:    "empty_output_yields_nil",
			output:  "",
			include: func(string) bool { return true },
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitGitLines(tt.output, tt.include)

			if len(got) != len(tt.want) {
				t.Fatalf("splitGitLines() = %v, want %v", got, tt.want)
			}

			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("splitGitLines()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
