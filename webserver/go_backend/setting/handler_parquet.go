package setting

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	reportsDir = "/usr/local/sv500/reports"
	trendDir   = "/usr/local/sv500/trendcsv"
)

func (h *Handler) registerParquetRoutes(rg *gin.RouterGroup) {
	rg.GET("/download", h.downloadSetupJSON)
	rg.GET("/backup/parquet/report/list", h.listReportFiles)
	rg.GET("/backup/parquet/report/download", h.downloadReportParquet)
	rg.GET("/backup/parquet/report/download-all", h.downloadAllReports)
	rg.GET("/backup/parquet/trend/list", h.listTrendFiles)
	rg.GET("/backup/parquet/trend/download", h.downloadTrendFile)
	rg.GET("/backup/parquet/trend/download-all", h.downloadAllTrends)
}

// GET /setting/download — download setup.json
func (h *Handler) downloadSetupJSON(c *gin.Context) {
	path := filepath.Join(h.deps.Config.ConfigDir, "setup.json")
	if _, err := os.Stat(path); err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "error": "setting file not found"})
		return
	}
	c.Header("Content-Disposition", `attachment; filename="setup.json"`)
	c.File(path)
}

// GET /setting/backup/parquet/report/list
func (h *Handler) listReportFiles(c *gin.Context) {
	if _, err := os.Stat(reportsDir); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Reports directory not found"})
		return
	}

	chEntries, _ := os.ReadDir(reportsDir)
	result := map[string][]gin.H{}
	en50160Re := regexp.MustCompile(`^en50160_weekly_(.+)\.parquet$`)
	diagRe := func(date string) *regexp.Regexp {
		return regexp.MustCompile(`^diagnosis_report_.+_` + regexp.QuoteMeta(date) + `\.parquet$`)
	}

	for _, ch := range chEntries {
		if !ch.IsDir() {
			continue
		}
		chPath := filepath.Join(reportsDir, ch.Name())
		files, _ := os.ReadDir(chPath)

		// Collect en50160 dates and diagnosis filenames per date.
		dates := []string{}
		diagByDate := map[string]string{}
		en50160ByDate := map[string]string{}

		for _, f := range files {
			if f.IsDir() {
				continue
			}
			if m := en50160Re.FindStringSubmatch(f.Name()); m != nil {
				date := m[1]
				dates = append(dates, date)
				en50160ByDate[date] = f.Name()
			}
		}
		sort.Strings(dates)

		for _, d := range dates {
			re := diagRe(d)
			for _, f := range files {
				if f.IsDir() {
					continue
				}
				if re.MatchString(f.Name()) {
					diagByDate[d] = f.Name()
					break
				}
			}
		}

		var items []gin.H
		for _, d := range dates {
			items = append(items, gin.H{
				"date":      d,
				"en50160":   en50160ByDate[d],
				"diagnosis": diagByDate[d],
			})
		}
		if len(items) > 0 {
			result[ch.Name()] = items
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// GET /setting/backup/parquet/report/download?channel=&date=
func (h *Handler) downloadReportParquet(c *gin.Context) {
	channel := c.Query("channel")
	date := c.Query("date")
	chDir := filepath.Join(reportsDir, channel)
	if _, err := os.Stat(chDir); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Channel directory not found: " + channel})
		return
	}

	en50160 := filepath.Join(chDir, fmt.Sprintf("en50160_weekly_%s.parquet", date))
	if _, err := os.Stat(en50160); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "EN50160 file not found: " + date})
		return
	}

	entries, _ := os.ReadDir(chDir)
	re := regexp.MustCompile(`^diagnosis_report_.+_` + regexp.QuoteMeta(date) + `\.parquet$`)
	var diagFiles []string
	for _, e := range entries {
		if !e.IsDir() && re.MatchString(e.Name()) {
			diagFiles = append(diagFiles, filepath.Join(chDir, e.Name()))
		}
	}

	// EN50160 only → stream the single file.
	if len(diagFiles) == 0 {
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="en50160_weekly_%s.parquet"`, date))
		c.File(en50160)
		return
	}

	// Combine into tar.gz.
	timestamp := time.Now().Format("20060102_150405")
	tempDir, err := os.MkdirTemp("", "report_")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	defer os.RemoveAll(tempDir)

	collect := filepath.Join(tempDir, fmt.Sprintf("report_%s_%s", channel, date))
	os.MkdirAll(collect, 0755)
	copyFile(en50160, filepath.Join(collect, filepath.Base(en50160)))
	for _, df := range diagFiles {
		copyFile(df, filepath.Join(collect, filepath.Base(df)))
	}

	backupName := fmt.Sprintf("report_%s_%s_%s", channel, date, timestamp)
	tarFile := filepath.Join(os.TempDir(), backupName+".tar.gz")
	out, err := exec.Command("tar", "--ignore-failed-read", "-czf", tarFile, "-C", tempDir, filepath.Base(collect)).CombinedOutput()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "tar failed: " + strings.TrimSpace(string(out))})
		return
	}
	defer os.Remove(tarFile)

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.tar.gz"`, backupName))
	c.Header("Content-Type", "application/gzip")
	c.File(tarFile)
}

