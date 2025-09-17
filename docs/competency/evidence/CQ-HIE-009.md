# CQ-HIE-009 Evidence Bundle

**Question.** Which validators are participating in each Hiero shard and what onboarding state (pending, active, slashed) do they occupy?

## Documentation anchors

* Hiero validator lifecycle and onboarding states. [Hiero validators](https://docs.hedera.com/hiero/concepts/validators)
* Shard architecture and layer roles. [Hiero shards](https://docs.hedera.com/hiero/concepts/shards)
* Roadmap milestones for progressive decentralisation. [Hiero roadmap](https://docs.hedera.com/hiero/roadmap)

## Data sources

| Asset | Purpose | Access method |
| ----- | ------- | ------------- |
| Validator onboarding registry (forthcoming) | Publish validator applications, scores, and approvals. | Hiero governance portal |
| Hiero state proofs | Provide cryptographic confirmation of validator sets per shard. | Hiero state proof API (planned) |
| Mirror node shard metrics | Telemetry for shard participation and health. | Planned Hiero telemetry endpoints |

## Modelling updates

1. Added `hedera:HieroParticipation` bridge class linking validators, shards, onboarding states, and support layers.
2. Introduced `hedera:HieroOnboardingState` to classify participation statuses and support future enumerations (pending, active, slashed).
3. Configured SHACL rules to ensure participation records include validators, shards, and onboarding metadata.

## Validation approach

* **SPARQL query:** `tests/queries/cq-hie-009.rq` extracts validators with onboarding state labels and readiness scores.
* **SHACL:** `ontology/shapes/hiero.shacl.ttl` enforces structural requirements for shard participation records.

## Outstanding tasks

* Integrate live Hiero onboarding registry feeds once available.
* Add competency queries comparing shard participation counts against target decentralisation metrics.
