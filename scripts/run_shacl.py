#!/usr/bin/env python3
"""Run SHACL validation over example datasets."""

# DEPRECATED: use ``bhashctl shacl`` (Go implementation) instead.
from __future__ import annotations

import pathlib
import sys

from pyshacl import validate
from rdflib import Graph


REPO_ROOT = pathlib.Path(__file__).resolve().parent.parent
SHAPES_DIR = REPO_ROOT / "ontology" / "shapes"
DATASETS = [
    REPO_ROOT / "ontology" / "examples" / "core-consensus.ttl",
    REPO_ROOT / "ontology" / "examples" / "token-compliance.ttl",
    REPO_ROOT / "ontology" / "examples" / "smart-contracts.ttl",
    REPO_ROOT / "ontology" / "examples" / "file-schedule.ttl",
    REPO_ROOT / "ontology" / "examples" / "mirror-analytics.ttl",
    REPO_ROOT / "ontology" / "examples" / "hiero.ttl",
    REPO_ROOT / "ontology" / "examples" / "alignment-impact.ttl",
]
DATASETS.extend(sorted((REPO_ROOT / "tests" / "fixtures" / "datasets").glob("*.ttl")))


def load_graph(paths: list[pathlib.Path]) -> Graph:
    graph = Graph()
    for path in paths:
        if not path.exists():
            continue
        graph.parse(path)
    return graph


def main() -> int:
    data_graph = load_graph(DATASETS)
    shapes_graph = load_graph(sorted(SHAPES_DIR.glob("*.ttl")))

    conforms, report_graph, report_text = validate(
        data_graph,
        shacl_graph=shapes_graph,
        inference="rdfs",
        serialize_report_graph=False,
    )

    print(report_text)
    if not conforms:
        output_dir = REPO_ROOT / "build" / "reports"
        output_dir.mkdir(parents=True, exist_ok=True)
        report_path = output_dir / "shacl-report.ttl"
        report_graph.serialize(destination=report_path, format="turtle")
        print(f"SHACL validation failed. Report written to {report_path}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
