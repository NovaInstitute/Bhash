# Fluree Cloud integration blueprint

## 1. Cloud HTTP API review

### 1.1 Authentication and base URL
- All Fluree Cloud calls use the `https://data.flur.ee` domain. Requests must include an `Authorization: Bearer <API token>` header and the tenant handle in the `x-user-handle` header.
- Dataset-scoped endpoints expect the dataset owner handle as a path segment: `/api/{handle}/...`.

### 1.2 Dataset lifecycle endpoints
- **Create dataset** – `POST /api/{handle}/create-dataset` with JSON payload (`datasetName`, `storageType`, `description`, `visibility`, optional `tags`). Returns confirmation payload on success.
- **List or manage datasets** – not yet documented in Cloud API; continue using console for manual inspection and track docs for updates.
- **Transact data** – `POST /fluree/transact` accepts JSON-LD context, `ledger` identifier (usually `{handle}/{dataset}`), and `insert` / `delete` / `where` objects for immutable commit semantics. Use this endpoint for seeding ontology-derived triples and test fixtures.

### 1.3 AI-assisted query endpoints
- **`POST /api/{handle}/generate-prompt`** – expands a natural language question into a SPARQL-ready prompt for the LLM agent. Request body includes a `datasets` array and `prompt` string.
- **`POST /api/{handle}/generate-sparql`** – asks the hosted model to emit a SPARQL query given datasets + natural language question; returns a `sparql` string.
- **`POST /api/{handle}/generate-answer`** – runs the generated SPARQL against the named datasets and returns an LLM-formatted answer payload.
- These endpoints power Fluree's chat-style UX and will back our automated examples once we publish an ontology-derived dataset.

## 2. Testing & example integration strategy

### 2.1 Local automation guardrails
- Add environment-aware pytest fixtures that read `FLUREE_API_TOKEN`, `FLUREE_HANDLE`, and optional dataset identifiers. Skip network tests when secrets are absent.
- Create a shared helper in `scripts/` (e.g., `scripts/fluree_client.py`) that wraps Cloud API calls, handles headers, and centralises logging for debugging.
- Use `responses` or `pytest-vcr` to mock Cloud API responses for unit-level tests so CI can run without external dependencies; integration suites can run selectively with `--run-fluree` marker.

### 2.2 Dataset provisioning for tests
- Automate dataset creation during integration tests only when a `FLUREE_TEST_DATASET` variable is absent; otherwise reuse configured dataset to avoid quota exhaustion.
- Seed dataset contents by transforming existing Turtle fixtures (e.g., `ontology/examples/*.ttl`) into JSON-LD inserts via ROBOT or rdflib conversion scripts, then POST batches to `/fluree/transact`.
- Record dataset IDs/commits in `data/fluree/fixtures.json` so tests can query deterministic content.

### 2.3 Example notebooks and documentation
- Extend `docs/examples/` with a Jupyter notebook or Markdown walk-through showing manual dataset creation, data loading, and use of `generate-sparql` / `generate-answer` for a canonical competency question.
- Provide CLI examples (via `Makefile` targets) to run SPARQL prompts against Fluree using stored dataset handles, enabling quick smoke tests for ontology updates.

## 3. Building an `f:DataModel` from our ontologies

### 3.1 Understand required shape
- Fluree chat experiences expect a top-level JSON-LD resource typed as `f:DataModel` describing datasets, context vocabularies, and entry points. The model must include:
  - `@context` mapping `f` to `https://flur.ee/ontology#` and linking to imported ontology prefixes.
  - `f:datasets` array referencing ledger identifiers and summary metadata.
  - `f:promptTemplates` (optional) to steer chat behaviour.

### 3.2 Conversion workflow
- Generate canonical JSON-LD contexts from our OWL/Turtle files using ROBOT (`robot convert --input ontology/src/<module>.ttl --format jsonld`), ensuring compacted prefixes align with Fluree requirements.
- Compose a `f:DataModel` document that imports these contexts and enumerates competency-question-aligned SPARQL templates.
- Validate the document locally by transacting it into a staging dataset and invoking `generate-sparql` to verify LLM awareness of ontology terms.

### 3.3 Governance and versioning
- Store the generated `f:DataModel` artifact under `data/fluree/<dataset>/data-model.jsonld` with provenance metadata (ontology commit hash, generation timestamp).
- Introduce a `make fluree-data-model` target that regenerates the artifact when ontology modules change, and add CI checks that diff the committed model to catch drift.

## 4. Fluree integration workplan

| Phase | Focus | Key tasks | Owners |
| --- | --- | --- | --- |
| Phase A (Week 1) | Enable tooling | - Implement `scripts/fluree_client.py` and pytest fixtures.<br>- Document required secrets in `docs/tooling/toolchain.md`.<br>- Create smoke test using a mock dataset. | Ontology tooling team |
| Phase B (Week 2) | Dataset bootstrapping | - Convert existing example graphs to JSON-LD inserts.<br>- Automate `/fluree/transact` seeding for staging dataset.<br>- Capture dataset metadata in `data/fluree/fixtures.json`. | Data engineering |
| Phase C (Week 3) | `f:DataModel` build | - Generate initial data model referencing core ontology modules.<br>- Validate LLM chat responses against competency Qs.<br>- Iterate on prompt templates to improve answers. | Ontology modelling |
| Phase D (Week 4) | Documentation & CI | - Publish tutorial notebook + CLI walkthrough.<br>- Add integration test marker to CI.<br>- Schedule weekly smoke tests with real API credentials. | DevRel & QA |
| Phase E (Ongoing) | Iteration & expansion | - Extend dataset coverage as new ontology modules land.<br>- Monitor API changes (Cloud HTTP release notes).<br>- Gather feedback from internal users of Fluree chat. | Cross-functional |

This blueprint should be revisited after the first integration sprint to capture lessons learned and refine automation guardrails.
