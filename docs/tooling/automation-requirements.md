# Automation Requirements Backlog

This backlog translates the toolchain expectations from `toolchain.md` into actionable automation tasks. It organises near-term build, reasoning, and validation work so Phase 2 ontology development can rely on repeatable CI checks.

## Goals

1. Provide reproducible commands for reasoning, integrity reports, and documentation exports using ROBOT.
2. Establish local developer ergonomics (Go CLI plus supporting Makefile/Justfile targets) and containerised workflows to minimise setup friction.
3. Define CI pipelines that exercise ontology reasoning, SHACL validation, and SPARQL regression suites on every push.

## Immediate actions

| ID | Task | Description | Dependencies |
| -- | ---- | ----------- | ------------ |
| AUT-001 | ROBOT reason target | Create a `Makefile` (or `justfile`) target that runs `robot reason --reasoner ELK --input ontology/src/core.ttl --output build/core-reasoned.ttl`. | ROBOT installed; Java 11+. | ✅ Implemented via `make reason-core` |
| AUT-002 | ROBOT report target | Add command `robot report --input ontology/src/core.ttl --output build/reports/core-report.tsv` and document how to interpret unsatisfiable classes. | AUT-001 | ✅ Implemented via `make report-core` |
| AUT-003 | Template pipeline | Scaffold a `templates/` directory with an example CSV + ROBOT template command to demonstrate module generation. | ROBOT; CSV seed. | ✅ `templates/example.csv` + `make template-example` |
| AUT-004 | SHACL harness | Introduce a Go-based command that executes SHACL shapes in `ontology/shapes/` against sample data via the TopBraid CLI. | TopBraid validator cached by `go run ./cmd/bhashctl install`. | ✅ `go run ./cmd/bhashctl shacl` |
| AUT-005 | SPARQL regression suite | Configure a `tests/queries/` directory and automation command that runs SPARQL queries against prepared datasets via ROBOT. | Dataset fixtures; ROBOT CLI. | ✅ `go run ./cmd/bhashctl sparql` |
| AUT-006 | Docs generation | Evaluate `robot export` or Widoco for generating HTML documentation from ontology modules; add placeholder command to build pipeline. | ROBOT export configured. |

## Actionable task breakdown

Translate each AUT item into a concrete GitHub issue (label: `automation`) so that implementation can be scheduled alongside ontology work. Suggested acceptance criteria are captured below.

### AUT-001 – ROBOT reason target

* **Issue stub:** `AUT-001: Add ROBOT reasoning Makefile target`
* **Implementation steps:**
  - Create `Makefile` (or `justfile`) entry `reason-core` executing `robot reason --reasoner ELK --input ontology/src/core.ttl --output build/core-reasoned.ttl`.
  - Ensure the `build/` directory is ignored by Git and created automatically when the command runs.
  - Document the command in `docs/tooling/toolchain.md` under the ROBOT workflow section.
* **Definition of done:** Running `make reason-core` (or `just reason-core`) succeeds locally and produces `build/core-reasoned.ttl`; command usage is documented.

### AUT-002 – ROBOT report target

* **Issue stub:** `AUT-002: Provide ROBOT report for core module`
* **Implementation steps:**
  - Extend the automation script/Makefile with `report-core` invoking `robot report --input ontology/src/core.ttl --output build/reports/core-report.tsv`.
  - Parse the report output in CI to fail on `ERROR` entries; store the TSV under `build/reports/`.
  - Reference the task from decision log D-0001 to maintain traceability for ROBOT adoption.
* **Definition of done:** Command executes locally, generates the report file, and documentation explains how to interpret unsatisfiable classes or property violations.

### AUT-003 – Template pipeline

* **Issue stub:** `AUT-003: Scaffold ROBOT template workflow`
* **Implementation steps:**
  - Create `templates/example.csv` and accompanying `templates/example-ontology.tsv` (or `.csv`) illustrating how module rows map to ontology terms.
  - Add Makefile target `template-example` running `robot template --template templates/example.csv --output build/templates/example.ttl`.
  - Capture usage instructions in `docs/tooling/toolchain.md` and link to the template from relevant ontology module READMEs.
* **Definition of done:** Sample template renders successfully; generated TTL is ignored by Git but previewed in documentation; onboarding guides reference the workflow.

### AUT-004 – SHACL harness

* **Issue stub:** `AUT-004: Introduce reusable SHACL validation command`
* **Implementation steps:**
  - Implement (or extend) a Go subcommand that aggregates example datasets and shapes before invoking the TopBraid SHACL CLI.
  - Ensure `go run ./cmd/bhashctl install` downloads the validator and caches it under `build/tools/`.
  - Document the CLI workflow in `docs/tooling/toolchain.md` and link to it from competency guides.
* **Definition of done:** Running `go run ./cmd/bhashctl shacl` fails on validation errors and can be executed in CI without manual setup.

### AUT-005 – SPARQL regression suite

* **Issue stub:** `AUT-005: Configure SPARQL regression harness`
* **Implementation steps:**
  - Create `tests/queries/` directory with placeholder query (e.g., `smoke-test.rq`) and expected results under `tests/fixtures/`.
  - Add a Go subcommand that merges datasets with ROBOT and executes every query, writing CSV outputs under `build/queries/`.
  - Surface mismatches against the stored fixtures and document the workflow in `docs/tooling/toolchain.md`.
* **Definition of done:** Running `go run ./cmd/bhashctl sparql` processes all `.rq` files, compares results to fixtures, and returns a non-zero exit code when differences appear.

### AUT-006 – Documentation generation

* **Issue stub:** `AUT-006: Add ontology documentation build pipeline`
* **Implementation steps:**
  - Evaluate `robot export` versus Widoco; select tooling and capture reasoning in `docs/decisions/log.md` if change is required.
  - Add Makefile target `docs` that emits HTML into `build/docs/` and references ontology modules (core + service-specific).
  - Update repository `.gitignore` for generated docs and note publication expectations (e.g., to GitHub Pages) in `docs/tooling/toolchain.md`.
* **Definition of done:** Local command produces browsable HTML documentation; CI plan includes publishing artefacts or marking as manual release step.

## CI/CD roadmap

1. **Bootstrap GitHub Actions workflow** executing AUT-001 through AUT-005 on every push/PR.
2. Cache ROBOT jar and TopBraid downloads to reduce runtime (<5 minutes target).
3. Publish reasoning and SHACL reports as workflow artefacts for reviewer visibility.
4. Gate merges on clean ROBOT reports and SHACL validation; allow documentation generation to run on release tags.

## Tracking & ownership

* Capture each `AUT-*` task as a GitHub issue (label: `automation`) with acceptance criteria mirroring the table above.
* Reference this document in the decision log entry D-0001 so automation scope remains aligned with the ROBOT adoption decision.
* Review automation coverage at the start of Phase 2 to confirm reasoning and validation scripts exist for each module.
