package handlers

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	influxapi "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

// AES key for token decryption
var aesKey = []byte("ntekSystem_20250721_mypark_caner")

const influxConfigPath = "/home/root/config/influx.json"

// ---------------------------------------------------------------------------
// AES Decryption
// ---------------------------------------------------------------------------

func decryptToken(encryptedB64 string) (string, error) {
	hash := sha256.Sum256(aesKey)
	iv := hash[:aes.BlockSize]

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedB64)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %w", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", fmt.Errorf("AES cipher error: %w", err)
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", fmt.Errorf("ciphertext length %d is not a multiple of block size", len(ciphertext))
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	plaintext, err = pkcs7Unpad(plaintext, aes.BlockSize)
	if err != nil {
		return "", fmt.Errorf("PKCS7 unpad error: %w", err)
	}

	return string(plaintext), nil
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}
	padLen := int(data[len(data)-1])
	if padLen == 0 || padLen > blockSize || padLen > len(data) {
		return nil, fmt.Errorf("invalid padding length %d", padLen)
	}
	for i := len(data) - padLen; i < len(data); i++ {
		if data[i] != byte(padLen) {
			return nil, fmt.Errorf("invalid padding byte")
		}
	}
	return data[:len(data)-padLen], nil
}

// ---------------------------------------------------------------------------
// Config
// ---------------------------------------------------------------------------

type influxConfig struct {
	URL   string `json:"url"`
	Token string `json:"token"`
	Org   string `json:"org"`
}

func loadInfluxConfig() (url, token, org string, err error) {
	data, err := os.ReadFile(influxConfigPath)
	if err != nil {
		return "", "", "", fmt.Errorf("read influx config: %w", err)
	}

	var cfg influxConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", "", "", fmt.Errorf("parse influx config: %w", err)
	}

	if cfg.URL == "" {
		cfg.URL = "http://127.0.0.1:8086"
	}

	plainToken, err := decryptToken(cfg.Token)
	if err != nil {
		return "", "", "", fmt.Errorf("decrypt token: %w", err)
	}

	return cfg.URL, plainToken, cfg.Org, nil
}

// ---------------------------------------------------------------------------
// InfluxConnectionPool
// ---------------------------------------------------------------------------

// PoolEntry wraps a single InfluxDB client connection.
type PoolEntry struct {
	Client    influxdb2.Client
	WriteAPI  influxapi.WriteAPIBlocking
	QueryAPI  influxapi.QueryAPI
	CreatedAt time.Time
}

type InfluxConnectionPool struct {
	mu           sync.Mutex
	pool         chan *PoolEntry
	size         int
	connLifetime time.Duration
	url          string
	token        string
	Org          string
}

func NewInfluxConnectionPool(size int) (*InfluxConnectionPool, error) {
	url, token, org, err := loadInfluxConfig()
	if err != nil {
		return nil, err
	}

	p := &InfluxConnectionPool{
		pool:         make(chan *PoolEntry, size),
		size:         size,
		connLifetime: 3600 * time.Second,
		url:          url,
		token:        token,
		Org:          org,
	}

	successCount := 0
	for i := 0; i < size; i++ {
		entry, err := p.createEntry()
		if err != nil {
			log.Printf("InfluxDB connection %d/%d failed: %v", i+1, size, err)
			continue
		}
		p.pool <- entry
		successCount++
		log.Printf("InfluxDB connection %d/%d created", i+1, size)
	}

	if successCount == 0 {
		return nil, fmt.Errorf("InfluxDB connection pool init failed: no connections")
	}

	log.Printf("InfluxDB connection pool initialized: %d/%d connections", successCount, size)
	return p, nil
}

func (p *InfluxConnectionPool) createEntry() (*PoolEntry, error) {
	opts := influxdb2.DefaultOptions().
		SetHTTPRequestTimeout(30)
	client := influxdb2.NewClientWithOptions(p.url, p.token, opts)

	// Health check
	health, err := client.Health(context.Background())
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("health check failed: %w", err)
	}
	if health.Status != "pass" {
		client.Close()
		return nil, fmt.Errorf("health check status: %s", health.Status)
	}

	return &PoolEntry{
		Client:    client,
		WriteAPI:  client.WriteAPIBlocking(p.Org, ""),
		QueryAPI:  client.QueryAPI(p.Org),
		CreatedAt: time.Now(),
	}, nil
}

