# go-codec — Library Adoption Audit

| Field | Value |
| --- | --- |
| Library | `github.com/larsartmann/go-codec` |
| Local checkout | `/home/lars/projects/go-codec` @ `98bb1bc` (2026-11-15) |
| Latest local tag | `v0.3.3` |
| Consumer | GoReleaser-Wizard @ `a493950` (2026-11-22) |
| Used? | **No** — `rg 'larsartmann/go-codec'` → 0 hits |
| Verdict | NOT-ADOPTED → **STAY-OUT** |
| Score | n/a — 0 applicable capabilities (7 rows, all N/A — decisive: no YAML codec exists) |
| Provenance | Fleet audit Group C, 2026-11-22; citations QC-sampled by orchestrator |

## What the library offers

A codec abstraction over **CBOR, JSON, and Raw** — canonical RFC 7049 CBOR, strict RFC 8949 compact CBOR, stdlib JSON (+ opt-in `encoding/json/v2` variant) — with streaming encoders/decoders, a generic `DecodeEnvelopeOrLegacy[T]`, `TranscodeToJSON`, format autodetection, an observability wrapper, and deterministic-codec contracts. Dependencies: fxamacker/cbor/v2, go-error-family (go-codec/go.mod:5-9).

| Capability | What it does | Citation | Applicable here? |
| --- | --- | --- | --- |
| CBOR (canonical/compact) | Signing-safe deterministic binary encoding | go-codec/cbor.go:19, cbor_compact.go | N/A — consumer emits no binary/event payloads |
| JSON codec + streaming | Stdlib JSON behind the `Codec` interface | go-codec/json.go:6, streaming.go:32, :42 | N/A — consumer's one JSON use is stdlib `json.Unmarshal` at internal/domain/config_core.go:393 |
| Raw passthrough | `RawCodec` | go-codec/raw.go:9 | N/A |
| **YAML** | **Does not exist** — `Encoding` consts are exactly `json\|cbor\|raw` | go-codec/codec.go:14-17, :38 | N/A — and this is decisive: the consumer's entire real encoding surface is YAML |
| Envelope/generics | `DecodeEnvelopeOrLegacy[T]`, `WrapEncode`/`UnwrapDecode` | go-codec/envelope.go:85, :23, :48 | N/A — no mixed-format streams, no signing |
| Transcode / autodetect | `TranscodeToJSON`, `AutoDetect` | go-codec/transcode.go:37, autodetect.go:63 | N/A |
| Observability | `ObserveCodec` metrics wrapper, `Size` | go-codec/observability.go:146, size.go:17 | N/A |

## What GoReleaser-Wizard does today (its real encoding work)

- `internal/domain/config_core.go:387-389` — `ToYAML()` → `yaml.Marshal` (go-faster/yaml, import at :8).
- `internal/domain/config_core.go:401-403` — `FromYAML()` → `yaml.Unmarshal` (go-faster/yaml).
- `internal/domain/config_core.go:392-395` — `FromJSON()` → stdlib `json.Unmarshal`.
- `cmd/goreleaser-wizard/validate_deprecations.go:34-36` — `yaml.Unmarshal` into a `yaml.Node` for deprecated-key detection (`go.yaml.in/yaml/v3`, import at :10).
- `internal/config/config.go:28, :37, :199` — Koanf v2 + `knadh/koanf/parsers/yaml`.
- Actual file generation is **Go text/template rendering** (cmd/goreleaser-wizard/templates/embedded.go) — not data serialization; no codec applies.

## Why STAY-OUT

- **Format mismatch is total**: the consumer encodes YAML (fleet-mandated go-faster/yaml per `how-to-golang` policy, present at go.mod); go-codec speaks CBOR/JSON/Raw only.
- The one stdlib JSON decode (config_core.go:393) is trivial and already optimal.
- Template rendering — the consumer's actual output path — is outside any codec's scope.
- Adds two transitive deps (fxamacker/cbor, go-error-family) for zero reachable capability.

## RECONSIDER-WHEN

Either (1) the wizard introduces content-addressed/deterministic payloads (signed or hash-pinned config snapshots, where `DeterministicCodec` + canonical CBOR matter), or (2) go-codec ships a YAML `Encoding` and fleet policy retires the go-faster/yaml mandate — then re-evaluate as a unification play.

## How to verify

```bash
cd /home/lars/projects/GoReleaser-Wizard
rg -n 'yaml\.(Marshal|Unmarshal)' --glob '!*_test.go'   # config_core.go:388, :402; validate_deprecations.go:36
rg -c 'larsartmann/go-codec' .                          # 0 hits
sed -n '14,17p' /home/lars/projects/go-codec/codec.go   # Encoding consts: json|cbor|raw — no YAML
git -C /home/lars/projects/go-codec tag --list 'v*' --sort=-v:refname | head -1   # v0.3.3
```
