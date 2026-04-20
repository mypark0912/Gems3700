package setting

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func (h *Handler) registerBackupRoutes(rg *gin.RouterGroup) {
	rg.GET("/Reset", h.resetSetting)
	rg.GET("/ResetAll", h.resetAll)
	rg.POST("/restoreSetting", h.restoreSetting)
	rg.POST("/upload", h.uploadSetting)
}

func (h *Handler) resetSetting(c *gin.Context) {
	ctx := context.Background()
	client := h.deps.Redis.Client0

	// Check modbus lock
	if flag, err := client.HGet(ctx, "Service", "setting").Result(); err == nil {
		if flag == "1" {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "Modbus setting is activated"})
			return
		}
	}

	configDir := h.deps.Config.ConfigDir
	settingPath := filepath.Join(configDir, "setup.json")
	backupPath := filepath.Join(configDir, "setup_backup.json")
	defaultPath := filepath.Join(configDir, "default.json")

	msg := ""

	// Backup current setup.json
	if _, err := os.Stat(settingPath); err == nil {
		copyFile(settingPath, backupPath)
		os.Remove(settingPath)
	} else {
		msg = "No exist setup.json"
	}

	// Reset to default with current mode and MAC
	if setupJSON, err := client.HGet(ctx, "System", "setup").Result(); err == nil {
		var nowSetup map[string]interface{}
		json.Unmarshal([]byte(setupJSON), &nowSetup)

		if defaultData, err := os.ReadFile(defaultPath); err == nil {
			var defaults map[string]interface{}
			json.Unmarshal(defaultData, &defaults)

			defaults["mode"] = nowSetup["mode"]
			if gen, ok := defaults["General"].(map[string]interface{}); ok {
				if di, ok := gen["deviceInfo"].(map[string]interface{}); ok {
					di["mac_address"] = getMAC()
					di["serial_number"] = getMAC()
				}
			}

			out, _ := json.MarshalIndent(defaults, "", "  ")
			os.WriteFile(defaultPath, out, 0644)
		}

		client.HDel(ctx, "System", "setup")
	}

	if msg == "" {
		c.JSON(http.StatusOK, gin.H{"success": true, "msg": "Reset success"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "msg": msg})
	}
}

func (h *Handler) resetAll(c *gin.Context) {
	ctx := context.Background()
	client := h.deps.Redis.Client0
	configDir := h.deps.Config.ConfigDir

	// Modbus 잠금 확인
	if flag, err := client.HGet(ctx, "Service", "setting").Result(); err == nil {
		if flag == "1" {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "Modbus setting is activated"})
			return
		}
	}

	msg := ""

	settingPath := filepath.Join(configDir, "setup.json")
	backupPath := filepath.Join(configDir, "setup_backup.json")
	defaultPath := filepath.Join(configDir, "default.json")
	userDBPath := filepath.Join(configDir, "user.db")
	bearingDBPath := filepath.Join(configDir, "bearing.db")
	mtDBPath := filepath.Join(configDir, "maintenance.db")

	// setup.json 백업 후 삭제
	if _, err := os.Stat(settingPath); err == nil {
		copyFile(settingPath, backupPath)
		os.Remove(settingPath)
	} else {
		msg = "No exist setup.json"
	}

	// DB 파일 삭제
	if _, err := os.Stat(userDBPath); err == nil {
		os.Remove(userDBPath)
	} else {
		msg += "No exist user.db"
	}
	if _, err := os.Stat(bearingDBPath); err == nil {
		os.Remove(bearingDBPath)
	}
	if _, err := os.Stat(mtDBPath); err == nil {
		os.Remove(mtDBPath)
	}

	// default.json 초기화
	if defaultData, err := os.ReadFile(defaultPath); err == nil {
		var defaults map[string]interface{}
		json.Unmarshal(defaultData, &defaults)

		defaults["mode"] = "device0"
		if gen, ok := defaults["General"].(map[string]interface{}); ok {
			if di, ok := gen["deviceInfo"].(map[string]interface{}); ok {
				di["mac_address"] = ""
				di["serial_number"] = ""
			}
		}

		out, _ := json.MarshalIndent(defaults, "", "  ")
		os.WriteFile(defaultPath, out, 0644)
	}

	// Redis 정리
	client.HDel(ctx, "System", "setup")
	client.HDel(ctx, "System", "mode")
	client.Del(ctx, "Equipment")

	// Core 서비스 시작
	exec.Command("systemctl", "start", "Core").Run()

	if msg == "" {
		c.JSON(http.StatusOK, gin.H{"success": true, "msg": "Reset success"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "msg": msg})
	}
}

func (h *Handler) restoreSetting(c *gin.Context) {
	ctx := context.Background()
	client := h.deps.Redis.Client0

	// Check modbus lock
	if flag, err := client.HGet(ctx, "Service", "setting").Result(); err == nil {
		if flag == "1" {
			c.JSON(http.StatusOK, gin.H{"status": "0", "error": "Modbus setting is activated"})
			return
		}
	}

	configDir := h.deps.Config.ConfigDir
	settingPath := filepath.Join(configDir, "setup.json")
	backupPath := filepath.Join(configDir, "setup_backup.json")

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		c.JSON(http.StatusBadRequest, gin.H{"passOK": "0", "error": "setup_backup.json not found"})
		return
	}

	// Copy backup → setup
	if err := copyFile(backupPath, settingPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"passOK": "0", "error": err.Error()})
		return
	}

	// Load backup into Redis
	data, err := os.ReadFile(backupPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"passOK": "0", "error": err.Error()})
		return
	}

	client.HSet(ctx, "System", "setup", string(data))
	client.HSet(ctx, "Service", "save", 1)
	client.HSet(ctx, "Service", "restart", 1)

	c.JSON(http.StatusOK, gin.H{"passOK": "1", "message": "Restore successful"})
}

func (h *Handler) uploadSetting(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil || file.Filename == "" {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "error": "No selected file"})
		return
	}

	ctx := context.Background()
	configDir := h.deps.Config.ConfigDir
	settingPath := filepath.Join(configDir, "setup.json")
	backupPath := filepath.Join(configDir, "setup_backup.json")

	// Backup existing
	if _, err := os.Stat(settingPath); err == nil {
		copyFile(settingPath, backupPath)
	}

	// Save uploaded file
	if err := c.SaveUploadedFile(file, settingPath); err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": "0", "error": err.Error()})
		return
	}

	// Read and store in Redis
	data, _ := os.ReadFile(settingPath)
	h.deps.Redis.Client0.HSet(ctx, "System", "setup", string(data))
	h.deps.Redis.Client0.HSet(ctx, "Service", "save", 1)
	h.deps.Redis.Client0.HSet(ctx, "Service", "restart", 1)

	c.JSON(http.StatusOK, gin.H{"passOK": "1", "file_path": settingPath})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
