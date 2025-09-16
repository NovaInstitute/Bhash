# Bhash: Hedera Network Ontology

Bhash is an open knowledge engineering effort to describe the Hedera Network in a machine-interpretable way. The goal is to capture the actors, assets, services, processes, and governance rules that appear across the Hedera Consensus Service (HCS), Token Service (HTS), Smart Contract Service (HSCS), File Service, and related tooling in a unified OWL/RDF ontology. By modelling Hedera's public documentation and implementation guidance—including the recently announced Hiero network architecture—we aim to provide reusable semantics for analytics, compliance, education, and integration projects.

## Why an ontology?

* **Shared vocabulary** – establish stable identifiers and definitions for Hedera-specific notions such as accounts, topics, tokens, scheduled transactions, staking nodes, and fee schedules.
* **Interoperability** – align Hedera concepts with standard vocabularies (PROV-O for provenance, DCAT for data catalogues, W3C DID Core for identities, etc.) so that data from mirror nodes, dApps, and compliance tools can interoperate.
* **Reasoning** – enable automated validation (via OWL reasoning and SHACL constraints) for network states, policy rules, token compliance requirements, and smart contract metadata.
* **Documentation** – provide an authoritative reference that augments Hedera/Hiero manuals with explicit relationships that are otherwise scattered across prose and code.

## Scope overview

The ontology will ultimately cover four complementary viewpoints:

1. **Network Topology** – councils, permissioned and permissionless node operators, consensus and mirror nodes, network shards/regions introduced with Hiero, and environments (mainnet, testnet, previewnet, devnet).
2. **Identity & Accounts** – accounts, keys, roles, staking metadata, scheduled/smart contract-controlled accounts, treasury accounts, HIP-defined standards.
3. **Services & Artefacts** – topics, messages, file objects, fungible/non-fungible tokens, token classes (HIP-412, HIP-540), token relationships, smart contracts, contract bytecode/files, scheduled transactions, system contracts.
4. **Process & Governance** – transaction lifecycle, consensus flow, staking and reward cycles, network upgrade procedures, compliance/regulatory processes, fee schedules, and telemetry exposed through mirror node APIs.

Hiero introduces modular network layers, virtualization of consensus participation, and the path toward permissionless validator onboarding. These aspects will be modelled alongside legacy Hedera concepts so that the ontology can express both historical and forward-looking architecture.

## Planned repository structure

```text
.
├── README.md              # Project overview and quick links
├── LICENSE                # Apache 2.0 (inherited)
├── docs/
│   └── workplan.md        # Detailed execution plan
└── ontology/
    ├── src/               # OWL/RDF source files (TBD)
    ├── shapes/            # SHACL shapes for data validation (TBD)
    └── examples/          # Sample RDF graphs and competency questions (TBD)
```

> **Note:** Only documentation is present initially. Ontology artefacts will be added as the modelling work progresses following the workplan.

## Tooling & conventions

* **Ontology language** – OWL 2 DL serialised in Turtle (`.ttl`) for human readability; JSON-LD exports will be generated for integration scenarios.
* **Namespace strategy** – canonical namespace `https://bhash.dev/hedera/` with modular sub-namespaces per service (`.../hcs/`, `.../hts/`, `.../governance/`, etc.). Stable URIs will be reserved for Hedera/Hiero concepts; ephemeral or speculative notions will live in a `/draft/` namespace until ratified.
* **Competency questions** – each modelling sprint will begin by documenting competency questions (e.g., "Which accounts control the treasury of a given token?", "What consensus nodes participate in shard X under Hiero?").
* **Alignment artifacts** – crosswalk sheets mapping Hedera documentation terms to ontology classes/properties will live under `docs/mappings/` (to be created).
* **Validation** – SHACL shapes will express structural constraints derived from service APIs; automated reasoning (e.g., via the ELK or HermiT reasoner) will run in CI before releases.

## Working practices

1. **Document-first research** – extract canonical definitions from Hedera documentation, Hiero architectural notes, HIPs, and mirror node API references before introducing new classes.
2. **Iterative modelling** – deliver scoped ontology modules per Hedera service, validated with sample graphs and SPARQL competency queries.
3. **Community alignment** – involve Hedera developer relations, HIP authors, and compliance experts for terminology approval and governance modelling.
4. **Automation** – integrate Python notebooks (RDFlib, OWLReady2) or Java tooling (OWL API, ROBOT) to generate derived artefacts, run reasoning, and publish documentation.
5. **Versioning** – use semantic versioning for ontology releases (e.g., `v0.1.0` for initial scaffolding) with changelogs capturing class/property additions, modifications, and deprecations.

## Getting started

1. **Review the workplan** in [`docs/workplan.md`](docs/workplan.md) to understand the staged deliverables.
2. **Set up tooling:**
   * Install [Protégé](https://protege.stanford.edu/) for authoring OWL axioms.
   * Optionally configure a Python environment with [`rdflib`](https://rdflib.readthedocs.io/), [`pyshacl`](https://github.com/RDFLib/pySHACL), and [`owlready2`](https://owlready2.readthedocs.io/).
3. **Collect source materials:** maintain a shared bibliography of Hedera, Hiero, HIP, and mirror node references under `docs/references/` (to be established in Phase 0).
4. **Model & validate:** start with the seed module defined in the workplan (accounts & governance), ensuring each class/property is backed by documentation quotes and competency questions.

## Contributing

Contributions are welcome! Please open an issue describing the Hedera/Hiero concept you plan to model, include citations to the relevant documentation, and propose competency questions your addition should satisfy. Pull requests should include:

* Updated OWL/Turtle files
* Updated documentation (glossary entries, competency questions, or mapping tables)
* Validation evidence (reasoner logs, SHACL reports, or sample SPARQL answers)

## Reference materials

* [Hedera documentation](https://docs.hedera.com/hedera) – canonical service specifications, API references, and governance processes.
* [Hiero architecture overview](https://docs.hedera.com/hiero) – modular network evolution, shard design, and permissionless validator roadmap.
* [Hedera Improvement Proposals (HIPs)](https://hips.hedera.com/) – normative standards for tokens, accounts, mirror node APIs, and upcoming Hiero capabilities.

## License

This repository is distributed under the terms of the [Apache License 2.0](LICENSE).
