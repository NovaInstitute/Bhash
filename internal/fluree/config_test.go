package fluree

import "testing"

func TestEnvConfigFromLookup(t *testing.T) {
	t.Parallel()

	lookup := func(key string) (string, bool) {
		switch key {
		case "FLUREE_API_TOKEN":
			return "token", true
		case "FLUREE_HANDLE":
			return "tenant", true
		case "FLUREE_BASE_URL":
			return "https://example.local", true
		default:
			return "", false
		}
	}

	cfg, err := EnvConfigFromLookup(lookup)
	if err != nil {
		t.Fatalf("env config: %v", err)
	}
	if cfg.APIToken != "token" || cfg.TenantHandle != "tenant" {
		t.Fatalf("unexpected config: %#v", cfg)
	}
	if cfg.BaseURL != "https://example.local" {
		t.Fatalf("unexpected base url: %s", cfg.BaseURL)
	}
}

func TestEnvConfigFromLookupMissing(t *testing.T) {
	t.Parallel()

	_, err := EnvConfigFromLookup(func(string) (string, bool) { return "", false })
	if err == nil {
		t.Fatalf("expected error when environment variables are missing")
	}
}

func TestWithOverrides(t *testing.T) {
	t.Parallel()

	cfg := Config{APIToken: "token", TenantHandle: "tenant", BaseURL: "https://data.flur.ee"}
	updated := cfg.WithOverrides("override", "owner", "https://custom")
	if updated.APIToken != "override" || updated.TenantHandle != "owner" {
		t.Fatalf("unexpected overrides: %#v", updated)
	}
	if updated.BaseURL != "https://custom" {
		t.Fatalf("expected custom base url, got %s", updated.BaseURL)
	}
	if cfg.APIToken != "token" {
		t.Fatalf("expected original config unchanged")
	}
}
