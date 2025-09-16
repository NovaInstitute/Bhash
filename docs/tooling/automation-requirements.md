# Automation Requirements Backlog

This backlog translates the toolchain expectations from `toolchain.md` into actionable automation tasks. It organises near-term build, reasoning, and validation work so Phase 2 ontology development can rely on repeatable CI checks.

## Goals

1. Provide reproducible commands for reasoning, integrity reports, and documentation exports using ROBOT.
2. Establish local developer ergonomics (Makefile/Justfile) and containerised workflows to minimise setup friction.
3. Define CI pipelines that exercise ontology reasoning, SHACL validation, and SPARQL regression suites on every push.

## Immediate actions

| ID | Task | Description | Dependencies |
| -- | ---- | ----------- | ------------ |
| AUT-001 | ROBOT reason target | Create a `Makefile` (or `justfile`) target that runs `robot reason --reasoner ELK --input ontology/src/core.ttl --output build/core-reasoned.ttl`. | ROBOT installed; Java 11+. |
| AUT-002 | ROBOT report target | Add command `robot report --input ontology/src/core.ttl --output build/reports/core-report.tsv` and document how to interpret unsatisfiable classes. | AUT-001 |
| AUT-003 | Template pipeline | Scaffold a `templates/` directory with an example CSV + ROBOT template command to demonstrate module generation. | ROBOT; CSV seed. |
| AUT-004 | SHACL harness | Introduce a Python script or ROBOT call that executes SHACL shapes in `ontology/shapes/` against sample data. Prefer `pyshacl` CLI within a reusable virtualenv. | pySHACL available. |
| AUT-005 | SPARQL regression suite | Configure a `tests/queries/` directory and shell wrapper that runs `sparql` queries against prepared datasets (e.g., Apache Jena `arq`). | Dataset fixture; Jena CLI. |
| AUT-006 | Docs generation | Evaluate `robot export` or Widoco for generating HTML documentation from ontology modules; add placeholder command to build pipeline. | ROBOT export configured. |

## CI/CD roadmap

1. **Bootstrap GitHub Actions workflow** executing AUT-001 through AUT-005 on every push/PR.
2. Cache ROBOT jar and Python dependencies to reduce runtime (<5 minutes target).
3. Publish reasoning and SHACL reports as workflow artefacts for reviewer visibility.
4. Gate merges on clean ROBOT reports and SHACL validation; allow documentation generation to run on release tags.

## Tracking & ownership

* Capture each `AUT-*` task as a GitHub issue (label: `automation`) with acceptance criteria mirroring the table above.
* Reference this document in the decision log entry D-0001 so automation scope remains aligned with the ROBOT adoption decision.
* Review automation coverage at the start of Phase 2 to confirm reasoning and validation scripts exist for each module.
