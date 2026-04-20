package processors

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"sv500_core/data"
	"sv500_core/handlers"
)

// EventType constants matching Python: SAG=1, SWELL=2, Short_Interrupt=3,
// Long_Interrupt=4, Over_Current=5, Under_Current=6, Voltage_Transient=7,
// Current_Transient=8.
const (
	EventTypeSAG              = 1
	EventTypeSWELL            = 2
	EventTypeShortInterrupt   = 3
	EventTypeLongInterrupt    = 4
	EventTypeOverCurrent      = 5
	EventTypeUnderCurrent     = 6
	EventTypeVoltageTransient = 7
	EventTypeCurrentTransient = 8
)

// EventTypeNames maps event type codes to human-readable names.
var EventTypeNames = map[int]string{
	EventTypeSAG:              "SAG",
	EventTypeSWELL:            "SWELL",
	EventTypeShortInterrupt:   "SHORT INTERRUPT",
	EventTypeLongInterrupt:    "LONG INTERRUPT",
	EventTypeOverCurrent:      "OVER CURRENT",
	EventTypeUnderCurrent:     "UNDER CURRENT",
	EventTypeVoltageTransient: "VOLTAGE TRANSIENT",
	EventTypeCurrentTransient: "CURRENT TRANSIENT",
}

// EventTypeToText returns the text name for an event type code.
func EventTypeToText(eventType int) string {
	if name, ok := EventTypeNames[eventType]; ok {
		return name
	}
	return fmt.Sprintf("UNKNOWN_%d", eventType)
}

// Channel mapping: {0: "Main", 1: "Sub"}
var eventChannelMapping = map[int]string{
	0: "Main",
	1: "Sub",
}

// EventRecordSize is 28 bytes: int32 + uint32 + 4*uint16 + 3*float32
// Format: '<iIHHHHfff' = id(4) + startTs(4) + msec(2) + duration(2) + type(2) + mask(2) + level_0(4) + level_1(4) + level_2(4)
const EventRecordSize = 28

