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

	"sv500_core/handlers"
)

// PQ10MinRecordSize is the total binary struct size in bytes.
// Layout (little-endian, packed):
//
//	i(4) + IIII(16) + III(12) + 60H(120) + 60f(240)
//	+ 3f(12) + 3f(12) + 2f(8) + 72f(288) + 3f(12) + 3f(12) + 3f(12)
//	+ 20H(40)
//	+ 3f(12) + f(4) + 2f(8) + 5f(20) + ff(8) + 2f2ff(20) + 6Q(48)
//	+ 20I(80)
//	+ If*2(16) + If*6(48) + If(8) + If*3(24) + If*3(24) + If*3(24) + If*3(24) + If*3(24)
//	= 1180 bytes
const PQ10MinRecordSize = 1180

// pq10minChannelMapping maps channel IDs to names.
var pq10minChannelMapping = map[int32]string{
	0: "Main",
	1: "Sub",
}

// pq10minEventTypes lists the five event types in struct order.
var pq10minEventTypes = []string{"sag", "swell", "short_intr", "long_intr", "rvc"}

// ---------------------------------------------------------------------------
// Raw parsed struct
// ---------------------------------------------------------------------------

type pq10minRaw struct {
	// Header
	ID       int32
	Ts10m    uint32
	Ts10s    uint32
	Count10m uint32
	Count10s uint32

	// QualAvgData
	StartTs    uint32
	EndTs      uint32
	AvgCount   uint32
	Ts10       [60]uint16
	Freq       [60]float32
	U          [3]float32
	Uthd       [3]float32
	Ubal       [2]float32
	Uhd        [3][24]float32 // 3 phases x 24 harmonics
	Pst        [3]float32
	Plt        [3]float32
	Svolt      [3]float32
	AvgEvents  [5][4]uint16 // sag, swell, short_intr, long_intr, rvc

	// QualAvgExpData
	I      [3]float32
	In     float32
	Ibal   [2]float32
	Temp   [5]float32
	PF     float32
	DPF    float32
	P      [2]float32
	Q      [2]float32
	S      float32
	Energy [3][2]uint64 // eh[3][2]

	// QualVarData
	VarEvents  [5][4]uint32 // sag, swell, short_intr, long_intr, rvc
	FreqVar    [2]struct{ Err uint32; Val float32 }
	VoltVar    [6]struct{ Err uint32; Val float32 } // volt1[3], volt2[3]
	VoltbalVar struct{ Err uint32; Val float32 }
	VoltThdVar [3]struct{ Err uint32; Val float32 }
	VoltHdVar  [3]struct{ Err uint32; Val float32 }
	PstVar     [3]struct{ Err uint32; Val float32 }
	PltVar     [3]struct{ Err uint32; Val float32 }
	SvoltVar   [3]struct{ Err uint32; Val float32 }
}

// ---------------------------------------------------------------------------
// Binary parser
// ---------------------------------------------------------------------------

func readInt32(buf []byte, off *int) int32 {
	v := int32(binary.LittleEndian.Uint32(buf[*off : *off+4]))
	*off += 4
	return v
}

func readUint32(buf []byte, off *int) uint32 {
	v := binary.LittleEndian.Uint32(buf[*off : *off+4])
	*off += 4
	return v
}

func readUint16(buf []byte, off *int) uint16 {
	v := binary.LittleEndian.Uint16(buf[*off : *off+2])
	*off += 2
	return v
}

func readFloat32(buf []byte, off *int) float32 {
	v := math.Float32frombits(binary.LittleEndian.Uint32(buf[*off : *off+4]))
	*off += 4
	return v
}

func readUint64(buf []byte, off *int) uint64 {
	v := binary.LittleEndian.Uint64(buf[*off : *off+8])
	*off += 8
	return v
}

// readQualVariation reads one (uint32, float32) pair.
func readQualVariation(buf []byte, off *int) struct{ Err uint32; Val float32 } {
	return struct{ Err uint32; Val float32 }{
		Err: readUint32(buf, off),
		Val: readFloat32(buf, off),
	}
}

