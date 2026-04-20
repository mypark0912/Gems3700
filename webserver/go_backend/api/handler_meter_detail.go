package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"serverGO/binary"
	"serverGO/metermap"
)

// normalizeChannel maps path param to "Main" / "Sub".
func normalizeChannel(channel string) string {
	if channel == "Main" || channel == "main" {
		return "Main"
	}
	return "Sub"
}

// GET /api/getInterval/:mode/:channel
func (h *Handler) getInterval(c *gin.Context) {
	mode := c.Param("mode")
	channel := c.Param("channel")
	ctx := context.Background()

	chName := "main"
	if channel != "Main" && channel != "main" {
		chName = "sub"
	}

	hashField := "SamplingPeriod"
	if mode == "demand" {
		hashField = "DemandInterval"
	}

	raw, err := h.deps.Redis.Client0.HGet(ctx, "Equipment", hashField).Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}

	var intervalDict map[string]interface{}
	if json.Unmarshal([]byte(raw), &intervalDict) != nil {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": intervalDict[chName]})
}

// GET /api/getOnesfromRedis/:channel/:unbal
// Mirrors Python get_OneSecond: MAXMIN_<channel> "1sec" + "15min" + meter flat.
func (h *Handler) getOnesfromRedis(c *gin.Context) {
	channel := normalizeChannel(c.Param("channel"))
	unbal := c.Param("unbal")
	ctx := context.Background()
	client := h.deps.Redis.Client0

	if unbal == "null" || unbal == "" || unbal == "-1" {
		unbal = "0"
	}

	keyname := "meter"
	exists, _ := client.Exists(ctx, keyname).Result()
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": "No exist key"})
		return
	}

	flatData, _ := client.HGetAll(ctx, keyname).Result()
	meters := make(map[string]float64, len(flatData))
	for k, v := range flatData {
		meters[k] = tryFloat(v)
	}

	oneSec, _ := binary.FetchMaxMin1Sec(ctx, client, channel)
	fifteenMin, _ := binary.FetchMaxMin15Min(ctx, client, channel)
	maxmin1s := binary.Flat1Sec(oneSec)
	maxmin15m := binary.Flat15Min(fifteenMin)

	pVoltageData := metermap.GetDataDict(meters, maxmin1s, metermap.PVoltageKeys, "V")
	freqData := metermap.GetDataDict(meters, maxmin1s, metermap.FreqKeys, "Hz")
	lVoltageData := metermap.GetDataDict(meters, maxmin1s, metermap.LVoltageKeys, "V")
	currentData := metermap.GetDataDict(meters, maxmin1s, metermap.CurrentKeys, "A")
	pfData := metermap.GetDataDict(meters, maxmin1s, metermap.PFKeys, "")

	var unbalData []map[string]interface{}
	if unbal == "1" {
		unbalData = metermap.GetDataDict(meters, maxmin1s, metermap.UnbalKeys, "%")
	} else {
		unbalData = metermap.GetDataDict(meters, maxmin1s, metermap.UnbalKeys2, "%")
	}

	aPowerData := metermap.GetDataDict(meters, maxmin15m, metermap.APowerKeys, "kW")
	rPowerData := metermap.GetDataDict(meters, maxmin15m, metermap.RPowerKeys, "kVar")
	apPowerData := metermap.GetDataDict(meters, maxmin15m, metermap.APPowerKeys, "kVA")
	thduData := metermap.GetDataDict(meters, maxmin15m, metermap.THDUKeys, "%")
	thdiData := metermap.GetDataDict(meters, maxmin15m, metermap.THDIKeys, "%")
	tddiData := metermap.GetDataDict(meters, maxmin15m, metermap.TDDIKeys, "%")

	result := gin.H{
		"success": true,
		"retData": gin.H{
			"unbalData": []gin.H{
				{"subTitle": "Voltage", "data": []gin.H{
					{"id": 0, "subTitle": "-", "value": meters["Ubal1"], "max": "-", "min": "-", "unit": "%"},
				}},
				{"subTitle": "Current", "data": []gin.H{
					{"id": 0, "subTitle": "-", "value": meters["Ibal1"], "max": "-", "min": "-", "unit": "%"},
				}},
			},
			"meterData": []gin.H{
				{"subTitle": "Phase Voltage", "data": pVoltageData},
				{"subTitle": "Line Voltage", "data": lVoltageData},
				{"subTitle": "Current", "data": currentData},
				{"subTitle": "Frequency", "data": freqData},
				{"subTitle": "PF", "data": pfData},
				{"subTitle": "Unbalance", "data": unbalData},
			},
			"powerData": []gin.H{
				{"subTitle": "Active Power", "data": aPowerData},
				{"subTitle": "Reactive Power", "data": rPowerData},
				{"subTitle": "Apparent Power", "data": apPowerData},
			},
			"thdData": []gin.H{
				{"subTitle": "THD-U", "data": thduData},
				{"subTitle": "THD-I", "data": thdiData},
				{"subTitle": "TDD-I", "data": tddiData},
			},
		},
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "retData": result})
}

