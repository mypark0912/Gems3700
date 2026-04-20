package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"serverGO/binary"
)

// pqRedisKey returns "pq_main" or "pq_sub" matching Python get_RedisKey(channel, 'pq').
func pqRedisKey(channel string) string {
	ch := strings.ToLower(channel)
	if ch != "main" && ch != "sub" {
		ch = "main"
	}
	return "pq_" + ch
}

// GET /api/getHarmonics/:channel
// Mirrors Python getHarmonics: hget(pq_<channel>, "harmonics") → JSON map → values /100.
func (h *Handler) getHarmonics(c *gin.Context) {
	channel := c.Param("channel")
	ctx := context.Background()

	raw, err := h.deps.Redis.Client0.HGet(ctx, pqRedisKey(channel), "harmonics").Result()
	if err != nil || raw == "" {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": "Redis Read Error"})
		return
	}

	var harmDict map[string][]float64
	if err := json.Unmarshal([]byte(raw), &harmDict); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": "Redis Read Error"})
		return
	}

	scaled := make(map[string][]float64, len(harmDict))
	for key, arr := range harmDict {
		vals := make([]float64, len(arr))
		for i, v := range arr {
			vals[i] = v / 100
		}
		scaled[key] = vals
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": scaled})
}

// GET /api/getWave/:channel
// Mirrors Python get_waveform: hget(pq_<channel>, "waveform") → WAVEFORM_L16 3000-byte blob.
func (h *Handler) getWave(c *gin.Context) {
	channel := c.Param("channel")
	ctx := context.Background()

	data, err := binary.FetchWaveform(ctx, h.deps.Redis.Client0, channel)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "No waveform data found for channel: " + channel,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}
