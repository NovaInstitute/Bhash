# Phase 1 Competency Question Backlog

This backlog captures high-priority competency questions (CQs) that will guide ontology modelling, SHACL validation, and SPARQL query development.  Each entry references the stakeholder motivation, expected data sources, and relevant documentation anchors.

## How to use this backlog

1. Reference IDs in issues, ontology modules, and SHACL shapes so traceability is preserved.
2. Update the **Status** column as questions move from draft to validated (via SPARQL or SHACL evidence).
3. Link answers or query artefacts once produced to accelerate regression testing.

| ID | Stakeholder | Priority | Question | Source(s) | Notes | Status |
| -- | ----------- | -------- | -------- | --------- | ----- | ------ |
| CQ-GOV-001 | Governance | High | Which Hedera Council members currently steward validator onboarding decisions and what HIPs underpin their mandates? | [Hedera Council](https://hedera.com/council); [HIP-840](https://hips.hedera.com/hip/hip-840) | Requires modelling council committees, meeting records, and HIP provenance; evidence bundle captured in [`evidence/CQ-GOV-001.md`](evidence/CQ-GOV-001.md). | In review |
| CQ-GOV-002 | Governance | High | What quorum of council votes authorised the latest network fee schedule update? | [Council governance](https://hedera.com/council); [Fee schedule docs](https://docs.hedera.com/hedera/core-concepts/fees) | Dependent on availability of public vote records or release notes. | Draft |
| CQ-COMP-003 | Compliance | High | Which tokens classified as stablecoins (HIP-540) enforce KYC and freeze keys, and who controls those keys? | [HIP-540](https://hips.hedera.com/hip/hip-540); [Token Service](https://docs.hedera.com/hedera/sdks-and-apis/token-service/introduction) | Drives modelling of token compliance attributes and account roles; evidence bundle stored in [`evidence/CQ-COMP-003.md`](evidence/CQ-COMP-003.md). | In review |
| CQ-COMP-004 | Compliance | Medium | Which scheduled transactions remain pending beyond 24 hours due to missing signatures? | [Scheduled transactions](https://docs.hedera.com/hedera/core-concepts/scheduled-transactions); [Mirror REST schedules](https://docs.hedera.com/hedera/mirror-node/sdks-and-apis/rest-api#tag/Schedule) | Requires temporal reasoning and mirror data ingestion. | Draft |
| CQ-DEV-005 | Developer tooling | High | Which smart contracts invoke HTS system contract precompiles and what gas usage patterns do they exhibit? | [HSCS system contracts](https://docs.hedera.com/hedera/core-concepts/smart-contracts/system-contracts); [Mirror contract logs](https://docs.hedera.com/hedera/mirror-node/sdks-and-apis/rest-api#tag/Contracts) | Supports optimisation guidance and ontology alignment with execution traces; evidence bundle recorded in [`evidence/CQ-DEV-005.md`](evidence/CQ-DEV-005.md). | In review |
| CQ-DEV-006 | Developer tooling | Medium | Which topics enforce both admin and submit key signatures for configuration updates? | [HCS configuration](https://docs.hedera.com/hedera/core-concepts/hedera-consensus-service/manage-topics); [Mirror topics endpoint](https://docs.hedera.com/hedera/mirror-node/sdks-and-apis/rest-api#tag/Topics) | Connects governance metadata with operational artefacts. | Draft |
| CQ-ANL-007 | Analytics | High | How is total supply of fungible tokens distributed across treasury and external accounts over time? | [Token Service](https://docs.hedera.com/hedera/sdks-and-apis/token-service/token-relationships); [Mirror balances](https://docs.hedera.com/hedera/mirror-node/sdks-and-apis/rest-api#tag/Balances) | Requires longitudinal aggregation and PROV links to mirror exports. | Draft |
| CQ-ANL-008 | Analytics | Medium | Which consensus topics have retention periods shorter than 30 days and which dApps depend on them? | [HCS retention](https://docs.hedera.com/hedera/core-concepts/hedera-consensus-service/message-retention); [Mirror subscriptions](https://docs.hedera.com/hedera/mirror-node/sdks-and-apis/topics) | Encourages modelling of topic consumers and SLAs. | Draft |
| CQ-HIE-009 | Hiero transition | High | Which validators are participating in each Hiero shard and what onboarding state (pending, active, slashed) do they occupy? | [HIP-840](https://hips.hedera.com/hip/hip-840); [Hiero docs](https://docs.hedera.com/hiero) | Depends on forthcoming public datasets; placeholder for future integration. | Draft |
| CQ-HIE-010 | Hiero transition | Medium | How do shard-level consensus finality times compare to mainnet benchmarks for a given epoch? | [Hiero docs](https://docs.hedera.com/hiero); [Mirror performance metrics](https://docs.hedera.com/hedera/mirror-node/architecture/metrics) | Will require time-series data capture and modelling of performance indicators. | Draft |

## Next actions

* Prioritise CQ-GOV-001, CQ-COMP-003, and CQ-DEV-005 for initial ontology modelling, as they cover governance, compliance, and developer tooling focal areas.
* Identify authoritative data sources (council minutes, HIP releases, mirror node datasets) to ground each question before drafting SPARQL queries.
* Add lower-priority questions (e.g., education, ecosystem tooling) once foundational modules stabilise.
