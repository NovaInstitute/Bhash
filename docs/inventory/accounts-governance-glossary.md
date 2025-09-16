# Accounts, Governance, and Node Role Glossary

This glossary expands the Phase 1 content audit with dedicated coverage of Hedera account constructs, governing bodies, and network roles. Each entry captures modelling cues and authoritative documentation links so the ontology can represent permissions, decision flows, and validator onboarding states with evidence.

## Accounts and credentials

| Term | Definition | Source | Modelling considerations |
| ---- | ---------- | ------ | ------------------------ |
| Hedera account | On-ledger identity created via transaction submission and funded by an existing account; stores keys, balance, and configuration state. | [Account creation](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/accounts/account-creation.md) | Model as subclass of `hedera:Actor` with required relationships to key material, payer account, and auto-renew configuration. |
| Account key list | Composite key (threshold or list) that controls authorisation on an account and downstream artefacts. | [Account properties – Keys](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/accounts/account-properties.md#keys) | Represent as reusable credential structure linked to accounts, files, topics, and tokens for governance pattern reuse. |
| Auto-create alias | Derived public-key or EVM address alias that allows creation of accounts upon first transfer. | [Auto account creation](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/accounts/auto-account-creation.md) | Capture alias identifiers and link to provenance of initial transfer event. |
| Treasury account | Account designated as token treasury with authority over supply and fee distributions. | [Tokenisation on Hedera](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/tokens/tokenization-on-hedera.md) | Relate to token artefacts via `hedera:managesArtefact`; include constraints for custom fee collectors. |
| Staking metadata account | Account that configures staking preferences, reward decline flag, and delegation target. | [Staking program](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/staking/staking.md#staking-nodes) | Model stake preferences as properties enabling inference of validator support and eligibility for rewards. |

## Governance bodies and decision actors

| Term | Definition | Source | Modelling considerations |
| ---- | ---------- | ------ | ------------------------ |
| Hedera Governing Council | Collective that operates permissioned consensus nodes, sets staking policy, and approves network-wide changes. | [Mainnet consensus nodes](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/networks/mainnet/mainnet-nodes/README.md) · [Staking program](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/staking/staking.md#phase-iii-staking-rewards-program-launch) | Represent as `hedera:GovernanceBody` with decision events tied to staking reward votes and validator onboarding approvals. |
| CoinCom (Council finance committee) | Council committee responsible for proposing and voting on staking reward rate updates. | [Staking program](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/staking/staking.md#phase-iii-staking-rewards-program-launch) | Model as specialised governance role connected to reward-rate modification processes and HIP mandates. |
| Compliance steward | Role accountable for ensuring token administrators manage keys (KYC, freeze, wipe) in line with HIP-540 immutability guarantees. | [HIP-540](https://raw.githubusercontent.com/hiero-ledger/hiero-improvement-proposals/main/HIP/hip-540.md) | Instantiate as `hedera:Role` associated with token admin accounts and referenced in competency questions on compliance controls. |

## Node and validator roles

| Term | Definition | Source | Modelling considerations |
| ---- | ---------- | ------ | ------------------------ |
| Consensus node operator | Council-operated node that submits transactions, participates in consensus, and enforces throttles. | [Mainnet consensus nodes](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/networks/mainnet/mainnet-nodes/README.md) | Link operators to node instances, staking weight, and governance approvals for onboarding or retirement. |
| Mirror node operator | Permissionless participant replicating consensus history for query APIs. | [Mainnet consensus nodes](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/networks/mainnet/mainnet-nodes/README.md) | Model as data publisher providing DCAT datasets and provenance for downstream analytics. |
| Validator onboarding candidate | Prospective Hiero validator subject to council review and HIP-540-style governance controls. | [Staking program](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/staking/staking.md#staking-nodes) · [HIP-540](https://raw.githubusercontent.com/hiero-ledger/hiero-improvement-proposals/main/HIP/hip-540.md) | Capture lifecycle states (applied, approved, active) and link to evidence bundles for CQ-GOV-001. |
| Reward account | System account 0.0.800 that disburses staking rewards once council-defined thresholds are met. | [Staking program](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/staking/staking.md#staking-reward-account) | Represent as managed artefact with policy constraints on balances and reward distribution triggers. |

These entries now provide the dedicated account and governance vocabulary required by Immediate Action 2. Subsequent ontology iterations can attach SKOS notes and provenance citations using `hedera:sourceDocument` annotations aligned with this glossary.
