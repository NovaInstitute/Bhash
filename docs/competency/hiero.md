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

1. Execute the SPARQL query to list validators, shards, onboarding state labels, and readiness scores:
   ```bash
   arq --data ontology/src/core.ttl \
       --data ontology/src/hiero.ttl \
       --data ontology/examples/hiero.ttl \
       --query tests/queries/cq-hie-009.rq
   ```
2. The query outputs both raw onboarding state IRIs and human-readable labels, supporting dashboards that visualise shard readiness.
3. Run SHACL validation to ensure each participation record captures validator, shard, onboarding state, and status text:
   ```bash
   python -m pyshacl --data-file ontology/examples/hiero.ttl \
                     --shacl-file ontology/shapes/hiero.shacl.ttl \
                     --inference rdfs
   ```

### Sample result (derived from `ontology/examples/hiero.ttl`)

| validator | shard | state | status | score |
| --------- | ----- | ----- | ------ | ----- |
| `exh:ValidatorAtlas` | `exh:ShardAlpha` | `Active` | `active` | `0.92` |
| `exh:ValidatorBeacon` | `exh:ShardBeta` | `Pending` | `pending` | `0.75` |

### Evidence bundle

See [`docs/competency/evidence/CQ-HIE-009.md`](evidence/CQ-HIE-009.md) for sourcing notes and future onboarding automation tasks.
