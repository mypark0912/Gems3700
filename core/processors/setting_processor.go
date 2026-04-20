package processors

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"strings"
	"unicode"

	"sv500_core/handlers"

	"github.com/redis/go-redis/v9"
)

const settingsTotalSize = 1280

// SettingsParser parses a 1280-byte SETTINGS binary blob.
type SettingsParser struct{}

// NewSettingsParser creates a new SettingsParser.
func NewSettingsParser() *SettingsParser {
	return &SettingsParser{}
}

// ParseSettings parses the full 1280-byte SETTINGS binary blob and returns a
// map whose keys correspond to the section names from the Python original:
// comm, pt, ct, etc, pqevt, transient, rcrd, extio, pqRpt, alarm, ts, ftp, sntp.
func (sp *SettingsParser) ParseSettings(data []byte) (map[string]interface{}, error) {
	if len(data) != settingsTotalSize {
		log.Printf("[SettingsParser] Data size: %d bytes (expected %d)", len(data), settingsTotalSize)
	}

	offset := 0
	result := make(map[string]interface{})

	// COMM_CFG (120 bytes)
	comm, err := sp.parseCommCfg(data[offset : offset+120])
	if err != nil {
		return nil, fmt.Errorf("parse COMM_CFG at offset %d: %w", offset, err)
	}
	result["comm"] = comm
	offset += 120

	// PT_DEF (20 bytes)
	pt, err := sp.parsePtDef(data[offset : offset+20])
	if err != nil {
		return nil, fmt.Errorf("parse PT_DEF at offset %d: %w", offset, err)
	}
	result["pt"] = pt
	offset += 20

	// CT_DEF (40 bytes)
	ct, err := sp.parseCtDef(data[offset : offset+40])
	if err != nil {
		return nil, fmt.Errorf("parse CT_DEF at offset %d: %w", offset, err)
	}
	result["ct"] = ct
	offset += 40

	// ETC_DEF (40 bytes)
	etc, err := sp.parseEtcDef(data[offset : offset+40])
	if err != nil {
		return nil, fmt.Errorf("parse ETC_DEF at offset %d: %w", offset, err)
	}
	result["etc"] = etc
	offset += 40

	// PQEVENT[5] (5 x 8 = 40 bytes)
	pqevt := make([]map[string]interface{}, 5)
	for i := 0; i < 5; i++ {
		pqevt[i] = sp.parsePQEvent(data[offset : offset+8])
		offset += 8
	}
	result["pqevt"] = pqevt

	// TRANSIENT_DEF[2] (2 x 8 = 16 bytes)
	transient := make([]map[string]interface{}, 2)
	for i := 0; i < 2; i++ {
		transient[i] = sp.parseTransientDef(data[offset : offset+8])
		offset += 8
	}
	result["transient"] = transient

	// _r1[12] skip (24 bytes)
	offset += 24

	// RECORDER_DEF[2] (2 x 8 = 16 bytes)
	rcrd := make([]map[string]interface{}, 2)
	for i := 0; i < 2; i++ {
		rcrd[i] = sp.parseRecorderDef(data[offset : offset+8])
		offset += 8
	}
	result["rcrd"] = rcrd

	// EXT_IO_DEF (84 bytes)
	extio, err := sp.parseExtIO(data[offset : offset+84])
	if err != nil {
		return nil, fmt.Errorf("parse EXT_IO_DEF at offset %d: %w", offset, err)
	}
	result["extio"] = extio
	offset += 84

	// _r3[50] skip (100 bytes)
	offset += 100

	// PQREPORT_DEF (20 bytes)
	pqRpt := sp.parsePQReportDef(data[offset : offset+20])
	result["pqRpt"] = pqRpt
	offset += 20

	// ALARM_DEF (400 bytes)
	alarm, err := sp.parseAlarm(data[offset : offset+400])
	if err != nil {
		return nil, fmt.Errorf("parse ALARM_DEF at offset %d: %w", offset, err)
	}
	result["alarm"] = alarm
	offset += 400

	// ts[4] (8 bytes)
	ts := make([]uint16, 4)
	for i := 0; i < 4; i++ {
		ts[i] = binary.LittleEndian.Uint16(data[offset : offset+2])
		offset += 2
	}
	result["ts"] = ts

	// FTP_DEF (124 bytes)
	if len(data)-offset >= 124 {
		ftp := sp.parseFtpDef(data[offset : offset+124])
		result["ftp"] = ftp
		offset += 124
	} else {
		log.Printf("[SettingsParser] Not enough data for FTP_DEF")
	}

	// SNTP_DEF (68 bytes)
	if len(data)-offset >= 68 {
		sntp := sp.parseSntpDef(data[offset : offset+68])
		result["sntp"] = sntp
		offset += 68
	}

	// _r4[80] (160 bytes) -- skip

	return result, nil
}

