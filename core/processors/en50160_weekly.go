package processors

/*
PQ 10min data weekly Parquet storage module (Go port of en50160_weekly.py)

- Queries data via InfluxDBHandler
- Saves as Parquet file
- Supports standalone execution/testing
*/

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/parquet-go/parquet-go"

	"sv500_core/handlers"
)

// ---------------------------------------------------------------------------
// Helper functions for map value extraction
// ---------------------------------------------------------------------------

func getStringValue(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		if s, ok := v.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func getIntValue(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok && v != nil {
		switch n := v.(type) {
		case int:
			return n
		case int64:
			return int(n)
		case float64:
			return int(n)
		case json.Number:
			i, _ := n.Int64()
			return int(i)
		}
	}
	return 0
}

func getFloat64Value(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok && v != nil {
		switch n := v.(type) {
		case float64:
			return n
		case float32:
			return float64(n)
		case int:
			return float64(n)
		case int64:
			return float64(n)
		case json.Number:
			f, _ := n.Float64()
			return f
		}
	}
	return 0
}

func getTimeValue(m map[string]interface{}, key string) time.Time {
	if v, ok := m[key]; ok && v != nil {
		if t, ok := v.(time.Time); ok {
			return t
		}
	}
	return time.Time{}
}

// ---------------------------------------------------------------------------
// Configuration
// ---------------------------------------------------------------------------

// WeeklyReportConfig holds settings for EN50160 weekly Parquet generation.
type WeeklyReportConfig struct {
	OutputDir      string
	Bucket         string
	Measurement    string
	Channels       []string
	RetentionWeeks int
}

// DefaultWeeklyReportConfig returns the standard configuration.
func DefaultWeeklyReportConfig() WeeklyReportConfig {
	return WeeklyReportConfig{
		OutputDir:      "/usr/local/sv500/reports",
		Bucket:         "ntek30",
		Measurement:    "en10min",
		Channels:       []string{"Main", "Sub"},
		RetentionWeeks: 12,
	}
}

// ---------------------------------------------------------------------------
// Parquet row schema
// ---------------------------------------------------------------------------

// PQWeeklyParquetRow defines the Parquet schema for EN50160 10-minute data.
// Matches the Python schema: voltage, current, power, frequency, THD,
// unbalance, harmonics, flicker, signal voltage, temperature, events, variations.
type PQWeeklyParquetRow struct {
	// Basic info
	Timestamp int64  `parquet:"timestamp,timestamp(microsecond)"`
	ChannelID int32  `parquet:"channel_id"`
	Channel   string `parquet:"channel"`

	// Voltage (3-phase)
	VoltageL1 float32 `parquet:"voltage_l1"`
	VoltageL2 float32 `parquet:"voltage_l2"`
	VoltageL3 float32 `parquet:"voltage_l3"`

	// Current (3-phase + neutral)
	CurrentL1 float32 `parquet:"current_l1"`
	CurrentL2 float32 `parquet:"current_l2"`
	CurrentL3 float32 `parquet:"current_l3"`
	CurrentN  float32 `parquet:"current_n"`

	// Power
	PowerActiveTotal   float32 `parquet:"power_active_total"`
	PowerActiveAvg     float32 `parquet:"power_active_avg"`
	PowerReactiveTotal float32 `parquet:"power_reactive_total"`
	PowerReactiveAvg   float32 `parquet:"power_reactive_avg"`
	PowerApparent      float32 `parquet:"power_apparent"`
	PowerFactor        float32 `parquet:"power_factor"`
	DisplacementPF     float32 `parquet:"displacement_pf"`

	// Frequency
	FrequencyAvg float32 `parquet:"frequency_avg"`

	// THD (3-phase)
	VoltageTHDL1 float32 `parquet:"voltage_thd_l1"`
	VoltageTHDL2 float32 `parquet:"voltage_thd_l2"`
	VoltageTHDL3 float32 `parquet:"voltage_thd_l3"`

	// Unbalance
	VoltageUnbalance0 float32 `parquet:"voltage_unbalance_0"`
	VoltageUnbalance1 float32 `parquet:"voltage_unbalance_1"`
	CurrentUnbalance0 float32 `parquet:"current_unbalance_0"`
	CurrentUnbalance1 float32 `parquet:"current_unbalance_1"`

	// Harmonics (3-phase x 24, H2~H25) - stored as JSON string
	HarmonicsL1 string `parquet:"harmonics_l1"`
	HarmonicsL2 string `parquet:"harmonics_l2"`
	HarmonicsL3 string `parquet:"harmonics_l3"`

	// Flicker
	PstL1 float32 `parquet:"pst_l1"`
	PstL2 float32 `parquet:"pst_l2"`
	PstL3 float32 `parquet:"pst_l3"`
	PltL1 float32 `parquet:"plt_l1"`
	PltL2 float32 `parquet:"plt_l2"`
	PltL3 float32 `parquet:"plt_l3"`

	// Signal Voltage
	SignalVoltageL1 float32 `parquet:"signal_voltage_l1"`
	SignalVoltageL2 float32 `parquet:"signal_voltage_l2"`
	SignalVoltageL3 float32 `parquet:"signal_voltage_l3"`

	// Temperature (5)
	Temp0 float32 `parquet:"temp_0"`
	Temp1 float32 `parquet:"temp_1"`
	Temp2 float32 `parquet:"temp_2"`
	Temp3 float32 `parquet:"temp_3"`
	Temp4 float32 `parquet:"temp_4"`

	// Event totals
	EventSagTotal       int32 `parquet:"event_sag_total"`
	EventSwellTotal     int32 `parquet:"event_swell_total"`
	EventShortIntrTotal int32 `parquet:"event_short_intr_total"`
	EventLongIntrTotal  int32 `parquet:"event_long_intr_total"`
	EventRVCTotal       int32 `parquet:"event_rvc_total"`

	// Variation
	FreqVar1  float32 `parquet:"freq_var1"`
	FreqVar2  float32 `parquet:"freq_var2"`
	VoltVar1L1 float32 `parquet:"volt_var1_l1"`
	VoltVar1L2 float32 `parquet:"volt_var1_l2"`
	VoltVar1L3 float32 `parquet:"volt_var1_l3"`
	VoltVar2L1 float32 `parquet:"volt_var2_l1"`
	VoltVar2L2 float32 `parquet:"volt_var2_l2"`
	VoltVar2L3 float32 `parquet:"volt_var2_l3"`
	THDVarL1   float32 `parquet:"thd_var_l1"`
	THDVarL2   float32 `parquet:"thd_var_l2"`
	THDVarL3   float32 `parquet:"thd_var_l3"`
}

// ---------------------------------------------------------------------------
// PQWeeklyParquetWriter
// ---------------------------------------------------------------------------

// PQWeeklyParquetWriter generates EN50160 10-minute weekly Parquet files.
type PQWeeklyParquetWriter struct {
	Config        WeeklyReportConfig
	influxHandler *handlers.InfluxDBHandler
}

// NewPQWeeklyParquetWriter creates and initializes a new writer.
func NewPQWeeklyParquetWriter(config WeeklyReportConfig) *PQWeeklyParquetWriter {
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Printf("[EN50160Weekly] failed to create output dir %s: %v", config.OutputDir, err)
	}
	log.Printf("[EN50160Weekly] output dir: %s", config.OutputDir)

	return &PQWeeklyParquetWriter{
		Config:        config,
		influxHandler: handlers.NewInfluxDBHandler(config.Bucket),
	}
}

