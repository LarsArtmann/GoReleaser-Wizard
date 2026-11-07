package main

import (
	"errors"
	"io/fs"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// TestFlakeGoVersionMeetsGoModRequirement guards the toolchain pin drift that
// gotcha 20 documents: the flake must pin a Go version at least as new as the
// go directive in go.mod, because buildGoModule sets GOTOOLCHAIN=local and a
// stale pin fails every nix build with "go.mod requires go >= X".
func TestFlakeGoVersionMeetsGoModRequirement(t *testing.T) {
	flakeBytes, err := os.ReadFile("../../flake.nix")
	if errors.Is(err, fs.ErrNotExist) {
		t.Skip(
			"flake.nix is absent from this source snapshot (nix builds exclude it); the pin guard runs where the flake lives",
		)
	}

	if err != nil {
		t.Fatalf("read flake.nix: %v", err)
	}

	goModBytes, err := os.ReadFile("../../go.mod")
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}

	goDirective := parseGoDirective(t, string(goModBytes))
	flakeGo := parseFlakeGoPin(t, string(flakeBytes))

	meets := flakeGo.major == goDirective.major && flakeGo.minor >= goDirective.minor
	if !meets {
		t.Errorf(
			"flake pins go_%d_%d but go.mod requires go >= %d.%d; bump the flake pin (buildGoModule.override { go = pkgs.go_1_%d; })",
			flakeGo.major,
			flakeGo.minor,
			goDirective.major,
			goDirective.minor,
			goDirective.minor,
		)
	}
}

type goVersion struct {
	major int
	minor int
}

func parseGoDirective(t *testing.T, goMod string) goVersion {
	t.Helper()

	re := regexp.MustCompile(`(?m)^go\s+(\d+)\.(\d+)$`)
	match := re.FindStringSubmatch(goMod)

	if match == nil {
		t.Fatalf("no `go <major>.<minor>` directive found in go.mod")
	}

	return parseVersionParts(t, match[1], match[2])
}

// parseFlakeGoPin extracts the version from the flake's buildGoModule.override
// go pin, e.g. `pkgs.buildGoModule.override { go = pkgs.go_1_27; }`.
func parseFlakeGoPin(t *testing.T, flake string) goVersion {
	t.Helper()

	re := regexp.MustCompile(`buildGoModule\.override\s*\{\s*go\s*=\s*pkgs\.go_(\d+)_(\d+)`)
	match := re.FindStringSubmatch(flake)

	if match == nil {
		t.Fatalf("no buildGoModule.override { go = pkgs.go_X_Y; } pin found in flake.nix")
	}

	return parseVersionParts(t, match[1], match[2])
}

func parseVersionParts(t *testing.T, majorStr, minorStr string) goVersion {
	t.Helper()

	major, err := strconv.Atoi(strings.TrimSpace(majorStr))
	if err != nil {
		t.Fatalf("parse major version %q: %v", majorStr, err)
	}

	minor, err := strconv.Atoi(strings.TrimSpace(minorStr))
	if err != nil {
		t.Fatalf("parse minor version %q: %v", minorStr, err)
	}

	return goVersion{major: major, minor: minor}
}
