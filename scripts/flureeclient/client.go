package flureeclient

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/hashgraph/bhash/internal/fluree"
)

// Client wraps the internal Fluree HTTP client with structured logging.
type Client struct {
	logger *slog.Logger
	inner  *fluree.Client
}

// New returns a logging-aware client using the provided configuration.
// If logger is nil a logger writing to io.Discard is used. When httpClient
// is nil the default HTTP client is wrapped with logging instrumentation.
func New(cfg fluree.Config, logger *slog.Logger, httpClient *http.Client) *Client {
	log := logger
	if log == nil {
		log = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	client := wrapHTTPClient(httpClient, log)
	return &Client{
		logger: log,
		inner:  fluree.NewClient(cfg, client),
	}
}

// NewFromEnv builds a Config from environment variables using the supplied
// lookup function and returns a logging-aware client. The lookup function
// mirrors os.LookupEnv to simplify testing.
func NewFromEnv(lookup func(string) (string, bool), logger *slog.Logger, httpClient *http.Client) (*Client, error) {
	cfg, err := fluree.EnvConfigFromLookup(lookup)
	if err != nil {
		return nil, err
	}
	return New(cfg, logger, httpClient), nil
}

// CreateDataset proxies to the underlying Fluree client and emits request
// metadata helpful for debugging.
func (c *Client) CreateDataset(ctx context.Context, ownerHandle string, req fluree.CreateDatasetRequest) (any, error) {
	start := time.Now()
	c.logger.Info("fluree create-dataset", "owner", ownerHandle, "dataset", req.DatasetName)
	resp, err := c.inner.CreateDataset(ctx, ownerHandle, req)
	c.logResult("create-dataset", start, err)
	return resp, err
}

// Transact proxies to the Fluree client with request/response logging.
func (c *Client) Transact(ctx context.Context, req fluree.TransactionRequest) (any, error) {
	start := time.Now()
	c.logger.Info(
		"fluree transact",
		"ledger", req.Ledger,
		"insert", len(req.Insert),
		"delete", len(req.Delete),
		"where", len(req.Where),
	)
	resp, err := c.inner.Transact(ctx, req)
	c.logResult("transact", start, err)
	return resp, err
}

// GeneratePrompt proxies to the generate-prompt endpoint and logs request metadata.
func (c *Client) GeneratePrompt(ctx context.Context, ownerHandle string, req fluree.PromptRequest) (any, error) {
	start := time.Now()
	c.logger.Info("fluree generate-prompt", "owner", ownerHandle, "datasets", req.Datasets)
	resp, err := c.inner.GeneratePrompt(ctx, ownerHandle, req)
	c.logResult("generate-prompt", start, err)
	return resp, err
}

// GenerateSPARQL proxies to the generate-sparql endpoint.
func (c *Client) GenerateSPARQL(ctx context.Context, ownerHandle string, req fluree.PromptRequest) (any, error) {
	start := time.Now()
	c.logger.Info("fluree generate-sparql", "owner", ownerHandle, "datasets", req.Datasets)
	resp, err := c.inner.GenerateSPARQL(ctx, ownerHandle, req)
	c.logResult("generate-sparql", start, err)
	return resp, err
}

// GenerateAnswer proxies to the generate-answer endpoint.
func (c *Client) GenerateAnswer(ctx context.Context, ownerHandle string, req fluree.PromptRequest) (any, error) {
	start := time.Now()
	c.logger.Info("fluree generate-answer", "owner", ownerHandle, "datasets", req.Datasets)
	resp, err := c.inner.GenerateAnswer(ctx, ownerHandle, req)
	c.logResult("generate-answer", start, err)
	return resp, err
}

func (c *Client) logResult(operation string, start time.Time, err error) {
	duration := time.Since(start)
	if err != nil {
		c.logger.Error("fluree request failed", "operation", operation, "duration", duration, "error", err)
		return
	}
	c.logger.Info("fluree request completed", "operation", operation, "duration", duration)
}

// wrapHTTPClient ensures loggingRoundTripper is wired even when the caller
// supplies nil (to use the default HTTP client).
func wrapHTTPClient(client *http.Client, logger *slog.Logger) *http.Client {
	if client == nil {
		return &http.Client{Timeout: 30 * time.Second, Transport: loggingRoundTripper{logger: logger, next: http.DefaultTransport}}
	}
	transport := client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	client.Transport = loggingRoundTripper{logger: logger, next: transport}
	if client.Timeout == 0 {
		client.Timeout = 30 * time.Second
	}
	return client
}

type loggingRoundTripper struct {
	logger *slog.Logger
	next   http.RoundTripper
}

func (l loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := l.next
	if transport == nil {
		transport = http.DefaultTransport
	}
	start := time.Now()
	l.logger.Debug("fluree http request", "method", req.Method, "url", req.URL.Redacted())
	resp, err := transport.RoundTrip(req)
	if err != nil {
		l.logger.Error("fluree http error", "method", req.Method, "url", req.URL.Redacted(), "error", err)
		return nil, err
	}
	l.logger.Debug("fluree http response", "status", resp.StatusCode, "duration", time.Since(start))
	return resp, nil
}
