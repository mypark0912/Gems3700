package processors

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/redis/go-redis/v9"

	"sv500_core/data"
	"sv500_core/handlers"
)

const (
	diagnosisBaseURL    = "http://127.0.0.1:5000/api"
	diagnosisHTTPTimeout = 20 * time.Second
)

var apiTimeout = &http.Client{
	Timeout: diagnosisHTTPTimeout,
}

// ---------------------------------------------------------------------------
// Helper: HTTP GET with timeout
// ---------------------------------------------------------------------------

func fetchAPI(endpoint string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s", diagnosisBaseURL, endpoint)
	resp, err := apiTimeout.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP GET %s: status %d", endpoint, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON parse error: %w", err)
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// Helper: Empty shell data check
// ---------------------------------------------------------------------------

func isEmptyShellData(datas map[string]interface{}) bool {
	treeList, ok := datas["TreeList"].([]interface{})
	if !ok || len(treeList) == 0 {
		return true
	}

	allNaN := true
	allStatusZero := true
	for _, item := range treeList {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if fmt.Sprintf("%v", m["Value"]) != "NaN" {
			allNaN = false
		}
		if data.GetIntFromInterface(m["Status"]) != 0 {
			allStatusZero = false
		}
	}

	return allNaN && allStatusZero
}

// ---------------------------------------------------------------------------
// ProcessAllStatusData: API → Redis (BarGraph 데이터 저장)
// Python: process_allstatus_data()
// ---------------------------------------------------------------------------

func ProcessAllStatusData(chName, asset, assetType string, confStatus map[string]interface{}, redisInst *redis.Client) bool {
	ctx := context.Background()

	datas, err := fetchAPI(fmt.Sprintf("getBarGraphs?name=%s", asset))
	if err != nil {
		log.Printf("Diagnosis API failed (%s): %v", chName, err)
		return false
	}
	if datas == nil {
		return false
	}

	// Process each state type and save to Redis
	stateList := []string{"Diagnostic", "PQ", "Events", "Faults"}
	redisfiedList := []string{"Diagnosis", "PQ", "Event", "Fault"}

	for i, stateName := range stateList {
		stData, ok := datas[stateName].([]interface{})
		if !ok {
			continue
		}

		var updated map[string]interface{}
		switch stateName {
		case "Diagnostic":
			diag := data.NewDiagnosis(assetType)
			updated = diag.UpdateFromBargraphData(stData)
		case "PQ":
			pq := &data.PowerQuality{Data: &data.PQ{}}
			data.InitPQDefaults(pq.Data)
			updated = pq.UpdateFromBargraphData(stData)
		case "Events":
			ev := &data.EventWrapper{Data: &data.EV{}}
			data.InitEVDefaults(ev.Data)
			updated = ev.UpdateFromBargraphData(stData)
		case "Faults":
			ft := &data.FaultWrapper{Data: &data.FT{}}
			data.InitFTDefaults(ft.Data)
			updated = ft.UpdateFromBargraphData(stData)
		}

		if updated != nil {
			jsonBytes, _ := json.Marshal(updated)
			redisInst.HSet(ctx, fmt.Sprintf("SmartSystem:%s", chName), redisfiedList[i], string(jsonBytes))
		}
	}

	// Process DO alarm status
	doAlarm := procDashAlarm(datas, chName, confStatus)
	redisInst.HSet(ctx, "DOStatus", chName, doAlarm)

	log.Printf("Diagnosis BarGraph to Redis saved: %s", chName)
	return true
}

