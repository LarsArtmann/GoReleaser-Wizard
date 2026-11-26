package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

// writeGoModFile writes the go.mod file content.
func writeGoModFile(t *testing.T, dir, moduleName string) {
	t.Helper()

	goMod := "module " + moduleName + "\ngo 1.21\n"

	err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0o644)
	if err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}
}

// createBasicGoProject creates a basic Go project structure for testing.
func createBasicGoProject(t *testing.T, dir, moduleName string) {
	t.Helper()

	writeGoModFile(t, dir, moduleName)

	mainContent := "package main\n\nfunc main() {}\n"

	err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainContent), 0o644)
	if err != nil {
		t.Fatalf("failed to create main.go: %v", err)
	}
}

// inTestDir executes the test function after changing to the test directory.
// It automatically restores the original directory after the test completes.
func inTestDir(t *testing.T, testDir string, testFunc func()) {
	t.Helper()

	t.Chdir(testDir)

	testFunc()
}

// AssertErr checks if an error matches the expected error state.
// This deduplicates the common pattern: if (err != nil) != tt.wantErr.
// Handles typed nil errors properly (a nil *DomainError stored as error interface is not == nil).
func AssertErr(t *testing.T, fnName string, err error, wantErr bool) {
	t.Helper()

	// Properly check for nil - a typed nil interface is not equal to nil
	isNil := err == nil || reflect.ValueOf(err).IsNil()
	// hasError is true when we have an actual error (isNil=false)
	hasError := !isNil

	if hasError != wantErr {
		t.Errorf("%s() error = %v, wantErr %v", fnName, err, wantErr)
	}
}

// AssertFileExists checks if a file exists at the given path.
func AssertFileExists(t *testing.T, path, msg string) {
	t.Helper()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error(msg)
	}
}

// CreateTempDir creates a temporary directory with the given prefix.
// It returns the directory path.
func CreateTempDir(prefix string) string {
	dir, err := os.MkdirTemp("", prefix)
	if err != nil {
		return ""
	}

	return dir
}

func initGitRepository(t *testing.T, dir string) {
	t.Helper()

	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git not found in PATH: %v", err)
	}

	commands := [][]string{
		{"init", "-q"},
		{"config", "user.email", "test@example.com"},
		{"config", "user.name", "Test User"},
		{"add", "."},
		{"commit", "-q", "-m", "init"},
	}

	for _, args := range commands {
		gitCmd := exec.CommandContext(t.Context(), "git", args...)
		gitCmd.Dir = dir

		if out, err := gitCmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v in %s failed: %v: %s", args, dir, err, out)
		}
	}
}
