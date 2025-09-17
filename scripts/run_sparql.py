#!/usr/bin/env python3
"""Execute repository SPARQL regression queries."""
from __future__ import annotations

import csv
import pathlib
import sys

from rdflib import Graph

REPO_ROOT = pathlib.Path(__file__).resolve().parent.parent
QUERIES_DIR = REPO_ROOT / "tests" / "queries"
RESULTS_DIR = REPO_ROOT / "tests" / "fixtures" / "results"
OUTPUT_DIR = REPO_ROOT / "build" / "queries"
DATASETS = [
    REPO_ROOT / "ontology" / "examples" / "core-consensus.ttl",
    REPO_ROOT / "ontology" / "examples" / "token-compliance.ttl",
    REPO_ROOT / "ontology" / "examples" / "smart-contracts.ttl",
    REPO_ROOT / "ontology" / "examples" / "file-schedule.ttl",
    REPO_ROOT / "ontology" / "examples" / "mirror-analytics.ttl",
    REPO_ROOT / "ontology" / "examples" / "hiero.ttl",
]
DATASETS.extend(sorted((REPO_ROOT / "tests" / "fixtures" / "datasets").glob("*.ttl")))


def load_graph() -> Graph:
    graph = Graph()
    for dataset in DATASETS:
        if dataset.exists():
            graph.parse(dataset)
    return graph


def write_results(path: pathlib.Path, headers: list[str], rows: list[tuple]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w", newline="", encoding="utf-8") as handle:
        writer = csv.writer(handle)
        writer.writerow(headers)
        for row in rows:
            writer.writerow(row)


def compare_results(expected: pathlib.Path, actual: pathlib.Path) -> bool:
    if not expected.exists():
        print(f"No expected results for {actual.name}; skipping comparison.")
        return True

    with expected.open("r", encoding="utf-8") as handle:
        expected_lines = [line.strip() for line in handle.readlines() if line.strip()]
    with actual.open("r", encoding="utf-8") as handle:
        actual_lines = [line.strip() for line in handle.readlines() if line.strip()]

    if expected_lines != actual_lines:
        print(f"Mismatch for {actual.name}:")
        print("Expected:")
        for line in expected_lines:
            print(f"  {line}")
        print("Actual:")
        for line in actual_lines:
            print(f"  {line}")
        return False
    return True


def run_query(graph: Graph, query_path: pathlib.Path) -> bool:
    query = query_path.read_text(encoding="utf-8")
    result = graph.query(query)
    headers = [str(var) for var in result.vars]
    rows = [tuple(str(value) if value is not None else "" for value in row) for row in result]

    output_path = OUTPUT_DIR / f"{query_path.stem}.csv"
    write_results(output_path, headers, rows)

    expected_path = RESULTS_DIR / f"{query_path.stem}.csv"
    return compare_results(expected_path, output_path)


def main() -> int:
    graph = load_graph()
    success = True
    for query_path in sorted(QUERIES_DIR.glob("*.rq")):
        print(f"Running {query_path.name}...")
        if not run_query(graph, query_path):
            success = False
    return 0 if success else 1


if __name__ == "__main__":
    sys.exit(main())
