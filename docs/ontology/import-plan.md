# PROV-O and DCAT Import Plan

This plan converts the ontology landscape review into concrete actions for profiling PROV-O and DCAT within the Bhash core module. The objective is to introduce external vocabulary with clear scope, version management, and testing hooks during Phase 2.

## Import strategy

1. **Scoped modules** – create dedicated import files (`ontology/src/imports/provo.ttl`, `ontology/src/imports/dcat.ttl`) that restrict the imported terms to those used by the core ontology.
2. **Version pinning** – reference stable IRIs (`https://www.w3.org/ns/prov#`, `https://www.w3.org/ns/dcat#`) alongside `owl:imports` annotations capturing the release date and checksum.
3. **Annotation alignment** – expose common annotations (`dcterms:title`, `dcterms:source`) via a shared prefixes file to maintain consistent metadata across native and imported terms.

## Task breakdown

| ID | Task | Description | Output |
| -- | ---- | ----------- | ------ |
| IMP-001 | Term inventory | Compile list of PROV-O classes/properties already referenced in `core.ttl` (e.g., `prov:Agent`, `prov:Activity`, `prov:Entity`, `prov:Role`). | Markdown checklist in `docs/ontology/import-plan.md`. |
| IMP-002 | ROBOT extract template | Use `robot extract --method MIREOT` to build trimmed PROV-O and DCAT files containing only required classes/properties. | `ontology/src/imports/provo.ttl`, `ontology/src/imports/dcat.ttl`. |
| IMP-003 | Namespace configuration | Update `ontology/src/core.ttl` to import the trimmed modules and declare `prov:` and `dcat:` prefixes consistently. | Revised `core.ttl`. |
| IMP-004 | Reasoning verification | Extend automation tasks (AUT-001/002) to include imported modules in reasoning/report runs. | Updated automation docs & scripts. |
| IMP-005 | Alignment documentation | Document mapping rationale and usage patterns in `docs/research/ontology-landscape.md` and create crosswalk tables for competency questions. | Updated research notes. |

## Profiling guidelines

* **PROV-O:** Limit import to classes (`prov:Entity`, `prov:Activity`, `prov:Agent`, `prov:Role`) and key properties (`prov:wasAssociatedWith`, `prov:wasGeneratedBy`, `prov:used`). Additional terms added only when required by competency evidence.
* **DCAT:** Focus on dataset discovery terms (`dcat:Dataset`, `dcat:Distribution`, `dcat:accessURL`, `dcat:downloadURL`, `dcat:spatial`) to model mirror node exports.
* **Validation:** Ensure imported terms remain within OWL DL by running `robot verify-profile --profile DL` as part of AUT-002.

## Open questions

1. Should we introduce a shared upper ontology module (e.g., `bhash-upper.ttl`) that houses all external imports before service-specific modules consume them?
2. How should we handle PROV-O qualified relations (e.g., `prov:qualifiedAssociation`) without complicating the competency query model?
3. Do we need lightweight application profiles (SHACL shapes) to enforce minimal metadata for DCAT datasets (title, description, access URL)?
