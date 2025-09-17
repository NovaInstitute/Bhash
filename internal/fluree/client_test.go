package fluree

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateDataset(t *testing.T) {
	var received map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/owner/create-dataset" {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer token" {
			t.Fatalf("unexpected Authorization header %q", got)
		}
		if got := r.Header.Get("x-user-handle"); got != "tenant" {
			t.Fatalf("unexpected tenant header %q", got)
		}
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	cfg := Config{APIToken: "token", TenantHandle: "tenant", BaseURL: server.URL}
	client := NewClient(cfg, server.Client())
	resp, err := client.CreateDataset(context.Background(), "owner", CreateDatasetRequest{
		DatasetName: "example",
		StorageType: "sparql",
		Description: "demo",
		Visibility:  "private",
		Tags:        []string{"alpha"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["datasetName"] != "example" || received["visibility"] != "private" {
		t.Fatalf("unexpected payload %+v", received)
	}
	data, ok := resp.(map[string]any)
	if !ok || data["ok"].(bool) != true {
		t.Fatalf("unexpected response %#v", resp)
	}
}

func TestTransactError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"boom"}`))
	}))
	defer server.Close()

	cfg := Config{APIToken: "token", TenantHandle: "tenant", BaseURL: server.URL}
	client := NewClient(cfg, server.Client())
	_, err := client.Transact(context.Background(), TransactionRequest{Ledger: "ledger"})
	if err == nil {
		t.Fatalf("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Message != "boom" {
		t.Fatalf("unexpected message %q", apiErr.Message)
	}
}

func TestGenerateAnswer(t *testing.T) {
	var received map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/owner/generate-answer" {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		defer r.Body.Close()
		json.NewDecoder(r.Body).Decode(&received)
		w.Write([]byte(`{"answer":"42"}`))
	}))
	defer server.Close()

	cfg := Config{APIToken: "token", TenantHandle: "tenant", BaseURL: server.URL}
	client := NewClient(cfg, server.Client())
	resp, err := client.GenerateAnswer(context.Background(), "owner", PromptRequest{
		Datasets: []string{"a", "b"},
		Prompt:   "question",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["prompt"] != "question" {
		t.Fatalf("unexpected payload %#v", received)
	}
	data := resp.(map[string]any)
	if data["answer"].(string) != "42" {
		t.Fatalf("unexpected response %#v", data)
	}
}
