package processors

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// UploadManager manages FTP upload clients per channel.
type UploadManager struct {
	clients map[string]*FTPProcessor
	mu      sync.Mutex
}

// NewUploadManager creates a new UploadManager.
func NewUploadManager() *UploadManager {
	return &UploadManager{
		clients: make(map[string]*FTPProcessor),
	}
}

// ApplySetup configures FTP uploads based on the provided setup map.
// Expected keys in setup:
//
//	"ftp"        - bool: global FTP enable flag
//	"main"       - bool: Main channel enable
//	"sub"        - bool: Sub channel enable
//	"ftpInfo"    - map with host, port, id, pass, isManual, upload_main, upload_sub
//	"deviceInfo" - map with mac_address
//	"sampling"   - map with Main/Sub intervals (int, seconds)
func (u *UploadManager) ApplySetup(setup map[string]interface{}) {
	u.mu.Lock()
	defer u.mu.Unlock()

	// Check global FTP enable flag.
	ftpEnabled, _ := setup["ftp"].(bool)
	if !ftpEnabled {
		u.stopChannelLocked("Main")
		u.stopChannelLocked("Sub")
		return
	}

	// Check FTP connection before proceeding.
	if !u.checkConnection(setup) {
		log.Println("[UploadManager] FTP connection failed")
		u.stopAllChannelsLocked()
		return
	}

	// Stop all existing channels before restarting.
	u.stopAllChannelsLocked()
	time.Sleep(2 * time.Second)

	// Manage each channel.
	mainEnabled, _ := setup["main"].(bool)
	subEnabled, _ := setup["sub"].(bool)

	u.manageChannelLocked("Main", mainEnabled, setup)
	u.manageChannelLocked("Sub", subEnabled, setup)
}

// StopChannel stops the upload client for a specific channel.
func (u *UploadManager) StopChannel(channelName string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.stopChannelLocked(channelName)
}

// StopAllChannels stops all active upload channels.
func (u *UploadManager) StopAllChannels() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.stopAllChannelsLocked()
}

// Start is a convenience alias; actual start happens via ApplySetup.
func (u *UploadManager) Start() {
	// No-op: channels are started through ApplySetup.
}

// Stop stops all channels and clears state.
func (u *UploadManager) Stop() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.stopAllChannelsLocked()
}

// GetClients returns a snapshot of current client references (for status queries).
func (u *UploadManager) GetClients() map[string]*FTPProcessor {
	u.mu.Lock()
	defer u.mu.Unlock()
	snapshot := make(map[string]*FTPProcessor, len(u.clients))
	for k, v := range u.clients {
		snapshot[k] = v
	}
	return snapshot
}

// ---------------------------------------------------------------------------
// Internal helpers (must be called with u.mu held)
// ---------------------------------------------------------------------------

// checkConnection tests FTP connectivity using raw FTP login.
func (u *UploadManager) checkConnection(setup map[string]interface{}) bool {
	ftpInfoRaw, ok := setup["ftpInfo"].(map[string]interface{})
	if !ok {
		log.Println("[UploadManager] missing or invalid ftpInfo")
		return false
	}

	host, _ := ftpInfoRaw["host"].(string)
	portRaw, _ := ftpInfoRaw["port"]
	id, _ := ftpInfoRaw["id"].(string)
	pass, _ := ftpInfoRaw["pass"].(string)

	port := ftpToInt(portRaw)

	// Use a temporary FTPProcessor just to test the connection.
	tmp := &FTPProcessor{
		Host:     host,
		Port:     port,
		User:     id,
		Password: pass,
	}
	if tmp.CheckFTP() {
		return true
	}
	log.Printf("[UploadManager] login failed: %s", id)
	return false
}

// manageChannelLocked starts or stops a channel based on the enabled flag.
func (u *UploadManager) manageChannelLocked(channel string, enabled bool, setup map[string]interface{}) {
	if enabled {
		u.startChannelLocked(channel, setup)
	} else {
		u.stopChannelLocked(channel)
	}
}

// startChannelLocked starts an upload client for the given channel.
func (u *UploadManager) startChannelLocked(channel string, setup map[string]interface{}) {
	// If an old client exists, stop it first.
	if old, exists := u.clients[channel]; exists {
		delete(u.clients, channel)
		old.Stop()
		time.Sleep(2 * time.Second)
	}

	log.Printf("[UploadManager] %s channel upload starting", channel)

	ftpInfo := parseFTPInfo(setup["ftpInfo"])
	mac := parseMAC(setup["deviceInfo"])
	sampling := parseSampling(setup["sampling"])

	ftp := NewFTPProcessor(ftpInfo, mac, sampling)
	ftp.CheckChannel(channel, true)
	ftp.Start()

	u.clients[channel] = ftp
}

// stopChannelLocked stops the client for a channel if it exists.
func (u *UploadManager) stopChannelLocked(channel string) {
	client, exists := u.clients[channel]
	if !exists {
		return
	}
	log.Printf("[UploadManager] %s channel upload stopping", channel)
	delete(u.clients, channel)
	client.Stop()
}

// stopAllChannelsLocked stops all active channels.
func (u *UploadManager) stopAllChannelsLocked() {
	for channel := range u.clients {
		u.stopChannelLocked(channel)
	}
}

// ---------------------------------------------------------------------------
// Setup parsing helpers
// ---------------------------------------------------------------------------

// parseFTPInfo extracts FTPInfo from a raw map.
func parseFTPInfo(raw interface{}) FTPInfo {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return FTPInfo{}
	}
	return FTPInfo{
		Host:       toString(m["host"]),
		Port:       ftpToInt(m["port"]),
		ID:         toString(m["id"]),
		Password:   toString(m["pass"]),
		IsManual:   toBool(m["isManual"]),
		UploadMain: toString(m["upload_main"]),
		UploadSub:  toString(m["upload_sub"]),
	}
}

// parseMAC extracts mac_address from deviceInfo.
func parseMAC(raw interface{}) string {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return ""
	}
	return toString(m["mac_address"])
}

// parseSampling extracts sampling intervals.
func parseSampling(raw interface{}) map[string]int {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return map[string]int{"Main": 60, "Sub": 60}
	}
	return map[string]int{
		"Main": ftpToInt(m["Main"]),
		"Sub":  ftpToInt(m["Sub"]),
	}
}

// toString safely converts an interface{} to string.
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

// toInt safely converts an interface{} to int.
func ftpToInt(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case float64:
		return int(val)
	case int64:
		return int(val)
	case string:
		var n int
		fmt.Sscanf(val, "%d", &n)
		return n
	default:
		return 0
	}
}

// toBool safely converts an interface{} to bool.
func toBool(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}
