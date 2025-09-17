# CQ-DEV-005 Evidence Bundle

**Question.** Which smart contracts invoke HTS system contract precompiles and what gas usage patterns do they exhibit?

## Documentation anchors

* Hedera system smart contract documentation enumerates the HTS precompile (address `0x167`) and explains how on-ledger contracts invoke token operations via reserved selectors. [System smart contracts overview](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/smart-contracts/system-smart-contracts/README.md)
* Hedera gas reference outlines intrinsic gas, opcode gas, and HTS-specific surcharge calculations needed to interpret execution cost profiles. [Gas and fees guide](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/smart-contracts/gas-and-fees.md)
* Mirror node REST APIs provide contract execution results, logs, and opcode traces for analysing which contracts called the HTS precompile and the gas spent. [Smart contract REST reference](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/sdks-and-apis/rest-api/smart-contracts.md)

## Data sources

| Asset | Purpose | Access method |
| ----- | ------- | ------------- |
| Mirror node `GET /api/v1/contracts/results?function=HTS` | Filter contract executions that include HTS system contract selectors (via `result.logs.topic` or call trace). | [Contracts results endpoint](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/sdks-and-apis/rest-api/smart-contracts.md) |
| Mirror node `GET /api/v1/contracts/results/{transactionIdOrHash}/actions` | Retrieves call action trace including system contract addresses invoked, enabling attribution to `0x167`. | [Contracts actions endpoint](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/sdks-and-apis/rest-api/smart-contracts.md) |
| Mirror node `GET /api/v1/contracts/results/{transactionIdOrHash}/opcodes` | Supplies per-opcode gas usage to observe patterns for HTS precompile invocations. | [Contracts opcodes endpoint](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/sdks-and-apis/rest-api/smart-contracts.md) |
| Exchange rate service (`GET /api/v1/network/exchangerate`) | Converts gas to USD/HBAR for aggregated cost reporting. | [Network REST reference](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/sdks-and-apis/rest-api/network.md) |

## Proposed modelling updates

1. Introduced `hedera:ContractExecution`, `hedera:PrecompileInvocation`, and supporting object properties (`hedera:executesContract`, `hedera:includesInvocation`, `hedera:targetsSystemContract`) to capture invocation provenance.
2. Linked smart contracts to stored bytecode files and gas fee schedules so that execution metadata can be contextualised.
3. Added datatype properties (`hedera:hasGasLimit`, `hedera:hasGasUsed`, `hedera:hasFunctionSelector`, `hedera:hasGasBreakdown`) to express analytics-ready metrics from mirror node opcode traces.

## Validation approach

* **SPARQL query:** `tests/queries/cq-dev-005.rq` enumerates contracts, their precompile invocations, target precompiles, and gas usage.
* **SHACL:** `ontology/shapes/smart-contracts.shacl.ttl` verifies that executions and invocations contain contract references, gas metrics, and selectors.

## Outstanding tasks

* Confirm Mirror node filters for `actions[].call_type = "SYSTEM"` reliably capture HTS precompile calls; if not, add log signature heuristics.
* Define ETL mapping from mirror node JSON to RDF individuals capturing gas metrics and call hierarchy.
* Align gas units (weibar vs gas) with exchange-rate conversion so reports include both HBAR and USD equivalents.
