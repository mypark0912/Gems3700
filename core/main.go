package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"sv500_core/data"
	"sv500_core/handlers"
	"sv500_core/processors"

	"github.com/redis/go-redis/v9"
)

const (
	SettingFolder = "/home/root/config"
	LogDir        = "/usr/local/sv500/logs/core"
	LogFile       = "core"
)

var (
	currentChannels    map[string]*data.Channel
	trendCancels       = make(map[string]context.CancelFunc)
	diagnosisCancels   = make(map[string]context.CancelFunc)
	diagnosis1hCancels = make(map[string]context.CancelFunc)
	energyTrendCancels = make(map[string]context.CancelFunc)
	threadMu           sync.Mutex
	mainStopCtx        context.Context
	mainStopCancel     context.CancelFunc
)

// ---------------------------------------------------------------------------
// Logging
// ---------------------------------------------------------------------------

func setupLogging() {
	os.MkdirAll(LogDir, 0755)
	logFilename := filepath.Join(LogDir, fmt.Sprintf("%s_%s.log", LogFile, time.Now().Format("20060102")))

	f, err := os.OpenFile(logFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}

	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// ---------------------------------------------------------------------------
// Memory monitoring
// ---------------------------------------------------------------------------

func logMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("Memory: Alloc=%dMB, Sys=%dMB, NumGoroutine=%d",
		m.Alloc/1024/1024, m.Sys/1024/1024, runtime.NumGoroutine())
}

// ---------------------------------------------------------------------------
// Service monitoring
// ---------------------------------------------------------------------------

func monitorServices(ctx context.Context, services []map[string]string, redisInst *redis.Client) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		statusResult := make(map[string]int)
		for _, svc := range services {
			// On Linux, check systemctl is-active
			statusResult[svc["label"]] = 1 // simplified
		}

		jsonData, _ := json.Marshal(statusResult)
		redisInst.HSet(context.Background(), "System", "Status", string(jsonData))
		logMemoryUsage()

		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Minute):
		}
	}
}

// ---------------------------------------------------------------------------
// Trend threads
// ---------------------------------------------------------------------------

func runTrendWithStop(ctx context.Context, channel *data.Channel, redisInst *redis.Client) {
	trend := data.NewTrend(channel.Name, channel.AssetDrive, channel.VoltageType)
	trend.SetCollectKeys(channel.TrendList)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		start := time.Now()

		dataDict := trend.CopyHash(trend.CollectList, redisInst, "", false)
		trend.SaveInflux(dataDict, channel.Name, "trend")

		elapsed := time.Since(start)
		sleepTime := time.Duration(channel.Period)*time.Second - elapsed
		if sleepTime < 500*time.Millisecond {
			sleepTime = 500 * time.Millisecond
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(sleepTime):
		}
	}
}

func runEnergyTrendWithStop(ctx context.Context, channel *data.Channel, redisInst *redis.Client) {
	period := 15 * time.Minute

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		start := time.Now()

		redisKey := "energy_main"
		if channel.Name != "Main" {
			redisKey = "energy_sub"
		}
		energyTrend := data.NewTrend(channel.Name, channel.AssetDrive, channel.VoltageType)
		energyTrend.SetCollectKeys([]string{"Energy"})
		dataDict := energyTrend.CopyHash(energyTrend.CollectList, redisInst, redisKey, true)
		energyTrend.SaveInflux(dataDict, channel.Name, "energy_trend")

		elapsed := time.Since(start)
		sleepTime := period - elapsed
		if sleepTime < 500*time.Millisecond {
			sleepTime = 500 * time.Millisecond
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(sleepTime):
		}
	}
}

// ---------------------------------------------------------------------------
// Diagnosis threads
// ---------------------------------------------------------------------------

func runDiagnosisWithStop(ctx context.Context, channel *data.Channel, redisInst *redis.Client) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		processors.ProcessAllDiagnosisData(channel, redisInst)

		// Wait 180 seconds with cancellation check
		select {
		case <-ctx.Done():
			return
		case <-time.After(180 * time.Second):
		}
	}
}

