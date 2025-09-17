package fluree

import (
	"fmt"
	"strings"
)

const defaultBaseURL = "https://data.flur.ee"

// Config contains the credentials and endpoint information required to talk
// to the Fluree Cloud API.
type Config struct {
	APIToken     string
	TenantHandle string
	BaseURL      string
}

// EnvConfigFromLookup builds a Config using the provided lookup function.
// The lookup function typically wraps os.LookupEnv.
func EnvConfigFromLookup(lookup func(string) (string, bool)) (Config, error) {
	cfg := Config{BaseURL: defaultBaseURL}
	if v, ok := lookup("FLUREE_API_TOKEN"); ok {
		cfg.APIToken = strings.TrimSpace(v)
	}
	if v, ok := lookup("FLUREE_HANDLE"); ok {
		cfg.TenantHandle = strings.TrimSpace(v)
	}
	if v, ok := lookup("FLUREE_BASE_URL"); ok && strings.TrimSpace(v) != "" {
		cfg.BaseURL = strings.TrimRight(strings.TrimSpace(v), "/")
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Validate ensures the configuration is complete.
func (c Config) Validate() error {
	var missing []string
	if c.APIToken == "" {
		missing = append(missing, "FLUREE_API_TOKEN")
	}
	if c.TenantHandle == "" {
		missing = append(missing, "FLUREE_HANDLE")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing environment variables: %s", strings.Join(missing, ", "))
	}
	if c.BaseURL == "" {
		return fmt.Errorf("Fluree base URL cannot be empty")
	}
	return nil
}

// WithOverrides returns a copy of the configuration with optional overrides
// applied. Empty override values are ignored.
func (c Config) WithOverrides(apiToken, tenant, baseURL string) Config {
	clone := c
	if strings.TrimSpace(apiToken) != "" {
		clone.APIToken = strings.TrimSpace(apiToken)
	}
	if strings.TrimSpace(tenant) != "" {
		clone.TenantHandle = strings.TrimSpace(tenant)
	}
	if strings.TrimSpace(baseURL) != "" {
		clone.BaseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	}
	return clone
}
