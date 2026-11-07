package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
)

func newTestFlagSet(t *testing.T) *pflag.FlagSet {
	t.Helper()

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.Bool("debug", false, "debug mode")
	flags.Bool("no-color", false, "no color")

	return flags
}

func TestManagerLoadDefaults(t *testing.T) {
	t.Setenv("GORELEASER_WIZARD_DEBUG", "")
	t.Setenv("GORELEASER_WIZARD_NO_COLOR", "")

	m := NewManager()

	if err := m.Load(""); err != nil {
		t.Fatalf("Load() = %v", err)
	}

	if m.IsDebug() {
		t.Error("IsDebug() = true with defaults, want false")
	}

	if m.NoColors() {
		t.Error("NoColors() = true with defaults, want false")
	}
}

func TestManagerLoadFromFile(t *testing.T) {
	t.Setenv("GORELEASER_WIZARD_DEBUG", "")
	t.Setenv("GORELEASER_WIZARD_NO_COLOR", "")

	path := filepath.Join(t.TempDir(), "gw.yaml")

	if err := os.WriteFile(path, []byte("debug: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := NewManager()

	if err := m.Load(path); err != nil {
		t.Fatalf("Load(%q) = %v", path, err)
	}

	if !m.IsDebug() {
		t.Error("IsDebug() = false with file debug:true, want true")
	}
}

func TestManagerLoadMissingFileErrors(t *testing.T) {
	m := NewManager()

	if err := m.Load(filepath.Join(t.TempDir(), "nope.yaml")); err == nil {
		t.Error("Load(missing) = nil error, want file-not-found")
	}
}

func TestManagerEnvOverridesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gw.yaml")

	if err := os.WriteFile(path, []byte("debug: false\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("GORELEASER_WIZARD_DEBUG", "true")

	m := NewManager()

	if err := m.Load(path); err != nil {
		t.Fatalf("Load() = %v", err)
	}

	if !m.IsDebug() {
		t.Error("env GORELEASER_WIZARD_DEBUG=true must override file debug:false")
	}
}

func TestManagerFlagsOverrideEnv(t *testing.T) {
	t.Setenv("GORELEASER_WIZARD_DEBUG", "true")

	flags := newTestFlagSet(t)

	if err := flags.Set("debug", "false"); err != nil {
		t.Fatal(err)
	}

	m := NewManager()
	m.RegisterFlags(flags)

	if err := m.Load(""); err != nil {
		t.Fatalf("Load() = %v", err)
	}

	if m.IsDebug() {
		t.Error("explicit --debug=false must override env GORELEASER_WIZARD_DEBUG=true")
	}
}

func TestManagerSetGetResetRoundTrip(t *testing.T) {
	m := NewManager()

	if err := m.Set("debug", true); err != nil {
		t.Fatalf("Set(debug) = %v", err)
	}

	if !m.IsDebug() {
		t.Error("IsDebug() = false after Set(debug, true)")
	}

	if got := m.GetRaw("debug"); got != true {
		t.Errorf("GetRaw(debug) = %v, want true", got)
	}

	m.Reset()

	if m.IsDebug() {
		t.Error("IsDebug() = true after Reset")
	}

	if m.GetRaw("debug") != nil {
		t.Errorf("GetRaw(debug) after Reset = %v, want nil", m.GetRaw("debug"))
	}
}

func TestGlobalManagerRoundTrip(t *testing.T) {
	original := GetManager()

	t.Cleanup(func() { SetGlobalManager(original) })

	replacement := NewManager()

	SetGlobalManager(replacement)

	if GetManager() != replacement {
		t.Error("GetManager() did not return the manager SetGlobalManager installed")
	}
}

func TestManagerEmptyEnvDoesNotOverrideFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gw.yaml")

	if err := os.WriteFile(path, []byte("debug: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("GORELEASER_WIZARD_DEBUG", "")

	m := NewManager()

	if err := m.Load(path); err != nil {
		t.Fatalf("Load() = %v", err)
	}

	if !m.IsDebug() {
		t.Error("empty GORELEASER_WIZARD_DEBUG must not override file debug:true")
	}
}
