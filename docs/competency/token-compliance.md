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

1. Load the ontology and example graph into an RDF store and execute the SPARQL query:
   ```bash
   arq --data ontology/src/core.ttl \
       --data ontology/src/token.ttl \
       --data ontology/examples/token-compliance.ttl \
       --query tests/queries/cq-comp-003.rq
   ```
2. Results enumerate each `hedera:StablecoinToken` alongside treasury accounts and controller actors who hold
   `hedera:KYCControllerRole` and `hedera:FreezeControllerRole` assignments.
3. Run SHACL validation to ensure stablecoin tokens declare both KYC and freeze controllers:
   ```bash
   python -m pyshacl --data-file ontology/examples/token-compliance.ttl \
                     --shacl-file ontology/shapes/token.shacl.ttl \
                     --inference rdfs
   ```

### Sample result (derived from `ontology/examples/token-compliance.ttl`)

| token | tokenLabel | treasuryAccount | kycController | freezeController |
| ----- | ---------- | ---------------- | ------------- | ---------------- |
| `ex:USDH` | "USDH Stablecoin" | `ex:IssuerAccount` | `ex:StablecoinIssuer` | `ex:CompliancePartner` |

### Evidence bundle

Modelling and data sourcing notes are detailed in
[`docs/competency/evidence/CQ-COMP-003.md`](evidence/CQ-COMP-003.md).
