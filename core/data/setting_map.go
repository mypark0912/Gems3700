package data

import (
	"fmt"
	"log"
	"strings"
)

// ensureMap safely retrieves a nested map[string]interface{} from a parent map.
// If the key does not exist, it creates and inserts an empty map.
func ensureMap(parent map[string]interface{}, key string) map[string]interface{} {
	if v, ok := parent[key]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			return m
		}
	}
	m := make(map[string]interface{})
	parent[key] = m
	return m
}

// getSlice safely retrieves a []interface{} from a map.
func getSlice(m map[string]interface{}, key string) []interface{} {
	if v, ok := m[key]; ok {
		if s, ok := v.([]interface{}); ok {
			return s
		}
	}
	return nil
}

// getMapVal safely retrieves a value from a nested map path like settings["comm"]["ip0"].
func getMapVal(m map[string]interface{}, keys ...string) (interface{}, bool) {
	current := m
	for i, key := range keys {
		if i == len(keys)-1 {
			v, ok := current[key]
			return v, ok
		}
		if v, ok := current[key]; ok {
			if next, ok := v.(map[string]interface{}); ok {
				current = next
			} else {
				return nil, false
			}
		} else {
			return nil, false
		}
	}
	return nil, false
}

// joinIPParts converts a slice of interface{} values (numbers or strings) into a dotted IP string.
func joinIPParts(parts interface{}) string {
	if s, ok := parts.([]interface{}); ok {
		strs := make([]string, len(s))
		for i, v := range s {
			strs[i] = fmt.Sprintf("%v", v)
		}
		return strings.Join(strs, ".")
	}
	return ""
}

