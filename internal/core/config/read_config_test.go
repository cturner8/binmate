package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// writeConfigFile writes a Config struct as JSON to a temp file and returns the path.
func writeConfigFile(t *testing.T, cfg any) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("failed to marshal test config: %v", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}
	return path
}

func TestReadConfig_Default(t *testing.T) {
	t.Setenv("BINMATE_CONFIG_PATH", "")

	cfg := ReadConfig(ConfigFlags{})

	if cfg.Version != 1 {
		t.Errorf("default version = %d, want 1", cfg.Version)
	}
	if cfg.LogLevel != "silent" {
		t.Errorf("default logLevel = %q, want %q", cfg.LogLevel, "silent")
	}
	if len(cfg.Binaries) != 0 {
		t.Errorf("default binaries length = %d, want 0", len(cfg.Binaries))
	}
	if cfg.Global.InstallPath != "" {
		t.Errorf("default global.installPath = %q, want empty string", cfg.Global.InstallPath)
	}
	if len(cfg.Global.Providers) != 0 {
		t.Errorf("default global.providers length = %d, want 0", len(cfg.Global.Providers))
	}
	if cfg.TUI.AskpassMode != "startup" {
		t.Errorf("default tui.askpassMode = %q, want %q", cfg.TUI.AskpassMode, "startup")
	}

}

func TestReadConfig_GlobalProviderAuthenticated(t *testing.T) {
	// Config with global.providers.github.authenticated=true; binary omits authenticated.
	raw := map[string]any{
		"global": map[string]any{
			"providers": map[string]any{
				"github": map[string]any{
					"authenticated": true,
				},
			},
		},
		"binaries": []any{
			map[string]any{
				"id":       "gh",
				"name":     "gh",
				"provider": "github",
				"path":     "cli/cli",
				"format":   ".tar.gz",
			},
		},
	}

	path := writeConfigFile(t, raw)
	cfg := ReadConfig(ConfigFlags{ConfigPath: path})

	// Verify the global provider default is read correctly.
	providerDefaults, exists := cfg.Global.Providers["github"]
	if !exists {
		t.Fatal("global.providers[\"github\"] not found in parsed config; global.providers.github.authenticated is being ignored")
	}
	if !providerDefaults.Authenticated {
		t.Error("global.providers[\"github\"].authenticated = false, want true")
	}

	// Verify merge applies the provider default to the binary.
	if len(cfg.Binaries) == 0 {
		t.Fatal("no binaries parsed from config")
	}
	if cfg.Binaries[0].Authenticated {
		t.Error("binary.authenticated = true, want false (should not inherit from global provider default)")
	}
}

func TestReadConfig_Binary(t *testing.T) {
	// Config with global.providers.github.authenticated=true; binary omits authenticated.
	raw := map[string]any{
		"binaries": []any{
			map[string]any{
				"id":            "gh",
				"name":          "gh",
				"provider":      "github",
				"path":          "cli/cli",
				"format":        ".tar.gz",
				"authenticated": true,
			},
		},
	}

	path := writeConfigFile(t, raw)
	cfg := ReadConfig(ConfigFlags{ConfigPath: path})

	// Verify merge applies the provider default to the binary.
	if len(cfg.Binaries) == 0 {
		t.Fatal("no binaries parsed from config")
	}
	if !cfg.Binaries[0].Authenticated {
		t.Error("binary.authenticated = false, want true")
	}
}

func TestReadConfig_GlobalInstallPath(t *testing.T) {
	raw := map[string]any{
		"global": map[string]any{
			"installPath": "/usr/local/bin",
		},
	}

	path := writeConfigFile(t, raw)
	cfg := ReadConfig(ConfigFlags{ConfigPath: path})

	if cfg.Global.InstallPath != "/usr/local/bin" {
		t.Errorf("global.installPath = %q, want %q", cfg.Global.InstallPath, "/usr/local/bin")
	}
}

func TestReadConfig_AskpassMode(t *testing.T) {
	raw := map[string]any{
		"tui": map[string]any{
			"askpassMode": "always",
		},
	}

	path := writeConfigFile(t, raw)
	cfg := ReadConfig(ConfigFlags{ConfigPath: path})

	if cfg.TUI.AskpassMode != "always" {
		t.Errorf("tui.askpassMode = %q, want %q", cfg.TUI.AskpassMode, "always")
	}
}

func TestReadConfig_AskpassModeDefault(t *testing.T) {
	raw := map[string]any{
		// "tui" is present but askpassMode is omitted; default viper marshal behavior won't set the default value
		"tui": map[string]any{},
	}

	path := writeConfigFile(t, raw)
	cfg := ReadConfig(ConfigFlags{ConfigPath: path})

	if cfg.TUI.AskpassMode != "startup" {
		t.Errorf("tui.askpassMode = %q, want %q", cfg.TUI.AskpassMode, "startup")
	}
}

func TestReadConfig_DateFormat(t *testing.T) {
	raw := map[string]any{
		"tui": map[string]any{
			"dateFormat": "2006-01-02",
		},
	}

	path := writeConfigFile(t, raw)
	cfg := ReadConfig(ConfigFlags{ConfigPath: path})

	if cfg.TUI.DateFormat != "2006-01-02" {
		t.Errorf("tui.dateFormat = %q, want %q", cfg.TUI.DateFormat, "2006-01-02")
	}
}

func TestReadConfig_LogLevel(t *testing.T) {
	raw := map[string]any{
		"logLevel": "debug",
	}

	path := writeConfigFile(t, raw)
	cfg := ReadConfig(ConfigFlags{ConfigPath: path})

	if cfg.LogLevel != "debug" {
		t.Errorf("logLevel = %q, want %q", cfg.LogLevel, "debug")
	}
}
