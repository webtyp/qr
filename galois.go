package qr

var (
	expTable [256]byte
	logTable [256]byte
)

func init() {
	x := 1
	for i := 0; i < 255; i++ {
		expTable[i] = byte(x)
		logTable[x] = byte(i)
		x <<= 1
		if x&0x100 != 0 {
			x ^= 0x11D
		}
	}
	expTable[255] = expTable[0]
}

func gfMul(a, b byte) byte {
	if a == 0 || b == 0 {
		return 0
	}
	sum := int(logTable[a]) + int(logTable[b])
	if sum >= 255 {
		sum -= 255
	}
	return expTable[sum]
}

// polyMul multiplies poly a and poly b over GF(256).
func polyMul(a, b []byte) []byte {
	res := make([]byte, len(a)+len(b)-1)
	for i := 0; i < len(a); i++ {
		if a[i] == 0 {
			continue
		}
		for j := 0; j < len(b); j++ {
			if b[j] == 0 {
				continue
			}
			res[i+j] ^= gfMul(a[i], b[j])
		}
	}
	return res
}

// rsGeneratorPoly builds the generator polynomial for numEC error correction codewords.
func rsGeneratorPoly(numEC int) []byte {
	g := []byte{1}
	for i := 0; i < numEC; i++ {
		factor := []byte{1, expTable[i]}
		g = polyMul(g, factor)
	}
	return g
}

// rsEncode generates numEC error correction codewords for the given data codewords.
func rsEncode(data []byte, numEC int) []byte {
	gen := rsGeneratorPoly(numEC)
	buf := make([]byte, len(data)+numEC)
	copy(buf, data)

	for i := 0; i < len(data); i++ {
		coef := buf[i]
		if coef == 0 {
			continue
		}
		for j := 0; j < len(gen); j++ {
			buf[i+j] ^= gfMul(gen[j], coef)
		}
	}

	ec := make([]byte, numEC)
	copy(ec, buf[len(data):])
	return ec
}
