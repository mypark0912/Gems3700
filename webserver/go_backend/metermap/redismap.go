package metermap

import "math"

type KeyDef struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
	Key   string `json:"key"`
}

// RedisMapDetail2 — main meter field mappings (used by getMeterRedisNew)
var FreqKeys = []KeyDef{{0, "-", "Freq"}}

var PVoltageKeys = []KeyDef{
	{0, "L1", "U1"}, {1, "L2", "U2"}, {2, "L3", "U3"},
}
var PVoltAngleKeys = []KeyDef{
	{0, "L1", "Uangle1"}, {1, "L2", "Uangle2"}, {2, "L3", "Uangle3"},
}
var LVoltageKeys = []KeyDef{
	{0, "L1-L2", "Upp1"}, {1, "L2-L3", "Upp2"}, {2, "L3-L1", "Upp3"},
}
var CurrentKeys = []KeyDef{
	{0, "L1", "I1"}, {1, "L2", "I2"}, {2, "L3", "I3"}, {3, "Ground", "Ig"},
}
var CurrentAngleKeys = []KeyDef{
	{0, "L1", "Iangle1"}, {1, "L2", "Iangle2"}, {2, "L3", "Iangle3"},
}
var PFKeys = []KeyDef{{0, "Total", "PF4"}}
var APowerKeys = []KeyDef{{0, "Total", "P4"}}
var RPowerKeys = []KeyDef{{0, "Total", "Q4"}}
var APPowerKeys = []KeyDef{{0, "Total", "S4"}}

var THDUKeys = []KeyDef{
	{0, "L1", "THD_U1"}, {1, "L2", "THD_U2"}, {2, "L3", "THD_U3"},
}
var THDIKeys = []KeyDef{
	{0, "L1", "THD_I1"}, {1, "L2", "THD_I2"}, {2, "L3", "THD_I3"},
}
var TDDIKeys = []KeyDef{
	{0, "L1", "TDD_I1"}, {1, "L2", "TDD_I2"}, {2, "L3", "TDD_I3"},
}
var KWhKeys = []KeyDef{
	{0, "Import", "total_kwh_import"}, {1, "Export", "total_kwh_export"},
	{2, "Import", "thismonth_kwh_import"}, {3, "Export", "thismonth_kwh_export"},
}
var UnbalKeys = []KeyDef{
	{0, "Voltage", "Ubal1"}, {1, "Current", "Ibal1"},
}
var UnbalKeys2 = []KeyDef{
	{0, "Voltage", "Ubal_nema"}, {1, "Current", "Ibal_nema"},
}

// RedisMapCalibrate — keys used by /calibrateNow (per-phase, no totals)
var CalPVoltageKeys = []KeyDef{
	{0, "L1", "U1"}, {1, "L2", "U2"}, {2, "L3", "U3"},
}
var CalPAngleKeys = []KeyDef{
	{0, "L1", "Pangle1"}, {1, "L2", "Pangle2"}, {2, "L3", "Pangle3"},
}
var CalCurrentKeys = []KeyDef{
	{0, "L1", "I1"}, {1, "L2", "I2"}, {2, "L3", "I3"}, {3, "Ground", "Ig"},
}
var CalAPowerKeys = []KeyDef{
	{0, "L1", "P1"}, {1, "L2", "P2"}, {2, "L3", "P3"},
}
var CalRPowerKeys = []KeyDef{
	{0, "L1", "Q1"}, {1, "L2", "Q2"}, {2, "L3", "Q3"},
}
var CalAPPowerKeys = []KeyDef{
	{0, "L1", "S1"}, {1, "L2", "S2"}, {2, "L3", "S3"},
}

// RedisMapped — dashboard field sets
var DashboardMeter = map[string][]string{
	"meter":  {"U1", "U2", "U3", "U4", "Upp1", "Upp2", "Upp3", "Upp4", "I1", "I2", "I3", "I4", "Itot", "Ubal1", "Ibal1", "Freq", "Ig", "Temp", "Ubal_nema", "Ibal_nema"},
	"power":  {"P1", "P2", "P3", "P4", "PF4", "PF1", "PF2", "PF3", "S4", "Q4"},
	"energy": {"total_kwh_import", "thismonth_kwh_import"},
	"thd":    {"THD_U1", "THD_U2", "THD_U3", "THD_I1", "THD_I2", "THD_I3", "TDD_I1", "TDD_I2", "TDD_I3"},
}

var DashboardTrans = map[string][]string{
	"meter": {"Freq", "Ig", "Temp"},
}

var EnergyReport = []string{
	"import kwh", "export kwh", "import kwh this month", "export kwh this month",
	"import kwh last month", "export kwh last month",
}

// ParameterOptions — 48 electrical measurement parameters for trend selection
var ParameterOptions = []string{
	"None", "Temperature", "Frequency",
	"Phase Voltage L1", "Phase Voltage L2", "Phase Voltage L3", "Phase Voltage Average",
	"Line Voltage L12", "Line Voltage L23", "Line Voltage L31", "Line Voltage Average",
	"Voltage Unbalance(Uo)", "Voltage Unbalance(Uu)",
	"Phase Current L1", "Phase Current L2", "Phase Current L3",
	"Phase Current Average", "Phase Current Total", "Phase Current Neutral",
	"Active Power L1", "Active Power L2", "Active Power L3", "Active Power Total",
	"Reactive Power L1", "Reactive Power L2", "Reactive Power L3", "Reactive Power Total",
	"Distortion Power L1", "Distortion Power L2", "Distortion Power L3", "Distortion Power Total",
	"Apparent Power L1", "Apparent Power L2", "Apparent Power L3", "Apparent Power Total",
	"Power Factor L1", "Power Factor L2", "Power Factor L3", "Power Factor Total",
	"THD Voltage L1", "THD Voltage L2", "THD Voltage L3",
	"THD Voltage L12", "THD Voltage L23", "THD Voltage L31",
	"THD Current L1", "THD Current L2", "THD Current L3",
}

// GetDataDict formats meter values with max/min for frontend
func GetDataDict(meters map[string]float64, maxmin map[string]interface{}, keys []KeyDef, unit string) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(keys))
	for _, item := range keys {
		key := item.Key
		entry := map[string]interface{}{
			"id":       item.ID,
			"subTitle": item.Label,
			"value":    safeRound(meters[key], 2),
			"max":      safeRound(toFloat(maxmin[key+"_max"]), 2),
			"maxTime":  maxmin[key+"_maxTime"],
			"min":      safeRound(toFloat(maxmin[key+"_min"]), 2),
			"minTime":  maxmin[key+"_minTime"],
			"unit":     unit,
		}
		result = append(result, entry)
	}
	return result
}

func safeRound(v float64, decimals int) float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	pow := math.Pow(10, float64(decimals))
	return math.Round(v*pow) / pow
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
	}
	return 0
}
