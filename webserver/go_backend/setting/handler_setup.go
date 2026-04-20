package setting

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"serverGO/infra"
)

type Handler struct {
	deps *infra.Dependencies
}

func RegisterRoutes(rg *gin.RouterGroup, deps *infra.Dependencies) {
	h := &Handler{deps: deps}
	cmdMgr := NewCMDManager(deps)

	// RS-485 command (WebSocket)
	rg.GET("/ws", cmdMgr.handleWS)
	rg.GET("/testScanResp", cmdMgr.testScanResp)
	rg.GET("/testPingResp", cmdMgr.testPingResp)

	// Setup core
	rg.GET("/getMac", h.getMac)
	rg.GET("/getSetting", h.getSetting)
	rg.POST("/savefileNew", h.savefileNew)
	rg.GET("/apply", h.apply)
	rg.GET("/releaseSetupMode", h.releaseSetupMode)
	rg.GET("/checkSetupMode", h.checkSetupMode)
	rg.GET("/getMode", h.getMode)
	rg.GET("/checkSettingFile", h.checkSettingFile)
	rg.GET("/getSettingData/:channel", h.getSettingData)
	rg.GET("/checkIbsm", h.checkIbsm)

	// Network
	h.registerNetworkRoutes(rg)
	// Service control
	h.registerServiceRoutes(rg)
	// Backup / Restore / Upload
	h.registerBackupRoutes(rg)
	// Parquet backup endpoints + setup.json download
	h.registerParquetRoutes(rg)
	// Bearing DB (SQLite)
	h.registerBearingRoutes(rg)
	// Misc: savefile/channel, restart, harm trigger, setSystemTime
	h.registerMiscRoutes(rg)
	// MQTT / Certs
	h.registerMQTTRoutes(rg)
}

func (h *Handler) getMac(c *gin.Context) {
	mac := getMAC()
	c.JSON(http.StatusOK, gin.H{"success": true, "mac": mac})
}

func (h *Handler) getSetting(c *gin.Context) {
	ctx := context.Background()
	client := h.deps.Redis.Client0

	var setting map[string]interface{}

	if setupJSON, err := client.HGet(ctx, "System", "setup").Result(); err == nil {
		json.Unmarshal([]byte(setupJSON), &setting)
	} else {
		setupPath := filepath.Join(h.deps.Config.ConfigDir, "setup.json")
		data, err := os.ReadFile(setupPath)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"passOK": 0, "error": "setting file not found"})
			return
		}
		if err := json.Unmarshal(data, &setting); err != nil {
			c.JSON(http.StatusOK, gin.H{"passOK": 0})
			return
		}
	}

	// Determine setup mode
	setupMode := 1
	if localVal, err := client.HGet(ctx, "System", "setupLocal").Result(); err == nil {
		if v, _ := strconv.Atoi(localVal); v == 1 {
			setupMode = 0
		}
	}

	if setupMode == 1 {
		client.HSet(ctx, "System", "setupRemote", 1)
	} else {
		client.HSet(ctx, "System", "setupRemote", 0)
	}

	// Convert ctInfo.inorminal from int to float (/1000)
	if channels, ok := setting["channel"].([]interface{}); ok {
		for _, ch := range channels {
			if chMap, ok := ch.(map[string]interface{}); ok {
				if ct, ok := chMap["ctInfo"].(map[string]interface{}); ok {
					if inorm, ok := ct["inorminal"]; ok {
						ct["inorminal"] = toFloat(inorm) / 1000.0
					}
				}
			}
		}
	}

	// Build response: split channels into main/sub/ibsm/ipsm72
	setupDict := map[string]interface{}{
		"mode":      setting["mode"],
		"lang":      setting["lang"],
		"setupMode": setupMode,
		"General":   setting["General"],
		"main":      map[string]interface{}{},
		"sub":       map[string]interface{}{},
		"ibsm":      map[string]interface{}{"channel": "ibsm", "Enable": 1, "tapboxs": []interface{}{}},
		"ipsm72":    map[string]interface{}{"channel": "ipsm72", "Enable": 1, "ipsm72_1": 0, "ipsm72_2": 0, "di60_1": 0, "di60_2": 0},
	}

	if channels, ok := setting["channel"].([]interface{}); ok {
		for _, ch := range channels {
			if chMap, ok := ch.(map[string]interface{}); ok {
				name, _ := chMap["channel"].(string)
				switch name {
				case "Main":
					setupDict["main"] = chMap
				case "Sub":
					setupDict["sub"] = chMap
				case "ibsm":
					setupDict["ibsm"] = chMap
				case "ipsm72":
					setupDict["ipsm72"] = chMap
				}
			}
		}
	}

	// channel에 ibsm이 없으면 Redis System > ibsmSet에서 가져오기
	if ibsmData, ok := setupDict["ibsm"].(map[string]interface{}); ok {
		if tapboxs, ok := ibsmData["tapboxs"].([]interface{}); ok && len(tapboxs) == 0 {
			if ibsmJSON, err := client.HGet(ctx, "System", "ibsmSet").Result(); err == nil {
				var ibsm map[string]interface{}
				if json.Unmarshal([]byte(ibsmJSON), &ibsm) == nil {
					ibsm["channel"] = "ibsm"
					if _, ok := ibsm["Enable"]; !ok {
						ibsm["Enable"] = 1
					}
					setupDict["ibsm"] = ibsm
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"passOK": 1, "data": setupDict})
}

func (h *Handler) savefileNew(c *gin.Context) {
	ctx := context.Background()
	client := h.deps.Redis.Client0

	// Check if modbus setting is active
	if flag, err := client.HGet(ctx, "Service", "setting").Result(); err == nil {
		if v, _ := strconv.Atoi(flag); v == 1 {
			c.JSON(http.StatusOK, gin.H{"status": "0", "error": "Modbus setting is activated"})
			return
		}
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "0", "error": "No data provided"})
		return
	}

	if _, ok := data["General"]; !ok {
		data["General"] = map[string]interface{}{}
	}
	if _, ok := data["channel"]; !ok {
		data["channel"] = []interface{}{}
	}

	// Set MAC address
	mac := getMAC()
	if gen, ok := data["General"].(map[string]interface{}); ok {
		if di, ok := gen["deviceInfo"].(map[string]interface{}); ok {
			di["mac_address"] = mac
		}
	}

	// Convert ctInfo.inorminal to int (*1000)
	if channels, ok := data["channel"].([]interface{}); ok {
		for _, ch := range channels {
			if chMap, ok := ch.(map[string]interface{}); ok {
				if ct, ok := chMap["ctInfo"].(map[string]interface{}); ok {
					if inorm, ok := ct["inorminal"]; ok {
						ct["inorminal"] = int(toFloat(inorm) * 1000)
					}
				}
			}
		}
	}

	out, _ := json.Marshal(data)
	client.HSet(ctx, "System", "config", string(out))
	c.JSON(http.StatusOK, gin.H{"status": "1"})
}