// GET /setting/backup/parquet/report/download-all
func (h *Handler) downloadAllReports(c *gin.Context) {
	if _, err := os.Stat(reportsDir); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "No report files found"})
		return
	}

	timestamp := time.Now().Format("20060102_150405")
	tempDir, err := os.MkdirTemp("", "reports_")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	defer os.RemoveAll(tempDir)

	collect := filepath.Join(tempDir, "reports")
	channels, _ := os.ReadDir(reportsDir)
	hasAny := false
	for _, ch := range channels {
		if !ch.IsDir() {
			continue
		}
		src := filepath.Join(reportsDir, ch.Name())
		files, _ := os.ReadDir(src)
		copied := false
		for _, f := range files {
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".parquet") {
				continue
			}
			if !copied {
				os.MkdirAll(filepath.Join(collect, ch.Name()), 0755)
				copied = true
			}
			copyFile(filepath.Join(src, f.Name()), filepath.Join(collect, ch.Name(), f.Name()))
			hasAny = true
		}
	}
	if !hasAny {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "No report files found"})
		return
	}

	backupName := "backup_report_all_" + timestamp
	tarFile := filepath.Join(os.TempDir(), backupName+".tar.gz")
	out, err := exec.Command("tar", "--ignore-failed-read", "-czf", tarFile, "-C", tempDir, "reports").CombinedOutput()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "tar failed: " + strings.TrimSpace(string(out))})
		return
	}
	defer os.Remove(tarFile)

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.tar.gz"`, backupName))
	c.Header("Content-Type", "application/gzip")
	c.File(tarFile)
}

// GET /setting/backup/parquet/trend/list
func (h *Handler) listTrendFiles(c *gin.Context) {
	if _, err := os.Stat(trendDir); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Trend directory not found"})
		return
	}
	channels, _ := os.ReadDir(trendDir)
	result := map[string][]string{}
	for _, ch := range channels {
		if !ch.IsDir() {
			continue
		}
		files, _ := os.ReadDir(filepath.Join(trendDir, ch.Name()))
		var names []string
		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".parquet") {
				names = append(names, f.Name())
			}
		}
		if len(names) > 0 {
			sort.Sort(sort.Reverse(sort.StringSlice(names)))
			result[ch.Name()] = names
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// GET /setting/backup/parquet/trend/download?filename=&channel=
func (h *Handler) downloadTrendFile(c *gin.Context) {
	filename := c.Query("filename")
	channel := c.Query("channel")

	var path string
	if channel != "" {
		path = filepath.Join(trendDir, channel, filename)
	} else {
		path = filepath.Join(trendDir, filename)
	}
	if _, err := os.Stat(path); err != nil || !strings.HasSuffix(filename, ".parquet") {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "File not found: " + filename})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.File(path)
}

// GET /setting/backup/parquet/trend/download-all
func (h *Handler) downloadAllTrends(c *gin.Context) {
	if _, err := os.Stat(trendDir); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "No trend files found"})
		return
	}

	timestamp := time.Now().Format("20060102_150405")
	tempDir, err := os.MkdirTemp("", "trend_")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	defer os.RemoveAll(tempDir)

	collect := filepath.Join(tempDir, "trendcsv")
	os.MkdirAll(collect, 0755)

	channels, _ := os.ReadDir(trendDir)
	hasAny := false
	for _, ch := range channels {
		if !ch.IsDir() {
			continue
		}
		src := filepath.Join(trendDir, ch.Name())
		files, _ := os.ReadDir(src)
		if len(files) == 0 {
			continue
		}
		dst := filepath.Join(collect, ch.Name())
		os.MkdirAll(dst, 0755)
		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".parquet") {
				copyFile(filepath.Join(src, f.Name()), filepath.Join(dst, f.Name()))
				hasAny = true
			}
		}
	}
	if !hasAny {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "No trend files found"})
		return
	}

	backupName := "backup_trend_all_" + timestamp
	tarFile := filepath.Join(os.TempDir(), backupName+".tar.gz")
	out, err := exec.Command("tar", "--ignore-failed-read", "-czf", tarFile, "-C", tempDir, "trendcsv").CombinedOutput()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "tar failed: " + strings.TrimSpace(string(out))})
		return
	}
	defer os.Remove(tarFile)

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.tar.gz"`, backupName))
	c.Header("Content-Type", "application/gzip")
	c.File(tarFile)
}
