package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gin-gonic/gin"
	"serverGO/api"
	"serverGO/auth"
"serverGO/config"
	"serverGO/crypto"
	"serverGO/infra"
	"serverGO/middleware"
	"serverGO/report"
	"serverGO/setting"
)

func main() {
	// 1. Load configuration
	cfg := config.Load()
	log.Printf("OS: %s, Mode: %d, Redis: %s, ConfigDir: %s", cfg.OS, cfg.Mode, cfg.RedisIP, cfg.ConfigDir)

	// 2. Initialize AES cipher
	cipher, err := crypto.NewAESCipher(cfg.ConfigDir)
	if err != nil {
		log.Printf("Warning: AES cipher init: %v", err)
	}

	// 3. Initialize Redis
	redisState, err := infra.NewRedisState(cfg)
	if err != nil {
		log.Fatalf("Redis init failed: %v", err)
	}
	defer redisState.Close()

	// 4. Initialize InfluxDB (lazy)
	influxState := infra.NewInfluxState(cfg, cipher)
	defer influxState.Close()

	// 5. Wire dependencies
	deps := &infra.Dependencies{
		Config: cfg,
		Redis:  redisState,
		Influx: influxState,
		Crypto: cipher,
	}

	// 6. Init setup (load setup.json → Redis)
	setting.InitSetup(deps)

	// 7. Create Gin engine
	r := gin.Default()

	// 8. Middleware
	r.Use(middleware.Session())
	r.Use(middleware.ClearOldCookies())
	r.Use(middleware.CORS())
	r.Use(middleware.NoCache())
	r.Use(middleware.SessionRefresh())

	// 9. Register routes
	auth.RegisterRoutes(r.Group("/auth"), deps)
	setting.RegisterRoutes(r.Group("/setting"), deps)
	api.RegisterRoutes(r.Group("/api"), deps)
	setting.RegisterConfigRoutes(r.Group("/config"), deps)
	report.RegisterRoutes(r.Group("/report"), deps)

	// 10. Static file serving (Vue.js SPA)
	distDir := "/home/root/webserver/frontend/dist"
	assetsDir := filepath.Join(distDir, "assets")
	indexFile := filepath.Join(distDir, "index.html")

	log.Printf("Vue dist: %s", distDir)

	if _, err := os.Stat(assetsDir); err == nil {
		r.Static("/assets", assetsDir)
	}

	// Root
	r.GET("/", func(c *gin.Context) {
		if _, err := os.Stat(indexFile); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "index.html not found"})
			return
		}
		c.File(indexFile)
	})

	// Vue Router history mode catch-all
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Skip API routes
		if len(path) > 4 && path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		// Try static file first
		staticPath := filepath.Join(distDir, path)
		if info, err := os.Stat(staticPath); err == nil && !info.IsDir() {
			c.File(staticPath)
			return
		}

		// Fallback to index.html
		if _, err := os.Stat(indexFile); err == nil {
			c.File(indexFile)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	// 11. Graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")
		redisState.Close()
		influxState.Close()
		os.Exit(0)
	}()

	// 12. Start server
	log.Println("Server starting on :9000")
	if err := r.Run(":9000"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
