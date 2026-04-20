package setting

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

var serviceMap = map[string]string{
	"Redis":                "redis",
	"InfluxDB":             "influxdb",
	"Core":                 "core",
	"WebServer":            "webserver",
	"A35":                  "sv500A35",
	"MQTTClient":           "mqClient",
	"frpc":                 "frpc",
	"frpc-restart-monitor": "frpc-restart-monitor",
}

func (h *Handler) registerServiceRoutes(rg *gin.RouterGroup) {
	rg.GET("/SysService/:cmd/:item", h.sysService)
	rg.GET("/SysCheck", h.sysCheck)
	rg.GET("/ServiceStatus", h.serviceStatus)
	rg.POST("/command", h.pushCommand)
	rg.GET("/checkFrp", h.checkFrp)
}

func (h *Handler) sysService(c *gin.Context) {
	cmd := c.Param("cmd")
	item := c.Param("item")
	result := runSysService(cmd, item)
	c.JSON(http.StatusOK, result)
}

func (h *Handler) sysCheck(c *gin.Context) {
	// Version info
	versionDict := getVersions()
	versions := map[string]string{}
	if v, ok := versionDict["web"]; ok {
		versions["WebServer"] = v
	}

	// Service status
	services := map[string]string{
		"redis":     "redis",
		"influxdb":  "influxdb",
		"webserver": "webserver",
	}
	serviceStatus := map[string]bool{}
	for key, name := range services {
		serviceStatus[key] = isServiceActive(name)
	}

	// Disk usage
	disk1 := getDiskUsage("/")
	disk2 := getDiskUsage("/usr/local")

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"data":     serviceStatus,
		"disk":     []map[string]interface{}{disk1, disk2},
		"versions": versions,
	})
}

func (h *Handler) serviceStatus(c *gin.Context) {
	ctx := context.Background()
	services := map[string]string{
		"redis":     "redis",
		"influxdb":  "influxdb",
		"core":      "core",
		"webserver": "webserver",
		"a35":       "a35",
	}

	serviceStatus := map[string]bool{}
	abnormal := false
	for key, name := range services {
		active := isServiceActive(name)
		serviceStatus[key] = active
		if !active {
			abnormal = true
		}
	}

	statusList := getSystemStatus()

	result := gin.H{}
	if abnormal {
		result["service"] = false
		result["serviceDict"] = serviceStatus
	} else {
		result["service"] = true
	}

	if len(statusList) > 0 {
		result["data"] = statusList
		result["system"] = false
	} else {
		result["system"] = true
	}

	_ = ctx
	c.JSON(http.StatusOK, result)
}

// Command model
type CommandRequest struct {
	Type int `json:"type"`
	Cmd  int `json:"cmd"`
	Item int `json:"item"`
}

func (h *Handler) pushCommand(c *gin.Context) {
	var cmd CommandRequest
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Pack as 3 x int32 little-endian (12 bytes)
	buf := make([]byte, 12)
	binary.LittleEndian.PutUint32(buf[0:4], uint32(cmd.Type))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(cmd.Cmd))
	binary.LittleEndian.PutUint32(buf[8:12], uint32(cmd.Item))

	ctx := context.Background()
	h.deps.Redis.Client0.LPush(ctx, "command", buf)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Command pushed to left of list",
		"data":    cmd,
	})
}

func (h *Handler) checkFrp(c *gin.Context) {
	if !serviceExists("frpc.service") {
		c.JSON(http.StatusOK, gin.H{"exist": false, "status": false})
		return
	}
	active := isServiceActive("frpc")
	c.JSON(http.StatusOK, gin.H{"exist": true, "status": active})
}

// helpers

func runSysService(cmd, item string) gin.H {
	serviceName, ok := serviceMap[item]
	if !ok {
		return gin.H{"success": false, "error": "unknown service: " + item}
	}

	result, err := exec.Command("sudo", "systemctl", cmd, serviceName).CombinedOutput()
	if err != nil {
		return gin.H{
			"success":    false,
			"error":      string(result),
			"service":    serviceName,
			"action":     cmd,
			"returncode": -1,
		}
	}
	return gin.H{
		"success":    true,
		"stdout":     strings.TrimSpace(string(result)),
		"service":    serviceName,
		"action":     cmd,
		"returncode": 0,
	}
}

func isServiceActive(name string) bool {
	out, err := exec.Command("systemctl", "is-active", name).Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "active"
}

func serviceExists(name string) bool {
	err := exec.Command("systemctl", "status", name).Run()
	// status returns 0 for active, 3 for inactive but exists; 4 for not found
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode() != 4
		}
		return false
	}
	return true
}

func isServiceEnabled(name string) bool {
	out, err := exec.Command("systemctl", "is-enabled", name).Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "enabled"
}

