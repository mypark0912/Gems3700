package processors

/*
Diagnosis data + trend Parquet storage module (Go port of diagnosis_weekly.py)

Called alongside EN50160 weekly report generation to:
1. Query the closest diagnosis/powerquality data to endTime from InfluxDB
2. Query 30-day trend for items with status > 1
3. Save results as Parquet
*/

import (
	"context"
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
// Configuration
// ---------------------------------------------------------------------------

// DiagnosisReportConfig holds settings for diagnosis Parquet generation.
type DiagnosisReportConfig struct {
	OutputDir      string
	Bucket         string
	Modes          []string
	RetentionWeeks int
}

// DefaultDiagnosisReportConfig returns the standard configuration.
func DefaultDiagnosisReportConfig() DiagnosisReportConfig {
	return DiagnosisReportConfig{
		OutputDir:      "/usr/local/sv500/reports",
		Bucket:         "ntek",
		Modes:          []string{"diagnosis", "powerquality"},
		RetentionWeeks: 52,
	}
}

// ---------------------------------------------------------------------------
// Parquet row schemas
// ---------------------------------------------------------------------------

// DiagnosisParquetRow stores one diagnosis report as a single Parquet row.
// Complex nested data (main, detail, trends) is JSON-serialized.
type DiagnosisParquetRow struct {
	AssetName    string `parquet:"asset_name"`
	ChannelName  string `parquet:"channel_name"`
	EndTime      string `parquet:"end_time"`
	GeneratedAt  string `parquet:"generated_at"`
	Diagnosis    string `parquet:"diagnosis"`    // JSON
	PowerQuality string `parquet:"powerquality"` // JSON
}

// TrendTrainingRow stores one time-series data point for LightGBM training.
type TrendTrainingRow struct {
	Asset       string  `parquet:"asset"`
	Channel     string  `parquet:"channel"`
	Group       string  `parquet:"group"`
	Title       string  `parquet:"title"`
	Unit        string  `parquet:"unit"`
	Timestamp   string  `parquet:"timestamp"`
	Value       float64 `parquet:"value"`
	Status      int64   `parquet:"status"`
	IStatus     int64   `parquet:"istatus"`
	CollectedAt string  `parquet:"collected_at"`
}

// ---------------------------------------------------------------------------
// DiagnosisReportWriter
// ---------------------------------------------------------------------------

// DiagnosisReportWriter generates diagnosis + trend Parquet files.
type DiagnosisReportWriter struct {
	Config        DiagnosisReportConfig
	influxHandler *handlers.InfluxDBHandler
}

// NewDiagnosisReportWriter creates and initializes a new writer.
func NewDiagnosisReportWriter(config DiagnosisReportConfig) *DiagnosisReportWriter {
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Printf("[DiagnosisWeekly] failed to create output dir %s: %v", config.OutputDir, err)
	}
	log.Printf("[DiagnosisWeekly] output dir: %s", config.OutputDir)

	return &DiagnosisReportWriter{
		Config:        config,
		influxHandler: handlers.NewInfluxDBHandler(config.Bucket),
	}
}

// ---------------------------------------------------------------------------
// localToUTC converts a local (no-timezone) time to UTC.
// ---------------------------------------------------------------------------

func localToUTC(local time.Time) time.Time {
	if local.Location() == time.UTC {
		return local
	}
	// Treat as local system time and convert to UTC.
	return local.UTC()
}

// ---------------------------------------------------------------------------
// GetAssetInfo reads channel asset info from Redis.
// Returns (assetName, assetType) or empty strings if not found/disabled.
// ---------------------------------------------------------------------------

