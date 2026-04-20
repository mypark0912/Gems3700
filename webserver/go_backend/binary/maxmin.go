package binary

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// MaxMinItemSize is the size of MAXMIN_DATA_t { float max, uint32 max_ts, float min, uint32 min_ts }.
const MaxMinItemSize = 16

// Python layout: MAXMIN_{main,sub} HASH, field "1sec" / "15min" → 21 × MAXMIN_DATA_t = 336 bytes.
const MaxMinBlobCount = 21
const MaxMinBlobSize = MaxMinBlobCount * MaxMinItemSize // 336

// MaxMinItem mirrors Python's parsed MAXMIN_DATA dict.
type MaxMinItem struct {
	Max          float64 `json:"max"`
	MaxTimestamp string  `json:"max_timestamp"`
	Min          float64 `json:"min"`
	MinTimestamp string  `json:"min_timestamp"`
}

var EmptyMaxMinItem = MaxMinItem{}

// MaxMin1SecData mirrors Python MaxMin1SecParser.field_mapping order.
type MaxMin1SecData struct {
	Freq MaxMinItem     `json:"freq"`
	Temp MaxMinItem     `json:"temp"`
	U    [4]MaxMinItem  `json:"u"`
	Upp  [4]MaxMinItem  `json:"upp"`
	I    [4]MaxMinItem  `json:"i"`
	Itot MaxMinItem     `json:"itot"`
	In   MaxMinItem     `json:"in"`
	Isum MaxMinItem     `json:"isum"`
	PF   [4]MaxMinItem  `json:"pf"`
}

// MaxMin15MinData mirrors Python MaxMin15MinParser.field_mapping order.
type MaxMin15MinData struct {
	P      [4]MaxMinItem `json:"p"`
	Q      [4]MaxMinItem `json:"q"`
	S      [4]MaxMinItem `json:"s"`
	ThdU   [3]MaxMinItem `json:"thd_u"`
	ThdUpp [3]MaxMinItem `json:"thd_upp"`
	ThdI   [3]MaxMinItem `json:"thd_i"`
}

// ParseMaxMinItem parses a single 16-byte MAXMIN_DATA_t.
func ParseMaxMinItem(data []byte) (MaxMinItem, error) {
	if len(data) < MaxMinItemSize {
		return EmptyMaxMinItem, fmt.Errorf("maxmin item: expected %d bytes, got %d", MaxMinItemSize, len(data))
	}
	maxVal := math.Float32frombits(binary.LittleEndian.Uint32(data[0:4]))
	maxTS := binary.LittleEndian.Uint32(data[4:8])
	minVal := math.Float32frombits(binary.LittleEndian.Uint32(data[8:12]))
	minTS := binary.LittleEndian.Uint32(data[12:16])

	item := MaxMinItem{Max: float64(maxVal), Min: float64(minVal)}
	if maxTS > 0 {
		item.MaxTimestamp = time.Unix(int64(maxTS), 0).Format("2006-01-02 15:04:05")
	}
	if minTS > 0 {
		item.MinTimestamp = time.Unix(int64(minTS), 0).Format("2006-01-02 15:04:05")
	}
	return item, nil
}

func readMaxMinArray(data []byte, offset, count int, out []MaxMinItem) int {
	for i := 0; i < count; i++ {
		out[i], _ = ParseMaxMinItem(data[offset : offset+MaxMinItemSize])
		offset += MaxMinItemSize
	}
	return offset
}

// ParseMaxMin1Sec parses the 336-byte 1-second MaxMin blob.
func ParseMaxMin1Sec(data []byte) (*MaxMin1SecData, error) {
	if len(data) < MaxMinBlobSize {
		return nil, fmt.Errorf("maxmin 1sec: expected %d bytes, got %d", MaxMinBlobSize, len(data))
	}
	d := &MaxMin1SecData{}
	off := 0
	d.Freq, _ = ParseMaxMinItem(data[off : off+MaxMinItemSize])
	off += MaxMinItemSize
	d.Temp, _ = ParseMaxMinItem(data[off : off+MaxMinItemSize])
	off += MaxMinItemSize
	off = readMaxMinArray(data, off, 4, d.U[:])
	off = readMaxMinArray(data, off, 4, d.Upp[:])
	off = readMaxMinArray(data, off, 4, d.I[:])
	d.Itot, _ = ParseMaxMinItem(data[off : off+MaxMinItemSize])
	off += MaxMinItemSize
	d.In, _ = ParseMaxMinItem(data[off : off+MaxMinItemSize])
	off += MaxMinItemSize
	d.Isum, _ = ParseMaxMinItem(data[off : off+MaxMinItemSize])
	off += MaxMinItemSize
	readMaxMinArray(data, off, 4, d.PF[:])
	return d, nil
}

