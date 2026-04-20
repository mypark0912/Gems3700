package data

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

// DataType constants
const (
	DataTypeEnergy      = "energy"
	DataTypeAlarm       = "alarm"
	DataTypeEvent       = "event"
	DataTypeEN50160     = "en50160"
	DataTypeEN10minData = "en10minData"
	DataTypeDemand      = "demand"
)

// BinaryDataConfig describes how to parse a binary data record.
type BinaryDataConfig struct {
	DataType            string
	Format              string
	FieldNames          []string
	Fields              []string // alias for FieldNames (used by some processors)
	ChannelMapping      map[int]string
	ScaleFactors        map[string]float64
	TimestampFieldIndex int
	ChannelFieldIndex   int
	StructSize          int
	Size                int // alias for StructSize (used by some processors)
}

// GetFieldNames returns FieldNames if set, otherwise Fields.
func (c *BinaryDataConfig) GetFieldNames() []string {
	if len(c.FieldNames) > 0 {
		return c.FieldNames
	}
	return c.Fields
}

// GetStructSize returns StructSize if set, otherwise Size.
func (c *BinaryDataConfig) GetStructSize() int {
	if c.StructSize > 0 {
		return c.StructSize
	}
	return c.Size
}

// ParsedData holds a single parsed binary record.
type ParsedData struct {
	Timestamp    time.Time
	ChannelID    int
	ChannelName  string
	DataType     string
	RawValues    map[string]interface{}
	ScaledValues map[string]interface{}
	Values       map[string]interface{} // convenience alias (populated from ScaledValues or RawValues)
}

// StandardBinaryParser parses binary blobs according to a BinaryDataConfig.
type StandardBinaryParser struct {
	Config BinaryDataConfig
}

// NewStandardBinaryParser creates a new StandardBinaryParser with the given config.
func NewStandardBinaryParser(config BinaryDataConfig) *StandardBinaryParser {
	return &StandardBinaryParser{Config: config}
}

// Parse reads one record from data using the supplied config and returns a ParsedData.
func (p *StandardBinaryParser) Parse(data []byte, configs ...*BinaryDataConfig) (*ParsedData, error) {
	var config *BinaryDataConfig
	if len(configs) > 0 && configs[0] != nil {
		config = configs[0]
	} else {
		cfg := p.Config
		config = &cfg
	}

	structSize := config.GetStructSize()
	if len(data) < structSize {
		return nil, fmt.Errorf("data too short: got %d bytes, need %d", len(data), structSize)
	}

	var parsed *ParsedData
	var err error

	switch config.DataType {
	case DataTypeEnergy:
		parsed, err = p.parseEnergy(data, config)
	case DataTypeAlarm:
		parsed, err = p.parseAlarm(data, config)
	case DataTypeEvent:
		parsed, err = p.parseEvent(data, config)
	default:
		parsed, err = p.parseGeneric(data, config)
	}

	if err == nil && parsed != nil && parsed.Values == nil {
		// Populate Values from ScaledValues (or RawValues if no scaling).
		if len(parsed.ScaledValues) > 0 {
			parsed.Values = parsed.ScaledValues
		} else {
			parsed.Values = parsed.RawValues
		}
	}

	return parsed, err
}

// parseEnergy: int32 (channel_id), uint32 (timestamp), 6x float32
func (p *StandardBinaryParser) parseEnergy(data []byte, config *BinaryDataConfig) (*ParsedData, error) {
	r := bytes.NewReader(data)

	var channelID int32
	if err := binary.Read(r, binary.LittleEndian, &channelID); err != nil {
		return nil, fmt.Errorf("read channel_id: %w", err)
	}

	var ts uint32
	if err := binary.Read(r, binary.LittleEndian, &ts); err != nil {
		return nil, fmt.Errorf("read timestamp: %w", err)
	}

	floats := make([]float32, 6)
	for i := range floats {
		if err := binary.Read(r, binary.LittleEndian, &floats[i]); err != nil {
			return nil, fmt.Errorf("read float[%d]: %w", i, err)
		}
	}

	raw := make(map[string]interface{})
	raw["channel_id"] = channelID
	raw["timestamp"] = ts
	for i, f := range floats {
		name := fmt.Sprintf("value_%d", i)
		if i+2 < len(config.FieldNames) {
			name = config.FieldNames[i+2]
		}
		raw[name] = f
	}

	scaled := applyScaleFactors(raw, config.ScaleFactors)

	channelName := ""
	if name, ok := config.ChannelMapping[int(channelID)]; ok {
		channelName = name
	}

	return &ParsedData{
		Timestamp:    time.Unix(int64(ts), 0),
		ChannelID:    int(channelID),
		ChannelName:  channelName,
		DataType:     DataTypeEnergy,
		RawValues:    raw,
		ScaledValues: scaled,
	}, nil
}