// ---------- individual section parsers ----------

func (sp *SettingsParser) parseCommCfg(d []byte) (map[string]interface{}, error) {
	if len(d) < 120 {
		return nil, fmt.Errorf("COMM_CFG data too short: %d", len(d))
	}
	o := 0
	r := make(map[string]interface{})

	u16 := func() uint16 {
		v := binary.LittleEndian.Uint16(d[o : o+2])
		o += 2
		return v
	}

	r["rs485Enable"] = u16()
	r["baud"] = u16()
	r["parity"] = u16()
	r["devId"] = u16()

	// ip0, sm0, gw0, dns0, sntp  (5 x 4 uint16)
	ipFields := []string{"ip0", "sm0", "gw0", "dns0", "sntp"}
	for _, name := range ipFields {
		arr := make([]uint16, 4)
		for j := 0; j < 4; j++ {
			arr[j] = u16()
		}
		r[name] = arr
	}

	// host (32-byte string)
	r["host"] = cleanString(d[o : o+32])
	o += 32

	r["tcpPort"] = u16()
	r["dhcpEn"] = u16()

	// r2 -> ftpEnable, r3 -> sntpEnable
	r["ftpEnable"] = u16()
	r["sntpEnable"] = u16()

	r["r1_0"] = u16()
	r["r1_1"] = u16()

	// daq_ip (4 uint16)
	daqIP := make([]uint16, 4)
	for j := 0; j < 4; j++ {
		daqIP[j] = u16()
	}
	r["daq_ip"] = daqIP

	r["RS485MasterMode"] = u16()
	r["daq_srate"] = u16()
	r["daq_length"] = u16()
	r["daq_interval"] = u16()
	r["daq_format"] = u16()
	r["daq_id"] = u16()

	// daq_bitpersample: uint8 + 1 byte padding ('Bx')
	r["daq_bitpersample"] = d[o]
	o += 2 // Bx = 1 byte + 1 pad

	r["pullup_485"] = u16()
	r["r0_0"] = u16()
	r["r0_1"] = u16()

	return r, nil
}

func (sp *SettingsParser) parsePtDef(d []byte) (map[string]interface{}, error) {
	if len(d) < 20 {
		return nil, fmt.Errorf("PT_DEF data too short: %d", len(d))
	}
	// format: <HHIIHHHH
	r := make(map[string]interface{})
	r["wiring"] = binary.LittleEndian.Uint16(d[0:2])
	r["freq"] = binary.LittleEndian.Uint16(d[2:4])
	r["vnorm"] = binary.LittleEndian.Uint32(d[4:8])
	r["PT1"] = binary.LittleEndian.Uint32(d[8:12])
	r["PT2"] = binary.LittleEndian.Uint16(d[12:14])
	r["r1_0"] = binary.LittleEndian.Uint16(d[14:16])
	r["r1_1"] = binary.LittleEndian.Uint16(d[16:18])
	r["r1_2"] = binary.LittleEndian.Uint16(d[18:20])
	return r, nil
}

