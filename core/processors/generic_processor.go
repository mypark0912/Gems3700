package processors

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"sv500_core/data"
	"sv500_core/handlers"
)

// ProcessorConfig holds configuration for a generic data processor.
type ProcessorConfig struct {
	DataConfig        data.BinaryDataConfig
	RedisKey          string
	MeasurementName   string
	PollInterval      time.Duration
	EnableConsumption bool
}

// ProcessorStatus represents the running state of a processor.
type ProcessorStatus struct {
	Name      string `json:"name"`
	Running   bool   `json:"running"`
	Processed int64  `json:"processed"`
	Errors    int64  `json:"errors"`
}

// DataProcessor is the interface that all processors implement.
type DataProcessor interface {
	ProcessSingleData(rawData []byte) error
	Start(ctx context.Context) error
	Stop()
	GetStatus() ProcessorStatus
}

// GenericDataProcessor provides base functionality for queue-based binary data processing.
type GenericDataProcessor struct {
	Config        ProcessorConfig
	RedisHandler  *handlers.RedisHandler
	InfluxHandler *handlers.InfluxDBHandler
	Parser        *data.StandardBinaryParser
	Name          string

	mu        sync.RWMutex
	running   bool
	processed int64
	errors    int64
	cancel    context.CancelFunc
	wg        sync.WaitGroup

	// processors is a map of named sub-processor functions that run as goroutines.
	processors map[string]func(ctx context.Context)

	// processFunc allows derived processors to override ProcessSingleData behavior.
	processFunc func(rawData []byte) error
}

// NewGenericDataProcessor creates a new GenericDataProcessor.
func NewGenericDataProcessor(name string, config ProcessorConfig, redis *handlers.RedisHandler, influx *handlers.InfluxDBHandler) *GenericDataProcessor {
	p := &GenericDataProcessor{
		Config:        config,
		RedisHandler:  redis,
		InfluxHandler: influx,
		Parser:        data.NewStandardBinaryParser(config.DataConfig),
		Name:          name,
		processors:    make(map[string]func(ctx context.Context)),
	}
	p.processFunc = p.defaultProcessSingleData
	return p
}

// SetProcessFunc allows derived processors to override the processing function.
func (p *GenericDataProcessor) SetProcessFunc(fn func(rawData []byte) error) {
	p.processFunc = fn
}

// RegisterProcessor adds a named sub-processor goroutine.
func (p *GenericDataProcessor) RegisterProcessor(name string, fn func(ctx context.Context)) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.processors[name] = fn
}

// ProcessSingleData parses and processes a single binary data payload.
func (p *GenericDataProcessor) ProcessSingleData(rawData []byte) error {
	return p.processFunc(rawData)
}

func (p *GenericDataProcessor) defaultProcessSingleData(rawData []byte) error {
	parsed, err := p.Parser.Parse(rawData)
	if err != nil {
		p.mu.Lock()
		p.errors++
		p.mu.Unlock()
		return fmt.Errorf("parse error: %w", err)
	}

	if p.InfluxHandler != nil {
		fields := make(map[string]interface{})
		for k, v := range parsed.Values {
			fields[k] = v
		}
		tags := map[string]string{
			"source": p.Name,
		}
		if err := p.InfluxHandler.Write(p.Config.MeasurementName, tags, fields, parsed.Timestamp); err != nil {
			p.mu.Lock()
			p.errors++
			p.mu.Unlock()
			return fmt.Errorf("influx write error: %w", err)
		}
	}

	p.mu.Lock()
	p.processed++
	p.mu.Unlock()
	return nil
}

// ProcessQueue continuously pops data from the Redis queue and processes it.
func (p *GenericDataProcessor) ProcessQueue(ctx context.Context) {
	ticker := time.NewTicker(p.Config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			rawData, err := p.RedisHandler.LPop(ctx, p.Config.RedisKey)
			if err != nil {
				continue
			}
			if len(rawData) == 0 {
				continue
			}
			if err := p.ProcessSingleData([]byte(rawData)); err != nil {
				log.Printf("[%s] process error: %v", p.Name, err)
			}
		}
	}
}

// StartProcessor starts the queue processor goroutine.
func (p *GenericDataProcessor) StartProcessor(ctx context.Context) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.ProcessQueue(ctx)
	}()
}

// StartAll starts the queue processor and all registered sub-processors.
func (p *GenericDataProcessor) StartAll() context.Context {
	p.mu.Lock()
	if p.running {
		p.mu.Unlock()
		return nil
	}
	p.running = true
	p.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	p.StartProcessor(ctx)

	p.mu.RLock()
	for name, fn := range p.processors {
		processorFn := fn
		processorName := name
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			log.Printf("[%s] sub-processor %s started", p.Name, processorName)
			processorFn(ctx)
			log.Printf("[%s] sub-processor %s stopped", p.Name, processorName)
		}()
	}
	p.mu.RUnlock()

	log.Printf("[%s] processor started", p.Name)
	return ctx
}

// Stop cancels all goroutines and waits for them to finish.
func (p *GenericDataProcessor) Stop() {
	p.mu.Lock()
	if !p.running {
		p.mu.Unlock()
		return
	}
	p.running = false
	p.mu.Unlock()

	if p.cancel != nil {
		p.cancel()
	}
	p.wg.Wait()
	log.Printf("[%s] processor stopped", p.Name)
}

// GetStatus returns the current processor status.
func (p *GenericDataProcessor) GetStatus() ProcessorStatus {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return ProcessorStatus{
		Name:      p.Name,
		Running:   p.running,
		Processed: p.processed,
		Errors:    p.errors,
	}
}

// GetStatusJSON returns the status as a JSON string.
func (p *GenericDataProcessor) GetStatusJSON() string {
	status := p.GetStatus()
	b, _ := json.Marshal(status)
	return string(b)
}
