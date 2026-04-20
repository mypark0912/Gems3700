package data

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/redis/go-redis/v9"
)

// MaxDemandEntry holds a single max-demand record (uint32 timestamp + float32 value).
type MaxDemandEntry struct {
	Timestamp uint32  `json:"timestamp"`
	Value     float32 `json:"value"`
}

// DemandData holds the fully parsed 124-byte demand struct.
type DemandData struct {
	// Max Demand
	MD_P [2]MaxDemandEntry // import / export
	MD_Q [2]MaxDemandEntry
	MD_S MaxDemandEntry
	MD_I [3]MaxDemandEntry // L1 / L2 / L3

	// Dynamic Demand
	DDTimestamp uint32
	DD_P       [2]float32
	DD_Q       [2]float32
	DD_S       float32
	DD_I       [3]float32

	// Present (Current) Demand
	CD_P [2]float32
	CD_Q [2]float32
	CD_S float32

	// Predict Demand
	PD_P float32
}

// DemandParser parses a 124-byte little-endian binary struct into DemandData.
type DemandParser struct{}

const demandTotalSize = 124

// NewDemandParser creates a new DemandParser.
func NewDemandParser() *DemandParser {
	return &DemandParser{}
}

// Parse decodes a 124-byte binary blob into *DemandData.
func (dp *DemandParser) Parse(data []byte) (*DemandData, error) {
	if len(data) != demandTotalSize {
		return nil, fmt.Errorf("Demand: expected %d bytes, got %d", demandTotalSize, len(data))
	}

	d := &DemandData{}
	off := 0

	// helper: read uint32 + float32 pair (8 bytes)
	readEntry := func() MaxDemandEntry {
		ts := binary.LittleEndian.Uint32(data[off : off+4])
		val := math.Float32frombits(binary.LittleEndian.Uint32(data[off+4 : off+8]))
		off += 8
		return MaxDemandEntry{Timestamp: ts, Value: val}
	}

	// helper: read a single little-endian float32
	readFloat := func() float32 {
		v := math.Float32frombits(binary.LittleEndian.Uint32(data[off : off+4]))
		off += 4
		return v
	}

	// --- Max Demand ---
	for i := 0; i < 2; i++ {
		d.MD_P[i] = readEntry()
	}
	for i := 0; i < 2; i++ {
		d.MD_Q[i] = readEntry()
	}
	d.MD_S = readEntry()
	for i := 0; i < 3; i++ {
		d.MD_I[i] = readEntry()
	}

	// --- Dynamic Demand ---
	d.DDTimestamp = binary.LittleEndian.Uint32(data[off : off+4])
	off += 4

	for i := 0; i < 2; i++ {
		d.DD_P[i] = readFloat()
	}
	for i := 0; i < 2; i++ {
		d.DD_Q[i] = readFloat()
	}
	d.DD_S = readFloat()
	for i := 0; i < 3; i++ {
		d.DD_I[i] = readFloat()
	}

	// --- Present (Current) Demand ---
	for i := 0; i < 2; i++ {
		d.CD_P[i] = readFloat()
	}
	for i := 0; i < 2; i++ {
		d.CD_Q[i] = readFloat()
	}
	d.CD_S = readFloat()

	// --- Predict Demand ---
	d.PD_P = readFloat()

	return d, nil
}

// FlattenDemand converts a parsed DemandData into a flat map suitable for InfluxDB fields.
func FlattenDemand(d *DemandData) map[string]interface{} {
	flat := make(map[string]interface{}, 32)

	ieLabels := [2]string{"import", "export"}
	phaseLabels := [3]string{"L1", "L2", "L3"}

	// Max Demand — P, Q
	for i, label := range ieLabels {
		flat[fmt.Sprintf("MD_P_%s_value", label)] = d.MD_P[i].Value
		flat[fmt.Sprintf("MD_P_%s_ts", label)] = d.MD_P[i].Timestamp
	}
	for i, label := range ieLabels {
		flat[fmt.Sprintf("MD_Q_%s_value", label)] = d.MD_Q[i].Value
		flat[fmt.Sprintf("MD_Q_%s_ts", label)] = d.MD_Q[i].Timestamp
	}
	flat["MD_S_value"] = d.MD_S.Value
	flat["MD_S_ts"] = d.MD_S.Timestamp
	for i, label := range phaseLabels {
		flat[fmt.Sprintf("MD_I_%s_value", label)] = d.MD_I[i].Value
		flat[fmt.Sprintf("MD_I_%s_ts", label)] = d.MD_I[i].Timestamp
	}

	// Dynamic Demand
	flat["ddTimestamp"] = d.DDTimestamp
	for i, label := range ieLabels {
		flat[fmt.Sprintf("DD_P_%s", label)] = d.DD_P[i]
		flat[fmt.Sprintf("DD_Q_%s", label)] = d.DD_Q[i]
	}
	flat["DD_S"] = d.DD_S
	for i, label := range phaseLabels {
		flat[fmt.Sprintf("DD_I_%s", label)] = d.DD_I[i]
	}

	// Present (Current) Demand
	for i, label := range ieLabels {
		flat[fmt.Sprintf("CD_P_%s", label)] = d.CD_P[i]
		flat[fmt.Sprintf("CD_Q_%s", label)] = d.CD_Q[i]
	}
	flat["CD_S"] = d.CD_S

	// Predict Demand
	flat["PD_P"] = d.PD_P

	return flat
}

// BinaryCollector reads and parses binary data from Redis.
type BinaryCollector struct {
	redis         *redis.Client
	demandParser  *DemandParser
}

// NewBinaryCollector creates a new BinaryCollector backed by the given Redis client.
func NewBinaryCollector(redisClient *redis.Client) *BinaryCollector {
	return &BinaryCollector{
		redis:        redisClient,
		demandParser: NewDemandParser(),
	}
}

// GetDemand reads the binary demand blob from Redis HGET "Demand" <channel>
// and parses it into a DemandData struct.
func (bc *BinaryCollector) GetDemand(channel string) (*DemandData, error) {
	data, err := bc.redis.HGet(context.Background(), "Demand", channel).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("Redis hget(Demand, %s) error: %w", channel, err)
	}
	return bc.demandParser.Parse(data)
}

// GetDemandFlat reads and parses demand data, then flattens it into a map for InfluxDB.
func (bc *BinaryCollector) GetDemandFlat(channel string) (map[string]interface{}, error) {
	parsed, err := bc.GetDemand(channel)
	if err != nil {
		return nil, err
	}
	if parsed == nil {
		return nil, nil
	}
	return FlattenDemand(parsed), nil
}
