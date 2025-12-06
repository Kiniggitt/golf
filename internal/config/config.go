package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the application configuration.
type Config struct {
	// AdminPassword is the password required to access the admin panel.
	AdminPassword string `json:"admin_password"`

	// DataDir is the directory where JSON data files are stored.
	DataDir string `json:"data_dir"`

	// UsersFile is the filename for user data storage.
	UsersFile string `json:"users_file"`

	// StandardFeesFile is the filename for standard fees data storage.
	StandardFeesFile string `json:"standard_fees_file"`

	// Window configuration
	Window WindowConfig `json:"window"`
}

// WindowConfig holds window display settings.
type WindowConfig struct {
	// DefaultWidth is the initial window width.
	DefaultWidth float32 `json:"default_width"`

	// DefaultHeight is the initial window height.
	DefaultHeight float32 `json:"default_height"`

	// AdminWidth is the window width when in admin mode.
	AdminWidth float32 `json:"admin_width"`

	// AdminHeight is the window height when in admin mode.
	AdminHeight float32 `json:"admin_height"`
}

// Default returns a Config with sensible default values.
func Default() *Config {
	return &Config{
		AdminPassword:    "admin123",
		DataDir:          "data",
		UsersFile:        "users.json",
		StandardFeesFile: "standard_fees.json",
		Window: WindowConfig{
			DefaultWidth:  600,
			DefaultHeight: 400,
			AdminWidth:    800,
			AdminHeight:   600,
		},
	}
}

// Load reads configuration from a file, falling back to defaults if the file doesn't exist.
func Load(path string) (*Config, error) {
	// If file doesn't exist, return defaults
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return Default(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate required fields
	if cfg.AdminPassword == "" {
		return nil, fmt.Errorf("admin_password cannot be empty")
	}
	if cfg.DataDir == "" {
		cfg.DataDir = "data"
	}
	if cfg.UsersFile == "" {
		cfg.UsersFile = "users.json"
	}
	if cfg.StandardFeesFile == "" {
		cfg.StandardFeesFile = "standard_fees.json"
	}

	// Validate window settings
	if cfg.Window.DefaultWidth <= 0 {
		cfg.Window.DefaultWidth = 600
	}
	if cfg.Window.DefaultHeight <= 0 {
		cfg.Window.DefaultHeight = 400
	}
	if cfg.Window.AdminWidth <= 0 {
		cfg.Window.AdminWidth = 800
	}
	if cfg.Window.AdminHeight <= 0 {
		cfg.Window.AdminHeight = 600
	}

	return &cfg, nil
}

// Save writes the configuration to a file.
func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetUsersPath returns the full path to the users data file.
func (c *Config) GetUsersPath() string {
	return fmt.Sprintf("%s/%s", c.DataDir, c.UsersFile)
}

// GetStandardFeesPath returns the full path to the standard fees data file.
func (c *Config) GetStandardFeesPath() string {
	return fmt.Sprintf("%s/%s", c.DataDir, c.StandardFeesFile)
}
