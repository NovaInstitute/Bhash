# Python Helper Replacement Plan

This plan captures the remaining work needed to replace the Python utilities under
`scripts/` with Go implementations so that all developer tooling can be managed by
`bhashctl`.

## 1. `scripts/fluree_client.py`

### Current responsibilities
- Provides helper functions for interacting with the Fluree ledger HTTP API.
- Handles dataset uploads and query execution for regression harnesses.
- Contains ad-hoc configuration handling for credentials and server URLs.

### Go replacement strategy
1. **Design a Fluree client package** in `internal/fluree` that exposes strongly typed
   helpers for authentication, dataset management, and query execution.
2. **Integrate configuration with `bhashctl`** by extending the existing CLI config
   loading so the same flags/env vars can be reused when calling the new Go commands.
3. **Re-create the Python workflows** (dataset upload, query run, teardown) as subcommands
   under `bhashctl fluree` using the new client package.
4. **Add integration tests** that spin up the Fluree test environment (matching whatever
   the Python script expects) and run through the command to ensure parity.
5. **Document the migration** in `docs/` and update any developer instructions to stop
   referencing the Python helper once the Go feature is stable.

## 2. `scripts/hedera_topic_to_fluree.py`

### Current responsibilities
- Bridges Hedera Consensus Service topics into Fluree collections.
- Handles Hedera authentication, message retrieval, and transformation into Fluree schema.

### Go replacement strategy
1. **Inventory existing Hedera Go clients** (or add a lightweight wrapper) inside
   `internal/hedera` so we can reuse auth/session logic across commands.
2. **Model the bridging workflow** as a Go pipeline: read topic messages, transform to the
   Fluree schema, and persist via the Go Fluree client (from section 1).
3. **Expose the workflow via `bhashctl hedera bridge`** (or similar) with flags matching
   the Python script options.
4. **Add fault-tolerance features** such as retry/backoff and idempotent writes that are
   harder to manage in the Python script today.
5. **Provide end-to-end tests** using mocked Hedera/Fluree endpoints so CI can verify the
   bridge behaviour without real network access.

## 3. `scripts/run_phase4_pilot.py`

### Current responsibilities
- Orchestrates the Phase 4 pilot by coordinating Oxigraph queries, dataset loading, and
  result validation.

### Go replacement strategy
1. **Extract reusable Oxigraph interaction code** (possibly under `internal/oxigraph`) to
   manage dataset loading and SPARQL execution.
2. **Create a dedicated `bhashctl pilot` command** that mirrors the Python script options
   but leverages shared Go utilities for dataset discovery and validation.
3. **Reuse the existing regression fixtures** so the Go implementation can compare outputs
   exactly like the Python script.
4. **Add logging and metrics hooks** to align with the broader Go tooling observability
   story.
5. **Run regression tests** that compare Python vs. Go outputs during the migration phase
   until the Python helper can be fully retired.

## Cross-cutting tasks
- Update release and onboarding documentation once Go replacements are complete.
- Remove the deprecated Python scripts after a suitable grace period.
- Ensure CI pipelines invoke the new Go commands and eliminate Python-only dependencies.
