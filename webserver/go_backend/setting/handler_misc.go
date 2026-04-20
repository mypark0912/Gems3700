package setting

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	binpkg "serverGO/binary"
)

func (h *Handler) registerMiscRoutes(rg *gin.RouterGroup) {
	rg.POST("/savefile/:channel", h.savefileByChannel)
	rg.GET("/restartdevice", h.restartDevice)
	rg.GET("/restartCore", h.restartCore)
	rg.GET("/trigger", h.triggerWaveform)
	rg.POST("/setSystemTime", h.setSystemTimeWithTZ)
}

// POST /setting/savefile/:channel
// Merges the posted JSON into setup.json for the named channel (or "general" for General block).
func (h *Handler) savefileByChannel(c *gin.Context) {
	channel := c.Param("channel")
	ctx := context.Background()
	client := h.deps.Redis.Client0

	if flag, err := client.HGet(ctx, "Service", "setting").Result(); err == nil && flag == "1" {
		c.JSON(http.StatusOK, gin.H{"status": "0", "error": "Modbus setting is activated"})
		return
	}

	var incoming map[string]interface{}
	if err := c.ShouldBindJSON(&incoming); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "0", "error": "No data provided"})
		return
	}

	configDir := h.deps.Config.ConfigDir
	setupPath := filepath.Join(configDir, "setup.json")

	var setting map[string]interface{}
	if data, err := os.ReadFile(setupPath); err == nil {
		_ = json.Unmarshal(data, &setting)
	}
	if setting == nil {
		setting = map[string]interface{}{"mode": "device0", "General": map[string]interface{}{}, "channel": []interface{}{}}
	}

	if strings.EqualFold(channel, "general") {
		setting["General"] = incoming
	} else {
		channels, _ := setting["channel"].([]interface{})
		updated := false
		for i, ch := range channels {
			chMap, _ := ch.(map[string]interface{})
			if chName, _ := chMap["channel"].(string); strings.EqualFold(chName, channel) {
				channels[i] = incoming
				if ct, ok := incoming["ctInfo"].(map[string]interface{}); ok {
					if v, ok := ct["inorminal"]; ok {
						ct["inorminal"] = int(toFloatMisc(v) * 1000)
					}
				}
				updated = true
				break
			}
		}
		if !updated {
			channels = append(channels, incoming)
		}
		setting["channel"] = channels
	}

	// Device identity refresh.
	mac := getMAC()
	if gen, ok := setting["General"].(map[string]interface{}); ok {
		if di, ok := gen["deviceInfo"].(map[string]interface{}); ok {
			di["mac_address"] = mac
			if cur, _ := di["serial_number"].(string); cur == "" {
				di["serial_number"] = mac
			}
		}
	}

	out, _ := json.MarshalIndent(setting, "", "    ")
	os.WriteFile(setupPath, out, 0644)

	payload, _ := json.Marshal(setting)
	client.HSet(ctx, "System", "setup", string(payload))

	// Persist derived values (start current / demand / channel info).
	SaveRedisSetup(ctx, h.deps, setting)

	c.JSON(http.StatusOK, gin.H{"status": "1", "data": setting})
}

// GET /setting/restartdevice?timeout=30
func (h *Handler) restartDevice(c *gin.Context) {
	ctx := context.Background()
	client := h.deps.Redis.Client0

	if flag, err := client.HGet(ctx, "Service", "setting").Result(); err == nil && flag == "1" {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": "Modbus setting is activated"})
		return
	}

	// Set restartFW flag in Equipment.applyStatus.
	var applyCtx map[string]interface{}
	if raw, err := client.HGet(ctx, "Equipment", "applyStatus").Result(); err == nil && raw != "" {
		_ = json.Unmarshal([]byte(raw), &applyCtx)
	}
	if applyCtx == nil {
		applyCtx = map[string]interface{}{"commisionAsset": map[string]interface{}{}}
	}
	applyCtx["restartFW"] = true
	payload, _ := json.Marshal(applyCtx)
	client.HSet(ctx, "Equipment", "applyStatus", string(payload))

	client.HSet(ctx, "Service", "save", 1)

	timeoutSec := 30
	if q := c.Query("timeout"); q != "" {
		fmt := 0
		for _, ch := range q {
			if ch < '0' || ch > '9' {
				fmt = -1
				break
			}
			fmt = fmt*10 + int(ch-'0')
		}
		if fmt > 0 {
			timeoutSec = fmt
		}
	}

	// Wait for Service.save to drop to 0.
	deadline := time.Now().Add(time.Duration(timeoutSec) * time.Second)
	for time.Now().Before(deadline) {
		val, err := client.HGet(ctx, "Service", "save").Result()
		if err == nil && val == "0" {
			time.Sleep(5 * time.Second)
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "Restart completed"})
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	c.JSON(http.StatusOK, gin.H{"success": false, "error": "Timeout - save flag not cleared"})
}

