package processors

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// WaveFolder is the base directory for local data files.
const WaveFolder = "/sv500/"

// SafeDelay is the minimum age (in seconds) a file must have before upload.
const SafeDelay = 2

// FTPInfo holds FTP connection parameters.
type FTPInfo struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	ID         string `json:"id"`
	Password   string `json:"pass"`
	IsManual   bool   `json:"isManual"`
	UploadMain string `json:"upload_main"`
	UploadSub  string `json:"upload_sub"`
}

// FTPProcessor handles FTP file transfers for a single channel pair.
type FTPProcessor struct {
	Host     string
	Port     int
	User     string
	Password string
	MAC      string

	// Remote upload paths per channel.
	RemotePaths map[string]string

	// Local directories per channel and data type.
	LocalDirs map[string]map[string]string

	// Sampling intervals per channel (seconds).
	Sampling map[string]int

	// Channel enable flags.
	MainEnabled bool
	SubEnabled  bool

	// Runtime state.
	enable bool
	mu     sync.RWMutex

	// Tracks sent files per channel/datatype to avoid re-uploads.
	sentFiles     map[string]map[string][]string
	sentFilesLock sync.Mutex

	// Thread management.
	threads map[string]chan struct{} // stop signals
	wg      sync.WaitGroup
}

// NewFTPProcessor creates a new FTPProcessor from configuration.
func NewFTPProcessor(ftpInfo FTPInfo, mac string, sampling map[string]int) *FTPProcessor {
	remotePaths := map[string]string{}
	if ftpInfo.IsManual {
		remotePaths["Main"] = ftpInfo.UploadMain
		remotePaths["Sub"] = ftpInfo.UploadSub
	} else {
		remotePaths["Main"] = fmt.Sprintf("/%s/ch1", mac)
		remotePaths["Sub"] = fmt.Sprintf("/%s/ch2", mac)
	}

	localDirs := map[string]map[string]string{
		"Main": {
			"waveform": filepath.Join(WaveFolder, "ch1/waveform"),
			"event":    filepath.Join(WaveFolder, "ch1/event"),
		},
		"Sub": {
			"waveform": filepath.Join(WaveFolder, "ch2/waveform"),
			"event":    filepath.Join(WaveFolder, "ch2/event"),
		},
	}

	sentFiles := map[string]map[string][]string{
		"Main": {"waveform": {}, "event": {}},
		"Sub":  {"waveform": {}, "event": {}},
	}

	return &FTPProcessor{
		Host:        ftpInfo.Host,
		Port:        ftpInfo.Port,
		User:        ftpInfo.ID,
		Password:    ftpInfo.Password,
		MAC:         mac,
		RemotePaths: remotePaths,
		LocalDirs:   localDirs,
		Sampling:    sampling,
		enable:      true,
		sentFiles:   sentFiles,
		threads:     make(map[string]chan struct{}),
	}
}

// CheckChannel sets the enabled state of a channel.
func (f *FTPProcessor) CheckChannel(channel string, value bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if channel == "Main" {
		f.MainEnabled = value
	} else {
		f.SubEnabled = value
	}
}

// CheckFTP tests whether FTP login succeeds.
func (f *FTPProcessor) CheckFTP() bool {
	conn, err := f.connect()
	if err != nil {
		log.Printf("[FTP] login failed: %s - %v", f.User, err)
		return false
	}
	f.quit(conn)
	return true
}

// Start begins upload goroutines for enabled channels.
func (f *FTPProcessor) Start() {
	f.mu.RLock()
	mainEnabled := f.MainEnabled
	subEnabled := f.SubEnabled
	f.mu.RUnlock()

	log.Printf("[FTP] start() - enable:%v, main:%v, sub:%v", f.enable, mainEnabled, subEnabled)

	if mainEnabled {
		f.startChannelThread("Main")
	}
	if subEnabled {
		f.startChannelThread("Sub")
	}
}

