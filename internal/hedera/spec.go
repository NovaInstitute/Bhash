package hedera

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// BootstrapSpec describes the artefacts that should be created during a Phase 3G
// operational bootstrap run.
type BootstrapSpec struct {
	Network  string        `json:"network"`
	Ledger   string        `json:"ledger"`
	Accounts []AccountSpec `json:"accounts"`
	Topics   []TopicSpec   `json:"topics"`
	Tokens   []TokenSpec   `json:"tokens"`
}

// AccountSpec configures how an account should be provisioned.
type AccountSpec struct {
	Alias                 string   `json:"alias"`
	Memo                  string   `json:"memo"`
	PublicKey             string   `json:"publicKey"`
	InitialBalanceTinybar int64    `json:"initialBalanceTinybar"`
	Tags                  []string `json:"tags"`
}

// TopicSpec configures a consensus topic.
type TopicSpec struct {
	Alias     string   `json:"alias"`
	Memo      string   `json:"memo"`
	AdminKey  string   `json:"adminKey"`
	SubmitKey string   `json:"submitKey"`
	Tags      []string `json:"tags"`
}

// TokenSpec configures a token to mint.
type TokenSpec struct {
	Alias             string   `json:"alias"`
	Name              string   `json:"name"`
	Symbol            string   `json:"symbol"`
	Memo              string   `json:"memo"`
	TreasuryAlias     string   `json:"treasuryAlias"`
	TreasuryAccountID string   `json:"treasuryAccountId"`
	Decimals          uint     `json:"decimals"`
	InitialSupply     uint64   `json:"initialSupply"`
	MaxSupply         int64    `json:"maxSupply"`
	SupplyType        string   `json:"supplyType"`
	TokenType         string   `json:"tokenType"`
	AdminKey          string   `json:"adminKey"`
	SupplyKey         string   `json:"supplyKey"`
	KYCKey            string   `json:"kycKey"`
	FreezeKey         string   `json:"freezeKey"`
	WipeKey           string   `json:"wipeKey"`
	PauseKey          string   `json:"pauseKey"`
	FreezeDefault     *bool    `json:"freezeDefault"`
	Tags              []string `json:"tags"`
}

// LoadBootstrapSpec reads a bootstrap specification from disk.
func LoadBootstrapSpec(path string) (BootstrapSpec, error) {
	file, err := os.Open(path)
	if err != nil {
		return BootstrapSpec{}, err
	}
	defer file.Close()

	var spec BootstrapSpec
	if err := json.NewDecoder(file).Decode(&spec); err != nil {
		return BootstrapSpec{}, fmt.Errorf("decode bootstrap spec: %w", err)
	}
	spec.Network = strings.TrimSpace(spec.Network)
	spec.Ledger = strings.TrimSpace(spec.Ledger)
	return spec, nil
}
