package hedera

import (
	"fmt"
	"strings"
)

const defaultNetwork = "testnet"

// Config captures the credentials and target network required to interact with the
// Hedera SDK.
type Config struct {
	Network            string
	OperatorAccountID  string
	OperatorPrivateKey string
	MirrorNetworkURL   string
}

// EnvConfigFromLookup constructs a Config using the provided lookup function.
// The lookup function typically wraps os.LookupEnv.
func EnvConfigFromLookup(lookup func(string) (string, bool)) (Config, error) {
	cfg := Config{Network: defaultNetwork}
	if v, ok := lookup("HEDERA_NETWORK"); ok && strings.TrimSpace(v) != "" {
		cfg.Network = strings.TrimSpace(v)
	}
	if v, ok := lookup("HEDERA_OPERATOR_ID"); ok {
		cfg.OperatorAccountID = strings.TrimSpace(v)
	}
	if v, ok := lookup("HEDERA_OPERATOR_KEY"); ok {
		cfg.OperatorPrivateKey = strings.TrimSpace(v)
	}
	if v, ok := lookup("HEDERA_MIRROR_URL"); ok {
		cfg.MirrorNetworkURL = strings.TrimSpace(v)
	}
	return cfg, cfg.Validate()
}

// Validate ensures the configuration is internally consistent.
func (c Config) Validate() error {
	if strings.TrimSpace(c.Network) == "" {
		return fmt.Errorf("hedera network cannot be empty")
	}
	if (c.OperatorAccountID == "") != (c.OperatorPrivateKey == "") {
		return fmt.Errorf("hedera operator id and key must both be provided or omitted")
	}
	return nil
}

// HasOperator reports whether operator credentials are configured.
func (c Config) HasOperator() bool {
	return c.OperatorAccountID != "" && c.OperatorPrivateKey != ""
}

// WithOverrides returns a copy of the configuration with optional overrides
// applied. Empty override values are ignored.
func (c Config) WithOverrides(network, operatorID, operatorKey, mirrorURL string) Config {
	clone := c
	if strings.TrimSpace(network) != "" {
		clone.Network = strings.TrimSpace(network)
	}
	if strings.TrimSpace(operatorID) != "" {
		clone.OperatorAccountID = strings.TrimSpace(operatorID)
	}
	if strings.TrimSpace(operatorKey) != "" {
		clone.OperatorPrivateKey = strings.TrimSpace(operatorKey)
	}
	if strings.TrimSpace(mirrorURL) != "" {
		clone.MirrorNetworkURL = strings.TrimSpace(mirrorURL)
	}
	return clone
}