// Stop halts all upload goroutines and waits for them to finish.
func (f *FTPProcessor) Stop() {
	log.Printf("[FTP] stop() called")
	f.mu.Lock()
	f.enable = false
	// Signal all channel threads to stop.
	for ch, stopCh := range f.threads {
		close(stopCh)
		delete(f.threads, ch)
	}
	f.mu.Unlock()

	f.wg.Wait()
	log.Println("[FTP] client fully stopped")
}

// GetStatus returns the current status of the FTP processor.
func (f *FTPProcessor) GetStatus() map[string]interface{} {
	f.sentFilesLock.Lock()
	sentCounts := map[string]map[string]int{
		"Main": {
			"waveform": len(f.sentFiles["Main"]["waveform"]),
			"event":    len(f.sentFiles["Main"]["event"]),
		},
		"Sub": {
			"waveform": len(f.sentFiles["Sub"]["waveform"]),
			"event":    len(f.sentFiles["Sub"]["event"]),
		},
	}
	f.sentFilesLock.Unlock()

	f.mu.RLock()
	defer f.mu.RUnlock()

	return map[string]interface{}{
		"host":     f.Host,
		"port":     f.Port,
		"enabled":  f.enable,
		"channels": map[string]bool{"main": f.MainEnabled, "sub": f.SubEnabled},
		"sent_files_count": sentCounts,
	}
}

// GetUploadStatistics returns upload counts.
func (f *FTPProcessor) GetUploadStatistics() map[string]int {
	f.sentFilesLock.Lock()
	mw := len(f.sentFiles["Main"]["waveform"])
	me := len(f.sentFiles["Main"]["event"])
	sw := len(f.sentFiles["Sub"]["waveform"])
	se := len(f.sentFiles["Sub"]["event"])
	f.sentFilesLock.Unlock()

	return map[string]int{
		"total_uploaded": mw + me + sw + se,
		"main_total":     mw + me,
		"sub_total":      sw + se,
		"main_waveform":  mw,
		"main_event":     me,
		"sub_waveform":   sw,
		"sub_event":      se,
	}
}

// ClearSentFilesHistory clears recorded sent file paths.
// Pass empty strings to clear all.
func (f *FTPProcessor) ClearSentFilesHistory(channel, dataType string) {
	f.sentFilesLock.Lock()
	defer f.sentFilesLock.Unlock()

	if channel != "" && dataType != "" {
		f.sentFiles[channel][dataType] = nil
	} else if channel != "" {
		f.sentFiles[channel]["waveform"] = nil
		f.sentFiles[channel]["event"] = nil
	} else {
		for _, ch := range []string{"Main", "Sub"} {
			f.sentFiles[ch]["waveform"] = nil
			f.sentFiles[ch]["event"] = nil
		}
	}
	log.Println("[FTP] sent files history cleared")
}

// ---------------------------------------------------------------------------
// Internal: FTP connection via net/textproto
// ---------------------------------------------------------------------------

// ftpConn wraps a textproto connection for basic FTP operations.
type ftpConn struct {
	conn *textproto.Conn
	host string
	port int
}

// connect establishes an FTP control connection and logs in.
func (f *FTPProcessor) connect() (*ftpConn, error) {
	addr := fmt.Sprintf("%s:%d", f.Host, f.Port)
	raw, err := textproto.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("connect %s: %w", addr, err)
	}

	c := &ftpConn{conn: raw, host: f.Host, port: f.Port}

	// Read greeting.
	if _, _, err := raw.ReadResponse(220); err != nil {
		raw.Close()
		return nil, fmt.Errorf("greeting: %w", err)
	}

	// Login.
	if err := c.cmd(331, "USER %s", f.User); err != nil {
		raw.Close()
		return nil, fmt.Errorf("USER: %w", err)
	}
	if err := c.cmd(230, "PASS %s", f.Password); err != nil {
		raw.Close()
		return nil, fmt.Errorf("PASS: %w", err)
	}

	// Binary mode.
	if err := c.cmd(200, "TYPE I"); err != nil {
		// Some servers return 250; try again allowing that.
		_ = c.cmdMulti([]int{200, 250}, "TYPE I")
	}

	return c, nil
}

