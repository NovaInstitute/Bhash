# CQ-GOV-001 Evidence Bundle

**Question.** Which Hedera Council members currently steward validator onboarding decisions and what HIPs underpin their mandates?

## Documentation anchors

* Hedera mainnet consensus nodes are presently permissioned and operated by the Hedera Council, establishing the governance scope for onboarding validators. [Mainnet consensus nodes reference](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/networks/mainnet/mainnet-nodes/README.md)
* The staking programme documents that the Council—through its CoinCom committee—votes on staking reward parameters, indicating formal decision points tied to validator participation. [Staking program reference](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/staking/staking.md#phase-iii-staking-rewards-program-launch)
* HIP-840 (Hedera Hiero Network) describes the transition to community validators and associated onboarding requirements. [HIP-840 landing page](https://hips.hedera.com/hip/hip-840) *(pending mirrored copy because the primary site is gated by Cloudflare challenges in this environment).* 

## Data sources

| Asset | Purpose | Access method |
| ----- | ------- | ------------- |
| Mirror node `GET /api/v1/network/nodes` | Provides canonical list of current consensus node operators, their account IDs, staking weights, and service endpoints for cross-checking Council stewardship. | [REST API network endpoint](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/sdks-and-apis/rest-api/network.md) |
| Council membership roster | Associates organisational members with council committees (e.g., CoinCom) for stewardship mapping. | [hedera.com/council](https://hedera.com/council) *(requires scraping/export into structured CSV for reproducibility).* | ✅ Cached under [`data/council-roster.csv`](../../../data/council-roster.csv) |
| HIP repository | Supplies normative mandates (e.g., validator onboarding criteria, committee responsibilities) for traceability. | [HIP GitHub archive](https://github.com/hiero-ledger/hiero-improvement-proposals) |

## Implemented modelling updates (2025-09-30)

1. Added `hedera:CouncilCommittee` and `hedera:ValidatorOnboardingCommittee` subclasses to the Consensus Service ontology so committee responsibilities are explicit in the TBox and annotated with `hedera:sourceDocument` citations.
2. Introduced the `hedera:hasSteward`/`hedera:stewardsProcess` object property pair to bind validator onboarding processes to their stewarding committees, ensuring provenance back to HIP mandates through `hedera:supportedBy` and `hedera:hasMandate`.
3. Refreshed fixtures and evidence—`tests/fixtures/datasets/cq-gov-001.ttl` now tags CoinCom and the Membership Committee as onboarding committees with `hedera:stewardsProcess` links, while `data/council-roster.csv` provides the roster-to-mandate join used by the SPARQL query.

## Validation approach

* **SPARQL query prototype:**
  `tests/queries/cq-gov-001.rq` now implements the query, executing against fixtures in `tests/fixtures/datasets/cq-gov-001.ttl` and diffing against the expected results snapshot in `tests/fixtures/results/cq-gov-001.csv` through the automation harness.
* **Mirror node reconciliation:** `ontology/shapes/consensus.shacl.ttl` now validates both validator onboarding processes (`hedera:ValidatorOnboardingShape`) and the committees that steward them (`hedera:ValidatorOnboardingCommitteeShape`) to guarantee mandates are attached to every steward listed in the query output.

## Next follow-ups

* Monitor HIP-840 revisions by tracking the public GitHub source ([`hiero-improvement-proposals/HIP/hip-840.md`](https://raw.githubusercontent.com/hiero-ledger/hiero-improvement-proposals/main/HIP/hip-840.md)) for changes that would alter onboarding states.
* Align validator onboarding states with forthcoming Hiero validator lifecycle documentation once public.