func GetAssetInfo(channelName string) (name string, assetType string) {
	redisClient := handlers.GetRedisClient("127.0.0.1", 0)
	if redisClient == nil {
		log.Printf("[DiagnosisWeekly] Redis client unavailable")
		return "", ""
	}

	ctx := context.Background()
	setupJSON, err := redisClient.HGet(ctx, "System", "setup").Result()
	if err != nil {
		log.Printf("[DiagnosisWeekly] Redis setup read error: %v", err)
		return "", ""
	}

	var setup map[string]interface{}
	if err := json.Unmarshal([]byte(setupJSON), &setup); err != nil {
		log.Printf("[DiagnosisWeekly] setup JSON parse error: %v", err)
		return "", ""
	}

	channels, ok := setup["channel"].([]interface{})
	if !ok {
		return "", ""
	}

	for _, ch := range channels {
		chMap, ok := ch.(map[string]interface{})
		if !ok {
			continue
		}
		if getStringValue(chMap, "channel") != channelName {
			continue
		}
		if getIntValue(chMap, "Enable") != 1 {
			continue
		}

		// Check diagnosis enabled.
		general, _ := setup["General"].(map[string]interface{})
		useFunction, _ := general["useFuction"].(map[string]interface{})

		var diagEnabled bool
		if channelName == "Main" {
			diagEnabled, _ = useFunction["diagnosis_main"].(bool)
		} else {
			diagEnabled, _ = useFunction["diagnosis_sub"].(bool)
		}

		if !diagEnabled {
			return "", ""
		}

		assetInfo, _ := chMap["assetInfo"].(map[string]interface{})
		return getStringValue(assetInfo, "name"), getStringValue(assetInfo, "type")
	}

	return "", ""
}

// ---------------------------------------------------------------------------
// GetClosestTimestamp finds the nearest stored timestamp to targetTime.
// ---------------------------------------------------------------------------

func (w *DiagnosisReportWriter) GetClosestTimestamp(mode, assetName string, targetTime time.Time) string {
	start := targetTime.Add(-24 * time.Hour).UTC().Format(time.RFC3339)
	end := targetTime.Add(1 * time.Hour).UTC().Format(time.RFC3339)

	query := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r["_measurement"] == "%s")
			|> filter(fn: (r) => r["asset_name"] == "%s")
			|> filter(fn: (r) => r["data_type"] == "main")
			|> filter(fn: (r) => r["_field"] == "status")
			|> sort(columns: ["_time"], desc: true)
			|> limit(n: 1)
	`, w.Config.Bucket, start, end, mode, assetName)

	results, err := w.influxHandler.ExecuteQuery(query)
	if err != nil {
		log.Printf("[DiagnosisWeekly] timestamp query failed: %v", err)
		return ""
	}
	if len(results) == 0 {
		return ""
	}

	if t, ok := results[0]["_time"].(time.Time); ok {
		localTime := t.Local()
		return localTime.Format("2006-01-02T15:04:05")
	}
	return ""
}

// ---------------------------------------------------------------------------
// GetReportData retrieves main/detail data at a specific timestamp.
// ---------------------------------------------------------------------------

func (w *DiagnosisReportWriter) GetReportData(mode, assetName, timestamp string) map[string]interface{} {
	localTime, err := time.ParseInLocation("2006-01-02T15:04:05", timestamp, time.Local)
	if err != nil {
		log.Printf("[DiagnosisWeekly] timestamp parse error: %v", err)
		return nil
	}

	utcTime := localToUTC(localTime)
	startTime := utcTime.Add(-1 * time.Second).Format("2006-01-02T15:04:05.000000Z")
	endTime := utcTime.Add(1 * time.Second).Format("2006-01-02T15:04:05.000000Z")

	log.Printf("[DiagnosisWeekly] [get_report_data] local: %s -> UTC: %s", timestamp, utcTime.Format(time.RFC3339))

	query := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r["_measurement"] == "%s")
			|> filter(fn: (r) => r["asset_name"] == "%s")
	`, w.Config.Bucket, startTime, endTime, mode, assetName)

	results, err := w.influxHandler.ExecuteQuery(query)
	if err != nil {
		log.Printf("[DiagnosisWeekly] report data query failed: %v", err)
		return nil
	}

	log.Printf("[DiagnosisWeekly] [get_report_data] query result: %d records", len(results))

	// Group by data_type and item_name.
	mainItems := make(map[string]map[string]interface{})
	detailItems := make(map[string]map[string]interface{})

	for _, rec := range results {
		dataType := getStringValue(rec, "data_type")
		itemName := getStringValue(rec, "item_name")

		if dataType != "main" && dataType != "detail" {
			continue
		}
		if itemName == "" {
			continue
		}

		var targetMap map[string]map[string]interface{}
		if dataType == "main" {
			targetMap = mainItems
		} else {
			targetMap = detailItems
		}

		if _, exists := targetMap[itemName]; !exists {
			item := map[string]interface{}{
				"item_name":  itemName,
				"data_type":  dataType,
				"asset_name": getStringValue(rec, "asset_name"),
				"channel":    getStringValue(rec, "channel"),
			}
			if dataType == "detail" {
				item["parent_name"] = getStringValue(rec, "parent_name")
				item["assembly_id"] = getStringValue(rec, "assembly_id")
			}
			targetMap[itemName] = item
		}

		// Add field values.
		fieldName := getStringValue(rec, "_field")
		if fieldName != "" {
			targetMap[itemName][fieldName] = rec["_value"]
		}
	}

	// Convert maps to slices.
	mainList := make([]interface{}, 0, len(mainItems))
	for _, v := range mainItems {
		mainList = append(mainList, v)
	}
	detailList := make([]interface{}, 0, len(detailItems))
	for _, v := range detailItems {
		detailList = append(detailList, v)
	}

	return map[string]interface{}{
		"main":   mainList,
		"detail": detailList,
	}
}

