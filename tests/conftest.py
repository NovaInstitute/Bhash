"""Pytest fixtures for the Fluree integration."""

from __future__ import annotations

import os
from typing import Dict, Optional

import pytest

from scripts.fluree_client import FlureeClient, FlureeConfig, FlureeConfigurationError

FLUREE_ENV_VARS = ("FLUREE_API_TOKEN", "FLUREE_HANDLE")


def pytest_addoption(parser: pytest.Parser) -> None:
    parser.addoption(
        "--run-fluree",
        action="store_true",
        default=False,
        help="Run tests that communicate with Fluree Cloud",
    )


def pytest_configure(config: pytest.Config) -> None:
    config.addinivalue_line("markers", "fluree_live: marks tests requiring Fluree Cloud")
    config.addinivalue_line(
        "markers",
        "fluree_smoke: quick offline smoke tests for the Fluree client",
    )


@pytest.fixture(scope="session")
def fluree_credentials() -> Optional[Dict[str, str]]:
    """Return configured Fluree credentials or ``None`` when absent."""

    values = {name: os.getenv(name) for name in FLUREE_ENV_VARS}
    missing = [name for name, value in values.items() if not value]
    if missing:
        return None
    values["FLUREE_BASE_URL"] = os.getenv("FLUREE_BASE_URL", "https://data.flur.ee")
    return values  # type: ignore[return-value]


@pytest.fixture
def fluree_client_live(fluree_credentials: Optional[Dict[str, str]], request: pytest.FixtureRequest) -> FlureeClient:
    """Instantiate a :class:`FlureeClient` when credentials are available."""

    if not request.config.getoption("--run-fluree"):
        pytest.skip("Pass --run-fluree to execute Fluree Cloud tests")

    if not fluree_credentials:
        pytest.skip("Fluree Cloud credentials are not configured")

    try:
        config = FlureeConfig(
            api_token=fluree_credentials["FLUREE_API_TOKEN"],
            tenant_handle=fluree_credentials["FLUREE_HANDLE"],
            base_url=fluree_credentials["FLUREE_BASE_URL"],
        )
    except KeyError as exc:  # pragma: no cover - defensive guardrail
        raise FlureeConfigurationError(f"Missing credential key: {exc}") from exc

    return FlureeClient(config=config)

