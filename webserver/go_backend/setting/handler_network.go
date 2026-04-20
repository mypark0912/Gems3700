package setting

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) registerNetworkRoutes(rg *gin.RouterGroup) {
	rg.GET("/applyNetwork", h.applyNetwork)
	rg.POST("/setDefaultIP", h.setDefaultIP)
}

func (h *Handler) applyNetwork(c *gin.Context) {
	ctx := context.Background()
	setupJSON, err := h.deps.Redis.Client0.HGet(ctx, "System", "setup").Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": false})
		return
	}

	var setupData map[string]interface{}
	json.Unmarshal([]byte(setupJSON), &setupData)

	gen, _ := setupData["General"].(map[string]interface{})
	netData, _ := gen["tcpip"].(map[string]interface{})
	sntpData, _ := gen["sntpInfo"].(map[string]interface{})

	// Run network + SNTP apply in background
	go applyNetworkSetting(netData)
	go applySNTPSetting(sntpData)

	c.JSON(http.StatusOK, gin.H{"result": true, "message": "적용 중..."})
}

func (h *Handler) setDefaultIP(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}

	defaultPath := h.deps.Config.ConfigDir + "/default.json"
	raw, err := os.ReadFile(defaultPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}

	var defaults map[string]interface{}
	json.Unmarshal(raw, &defaults)

	if gen, ok := defaults["General"].(map[string]interface{}); ok {
		if tcpip, ok := gen["tcpip"].(map[string]interface{}); ok {
			tcpip["ip_address"] = data["ip"]
		}
	}

	out, _ := json.MarshalIndent(defaults, "", "  ")
	os.WriteFile(defaultPath, out, 0644)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// applyNetworkSetting configures network interface (Linux systemd-networkd)
func applyNetworkSetting(netData map[string]interface{}) {
	const networkFile = "/etc/systemd/network/10-static-end1.network"
	const iface = "end1"

	dhcp := int(toFloat(netData["dhcp"]))
	ip, _ := netData["ip_address"].(string)
	mask, _ := netData["subnet_mask"].(string)
	gateway, _ := netData["gateway"].(string)
	cidr := maskToCIDR(mask)

	dhcpTemplate := `[Match]
Name=end1
[Network]
DHCP=ipv4
[DHCP]
UseDNS=true
UseRoutes=true
`

	staticTemplate := fmt.Sprintf(`[Match]
Name=end1
[Network]
#DHCP=ipv4
Address=%s/%d
Gateway=%s
[DHCP]
UseDNS=true
UseRoutes=true
`, ip, cidr, gateway)

	var newContent string
	if dhcp == 1 {
		newContent = dhcpTemplate
	} else {
		newContent = staticTemplate
	}

	// Check if unchanged
	if current, err := os.ReadFile(networkFile); err == nil {
		if strings.TrimSpace(string(current)) == strings.TrimSpace(newContent) {
			return
		}
	}

	os.WriteFile(networkFile, []byte(newContent), 0644)

	exec.Command("systemctl", "stop", "frpc-restart-monitor").Run()
	time.Sleep(500 * time.Millisecond)
	exec.Command("systemctl", "restart", "systemd-networkd").Run()

	if dhcp == 1 {
		for i := 0; i < 10; i++ {
			time.Sleep(2 * time.Second)
			if cur := getCurrentIP(iface); cur != nil {
				exec.Command("systemctl", "start", "frpc-restart-monitor").Run()
				return
			}
		}
		// DHCP failed, fallback to static
		os.WriteFile(networkFile, []byte(staticTemplate), 0644)
		exec.Command("systemctl", "restart", "systemd-networkd").Run()
		time.Sleep(3 * time.Second)
		exec.Command("systemctl", "start", "frpc-restart-monitor").Run()
	} else {
		time.Sleep(3 * time.Second)
		exec.Command("systemctl", "start", "frpc-restart-monitor").Run()
	}
}

// applySNTPSetting configures timezone and NTP
func applySNTPSetting(sntpData map[string]interface{}) {
	tz, _ := sntpData["timezone"].(string)
	ntpServer, _ := sntpData["host"].(string)

	if tz != "" {
		out, _ := exec.Command("timedatectl", "show", "--property=Timezone").Output()
		currentTZ := strings.TrimSpace(strings.TrimPrefix(string(out), "Timezone="))
		if currentTZ != tz {
			exec.Command("timedatectl", "set-timezone", tz).Run()
		}
	}

	if ntpServer == "" {
		// Disable NTP
		out, _ := exec.Command("systemctl", "is-enabled", "systemd-timesyncd").Output()
		if strings.TrimSpace(string(out)) == "enabled" {
			exec.Command("systemctl", "stop", "systemd-timesyncd").Run()
			exec.Command("systemctl", "disable", "systemd-timesyncd").Run()
		}
		return
	}

	const confPath = "/etc/systemd/timesyncd.conf"
	newContent := fmt.Sprintf("[Time]\nNTP=%s\n", ntpServer)

	ntpChanged := false
	if current, err := os.ReadFile(confPath); err != nil || strings.TrimSpace(string(current)) != strings.TrimSpace(newContent) {
		os.WriteFile(confPath, []byte(newContent), 0644)
		ntpChanged = true
	}

	if ntpChanged {
		exec.Command("systemctl", "enable", "systemd-timesyncd").Run()
		exec.Command("systemctl", "restart", "systemd-timesyncd").Run()
	}
}

func maskToCIDR(mask string) int {
	cidr := 0
	for _, octet := range strings.Split(mask, ".") {
		v, _ := strconv.Atoi(octet)
		for v > 0 {
			cidr += v & 1
			v >>= 1
		}
	}
	return cidr
}

func getCurrentIP(iface string) map[string]string {
	out, err := exec.Command("ip", "addr", "show", iface).Output()
	if err != nil {
		return nil
	}
	re := regexp.MustCompile(`inet (\d+\.\d+\.\d+\.\d+)/(\d+)`)
	match := re.FindStringSubmatch(string(out))
	if len(match) >= 3 {
		return map[string]string{"ip": match[1], "cidr": match[2]}
	}
	return nil
}
