from __future__ import annotations

import json

import pytest
import responses

from scripts.fluree_client import (
    FlureeClient,
    FlureeClientError,
    FlureeConfig,
    FlureeConfigurationError,
)


@pytest.fixture
def client() -> FlureeClient:
    config = FlureeConfig(api_token="token", tenant_handle="tenant")
    return FlureeClient(config)


@responses.activate
def test_generate_prompt_posts_expected_payload(client: FlureeClient) -> None:
    url = "https://data.flur.ee/api/tenant/generate-prompt"
    responses.post(url, json={"prompt": "SELECT *"}, status=200)

    result = client.generate_prompt("tenant", datasets=["tenant/sample"], prompt="List")

    assert result == {"prompt": "SELECT *"}
    call = responses.calls[0]
    assert json.loads(call.request.body) == {
        "datasets": ["tenant/sample"],
        "prompt": "List",
    }
    assert call.request.headers["Authorization"] == "Bearer token"
    assert call.request.headers["x-user-handle"] == "tenant"


@responses.activate
def test_transact_includes_optional_blocks(client: FlureeClient) -> None:
    url = "https://data.flur.ee/fluree/transact"
    responses.post(url, json={"status": "ok"}, status=200)

    payload = {
        "ledger": "tenant/sample",
        "context": {"@context": {"ex": "https://example.com/"}},
        "insert": [{"@id": "ex:thing", "@type": "ex:Class"}],
        "delete": [],
        "where": [{"@id": "ex:thing"}],
    }
    client.transact(
        ledger="tenant/sample",
        context=payload["context"],
        insert=payload["insert"],
        delete=payload["delete"],
        where=payload["where"],
    )

    sent_body = json.loads(responses.calls[0].request.body)
    assert sent_body == payload


@responses.activate
def test_error_response_raises(client: FlureeClient) -> None:
    url = "https://data.flur.ee/api/tenant/generate-sparql"
    responses.post(url, json={"message": "No dataset"}, status=404)

    with pytest.raises(FlureeClientError) as excinfo:
        client.generate_sparql("tenant", datasets=["missing"], question="Where")

    assert "HTTP 404" in str(excinfo.value)


def test_config_from_env(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv("FLUREE_API_TOKEN", "abc")
    monkeypatch.setenv("FLUREE_HANDLE", "tenant")
    monkeypatch.setenv("FLUREE_BASE_URL", "https://override")

    config = FlureeConfig.from_env()
    assert config.api_token == "abc"
    assert config.tenant_handle == "tenant"
    assert config.base_url == "https://override"


def test_config_from_env_missing(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.delenv("FLUREE_API_TOKEN", raising=False)
    monkeypatch.delenv("FLUREE_HANDLE", raising=False)

    with pytest.raises(FlureeConfigurationError):
        FlureeConfig.from_env()


@responses.activate
@pytest.mark.fluree_smoke
def test_smoke_generate_answer(client: FlureeClient) -> None:
    url = "https://data.flur.ee/api/tenant/generate-answer"
    responses.post(url, json={"answer": "Sample"}, status=200)

    result = client.generate_answer("tenant", datasets=["tenant/sample"], question="Q")

    assert result == {"answer": "Sample"}

