package data

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/redis/go-redis/v9"

	"sv500_core/handlers"
)

// KEY_MAPPING maps parameter names to their Redis key lists.
var KEY_MAPPING = map[string][]string{
	"Phase Voltage": {"U1", "U2", "U3", "U4"},
	"Frequency":     {"Freq"},
	"Line Voltage":  {"Upp1", "Upp2", "Upp3", "Upp4"},
	"Current":       {"I1", "I2", "I3", "I4", "Itot", "In", "Ig"},
	"Unbalance":     {"Ubal1", "Ubal2", "Ibal1", "Ibal2"},
	"Power":         {"P1", "P2", "P3", "P4", "Q1", "Q2", "Q3", "Q4", "S1", "S2", "S3", "S4"},
	"PF":            {"PF1", "PF2", "PF3", "PF4"},
	"THD":           {"THD_U1", "THD_U2", "THD_U3", "THD_Upp1", "THD_Upp2", "THD_Upp3", "THD_I1", "THD_I2", "THD_I3"},
	"TDD":           {"TDD_I1", "TDD_I2", "TDD_I3"},
	"Energy":        {"total_kwh_import", "total_kvarh_import", "total_kvah_import", "total_kwh_export", "total_kvarh_export", "total_kvah_export"},
	"THDTDD_avg":    {"THD_V", "THD_I", "TDD_I"},
}

// PARAMETER_OPTIONS lists the available parameter option names.
var PARAMETER_OPTIONS = []string{
	"Freq.", "U1", "U2", "U3", "U~", "U12", "U23", "U31", "Upp~", "Uu", "Uo",
	"I1", "I2", "I3", "I~", "Itotal", "In", "Ig", "P1", "P2", "P3", "Ptotal", "Q1", "Q2", "Q3",
	"Qtotal", "D1", "D2", "D3", "Dtotal", "S1", "S2", "S3", "Stotal", "PF1", "PF2", "PF3",
	"PFtotal", "THD U1", "THD U2", "THD U3", "THD U12", "THD U23", "THD U31", "THD I1",
	"THD I2", "THD I3",
}

// Trend manages trend data collection for a channel.
type Trend struct {
	Channel     string
	AssetDrive  bool
	VoltageType int
	CollectList []string
}

// NewTrend creates a new Trend for the given channel.
func NewTrend(channel string, assetDrive bool, voltageType int) *Trend {
	return &Trend{
		Channel:     channel,
		AssetDrive:  assetDrive,
		VoltageType: voltageType,
		CollectList: make([]string, 0),
	}
}

// SetCollectKeys sets the list of Redis keys to collect based on the given parameter names.
func (t *Trend) SetCollectKeys(parameters []string) {
	t.CollectList = make([]string, 0)
	for _, param := range parameters {
		keys := t.GetKeyList(param)
		t.CollectList = append(t.CollectList, keys...)
	}
}

// GetKeyList returns the Redis keys for a parameter name.
func (t *Trend) GetKeyList(param string) []string {
	if keys, ok := KEY_MAPPING[param]; ok {
		return keys
	}
	return []string{param}
}