func (h *Handler) apply(c *gin.Context) {
	ctx := context.Background()
	client := h.deps.Redis.Client0

	// Check modbus lock
	if flag, err := client.HGet(ctx, "Service", "setting").Result(); err == nil {
		if v, _ := strconv.Atoi(flag); v == 1 {
			c.JSON(http.StatusOK, gin.H{"status": "0", "error": "Modbus setting is activated"})
			return
		}
	}

	// Load saved config and previous setup
	saveSetup, err := client.HGet(ctx, "System", "config").Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "0", "error": "No config data"})
		return
	}
	var saveData map[string]interface{}
	json.Unmarshal([]byte(saveSetup), &saveData)

	prevSetup, _ := client.HGet(ctx, "System", "setup").Result()
	var prevData map[string]interface{}
	json.Unmarshal([]byte(prevSetup), &prevData)

	// Write to file
	setupPath := filepath.Join(h.deps.Config.ConfigDir, "setup.json")
	out, _ := json.MarshalIndent(saveData, "", "    ")
	os.WriteFile(setupPath, out, 0644)

	// Process and save to Redis
	SaveRedisSetup(ctx, h.deps, saveData)

	// Compare changes
	result := compareChannelChanges(prevData, saveData)

	restartDevice := false
	if getStatus(result, "General") == "config_changed" ||
		getStatus(result, "Main") == "config_changed" ||
		getStatus(result, "Sub") == "config_changed" {
		restartDevice = true
	}

	// ibsm 분리: channel 배열에서 ibsm 추출하여 Redis에 별도 저장
	if channels, ok := saveData["channel"].([]interface{}); ok {
		for _, ch := range channels {
			if chMap, ok := ch.(map[string]interface{}); ok {
				if name, _ := chMap["channel"].(string); name == "ibsm" {
					ibsmJSON, _ := json.Marshal(chMap)
					client.HSet(ctx, "System", "ibsmSet", string(ibsmJSON))
					break
				}
			}
		}
	}

	// Update Redis setup
	setupJSON, _ := json.Marshal(saveData)
	client.HSet(ctx, "System", "setup", string(setupJSON))

	// Release setup mode
	if exists, _ := client.HExists(ctx, "System", "setupRemote").Result(); exists {
		client.HSet(ctx, "System", "setupRemote", 0)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        "1",
		"data":          saveData,
		"restartDevice": restartDevice,
	})
}

