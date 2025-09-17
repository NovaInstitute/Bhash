# File and Schedule Service Competency Answers

Phase 3 introduces executable artefacts for CQ-COMP-004 to monitor deferred transactions awaiting signatures.

## CQ-COMP-004 – Which scheduled transactions remain pending beyond 24 hours due to missing signatures?

* **Status:** Example data, SHACL validation, and SPARQL query committed.
* **Scope:** Identify `hedera:ScheduledTransaction` instances marked `pending-signatures`, quantify missing approvals, and surface expirations.
* **Inputs:**
  * Ontology: `ontology/src/file-schedule.ttl`
  * Sample graph: `ontology/examples/file-schedule.ttl`
  * Query: `tests/queries/cq-comp-004.rq`
  * SHACL: `ontology/shapes/file-schedule.shacl.ttl`

### Execution notes

1. Execute the SPARQL query to list pending schedules with missing signatures and expirations:
   ```bash
   arq --data ontology/src/core.ttl \
       --data ontology/src/file-schedule.ttl \
       --data ontology/src/smart-contracts.ttl \
       --data ontology/examples/file-schedule.ttl \
       --query tests/queries/cq-comp-004.rq
   ```
2. The query computes `missingSignatures` by subtracting collected signatures from the required total, enabling alerting pipelines.
3. Validate the dataset with SHACL to ensure pending schedules track signature counts and expiration timestamps:
   ```bash
   python -m pyshacl --data-file ontology/examples/file-schedule.ttl \
                     --shacl-file ontology/shapes/file-schedule.shacl.ttl \
                     --inference rdfs
   ```

### Sample result (derived from `ontology/examples/file-schedule.ttl`)

| schedule | status | missingSignatures | expires |
| -------- | ------ | ----------------- | ------- |
| `exfs:EmergencyPauseSchedule` | `pending-signatures` | `1` | `2025-10-01T00:00:00Z` |

### Evidence bundle

Background research and monitoring requirements are recorded in [`docs/competency/evidence/CQ-COMP-004.md`](evidence/CQ-COMP-004.md).
