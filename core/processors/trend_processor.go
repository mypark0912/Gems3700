package processors

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"sv500_core/data"
	"sv500_core/handlers"
)

// KeyMapping maps parameter names to their Redis hash keys.
var KeyMapping = map[string]string{
	"V_AN":        "voltage:an",
	"V_BN":        "voltage:bn",
	"V_CN":        "voltage:cn",
	"V_AB":        "voltage:ab",
	"V_BC":        "voltage:bc",
	"V_CA":        "voltage:ca",
	"I_A":         "current:a",
	"I_B":         "current:b",
	"I_C":         "current:c",
	"I_N":         "current:n",
	"kW_Total":    "power:kw_total",
	"kW_A":        "power:kw_a",
	"kW_B":        "power:kw_b",
	"kW_C":        "power:kw_c",
	"kvar_Total":  "power:kvar_total",
	"kvar_A":      "power:kvar_a",
	"kvar_B":      "power:kvar_b",
	"kvar_C":      "power:kvar_c",
	"kVA_Total":   "power:kva_total",
	"kVA_A":       "power:kva_a",
	"kVA_B":       "power:kva_b",
	"kVA_C":       "power:kva_c",
	"PF_Total":    "pf:total",
	"PF_A":        "pf:a",
	"PF_B":        "pf:b",
	"PF_C":        "pf:c",
	"Freq":        "frequency",
	"THD_V_AN":    "thd:v_an",
	"THD_V_BN":    "thd:v_bn",
	"THD_V_CN":    "thd:v_cn",
	"THD_I_A":     "thd:i_a",
	"THD_I_B":     "thd:i_b",
	"THD_I_C":     "thd:i_c",
	"Unbal_V":     "unbalance:v",
	"Unbal_I":     "unbalance:i",
	"kWh_Imp":     "energy:kwh_imp",
	"kWh_Exp":     "energy:kwh_exp",
	"kvarh_Imp":   "energy:kvarh_imp",
	"kvarh_Exp":   "energy:kvarh_exp",
	"Temp":        "temperature",
}

// ParseTimestamp parses a timestamp string into a time.Time.
// Supports Unix epoch (integer or float) and RFC3339 formats.
func ParseTimestamp(ts string) (time.Time, error) {
	ts = strings.TrimSpace(ts)

	// Try Unix epoch (integer).
	if epochSec, err := strconv.ParseInt(ts, 10, 64); err == nil {
		return time.Unix(epochSec, 0), nil
	}

	// Try Unix epoch (float).
	if epochFloat, err := strconv.ParseFloat(ts, 64); err == nil {
		sec := int64(epochFloat)
		nsec := int64((epochFloat - float64(sec)) * 1e9)
		return time.Unix(sec, nsec), nil
	}

	// Try RFC3339.
	if t, err := time.Parse(time.RFC3339, ts); err == nil {
		return t, nil
	}

	// Try common datetime format.
	if t, err := time.Parse("2006-01-02 15:04:05", ts); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("unable to parse timestamp: %s", ts)
}

// SaveTrendToInflux writes trend data to InfluxDB using the Trend's SaveInflux method.
func SaveTrendToInflux(trend *data.Trend, dataDict map[string]string, channel string, measurement string) {
	trend.SaveInflux(dataDict, channel, measurement)
}

// CopyHash copies all fields from one Redis hash to another.
func CopyHash(ctx context.Context, redis *handlers.RedisHandler, srcKey, dstKey string) error {
	vals, err := redis.HGetAll(ctx, srcKey)
	if err != nil {
		return fmt.Errorf("failed to read source hash %s: %w", srcKey, err)
	}

	if len(vals) == 0 {
		return nil
	}

	// Convert map[string]string to map[string]interface{} for HMSet.
	fields := make(map[string]interface{}, len(vals))
	for k, v := range vals {
		fields[k] = v
	}

	return redis.HMSet(ctx, dstKey, fields)
}

// TrendSaver periodically collects trend data from Redis and saves to InfluxDB.
type TrendSaver struct {
	trend       *data.Trend
	redisClient *redis.Client
	redisKey    string
	measurement string
	interval    time.Duration

	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewTrendSaver creates a new TrendSaver.
func NewTrendSaver(trend *data.Trend, redisClient *redis.Client, redisKey string, measurement string, interval time.Duration) *TrendSaver {
	return &TrendSaver{
		trend:       trend,
		redisClient: redisClient,
		redisKey:    redisKey,
		measurement: measurement,
		interval:    interval,
	}
}

// Start begins periodic trend collection and saving.
func (ts *TrendSaver) Start(ctx context.Context) {
	childCtx, cancel := context.WithCancel(ctx)
	ts.cancel = cancel

	ts.wg.Add(1)
	go func() {
		defer ts.wg.Done()
		ticker := time.NewTicker(ts.interval)
		defer ticker.Stop()

		for {
			select {
			case <-childCtx.Done():
				return
			case <-ticker.C:
				dataDict := ts.trend.CopyHash(ts.trend.CollectList, ts.redisClient, ts.redisKey, false)
				ts.trend.SaveInflux(dataDict, ts.trend.Channel, ts.measurement)
			}
		}
	}()
}

// Stop halts the trend saver.
func (ts *TrendSaver) Stop() {
	if ts.cancel != nil {
		ts.cancel()
	}
	ts.wg.Wait()
}
