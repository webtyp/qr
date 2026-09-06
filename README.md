# qr
<img src="docs/img/badges.svg">

QR code generation for WebTyp — encodes to a module matrix, renders to SVG; neutral Go, no build tags

## Status

Under construction. The API is specified in [docs/PLAN.md](docs/PLAN.md);
nothing is implemented yet.

```go
m, err := qr.Encode("https://192.168.1.5:8080/__webtyp/ca", qr.LevelM)
if err != nil {
    return err
}
svg := m.SVG(qr.SVGOptions{})
```

Encoding and rendering are separate: `Encode` returns the module matrix, and a
renderer draws it. Neutral Go with no build tags — the same code runs on the
server and in a WASM client.
