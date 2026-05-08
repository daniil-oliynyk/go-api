package config

import (
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/toronto_str")
	t.Setenv("PORT", "9090")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000, https://example.com, ")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected config to load, got error: %v", err)
	}

	if cfg.DatabaseURL != "postgres://user:pass@localhost:5432/toronto_str" {
		t.Fatalf("unexpected database url: %q", cfg.DatabaseURL)
	}
	if cfg.Port != "9090" {
		t.Fatalf("unexpected port: %q", cfg.Port)
	}

	expectedOrigins := []string{"http://localhost:3000", "https://example.com"}
	if !reflect.DeepEqual(cfg.CORSAllowedOrigins, expectedOrigins) {
		t.Fatalf("expected origins %v, got %v", expectedOrigins, cfg.CORSAllowedOrigins)
	}
}

func TestLoadDefaultsPort(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/toronto_str")
	t.Setenv("PORT", "")
	t.Setenv("CORS_ALLOWED_ORIGINS", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected config to load, got error: %v", err)
	}

	if cfg.Port != "8080" {
		t.Fatalf("expected default port 8080, got %q", cfg.Port)
	}
	if len(cfg.CORSAllowedOrigins) != 0 {
		t.Fatalf("expected no CORS origins, got %v", cfg.CORSAllowedOrigins)
	}
}

func TestLoadRequiresDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")

	if _, err := Load(); err == nil {
		t.Fatal("expected missing DATABASE_URL error")
	}
}

func TestLoadRejectsInvalidDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "not a url")

	if _, err := Load(); err == nil {
		t.Fatal("expected invalid DATABASE_URL error")
	}
}

func TestLoadRejectsInvalidPort(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/toronto_str")
	t.Setenv("PORT", "70000")

	if _, err := Load(); err == nil {
		t.Fatal("expected invalid PORT error")
	}
}