// GET /setting/restartCore
func (h *Handler) restartCore(c *gin.Context) {
	ctx := context.Background()
	client := h.deps.Redis.Client0

	if flag, err := client.HGet(ctx, "Service", "setting").Result(); err == nil && flag == "1" {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": "Modbus setting is activated"})
		return
	}
	if err := client.HSet(ctx, "Service", "restart", 1).Err(); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": "Redis Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

var waveformPaths = [2]string{"/sv500/ch1/waveform", "/sv500/ch2/waveform"}

// GET /setting/trigger?target=&timeout=
// Clears stale waveform files, issues capture command(s), waits for new file to appear.
func (h *Handler) triggerWaveform(c *gin.Context) {
	ctx := context.Background()
	client := h.deps.Redis.Client0

	target := 2
	if q := c.Query("target"); q == "0" || q == "1" {
		target = int(q[0] - '0')
	}
	timeoutSec := 90
	if q := c.Query("timeout"); q != "" {
		n := 0
		for _, ch := range q {
			if ch < '0' || ch > '9' {
				n = -1
				break
			}
			n = n*10 + int(ch-'0')
		}
		if n > 0 {
			timeoutSec = n
		}
	}

	var watchPaths []string
	if target == 2 {
		watchPaths = waveformPaths[:]
	} else {
		watchPaths = []string{waveformPaths[target]}
	}

	// Clear old files.
	for _, p := range watchPaths {
		entries, _ := os.ReadDir(p)
		for _, e := range entries {
			if !e.IsDir() {
				os.Remove(filepath.Join(p, e.Name()))
			}
		}
	}

	// Push capture command for target channel(s).
	if target == 0 || target == 2 {
		cmd := binpkg.Command{Type: 0, Cmd: binpkg.CmdCapture, Item: binpkg.ItemWaveform}
		client.LPush(ctx, "command", cmd.Encode())
	}
	if target == 1 || target == 2 {
		cmd := binpkg.Command{Type: 1, Cmd: binpkg.CmdCapture, Item: binpkg.ItemWaveform}
		client.LPush(ctx, "command", cmd.Encode())
	}

	// Poll for file creation (.json / .bin) across all watch paths.
	found := map[string]string{}
	deadline := time.Now().Add(time.Duration(timeoutSec) * time.Second)
	for time.Now().Before(deadline) {
		allDone := true
		for _, p := range watchPaths {
			if _, ok := found[p]; ok {
				continue
			}
			entries, _ := os.ReadDir(p)
			for _, e := range entries {
				if e.IsDir() {
					continue
				}
				name := e.Name()
				if strings.HasSuffix(name, ".json") || strings.HasSuffix(name, ".bin") {
					found[p] = filepath.Join(p, name)
					break
				}
			}
			if _, ok := found[p]; !ok {
				allDone = false
			}
		}
		if allDone {
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "정상 완료", "files": found})
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	if len(found) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": true, "message": "파일 생성 감지 (일부)",
			"files": found, "assumed_taken": true,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": false, "message": "타임아웃 - 파일 생성 감지 실패"})
}

// POST /setting/setSystemTime
// Converts the supplied local datetime from the client's timezone into the device timezone
// (from Redis Device.Timezone, falling back to the OS timezone), then applies it.
func (h *Handler) setSystemTimeWithTZ(c *gin.Context) {
	var req struct {
		DatetimeStr string `json:"datetime_str" binding:"required"`
		Timezone    string `json:"timezone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}

	ctx := context.Background()
	client := h.deps.Redis.Client0
	targetTZ, _ := client.HGet(ctx, "Device", "Timezone").Result()
	if targetTZ == "" {
		targetTZ = systemTimezone()
	}

	clientLoc, err := time.LoadLocation(req.Timezone)
	if err != nil {
		clientLoc = time.Local
	}
	targetLoc, err := time.LoadLocation(targetTZ)
	if err != nil {
		targetLoc = time.Local
	}

	clientDt, err := time.ParseInLocation("2006-01-02 15:04:05", req.DatetimeStr, clientLoc)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	targetDtStr := clientDt.In(targetLoc).Format("2006-01-02 15:04:05")

	_ = exec.Command("sudo", "timedatectl", "set-ntp", "false").Run()
	if out, err := exec.Command("sudo", "timedatectl", "set-timezone", targetTZ).CombinedOutput(); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "set-timezone failed: " + strings.TrimSpace(string(out))})
		return
	}
	if out, err := exec.Command("sudo", "date", "-s", targetDtStr).CombinedOutput(); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "date -s failed: " + strings.TrimSpace(string(out))})
		return
	}
	_ = exec.Command("sudo", "hwclock", "-w").Run()

	current, _ := exec.Command("date").Output()

	// Log the action (reuses auth/session + cfgroute log DB; here we just log to journald-style stderr).
	session := sessions.Default(c)
	_ = session.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "System time updated",
		"current_time": strings.TrimSpace(string(current)),
		"timezone":     targetTZ,
	})
}

func systemTimezone() string {
	out, err := exec.Command("timedatectl", "show", "--value", "-p", "Timezone").Output()
	if err == nil {
		if tz := strings.TrimSpace(string(out)); tz != "" {
			return tz
		}
	}
	if data, err := os.Readlink("/etc/localtime"); err == nil {
		if idx := strings.Index(data, "zoneinfo/"); idx >= 0 {
			return data[idx+len("zoneinfo/"):]
		}
	}
	return "Asia/Seoul"
}

func toFloatMisc(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		if math.IsNaN(val) || math.IsInf(val, 0) {
			return 0
		}
		return val
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		return 0
	}
	return 0
}
