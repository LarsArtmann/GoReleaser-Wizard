package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// BenchmarkProjectDetection benchmarks project detection performance
func BenchmarkProjectDetection(b *testing.B) {
	// Create a temporary project for benchmarking
	tmpDir, _ := os.MkdirTemp("", "wizard-benchmark")
	defer os.RemoveAll(tmpDir)

	// Create a moderately complex project structure
	goMod := `module github.com/user/benchmark-test
go 1.21
require github.com/charmbracelet/huh v0.7.0
require github.com/charmbracelet/lipgloss v1.1.0
`
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644)

	// Create main and cmd structure
	os.MkdirAll(filepath.Join(tmpDir, "cmd", "benchmark-test"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\n\nfunc main() {}"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "cmd", "benchmark-test", "main.go"), []byte("package main\n\nfunc main() {}"), 0644)

	// Change to test directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		config := &ProjectConfig{}
		detectProjectInfo(config)
	}
}

// BenchmarkConfigGeneration benchmarks GoReleaser config generation
func BenchmarkConfigGeneration(b *testing.B) {
	// Create a temporary project for benchmarking
	tmpDir, _ := os.MkdirTemp("", "wizard-config-benchmark")
	defer os.RemoveAll(tmpDir)

	// Create basic project
	goMod := `module github.com/user/config-benchmark
go 1.21
`
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644)
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\n\nfunc main() {}"), 0644)

	// Change to test directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	config := &ProjectConfig{
		ProjectName:        "config-benchmark",
		ProjectDescription: "A benchmark test project",
		ProjectType:        "CLI Application",
		BinaryName:         "config-benchmark",
		MainPath:           ".",
		Platforms:          []string{"linux", "darwin", "windows"},
		Architectures:      []string{"amd64", "arm64"},
		CGOEnabled:         false,
		GitProvider:        "GitHub",
		DockerEnabled:      true,
		DockerRegistry:     "ghcr.io/user",
		Signing:            true,
		Homebrew:           true,
		GenerateActions:    true,
		ActionsOn:          []string{"On version tags (v*)"},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := generateGoReleaserConfig(config)
		if err != nil {
			b.Fatalf("Config generation failed: %v", err)
		}
		os.Remove(".goreleaser.yaml") // Clean up for next iteration
	}
}

// BenchmarkGitHubActionsGeneration benchmarks GitHub Actions workflow generation
func BenchmarkGitHubActionsGeneration(b *testing.B) {
	// Create a temporary project for benchmarking
	tmpDir, _ := os.MkdirTemp("", "wizard-actions-benchmark")
	defer os.RemoveAll(tmpDir)

	// Create basic project
	goMod := `module github.com/user/actions-benchmark
go 1.21
`
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644)
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\n\nfunc main() {}"), 0644)

	// Change to test directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	config := &ProjectConfig{
		ProjectName:     "actions-benchmark",
		BinaryName:      "actions-benchmark",
		GenerateActions: true,
		DockerEnabled:   true,
		Signing:         true,
		ActionsOn:       []string{"On all tags"},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := generateGitHubActions(config)
		if err != nil {
			b.Fatalf("GitHub Actions generation failed: %v", err)
		}
		os.RemoveAll(".github") // Clean up for next iteration
	}
}

// BenchmarkFileOperations benchmarks file operation performance
func BenchmarkFileOperations(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "wizard-fileops-benchmark")
	defer os.RemoveAll(tmpDir)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	testContent := []string{
		"Small content for testing file operations",
		strings.Repeat("Larger content for testing file operations with more data. ", 50),
		strings.Repeat("Very large content for testing file operations with much more data. ", 200),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		content := testContent[i%len(testContent)]
		filename := fmt.Sprintf("benchmark-file-%d.txt", i)

		// Test write operation
		err := SafeFileWrite(filename, []byte(content), 0644)
		if err != nil {
			b.Fatalf("SafeFileWrite failed: %v", err)
		}

		// Test read operation
		readContent, err := SafeReadFile(filename)
		if err != nil {
			b.Fatalf("SafeReadFile failed: %v", err)
		}

		if string(readContent) != content {
			b.Fatalf("Content mismatch")
		}

		// Clean up
		os.Remove(filename)
	}
}

// TestPerformanceCharacteristics tests performance characteristics under different conditions
func TestPerformanceCharacteristics(t *testing.T) {
	tests := []struct {
		name          string
		complexity    int
		expectedMaxMs int64
	}{
		{"simple_project", 1, 100},    // Simple project should complete in <100ms
		{"medium_project", 5, 500},    // Medium project should complete in <500ms
		{"complex_project", 10, 2000}, // Complex project should complete in <2s
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()

			// Create project with specified complexity
			tmpDir, _ := os.MkdirTemp("", fmt.Sprintf("wizard-perf-%s", tt.name))
			defer os.RemoveAll(tmpDir)

			createBenchmarkProject(t, tmpDir, tt.complexity)

			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)
			os.Chdir(tmpDir)

			// Run full wizard workflow
			config := &ProjectConfig{}
			detectProjectInfo(config)

			err := generateGoReleaserConfig(config)
			if err != nil {
				t.Errorf("Config generation failed: %v", err)
			}

			err = generateGitHubActions(config)
			if err != nil {
				t.Errorf("GitHub Actions generation failed: %v", err)
			}

			duration := time.Since(start)

			// Check performance requirements
			if duration.Milliseconds() > tt.expectedMaxMs {
				t.Errorf("Performance exceeded threshold: %v > %dms", duration, tt.expectedMaxMs)
			}

			t.Logf("Performance: %v for %s (threshold: %dms)", duration, tt.name, tt.expectedMaxMs)
		})
	}
}

