# Ontology workspace

This directory hosts the machine-readable artefacts for the Bhash Hedera ontology.  The structure mirrors the layout described
in the project README and will expand as modelling sprints deliver additional modules.

## Layout

- `src/` – OWL/Turtle source files.  Each module is published as a distinct `.ttl` file that can be imported into Protégé or
  tooling pipelines.
- `shapes/` – SHACL constraint files that validate RDF data against the ontology (to be populated in later phases).
- `examples/` – Example graphs, competency question answers, and sample data extracts used for validation.

## Core module

`src/core.ttl` introduces the upper-level vocabulary shared by all service-specific modules.  The initial scaffold defines
classes for actors, networks, services, artefacts, processes, and events together with the relationships that connect them.  The
module is intentionally lightweight so that it can be extended incrementally as governance, token, and smart contract modules
contribute richer detail.

All IRIs follow the canonical namespace `https://bhash.dev/hedera/`.  Additional namespaces (e.g., `/governance/`, `/hts/`)
will be introduced as the project migrates from the core to service-focused modelling work.
