package qr

import (
	"webtyp.com/fmt"
)

const (
	defaultModuleSize = 4
	minQuietZone      = 4
	defaultDarkColor  = "#000000"
	defaultLightColor = "#ffffff"

	svgHeaderFormat = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">`
	svgRectFormat   = `<rect width="%d" height="%d" fill="%s"/>`
	svgPathStart    = `<path fill="%s" shape-rendering="crispEdges" d="`
	svgPathEnd      = `"/></svg>`
)

func (m *Matrix) renderSVG(opts SVGOptions) string {
	if m == nil || m.size == 0 {
		return ""
	}

	moduleSize := opts.ModuleSize
	if moduleSize <= 0 {
		moduleSize = defaultModuleSize
	}

	quietZone := opts.QuietZone
	if quietZone < minQuietZone {
		quietZone = minQuietZone
	}

	darkColor := opts.Dark
	if darkColor == "" {
		darkColor = defaultDarkColor
	}

	lightColor := opts.Light
	if lightColor == "" {
		lightColor = defaultLightColor
	}

	totalModules := m.size + 2*quietZone
	totalPixels := totalModules * moduleSize

	header := fmt.Sprintf(svgHeaderFormat, totalPixels, totalPixels, totalPixels, totalPixels)
	bg := fmt.Sprintf(svgRectFormat, totalPixels, totalPixels, lightColor)
	pathHead := fmt.Sprintf(svgPathStart, darkColor)

	var dBuf []byte

	for y := 0; y < m.size; y++ {
		for x := 0; x < m.size; x++ {
			if m.Dark(x, y) {
				px := (x + quietZone) * moduleSize
				py := (y + quietZone) * moduleSize
				cmd := fmt.Sprintf("M%d %dh%dv%dh-%dz", px, py, moduleSize, moduleSize, moduleSize)
				dBuf = append(dBuf, []byte(cmd)...)
			}
		}
	}

	return header + bg + pathHead + string(dBuf) + svgPathEnd
}