func (p *InfluxConnectionPool) GetConnection(timeout time.Duration) (*PoolEntry, error) {
	select {
	case entry := <-p.pool:
		if time.Since(entry.CreatedAt) > p.connLifetime {
			log.Println("Recycling old InfluxDB connection")
			entry.Client.Close()
			newEntry, err := p.createEntry()
			if err != nil {
				return nil, fmt.Errorf("connection recreation failed: %w", err)
			}
			return newEntry, nil
		}
		return entry, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("InfluxDB connection pool timeout after %v", timeout)
	}
}

func (p *InfluxConnectionPool) ReturnConnection(entry *PoolEntry) {
	select {
	case p.pool <- entry:
	default:
		entry.Client.Close()
	}
}

func (p *InfluxConnectionPool) WithConnection(fn func(entry *PoolEntry) error) error {
	const maxRetries = 3
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		entry, err := p.GetConnection(30 * time.Second)
		if err != nil {
			return err
		}

		err = fn(entry)
		if err == nil {
			p.ReturnConnection(entry)
			return nil
		}

		lastErr = err
		if isRetriableError(err) && attempt < maxRetries-1 {
			p.ReturnConnection(entry)
			waitTime := time.Duration(attempt+1) * 2 * time.Second
			log.Printf("InfluxDB retriable error, retrying in %v (attempt %d/%d): %v",
				waitTime, attempt+1, maxRetries, err)
			time.Sleep(waitTime)
			continue
		}

		// Non-retriable: replace connection
		entry.Client.Close()
		newEntry, createErr := p.createEntry()
		if createErr != nil {
			log.Printf("Connection recreation failed: %v", createErr)
		} else {
			p.ReturnConnection(newEntry)
		}
		return err
	}

	return lastErr
}

func isRetriableError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "500") ||
		strings.Contains(msg, "Internal Server Error") ||
		err == io.EOF ||
		strings.Contains(msg, "EOF")
}

func (p *InfluxConnectionPool) GetPoolStatus() map[string]interface{} {
	return map[string]interface{}{
		"pool_size":           p.size,
		"available":           len(p.pool),
		"url":                 p.url,
		"org":                 p.Org,
		"connection_lifetime": p.connLifetime.Seconds(),
	}
}

func (p *InfluxConnectionPool) CloseAllConnections() {
	closed := 0
	for {
		select {
		case entry := <-p.pool:
			entry.Client.Close()
			closed++
		default:
			log.Printf("Closed %d InfluxDB connections", closed)
			return
		}
	}
}

// ---------------------------------------------------------------------------
// Singleton pool
// ---------------------------------------------------------------------------

var (
	influxPoolOnce sync.Once
	influxPoolInst *InfluxConnectionPool
	influxPoolErr  error
)

func GetInfluxPool() (*InfluxConnectionPool, error) {
	influxPoolOnce.Do(func() {
		influxPoolInst, influxPoolErr = NewInfluxConnectionPool(12)
	})
	return influxPoolInst, influxPoolErr
}

// ---------------------------------------------------------------------------
// InfluxDBHandler
// ---------------------------------------------------------------------------

type InfluxDBHandler struct {
	Bucket     string
	pool       *InfluxConnectionPool
	writeCount int64
	lastLogAt  time.Time
	mu         sync.Mutex
}

func NewInfluxDBHandler(bucket string) *InfluxDBHandler {
	return &InfluxDBHandler{
		Bucket:    bucket,
		lastLogAt: time.Now(),
	}
}

func (h *InfluxDBHandler) getPool() (*InfluxConnectionPool, error) {
	if h.pool != nil {
		return h.pool, nil
	}
	pool, err := GetInfluxPool()
	if err != nil {
		return nil, err
	}
	h.pool = pool
	return pool, nil
}

func (h *InfluxDBHandler) logWriteStats() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.writeCount++
	if time.Since(h.lastLogAt) > 5*time.Minute {
		log.Printf("InfluxDB write stats: %d writes", h.writeCount)
		h.lastLogAt = time.Now()
	}
}

