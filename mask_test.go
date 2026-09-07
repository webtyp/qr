package qr

import (
	"testing"
)

func TestMaskPenalties(t *testing.T) {
	// Rule 1: N1 penalty for 5+ consecutive same-color modules
	g1 := newModuleGrid(21)
	for x := 0; x < 5; x++ {
		g1.modules[x] = true
	}
	if p1 := penaltyN1(g1); p1 < 3 {
		t.Fatalf("penaltyN1 = %d; want >= 3", p1)
	}

	// Rule 2: N2 penalty for 2x2 blocks of same color
	g2 := newModuleGrid(21)
	for y := 0; y < 21; y++ {
		for x := 0; x < 21; x++ {
			g2.modules[y*21+x] = (x+y)%2 == 0
		}
	}
	g2.modules[0*21+0] = true
	g2.modules[0*21+1] = true
	g2.modules[1*21+0] = true
	g2.modules[1*21+1] = true
	if p2 := penaltyN2(g2); p2 != 3 {
		t.Fatalf("penaltyN2 = %d; want 3", p2)
	}

	// Rule 3: N3 penalty for 1:1:3:1:1 pattern with 4 light modules
	g3 := newModuleGrid(21)
	g3.modules[4] = true
	g3.modules[5] = false
	g3.modules[6] = true
	g3.modules[7] = true
	g3.modules[8] = true
	g3.modules[9] = false
	g3.modules[10] = true
	if p3 := penaltyN3(g3); p3 < 40 {
		t.Fatalf("penaltyN3 = %d; want >= 40", p3)
	}

	// Rule 4: N4 penalty for dark module proportion
	g4 := newModuleGrid(21)
	if p4 := penaltyN4(g4); p4 < 100 {
		t.Fatalf("penaltyN4 = %d; want >= 100", p4)
	}
}
