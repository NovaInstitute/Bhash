"""Utilities for interacting with the Fluree Cloud HTTP API.

This module centralises authentication headers, request/response logging,
and the small set of API endpoints we rely on for the initial integration
work.  The design keeps dependencies intentionally lightweight so it can be
used by both command-line scripts and pytest fixtures without pulling in
framework-specific helpers.
"""

from __future__ import annotations

import json
import logging
import os
from dataclasses import dataclass
from typing import Any, Dict, Iterable, Mapping, MutableMapping, Optional

import requests

logger = logging.getLogger(__name__)


class FlureeClientError(RuntimeError):
    """Raised when the Fluree API responds with an error."""


class FlureeConfigurationError(FlureeClientError):
    """Raised when the client cannot be configured from the environment."""


@dataclass(frozen=True)
class FlureeConfig:
    """Configuration bundle for :class:`FlureeClient`."""

    api_token: str
    tenant_handle: str
    base_url: str = "https://data.flur.ee"

    @classmethod
    def from_env(cls) -> "FlureeConfig":
        """Build a configuration object from environment variables.

        Expected variables:

        ``FLUREE_API_TOKEN``
            API token issued via the Fluree Cloud console.

        ``FLUREE_HANDLE``
            Tenant handle used as the path prefix for Cloud HTTP endpoints.

        ``FLUREE_BASE_URL`` (optional)
            Override the API base URL; defaults to ``https://data.flur.ee``.
        """

        api_token = os.getenv("FLUREE_API_TOKEN")
        tenant_handle = os.getenv("FLUREE_HANDLE")
        base_url = os.getenv("FLUREE_BASE_URL", cls.base_url)

        missing = [name for name, value in {
            "FLUREE_API_TOKEN": api_token,
            "FLUREE_HANDLE": tenant_handle,
        }.items() if not value]

        if missing:
            raise FlureeConfigurationError(
                "Missing environment variables: " + ", ".join(sorted(missing))
            )

        return cls(api_token=api_token or "", tenant_handle=tenant_handle or "", base_url=base_url)


class FlureeClient:
    """Minimal Fluree Cloud client with JSON helper methods."""

    def __init__(self, config: FlureeConfig, session: Optional[requests.Session] = None) -> None:
        self.config = config
        self._session = session or requests.Session()

    # -- HTTP helpers -------------------------------------------------

    def _headers(self) -> Dict[str, str]:
        return {
            "Authorization": f"Bearer {self.config.api_token}",
            "x-user-handle": self.config.tenant_handle,
            "Content-Type": "application/json",
        }

    def _post(self, path: str, payload: Mapping[str, Any]) -> Any:
        url = f"{self.config.base_url.rstrip('/')}/{path.lstrip('/')}"
        logger.debug("POST %s payload=%s", url, json.dumps(payload, ensure_ascii=False))
        try:
            response = self._session.post(url, headers=self._headers(), json=payload, timeout=30)
        except requests.RequestException as exc:  # pragma: no cover - network errors
            raise FlureeClientError(f"Error connecting to {url}: {exc}") from exc

        if response.status_code >= 400:
            message = self._format_error(response)
            logger.error("Fluree API error %s: %s", response.status_code, message)
            raise FlureeClientError(message)

        if not response.content:
            return None
        try:
            return response.json()
        except ValueError:  # pragma: no cover - unexpected payloads
            logger.debug("Response not JSON, returning text")
            return response.text

    @staticmethod
    def _format_error(response: requests.Response) -> str:
        try:
            data = response.json()
        except ValueError:
            return f"HTTP {response.status_code}: {response.text}"
        if isinstance(data, Mapping):
            message = data.get("message") or data.get("error") or data
            return f"HTTP {response.status_code}: {message}"
        return f"HTTP {response.status_code}: {data}"

    # -- API surface --------------------------------------------------

    def create_dataset(self, owner_handle: str, *, dataset_name: str, storage_type: str,
                       description: str, visibility: str = "private",
                       tags: Optional[Iterable[str]] = None) -> Any:
        payload: MutableMapping[str, Any] = {
            "datasetName": dataset_name,
            "storageType": storage_type,
            "description": description,
            "visibility": visibility,
        }
        if tags:
            payload["tags"] = list(tags)
        path = f"api/{owner_handle}/create-dataset"
        return self._post(path, payload)

    def transact(self, *, ledger: str, insert: Optional[Iterable[Mapping[str, Any]]] = None,
                 delete: Optional[Iterable[Mapping[str, Any]]] = None,
                 where: Optional[Iterable[Mapping[str, Any]]] = None,
                 context: Optional[Mapping[str, Any]] = None) -> Any:
        payload: Dict[str, Any] = {"ledger": ledger}
        if context is not None:
            payload["context"] = context
        if insert is not None:
            payload["insert"] = list(insert)
        if delete is not None:
            payload["delete"] = list(delete)
        if where is not None:
            payload["where"] = list(where)
        return self._post("fluree/transact", payload)

    def generate_prompt(self, owner_handle: str, *, datasets: Iterable[str], prompt: str) -> Any:
        payload = {"datasets": list(datasets), "prompt": prompt}
        return self._post(f"api/{owner_handle}/generate-prompt", payload)

    def generate_sparql(self, owner_handle: str, *, datasets: Iterable[str], question: str) -> Any:
        payload = {"datasets": list(datasets), "prompt": question}
        return self._post(f"api/{owner_handle}/generate-sparql", payload)

    def generate_answer(self, owner_handle: str, *, datasets: Iterable[str], question: str) -> Any:
        payload = {"datasets": list(datasets), "prompt": question}
        return self._post(f"api/{owner_handle}/generate-answer", payload)


__all__ = [
    "FlureeClient",
    "FlureeClientError",
    "FlureeConfigurationError",
    "FlureeConfig",
]

