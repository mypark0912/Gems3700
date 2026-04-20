package processors

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"sv500_core/handlers"
)

// EN50160RecordSize is the total binary struct size: 104 bytes.
// Format: '<' + 'i' + 'II' + 'I' + 'HH' + 'HHH' + 'HHH' + 'H' + 'HHH' + 'HHH' + 'HHH' + 'HHH' + 'HHH' + 'HHHH' + 'HHHH' + 'HHHH' + 'HHHH' + 'HHHH'
const EN50160RecordSize = 104

// en50160ChannelMapping maps channel IDs to names.
var en50160ChannelMapping = map[int32]string{
	0: "Main",
	1: "Sub",
}

// en50160RawReport represents the raw binary struct fields after parsing.
type en50160RawReport struct {
	ID         int32
	STime      uint32
	ETime      uint32
	Compliance uint32
	Fvar1      uint16
	Fvar2      uint16
	Volt1      [3]uint16
	Volt2      [3]uint16
	Voltbal    uint16
	VoltThd    [3]uint16
	VoltHd     [3]uint16
	Pst        [3]uint16
	Plt        [3]uint16
	Svolt      [3]uint16
	Sag        [4]uint16
	Swell      [4]uint16
	ShortIntr  [4]uint16
	LongIntr   [4]uint16
	Rvc        [4]uint16
}

// EN50160StatusManager manages EN50160 status in Redis via HSET "en50160_status".
type EN50160StatusManager struct {
	Redis *handlers.RedisHandler
}

// UpdateReportStatus sets the channel status in Redis hash "en50160_status".
func (m *EN50160StatusManager) UpdateReportStatus(channelName string, updateData map[string]interface{}) error {
	jsonBytes, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("marshal status data: %w", err)
	}
	return m.Redis.HSet(context.Background(), "en50160_status", channelName, string(jsonBytes))
}

// EN50160DataProcessor processes EN50160 report data from a Redis queue.
type EN50160DataProcessor struct {
	*GenericDataProcessor

	StatusManager *EN50160StatusManager

	mu      sync.RWMutex
	running bool
	cancel  context.CancelFunc
	wg      sync.WaitGroup

	// Statistics
	reportStats en50160Stats
}

type en50160Stats struct {
	TotalReports   int64
	ByChannel      map[string]int64
	CompliancePass int64
	ComplianceFail int64
	TotalEvents    map[string]int64
}

// NewEN50160DataProcessor creates a new EN50160DataProcessor.
func NewEN50160DataProcessor(redis *handlers.RedisHandler, influx *handlers.InfluxDBHandler) *EN50160DataProcessor {
	config := ProcessorConfig{
		RedisKey:        "en50160_report",
		MeasurementName: "en50160",
		PollInterval:    3600 * 24 * time.Second, // 24 hours
	}

	ep := &EN50160DataProcessor{
		GenericDataProcessor: NewGenericDataProcessor("EN50160Processor", config, redis, influx),
		StatusManager:        &EN50160StatusManager{Redis: redis},
		reportStats: en50160Stats{
			ByChannel: map[string]int64{"Main": 0, "Sub": 0},
			TotalEvents: map[string]int64{
				"sag":                0,
				"swell":             0,
				"short_interruption": 0,
				"long_interruption":  0,
				"rvc":               0,
			},
		},
	}

	ep.GenericDataProcessor.SetProcessFunc(ep.ProcessSingleData)

	return ep
}