func parsePQ10MinBinary(buf []byte) (*pq10minRaw, error) {
	if len(buf) < PQ10MinRecordSize {
		return nil, fmt.Errorf("PQ10Min data too short: got %d bytes, need %d", len(buf), PQ10MinRecordSize)
	}

	r := &pq10minRaw{}
	off := 0

	// Header
	r.ID = readInt32(buf, &off)
	r.Ts10m = readUint32(buf, &off)
	r.Ts10s = readUint32(buf, &off)
	r.Count10m = readUint32(buf, &off)
	r.Count10s = readUint32(buf, &off)

	// QualAvgData
	r.StartTs = readUint32(buf, &off)
	r.EndTs = readUint32(buf, &off)
	r.AvgCount = readUint32(buf, &off)

	for i := 0; i < 60; i++ {
		r.Ts10[i] = readUint16(buf, &off)
	}
	for i := 0; i < 60; i++ {
		r.Freq[i] = readFloat32(buf, &off)
	}
	for i := 0; i < 3; i++ {
		r.U[i] = readFloat32(buf, &off)
	}
	for i := 0; i < 3; i++ {
		r.Uthd[i] = readFloat32(buf, &off)
	}
	for i := 0; i < 2; i++ {
		r.Ubal[i] = readFloat32(buf, &off)
	}
	for phase := 0; phase < 3; phase++ {
		for h := 0; h < 24; h++ {
			r.Uhd[phase][h] = readFloat32(buf, &off)
		}
	}
	for i := 0; i < 3; i++ {
		r.Pst[i] = readFloat32(buf, &off)
	}
	for i := 0; i < 3; i++ {
		r.Plt[i] = readFloat32(buf, &off)
	}
	for i := 0; i < 3; i++ {
		r.Svolt[i] = readFloat32(buf, &off)
	}

	// Events (5 types x 4 uint16)
	for t := 0; t < 5; t++ {
		for i := 0; i < 4; i++ {
			r.AvgEvents[t][i] = readUint16(buf, &off)
		}
	}

	// QualAvgExpData
	for i := 0; i < 3; i++ {
		r.I[i] = readFloat32(buf, &off)
	}
	r.In = readFloat32(buf, &off)
	for i := 0; i < 2; i++ {
		r.Ibal[i] = readFloat32(buf, &off)
	}
	for i := 0; i < 5; i++ {
		r.Temp[i] = readFloat32(buf, &off)
	}
	r.PF = readFloat32(buf, &off)
	r.DPF = readFloat32(buf, &off)
	for i := 0; i < 2; i++ {
		r.P[i] = readFloat32(buf, &off)
	}
	for i := 0; i < 2; i++ {
		r.Q[i] = readFloat32(buf, &off)
	}
	r.S = readFloat32(buf, &off)
	for phase := 0; phase < 3; phase++ {
		for idx := 0; idx < 2; idx++ {
			r.Energy[phase][idx] = readUint64(buf, &off)
		}
	}

	// QualVarData - event variations (5 types x 4 uint32)
	for t := 0; t < 5; t++ {
		for i := 0; i < 4; i++ {
			r.VarEvents[t][i] = readUint32(buf, &off)
		}
	}

	// Freq variations (2)
	for i := 0; i < 2; i++ {
		r.FreqVar[i] = readQualVariation(buf, &off)
	}
	// Volt variations (6: volt1[3] then volt2[3])
	for i := 0; i < 6; i++ {
		r.VoltVar[i] = readQualVariation(buf, &off)
	}
	// Voltbal
	r.VoltbalVar = readQualVariation(buf, &off)
	// VoltThd (3)
	for i := 0; i < 3; i++ {
		r.VoltThdVar[i] = readQualVariation(buf, &off)
	}
	// VoltHd (3)
	for i := 0; i < 3; i++ {
		r.VoltHdVar[i] = readQualVariation(buf, &off)
	}
	// Pst (3)
	for i := 0; i < 3; i++ {
		r.PstVar[i] = readQualVariation(buf, &off)
	}
	// Plt (3)
	for i := 0; i < 3; i++ {
		r.PltVar[i] = readQualVariation(buf, &off)
	}
	// Svolt (3)
	for i := 0; i < 3; i++ {
		r.SvoltVar[i] = readQualVariation(buf, &off)
	}

	return r, nil
}