// parseAlarm: 5x uint16, 1x uint32, 1x float32, 2x uint16, 1x float32 (HHHHHIfHHf)
func (p *StandardBinaryParser) parseAlarm(data []byte, config *BinaryDataConfig) (*ParsedData, error) {
	r := bytes.NewReader(data)

	fields := []interface{}{
		new(uint16), new(uint16), new(uint16), new(uint16), new(uint16),
		new(uint32),
		new(float32),
		new(uint16), new(uint16),
		new(float32),
	}

	for i, f := range fields {
		if err := binary.Read(r, binary.LittleEndian, f); err != nil {
			return nil, fmt.Errorf("read field[%d]: %w", i, err)
		}
	}

	raw := make(map[string]interface{})
	for i, f := range fields {
		name := fmt.Sprintf("field_%d", i)
		if i < len(config.FieldNames) {
			name = config.FieldNames[i]
		}
		switch v := f.(type) {
		case *uint16:
			raw[name] = *v
		case *uint32:
			raw[name] = *v
		case *float32:
			raw[name] = *v
		}
	}

	scaled := applyScaleFactors(raw, config.ScaleFactors)

	channelID := 0
	if config.ChannelFieldIndex >= 0 && config.ChannelFieldIndex < len(fields) {
		if v, ok := fields[config.ChannelFieldIndex].(*uint16); ok {
			channelID = int(*v)
		}
	}

	var ts time.Time
	if config.TimestampFieldIndex >= 0 && config.TimestampFieldIndex < len(fields) {
		if v, ok := fields[config.TimestampFieldIndex].(*uint32); ok {
			ts = time.Unix(int64(*v), 0)
		}
	}

	channelName := ""
	if name, ok := config.ChannelMapping[channelID]; ok {
		channelName = name
	}

	return &ParsedData{
		Timestamp:    ts,
		ChannelID:    channelID,
		ChannelName:  channelName,
		DataType:     DataTypeAlarm,
		RawValues:    raw,
		ScaledValues: scaled,
	}, nil
}

// parseEvent: int32, uint32, 4x uint16, 3x float32 (iIHHHHfff)
func (p *StandardBinaryParser) parseEvent(data []byte, config *BinaryDataConfig) (*ParsedData, error) {
	r := bytes.NewReader(data)

	var channelID int32
	if err := binary.Read(r, binary.LittleEndian, &channelID); err != nil {
		return nil, fmt.Errorf("read channel_id: %w", err)
	}

	var ts uint32
	if err := binary.Read(r, binary.LittleEndian, &ts); err != nil {
		return nil, fmt.Errorf("read timestamp: %w", err)
	}

	u16s := make([]uint16, 4)
	for i := range u16s {
		if err := binary.Read(r, binary.LittleEndian, &u16s[i]); err != nil {
			return nil, fmt.Errorf("read uint16[%d]: %w", i, err)
		}
	}

	f32s := make([]float32, 3)
	for i := range f32s {
		if err := binary.Read(r, binary.LittleEndian, &f32s[i]); err != nil {
			return nil, fmt.Errorf("read float32[%d]: %w", i, err)
		}
	}

	raw := make(map[string]interface{})
	raw["channel_id"] = channelID
	raw["timestamp"] = ts

	fieldIdx := 2
	for i, v := range u16s {
		name := fmt.Sprintf("u16_%d", i)
		if fieldIdx < len(config.FieldNames) {
			name = config.FieldNames[fieldIdx]
		}
		raw[name] = v
		fieldIdx++
	}
	for i, v := range f32s {
		name := fmt.Sprintf("f32_%d", i)
		if fieldIdx < len(config.FieldNames) {
			name = config.FieldNames[fieldIdx]
		}
		raw[name] = v
		fieldIdx++
	}

	scaled := applyScaleFactors(raw, config.ScaleFactors)

	channelName := ""
	if name, ok := config.ChannelMapping[int(channelID)]; ok {
		channelName = name
	}

	return &ParsedData{
		Timestamp:    time.Unix(int64(ts), 0),
		ChannelID:    int(channelID),
		ChannelName:  channelName,
		DataType:     DataTypeEvent,
		RawValues:    raw,
		ScaledValues: scaled,
	}, nil
}