// UpdateGeneralFromSettings merges general settings into jsonData for the given channel.
// Python equivalent: SettingMap.update_general_from_settings
func UpdateGeneralFromSettings(jsonData, settings map[string]interface{}, channelName string) map[string]interface{} {
	if jsonData == nil {
		jsonData = make(map[string]interface{})
	}

	general := ensureMap(jsonData, "General")

	if channelName == "Main" {
		// COMM_CFG data
		if comm, ok := settings["comm"]; ok {
			commMap := comm.(map[string]interface{})

			useFunction := ensureMap(general, "useFuction")

			if v, ok := commMap["ftpEnable"]; ok {
				useFunction["ftp"] = v
			}
			if v, ok := commMap["sntpEnable"]; ok {
				useFunction["sntp"] = v
			}

			tcpip := ensureMap(general, "tcpip")

			if v, ok := commMap["ip0"]; ok {
				tcpip["ip_address"] = joinIPParts(v)
			}
			if v, ok := commMap["sm0"]; ok {
				tcpip["subnet_mask"] = joinIPParts(v)
			}
			if v, ok := commMap["gw0"]; ok {
				tcpip["gateway"] = joinIPParts(v)
			}
			if v, ok := commMap["dns0"]; ok {
				tcpip["dnsserver"] = joinIPParts(v)
			}

			// Modbus settings
			modbus := ensureMap(general, "modbus")

			if v, ok := commMap["rs485Enable"]; ok {
				modbus["rtu_use"] = toInt(v)
			}
			if v, ok := commMap["baud"]; ok {
				modbus["baud_rate"] = toInt(v)
			}
			if v, ok := commMap["devId"]; ok {
				modbus["modbus_id"] = toInt(v)
			}
			if v, ok := commMap["tcpPort"]; ok {
				modbus["tcp_port"] = toInt(v)
			}
			if v, ok := commMap["parity"]; ok {
				modbus["parity"] = toInt(v)
			}
		}

		// ETC_DEF data
		if etc, ok := settings["etc"]; ok {
			etcMap := etc.(map[string]interface{})

			if v, ok := etcMap["VA_type"]; ok {
				general["va_type"] = v
			}
			if v, ok := etcMap["PF_sign"]; ok {
				general["pf_sign"] = v
			}

			// Timezone conversion (minutes -> timezone string)
			if v, ok := etcMap["timezone"]; ok {
				sntpInfo := ensureMap(general, "sntpInfo")
				tzMinutes := toInt(v)
				switch tzMinutes {
				case 540: // UTC+9
					sntpInfo["timezone"] = "Asia/Seoul"
				case 0:
					sntpInfo["timezone"] = "UTC"
				case -300: // UTC-5
					sntpInfo["timezone"] = "America/New_York"
				case -480: // UTC-8
					sntpInfo["timezone"] = "America/Los_Angeles"
				default:
					hours := tzMinutes / 60
					sntpInfo["timezone"] = fmt.Sprintf("UTC%+d", hours)
				}
			}
		}

		// SNTP_DEF data
		if sntp, ok := settings["sntp"]; ok {
			sntpMap := sntp.(map[string]interface{})
			sntpInfo := ensureMap(general, "sntpInfo")

			if v, ok := sntpMap["host"]; ok {
				if s, ok := v.(string); ok && s != "" {
					sntpInfo["host"] = s
				}
			}
			if v, ok := sntpMap["timezone"]; ok {
				if s, ok := v.(string); ok && s != "" {
					sntpInfo["timezone"] = s
				}
			}
			if v, ok := sntpMap["sntpEnable"]; ok {
				useFunction := ensureMap(general, "useFuction")
				useFunction["sntp"] = v
			}
		}

		// SNTP server IP from COMM_CFG
		if v, found := getMapVal(settings, "comm", "sntp"); found {
			sntpInfo := ensureMap(general, "sntpInfo")
			sntpInfo["server_ip"] = joinIPParts(v)
		}

		// FTP_DEF data
		if ftp, ok := settings["ftp"]; ok {
			ftpMap := ftp.(map[string]interface{})
			ftpInfo := ensureMap(general, "ftpInfo")

			if v, ok := ftpMap["host"]; ok {
				ftpInfo["host"] = joinIPParts(v)
			}
			if v, ok := ftpMap["port"]; ok {
				ftpInfo["port"] = toInt(v)
			}
			if v, ok := ftpMap["id"]; ok {
				if s, ok := v.(string); ok && s != "" {
					ftpInfo["id"] = s
				}
			}
			if v, ok := ftpMap["pass"]; ok {
				if s, ok := v.(string); ok && s != "" {
					ftpInfo["pass"] = s
				}
			}
			if v, ok := ftpMap["enable"]; ok {
				useFunction := ensureMap(general, "useFuction")
				useFunction["ftp"] = v
			}
			if v, ok := ftpMap["dir"]; ok {
				if s, ok := v.(string); ok && s != "" {
					ftpInfo["upload_main"] = s
				}
			}
		}

		// Main channel sampling info
		if channels := getSlice(jsonData, "channel"); channels != nil {
			for _, ch := range channels {
				chMap, ok := ch.(map[string]interface{})
				if !ok {
					continue
				}
				if chMap["channel"] != "Main" {
					continue
				}
				if _, ok := chMap["sampling"]; !ok {
					continue
				}
				samplingInfo := chMap["sampling"].(map[string]interface{})
				if commMap, ok := settings["comm"]; ok {
					comm := commMap.(map[string]interface{})
					if v, ok := comm["daq_srate"]; ok {
						samplingInfo["rate"] = v
					}
					if v, ok := comm["daq_length"]; ok {
						samplingInfo["duration"] = v
					}
					if v, ok := comm["daq_interval"]; ok {
						samplingInfo["period"] = v
					}
				}
				break
			}
		}

	} else {
		// Sub channel
		if ftpVal, ok := settings["ftp"]; ok {
			ftpMap := ftpVal.(map[string]interface{})
			if v, ok := ftpMap["dir"]; ok {
				if s, ok := v.(string); ok && s != "" {
					ftpInfo := ensureMap(general, "ftpInfo")
					ftpInfo["upload_sub"] = s
				}
			}
		}

		// Sub channel sampling info
		if channels := getSlice(jsonData, "channel"); channels != nil {
			for _, ch := range channels {
				chMap, ok := ch.(map[string]interface{})
				if !ok {
					continue
				}
				if chMap["channel"] != "Sub" {
					continue
				}
				if _, ok := chMap["sampling"]; !ok {
					continue
				}
				samplingInfo := chMap["sampling"].(map[string]interface{})
				if commMap, ok := settings["comm"]; ok {
					comm := commMap.(map[string]interface{})
					if v, ok := comm["daq_srate"]; ok {
						samplingInfo["rate"] = v
					}
					if v, ok := comm["daq_length"]; ok {
						samplingInfo["duration"] = v
					}
					if v, ok := comm["daq_interval"]; ok {
						samplingInfo["period"] = v
					}
				}
				break
			}
		}
	}

	return jsonData
}