func (h *Handler) releaseSetupMode(c *gin.Context) {
	ctx := context.Background()
	h.deps.Redis.Client0.HSet(ctx, "System", "setupRemote", 0)
	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (h *Handler) checkSetupMode(c *gin.Context) {
	ctx := context.Background()
	client := h.deps.Redis.Client0

	if localVal, err := client.HGet(ctx, "System", "setupLocal").Result(); err == nil {
		if v, _ := strconv.Atoi(localVal); v == 1 {
			c.JSON(http.StatusOK, gin.H{"status": false})
			return
		}
	}

	client.HSet(ctx, "System", "setupRemote", 1)
	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (h *Handler) getMode(c *gin.Context) {
	ctx := context.Background()
	mode := ""
	if m, err := h.deps.Redis.Client0.HGet(ctx, "System", "mode").Result(); err == nil {
		mode = m
	}
	c.JSON(http.StatusOK, gin.H{"mode": mode})
}

func (h *Handler) checkSettingFile(c *gin.Context) {
	ctx := context.Background()
	client := h.deps.Redis.Client0

	var setting map[string]interface{}

	// Redis에서 읽기, 없으면 파일에서 읽기
	if setupJSON, err := client.HGet(ctx, "System", "setup").Result(); err == nil {
		json.Unmarshal([]byte(setupJSON), &setting)
	} else {
		setupPath := filepath.Join(h.deps.Config.ConfigDir, "setup.json")
		data, err := os.ReadFile(setupPath)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"result": "0", "error": "setting file not found"})
			return
		}
		json.Unmarshal(data, &setting)
		// Redis에 저장
		client.HSet(ctx, "System", "setup", string(data))
		if mode, ok := setting["mode"].(string); ok {
			client.HSet(ctx, "System", "mode", mode)
		}
	}

	result := gin.H{
		"result":      "1",
		"mode":        setting["mode"],
		"lang":        setting["lang"],
		"location":    "",
		"enable_main": false,
		"enable_sub":  false,
		"main_kva":    -1,
		"sub_kva":     -1,
		"pf_sign":     -1,
		"va_type":     -1,
		"unbalance":   -1,
	}

	channels, ok := setting["channel"].([]interface{})
	if !ok {
		result["result"] = "0"
		c.JSON(http.StatusOK, result)
		return
	}

	modeStr, _ := setting["mode"].(string)

	if general, ok := setting["General"].(map[string]interface{}); ok && len(modeStr) > 0 {
		if devInfo, ok := general["deviceInfo"].(map[string]interface{}); ok {
			result["location"] = devInfo["location"]
		}
	}

	for _, ch := range channels {
		chMap, ok := ch.(map[string]interface{})
		if !ok {
			continue
		}
		chName, _ := chMap["channel"].(string)
		if chName == "Main" {
			result["enable_main"] = toBool(chMap["Enable"])
			if opt, ok := chMap["opt"].(map[string]interface{}); ok {
				result["pf_sign"] = opt["pf_sign"]
				result["va_type"] = opt["va_type"]
				result["unbalance"] = opt["unbalance"]
			}
		} else if chName == "Sub" {
			result["enable_sub"] = toBool(chMap["Enable"])
		}
	}

	c.JSON(http.StatusOK, result)
}


func (h *Handler) checkIbsm(c *gin.Context) {
	ctx := context.Background()
	client := h.deps.Redis.Client0

	ibsmJSON, err := client.HGet(ctx, "System", "ibsmSet").Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": "0", "error": "ibsmSet not found"})
		return
	}

	var ibsmSet map[string]interface{}
	if err := json.Unmarshal([]byte(ibsmJSON), &ibsmSet); err != nil {
		c.JSON(http.StatusOK, gin.H{"result": "0", "error": "ibsmSet parse error"})
		return
	}

	tapboxs, ok := ibsmSet["tapboxs"].([]interface{})
	if !ok {
		c.JSON(http.StatusOK, gin.H{"result": "1", "tapboxCount": 0, "tapboxs": []interface{}{}})
		return
	}

	type tapboxInfo struct {
		Index   int `json:"index"`
		CbCount int `json:"cbcount"`
	}

	infos := make([]tapboxInfo, 0, len(tapboxs))
	for _, tb := range tapboxs {
		if tbMap, ok := tb.(map[string]interface{}); ok {
			infos = append(infos, tapboxInfo{
				Index:   toInt(tbMap["index"]),
				CbCount: toInt(tbMap["cbcount"]),
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"result":       "1",
		"tapboxCount":  len(tapboxs),
		"tapboxs":      infos,
	})
}

func toBool(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case float64:
		return val != 0
	case int:
		return val != 0
	case string:
		return val == "1" || val == "true"
	}
	return false
}

