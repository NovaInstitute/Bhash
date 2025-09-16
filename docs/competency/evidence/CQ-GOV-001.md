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
| Council membership roster | Associates organisational members with council committees (e.g., CoinCom) for stewardship mapping. | [hedera.com/council](https://hedera.com/council) *(requires scraping/export into structured CSV for reproducibility).* |
| HIP repository | Supplies normative mandates (e.g., validator onboarding criteria, committee responsibilities) for traceability. | [HIP GitHub archive](https://github.com/hiero-ledger/hiero-improvement-proposals) |

## Proposed modelling updates

1. Extend `hedera:GovernanceBody` with subclasses for council committees (CoinCom, Membership Committee) and annotate them with HIP identifiers (`dcterms:identifier`).
2. Introduce a `hedera:ValidatorOnboardingProcess` class linking governance bodies to validator candidates via `prov:wasInformedBy` relations referencing HIP decisions.
3. Attach `hedera:hasMandate` annotations (datatype property) capturing HIP references for each committee and encode them as SKOS notes for downstream documentation.

## Validation approach

* **SPARQL query prototype:**
  ```sparql
  PREFIX hedera: <https://bhash.dev/hedera/core/>
  PREFIX dcterms: <http://purl.org/dc/terms/>
  SELECT ?councilMember ?committee ?hip
  WHERE {
    ?committee a hedera:GovernanceBody ;
               hedera:hasRole hedera:ValidatorOnboardingSteward ;
               dcterms:identifier ?hip .
    ?membership hedera:hasParticipant ?councilMember ;
                hedera:hasParticipant ?committee ;
                a hedera:ValidatorOnboardingProcess .
  }
  ```
  This returns council organisations linked to onboarding processes and surfaces the governing HIP identifiers.
* **Mirror node reconciliation:** create a SHACL shape verifying that every active consensus node (`hedera:ConsensusNode`) is `isOperatedBy` an actor that participates in a governance process endorsed by a HIP reference.

## Outstanding tasks

* Capture a cached copy of the HIP-840 specification to avoid Cloudflare access barriers and extract explicit onboarding state definitions.
* Build or ingest a structured council roster (CSV/JSON) with committee assignments so the ontology can populate the `hedera:ValidatorOnboardingSteward` role.
* Align validator onboarding states with forthcoming Hiero validator lifecycle documentation once public.