func (sp *SettingsParser) parseCtDef(d []byte) (map[string]interface{}, error) {
	if len(d) < 40 {
		return nil, fmt.Errorf("CT_DEF data too short: %d", len(d))
	}
	// format: <HHHHHHHHHhhhHHIHHHH
	o := 0
	r := make(map[string]interface{})

	u16 := func() uint16 {
		v := binary.LittleEndian.Uint16(d[o : o+2])
		o += 2
		return v
	}
	i16 := func() int16 {
		v := int16(binary.LittleEndian.Uint16(d[o : o+2]))
		o += 2
		return v
	}

	r["_r1"] = u16()
	r["CT1"] = u16()
	r["CT2"] = u16()
	r["nTurns"] = u16()
	r["I_start"] = u16()

	// ct_dir[3] -> reconstruct to array
	ctDir := make([]uint16, 3)
	ctDir[0] = u16()
	ctDir[1] = u16()
	ctDir[2] = u16()
	r["ct_dir"] = ctDir

	r["rogowskiParam"] = u16()

	// phaseOfs[3] -> reconstruct to array (signed)
	phaseOfs := make([]int16, 3)
	phaseOfs[0] = i16()
	phaseOfs[1] = i16()
	phaseOfs[2] = i16()
	r["phaseOfs"] = phaseOfs

	r["zctType"] = u16()
	r["zctScale"] = u16()
	r["inorm"] = binary.LittleEndian.Uint32(d[o : o+4])
	o += 4
	r["r2_0"] = u16()
	r["r2_1"] = u16()
	r["r2_2"] = u16()
	r["r2_3"] = u16()

	return r, nil
}

func (sp *SettingsParser) parseEtcDef(d []byte) (map[string]interface{}, error) {
	if len(d) < 40 {
		return nil, fmt.Errorf("ETC_DEF data too short: %d", len(d))
	}
	// format: <HHHHIHHHhHHHHHHHHHH
	o := 0
	r := make(map[string]interface{})

	u16 := func() uint16 {
		v := binary.LittleEndian.Uint16(d[o : o+2])
		o += 2
		return v
	}

	r["VA_type"] = u16()
	r["PF_sign"] = u16()
	r["interval"] = u16()
	r["Iload"] = u16()
	r["P_target"] = binary.LittleEndian.Uint32(d[o : o+4])
	o += 4
	r["backlightTime"] = u16()
	r["brightness"] = u16()
	r["autorotation"] = u16()
	r["timezone"] = int16(binary.LittleEndian.Uint16(d[o : o+2]))
	o += 2
	r["maxminItv"] = u16()
	r["doEnable"] = u16()
	r["doType"] = u16()
	r["doMode"] = u16()
	r["doTimer"] = u16()
	r["doContact"] = u16()
	r["r3_0"] = u16()
	r["r3_1"] = u16()
	r["r3_2"] = u16()
	r["r3_3"] = u16()

	return r, nil
}

func (sp *SettingsParser) parsePQEvent(d []byte) map[string]interface{} {
	return map[string]interface{}{
		"level":      binary.LittleEndian.Uint16(d[0:2]),
		"nCyc":       binary.LittleEndian.Uint16(d[2:4]),
		"action":     binary.LittleEndian.Uint16(d[4:6]),
		"holdOffCyc": binary.LittleEndian.Uint16(d[6:8]),
	}
}

func (sp *SettingsParser) parseTransientDef(d []byte) map[string]interface{} {
	return map[string]interface{}{
		"holdOff":    binary.LittleEndian.Uint16(d[0:2]),
		"level":      binary.LittleEndian.Uint16(d[2:4]),
		"fastChange": binary.LittleEndian.Uint16(d[4:6]),
		"action":     binary.LittleEndian.Uint16(d[6:8]),
	}
}

