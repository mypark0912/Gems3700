package setting

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) registerMQTTRoutes(rg *gin.RouterGroup) {
	rg.GET("/checkMQTT", h.checkMQTT)
	rg.POST("/uploadCerts", h.uploadCerts)
	rg.GET("/listCerts", h.listCerts)
	rg.DELETE("/deleteCert/:filename", h.deleteCert)
}

func (h *Handler) checkMQTT(c *gin.Context) {
	ctx := context.Background()
	var setup map[string]interface{}

	if setupJSON, err := h.deps.Redis.Client0.HGet(ctx, "System", "setup").Result(); err == nil {
		json.Unmarshal([]byte(setupJSON), &setup)
	} else {
		setupPath := filepath.Join(h.deps.Config.ConfigDir, "setup.json")
		data, err := os.ReadFile(setupPath)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"passOK": 0, "error": "setting file not found"})
			return
		}
		json.Unmarshal(data, &setup)
	}

	gen, _ := setup["General"].(map[string]interface{})
	if gen == nil || gen["MQTT"] == nil {
		c.JSON(http.StatusOK, gin.H{"passOK": 1, "result": map[string]interface{}{}})
		return
	}

	// Apply MQTT services in background
	go h.applyMQTTServices(setup)

	c.JSON(http.StatusOK, gin.H{"passOK": 1, "result": "processing"})
}

func (h *Handler) uploadCerts(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": 0, "error": "No files provided"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusOK, gin.H{"passOK": 0, "error": "No files provided"})
		return
	}

	allowedExts := map[string]bool{".pem": true, ".crt": true, ".key": true, ".cert": true}
	savePath := filepath.Join(h.deps.Config.ConfigDir, "certs")
	os.MkdirAll(savePath, 0755)

	var uploaded []string
	var errors []map[string]string

	for _, file := range files {
		if file.Filename == "" {
			continue
		}
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedExts[ext] {
			errors = append(errors, map[string]string{"filename": file.Filename, "error": "Invalid file type"})
			continue
		}

		dst := filepath.Join(savePath, file.Filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			errors = append(errors, map[string]string{"filename": file.Filename, "error": err.Error()})
			continue
		}
		os.Chmod(dst, 0600)
		uploaded = append(uploaded, file.Filename)
	}

	passOK := 0
	if len(uploaded) > 0 {
		passOK = 1
	}
	c.JSON(http.StatusOK, gin.H{"passOK": passOK, "uploaded": uploaded, "errors": errors})
}

func (h *Handler) listCerts(c *gin.Context) {
	certPath := filepath.Join(h.deps.Config.ConfigDir, "certs")

	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{"passOK": 1, "files": []interface{}{}})
		return
	}

	entries, _ := os.ReadDir(certPath)
	var files []map[string]interface{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, _ := entry.Info()
		files = append(files, map[string]interface{}{
			"filename": entry.Name(),
			"size":     info.Size(),
			"modified": info.ModTime().Unix(),
		})
	}

	c.JSON(http.StatusOK, gin.H{"passOK": 1, "files": files})
}

func (h *Handler) deleteCert(c *gin.Context) {
	filename := c.Param("filename")

	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		c.JSON(http.StatusOK, gin.H{"passOK": 0, "error": "Invalid filename"})
		return
	}

	filePath := filepath.Join(h.deps.Config.ConfigDir, "certs", filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{"passOK": 0, "error": "File not found"})
		return
	}

	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusOK, gin.H{"passOK": 0, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"passOK": 1, "filename": filename})
}