// EventDataProcessor processes event binary data from Redis queue.
type EventDataProcessor struct {
	*GenericDataProcessor

	queueName   string
	summaryKey  string
	measurement string

	mu      sync.RWMutex
	running bool
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

// NewEventDataProcessor creates a new EventDataProcessor.
func NewEventDataProcessor(redis *handlers.RedisHandler, influx *handlers.InfluxDBHandler) *EventDataProcessor {
	eventConfig := data.CreateEventConfig()
	// Override field names to match Python: '<iIHHHHfff'
	// Fields: id, startTs, msec, duration, type, mask, level_0, level_1, level_2
	eventConfig.FieldNames = []string{
		"channel_id", "timestamp", "msec", "duration", "type", "mask",
		"level_0", "level_1", "level_2",
	}
	eventConfig.ChannelMapping = eventChannelMapping
	eventConfig.StructSize = EventRecordSize
	eventConfig.Size = EventRecordSize

	config := ProcessorConfig{
		DataConfig:      *eventConfig,
		RedisKey:        "event_log",
		MeasurementName: "events",
		PollInterval:    10 * time.Second,
	}

	ep := &EventDataProcessor{
		GenericDataProcessor: NewGenericDataProcessor("EventProcessor", config, redis, influx),
		queueName:            "event_log",
		summaryKey:           "event_summary",
		measurement:          "events",
	}

	ep.GenericDataProcessor.SetProcessFunc(ep.processEventSingleData)

	return ep
}

// processEventSingleData processes a single event data payload.
func (ep *EventDataProcessor) processEventSingleData(rawData []byte) error {
	// 1. Parse binary data using StandardBinaryParser (handles '<iIHHHHfff').
	parsed, err := ep.Parser.Parse(rawData)
	if err != nil {
		return fmt.Errorf("event parse error: %w", err)
	}

	// 2. Extract event fields from parsed data.
	eventData := ep.extractEventData(parsed)

	// 3. Write event to InfluxDB.
	if ep.InfluxHandler != nil {
		if err := ep.InfluxHandler.WriteEventData(eventData, ep.measurement); err != nil {
			return fmt.Errorf("influx write event error: %w", err)
		}
	}

	// 4. Save event summary to Redis.
	ep.saveEventSummary(eventData)

	log.Printf("[EventProcessor] Event processed: id=%v type=%s channel=%s duration=%v phases=%v",
		eventData["event_id"], eventData["event_type_text"], eventData["channel_name"],
		eventData["duration"], eventData["affected_phases"])

	return nil
}

// extractEventData extracts event information from parsed binary data.
func (ep *EventDataProcessor) extractEventData(parsed *data.ParsedData) map[string]interface{} {
	raw := parsed.RawValues

	eventID := toInt(raw["channel_id"])
	eventType := toInt(raw["type"])
	msec := toInt(raw["msec"])
	duration := toInt(raw["duration"])
	mask := toInt(raw["mask"])

	// Channel determined by parser from channel_id field via ChannelMapping.
	channelName := parsed.ChannelName

	// Calculate exact timestamp: startTs + milliseconds offset.
	exactTimestamp := parsed.Timestamp.Add(time.Duration(msec) * time.Millisecond)

	// Decode phase mask: bit 0=L1, bit 1=L2, bit 2=L3.
	affectedPhases := DecodePhaseMask(uint16(mask))

	return map[string]interface{}{
		"timestamp":       exactTimestamp,
		"event_id":        eventID,
		"channel_name":    channelName,
		"milliseconds":    msec,
		"duration":        duration,
		"event_type":      eventType,
		"event_type_text": EventTypeToText(eventType),
		"mask":            mask,
		"level_l1":        toFloat64(raw["level_0"]),
		"level_l2":        toFloat64(raw["level_1"]),
		"level_l3":        toFloat64(raw["level_2"]),
		"affected_phases": affectedPhases,
	}
}

// saveEventSummary saves event summary data to Redis.
func (ep *EventDataProcessor) saveEventSummary(eventData map[string]interface{}) {
	channelName := fmt.Sprintf("%v", eventData["channel_name"])

	ts, _ := eventData["timestamp"].(time.Time)

	eventSummary := map[string]interface{}{
		"timestamp":  ts.Format(time.RFC3339Nano),
		"event_id":   eventData["event_id"],
		"event_type": eventData["event_type_text"],
		"duration":   eventData["duration"],
		"levels": map[string]interface{}{
			"L1": eventData["level_l1"],
			"L2": eventData["level_l2"],
			"L3": eventData["level_l3"],
		},
		"affected_phases": eventData["affected_phases"],
	}

	if err := ep.RedisHandler.SaveSummary(ep.summaryKey, channelName, eventSummary); err != nil {
		log.Printf("[EventProcessor] SaveSummary error for %s: %v", channelName, err)
	}
}

// DecodePhaseMask decodes a phase bitmask into phase names.
// bit 0 = L1, bit 1 = L2, bit 2 = L3 (matching Python).
func DecodePhaseMask(mask uint16) []string {
	var phases []string
	if mask&0x01 != 0 {
		phases = append(phases, "L1")
	}
	if mask&0x02 != 0 {
		phases = append(phases, "L2")
	}
	if mask&0x04 != 0 {
		phases = append(phases, "L3")
	}
	if len(phases) == 0 {
		phases = append(phases, "None")
	}
	return phases
}

// ProcessQueue runs the event data queue processor loop.
// Uses RedisHandler.GetBinaryData to RPOP from "event_log" on DB 1.
func (ep *EventDataProcessor) ProcessQueue(ctx context.Context) {
	log.Printf("[EventProcessor] Queue processor started (poll interval: %v)", ep.Config.PollInterval)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			rawData, err := ep.RedisHandler.GetBinaryData(ep.queueName)
			if err != nil {
				log.Printf("[EventProcessor] GetBinaryData error: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			if rawData != nil && len(rawData) > 0 {
				log.Printf("[EventProcessor] Received event data - size: %d bytes", len(rawData))
				if err := ep.processEventSingleData(rawData); err != nil {
					log.Printf("[EventProcessor] process error: %v", err)
				}
			} else {
				select {
				case <-ctx.Done():
					return
				case <-time.After(ep.Config.PollInterval):
				}
			}
		}
	}
}

// Start begins the event data processor.
func (ep *EventDataProcessor) Start(ctx context.Context) error {
	ep.mu.Lock()
	if ep.running {
		ep.mu.Unlock()
		return fmt.Errorf("event processor already running")
	}
	ep.running = true
	ep.mu.Unlock()

	childCtx, cancel := context.WithCancel(ctx)
	ep.cancel = cancel

	ep.wg.Add(1)
	go func() {
		defer ep.wg.Done()
		ep.ProcessQueue(childCtx)
	}()

	log.Println("[EventProcessor] started")
	return nil
}

// Stop halts the event data processor.
func (ep *EventDataProcessor) Stop() {
	ep.mu.Lock()
	if !ep.running {
		ep.mu.Unlock()
		return
	}
	ep.running = false
	ep.mu.Unlock()

	if ep.cancel != nil {
		ep.cancel()
	}
	ep.wg.Wait()
	log.Println("[EventProcessor] stopped")
}

// Helper functions for type conversion.

func toInt(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case int32:
		return int(val)
	case int64:
		return int(val)
	case uint16:
		return int(val)
	case uint32:
		return int(val)
	case float32:
		return int(val)
	case float64:
		return int(val)
	default:
		return 0
	}
}

func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int32:
		return float64(val)
	case uint16:
		return float64(val)
	default:
		return 0.0
	}
}
