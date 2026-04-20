// Package report ports the EN50160-only endpoints from python_backend/routes/report.py.
// Excluded:
//   - /generate, /downloadDiagnosisReport, /downloadWeeklyReport (Word doc)
//   - /lastReportData, /reportTimes, /reportDataByTime, /status_trend (InfluxDB)
//   - /getReportDiagnosis (diagnosis parquet, not needed for this build)
package report

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"serverGO/en50160"
	"serverGO/infra"
)

// Handler wraps shared deps plus a processor that we reuse across requests.
type Handler struct {
	deps       *infra.Dependencies
	processor  *en50160.Processor
	reportsDir string
}

// RegisterRoutes mounts the report endpoints under the supplied group.
func RegisterRoutes(rg *gin.RouterGroup, deps *infra.Dependencies) {
	proc := en50160.NewProcessor(nil) // defaults: /usr/local/sv500/reports, 60Hz, 22900V
	h := &Handler{
		deps:       deps,
		processor:  proc,
		reportsDir: proc.Config.OutputDir,
	}

	rg.GET("/week/:channel/:filename", h.getWeekly)
	rg.GET("/list/:channel", h.getFileList)
	rg.GET("/weeklyReportData/:channel/:date", h.getWeeklyReportData)
}

// GET /report/week/:channel/:filename — EN50160 chart bundle for the given file.
func (h *Handler) getWeekly(c *gin.Context) {
	channel := c.Param("channel")
	filename := c.Param("filename")
	ctx := context.Background()

	h.setLimit(ctx, channel)

	data, err := h.processor.GetAllChartData(channel, filename)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GET /report/list/:channel — sorted date list for the channel's reports folder.
func (h *Handler) getFileList(c *gin.Context) {
	channel := c.Param("channel")
	dates, err := h.processor.ListFiles(channel)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": dates})
}

// GET /report/weeklyReportData/:channel/:date — EN50160 chart bundle for one week.
func (h *Handler) getWeeklyReportData(c *gin.Context) {
	channel := c.Param("channel")
	date := c.Param("date")

	start, end, err := getWeekRangeFromDate(date)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "msg": err.Error()})
		return
	}
	reportPeriod := fmt.Sprintf("%s ~ %s", start.Format("2006-01-02"), end.Format("2006-01-02"))

	ctx := context.Background()
	h.setLimit(ctx, channel)

	var en50160Data interface{}
	enPath := filepath.Join(h.reportsDir, channel, fmt.Sprintf("en50160_weekly_%s.parquet", date))
	if _, err := os.Stat(enPath); err == nil {
		if data, err := h.processor.GetAllChartData(channel, fmt.Sprintf("en50160_weekly_%s.parquet", date)); err == nil {
			en50160Data = data
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"channel":      channel,
			"date":         date,
			"reportPeriod": reportPeriod,
			"en50160":      en50160Data,
		},
	})
}

// setLimit ports Python set_limit — applies nominal voltage/current/frequency from Redis
// Equipment.ChannelData (fallback: System.setup's ptInfo/ctInfo).
func (h *Handler) setLimit(ctx context.Context, channel string) {
	client := h.deps.Redis.Client0
	chKey := "main"
	if !strings.EqualFold(channel, "Main") {
		chKey = "sub"
	}

	// Primary: Equipment.ChannelData.
	if raw, err := client.HGet(ctx, "Equipment", "ChannelData").Result(); err == nil && raw != "" {
		var chDict map[string]interface{}
		if json.Unmarshal([]byte(raw), &chDict) == nil {
			if sub, ok := chDict[chKey].(map[string]interface{}); ok {
				v := toFloat(sub["RatedVoltage"])
				cur := toFloat(sub["RatedCurrent"])
				f := toFloat(sub["RatedFrequency"])
				h.processor.SetLimits(&v, &cur, &f)
				return
			}
		}
	}

	// Fallback: System.setup channel array.
	if raw, err := client.HGet(ctx, "System", "setup").Result(); err == nil && raw != "" {
		var setup map[string]interface{}
		if json.Unmarshal([]byte(raw), &setup) == nil {
			if arr, ok := setup["channel"].([]interface{}); ok {
				for _, item := range arr {
					chMap, _ := item.(map[string]interface{})
					name, _ := chMap["channel"].(string)
					if !strings.EqualFold(name, channel) {
						continue
					}
					pt, _ := chMap["ptInfo"].(map[string]interface{})
					ct, _ := chMap["ctInfo"].(map[string]interface{})
					v := toFloat(pt["vnorminal"])
					cur := toFloat(ct["inorminal"])
					f := toFloat(pt["linefrequency"])
					if f == 0 {
						f = toFloat(pt["fnominal"])
					}
					h.processor.SetLimits(&v, &cur, &f)
					return
				}
			}
		}
	}
}

// getWeekRangeFromDate mirrors Python get_week_range_from_date — returns the Mon-Sun range
// for the week containing the given YYYYMMDD date.
func getWeekRangeFromDate(dateStr string) (time.Time, time.Time, error) {
	target, err := time.ParseInLocation("20060102", dateStr, time.Local)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	// Python uses Mon=0 … Sun=6. Go uses Sun=0 … Sat=6 — convert.
	weekday := (int(target.Weekday()) + 6) % 7
	monday := target.AddDate(0, 0, -weekday)
	sunday := monday.AddDate(0, 0, 6)
	return monday, sunday, nil
}

func toFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		var f float64
		fmt.Sscanf(val, "%f", &f)
		return f
	}
	return 0
}