func (sp *SettingsParser) parseRecorderDef(d []byte) map[string]interface{} {
	return map[string]interface{}{
		"resolution": binary.LittleEndian.Uint16(d[0:2]),
		"params":     binary.LittleEndian.Uint16(d[2:4]),
		"pre":        binary.LittleEndian.Uint16(d[4:6]),
		"post":       binary.LittleEndian.Uint16(d[6:8]),
	}
}

func (sp *SettingsParser) parseExtIO(d []byte) (map[string]interface{}, error) {
	if len(d) < 84 {
		return nil, fmt.Errorf("EXT_IO_DEF data too short: %d", len(d))
	}
	r := make(map[string]interface{})
	r["aiEnable"] = binary.LittleEndian.Uint16(d[0:2])
	r["_r"] = binary.LittleEndian.Uint16(d[2:4])

	aiscan := make([]map[string]interface{}, 4)
	off := 4
	for i := 0; i < 4; i++ {
		s := d[off : off+20]
		aiscan[i] = map[string]interface{}{
			"enable":   binary.LittleEndian.Uint16(s[0:2]),
			"devID":    binary.LittleEndian.Uint16(s[2:4]),
			"start":    binary.LittleEndian.Uint16(s[4:6]),
			"count":    binary.LittleEndian.Uint16(s[6:8]),
			"dataType": binary.LittleEndian.Uint16(s[8:10]),
			"dest":     binary.LittleEndian.Uint16(s[10:12]),
			"period":   binary.LittleEndian.Uint16(s[12:14]),
			"offset":   binary.LittleEndian.Uint16(s[14:16]),
			"scale":    math.Float32frombits(binary.LittleEndian.Uint32(s[16:20])),
		}
		off += 20
	}
	r["aiscan"] = aiscan
	return r, nil
}

func (sp *SettingsParser) parsePQReportDef(d []byte) map[string]interface{} {
	// format: <HH + H*8 = 20 bytes
	r := make(map[string]interface{})
	r["active"] = binary.LittleEndian.Uint16(d[0:2])
	r["startDay"] = binary.LittleEndian.Uint16(d[2:4])
	for i := 0; i < 8; i++ {
		r[fmt.Sprintf("r0_%d", i)] = binary.LittleEndian.Uint16(d[4+i*2 : 6+i*2])
	}
	return r
}

func (sp *SettingsParser) parseAlarm(d []byte) (map[string]interface{}, error) {
	if len(d) < 400 {
		return nil, fmt.Errorf("ALARM_DEF data too short: %d", len(d))
	}
	r := make(map[string]interface{})
	o := 0

	// header: <HHHH (8 bytes)
	r["delay"] = binary.LittleEndian.Uint16(d[o : o+2])
	r["r0"] = []uint16{
		binary.LittleEndian.Uint16(d[o+2 : o+4]),
		binary.LittleEndian.Uint16(d[o+4 : o+6]),
		binary.LittleEndian.Uint16(d[o+6 : o+8]),
	}
	o += 8

	// 32 alarm sets, each <HHHHf (12 bytes)
	sets := make([]map[string]interface{}, 32)
	for i := 0; i < 32; i++ {
		sets[i] = map[string]interface{}{
			"chan":   binary.LittleEndian.Uint16(d[o : o+2]),
			"cond":   binary.LittleEndian.Uint16(d[o+2 : o+4]),
			"dband":  binary.LittleEndian.Uint16(d[o+4 : o+6]),
			"action": binary.LittleEndian.Uint16(d[o+6 : o+8]),
			"level":  math.Float32frombits(binary.LittleEndian.Uint32(d[o+8 : o+12])),
		}
		o += 12
	}
	r["set"] = sets

	// footer: <HHHH (8 bytes)
	r["r1"] = []uint16{
		binary.LittleEndian.Uint16(d[o : o+2]),
		binary.LittleEndian.Uint16(d[o+2 : o+4]),
		binary.LittleEndian.Uint16(d[o+4 : o+6]),
		binary.LittleEndian.Uint16(d[o+6 : o+8]),
	}

	return r, nil
}

