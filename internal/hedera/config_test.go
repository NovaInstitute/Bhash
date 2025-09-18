package hedera

import "testing"

func TestConfigEnv(t *testing.T) {
	lookup := func(values map[string]string) func(string) (string, bool) {
		return func(key string) (string, bool) {
			v, ok := values[key]
			return v, ok
		}
	}

	cfg, err := EnvConfigFromLookup(lookup(map[string]string{
		"HEDERA_NETWORK":      "previewnet",
		"HEDERA_OPERATOR_ID":  "0.0.1001",
		"HEDERA_OPERATOR_KEY": "302a",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Network != "previewnet" || cfg.OperatorAccountID != "0.0.1001" {
		t.Fatalf("unexpected config: %+v", cfg)
	}

	cfg = cfg.WithOverrides("mainnet", "0.0.2002", "abcd", "")
	if cfg.Network != "mainnet" || cfg.OperatorAccountID != "0.0.2002" {
		t.Fatalf("unexpected override: %+v", cfg)
	}

	if !cfg.HasOperator() {
		t.Fatalf("expected operator credentials to be present")
	}
}

func TestConfigValidateOperatorMismatch(t *testing.T) {
	cfg := Config{Network: "testnet", OperatorAccountID: "0.0.2"}
	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected validation error")
	}
}
