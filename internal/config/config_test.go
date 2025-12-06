package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.AdminPassword != "admin123" {
		t.Errorf("Expected default admin password 'admin123', got %s", cfg.AdminPassword)
	}

	if cfg.DataDir != "data" {
		t.Errorf("Expected default data dir 'data', got %s", cfg.DataDir)
	}

	if cfg.UsersFile != "users.json" {
		t.Errorf("Expected default users file 'users.json', got %s", cfg.UsersFile)
	}

	if cfg.StandardFeesFile != "standard_fees.json" {
		t.Errorf("Expected default standard fees file 'standard_fees.json', got %s", cfg.StandardFeesFile)
	}

	if cfg.Window.DefaultWidth != 600 {
		t.Errorf("Expected default width 600, got %.0f", cfg.Window.DefaultWidth)
	}

	if cfg.Window.DefaultHeight != 400 {
		t.Errorf("Expected default height 400, got %.0f", cfg.Window.DefaultHeight)
	}

	if cfg.Window.AdminWidth != 800 {
		t.Errorf("Expected admin width 800, got %.0f", cfg.Window.AdminWidth)
	}

	if cfg.Window.AdminHeight != 600 {
		t.Errorf("Expected admin height 600, got %.0f", cfg.Window.AdminHeight)
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	cfg, err := Load("nonexistent_config.json")
	if err != nil {
		t.Fatalf("Expected no error when loading nonexistent file, got %v", err)
	}

	// Should return defaults
	if cfg.AdminPassword != "admin123" {
		t.Errorf("Expected default password when file doesn't exist")
	}
}

func TestLoadValidConfig(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.json")

	configData := `{
    "admin_password": "mypassword123",
    "data_dir": "testdata",
    "users_file": "test_users.json",
    "standard_fees_file": "test_fees.json",
    "window": {
        "default_width": 700,
        "default_height": 500,
        "admin_width": 900,
        "admin_height": 700
    }
}`

	err := os.WriteFile(configFile, []byte(configData), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.AdminPassword != "mypassword123" {
		t.Errorf("Expected password 'mypassword123', got %s", cfg.AdminPassword)
	}

	if cfg.DataDir != "testdata" {
		t.Errorf("Expected data dir 'testdata', got %s", cfg.DataDir)
	}

	if cfg.UsersFile != "test_users.json" {
		t.Errorf("Expected users file 'test_users.json', got %s", cfg.UsersFile)
	}

	if cfg.StandardFeesFile != "test_fees.json" {
		t.Errorf("Expected standard fees file 'test_fees.json', got %s", cfg.StandardFeesFile)
	}

	if cfg.Window.DefaultWidth != 700 {
		t.Errorf("Expected width 700, got %.0f", cfg.Window.DefaultWidth)
	}

	if cfg.Window.DefaultHeight != 500 {
		t.Errorf("Expected height 500, got %.0f", cfg.Window.DefaultHeight)
	}

	if cfg.Window.AdminWidth != 900 {
		t.Errorf("Expected admin width 900, got %.0f", cfg.Window.AdminWidth)
	}

	if cfg.Window.AdminHeight != 700 {
		t.Errorf("Expected admin height 700, got %.0f", cfg.Window.AdminHeight)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "invalid.json")

	err := os.WriteFile(configFile, []byte("{invalid json}"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	_, err = Load(configFile)
	if err == nil {
		t.Error("Expected error when loading invalid JSON")
	}
}

func TestLoadEmptyPassword(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.json")

	configData := `{
    "admin_password": "",
    "data_dir": "data"
}`

	err := os.WriteFile(configFile, []byte(configData), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	_, err = Load(configFile)
	if err == nil {
		t.Error("Expected error when admin_password is empty")
	}
}

func TestLoadMissingFields(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.json")

	// Minimal config with only password
	configData := `{
    "admin_password": "test123"
}`

	err := os.WriteFile(configFile, []byte(configData), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should use defaults for missing fields
	if cfg.DataDir != "data" {
		t.Errorf("Expected default data dir 'data', got %s", cfg.DataDir)
	}

	if cfg.UsersFile != "users.json" {
		t.Errorf("Expected default users file 'users.json', got %s", cfg.UsersFile)
	}

	if cfg.Window.DefaultWidth != 600 {
		t.Errorf("Expected default width 600, got %.0f", cfg.Window.DefaultWidth)
	}
}

func TestSave(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "saved_config.json")

	cfg := &Config{
		AdminPassword:    "savetest123",
		DataDir:          "mydata",
		UsersFile:        "my_users.json",
		StandardFeesFile: "my_fees.json",
		Window: WindowConfig{
			DefaultWidth:  650,
			DefaultHeight: 450,
			AdminWidth:    850,
			AdminHeight:   650,
		},
	}

	err := cfg.Save(configFile)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatal("Expected config file to be created")
	}

	// Load it back and verify
	loadedCfg, err := Load(configFile)
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if loadedCfg.AdminPassword != cfg.AdminPassword {
		t.Errorf("Password mismatch: expected %s, got %s", cfg.AdminPassword, loadedCfg.AdminPassword)
	}

	if loadedCfg.DataDir != cfg.DataDir {
		t.Errorf("DataDir mismatch: expected %s, got %s", cfg.DataDir, loadedCfg.DataDir)
	}

	if loadedCfg.Window.DefaultWidth != cfg.Window.DefaultWidth {
		t.Errorf("DefaultWidth mismatch: expected %.0f, got %.0f", cfg.Window.DefaultWidth, loadedCfg.Window.DefaultWidth)
	}
}

func TestGetUsersPath(t *testing.T) {
	cfg := &Config{
		DataDir:   "testdata",
		UsersFile: "test_users.json",
	}

	expected := "testdata/test_users.json"
	result := cfg.GetUsersPath()

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestGetStandardFeesPath(t *testing.T) {
	cfg := &Config{
		DataDir:          "testdata",
		StandardFeesFile: "test_fees.json",
	}

	expected := "testdata/test_fees.json"
	result := cfg.GetStandardFeesPath()

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestValidationZeroWindowSize(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.json")

	configData := `{
    "admin_password": "test123",
    "window": {
        "default_width": 0,
        "default_height": 0,
        "admin_width": 0,
        "admin_height": 0
    }
}`

	err := os.WriteFile(configFile, []byte(configData), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should use defaults for zero/negative values
	if cfg.Window.DefaultWidth != 600 {
		t.Errorf("Expected default width 600 for zero value, got %.0f", cfg.Window.DefaultWidth)
	}

	if cfg.Window.DefaultHeight != 400 {
		t.Errorf("Expected default height 400 for zero value, got %.0f", cfg.Window.DefaultHeight)
	}

	if cfg.Window.AdminWidth != 800 {
		t.Errorf("Expected admin width 800 for zero value, got %.0f", cfg.Window.AdminWidth)
	}

	if cfg.Window.AdminHeight != 600 {
		t.Errorf("Expected admin height 600 for zero value, got %.0f", cfg.Window.AdminHeight)
	}
}
