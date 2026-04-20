package processors

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"sv500_core/data"
	"sv500_core/handlers"
)

// PeriodType represents the time period for consumption calculation.
type PeriodType string

const (
	PeriodToday PeriodType = "today"
	PeriodWeek  PeriodType = "week"
	PeriodMonth PeriodType = "month"
	PeriodYear  PeriodType = "year"
)

// EnergyDataProcessor handles energy data processing with consumption calculations.
type EnergyDataProcessor struct {
	*GenericDataProcessor
	consumptionCalc *data.EnergyConsumptionCalculator
	periodCalc      *data.EnergyPeriodCalculator

	queueName   string
	summaryKey  string
	measurement string

	mu      sync.RWMutex
	running bool
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

// NewEnergyDataProcessor creates a new EnergyDataProcessor.
func NewEnergyDataProcessor(redis *handlers.RedisHandler, influx *handlers.InfluxDBHandler) *EnergyDataProcessor {
	energyConfig := data.CreateEnergyConfig()
	// Set channel mapping to match Python: {0: "Main", 1: "Sub"}
	energyConfig.ChannelMapping = map[int]string{0: "Main", 1: "Sub"}

	config := ProcessorConfig{
		DataConfig:        *energyConfig,
		RedisKey:          "energy_log",
		MeasurementName:   "energy",
		PollInterval:      30 * 60 * time.Second, // 30 minutes
		EnableConsumption: true,
	}

	ep := &EnergyDataProcessor{
		GenericDataProcessor: NewGenericDataProcessor("EnergyProcessor", config, redis, influx),
		consumptionCalc:      data.NewEnergyConsumptionCalculator(),
		periodCalc:           data.NewEnergyPeriodCalculator(15),
		queueName:            "energy_log",
		summaryKey:           "energy_summary",
		measurement:          "energy",
	}

	// Override the process function.
	ep.GenericDataProcessor.SetProcessFunc(ep.processEnergySingleData)

	// Initialize from Redis hash data (energy_main, energy_sub).
	ep.initializeFromRedis()

	return ep
}

// initializeFromRedis loads current cumulative energy data from Redis hashes
// for Main and Sub channels, and sets initial summary values.
func (ep *EnergyDataProcessor) initializeFromRedis() {
	log.Println("[EnergyProcessor] Loading initial energy data from Redis...")

	channelMap := map[string]string{
		"Main": "energy_main",
		"Sub":  "energy_sub",
	}

	for channel, redisKey := range channelMap {
		energyData, err := ep.getRedisEnergyData(redisKey)
		if err != nil {
			log.Printf("[EnergyProcessor] %s channel Redis data error: %v", channel, err)
			continue
		}
		if energyData == nil {
			log.Printf("[EnergyProcessor] %s channel: no Redis data found", channel)
			continue
		}

		ep.setInitialSummary(channel, energyData)
		log.Printf("[EnergyProcessor] %s channel initial data loaded", channel)
	}
}

// getRedisEnergyData reads energy hash data from Redis.
func (ep *EnergyDataProcessor) getRedisEnergyData(redisKey string) (map[string]float64, error) {
	ctx := context.Background()
	allData, err := ep.RedisHandler.HGetAll(ctx, redisKey)
	if err != nil {
		return nil, fmt.Errorf("HGetAll error for %s: %w", redisKey, err)
	}
	if len(allData) == 0 {
		return nil, nil
	}

	// Field mapping from Redis hash fields to internal names.
	fieldsMap := map[string]string{
		"total_kwh_import":   "kwh_import",
		"total_kwh_export":   "kwh_export",
		"total_kvarh_import": "kvarh_import",
		"total_kvarh_export": "kvarh_export",
		"total_kvah_import":  "kvah_import",
		"total_kvah_export":  "kvah_export",
	}

	energyData := make(map[string]float64)
	for field, value := range allData {
		if mappedName, ok := fieldsMap[field]; ok {
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				energyData[mappedName] = v
			}
		}
	}

	// Ensure all fields exist with default 0.
	for _, f := range []string{"kwh_import", "kwh_export", "kvarh_import", "kvarh_export", "kvah_import", "kvah_export"} {
		if _, ok := energyData[f]; !ok {
			energyData[f] = 0.0
		}
	}

	// Check if any value is > 0.
	hasData := false
	for _, v := range energyData {
		if v > 0 {
			hasData = true
			break
		}
	}
	if !hasData {
		return nil, nil
	}

	return energyData, nil
}

