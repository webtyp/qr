package qr

func isMasked(mask int, x, y int) bool {
	switch mask {
	case 0:
		return (x+y)%2 == 0
	case 1:
		return y%2 == 0
	case 2:
		return x%3 == 0
	case 3:
		return (x+y)%3 == 0
	case 4:
		return (y/2+x/3)%2 == 0
	case 5:
		return ((x*y)%2)+((x*y)%3) == 0
	case 6:
		return (((x*y)%2)+((x*y)%3))%2 == 0
	case 7:
		return (((x+y)%2)+((x*y)%3))%2 == 0
	default:
		return false
	}
}

// Penalty N1: 5 or more consecutive modules of same color in row or column.
func penaltyN1(g *moduleGrid) int {
	s := g.size
	score := 0

	// Horizontal
	for y := 0; y < s; y++ {
		count := 0
		var lastColor bool
		for x := 0; x < s; x++ {
			color := g.modules[y*s+x]
			if x == 0 || color != lastColor {
				if count >= 5 {
					score += 3 + (count - 5)
				}
				count = 1
				lastColor = color
			} else {
				count++
			}
		}
		if count >= 5 {
			score += 3 + (count - 5)
		}
	}

	// Vertical
	for x := 0; x < s; x++ {
		count := 0
		var lastColor bool
		for y := 0; y < s; y++ {
			color := g.modules[y*s+x]
			if y == 0 || color != lastColor {
				if count >= 5 {
					score += 3 + (count - 5)
				}
				count = 1
				lastColor = color
			} else {
				count++
			}
		}
		if count >= 5 {
			score += 3 + (count - 5)
		}
	}

	return score
}

// Penalty N2: 2x2 blocks of same color modules.
func penaltyN2(g *moduleGrid) int {
	s := g.size
	score := 0
	for y := 0; y < s-1; y++ {
		for x := 0; x < s-1; x++ {
			c := g.modules[y*s+x]
			if g.modules[y*s+x+1] == c && g.modules[(y+1)*s+x] == c && g.modules[(y+1)*s+x+1] == c {
				score += 3
			}
		}
	}
	return score
}

// Penalty N3: 1:1:3:1:1 pattern flanked by 4 light modules on either side.
func penaltyN3(g *moduleGrid) int {
	s := g.size
	score := 0

	patternMatch := func(m0, m1, m2, m3, m4, m5, m6 bool) bool {
		return m0 && !m1 && m2 && m3 && m4 && !m5 && m6
	}

	// Horizontal
	for y := 0; y < s; y++ {
		for x := 0; x <= s-7; x++ {
			if patternMatch(
				g.modules[y*s+x],
				g.modules[y*s+x+1],
				g.modules[y*s+x+2],
				g.modules[y*s+x+3],
				g.modules[y*s+x+4],
				g.modules[y*s+x+5],
				g.modules[y*s+x+6],
			) {
				beforeLight := true
				if x >= 4 {
					for i := 1; i <= 4; i++ {
						if g.modules[y*s+x-i] {
							beforeLight = false
							break
						}
					}
				} else {
					beforeLight = false
				}

				afterLight := true
				if x+7+4 <= s {
					for i := 0; i < 4; i++ {
						if g.modules[y*s+x+7+i] {
							afterLight = false
							break
						}
					}
				} else {
					afterLight = false
				}

				if beforeLight || afterLight {
					score += 40
				}
			}
		}
	}

	// Vertical
	for x := 0; x < s; x++ {
		for y := 0; y <= s-7; y++ {
			if patternMatch(
				g.modules[y*s+x],
				g.modules[(y+1)*s+x],
				g.modules[(y+2)*s+x],
				g.modules[(y+3)*s+x],
				g.modules[(y+4)*s+x],
				g.modules[(y+5)*s+x],
				g.modules[(y+6)*s+x],
			) {
				beforeLight := true
				if y >= 4 {
					for i := 1; i <= 4; i++ {
						if g.modules[(y-i)*s+x] {
							beforeLight = false
							break
						}
					}
				} else {
					beforeLight = false
				}

				afterLight := true
				if y+7+4 <= s {
					for i := 0; i < 4; i++ {
						if g.modules[(y+7+i)*s+x] {
							afterLight = false
							break
						}
					}
				} else {
					afterLight = false
				}

				if beforeLight || afterLight {
					score += 40
				}
			}
		}
	}

	return score
}

// Penalty N4: Ratio of dark modules to total modules.
func penaltyN4(g *moduleGrid) int {
	s := g.size
	darkCount := 0
	total := s * s
	for i := 0; i < total; i++ {
		if g.modules[i] {
			darkCount++
		}
	}
	percent := (darkCount * 100) / total
	diff1 := percent - 50
	if diff1 < 0 {
		diff1 = -diff1
	}

	p1 := diff1 / 5
	return p1 * 10
}

func evaluatePenalty(g *moduleGrid) int {
	return penaltyN1(g) + penaltyN2(g) + penaltyN3(g) + penaltyN4(g)
}

func applyMask(g *moduleGrid, mask int) {
	s := g.size
	for y := 0; y < s; y++ {
		for x := 0; x < s; x++ {
			if !g.isReserved(x, y) {
				if isMasked(mask, x, y) {
					idx := y*s + x
					g.modules[idx] = !g.modules[idx]
				}
			}
		}
	}
}

func selectBestMask(codewords []byte, version int, level Level) (*moduleGrid, int) {
	bestScore := -1
	var bestGrid *moduleGrid
	bestMask := 0

	for mask := 0; mask < 8; mask++ {
		grid := newModuleGrid(21 + (version-1)*4)
		grid.placeFunctionPatterns(version)
		grid.placeVersionInfo(version)
		grid.placeData(codewords)

		applyMask(grid, mask)
		grid.placeFormatInfo(level, mask)

		score := evaluatePenalty(grid)
		if bestScore == -1 || score < bestScore {
			bestScore = score
			bestGrid = grid
			bestMask = mask
		}
	}

	return bestGrid, bestMask
}

func (m *Matrix) fillFromGrid(g *moduleGrid) {
	m.version = (g.size - 17) / 4
	m.size = g.size
	m.data = make([]bool, len(g.modules))
	copy(m.data, g.modules)
}

func buildMatrixComplete(codewords []byte, version int, level Level) (*Matrix, error) {
	interleaved := interleave(codewords, version, level)
	grid, _ := selectBestMask(interleaved, version, level)

	m := &Matrix{}
	m.fillFromGrid(grid)
	return m, nil
}
