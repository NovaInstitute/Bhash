package fluree

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

var runFluree = flag.Bool("run-fluree", false, "run Fluree Cloud integration tests")

func TestIntegrationGeneratePrompt(t *testing.T) {
	if !*runFluree {
		t.Skip("pass -run-fluree to execute Fluree Cloud integration tests")
	}

	dataset := strings.TrimSpace(os.Getenv("FLUREE_DATASET"))
	if dataset == "" {
		t.Skip("set FLUREE_DATASET to enable Fluree Cloud integration tests")
	}

	cfg, err := EnvConfigFromLookup(os.LookupEnv)
	if err != nil {
		t.Skipf("Fluree credentials not configured: %v", err)
	}

	client := NewClient(cfg, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = client.GeneratePrompt(ctx, cfg.TenantHandle, PromptRequest{
		Datasets: []string{dataset},
		Prompt:   "health check",
	})
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
			t.Skipf("dataset %s unavailable: %v", dataset, err)
		}
		t.Fatalf("generate prompt: %v", err)
	}
}
