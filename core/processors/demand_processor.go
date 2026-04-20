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

// DemandDataProcessor handles demand data processing with per-channel support.
type DemandDataProcessor struct {
	RedisHandler  *handlers.RedisHandler
	InfluxHandler *handlers.InfluxDBHandler
	Collector     *data.BinaryCollector
	Parser        *data.DemandParser
	Channels      []data.Channel
	PollInterval  time.Duration

	mu             sync.RWMutex
	running        bool
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	lastTimestamps map[string]uint32 // dedup: channel -> last ddTimestamp
}

// NewDemandDataProcessor creates a new DemandDataProcessor.
func NewDemandDataProcessor(redis *handlers.RedisHandler, influx *handlers.InfluxDBHandler, channels []data.Channel) *DemandDataProcessor {
	return &DemandDataProcessor{
		RedisHandler:   redis,
		InfluxHandler:  influx,
		Collector:      data.NewBinaryCollector(redis.Client),
		Parser:         data.NewDemandParser(),
		Channels:       channels,
		PollInterval:   1 * time.Second,
		lastTimestamps: make(map[string]uint32),
	}
}

// SetPollInterval updates the polling interval for the demand processor.
func (dp *DemandDataProcessor) SetPollInterval(interval time.Duration) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	dp.PollInterval = interval
}

// isDuplicate checks if a ddTimestamp has already been processed for a channel.
func (dp *DemandDataProcessor) isDuplicate(channelID string, ts uint32) bool {
	dp.mu.RLock()
	lastTS, exists := dp.lastTimestamps[channelID]
	dp.mu.RUnlock()

	if exists && ts == lastTS {
		return true
	}

	dp.mu.Lock()
	dp.lastTimestamps[channelID] = ts
	dp.mu.Unlock()

	return false
}

// processChannel processes demand data for a single channel.
// Python: reads from Redis hash "Demand" field=channel via BinaryCollector.get_demand()
func (dp *DemandDataProcessor) processChannel(ctx context.Context, channel data.Channel) {
	// Use BinaryCollector to read from Redis HGET "Demand" <channel>
	demandData, err := dp.Collector.GetDemand(channel.Name)
	if err != nil || demandData == nil {
		return
	}

	// Flatten the demand data for InfluxDB.
	flatData := data.FlattenDemand(demandData)

	// Deduplicate by ddTimestamp.
	ts := demandData.DDTimestamp
	if ts != 0 && dp.isDuplicate(channel.Name, ts) {
		return
	}

	// Write to InfluxDB.
	if dp.InfluxHandler != nil {
		fields := make(map[string]interface{})
		for k, v := range flatData {
			fields[k] = v
		}
		tags := map[string]string{
			"channel": channel.Name,
		}
		if err := dp.InfluxHandler.Write("demand", tags, fields, time.Now()); err != nil {
			log.Printf("[DemandProcessor] channel %s influx write error: %v", channel.Name, err)
		}
	}

	// Save summary to Redis
	if dp.RedisHandler != nil {
		summary := map[string]interface{}{
			"channel":     channel.Name,
			"ddTimestamp":  ts,
			"last_update":  time.Now().Format(time.RFC3339),
		}
		for k, v := range flatData {
			summary[k] = v
		}
		dp.RedisHandler.SaveSummary("demand_summary", channel.Name, summary)
	}
}

// workerLoop is the main processing loop for all demand channels.
func (dp *DemandDataProcessor) workerLoop(ctx context.Context) {
	ticker := time.NewTicker(dp.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for _, ch := range dp.Channels {
				dp.processChannel(ctx, ch)
			}
		}
	}
}

// Start begins the demand data processor.
func (dp *DemandDataProcessor) Start(ctx context.Context) error {
	dp.mu.Lock()
	if dp.running {
		dp.mu.Unlock()
		return fmt.Errorf("demand processor already running")
	}
	dp.running = true
	dp.mu.Unlock()

	childCtx, cancel := context.WithCancel(ctx)
	dp.cancel = cancel

	dp.wg.Add(1)
	go func() {
		defer dp.wg.Done()
		dp.workerLoop(childCtx)
	}()

	log.Println("[DemandProcessor] started")
	return nil
}

// Stop halts the demand data processor.
func (dp *DemandDataProcessor) Stop() {
	dp.mu.Lock()
	if !dp.running {
		dp.mu.Unlock()
		return
	}
	dp.running = false
	dp.mu.Unlock()

	if dp.cancel != nil {
		dp.cancel()
	}
	dp.wg.Wait()
	log.Println("[DemandProcessor] stopped")
}

// GetStatus returns the current processor status.
func (dp *DemandDataProcessor) GetStatus() ProcessorStatus {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	return ProcessorStatus{
		Name:    "DemandProcessor",
		Running: dp.running,
	}
}

// UpdateSettings updates demand processor settings at runtime.
func (dp *DemandDataProcessor) UpdateSettings(pollInterval time.Duration, channels []data.Channel) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	dp.PollInterval = pollInterval
	dp.Channels = channels
}