// parseEN50160Binary parses 104 bytes of little-endian binary data into an en50160RawReport.
func parseEN50160Binary(buf []byte) (*en50160RawReport, error) {
	if len(buf) < EN50160RecordSize {
		return nil, fmt.Errorf("EN50160 data too short: got %d bytes, need %d", len(buf), EN50160RecordSize)
	}

	r := &en50160RawReport{}
	offset := 0

	// id: int32 (4 bytes)
	r.ID = int32(binary.LittleEndian.Uint32(buf[offset : offset+4]))
	offset += 4

	// sTime: uint32 (4 bytes)
	r.STime = binary.LittleEndian.Uint32(buf[offset : offset+4])
	offset += 4

	// eTime: uint32 (4 bytes)
	r.ETime = binary.LittleEndian.Uint32(buf[offset : offset+4])
	offset += 4

	// compliance: uint32 (4 bytes)
	r.Compliance = binary.LittleEndian.Uint32(buf[offset : offset+4])
	offset += 4

	// Fvar1, Fvar2: uint16 x 2
	r.Fvar1 = binary.LittleEndian.Uint16(buf[offset : offset+2])
	offset += 2
	r.Fvar2 = binary.LittleEndian.Uint16(buf[offset : offset+2])
	offset += 2

	// Volt1[3]: uint16 x 3
	for i := 0; i < 3; i++ {
		r.Volt1[i] = binary.LittleEndian.Uint16(buf[offset : offset+2])
		offset += 2
	}

	// Volt2[3]: uint16 x 3
	for i := 0; i < 3; i++ {
		r.Volt2[i] = binary.LittleEndian.Uint16(buf[offset : offset+2])
		offset += 2
	}

	// Voltbal: uint16
	r.Voltbal = binary.LittleEndian.Uint16(buf[offset : offset+2])
	offset += 2

	// VoltThd[3]: uint16 x 3
	for i := 0; i < 3; i++ {
		r.VoltThd[i] = binary.LittleEndian.Uint16(buf[offset : offset+2])
		offset += 2
	}

	// VoltHd[3]: uint16 x 3
	for i := 0; i < 3; i++ {
		r.VoltHd[i] = binary.LittleEndian.Uint16(buf[offset : offset+2])
		offset += 2
	}

	// Pst[3]: uint16 x 3
	for i := 0; i < 3; i++ {
		r.Pst[i] = binary.LittleEndian.Uint16(buf[offset : offset+2])
		offset += 2
	}

	// Plt[3]: uint16 x 3
	for i := 0; i < 3; i++ {
		r.Plt[i] = binary.LittleEndian.Uint16(buf[offset : offset+2])
		offset += 2
	}

	// Svolt[3]: uint16 x 3
	for i := 0; i < 3; i++ {
		r.Svolt[i] = binary.LittleEndian.Uint16(buf[offset : offset+2])
		offset += 2
	}

	// Sag[4]: uint16 x 4
	for i := 0; i < 4; i++ {
		r.Sag[i] = binary.LittleEndian.Uint16(buf[offset : offset+2])
		offset += 2
	}

	// Swell[4]: uint16 x 4
	for i := 0; i < 4; i++ {
		r.Swell[i] = binary.LittleEndian.Uint16(buf[offset : offset+2])
		offset += 2
	}

	// ShortIntr[4]: uint16 x 4
	for i := 0; i < 4; i++ {
		r.ShortIntr[i] = binary.LittleEndian.Uint16(buf[offset : offset+2])
		offset += 2
	}

	// LongIntr[4]: uint16 x 4
	for i := 0; i < 4; i++ {
		r.LongIntr[i] = binary.LittleEndian.Uint16(buf[offset : offset+2])
		offset += 2
	}

	// Rvc[4]: uint16 x 4
	for i := 0; i < 4; i++ {
		r.Rvc[i] = binary.LittleEndian.Uint16(buf[offset : offset+2])
		offset += 2
	}

	return r, nil
}

// uint16SliceToIntSlice converts a fixed-size uint16 array slice to []int.
func uint16ArrayToIntSlice(arr []uint16) []int {
	result := make([]int, len(arr))
	for i, v := range arr {
		result[i] = int(v)
	}
	return result
}

// sumIntSlice returns the sum of an int slice.
func sumIntSlice(s []int) int {
	total := 0
	for _, v := range s {
		total += v
	}
	return total
}

