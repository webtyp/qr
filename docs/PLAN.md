---
PLAN: "feat: QR encoding to a module matrix with SVG rendering"
EXECUTOR: jules
REVIEWER: none
STATUS: review
SESSION: 16194427645207130646
PR: https://github.com/webtyp/qr/pull/1
---

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.

## Prerequisite — install the test runner

External agents run in isolated environments where `gotest` is not installed.
Run this **before anything else**; the acceptance criteria depend on it:

```bash
go install webtyp.com/devflow/cmd/gotest@latest
```

Then use `gotest` for the whole suite and `gotest -run TestName` for one test.
Never call `go test` directly.

# Plan — `webtyp.com/qr`

## Context (the executing agent has none — read this fully)

This repository is new and empty apart from a placeholder `qr.go`, which this
plan replaces entirely.

It generates QR codes for the WebTyp ecosystem. The first consumer is the
development server: `webtyp dev` serves its TLS certificate authority at
`/__webtyp/ca` and prints a QR of that URL so a phone on the LAN can install it
and test a PWA over HTTPS.

### The rules of this ecosystem — violating any of these fails review

- **This package is neutral: no build tags.** The same code runs on the server
  (the TUI printing a QR) and in the browser (a WASM client rendering one). That
  is the point of it being a library.
- **Never import the Go standard library** for what has a replacement: use
  `webtyp.com/fmt` instead of `fmt`, `errors`, `strconv` and `strings`. `math`
  is not on that list, but this package needs none — see below.
- **No `map` declarations.** Skill **wasm**: maps inflate the WASM binary. The
  GF(256) tables are two 256-byte arrays, and every capacity table is an array
  literal.
- **No generic holes.** No `func(...any)`, no `interface{}` in the public API.
- **Minimal surface.** Export only what a consumer calls.

### Why this is its own repository

Only the applications that draw a QR should pay for one in their binary. The
ecosystem splits one concern per module for exactly this reason.

### There is no floating-point arithmetic here

QR encoding is data segmentation, Reed–Solomon error correction over GF(256),
bit interleaving and mask selection: XOR, table lookups and bit shifts. Nothing
in this package needs `math`, and reaching for it is a sign of a wrong turn.

## Design gate

**1. Prior art.** `skip2/go-qrcode` and `boombuler/barcode` are the Go
references; `qrcode.js` and `node-qrcode` are the JavaScript ones. All of them
couple encoding to rendering: their entry point returns a PNG, a data URL or a
canvas. That coupling is why a project that wants SVG ends up either pulling an
image encoder it does not need, or writing the encoder again.

This package separates them: `Encode` returns the **matrix**, and rendering is a
method over it. The consumer picks the output, and a future PNG or canvas
renderer is additive rather than a second encoder.

**2. Novice-name test.** `qr.Encode("https://…", qr.LevelM)` reads as "encode
this as a QR at level M". `m.Dark(x, y)` asks whether a module is dark —
the term the specification uses, so a reader checking the code against the spec
finds the same word. Rejected: `New` (says nothing about what it makes),
`Generate` (ambiguous between the matrix and the image), `Get` (meaningless).

**3. Ledger.**

```
Concepts to learn                  +3   (Encode, Matrix, Level)
Repositories in the org            +1
Bytes in an app that draws no QR    0   (separate module)
Renderers coupled to the encoder   −1   (vs every existing library)
Ways to produce a QR               +1   (from 0 — new capability)
```

**4. Where it belongs.** QR encoding is one concern and owns this repository.
Rendering to SVG lives here because a matrix with no way to draw it is not
usable; a PNG renderer would be a second concern and is out of scope.

**5. What it deletes.** The placeholder `qr.go` and its `Qr` type.

## Stage 1 — the public surface

Replace `qr.go` entirely.

```go
// Level is the error-correction level. The zero value is LevelM, the level
// every general-purpose encoder defaults to: it survives about 15% damage,
// which is what a phone camera needs off a screen or a printed page.
type Level uint8

const (
    LevelM Level = iota // 15% — the default
    LevelL              //  7% — largest capacity, least tolerance
    LevelQ              // 25%
    LevelH              // 30% — smallest capacity, most tolerance
)

// Encode builds the QR matrix for data.
//
// It selects the smallest version (1–40) that holds data at the given level,
// and the mask that scores best under the specification's penalty rules.
func Encode(data string, level Level) (*Matrix, error)

// Matrix is a square grid of QR modules.
type Matrix struct { /* unexported */ }

// Size is the side length in modules: 21 for version 1, up to 177 for 40.
func (m *Matrix) Size() int

// Dark reports whether the module at (x, y) is dark. Out-of-range coordinates
// report false — the quiet zone is light.
func (m *Matrix) Dark(x, y int) bool

// Version is the QR version, 1 to 40.
func (m *Matrix) Version() int
```

Errors, each an exported value:

| Error | When |
|---|---|
| `ErrEmptyData` | `data` is empty — an empty QR is a programming error, not an outcome |
| `ErrTooLong` | the data exceeds version 40 at that level |
| `ErrUnknownLevel` | `level` is not one of the four constants |

`Level`'s zero value being `LevelM` is deliberate: a caller who writes
`qr.Encode(url, 0)` — or who forgets the argument in a struct literal — gets the
correct general-purpose level rather than the least tolerant one. Making
`LevelL` the zero value would hand the weakest error correction to whoever
thought about it least.

## Stage 2 — encoding

Byte mode only, over UTF-8. Numeric and alphanumeric modes compress digits and
uppercase text further; they are a size optimisation, not a capability, and are
out of scope. Do not add them.

The pipeline, in order — each step is the specification's, do not improvise:

