package setting

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

func (h *Handler) registerBearingRoutes(rg *gin.RouterGroup) {
	rg.GET("/checkBearing", h.checkBearing)
	rg.POST("/uploadBearing", h.uploadBearing)
}

type bearingRow struct {
	Name string  `json:"Name"`
	BPFO float64 `json:"BPFO"`
	BPFI float64 `json:"BPFI"`
	BSF  float64 `json:"BSF"`
	FTF  float64 `json:"FTF"`
}

// GET /setting/checkBearing
func (h *Handler) checkBearing(c *gin.Context) {
	dbPath := filepath.Join(h.deps.Config.ConfigDir, "bearings.db")
	db, err := openBearingDB(dbPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "msg": "Database error: " + err.Error()})
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT Name, BPFO, BPFI, BSF, FTF FROM bearings ORDER BY Name`)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "msg": err.Error()})
		return
	}
	defer rows.Close()

	var result []bearingRow
	for rows.Next() {
		var r bearingRow
		if err := rows.Scan(&r.Name, &r.BPFO, &r.BPFI, &r.BSF, &r.FTF); err != nil {
			c.JSON(http.StatusOK, gin.H{"passOK": "0", "msg": err.Error()})
			return
		}
		result = append(result, r)
	}

	if len(result) == 0 {
		c.JSON(http.StatusOK, gin.H{"passOK": "1", "data": []bearingRow{}, "msg": "No bearing data in database"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"passOK": "1", "data": result})
}

// POST /setting/uploadBearing
func (h *Handler) uploadBearing(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil || header.Filename == "" {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "error": "No selected file"})
		return
	}
	defer file.Close()

	configDir := h.deps.Config.ConfigDir
	tempPath := filepath.Join(configDir, header.Filename)
	dst, err := os.Create(tempPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "error": err.Error()})
		return
	}
	if _, err := io.Copy(dst, file); err != nil {
		dst.Close()
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "error": err.Error()})
		return
	}
	dst.Close()

	rows, err := getBearingFromFile(tempPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "error": err.Error()})
		return
	}

	dbPath := filepath.Join(configDir, "bearings.db")
	db, err := openBearingDB(dbPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "error": err.Error()})
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(`INSERT INTO bearings (Name, BPFO, BPFI, BSF, FTF) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "error": err.Error()})
		return
	}
	defer stmt.Close()

	var inserted []bearingRow
	skipped := 0
	for _, r := range rows {
		if _, err := stmt.Exec(r.Name, r.BPFO, r.BPFI, r.BSF, r.FTF); err != nil {
			skipped++
			continue
		}
		inserted = append(inserted, r)
	}

	// Timestamped backup copy, then remove the temp upload.
	timestamp := time.Now().Format("20060102_15")
	base := header.Filename
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	backupPath := filepath.Join(configDir, fmt.Sprintf("%s_%s%s", stem, timestamp, ext))
	copyFile(tempPath, backupPath)
	os.Remove(tempPath)

	if len(inserted) == 0 && skipped > 0 {
		c.JSON(http.StatusOK, gin.H{
			"passOK": "1", "data": []bearingRow{},
			"msg": fmt.Sprintf("All %d bearings already exist in database", skipped),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"passOK":   "1",
		"data":     inserted,
		"inserted": len(inserted),
		"skipped":  skipped,
	})
}

func openBearingDB(path string) (*sql.DB, error) {
	os.MkdirAll(filepath.Dir(path), 0755)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS bearings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		Name TEXT UNIQUE NOT NULL,
		BPFO REAL NOT NULL, BPFI REAL NOT NULL,
		BSF REAL NOT NULL, FTF REAL NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// getBearingFromFile mirrors Python get_Bearing: CSV → rows; non-CSV → try JSON parse.
func getBearingFromFile(path string) ([]bearingRow, error) {
	if !strings.EqualFold(filepath.Ext(path), ".csv") {
		return nil, fmt.Errorf("only .csv files are supported")
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) < 2 {
		return nil, nil
	}

	header := records[0]
	idx := map[string]int{}
	for i, h := range header {
		idx[strings.TrimSpace(h)] = i
	}

	toF := func(row []string, key string) float64 {
		i, ok := idx[key]
		if !ok || i >= len(row) {
			return 0
		}
		v, _ := strconv.ParseFloat(strings.TrimSpace(row[i]), 64)
		return v
	}

	out := make([]bearingRow, 0, len(records)-1)
	for _, rec := range records[1:] {
		name := ""
		if i, ok := idx["Name"]; ok && i < len(rec) {
			name = strings.TrimSpace(rec[i])
		}
		out = append(out, bearingRow{
			Name: name,
			BPFO: toF(rec, "BPFO"),
			BPFI: toF(rec, "BPFI"),
			BSF:  toF(rec, "BSF"),
			FTF:  toF(rec, "FTF"),
		})
	}
	return out, nil
}
