# Toolchain Overview (Phase 0)

Phase 0 establishes the shared tooling required to research, author, and continuously validate the Bhash ontology.  The stack combines human-in-the-loop modelling in Protégé with automated build and test pipelines orchestrated by ROBOT, while Codex acts as the primary AI partner for research synthesis and boilerplate generation.

## Core tools

| Role | Tool | Purpose | Setup notes |
| ---- | ---- | ------- | ----------- |
| AI-assisted research & drafting | **Codex** | Summarise Hedera/Hiero documentation, draft competency questions, and bootstrap ontology skeletons or SHACL templates under human review. | Engage Codex through the repository issue/PR workflow. Capture prompts and generated artefacts in decision records when they influence modelling. |
| Ontology authoring | Protégé | Interactive OWL editing, class hierarchy management, annotation authoring. | Install Protégé 5.5+; configure the Bhash namespace prefix and enable reasoning with HermiT/ELK for spot checks. |
| Automated build/test | **ROBOT** | Command-line automation for OWL workflows: template expansion, reasoning, integrity reports, release packaging. | Requires Java 11+. Install via [ROBOT releases](https://github.com/ontodev/robot/releases) or Homebrew (`brew tap obolibrary/tools && brew install robot`). Documented below. |
| Data scripting | RDFlib (optional) | Python-based RDF manipulation for ad-hoc scripts, dataset ingestion, or SPARQL prototyping. | Managed in a `requirements.txt` when scripts are introduced. |
| Validation | pySHACL (optional) | Execute SHACL constraints against sample data exports. | Install via pip in the ontology tooling environment. |

## Why ROBOT for automation?

ROBOT provides first-class support for ontology engineering pipelines:

* **Reasoning & reports** – run ELK/HermiT, detect unsatisfiable classes, and generate validation reports (`robot reason`, `robot report`).
* **Template-driven authoring** – convert CSV design tables into OWL modules, aligning with our plan for competency-question-derived schemas.
* **Release assembly** – merge modules, extract subsets, and publish versioned artifacts with provenance metadata (`robot annotate`, `robot export`).
* **CI integration** – straightforward CLI invocation that fits GitHub Actions, enabling continuous checks on pull requests.

RDFlib remains useful for Python-based data wrangling, but ROBOT's ontology-specific automation better satisfies the "automatic building and testing" requirement.

## Installation & environment

1. **Java runtime** – install OpenJDK 11 or 17.
2. **ROBOT** – download the latest release or install via package manager.
   ```bash
   # macOS (Homebrew)
   brew tap obolibrary/tools
   brew install robot

   # Linux (manual)
   curl -L -o robot.jar https://github.com/ontodev/robot/releases/download/v1.9.5/robot.jar
   mkdir -p ~/opt/robot ~/bin
   mv robot.jar ~/opt/robot/robot.jar
   printf '#!/usr/bin/env bash\nexec java -jar "$(dirname "$0")/../opt/robot/robot.jar" "$@"\n' > ~/bin/robot
   chmod +x ~/bin/robot
   ```
   Ensure `~/bin` is on your `PATH` (e.g., `echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bashrc`).
3. **Environment variables** – set `ROBOT_JAR` if invoking the jar directly; configure `JAVA_OPTS` for memory-intensive operations.
4. **Project scripts** – upcoming `Makefile`/`justfile` targets will wrap ROBOT commands for reasoning, report generation, and release packaging.

## Codex collaboration guidelines

1. **Prompt hygiene** – share relevant documentation excerpts when asking Codex for modelling assistance to maintain traceability.
2. **Review AI output** – treat generated classes or annotations as drafts; verify them against the bibliography before committing.
3. **Capture provenance** – if Codex materially influences a decision, reference the interaction in the decision log and link to the supporting documentation.
4. **Security** – avoid sharing private keys, credentials, or unreleased Hedera material in prompts.

## Next automation steps

* Draft ROBOT command recipes for reasoning (`robot reason`), profile validation (`robot verify-profile`), and release generation (`robot export`).
* Create CI workflows that install ROBOT, run reasoning checks on ontology modules, and publish reports as build artefacts.
* Provide Python notebooks leveraging RDFlib for data-driven validation once sample mirror node datasets are introduced.