1. **Segment**: the byte-mode header (mode indicator, character count) plus the
   data, padded with the alternating `0xEC 0x11` pad bytes.
2. **Version selection**: the smallest version whose data capacity at `level`
   holds the segment. Capacity is an array literal indexed by version and level.
3. **Reed–Solomon**: generate the error-correction codewords over GF(256) with
   the primitive polynomial `0x11D`. Build the log and antilog tables as two
   `[256]byte` arrays at init.
4. **Interleave**: split data and EC codewords into blocks per the version/level
   block table and interleave them as the spec prescribes.
5. **Place**: finder patterns, separators, timing patterns, alignment patterns,
   the dark module, format and version information, then the data in the
   zig-zag order.
6. **Mask**: apply each of the 8 mask patterns to the data modules, score each
   with the four penalty rules, keep the lowest. **Do not skip this and hard-code
   a mask** — some scanners fail on a badly masked symbol, and the failure looks
   like "the QR does not work sometimes", which is unfindable.

Put each step in its own file: `encode.go`, `version.go`, `galois.go`,
`interleave.go`, `place.go`, `mask.go`. Skill **core-principles** caps a file at
500 lines and requires one responsibility per file.

## Stage 3 — SVG rendering

```go
// SVGOptions controls the rendered SVG. The zero value renders a scannable
// code: 4-module quiet zone, black on white, sized to the module count.
type SVGOptions struct {
    ModuleSize int    // pixels per module; 0 = 4
    QuietZone  int    // modules of margin; 0 = 4, the specification's minimum
    Dark       string // CSS colour for dark modules; "" = "#000000"
    Light      string // CSS colour for the background; "" = "#ffffff"
}

// SVG renders the matrix as a standalone SVG document.
func (m *Matrix) SVG(opts SVGOptions) string
```

Emit the dark modules as **one `<path>`** whose `d` concatenates `M<x> <y>h1v1h-1z`
per module, with `shape-rendering="crispEdges"`, not as one `<rect>` per module.
A version-40 code has up to 31,329 modules; one path is a fraction of the size
and renders without seams between adjacent squares.

The quiet zone is not decoration: below 4 modules many scanners fail. `QuietZone`
may be raised, and a value below 4 is clamped to 4 rather than honoured — a
silently unscannable code is worse than an ignored option.

## Constraints

- **No hardcoded strings.** Every SVG attribute name, every default colour and
  every message is a named constant.
- **No `map`.** Arrays for every table.
- `webtyp.com/fmt` only; no `fmt`, `strings`, `strconv`, `errors`.
- No dependency other than `webtyp.com/fmt`.

## Tests

The pieces are deterministic and get exact assertions. The finished symbol is a
visual artifact and follows skill **testing**'s rule for binary output: write it
to disk for the developer to scan.

Exact, table-driven:

1. GF(256): `mul(0x53, 0xCA) == 0x01`, `mul(x, 0) == 0`, `mul(x, 1) == x`, and
   log/antilog round-trip for all 255 non-zero elements.
2. Reed–Solomon over the specification's worked example: the 16 data codewords
   of the version-1-M symbol for `01234567` produce the spec's 10 EC codewords.
   These values are published in ISO/IEC 18004 Annex I; assert them literally.
3. Version selection: a 1-byte payload → version 1; a payload one byte over
   version 1 level M capacity → version 2; over version 40 level H → `ErrTooLong`.
4. `Encode("", LevelM)` → `ErrEmptyData`.
5. `Encode(data, Level(99))` → `ErrUnknownLevel`.
6. `Level(0) == LevelM` — a guard, because reordering the constants would
   silently change every caller's default.
7. Mask penalty: each of the four rules scored against a hand-built matrix with
   a known violation.
8. `Dark(-1, 0)` and `Dark(size, 0)` → false, no panic.
9. `SVG(SVGOptions{QuietZone: 1})` → the rendered margin is 4 modules, not 1.
10. `SVG` output parses as XML and contains exactly one `<path>`.

Visual, per skill **testing**:

11. A test writes `testdata/out/*.svg` for a short URL at each of the four
    levels and for a payload large enough to reach a high version. It always
    passes; the developer scans the files with a phone. Document that in the
    test's comment so nobody mistakes it for a weak assertion — the exact
    assertions above are what guard correctness.

## Acceptance criteria

1. `grep -rnE '"(fmt|errors|strings|strconv|math)"' --include='*.go' .` → empty.
2. `grep -rn "map\[" --include='*.go' . | grep -v _test` → empty.
3. `grep -rn "type Qr struct" .` → empty (the placeholder is gone).
4. `gotest` passes.
5. Tests 2 and 6 pass — they are the ones that catch a silently wrong encoder.
6. No file exceeds 500 lines.

## Stages

| # | Stage | File(s) | Gate |
|---|---|---|---|
| 1 | public surface, `Level`, errors | `qr.go` | tests 4, 5, 6 |
| 2 | GF(256) + Reed–Solomon | `galois.go` | tests 1, 2 |
| 3 | version tables + segmentation | `version.go`, `encode.go` | test 3 |
| 4 | interleave + placement | `interleave.go`, `place.go` | builds a full matrix |
| 5 | mask selection | `mask.go` | test 7 |
| 6 | SVG | `svg.go` | tests 8, 9, 10, 11 |

Sequential. Stage 5 is the one most likely to be skipped or faked; criterion 5
and test 7 exist to make that impossible to hide.

## Out of scope

- Numeric, alphanumeric and Kanji modes.
- Micro QR.
- PNG or canvas rendering — additive later, over the same `Matrix`.
- Decoding.

Adding any of them is scope creep: report the need, do not build it.
