package data

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var bgCtx = context.Background()

// asMap safely extracts a map from an interface{}.
func asMap(v interface{}) map[string]interface{} {
	if m, ok := v.(map[string]interface{}); ok {
		return m
	}
	return nil
}

// asStr safely extracts a string from an interface{}.
func asStr(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// asFloat safely extracts a float64 from an interface{}.
func asFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	default:
		return 0
	}
}

// GetSetup reads System/setup from Redis, parses JSON, and creates Channel map.
// Matches Python: data/SetupInfo.py get_setup()
func GetSetup(redisClient *redis.Client) map[string]*Channel {
	chdict := make(map[string]*Channel)

	redisClient.Do(bgCtx, "SELECT", 0)
	setup, err := redisClient.HGet(bgCtx, "System", "setup").Result()
	if err != nil {
		log.Printf("Error getting setup from Redis: %v", err)
		return chdict
	}

	var setupInfo map[string]interface{}
	if err := json.Unmarshal([]byte(setup), &setupInfo); err != nil {
		log.Printf("Error parsing setup JSON: %v", err)
		return chdict
	}

	channelList, ok := setupInfo["channel"].([]interface{})
	if !ok {
		log.Println("No 'channel' array found in setup data")
		return chdict
	}

	general := asMap(setupInfo["General"])
	var useFunction map[string]interface{}
	if general != nil {
		useFunction = asMap(general["useFuction"])
	}

	for i := 0; i < len(channelList); i++ {
		chData := asMap(channelList[i])
		if chData == nil {
			continue
		}

		// Check Enable
		if toInt(chData["Enable"]) != 1 {
			continue
		}

		keyName := asStr(chData["channel"])
		if keyName == "" {
			continue
		}

		ch := NewChannel(keyName)

		// Diagnosis enable from General.useFuction
		if useFunction != nil {
			if keyName == "Main" {
				ch.Diagnosis = toInt(useFunction["diagnosis_main"]) == 1
			} else {
				ch.Diagnosis = toInt(useFunction["diagnosis_sub"]) == 1
			}
		}

		// Asset info (only if diagnosis enabled)
		if ch.Diagnosis {
			assetInfo := asMap(chData["assetInfo"])
			if assetInfo != nil {
				ch.AssetName = asStr(assetInfo["name"])
				ch.AssetType = asStr(assetInfo["type"])
				if asStr(assetInfo["driveType"]) == "VFD" {
					ch.AssetDrive = true
				}
			}
		}

		// confStatus
		if confStatus, ok := chData["confStatus"]; ok {
			ch.UseConfStatus = toInt(confStatus) == 1
		}

		// useDO
		if useDO, ok := chData["useDO"]; ok {
			ch.UseDO = toInt(useDO) == 1
		}

		// confStatus detail
		if ch.UseDO && ch.UseConfStatus {
			if statusInfo := asMap(chData["status_Info"]); statusInfo != nil {
				ch.ConfStatus = statusInfo
			}
		}

		// Trend info
		if trendInfo := asMap(chData["trendInfo"]); trendInfo != nil {
			params, ok := trendInfo["params"].([]interface{})
			if ok {
				nonCount := 0
				for _, p := range params {
					ps := asStr(p)
					if ps == "None" || ps == "" {
						nonCount++
					} else {
						ch.TrendList = append(ch.TrendList, ps)
					}
				}
				if len(params) != nonCount {
					ch.Period = toInt(trendInfo["period"]) * 60
					ch.Trend = true
				} else {
					ch.Trend = false
				}
			}
		}

		// Voltage type from ptInfo.wiringmode
		if ptInfo := asMap(chData["ptInfo"]); ptInfo != nil {
			if toInt(ptInfo["wiringmode"]) == 0 {
				ch.VoltageType = 0 // 3P4W -> phase voltage
			} else {
				ch.VoltageType = 1 // 3P3W -> line voltage
			}
		}

		// Demand settings
		if demand := asMap(chData["demand"]); demand != nil {
			ch.DemandPeriod = toInt(demand["demand_interval"])
			if collect, ok := demand["collect"]; ok {
				ch.DemandTrend = toInt(collect)
			}
		}

		// Alarm settings
		if alarm := asMap(chData["alarm"]); alarm != nil {
			nonCount := 0
			for j := 1; j <= 32; j++ {
				keyStr := fmt.Sprintf("%d", j)
				if arr, ok := alarm[keyStr].([]interface{}); ok && len(arr) > 0 {
					if asFloat(arr[0]) == 0 {
						nonCount++
					}
				} else {
					nonCount++
				}
			}
			if nonCount != 32 {
				ch.Alarm = true
				if delay, ok := alarm["CompareTimeDelay"]; ok {
					ch.AlarmComDelay = toInt(delay)
				}
			} else {
				ch.Alarm = false
			}
		}

		chdict[keyName] = ch
	}

	return chdict
}

// GetFtp reads FTP/system configuration from Redis.
// Matches Python: data/SetupInfo.py getFtp()
func GetFtp(redisClient *redis.Client) map[string]interface{} {
	redisClient.Do(bgCtx, "SELECT", 0)
	setup, err := redisClient.HGet(bgCtx, "System", "setup").Result()
	if err != nil {
		log.Printf("Error getting setup from Redis: %v", err)
		return map[string]interface{}{}
	}

	var setupInfo map[string]interface{}
	if err := json.Unmarshal([]byte(setup), &setupInfo); err != nil {
		log.Printf("Error parsing setup JSON: %v", err)
		return map[string]interface{}{}
	}

	general := asMap(setupInfo["General"])
	if general == nil {
		return map[string]interface{}{"mode": asStr(setupInfo["mode"])}
	}

	useFunction := asMap(general["useFuction"])
	ftpEnabled := false
	if useFunction != nil {
		ftpEnabled = toInt(useFunction["ftp"]) == 1
	}

	mainEnabled, subEnabled := CheckFtpFromSetup(setupInfo)
	timingDict := GetTimingFromSetup(setupInfo)

	return map[string]interface{}{
		"mode":       asStr(setupInfo["mode"]),
		"ftp":        ftpEnabled,
		"ftpInfo":    asMap(general["ftpInfo"]),
		"deviceInfo": asMap(general["deviceInfo"]),
		"main":       mainEnabled,
		"sub":        subEnabled,
		"sampling":   timingDict,
	}
}

// CheckFtpFromSetup checks which channels are enabled.
func CheckFtpFromSetup(setupInfo map[string]interface{}) (mainEnabled, subEnabled bool) {
	channelList, ok := setupInfo["channel"].([]interface{})
	if !ok {
		return false, false
	}

	for _, chRaw := range channelList {
		ch := asMap(chRaw)
		if ch == nil {
			continue
		}
		chName := asStr(ch["channel"])
		enabled := toInt(ch["Enable"]) == 1
		if chName == "Main" {
			mainEnabled = enabled
		} else if chName == "Sub" {
			subEnabled = enabled
		}
	}
	return
}

// GetTimingFromSetup extracts sampling timing per channel.
func GetTimingFromSetup(setupInfo map[string]interface{}) map[string]int {
	timeDict := make(map[string]int)
	channelList, ok := setupInfo["channel"].([]interface{})
	if !ok {
		return timeDict
	}

	for _, chRaw := range channelList {
		ch := asMap(chRaw)
		if ch == nil {
			continue
		}
		chName := asStr(ch["channel"])
		if toInt(ch["Enable"]) != 1 {
			continue
		}
		if sampling := asMap(ch["sampling"]); sampling != nil {
			timeDict[chName] = toInt(sampling["period"]) * 60
		}
	}
	return timeDict
}