// extractReportData converts a parsed raw report into the report data map
// matching the Python _extract_report_data output.
func extractReportData(raw *en50160RawReport) map[string]interface{} {
	channelID := raw.ID
	channelName, ok := en50160ChannelMapping[channelID]
	if !ok {
		channelName = fmt.Sprintf("Channel_%d", channelID)
	}

	startTime := time.Unix(int64(raw.STime), 0)
	endTime := time.Unix(int64(raw.ETime), 0)

	volt1Var := uint16ArrayToIntSlice(raw.Volt1[:])
	volt2Var := uint16ArrayToIntSlice(raw.Volt2[:])
	voltThdVar := uint16ArrayToIntSlice(raw.VoltThd[:])
	voltHdVar := uint16ArrayToIntSlice(raw.VoltHd[:])
	pstVar := uint16ArrayToIntSlice(raw.Pst[:])
	pltVar := uint16ArrayToIntSlice(raw.Plt[:])
	svoltVar := uint16ArrayToIntSlice(raw.Svolt[:])
	sag := uint16ArrayToIntSlice(raw.Sag[:])
	swell := uint16ArrayToIntSlice(raw.Swell[:])
	shortIntr := uint16ArrayToIntSlice(raw.ShortIntr[:])
	longIntr := uint16ArrayToIntSlice(raw.LongIntr[:])
	rvc := uint16ArrayToIntSlice(raw.Rvc[:])

	eventCount := sumIntSlice(sag) + sumIntSlice(swell) + sumIntSlice(shortIntr) + sumIntSlice(longIntr) + sumIntSlice(rvc)

	complianceStatus := "PASS"
	if raw.Compliance == 0 {
		complianceStatus = "FAIL"
	}

	log.Printf("[EN50160Processor] Report extracted: channel=%s, period=%s ~ %s, compliance=%s, total_events=%d",
		channelName,
		startTime.Format("2006-01-02 15:04"),
		endTime.Format("2006-01-02 15:04"),
		complianceStatus,
		eventCount,
	)

	return map[string]interface{}{
		"timestamp":          endTime,
		"channel_id":         int(channelID),
		"channel_name":       channelName,
		"start_time":         startTime,
		"end_time":           endTime,
		"compliance":         int(raw.Compliance),
		"freq_var1":          int(raw.Fvar1),
		"freq_var2":          int(raw.Fvar2),
		"volt1_var":          volt1Var,
		"volt2_var":          volt2Var,
		"voltbal_var":        int(raw.Voltbal),
		"volt_thd_var":       voltThdVar,
		"volt_hd_var":        voltHdVar,
		"pst_var":            pstVar,
		"plt_var":            pltVar,
		"svolt_var":          svoltVar,
		"sag":                sag,
		"swell":              swell,
		"short_interruption": shortIntr,
		"long_interruption":  longIntr,
		"rvc":                rvc,
		"event_count":        eventCount,
	}
}

// makeSummaryData creates the summary/status data map matching the Python _make_summary_data.
func makeSummaryData(reportData map[string]interface{}) map[string]interface{} {
	startTime := reportData["start_time"].(time.Time)
	endTime := reportData["end_time"].(time.Time)
	sag := reportData["sag"].([]int)
	swell := reportData["swell"].([]int)
	shortIntr := reportData["short_interruption"].([]int)
	longIntr := reportData["long_interruption"].([]int)
	rvc := reportData["rvc"].([]int)

	return map[string]interface{}{
		"last_report_time":         startTime.Format(time.RFC3339),
		"end_time":                 endTime.Format(time.RFC3339),
		"compliance":               reportData["compliance"],
		"freq_var1":                reportData["freq_var1"],
		"freq_var2":                reportData["freq_var2"],
		"volt1_var":                reportData["volt1_var"],
		"volt2_var":                reportData["volt2_var"],
		"voltbal_var":              reportData["voltbal_var"],
		"volt_thd_var":             reportData["volt_thd_var"],
		"volt_hd_var":              reportData["volt_hd_var"],
		"pst_var":                  reportData["pst_var"],
		"plt_var":                  reportData["plt_var"],
		"svolt_var":                reportData["svolt_var"],
		"sag":                      sag,
		"swell":                    swell,
		"short_interruption":       shortIntr,
		"long_interruption":        longIntr,
		"rvc":                      rvc,
		"total_sag":                sumIntSlice(sag),
		"total_swell":              sumIntSlice(swell),
		"total_short_interruption": sumIntSlice(shortIntr),
		"total_long_interruption":  sumIntSlice(longIntr),
		"total_rvc":                sumIntSlice(rvc),
		"last_update":              time.Now().Unix(),
	}
}

