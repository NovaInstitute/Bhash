# ESG Alignment Competency Answers

Phase 4 introduces the cross-ontology assets required to exercise CQ-ESG-001 and demonstrate AIAO-aligned reporting with Hedera data.

## CQ-ESG-001 – How do scheduled sustainability transactions translate into AIAO impact assertions with supporting KPIs and telemetry provenance?

* **Status:** Prototype ontology bridges, sample graph, and SPARQL query committed.
* **Scope:** Show how scheduled mitigation actions funded on Hedera surface as AIAO impact assertions with provenance to mirror datasets and sustainability KPIs.
* **Inputs:**
  * Ontology: `ontology/src/alignment/aiao.ttl`, `ontology/src/alignment/claimont.ttl`, `ontology/src/alignment/impactont.ttl`, `ontology/src/alignment/infocomm.ttl`
  * Sample graph: `ontology/examples/alignment-esg.ttl`
  * Query: `tests/queries/cq-esg-001.rq`

### Execution notes

1. Load the core and alignment ontologies alongside the ESG sample graph into an RDF store.
2. Execute the SPARQL query to retrieve each `aiao:ImpactAssertion`, its Hedera subject, associated KPI metrics, optional SDG goal, related scheduled mitigation action, and the communication node exposing evidence datasets.
3. Review the results to verify that Hedera concepts inherit semantics from AIAO, ClaimOnt, ImpactOnt, and InfoComm through the bridge axioms.

### Sample result (derived from `ontology/examples/alignment-esg.ttl`)

| assertion | subject | indicator | value | sdgGoal | schedule | node |
| --------- | ------- | --------- | ----- | ------- | -------- | ---- |
| `ex:ReforestationAssertion` | `ex:CommunityTreasury` | `ex:ReforestationMetric` | `"1250"^^xsd:decimal` | `impactont:SDG15LifeOnLand` | `ex:ReforestationSchedule` | `ex:MirrorNodeAlpha` |

### Next steps

* Expand sample data with additional SDG targets and multiple schedules to test aggregation.
* Create SHACL shapes verifying that each `hedera:StablecoinReserveAssertion` declares impact window, evidence dataset, and KPI linkages before asserting equivalence in production.
