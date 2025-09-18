package hedera

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

// MockOption customises the behaviour of the mock network.
type MockOption func(*MockNetwork)

// WithStartingIDs sets the initial identifier counters for the mock network.
func WithStartingIDs(account, topic, token int64) MockOption {
	return func(m *MockNetwork) {
		m.nextAccount = account
		m.nextTopic = topic
		m.nextToken = token
	}
}

// WithNowFunc allows overriding the clock used when generating metadata.
func WithNowFunc(fn func() time.Time) MockOption {
	return func(m *MockNetwork) {
		m.now = fn
	}
}

// MockNetwork provides a deterministic in-memory implementation of the Network
// interface. It is used for local development and tests when a live Hedera
// network is not available.
type MockNetwork struct {
	networkName string
	mu          sync.Mutex
	nextAccount int64
	nextTopic   int64
	nextToken   int64
	now         func() time.Time
}

// NewMockNetwork constructs a mock network with optional customisation.
func NewMockNetwork(network string, opts ...MockOption) *MockNetwork {
	m := &MockNetwork{
		networkName: network,
		nextAccount: 1000,
		nextTopic:   2000,
		nextToken:   3000,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *MockNetwork) CreateAccount(_ context.Context, spec AccountSpec) (AccountRecord, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := m.nextAccount
	m.nextAccount++
	record := AccountRecord{
		Alias:     spec.Alias,
		AccountID: fmt.Sprintf("0.0.%d", id),
		PublicKey: spec.PublicKey,
		Memo:      spec.Memo,
		Tags:      append([]string(nil), spec.Tags...),
		CreatedAt: m.now(),
	}
	return record, nil
}

func (m *MockNetwork) CreateTopic(_ context.Context, spec TopicSpec) (TopicRecord, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := m.nextTopic
	m.nextTopic++
	record := TopicRecord{
		Alias:     spec.Alias,
		TopicID:   fmt.Sprintf("0.0.%d", id),
		Memo:      spec.Memo,
		Tags:      append([]string(nil), spec.Tags...),
		CreatedAt: m.now(),
	}
	return record, nil
}

func (m *MockNetwork) CreateToken(_ context.Context, spec TokenSpec) (TokenRecord, error) {
	if strings.TrimSpace(spec.TreasuryAccountID) == "" {
		return TokenRecord{}, fmt.Errorf("treasury account id is required for token %q", spec.Alias)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	id := m.nextToken
	m.nextToken++
	record := TokenRecord{
		Alias:             spec.Alias,
		TokenID:           fmt.Sprintf("0.0.%d", id),
		Name:              spec.Name,
		Symbol:            spec.Symbol,
		Memo:              spec.Memo,
		TreasuryAccountID: spec.TreasuryAccountID,
		Decimals:          spec.Decimals,
		InitialSupply:     spec.InitialSupply,
		MaxSupply:         spec.MaxSupply,
		SupplyType:        spec.SupplyType,
		TokenType:         spec.TokenType,
		Tags:              append([]string(nil), spec.Tags...),
		CreatedAt:         m.now(),
	}
	return record, nil
}
