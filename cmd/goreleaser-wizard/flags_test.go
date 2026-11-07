package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
	"github.com/spf13/cobra"
)

func TestGetStringFlagReturnsRegisteredValue(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("project-name", "wizard", "test flag")

	if got := getStringFlag(cmd, "project-name"); got != "wizard" {
		t.Errorf("getStringFlag = %q, want %q", got, "wizard")
	}
}

func TestGetBoolFlagReturnsRegisteredValue(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("force", true, "test flag")

	if got := getBoolFlag(cmd, "force"); !got {
		t.Error("getBoolFlag = false, want true")
	}
}

func TestFlagLookupUnknownNamePanics(t *testing.T) {
	tests := []struct {
		name              string
		read              func(*cobra.Command, string)
		flag              string
		wantPanicContains string
	}{
		{
			name:              "unknown_string_flag_panics",
			read:              func(cmd *cobra.Command, name string) { getStringFlag(cmd, name) },
			flag:              "no-such-string-flag",
			wantPanicContains: `flag "no-such-string-flag" not defined`,
		},
		{
			name:              "unknown_bool_flag_panics",
			read:              func(cmd *cobra.Command, name string) { getBoolFlag(cmd, name) },
			flag:              "no-such-bool-flag",
			wantPanicContains: `flag "no-such-bool-flag" not defined`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}

			defer func() {
				r := recover()
				if r == nil {
					t.Fatal("expected panic for undefined flag, got none")
				}

				msg, ok := r.(string)
				if !ok {
					t.Fatalf("panic value is %T, want string", r)
				}

				if !strings.Contains(msg, tt.wantPanicContains) {
					t.Errorf("panic message %q does not contain %q", msg, tt.wantPanicContains)
				}
			}()

			tt.read(cmd, tt.flag)
		})
	}
}

func TestDetectAlternativeStructureGlobErrorPropagates(t *testing.T) {
	original := globProjects
	injected := errors.New("bad pattern")

	globProjects = func(string, string) ([]string, error) {
		return nil, injected
	}
	defer func() { globProjects = original }()

	mainPath, binaryName, projectType, err := detectAlternativeStructure(t.TempDir())
	if err == nil {
		t.Fatal("detectAlternativeStructure = nil error, want glob failure propagated")
	}

	if !errors.Is(err, injected) {
		t.Errorf("error chain lost injected cause: %v", err)
	}

	if mainPath != "" || binaryName != "" || projectType != "" {
		t.Errorf("error path returned values (%q, %q, %q), want empty", mainPath, binaryName, projectType)
	}

	var domainErr *domain.DomainError
	if !errors.As(err, &domainErr) {
		t.Fatalf("error is %T, want *domain.DomainError", err)
	}

	if domainErr.Code != domain.ErrInvalidProjectStructure {
		t.Errorf("error code = %q, want %q", domainErr.Code, domain.ErrInvalidProjectStructure)
	}
}

func TestDetectAlternativeStructureHappyPaths(t *testing.T) {
	t.Run("root_go_files_fall_back_to_cli", func(t *testing.T) {
		dir := t.TempDir()

		mainPath, binaryName, projectType, err := detectAlternativeStructure(dir)
		if err != nil {
			t.Fatalf("detectAlternativeStructure(%q) = %v, want nil", dir, err)
		}

		if mainPath != "." || projectType != "cli" {
			t.Errorf("got (%q, %q), want (\".\", \"cli\")", mainPath, projectType)
		}

		wantBinary := filepath.Base(dir)
		if binaryName != wantBinary {
			t.Errorf("binaryName = %q, want %q", binaryName, wantBinary)
		}
	})

	t.Run("src_layout_extracts_binary", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, "src", "mytool"), 0o755); err != nil {
			t.Fatalf("mkdir failed: %v", err)
		}

		mainGo := filepath.Join(dir, "src", "mytool", "main.go")
		if err := os.WriteFile(mainGo, []byte("package main\n"), 0o644); err != nil {
			t.Fatalf("write failed: %v", err)
		}

		mainPath, binaryName, projectType, err := detectAlternativeStructure(dir)
		if err != nil {
			t.Fatalf("detectAlternativeStructure(%q) = %v, want nil", dir, err)
		}

		if mainPath != "src/*/main.go" || binaryName != "mytool" || projectType != "library" {
			t.Errorf("got (%q, %q, %q), want (src/*/main.go, mytool, library)", mainPath, binaryName, projectType)
		}
	})
}
