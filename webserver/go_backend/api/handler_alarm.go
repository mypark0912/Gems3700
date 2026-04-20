package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /api/getAlarmStatus/:channel
// Mirrors Python getAlarmStatus: HGETALL on DB1 key "alarm_status:<channel>",
// returns the raw hash as-is.
func (h *Handler) getAlarmStatus(c *gin.Context) {
	channel := c.Param("channel")
	ctx := context.Background()
	client := h.deps.Redis.Client1

	key := "alarm_status:" + channel
	exists, _ := client.Exists(ctx, key).Result()
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": "No Data"})
		return
	}

	data, err := client.HGetAll(ctx, key).Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": "Redis Read Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}
