package setting

import (
	"context"
	"database/sql"
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"

	binpkg "serverGO/binary"
	"serverGO/infra"
	"serverGO/metermap"
)

// -----------------------------------------------------------------------------
// Models (mirror python_backend/routes/config.py pydantic models)
// -----------------------------------------------------------------------------

type CaliSet struct {
	Channel string `json:"channel"`
	Type    string `json:"type"`
	Cmd     string `json:"cmd"`
	CmdNum  string `json:"cmdnum"`
	Ref1    string `json:"ref1"`
	Ref2    string `json:"ref2"`
	Param   string `json:"param"`
}

type CaliRef struct {
	U     string `json:"U"`
	I     string `json:"I"`
	In    string `json:"In"`
	P     string `json:"P"`
	Error string `json:"Error"`
}

type TimeSetReq struct {
	DatetimeStr string `json:"datetime_str" binding:"required"`
	Timezone    string `json:"timezone"`
}

type PostData struct {
	ID           int    `json:"id,omitempty"`
	Title        string `json:"title"`
	Context      string `json:"context"`
	Mtype        int    `json:"mtype"`
	Utype        string `json:"utype"`
	FVersion     string `json:"f_version"`
	AVersion     string `json:"a_version"`
	WVersion     string `json:"w_version"`
	CVersion     string `json:"c_version"`
	SmartVersion string `json:"smart_version"`
	BuildVersion string `json:"build_version"`
	Date         string `json:"date"`
}

// -----------------------------------------------------------------------------
// Handler and routes
// -----------------------------------------------------------------------------

type ConfigHandler struct {
	deps          *infra.Dependencies
	maintenanceDB *sql.DB
	logDB         *sql.DB
}

// RegisterConfigRoutes mounts every Python routes/config.py endpoint on the
// supplied group. Skipped on purpose: /getTrain, /backup/restore/influxdb.
func RegisterConfigRoutes(rg *gin.RouterGroup, deps *infra.Dependencies) {
	mdb, err := openMaintenanceDBCfg(deps.Config.ConfigDir)
	if err != nil {
		log.Printf("config: open maintenance.db: %v", err)
	}
	ldb, err := openLogDBCfg(deps.Config.ConfigDir)
	if err != nil {
		log.Printf("config: open log.db: %v", err)
	}

	h := &ConfigHandler{deps: deps, maintenanceDB: mdb, logDB: ldb}

	rg.POST("/upload", h.upload)
	rg.GET("/checkSetup", h.checkSetup)
	rg.GET("/calibrate/getUnbal", h.getUnbal)
	rg.GET("/calibrateNow", h.calibrateNow)
	rg.GET("/calibrate/applySetup", h.applySetup)
	rg.POST("/calibrate/saveRef", h.saveRef)
	rg.GET("/calibrate/start", h.caliStart)
	rg.GET("/calibrate/end", h.caliEnd)
	rg.POST("/calibrate/cmd", h.caliCmd)

	rg.GET("/calibrate/gettime", h.getDeviceTime)
	rg.POST("/calibrate/setSystemTime", h.setSystemTime)
	rg.POST("/checktime", h.checkDeviceTime)

	rg.POST("/savePost/:mode/:idx", h.savePost)
	rg.GET("/getPost", h.getPost)
	rg.GET("/getLastPost", h.getLastPost)
	rg.GET("/deletePost/:idx", h.deletePost)

	rg.GET("/getReleaseNotes", h.getReleaseNotes)
	rg.GET("/getReleaseNote/:lang/:version", h.getReleaseNote)

	rg.GET("/getLog", h.getLogPaged)
	rg.DELETE("/deleteLog", h.deleteAllLogs)
	rg.GET("/applog/recent/:item", h.getRecentLogs)
}

// -----------------------------------------------------------------------------
// Calibration: upload / checkSetup / getUnbal / applySetup / saveRef / start / end / cmd
// -----------------------------------------------------------------------------

