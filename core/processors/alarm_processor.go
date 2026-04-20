package processors

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"sv500_core/data"
	"sv500_core/handlers"
)

// AlarmStatus constants.
const (
	AlarmStatusCleared  = 0 // CLEARED
	AlarmStatusOccurred = 1 // OCCURRED
)

// AlarmCondition constants.
const (
	AlarmConditionUnder = 0
	AlarmConditionOver  = 1
)

// AlarmRecord represents a single parsed alarm record.
// Format: '<HHHHIfHHf' = alarm_id, meter_id, status, count, ts, level, chan, cond, value
type AlarmRecord struct {
	AlarmID uint16  // alarm_id (chid, 0-31)
	MeterID uint16  // meter_id (0=Main, 1=Sub)
	Status  uint16  // 0=CLEARED, 1=OCCURRED
	Count   uint16  // count
	Ts      uint32  // timestamp
	Level   float32 // level
	Chan    uint16  // channel parameter index
	Cond    uint16  // 0=UNDER, 1=OVER
	Value   float32 // value
}

// AlarmRecordSize is the byte size of a single alarm record.
const AlarmRecordSize = 2 + 2 + 2 + 2 + 4 + 4 + 2 + 2 + 4 // = 24 bytes

// parameterOptions lists the 48 alarm parameter names (matching Python original).
var parameterOptions = []string{
	"None",
	"Temperature",
	"Frequency",
	"Phase Voltage L1",
	"Phase Voltage L2",
	"Phase Voltage L3",
	"Phase Voltage Average",
	"Phase Voltage L12",
	"Phase Voltage L23",
	"Phase Voltage L31",
	"Line Voltage Average",
	"Voltage Unbalance(Uo)",
	"Voltage Unbalance(Uu)",
	"Phase Current L1",
	"Phase Current L2",
	"Phase Current L3",
	"Phase Current Average",
	"Phase Current Total",
	"Phase Current Neutral",
	"Active Power L1",
	"Active Power L2",
	"Active Power L3",
	"Active Power Total",
	"Reactive Power L1",
	"Reactive Power L2",
	"Reactive Power L3",
	"Reactive Power Total",
	"D1",
	"D2",
	"D3",
	"D",
	"Apparent Power L1",
	"Apparent Power L2",
	"Apparent Power L3",
	"Apparent Power Total",
	"Power Factor L1",
	"Power Factor L2",
	"Power Factor L3",
	"Power Factor Total",
	"THD Voltage L1",
	"THD Voltage L2",
	"THD Voltage L3",
	"THD Voltage L12",
	"THD Voltage L23",
	"THD Voltage L31",
	"THD Current L1",
	"THD Current L2",
	"THD Current L3",
}

// AlarmStatusManager manages alarm status updates in Redis.
type AlarmStatusManager struct {
	redis *handlers.RedisHandler
}

// NewAlarmStatusManager creates a new AlarmStatusManager.
func NewAlarmStatusManager(redis *handlers.RedisHandler) *AlarmStatusManager {
	return &AlarmStatusManager{
		redis: redis,
	}
}

// UpdateAlarmStatus updates the alarm status in Redis using HSET.
// Key: "alarm_status:{channel}", Field: chid+1, Value: JSON of alarm data.
func (m *AlarmStatusManager) UpdateAlarmStatus(ctx context.Context, alarmData map[string]interface{}) error {
	channel := alarmData["channel_name"].(string)
	chid := alarmData["chid"].(int)

	key := fmt.Sprintf("alarm_status:%s", channel)

	ts := alarmData["timestamp"]
	var timestampInt int64
	switch v := ts.(type) {
	case time.Time:
		timestampInt = v.Unix()
	case int64:
		timestampInt = v
	case uint32:
		timestampInt = int64(v)
	}

	updateData := map[string]interface{}{
		"status":      alarmData["status"],
		"condition":   alarmData["condition_text"],
		"level":       alarmData["level"],
		"value":       alarmData["value"],
		"chan":         alarmData["chan"],
		"cond":        alarmData["condition"],
		"chan_text":    alarmData["chan_text"],
		"last_update":  timestampInt,
		"count":       alarmData["count"],
		"status_text": alarmData["status_text"],
	}

	jsonBytes, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	field := fmt.Sprintf("%d", chid+1)
	return m.redis.HSet(ctx, key, field, string(jsonBytes))
}

// CreateAlarmConfig returns a BinaryDataConfig for alarm data.
// Python format: '<HHHHIfHHf' = alarm_id, meter_id, status, count, ts, level, chan, cond, value
func CreateAlarmConfig() data.BinaryDataConfig {
	return data.BinaryDataConfig{
		Format: "<HHHHIfHHf",
		Fields: []string{
			"alarm_id", "meter_id", "status", "count",
			"ts", "level",
			"chan", "cond", "value",
		},
		Size: AlarmRecordSize,
	}
}