func run6hDiagnosisWithStop(ctx context.Context, channel *data.Channel) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		processors.Process1HDiagnosisData(channel.Name, channel.AssetName, channel.AssetType)

		// Wait 6 hours
		select {
		case <-ctx.Done():
			return
		case <-time.After(6 * time.Hour):
		}
	}
}

// ---------------------------------------------------------------------------
// Thread launchers
// ---------------------------------------------------------------------------

func launchTrendThreads(channelDict map[string]*data.Channel, redisInst *redis.Client) {
	threadMu.Lock()
	defer threadMu.Unlock()

	channelOrder := []string{"Main", "Sub"}

	for _, name := range channelOrder {
		ch, ok := channelDict[name]
		if !ok || !ch.Trend {
			continue
		}

		threadName := "trend_" + name
		if _, exists := trendCancels[threadName]; exists {
			log.Printf("Channel '%s' trend thread already running", name)
			continue
		}

		log.Printf("Starting channel '%s' trend thread", name)
		ctx, cancel := context.WithCancel(mainStopCtx)
		trendCancels[threadName] = cancel

		go func(ch *data.Channel) {
			runTrendWithStop(ctx, ch, redisInst)
		}(ch)
	}
}

func launchDiagnosisThreads(channelDict map[string]*data.Channel, redisInst *redis.Client) {
	threadMu.Lock()
	defer threadMu.Unlock()

	channelOrder := []string{"Main", "Sub"}

	for _, name := range channelOrder {
		ch, ok := channelDict[name]
		if !ok || !ch.Diagnosis {
			continue
		}

		threadName := "diagnosis_" + name
		if _, exists := diagnosisCancels[threadName]; exists {
			continue
		}

		log.Printf("Starting channel '%s' diagnosis thread", name)
		ctx, cancel := context.WithCancel(mainStopCtx)
		diagnosisCancels[threadName] = cancel

		go func(ch *data.Channel) {
			runDiagnosisWithStop(ctx, ch, redisInst)
		}(ch)
	}
}

func launchDiagnosis6hThreads(channelDict map[string]*data.Channel) {
	threadMu.Lock()
	defer threadMu.Unlock()

	channelOrder := []string{"Main", "Sub"}

	for _, name := range channelOrder {
		ch, ok := channelDict[name]
		if !ok || !ch.Diagnosis {
			continue
		}

		threadName := "diagnosis_1h_" + name
		if _, exists := diagnosis1hCancels[threadName]; exists {
			continue
		}

		log.Printf("Starting channel '%s' 6h diagnosis thread", name)
		ctx, cancel := context.WithCancel(mainStopCtx)
		diagnosis1hCancels[threadName] = cancel

		go func(ch *data.Channel) {
			run6hDiagnosisWithStop(ctx, ch)
		}(ch)
	}
}

func launchEnergyTrendThreads(channelDict map[string]*data.Channel, redisInst *redis.Client) {
	threadMu.Lock()
	defer threadMu.Unlock()

	channelOrder := []string{"Main", "Sub"}

	for _, name := range channelOrder {
		ch, ok := channelDict[name]
		if !ok {
			continue
		}

		threadName := "energy_trend_" + name
		if _, exists := energyTrendCancels[threadName]; exists {
			continue
		}

		log.Printf("Starting channel '%s' energy trend thread (15min period)", name)
		ctx, cancel := context.WithCancel(mainStopCtx)
		energyTrendCancels[threadName] = cancel

		go func(ch *data.Channel) {
			runEnergyTrendWithStop(ctx, ch, redisInst)
		}(ch)
	}
}

// ---------------------------------------------------------------------------
// Stop helpers
// ---------------------------------------------------------------------------

func stopCancelMap(name string, cancels map[string]context.CancelFunc) {
	threadMu.Lock()
	defer threadMu.Unlock()

	for key, cancel := range cancels {
		cancel()
		delete(cancels, key)
	}
	log.Printf("All %s threads stopped", name)
}

