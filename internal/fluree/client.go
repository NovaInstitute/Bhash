package fluree

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// Client provides strongly typed helpers for interacting with the Fluree HTTP
// API.
type Client struct {
	httpClient *http.Client
	config     Config
}

// NewClient returns a Client configured with the supplied credentials. When
// httpClient is nil the default http.Client with a 30s timeout is used.
func NewClient(cfg Config, httpClient *http.Client) *Client {
	client := httpClient
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	return &Client{httpClient: client, config: cfg}
}

// APIError captures the HTTP status code and message returned by Fluree.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("fluree: unexpected HTTP status %d", e.StatusCode)
	}
	return fmt.Sprintf("fluree: %s (status %d)", e.Message, e.StatusCode)
}

func (c *Client) doPost(ctx context.Context, endpoint string, payload any) (any, error) {
	base, err := url.Parse(c.config.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("fluree: invalid base URL %q: %w", c.config.BaseURL, err)
	}
	base.Path = path.Join(strings.TrimSuffix(base.Path, "/"), endpoint)

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("fluree: encode payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("fluree: build request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIToken))
	req.Header.Set("x-user-handle", c.config.TenantHandle)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fluree: perform request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fluree: read response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, &APIError{StatusCode: resp.StatusCode, Message: parseErrorMessage(data, resp.Status)}
	}
	if len(data) == 0 {
		return nil, nil
	}
	var decoded any
	if err := json.Unmarshal(data, &decoded); err != nil {
		return string(data), nil
	}
	return decoded, nil
}

func parseErrorMessage(data []byte, fallback string) string {
	var decoded any
	if err := json.Unmarshal(data, &decoded); err != nil {
		return fallback
	}
	switch v := decoded.(type) {
	case map[string]any:
		if msg, ok := v["message"].(string); ok && msg != "" {
			return msg
		}
		if msg, ok := v["error"].(string); ok && msg != "" {
			return msg
		}
		return fmt.Sprint(v)
	case string:
		return v
	default:
		return fmt.Sprint(v)
	}
}

// CreateDatasetRequest describes the payload sent to the create-dataset endpoint.
type CreateDatasetRequest struct {
	DatasetName string
	StorageType string
	Description string
	Visibility  string
	Tags        []string
}

// CreateDataset creates a dataset owned by the provided Fluree handle.
func (c *Client) CreateDataset(ctx context.Context, ownerHandle string, req CreateDatasetRequest) (any, error) {
	if ownerHandle == "" {
		return nil, fmt.Errorf("fluree: owner handle is required")
	}
	payload := map[string]any{
		"datasetName": req.DatasetName,
		"storageType": req.StorageType,
		"description": req.Description,
	}
	if req.Visibility != "" {
		payload["visibility"] = req.Visibility
	}
	if len(req.Tags) > 0 {
		payload["tags"] = req.Tags
	}
	endpoint := path.Join("api", ownerHandle, "create-dataset")
	return c.doPost(ctx, endpoint, payload)
}

// TransactionRequest represents the payload for a Fluree transact request.
type TransactionRequest struct {
	Ledger  string
	Insert  []map[string]any
	Delete  []map[string]any
	Where   []map[string]any
	Context map[string]any
}

// Transact executes a ledger transaction.
func (c *Client) Transact(ctx context.Context, req TransactionRequest) (any, error) {
	if strings.TrimSpace(req.Ledger) == "" {
		return nil, fmt.Errorf("fluree: ledger is required")
	}
	payload := map[string]any{"ledger": req.Ledger}
	if req.Context != nil {
		payload["context"] = req.Context
	}
	if len(req.Insert) > 0 {
		payload["insert"] = req.Insert
	}
	if len(req.Delete) > 0 {
		payload["delete"] = req.Delete
	}
	if len(req.Where) > 0 {
		payload["where"] = req.Where
	}
	return c.doPost(ctx, "fluree/transact", payload)
}

// PromptRequest describes a request that renders natural language responses.
type PromptRequest struct {
	Datasets []string
	Prompt   string
}

// GeneratePrompt calls the Fluree generate-prompt endpoint.
func (c *Client) GeneratePrompt(ctx context.Context, ownerHandle string, req PromptRequest) (any, error) {
	return c.generateHelper(ctx, ownerHandle, "generate-prompt", req)
}

// GenerateSPARQL calls the Fluree generate-sparql endpoint.
func (c *Client) GenerateSPARQL(ctx context.Context, ownerHandle string, req PromptRequest) (any, error) {
	return c.generateHelper(ctx, ownerHandle, "generate-sparql", req)
}

// GenerateAnswer calls the Fluree generate-answer endpoint.
func (c *Client) GenerateAnswer(ctx context.Context, ownerHandle string, req PromptRequest) (any, error) {
	return c.generateHelper(ctx, ownerHandle, "generate-answer", req)
}

func (c *Client) generateHelper(ctx context.Context, ownerHandle, suffix string, req PromptRequest) (any, error) {
	if ownerHandle == "" {
		return nil, fmt.Errorf("fluree: owner handle is required")
	}
	payload := map[string]any{
		"datasets": req.Datasets,
		"prompt":   req.Prompt,
	}
	endpoint := path.Join("api", ownerHandle, suffix)
	return c.doPost(ctx, endpoint, payload)
}
