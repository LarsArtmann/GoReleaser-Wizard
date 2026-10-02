// Package config provides configuration management using koanf.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

// Config holds application configuration.
type Config struct {
	Debug   bool `koanf:"debug"`
	NoColor bool `koanf:"no-color"`
}

// Manager manages configuration loading and access.
type Manager struct {
	k     *koanf.Koanf
	mu    sync.RWMutex
	cfg   Config
	flags *pflag.FlagSet
}

// NewManager creates a new configuration manager.
func NewManager() *Manager {
	return &Manager{
		k: koanf.New("."),
		cfg: Config{
			Debug:   false,
			NoColor: false,
		},
	}
}

// RegisterFlags registers command-line flags for configuration binding.
func (m *Manager) RegisterFlags(flags *pflag.FlagSet) {
	m.flags = flags
}

// Load loads configuration from defaults, file, environment, and flags.
func (m *Manager) Load(cfgFile string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 1. Load defaults (lowest priority)
	err := m.k.Load(confmap.Provider(map[string]any{
		"debug":    false,
		"no-color": false,
	}, "."), nil)
	if err != nil {
		return fmt.Errorf("failed to load defaults: %w", err)
	}

	// 2. Load config file if specified or found
	if cfgFile != "" {
		return m.loadFile(cfgFile)
	}

	// Try to find config in home directory
	m.tryLoadHomeConfig()

	// 3. Load environment variables (higher priority)
	err = m.k.Load(env.Provider(".", env.Opt{
		Prefix:      "GORELEASER_WIZARD_",
		EnvironFunc: os.Environ,
		TransformFunc: func(key, value string) (string, any) {
			key = strings.ToLower(key)
			key = strings.ReplaceAll(key, "_", "-")

			return key, value
		},
	}), nil)
	if err != nil {
		return fmt.Errorf("failed to load environment: %w", err)
	}

	// 4. Load flags (highest priority)
	if m.flags != nil {
		err = m.k.Load(posflag.Provider(m.flags, ".", m.k), nil)
		if err != nil {
			return fmt.Errorf("failed to load flags: %w", err)
		}
	}

	// Unmarshal into config struct
	err = m.k.Unmarshal("", &m.cfg)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// loadFile loads configuration from a specific file.
func (m *Manager) loadFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("config file not found: %s", path)
	}

	_ = info

	err = m.k.Load(file.Provider(path), yaml.Parser())
	if err != nil {
		return fmt.Errorf("failed to load config file: %w", err)
	}

	return nil
}

// tryLoadHomeConfig attempts to load config from home directory.
func (m *Manager) tryLoadHomeConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	configPath := filepath.Join(home, ".goreleaser-wizard.yaml")

	_, err = os.Stat(configPath)
	if err != nil {
		return
	}

	err = m.loadFile(configPath)
	if err != nil {
		// Config file exists but couldn't be read - warn but don't fail
		fmt.Fprintf(os.Stderr, "Warning: could not read config file: %v\n", err)
	}
}

// Get returns the current configuration.
func (m *Manager) Get() Config {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.cfg
}

// IsDebug returns true if debug mode is enabled.
func (m *Manager) IsDebug() bool {
	return m.Get().Debug
}

// NoColors returns true if color output is disabled.
func (m *Manager) NoColors() bool {
	return m.Get().NoColor
}

// Set allows programmatic configuration setting (mainly for tests).
func (m *Manager) Set(key string, value any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	values := map[string]any{key: value}

	err := m.k.Load(confmap.Provider(values, "."), nil)
	if err != nil {
		return fmt.Errorf("failed to load configuration value (key=%q): %w", key, err)
	}

	return fmt.Errorf("failed to unmarshal configuration: %w", m.k.Unmarshal("", &m.cfg))
}

// Reset clears all configuration (mainly for tests).
func (m *Manager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.k = koanf.New(".")
	m.cfg = Config{}
}

// GetRaw returns a raw value from the config (mainly for tests).
func (m *Manager) GetRaw(key string) any {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.k.Get(key)
}

// Global config manager instance.
var globalManager = NewManager()

// GetManager returns the global config manager.
func GetManager() *Manager {
	return globalManager
}

// SetGlobalManager sets the global config manager (for testing).
func SetGlobalManager(m *Manager) {
	globalManager = m
}