// ---------------------------------------------------------------------------
// GetTrendData retrieves 30-day status trend for a specific item.
// ---------------------------------------------------------------------------

func (w *DiagnosisReportWriter) GetTrendData(mode, assetName, itemName string) []map[string]interface{} {
	query := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: -30d)
			|> filter(fn: (r) => r["_measurement"] == "%s")
			|> filter(fn: (r) => r["asset_name"] == "%s")
			|> filter(fn: (r) => r["data_type"] == "main")
			|> filter(fn: (r) => r["item_name"] == "%s")
			|> filter(fn: (r) => r["_field"] == "status")
			|> sort(columns: ["_time"])
	`, w.Config.Bucket, mode, assetName, itemName)

	results, err := w.influxHandler.ExecuteQuery(query)
	if err != nil {
		log.Printf("[DiagnosisWeekly] trend query failed: %v", err)
		return nil
	}

	trend := make([]map[string]interface{}, 0, len(results))
	for _, rec := range results {
		t := getTimeValue(rec, "_time")
		localTime := t.Local()
		trend = append(trend, map[string]interface{}{
			"timestamp": localTime.Format("2006-01-02T15:04:05"),
			"status":    rec["_value"],
		})
	}

	return trend
}

// ---------------------------------------------------------------------------
// GenerateDiagnosisParquet creates a diagnosis + trend Parquet file.
// ---------------------------------------------------------------------------

func (w *DiagnosisReportWriter) GenerateDiagnosisParquet(endTime time.Time, assetName, channelName string) string {
	log.Printf("[DiagnosisWeekly] ================================================")
	log.Printf("[DiagnosisWeekly] diagnosis report Parquet generation start")
	log.Printf("[DiagnosisWeekly] reference time: %s, asset: %s", endTime.Format(time.RFC3339), assetName)
	log.Printf("[DiagnosisWeekly] ================================================")

	result := map[string]interface{}{
		"asset_name":   assetName,
		"channel_name": channelName,
		"end_time":     endTime.Format(time.RFC3339),
		"generated_at": time.Now().Format(time.RFC3339),
		"diagnosis":    nil,
		"powerquality": nil,
	}

	for _, mode := range w.Config.Modes {
		log.Printf("[DiagnosisWeekly] [%s] querying data...", mode)

		// 1. Find closest timestamp.
		closestTS := w.GetClosestTimestamp(mode, assetName, endTime)
		if closestTS == "" {
			log.Printf("[DiagnosisWeekly] [%s] no stored data found", mode)
			continue
		}
		log.Printf("[DiagnosisWeekly] [%s] closest timestamp: %s", mode, closestTS)

		// 2. Retrieve main/detail data.
		reportData := w.GetReportData(mode, assetName, closestTS)
		if reportData == nil {
			log.Printf("[DiagnosisWeekly] [%s] report data query failed", mode)
			continue
		}

		// 3. Find warning items (status > 1).
		mainItems, _ := reportData["main"].([]interface{})
		var warningItems []map[string]interface{}
		for _, item := range mainItems {
			m, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if getIntValue(m, "status") > 1 {
				warningItems = append(warningItems, m)
			}
		}

		log.Printf("[DiagnosisWeekly] [%s] main: %d, warning(status>1): %d",
			mode, len(mainItems), len(warningItems))

		// 4. Query 30-day trend for each warning item.
		var trends map[string]interface{}
		if len(warningItems) > 0 {
			trends = make(map[string]interface{})
			for _, item := range warningItems {
				itemName := strings.ReplaceAll(getStringValue(item, "item_name"), " ", "")
				trend := w.GetTrendData(mode, assetName, itemName)
				if len(trend) > 0 {
					trends[itemName] = trend
					log.Printf("[DiagnosisWeekly] [%s] %s trend: %d points", mode, itemName, len(trend))
				}
			}
		}

		modeResult := map[string]interface{}{
			"timestamp": closestTS,
			"main":      reportData["main"],
			"detail":    reportData["detail"],
		}
		if len(trends) > 0 {
			modeResult["trends"] = trends
		}

		result[mode] = modeResult
	}

	// 5. Write Parquet file.
	dateStr := endTime.Format("20060102")
	filename := fmt.Sprintf("diagnosis_report_%s_%s.parquet", assetName, dateStr)
	filePath := filepath.Join(w.Config.OutputDir, channelName, filename)

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		log.Printf("[DiagnosisWeekly] mkdir error: %v", err)
		return ""
	}

	// Serialize complex fields to JSON.
	diagJSON, _ := json.Marshal(result["diagnosis"])
	pqJSON, _ := json.Marshal(result["powerquality"])

	row := DiagnosisParquetRow{
		AssetName:    assetName,
		ChannelName:  channelName,
		EndTime:      endTime.Format(time.RFC3339),
		GeneratedAt:  time.Now().Format(time.RFC3339),
		Diagnosis:    string(diagJSON),
		PowerQuality: string(pqJSON),
	}

	f, err := os.Create(filePath)
	if err != nil {
		log.Printf("[DiagnosisWeekly] file create error: %v", err)
		return ""
	}
	defer f.Close()

	pw := parquet.NewGenericWriter[DiagnosisParquetRow](f,
		parquet.Compression(&parquet.Snappy),
	)
	if _, err := pw.Write([]DiagnosisParquetRow{row}); err != nil {
		log.Printf("[DiagnosisWeekly] parquet write error: %v", err)
		return ""
	}
	if err := pw.Close(); err != nil {
		log.Printf("[DiagnosisWeekly] parquet close error: %v", err)
		return ""
	}

	info, _ := os.Stat(filePath)
	sizeKB := float64(0)
	if info != nil {
		sizeKB = float64(info.Size()) / 1024
	}

	log.Printf("[DiagnosisWeekly] diagnosis report Parquet saved: %s", filePath)
	log.Printf("[DiagnosisWeekly]   size: %.1f KB", sizeKB)

	return filePath
}

// ---------------------------------------------------------------------------
// CleanupOldDiagnosisFiles removes old diagnosis/trend Parquet files.
// ---------------------------------------------------------------------------

func (w *DiagnosisReportWriter) CleanupOldDiagnosisFiles(channel string) int {
	var targetDirs []string
	if channel != "" {
		targetDirs = []string{filepath.Join(w.Config.OutputDir, channel)}
	} else {
		entries, err := os.ReadDir(w.Config.OutputDir)
		if err != nil {
			log.Printf("[DiagnosisWeekly] cleanup read dir error: %v", err)
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

	patterns := []string{"diagnosis_report_*.parquet", "trend_training_*.parquet"}

	for _, dir := range targetDirs {
		for _, pattern := range patterns {
			matches, err := filepath.Glob(filepath.Join(dir, pattern))
			if err != nil {
				continue
			}
			for _, file := range matches {
				base := filepath.Base(file)
				parts := strings.Split(strings.TrimSuffix(base, ".parquet"), "_")
				if len(parts) == 0 {
					continue
				}
				dateStr := parts[len(parts)-1]

				fileDate, err := time.Parse("20060102", dateStr)
				if err != nil {
					log.Printf("[DiagnosisWeekly] date parse failed (skip): %s", base)
					continue
				}
				if fileDate.Before(cutoff) {
					if err := os.Remove(file); err != nil {
						log.Printf("[DiagnosisWeekly] delete failed: %s - %v", base, err)
					} else {
						totalDeleted++
						log.Printf("[DiagnosisWeekly] deleted: %s", base)
					}
				}
			}
		}
	}

	if totalDeleted > 0 {
		log.Printf("[DiagnosisWeekly] cleaned up %d files", totalDeleted)
	} else {
		log.Printf("[DiagnosisWeekly] no files to clean up")
	}

	return totalDeleted
}

// ---------------------------------------------------------------------------
// GenerateTrendTrainingParquet creates a flat Parquet for LightGBM training.
// ---------------------------------------------------------------------------

func GenerateTrendTrainingParquet(endTime time.Time, assetName, channelName string, days int, outputDir string) string {
	if outputDir == "" {
		outputDir = "/usr/local/sv500/trendcsv"
	}
	if days <= 0 {
		days = 7 // default 1 week
	}

	startDate := endTime.AddDate(0, 0, -days).Format("2006-01-02 15:04:05")
	endDate := endTime.Format("2006-01-02 15:04:05")

	log.Printf("[DiagnosisWeekly] trend training data collection: %s (%s ~ %s)", assetName, startDate, endDate)

	// Query trend data via the diagnosis API.
	endpoint := fmt.Sprintf("collectAllGroupTrends?name=%s&start=%s&end=%s", assetName, startDate, endDate)
	result, err := fetchAPI(endpoint)
	if err != nil {
		log.Printf("[DiagnosisWeekly] trend data collection failed: %v", err)
		return ""
	}

	success, _ := result["success"].(bool)
	if !success {
		log.Printf("[DiagnosisWeekly] trend data collection failed: %s", assetName)
		return ""
	}

	dataMap, _ := result["data"].(map[string]interface{})
	if dataMap == nil {
		log.Printf("[DiagnosisWeekly] no trend data: %s", assetName)
		return ""
	}

	// Build flat rows.
	var rows []TrendTrainingRow
	for groupName, groupResultRaw := range dataMap {
		groupResult, ok := groupResultRaw.(map[string]interface{})
		if !ok {
			continue
		}
		groupSuccess, _ := groupResult["success"].(bool)
		if !groupSuccess {
			continue
		}
		groupData, _ := groupResult["data"].(map[string]interface{})
		if groupData == nil {
			continue
		}

		for title, paramRaw := range groupData {
			param, ok := paramRaw.(map[string]interface{})
			if !ok {
				continue
			}
			unit := getStringValue(param, "Unit")
			dataPoints, _ := param["data"].([]interface{})

			for _, pointRaw := range dataPoints {
				point, ok := pointRaw.(map[string]interface{})
				if !ok {
					continue
				}
				rows = append(rows, TrendTrainingRow{
					Asset:       assetName,
					Channel:     channelName,
					Group:       groupName,
					Title:       title,
					Unit:        unit,
					Timestamp:   getStringValue(point, "XAxis"),
					Value:       getFloat64Value(point, "YAxis"),
					Status:      int64(getIntValue(point, "Status")),
					IStatus:     int64(getIntValue(point, "iStatus")),
					CollectedAt: endTime.Format(time.RFC3339),
				})
			}
		}
	}

	if len(rows) == 0 {
		log.Printf("[DiagnosisWeekly] no trend training data to save: %s", assetName)
		return ""
	}

	// Write Parquet.
	dateStr := endTime.Format("20060102")
	filename := fmt.Sprintf("trend_training_%s_%s.parquet", assetName, dateStr)
	filePath := filepath.Join(outputDir, channelName, filename)

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		log.Printf("[DiagnosisWeekly] mkdir error: %v", err)
		return ""
	}

	f, err := os.Create(filePath)
	if err != nil {
		log.Printf("[DiagnosisWeekly] file create error: %v", err)
		return ""
	}
	defer f.Close()

	pw := parquet.NewGenericWriter[TrendTrainingRow](f,
		parquet.Compression(&parquet.Snappy),
	)
	if _, err := pw.Write(rows); err != nil {
		log.Printf("[DiagnosisWeekly] parquet write error: %v", err)
		return ""
	}
	if err := pw.Close(); err != nil {
		log.Printf("[DiagnosisWeekly] parquet close error: %v", err)
		return ""
	}

	info, _ := os.Stat(filePath)
	sizeKB := float64(0)
	if info != nil {
		sizeKB = float64(info.Size()) / 1024
	}

	log.Printf("[DiagnosisWeekly] trend training Parquet saved: %s (%.1f KB), %d rows",
		filePath, sizeKB, len(rows))

	return filePath
}

// ---------------------------------------------------------------------------
// Entry points
// ---------------------------------------------------------------------------

// RunDiagnosisReportJob generates a diagnosis report Parquet for the given channel.
func RunDiagnosisReportJob(endTime time.Time, channelName string) string {
	assetName, _ := GetAssetInfo(channelName)
	if assetName == "" {
		log.Printf("[DiagnosisWeekly] %s channel: no asset info or diagnosis disabled", channelName)
		return ""
	}

	log.Printf("[DiagnosisWeekly] diagnosis report Parquet job started")

	writer := NewDiagnosisReportWriter(DefaultDiagnosisReportConfig())
	filePath := writer.GenerateDiagnosisParquet(endTime, assetName, channelName)

	if filePath != "" {
		log.Printf("[DiagnosisWeekly] diagnosis report generated: %s", filePath)
	} else {
		log.Printf("[DiagnosisWeekly] diagnosis report generation failed or no data")
	}

	writer.CleanupOldDiagnosisFiles(channelName)
	return filePath
}

// RunTrendTrainingJob generates a trend training Parquet for the given channel.
func RunTrendTrainingJob(endTime time.Time, channelName string) string {
	assetName, _ := GetAssetInfo(channelName)
	if assetName == "" {
		log.Printf("[DiagnosisWeekly] %s channel: no asset info or diagnosis disabled", channelName)
		return ""
	}

	return GenerateTrendTrainingParquet(endTime, assetName, channelName, 7, "")
}

// RunTrendTrainingManual is for initial training after commissioning.
// Collects the most recent `days` of data.
func RunTrendTrainingManual(channelName string, days int) string {
	assetName, _ := GetAssetInfo(channelName)
	if assetName == "" {
		log.Printf("[DiagnosisWeekly] %s channel: no asset info or diagnosis disabled", channelName)
		return ""
	}

	if days <= 0 {
		days = 14
	}

	endTime := time.Now()
	log.Printf("[DiagnosisWeekly] [manual] trend training data: %s/%s, last %d days",
		channelName, assetName, days)

	return GenerateTrendTrainingParquet(endTime, assetName, channelName, days, "/usr/local/sv500/train")
}