// AlarmDataProcessor extends GenericDataProcessor for alarm data.
type AlarmDataProcessor struct {
	*GenericDataProcessor
	statusManager *AlarmStatusManager

	mu      sync.RWMutex
	running bool
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

// NewAlarmDataProcessor creates a new AlarmDataProcessor.
func NewAlarmDataProcessor(redis *handlers.RedisHandler, influx *handlers.InfluxDBHandler) *AlarmDataProcessor {
	config := ProcessorConfig{
		DataConfig:      CreateAlarmConfig(),
		RedisKey:        "alarm_log",
		MeasurementName: "alarms",
		PollInterval:    1 * time.Second,
	}

	ap := &AlarmDataProcessor{
		GenericDataProcessor: NewGenericDataProcessor("AlarmProcessor", config, redis, influx),
		statusManager: NewAlarmStatusManager(redis),
	}

	ap.GenericDataProcessor.SetProcessFunc(ap.processAlarmSingleData)

	return ap
}

// parseAlarmRecord parses a single alarm record from binary data.
func parseAlarmRecord(buf []byte) (*AlarmRecord, error) {
	if len(buf) < AlarmRecordSize {
		return nil, fmt.Errorf("alarm record too short: %d < %d", len(buf), AlarmRecordSize)
	}

	r := &AlarmRecord{
		AlarmID: binary.LittleEndian.Uint16(buf[0:2]),
		MeterID: binary.LittleEndian.Uint16(buf[2:4]),
		Status:  binary.LittleEndian.Uint16(buf[4:6]),
		Count:   binary.LittleEndian.Uint16(buf[6:8]),
		Ts:      binary.LittleEndian.Uint32(buf[8:12]),
		Level:   math.Float32frombits(binary.LittleEndian.Uint32(buf[12:16])),
		Chan:    binary.LittleEndian.Uint16(buf[16:18]),
		Cond:    binary.LittleEndian.Uint16(buf[18:20]),
		Value:   math.Float32frombits(binary.LittleEndian.Uint32(buf[20:24])),
	}
	return r, nil
}

// alarmStatusText returns "OCCURRED" or "CLEARED" based on status value.
func alarmStatusText(status uint16) string {
	if status == AlarmStatusOccurred {
		return "OCCURRED"
	}
	return "CLEARED"
}

// alarmConditionText returns "UNDER" or "OVER" based on condition value.
func alarmConditionText(cond uint16) string {
	if cond == AlarmConditionUnder {
		return "UNDER"
	}
	return "OVER"
}

// extractAlarmData converts a parsed AlarmRecord into an alarm data map (matching Python _extract_alarm_data).
func extractAlarmData(record *AlarmRecord) map[string]interface{} {
	// meter_id: 0=Main, 1=Sub
	channelName := "Main"
	if record.MeterID != 0 {
		channelName = "Sub"
	}

	// Look up chan_text from parameterOptions
	chanText := "unknown"
	if int(record.Chan) < len(parameterOptions) {
		chanText = parameterOptions[record.Chan]
	}

	ts := time.Unix(int64(record.Ts), 0)

	return map[string]interface{}{
		"timestamp":      ts,
		"id":             int(record.MeterID),
		"channel_name":   channelName,
		"chid":           int(record.AlarmID),
		"chan":            int(record.Chan),
		"chan_text":       chanText,
		"status":         int(record.Status),
		"status_text":    alarmStatusText(record.Status),
		"condition":      int(record.Cond),
		"condition_text": alarmConditionText(record.Cond),
		"level":          float64(record.Level),
		"value":          float64(record.Value),
		"is_alarm":       record.Status == AlarmStatusOccurred,
		"count":          int(record.Count),
	}
}

// processAlarmSingleData processes a single alarm data payload.
func (ap *AlarmDataProcessor) processAlarmSingleData(rawData []byte) error {
	if len(rawData) < AlarmRecordSize {
		return fmt.Errorf("alarm data too short: %d bytes", len(rawData))
	}

	record, err := parseAlarmRecord(rawData)
	if err != nil {
		return fmt.Errorf("alarm parse error: %w", err)
	}

	// Extract alarm data (matching Python _extract_alarm_data).
	alarmData := extractAlarmData(record)

	ctx := context.Background()

	// Write to InfluxDB if status != 2.
	if alarmData["status"].(int) != 2 && ap.InfluxHandler != nil {
		ts := alarmData["timestamp"].(time.Time)
		fields := map[string]interface{}{
			"alarm_id":       alarmData["chid"],
			"meter_id":       alarmData["id"],
			"status":         alarmData["status"],
			"count":          alarmData["count"],
			"level":          alarmData["level"],
			"chan":            alarmData["chan"],
			"cond":           alarmData["condition"],
			"value":          alarmData["value"],
			"status_text":    alarmData["status_text"],
			"condition_text": alarmData["condition_text"],
			"chan_text":       alarmData["chan_text"],
		}
		tags := map[string]string{
			"channel": alarmData["channel_name"].(string),
		}
		if err := ap.InfluxHandler.Write(ap.Config.MeasurementName, tags, fields, ts); err != nil {
			log.Printf("[AlarmProcessor] influx write error: %v", err)
		}
	}

	// Update alarm status in Redis via HSET.
	if err := ap.statusManager.UpdateAlarmStatus(ctx, alarmData); err != nil {
		log.Printf("[AlarmProcessor] status update error: %v", err)
	}

	return nil
}

// Start begins the alarm data processor.
func (ap *AlarmDataProcessor) Start(ctx context.Context) error {
	ap.mu.Lock()
	if ap.running {
		ap.mu.Unlock()
		return fmt.Errorf("alarm processor already running")
	}
	ap.running = true
	ap.mu.Unlock()

	childCtx, cancel := context.WithCancel(ctx)
	ap.cancel = cancel

	ap.wg.Add(1)
	go func() {
		defer ap.wg.Done()
		ap.GenericDataProcessor.ProcessQueue(childCtx)
	}()

	log.Println("[AlarmProcessor] started")
	return nil
}

// Stop halts the alarm data processor.
func (ap *AlarmDataProcessor) Stop() {
	ap.mu.Lock()
	if !ap.running {
		ap.mu.Unlock()
		return
	}
	ap.running = false
	ap.mu.Unlock()

	if ap.cancel != nil {
		ap.cancel()
	}
	ap.wg.Wait()
	log.Println("[AlarmProcessor] stopped")
}