func getVersions() map[string]string {
	data, err := os.ReadFile("/home/root/versionInfo.txt")
	if err != nil {
		return nil
	}
	result := map[string]string{}
	for _, line := range strings.Split(string(data), "\n") {
		parts := strings.SplitN(strings.TrimSpace(line), "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

func getDiskUsage(path string) map[string]interface{} {
	out, err := exec.Command("df", "-B1", path).Output()
	if err != nil {
		return map[string]interface{}{"drive": path, "totalGB": 0, "freeGB": 0, "status": "error"}
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return map[string]interface{}{"drive": path, "totalGB": 0, "freeGB": 0, "status": "error"}
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 4 {
		return map[string]interface{}{"drive": path, "totalGB": 0, "freeGB": 0, "status": "error"}
	}

	total := parseBytes(fields[1])
	free := parseBytes(fields[3])
	percent := 0.0
	if total > 0 {
		percent = float64(total-free) / float64(total) * 100
	}

	status := "ok"
	if percent > 90 {
		status = "warning"
	}

	return map[string]interface{}{
		"drive":   path,
		"totalGB": math.Round(float64(total)/1073741824*100) / 100,
		"freeGB":  math.Round(float64(free)/1073741824*100) / 100,
		"status":  status,
	}
}

func parseBytes(s string) int64 {
	var v int64
	fmt.Sscanf(s, "%d", &v)
	return v
}

func getSystemStatus() []string {
	var highList []string
	// Simple check using /proc or exec
	if out, err := exec.Command("sh", "-c", "cat /proc/meminfo | grep -E 'MemTotal|MemAvailable'").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		var total, avail int64
		for _, line := range lines {
			if strings.HasPrefix(line, "MemTotal:") {
				fmt.Sscanf(line, "MemTotal: %d", &total)
			}
			if strings.HasPrefix(line, "MemAvailable:") {
				fmt.Sscanf(line, "MemAvailable: %d", &avail)
			}
		}
		if total > 0 {
			memPercent := float64(total-avail) / float64(total) * 100
			if memPercent > 90 {
				highList = append(highList, "MEMORY")
			}
		}
	}
	return highList
}

func execService(cmd string) gin.H {
	out, err := exec.Command("sudo", "systemctl", cmd).CombinedOutput()
	if err != nil {
		return gin.H{"success": false, "error": string(out)}
	}
	return gin.H{"success": true}
}

// saveFrpcConfig writes frpc.toml
func saveFrpcConfig(subdomain, prefix string) bool {
	const filePath = "/home/root/frp_0.66.0_linux_arm64/frpc.toml"
	dir := filepath.Dir(filePath)
	os.MkdirAll(dir, 0755)

	content := fmt.Sprintf(`serverAddr = "13.125.5.143"
serverPort = 7000
auth.token = "NTEK_system_20260116_mypark"
transport.tls.enable = true

[[proxies]]
name = "device-web-%s"
type = "http"
localIP = "127.0.0.1"
localPort = 4000
subdomain = "%s"
`, prefix, subdomain)

	return os.WriteFile(filePath, []byte(content), 0644) == nil
}

func createServiceFile(name, desc, execStart, workDir, after, restart string) gin.H {
	content := fmt.Sprintf(`[Unit]
Description=%s
After=%s

[Service]
Type=simple
ExecStart=%s
WorkingDirectory=%s
Restart=%s
RestartSec=10
User=root

[Install]
WantedBy=multi-user.target
`, desc, after, execStart, workDir, restart)

	path := fmt.Sprintf("/etc/systemd/system/%s.service", name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return gin.H{"success": false, "error": err.Error()}
	}
	return gin.H{"success": true}
}

// Helper for MQTT background processing
func (h *Handler) applyMQTTServices(setup map[string]interface{}) {
	gen, _ := setup["General"].(map[string]interface{})
	mqtt, _ := gen["MQTT"].(map[string]interface{})
	if mqtt == nil {
		return
	}

	ctx := context.Background()
	useStr := fmt.Sprintf("%v", mqtt["Use"])
	if useStr == "1" {
		h.deps.Redis.Client0.HSet(ctx, "System", "MQTT", 1)

		// Save mqtt.json
		mqttPath := filepath.Join(h.deps.Config.ConfigDir, "mqtt.json")
		jsonData := map[string]interface{}{
			"host":      mqtt["host"],
			"port":      mqtt["port"],
			"device_id": mqtt["device_id"],
			"username":  mqtt["username"],
			"password":  mqtt["password"],
		}
		out, _ := json.MarshalIndent(jsonData, "", "  ")
		os.WriteFile(mqttPath, out, 0644)

		// Start/restart mqClient service
		if serviceExists("mqClient") {
			if isServiceEnabled("mqClient") {
				if isServiceActive("mqClient") {
					runSysService("restart", "MQTTClient")
				} else {
					runSysService("start", "MQTTClient")
				}
			} else {
				runSysService("enable", "MQTTClient")
				runSysService("start", "MQTTClient")
			}
		}
	} else {
		h.deps.Redis.Client0.HSet(ctx, "System", "MQTT", 0)
		if serviceExists("mqClient.service") {
			if isServiceActive("mqClient") {
				runSysService("stop", "MQTTClient")
			}
			if isServiceEnabled("mqClient") {
				runSysService("disable", "MQTTClient")
			}
		}
	}
}
