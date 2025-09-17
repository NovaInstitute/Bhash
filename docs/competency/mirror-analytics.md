# Mirror and Analytics Competency Answers

Phase 3 delivers executable artefacts for CQ-ANL-007 to track dataset retention supporting treasury analytics.

## CQ-ANL-007 – How is total supply of fungible tokens distributed across treasury and external accounts over time?

* **Status:** Example data, SHACL validation, and SPARQL query committed for retention coverage; quantitative distribution analytics remain for future data drops.
* **Scope:** Ensure mirror datasets powering treasury dashboards declare retention policies and workspace consumers so that time-series analysis can be automated once live data arrives.
* **Inputs:**
  * Ontology: `ontology/src/mirror-analytics.ttl`
  * Sample graph: `ontology/examples/mirror-analytics.ttl`
  * Query: `tests/queries/cq-anl-007.rq`
  * SHACL: `ontology/shapes/mirror-analytics.shacl.ttl`

### Execution notes

1. Use the Go CLI to install tooling (first run only) and execute the SPARQL regression suite that includes this retention query:
   ```bash
   go run ./cmd/bhashctl install
   go run ./cmd/bhashctl sparql
   ```
2. The generated `cq-anl-007.csv` in `build/queries/` enumerates mirror datasets, covered services, retention windows, and downstream workspaces, with diffs surfaced against fixtures.
3. Use SHACL validation via the Go CLI to ensure datasets expose retention metadata required by compliance teams:
   ```bash
   go run ./cmd/bhashctl shacl
   ```
4. Future automation will layer token balance aggregates on top of the validated retention model.

### Sample result (derived from `ontology/examples/mirror-analytics.ttl`)

| dataset | retentionDays | workspace |
| ------- | -------------- | --------- |
| `exma:TokenBalancesDataset` | `30` | `exma:TreasuryDashboard` |
| `exma:ContractResultsDataset` | `7` | `exma:ExecutionInsightsWorkspace` |

### Evidence bundle

See [`docs/competency/evidence/CQ-ANL-007.md`](evidence/CQ-ANL-007.md) for detailed sourcing and analytics requirements.