// cmd sends a command and expects a specific response code.
func (c *ftpConn) cmd(expectCode int, format string, args ...interface{}) error {
	id, err := c.conn.Cmd(format, args...)
	if err != nil {
		return err
	}
	c.conn.StartResponse(id)
	defer c.conn.EndResponse(id)
	_, _, err = c.conn.ReadResponse(expectCode)
	return err
}

// cmdMulti sends a command and accepts any of the given response codes.
func (c *ftpConn) cmdMulti(codes []int, format string, args ...interface{}) error {
	id, err := c.conn.Cmd(format, args...)
	if err != nil {
		return err
	}
	c.conn.StartResponse(id)
	defer c.conn.EndResponse(id)
	code, _, err := c.conn.ReadCodeLine(0)
	if err != nil {
		return err
	}
	for _, c := range codes {
		if code == c {
			return nil
		}
	}
	return fmt.Errorf("unexpected code %d", code)
}

// quit sends QUIT and closes the connection.
func (f *FTPProcessor) quit(c *ftpConn) {
	if c == nil {
		return
	}
	_ = c.cmd(221, "QUIT")
	c.conn.Close()
}

// ---------------------------------------------------------------------------
// Internal: FTP file operations using github.com/jlaffaye/ftp-style raw FTP
// ---------------------------------------------------------------------------

// uploadViaFTP uploads a single file to the FTP server.
func (f *FTPProcessor) uploadViaFTP(filePath, channel, dataType string) bool {
	conn, err := f.connect()
	if err != nil {
		log.Printf("[FTP] connection failed: %v", err)
		return false
	}
	defer func() {
		_, _ = conn.conn.Cmd("QUIT")
		conn.conn.Close()
	}()

	filename := filepath.Base(filePath)
	targetPath := fmt.Sprintf("%s/%s", f.RemotePaths[channel], dataType)

	// Ensure remote directory exists.
	f.ensureFTPPath(conn, targetPath)

	// CWD to target directory.
	_ = conn.cmd(250, "CWD %s", targetPath)

	// Open local file.
	localFile, err := os.Open(filePath)
	if err != nil {
		log.Printf("[FTP] cannot open local file %s: %v", filePath, err)
		return false
	}
	defer localFile.Close()

	// Enter passive mode and get data connection.
	dataConn, err := f.openPassiveDataConn(conn)
	if err != nil {
		log.Printf("[FTP] passive mode failed: %v", err)
		return false
	}

	// Send STOR command.
	id, err := conn.conn.Cmd("STOR %s", filename)
	if err != nil {
		dataConn.Close()
		log.Printf("[FTP] STOR command failed: %v", err)
		return false
	}

	conn.conn.StartResponse(id)
	// Read 150 preliminary response.
	code, _, err := conn.conn.ReadCodeLine(0)
	if err != nil || (code != 150 && code != 125) {
		conn.conn.EndResponse(id)
		dataConn.Close()
		log.Printf("[FTP] unexpected response to STOR: code=%d err=%v", code, err)
		return false
	}

	// Transfer data.
	_, err = io.Copy(dataConn, localFile)
	dataConn.Close()

	if err != nil {
		conn.conn.EndResponse(id)
		log.Printf("[FTP] data transfer failed: %v", err)
		return false
	}

	// Read final 226 response.
	code, _, err = conn.conn.ReadCodeLine(0)
	conn.conn.EndResponse(id)

	if code == 226 {
		savePath := fmt.Sprintf("%s/%s", targetPath, filename)
		log.Printf("[FTP] upload complete: %s %s : %s -> %s", channel, dataType, filePath, savePath)
		return true
	}

	if err != nil {
		log.Printf("[FTP] transfer response error: %v", err)
	} else {
		log.Printf("[FTP] unexpected transfer response: %d", code)
	}

	// For 421 errors, verify upload after error.
	if code == 421 {
		if f.verifyUploadAfterError(filename, targetPath, channel, dataType) {
			log.Printf("[FTP] 421 error but file verified on server")
			return true
		}
	}

	return false
}

