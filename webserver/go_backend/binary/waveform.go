package binary

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"strings"

	"github.com/redis/go-redis/v9"
)

// WAVEFORM_L16 Python layout (3000 bytes):
//   int16 U[3][160]     : 960 bytes
//   int16 Upp[3][160]   : 960 bytes
//   int16 I[3][160]     : 960 bytes
//   float vscale        : 4 bytes
//   float iscale        : 4 bytes
//   int16 r[56]         : 112 bytes (reserved)
const (
	WaveformSize      = 3000
	WaveformSamples   = 160
	WaveformPhaseCnt  = 3
	waveformPhaseSize = WaveformSamples * 2 // bytes per phase row
	waveformBlockSize = WaveformPhaseCnt * waveformPhaseSize
)

// ParseWaveform parses the 3000-byte WAVEFORM_L16 blob.
// Returned map matches Python WaveformParser output keys.
func ParseWaveform(data []byte) (map[string]interface{}, error) {
	if len(data) < WaveformSize {
		return nil, fmt.Errorf("waveform: expected %d bytes, got %d", WaveformSize, len(data))
	}

	uData := readWaveBlock(data, 0)
	uppData := readWaveBlock(data, waveformBlockSize)
	iData := readWaveBlock(data, 2*waveformBlockSize)

	scaleOff := 3 * waveformBlockSize
	vscale := math.Float32frombits(binary.LittleEndian.Uint32(data[scaleOff : scaleOff+4]))
	iscale := math.Float32frombits(binary.LittleEndian.Uint32(data[scaleOff+4 : scaleOff+8]))

	return map[string]interface{}{
		"Wave Form V1":  uData[0],
		"Wave Form V2":  uData[1],
		"Wave Form V3":  uData[2],
		"Wave Form V12": uppData[0],
		"Wave Form V23": uppData[1],
		"Wave Form V31": uppData[2],
		"Wave Form I1":  iData[0],
		"Wave Form I2":  iData[1],
		"Wave Form I3":  iData[2],
		"vscale":        float64(vscale),
		"iscale":        float64(iscale),
	}, nil
}

func readWaveBlock(data []byte, offset int) [WaveformPhaseCnt][]int16 {
	var phases [WaveformPhaseCnt][]int16
	for p := 0; p < WaveformPhaseCnt; p++ {
		samples := make([]int16, WaveformSamples)
		base := offset + p*waveformPhaseSize
		for s := 0; s < WaveformSamples; s++ {
			samples[s] = int16(binary.LittleEndian.Uint16(data[base+s*2:]))
		}
		phases[p] = samples
	}
	return phases
}

// waveformRedisKey builds "pq_main" / "pq_sub" from channel name.
func waveformRedisKey(channel string) string {
	ch := strings.ToLower(channel)
	if ch != "main" && ch != "sub" {
		ch = "main"
	}
	return "pq_" + ch
}

// FetchWaveform fetches HGET pq_<channel> "waveform" and parses it.
func FetchWaveform(ctx context.Context, client *redis.Client, channel string) (map[string]interface{}, error) {
	raw, err := client.HGet(ctx, waveformRedisKey(channel), "waveform").Bytes()
	if err != nil {
		return nil, err
	}
	return ParseWaveform(raw)
}