func (h *InfluxDBHandler) writePoint(bucket string, point *write.Point) error {
	pool, err := h.getPool()
	if err != nil {
		return err
	}

	err = pool.WithConnection(func(entry *PoolEntry) error {
		writeAPI := entry.Client.WriteAPIBlocking(pool.Org, bucket)
		return writeAPI.WritePoint(context.Background(), point)
	})

	if err == nil {
		h.logWriteStats()
	}
	return err
}

func (h *InfluxDBHandler) writePoints(bucket string, points []*write.Point) error {
	pool, err := h.getPool()
	if err != nil {
		return err
	}

	return pool.WithConnection(func(entry *PoolEntry) error {
		writeAPI := entry.Client.WriteAPIBlocking(pool.Org, bucket)
		for _, pt := range points {
			if err := writeAPI.WritePoint(context.Background(), pt); err != nil {
				return err
			}
		}
		return nil
	})
}

// Write writes a data point to the handler's default bucket.
func (h *InfluxDBHandler) Write(measurement string, tags map[string]string, fields map[string]interface{}, ts time.Time) error {
	return h.WriteDataPoint(h.Bucket, measurement, tags, fields, ts)
}

// WriteDataPoint writes a generic data point.
func (h *InfluxDBHandler) WriteDataPoint(bucket, measurement string, tags map[string]string, fields map[string]interface{}, ts time.Time) error {
	point := influxdb2.NewPoint(measurement, tags, fields, ts)
	return h.writePoint(bucket, point)
}

// WriteCumulativeData writes cumulative energy data.
func (h *InfluxDBHandler) WriteCumulativeData(data map[string]interface{}, measurement string) error {
	tags := map[string]string{
		"channel": fmt.Sprintf("%v", data["channel_name"]),
	}
	if mid, ok := data["meter_id"]; ok {
		tags["meter_id"] = fmt.Sprintf("%v", mid)
	}

	fields := make(map[string]interface{})
	for _, f := range []string{"kwh_import", "kwh_export", "kvarh_import", "kvarh_export", "kvah_import", "kvah_export"} {
		if v, ok := data[f]; ok {
			fields[f] = v
		}
	}

	var ts time.Time
	if t, ok := data["timestamp"].(time.Time); ok {
		ts = t
	} else {
		ts = time.Now()
	}

	point := influxdb2.NewPoint(measurement, tags, fields, ts)
	return h.writePoint(h.Bucket, point)
}

// WriteConsumptionData writes energy consumption data.
func (h *InfluxDBHandler) WriteConsumptionData(data map[string]interface{}, consumption map[string]float64, measurement string) error {
	if len(consumption) == 0 {
		return nil
	}

	tags := map[string]string{
		"channel": fmt.Sprintf("%v", data["channel_name"]),
	}
	if mid, ok := data["meter_id"]; ok {
		tags["meter_id"] = fmt.Sprintf("%v", mid)
	}
	if period, ok := data["period"]; ok {
		tags["period"] = fmt.Sprintf("%v", period)
	}

	fields := make(map[string]interface{})
	for k, v := range consumption {
		fields[k] = v
	}

	var ts time.Time
	if t, ok := data["timestamp"].(time.Time); ok {
		ts = t
	} else {
		ts = time.Now()
	}

	point := influxdb2.NewPoint(measurement, tags, fields, ts)
	return h.writePoint(h.Bucket, point)
}

// WriteEventData writes event data.
func (h *InfluxDBHandler) WriteEventData(eventData map[string]interface{}, measurement string) error {
	tags := map[string]string{
		"channel":    fmt.Sprintf("%v", eventData["channel_name"]),
		"event_type": fmt.Sprintf("%v", eventData["event_type_text"]),
		"event_id":   fmt.Sprintf("%v", eventData["event_id"]),
	}

	fields := map[string]interface{}{
		"duration":     eventData["duration"],
		"mask":         eventData["mask"],
		"level_l1":     eventData["level_l1"],
		"level_l2":     eventData["level_l2"],
		"level_l3":     eventData["level_l3"],
		"milliseconds": eventData["milliseconds"],
	}

	var ts time.Time
	if t, ok := eventData["timestamp"].(time.Time); ok {
		ts = t
	} else {
		ts = time.Now()
	}

	point := influxdb2.NewPoint(measurement, tags, fields, ts)
	return h.writePoint(h.Bucket, point)
}

