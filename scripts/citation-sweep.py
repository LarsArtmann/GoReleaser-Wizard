#!/usr/bin/env python3
"""Citation sweep for docs/review/libs/ reports.

Extracts every `path/file.ext:LINE` (or :LINE-LINE) citation from the audit
reports and verifies file existence + line range against the consumer repo
AND the 20 library checkouts under ~/projects. Exit 0 = ship gate passed.

Resolution order per citation:
  1. repo-relative (consumer root)
  2. consumer-context prefixes (cmd layer, generators, internal packages)
  3. <lib-name>/-qualified: first path segment is a known library checkout
  4. full path under each library checkout
  5. FAIL
"""

from __future__ import annotations

import re
import sys
from pathlib import Path

REPO = Path(__file__).resolve().parent.parent
LIBS_ROOT = Path("/home/lars/projects")
LIB_NAMES = [
    "go-error-family",
    "go-idempotency",
    "go-retry",
    "go-etag",
    "httputil",
    "go-codec",
    "go-sse",
    "go-datastar",
    "go-atomic-write",
    "go-filewatcher",
    "go-daemon",
    "go-output",
    "go-health",
    "go-health-dashboard",
    "go-credentials-store",
    "cmdguard",
    "templ-components",
    "go-appkit",
    "go-cqrs-lite",
    "cqrs-htmx",
]
CONTEXT_PREFIXES = [
    "cmd/goreleaser-wizard",
    "cmd/goreleaser-wizard/generators",
    "internal/domain",
    "internal",
    "cmd",
]

CITE = re.compile(
    r"(?<![\w/.-])([A-Za-z0-9_.][A-Za-z0-9_./-]*\.(?:go|md|mod)):(\d+)(?:-(\d+))?"
)
EN_DASH = "\u2013"


def candidates(rel: str):
    p = Path(rel)
    yield REPO / p, "repo"
    for prefix in CONTEXT_PREFIXES:
        yield (
            REPO / prefix / p.name if p.parent == Path(".") else REPO / prefix / rel,
            f"ctx:{prefix}",
        )
    first = rel.split("/", 1)[0]
    if first in LIB_NAMES and "/" in rel:
        rest = rel.split("/", 1)[1]
        yield LIBS_ROOT / first / rest, f"lib:{first}"
    for lib in LIB_NAMES:
        yield LIBS_ROOT / lib / rel, f"scan:{lib}"


def resolve(rel: str):
    seen = set()
    for path, via in candidates(rel):
        s = str(path)
        if s in seen:
            continue
        seen.add(s)
        if path.is_file():
            return path, via
    return None, None


def main() -> int:
    report_dir = REPO / "docs" / "review" / "libs"
    reports = sorted(report_dir.glob("*.md"))
    if len(reports) != 20:
        print(f"FATAL: expected 20 reports, found {len(reports)}")
        return 2
    failures = 0
    checked = 0
    for report in reports:
        text = report.read_text(encoding="utf-8").replace(EN_DASH, "-")
        for match in CITE.finditer(text):
            rel, start, end = match.group(1), int(match.group(2)), match.group(3)
            end = int(end) if end else start
            checked += 1
            path, via = resolve(rel)
            if path is None:
                print(f"FAIL {report.name}: unresolved {rel}:{start}-{end}")
                failures += 1
                continue
            total = sum(1 for _ in path.open(encoding="utf-8", errors="replace"))
            if start > total or end > total:
                print(
                    f"FAIL {report.name}: {rel}:{start}-{end} out of range "
                    f"(resolved {path} via {via}, {total} lines)"
                )
                failures += 1
            else:
                print(
                    f"OK   {report.name}: {rel}:{start}-{end} -> {via}:{path.relative_to(REPO if REPO in path.parents else LIBS_ROOT)}"
                )
    print(f"\n{checked} citations checked, {failures} failures")
    return 1 if failures else 0


if __name__ == "__main__":
    sys.exit(main())
