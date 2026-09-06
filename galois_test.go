package qr

import (
	"testing"
)

func TestGF256(t *testing.T) {
	// mul(x, 0) == 0, mul(x, 1) == x
	for x := 0; x < 256; x++ {
		bx := byte(x)
		if res := gfMul(bx, 0); res != 0 {
			t.Fatalf("gfMul(0x%02X, 0) = 0x%02X; want 0", bx, res)
		}
		if res := gfMul(bx, 1); res != bx {
			t.Fatalf("gfMul(0x%02X, 1) = 0x%02X; want 0x%02X", bx, res, bx)
		}
	}

	// log/antilog round-trip for all 255 non-zero elements
	for i := 1; i < 256; i++ {
		bi := byte(i)
		logVal := logTable[bi]
		expVal := expTable[logVal]
		if expVal != bi {
			t.Fatalf("expTable[logTable[0x%02X]] = 0x%02X; want 0x%02X", bi, expVal, bi)
		}
	}
}

func TestReedSolomonISO18004Example(t *testing.T) {
	// Worked example from ISO/IEC 18004 Annex I:
	// Version 1-M symbol for "01234567"
	// Data codewords (16 bytes):
	data := []byte{
		0x10, 0x20, 0x0C, 0x56, 0x61, 0x80, 0xEC, 0x11,
		0xEC, 0x11, 0xEC, 0x11, 0xEC, 0x11, 0xEC, 0x11,
	}
	// Expected 10 EC codewords generated via RS generator polynomial (2^0..2^9):
	expectedEC := []byte{
		0xA5, 0x24, 0xD4, 0xC1, 0xED, 0x36, 0xC7, 0x87, 0x2C, 0x55,
	}

	ec := rsEncode(data, 10)
	if len(ec) != len(expectedEC) {
		t.Fatalf("got %d EC bytes; want %d", len(ec), len(expectedEC))
	}
	for i := range ec {
		if ec[i] != expectedEC[i] {
			t.Errorf("EC byte %d = 0x%02X; want 0x%02X", i, ec[i], expectedEC[i])
		}
	}
}