// POST /config/upload
// Mirrors python upload_file — CSV → DictReader → hset("calibration","setup", json)
func (h *ConfigHandler) upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil || header.Filename == "" {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "error": "No selected file"})
		return
	}
	defer file.Close()

	configDir := h.deps.Config.ConfigDir
	destPath := filepath.Join(configDir, header.Filename)
	dst, err := os.Create(destPath)
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

	rows, err := parseCSVDict(destPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "error": err.Error()})
		return
	}

	payload, _ := json.Marshal(rows)
	ctx := context.Background()
	if err := h.deps.Redis.Client0.HSet(ctx, "calibration", "setup", string(payload)).Err(); err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"passOK": "1", "data": rows, "file_path": header.Filename})
}

// GET /config/checkSetup
func (h *ConfigHandler) checkSetup(c *gin.Context) {
	ctx := context.Background()
	mac := getMACPlain()

	raw, err := h.deps.Redis.Client0.HGet(ctx, "calibration", "setup").Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "error": "No exist calibration setup", "mac": mac})
		return
	}
	var data interface{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "error": "No exist calibration setup", "mac": mac})
		return
	}
	c.JSON(http.StatusOK, gin.H{"passOK": "1", "data": data, "mac": mac})
}

// GET /config/calibrate/getUnbal
func (h *ConfigHandler) getUnbal(c *gin.Context) {
	ctx := context.Background()
	configDir := h.deps.Config.ConfigDir
	settingPath := filepath.Join(configDir, "setup.json")
	defaultPath := filepath.Join(configDir, "default.json")

	var setting map[string]interface{}

	if raw, err := h.deps.Redis.Client0.HGet(ctx, "System", "setup").Result(); err == nil {
		_ = json.Unmarshal([]byte(raw), &setting)
	}

	if setting == nil {
		data, err := os.ReadFile(settingPath)
		if err != nil || json.Unmarshal(data, &setting) != nil {
			// setup.json 손상 시 default.json으로 복구
			if def, err2 := os.ReadFile(defaultPath); err2 == nil {
				_ = os.WriteFile(settingPath, def, 0644)
				_ = json.Unmarshal(def, &setting)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"unbal": setting["unbalance"]})
}

// GET /config/calibrateNow
// Reads per-phase calibration meter data for Main and Sub channels.
func (h *ConfigHandler) calibrateNow(c *gin.Context) {
	ctx := context.Background()
	main := buildCalibrateChannel(ctx, h.deps, "Main")
	sub := buildCalibrateChannel(ctx, h.deps, "Sub")

	c.JSON(http.StatusOK, gin.H{
		"mainStatus": main["success"],
		"mainData":   dig(main, "retData", "meterData"),
		"mainRef":    dig(main, "retData", "refData"),
		"subStatus":  sub["success"],
		"subData":    dig(sub, "retData", "meterData"),
		"subRef":     dig(sub, "retData", "refData"),
	})
}

// buildCalibrateChannel mirrors Python get_Calibrate(channel).
func buildCalibrateChannel(ctx context.Context, deps *infra.Dependencies, channel string) map[string]interface{} {
	client := deps.Redis.Client0

	var refDict interface{}
	if refStr, err := client.HGet(ctx, "calibration", "ref").Result(); err == nil {
		var ref map[string]interface{}
		if json.Unmarshal([]byte(refStr), &ref) == nil {
			refDict = ref
		}
	}

	keyname := "meter"
	exists, _ := client.Exists(ctx, keyname).Result()
	if exists == 0 {
		return map[string]interface{}{"success": false, "error": "No exist key"}
	}

	flat, _ := client.HGetAll(ctx, keyname).Result()
	meters := make(map[string]float64, len(flat))
	for k, v := range flat {
		meters[k] = tryFloatCfg(v)
	}

	oneSec, _ := binpkg.FetchMaxMin1Sec(ctx, client, channel)
	maxmin := binpkg.Flat1Sec(oneSec)

	pVolt := metermap.GetDataDict(meters, maxmin, metermap.CalPVoltageKeys, "V")
	curr := metermap.GetDataDict(meters, maxmin, metermap.CalCurrentKeys, "A")
	angle := metermap.GetDataDict(meters, maxmin, metermap.CalPAngleKeys, "%")
	p := metermap.GetDataDict(meters, maxmin, metermap.CalAPowerKeys, "kw")
	q := metermap.GetDataDict(meters, maxmin, metermap.CalRPowerKeys, "kvar")
	s := metermap.GetDataDict(meters, maxmin, metermap.CalAPPowerKeys, "kVA")

	meterData := []gin.H{
		{"subTitle": "Phase Voltage", "data": pVolt},
		{"subTitle": "Current", "data": curr},
		{"subTitle": "PowerAngle", "data": angle},
		{"subTitle": "ActivePower", "data": p},
		{"subTitle": "ReactivePower", "data": q},
		{"subTitle": "ApparentPower", "data": s},
	}

	return map[string]interface{}{
		"success": true,
		"retData": gin.H{
			"success":   true,
			"meterData": meterData,
			"refData":   refDict,
		},
	}
}

// GET /config/calibrate/applySetup
func (h *ConfigHandler) applySetup(c *gin.Context) {
	ctx := context.Background()
	client := h.deps.Redis.Client0
	configDir := h.deps.Config.ConfigDir
	settingFile := filepath.Join(configDir, "setup.json")
	backupFile := filepath.Join(configDir, "setup_backup.json")

	raw, err := client.HGet(ctx, "calibration", "setup").Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "error": "No exist calibration setup"})
		return
	}

	// 현재 setting을 backup으로 백업
	if _, err := os.Stat(backupFile); err == nil {
		_ = copyFile(backupFile, settingFile)
	}

	var list interface{}
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "error": err.Error()})
		return
	}

	payload, _ := json.Marshal(list)
	client.HSet(ctx, "System", "setup", string(payload))

	// setup.json 기록
	_ = os.WriteFile(settingFile, payload, 0644)

	client.HSet(ctx, "service", "cflag", 1)
	client.HSet(ctx, "Service", "save", 1)

	c.JSON(http.StatusOK, gin.H{"passOK": "1"})
}

