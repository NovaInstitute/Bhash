package hedera

import (
	"fmt"
	"time"

	"github.com/hashgraph/bhash/internal/fluree"
)

var defaultContext = map[string]any{
	"@vocab":                  "https://hashgraph.github.io/bhash/hedera#",
	"hedera":                  "https://hashgraph.github.io/bhash/hedera#",
	"prov":                    "http://www.w3.org/ns/prov#",
	"schema":                  "http://schema.org/",
	"xsd":                     "http://www.w3.org/2001/XMLSchema#",
	"prov:generatedAtTime":    map[string]any{"@type": "xsd:dateTime"},
	"hedera:belongsToNetwork": map[string]any{"@type": "xsd:string"},
	"hedera:accountId":        map[string]any{"@type": "xsd:string"},
	"hedera:topicId":          map[string]any{"@type": "xsd:string"},
	"hedera:tokenId":          map[string]any{"@type": "xsd:string"},
	"hedera:treasuryAccount":  map[string]any{"@type": "@id"},
	"schema:keywords":         map[string]any{"@container": "@set"},
}

// Transaction builds a Fluree transaction that inserts JSON-LD nodes for every
// artefact recorded in the result.
func (r BootstrapResult) Transaction(ledger string) fluree.TransactionRequest {
	ctx := make(map[string]any, len(defaultContext))
	for k, v := range defaultContext {
		ctx[k] = v
	}

	req := fluree.TransactionRequest{Ledger: ledger, Context: ctx}
	for _, account := range r.Accounts {
		req.Insert = append(req.Insert, account.asJSONLD(r.Network))
	}
	for _, topic := range r.Topics {
		req.Insert = append(req.Insert, topic.asJSONLD(r.Network))
	}
	for _, token := range r.Tokens {
		req.Insert = append(req.Insert, token.asJSONLD(r.Network))
	}
	return req
}

func (a AccountRecord) asJSONLD(network string) map[string]any {
	node := map[string]any{
		"@id":                     urn("account", a.AccountID),
		"@type":                   []string{"hedera:Account", "prov:Agent"},
		"hedera:accountId":        a.AccountID,
		"hedera:belongsToNetwork": network,
	}
	if !a.CreatedAt.IsZero() {
		node["prov:generatedAtTime"] = formatTime(a.CreatedAt)
	}
	if a.Alias != "" {
		node["schema:name"] = a.Alias
	}
	if a.Memo != "" {
		node["schema:description"] = a.Memo
	}
	if a.PublicKey != "" {
		node["hedera:publicKey"] = a.PublicKey
	}
	if len(a.Tags) > 0 {
		node["schema:keywords"] = append([]string(nil), a.Tags...)
	}
	return node
}

func (t TopicRecord) asJSONLD(network string) map[string]any {
	node := map[string]any{
		"@id":                     urn("topic", t.TopicID),
		"@type":                   []string{"hedera:ConsensusTopic", "prov:Entity"},
		"hedera:topicId":          t.TopicID,
		"hedera:belongsToNetwork": network,
	}
	if !t.CreatedAt.IsZero() {
		node["prov:generatedAtTime"] = formatTime(t.CreatedAt)
	}
	if t.Memo != "" {
		node["schema:description"] = t.Memo
	}
	if len(t.Tags) > 0 {
		node["schema:keywords"] = append([]string(nil), t.Tags...)
	}
	if t.Sequence > 0 {
		node["hedera:initialSequence"] = t.Sequence
	}
	return node
}

func (t TokenRecord) asJSONLD(network string) map[string]any {
	node := map[string]any{
		"@id":                     urn("token", t.TokenID),
		"@type":                   []string{"hedera:Token", "prov:Entity"},
		"hedera:tokenId":          t.TokenID,
		"hedera:belongsToNetwork": network,
	}
	if t.Name != "" {
		node["schema:name"] = t.Name
	}
	if t.Symbol != "" {
		node["schema:identifier"] = t.Symbol
	}
	if !t.CreatedAt.IsZero() {
		node["prov:generatedAtTime"] = formatTime(t.CreatedAt)
	}
	if t.Memo != "" {
		node["schema:description"] = t.Memo
	}
	if t.TreasuryAccountID != "" {
		node["hedera:treasuryAccount"] = urn("account", t.TreasuryAccountID)
	}
	if t.Decimals > 0 {
		node["hedera:decimals"] = t.Decimals
	}
	if t.InitialSupply > 0 {
		node["hedera:initialSupply"] = t.InitialSupply
	}
	if t.MaxSupply != 0 {
		node["hedera:maxSupply"] = t.MaxSupply
	}
	if t.SupplyType != "" {
		node["hedera:supplyType"] = t.SupplyType
	}
	if t.TokenType != "" {
		node["hedera:tokenType"] = t.TokenType
	}
	if len(t.Tags) > 0 {
		node["schema:keywords"] = append([]string(nil), t.Tags...)
	}
	return node
}

func urn(kind, id string) string {
	return fmt.Sprintf("urn:hedera:%s:%s", kind, id)
}

func formatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}