// UpdateChannelFromSetting merges channel-specific settings (PT, CT, demand, alarm) into jsonData.
// Python equivalent: SettingMap.update_channel_from_setting
func UpdateChannelFromSetting(jsonData, settings map[string]interface{}, channelName string) map[string]interface{} {
	if jsonData == nil {
		jsonData = make(map[string]interface{})
	}

	channels := getSlice(jsonData, "channel")
	if channels == nil {
		return jsonData
	}

	found := false
	for _, ch := range channels {
		chMap, ok := ch.(map[string]interface{})
		if !ok {
			continue
		}
		if chMap["channel"] != channelName {
			continue
		}

		found = true
		log.Printf("'%s' channel updating...", channelName)

		// PT_DEF update
		if ptSettings, ok := settings["pt"]; ok {
			pt := ptSettings.(map[string]interface{})
			if ptInfoVal, ok := chMap["ptInfo"]; ok {
				ptInfo := ptInfoVal.(map[string]interface{})

				if v, ok := pt["wiring"]; ok {
					ptInfo["wiringmode"] = v
				}
				if v, ok := pt["freq"]; ok {
					ptInfo["linefrequency"] = v
				}
				if v, ok := pt["vnorm"]; ok {
					ptInfo["vnorminal"] = v
				}
				if v, ok := pt["PT1"]; ok {
					ptInfo["pt1"] = v
				}
				if v, ok := pt["PT2"]; ok {
					ptInfo["pt2"] = v
				}
			}
		}

		// CT_DEF update
		if ctSettings, ok := settings["ct"]; ok {
			ct := ctSettings.(map[string]interface{})
			if ctInfoVal, ok := chMap["ctInfo"]; ok {
				ctInfo := ctInfoVal.(map[string]interface{})

				if v, ok := ct["inorm"]; ok {
					ctInfo["inorminal"] = v
				}
				if v, ok := ct["CT1"]; ok {
					ctInfo["ct1"] = v
				}
				if v, ok := ct["CT2"]; ok {
					ctInfo["ct2"] = v
				}
				if v, ok := ct["ct_dir"]; ok {
					ctInfo["direction"] = v
				}
				if v, ok := ct["I_start"]; ok {
					ctInfo["startingcurrent"] = v
				}
				if v, ok := ct["zctScale"]; ok {
					ctInfo["zctscale"] = v
				}
				if v, ok := ct["zctType"]; ok {
					ctInfo["zcttpye"] = v
				}
			}
		}

		// Demand update (from etc settings)
		if etcSettings, ok := settings["etc"]; ok {
			etc := etcSettings.(map[string]interface{})
			if demandVal, ok := chMap["demand"]; ok {
				demandInfo := demandVal.(map[string]interface{})

				if v, ok := etc["interval"]; ok {
					demandInfo["demand_interval"] = toInt(v)
				}
				if v, ok := etc["P_target"]; ok {
					demandInfo["target"] = v
				}
			}
		}

		// ALARM_DEF update
		if alarmSettings, ok := settings["alarm"]; ok {
			alarm := alarmSettings.(map[string]interface{})
			if alarmVal, ok := chMap["alarm"]; ok {
				alarmMap := alarmVal.(map[string]interface{})

				if v, ok := alarm["delay"]; ok {
					alarmMap["CompareTimeDelay"] = v
				}

				if setVal, ok := alarm["set"]; ok {
					alarmSets := setVal.([]interface{})
					for alarmIdx, alarmSetVal := range alarmSets {
						alarmKey := fmt.Sprintf("%d", alarmIdx+1)
						if _, ok := alarmMap[alarmKey]; ok {
							alarmSet := alarmSetVal.(map[string]interface{})
							chanVal := 0
							condVal := 0
							dbandVal := 1
							levelVal := 0
							if v, ok := alarmSet["chan"]; ok {
								chanVal = toInt(v)
							}
							if v, ok := alarmSet["cond"]; ok {
								condVal = toInt(v)
							}
							if v, ok := alarmSet["dband"]; ok {
								dbandVal = toInt(v)
							}
							if v, ok := alarmSet["level"]; ok {
								// level can be float, keep as-is
								levelVal = toInt(v)
							}
							alarmMap[alarmKey] = []interface{}{chanVal, condVal, dbandVal, levelVal}
						}
					}
				}
			}
		}

		break
	}

	if !found {
		log.Printf("Warning: '%s' channel not found.", channelName)
	}

	return jsonData
}