// POST /config/calibrate/saveRef
func (h *ConfigHandler) saveRef(c *gin.Context) {
	var ref CaliRef
	if err := c.ShouldBindJSON(&ref); err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0"})
		return
	}

	uInt, _ := strconv.Atoi(ref.U)
	iInt, _ := strconv.Atoi(ref.I)
	inInt, _ := strconv.Atoi(ref.In)
	pInt, _ := strconv.Atoi(ref.P)
	errFloat, _ := strconv.ParseFloat(ref.Error, 64)

	refData := map[string]interface{}{
		"U": uInt, "I": iInt, "In": inInt, "P": pInt, "Error": errFloat,
	}
	payload, _ := json.Marshal(refData)

	ctx := context.Background()
	if err := h.deps.Redis.Client0.HSet(ctx, "calibration", "ref", string(payload)).Err(); err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"passOK": "1"})
}

// GET /config/calibrate/start
func (h *ConfigHandler) caliStart(c *gin.Context) {
	ctx := context.Background()
	if err := h.deps.Redis.Client0.HSet(ctx, "Service", "calibration", 1).Err(); err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"passOK": "1"})
}

// GET /config/calibrate/end
func (h *ConfigHandler) caliEnd(c *gin.Context) {
	ctx := context.Background()
	client := h.deps.Redis.Client0
	configDir := h.deps.Config.ConfigDir
	settingFile := filepath.Join(configDir, "setup.json")
	backupFile := filepath.Join(configDir, "setup_backup.json")

	client.HSet(ctx, "Service", "calibration", 0)

	endflag := false
	hasSetup, _ := client.HExists(ctx, "calibration", "setup").Result()
	hasCflag, _ := client.HExists(ctx, "calibration", "cflag").Result()

	if hasSetup && hasCflag {
		if cflag, err := client.HGet(ctx, "calibration", "cflag").Result(); err == nil && cflag == "1" {
			endflag = true
			client.HSet(ctx, "calibration", "cflag", 0)
		}
	}

	if endflag {
		if err := copyFile(settingFile, backupFile); err != nil {
			c.JSON(http.StatusOK, gin.H{"passOK": "0", "error": "No exist setup file"})
			return
		}
		data, err := os.ReadFile(settingFile)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"passOK": "0", "error": "No exist setup file"})
			return
		}
		var setting interface{}
		_ = json.Unmarshal(data, &setting)
		payload, _ := json.Marshal(setting)
		client.HSet(ctx, "System", "setup", string(payload))
		client.HDel(ctx, "calibration", "setup")
		client.HSet(ctx, "Service", "save", 1)
		c.JSON(http.StatusOK, gin.H{"passOK": "1"})
		return
	}

	if hasRef, _ := client.HExists(ctx, "calibration", "ref").Result(); hasRef {
		client.HDel(ctx, "calibration", "ref")
	}
	c.JSON(http.StatusOK, gin.H{"passOK": "1"})
}