func procDashAlarm(datas map[string]interface{}, chName string, channelConfig map[string]interface{}) int {
	if channelConfig == nil {
		return 0
	}

	stateList := []string{"Diagnostic", "PQ"}
	matcher := data.NewAlarmStatusMatcher()

	emptyConfig := 0
	retDict := make(map[string]map[string]interface{})

	for _, category := range stateList {
		var configKey string
		if category == "Diagnostic" {
			configKey = "diagnosis"
		} else {
			configKey = "pq"
		}

		configListRaw, ok := channelConfig[configKey].([]interface{})
		if !ok || len(configListRaw) == 0 {
			emptyConfig++
			continue
		}

		// Convert []interface{} to []map[string]interface{}
		configList := make([]map[string]interface{}, 0, len(configListRaw))
		for _, item := range configListRaw {
			if m, ok := item.(map[string]interface{}); ok {
				configList = append(configList, m)
			}
		}

		barGraphRaw, _ := datas[category].([]interface{})
		barGraph := make([]map[string]interface{}, 0, len(barGraphRaw))
		for _, item := range barGraphRaw {
			if m, ok := item.(map[string]interface{}); ok {
				barGraph = append(barGraph, m)
			}
		}

		result := matcher.Diagnose(configList, barGraph)

		finalStatus := result.FinalStatus
		matchedParams := result.MatchedParameters

		if len(matchedParams) == 0 {
			if finalStatus == 0 {
				retDict[category] = map[string]interface{}{"status": 0, "item": "All"}
			} else {
				retDict[category] = map[string]interface{}{"status": 1, "item": "All"}
			}
			continue
		}

		maxStatus := 0
		var topItems []data.MatchedParameter
		for _, p := range matchedParams {
			if p.ActualStatus > maxStatus {
				maxStatus = p.ActualStatus
				topItems = nil
			}
			if p.ActualStatus == maxStatus {
				topItems = append(topItems, p)
			}
		}

		itemLabel := ""
		if len(topItems) == 1 {
			itemLabel = topItems[0].Name
		} else if len(topItems) == 2 {
			itemLabel = fmt.Sprintf("%s, %s", topItems[0].Name, topItems[1].Name)
		} else if len(topItems) > 2 {
			itemLabel = fmt.Sprintf("%s ... +%d", topItems[0].Name, len(topItems)-1)
		}

		retDict[category] = map[string]interface{}{"status": finalStatus, "item": itemLabel}
	}

	if emptyConfig == len(stateList) {
		return 0
	}

	// If PQ or Diagnostic final status >= 2, alarm
	for _, category := range stateList {
		if cat, ok := retDict[category]; ok {
			if data.GetIntFromInterface(cat["status"]) >= 2 {
				return 1
			}
		}
	}
	return 0
}

// ---------------------------------------------------------------------------
// ProcessDiagnosisData: API → InfluxDB (진단 상세 데이터)
// Python: process_diagnosis_data()
// ---------------------------------------------------------------------------

func ProcessDiagnosisData(chName, asset, assetType string) bool {
	datas, err := fetchAPI(fmt.Sprintf("getDiagnostic?name=%s", asset))
	if err != nil {
		log.Printf("Diagnosis Influx failed (%s): %v", chName, err)
		return false
	}
	if datas == nil || isEmptyShellData(datas) {
		return false
	}

	bargraph, _ := datas["BarGraph"].([]interface{})
	lastRecordTime, _ := datas["LastRecordDateTime"].(string)

	diag := data.NewDiagnosis(assetType)
	updated := diag.UpdateFromBargraphDataForInflux(bargraph)

	filteredTree := procAndFilterTree(datas)

	saveAPIToInfluxDB(chName, asset, assetType, lastRecordTime, updated, filteredTree, "diagnosis")
	log.Printf("Diagnosis API to InfluxDB saved: %s", chName)
	return true
}

// ProcessPQData: PQ API → InfluxDB
func ProcessPQData(chName, asset, assetType string) bool {
	datas, err := fetchAPI(fmt.Sprintf("getPQ?name=%s", asset))
	if err != nil {
		log.Printf("PQ Influx failed (%s): %v", chName, err)
		return false
	}
	if datas == nil || isEmptyShellData(datas) {
		return false
	}

	bargraph, _ := datas["BarGraph"].([]interface{})
	lastRecordTime, _ := datas["LastRecordDateTime"].(string)

	pq := &data.PowerQuality{Data: &data.PQ{}}
	data.InitPQDefaults(pq.Data)
	updated := pq.UpdateFromBargraphDataForInflux(bargraph)

	filteredTree := procAndFilterTree(datas)

	saveAPIToInfluxDB(chName, asset, assetType, lastRecordTime, updated, filteredTree, "powerquality")
	log.Printf("PQ API to InfluxDB saved: %s", chName)
	return true
}

