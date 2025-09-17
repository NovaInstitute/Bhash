# CQ-ANL-007 Evidence Bundle

**Question.** How is total supply of fungible tokens distributed across treasury and external accounts over time?

## Documentation anchors

* Mirror node REST balance endpoint provides token balance snapshots per account. [Mirror balances reference](https://docs.hedera.com/hedera/mirror-node/sdks-and-apis/rest-api#tag/Balances)
* Record and balance files define retention expectations for historical data. [Record and balance files](https://docs.hedera.com/hedera/mirror-node/architecture/record-and-balance-files)
* Treasury analytics require mapping datasets to downstream BI tooling. [Analytics overview](https://docs.hedera.com/hedera/mirror-node/architecture/overview)

## Data sources

| Asset | Purpose | Access method |
| ----- | ------- | ------------- |
| Mirror node `GET /api/v1/balances` | Retrieve per-account token balances for given timestamps. | REST call with `timestamp` query |
| Mirror node `GET /api/v1/tokens/{id}` | Obtain token metadata (supply, decimals) for contextualisation. | REST call per token |
| Mirror node `GET /api/v1/accounts/{id}` | Identify treasury vs external accounts. | REST call per account |
| Analytics workspace metadata | Describe dashboards consuming the datasets. | Internal documentation |

## Modelling updates

1. Introduced `hedera:MirrorDataset`, `hedera:DatasetRetentionPolicy`, and `hedera:AnalyticsWorkspace` classes with supporting properties to document dataset lineage.
2. Added SHACL shapes to enforce dataset type, retention, and coverage metadata.
3. Prepared query `cq-anl-007.rq` to verify dataset/workspace wiring ahead of live metric ingestion.

## Validation approach

* **SPARQL query:** `tests/queries/cq-anl-007.rq` ensures each dataset enumerates its retention period and consumer workspace.
* **SHACL:** `ontology/shapes/mirror-analytics.shacl.ttl` validates dataset metadata and retention definitions.

## Outstanding tasks

* Ingest real token balance extracts and join to treasury account classifications to answer the quantitative portion of the competency question.
* Capture retention SLAs from mirror operators to monitor for data drift.
