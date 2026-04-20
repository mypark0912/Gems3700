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

// DataProcessorManager manages all data processors.
type DataProcessorManager struct {
	RedisHandler  *handlers.RedisHandler
	InfluxHandler *handlers.InfluxDBHandler

	energyProcessor    *EnergyDataProcessor
	alarmProcessor     *AlarmDataProcessor
	eventProcessor     *EventDataProcessor
	en50160Processor   *EN50160DataProcessor
	pq10MinProcessor   *PQ10MinDataProcessor
	demandProcessor    *DemandDataProcessor
	diagnosisProcessor *DiagnosisProcessor

	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

// NewDataProcessorManager creates a new DataProcessorManager.
func NewDataProcessorManager(redis *handlers.RedisHandler, influx *handlers.InfluxDBHandler) *DataProcessorManager {
	return &DataProcessorManager{
		RedisHandler:  redis,
		InfluxHandler: influx,
	}
}

// Initialize creates and configures all processors.
func (m *DataProcessorManager) Initialize(channels []data.Channel) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.Println("[Manager] initializing processors...")

	// Energy processor.
	m.energyProcessor = NewEnergyDataProcessor(m.RedisHandler, m.InfluxHandler)

	// Alarm processor.
	m.alarmProcessor = NewAlarmDataProcessor(m.RedisHandler, m.InfluxHandler)

	// Event processor.
	m.eventProcessor = NewEventDataProcessor(m.RedisHandler, m.InfluxHandler)

	// EN50160 processor.
	m.en50160Processor = NewEN50160DataProcessor(m.RedisHandler, m.InfluxHandler)

	// PQ 10-minute data processor.
	m.pq10MinProcessor = NewPQ10MinDataProcessor(m.RedisHandler, m.InfluxHandler)

	// Demand processor.
	if channels == nil {
		channels = []data.Channel{}
	}
	m.demandProcessor = NewDemandDataProcessor(m.RedisHandler, m.InfluxHandler, channels)

	// Diagnosis processor.
	m.diagnosisProcessor = NewDiagnosisProcessor(m.InfluxHandler)

	log.Println("[Manager] all processors initialized")
	return nil
}

// StartAll starts all processors.
func (m *DataProcessorManager) StartAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	m.ctx = ctx
	m.cancel = cancel

	log.Println("[Manager] starting all processors...")

	if m.energyProcessor != nil {
		if err := m.energyProcessor.Start(ctx); err != nil {
			log.Printf("[Manager] energy processor start error: %v", err)
		}
	}

	if m.alarmProcessor != nil {
		if err := m.alarmProcessor.Start(ctx); err != nil {
			log.Printf("[Manager] alarm processor start error: %v", err)
		}
	}

	if m.eventProcessor != nil {
		if err := m.eventProcessor.Start(ctx); err != nil {
			log.Printf("[Manager] event processor start error: %v", err)
		}
	}

	if m.en50160Processor != nil {
		if err := m.en50160Processor.Start(ctx); err != nil {
			log.Printf("[Manager] EN50160 processor start error: %v", err)
		}
	}

	if m.pq10MinProcessor != nil {
		if err := m.pq10MinProcessor.Start(ctx); err != nil {
			log.Printf("[Manager] PQ10Min processor start error: %v", err)
		}
	}

	if m.demandProcessor != nil {
		if err := m.demandProcessor.Start(ctx); err != nil {
			log.Printf("[Manager] demand processor start error: %v", err)
		}
	}

	if m.diagnosisProcessor != nil {
		m.diagnosisProcessor.Start(ctx)
	}

	log.Println("[Manager] all processors started")
	return nil
}

// StopAll stops all processors.
func (m *DataProcessorManager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.Println("[Manager] stopping all processors...")

	if m.cancel != nil {
		m.cancel()
	}

	if m.energyProcessor != nil {
		m.energyProcessor.Stop()
	}
	if m.alarmProcessor != nil {
		m.alarmProcessor.Stop()
	}
	if m.eventProcessor != nil {
		m.eventProcessor.Stop()
	}
	if m.en50160Processor != nil {
		m.en50160Processor.Stop()
	}
	if m.pq10MinProcessor != nil {
		m.pq10MinProcessor.Stop()
	}
	if m.demandProcessor != nil {
		m.demandProcessor.Stop()
	}
	if m.diagnosisProcessor != nil {
		m.diagnosisProcessor.Stop()
	}

	log.Println("[Manager] all processors stopped")
}

// GetStatus returns the status of all processors.
func (m *DataProcessorManager) GetStatus() map[string]ProcessorStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	statuses := make(map[string]ProcessorStatus)

	if m.energyProcessor != nil {
		statuses["energy"] = m.energyProcessor.GenericDataProcessor.GetStatus()
	}
	if m.alarmProcessor != nil {
		statuses["alarm"] = m.alarmProcessor.GenericDataProcessor.GetStatus()
	}
	if m.eventProcessor != nil {
		statuses["event"] = m.eventProcessor.GenericDataProcessor.GetStatus()
	}
	if m.en50160Processor != nil {
		statuses["en50160"] = m.en50160Processor.GenericDataProcessor.GetStatus()
	}
	if m.pq10MinProcessor != nil {
		statuses["pq_10min"] = m.pq10MinProcessor.GenericDataProcessor.GetStatus()
	}
	if m.demandProcessor != nil {
		statuses["demand"] = m.demandProcessor.GetStatus()
	}

	return statuses
}

// UpdateDemandSettings updates demand processor settings at runtime.
func (m *DataProcessorManager) UpdateDemandSettings(pollInterval time.Duration, channels []data.Channel) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.demandProcessor == nil {
		return fmt.Errorf("demand processor not initialized")
	}

	m.demandProcessor.UpdateSettings(pollInterval, channels)
	log.Printf("[Manager] demand settings updated: interval=%v, channels=%d", pollInterval, len(channels))
	return nil
}

// GetEnergyProcessor returns the energy processor instance.
func (m *DataProcessorManager) GetEnergyProcessor() *EnergyDataProcessor {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.energyProcessor
}

// GetAlarmProcessor returns the alarm processor instance.
func (m *DataProcessorManager) GetAlarmProcessor() *AlarmDataProcessor {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.alarmProcessor
}

// GetDemandProcessor returns the demand processor instance.
func (m *DataProcessorManager) GetDemandProcessor() *DemandDataProcessor {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.demandProcessor
}
