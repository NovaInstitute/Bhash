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

1. Use the Go CLI to install tooling (first run only) and execute the SPARQL regression suite that covers this competency query:
   ```bash
   go run ./cmd/bhashctl install
   go run ./cmd/bhashctl sparql
   ```
2. The resulting `cq-comp-004.csv` under `build/queries/` lists pending schedules with missing signatures and expirations, and the command reports deviations from fixtures.
3. Validate the dataset with SHACL via the Go CLI to ensure pending schedules track signature counts and expiration timestamps:
   ```bash
   go run ./cmd/bhashctl shacl
   ```

### Sample result (derived from `ontology/examples/file-schedule.ttl`)

| schedule | status | missingSignatures | expires |
| -------- | ------ | ----------------- | ------- |
| `exfs:EmergencyPauseSchedule` | `pending-signatures` | `1` | `2025-10-01T00:00:00Z` |

### Evidence bundle

Background research and monitoring requirements are recorded in [`docs/competency/evidence/CQ-COMP-004.md`](evidence/CQ-COMP-004.md).