// ---------------------------------------------------------------------------
// GetWeekRange computes Monday 00:00 to Sunday 23:59:59
// ---------------------------------------------------------------------------

// GetWeekRange calculates the weekly range (Monday 00:00 ~ Sunday 23:59:59).
// If targetDate is zero, it returns the previous week's range.
func GetWeekRange(targetDate time.Time) (start, end time.Time) {
	if targetDate.IsZero() {
		// Previous week calculation
		today := time.Now().Truncate(24 * time.Hour)
		weekday := int(today.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		thisMonday := today.AddDate(0, 0, -(weekday - 1))
		start = thisMonday.AddDate(0, 0, -7)
	} else {
		// Monday of the week containing targetDate
		t := targetDate.Truncate(24 * time.Hour)
		weekday := int(t.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start = t.AddDate(0, 0, -(weekday - 1))
	}

	// Sunday 23:59:59
	end = start.AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	return start, end
}

// ---------------------------------------------------------------------------
// QueryWeekData queries InfluxDB for one week of data.
// ---------------------------------------------------------------------------

func (w *PQWeeklyParquetWriter) QueryWeekData(start, end time.Time, channel string) ([]map[string]interface{}, error) {
	log.Printf("[EN50160Weekly] data query: %s ~ %s",
		start.Format("2006-01-02 15:04"), end.Format("2006-01-02 15:04"))

	channelFilter := ""
	if channel != "" {
		channelFilter = fmt.Sprintf(`|> filter(fn: (r) => r["channel"] == "%s")`, channel)
	}

	query := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r["_measurement"] == "%s")
			%s
			|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
			|> sort(columns: ["_time"])
	`, w.Config.Bucket, start.UTC().Format(time.RFC3339), end.UTC().Format(time.RFC3339),
		w.Config.Measurement, channelFilter)

	results, err := w.influxHandler.ExecuteQuery(query)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}

	if len(results) == 0 {
		log.Printf("[EN50160Weekly] no data returned")
		return nil, nil
	}

	log.Printf("[EN50160Weekly] query complete: %d records", len(results))
	return results, nil
}

// ---------------------------------------------------------------------------
// prepareParquetRows converts query results to Parquet row structs.
// ---------------------------------------------------------------------------

func (w *PQWeeklyParquetWriter) prepareParquetRows(records []map[string]interface{}) []PQWeeklyParquetRow {
	rows := make([]PQWeeklyParquetRow, 0, len(records))

	for _, rec := range records {
		ts := getTimeValue(rec, "_time")
		if ts.IsZero() {
			ts = getTimeValue(rec, "timestamp")
		}

		// Build harmonics arrays as JSON strings
		harmonicsL1 := make([]float32, 24)
		harmonicsL2 := make([]float32, 24)
		harmonicsL3 := make([]float32, 24)
		for h := 1; h <= 24; h++ {
			harmonicsL1[h-1] = getFloat32FromMap(rec, fmt.Sprintf("harmonic_l1_h%d", h))
			harmonicsL2[h-1] = getFloat32FromMap(rec, fmt.Sprintf("harmonic_l2_h%d", h))
			harmonicsL3[h-1] = getFloat32FromMap(rec, fmt.Sprintf("harmonic_l3_h%d", h))
		}

		hL1JSON, _ := json.Marshal(harmonicsL1)
		hL2JSON, _ := json.Marshal(harmonicsL2)
		hL3JSON, _ := json.Marshal(harmonicsL3)

		row := PQWeeklyParquetRow{
			Timestamp: ts.UnixMicro(),
			ChannelID: int32(getIntValue(rec, "channel_id")),
			Channel:   getStringValue(rec, "channel"),

			// Voltage
			VoltageL1: getFloat32FromMap(rec, "voltage_l1"),
			VoltageL2: getFloat32FromMap(rec, "voltage_l2"),
			VoltageL3: getFloat32FromMap(rec, "voltage_l3"),

			// Current
			CurrentL1: getFloat32FromMap(rec, "current_l1"),
			CurrentL2: getFloat32FromMap(rec, "current_l2"),
			CurrentL3: getFloat32FromMap(rec, "current_l3"),
			CurrentN:  getFloat32FromMap(rec, "current_n"),

			// Power
			PowerActiveTotal:   getFloat32FromMap(rec, "power_active_total"),
			PowerActiveAvg:     getFloat32FromMap(rec, "power_active_avg"),
			PowerReactiveTotal: getFloat32FromMap(rec, "power_reactive_total"),
			PowerReactiveAvg:   getFloat32FromMap(rec, "power_reactive_avg"),
			PowerApparent:      getFloat32FromMap(rec, "power_apparent"),
			PowerFactor:        getFloat32FromMap(rec, "power_factor"),
			DisplacementPF:     getFloat32FromMap(rec, "displacement_pf"),

			// Frequency
			FrequencyAvg: getFloat32FromMap(rec, "frequency_avg"),

			// THD
			VoltageTHDL1: getFloat32FromMap(rec, "voltage_thd_l1"),
			VoltageTHDL2: getFloat32FromMap(rec, "voltage_thd_l2"),
			VoltageTHDL3: getFloat32FromMap(rec, "voltage_thd_l3"),

			// Unbalance
			VoltageUnbalance0: getFloat32FromMap(rec, "voltage_unbalance_0"),
			VoltageUnbalance1: getFloat32FromMap(rec, "voltage_unbalance_1"),
			CurrentUnbalance0: getFloat32FromMap(rec, "current_unbalance_0"),
			CurrentUnbalance1: getFloat32FromMap(rec, "current_unbalance_1"),

			// Harmonics (JSON)
			HarmonicsL1: string(hL1JSON),
			HarmonicsL2: string(hL2JSON),
			HarmonicsL3: string(hL3JSON),

			// Flicker
			PstL1: getFloat32FromMap(rec, "pst_l1"),
			PstL2: getFloat32FromMap(rec, "pst_l2"),
			PstL3: getFloat32FromMap(rec, "pst_l3"),
			PltL1: getFloat32FromMap(rec, "plt_l1"),
			PltL2: getFloat32FromMap(rec, "plt_l2"),
			PltL3: getFloat32FromMap(rec, "plt_l3"),

			// Signal Voltage
			SignalVoltageL1: getFloat32FromMap(rec, "signal_voltage_l1"),
			SignalVoltageL2: getFloat32FromMap(rec, "signal_voltage_l2"),
			SignalVoltageL3: getFloat32FromMap(rec, "signal_voltage_l3"),

			// Temperature
			Temp0: getFloat32FromMap(rec, "temp_0"),
			Temp1: getFloat32FromMap(rec, "temp_1"),
			Temp2: getFloat32FromMap(rec, "temp_2"),
			Temp3: getFloat32FromMap(rec, "temp_3"),
			Temp4: getFloat32FromMap(rec, "temp_4"),

			// Events
			EventSagTotal:       int32(getIntValue(rec, "event_sag_total")),
			EventSwellTotal:     int32(getIntValue(rec, "event_swell_total")),
			EventShortIntrTotal: int32(getIntValue(rec, "event_short_intr_total")),
			EventLongIntrTotal:  int32(getIntValue(rec, "event_long_intr_total")),
			EventRVCTotal:       int32(getIntValue(rec, "event_rvc_total")),

			// Variation
			FreqVar1:   getFloat32FromMap(rec, "freq_var1"),
			FreqVar2:   getFloat32FromMap(rec, "freq_var2"),
			VoltVar1L1: getFloat32FromMap(rec, "volt_var1_l1"),
			VoltVar1L2: getFloat32FromMap(rec, "volt_var1_l2"),
			VoltVar1L3: getFloat32FromMap(rec, "volt_var1_l3"),
			VoltVar2L1: getFloat32FromMap(rec, "volt_var2_l1"),
			VoltVar2L2: getFloat32FromMap(rec, "volt_var2_l2"),
			VoltVar2L3: getFloat32FromMap(rec, "volt_var2_l3"),
			THDVarL1:   getFloat32FromMap(rec, "thd_var_l1"),
			THDVarL2:   getFloat32FromMap(rec, "thd_var_l2"),
			THDVarL3:   getFloat32FromMap(rec, "thd_var_l3"),
		}

		rows = append(rows, row)
	}

	return rows
}

// getFloat32FromMap extracts a float32 from a map, returning 0 on missing/nil.
func getFloat32FromMap(m map[string]interface{}, key string) float32 {
	return float32(getFloat64Value(m, key))
}

// ---------------------------------------------------------------------------
// GenerateWeeklyParquet creates the weekly Parquet file.
// ---------------------------------------------------------------------------

func (w *PQWeeklyParquetWriter) GenerateWeeklyParquet(start, end time.Time, channel string, reportSummary map[string]interface{}) string {
	// 1. Determine time range
	if start.IsZero() || end.IsZero() {
		start, end = GetWeekRange(time.Time{})
	}

	weekStr := end.Format("20060102")
	log.Printf("[EN50160Weekly] ==================================================")
	log.Printf("[EN50160Weekly] weekly Parquet generation start: %s", weekStr)
	log.Printf("[EN50160Weekly] period: %s ~ %s", start.Format("2006-01-02"), end.Format("2006-01-02"))
	log.Printf("[EN50160Weekly] ==================================================")

	// 2. Query data
	records, err := w.QueryWeekData(start, end, channel)
	if err != nil {
		log.Printf("[EN50160Weekly] query error: %v", err)
		return ""
	}
	if len(records) == 0 {
		log.Printf("[EN50160Weekly] no data: %s", weekStr)
		return ""
	}

	// 3. Prepare Parquet rows
	rows := w.prepareParquetRows(records)
	log.Printf("[EN50160Weekly] conversion complete: %d records", len(rows))

	// 4. Determine file path
	filename := fmt.Sprintf("en50160_weekly_%s.parquet", weekStr)
	var filePath string
	if channel != "" {
		filePath = filepath.Join(w.Config.OutputDir, channel, filename)
	} else {
		filePath = filepath.Join(w.Config.OutputDir, filename)
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		log.Printf("[EN50160Weekly] mkdir error: %v", err)
		return ""
	}

	// 5. Write Parquet file (snappy compression)
	f, err := os.Create(filePath)
	if err != nil {
		log.Printf("[EN50160Weekly] file create error: %v", err)
		return ""
	}
	defer f.Close()

	// Build writer options
	writerOpts := []parquet.WriterOption{
		parquet.Compression(&parquet.Snappy),
	}
	if reportSummary != nil {
		summaryJSON, err := json.Marshal(reportSummary)
		if err == nil {
			writerOpts = append(writerOpts, parquet.KeyValueMetadata("en50160_summary", string(summaryJSON)))
			compliance := reportSummary["compliance"]
			log.Printf("[EN50160Weekly] summary metadata added: compliance=%v", compliance)
		}
	}

	pw := parquet.NewGenericWriter[PQWeeklyParquetRow](f, writerOpts...)

	if _, err := pw.Write(rows); err != nil {
		log.Printf("[EN50160Weekly] parquet write error: %v", err)
		return ""
	}
	if err := pw.Close(); err != nil {
		log.Printf("[EN50160Weekly] parquet close error: %v", err)
		return ""
	}

	info, _ := os.Stat(filePath)
	sizeKB := float64(0)
	if info != nil {
		sizeKB = float64(info.Size()) / 1024
	}

	log.Printf("[EN50160Weekly] Parquet saved!")
	log.Printf("[EN50160Weekly]   file: %s", filePath)
	log.Printf("[EN50160Weekly]   records: %d", len(rows))
	log.Printf("[EN50160Weekly]   size: %.1f KB", sizeKB)

	return filePath
}

// ---------------------------------------------------------------------------
// CleanupOldFiles removes old Parquet files beyond retention period.
// ---------------------------------------------------------------------------

func (w *PQWeeklyParquetWriter) CleanupOldFiles(channel string) int {
	var targetDirs []string
	if channel != "" {
		targetDirs = []string{filepath.Join(w.Config.OutputDir, channel)}
	} else {
		entries, err := os.ReadDir(w.Config.OutputDir)
		if err != nil {
			log.Printf("[EN50160Weekly] cleanup read dir error: %v", err)
			return 0
		}
		for _, e := range entries {
			if e.IsDir() {
				targetDirs = append(targetDirs, filepath.Join(w.Config.OutputDir, e.Name()))
			}
		}
	}

	cutoff := time.Now().AddDate(0, 0, -w.Config.RetentionWeeks*7)
	totalDeleted := 0

	for _, dir := range targetDirs {
		matches, err := filepath.Glob(filepath.Join(dir, "en50160_weekly_*.parquet"))
		if err != nil {
			continue
		}

		for _, file := range matches {
			base := filepath.Base(file)
			// Filename: en50160_weekly_20251221.parquet
			parts := strings.Split(strings.TrimSuffix(base, ".parquet"), "_")
			if len(parts) == 0 {
				continue
			}
			dateStr := parts[len(parts)-1]

			fileDate, err := time.Parse("20060102", dateStr)
			if err != nil {
				log.Printf("[EN50160Weekly] date parse failed (skip): %s", base)
				continue
			}
			if fileDate.Before(cutoff) {
				if err := os.Remove(file); err != nil {
					log.Printf("[EN50160Weekly] delete failed: %s - %v", base, err)
				} else {
					totalDeleted++
					log.Printf("[EN50160Weekly] deleted: %s/%s", filepath.Base(filepath.Dir(file)), base)
				}
			}
		}
	}

	if totalDeleted > 0 {
		log.Printf("[EN50160Weekly] cleaned up %d files", totalDeleted)
	} else {
		log.Printf("[EN50160Weekly] no files to clean up")
	}

	return totalDeleted
}

// ---------------------------------------------------------------------------
// Singleton writer
// ---------------------------------------------------------------------------

var defaultWeeklyWriter *PQWeeklyParquetWriter

// GetParquetWriter returns a singleton PQWeeklyParquetWriter.
func GetParquetWriter() *PQWeeklyParquetWriter {
	if defaultWeeklyWriter == nil {
		config := DefaultWeeklyReportConfig()
		config.RetentionWeeks = 52
		defaultWeeklyWriter = NewPQWeeklyParquetWriter(config)
	}
	return defaultWeeklyWriter
}

// ---------------------------------------------------------------------------
// Entry point
// ---------------------------------------------------------------------------

// RunWeeklyJob generates the weekly EN50160 Parquet file.
func RunWeeklyJob(start, end time.Time, channel string, reportSummary map[string]interface{}) string {
	log.Printf("[EN50160Weekly] weekly Parquet generation job started")

	writer := GetParquetWriter()

	filePath := writer.GenerateWeeklyParquet(start, end, channel, reportSummary)

	if filePath != "" {
		log.Printf("[EN50160Weekly] weekly Parquet generation success: %s", filePath)
	} else {
		log.Printf("[EN50160Weekly] weekly Parquet generation failed or no data")
	}

	// Clean up old files
	writer.CleanupOldFiles(channel)
	log.Printf("[EN50160Weekly] weekly Parquet generation job complete")

	return filePath
}