// ProcessAssetData: THD/TDD Redis 저장 (VFD용)
func ProcessAssetData(chName, asset string, redisInst *redis.Client) bool {
	ctx := context.Background()
	datas, err := fetchAPI(fmt.Sprintf("getRealTimeData?name=%s", asset))
	if err != nil {
		log.Printf("Asset Data failed (%s): %v", chName, err)
		return false
	}
	if datas == nil {
		return false
	}

	jsonBytes, _ := json.Marshal(datas)
	redisInst.HSet(ctx, fmt.Sprintf("SmartAPI:%s", chName), "AssetData", string(jsonBytes))

	// THD/TDD to meter hash
	thdMapping := map[string]string{"thdv": "THD_V", "thdi": "THD_I", "tddi": "TDD_I"}
	meterKey := "meter_main"
	if chName != "Main" {
		meterKey = "meter_sub"
	}

	if dataList, ok := datas["Data"].([]interface{}); ok {
		for _, item := range dataList {
			m, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			name := fmt.Sprintf("%v", m["Name"])
			if redisField, ok := thdMapping[name]; ok {
				redisInst.HSet(ctx, meterKey, redisField, fmt.Sprintf("%v", m["Value"]))
			}
		}
	}

	log.Printf("Asset Data saved: %s", chName)
	return true
}

// ---------------------------------------------------------------------------
// ProcessAllDiagnosisData: 메인 진단 루프 (3분마다)
// Python: process_all_diagnosis_data()
// ---------------------------------------------------------------------------

func ProcessAllDiagnosisData(channel *data.Channel, redisInst *redis.Client) {
	chName := channel.Name
	if strings.ToLower(chName) == "main" {
		chName = "Main"
	} else {
		chName = "Sub"
	}

	log.Printf("Starting %s channel diagnosis processing", chName)

	ProcessAllStatusData(chName, channel.AssetName, channel.AssetType, channel.ConfStatus, redisInst)

	if channel.AssetDrive {
		ProcessAssetData(chName, channel.AssetName, redisInst)
	}

	log.Printf("Completed %s diagnosis processing", chName)
	runtime.GC()
}

// Process1HDiagnosisData: 6시간 주기 InfluxDB 저장
// Python: process_1h_diagnosis_data()
func Process1HDiagnosisData(channelName, asset, assetType string) {
	chName := channelName
	if strings.ToLower(chName) == "main" {
		chName = "Main"
	} else {
		chName = "Sub"
	}

	log.Printf("Starting %s 6h diagnosis processing", chName)

	results := map[string]bool{
		"diagnosis": ProcessDiagnosisData(chName, asset, assetType),
		"pq":        ProcessPQData(chName, asset, assetType),
	}

	successCount := 0
	for _, v := range results {
		if v {
			successCount++
		}
	}
	log.Printf("Completed %s - %d diagnosis API to InfluxDB", chName, successCount)
	runtime.GC()
}

// ---------------------------------------------------------------------------
// Tree processing helpers
// ---------------------------------------------------------------------------

func procAndFilterTree(datas map[string]interface{}) []map[string]interface{} {
	// Build tree from BarGraph + TreeList
	bargraph, _ := datas["BarGraph"].([]interface{})
	treeList, _ := datas["TreeList"].([]interface{})

	if len(bargraph) == 0 || len(treeList) == 0 {
		return nil
	}

	// Build superNodes set from BarGraph NodeTypes
	superNodes := make(map[int]bool)
	for _, item := range bargraph {
		m, ok := item.(map[string]interface{})
		if ok {
			superNodes[data.GetIntFromInterface(m["NodeType"])] = true
		}
	}

	// Build children index
	childrenByParent := make(map[int][]map[string]interface{})
	for _, item := range treeList {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		parentID := data.GetIntFromInterface(m["ParentID"])
		childrenByParent[parentID] = append(childrenByParent[parentID], m)
	}

	// Build super list
	var superList []map[string]interface{}
	for _, item := range treeList {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		nodeType := data.GetIntFromInterface(m["NodeType"])
		if !superNodes[nodeType] {
			continue
		}

		node := copyNodeFields(m)
		node["isParent"] = true

		itemID := data.GetIntFromInterface(m["ID"])
		if children, ok := childrenByParent[itemID]; ok {
			var childList []map[string]interface{}
			for _, child := range children {
				childNode := copyNodeFields(child)
				childNode["NodeType"] = data.GetIntFromInterface(child["NodeType"])
				childNode["isSub"] = data.GetIntFromInterface(child["NodeType"]) == 10

				// NodeType 10 children (NodeType 11)
				childID := data.GetIntFromInterface(child["ID"])
				if data.GetIntFromInterface(child["NodeType"]) == 10 {
					if subChildren, ok := childrenByParent[childID]; ok {
						var subChildList []map[string]interface{}
						for _, sc := range subChildren {
							if data.GetIntFromInterface(sc["NodeType"]) == 11 {
								subChildList = append(subChildList, copyNodeFields(sc))
							}
						}
						if len(subChildList) > 0 {
							childNode["children"] = subChildList
						}
					}
				}
				childList = append(childList, childNode)
			}
			node["children"] = childList
		}
		superList = append(superList, node)
	}

	// Filter abnormal (Status >= 2)
	return filterAbnormal(superList)
}

