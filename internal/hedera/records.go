package hedera

import "time"

// AccountRecord captures the Hedera identifiers required to reference an account.
type AccountRecord struct {
	Alias     string    `json:"alias"`
	AccountID string    `json:"accountId"`
	PublicKey string    `json:"publicKey"`
	Memo      string    `json:"memo"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
}

// TopicRecord captures metadata about a consensus topic.
type TopicRecord struct {
	Alias     string    `json:"alias"`
	TopicID   string    `json:"topicId"`
	Memo      string    `json:"memo"`
	Sequence  uint64    `json:"sequence"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
}

// TokenRecord captures metadata about a token.
type TokenRecord struct {
	Alias             string    `json:"alias"`
	TokenID           string    `json:"tokenId"`
	Name              string    `json:"name"`
	Symbol            string    `json:"symbol"`
	Memo              string    `json:"memo"`
	TreasuryAccountID string    `json:"treasuryAccountId"`
	Decimals          uint      `json:"decimals"`
	InitialSupply     uint64    `json:"initialSupply"`
	MaxSupply         int64     `json:"maxSupply"`
	SupplyType        string    `json:"supplyType"`
	TokenType         string    `json:"tokenType"`
	Tags              []string  `json:"tags"`
	CreatedAt         time.Time `json:"createdAt"`
}

// BootstrapResult aggregates the artefacts created during a bootstrap run.
type BootstrapResult struct {
	Network  string          `json:"network"`
	Accounts []AccountRecord `json:"accounts"`
	Topics   []TopicRecord   `json:"topics"`
	Tokens   []TokenRecord   `json:"tokens"`
}
