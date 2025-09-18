package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hashgraph/bhash/internal/fluree"
	bhedera "github.com/hashgraph/bhash/internal/hedera"
)

type flureeClientFunc func(context.Context, fluree.TransactionRequest) (any, error)

func (f flureeClientFunc) Transact(ctx context.Context, req fluree.TransactionRequest) (any, error) {
	return f(ctx, req)
}

func TestRunHederaBootstrap(t *testing.T) {
	tempDir := t.TempDir()
	specPath := filepath.Join(tempDir, "spec.json")
	spec := bhedera.BootstrapSpec{
		Network:  "testnet",
		Ledger:   "tenant/spec-ledger",
		Accounts: []bhedera.AccountSpec{{Alias: "treasury"}},
		Tokens: []bhedera.TokenSpec{{
			Alias:         "token",
			Name:          "Demo",
			Symbol:        "DEM",
			TreasuryAlias: "treasury",
		}},
	}
	file, err := os.Create(specPath)
	if err != nil {
		t.Fatalf("create spec: %v", err)
	}
	if err := json.NewEncoder(file).Encode(spec); err != nil {
		t.Fatalf("encode spec: %v", err)
	}
	file.Close()

	originalFactory := hederaNetworkFactory
	defer func() { hederaNetworkFactory = originalFactory }()
	fixedTime := time.Date(2024, 9, 1, 12, 0, 0, 0, time.UTC)
	hederaNetworkFactory = func(cfg bhedera.Config, simulate bool) (bhedera.Network, func(), error) {
		network := bhedera.NewMockNetwork(cfg.Network, bhedera.WithStartingIDs(1000, 2000, 3000), bhedera.WithNowFunc(func() time.Time { return fixedTime }))
		return network, func() {}, nil
	}

	originalFluree := flureeClientFactory
	defer func() { flureeClientFactory = originalFluree }()
	var capturedLedger string
	flureeClientFactory = func(cfg fluree.Config) flureeTransactor {
		return flureeClientFunc(func(ctx context.Context, req fluree.TransactionRequest) (any, error) {
			capturedLedger = req.Ledger
			return map[string]any{"status": "ok"}, nil
		})
	}

	buf := &bytes.Buffer{}
	originalWriter := outputWriter
	outputWriter = buf
	defer func() { outputWriter = originalWriter }()

	runHederaBootstrap([]string{"--spec", specPath, "--ledger", "tenant/dataset"})

	var output map[string]any
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Fatalf("decode output: %v", err)
	}
	if output["ledger"].(string) != "tenant/dataset" {
		t.Fatalf("unexpected ledger: %#v", output["ledger"])
	}
	accounts := output["accounts"].([]any)
	if len(accounts) != 1 {
		t.Fatalf("expected 1 account, got %d", len(accounts))
	}
	transaction := output["transaction"].(map[string]any)
	inserts := transaction["Insert"].([]any)
	if len(inserts) != 2 {
		t.Fatalf("expected 2 inserts, got %d", len(inserts))
	}
	if capturedLedger != "" {
		t.Fatalf("did not expect fluree transaction when commit=false")
	}
}

func TestRunHederaBootstrapCommit(t *testing.T) {
	tempDir := t.TempDir()
	specPath := filepath.Join(tempDir, "spec.json")
	spec := bhedera.BootstrapSpec{Accounts: []bhedera.AccountSpec{{Alias: "treasury"}}}
	file, err := os.Create(specPath)
	if err != nil {
		t.Fatalf("create spec: %v", err)
	}
	if err := json.NewEncoder(file).Encode(spec); err != nil {
		t.Fatalf("encode spec: %v", err)
	}
	file.Close()

	originalFactory := hederaNetworkFactory
	defer func() { hederaNetworkFactory = originalFactory }()
	hederaNetworkFactory = func(cfg bhedera.Config, simulate bool) (bhedera.Network, func(), error) {
		return bhedera.NewMockNetwork(cfg.Network, bhedera.WithStartingIDs(1, 1, 1)), func() {}, nil
	}

	var committed bool
	originalFluree := flureeClientFactory
	defer func() { flureeClientFactory = originalFluree }()
	flureeClientFactory = func(cfg fluree.Config) flureeTransactor {
		return flureeClientFunc(func(ctx context.Context, req fluree.TransactionRequest) (any, error) {
			committed = true
			return map[string]any{"status": "ok"}, nil
		})
	}

	t.Setenv("FLUREE_API_TOKEN", "env-token")
	t.Setenv("FLUREE_HANDLE", "env-tenant")
	t.Setenv("FLUREE_BASE_URL", "http://env.example")

	buf := &bytes.Buffer{}
	originalWriter := outputWriter
	outputWriter = buf
	defer func() { outputWriter = originalWriter }()

	runHederaBootstrap([]string{"--spec", specPath, "--ledger", "tenant/dataset", "--commit", "--api-token", "token", "--tenant", "tenant", "--base-url", "http://example"})

	if !committed {
		t.Fatalf("expected fluree transact to be called")
	}
}
