package flureeclient_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashgraph/bhash/internal/fluree"
	"github.com/hashgraph/bhash/scripts/flureeclient"
)

func TestGeneratePromptLogsAndSetsHeaders(t *testing.T) {
	t.Parallel()

	var capturedBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer token" {
			t.Fatalf("unexpected Authorization header: %q", got)
		}
		if got := r.Header.Get("x-user-handle"); got != "tenant" {
			t.Fatalf("unexpected x-user-handle header: %q", got)
		}
		var err error
		capturedBody, err = io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"prompt": "SELECT *"})
	}))
	defer server.Close()

	cfg := fluree.Config{APIToken: "token", TenantHandle: "tenant", BaseURL: server.URL}
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	client := flureeclient.New(cfg, logger, nil)

	resp, err := client.GeneratePrompt(context.Background(), "tenant", fluree.PromptRequest{
		Datasets: []string{"tenant/sample"},
		Prompt:   "List accounts",
	})
	if err != nil {
		t.Fatalf("generate prompt: %v", err)
	}
	if !strings.Contains(logBuf.String(), "generate-prompt") {
		t.Fatalf("expected logging output, got %q", logBuf.String())
	}

	var payload map[string]any
	if err := json.Unmarshal(capturedBody, &payload); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	datasets := payload["datasets"].([]any)
	if len(datasets) != 1 || datasets[0].(string) != "tenant/sample" {
		t.Fatalf("unexpected datasets payload: %#v", payload["datasets"])
	}

	decoded := resp.(map[string]any)
	if decoded["prompt"].(string) != "SELECT *" {
		t.Fatalf("unexpected response payload: %#v", decoded)
	}
}

func TestNewFromEnvMissingSecrets(t *testing.T) {
	t.Parallel()

	_, err := flureeclient.NewFromEnv(func(string) (string, bool) { return "", false }, nil, nil)
	if err == nil {
		t.Fatalf("expected error when environment variables are missing")
	}
	if !strings.Contains(err.Error(), "FLUREE_API_TOKEN") {
		t.Fatalf("expected missing token in error: %v", err)
	}
}