func copyNodeFields(m map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"ID":           m["ID"],
		"Name":         m["Name"],
		"Title":        m["Title"],
		"Titles":       m["Titles"],
		"AssemblyID":   m["AssemblyID"],
		"Description":  m["Description"],
		"Descriptions": m["Descriptions"],
		"Path":         m["Path"],
		"Status":       m["Status"],
		"Value":        m["Value"],
	}
}

func filterAbnormal(dataTree []map[string]interface{}) []map[string]interface{} {
	var filtered []map[string]interface{}
	for _, parent := range dataTree {
		if data.GetIntFromInterface(parent["Status"]) < 2 {
			continue
		}
		filteredParent := copyNodeFields(parent)
		filteredParent["isParent"] = true

		if children, ok := parent["children"].([]map[string]interface{}); ok {
			var filteredChildren []map[string]interface{}
			for _, child := range children {
				if data.GetIntFromInterface(child["Status"]) < 2 {
					continue
				}
				fc := copyNodeFields(child)
				fc["NodeType"] = child["NodeType"]
				fc["isSub"] = child["isSub"]

				if subChildren, ok := child["children"].([]map[string]interface{}); ok {
					var filteredSub []map[string]interface{}
					for _, sc := range subChildren {
						if data.GetIntFromInterface(sc["Status"]) >= 2 {
							filteredSub = append(filteredSub, copyNodeFields(sc))
						}
					}
					if len(filteredSub) > 0 {
						fc["children"] = filteredSub
					}
				}
				filteredChildren = append(filteredChildren, fc)
			}
			if len(filteredChildren) > 0 {
				filteredParent["children"] = filteredChildren
			}
		}
		filtered = append(filtered, filteredParent)
	}
	return filtered
}

// ---------------------------------------------------------------------------
// saveAPIToInfluxDB: batch InfluxDB write
// Python: save_api_to_influxdb()
// ---------------------------------------------------------------------------

