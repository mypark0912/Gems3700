package processors

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"sv500_core/handlers"
)

// ArchiveManager manages PQ 10-minute data archiving and retrieval.
type ArchiveManager struct {
	InfluxHandler   *handlers.InfluxDBHandler
	ArchiveBasePath string
}

// NewArchiveManager creates a new ArchiveManager.
func NewArchiveManager(influxHandler *handlers.InfluxDBHandler, archiveBasePath string) *ArchiveManager {
	if archiveBasePath == "" {
		archiveBasePath = "/data/pq_archive"
	}

	// Ensure base directory exists.
	if err := os.MkdirAll(archiveBasePath, 0755); err != nil {
		log.Printf("[Archive] failed to create base path %s: %v", archiveBasePath, err)
	}

	return &ArchiveManager{
		InfluxHandler:   influxHandler,
		ArchiveBasePath: archiveBasePath,
	}
}

// ArchiveWeekData archives a specific week's data from InfluxDB to a Parquet file.
func (a *ArchiveManager) ArchiveWeekData(year, week int, channel string) bool {
	// Calculate the start/end dates for the ISO week.
	startDate := isoWeekStart(year, week)
	endDate := startDate.AddDate(0, 0, 7)

	log.Printf("[Archive] archiving: %d-W%02d %s (%s ~ %s)",
		year, week, channel,
		startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	// Query InfluxDB for the week's data.
	query := fmt.Sprintf(`
from(bucket: "ntek_pq_raw")
    |> range(start: %sZ, stop: %sZ)
    |> filter(fn: (r) => r._measurement == "en10min")
    |> filter(fn: (r) => r.channel == "%s")
    |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
`,
		startDate.Format(time.RFC3339),
		endDate.Format(time.RFC3339),
		channel,
	)

	data, err := a.InfluxHandler.ExecuteQuery(query)
	if err != nil {
		log.Printf("[Archive] query failed: %d-W%02d %s - %v", year, week, channel, err)
		return false
	}

	if len(data) == 0 {
		log.Printf("[Archive] no data: %d-W%02d %s", year, week, channel)
		return false
	}

	log.Printf("[Archive]   data points: %d", len(data))

	// Filter out internal InfluxDB metadata columns (starting with '_').
	filteredData := filterMetadataColumns(data)

	// Create year directory.
	yearDir := filepath.Join(a.ArchiveBasePath, fmt.Sprintf("%d", year))
	if err := os.MkdirAll(yearDir, 0755); err != nil {
		log.Printf("[Archive] failed to create year dir %s: %v", yearDir, err)
		return false
	}

	// Write data as JSON (since Go doesn't have a built-in Parquet writer
	// without additional dependencies; the parquet-go library in go.mod
	// can be used for a full implementation).
	parquetFile := filepath.Join(yearDir, fmt.Sprintf("W%02d_%s.json", week, channel))

	if err := writeArchiveData(parquetFile, filteredData); err != nil {
		log.Printf("[Archive] save failed: %s - %v", parquetFile, err)
		return false
	}

	info, err := os.Stat(parquetFile)
	if err == nil {
		sizeMB := float64(info.Size()) / (1024 * 1024)
		log.Printf("[Archive] saved: %s (%.2f MB)", filepath.Base(parquetFile), sizeMB)
	}

	return true
}

// LoadWeekData loads archived week data from file.
func (a *ArchiveManager) LoadWeekData(year, week int, channel string) ([]map[string]interface{}, error) {
	archiveFile := filepath.Join(a.ArchiveBasePath, fmt.Sprintf("%d", year),
		fmt.Sprintf("W%02d_%s.json", week, channel))

	if _, err := os.Stat(archiveFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", archiveFile)
	}

	data, err := readArchiveData(archiveFile)
	if err != nil {
		return nil, fmt.Errorf("load failed: %w", err)
	}

	log.Printf("[Archive] loaded: %s (%d records)", filepath.Base(archiveFile), len(data))
	return data, nil
}

// LoadColumns loads specific columns from archived data.
func (a *ArchiveManager) LoadColumns(year, week int, channel string, columns []string) ([]map[string]interface{}, error) {
	data, err := a.LoadWeekData(year, week, channel)
	if err != nil {
		return nil, err
	}

	// Filter to only requested columns.
	colSet := make(map[string]struct{}, len(columns))
	for _, c := range columns {
		colSet[c] = struct{}{}
	}

	filtered := make([]map[string]interface{}, len(data))
	for i, row := range data {
		newRow := make(map[string]interface{})
		for k, v := range row {
			if _, ok := colSet[k]; ok {
				newRow[k] = v
			}
		}
		filtered[i] = newRow
	}

	return filtered, nil
}

// CleanupOldFiles removes archive files older than the given retention period.
func (a *ArchiveManager) CleanupOldFiles(retentionDays int) error {
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	cutoffYear := cutoff.Year()

	entries, err := os.ReadDir(a.ArchiveBasePath)
	if err != nil {
		return fmt.Errorf("failed to read archive base: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		var dirYear int
		if _, err := fmt.Sscanf(entry.Name(), "%d", &dirYear); err != nil {
			continue
		}

		// Remove entire year directories older than cutoff.
		if dirYear < cutoffYear {
			dirPath := filepath.Join(a.ArchiveBasePath, entry.Name())
			log.Printf("[Archive] removing old archive directory: %s", dirPath)
			if err := os.RemoveAll(dirPath); err != nil {
				log.Printf("[Archive] failed to remove %s: %v", dirPath, err)
			}
			continue
		}

		// For the cutoff year, check individual files.
		if dirYear == cutoffYear {
			yearDir := filepath.Join(a.ArchiveBasePath, entry.Name())
			a.cleanupYearDir(yearDir, cutoff)
		}
	}

	return nil
}

// ArchiveDataFiles archives data files from source directory to archive directory.
func (a *ArchiveManager) ArchiveDataFiles(srcDir, channel string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("failed to read source dir %s: %w", srcDir, err)
	}

	now := time.Now()
	yearDir := filepath.Join(a.ArchiveBasePath, fmt.Sprintf("%d", now.Year()))
	if err := os.MkdirAll(yearDir, 0755); err != nil {
		return fmt.Errorf("failed to create year dir: %w", err)
	}

	archived := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(yearDir, fmt.Sprintf("%s_%s", channel, entry.Name()))

		if err := moveFile(srcPath, dstPath); err != nil {
			log.Printf("[Archive] failed to archive %s: %v", entry.Name(), err)
			continue
		}
		archived++
	}

	if archived > 0 {
		log.Printf("[Archive] archived %d files from %s", archived, srcDir)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// isoWeekStart returns the Monday of the given ISO year/week.
func isoWeekStart(year, week int) time.Time {
	// January 4 is always in week 1 of the ISO year.
	jan4 := time.Date(year, 1, 4, 0, 0, 0, 0, time.UTC)
	// Find the Monday of week 1.
	daysSinceMonday := (int(jan4.Weekday()) + 6) % 7
	week1Monday := jan4.AddDate(0, 0, -daysSinceMonday)
	// Add (week-1) weeks to get the Monday of the desired week.
	return week1Monday.AddDate(0, 0, (week-1)*7)
}

// filterMetadataColumns removes columns starting with '_' from query results.
func filterMetadataColumns(data []map[string]interface{}) []map[string]interface{} {
	filtered := make([]map[string]interface{}, len(data))
	for i, row := range data {
		newRow := make(map[string]interface{})
		for k, v := range row {
			if !strings.HasPrefix(k, "_") {
				newRow[k] = v
			}
		}
		filtered[i] = newRow
	}
	return filtered
}

// writeArchiveData writes data to a JSON archive file.
func writeArchiveData(filePath string, data []map[string]interface{}) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write as JSON lines for efficient streaming reads.
	for _, row := range data {
		line := mapToJSONLine(row)
		if _, err := f.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return nil
}

// readArchiveData reads data from a JSON archive file.
func readArchiveData(filePath string) ([]map[string]interface{}, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	var result []map[string]interface{}

	for _, line := range lines {
		if line == "" {
			continue
		}
		row, err := jsonLineToMap(line)
		if err != nil {
			log.Printf("[Archive] skipping malformed line: %v", err)
			continue
		}
		result = append(result, row)
	}

	return result, nil
}

// mapToJSONLine converts a map to a JSON string.
func mapToJSONLine(m map[string]interface{}) string {
	b, _ := json.Marshal(m)
	return string(b)
}

// jsonLineToMap converts a JSON string to a map.
func jsonLineToMap(line string) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(line), &m)
	return m, err
}

// moveFile moves a file from src to dst (copy + delete for cross-device moves).
func moveFile(src, dst string) error {
	// Try rename first (same filesystem).
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	// Fall back to copy + delete.
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.WriteFile(dst, data, 0644); err != nil {
		return err
	}
	return os.Remove(src)
}

// cleanupYearDir removes individual archive files older than cutoff within a year directory.
func (a *ArchiveManager) cleanupYearDir(yearDir string, cutoff time.Time) {
	entries, err := os.ReadDir(yearDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			filePath := filepath.Join(yearDir, entry.Name())
			log.Printf("[Archive] removing old file: %s", filePath)
			os.Remove(filePath)
		}
	}

	// Remove year directory if empty.
	remaining, _ := os.ReadDir(yearDir)
	if len(remaining) == 0 {
		os.Remove(yearDir)
	}
}