// parseGeneric reads values according to the field count as float32.
func (p *StandardBinaryParser) parseGeneric(data []byte, config *BinaryDataConfig) (*ParsedData, error) {
	r := bytes.NewReader(data)

	raw := make(map[string]interface{})
	for i, name := range config.FieldNames {
		var val float32
		if err := binary.Read(r, binary.LittleEndian, &val); err != nil {
			return nil, fmt.Errorf("read field[%d] %s: %w", i, name, err)
		}
		raw[name] = val
	}

	scaled := applyScaleFactors(raw, config.ScaleFactors)

	return &ParsedData{
		DataType:     config.DataType,
		RawValues:    raw,
		ScaledValues: scaled,
	}, nil
}

func applyScaleFactors(raw map[string]interface{}, factors map[string]float64) map[string]interface{} {
	scaled := make(map[string]interface{}, len(raw))
	for k, v := range raw {
		if factor, ok := factors[k]; ok {
			switch val := v.(type) {
			case float32:
				scaled[k] = float64(val) * factor
			case float64:
				scaled[k] = val * factor
			default:
				scaled[k] = v
			}
		} else {
			scaled[k] = v
		}
	}
	return scaled
}

// CreateEnergyConfig returns a BinaryDataConfig for energy data.
func CreateEnergyConfig() *BinaryDataConfig {
	return &BinaryDataConfig{
		DataType:            DataTypeEnergy,
		FieldNames:          []string{"channel_id", "timestamp", "kwh_import", "kwh_export", "kvarh_import", "kvarh_export", "kvah_import", "kvah_export"},
		ChannelMapping:      map[int]string{0: "Main", 1: "Sub"},
		ScaleFactors:        make(map[string]float64),
		TimestampFieldIndex: 1,
		ChannelFieldIndex:   0,
		StructSize:          32, // 4 + 4 + 6*4
	}
}

// CreateAlarmConfig returns a BinaryDataConfig for alarm data.
func CreateAlarmConfig() *BinaryDataConfig {
	return &BinaryDataConfig{
		DataType:            DataTypeAlarm,
		FieldNames:          []string{"alarm_type", "alarm_code", "channel_id", "phase", "status", "timestamp", "value", "threshold_high", "threshold_low", "duration"},
		ChannelMapping:      make(map[int]string),
		ScaleFactors:        make(map[string]float64),
		TimestampFieldIndex: 5,
		ChannelFieldIndex:   2,
		StructSize:          24, // 5*2 + 4 + 4 + 2*2 + 4
	}
}

// CreateEventConfig returns a BinaryDataConfig for event data.
func CreateEventConfig() *BinaryDataConfig {
	return &BinaryDataConfig{
		DataType:            DataTypeEvent,
		FieldNames:          []string{"channel_id", "timestamp", "event_type", "event_code", "phase", "status", "value", "pre_value", "post_value"},
		ChannelMapping:      make(map[int]string),
		ScaleFactors:        make(map[string]float64),
		TimestampFieldIndex: 1,
		ChannelFieldIndex:   0,
		StructSize:          28, // 4 + 4 + 4*2 + 3*4
	}
}
