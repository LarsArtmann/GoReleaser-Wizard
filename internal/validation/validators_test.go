package validation

import (
	"strings"
	"testing"
)

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid simple name", "myproject", false},
		{"Valid with hyphens", "my-project", false},
		{"Valid with underscores", "my_project", false},
		{"Valid with dots", "my.project", false},
		{"Valid mixed", "my-project.v2", false},
		{"Empty string", "", true},
		{"Too long", strings.Repeat("a", 64), true},
		{"Starts with hyphen", "-project", true},
		{"Ends with hyphen", "project-", true},
		{"Starts with dot", ".project", true},
		{"Ends with dot", "project.", true},
		{"Consecutive hyphens", "project--name", true},
		{"Consecutive dots", "project..name", true},
		{"Reserved name - go", "go", true},
		{"Reserved name - test", "test", true},
		{"Reserved name - con", "con", true},
		{"Reserved name - aux", "aux", true},
		{"Special characters", "project@name", true},
		{"Spaces", "project name", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProjectName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProjectName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateBinaryName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid simple name", "myapp", false},
		{"Valid with hyphens", "my-app", false},
		{"Valid with underscores", "my_app", false},
		{"Empty string", "", true},
		{"Too long", strings.Repeat("a", 256), true},
		{"Starts with number", "123app", true},
		{"Reserved name - con", "con", true},
		{"Reserved name - aux", "aux", true},
		{"Shell metacharacters", "my;app", true},
		{"File separators", "my/app", true},
		{"Windows separators", "my\\app", true},
		{"Quotes", "my'app", true},
		{"Pipes", "my|app", true},
		{"Angle brackets", "my<app>", true},
		{"Dangerous extension", "myapp.exe", true},
		{"Dangerous extension bat", "myapp.bat", true},
		{"Valid number in middle", "my2app", false},
		{"Valid number at end", "myapp2", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBinaryName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBinaryName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateMainPath(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid current dir", ".", false},
		{"Valid relative path", "./cmd/app", false},
		{"Valid deeper path", "cmd/myapp", false},
		{"Empty string", "", true},
		{"Path traversal", "../../../etc/passwd", true},
		{"Path traversal with backslash", "..\\..\\windows\\system32", true},
		{"Absolute path", "/usr/local/bin", true},
		{"Windows absolute path", "C:\\Windows\\System32", true},
		{"Shell metacharacters", "cmd;rm -rf /", true},
		{"Script injection", "./cmd/app && rm -rf", true},
		{"Reserved directory", "etc/passwd", true},
		{"Reserved directory - bin", "bin/app", true},
		{"Reserved directory - usr", "usr/bin/app", true},
		{"Current dir with reserved", "./con", true},
		{"Clean path should be valid", "./cmd", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMainPath(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMainPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateProjectDescription(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid description", "A great Go application", false},
		{"Empty description", "", false},
		{"Max length description", strings.Repeat("a", 255), false},
		{"Too long", strings.Repeat("a", 256), true},
		{"Script injection", "<script>alert('xss')</script>", true},
		{"JavaScript injection", "javascript:alert('xss')", true},
		{"Valid with punctuation", "My app, version 2.0!", false},
		{"Valid with newlines", "Line 1\nLine 2", false},
		{"Suspicious long whitespace", strings.Repeat(" ", 200), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProjectDescription(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProjectDescription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateBuildTags(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		wantErr bool
	}{
		{"Valid single tag", []string{"prod"}, false},
		{"Valid multiple tags", []string{"prod", "linux", "amd64"}, false},
		{"Valid complex tag", []string{"my_custom_tag_123"}, false},
		{"Empty tags", []string{}, false},
		{"Empty tag in slice", []string{"prod", ""}, true},
		{"Too long tag", []string{strings.Repeat("a", 51)}, true},
		{"Invalid characters", []string{"prod-tag"}, true},
		{"Path traversal", []string{"../../etc"}, true},
		{"Shell metacharacters", []string{"prod;rm"}, true},
		{"Valid camelCase", []string{"myCustomTag"}, false},
		{"Valid with numbers", []string{"v2"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBuildTags(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBuildTags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateDockerRegistry(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid Docker Hub", "docker.io/username", false},
		{"Valid GitHub Registry", "ghcr.io/username/app", false},
		{"Valid GitLab Registry", "registry.gitlab.com/username/app", false},
		{"Valid localhost HTTP", "localhost:5000/app", false},
		{"Valid 127.0.0.1 HTTP", "127.0.0.1:5000/app", false},
		{"Valid with port", "ghcr.io:443/username/app", false},
		{"Empty string", "", true},
		{"Non-localhost HTTP", "http://registry.example.com/app", true},
		{"With credentials", "username:password@registry.com/app", true},
		{"Invalid format", "invalid..registry", true},
		{"With query params", "registry.com/app?tag=v1", true},
		{"With fragment", "registry.com/app#tag", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDockerRegistry(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDockerRegistry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGitProvider(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid GitHub", "github", false},
		{"Valid GitHub uppercase", "GitHub", false},
		{"Valid GitLab", "gitlab", false},
		{"Valid Bitbucket", "bitbucket", false},
		{"Valid Gitea", "gitea", false},
		{"Valid Self-hosted", "self-hosted", false},
		{"Valid Local", "local", false},
		{"Invalid provider", "gitfake", true},
		{"Empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGitProvider(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGitProvider() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Normal text", "hello world", "hello world"},
		{"With null bytes", "hello\x00world", "helloworld"},
		{"With control chars", "hello\x01world", "helloworld"},
		{"With tabs", "hello\tworld", "hello\tworld"},
		{"With newlines", "hello\nworld", "hello\nworld"},
		{"Extra whitespace", "  hello world  ", "hello world"},
		{"Mixed", "  hello\x00\tworld\n  ", "hello\tworld\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeInput(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeInput() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Fuzzing tests for security validation
func FuzzValidateProjectName(f *testing.F) {
	seed := []string{"myproject", "my-project", "my_project", "my.project", "con", "aux", "test"}
	for _, s := range seed {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		err := ValidateProjectName(input)
		// The function should not panic with any input
		if err == nil {
			// If valid, should be reproducible
			err2 := ValidateProjectName(input)
			if err2 != nil {
				t.Errorf("Inconsistent validation for %s: first pass ok, second pass failed", input)
			}
		}
	})
}

func FuzzValidateBinaryName(f *testing.F) {
	seed := []string{"myapp", "my-app", "my_app", "con", "aux", "test", "myapp.exe"}
	for _, s := range seed {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		err := ValidateBinaryName(input)
		// Should not panic with any input
		_ = err
	})
}

func FuzzValidateMainPath(f *testing.F) {
	seed := []string{".", "./cmd/app", "../../../etc/passwd", "/usr/bin", "C:\\Windows"}
	for _, s := range seed {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		err := ValidateMainPath(input)
		// Should not panic with any input
		_ = err
	})
}

// Benchmark tests to ensure validation is efficient
func BenchmarkValidateProjectName(b *testing.B) {
	for b.Loop() {
		_ = ValidateProjectName("my-test-project")
	}
}

func BenchmarkValidateBinaryName(b *testing.B) {
	for b.Loop() {
		_ = ValidateBinaryName("myapp")
	}
}

func BenchmarkValidateMainPath(b *testing.B) {
	for b.Loop() {
		_ = ValidateMainPath("./cmd/myapp")
	}
}

func BenchmarkSanitizeInput(b *testing.B) {
	for b.Loop() {
		_ = SanitizeInput("  hello world with whitespace  ")
	}
}