// POST /config/calibrate/cmd
func (h *ConfigHandler) caliCmd(c *gin.Context) {
	var req CaliSet
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0"})
		return
	}

	var mode int32
	switch req.Channel {
	case "Main":
		mode = 0
	case "Sub":
		mode = 1
	default:
		mode = 2
	}

	var val1, val2 float32
	if req.Type == "SET" {
		if req.Ref1 != "None" {
			if v, err := strconv.ParseFloat(req.Ref1, 32); err == nil {
				val1 = float32(v)
			}
		}
		if req.Ref2 != "None" {
			if v, err := strconv.ParseFloat(req.Ref2, 32); err == nil {
				val2 = float32(v)
			}
		}
	}

	cmdNum, _ := strconv.Atoi(req.CmdNum)

	ctx := context.Background()
	client := h.deps.Redis.Client0

	// calibration.ref 업데이트 (param 지정 시)
	if hasRef, _ := client.HExists(ctx, "calibration", "ref").Result(); hasRef && req.Param != "None" {
		if refStr, err := client.HGet(ctx, "calibration", "ref").Result(); err == nil {
			var refDict map[string]interface{}
			if json.Unmarshal([]byte(refStr), &refDict) == nil {
				r1, _ := strconv.ParseFloat(req.Ref1, 64)
				r2, _ := strconv.ParseFloat(req.Ref2, 64)
				if strings.Contains(req.Param, ",") {
					refDict["U"] = r1
					refDict["I"] = r2
				} else {
					refDict[req.Param] = r1
				}
				payload, _ := json.Marshal(refDict)
				client.HSet(ctx, "calibration", "ref", string(payload))
			}
		}
	}

	// binary 'iiff' pack
	buf := make([]byte, 16)
	binary.LittleEndian.PutUint32(buf[0:4], uint32(mode))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(cmdNum))
	binary.LittleEndian.PutUint32(buf[8:12], math.Float32bits(val1))
	binary.LittleEndian.PutUint32(buf[12:16], math.Float32bits(val2))

	// Python uses redis_state.binary_client (DB 0)
	if err := client.LPush(ctx, "cali_command", buf).Err(); err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0"})
		return
	}

	if req.Type == "SAVE" {
		post := PostData{Title: "Calibration", Context: "Calibration", Mtype: 3}
		_ = insertPostCfg(h.maintenanceDB, post)
	}

	c.JSON(http.StatusOK, gin.H{"passOK": "1"})
}

// -----------------------------------------------------------------------------
// Time endpoints
// -----------------------------------------------------------------------------

// GET /config/calibrate/gettime
func (h *ConfigHandler) getDeviceTime(c *gin.Context) {
	now := time.Now().Format("2006/01/02 15:04:05")
	c.JSON(http.StatusOK, gin.H{"success": true, "deviceTime": now})
}

// POST /config/calibrate/setSystemTime
func (h *ConfigHandler) setSystemTime(c *gin.Context) {
	var req TimeSetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}
	if req.Timezone == "" {
		req.Timezone = "Asia/Seoul"
	}

	_ = exec.Command("sudo", "timedatectl", "set-ntp", "false").Run()

	if out, err := exec.Command("sudo", "timedatectl", "set-timezone", req.Timezone).CombinedOutput(); err != nil {
		log.Printf("set-timezone failed: %s", string(out))
	}

	if out, err := exec.Command("sudo", "date", "-s", req.DatetimeStr).CombinedOutput(); err != nil {
		log.Printf("date -s failed: %s", string(out))
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "date -s failed: " + strings.TrimSpace(string(out))})
		return
	}

	if out, err := exec.Command("sudo", "hwclock", "-w").CombinedOutput(); err != nil {
		log.Printf("hwclock -w failed: %s", string(out))
	}

	current, _ := exec.Command("date").Output()

	session := sessions.Default(c)
	account, _ := session.Get("user").(string)
	role, _ := session.Get("userRole").(string)
	if h.logDB != nil && account != "" {
		_ = upsertLogCfg(h.logDB, account, role, "Set Time")
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "System time updated",
		"current_time": strings.TrimSpace(string(current)),
		"timezone":     req.Timezone,
	})
}

