# Modelling Decision Log

This log captures notable decisions, trade-offs, and open questions encountered while modelling the Hedera/Hiero ontology.  Each entry should link to the supporting documentation (e.g., bibliography citations, competency questions, pull requests) and record the status for traceability.

| ID | Date | Decision | Context & Rationale | Status | Owner |
| -- | ---- | -------- | ------------------- | ------ | ----- |
| D-0001 | 2024-05-01 | Adopt ROBOT for ontology automation | Compared ROBOT and RDFlib for build/test automation. ROBOT provides purpose-built OWL workflows (templates, reasoning, report generation) and integrates with CI pipelines for validation. RDFlib remains available for data scripting, but ROBOT will anchor automated builds. | Accepted | Ontology Engineering Team |

## How to propose updates

1. **Draft entry** – add a new row with a unique ID (`D-XXXX`), provisional status (`Proposed`), and reference links.
2. **Review** – discuss in issues or pull requests; capture alternatives considered.
3. **Finalise** – once approved, update the status to `Accepted` or `Rejected` and reference the decision artefact.

Maintain chronological order and ensure each decision points back to the relevant sources in `docs/references/bibliography.md`.
