# Python Helper Replacement Status

All developer tooling that previously relied on the Python helpers under `scripts/`
now has first-class support inside the Go-based toolchain. The sections below
summarise the new commands and packages that replace the legacy utilities and
highlight any remaining Python entry points that still need attention.

## Fluree client tooling

The responsibilities of `scripts/fluree_client.py` are now covered by the
`internal/fluree` package, the `scripts/flureeclient` Go helper (exposing
structured logging for Cloud debugging), and the `bhashctl fluree` command
family.

### Highlights

* **Typed client** – `internal/fluree` exposes a strongly typed HTTP client with
  helpers for dataset management, ledger transactions, and prompt-generation
  endpoints. Configuration is sourced from the `FLUREE_*` environment variables
  that the Python helper required, but can also be overridden via CLI flags.
* **Unified CLI** – `bhashctl fluree` offers the following subcommands:
  * `create-dataset` – create datasets on a tenant with optional tags and
    visibility.
  * `transact` – apply JSON-LD transactions pulled from local fixture files.
  * `generate-sparql`, `generate-answer`, and `generate-prompt` – invoke the
    Fluree assistant endpoints used by regression harnesses.
* **Regression parity** – the Go client and the logging wrapper are covered by
  unit tests that verify request/response handling and error propagation. Live
  integration checks can be enabled with `go test ./internal/fluree -run-fluree`
  once the `FLUREE_*` secrets and `FLUREE_DATASET` identifiers are exported.

### Usage

```bash
export FLUREE_API_TOKEN=...  # credentials issued via Fluree Cloud
export FLUREE_HANDLE=...     # tenant handle

go run ./cmd/bhashctl fluree create-dataset \
  --owner my-tenant \
  --dataset-name pilot-ledger \
  --description "Phase 4 pilot dataset"

go run ./cmd/bhashctl fluree transact \
  --ledger my-tenant/pilot-ledger \
  --insert build/fixtures/pilot-insert.json

go run ./cmd/bhashctl fluree generate-sparql \
  --owner my-tenant --dataset pilot-ledger \
  --prompt "Which claims reference HIP-540?"
```

## Hedera topic bridge

The Hedera Consensus Service bridge previously implemented in
`scripts/hedera_topic_to_fluree.py` is tracked for migration into Go. The
existing Python helper remains available during the transition while the Go
implementation matures.

## Phase 4 pilot harness

The orchestration logic encapsulated by `scripts/run_phase4_pilot.py` will be
ported to a dedicated `bhashctl pilot` command after the Hedera bridge is
completed.

## Next steps

* Port the Hedera bridge to Go with retry/idempotency semantics and tests.
* Reimplement the Phase 4 pilot harness on top of shared Go utilities so the
  entire developer experience is consolidated within `bhashctl`.
* Replace the remaining Python utilities (`run_phase4_pilot.py`,
  `run_shacl.py`, and `run_sparql.py`) with Go equivalents and update the
  Makefile targets accordingly.
* Once the replacements land, remove the Python runtime dependency and delete
  the legacy scripts.
