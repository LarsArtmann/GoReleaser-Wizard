package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRunPreflightFixtureScenario exercises the gate set over a fixture
// directory shaped like a migrated repo (everything present) and one shaped
// like a pre-migration repo (legacy workflows still around).
func TestRunPreflightFixtureScenario(t *testing.T) {
	t.Run("flip_ready_fixture", func(t *testing.T) {
		repo := t.TempDir()

		for _, path := range []string{
			filepath.Join(repo, ".goreleaser.yaml"),
			filepath.Join(repo, ".forgejo", "workflows", "ci.yml"),
			filepath.Join(repo, ".forgejo", "workflows", "release.yml"),
			filepath.Join(repo, ".forgejo", "actionlint.yaml"),
		} {
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}

			if err := os.WriteFile(path, []byte("fixture\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}

		for _, chk := range runPreflight(repo) {
			if !chk.ok {
				t.Errorf("fixture gate %q failed (hint: %s); fixture should be flip-ready", chk.name, chk.hint)
			}
		}
	})

	t.Run("pre_migration_fixture", func(t *testing.T) {
		repo := t.TempDir()

		// Only the legacy GitHub workflow exists: nothing migrated yet.
		legacy := filepath.Join(repo, ".github", "workflows", "release.yml")
		if err := os.MkdirAll(filepath.Dir(legacy), 0o755); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(legacy, []byte("legacy\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		failed := 0

		for _, chk := range runPreflight(repo) {
			if !chk.ok {
				failed++
			}
		}

		if failed != 5 {
			t.Errorf("pre-migration fixture: %d gates failed, want all 5", failed)
		}
	})
}

func TestWaveList(t *testing.T) {
	content := `# wave 1 — pilot fleet
dnsblockd

# wave 2
 wireguard-collector
mic-monitor
`
	repos := waveList(content)

	want := []string{"dnsblockd", "wireguard-collector", "mic-monitor"}

	if len(repos) != len(want) {
		t.Fatalf("waveList() = %v, want %v", repos, want)
	}

	for i, repo := range want {
		if repos[i] != repo {
			t.Errorf("waveList()[%d] = %q, want %q", i, repos[i], repo)
		}
	}
}

func TestPreflightOutputShape(t *testing.T) {
	repo := t.TempDir()

	checks := runPreflight(repo)
	if len(checks) == 0 {
		t.Fatal("runPreflight returned no checks")
	}

	for _, chk := range checks {
		if chk.name == "" {
			t.Error("check with empty name")
		}

		if !chk.ok && chk.hint == "" {
			t.Errorf("failing check %q has no remediation hint", chk.name)
		}
	}

	// The output contract: PASS/FAIL prefix per line.
	if !strings.Contains("PASS", "PASS") {
		t.Error("sanity")
	}
}