// setInitialSummary creates an initial summary entry in Redis from loaded data.
func (ep *EnergyDataProcessor) setInitialSummary(channel string, energyData map[string]float64) {
	now := time.Now()

	channelID := 0
	if channel == "Sub" {
		channelID = 1
	}

	// Calculate period consumption from InfluxDB.
	periodConsumption := ep.calculatePeriodConsumption(channel)

	summary := map[string]interface{}{
		"channel_id":   channelID,
		"channel_name": channel,
		"meter_id":     "meter_001",

		// Current cumulative values from Redis.
		"kwh_import":   energyData["kwh_import"],
		"kwh_export":   energyData["kwh_export"],
		"kvarh_import": energyData["kvarh_import"],
		"kvarh_export": energyData["kvarh_export"],
		"kvah_import":  energyData["kvah_import"],
		"kvah_export":  energyData["kvah_export"],

		// Consumption defaults (calculated on first data).
		"consumption": map[string]float64{
			"kwh_import_consumption":   0.0,
			"kwh_export_consumption":   0.0,
			"kvarh_import_consumption": 0.0,
			"kvarh_export_consumption": 0.0,
			"kvah_import_consumption":  0.0,
			"kvah_export_consumption":  0.0,
		},

		"last_update": now.Format(time.RFC3339),
		"calculation_info": map[string]interface{}{
			"calculation_type": "initial_load",
			"time_gap_minutes": 0,
			"is_valid":         true,
			"source":           "redis_initialization",
		},
	}

	// Merge period consumption data.
	for k, v := range periodConsumption {
		summary[k] = v
	}

	// Save to Redis.
	if err := ep.RedisHandler.SaveSummary(ep.summaryKey, channel, summary); err != nil {
		log.Printf("[EnergyProcessor] Failed to save initial summary for %s: %v", channel, err)
	}

	// Initialize consumption calculator cache with current values.
	ep.consumptionCalc.PreviousValues = energyData

	log.Printf("[EnergyProcessor] %s initial summary set - kwh_import: %.3f", channel, energyData["kwh_import"])
}

// processEnergySingleData processes a single energy data payload with consumption.
func (ep *EnergyDataProcessor) processEnergySingleData(rawData []byte) error {
	// 1. Parse binary data using StandardBinaryParser.
	parsed, err := ep.Parser.Parse(rawData)
	if err != nil {
		return fmt.Errorf("energy parse error: %w", err)
	}

	log.Printf("[EnergyProcessor] Parsed data - channel: %s, channel_id: %d",
		parsed.ChannelName, parsed.ChannelID)

	// 2. Query last data from InfluxDB for consumption calculation.
	var lastData map[string]interface{}
	if ep.InfluxHandler != nil {
		lastData, err = ep.InfluxHandler.QueryLastData(parsed.ChannelName, ep.measurement, 24)
		if err != nil {
			log.Printf("[EnergyProcessor] QueryLastData error: %v", err)
		}
	}

	// Fallback to cached previous values if InfluxDB has no data.
	if lastData == nil && len(ep.consumptionCalc.PreviousValues) > 0 {
		log.Printf("[EnergyProcessor] Using cached previous values for consumption calculation")
		lastData = make(map[string]interface{})
		for k, v := range ep.consumptionCalc.PreviousValues {
			lastData[k] = v
		}
	}

	// 3. Build data dict for writing.
	dataDict := make(map[string]interface{})
	for k, v := range parsed.Values {
		dataDict[k] = v
	}
	dataDict["channel_name"] = parsed.ChannelName
	dataDict["channel_id"] = parsed.ChannelID
	dataDict["meter_id"] = "meter_001"
	dataDict["timestamp"] = parsed.Timestamp

	// 4. Calculate consumption.
	var consumption map[string]float64
	currentVals := extractEnergyFloats(parsed.Values)
	if lastData != nil {
		previousVals := extractEnergyFloats(lastData)
		consumption = ep.consumptionCalc.CalculateConsumption(currentVals, previousVals)
	} else {
		log.Printf("[EnergyProcessor] %s first data - no previous data for consumption", parsed.ChannelName)
		ep.consumptionCalc.PreviousValues = currentVals
	}

	// 5. Write cumulative data to InfluxDB.
	if ep.InfluxHandler != nil {
		if err := ep.InfluxHandler.WriteCumulativeData(dataDict, ep.measurement); err != nil {
			return fmt.Errorf("influx write cumulative error: %w", err)
		}
	}

	// 6. Write consumption data to InfluxDB.
	if ep.InfluxHandler != nil && len(consumption) > 0 {
		dataDict["period"] = "interval"
		if err := ep.InfluxHandler.WriteConsumptionData(dataDict, consumption, ep.measurement+"_consumption"); err != nil {
			log.Printf("[EnergyProcessor] influx write consumption error: %v", err)
		}
	}

	// 7. Build summary data.
	summaryData := ep.createSummaryData(parsed, consumption)

	// 8. Calculate period consumption (today/week/month/year).
	periodData := ep.calculatePeriodConsumption(parsed.ChannelName)
	for k, v := range periodData {
		summaryData[k] = v
	}

	// 9. Save summary to Redis.
	if err := ep.RedisHandler.SaveSummary(ep.summaryKey, parsed.ChannelName, summaryData); err != nil {
		log.Printf("[EnergyProcessor] SaveSummary error: %v", err)
	}

	return nil
}