// ---------------------------------------------------------------------------
// InfluxDB field builder
// ---------------------------------------------------------------------------

func buildPQ10MinFields(r *pq10minRaw) map[string]interface{} {
	fields := make(map[string]interface{}, 180)

	// Voltage (3 phases)
	fields["voltage_l1"] = float64(r.U[0])
	fields["voltage_l2"] = float64(r.U[1])
	fields["voltage_l3"] = float64(r.U[2])

	// Current (3 phases + neutral)
	fields["current_l1"] = float64(r.I[0])
	fields["current_l2"] = float64(r.I[1])
	fields["current_l3"] = float64(r.I[2])
	fields["current_n"] = float64(r.In)

	// Power
	fields["power_active_total"] = float64(r.P[0])
	fields["power_active_avg"] = float64(r.P[1])
	fields["power_reactive_total"] = float64(r.Q[0])
	fields["power_reactive_avg"] = float64(r.Q[1])
	fields["power_apparent"] = float64(r.S)
	fields["power_factor"] = float64(r.PF)
	fields["displacement_pf"] = float64(r.DPF)

	// Frequency average (average of non-zero samples)
	var freqSum float64
	var freqCount int
	for i := 0; i < 60; i++ {
		if r.Freq[i] > 0 {
			freqSum += float64(r.Freq[i])
			freqCount++
		}
	}
	freqAvg := 0.0
	if freqCount > 0 {
		freqAvg = freqSum / float64(freqCount)
	}
	fields["frequency_avg"] = freqAvg

	// Voltage THD
	fields["voltage_thd_l1"] = float64(r.Uthd[0])
	fields["voltage_thd_l2"] = float64(r.Uthd[1])
	fields["voltage_thd_l3"] = float64(r.Uthd[2])

	// Voltage unbalance
	fields["voltage_unbalance_0"] = float64(r.Ubal[0])
	fields["voltage_unbalance_1"] = float64(r.Ubal[1])

	// Current unbalance
	fields["current_unbalance_0"] = float64(r.Ibal[0])
	fields["current_unbalance_1"] = float64(r.Ibal[1])

	// Harmonics: 3 phases x 24
	for phase := 0; phase < 3; phase++ {
		for h := 0; h < 24; h++ {
			key := fmt.Sprintf("harmonic_l%d_h%d", phase+1, h+1)
			fields[key] = float64(r.Uhd[phase][h])
		}
	}

	// Flicker
	fields["pst_l1"] = float64(r.Pst[0])
	fields["pst_l2"] = float64(r.Pst[1])
	fields["pst_l3"] = float64(r.Pst[2])
	fields["plt_l1"] = float64(r.Plt[0])
	fields["plt_l2"] = float64(r.Plt[1])
	fields["plt_l3"] = float64(r.Plt[2])

	// Signal voltage
	fields["signal_voltage_l1"] = float64(r.Svolt[0])
	fields["signal_voltage_l2"] = float64(r.Svolt[1])
	fields["signal_voltage_l3"] = float64(r.Svolt[2])

	// Temperature
	for i := 0; i < 5; i++ {
		fields[fmt.Sprintf("temp_%d", i)] = float64(r.Temp[i])
	}

	// Events (avg 10min counts) + totals
	for t, evtName := range pq10minEventTypes {
		total := 0
		for i := 0; i < 4; i++ {
			v := int(r.AvgEvents[t][i])
			fields[fmt.Sprintf("event_%s_%d", evtName, i)] = v
			total += v
		}
		fields[fmt.Sprintf("event_%s_total", evtName)] = total
	}

	// Variation fields (key values only, matching Python)
	fields["freq_var1"] = float64(r.FreqVar[0].Val)
	fields["freq_var2"] = float64(r.FreqVar[1].Val)

	for i := 0; i < 3; i++ {
		fields[fmt.Sprintf("volt_var1_l%d", i+1)] = float64(r.VoltVar[i].Val)
		fields[fmt.Sprintf("volt_var2_l%d", i+1)] = float64(r.VoltVar[3+i].Val)
		fields[fmt.Sprintf("thd_var_l%d", i+1)] = float64(r.VoltThdVar[i].Val)
	}

	return fields
}

