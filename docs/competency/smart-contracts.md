# Smart Contract Service Competency Answers

Phase 3 delivers executable artefacts for CQ-DEV-005 to trace Hedera Token Service (HTS) precompile usage from contract executions.

## CQ-DEV-005 – Which smart contracts invoke HTS system contract precompiles and what gas usage patterns do they exhibit?

* **Status:** Example data, SHACL validation, and SPARQL query committed.
* **Scope:** Surface contract executions that include `hedera:PrecompileInvocation` events targeting the HTS precompile along with gas consumption.
* **Inputs:**
  * Ontology: `ontology/src/smart-contracts.ttl`
  * Sample graph: `ontology/examples/smart-contracts.ttl`
  * Query: `tests/queries/cq-dev-005.rq`
  * SHACL: `ontology/shapes/smart-contracts.shacl.ttl`

### Execution notes

1. Load the ontology and example graph into an RDF store and execute the SPARQL query:
   ```bash
   arq --data ontology/src/core.ttl \
       --data ontology/src/token.ttl \
       --data ontology/src/smart-contracts.ttl \
       --data ontology/examples/smart-contracts.ttl \
       --query tests/queries/cq-dev-005.rq
   ```
2. Results list each `hedera:ContractExecution` paired with the `hedera:Precompile` it invoked and the gas consumed per invocation.
3. Run SHACL validation to ensure executions and invocations include the metadata required for analytics tooling:
   ```bash
   python -m pyshacl --data-file ontology/examples/smart-contracts.ttl \
                     --shacl-file ontology/shapes/smart-contracts.shacl.ttl \
                     --inference rdfs
   ```

### Sample result (derived from `ontology/examples/smart-contracts.ttl`)

| contract | precompile | gasUsed |
| -------- | ---------- | ------- |
| `exsc:DEXContract` | `exsc:HTSPrecompile` | `42000` |

### Evidence bundle

Source citations and ETL considerations are detailed in [`docs/competency/evidence/CQ-DEV-005.md`](evidence/CQ-DEV-005.md).