// openPassiveDataConn parses PASV response and opens a raw data connection.
func (f *FTPProcessor) openPassiveDataConn(c *ftpConn) (net.Conn, error) {
	id, err := c.conn.Cmd("PASV")
	if err != nil {
		return nil, err
	}
	c.conn.StartResponse(id)
	_, msg, err := c.conn.ReadResponse(227)
	c.conn.EndResponse(id)
	if err != nil {
		return nil, fmt.Errorf("PASV: %w", err)
	}

	// Parse PASV response: 227 Entering Passive Mode (h1,h2,h3,h4,p1,p2).
	start := strings.Index(msg, "(")
	end := strings.Index(msg, ")")
	if start < 0 || end < 0 {
		return nil, fmt.Errorf("cannot parse PASV response: %s", msg)
	}
	parts := strings.Split(msg[start+1:end], ",")
	if len(parts) != 6 {
		return nil, fmt.Errorf("invalid PASV format: %s", msg)
	}

	var nums [6]int
	for i, p := range parts {
		_, _ = fmt.Sscanf(strings.TrimSpace(p), "%d", &nums[i])
	}
	dataAddr := fmt.Sprintf("%d.%d.%d.%d:%d", nums[0], nums[1], nums[2], nums[3], nums[4]*256+nums[5])

	dataConn, err := net.Dial("tcp", dataAddr)
	if err != nil {
		return nil, fmt.Errorf("data connect %s: %w", dataAddr, err)
	}
	return dataConn, nil
}

// ensureFTPPath creates directories along a remote path if they don't exist.
func (f *FTPProcessor) ensureFTPPath(c *ftpConn, path string) {
	path = "/" + strings.Trim(path, "/")

	// Try direct CWD first.
	if err := c.cmd(250, "CWD %s", path); err == nil {
		return
	}

	// Build path component by component.
	_ = c.cmd(250, "CWD /")
	for _, part := range strings.Split(strings.Trim(path, "/"), "/") {
		if part == "" {
			continue
		}
		if err := c.cmd(250, "CWD %s", part); err != nil {
			// Try to create directory.
			_ = c.cmd(257, "MKD %s", part)
			_ = c.cmd(250, "CWD %s", part)
		}
	}
}

// verifyUploadAfterError checks if a file exists on the server after a 421 error.
func (f *FTPProcessor) verifyUploadAfterError(filename, targetPath, channel, dataType string) bool {
	time.Sleep(1 * time.Second)

	conn, err := f.connect()
	if err != nil {
		log.Printf("[FTP] verify connection failed: %v", err)
		return false
	}
	defer func() {
		_, _ = conn.conn.Cmd("QUIT")
		conn.conn.Close()
	}()

	// Navigate to target directory.
	for _, part := range strings.Split(strings.Trim(targetPath, "/"), "/") {
		if part != "" {
			_ = conn.cmd(250, "CWD %s", part)
		}
	}

	// List files and check for our file.
	id, err := conn.conn.Cmd("NLST")
	if err != nil {
		return false
	}
	conn.conn.StartResponse(id)
	_, _, err = conn.conn.ReadCodeLine(0)
	conn.conn.EndResponse(id)

	// A simplified check: try SIZE command.
	sizeID, err := conn.conn.Cmd("SIZE %s", filename)
	if err != nil {
		return false
	}
	conn.conn.StartResponse(sizeID)
	code, _, err := conn.conn.ReadCodeLine(0)
	conn.conn.EndResponse(sizeID)

	if code == 213 {
		// File exists on server.
		log.Printf("[FTP] file verified on server after error: %s", filename)
		return true
	}
	return false
}

// uploadFile handles upload and post-upload deletion.
func (f *FTPProcessor) uploadFile(filePath, channel, dataType string) bool {
	if f.uploadViaFTP(filePath, channel, dataType) {
		f.deleteFile(filePath, channel, dataType)
		return true
	}
	return false
}

