package hedera

import (
	"context"
	"fmt"
)

// Network exposes the subset of Hedera SDK functionality required by the
// bootstrap workflow.
type Network interface {
	CreateAccount(context.Context, AccountSpec) (AccountRecord, error)
	CreateTopic(context.Context, TopicSpec) (TopicRecord, error)
	CreateToken(context.Context, TokenSpec) (TokenRecord, error)
}

// Bootstrapper orchestrates creation of Hedera artefacts before exporting the
// resulting metadata to Fluree.
type Bootstrapper struct {
	network     Network
	networkName string
}

// NewBootstrapper returns a Bootstrapper backed by the supplied network implementation.
func NewBootstrapper(network Network, networkName string) *Bootstrapper {
	return &Bootstrapper{network: network, networkName: networkName}
}

// Execute provisions the artefacts described by spec and returns the metadata
// required to build a Fluree transaction.
func (b *Bootstrapper) Execute(ctx context.Context, spec BootstrapSpec) (BootstrapResult, error) {
	result := BootstrapResult{Network: b.networkName}
	if spec.Network != "" {
		result.Network = spec.Network
	}

	accountByAlias := make(map[string]AccountRecord)
	for _, account := range spec.Accounts {
		record, err := b.network.CreateAccount(ctx, account)
		if err != nil {
			return BootstrapResult{}, fmt.Errorf("create account %q: %w", account.Alias, err)
		}
		if record.Alias == "" {
			record.Alias = account.Alias
		}
		result.Accounts = append(result.Accounts, record)
		if record.Alias != "" {
			accountByAlias[record.Alias] = record
		}
	}

	for _, topic := range spec.Topics {
		record, err := b.network.CreateTopic(ctx, topic)
		if err != nil {
			return BootstrapResult{}, fmt.Errorf("create topic %q: %w", topic.Alias, err)
		}
		if record.Alias == "" {
			record.Alias = topic.Alias
		}
		result.Topics = append(result.Topics, record)
	}

	for _, token := range spec.Tokens {
		resolved := token
		if resolved.TreasuryAccountID == "" && resolved.TreasuryAlias != "" {
			account, ok := accountByAlias[resolved.TreasuryAlias]
			if !ok {
				return BootstrapResult{}, fmt.Errorf("treasury alias %q not found", resolved.TreasuryAlias)
			}
			resolved.TreasuryAccountID = account.AccountID
		}
		record, err := b.network.CreateToken(ctx, resolved)
		if err != nil {
			return BootstrapResult{}, fmt.Errorf("create token %q: %w", token.Alias, err)
		}
		if record.Alias == "" {
			record.Alias = token.Alias
		}
		result.Tokens = append(result.Tokens, record)
	}

	return result, nil
}
