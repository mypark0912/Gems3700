package binary

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
)

// DemandRedisSize mirrors Python DemandRedis: 124 bytes single blob.
const DemandRedisSize = 124

// MaxDemandItem mirrors MAX_DEMAND_t { uint32 mdTime, float value } (8 bytes).
type MaxDemandItem struct {
	Time      string  `json:"time"`
	Timestamp uint32  `json:"timestamp"`
	Value     float64 `json:"value"`
}

var EmptyDemandItem = MaxDemandItem{}

// DemandData mirrors the parsed result of Python's DemandParser.parse().
type DemandData struct {
	MDP         [2]MaxDemandItem `json:"MD_P"`
	MDQ         [2]MaxDemandItem `json:"MD_Q"`
	MDS         MaxDemandItem    `json:"MD_S"`
	MDI         [3]MaxDemandItem `json:"MD_I"`
	DDTime      string           `json:"ddTime"`
	DDTimestamp uint32           `json:"ddTimestamp"`
	DDP         [2]float64       `json:"DD_P"`
	DDQ         [2]float64       `json:"DD_Q"`
	DDS         float64          `json:"DD_S"`
	DDI         [3]float64       `json:"DD_I"`
	CDP         [2]float64       `json:"CD_P"`
	CDQ         [2]float64       `json:"CD_Q"`
	CDS         float64          `json:"CD_S"`
	PDP         float64          `json:"PD_P"`
}

func parseMaxDemandItem(data []byte) MaxDemandItem {
	ts := binary.LittleEndian.Uint32(data[0:4])
	val := math.Float32frombits(binary.LittleEndian.Uint32(data[4:8]))
	item := MaxDemandItem{Timestamp: ts, Value: float64(val)}
	if ts > 0 {
		item.Time = time.Unix(int64(ts), 0).Format("2006-01-02 15:04:05")
	}
	return item
}

func readFloat32LE(data []byte) float64 {
	return float64(math.Float32frombits(binary.LittleEndian.Uint32(data)))
}

// ParseDemand parses the 124-byte Demand blob exactly as Python DemandParser does.
func ParseDemand(data []byte) (*DemandData, error) {
	if len(data) != DemandRedisSize {
		return nil, fmt.Errorf("demand: expected %d bytes, got %d", DemandRedisSize, len(data))
	}

	d := &DemandData{}
	off := 0

	for i := 0; i < 2; i++ {
		d.MDP[i] = parseMaxDemandItem(data[off : off+8])
		off += 8
	}
	for i := 0; i < 2; i++ {
		d.MDQ[i] = parseMaxDemandItem(data[off : off+8])
		off += 8
	}
	d.MDS = parseMaxDemandItem(data[off : off+8])
	off += 8
	for i := 0; i < 3; i++ {
		d.MDI[i] = parseMaxDemandItem(data[off : off+8])
		off += 8
	}

	d.DDTimestamp = binary.LittleEndian.Uint32(data[off : off+4])
	if d.DDTimestamp > 0 {
		d.DDTime = time.Unix(int64(d.DDTimestamp), 0).Format("2006-01-02 15:04:05")
	}
	off += 4

	for i := 0; i < 2; i++ {
		d.DDP[i] = readFloat32LE(data[off : off+4])
		off += 4
	}
	for i := 0; i < 2; i++ {
		d.DDQ[i] = readFloat32LE(data[off : off+4])
		off += 4
	}
	d.DDS = readFloat32LE(data[off : off+4])
	off += 4
	for i := 0; i < 3; i++ {
		d.DDI[i] = readFloat32LE(data[off : off+4])
		off += 4
	}

	for i := 0; i < 2; i++ {
		d.CDP[i] = readFloat32LE(data[off : off+4])
		off += 4
	}
	for i := 0; i < 2; i++ {
		d.CDQ[i] = readFloat32LE(data[off : off+4])
		off += 4
	}
	d.CDS = readFloat32LE(data[off : off+4])
	off += 4
	d.PDP = readFloat32LE(data[off : off+4])

	return d, nil
}

