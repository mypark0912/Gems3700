package setting

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"serverGO/infra"
)

// InitSetup loads setup.json into Redis at boot (equivalent to Python init_setup)
func InitSetup(deps *infra.Dependencies) {
	ctx := context.Background()
	client := deps.Redis.Client0

	// Check if already loaded
	exists, _ := client.HExists(ctx, "System", "setup").Result()
	if exists {
		log.Println("System setup already exists in Redis, skipping init")
		return
	}

	configDir := deps.Config.ConfigDir
	setupPath := filepath.Join(configDir, "setup.json")
	defaultPath := filepath.Join(configDir, "default.json")

	// If setup.json doesn't exist, copy from default
	if _, err := os.Stat(setupPath); os.IsNotExist(err) {
		if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
			log.Println("setup.json and default.json both missing")
			return
		}
		data, _ := os.ReadFile(defaultPath)
		os.WriteFile(setupPath, data, 0644)
	}

	data, err := os.ReadFile(setupPath)
	if err != nil {
		log.Printf("Failed to read setup.json: %v", err)
		return
	}

	var setting map[string]interface{}
	if err := json.Unmarshal(data, &setting); err != nil {
		// Try fallback to default
		data, _ = os.ReadFile(defaultPath)
		os.WriteFile(setupPath, data, 0644)
		if err := json.Unmarshal(data, &setting); err != nil {
			log.Printf("Failed to parse setup.json: %v", err)
			return
		}
	}

	// Update MAC address
	mac := getMAC()
	if gen, ok := setting["General"].(map[string]interface{}); ok {
		if di, ok := gen["deviceInfo"].(map[string]interface{}); ok {
			if di["mac_address"] != mac {
				di["mac_address"] = mac
			}
		}
	}

	// Read serial number from file (Linux only)
	if deps.Config.OS != "Windows" && deps.Config.OS != "macOS" {
		serPath := filepath.Join(configDir, "serial_num_do_not_modify.txt")
		if serData, err := os.ReadFile(serPath); err == nil {
			ser := strings.TrimSpace(string(serData))
			if ser != "" {
				if gen, ok := setting["General"].(map[string]interface{}); ok {
					if di, ok := gen["deviceInfo"].(map[string]interface{}); ok {
						di["serial_number"] = ser
					}
				}
			}
		}
	}

	SaveRedisSetup(ctx, deps, setting)

	out, _ := json.Marshal(setting)
	client.HSet(ctx, "System", "setup", string(out))

	if mode, ok := setting["mode"].(string); ok {
		client.HSet(ctx, "System", "mode", mode)
	}

	// ibsm 분리 저장: channel 배열에서 ibsm 추출
	if channels, ok := setting["channel"].([]interface{}); ok {
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

	log.Println("Boot setup loaded into Redis")
}

// SaveRedisSetup processes setup data and saves derived values to Redis
func SaveRedisSetup(ctx context.Context, deps *infra.Dependencies, setupData map[string]interface{}) {
	channels, _ := setupData["channel"].([]interface{})
	if len(channels) == 0 {
		return
	}

	// Process channels for StartCurrent, Demand, ChannelData
	mainData := map[string]interface{}{}
	subData := map[string]interface{}{}

	for _, ch := range channels {
		chMap, ok := ch.(map[string]interface{})
		if !ok {
			continue
		}
		chName, _ := chMap["channel"].(string)
		switch chName {
		case "Main":
			mainData = processChannelData(chMap)
		case "Sub":
			subData = processChannelData(chMap)
		}
	}

	startCurrent := map[string]interface{}{"main": mainData["startCurrent"], "sub": subData["startCurrent"]}
	demand := map[string]interface{}{"main": mainData["demand"], "sub": subData["demand"]}
	channelInfo := map[string]interface{}{"main": mainData["info"], "sub": subData["info"]}

	hsetJSON(ctx, deps, "Equipment", "StartingCurrent", startCurrent)
	hsetJSON(ctx, deps, "Equipment", "DemandInterval", demand)
	hsetJSON(ctx, deps, "Equipment", "ChannelData", channelInfo)
}

func processChannelData(ch map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"startCurrent": 0,
		"demand":       0,
		"info":         ch,
	}

	if ct, ok := ch["ctInfo"].(map[string]interface{}); ok {
		if sc, ok := ct["startingcurrent"]; ok {
			result["startCurrent"] = sc
		}
	}

	if d, ok := ch["demand"].(map[string]interface{}); ok {
		if di, ok := d["demand_interval"]; ok {
			result["demand"] = di
		}
	}

	return result
}

func hsetJSON(ctx context.Context, deps *infra.Dependencies, hash, field string, value interface{}) {
	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("hsetJSON marshal error: %v", err)
		return
	}
	deps.Redis.Client0.HSet(ctx, hash, field, string(data))
}

func getMAC() string {
	ifaces, err := os.ReadDir("/sys/class/net")
	if err != nil {
		return "00:00:00:00:00:00"
	}
	for _, iface := range ifaces {
		name := iface.Name()
		if name == "sw0ep" || name == "end1" || name == "eth0" {
			data, err := os.ReadFile(fmt.Sprintf("/sys/class/net/%s/address", name))
			if err == nil {
				return strings.TrimSpace(string(data))
			}
		}
	}
	return "00:00:00:00:00:00"
}
