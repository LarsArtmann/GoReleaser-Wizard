# templ-components — Library Adoption Audit

| Field | Value |
| --- | --- |
| Library | `github.com/larsartmann/templ-components` (+ htmx/datastar/icons/charts/errorpage/utils submodules) |
| Local checkout | `/home/lars/projects/templ-components` @ `bc369053` (2026-11-20) |
| Latest local tag | `v1.27.2` (all submodules) |
| Consumer | GoReleaser-Wizard @ `a493950` (2026-11-22) |
| Used? | **No** — `rg 'larsartmann/templ-components'` → 0 hits |
| Verdict | NOT-ADOPTED → **STAY-OUT** |
| Score | n/a — 0 applicable capabilities (8 rows, all N/A) |
| Provenance | Fleet audit Group D, 2026-11-22; citations QC-sampled by orchestrator |

## What the library offers

A server-rendered Go component library over templ + HTMX + Tailwind v4: 53 display components (cards, tables, tabs, modals, stat cards, kanban, native SVG charts), form/feedback/navigation/layout families with CSP-nonce plumbing, plus htmx/datastar integration, type-safe icons, and pretty HTTP error pages.

| Capability | What it does | Citation | Applicable here? |
| --- | --- | --- | --- |
| Server-rendered components | HTML over templ + HTMX + Tailwind v4 | templ-components/README.md:7 | N/A — no HTTP server, no templ files, no HTML output |
| `display` package (53 components) | Cards, tables, tabs, modals, charts, kanban | templ-components/README.md:151-153 | N/A — nothing to display in a browser |
| `forms` / `feedback` / `navigation` / `layout` | Form controls, toasts, nav, page scaffolds | templ-components/README.md quick start | N/A — the consumer's "forms" are huh TUI forms (tui_wizard.go) |
| `htmx` submodule | HTMX integration helpers | templ-components/htmx/doc.go:1 | N/A — no AJAX/hypermedia layer |
| `datastar` submodule | Datastar (SSE hypermedia) helpers | templ-components/datastar/ (attributes.go, cancellation.go) | N/A — no server-sent events, no server |
| `icons` submodule | Type-safe named icon constants + SVG | templ-components/icons/doc.go:1-2 | N/A — terminal UX uses emoji (✅❌⚠️💡 in validation_display.go) |
| `charts` submodule | Native SVG line/area/pie charts | templ-components/README.md:153 | N/A — no data-visualization surface |
| `errorpage` submodule | Pretty HTTP error pages | templ-components/errorpage/constructors.go:1 | N/A — CLI errors render via fang + lipgloss (main.go:138-146, :168-212) |

## What GoReleaser-Wizard does today

Nothing in this space — and explicitly refuses it: ROADMAP.md:67 lists **"Web UI" under Non-goals** ("the CLI/TUI is the product; a web layer duplicates it without a consumer"). go.mod has no templ, no HTMX, no Tailwind tooling; the only interactive surface is the huh wizard (tui_wizard.go:63-80).

## Why STAY-OUT

- The entire family is web-rendering shaped; the consumer is a terminal product by explicit ruling (ROADMAP.md:67), not by accident.
- Adopting any submodule would drag templ codegen + Tailwind + a component pipeline into a repo with zero HTML paths.
- Even the "closest" candidates (icons, errorpage) target SVG/HTML, while the consumer's equivalents (emoji, fang) are already terminal-native.
- No dependency overlap, no shared infrastructure, no migration path worth pricing.

## RECONSIDER-WHEN

Only if the Non-goal is revoked — a hosted/web version of the wizard is green-lit. Until ROADMAP.md:67 changes, this is a permanent STAY-OUT.

## How to verify

```bash
cd /home/lars/projects/GoReleaser-Wizard
rg -c 'larsartmann/templ-components|templ' go.mod   # 0 hits
sed -n '67p' ROADMAP.md                              # Web UI non-goal
git -C /home/lars/projects/templ-components tag --list 'v*' --sort=-v:refname | head -1   # v1.27.2
```
