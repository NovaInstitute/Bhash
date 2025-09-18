package fluree

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeneratePromptPostsPayload(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/tenant/generate-prompt" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer token" {
			t.Fatalf("unexpected Authorization header: %q", got)
		}
		if got := r.Header.Get("x-user-handle"); got != "tenant" {
			t.Fatalf("unexpected handle header: %q", got)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		datasets := payload["datasets"].([]any)
		if len(datasets) != 1 || datasets[0] != "tenant/sample" {
			t.Fatalf("unexpected datasets: %#v", datasets)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"prompt": "SELECT *"})
	}))
	defer server.Close()

	cfg := Config{APIToken: "token", TenantHandle: "tenant", BaseURL: server.URL}
	client := NewClient(cfg, server.Client())

	resp, err := client.GeneratePrompt(context.Background(), "tenant", PromptRequest{
		Datasets: []string{"tenant/sample"},
		Prompt:   "List",
	})
	if err != nil {
		t.Fatalf("generate prompt: %v", err)
	}
	decoded := resp.(map[string]any)
	if decoded["prompt"].(string) != "SELECT *" {
		t.Fatalf("unexpected response: %#v", decoded)
	}
}

func TestTransactIncludesOptionalBlocks(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["ledger"].(string) != "tenant/sample" {
			t.Fatalf("unexpected ledger: %s", payload["ledger"])
		}
		if _, ok := payload["context"]; !ok {
			t.Fatalf("expected context block: %#v", payload)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
	}))
	defer server.Close()

	cfg := Config{APIToken: "token", TenantHandle: "tenant", BaseURL: server.URL}
	client := NewClient(cfg, server.Client())

	_, err := client.Transact(context.Background(), TransactionRequest{
		Ledger:  "tenant/sample",
		Context: map[string]any{"@context": map[string]any{"ex": "https://example.com/"}},
		Insert: []map[string]any{
			{"@id": "ex:thing", "@type": "ex:Class"},
		},
		Delete: []map[string]any{},
		Where: []map[string]any{
			{"@id": "ex:thing"},
		},
	})
	if err != nil {
		t.Fatalf("transact: %v", err)
	}
}

func TestErrorResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "No dataset"})
	}))
	defer server.Close()

	cfg := Config{APIToken: "token", TenantHandle: "tenant", BaseURL: server.URL}
	client := NewClient(cfg, server.Client())

	_, err := client.GenerateSPARQL(context.Background(), "tenant", PromptRequest{Datasets: []string{"missing"}})
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %v", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected status: %d", apiErr.StatusCode)
	}
}