// TestMemoryUsage tests memory usage patterns
func TestMemoryUsage(t *testing.T) {
	// Get initial memory stats
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// Create and configure multiple projects
	tmpDir, _ := os.MkdirTemp("", "wizard-memory-test")
	defer os.RemoveAll(tmpDir)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	for i := 0; i < 10; i++ {
		projectDir := filepath.Join(tmpDir, fmt.Sprintf("project-%d", i))
		os.MkdirAll(projectDir, 0755)
		os.Chdir(projectDir)

		// Create basic project
		goMod := fmt.Sprintf("module github.com/user/memory-test-%d\ngo 1.21\n", i)
		os.WriteFile("go.mod", []byte(goMod), 0644)
		os.WriteFile("main.go", []byte("package main\n\nfunc main() {}"), 0644)

		// Run wizard operations
		config := &ProjectConfig{}
		detectProjectInfo(config)
		generateGoReleaserConfig(config)
	}

	// Get final memory stats
	runtime.GC()
	runtime.ReadMemStats(&m2)

	// Calculate memory usage
	allocDiff := m2.Alloc - m1.Alloc
	totalAllocDiff := m2.TotalAlloc - m1.TotalAlloc

	t.Logf("Memory usage: Alloc diff = %d bytes, TotalAlloc diff = %d bytes", allocDiff, totalAllocDiff)

	// Memory usage should be reasonable (less than 50MB for 10 projects)
	if totalAllocDiff > 50*1024*1024 {
		t.Errorf("Memory usage too high: %d bytes (> 50MB)", totalAllocDiff)
	}
}

// TestConcurrentOperations tests concurrent wizard operations
func TestConcurrentOperations(t *testing.T) {
	// Test that wizard can handle concurrent operations safely
	concurrency := 5
	done := make(chan bool, concurrency)
	errors := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer func() {
				done <- true
			}()

			// Create temporary project
			tmpDir, _ := os.MkdirTemp("", fmt.Sprintf("wizard-concurrent-%d", id))
			defer os.RemoveAll(tmpDir)

			// Create project
			goMod := fmt.Sprintf("module github.com/user/concurrent-test-%d\ngo 1.21\n", id)
			os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644)
			os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\n\nfunc main() {}"), 0644)

			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)
			os.Chdir(tmpDir)

			// Run wizard operations
			config := &ProjectConfig{}
			detectProjectInfo(config)
			err := generateGoReleaserConfig(config)

			if err != nil {
				errors <- fmt.Errorf("project %d: %v", id, err)
				return
			}
		}(i)
	}

	// Wait for all operations to complete
	for i := 0; i < concurrency; i++ {
		<-done
	}

	// Check for errors
	close(errors)
	for err := range errors {
		t.Errorf("Concurrent operation error: %v", err)
	}

	t.Logf("Successfully completed %d concurrent operations", concurrency)
}

// createBenchmarkProject creates a project with specified complexity
func createBenchmarkProject(t *testing.T, dir string, complexity int) {
	// Create basic structure
	goMod := `module github.com/user/benchmark-project
go 1.21
`
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n\nfunc main() {}"), 0644)

	// Add complexity based on level
	if complexity >= 1 {
		// Add cmd structure
		os.MkdirAll(filepath.Join(dir, "cmd", "benchmark-project"), 0755)
		os.WriteFile(filepath.Join(dir, "cmd", "benchmark-project", "main.go"), []byte("package main\n\nfunc main() {}"), 0644)
	}

	if complexity >= 3 {
		// Add internal structure
		os.MkdirAll(filepath.Join(dir, "internal", "app"), 0755)
		os.WriteFile(filepath.Join(dir, "internal", "app", "app.go"), []byte("package app\n\nfunc Run() {}"), 0644)
		os.WriteFile(filepath.Join(dir, "internal", "app", "config.go"), []byte("package app\n\ntype Config struct {}"), 0644)
	}

	if complexity >= 5 {
		// Add API structure
		os.MkdirAll(filepath.Join(dir, "api", "v1"), 0755)
		os.WriteFile(filepath.Join(dir, "api", "v1", "handler.go"), []byte("package v1\n\nfunc Handle() {}"), 0644)
		os.WriteFile(filepath.Join(dir, "api", "v1", "middleware.go"), []byte("package v1\n\nfunc Middleware() {}"), 0644)
	}

	if complexity >= 7 {
		// Add pkg structure
		os.MkdirAll(filepath.Join(dir, "pkg", "utils"), 0755)
		os.WriteFile(filepath.Join(dir, "pkg", "utils", "helper.go"), []byte("package utils\n\nfunc Helper() {}"), 0644)
		os.WriteFile(filepath.Join(dir, "pkg", "utils", "validator.go"), []byte("package utils\n\nfunc Validate() {}"), 0644)
	}

	if complexity >= 10 {
		// Add extensive structure
		for i := 0; i < 5; i++ {
			pkgName := fmt.Sprintf("pkg%02d", i)
			os.MkdirAll(filepath.Join(dir, pkgName), 0755)
			os.WriteFile(filepath.Join(dir, pkgName, fmt.Sprintf("%s.go", pkgName)), []byte(fmt.Sprintf("package %s\n\nfunc Func() {}", pkgName)), 0644)
		}
	}
}