func stopAllThreads() {
	if mainStopCancel != nil {
		mainStopCancel()
	}
	stopCancelMap("trend", trendCancels)
	stopCancelMap("diagnosis", diagnosisCancels)
	stopCancelMap("diagnosis_1h", diagnosis1hCancels)
	stopCancelMap("energy_trend", energyTrendCancels)

	// Re-create main stop context for potential restart
	mainStopCtx, mainStopCancel = context.WithCancel(context.Background())
	log.Println("All threads stopped")
}

func restartAllThreads(newChannelDict map[string]*data.Channel, redisInst *redis.Client, redisDiagnosis *redis.Client) {
	log.Println("Restarting all threads...")
	stopAllThreads()
	time.Sleep(2 * time.Second)

	launchTrendThreads(newChannelDict, redisInst)
	launchDiagnosisThreads(newChannelDict, redisDiagnosis)
	launchDiagnosis6hThreads(newChannelDict)
	launchEnergyTrendThreads(newChannelDict, redisInst)

	currentChannels = newChannelDict
	log.Println("All threads restarted")
}

// ---------------------------------------------------------------------------
// Alarm config init
// ---------------------------------------------------------------------------

func initAlarmConfig(redisD *redis.Client, setup map[string]interface{}) {
	channels, ok := setup["channel"].([]interface{})
	if !ok {
		return
	}

	ctx := context.Background()
	for _, chRaw := range channels {
		ch, ok := chRaw.(map[string]interface{})
		if !ok {
			continue
		}
		chName, _ := ch["channel"].(string)
		alarms, ok := ch["alarm"].(map[string]interface{})
		if !ok {
			continue
		}

		redisKey := fmt.Sprintf("alarm_status:%s", chName)
		redisD.Del(ctx, redisKey)

		for j := 1; j <= 32; j++ {
			key := fmt.Sprintf("%d", j)
			alarmArr, ok := alarms[key].([]interface{})
			if !ok || len(alarmArr) < 4 {
				continue
			}

			chanVal, _ := alarmArr[0].(float64)
			if chanVal == 0 {
				continue
			}

			condVal, _ := alarmArr[1].(float64)
			condStr := "OVER"
			if condVal == 0 {
				condStr = "UNDER"
			}
			levelVal, _ := alarmArr[3].(float64)

			initData := map[string]interface{}{
				"status":      0,
				"count":       0,
				"condition":   condStr,
				"value":       0,
				"chan":         int(chanVal),
				"cond":        int(condVal),
				"level":       levelVal,
				"chan_text":    "",
				"last_update": time.Now().Unix(),
				"status_text": "None",
			}

			jsonBytes, _ := json.Marshal(initData)
			redisD.HSet(ctx, redisKey, key, string(jsonBytes))
		}
	}
}

// ---------------------------------------------------------------------------
// Setup restart monitor
// ---------------------------------------------------------------------------

func monitorRestartFlag(redisInst *redis.Client, redisBinary *redis.Client, redisDiagnosis *redis.Client,
	ctx context.Context, processorManager *processors.DataProcessorManager) {

	restartCount := 0

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		bgCtx := context.Background()

		// Check setting flag
		if redisInst.HExists(bgCtx, "Service", "setting").Val() {
			setFlag, _ := redisInst.HGet(bgCtx, "Service", "setting").Int()
			if setFlag == 1 {
				log.Println("Receive: Modbus Setting Command")
				handleSettingChange(redisInst, redisBinary, redisDiagnosis)
			}
		}

		// Check factory reset flag
		if redisInst.HExists(bgCtx, "Service", "fdreset").Val() {
			resetFlag, _ := redisInst.HGet(bgCtx, "Service", "fdreset").Int()
			if resetFlag == 1 {
				log.Println("Received FD Restart flag")
				handleFactoryReset(redisInst)
			}
		}

		// Check restart flag
		if redisInst.HExists(bgCtx, "Service", "restart").Val() {
			flag, _ := redisInst.HGet(bgCtx, "Service", "restart").Int()
			if flag == 1 {
				restartCount++
				log.Printf("Setup change detected #%d", restartCount)
				redisInst.HSet(bgCtx, "Service", "restart", 0)

				newChannelDict := data.GetSetup(redisInst)
				if processorManager != nil {
					log.Println("Restarting data processors...")
					processorManager.StopAll()
					channels := make([]data.Channel, 0, len(newChannelDict))
					for _, ch := range newChannelDict {
						channels = append(channels, *ch)
					}
					processorManager.UpdateDemandSettings(1*time.Second, channels)
					time.Sleep(1 * time.Second)
					processorManager.StartAll()
					log.Println("Data processors restarted")
				}

				if newChannelDict != nil {
					restartAllThreads(newChannelDict, redisInst, redisDiagnosis)
				}

				log.Printf("Setup restart complete #%d", restartCount)
				time.Sleep(30 * time.Second)
			}
		}

		logMemoryUsage()
		time.Sleep(2 * time.Second)
	}
}

