# Bhash Ontology Workplan

This workplan outlines the research, modelling, validation, and delivery activities required to produce a comprehensive OWL ontology for the Hedera Network and the forthcoming Hiero architecture. The plan assumes an iterative delivery model with fortnightly review cycles and emphasises documentation traceability to the official Hedera manuals, HIPs, mirror node APIs, and Hiero technical notes.

## Guiding principles

* **Documentation-first** – every ontology element must be backed by a citation to Hedera or Hiero documentation, a HIP, or production API schema.
* **Modular delivery** – deliver service-specific modules that can evolve independently (e.g., consensus, tokens, smart contracts) but share a common upper ontology.
* **Competency-driven** – define competency questions before modelling and verify them via SPARQL queries or SHACL validation for each module.
* **Alignment-friendly** – reuse or map to established vocabularies (PROV-O, DCAT, W3C DID, schema.org, FIBO for financial concepts) when semantics align.
* **Automation-aware** – configure reasoning and validation scripts early so that ontology changes are continuously checked in CI.

## Phase 0 – Project setup (Week 0)


| Deliverable | Description | Location / Status |
| ----------- | ----------- | ----------------- |
| Bibliography | Seed `docs/references/` with links to Hedera docs, Hiero technical blog posts, HIP specifications, mirror node API references. | [`references/bibliography.md`](references/bibliography.md) initialised with core sources. |
| Governance log | Create `docs/decisions/` to capture modelling decisions, open questions, and approvals. | [`decisions/log.md`](decisions/log.md) created; ROBOT adoption decision recorded. |
| Toolchain | Decide on preferred tooling (Protégé, ROBOT, rdflib notebooks) and document setup scripts. | [`tooling/toolchain.md`](tooling/toolchain.md) documents Codex + ROBOT workflow and installation steps. |


### Phase 0 completion review (2024-05-02)

* ✅ Bibliography seeded with canonical Hedera and Hiero sources as captured in [`docs/references/bibliography.md`](references/bibliography.md).
* ✅ Governance and decision tracking process established via [`docs/decisions/log.md`](decisions/log.md), including tooling adoption records.
* ✅ Preferred research and automation toolchain documented in [`docs/tooling/toolchain.md`](tooling/toolchain.md) with actionable installation guidance.
* 🎯 Exit criteria met: documentation foundations are in place, enabling transition into the Phase 1 domain inventory tasks below.


## Phase 1 – Domain inventory (Weeks 1-2)

1. **Content audit:** Extract glossary candidates and process descriptions from Hedera docs (HCS, HTS, HSCS, File Service, Scheduled Transactions, Staking, Mirror Nodes, Token Service standards) and Hiero architectural posts.
2. **Context diagrams:** Draft context maps showing relationships between actors (council, node operators, wallets, dApps), services, and artefacts.
3. **Competency question backlog:** Capture high-priority queries grouped by stakeholder (governance, compliance, developer tooling, analytics).
4. **Existing ontology review:** Identify reusable patterns (e.g., ODRL for permissions, FIBO for financial instruments, IETF DID for identifiers).

_Current artefacts (2024-05-02):_ [`content-audit`](inventory/content-audit.md), [`context-maps`](inventory/context-maps.md), [`competency backlog`](competency/backlog.md), and [`ontology landscape review`](research/ontology-landscape.md).

### Phase 1 progress review (2025-09-16)

* ✅ The content audit now inventories major Hedera and Hiero services—including consensus, token, smart contract, file, schedule, staking, mirror, governance, and Hiero extensions—giving each an initial ontological description and documentation anchor to guide modelling decisions.
* ✅ Draft Mermaid-based context maps capture actor-to-service flows, Hiero modular layering, and end-to-end data lifecycles, clarifying how forthcoming OWL modules should interrelate.
* ✅ The competency question backlog tracks ten multi-stakeholder questions (governance, compliance, developer tooling, analytics, Hiero transition) that will drive prioritisation once validation assets are prepared, though all entries remain in "Draft" status pending evidence gathering.
* ✅ The external ontology landscape review summarises candidate imports (PROV-O, DCAT, DID Core, FIBO, ODRL, OpenTelemetry) and next actions for profiling them against Hedera requirements.
* ✅ Account and role coverage has been expanded in [`inventory/accounts-governance-glossary.md`](inventory/accounts-governance-glossary.md), the high-priority competency questions now have executable work packages in the [backlog](competency/backlog.md), and automation requirements are broken down into issue-ready tasks in [`tooling/automation-requirements.md`](tooling/automation-requirements.md).

## Phase 2 – Core ontology foundation (Weeks 3-4)

Focus on cross-cutting abstractions that every module depends on.