func saveAPIToInfluxDB(chName, asset, assetType, lastRecordTime string, updated map[string]interface{}, filteredTree []map[string]interface{}, meas string) {
	pool, err := handlers.GetInfluxPool()
	if err != nil {
		log.Printf("InfluxDB pool error: %v", err)
		return
	}

	// Parse timestamp
	timestamp := time.Now().UTC()
	if lastRecordTime != "" && !strings.HasPrefix(lastRecordTime, "0000") {
		if t, err := time.Parse(time.RFC3339, lastRecordTime); err == nil {
			timestamp = t.UTC()
		} else if t, err := time.Parse("2006-01-02T15:04:05-07:00", lastRecordTime); err == nil {
			timestamp = t.UTC()
		} else if t, err := time.Parse("2006-01-02T15:04:05", lastRecordTime); err == nil {
			timestamp = t.UTC()
		}
	}

	var allPoints []*write.Point

	// 1. BarGraph points (main data)
	equipmentType := 0
	if t, ok := updated["type"]; ok {
		equipmentType = data.GetIntFromInterface(t)
	}

	for fieldName, fieldDataRaw := range updated {
		if fieldName == "type" {
			continue
		}

		fieldData, ok := fieldDataRaw.(map[string]interface{})
		if !ok {
			continue
		}

		point := influxdb2.NewPointWithMeasurement(meas).
			AddTag("channel", chName).
			AddTag("asset_name", asset).
			AddTag("data_type", "main").
			AddTag("item_name", fieldName).
			AddField("item_id", data.GetIntFromInterface(fieldData["id"])).
			AddField("asset_type", assetType).
			AddField("equipment_type", equipmentType).
			AddField("status", data.GetIntFromInterface(fieldData["status"])).
			AddField("title_en", fmt.Sprintf("%v", fieldData["title_en"])).
			AddField("title_ko", fmt.Sprintf("%v", fieldData["title_ko"])).
			AddField("title_ja", fmt.Sprintf("%v", fieldData["title_ja"])).
			AddField("description_en", fmt.Sprintf("%v", fieldData["description_en"])).
			AddField("description_ko", fmt.Sprintf("%v", fieldData["description_ko"])).
			AddField("description_ja", fmt.Sprintf("%v", fieldData["description_ja"])).
			SetTime(timestamp)

		allPoints = append(allPoints, point)
	}

	// 2. Filtered tree points (detail)
	for _, parentNode := range filteredTree {
		children, ok := parentNode["children"].([]map[string]interface{})
		if !ok {
			continue
		}
		parentName := fmt.Sprintf("%v", parentNode["Name"])

		for _, child := range children {
			isSub, _ := child["isSub"].(bool)
			if isSub {
				if subChildren, ok := child["children"].([]map[string]interface{}); ok {
					for _, sc := range subChildren {
						allPoints = append(allPoints, createDetailPoint(meas, chName, asset, assetType, parentName, sc, timestamp))
					}
				}
			} else {
				allPoints = append(allPoints, createDetailPoint(meas, chName, asset, assetType, parentName, child, timestamp))
			}
		}
	}

	// 3. Batch write
	if len(allPoints) > 0 {
		err := pool.WithConnection(func(entry *handlers.PoolEntry) error {
			writeAPI := entry.Client.WriteAPIBlocking(pool.Org, "ntek")
			for _, pt := range allPoints {
				if err := writeAPI.WritePoint(context.Background(), pt); err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			log.Printf("%s InfluxDB save failed (%s): %v", strings.ToUpper(meas), chName, err)
		} else {
			log.Printf("%s InfluxDB saved: %s, %d points", strings.ToUpper(meas), chName, len(allPoints))
		}
	}
}

func createDetailPoint(meas, chName, asset, assetType, parentName string, item map[string]interface{}, timestamp time.Time) *write.Point {
	point := influxdb2.NewPointWithMeasurement(meas).
		AddTag("channel", chName).
		AddTag("asset_name", asset).
		AddTag("data_type", "detail").
		AddTag("item_name", fmt.Sprintf("%v", item["Name"])).
		AddTag("parent_name", parentName).
		AddField("item_id", data.GetIntFromInterface(item["ID"])).
		AddField("asset_type", assetType).
		AddField("assembly_id", fmt.Sprintf("%v", item["AssemblyID"])).
		AddField("status", data.GetIntFromInterface(item["Status"])).
		SetTime(timestamp)

	if title := item["Title"]; title != nil {
		point.AddField("title", fmt.Sprintf("%v", title))
	}
	if titles, ok := item["Titles"].(map[string]interface{}); ok {
		for _, lang := range []string{"en", "ko", "ja"} {
			if v, ok := titles[lang]; ok {
				point.AddField("title_"+lang, fmt.Sprintf("%v", v))
			}
		}
	}
	if desc := item["Description"]; desc != nil {
		point.AddField("description", fmt.Sprintf("%v", desc))
	}
	if descs, ok := item["Descriptions"].(map[string]interface{}); ok {
		for _, lang := range []string{"en", "ko", "ja"} {
			if v, ok := descs[lang]; ok {
				point.AddField("description_"+lang, fmt.Sprintf("%v", v))
			}
		}
	}
	if val := item["Value"]; val != nil && fmt.Sprintf("%v", val) != "NaN" {
		point.AddField("value", fmt.Sprintf("%v", val))
	}

	return point
}

// DiagnosisProcessor: manages periodic diagnosis processing as part of DataProcessorManager
type DiagnosisProcessor struct {
	InfluxHandler *handlers.InfluxDBHandler
	mu            sync.RWMutex
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

func NewDiagnosisProcessor(influx *handlers.InfluxDBHandler) *DiagnosisProcessor {
	return &DiagnosisProcessor{InfluxHandler: influx}
}

func (dp *DiagnosisProcessor) Start(ctx context.Context) {
	log.Println("[DiagnosisProcessor] started (managed via main goroutines)")
}

func (dp *DiagnosisProcessor) Stop() {
	log.Println("[DiagnosisProcessor] stopped")
}

func (dp *DiagnosisProcessor) GetStatus() map[string]interface{} {
	return map[string]interface{}{"status": "running"}
}