func handleSettingChange(redisInst *redis.Client, redisBinary *redis.Client, redisDiagnosis *redis.Client) {
	bgCtx := context.Background()

	existSetupJSON, _ := redisInst.HGet(bgCtx, "System", "setup").Result()
	var existSetup map[string]interface{}
	json.Unmarshal([]byte(existSetupJSON), &existSetup)

	reader := processors.NewSettingsRedisReader(handlers.NewRedisHandler(redisBinary))

	settingsMain, err := reader.ReadSettings(context.Background(), "SYS_CFG", "setup_main")
	if err != nil {
		log.Printf("Failed to read setup_main: %v", err)
		return
	}
	newSetup := data.UpdateGeneralFromSettings(existSetup, settingsMain, "Main")
	newSetup = data.UpdateChannelFromSetting(newSetup, settingsMain, "Main")
	newSetup = data.UpdateChannelEventFromSetting(newSetup, settingsMain, "Main")

	settingsSub, err := reader.ReadSettings(context.Background(), "SYS_CFG", "setup_sub")
	if err != nil {
		log.Printf("Failed to read setup_sub: %v", err)
	} else {
		newSetup = data.UpdateGeneralFromSettings(newSetup, settingsSub, "Sub")
		newSetup = data.UpdateChannelFromSetting(newSetup, settingsSub, "Sub")
		newSetup = data.UpdateChannelEventFromSetting(newSetup, settingsSub, "Sub")
	}

	// Save to file
	settingPath := filepath.Join(SettingFolder, "setup.json")
	jsonBytes, _ := json.MarshalIndent(newSetup, "", "  ")
	os.WriteFile(settingPath, jsonBytes, 0644)

	initAlarmConfig(redisDiagnosis, newSetup)

	setupJSON, _ := json.Marshal(newSetup)
	redisInst.HSet(bgCtx, "System", "setup", string(setupJSON))
	redisInst.HSet(bgCtx, "Service", "save", 1)
	redisInst.HSet(bgCtx, "Service", "restart", 1)
	redisInst.HSet(bgCtx, "Service", "setting", 0)
}

func handleFactoryReset(redisInst *redis.Client) {
	bgCtx := context.Background()

	settingPath := filepath.Join(SettingFolder, "setup.json")
	backupPath := filepath.Join(SettingFolder, "setup_backup.json")
	defaultPath := filepath.Join(SettingFolder, "default.json")

	// Backup current setup
	if _, err := os.Stat(settingPath); err == nil {
		data, _ := os.ReadFile(settingPath)
		os.WriteFile(backupPath, data, 0644)
		os.Remove(settingPath)
	}

	// Load default
	defaultData, err := os.ReadFile(defaultPath)
	if err != nil {
		log.Printf("Failed to read default.json: %v", err)
		return
	}

	var defaults map[string]interface{}
	json.Unmarshal(defaultData, &defaults)

	// Preserve mode and lang
	if redisInst.HExists(bgCtx, "System", "setup").Val() {
		nowSetupJSON, _ := redisInst.HGet(bgCtx, "System", "setup").Result()
		var nowSetup map[string]interface{}
		json.Unmarshal([]byte(nowSetupJSON), &nowSetup)

		if mode, ok := nowSetup["mode"]; ok {
			defaults["mode"] = mode
		}
		if lang, ok := nowSetup["lang"]; ok {
			defaults["lang"] = lang
		}
	}

	jsonBytes, _ := json.MarshalIndent(defaults, "", "  ")
	os.WriteFile(settingPath, jsonBytes, 0644)

	setupJSON, _ := json.Marshal(defaults)
	redisInst.HDel(bgCtx, "System", "setup")
	redisInst.HSet(bgCtx, "System", "setup", string(setupJSON))
	redisInst.HSet(bgCtx, "Service", "fdreset", 0)
	redisInst.HSet(bgCtx, "Service", "save", 1)
	redisInst.HSet(bgCtx, "Service", "restart", 1)
}

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