// GET /api/getonemfromRedis/:channel
// Mirrors Python get_onemfromRedis: phasor magnitude/angle + U/I max.
func (h *Handler) getonemfromRedis(c *gin.Context) {
	channel := normalizeChannel(c.Param("channel"))
	ctx := context.Background()
	client := h.deps.Redis.Client0

	keyname := "meter"
	exists, _ := client.Exists(ctx, keyname).Result()
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": "No exist key"})
		return
	}

	flatData, _ := client.HGetAll(ctx, keyname).Result()
	oneSec, _ := binary.FetchMaxMin1Sec(ctx, client, channel)

	maxVals := make([]interface{}, 6)
	fallbackKeys := []string{"U1", "U2", "U3", "I1", "I2", "I3"}
	for i := 0; i < 6; i++ {
		var v float64
		if oneSec != nil {
			if i < 3 {
				v = oneSec.U[i].Max
			} else {
				v = oneSec.I[i-3].Max
			}
		}
		if v == 0 {
			v = tryFloat(flatData[fallbackKeys[i]])
		}
		maxVals[i] = v
	}

	result := gin.H{
		"success": true,
		"retData": gin.H{
			"angleData": gin.H{
				"magnitude": []float64{
					tryFloat(flatData["U1"]), tryFloat(flatData["U2"]), tryFloat(flatData["U3"]),
					tryFloat(flatData["I1"]), tryFloat(flatData["I2"]), tryFloat(flatData["I3"]),
				},
				"degree": []float64{
					tryFloat(flatData["Uangle1"]), tryFloat(flatData["Uangle2"]), tryFloat(flatData["Uangle3"]),
					tryFloat(flatData["Iangle1"]), tryFloat(flatData["Iangle2"]), tryFloat(flatData["Iangle3"]),
				},
				"texts": []string{"V1", "V2", "V3", "I1", "I2", "I3"},
				"max":   maxVals,
			},
		},
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "retData": result})
}

// GET /api/getFifthMfromRedis/:channel
// Mirrors Python get_FifthMfromRedis: demand P/Q/S/I from Demand hash field <channel>.
func (h *Handler) getFifthMfromRedis(c *gin.Context) {
	channel := normalizeChannel(c.Param("channel"))
	ctx := context.Background()
	client := h.deps.Redis.Client0

	demand, err := binary.FetchDemand(ctx, client, channel)
	if err != nil || demand == nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": "No demand data"})
		return
	}

	formatted := binary.FormatDemandData(demand)

	result := gin.H{
		"success": true,
		"retData": gin.H{
			"demandDataP": []gin.H{
				{"subTitle": "Active", "data": formatted["power_demand"]},
				{"subTitle": "Reactive", "data": formatted["reactive_demand"]},
				{"subTitle": "Apparent", "data": formatted["apparent_demand"]},
			},
			"demandDataI": []gin.H{
				{"subTitle": "Current", "data": formatted["current_demand"]},
			},
		},
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "retData": result})
}