// UpdateChannelEventFromSetting merges PQEvent and Transient array data into jsonData.
// Python equivalent: SettingMap.update_channel_event_from_setting
func UpdateChannelEventFromSetting(jsonData, settings map[string]interface{}, channelName string) map[string]interface{} {
	if jsonData == nil {
		jsonData = make(map[string]interface{})
	}

	// PQEvent mapping (index -> type name)
	pqeventMapping := map[int]string{
		0: "oc",
		1: "sag",
		2: "swell",
		3: "inter",
		4: "reserved",
	}

	// Transient mapping (index -> type name)
	transientMapping := map[int]string{
		0: "tv",
		1: "tc",
	}

	channels := getSlice(jsonData, "channel")
	if channels == nil {
		return jsonData
	}

	for _, ch := range channels {
		chMap, ok := ch.(map[string]interface{})
		if !ok {
			continue
		}
		if chMap["channel"] != channelName {
			continue
		}

		eventInfo := ensureMap(chMap, "eventInfo")

		// PQEvent update
		if pqevtVal, ok := settings["pqevt"]; ok {
			pqevtSlice := pqevtVal.([]interface{})
			for idx, pqeventVal := range pqevtSlice {
				eventType, exists := pqeventMapping[idx]
				if !exists {
					continue
				}
				if eventType == "reserved" {
					continue
				}

				pqevent := pqeventVal.(map[string]interface{})

				actionVal := 0
				if v, ok := pqevent["action"]; ok {
					actionVal = toInt(v)
				}
				eventInfo[fmt.Sprintf("%s_action", eventType)] = actionVal

				levelVal := 0
				if v, ok := pqevent["level"]; ok {
					levelVal = toInt(v)
				}
				eventInfo[fmt.Sprintf("%s_level", eventType)] = levelVal

				holdoffVal := 0
				if v, ok := pqevent["holdOffCyc"]; ok {
					holdoffVal = toInt(v)
				}
				eventInfo[fmt.Sprintf("%s_holdofftime", eventType)] = holdoffVal

				if v, ok := pqevent["nCyc"]; ok {
					eventInfo[fmt.Sprintf("%s_ncyc", eventType)] = toInt(v)
				}
			}
		}

		// Transient update
		if transientVal, ok := settings["transient"]; ok {
			transientSlice := transientVal.([]interface{})
			for idx, transVal := range transientSlice {
				transType, exists := transientMapping[idx]
				if !exists {
					continue
				}

				trans := transVal.(map[string]interface{})

				actionVal := 0
				if v, ok := trans["action"]; ok {
					actionVal = toInt(v)
				}
				eventInfo[fmt.Sprintf("%s_action", transType)] = actionVal

				levelVal := 0
				if v, ok := trans["level"]; ok {
					levelVal = toInt(v)
				}
				eventInfo[fmt.Sprintf("%s_level", transType)] = levelVal

				holdoffVal := 0
				if v, ok := trans["holdOff"]; ok {
					holdoffVal = toInt(v)
				}
				eventInfo[fmt.Sprintf("%s_holdofftime", transType)] = holdoffVal

				fastchangeVal := 0
				if v, ok := trans["fastChange"]; ok {
					fastchangeVal = toInt(v)
				}
				eventInfo[fmt.Sprintf("%s_fastchange", transType)] = fastchangeVal
			}
		}

		break
	}

	return jsonData
}
