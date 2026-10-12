package sword

import (
	"encoding/binary"
	"fmt"
)

func readU16LE(buf []byte, off int) (uint16, error) {
	if off+2 > len(buf) {
		return 0, fmt.Errorf("read u16 at %d: short buffer", off)
	}

	return binary.LittleEndian.Uint16(buf[off:]), nil
}

func readU32LE(buf []byte, off int) (uint32, error) {
	if off+4 > len(buf) {
		return 0, fmt.Errorf("read u32 at %d: short buffer", off)
	}

	return binary.LittleEndian.Uint32(buf[off:]), nil
}

func readU48LE(buf []byte, off int) (uint64, error) {
	if off+6 > len(buf) {
		return 0, fmt.Errorf("read u48 at %d: short buffer", off)
	}
	b := buf[off : off+6]
	// cross-connect / swordjs zText index layout (not plain LE uint48).
	return uint64(b[1])<<40 |
		uint64(b[0])<<32 |
		uint64(b[5])<<24 |
		uint64(b[4])<<16 |
		uint64(b[3])<<8 |
		uint64(b[2]), nil
}