// GET /api/getOnehfromRedis/:channel
// Mirrors Python get_onehRedis: energy totals from energy_<channel>.
func (h *Handler) getOnehfromRedis(c *gin.Context) {
	channel := c.Param("channel")
	ctx := context.Background()
	client := h.deps.Redis.Client0

	suffix := "main"
	if channel != "Main" && channel != "main" {
		suffix = "sub"
	}
	keyname := "energy_" + suffix

	exists, _ := client.Exists(ctx, keyname).Result()
	if exists == 0 {
		emptyResult := gin.H{
			"success": false,
			"retData": gin.H{
				"energyData": []gin.H{
					{"subTitle": "Import", "data": []gin.H{
						{"id": 0, "subTitle": "-", "value": 0, "max": "-", "min": "-", "unit": "kWh"},
					}},
					{"subTitle": "Export", "data": []gin.H{
						{"id": 0, "subTitle": "-", "value": 0, "max": "-", "min": "-", "unit": "kWh"},
					}},
				},
			},
		}
		c.JSON(http.StatusOK, gin.H{"success": false, "error": "No exist key", "retData": emptyResult})
		return
	}

	fields := []string{"total_kwh_import", "total_kwh_export", "thismonth_kwh_import", "thismonth_kwh_export"}
	values, _ := client.HMGet(ctx, keyname, fields...).Result()

	flat := make(map[string]float64, len(fields))
	for i, f := range fields {
		if s, ok := values[i].(string); ok {
			flat[f] = tryFloat(s)
		}
	}

	result := gin.H{
		"success": true,
		"retData": gin.H{
			"energyData": []gin.H{
				{"subTitle": "Import", "data": []gin.H{
					{"id": 0, "subTitle": "-", "value": flat["total_kwh_import"], "max": "-", "min": "-", "unit": "kWh"},
					{"id": 1, "subTitle": "-", "value": flat["thismonth_kwh_import"], "max": "-", "min": "-", "unit": "kWh"},
				}},
				{"subTitle": "Export", "data": []gin.H{
					{"id": 0, "subTitle": "-", "value": flat["total_kwh_export"], "max": "-", "min": "-", "unit": "kWh"},
					{"id": 1, "subTitle": "-", "value": flat["thismonth_kwh_export"], "max": "-", "min": "-", "unit": "kWh"},
				}},
			},
		},
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "retData": result})
}

// GET /api/getEnergyRedis/:channel
func (h *Handler) getEnergyRedis(c *gin.Context) {
	channel := c.Param("channel")
	ctx := context.Background()
	client := h.deps.Redis.Client0

	emptyData := gin.H{
		"today": 0, "daily": 0, "weekly": 0, "monthly": 0, "yearly": 0,
		"daily_comparison": 0, "weekly_comparison": 0, "monthly_comparison": 0, "yearly_comparison": 0,
	}

	raw, err := client.HGet(ctx, "energy_summary", channel).Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": emptyData})
		return
	}

	var conData map[string]interface{}
	if json.Unmarshal([]byte(raw), &conData) != nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": emptyData})
		return
	}

	getFloat := func(m map[string]interface{}, key string) float64 {
		if v, ok := m[key]; ok {
			if val, ok := v.(float64); ok {
				return val
			}
		}
		return 0
	}

	getNestedFloat := func(m map[string]interface{}, outerKey, innerKey string) float64 {
		if sub, ok := m[outerKey].(map[string]interface{}); ok {
			return getFloat(sub, innerKey)
		}
		return 0
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{
		"today":              getNestedFloat(conData, "consumption", "kwh_import_consumption"),
		"daily":              getFloat(conData, "daily_kwh_import"),
		"weekly":             getFloat(conData, "weekly_kwh_import"),
		"monthly":            getFloat(conData, "monthly_kwh_import"),
		"yearly":             getFloat(conData, "yearly_kwh_import"),
		"daily_comparison":   getNestedFloat(conData, "daily_comparison", "change_percent"),
		"weekly_comparison":  getNestedFloat(conData, "weekly_comparison", "change_percent"),
		"monthly_comparison": getNestedFloat(conData, "monthly_comparison", "change_percent"),
		"yearly_comparison":  getNestedFloat(conData, "yearly_comparison", "change_percent"),
	}})
}
