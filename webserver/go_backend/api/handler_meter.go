package api

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"serverGO/infra"
	"serverGO/metermap"
)

func tryFloat(s string) float64 {
	if s == "" {
		return 0.0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsNaN(v) {
		return 0.0
	}
	return math.Round(v*100) / 100
}

type Handler struct {
	deps *infra.Dependencies
}

func RegisterRoutes(rg *gin.RouterGroup, deps *infra.Dependencies) {
	h := &Handler{deps: deps}

	// Meter
	rg.GET("/getMeterRedisNew/:channel/:mode", h.getMeterRedisNew)
	rg.GET("/getChannelSetting/:channel", h.getChannelSettingHandler)

	// Harmonics / Waveform
	rg.GET("/getHarmonics/:channel", h.getHarmonics)
	rg.GET("/getWave/:channel", h.getWave)

	// Meter detail
	rg.GET("/getOnesfromRedis/:channel/:unbal", h.getOnesfromRedis)
	rg.GET("/getonemfromRedis/:channel", h.getonemfromRedis)
	rg.GET("/getFifthMfromRedis/:channel", h.getFifthMfromRedis)
	rg.GET("/getOnehfromRedis/:channel", h.getOnehfromRedis)
	rg.GET("/getInteverval/:mode/:channel", h.getInterval)
	rg.GET("/getEnergyRedis/:channel", h.getEnergyRedis)

	// Alarm / Event
	rg.GET("/getAlarmStatus/:channel", h.getAlarmStatus)

	// System
	rg.GET("/getSystemStatus", h.getSystemStatus)
	rg.GET("/getTrendParameters/:channel", h.getTrendParameters)

	// Alarm Log (SQLite)
	h.registerAlarmLogRoutes(rg)
}

func (h *Handler) getMeterRedisNew(c *gin.Context) {
	channel := c.Param("channel")
	mode := c.Param("mode")
	ctx := context.Background()
	client := h.deps.Redis.Client0

	// Python 원본: hgetall("meter") — 키가 "meter" 그대로
	flatData, err := client.HGetAll(ctx, "meter").Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": "Redis Read Error"})
		return
	}

	// dashboard_meter 기준으로 플랫 딕셔너리 구성 (Python 형식과 동일)
	meterdata := make(map[string]interface{})
	var keySet map[string][]string
	if mode == "server" {
		keySet = metermap.DashboardTrans
	} else {
		keySet = metermap.DashboardMeter
	}

	for _, keys := range keySet {
		for _, key := range keys {
			meterdata[key] = tryFloat(flatData[key])
		}
	}

	// THD/TDD 평균값 추가
	thdData, _ := client.HGetAll(ctx, "meter").Result()
	thduTotal := (tryFloat(thdData["THD_U1"]) + tryFloat(thdData["THD_U2"]) + tryFloat(thdData["THD_U3"])) / 3
	thdiTotal := (tryFloat(thdData["THD_I1"]) + tryFloat(thdData["THD_I2"]) + tryFloat(thdData["THD_I3"])) / 3
	tddiTotal := (tryFloat(thdData["TDD_I1"]) + tryFloat(thdData["TDD_I2"]) + tryFloat(thdData["TDD_I3"])) / 3
	meterdata["thdu total"] = math.Round(thduTotal*100) / 100
	meterdata["thdi total"] = math.Round(thdiTotal*100) / 100
	meterdata["tddi total"] = math.Round(tddiTotal*100) / 100

	// DashPT 설정값 조회
	dashPT := 1
	chData := getChannelSetting(ctx, h.deps, channel)
	if chData != nil {
		if v, ok := chData["DashPT"]; ok {
			switch val := v.(type) {
			case float64:
				dashPT = int(val)
			case string:
				fmt.Sscanf(val, "%d", &dashPT)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": meterdata, "DashPT": dashPT})
}

func (h *Handler) getChannelSettingHandler(c *gin.Context) {
	channel := c.Param("channel")
	ctx := context.Background()

	chData := getChannelSetting(ctx, h.deps, channel)
	if chData == nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": "channel not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": chData})
}

func (h *Handler) getTrendParameters(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    metermap.ParameterOptions,
	})
}

func (h *Handler) getSystemStatus(c *gin.Context) {
	ctx := context.Background()

	services := gin.H{
		"Redis":    h.deps.Redis.Client0.Ping(ctx).Err() == nil,
		"InfluxDB": isServiceActive("influxdb"),
		"WebServer": true,
	}

	allOK := true
	for _, v := range services {
		if ok, _ := v.(bool); !ok {
			allOK = false
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": allOK, "services": services})
}
