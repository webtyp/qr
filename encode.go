package qr

type bitBuffer struct {
	bytes []byte
	bits  int
}

func (b *bitBuffer) writeBits(val uint32, numBits int) {
	for i := numBits - 1; i >= 0; i-- {
		bit := byte((val >> i) & 1)
		byteIdx := b.bits / 8
		if byteIdx >= len(b.bytes) {
			b.bytes = append(b.bytes, 0)
		}
		b.bytes[byteIdx] |= (bit << (7 - (b.bits % 8)))
		b.bits++
	}
}

func charCountBits(version int) int {
	if version <= 9 {
		return 8
	}
	return 16
}

func selectVersion(dataLen int, level Level) (int, error) {
	if level > LevelH {
		return 0, ErrUnknownLevel
	}
	for v := 1; v <= 40; v++ {
		totalBits := 4 + charCountBits(v) + dataLen*8
		reqBytes := (totalBits + 7) / 8
		if reqBytes <= dataCapacity(v, level) {
			return v, nil
		}
	}
	return 0, ErrTooLong
}

func buildDataCodewords(data string, version int, level Level) ([]byte, error) {
	capBytes := dataCapacity(version, level)
	if capBytes == 0 {
		return nil, ErrTooLong
	}

	buf := &bitBuffer{bytes: make([]byte, 0, capBytes)}

	// 1. Mode indicator (0100 for Byte mode)
	buf.writeBits(0x4, 4)

	// 2. Character count indicator
	ccBits := charCountBits(version)
	buf.writeBits(uint32(len(data)), ccBits)

	// 3. Payload
	for i := 0; i < len(data); i++ {
		buf.writeBits(uint32(data[i]), 8)
	}

	// 4. Terminator bits (up to 4 zero bits)
	maxDataBits := capBytes * 8
	termBits := maxDataBits - buf.bits
	if termBits > 4 {
		termBits = 4
	}
	if termBits > 0 {
		buf.writeBits(0, termBits)
	}

	// 5. Pad to byte boundary
	if buf.bits%8 != 0 {
		padBits := 8 - (buf.bits % 8)
		buf.writeBits(0, padBits)
	}

	// 6. Pad bytes 0xEC and 0x11
	padByte := byte(0xEC)
	for len(buf.bytes) < capBytes {
		buf.writeBits(uint32(padByte), 8)
		if padByte == 0xEC {
			padByte = 0x11
		} else {
			padByte = 0xEC
		}
	}

	return buf.bytes, nil
}

func encode(data string, level Level) (*Matrix, error) {
	if len(data) == 0 {
		return nil, ErrEmptyData
	}
	if level > LevelH {
		return nil, ErrUnknownLevel
	}

	version, err := selectVersion(len(data), level)
	if err != nil {
		return nil, err
	}

	codewords, err := buildDataCodewords(data, version, level)
	if err != nil {
		return nil, err
	}

	return buildMatrix(codewords, version, level)
}

func buildMatrix(codewords []byte, version int, level Level) (*Matrix, error) {
	return buildMatrixComplete(codewords, version, level)
}