// POST /config/checktime
func (h *ConfigHandler) checkDeviceTime(c *gin.Context) {
	var req TimeSetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "status": false})
		return
	}

	deviceNow := time.Now()
	clientDt, err := time.ParseInLocation("2006-01-02 15:04:05", req.DatetimeStr, time.Local)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "status": false})
		return
	}

	diff := deviceNow.Sub(clientDt)
	if diff < 0 {
		diff = -diff
	}
	status := diff <= 24*time.Hour

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"deviceTime":  deviceNow.Format(time.RFC3339),
		"status":      status,
		"diffSeconds": int(diff.Seconds()),
	})
}

// -----------------------------------------------------------------------------
// Maintenance post endpoints
// -----------------------------------------------------------------------------

// POST /config/savePost/:mode/:idx
func (h *ConfigHandler) savePost(c *gin.Context) {
	mode, _ := strconv.Atoi(c.Param("mode"))
	idx, _ := strconv.Atoi(c.Param("idx"))

	var post PostData
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "msg": err.Error()})
		return
	}

	var err error
	if mode == 0 {
		err = insertPostCfg(h.maintenanceDB, post)
	} else {
		err = updatePostCfg(h.maintenanceDB, post, idx)
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"passOK": "1"})
}

// GET /config/getPost
func (h *ConfigHandler) getPost(c *gin.Context) {
	posts, err := getAllPostsCfg(h.maintenanceDB)
	if err != nil || len(posts) == 0 {
		c.JSON(http.StatusOK, gin.H{"result": 0})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": 1, "data": posts})
}

// GET /config/getLastPost
func (h *ConfigHandler) getLastPost(c *gin.Context) {
	post, err := getLastPostCfg(h.maintenanceDB)
	if err != nil || post == nil {
		c.JSON(http.StatusOK, gin.H{"result": 0})
		return
	}

	versions := getVersions()
	if len(versions) == 0 {
		c.JSON(http.StatusOK, gin.H{"result": 1, "data": post})
		return
	}

	keyMap := map[string]*string{
		"fw":          &post.FVersion,
		"a35":         &post.AVersion,
		"web":         &post.WVersion,
		"core":        &post.CVersion,
		"smartsystem": &post.SmartVersion,
	}

	var updateList []string
	for src, dstPtr := range keyMap {
		if v, ok := versions[src]; ok && v != "" && *dstPtr != v {
			*dstPtr = v
			updateList = append(updateList, src)
		}
	}

	if len(updateList) > 0 {
		update := PostData{
			Title:        "Update SW",
			Context:      "Update SW",
			Mtype:        1,
			Utype:        strings.Join(updateList, ","),
			FVersion:     post.FVersion,
			AVersion:     post.AVersion,
			WVersion:     post.WVersion,
			CVersion:     post.CVersion,
			SmartVersion: post.SmartVersion,
			BuildVersion: post.BuildVersion,
		}
		c.JSON(http.StatusOK, gin.H{"result": 1, "data": update, "update": updateList})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": 1, "data": post})
}

// GET /config/deletePost/:idx
func (h *ConfigHandler) deletePost(c *gin.Context) {
	idx := c.Param("idx")
	if err := deletePostCfg(h.maintenanceDB, idx); err != nil {
		c.JSON(http.StatusOK, gin.H{"result": 0})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": 1})
}

// -----------------------------------------------------------------------------
// Release notes
// -----------------------------------------------------------------------------

// GET /config/getReleaseNotes
func (h *ConfigHandler) getReleaseNotes(c *gin.Context) {
	dir := filepath.Join(h.deps.Config.BaseDir, "release_notes", "ko")
	entries, err := os.ReadDir(dir)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "data": []string{}})
		return
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		names = append(names, strings.TrimSuffix(e.Name(), ".md"))
	}
	sort.Sort(sort.Reverse(sort.StringSlice(names)))
	c.JSON(http.StatusOK, gin.H{"success": true, "data": names})
}