// ---------------------------------------------------------------------------
// Redis status builder
// ---------------------------------------------------------------------------

func buildPQ10MinStatus(r *pq10minRaw, channelName string, startTime, endTime time.Time, freqAvg float64) map[string]interface{} {
	// Event totals
	eventTotals := make(map[string]int)
	for t, evtName := range pq10minEventTypes {
		total := 0
		for i := 0; i < 4; i++ {
			total += int(r.AvgEvents[t][i])
		}
		eventTotals[evtName] = total
	}

	return map[string]interface{}{
		"last_update": startTime.Format(time.RFC3339),
		"period": map[string]interface{}{
			"start": startTime.Format(time.RFC3339),
			"end":   endTime.Format(time.RFC3339),
		},
		"voltage": map[string]interface{}{
			"l1": float64(r.U[0]),
			"l2": float64(r.U[1]),
			"l3": float64(r.U[2]),
		},
		"current": map[string]interface{}{
			"l1": float64(r.I[0]),
			"l2": float64(r.I[1]),
			"l3": float64(r.I[2]),
			"n":  float64(r.In),
		},
		"power": map[string]interface{}{
			"active":   float64(r.P[0]),
			"reactive": float64(r.Q[0]),
			"apparent": float64(r.S),
			"factor":   float64(r.PF),
		},
		"frequency": freqAvg,
		"quality": map[string]interface{}{
			"voltage_thd":       []float64{float64(r.Uthd[0]), float64(r.Uthd[1]), float64(r.Uthd[2])},
			"voltage_unbalance": []float64{float64(r.Ubal[0]), float64(r.Ubal[1])},
			"current_unbalance": []float64{float64(r.Ibal[0]), float64(r.Ibal[1])},
			"pst":               []float64{float64(r.Pst[0]), float64(r.Pst[1]), float64(r.Pst[2])},
			"plt":               []float64{float64(r.Plt[0]), float64(r.Plt[1]), float64(r.Plt[2])},
		},
		"events_total": eventTotals,
		"timestamp":    time.Now().Unix(),
	}
}

// ---------------------------------------------------------------------------
// Processor
// ---------------------------------------------------------------------------

// PQ10MinDataProcessor processes PQ 10-minute data from the "en10min_data" Redis queue.
type PQ10MinDataProcessor struct {
	*GenericDataProcessor

	mu      sync.RWMutex
	running bool
	cancel  context.CancelFunc
	wg      sync.WaitGroup

	// Statistics
	stats pq10minStats
}

type pq10minStats struct {
	TotalProcessed int64
	ByChannel      map[string]int64
}

// NewPQ10MinDataProcessor creates a new PQ10MinDataProcessor.
func NewPQ10MinDataProcessor(redis *handlers.RedisHandler, influx *handlers.InfluxDBHandler) *PQ10MinDataProcessor {
	config := ProcessorConfig{
		RedisKey:        "en10min_data",
		MeasurementName: "en10min",
		PollInterval:    600 * time.Second, // 10 minutes
	}

	ep := &PQ10MinDataProcessor{
		GenericDataProcessor: NewGenericDataProcessor("PQ10MinProcessor", config, redis, influx),
		stats: pq10minStats{
			ByChannel: map[string]int64{"Main": 0, "Sub": 0},
		},
	}

	ep.GenericDataProcessor.SetProcessFunc(ep.ProcessSingleData)

	return ep
}

