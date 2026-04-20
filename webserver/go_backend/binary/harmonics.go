package binary

import (
	"encoding/binary"
	"fmt"
)

const (
	PhaseU    = 3
	ChannelI  = 4
	Orders    = 64
	HarmTotal = (PhaseU + ChannelI) * Orders // 448 uint16 = 896 bytes
	HarmSize  = HarmTotal * 2
)

// ParseHarmonics parses a 896-byte harmonics block.
// struct { uint16_t U[3][64]; uint16_t I[4][64]; }
func ParseHarmonics(data []byte) (map[string]interface{}, error) {
	return parseBlock(data)
}

// ParseInterharmonics parses a 896-byte interharmonics block (same layout).
func ParseInterharmonics(data []byte) (map[string]interface{}, error) {
	return parseBlock(data)
}

func parseBlock(data []byte) (map[string]interface{}, error) {
	if len(data) < HarmSize {
		return nil, fmt.Errorf("harmonics: expected %d bytes, got %d", HarmSize, len(data))
	}

	values := make([]uint16, HarmTotal)
	for i := 0; i < HarmTotal; i++ {
		values[i] = binary.LittleEndian.Uint16(data[i*2:])
	}

	u := map[string]interface{}{}
	for phase := 0; phase < PhaseU; phase++ {
		start := phase * Orders
		arr := make([]uint16, Orders)
		copy(arr, values[start:start+Orders])
		u[fmt.Sprintf("U%d", phase+1)] = arr
	}

	iMap := map[string]interface{}{}
	iOffset := PhaseU * Orders
	for ch := 0; ch < ChannelI; ch++ {
		start := iOffset + ch*Orders
		arr := make([]uint16, Orders)
		copy(arr, values[start:start+Orders])
		iMap[fmt.Sprintf("I%d", ch+1)] = arr
	}

	return map[string]interface{}{"U": u, "I": iMap}, nil
}
