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

1. Run the SPARQL query to enumerate mirror datasets, their covered services, retention windows, and downstream workspaces:
   ```bash
   arq --data ontology/src/core.ttl \
       --data ontology/src/mirror-analytics.ttl \
       --data ontology/examples/mirror-analytics.ttl \
       --query tests/queries/cq-anl-007.rq
   ```
2. Use SHACL validation to ensure datasets expose retention metadata required by compliance teams:
   ```bash
   python -m pyshacl --data-file ontology/examples/mirror-analytics.ttl \
                     --shacl-file ontology/shapes/mirror-analytics.shacl.ttl \
                     --inference rdfs
   ```
3. Future automation will layer token balance aggregates on top of the validated retention model.

### Sample result (derived from `ontology/examples/mirror-analytics.ttl`)

| dataset | retentionDays | workspace |
| ------- | -------------- | --------- |
| `exma:TokenBalancesDataset` | `30` | `exma:TreasuryDashboard` |
| `exma:ContractResultsDataset` | `7` | `exma:ExecutionInsightsWorkspace` |

### Evidence bundle

See [`docs/competency/evidence/CQ-ANL-007.md`](evidence/CQ-ANL-007.md) for detailed sourcing and analytics requirements.
