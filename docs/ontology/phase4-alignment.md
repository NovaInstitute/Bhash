# Phase 4 External Alignment Blueprint

This blueprint operationalises the Phase 4 integration goal of binding Hedera-native modules to Anthropogenic Impact Accounting ontologies. The focus is on four priority vocabularies: AIAO, CliaMont, ImpactOnt, and InfoComm.

## Objectives

1. **Traceability:** Preserve provenance from Hedera transactions and governance artefacts into impact reporting statements.
2. **Semantic interoperability:** Ensure ESG, climate, and infrastructure concepts reuse established classes/properties from the external ontologies without creating logical conflicts.
3. **Automation-readiness:** Leverage ROBOT and existing SPARQL/SHACL pipelines so bridge modules participate in automated reasoning and validation.

## Deliverables

| Artefact | Description | Owner | Status |
| -------- | ----------- | ----- | ------ |
| `ontology/src/alignment/aiao.ttl` | MIREOT-trimmed import/bridge linking Hedera consensus & token events to `aiao:ImpactAssertion` patterns. | Ontology modellers | 🔄 In progress |
| `ontology/src/alignment/cliamont.ttl` | Alignment axioms connecting sustainability commitments and scheduled actions to CliaMont mitigation/adaptation classes. | Ontology modellers | 🔄 In progress |
| `ontology/src/alignment/impactont.ttl` | Bridges HTS compliance & treasury analytics metrics to ImpactOnt KPIs with provenance annotations. | Ontology modellers | 🔄 In progress |
| `ontology/src/alignment/infocomm.ttl` | Maps Hedera/Hiero infrastructure classes (nodes, shards, pipelines) to InfoComm communication assets. | Ontology modellers | 🔄 In progress |
| `docs/mappings/aiao-alignment.csv` | Traceability matrix referencing Hedera docs/HIPs that justify each mapping. | Documentation lead | ✅ Complete |
| `tests/queries/cq-esg-001.rq` | Competency query demonstrating cross-ontology ESG reporting. | Tooling lead | ✅ Complete |

## Integration steps

1. **Term scoping:**
   * Use `robot extract --method MIREOT` to pull only the terms required to satisfy targeted competency questions.
   * Maintain a shared `alignment-prefixes.ttl` file to control namespace declarations and avoid duplication. *(Initial prefixes committed under `ontology/src/alignment/prefixes.ttl`.)*
2. **Bridge modelling:**
   * Start with `owl:equivalentClass`/`owl:subClassOf` assertions where Hedera classes are narrower than the external ontologies.
   * Annotate every bridge axiom with `dcterms:source` linking to Hedera documentation (HIPs, service manuals) and to the relevant ontology documentation.
3. **Validation hooks:**
   * Extend existing ROBOT reasoning tasks to include the new alignment modules.
   * Author SHACL shapes that verify minimum data requirements for AIAO impact statements (e.g., actor, activity, impact metric).
4. **Pilot datasets:**
   * Transform the ESG-focused sample data into RDF using `scripts/run_sparql.py` pipeline extensions, targeting AIAO and ImpactOnt classes.
   * Capture feedback from ESG stakeholders and iterate on mappings before locking down equivalence axioms.

## Open questions

* What governance process is needed to version-sync with future AIAO/CliaMont releases?
* Should bridge modules live in the main namespace or a dedicated `/alignment/` namespace until stabilised?
* Can existing Hedera sustainability disclosures provide enough evidence to assert strong (equivalent) mappings, or should we start with weaker `rdfs:seeAlso` references?

## References

* [AIAO](https://datadudes.xyz/aiao)
* [CliaMont](https://datadudes.xyz/cliamont)
* [ImpactOnt](https://datadudes.xyz/impactont)
* [InfoComm](https://datadudes.xyz/infocomm)
