package fluree

import "testing"

type mapLookup map[string]string

func (m mapLookup) Lookup(key string) (string, bool) {
	v, ok := m[key]
	return v, ok
}

func TestEnvConfigFromLookup(t *testing.T) {
	cfg, err := EnvConfigFromLookup(mapLookup{
		"FLUREE_API_TOKEN": " token\t",
		"FLUREE_HANDLE":    " tenant ",
	}.Lookup)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.APIToken != "token" {
		t.Fatalf("expected API token to be trimmed, got %q", cfg.APIToken)
	}
	if cfg.TenantHandle != "tenant" {
		t.Fatalf("expected tenant to be trimmed, got %q", cfg.TenantHandle)
	}
	if cfg.BaseURL != defaultBaseURL {
		t.Fatalf("expected default base URL, got %q", cfg.BaseURL)
	}
}

func TestEnvConfigFromLookup_WithBaseURL(t *testing.T) {
	cfg, err := EnvConfigFromLookup(mapLookup{
		"FLUREE_API_TOKEN": "token",
		"FLUREE_HANDLE":    "tenant",
		"FLUREE_BASE_URL":  " https://example.com/api/ ",
	}.Lookup)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "https://example.com/api"
	if cfg.BaseURL != expected {
		t.Fatalf("expected %q, got %q", expected, cfg.BaseURL)
	}
}

func TestEnvConfigFromLookup_Missing(t *testing.T) {
	if _, err := EnvConfigFromLookup(mapLookup{}.Lookup); err == nil {
		t.Fatalf("expected error for missing values")
	}
}

func TestWithOverrides(t *testing.T) {
	cfg := Config{APIToken: "a", TenantHandle: "b", BaseURL: "https://example"}
	updated := cfg.WithOverrides("", " tenant", "https://override/")
	if updated.APIToken != "a" {
		t.Fatalf("unexpected API token %q", updated.APIToken)
	}
	if updated.TenantHandle != "tenant" {
		t.Fatalf("unexpected tenant %q", updated.TenantHandle)
	}
	if updated.BaseURL != "https://override" {
		t.Fatalf("unexpected base URL %q", updated.BaseURL)
	}
}
