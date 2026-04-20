package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"serverGO/metermap"

	_ "modernc.org/sqlite"
)

const pageSizeAlarm = 10

var conditionMap = map[int]string{0: "EQUAL", 1: "UNDER", 2: "OVER"}
var eventMap = map[int]string{1: "OCCURRED", 0: "CLEARED"}

type AlarmLogSearch struct {
	StartDate string `json:"StartDate" binding:"required"`
	EndDate   string `json:"EndDate" binding:"required"`
}

func (h *Handler) registerAlarmLogRoutes(rg *gin.RouterGroup) {
	rg.GET("/getAlarmParms/:channel", h.getAlarmParms)
	rg.GET("/getAlarmLog/:page", h.getAlarmLog)
	rg.POST("/getAlarmLogSearch/:page", h.getAlarmLogSearch)
	rg.GET("/getRecentAlarmLog", h.getRecentAlarmLog)
}

func (h *Handler) getAlarmParms(c *gin.Context) {
	channel := c.Param("channel")
	ctx := context.Background()

	setupJSON, err := h.deps.Redis.Client0.HGet(ctx, "System", "setup").Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}

	var setting map[string]interface{}
	json.Unmarshal([]byte(setupJSON), &setting)

	var almSetup map[string]interface{}
	if channels, ok := setting["channel"].([]interface{}); ok {
		for _, ch := range channels {
			chMap, ok := ch.(map[string]interface{})
			if !ok {
				continue
			}
			if chMap["channel"] == channel {
				almSetup, _ = chMap["alarm"].(map[string]interface{})
				break
			}
		}
	}

	if len(almSetup) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}

	paramList := []gin.H{}
	seen := map[int]bool{}
	for i := 1; i <= 32; i++ {
		entry, ok := almSetup[fmt.Sprintf("%d", i)].(map[string]interface{})
		if !ok {
			continue
		}
		idx := 0
		if v, ok := entry["chan"].(float64); ok {
			idx = int(v)
		}
		if idx != 0 && !seen[idx] && idx < len(metermap.ParameterOptions) {
			paramList = append(paramList, gin.H{"id": idx, "label": metermap.ParameterOptions[idx]})
			seen[idx] = true
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": paramList})
}

func (h *Handler) openAlarmDB() (*sql.DB, error) {
	dbPath := filepath.Join(h.deps.Config.ConfigDir, "alarm.db")
	return sql.Open("sqlite", dbPath)
}

func formatAlarmRow(row map[string]interface{}) map[string]interface{} {
	tsMs, _ := row["ts_ms"].(int64)
	ts := time.Unix(tsMs/1000, (tsMs%1000)*int64(time.Millisecond))
	tsStr := ts.Format("2006-01-02 15:04:05")

	chanIdx, _ := row["chan"].(int64)
	chanText := fmt.Sprintf("CH%d", chanIdx)
	if int(chanIdx) >= 0 && int(chanIdx) < len(metermap.ParameterOptions) {
		chanText = metermap.ParameterOptions[int(chanIdx)]
	}

	condVal, _ := row["condition"].(int64)
	condText := conditionMap[int(condVal)]
	if condText == "" {
		condText = fmt.Sprintf("%d", condVal)
	}

	threshold, _ := row["threshold"].(float64)
	criteria := fmt.Sprintf("%s %s %g", chanText, condText, threshold)

	eventVal, _ := row["event"].(int64)
	statusText := eventMap[int(eventVal)]
	if statusText == "" {
		statusText = fmt.Sprintf("%d", eventVal)
	}

	value, _ := row["value"].(float64)

	return map[string]interface{}{
		"ts_formatted": tsStr,
		"criteria":     criteria,
		"status":       statusText,
		"value":        math.Round(value*100) / 100,
	}
}

func queryAlarmLog(db *sql.DB, page int, startMs, endMs *int64) ([]map[string]interface{}, error) {
	offset := (page - 1) * pageSizeAlarm
	limit := pageSizeAlarm + 1

	var rows *sql.Rows
	var err error

	if startMs != nil && endMs != nil {
		rows, err = db.Query(
			"SELECT * FROM alarm_log WHERE ts_ms >= ? AND ts_ms <= ? ORDER BY ts_ms DESC LIMIT ? OFFSET ?",
			*startMs, *endMs, limit, offset)
	} else {
		rows, err = db.Query(
			"SELECT * FROM alarm_log ORDER BY ts_ms DESC LIMIT ? OFFSET ?",
			limit, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	var result []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range values {
			ptrs[i] = &values[i]
		}
		rows.Scan(ptrs...)

		row := make(map[string]interface{})
		for i, col := range cols {
			row[col] = values[i]
		}
		result = append(result, formatAlarmRow(row))
	}
	return result, nil
}

func (h *Handler) getAlarmLog(c *gin.Context) {
	page, _ := strconv.Atoi(c.Param("page"))
	if page < 1 {
		page = 1
	}

	db, err := h.openAlarmDB()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error(), "data": []interface{}{}, "page": page, "hasNext": false, "hasPrev": false})
		return
	}
	defer db.Close()

	rows, err := queryAlarmLog(db, page, nil, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error(), "data": []interface{}{}, "page": page, "hasNext": false, "hasPrev": false})
		return
	}

	hasNext := len(rows) > pageSizeAlarm
	if hasNext {
		rows = rows[:pageSizeAlarm]
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    rows,
		"page":    page,
		"hasNext": hasNext,
		"hasPrev": page > 1,
	})
}

func (h *Handler) getAlarmLogSearch(c *gin.Context) {
	page, _ := strconv.Atoi(c.Param("page"))
	if page < 1 {
		page = 1
	}

	var req AlarmLogSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error(), "data": []interface{}{}, "page": page, "hasNext": false, "hasPrev": false})
		return
	}

	startTime, _ := time.Parse(time.RFC3339, req.StartDate)
	endTime, _ := time.Parse(time.RFC3339, req.EndDate)
	startMs := startTime.UnixMilli()
	endMs := endTime.UnixMilli()

	db, err := h.openAlarmDB()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error(), "data": []interface{}{}, "page": page, "hasNext": false, "hasPrev": false})
		return
	}
	defer db.Close()

	rows, err := queryAlarmLog(db, page, &startMs, &endMs)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error(), "data": []interface{}{}, "page": page, "hasNext": false, "hasPrev": false})
		return
	}

	hasNext := len(rows) > pageSizeAlarm
	if hasNext {
		rows = rows[:pageSizeAlarm]
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    rows,
		"page":    page,
		"hasNext": hasNext,
		"hasPrev": page > 1,
	})
}

func (h *Handler) getRecentAlarmLog(c *gin.Context) {
	db, err := h.openAlarmDB()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error(), "data": []interface{}{}})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM alarm_log ORDER BY ts_ms DESC LIMIT 5")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": err.Error(), "data": []interface{}{}})
		return
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	var result []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range values {
			ptrs[i] = &values[i]
		}
		rows.Scan(ptrs...)

		row := make(map[string]interface{})
		for i, col := range cols {
			row[col] = values[i]
		}
		result = append(result, formatAlarmRow(row))
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}