// FetchDemand fetches and parses HGET Demand <channel> ("Main" or "Sub").
func FetchDemand(ctx context.Context, client *redis.Client, channel string) (*DemandData, error) {
	raw, err := client.HGet(ctx, "Demand", channel).Bytes()
	if err != nil {
		return nil, err
	}
	return ParseDemand(raw)
}

// FormatDemandData mirrors Python DemandDataFormatter.format_demand_data.
func FormatDemandData(d *DemandData) map[string]interface{} {
	if d == nil {
		return map[string]interface{}{}
	}

	meters := map[string]float64{
		"p_import": d.DDP[0], "p_export": d.DDP[1],
		"q_import": d.DDQ[0], "q_export": d.DDQ[1],
		"s":   d.DDS,
		"i_a": d.DDI[0], "i_b": d.DDI[1], "i_c": d.DDI[2],
	}

	maxmin := map[string]interface{}{
		"p_import_max": d.MDP[0].Value, "p_import_maxTime": d.MDP[0].Time,
		"p_export_max": d.MDP[1].Value, "p_export_maxTime": d.MDP[1].Time,
		"q_import_max": d.MDQ[0].Value, "q_import_maxTime": d.MDQ[0].Time,
		"q_export_max": d.MDQ[1].Value, "q_export_maxTime": d.MDQ[1].Time,
		"s_max": d.MDS.Value, "s_maxTime": d.MDS.Time,
		"i_a_max": d.MDI[0].Value, "i_a_maxTime": d.MDI[0].Time,
		"i_b_max": d.MDI[1].Value, "i_b_maxTime": d.MDI[1].Time,
		"i_c_max": d.MDI[2].Value, "i_c_maxTime": d.MDI[2].Time,
	}
	for _, k := range []string{"p_import", "p_export", "q_import", "q_export", "s", "i_a", "i_b", "i_c"} {
		maxmin[k+"_min"] = 0.0
		maxmin[k+"_minTime"] = ""
	}

	powerKeys := []map[string]interface{}{
		{"id": 1, "key": "p_import", "label": "Import"},
		{"id": 2, "key": "p_export", "label": "Export"},
	}
	reactiveKeys := []map[string]interface{}{
		{"id": 1, "key": "q_import", "label": "Import"},
		{"id": 2, "key": "q_export", "label": "Export"},
	}
	apparentKeys := []map[string]interface{}{
		{"id": 1, "key": "s", "label": "Total"},
	}
	currentKeys := []map[string]interface{}{
		{"id": 1, "key": "i_a", "label": "L1"},
		{"id": 2, "key": "i_b", "label": "L2"},
		{"id": 3, "key": "i_c", "label": "L3"},
	}

	return map[string]interface{}{
		"power_demand":    getDataDict(meters, maxmin, powerKeys, "kW"),
		"reactive_demand": getDataDict(meters, maxmin, reactiveKeys, "kVar"),
		"apparent_demand": getDataDict(meters, maxmin, apparentKeys, "kVA"),
		"current_demand":  getDataDict(meters, maxmin, currentKeys, "A"),
		"predict_demand":  map[string]interface{}{"value": round2(d.PDP), "unit": "kW"},
		"dynamic_demand": map[string]interface{}{
			"timestamp": d.DDTime,
			"power":     d.DDP, "reactive": d.DDQ, "apparent": d.DDS,
			"current": d.DDI,
		},
	}
}

func getDataDict(meters map[string]float64, maxmin map[string]interface{}, keys []map[string]interface{}, unit string) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(keys))
	for _, item := range keys {
		key := item["key"].(string)
		result = append(result, map[string]interface{}{
			"id":       item["id"],
			"subTitle": item["label"],
			"value":    round2(meters[key]),
			"max":      round2(toF(maxmin[key+"_max"])),
			"maxTime":  maxmin[key+"_maxTime"],
			"min":      round2(toF(maxmin[key+"_min"])),
			"minTime":  maxmin[key+"_minTime"],
			"unit":     unit,
		})
	}
	return result
}

func round2(v float64) float64 {
	if !isFiniteFloat(v) {
		return 0
	}
	return math.Round(v*100) / 100
}

func isFiniteFloat(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}

func toF(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	}
	return 0
}