* **Upper-level classes:** Define `hedera:Network`, `hedera:Service`, `hedera:Actor`, `hedera:Artefact`, `hedera:Process`, `hedera:Event` with links to external vocabularies (e.g., `prov:Activity`).
* **Identity & governance:** Model accounts, keys, staking relationships, governance bodies (Hedera Council, community validators), and regulatory roles.
* **Transaction skeleton:** Capture the lifecycle of a Hedera transaction (initiation, consensus, execution, mirror export) and attach PROV relations for traceability.
* **Hiero overlay:** Introduce classes for Hiero modular layers (Consensus Layer, Execution Layer, Service Layer), shards, and virtualization constructs that extend the base `hedera:Network` hierarchy.
* **Artefact registry:** Define object properties for linking services to artefacts (e.g., `hedera:managesArtefact`, `hedera:consumesMessage`).
* **Outputs:** Core OWL module (`ontology/src/core.ttl`), glossary entries, competency question answers for foundational queries (e.g., "Which entities participate in consensus on mainnet?").

### Phase 2 progress review (2025-09-16)

* ✅ Core ontology updated with account/credential structures, validator onboarding processes, transaction lifecycle phases, and Hiero layer overlays, enabling downstream modules to extend shared semantics. [`ontology/src/core.ttl`](../ontology/src/core.ttl)
* ✅ Documentation refreshed with governance/transaction glossary entries and transaction lifecycle notes to anchor modelling decisions. [`docs/inventory/accounts-governance-glossary.md`](inventory/accounts-governance-glossary.md) · [`docs/inventory/transaction-lifecycle.md`](inventory/transaction-lifecycle.md)
* ✅ Example graph and SPARQL query published for CQ-CORE-001, demonstrating consensus participation tracing from governance mandates to validator accounts. [`ontology/examples/core-consensus.ttl`](../ontology/examples/core-consensus.ttl) · [`tests/queries/cq-core-001.rq`](../tests/queries/cq-core-001.rq) · [`docs/competency/core-foundational.md`](competency/core-foundational.md)
* 🎯 Exit criteria met: foundational ontology constructs, glossary coverage, and a validated competency question establish the platform for Phase 3 service-specific modules.

## Phase 3 – Service-specific modules (Weeks 5-10)

Deliver OWL modules iteratively; each sprint targets one service.

### Phase 3 progress review (2025-09-30)

* ✅ Consensus Service sprint delivered a dedicated ontology module with topic/message governance, validator onboarding roles, and mirror export artefacts plus SHACL coverage to keep key requirements executable. [`ontology/src/consensus.ttl`](../ontology/src/consensus.ttl) · [`ontology/shapes/consensus.shacl.ttl`](../ontology/shapes/consensus.shacl.ttl)
* ✅ Token Service sprint introduced HIP-540-compliant classes for stablecoins, token key assignments, and custom fees with supporting competency assets answering CQ-COMP-003. [`ontology/src/token.ttl`](../ontology/src/token.ttl) · [`docs/competency/token-compliance.md`](competency/token-compliance.md) · [`tests/queries/cq-comp-003.rq`](../tests/queries/cq-comp-003.rq)
* 🔄 Governance traceability for validator stewardship remains in progress to close CQ-GOV-001-B; consensus module scaffolding is in place but still needs end-to-end evidence bundles referencing council rosters. [`docs/competency/backlog.md`](competency/backlog.md)
* ⚠️ Remaining Phase 3 modules (Smart Contracts, File/Schedule, Mirror/Analytics, Hiero overlay) have design notes but no ontology commits yet; backlogged tasks depend on new data fixtures and automation reuse. [`docs/competency/backlog.md`](competency/backlog.md)

### 3A. Consensus Service (HCS)

* Model topics, messages, submission flow, ordering, finality guarantees, message retention policies, and associated throttling/fee schedule concepts.
* Represent links between topics and consuming applications (e.g., wallets, bridges, compliance pipelines).
* Encode mirror node export artefacts (gRPC, REST) and retention windows.

### 3B. Token Service (HTS)

* Capture fungible/non-fungible token types, supply controls, treasury accounts, KYC/freeze/kycKey semantics, custom fees, royalty rules (HIP-423), and token relationship state machines.
* Include HIP-driven standards (HIP-412 for NFTs, HIP-540 for stablecoins) as subclasses or profiles.
* Model compliance-relevant states (token freeze, KYC, wipe) and events (mint, burn, transfer).

### 3C. Smart Contract Service (HSCS)

* Represent contracts, bytecode, contract accounts, system contracts, EVM compatibility layers, scheduled transactions triggering contracts, and precompile interactions with HTS.
* Capture gas metering, state access patterns, and relationships to mirror node contract logs.