// ProcessSingleData parses binary PQ 10-minute data and writes to InfluxDB and Redis.
func (ep *PQ10MinDataProcessor) ProcessSingleData(rawData []byte) error {
	// 1. Parse binary data
	r, err := parsePQ10MinBinary(rawData)
	if err != nil {
		log.Printf("[PQ10MinProcessor] Parse error: %v", err)
		return fmt.Errorf("PQ10Min parse error: %w", err)
	}

	// 2. Channel info
	channelID := r.ID
	channelName, ok := pq10minChannelMapping[channelID]
	if !ok {
		channelName = fmt.Sprintf("Channel_%d", channelID)
	}

	// 3. Timestamps
	startTime := time.Unix(int64(r.StartTs), 0)
	endTime := time.Unix(int64(r.EndTs), 0)

	// 4. Compute frequency average
	var freqSum float64
	var freqCount int
	for i := 0; i < 60; i++ {
		if r.Freq[i] > 0 {
			freqSum += float64(r.Freq[i])
			freqCount++
		}
	}
	freqAvg := 0.0
	if freqCount > 0 {
		freqAvg = freqSum / float64(freqCount)
	}

	// 5. Build InfluxDB fields
	fields := buildPQ10MinFields(r)

	tags := map[string]string{
		"channel_id": fmt.Sprintf("%d", channelID),
		"channel":    channelName,
	}

	// 6. Write to InfluxDB bucket "ntek30"
	if ep.InfluxHandler != nil {
		if err := ep.InfluxHandler.WriteDataPoint("ntek30", ep.Config.MeasurementName, tags, fields, startTime); err != nil {
			return fmt.Errorf("PQ10Min InfluxDB write error: %w", err)
		}
	}

	// 7. Update Redis status (HSET "pq10min_status")
	if ep.RedisHandler != nil {
		statusData := buildPQ10MinStatus(r, channelName, startTime, endTime, freqAvg)
		if err := ep.RedisHandler.SaveSummary("pq10min_status", channelName, statusData); err != nil {
			log.Printf("[PQ10MinProcessor] Redis status update error: %v", err)
		}
	}

	// 8. Update statistics
	ep.mu.Lock()
	ep.stats.TotalProcessed++
	ep.stats.ByChannel[channelName]++
	ep.mu.Unlock()

	log.Printf("[PQ10MinProcessor] processed: channel=%s, time=%s-%s, THD_L1=%.2f%%, Freq_avg=%.3fHz",
		channelName,
		startTime.Format("15:04"),
		endTime.Format("15:04"),
		r.Uthd[0],
		freqAvg,
	)

	return nil
}

// GetLatestPQData retrieves the latest PQ data for a channel from Redis.
func (ep *PQ10MinDataProcessor) GetLatestPQData(channel string) (map[string]interface{}, error) {
	if ep.RedisHandler == nil {
		return nil, fmt.Errorf("Redis handler not available")
	}
	return ep.RedisHandler.GetSummary("pq10min_status", channel)
}

// GetStatistics returns processing statistics as JSON.
func (ep *PQ10MinDataProcessor) GetStatistics() string {
	ep.mu.RLock()
	defer ep.mu.RUnlock()
	b, _ := json.Marshal(ep.stats)
	return string(b)
}

// Start begins the PQ 10-minute data processor.
func (ep *PQ10MinDataProcessor) Start(ctx context.Context) error {
	ep.mu.Lock()
	if ep.running {
		ep.mu.Unlock()
		return fmt.Errorf("PQ10Min processor already running")
	}
	ep.running = true
	ep.mu.Unlock()

	childCtx, cancel := context.WithCancel(ctx)
	ep.cancel = cancel

	ep.wg.Add(1)
	go func() {
		defer ep.wg.Done()
		ep.GenericDataProcessor.ProcessQueue(childCtx)
	}()

	log.Println("[PQ10MinProcessor] started")
	return nil
}

// Stop halts the PQ 10-minute data processor.
func (ep *PQ10MinDataProcessor) Stop() {
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
	log.Println("[PQ10MinProcessor] stopped")
}
