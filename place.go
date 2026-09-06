package qr

// Alignment pattern centers for versions 1-40.
var alignmentPatternCenters = [40][]int{
	{},                       // V1 (none)
	{6, 18},                  // V2
	{6, 22},                  // V3
	{6, 26},                  // V4
	{6, 30},                  // V5
	{6, 34},                  // V6
	{6, 22, 38},              // V7
	{6, 24, 42},              // V8
	{6, 26, 46},              // V9
	{6, 28, 50},              // V10
	{6, 30, 54},              // V11
	{6, 32, 58},              // V12
	{6, 34, 62},              // V13
	{6, 26, 46, 66},          // V14
	{6, 26, 48, 70},          // V15
	{6, 26, 50, 74},          // V16
	{6, 30, 54, 78},          // V17
	{6, 30, 56, 82},          // V18
	{6, 30, 58, 86},          // V19
	{6, 34, 62, 90},          // V20
	{6, 28, 50, 72, 94},      // V21
	{6, 26, 50, 74, 98},      // V22
	{6, 30, 54, 78, 102},     // V23
	{6, 28, 54, 80, 106},     // V24
	{6, 32, 58, 84, 110},     // V25
	{6, 30, 58, 86, 114},     // V26
	{6, 34, 62, 90, 118},     // V27
	{6, 26, 50, 74, 98, 122}, // V28
	{6, 30, 54, 78, 102, 126},// V29
	{6, 26, 52, 78, 104, 130},// V30
	{6, 30, 56, 82, 108, 134},// V31
	{6, 34, 60, 86, 112, 138},// V32
	{6, 30, 58, 86, 114, 142},// V33
	{6, 34, 62, 90, 118, 146},// V34
	{6, 30, 54, 78, 102, 126, 150}, // V35
	{6, 24, 50, 76, 102, 128, 154}, // V36
	{6, 28, 54, 80, 106, 132, 158}, // V37
	{6, 32, 58, 84, 110, 136, 162}, // V38
	{6, 26, 54, 82, 110, 138, 166}, // V39
	{6, 30, 58, 86, 114, 142, 170}, // V40
}

type moduleGrid struct {
	size      int
	modules   []bool
	reserved  []bool
}

func newModuleGrid(size int) *moduleGrid {
	return &moduleGrid{
		size:     size,
		modules:  make([]bool, size*size),
		reserved: make([]bool, size*size),
	}
}

func (g *moduleGrid) set(x, y int, dark bool, reserve bool) {
	if x < 0 || y < 0 || x >= g.size || y >= g.size {
		return
	}
	idx := y*g.size + x
	g.modules[idx] = dark
	if reserve {
		g.reserved[idx] = true
	}
}

func (g *moduleGrid) isReserved(x, y int) bool {
	if x < 0 || y < 0 || x >= g.size || y >= g.size {
		return true
	}
	return g.reserved[y*g.size+x]
}

func (g *moduleGrid) drawFinderPattern(x, y int) {
	for dy := -1; dy <= 7; dy++ {
		for dx := -1; dx <= 7; dx++ {
			px, py := x+dx, y+dy
			if px < 0 || py < 0 || px >= g.size || py >= g.size {
				continue
			}
			if dx == -1 || dx == 7 || dy == -1 || dy == 7 {
				g.set(px, py, false, true) // separator
			} else if dx == 0 || dx == 6 || dy == 0 || dy == 6 || (dx >= 2 && dx <= 4 && dy >= 2 && dy <= 4) {
				g.set(px, py, true, true)
			} else {
				g.set(px, py, false, true)
			}
		}
	}
}

func (g *moduleGrid) drawAlignmentPattern(cx, cy int) {
	for dy := -2; dy <= 2; dy++ {
		for dx := -2; dx <= 2; dx++ {
			px, py := cx+dx, cy+dy
			if g.isReserved(px, py) {
				return
			}
		}
	}
	for dy := -2; dy <= 2; dy++ {
		for dx := -2; dx <= 2; dx++ {
			px, py := cx+dx, cy+dy
			if dx == -2 || dx == 2 || dy == -2 || dy == 2 || (dx == 0 && dy == 0) {
				g.set(px, py, true, true)
			} else {
				g.set(px, py, false, true)
			}
		}
	}
}

