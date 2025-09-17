"""Phase 4 triple store data pilot automation using Oxigraph."""

from __future__ import annotations

import csv
import json
import pathlib
import shutil
import sys
import time
from typing import Iterable

from oxrdflib import OxigraphStore
from pyshacl import validate
from rdflib import Graph


REPO_ROOT = pathlib.Path(__file__).resolve().parent.parent
PILOT_DIR = REPO_ROOT / "build" / "pilots" / "phase4"
STORE_PATH = PILOT_DIR / "oxigraph-store"


def dataset_paths() -> list[pathlib.Path]:
    module_paths = [
        pathlib.Path("ontology") / "src" / name
        for name in (
            "core.ttl",
            "consensus.ttl",
            "token.ttl",
            "smart-contracts.ttl",
            "file-schedule.ttl",
            "mirror-analytics.ttl",
            "hiero.ttl",
        )
    ]
    alignment_paths = [
        pathlib.Path("ontology") / "src" / "alignment" / name
        for name in ("aiao.ttl", "claimont.ttl", "impactont.ttl", "infocomm.ttl")
    ]
    example_paths = sorted((REPO_ROOT / "ontology" / "examples").glob("*.ttl"))

    absolute_paths = [REPO_ROOT / path for path in module_paths + alignment_paths]
    absolute_paths.extend(example_paths)
    return [path for path in absolute_paths if path.exists()]


def shapes_paths() -> list[pathlib.Path]:
    return sorted((REPO_ROOT / "ontology" / "shapes").glob("*.ttl"))


def load_store(paths: Iterable[pathlib.Path]) -> Graph:
    if STORE_PATH.exists():
        shutil.rmtree(STORE_PATH)

    store = OxigraphStore()
    store.open(str(STORE_PATH), create=True)
    graph = Graph(store=store)

    for path in paths:
        graph.parse(path)

    graph.commit()
    return graph


def run_query(graph: Graph, query_path: pathlib.Path, output_csv: pathlib.Path) -> list[tuple[str, ...]]:
    query_text = query_path.read_text(encoding="utf-8")
    result = graph.query(query_text)

    headers = [str(var) for var in result.vars]
    rows = [tuple(str(value) if value is not None else "" for value in row) for row in result]

    output_csv.parent.mkdir(parents=True, exist_ok=True)
    with output_csv.open("w", newline="", encoding="utf-8") as handle:
        writer = csv.writer(handle)
        writer.writerow(headers)
        writer.writerows(rows)
    return rows


def run_shacl(graph: Graph, shape_paths: Iterable[pathlib.Path]) -> dict[str, object]:
    shapes_graph = Graph()
    for path in shape_paths:
        shapes_graph.parse(path)

    conforms, report_graph, report_text = validate(
        graph,
        shacl_graph=shapes_graph,
        inference="rdfs",
        serialize_report_graph=True,
    )

    report_dir = PILOT_DIR
    report_dir.mkdir(parents=True, exist_ok=True)

    (report_dir / "shacl-report.txt").write_text(report_text, encoding="utf-8")
    report_path = report_dir / "shacl-report.ttl"
    if isinstance(report_graph, (bytes, bytearray)):
        report_path.write_bytes(report_graph)
    else:
        report_graph.serialize(destination=report_path, format="turtle")

    return {"conforms": conforms, "report": "build/pilots/phase4/shacl-report.ttl"}


def serialize_graph(graph: Graph, destination: pathlib.Path) -> None:
    destination.parent.mkdir(parents=True, exist_ok=True)
    graph.serialize(destination=destination, format="turtle")


def main() -> int:
    start = time.time()

    datasets = dataset_paths()
    graph = load_store(datasets)

    query_path = REPO_ROOT / "tests" / "queries" / "cq-impact-001.rq"
    query_output = PILOT_DIR / "cq-impact-001.csv"
    rows = run_query(graph, query_path, query_output)

    shacl_info = run_shacl(graph, shapes_paths())

    serialize_graph(graph, PILOT_DIR / "anthropogenic-impact-dump.ttl")

    summary = {
        "datasets": [str(path.relative_to(REPO_ROOT)) for path in datasets],
        "query_results": str(query_output.relative_to(REPO_ROOT)),
        "shacl": shacl_info,
        "runtime_seconds": round(time.time() - start, 2),
        "triple_count": len(graph),
        "result_rows": rows,
        "store_path": str(STORE_PATH.relative_to(REPO_ROOT)),
    }

    PILOT_DIR.mkdir(parents=True, exist_ok=True)
    (PILOT_DIR / "pilot-summary.json").write_text(json.dumps(summary, indent=2), encoding="utf-8")

    print(json.dumps(summary, indent=2))
    graph.close()
    return 0


if __name__ == "__main__":
    sys.exit(main())