// ParseMaxMin15Min parses the 336-byte 15-minute MaxMin blob.
func ParseMaxMin15Min(data []byte) (*MaxMin15MinData, error) {
	if len(data) < MaxMinBlobSize {
		return nil, fmt.Errorf("maxmin 15min: expected %d bytes, got %d", MaxMinBlobSize, len(data))
	}
	d := &MaxMin15MinData{}
	off := 0
	off = readMaxMinArray(data, off, 4, d.P[:])
	off = readMaxMinArray(data, off, 4, d.Q[:])
	off = readMaxMinArray(data, off, 4, d.S[:])
	off = readMaxMinArray(data, off, 3, d.ThdU[:])
	off = readMaxMinArray(data, off, 3, d.ThdUpp[:])
	readMaxMinArray(data, off, 3, d.ThdI[:])
	return d, nil
}

// maxminRedisKey builds "MAXMIN_main" / "MAXMIN_sub" from channel name.
func maxminRedisKey(channel string) string {
	ch := strings.ToLower(channel)
	if ch != "main" && ch != "sub" {
		ch = "main"
	}
	return "MAXMIN_" + ch
}

// FetchMaxMin1Sec fetches HGET MAXMIN_<channel> "1sec".
func FetchMaxMin1Sec(ctx context.Context, client *redis.Client, channel string) (*MaxMin1SecData, error) {
	raw, err := client.HGet(ctx, maxminRedisKey(channel), "1sec").Bytes()
	if err != nil {
		return nil, err
	}
	return ParseMaxMin1Sec(raw)
}

// FetchMaxMin15Min fetches HGET MAXMIN_<channel> "15min".
func FetchMaxMin15Min(ctx context.Context, client *redis.Client, channel string) (*MaxMin15MinData, error) {
	raw, err := client.HGet(ctx, maxminRedisKey(channel), "15min").Bytes()
	if err != nil {
		return nil, err
	}
	return ParseMaxMin15Min(raw)
}

// Flat1Sec returns a flat "U1_max"/"U1_maxTime"/... map matching RedisMapDetail legacy keys.
func Flat1Sec(d *MaxMin1SecData) map[string]interface{} {
	m := map[string]interface{}{}
	if d == nil {
		return m
	}
	putItem(m, "Freq", d.Freq)
	putItem(m, "Temp", d.Temp)
	for i, v := range d.U {
		putItem(m, fmt.Sprintf("U%d", i+1), v)
	}
	for i, v := range d.Upp {
		putItem(m, fmt.Sprintf("Upp%d", i+1), v)
	}
	for i, v := range d.I {
		putItem(m, fmt.Sprintf("I%d", i+1), v)
	}
	putItem(m, "Itot", d.Itot)
	putItem(m, "In", d.In)
	putItem(m, "Isum", d.Isum)
	for i, v := range d.PF {
		putItem(m, fmt.Sprintf("PF%d", i+1), v)
	}
	return m
}

// Flat15Min returns a flat "P1_max"/... map for 15-min aggregates.
func Flat15Min(d *MaxMin15MinData) map[string]interface{} {
	m := map[string]interface{}{}
	if d == nil {
		return m
	}
	for i, v := range d.P {
		putItem(m, fmt.Sprintf("P%d", i+1), v)
	}
	for i, v := range d.Q {
		putItem(m, fmt.Sprintf("Q%d", i+1), v)
	}
	for i, v := range d.S {
		putItem(m, fmt.Sprintf("S%d", i+1), v)
	}
	for i, v := range d.ThdU {
		putItem(m, fmt.Sprintf("THD_U%d", i+1), v)
	}
	for i, v := range d.ThdUpp {
		putItem(m, fmt.Sprintf("THD_Upp%d", i+1), v)
	}
	for i, v := range d.ThdI {
		putItem(m, fmt.Sprintf("THD_I%d", i+1), v)
	}
	return m
}

func putItem(m map[string]interface{}, prefix string, v MaxMinItem) {
	m[prefix+"_max"] = v.Max
	m[prefix+"_maxTime"] = v.MaxTimestamp
	m[prefix+"_min"] = v.Min
	m[prefix+"_minTime"] = v.MinTimestamp
}
