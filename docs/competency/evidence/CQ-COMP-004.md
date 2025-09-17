# CQ-COMP-004 Evidence Bundle

**Question.** Which scheduled transactions remain pending beyond 24 hours due to missing signatures?

## Documentation anchors

* Scheduled transactions accumulate signatures until executed or expired. [How scheduled transactions work](https://docs.hedera.com/hedera/core-concepts/scheduled-transactions/how-scheduled-transactions-work)
* Mirror REST provides schedule status and collected signature counts. [Mirror REST schedules endpoint](https://docs.hedera.com/hedera/mirror-node/sdks-and-apis/rest-api#tag/Schedule)
* Network governance guidance outlines best practices for emergency pauses via scheduled transactions. [Emergency pause operations](https://docs.hedera.com/hedera/core-concepts/network-governance/emergency-management)

## Data sources

| Asset | Purpose | Access method |
| ----- | ------- | ------------- |
| Mirror node `GET /api/v1/schedules` | Retrieve schedule status, collected signatures, and expiration timestamps. | REST call filtered by `status=pending` |
| Mirror node `GET /api/v1/schedules/{id}/signatures` | Enumerate collected signatures, enabling missing signature calculations. | REST call per schedule |
| Governance playbooks | Identify which accounts must sign emergency or compliance-related schedules. | Council documentation |

## Modelling updates

1. `hedera:PendingSchedule` subclass captures schedules awaiting signatures with explicit counts and expiry metadata.
2. Object properties `hedera:requiresSignature` and `hedera:hasCollectedSignature` maintain traceability to required signers and collected approvals.
3. SHACL shape `hedera:PendingScheduleShape` enforces presence of signature counts and expiration timestamps.

## Validation approach

* **SPARQL query:** `tests/queries/cq-comp-004.rq` filters schedules with `pending-signatures` status and computes missing signatures using `hasRequiredSignatureCount` and `hasCollectedSignatureCount`.
* **SHACL:** `ontology/shapes/file-schedule.shacl.ttl` validates schedule artefacts.

## Outstanding tasks

* Extend fixtures with real mirror node payloads once API sampling automation is available.
* Align missing signature thresholds with operational SLAs (e.g., escalate if pending more than 24 hours).
