package cfgroute

import (
	"database/sql"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"serverGO/infra"
)

type Handler struct {
	deps          *infra.Dependencies
	maintenanceDB *sql.DB
	logDB         *sql.DB
}

func RegisterRoutes(rg *gin.RouterGroup, deps *infra.Dependencies) {
	mdb, err := openMaintenanceDB(deps.Config.ConfigDir)
	if err != nil {
		log.Printf("cfgroute: open maintenance.db: %v", err)
	}
	ldb, err := openLogDB(deps.Config.ConfigDir)
	if err != nil {
		log.Printf("cfgroute: open log.db: %v", err)
	}

	h := &Handler{deps: deps, maintenanceDB: mdb, logDB: ldb}

	// Maintenance
	rg.POST("/savePost/:mode/:idx", h.savePostHandler)
	rg.GET("/getPost", h.getPost)
	rg.GET("/getLastPost", h.getLastPost)
	rg.GET("/deletePost/:idx", h.deletePost)

	// Log
	rg.GET("/getLog", h.getLog)

	// Time
	rg.GET("/calibrate/gettime", h.getTime)
	rg.POST("/calibrate/setSystemTime", h.setSystemTime)
	rg.POST("/checktime", h.checkTime)
}

func (h *Handler) savePostHandler(c *gin.Context) {
	mode, _ := strconv.Atoi(c.Param("mode"))
	idx, _ := strconv.Atoi(c.Param("idx"))

	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusOK, gin.H{"result": 0, "msg": err.Error()})
		return
	}

	if err := savePost(h.maintenanceDB, post, mode, idx); err != nil {
		c.JSON(http.StatusOK, gin.H{"result": 0, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": 1})
}

func (h *Handler) getPost(c *gin.Context) {
	posts, err := getAllPosts(h.maintenanceDB)
	if err != nil || len(posts) == 0 {
		c.JSON(http.StatusOK, gin.H{"result": 0})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": 1, "data": posts})
}

func (h *Handler) getLastPost(c *gin.Context) {
	post, err := getLastPost(h.maintenanceDB)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": 0})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": 1, "data": post})
}

func (h *Handler) deletePost(c *gin.Context) {
	idx := c.Param("idx")
	if err := deletePost(h.maintenanceDB, idx); err != nil {
		c.JSON(http.StatusOK, gin.H{"result": 0})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": 1})
}

func (h *Handler) getLog(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	logs, total, err := getLogsPaginated(h.logDB, page, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": 0, "msg": err.Error()})
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	result := 0
	if len(logs) > 0 {
		result = 1
	}

	c.JSON(http.StatusOK, gin.H{
		"result":      result,
		"data":        logs,
		"total":       total,
		"page":        page,
		"total_pages": totalPages,
	})
}

func (h *Handler) getTime(c *gin.Context) {
	now := time.Now()
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"deviceTime": now.Format(time.RFC3339),
		"timestamp":  now.Unix(),
	})
}

type TimeSetRequest struct {
	DatetimeStr string `json:"datetime_str" binding:"required"`
	Timezone    string `json:"timezone"`
}

func (h *Handler) setSystemTime(c *gin.Context) {
	var req TimeSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}

	if req.Timezone == "" {
		req.Timezone = "Asia/Seoul"
	}

	exec.Command("timedatectl", "set-timezone", req.Timezone).Run()
	exec.Command("date", "-s", req.DatetimeStr).Run()

	current, _ := exec.Command("date").Output()

	// Log action
	session := sessions.Default(c)
	account, _ := session.Get("user").(string)
	role, _ := session.Get("userRole").(string)
	if h.logDB != nil && account != "" {
		updateLog(h.logDB, account, role, "Set Time")
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "System time updated",
		"current_time": strings.TrimSpace(string(current)),
		"timezone":     req.Timezone,
	})
}

func (h *Handler) checkTime(c *gin.Context) {
	var req TimeSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "status": false})
		return
	}

	deviceNow := time.Now()
	clientDt, err := time.Parse("2006-01-02 15:04:05", req.DatetimeStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "status": false})
		return
	}

	diff := deviceNow.Sub(clientDt)
	if diff < 0 {
		diff = -diff
	}
	twentyFourHours := 24 * time.Hour
	status := diff <= twentyFourHours

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"deviceTime":  deviceNow.Format(time.RFC3339),
		"status":      status,
		"diffSeconds": int(diff.Seconds()),
	})
}