// CopyHash collects data from Redis hash matching the collectList keys.
// redisKey: source Redis hash key (empty = auto-detect by channel name)
// skipTotalFilter: if false, filters out keys containing "total"
func (t *Trend) CopyHash(collectList []string, redisInst *redis.Client, redisKey string, skipTotalFilter bool) map[string]string {
	ctx := context.Background()

	if redisKey == "" {
		if t.Channel == "Main" {
			redisKey = "meter_main"
		} else {
			redisKey = "meter_sub"
		}
	}

	// Filter keys containing "total" unless skipped
	var normalKeys []string
	if skipTotalFilter {
		normalKeys = make([]string, len(collectList))
		copy(normalKeys, collectList)
	} else {
		for _, key := range collectList {
			lower := strings.ToLower(key)
			if lower == "total" || strings.HasPrefix(lower, "total") || strings.Contains(lower, "total") {
				continue
			}
			normalKeys = append(normalKeys, key)
		}
	}

	dataDict := make(map[string]string)

	if len(normalKeys) > 0 {
		vals, err := redisInst.HMGet(ctx, redisKey, normalKeys...).Result()
		if err != nil {
			log.Printf("Redis data collection failed (%s): %v", t.Channel, err)
			return dataDict
		}

		for i, key := range normalKeys {
			if i < len(vals) && vals[i] != nil {
				dataDict[key] = fmt.Sprintf("%v", vals[i])
			}
		}
	}

	// THD/TDD average calculation
	if !t.AssetDrive {
		var thdVSources []string
		if t.VoltageType == 0 {
			thdVSources = []string{"THD_U1", "THD_U2", "THD_U3"}
		} else {
			thdVSources = []string{"THD_Upp1", "THD_Upp2", "THD_Upp3"}
		}

		avgSources := map[string][]string{
			"THD_V": thdVSources,
			"THD_I": {"THD_I1", "THD_I2", "THD_I3"},
			"TDD_I": {"TDD_I1", "TDD_I2", "TDD_I3"},
		}

		for avgKey, sourceKeys := range avgSources {
			var values []float64
			for _, k := range sourceKeys {
				if v, ok := dataDict[k]; ok {
					if fv, err := strconv.ParseFloat(v, 64); err == nil {
						values = append(values, fv)
					}
				}
			}
			if len(values) > 0 {
				sum := 0.0
				for _, v := range values {
					sum += v
				}
				avgVal := sum / float64(len(values))
				dataDict[avgKey] = fmt.Sprintf("%f", avgVal)
				redisInst.HSet(ctx, redisKey, avgKey, avgVal)
			}
		}
	} else {
		// VFD mode: read THD/TDD from Redis (stored by another thread)
		thdKeys := KEY_MAPPING["THD"]
		tddKeys := KEY_MAPPING["TDD"]

		hasThd := false
		hasTdd := false
		for _, k := range collectList {
			for _, tk := range thdKeys {
				if k == tk {
					hasThd = true
				}
			}
			for _, dk := range tddKeys {
				if k == dk {
					hasTdd = true
				}
			}
		}

		var thdtddKeys []string
		if hasThd {
			thdtddKeys = append(thdtddKeys, "THD_V", "THD_I")
		}
		if hasTdd {
			thdtddKeys = append(thdtddKeys, "TDD_I")
		}

		if len(thdtddKeys) > 0 {
			vals, err := redisInst.HMGet(ctx, redisKey, thdtddKeys...).Result()
			if err == nil {
				for i, k := range thdtddKeys {
					if i < len(vals) && vals[i] != nil {
						s := fmt.Sprintf("%v", vals[i])
						if fv, err := strconv.ParseFloat(s, 64); err == nil {
							if math.IsNaN(fv) {
								fv = 0.0
							}
							dataDict[k] = fmt.Sprintf("%f", fv)
						}
					}
				}
			}
		}
	}

	return dataDict
}

// SaveInflux saves the collected trend data to InfluxDB via connection pool.
func (t *Trend) SaveInflux(dataDict map[string]string, channel string, measurement string) {
	if len(dataDict) == 0 {
		log.Printf("No data to save (%s/%s), skipping", channel, measurement)
		return
	}

	pool, err := handlers.GetInfluxPool()
	if err != nil {
		log.Printf("InfluxDB pool error: %v", err)
		return
	}

	err = pool.WithConnection(func(entry *handlers.PoolEntry) error {
		point := influxdb2.NewPointWithMeasurement(measurement).
			AddTag("channel", channel)

		for k, v := range dataDict {
			if strings.Contains(k, "Time") || strings.Contains(k, "Timestamp") {
				if ts, err := ParseTimestamp(v); err == nil {
					point.AddField(k, float64(ts.Unix()))
				}
			} else {
				if fv, err := strconv.ParseFloat(v, 64); err == nil {
					if math.IsNaN(fv) {
						fv = 0.0
					}
					point.AddField(k, fv)
				}
			}
		}

		writeAPI := entry.Client.WriteAPIBlocking(pool.Org, "ntek")
		return writeAPI.WritePoint(context.Background(), point)
	})

	if err != nil {
		log.Printf("InfluxDB save failed (%s/%s): %v", channel, measurement, err)
	}
}

// ParseTimestamp parses a timestamp string in various common formats and returns a time.Time.
func ParseTimestamp(ts string) (time.Time, error) {
	ts = strings.TrimSpace(ts)

	if unixSec, err := strconv.ParseInt(ts, 10, 64); err == nil {
		return time.Unix(unixSec, 0), nil
	}

	if unixFloat, err := strconv.ParseFloat(ts, 64); err == nil {
		sec := int64(unixFloat)
		nsec := int64((unixFloat - float64(sec)) * 1e9)
		return time.Unix(sec, nsec), nil
	}

	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05.000",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, ts); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse timestamp: %s", ts)
}