### 3D. File & Scheduled Transaction Services

* Model File Service entities (files, keys, expiration, content types) and their link to other services (e.g., storing contract bytecode).
* Describe Scheduled Transactions, including scheduled signatures, expiration, and execution outcomes.

### 3E. Mirror & Analytics Ecosystem

* Capture mirror node types (community, managed, enterprise), data ingestion pipelines, REST/gRPC APIs, and exported datasets (transactions, balances, topics).
* Model integration touchpoints for Hiero telemetry and node observability.

### 3F. Hiero-specific Enhancements

* Extend modules with Hiero architectural primitives: shards, cross-shard messaging, node roles (consensus, execution, availability), virtualization/rollup constructs, and onboarding workflows for permissionless validators.
* Align terminology with Hiero documentation and capture transitions from current Hedera mainnet operations.

Each module should deliver:

* OWL file (`ontology/src/<module>.ttl`)
* SHACL constraints (`ontology/shapes/<module>.ttl`)
* Competency question log with SPARQL queries and sample answers (`docs/competency/<module>.md`)
* Mapping table to source documentation (`docs/mappings/<module>.csv`)

## Phase 4 – Integration & alignment (Weeks 11-12)

* **Cross-module consistency:** Run reasoning to ensure shared classes/properties do not conflict; refactor to maintain OWL DL compliance.
* **External alignment:** Create equivalence/subClassOf relations to external ontologies, document alignments, and resolve semantic gaps.
* **Data pilots:** Load sample data (mirror node exports, HIP reference payloads) into a triple store (e.g., GraphDB, Blazegraph) and validate SHACL constraints.

## Phase 5 – Publication & automation (Weeks 13-14)

* **Documentation site:** Generate HTML documentation (via Widoco or MkDocs) summarising classes, properties, and examples.
* **Release management:** Tag the repository (`v0.1.0`), publish packaged ontology files, and document release notes.
* **CI/CD:** Configure GitHub Actions (or equivalent) to run unit reasoning, SHACL validation, SPARQL regression tests, and produce documentation artifacts on each push.

## Phase 6 – Adoption & feedback (Ongoing)

* **Community feedback loops:** Present ontology modules to Hedera and Hiero community channels, gather feedback, and track change requests.
* **Pilot integrations:** Support proof-of-concept projects (compliance analytics, DeFi dashboards, educational materials) that consume the ontology.
* **Maintenance cadence:** Establish quarterly review cycles to align with Hedera network upgrades, new HIPs, and Hiero milestones.

## Immediate next actions (updated 2025-09-30)

| # | Action | Status | Notes |
| - | ------ | ------ | ----- |
| 7 | Complete governance traceability modelling to close CQ-GOV-001-B (Phase 3 consensus focus). | 🔄 In progress | Extend `ontology/src/consensus.ttl` with council committee relationships and link evidence in [`docs/competency/evidence/CQ-GOV-001.md`](competency/evidence/CQ-GOV-001.md). |
| 8 | Assemble live stablecoin key custody snapshot for CQ-COMP-003-A. | ⏳ Not started | Create extraction script for mirror REST data and populate `data/token-compliance.json` before expanding ontology annotations. |
| 9 | Kick off Smart Contract Service module covering HTS precompile usage (CQ-DEV-005-B). | ⏳ Not started | Introduce `ontology/src/smart-contracts.ttl` with `hedera:PrecompileInvocation` classes and align with automation hooks in [`scripts/run_sparql.py`](../scripts/run_sparql.py). |
| 10 | Plan File/Scheduled Transactions and Mirror/Analytics sprints with dependency on shared datasets. | 📝 Planned | Draft modelling outline and dataset requirements to unblock subsequent Phase 3 submodules once consensus and token tasks land. |

## Risk considerations

| Risk | Mitigation |
| ---- | ---------- |
| Rapid evolution of Hiero specifications | Maintain `/draft/` namespace for unstable concepts and schedule frequent documentation reviews. |
| Overlapping semantics with external ontologies | Use explicit mapping documents and consultation with domain experts before asserting equivalence axioms. |
| Data availability for validation | Engage mirror node operators early to obtain sample datasets and confirm SHACL constraints mirror production payloads. |
| Toolchain complexity | Automate environment setup via scripts/containers and document manual alternatives. |

## Success metrics

* Coverage of priority competency questions with passing SPARQL queries.
* Positive validation of sample datasets against SHACL shapes for each module.
* Traceability matrix linking ontology elements to Hedera/Hiero documentation and HIP identifiers.
* Community adoption indicators (issue reports, integration references, citations in HIP discussions).