// GET /config/getReleaseNote/:lang/:version
func (h *ConfigHandler) getReleaseNote(c *gin.Context) {
	lang := c.Param("lang")
	version := c.Param("version")
	path := filepath.Join(h.deps.Config.BaseDir, "release_notes", lang, version+".md")

	data, err := os.ReadFile(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "content": "", "message": "File not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "content": string(data)})
}

// -----------------------------------------------------------------------------
// Logs
// -----------------------------------------------------------------------------

// GET /config/getLog?page=&page_size=
func (h *ConfigHandler) getLogPaged(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	logs, total, err := getLogsPagedCfg(h.logDB, page, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": 0, "msg": err.Error()})
		return
	}
	totalPages := 0
	if pageSize > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}
	result := 0
	if len(logs) > 0 {
		result = 1
	}
	c.JSON(http.StatusOK, gin.H{
		"result":      result,
		"data":        logs,
		"total":       total,
		"page":        page,
		"total_pages": totalPages,
	})
}

// DELETE /config/deleteLog
func (h *ConfigHandler) deleteAllLogs(c *gin.Context) {
	if _, err := h.logDB.Exec("DELETE FROM log"); err != nil {
		c.JSON(http.StatusOK, gin.H{"result": 0, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": 1, "message": "All logs deleted"})
}

// GET /config/applog/recent/:item?lines=&log_type=
func (h *ConfigHandler) getRecentLogs(c *gin.Context) {
	item := c.Param("item")
	lines, _ := strconv.Atoi(c.DefaultQuery("lines", "5"))
	if lines <= 0 {
		lines = 5
	}
	logType := c.DefaultQuery("log_type", "all")

	journalServices := map[string]string{
		"frpc":     "frpc",
		"mqClient": "mqClient",
	}
	logPaths := map[string]string{
		"SmartSystems": "/usr/local/smartsystems/log",
		"Core":         "/usr/local/sv500/logs/core",
		"WebServer":    "/usr/local/sv500/logs/web",
		"A35":          "/usr/local/sv500/logs/a35",
	}

	if svc, ok := journalServices[item]; ok {
		if !serviceExists(svc + ".service") {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": item + " 서비스가 설치되지 않음"})
			return
		}
		out, _ := exec.Command("journalctl", "-u", svc, "-n", strconv.Itoa(lines), "--no-pager", "-o", "short").Output()
		var result []string
		for _, l := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			if l != "" && !strings.HasPrefix(l, "--") {
				result = append(result, l)
			}
		}
		if len(result) == 0 {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "로그 없음"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true, "source": "journal", "service": svc,
			"lines": result, "count": len(result),
		})
		return
	}

	dir, ok := logPaths[item]
	if !ok {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "잘못된 항목"})
		return
	}
	if _, err := os.Stat(dir); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "로그 디렉토리 없음"})
		return
	}

	entries, _ := os.ReadDir(dir)
	type logFile struct {
		name    string
		modTime time.Time
	}
	var files []logFile
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".log") {
			continue
		}
		info, err := e.Info()
		if err != nil || info.Size() == 0 {
			continue
		}
		if item == "SmartSystems" && logType != "all" {
			if logType == "ss" && !strings.HasPrefix(e.Name(), "SS") {
				continue
			}
			if logType == "api" && !strings.HasPrefix(e.Name(), "RestAPI") {
				continue
			}
		}
		files = append(files, logFile{e.Name(), info.ModTime()})
	}
	if len(files) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "로그 파일 없음"})
		return
	}

	latest := files[0]
	for _, f := range files[1:] {
		if f.modTime.After(latest.modTime) {
			latest = f
		}
	}

	filePath := filepath.Join(dir, latest.name)
	out, _ := exec.Command("tail", "-n", strconv.Itoa(lines), filePath).Output()
	var result []string
	if s := strings.TrimSpace(string(out)); s != "" {
		result = strings.Split(s, "\n")
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true, "source": "file", "file": latest.name,
		"lines": result, "count": len(result),
	})
}

// -----------------------------------------------------------------------------
// Helpers: MAC/IP
// -----------------------------------------------------------------------------

