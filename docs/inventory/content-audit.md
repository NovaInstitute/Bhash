# Phase 1 Content Audit

This document captures candidate glossary entries and modelling hooks extracted from Hedera and Hiero documentation.  The audit focuses on concepts that recur across manuals, HIPs, and API references so that subsequent ontology classes can be justified with citations.

## Methodology

1. Prioritise canonical sources enumerated in the shared [bibliography](../references/bibliography.md).
2. Capture the Hedera/Hiero term, a concise ontological description, and the documentation anchor for traceability.
3. Note modelling considerations (e.g., required properties, lifecycle events, alignment opportunities) to inform competency questions and SHACL shapes.
4. Flag follow-up research when documentation gaps or ambiguities appear.

## Service and domain glossary candidates

### Consensus Service (HCS)

| Concept | Description | Source | Notes |
| ------- | ----------- | ------ | ----- |
| Topic | Ordered message channel configured with admin/submit keys, access controls, and message retention settings. | [Hedera Consensus Service](https://docs.hedera.com/hedera/core-concepts/hedera-consensus-service) | Model as core artefact with governance on topic creation, key rotation, and deletion. |
| SubmitMessage transaction | Transaction that appends messages to a topic with optional chunking, memo, and submit key signature requirements. | [HCS Submit message flow](https://docs.hedera.com/hedera/core-concepts/hedera-consensus-service/publish-and-subscribe) | Requires provenance links between submitting account, topic, and resulting message event. |
| Topic message | Immutable payload recorded with consensus timestamp and sequence number, optionally mirrored via REST/gRPC. | [Mirror node message reference](https://docs.hedera.com/hedera/mirror-node/sdks-and-apis/topics) | Capture linkage to mirror exports and retention policies. |
| Topic retention policy | Configuration for how long messages remain queryable (per topic or network default). | [HCS retention docs](https://docs.hedera.com/hedera/core-concepts/hedera-consensus-service/message-retention) | Influences temporal reasoning and data availability constraints. |

### Token Service (HTS)

| Concept | Description | Source | Notes |
| ------- | ----------- | ------ | ----- |
| Fungible token | Token type with divisible supply, treasury account, and configurable supply keys/custom fees. | [Token Service overview](https://docs.hedera.com/hedera/sdks-and-apis/token-service/introduction) | Capture supply control properties and relationships to treasury accounts. |
| Non-fungible token (NFT) | Unique token instances with serial numbers, metadata, and HIP-412 compliance requirements. | [HIP-412 NFT standard](https://hips.hedera.com/hip/hip-412) | Model mint/burn lifecycle events and metadata schema constraints. |
| Token relationship | Per-account state tracking KYC, freeze, balance, and allowance settings for a specific token. | [Token relationship docs](https://docs.hedera.com/hedera/sdks-and-apis/token-service/token-relationships) | Requires SHACL constraints for allowed state combinations. |
| Custom fee (royalty/fixed) | Fee assessed during token transfers (e.g., HIP-423 royalties, fixed fees). | [HIP-423 Royalties](https://hips.hedera.com/hip/hip-423) | Represent as policy artefact referencing collector accounts and calculation formulas. |

### Smart Contract Service (HSCS)

| Concept | Description | Source | Notes |
| ------- | ----------- | ------ | ----- |
| Smart contract | Deployed EVM bytecode executed on Hedera with contract ID, linked admin keys, and gas usage. | [Smart contract overview](https://docs.hedera.com/hedera/core-concepts/smart-contracts) | Align with `prov:Activity` for executions and capture state linkage to file storage. |
| Contract bytecode file | File Service artefact storing the compiled contract bytecode prior to deployment. | [Contract deployment flow](https://docs.hedera.com/hedera/sdks-and-apis/smart-contracts/deploy-contracts) | Bridge HSCS with File Service; track provenance to source artefacts. |
| System contract | Hedera-provided precompiled contract that exposes native services (e.g., HTS precompiles). | [System contract reference](https://docs.hedera.com/hedera/core-concepts/smart-contracts/system-contracts) | Represent as managed artefacts with versioning tied to network releases. |
| Gas & fee schedule | Cost model for contract execution, measured in gas and convertible to hbar fees. | [Smart contract fees](https://docs.hedera.com/hedera/core-concepts/smart-contracts/gas-and-fees) | Link to network fee schedules and staking rewards accounting. |

### File Service

| Concept | Description | Source | Notes |
| ------- | ----------- | ------ | ----- |
| File object | Storage container for bytecode, configuration, or metadata with size limits and expiration. | [File Service concepts](https://docs.hedera.com/hedera/core-concepts/file-service) | Associate with checksum, content type, and governance on mutability. |
| File keys | Key list controlling file updates, deletion, and access. | [File permissions](https://docs.hedera.com/hedera/core-concepts/file-service/file-permissions) | Model as key relationships reused by contracts and topics. |
| File expiration & auto-renew | Lifecycle management specifying expiration timestamps and auto-renew accounts. | [File lifecycle](https://docs.hedera.com/hedera/core-concepts/file-service/file-lifecycle) | Important for retention policies and dependency tracking. |

### Scheduled Transactions

| Concept | Description | Source | Notes |
| ------- | ----------- | ------ | ----- |
| Schedule entity | Wrapper storing a future transaction, its payer, and execution conditions. | [Scheduled transactions overview](https://docs.hedera.com/hedera/core-concepts/scheduled-transactions) | Map to underlying transaction types and required signatures. |
| Scheduled signature | Individual signature contributions tracked toward execution thresholds. | [Scheduled signing flow](https://docs.hedera.com/hedera/core-concepts/scheduled-transactions/sign-a-scheduled-transaction) | Requires modelling of signer roles and expiration windows. |
| Execution window | Period before the schedule expires or is executed once signatures are collected. | [Schedule execution semantics](https://docs.hedera.com/hedera/core-concepts/scheduled-transactions/execute-a-scheduled-transaction) | Enables temporal reasoning about pending obligations. |

### Staking & Network Participation

| Concept | Description | Source | Notes |
| ------- | ----------- | ------ | ----- |
| Staking node | Consensus node with staking parameters, node ID, and stake weight. | [Staking overview](https://docs.hedera.com/hedera/core-concepts/staking) | Links to validator operators and reward calculation. |
| Staking reward period | Epoch in which rewards accrue and are distributed to staked accounts. | [Reward period docs](https://docs.hedera.com/hedera/core-concepts/staking/rewards) | Drives temporal events and payout provenance. |
| Staking metadata account | Account storing metadata about staking preferences (decline rewards, auto-stake). | [Stake account settings](https://docs.hedera.com/hedera/core-concepts/staking/stake-your-account) | Connects to governance rules and scheduled payouts. |

### Mirror Nodes & Data Exports

| Concept | Description | Source | Notes |
| ------- | ----------- | ------ | ----- |
| Mirror node | Off-ledger service replicating consensus data via REST/gRPC APIs. | [Mirror node overview](https://docs.hedera.com/hedera/mirror-node) | Distinguish community vs managed nodes and deployment topology. |
| REST API dataset | Resources exposing transactions, balances, topics, tokens, contracts, etc. | [Mirror node REST reference](https://docs.hedera.com/hedera/mirror-node/sdks-and-apis/rest-api) | Map to DCAT dataset vocabulary for catalogue alignment. |
| Record file stream | Consensus record export (protobuf) consumed by analytics tooling. | [Record file format](https://docs.hedera.com/hedera/mirror-node/architecture/record-and-balance-files) | Supports provenance modelling for downstream data products. |

### Governance & Actors

| Concept | Description | Source | Notes |
| ------- | ----------- | ------ | ----- |
| Hedera Council | Governing body overseeing network policy, membership, and upgrades. | [Hedera Council](https://hedera.com/council) | Represent membership, voting rules, and decision provenance. |
| Network node operator | Organisation operating consensus or mirror nodes under council agreements. | [Node operator program](https://hedera.com/council/members) | Capture contractual roles and compliance obligations. |
| Account key hierarchy | Multi-key structures (threshold, key lists) controlling accounts and services. | [Account keys](https://docs.hedera.com/hedera/core-concepts/accounts/keys) | Foundational for modelling permissions across services. |

### Hiero architecture extensions

| Concept | Description | Source | Notes |
| ------- | ----------- | ------ | ----- |
| Shard | Partition of the network providing consensus scalability in Hiero. | [HIP-840 Hiero network](https://hips.hedera.com/hip/hip-840) | Model as `hedera:NetworkSegment` linked to validators and services. |
| Consensus layer | Hiero module responsible for ordering transactions and managing validator sets. | [Hiero docs](https://docs.hedera.com/hiero) | Align with base `hedera:Service` hierarchy. |
| Execution layer | Layer handling transaction execution, smart contract processing, and state storage. | [Hiero architecture overview](https://docs.hedera.com/hiero) | Capture dependencies on consensus layer outputs and state commitments. |
| Validator role (permissionless) | Expanded node role enabling community participation with staking and slashing rules. | [HIP-840 Hiero network](https://hips.hedera.com/hip/hip-840) | Model onboarding workflows and compatibility with existing governance. |

## Cross-cutting observations

* Many artefacts (topics, tokens, contracts, schedules) share governance patterns involving admin/submit keys; modelling reusable permission structures will reduce duplication.
* Lifecycle events (creation, update, expiry, deletion) should be represented consistently—likely via `prov:Activity` specialisations referencing network timestamps and mirror exports.
* Mirror node datasets provide verifiable evidence for competency question validation; align them with DCAT and PROV to trace raw to processed data products.

## Open research questions

1. How should topic/message retention policies map to SHACL constraints for data consumers with partial history requirements?
2. What minimal metadata is required to represent Hiero shard topology and validator capabilities without speculative assumptions?
3. Which governance artefacts (council votes, HIP ratification steps) have structured data sources suitable for automation?

## Next steps

* Prioritise glossary items with high ontology impact (tokens, accounts, staking) for detailed definition drafts in Protégé.
* Coordinate with competency question authors to ensure each item has at least one motivating query in the [backlog](../competency/backlog.md).
* Augment this audit with direct quotations or paraphrased definitions during ontology authoring to maintain traceability.