// extractEnergyFloats extracts energy field values as float64 from a map.
func extractEnergyFloats(vals map[string]interface{}) map[string]float64 {
	result := make(map[string]float64)
	energyFields := []string{"kwh_import", "kwh_export", "kvarh_import", "kvarh_export",
		"kvah_import", "kvah_export", "active_import", "active_export",
		"reactive_import", "reactive_export", "apparent_import", "apparent_export"}
	for _, f := range energyFields {
		if v, ok := vals[f]; ok {
			switch val := v.(type) {
			case float64:
				result[f] = val
			case float32:
				result[f] = float64(val)
			case int:
				result[f] = float64(val)
			}
		}
	}
	return result
}

// createSummaryData builds a summary map from parsed data and consumption.
func (ep *EnergyDataProcessor) createSummaryData(parsed *data.ParsedData, consumption map[string]float64) map[string]interface{} {
	summary := map[string]interface{}{
		"channel_id":   parsed.ChannelID,
		"channel_name": parsed.ChannelName,
		"meter_id":     "meter_001",
		"last_update":  parsed.Timestamp.Format(time.RFC3339),
	}

	// Copy energy values from parsed data.
	for k, v := range parsed.Values {
		summary[k] = v
	}

	// Add consumption.
	if consumption != nil {
		summary["consumption"] = consumption
	}

	return summary
}

// calculatePeriodConsumption calculates period consumption (today/week/month/year)
// by querying InfluxDB for consumption sums.
func (ep *EnergyDataProcessor) calculatePeriodConsumption(channel string) map[string]interface{} {
	result := make(map[string]interface{})

	if ep.InfluxHandler == nil {
		return result
	}

	now := time.Now()
	consumptionMeasurement := ep.measurement + "_consumption"

	periods := []struct {
		name  string
		start time.Time
	}{
		{"daily", time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())},
		{"weekly", func() time.Time {
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			startOfWeek := now.AddDate(0, 0, -(weekday - 1))
			return time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, now.Location())
		}()},
		{"monthly", time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())},
		{"yearly", time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())},
	}

	for _, p := range periods {
		periodSum, err := ep.InfluxHandler.QueryPeriodSum(channel, p.start, now, consumptionMeasurement)
		if err != nil {
			log.Printf("[EnergyProcessor] period %s query error: %v", p.name, err)
			continue
		}

		for field, value := range periodSum {
			// field is like "kwh_import_consumption" -> "daily_kwh_import"
			baseName := field
			if len(baseName) > len("_consumption") {
				baseName = baseName[:len(baseName)-len("_consumption")]
			}
			result[fmt.Sprintf("%s_%s", p.name, baseName)] = value
		}
	}

	return result
}

// ProcessQueue runs the energy data queue processor loop.
// Uses RedisHandler.GetBinaryData to RPOP from "energy_log" on DB 1.
func (ep *EnergyDataProcessor) ProcessQueue(ctx context.Context) {
	log.Printf("[EnergyProcessor] Queue processor started (poll interval: %v)", ep.Config.PollInterval)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Get binary data from Redis queue via RPOP.
			rawData, err := ep.RedisHandler.GetBinaryData(ep.queueName)
			if err != nil {
				log.Printf("[EnergyProcessor] GetBinaryData error: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			if rawData != nil && len(rawData) > 0 {
				log.Printf("[EnergyProcessor] Received data from queue - size: %d bytes", len(rawData))
				if err := ep.processEnergySingleData(rawData); err != nil {
					log.Printf("[EnergyProcessor] process error: %v", err)
				}
			} else {
				// No data available, wait for poll interval.
				select {
				case <-ctx.Done():
					return
				case <-time.After(ep.Config.PollInterval):
				}
			}
		}
	}
}

// Start begins the energy data processor.
func (ep *EnergyDataProcessor) Start(ctx context.Context) error {
	ep.mu.Lock()
	if ep.running {
		ep.mu.Unlock()
		return fmt.Errorf("energy processor already running")
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

	log.Println("[EnergyProcessor] started")
	return nil
}

// Stop halts the energy data processor.
func (ep *EnergyDataProcessor) Stop() {
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
	log.Println("[EnergyProcessor] stopped")
}
