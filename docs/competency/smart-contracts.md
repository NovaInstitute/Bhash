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

1. Use the Go CLI to install tooling (first run only) and execute the SPARQL regression suite that includes this query:
   ```bash
   go run ./cmd/bhashctl install
   go run ./cmd/bhashctl sparql
   ```
2. The ROBOT-backed run writes results for `cq-dev-005.rq` to `build/queries/cq-dev-005.csv` and flags any differences from the expected fixture under `tests/fixtures/results/`.
3. Run SHACL validation with the same CLI to ensure executions and invocations include the metadata required for analytics tooling:
   ```bash
   go run ./cmd/bhashctl shacl
   ```

### Sample result (derived from `ontology/examples/smart-contracts.ttl`)

| contract | precompile | gasUsed |
| -------- | ---------- | ------- |
| `exsc:DEXContract` | `exsc:HTSPrecompile` | `42000` |

### Evidence bundle

Source citations and ETL considerations are detailed in [`docs/competency/evidence/CQ-DEV-005.md`](evidence/CQ-DEV-005.md).
