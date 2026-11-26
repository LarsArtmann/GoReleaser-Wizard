package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var throwawayCommandDirPattern = regexp.MustCompile(`exec\.Command\([^)]*\)\.Dir\s*=`)

func TestTestSourcesNeverConfigureThrowawayCommand(t *testing.T) {
	files, err := filepath.Glob("*_test.go")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat("main.go"); err != nil {
		t.Skipf("not running inside the package directory: %v", err)
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}

		for i, line := range strings.Split(string(content), "\n") {
			if throwawayCommandDirPattern.MatchString(line) {
				t.Errorf(
					"%s:%d configures a throwaway exec.Cmd: the assignment is discarded and the command runs in the process working directory. Assign it to a variable first or pass the directory as an argument:\n\t%s",
					file,
					i+1,
					strings.TrimSpace(line),
				)
			}
		}
	}
}

func TestPackageDirectoryHasNoGitRepository(t *testing.T) {
	if _, err := os.Stat(".git"); err == nil {
		t.Fatal(
			"a git repository exists inside the package directory; test setup must never run git init in the source tree. Use initGitRepository(t, dir) so commands run in a temporary directory.",
		)
	}
}
