package listing

import (
	"testing"
)

func TestConfigDefaults(t *testing.T) {
	cfg := &Config{}
	if cfg.RootPath != "" {
		t.Errorf("Expected RootPath to be empty by default, got %q", cfg.RootPath)
	}
	if cfg.ShowHeaderFooters != false {
		t.Errorf("Expected ShowHeaderFooters=false by default, got %v", cfg.ShowHeaderFooters)
	}
	// You can add more default checks if needed
}

func TestConfigAssignment(t *testing.T) {
	cfg := &Config{
		RootPath:          "/some/path",
		Include:           "*.go",
		Exclude:           "node_modules",
		GitTrackedOnly:    true,
		ShowTree:          false,
		ShowContent:       true,
		SeparateContent:   true,
		Output:            "out.txt",
		ShowLineNumbers:   true,
		ShowHeaderFooters: true,
	}
	if cfg.RootPath != "/some/path" {
		t.Error("RootPath not set correctly.")
	}
	if cfg.Include != "*.go" {
		t.Error("Include not set correctly.")
	}
	if !cfg.GitTrackedOnly {
		t.Error("GitTrackedOnly should be true.")
	}
	if cfg.ShowTree {
		t.Error("ShowTree should be false.")
	}
	// ... etc. for your fields
}