func toInt(v interface{}) int {
	switch val := v.(type) {
	case float64:
		return int(val)
	case int:
		return val
	case string:
		n, _ := strconv.Atoi(val)
		return n
	}
	return 0
}

func (h *Handler) getSettingData(c *gin.Context) {
	channel := c.Param("channel")
	ctx := context.Background()

	var setting map[string]interface{}
	if setupJSON, err := h.deps.Redis.Client0.HGet(ctx, "System", "setup").Result(); err == nil {
		json.Unmarshal([]byte(setupJSON), &setting)
	} else {
		setupPath := filepath.Join(h.deps.Config.ConfigDir, "setup.json")
		data, _ := os.ReadFile(setupPath)
		json.Unmarshal(data, &setting)
	}

	if channels, ok := setting["channel"].([]interface{}); ok {
		for _, ch := range channels {
			if chMap, ok := ch.(map[string]interface{}); ok {
				if ct, ok := chMap["ctInfo"].(map[string]interface{}); ok {
					if inorm, ok := ct["inorminal"]; ok {
						ct["inorminal"] = toFloat(inorm) / 1000.0
					}
				}
				name, _ := chMap["channel"].(string)
				if name == channel {
					c.JSON(http.StatusOK, gin.H{"passOK": 1, "data": chMap})
					return
				}
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"passOK": 0, "error": "channel not found"})
}

// compareChannelChanges compares previous and new setup data
func compareChannelChanges(redisData, postData map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"General": map[string]interface{}{"status": "no_change", "changed_fields": []string{}},
		"Main":    map[string]interface{}{"status": "disabled", "changed_fields": []string{}},
		"Sub":     map[string]interface{}{"status": "disabled", "changed_fields": []string{}},
	}

	// Compare General fields
	redisGeneral, _ := redisData["General"].(map[string]interface{})
	postGeneral, _ := postData["General"].(map[string]interface{})

	generalFields := []string{"etc", "deviceInfo", "tcpip", "modbus", "useFunction", "ftpInfo", "sntpInfo"}
	var generalChanged []string
	for _, field := range generalFields {
		rv := toJSON(redisGeneral[field])
		pv := toJSON(postGeneral[field])
		if rv != pv {
			generalChanged = append(generalChanged, field)
		}
	}
	if len(generalChanged) > 0 {
		result["General"] = map[string]interface{}{"status": "config_changed", "changed_fields": generalChanged}
	}

	// Build channel maps
	redisChannels := buildChannelMap(redisData)
	postChannels := buildChannelMap(postData)

	channelFields := []string{"Enable", "ctInfo", "ptInfo", "eventInfo", "demand", "trendInfo", "alarm", "opt"}

	for _, chName := range []string{"Main", "Sub"} {
		redisCh := redisChannels[chName]
		postCh := postChannels[chName]

		if toFloat(postCh["Enable"]) != 1 {
			continue
		}

		var changed []string
		for _, field := range channelFields {
			rv := toJSON(redisCh[field])
			pv := toJSON(postCh[field])
			if rv != pv {
				changed = append(changed, field)
			}
		}

		status := "no_change"
		if len(changed) > 0 {
			status = "config_changed"
		}
		result[chName] = map[string]interface{}{"status": status, "changed_fields": changed}
	}

	return result
}

func buildChannelMap(data map[string]interface{}) map[string]map[string]interface{} {
	result := map[string]map[string]interface{}{
		"Main": {},
		"Sub":  {},
	}
	channels, _ := data["channel"].([]interface{})
	for _, ch := range channels {
		chMap, ok := ch.(map[string]interface{})
		if !ok {
			continue
		}
		name, _ := chMap["channel"].(string)
		if name == "Main" || name == "Sub" {
			result[name] = chMap
		}
	}
	return result
}

func getStatus(result map[string]interface{}, key string) string {
	if r, ok := result[key].(map[string]interface{}); ok {
		if s, ok := r["status"].(string); ok {
			return s
		}
	}
	return ""
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
	case json.Number:
		f, _ := val.Float64()
		return f
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	}
	return 0
}

func toJSON(v interface{}) string {
	if v == nil {
		return "null"
	}
	b, _ := json.Marshal(v)
	return string(b)
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