// deleteFile removes the local file after successful upload.
func (f *FTPProcessor) deleteFile(filePath, channel, dataType string) {
	if err := os.Remove(filePath); err != nil {
		if !os.IsNotExist(err) {
			log.Printf("[FTP] failed to delete file %s: %v", filePath, err)
		}
		return
	}

	f.sentFilesLock.Lock()
	f.sentFiles[channel][dataType] = append(f.sentFiles[channel][dataType], filePath)
	// Cap the slice at 500 entries to limit memory usage.
	if len(f.sentFiles[channel][dataType]) > 500 {
		f.sentFiles[channel][dataType] = f.sentFiles[channel][dataType][len(f.sentFiles[channel][dataType])-500:]
	}
	f.sentFilesLock.Unlock()

	log.Printf("[FTP] file deleted: %s", filePath)
}

// handleUpload performs the upload if the file exists.
func (f *FTPProcessor) handleUpload(filePath, channel, dataType string) {
	if _, err := os.Stat(filePath); err == nil {
		f.uploadFile(filePath, channel, dataType)
	}
}

// safeGetMtime returns the modification time of a file, or zero time on error.
func safeGetMtime(path string) (time.Time, bool) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, false
	}
	return info.ModTime(), true
}

// processExistingFiles uploads pre-existing files in a channel/datatype directory.
func (f *FTPProcessor) processExistingFiles(channel, dataType string) {
	dirPath := f.LocalDirs[channel][dataType]
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("[FTP] failed to read directory %s: %v", dirPath, err)
		}
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".json" && ext != ".bin" {
			continue
		}

		fullPath := filepath.Join(dirPath, entry.Name())
		mtime, ok := safeGetMtime(fullPath)
		if !ok {
			continue
		}
		if time.Since(mtime).Seconds() < SafeDelay {
			continue
		}

		f.handleUpload(fullPath, channel, dataType)
	}
}

// startChannelThread runs a periodic upload loop for the given channel.
func (f *FTPProcessor) startChannelThread(channel string) {
	stopCh := make(chan struct{})
	f.mu.Lock()
	f.threads[channel] = stopCh
	f.mu.Unlock()

	f.wg.Add(1)
	go func() {
		defer f.wg.Done()
		log.Printf("[FTP] %s channel upload thread started", channel)

		// Process existing files first.
		f.processExistingFiles(channel, "waveform")
		f.processExistingFiles(channel, "event")

		interval := time.Duration(f.Sampling[channel]+60) * time.Second
		log.Printf("[FTP] %s channel periodic watch started (interval: %v)", channel, interval)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-stopCh:
				log.Printf("[FTP] %s stop signal received", channel)
				return
			case <-ticker.C:
				f.mu.RLock()
				enabled := f.enable
				f.mu.RUnlock()
				if !enabled {
					return
				}
				f.processChannelFiles(channel, "waveform")
				f.processChannelFiles(channel, "event")
			}
		}
	}()
}

// processChannelFiles scans and uploads new files for a channel/datatype.
func (f *FTPProcessor) processChannelFiles(channel, dataType string) {
	dirPath := f.LocalDirs[channel][dataType]
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("[FTP] %s %s directory read error: %v", channel, dataType, err)
		}
		return
	}

	// Build set of already-sent files.
	f.sentFilesLock.Lock()
	sentSet := make(map[string]struct{}, len(f.sentFiles[channel][dataType]))
	for _, p := range f.sentFiles[channel][dataType] {
		sentSet[p] = struct{}{}
	}
	f.sentFilesLock.Unlock()

	for _, entry := range entries {
		f.mu.RLock()
		enabled := f.enable
		f.mu.RUnlock()
		if !enabled {
			return
		}

		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".json" && ext != ".bin" {
			continue
		}

		fullPath := filepath.Join(dirPath, entry.Name())
		if _, exists := sentSet[fullPath]; exists {
			continue
		}

		mtime, ok := safeGetMtime(fullPath)
		if !ok || time.Since(mtime).Seconds() < SafeDelay {
			continue
		}

		f.handleUpload(fullPath, channel, dataType)
	}
}