// WriteAlarmData writes alarm data.
func (h *InfluxDBHandler) WriteAlarmData(alarmData map[string]interface{}, measurement string) error {
	tags := map[string]string{
		"channel": fmt.Sprintf("%v", alarmData["channel_name"]),
		"chan":    fmt.Sprintf("%v", alarmData["chan"]),
	}

	fields := map[string]interface{}{
		"status":          alarmData["status_text"],
		"alarm_config_id": alarmData["chid"],
		"chan_text":        alarmData["chan_text"],
		"condition":        alarmData["condition_text"],
		"level":            alarmData["level"],
		"value":            alarmData["value"],
		"status_code":      alarmData["status"],
		"condition_code":   alarmData["condition"],
		"count":            alarmData["count"],
	}

	var ts time.Time
	if t, ok := alarmData["timestamp"].(time.Time); ok {
		ts = t
	} else {
		ts = time.Now()
	}

	point := influxdb2.NewPoint(measurement, tags, fields, ts)
	return h.writePoint(h.Bucket, point)
}

// WriteEN50160Report writes EN50160 report data.
func (h *InfluxDBHandler) WriteEN50160Report(reportData map[string]interface{}, measurement string) error {
	tags := map[string]string{
		"channel":    fmt.Sprintf("%v", reportData["channel_name"]),
		"channel_id": fmt.Sprintf("%v", reportData["channel_id"]),
	}

	fields := make(map[string]interface{})
	for k, v := range reportData {
		if k == "channel_name" || k == "channel_id" || k == "timestamp" {
			continue
		}
		fields[k] = v
	}

	var ts time.Time
	if t, ok := reportData["timestamp"].(time.Time); ok {
		ts = t
	} else {
		ts = time.Now()
	}

	point := influxdb2.NewPoint(measurement, tags, fields, ts)
	return h.writePoint(h.Bucket, point)
}

// ExecuteQuery runs a Flux query and returns results.
func (h *InfluxDBHandler) ExecuteQuery(query string) ([]map[string]interface{}, error) {
	pool, err := h.getPool()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	err = pool.WithConnection(func(entry *PoolEntry) error {
		result, err := entry.QueryAPI.Query(context.Background(), query)
		if err != nil {
			return fmt.Errorf("query error: %w", err)
		}

		for result.Next() {
			record := result.Record()
			results = append(results, record.Values())
		}

		return result.Err()
	})

	return results, err
}

// QueryLastData queries the most recent data for a channel.
func (h *InfluxDBHandler) QueryLastData(channel, measurement string, lookbackHours int) (map[string]interface{}, error) {
	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -%dh)
		|> filter(fn: (r) => r["_measurement"] == "%s")
		|> filter(fn: (r) => r["channel"] == "%s")
		|> last()
		|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")`,
		h.Bucket, lookbackHours, measurement, channel)

	results, err := h.ExecuteQuery(query)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return results[0], nil
}

// QueryPeriodSum queries sum for a period.
func (h *InfluxDBHandler) QueryPeriodSum(channel string, startTime, endTime time.Time, measurement string) (map[string]float64, error) {
	startUTC := startTime.UTC().Format(time.RFC3339)
	endUTC := endTime.UTC().Format(time.RFC3339)

	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r["_measurement"] == "%s")
		|> filter(fn: (r) => r["channel"] == "%s")
		|> filter(fn: (r) => r["_field"] =~ /.*_consumption$/)
		|> group(columns: ["_field"])
		|> sum()`,
		h.Bucket, startUTC, endUTC, measurement, channel)

	results, err := h.ExecuteQuery(query)
	if err != nil {
		return nil, err
	}

	result := make(map[string]float64)
	for _, r := range results {
		if field, ok := r["_field"].(string); ok {
			if value, ok := r["_value"].(float64); ok {
				result[field] = value
			}
		}
	}
	return result, nil
}

// WriteBatchPoints writes multiple points at once.
func (h *InfluxDBHandler) WriteBatchPoints(bucket string, points []*write.Point) error {
	return h.writePoints(bucket, points)
}
