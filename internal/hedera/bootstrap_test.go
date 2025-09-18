package hedera

import (
	"context"
	"testing"
	"time"
)

func TestBootstrapperExecute(t *testing.T) {
	now := time.Date(2024, 9, 1, 12, 0, 0, 0, time.UTC)
	network := NewMockNetwork("testnet", WithStartingIDs(5000, 7000, 9000), WithNowFunc(func() time.Time { return now }))
	bootstrapper := NewBootstrapper(network, "testnet")
	spec := BootstrapSpec{
		Accounts: []AccountSpec{{Alias: "treasury", Memo: "Treasury account"}},
		Topics:   []TopicSpec{{Alias: "consensus", Memo: "Consensus topic"}},
		Tokens: []TokenSpec{{
			Alias:         "demo",
			Name:          "Demo Token",
			Symbol:        "DEM",
			TreasuryAlias: "treasury",
			SupplyType:    "INFINITE",
			TokenType:     "FUNGIBLE_COMMON",
		}},
	}
	result, err := bootstrapper.Execute(context.Background(), spec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Accounts) != 1 || result.Accounts[0].AccountID != "0.0.5000" {
		t.Fatalf("unexpected accounts: %+v", result.Accounts)
	}
	if len(result.Topics) != 1 || result.Topics[0].TopicID != "0.0.7000" {
		t.Fatalf("unexpected topics: %+v", result.Topics)
	}
	if len(result.Tokens) != 1 || result.Tokens[0].TreasuryAccountID != "0.0.5000" {
		t.Fatalf("unexpected tokens: %+v", result.Tokens)
	}
}

func TestBootstrapTransaction(t *testing.T) {
	now := time.Date(2024, 9, 1, 12, 0, 0, 0, time.UTC)
	result := BootstrapResult{
		Network: "testnet",
		Accounts: []AccountRecord{{
			Alias:     "treasury",
			AccountID: "0.0.1001",
			Memo:      "Treasury",
			CreatedAt: now,
			Tags:      []string{"governance"},
		}},
		Topics: []TopicRecord{{
			Alias:     "consensus",
			TopicID:   "0.0.2001",
			Memo:      "Consensus",
			CreatedAt: now,
		}},
		Tokens: []TokenRecord{{
			Alias:             "token",
			TokenID:           "0.0.3001",
			Name:              "Demo",
			Symbol:            "DEM",
			TreasuryAccountID: "0.0.1001",
			InitialSupply:     1000,
			CreatedAt:         now,
		}},
	}
	tx := result.Transaction("tenant/dataset")
	if tx.Ledger != "tenant/dataset" {
		t.Fatalf("unexpected ledger: %s", tx.Ledger)
	}
	if len(tx.Insert) != 3 {
		t.Fatalf("expected 3 inserts, got %d", len(tx.Insert))
	}
	accountNode := tx.Insert[0]
	if accountNode["@id"].(string) != "urn:hedera:account:0.0.1001" {
		t.Fatalf("unexpected account node: %+v", accountNode)
	}
	tokenNode := tx.Insert[2]
	if tokenNode["hedera:treasuryAccount"].(string) != "urn:hedera:account:0.0.1001" {
		t.Fatalf("expected treasury link, got %+v", tokenNode)
	}
}
