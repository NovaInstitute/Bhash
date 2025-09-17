# Hiero Overlay Competency Answers

Phase 3 delivers executable artefacts for CQ-HIE-009 to monitor validator onboarding across Hiero shards.

## CQ-HIE-009 – Which validators are participating in each Hiero shard and what onboarding state (pending, active, slashed) do they occupy?

* **Status:** Example data, SHACL validation, and SPARQL query committed.
* **Scope:** Record validator participation in Hiero shards, including role, onboarding state, and readiness scores.
* **Inputs:**
  * Ontology: `ontology/src/hiero.ttl`
  * Sample graph: `ontology/examples/hiero.ttl`
  * Query: `tests/queries/cq-hie-009.rq`
  * SHACL: `ontology/shapes/hiero.shacl.ttl`

### Execution notes

1. Use the Go CLI to install tooling (first run only) and execute the SPARQL regression suite containing this onboarding query:
   ```bash
   go run ./cmd/bhashctl install
   go run ./cmd/bhashctl sparql
   ```
2. The generated `cq-hie-009.csv` in `build/queries/` lists validators, shards, onboarding states, and readiness scores, with diffs reported against the stored fixture.
3. Run SHACL validation with the Go CLI to ensure each participation record captures validator, shard, onboarding state, and status text:
   ```bash
   go run ./cmd/bhashctl shacl
   ```

### Sample result (derived from `ontology/examples/hiero.ttl`)

| validator | shard | state | status | score |
| --------- | ----- | ----- | ------ | ----- |
| `exh:ValidatorAtlas` | `exh:ShardAlpha` | `Active` | `active` | `0.92` |
| `exh:ValidatorBeacon` | `exh:ShardBeta` | `Pending` | `pending` | `0.75` |

### Evidence bundle

See [`docs/competency/evidence/CQ-HIE-009.md`](evidence/CQ-HIE-009.md) for sourcing notes and future onboarding automation tasks.