func (g *moduleGrid) placeFunctionPatterns(version int) {
	s := g.size

	// 1. Finder patterns + separators
	g.drawFinderPattern(0, 0)
	g.drawFinderPattern(s-7, 0)
	g.drawFinderPattern(0, s-7)

	// 2. Alignment patterns
	centers := alignmentPatternCenters[version-1]
	for _, cy := range centers {
		for _, cx := range centers {
			g.drawAlignmentPattern(cx, cy)
		}
	}

	// 3. Timing patterns
	for i := 8; i < s-8; i++ {
		if !g.isReserved(i, 6) {
			g.set(i, 6, i%2 == 0, true)
		}
		if !g.isReserved(6, i) {
			g.set(6, i, i%2 == 0, true)
		}
	}

	// 4. Dark module
	g.set(8, 4*version+9, true, true)

	// 5. Reserve format info areas
	for i := 0; i < 9; i++ {
		if i != 6 {
			g.set(i, 8, false, true)
			g.set(8, i, false, true)
		}
	}
	for i := 0; i < 8; i++ {
		g.set(s-1-i, 8, false, true)
		g.set(8, s-1-i, false, true)
	}
	g.set(8, 8, false, true)

	// 6. Reserve version info areas (Version >= 7)
	if version >= 7 {
		for r := 0; r < 6; r++ {
			for c := 0; c < 3; c++ {
				g.set(c, s-11+r, false, true)
				g.set(s-11+r, c, false, true)
			}
		}
	}
}

// formatBits calculates the 15-bit format codeword with BCH(15,5) and XOR 0x5371.
func formatBits(level Level, mask int) uint16 {
	var levBits uint32
	switch level {
	case LevelL:
		levBits = 1
	case LevelM:
		levBits = 0
	case LevelQ:
		levBits = 3
	case LevelH:
		levBits = 2
	}

	data := (levBits << 3) | uint32(mask&7)
	rem := data << 10
	for i := 14; i >= 10; i-- {
		if (rem >> i) & 1 != 0 {
			rem ^= (0x537 << (i - 10))
		}
	}
	return uint16(((data << 10) | rem) ^ 0x5371)
}

func (g *moduleGrid) placeFormatInfo(level Level, mask int) {
	fmtVal := formatBits(level, mask)
	s := g.size

	for i := 0; i < 15; i++ {
		bit := ((fmtVal >> i) & 1) != 0

		if i < 6 {
			g.set(i, 8, bit, false)
		} else if i < 8 {
			g.set(i+1, 8, bit, false)
		} else if i == 8 {
			g.set(8, 7, bit, false)
		} else {
			g.set(8, 14-i, bit, false)
		}

		if i < 8 {
			g.set(8, s-1-i, bit, false)
		} else {
			g.set(s-15+i, 8, bit, false)
		}
	}
}

// versionBits calculates the 18-bit version codeword for version >= 7 using BCH(18,6).
func versionBits(version int) uint32 {
	if version < 7 {
		return 0
	}
	data := uint32(version)
	rem := data << 12
	for i := 17; i >= 12; i-- {
		if (rem >> i) & 1 != 0 {
			rem ^= (0x1F25 << (i - 12))
		}
	}
	return (data << 12) | rem
}

func (g *moduleGrid) placeVersionInfo(version int) {
	if version < 7 {
		return
	}
	vBits := versionBits(version)
	s := g.size

	for i := 0; i < 18; i++ {
		bit := ((vBits >> i) & 1) != 0
		r := i / 3
		c := i % 3

		g.set(c, s-11+r, bit, false)
		g.set(s-11+r, c, bit, false)
	}
}

func (g *moduleGrid) placeData(codewords []byte) {
	s := g.size
	bitIdx := 0
	totalBits := len(codewords) * 8

	dirUp := true
	x := s - 1

	for x > 0 {
		if x == 6 {
			x--
		}

		for yCount := 0; yCount < s; yCount++ {
			y := yCount
			if dirUp {
				y = s - 1 - yCount
			}

			for c := 0; c < 2; c++ {
				px := x - c
				if !g.isReserved(px, y) {
					var bit bool
					if bitIdx < totalBits {
						byteVal := codewords[bitIdx/8]
						bit = ((byteVal >> (7 - (bitIdx % 8))) & 1) != 0
						bitIdx++
					}
					g.set(px, y, bit, false)
				}
			}
		}

		dirUp = !dirUp
		x -= 2
	}
}
