# Token Compliance Competency Answers

Phase 3 introduces executable artefacts for CQ-COMP-003 to validate HIP-540 token governance requirements.

## CQ-COMP-003 – Which tokens classified as stablecoins (HIP-540) enforce KYC and freeze keys, and who controls those keys?

* **Status:** Example data, SHACL validation, and SPARQL query committed.
* **Scope:** Identify HIP-540-aligned stablecoins, their treasury accounts, and the actors stewarding KYC and freeze keys.
* **Inputs:**
  * Ontology: `ontology/src/token.ttl`
  * Sample graph: `ontology/examples/token-compliance.ttl`
  * Query: `tests/queries/cq-comp-003.rq`
  * SHACL: `ontology/shapes/token.shacl.ttl`

### Execution notes

1. Use the Go CLI to install tooling (first run only) and execute the SPARQL regression suite that includes this query:
   ```bash
   go run ./cmd/bhashctl install
   go run ./cmd/bhashctl sparql
   ```
2. The generated CSV for `cq-comp-003.rq` appears under `build/queries/cq-comp-003.csv`, and the command reports any drift from the fixture in `tests/fixtures/results/`.
3. Run SHACL validation with the Go CLI to ensure stablecoin tokens declare both KYC and freeze controllers:
   ```bash
   go run ./cmd/bhashctl shacl
   ```

### Sample result (derived from `ontology/examples/token-compliance.ttl`)

| token | tokenLabel | treasuryAccount | kycController | freezeController |
| ----- | ---------- | ---------------- | ------------- | ---------------- |
| `ex:USDH` | "USDH Stablecoin" | `ex:IssuerAccount` | `ex:StablecoinIssuer` | `ex:CompliancePartner` |

### Evidence bundle

Modelling and data sourcing notes are detailed in
[`docs/competency/evidence/CQ-COMP-003.md`](evidence/CQ-COMP-003.md).