// getMACPlain returns MAC without separators, lowercased — matches Python get_mac_address.
func getMACPlain() string {
	ifaces, err := net.Interfaces()
	if err == nil {
		for _, name := range []string{"sw0ep", "end1"} {
			for _, i := range ifaces {
				if i.Name == name && i.HardwareAddr != nil {
					mac := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(i.HardwareAddr.String(), ":", ""), "-", ""))
					if mac != "" && mac != "000000000000" {
						return mac
					}
				}
			}
		}
	}
	// Fallback — any non-loopback MAC
	if ifaces != nil {
		for _, i := range ifaces {
			if i.Flags&net.FlagLoopback == 0 && i.HardwareAddr != nil {
				mac := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(i.HardwareAddr.String(), ":", ""), "-", ""))
				if mac != "" && mac != "000000000000" {
					return mac
				}
			}
		}
	}
	return "000000000000"
}

// -----------------------------------------------------------------------------
// Helpers: CSV / misc
// -----------------------------------------------------------------------------

func parseCSVDict(path string) ([]map[string]string, error) {
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
	if len(records) < 1 {
		return []map[string]string{}, nil
	}
	header := records[0]
	rows := make([]map[string]string, 0, len(records)-1)
	for _, rec := range records[1:] {
		m := make(map[string]string, len(header))
		for i, h := range header {
			if i < len(rec) {
				m[h] = rec[i]
			} else {
				m[h] = ""
			}
		}
		rows = append(rows, m)
	}
	return rows, nil
}

func tryFloatCfg(s string) float64 {
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsNaN(v) {
		return 0
	}
	return v
}

func dig(m map[string]interface{}, keys ...string) interface{} {
	var cur interface{} = m
	for _, k := range keys {
		obj, ok := cur.(map[string]interface{})
		if !ok {
			obj2, ok2 := cur.(gin.H)
			if !ok2 {
				return nil
			}
			obj = map[string]interface{}(obj2)
		}
		cur = obj[k]
	}
	return cur
}

// -----------------------------------------------------------------------------
// SQLite stores (maintenance.db, log.db)
// -----------------------------------------------------------------------------

