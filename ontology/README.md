# Ontology workspace

This directory hosts the machine-readable artefacts for the Bhash Hedera ontology.  The structure mirrors the layout described
in the project README and will expand as modelling sprints deliver additional modules.

## Layout

- `src/` – OWL/Turtle source files.  Each module is published as a distinct `.ttl` file that can be imported into Protégé or
  tooling pipelines.
- `shapes/` – SHACL constraint files that validate RDF data against the ontology (to be populated in later phases).
- `examples/` – Example graphs, competency question answers, and sample data extracts used for validation.

## Core module

`src/core.ttl` introduces the upper-level vocabulary shared by all service-specific modules.  Phase 2 expanded the scaffold
with account/credential abstractions, validator onboarding processes, transaction lifecycle phases, mirror data publications,
and Hiero layer overlays so downstream modules can specialise these foundations instead of redefining them.  The module now imports trimmed PROV-O and DCAT profiles under `src/imports/` so automation can reason over external terms without dragging the full ontologies into the build.

`src/consensus.ttl` kicks off the Hedera Consensus Service module.  It adds explicit classes for topics, messages, node endpoints, role assignments, and validator stewardship roles aligned with HIP-840 so the governance competency questions have a dedicated namespace for service-specific concepts.

All IRIs follow the canonical namespace `https://bhash.dev/hedera/`.  Additional namespaces (e.g., `/governance/`, `/hts/`)
will be introduced as the project migrates from the core to service-focused modelling work.

### Example data

- `examples/core-consensus.ttl` provides a small worked example showing a mainnet validator, its staking account, the governing
  council mandate, transaction phases, and mirror datasets.  The sample underpins the foundational competency query
  `tests/queries/cq-core-001.rq` documented in `docs/competency/core-foundational.md`.
- `tests/fixtures/datasets/cq-gov-001.ttl` extends the example space with committee membership, role assignments, and HIP references used by the governance competency query harness.

## Validation assets

- `shapes/consensus.shacl.ttl` defines the initial SHACL shapes verifying that validator onboarding processes cite mandates and that consensus nodes declare their operators.
- Regression queries live under `tests/queries/`, with automation scripts in `scripts/` and expected results captured in `tests/fixtures/results/`.
