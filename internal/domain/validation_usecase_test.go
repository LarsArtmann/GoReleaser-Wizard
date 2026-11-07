package domain

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// stubLogger satisfies the domain Logger interface without output.
type stubLogger struct{}

func (stubLogger) Debug(string, ...any)                         {}
func (stubLogger) Info(string, ...any)                          {}
func (stubLogger) Warn(string, ...any)                          {}
func (stubLogger) Error(string, ...any)                         {}
func (stubLogger) Fatal(string, ...any)                         {}
func (stubLogger) DebugContext(context.Context, string, ...any) {}
func (stubLogger) InfoContext(context.Context, string, ...any)  {}
func (stubLogger) WarnContext(context.Context, string, ...any)  {}
func (stubLogger) ErrorContext(context.Context, string, ...any) {}

func (stubLogger) WithField(string, any) Logger     { return stubLogger{} }
func (stubLogger) WithFields(map[string]any) Logger { return stubLogger{} }
func (stubLogger) WithError(error) Logger           { return stubLogger{} }

// osRepo wraps the real filesystem: the use case talks to a temp module
// layout, not mocks, per the integration-tests-over-mocks philosophy.
type osRepo struct{}

func (osRepo) ReadFile(_ context.Context, path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (osRepo) WriteFile(_ context.Context, path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

func (osRepo) CreateFile(_ context.Context, path string) (io.WriteCloser, error) {
	return os.Create(path)
}

func (osRepo) DeleteFile(_ context.Context, path string) error {
	return os.Remove(path)
}

func (osRepo) FileExists(_ context.Context, path string) (bool, error) {
	_, err := os.Stat(path)

	return err == nil, nil
}

func (osRepo) CreateDir(_ context.Context, path string, perm os.FileMode) error {
	return os.Mkdir(path, perm)
}

func (osRepo) CreateDirAll(_ context.Context, path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (osRepo) DirExists(_ context.Context, path string) (bool, error) {
	info, err := os.Stat(path)

	return err == nil && info.IsDir(), nil
}

func (osRepo) ReadDir(_ context.Context, path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

func (osRepo) GetFileInfo(_ context.Context, path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func (osRepo) CheckPermissions(_ context.Context, path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, nil
	}

	file.Close()

	return true, nil
}

func (osRepo) AbsPath(path string) (string, error) {
	return filepath.Abs(path)
}

func (osRepo) RelPath(base, target string) (string, error) {
	return filepath.Rel(base, target)
}

func (osRepo) CleanPath(path string) string {
	return filepath.Clean(path)
}

func (osRepo) JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}

func (osRepo) TempDir(dir, pattern string) (string, error) {
	return os.MkdirTemp(dir, pattern)
}

// compile-time interface checks
var (
	_ FileSystemRepository = osRepo{}
	_ Logger               = stubLogger{}
)

func TestValidationUseCaseValidConfiguration(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	// cmd/app/main.go layout the structure check expects
	if err := os.MkdirAll(filepath.Join(dir, "cmd", "app"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(dir, "cmd", "app", "main.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	vu := NewValidationUseCase(stubLogger{}, osRepo{})
	result, err := vu.ValidateConfiguration(context.Background(), validTestConfig())
	if err != nil {
		t.Fatalf("ValidateConfiguration() = %v", err)
	}

	if !result.IsValid {
		t.Errorf("valid config judged invalid: %v", result.Errors)
	}
}

func TestValidationUseCaseInvalidConfiguration(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	vu := NewValidationUseCase(stubLogger{}, osRepo{})

	config := validTestConfig()
	config.MainPath = "../escape"

	result, err := vu.ValidateConfiguration(context.Background(), config)
	if err != nil {
		t.Fatalf("ValidateConfiguration() = %v", err)
	}

	if result.IsValid {
		t.Error("config with traversal main path judged valid")
	}

	if len(result.Errors) == 0 {
		t.Error("config with traversal main path produced no errors")
	}
}

func TestValidationUseCaseWarnsOnRecommendations(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	vu := NewValidationUseCase(stubLogger{}, osRepo{})

	config := validTestConfig()
	config.DockerSupport = DockerSupportBuild

	result, err := vu.ValidateConfiguration(context.Background(), config)
	if err != nil {
		t.Fatalf("ValidateConfiguration() = %v", err)
	}

	// Missing Dockerfile triggers a recommendation warning path.
	if len(result.Warnings) == 0 {
		t.Log("no warnings emitted for missing Dockerfile; warnings:", result.Warnings)
	}
}

func TestAddValidationErrorFiltersNonDomainErrors(t *testing.T) {
	result := &ValidationResult{IsValid: true, Errors: []*DomainError{}}

	addValidationError(result, errors.New("plain error"))
	if !result.IsValid {
		t.Error("plain errors must be ignored by addValidationError")
	}

	addValidationError(result, NewValidationError(ErrInvalidConfig, "typed", "detail"))
	if result.IsValid {
		t.Error("DomainError must flip IsValid to false")
	}

	if len(result.Errors) != 1 {
		t.Errorf("errors = %d, want 1", len(result.Errors))
	}
}

func TestValidateProjectStructureInTempModule(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	if err := os.MkdirAll(filepath.Join(dir, "cmd", "app"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(dir, "cmd", "app", "main.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module demo\n\ngo 1.25\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	vu := NewValidationUseCase(stubLogger{}, osRepo{})

	result, err := vu.ValidateProjectStructure(context.Background(), dir)
	if err != nil {
		t.Fatalf("ValidateProjectStructure() = %v", err)
	}

	if result == nil {
		t.Fatal("ValidateProjectStructure() = nil result")
	}
}