func openMaintenanceDBCfg(configDir string) (*sql.DB, error) {
	dbPath := filepath.Join(configDir, "maintenance.db")
	os.MkdirAll(configDir, 0755)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS maintenance (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		context TEXT NOT NULL,
		mtype INT NOT NULL,
		utype TEXT NOT NULL,
		f_version TEXT,
		a_version TEXT,
		w_version TEXT,
		c_version TEXT,
		smart_version TEXT,
		build_version TEXT,
		date TEXT
	)`)
	if err != nil {
		db.Close()
		return nil, err
	}

	// build_version migration
	rows, _ := db.Query("PRAGMA table_info(maintenance)")
	hasBuild := false
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		_ = rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk)
		if name == "build_version" {
			hasBuild = true
		}
	}
	rows.Close()
	if !hasBuild {
		db.Exec("ALTER TABLE maintenance ADD COLUMN build_version TEXT DEFAULT ''")
	}

	return db, nil
}

func openLogDBCfg(configDir string) (*sql.DB, error) {
	dbPath := "/usr/local/sv500/logs/web/log.db"
	if _, err := os.Stat(filepath.Dir(dbPath)); err != nil {
		// Fallback to configDir on non-target systems
		dbPath = filepath.Join(configDir, "log.db")
	}
	os.MkdirAll(filepath.Dir(dbPath), 0755)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		logdate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		account TEXT NOT NULL,
		userRole TEXT NOT NULL,
		action TEXT NOT NULL
	)`)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func insertPostCfg(db *sql.DB, p PostData) error {
	if db == nil {
		return fmt.Errorf("maintenance db not open")
	}
	today := time.Now().Format("2006-01-02")
	_, err := db.Exec(
		`INSERT INTO maintenance (title, context, mtype, utype, f_version, a_version, w_version, c_version, smart_version, build_version, date) VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
		p.Title, p.Context, p.Mtype, p.Utype, p.FVersion, p.AVersion, p.WVersion, p.CVersion, p.SmartVersion, p.BuildVersion, today,
	)
	return err
}

func updatePostCfg(db *sql.DB, p PostData, id int) error {
	if db == nil {
		return fmt.Errorf("maintenance db not open")
	}
	today := time.Now().Format("2006-01-02")
	_, err := db.Exec(
		`UPDATE maintenance SET title=?, context=?, mtype=?, utype=?, f_version=?, a_version=?, w_version=?, c_version=?, smart_version=?, build_version=?, date=? WHERE id=?`,
		p.Title, p.Context, p.Mtype, p.Utype, p.FVersion, p.AVersion, p.WVersion, p.CVersion, p.SmartVersion, p.BuildVersion, today, id,
	)
	return err
}

func getAllPostsCfg(db *sql.DB) ([]PostData, error) {
	if db == nil {
		return nil, fmt.Errorf("maintenance db not open")
	}
	rows, err := db.Query(`SELECT id, title, context, mtype, utype,
		COALESCE(f_version,''), COALESCE(a_version,''), COALESCE(w_version,''),
		COALESCE(c_version,''), COALESCE(smart_version,''), COALESCE(build_version,''),
		COALESCE(date,'') FROM maintenance ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []PostData
	for rows.Next() {
		var p PostData
		if err := rows.Scan(&p.ID, &p.Title, &p.Context, &p.Mtype, &p.Utype,
			&p.FVersion, &p.AVersion, &p.WVersion, &p.CVersion, &p.SmartVersion, &p.BuildVersion, &p.Date); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func getLastPostCfg(db *sql.DB) (*PostData, error) {
	if db == nil {
		return nil, fmt.Errorf("maintenance db not open")
	}
	row := db.QueryRow(`SELECT id, title, context, mtype, utype,
		COALESCE(f_version,''), COALESCE(a_version,''), COALESCE(w_version,''),
		COALESCE(c_version,''), COALESCE(smart_version,''), COALESCE(build_version,''),
		COALESCE(date,'') FROM maintenance ORDER BY id DESC LIMIT 1`)
	var p PostData
	if err := row.Scan(&p.ID, &p.Title, &p.Context, &p.Mtype, &p.Utype,
		&p.FVersion, &p.AVersion, &p.WVersion, &p.CVersion, &p.SmartVersion, &p.BuildVersion, &p.Date); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func deletePostCfg(db *sql.DB, idx string) error {
	if db == nil {
		return fmt.Errorf("maintenance db not open")
	}
	_, err := db.Exec("DELETE FROM maintenance WHERE id=?", idx)
	return err
}

type LogRow struct {
	ID       int    `json:"id"`
	LogDate  string `json:"logdate"`
	Account  string `json:"account"`
	UserRole string `json:"userRole"`
	Action   string `json:"action"`
}

func getLogsPagedCfg(db *sql.DB, page, pageSize int) ([]LogRow, int, error) {
	if db == nil {
		return nil, 0, fmt.Errorf("log db not open")
	}
	var total int
	_ = db.QueryRow("SELECT COUNT(*) FROM log").Scan(&total)

	if pageSize <= 0 {
		pageSize = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize

	rows, err := db.Query(
		`SELECT id, datetime(logdate,'localtime') as logdate, account, userRole, action
		 FROM log ORDER BY id DESC LIMIT ? OFFSET ?`, pageSize, offset)
	if err != nil {
		return nil, total, err
	}
	defer rows.Close()

	var out []LogRow
	for rows.Next() {
		var l LogRow
		if err := rows.Scan(&l.ID, &l.LogDate, &l.Account, &l.UserRole, &l.Action); err != nil {
			return nil, total, err
		}
		out = append(out, l)
	}
	return out, total, nil
}

func upsertLogCfg(db *sql.DB, account, role, action string) error {
	if db == nil {
		return fmt.Errorf("log db not open")
	}
	var id int
	err := db.QueryRow("SELECT id FROM log WHERE action=?", action).Scan(&id)
	if err == nil {
		_, err = db.Exec("UPDATE log SET logdate=CURRENT_TIMESTAMP, account=?, userRole=? WHERE id=?", account, role, id)
		return err
	}
	_, err = db.Exec("INSERT INTO log (account, userRole, action) VALUES (?,?,?)", account, role, action)
	return err
}