func (sp *SettingsParser) parseFtpDef(d []byte) map[string]interface{} {
	// format: <HHHHHH16s32s64s = 124 bytes
	r := make(map[string]interface{})
	host := []uint16{
		binary.LittleEndian.Uint16(d[0:2]),
		binary.LittleEndian.Uint16(d[2:4]),
		binary.LittleEndian.Uint16(d[4:6]),
		binary.LittleEndian.Uint16(d[6:8]),
	}
	r["host"] = host
	r["port"] = binary.LittleEndian.Uint16(d[8:10])
	r["enable"] = binary.LittleEndian.Uint16(d[10:12])
	r["id"] = cleanString(d[12:28])
	r["pass"] = cleanString(d[28:60])
	r["dir"] = cleanString(d[60:124])
	return r
}

func (sp *SettingsParser) parseSntpDef(d []byte) map[string]interface{} {
	// format: <HH32s32s = 68 bytes
	r := make(map[string]interface{})
	r["sntpEnable"] = binary.LittleEndian.Uint16(d[0:2])
	r["_r"] = binary.LittleEndian.Uint16(d[2:4])
	r["host"] = cleanString(d[4:36])
	r["timezone"] = cleanString(d[36:68])
	return r
}

// cleanString trims null bytes and non-printable characters from a byte slice,
// returning a clean UTF-8 string.
func cleanString(b []byte) string {
	// Find first null or trim trailing nulls/spaces
	s := string(b)

	// Trim trailing nulls and spaces
	s = strings.TrimRight(s, "\x00 ")
	if len(s) == 0 {
		return ""
	}

	// If first byte is null or non-printable control char, return empty
	if s[0] == 0 || (s[0] < 32 && s[0] != '\t' && s[0] != '\n' && s[0] != '\r') {
		return ""
	}

	// Split on null if embedded
	if idx := strings.IndexByte(s, 0); idx >= 0 {
		s = s[:idx]
	}

	// Keep only printable characters
	var sb strings.Builder
	for _, r := range s {
		if unicode.IsPrint(r) || r == '\n' || r == '\r' || r == '\t' {
			sb.WriteRune(r)
		}
	}

	return strings.TrimSpace(sb.String())
}

// SettingsRedisReader reads SETTINGS binary data from Redis and parses it.
type SettingsRedisReader struct {
	handler *handlers.RedisHandler
	parser  *SettingsParser
}

// NewSettingsRedisReader creates a new SettingsRedisReader.
func NewSettingsRedisReader(handler *handlers.RedisHandler) *SettingsRedisReader {
	return &SettingsRedisReader{
		handler: handler,
		parser:  NewSettingsParser(),
	}
}

// ReadSettings reads a 1280-byte settings blob from Redis hash "Service"
// (DB 0) at the given field key (e.g. "setup_main" or "setup_sub"),
// then parses and returns the structured result.
//
// The formatName parameter is accepted for backward compatibility but is
// not used; the parser always processes the full 1280-byte SETTINGS layout.
func (sr *SettingsRedisReader) ReadSettings(ctx context.Context, formatName string, key string) (map[string]interface{}, error) {
	client := sr.handler.Client

	// Use a pipeline to SELECT DB 0 and HGET in one round-trip.
	pipe := client.Pipeline()
	pipe.Do(ctx, "SELECT", 0)
	cmd := pipe.HGet(ctx, "Service", key)
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("redis pipeline error: %w", err)
	}

	data, err := cmd.Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("key '%s' not found in Redis hash 'Service'", key)
		}
		return nil, fmt.Errorf("redis HGET error: %w", err)
	}

	log.Printf("[SettingsRedisReader] Read %d bytes from Redis key: %s", len(data), key)

	settings, err := sr.parser.ParseSettings(data)
	if err != nil {
		return nil, fmt.Errorf("parse settings error: %w", err)
	}

	return settings, nil
}
