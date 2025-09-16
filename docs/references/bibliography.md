# Hedera & Hiero Bibliography

This living bibliography anchors every ontology modelling decision in authoritative Hedera and Hiero sources.  Each entry should cite the most stable URL available and, when possible, note the relevant Hedera Improvement Proposal (HIP), API specification, or whitepaper section.  Contributors are encouraged to extend the tables below as new artefacts are published.

## Core protocol references

| Topic | Reference | Notes |
| ----- | --------- | ----- |
| Hedera main documentation | https://docs.hedera.com/hedera | Canonical API and service behaviour reference (HCS, HTS, HSCS, File Service, Scheduled Transactions, Staking, Mirror Nodes). |
| Hedera API specifications | https://github.com/hashgraph/hedera-protobufs | Protobuf definitions for services and mirror node APIs. |
| Hedera mirror node docs | https://docs.hedera.com/hedera/mirror-node | REST and gRPC endpoints for consensus data export. |
| Hedera Council governance | https://hedera.com/council | Council composition, governance charter, voting processes. |

## Hiero architecture

| Topic | Reference | Notes |
| ----- | --------- | ----- |
| Hiero overview | https://docs.hedera.com/hiero | Modular layer descriptions, shard strategy, validator onboarding. |
| Hiero technical blogs | https://hedera.com/blog | Product roadmap articles; cite the specific post when referenced. |
| HIP-840 (Hiero network) | https://hips.hedera.com/hip/hip-840 | Formalizes Hiero network evolution; includes terminology for shards and validator roles. |

## Service-specific materials

| Service | Reference | Notes |
| ------- | --------- | ----- |
| Consensus Service (HCS) | https://docs.hedera.com/hedera/core-concepts/hedera-consensus-service | Concepts, ordering guarantees, topic configuration. |
| Token Service (HTS) | https://docs.hedera.com/hedera/sdks-and-apis/token-service | Token lifecycle, custom fees, treasury accounts, compliance semantics. |
| Smart Contract Service (HSCS) | https://docs.hedera.com/hedera/core-concepts/smart-contracts | EVM compatibility, gas metering, precompiles. |
| File Service | https://docs.hedera.com/hedera/core-concepts/file-service | File storage, permissions, expiry, linkage to contract bytecode. |
| Scheduled Transactions | https://docs.hedera.com/hedera/core-concepts/scheduled-transactions | Scheduling flow, signature collection, execution semantics. |
| Staking | https://docs.hedera.com/hedera/core-concepts/staking | Staking process, reward cycle, node roles. |
| HIP-412 (NFT standard) | https://hips.hedera.com/hip/hip-412 | Defines metadata requirements for NFTs. |
| HIP-423 (Royalties) | https://hips.hedera.com/hip/hip-423 | Describes royalty fee mechanics for NFTs. |
| HIP-540 (Stablecoins) | https://hips.hedera.com/hip/hip-540 | Stablecoin token classification and governance. |

## External ontologies & patterns

| Domain | Reference | Applicability |
| ------ | --------- | ------------- |
| Provenance | https://www.w3.org/TR/prov-o/ | Model transaction and process lineage (`prov:Activity`, `prov:Agent`). |
| Data catalogues | https://www.w3.org/TR/vocab-dcat-3/ | Represent mirror node exports and dataset metadata. |
| Decentralised identifiers | https://www.w3.org/TR/did-core/ | Capture account identifiers and wallet bindings. |
| Financial instruments | https://spec.edmcouncil.org/fibo/ | Potential alignment for token classifications and governance roles. |

## Usage guidelines

1. **Cite precisely** – when referencing a statement in the ontology or decision log, point to the relevant section ID (e.g., `HTS > Custom fees > Royalty fees`).
2. **Track versions** – include documentation release tags or commit hashes when Hedera publishes versioned updates.
3. **Archive critical artefacts** – store PDFs or release notes in `docs/references/archive/` when URLs are prone to change.
4. **Request additions** – open an issue or pull request with the citation and justification when new documents appear.

> _Maintainers:_ Ensure this bibliography is reviewed at each release cycle to incorporate newly ratified HIPs and Hiero updates.
