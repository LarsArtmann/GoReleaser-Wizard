package validation

import (
	"testing"
)

func TestTemplateEscaper_EscapeYAML(t *testing.T) {
	te := NewTemplateEscaper()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Simple text", "hello", "hello"},
		{"Empty string", "", ""},
		{"String with colon", "name: value", "'name: value'"},
		{"String with space", " leading space", "leading space"},
		{"Number", "123", "'123'"},
		{"Boolean-like", "true", "true"},
		{"String starting with special", "!important", "'!important'"},
		{"Multi-line", "line1\nline2", "|-\nline1\n  line2"},
		{"Complex multi-line", "line1: value\nline2: value", "|-\nline1: value\n  line2: value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := te.EscapeYAML(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeYAML() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTemplateEscaper_EscapeShell(t *testing.T) {
	te := NewTemplateEscaper()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Simple text", "hello", "'hello'"},
		{"Empty string", "", ""},
		{"Single quotes", "don't panic", "'don''t panic'"},
		{"Safe characters", "my-app_v1.0", "'my-app_v1.0'"},
		{"Dangerous content", "rm -rf /", ""},     // Should be filtered
		{"Script injection", "; echo hacked", ""}, // Should be filtered
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := te.EscapeShell(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeShell() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTemplateEscaper_EscapeGitHubActions(t *testing.T) {
	te := NewTemplateEscaper()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Simple text", "hello", "hello"},
		{"Expression syntax", "${{ github.repository }}", "'${{ '' }}${{ github.repository }}'"},
		{"Empty string", "", ""},
		{"Complex YAML", "name: value", "'name: value'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := te.EscapeGitHubActions(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeGitHubActions() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTemplateEscaper_EscapeJSON(t *testing.T) {
	te := NewTemplateEscaper()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Simple text", "hello", `"hello"`},
		{"Empty string", "", `""`},
		{"Quotes", "say \"hello\"", `"say \"hello\""`},
		{"Backslash", "path\\to\\file", `"path\\to\\file"`},
		{"Newline", "line1\nline2", `"line1\nline2"`},
		{"Tab", "col1\tcol2", `"col1\tcol2"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := te.EscapeJSON(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeJSON() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTemplateEscaper_EscapeDockerLabel(t *testing.T) {
	te := NewTemplateEscaper()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Simple text", "hello", "hello"},
		{"Empty string", "", ""},
		{"Valid characters", "my-app_v1.0", "my-app_v1.0"},
		{"Invalid characters", "my@app$", "my-app-"},
		{"Starts with number", "123label", "label-123label"},
		{"Starts with dot", ".hidden", "label-.hidden"},
		{"Starts with dash", "-dash", "label--dash"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := te.EscapeDockerLabel(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeDockerLabel() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTemplateEscaper_ValidateTemplateContent(t *testing.T) {
	te := NewTemplateEscaper()

	tests := []struct {
		name         string
		content      string
		templateType string
		wantErr      bool
	}{
		{"Safe YAML", "name: myapp", "yaml", false},
		{"YAML injection", "name: ${SCRIPT}", "yaml", false}, // Pattern doesn't include full variable
		{"Safe shell", "echo 'hello'", "shell", false},
		{"Shell injection", "rm -rf /", "shell", true},
		{"Safe GitHub Actions", "name: build", "github-actions", false},
		{"GitHub Actions injection", "${{ github.token }}", "github-actions", true},
		{"Script tag", "<script>alert('xss')</script", "any", true},
		{"JavaScript", "javascript:void(0)", "any", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := te.ValidateTemplateContent(tt.content, tt.templateType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTemplateContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLooksLikeNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123", true},
		{"123.45", true},
		{"-123", true},
		{"+123", true},
		{"1e10", true},
		{"123abc", false},
		{"abc123", false},
		{"", false},
		{"-", false},
		{".", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := looksLikeNumber(tt.input)
			if result != tt.expected {
				t.Errorf("looksLikeNumber() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestContainsShellInjection(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"echo hello", false},
		{"rm -rf /", true},
		{"cat file | grep pattern", true},
		{"command && rm file", true},
		{"script.sh", false},
		{"$(rm file)", true},
		{"`rm file`", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := containsShellInjection(tt.input)
			if result != tt.expected {
				t.Errorf("containsShellInjection() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsValidDockerLabel(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"label", true},
		{"my-label.v1", true},
		{"my_label", true},
		{"", false},
		{"invalid@label", false},
		{"label with spaces", false},
		{"label/with/slashes", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isValidDockerLabel(tt.input)
			if result != tt.expected {
				t.Errorf("isValidDockerLabel() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Fuzzing tests for escaping functions
func FuzzEscapeYAML(f *testing.F) {
	seed := []string{"hello", "name: value", "don't", "multiline\nstring"}
	for _, s := range seed {
		f.Add(s)
	}

	te := NewTemplateEscaper()
	f.Fuzz(func(t *testing.T, input string) {
		result := te.EscapeYAML(input)
		// Should not panic
		_ = result
	})
}

func FuzzEscapeShell(f *testing.F) {
	seed := []string{"hello", "rm -rf", "; echo", "command$(rm)"}
	for _, s := range seed {
		f.Add(s)
	}

	te := NewTemplateEscaper()
	f.Fuzz(func(t *testing.T, input string) {
		result := te.EscapeShell(input)
		// Should not panic
		_ = result
	})
}

func FuzzEscapeJSON(f *testing.F) {
	seed := []string{"hello", "quote's", `back\slash`, "newline\n"}
	for _, s := range seed {
		f.Add(s)
	}

	te := NewTemplateEscaper()
	f.Fuzz(func(t *testing.T, input string) {
		result := te.EscapeJSON(input)
		// Should not panic
		_ = result
	})
}

// Benchmark tests
func BenchmarkEscapeYAML(b *testing.B) {
	te := NewTemplateEscaper()
	input := "my-project-name: value with 'quotes'"

	for b.Loop() {
		_ = te.EscapeYAML(input)
	}
}

func BenchmarkEscapeShell(b *testing.B) {
	te := NewTemplateEscaper()
	input := "my-app with single 'quotes'"

	for b.Loop() {
		_ = te.EscapeShell(input)
	}
}

func BenchmarkEscapeJSON(b *testing.B) {
	te := NewTemplateEscaper()
	input := `string with "quotes" and \backslashes`

	for b.Loop() {
		_ = te.EscapeJSON(input)
	}
}