// ProcessSingleData parses binary EN50160 data and writes to InfluxDB and Redis.
func (ep *EN50160DataProcessor) ProcessSingleData(rawData []byte) error {
	// 1. Parse binary data
	raw, err := parseEN50160Binary(rawData)
	if err != nil {
		log.Printf("[EN50160Processor] Parse error: %v", err)
		return fmt.Errorf("EN50160 parse error: %w", err)
	}

	// 2. Extract report data
	reportData := extractReportData(raw)

	// 3. Write to InfluxDB
	if ep.InfluxHandler != nil {
		if err := ep.InfluxHandler.WriteEN50160Report(reportData, ep.Config.MeasurementName); err != nil {
			return fmt.Errorf("EN50160 InfluxDB write error: %w", err)
		}
	}

	// 4. Update Redis status via EN50160StatusManager
	channelName := reportData["channel_name"].(string)
	summaryData := makeSummaryData(reportData)
	if err := ep.StatusManager.UpdateReportStatus(channelName, summaryData); err != nil {
		log.Printf("[EN50160Processor] Status update error: %v", err)
	}

	// 5. Run weekly job
	startTime := reportData["start_time"].(time.Time)
	endTime := reportData["end_time"].(time.Time)
	RunWeeklyJob(startTime, endTime, channelName, summaryData)

	// 6. Run diagnosis report and trend training jobs
	RunDiagnosisReportJob(endTime, channelName)
	RunTrendTrainingJob(endTime, channelName)

	// 7. Update statistics
	ep.updateStatistics(reportData)

	log.Printf("[EN50160Processor] Report processed: channel=%s, id=%d, period=%s ~ %s",
		channelName,
		reportData["channel_id"],
		reportData["start_time"].(time.Time).Format("2006-01-02 15:04"),
		reportData["end_time"].(time.Time).Format("2006-01-02 15:04"),
	)

	return nil
}

// updateStatistics updates internal report statistics.
func (ep *EN50160DataProcessor) updateStatistics(reportData map[string]interface{}) {
	ep.mu.Lock()
	defer ep.mu.Unlock()

	ep.reportStats.TotalReports++

	channel := reportData["channel_name"].(string)
	ep.reportStats.ByChannel[channel]++

	compliance := reportData["compliance"].(int)
	if compliance != 0 {
		ep.reportStats.CompliancePass++
	} else {
		ep.reportStats.ComplianceFail++
	}

	ep.reportStats.TotalEvents["sag"] += int64(sumIntSlice(reportData["sag"].([]int)))
	ep.reportStats.TotalEvents["swell"] += int64(sumIntSlice(reportData["swell"].([]int)))
	ep.reportStats.TotalEvents["short_interruption"] += int64(sumIntSlice(reportData["short_interruption"].([]int)))
	ep.reportStats.TotalEvents["long_interruption"] += int64(sumIntSlice(reportData["long_interruption"].([]int)))
	ep.reportStats.TotalEvents["rvc"] += int64(sumIntSlice(reportData["rvc"].([]int)))
}

// Start begins the EN50160 data processor.
func (ep *EN50160DataProcessor) Start(ctx context.Context) error {
	ep.mu.Lock()
	if ep.running {
		ep.mu.Unlock()
		return fmt.Errorf("EN50160 processor already running")
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

	log.Println("[EN50160Processor] started")
	return nil
}

// Stop halts the EN50160 data processor.
func (ep *EN50160DataProcessor) Stop() {
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
	log.Println("[EN50160Processor] stopped")
}

// RunWeeklyJob, RunDiagnosisReportJob, RunTrendTrainingJob are implemented in
// en50160_weekly.go and diagnosis_weekly.go respectively.
