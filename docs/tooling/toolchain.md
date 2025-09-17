# Toolchain Overview (Phase 0)

Phase 0 establishes the shared tooling required to research, author, and continuously validate the Bhash ontology.  The stack combines human-in-the-loop modelling in Protégé with automated build and test pipelines orchestrated by ROBOT, while Codex acts as the primary AI partner for research synthesis and boilerplate generation.

## Core tools

| Role | Tool | Purpose | Setup notes |
| ---- | ---- | ------- | ----------- |
| AI-assisted research & drafting | **Codex** | Summarise Hedera/Hiero documentation, draft competency questions, and bootstrap ontology skeletons or SHACL templates under human review. | Engage Codex through the repository issue/PR workflow. Capture prompts and generated artefacts in decision records when they influence modelling. |
| Ontology authoring | Protégé | Interactive OWL editing, class hierarchy management, annotation authoring. | Install Protégé 5.5+; configure the Bhash namespace prefix and enable reasoning with HermiT/ELK for spot checks. |
| Automation & validation | **Go CLI (`cmd/bhashctl`)** | Orchestrates ROBOT-driven SPARQL suites and TopBraid SHACL validation while keeping fixtures in sync. | Install Go 1.21+. Run `go run ./cmd/bhashctl install` to download ROBOT and the TopBraid SHACL CLI into `build/tools/`, then reuse `go run ./cmd/bhashctl {sparql,shacl}` for checks. |
| Underlying ontology automation | ROBOT | CLI executed by `bhashctl` for reasoning, report generation, and release assembly. | Advanced users can call `robot` directly; the Go CLI fetches the jar automatically and exposes configuration under `internal/tools`. |
| Legacy data scripting | Python 3 + RDFlib (optional) | Historical ingestion helpers and Fluree prototypes awaiting Go ports. | Create a virtual environment manually (`python3 -m venv build/venv && build/venv/bin/pip install -r requirements.txt`) if you need to run scripts under `scripts/`. |

## Why ROBOT for automation?

ROBOT provides first-class support for ontology engineering pipelines:

* **Reasoning & reports** – run ELK/HermiT, detect unsatisfiable classes, and generate validation reports (`robot reason`, `robot report`).
* **Template-driven authoring** – convert CSV design tables into OWL modules, aligning with our plan for competency-question-derived schemas.
* **Release assembly** – merge modules, extract subsets, and publish versioned artifacts with provenance metadata (`robot annotate`, `robot export`).
* **CI integration** – straightforward CLI invocation that fits GitHub Actions, enabling continuous checks on pull requests.

The Go-based `bhashctl` CLI standardises ontology automation while still leveraging ROBOT under the hood. Python remains available for ad-hoc data wrangling, but new validation and regression work should target the Go tooling for consistency and testability.

## Installation & environment

1. **Go toolchain** – install Go 1.21+ so the repository CLI can compile and run.
2. **Java runtime** – install OpenJDK 11 or 17; both ROBOT and the TopBraid SHACL CLI require a JVM.
3. **Bootstrap automation** – from the repository root, run:
   ```bash
   go run ./cmd/bhashctl install
   ```
   The command downloads ROBOT and the TopBraid SHACL distribution into `build/tools/` and records paths in `.bhashctl.yaml`. Subsequent `go run ./cmd/bhashctl sparql` and `go run ./cmd/bhashctl shacl` invocations reuse the cached binaries.
4. **Optional manual ROBOT install** – if you prefer direct CLI access, you can still install ROBOT via package managers or the helper script below. Ensure `~/bin` is on your `PATH` (e.g., `echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bashrc`).
   ```bash
   # macOS (Homebrew)
   brew tap obolibrary/tools
   brew install robot

   # Linux (manual)
   curl -L -o robot.jar https://github.com/ontodev/robot/releases/download/v1.9.5/robot.jar
   mkdir -p ~/opt/robot ~/bin
   mv robot.jar ~/opt/robot/robot.jar
   printf '#!/usr/bin/env bash
exec java -jar "$(dirname "$0")/../opt/robot/robot.jar" "$@"
' > ~/bin/robot
   chmod +x ~/bin/robot

   # Repository helper (Debian/Ubuntu)
   scripts/install_robot_cli.sh
   ```

### Fluree Cloud credentials

Phase A introduces a lightweight Python client for the Fluree Cloud HTTP API. Networked tests and CLI helpers expect the following environment variables:

| Variable | Purpose |
| --- | --- |
| `FLUREE_API_TOKEN` | Bearer token generated in the Fluree Cloud console. |
| `FLUREE_HANDLE` | Tenant handle that prefixes Cloud API routes. |
| `FLUREE_BASE_URL` (optional) | Override for non-production tenants; defaults to `https://data.flur.ee`. |

Integration tests marked with `@pytest.mark.fluree_live` remain skipped unless both the credentials are present **and** `pytest` is invoked with `--run-fluree`. Offline smoke tests use mocked HTTP responses and are safe to run by default.

## Codex collaboration guidelines

1. **Prompt hygiene** – share relevant documentation excerpts when asking Codex for modelling assistance to maintain traceability.
2. **Review AI output** – treat generated classes or annotations as drafts; verify them against the bibliography before committing.
3. **Capture provenance** – if Codex materially influences a decision, reference the interaction in the decision log and link to the supporting documentation.
4. **Security** – avoid sharing private keys, credentials, or unreleased Hedera material in prompts.

## Automation targets (Phase 3 bootstrap)

The Go CLI now fronts the regression workflow, with Make targets reserved for lower-level ROBOT utilities:

| Command | Purpose |
| ------- | ------- |
| `go run ./cmd/bhashctl install` | Downloads ROBOT and the TopBraid SHACL CLI into `build/tools/` and records paths in `.bhashctl.yaml`. |
| `go run ./cmd/bhashctl sparql` | Merges example datasets with ROBOT and executes every query under `tests/queries/`, comparing outputs to `tests/fixtures/results/`. |
| `go run ./cmd/bhashctl shacl` | Aggregates example data and shapes before invoking the TopBraid validator; writes reports to `build/reports/` on failure. |
| `make reason-core` | `robot reason --reasoner ELK --input ontology/src/core.ttl --output build/core-reasoned.ttl` – run ELK reasoning over the core module. |
| `make report-core` | `robot report --input ontology/src/core.ttl --output build/reports/core-report.tsv` – generate integrity reports to catch unsatisfiable classes or warnings. |
| `make template-example` | `robot template --template templates/example.csv --output build/templates/example.ttl` – demonstrate the CSV-to-OWL workflow seeded for AUT-003. |

The Go commands cache their downloads under `build/tools/`; delete `build/` if you need to force a fresh install.

## Next automation steps

* Extend CI workflows so GitHub Actions runs `go run ./cmd/bhashctl install`, `go run ./cmd/bhashctl sparql`, `go run ./cmd/bhashctl shacl`, and the ROBOT reasoning/report targets on each pull request.
* Add ROBOT profile verification (`robot verify-profile --profile DL`) once PROV-O/DCAT imports stabilise and additional modules land.
* Migrate remaining data helpers from Python to Go so ingestion and validation share the same toolchain.
