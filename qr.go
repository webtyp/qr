package qr

import (
	"webtyp.com/fmt"
)

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

var (
	ErrEmptyData    = fmt.Errf("qr: empty data")
	ErrTooLong      = fmt.Errf("qr: data too long")
	ErrUnknownLevel = fmt.Errf("qr: unknown level")
)

// Matrix is a square grid of QR modules.
type Matrix struct {
	version int
	size    int
	data    []bool
}

// Size is the side length in modules: 21 for version 1, up to 177 for 40.
func (m *Matrix) Size() int {
	if m == nil {
		return 0
	}
	return m.size
}

// Dark reports whether the module at (x, y) is dark. Out-of-range coordinates
// report false — the quiet zone is light.
func (m *Matrix) Dark(x, y int) bool {
	if m == nil || x < 0 || y < 0 || x >= m.size || y >= m.size {
		return false
	}
	return m.data[y*m.size+x]
}

// Version is the QR version, 1 to 40.
func (m *Matrix) Version() int {
	if m == nil {
		return 0
	}
	return m.version
}

// SVGOptions controls the rendered SVG. The zero value renders a scannable
// code: 4-module quiet zone, black on white, sized to the module count.
type SVGOptions struct {
	ModuleSize int    // pixels per module; 0 = 4
	QuietZone  int    // modules of margin; 0 = 4, the specification's minimum
	Dark       string // CSS colour for dark modules; "" = "#000000"
	Light      string // CSS colour for the background; "" = "#ffffff"
}

// SVG renders the matrix as a standalone SVG document.
func (m *Matrix) SVG(opts SVGOptions) string {
	return m.renderSVG(opts)
}

// Encode builds the QR matrix for data.
//
// It selects the smallest version (1–40) that holds data at the given level,
// and the mask that scores best under the specification's penalty rules.
func Encode(data string, level Level) (*Matrix, error) {
	return encode(data, level)
}