func main() {
	setupLogging()

	// Initialize main stop context
	mainStopCtx, mainStopCancel = context.WithCancel(context.Background())

	// Redis clients
	redisInst := handlers.GetRedisClient("127.0.0.1", 0)
	redisDiagnosis := handlers.GetRedisClient("127.0.0.1", 1)
	redisBinary := handlers.GetRedisBinaryClient("127.0.0.1", 0)

	redisHandler := handlers.NewRedisHandler(redisBinary)
	_ = redisHandler

	// Get FTP setup
	setup := data.GetFtp(redisInst)
	mode, _ := setup["mode"].(string)

	// Service monitoring
	var services []map[string]string
	if mode == "device0" {
		services = []map[string]string{
			{"name": "influxdb", "label": "InfluxDB"},
			{"name": "redis", "label": "Redis"},
		}
	} else {
		services = []map[string]string{
			{"name": "influxdb", "label": "InfluxDB"},
			{"name": "redis", "label": "Redis"},
			{"name": "smartsystemsservice", "label": "SmartSystems"},
			{"name": "smartsystemsrestapiservice", "label": "SmartAPI"},
		}
	}

	// Initialize InfluxDB connection pool
	_, err := handlers.GetInfluxPool()
	if err != nil {
		log.Printf("InfluxDB connection pool init failed: %v", err)
		fmt.Println("InfluxDB connection failed, exiting.")
		handlers.CleanupRedisPools()
		os.Exit(1)
	}
	log.Println("InfluxDB connection pool initialized")

	// Get channel setup
	channelDict := data.GetSetup(redisInst)
	currentChannels = channelDict

	// Launch data processors
	influxHandler := handlers.NewInfluxDBHandler("ntek")
	binaryRedisHandler := handlers.NewRedisHandler(handlers.GetRedisBinaryClient("127.0.0.1", 1))
	processorManager := processors.NewDataProcessorManager(binaryRedisHandler, influxHandler)
	{
		channels := make([]data.Channel, 0)
		if channelDict != nil {
			for _, ch := range channelDict {
				channels = append(channels, *ch)
			}
		}
		processorManager.Initialize(channels)
	}
	processorManager.StartAll()
	log.Println("Data processor manager started")

	// Service monitor goroutine
	serviceCtx, serviceCancel := context.WithCancel(context.Background())
	go monitorServices(serviceCtx, services, redisInst)

	// Launch channel threads
	if channelDict != nil {
		launchTrendThreads(channelDict, redisDiagnosis)
		launchDiagnosisThreads(channelDict, redisDiagnosis)
		launchDiagnosis6hThreads(channelDict)
		launchEnergyTrendThreads(channelDict, redisDiagnosis)
		fmt.Printf("Channel data processing started (%s mode)\n", mode)
	}

	// Setup restart monitor
	monitorCtx, monitorCancel := context.WithCancel(context.Background())
	go monitorRestartFlag(redisInst, redisBinary, redisDiagnosis, monitorCtx, processorManager)
	fmt.Printf("Setup monitoring started (%s mode)\n", mode)

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logMemoryUsage()

	sig := <-sigChan
	fmt.Printf("\nReceived signal %v, shutting down...\n", sig)

	// Graceful shutdown
	monitorCancel()
	serviceCancel()
	stopAllThreads()

	if processorManager != nil {
		processorManager.StopAll()
	}

	if pool, err := handlers.GetInfluxPool(); err == nil {
		pool.CloseAllConnections()
	}

	handlers.CleanupRedisPools()

	_ = strings.TrimSpace("") // keep import
	_ = influxHandler

	fmt.Println("Shutdown complete")
}
