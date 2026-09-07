# qr
<img src="docs/img/badges.svg">

QR code generation for WebTyp — encodes to a module matrix, renders to SVG; neutral Go, no build tags.

## Getting Started

Add the module to your Go project:

```bash
go get webtyp.com/qr
```

### Basic Usage

Encode a URL or string data into a QR module matrix and render it as an SVG string:

```go
package main

import (
	"webtyp.com/fmt"
	"webtyp.com/qr"
)

func main() {
	// 1. Encode data into a Matrix using a specific error correction Level (LevelM is default zero value)
	matrix, err := qr.Encode("https://192.168.1.5:8080/__webtyp/ca", qr.LevelM)
	if err != nil {
		fmt.Printf("Encoding failed: %v\n", err)
		return
	}

	// 2. Render the Matrix to a standalone SVG string
	svg := matrix.SVG(qr.SVGOptions{
		ModuleSize: 4,         // 4px per module (default: 4)
		QuietZone:  4,         // 4-module quiet zone margin (minimum: 4)
		Dark:       "#000000", // Dark module CSS color (default: #000000)
		Light:      "#ffffff", // Background CSS color (default: #ffffff)
	})

	fmt.Println(svg)
}
```

### Public Surface

- `qr.Encode(data string, level qr.Level) (*qr.Matrix, error)`
- `qr.Level`: `qr.LevelM` (default, 15%), `qr.LevelL` (7%), `qr.LevelQ` (25%), `qr.LevelH` (30%).
- `matrix.Size() int`: Side length of the QR grid in modules.
- `matrix.Version() int`: QR version (1 to 40).
- `matrix.Dark(x, y int) bool`: Returns true if module at `(x, y)` is dark.
- `matrix.SVG(opts qr.SVGOptions) string`: Renders matrix as an SVG.
