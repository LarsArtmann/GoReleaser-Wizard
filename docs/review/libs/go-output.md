# go-output — Library Adoption Audit

| Field | Value |
| --- | --- |
| Library | `github.com/larsartmann/go-output` (+ submodules table/tree/markdown/d2/plantuml/nom/graph/daghtml/delimited/markup/escape/tui/bdd/serialization/testhelpers, all `v0.38.7`) |
| Local checkout | `/home/lars/projects/go-output` @ `90fb2a6` (2026-11-17) |
| Latest local tag | root `v0.38.7` (examples/integration at `v0.23.2`) |
| Consumer | GoReleaser-Wizard @ `a493950` (2026-11-22) |
| Used? | **No** — `rg 'larsartmann/go-output'` → 0 hits |
| Verdict | NOT-ADOPTED → **STAY-OUT** |
| Score | n/a — 0 applicable capabilities (9 rows, all N/A) |
| Provenance | Fleet audit Group D, 2026-11-22; citations QC-sampled by orchestrator |

## What the library offers

Structured terminal/file output: a core table data model, a 16-format `Format` enum with unified `RenderTable` dispatch, fluent Table/Tree builders, and per-format submodules (markdown, D2, PlantUML, nix-output-monitor live progress, diagrams, CSV/TSV, JSON/YAML/TOML serialization, bubbletea views, BDD/test helpers). The root module is zero-dependency (imports only `golang.org/x/term`).

| Capability | What it does | Citation | Applicable here? |
| --- | --- | --- | --- |
| Core Table data model | Headers/rows/footer with sentinel validation errors | go-output/tabledata.go:18, :30 | N/A — nothing prints tabular data; output is status lines, emoji headers, bullet lists |
| 16-format `Format` enum | Type-safe runtime format selection for CLI flags | go-output/format.go:4, :48 | N/A — consumer has no `--format` flag |
| Unified render dispatch | `RenderTable(data, format, opts)` | go-output/render_tabledata.go:45 | N/A |
| Fluent builders | `TableBuilder` / `TreeBuilder` | go-output/table_builder.go:27; tree_builder.go:24 | N/A — no structured datasets constructed at runtime |
| Tree rendering | Parent/child dependency trees | go-output/tree_builder.go:24 | N/A — job engine logs linearly |
| markdown/d2/plantuml/nom submodules | Markdown tables, D2/PlantUML diagrams, live progress | go-output/markdown/markdown.go:49, :59; d2/d2_convert.go:1; plantuml/plantuml.go:1; nom/doc.go:1-12 | N/A — no exports, no diagrams, no long-running progress in a sub-second CLI |
| ColorMode control | auto/always/never | go-output/color.go:14-16 | N/A — handled via lipgloss + `--no-color` through koanf (internal/config/config.go:23, :171-173) |
| serialization/delimited/tui/bdd | JSON/YAML/TOML serializers, CSV/TSV streaming, bubbletea views | go-output/serialization/doc.go:12 | N/A — the only serialization is YAML *file generation*, unrelated to output formatting |
| Zero-dep root module | Heavy renderers isolated in submodules | go-output/README.md:19 | N/A (moot) — though adoption cost would be low *if* a need appeared |

## What GoReleaser-Wizard does today

Terminal UX is charm-only: lipgloss styles for titles/success/error/info (main.go:94-112), `printCommandHeader` (validate.go:11-15), and validate findings rendered as **emoji section headers + bullet lists** — `displayValidationResults` (cmd/goreleaser-wizard/validation_display.go:8-18), `displayErrors`/`displayWarnings`/`displayRecommendations` (:50-102). Deliberately list-shaped, not table-shaped. No `lipgloss/table`, no table/tree/diagram imports anywhere.

## Why STAY-OUT

- The only multi-item output (validate findings) is intentionally human/emoji-styled; a table renderer adds nothing users asked for.
- No exports, no machine-readable output flag, no diagrams — every go-output shape has zero call sites.
- The root module is zero-dep, so staying out costs nothing either; there is no incidental pressure.
- If `validate` ever grows columns of data, the root `Table` + one submodule would be the right tool *then* — not now.

## RECONSIDER-WHEN

A `validate --format json|csv|markdown` (scriptable/CI-consumable findings) or a tabular findings report is requested — go-output's `Format` enum + table model is purpose-built for exactly that flag.

## How to verify

```bash
cd /home/lars/projects/GoReleaser-Wizard
sed -n '8,18p' cmd/goreleaser-wizard/validation_display.go   # bullet-list rendering, no table
rg -c 'larsartmann/go-output' .                               # 0 hits
git -C /home/lars/projects/go-output tag --list 'v*' --sort=-v:refname | head -1   # v0.38.7
```
