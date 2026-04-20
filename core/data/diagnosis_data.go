package data

import (
	"reflect"
	"strings"
)

// Equipment type constants
const (
	EquipmentNone               = 0
	EquipmentCompressor         = 1
	EquipmentFan                = 2
	EquipmentPump               = 3
	EquipmentMotor              = 4
	EquipmentPowerSupply        = 5
	EquipmentTransformer        = 6
	EquipmentMotorFeed          = 7
	EquipmentPrimaryTransformer = 8
)

// DiagnosisData is the interface for all equipment diagnosis data types.
type DiagnosisData interface {
	GetAllValues() map[string]int
}

// --- Equipment diagnosis structs ---

// CompressorFan holds diagnosis data for compressor and fan equipment.
type CompressorFan struct {
	Turbulence         int
	MechanicalUnbalance int
	SoftFoot           int
	TorqueRipple       int
	Load               int
	CableConnection    int
	NoiseVibration     int
	Heat               int
	IncomingVoltage    int
	Rotor              int
	Stator             int
	Bearing            int
	DCLink             int
}

func NewCompressorFan() *CompressorFan {
	return &CompressorFan{
		Turbulence: -1, MechanicalUnbalance: -1, SoftFoot: -1, TorqueRipple: -1,
		Load: -1, CableConnection: -1, NoiseVibration: -1, Heat: -1,
		IncomingVoltage: -1, Rotor: -1, Stator: -1, Bearing: -1, DCLink: -1,
	}
}

func (c *CompressorFan) GetAllValues() map[string]int {
	return structToMap(c)
}

// Pump holds diagnosis data for pump equipment.
type Pump struct {
	Cavitation          int
	MechanicalUnbalance int
	SoftFoot            int
	TorqueRipple        int
	Load                int
	CableConnection     int
	NoiseVibration      int
	Heat                int
	IncomingVoltage     int
	Rotor               int
	Stator              int
	Bearing             int
	DCLink              int
}

func NewPump() *Pump {
	return &Pump{
		Cavitation: -1, MechanicalUnbalance: -1, SoftFoot: -1, TorqueRipple: -1,
		Load: -1, CableConnection: -1, NoiseVibration: -1, Heat: -1,
		IncomingVoltage: -1, Rotor: -1, Stator: -1, Bearing: -1, DCLink: -1,
	}
}

func (p *Pump) GetAllValues() map[string]int {
	return structToMap(p)
}

// Motor holds diagnosis data for motor equipment.
type Motor struct {
	TorqueRipple        int
	MechanicalUnbalance int
	SoftFoot            int
	Load                int
	CableConnection     int
	NoiseVibration      int
	Heat                int
	IncomingVoltage     int
	Rotor               int
	Stator              int
	Bearing             int
	DCLink              int
}

func NewMotor() *Motor {
	return &Motor{
		TorqueRipple: -1, MechanicalUnbalance: -1, SoftFoot: -1,
		Load: -1, CableConnection: -1, NoiseVibration: -1, Heat: -1,
		IncomingVoltage: -1, Rotor: -1, Stator: -1, Bearing: -1, DCLink: -1,
	}
}

func (m *Motor) GetAllValues() map[string]int {
	return structToMap(m)
}

// MotorFeed holds diagnosis data for motor feed equipment.
type MotorFeed struct {
	Load            int
	CableConnection int
	NoiseVibration  int
	Heat            int
	IncomingVoltage int
	DCLink          int
}

func NewMotorFeed() *MotorFeed {
	return &MotorFeed{
		Load: -1, CableConnection: -1, NoiseVibration: -1,
		Heat: -1, IncomingVoltage: -1, DCLink: -1,
	}
}

func (mf *MotorFeed) GetAllValues() map[string]int {
	return structToMap(mf)
}

// PowerSupply holds diagnosis data for power supply equipment.
type PowerSupply struct {
	NoiseVibration  int
	Heat            int
	IncomingVoltage int
}

func NewPowerSupply() *PowerSupply {
	return &PowerSupply{
		NoiseVibration: -1, Heat: -1, IncomingVoltage: -1,
	}
}

func (ps *PowerSupply) GetAllValues() map[string]int {
	return structToMap(ps)
}

// Transformer holds diagnosis data for transformer equipment.
type Transformer struct {
	Core            int
	Load            int
	GroundFaultD    int
	Capacitor       int
	TapChanger      int
	Bushings        int
	Stress          int
	LoadUnbalance   int
	CableConnection int
	Winding         int
	NoiseVibration  int
	Heat            int
	IncomingVoltage int
}

func NewTransformer() *Transformer {
	return &Transformer{
		Core: -1, Load: -1, GroundFaultD: -1, Capacitor: -1, TapChanger: -1,
		Bushings: -1, Stress: -1, LoadUnbalance: -1, CableConnection: -1,
		Winding: -1, NoiseVibration: -1, Heat: -1, IncomingVoltage: -1,
	}
}

func (t *Transformer) GetAllValues() map[string]int {
	return structToMap(t)
}

// PrimaryTransformer holds diagnosis data for primary transformer equipment.
type PrimaryTransformer struct {
	Capacitor       int
	TapChanger      int
	Bushings        int
	Stress          int
	LoadUnbalance   int
	CableConnection int
	Winding         int
	NoiseVibration  int
	Heat            int
	IncomingVoltage int
}

func NewPrimaryTransformer() *PrimaryTransformer {
	return &PrimaryTransformer{
		Capacitor: -1, TapChanger: -1, Bushings: -1, Stress: -1,
		LoadUnbalance: -1, CableConnection: -1, Winding: -1,
		NoiseVibration: -1, Heat: -1, IncomingVoltage: -1,
	}
}

func (pt *PrimaryTransformer) GetAllValues() map[string]int {
	return structToMap(pt)
}

// --- PQ, FT, EV structs ---

// PQ holds power quality diagnosis data.
type PQ struct {
	VoltagePhaseAngle    int
	CurrentRMS           int
	CrestFactor          int
	Unbalance            int
	Harmonics            int
	ZeroSequence         int
	NegativeSequence     int
	CurrentPhaseAngle    int
	PhaseAngle           int
	PowerFactor          int
	TotalDemandDistortion int
	Power                int
	VoltageRMS           int
	DC                   int
	Events               int
}

func NewPQ() *PQ {
	return &PQ{
		VoltagePhaseAngle: -1, CurrentRMS: -1, CrestFactor: -1, Unbalance: -1,
		Harmonics: -1, ZeroSequence: -1, NegativeSequence: -1, CurrentPhaseAngle: -1,
		PhaseAngle: -1, PowerFactor: -1, TotalDemandDistortion: -1, Power: -1,
		VoltageRMS: -1, DC: -1, Events: -1,
	}
}

func (pq *PQ) GetAllValues() map[string]int {
	return structToMap(pq)
}

// FT (Fault) holds fault diagnosis data.
type FT struct {
	PhaseOrder   int
	NoLoad       int
	OverCurrent  int
	CF           int
	NoPower      int
	OverVoltage  int
	UnderVoltage int
	LowFrequency int
	VF           int
}

func NewFT() *FT {
	return &FT{
		PhaseOrder: -1, NoLoad: -1, OverCurrent: -1, CF: -1,
		NoPower: -1, OverVoltage: -1, UnderVoltage: -1, LowFrequency: -1, VF: -1,
	}
}

func (ft *FT) GetAllValues() map[string]int {
	return structToMap(ft)
}

// EV (Event) holds event diagnosis data.
type EV struct {
	TransientCurrentEvent int
	OverCurrentEvent      int
	UnderCurrentEvent     int
	SagEvent              int
	SwellEvent            int
	InterruptionEvent     int
	TransientVoltageEvent int
}

func NewEV() *EV {
	return &EV{
		TransientCurrentEvent: -1, OverCurrentEvent: -1, UnderCurrentEvent: -1,
		SagEvent: -1, SwellEvent: -1, InterruptionEvent: -1, TransientVoltageEvent: -1,
	}
}

func (ev *EV) GetAllValues() map[string]int {
	return structToMap(ev)
}

// --- Wrapper types ---

// Diagnosis wraps a diagnosis data instance with its equipment type.
type Diagnosis struct {
	Type int
	Data DiagnosisData
}

// NewDiagnosis creates a new Diagnosis for the given equipment type string.
func NewDiagnosis(equipmentType string) *Diagnosis {
	d := &Diagnosis{}
	d.Type = getTypeCode(equipmentType)
	d.initializeData()
	return d
}

func getTypeCode(equipmentType string) int {
	typeMapping := map[string]int{
		"Compressor":         EquipmentCompressor,
		"Fan":                EquipmentFan,
		"Pump":               EquipmentPump,
		"Motor":              EquipmentMotor,
		"PSupply":            EquipmentPowerSupply,
		"PowerSupply":        EquipmentPowerSupply,
		"Transformer":        EquipmentTransformer,
		"MotorFeed":          EquipmentMotorFeed,
		"PrimaryTransformer": EquipmentPrimaryTransformer,
	}
	if code, ok := typeMapping[equipmentType]; ok {
		return code
	}
	return EquipmentNone
}

func (d *Diagnosis) initializeData() {
	switch d.Type {
	case EquipmentCompressor, EquipmentFan:
		d.Data = NewCompressorFan()
	case EquipmentPump:
		d.Data = NewPump()
	case EquipmentMotor:
		d.Data = NewMotor()
	case EquipmentPowerSupply:
		d.Data = NewPowerSupply()
	case EquipmentTransformer:
		d.Data = NewTransformer()
	case EquipmentMotorFeed:
		d.Data = NewMotorFeed()
	case EquipmentPrimaryTransformer:
		d.Data = NewPrimaryTransformer()
	}
}

// UpdateFromBargraphData updates fields from bargraph data list and returns updated field->status map.
func (d *Diagnosis) UpdateFromBargraphData(bargraphData []interface{}) map[string]interface{} {
	if d.Data == nil {
		return map[string]interface{}{}
	}

	updatedFields := make(map[string]interface{})

	for _, item := range bargraphData {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		fieldName, _ := itemMap["Name"].(string)
		statusValue, hasStatus := itemMap["Status"]
		if fieldName == "" || !hasStatus || statusValue == nil {
			continue
		}
		fieldName = strings.ReplaceAll(fieldName, " ", "")
		statusInt := toInt(statusValue)

		if setStructField(d.Data, fieldName, statusInt) {
			updatedFields[fieldName] = statusInt
		}
	}
	updatedFields["type"] = d.Type
	return updatedFields
}

// UpdateFromBargraphDataForInflux updates fields and returns extended info including ID, Titles, Descriptions.
func (d *Diagnosis) UpdateFromBargraphDataForInflux(bargraphData []interface{}) map[string]interface{} {
	if d.Data == nil {
		return map[string]interface{}{}
	}

	updatedFields := make(map[string]interface{})

	for _, item := range bargraphData {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		fieldName, _ := itemMap["Name"].(string)
		statusValue, hasStatus := itemMap["Status"]
		if fieldName == "" || !hasStatus || statusValue == nil {
			continue
		}
		fieldName = strings.ReplaceAll(fieldName, " ", "")
		statusInt := toInt(statusValue)

		if setStructField(d.Data, fieldName, statusInt) {
			itemID := itemMap["ID"]
			titles := toStringMap(itemMap["Titles"])
			descriptions := toStringMap(itemMap["Descriptions"])

			updatedFields[fieldName] = map[string]interface{}{
				"id":             itemID,
				"status":         statusInt,
				"title_en":       mapGetDefault(titles, "en", fieldName),
				"title_ko":       mapGetDefault(titles, "ko", fieldName),
				"title_ja":       mapGetDefault(titles, "ja", fieldName),
				"description_en": mapGetDefault(descriptions, "en", ""),
				"description_ko": mapGetDefault(descriptions, "ko", ""),
				"description_ja": mapGetDefault(descriptions, "ja", ""),
			}
		}
	}
	updatedFields["type"] = d.Type
	return updatedFields
}

// GetAllValues returns all field values as a map.
func (d *Diagnosis) GetAllValues() map[string]int {
	if d.Data == nil {
		return map[string]int{}
	}
	return d.Data.GetAllValues()
}

// GetNonZeroValues returns only fields with non-zero values.
func (d *Diagnosis) GetNonZeroValues() map[string]int {
	all := d.GetAllValues()
	result := make(map[string]int)
	for k, v := range all {
		if v != 0 {
			result[k] = v
		}
	}
	return result
}

// GetConfiguredValues returns only fields that have been set (not -1).
func (d *Diagnosis) GetConfiguredValues() map[string]int {
	all := d.GetAllValues()
	result := make(map[string]int)
	for k, v := range all {
		if v != -1 {
			result[k] = v
		}
	}
	return result
}

// PowerQuality wraps PQ data.
type PowerQuality struct {
	Data *PQ
}

// NewPowerQuality creates a new PowerQuality instance.
func NewPowerQuality() *PowerQuality {
	return &PowerQuality{Data: NewPQ()}
}

// UpdateFromBargraphData updates PQ fields from bargraph data list.
func (pq *PowerQuality) UpdateFromBargraphData(bargraphData []interface{}) map[string]interface{} {
	if pq.Data == nil {
		return map[string]interface{}{}
	}

	updatedFields := make(map[string]interface{})

	for _, item := range bargraphData {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		fieldName, _ := itemMap["Name"].(string)
		statusValue, hasStatus := itemMap["Status"]
		if fieldName == "" || !hasStatus || statusValue == nil {
			continue
		}
		statusInt := toInt(statusValue)

		if setStructField(pq.Data, fieldName, statusInt) {
			updatedFields[fieldName] = statusInt
		}
	}
	return updatedFields
}

// UpdateFromBargraphDataForInflux updates PQ fields and returns extended info.
func (pq *PowerQuality) UpdateFromBargraphDataForInflux(bargraphData []interface{}) map[string]interface{} {
	if pq.Data == nil {
		return map[string]interface{}{}
	}

	updatedFields := make(map[string]interface{})

	for _, item := range bargraphData {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		fieldName, _ := itemMap["Name"].(string)
		statusValue, hasStatus := itemMap["Status"]
		if fieldName == "" || !hasStatus || statusValue == nil {
			continue
		}
		fieldName = strings.ReplaceAll(fieldName, " ", "")
		statusInt := toInt(statusValue)

		if setStructField(pq.Data, fieldName, statusInt) {
			itemID := itemMap["ID"]
			titles := toStringMap(itemMap["Titles"])
			descriptions := toStringMap(itemMap["Descriptions"])

			updatedFields[fieldName] = map[string]interface{}{
				"id":             itemID,
				"status":         statusInt,
				"title_en":       mapGetDefault(titles, "en", fieldName),
				"title_ko":       mapGetDefault(titles, "ko", fieldName),
				"title_ja":       mapGetDefault(titles, "ja", fieldName),
				"description_en": mapGetDefault(descriptions, "en", ""),
				"description_ko": mapGetDefault(descriptions, "ko", ""),
				"description_ja": mapGetDefault(descriptions, "ja", ""),
			}
		}
	}
	return updatedFields
}

// Fault wraps FT data.
type Fault struct {
	Data *FT
}

// NewFault creates a new Fault instance.
func NewFault() *Fault {
	return &Fault{Data: NewFT()}
}

// UpdateFromBargraphData updates FT fields from bargraph data list.
func (f *Fault) UpdateFromBargraphData(bargraphData []interface{}) map[string]interface{} {
	if f.Data == nil {
		return map[string]interface{}{}
	}

	updatedFields := make(map[string]interface{})

	for _, item := range bargraphData {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		fieldName, _ := itemMap["Name"].(string)
		statusValue, hasStatus := itemMap["Status"]
		if fieldName == "" || !hasStatus || statusValue == nil {
			continue
		}
		statusInt := toInt(statusValue)

		if setStructField(f.Data, fieldName, statusInt) {
			updatedFields[fieldName] = statusInt
		}
	}
	return updatedFields
}

// UpdateFromBargraphDataForInflux updates FT fields and returns extended info.
func (f *Fault) UpdateFromBargraphDataForInflux(bargraphData []interface{}) map[string]interface{} {
	if f.Data == nil {
		return map[string]interface{}{}
	}

	updatedFields := make(map[string]interface{})

	for _, item := range bargraphData {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		fieldName, _ := itemMap["Name"].(string)
		statusValue, hasStatus := itemMap["Status"]
		if fieldName == "" || !hasStatus || statusValue == nil {
			continue
		}
		fieldName = strings.ReplaceAll(fieldName, " ", "")
		statusInt := toInt(statusValue)

		if setStructField(f.Data, fieldName, statusInt) {
			itemID := itemMap["ID"]
			titles := toStringMap(itemMap["Titles"])
			descriptions := toStringMap(itemMap["Descriptions"])

			updatedFields[fieldName] = map[string]interface{}{
				"id":             itemID,
				"status":         statusInt,
				"title_en":       mapGetDefault(titles, "en", fieldName),
				"title_ko":       mapGetDefault(titles, "ko", fieldName),
				"title_ja":       mapGetDefault(titles, "ja", fieldName),
				"description_en": mapGetDefault(descriptions, "en", ""),
				"description_ko": mapGetDefault(descriptions, "ko", ""),
				"description_ja": mapGetDefault(descriptions, "ja", ""),
			}
		}
	}
	return updatedFields
}

// Event wraps EV data.
type Event struct {
	Data *EV
}

// NewEvent creates a new Event instance.
func NewEvent() *Event {
	return &Event{Data: NewEV()}
}

// UpdateFromBargraphData updates EV fields from bargraph data list.
func (e *Event) UpdateFromBargraphData(bargraphData []interface{}) map[string]interface{} {
	if e.Data == nil {
		return map[string]interface{}{}
	}

	updatedFields := make(map[string]interface{})

	for _, item := range bargraphData {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		fieldName, _ := itemMap["Name"].(string)
		statusValue, hasStatus := itemMap["Status"]
		if fieldName == "" || !hasStatus || statusValue == nil {
			continue
		}
		statusInt := toInt(statusValue)

		if setStructField(e.Data, fieldName, statusInt) {
			updatedFields[fieldName] = statusInt
		}
	}
	return updatedFields
}

// UpdateFromBargraphDataForInflux updates EV fields and returns extended info.
func (e *Event) UpdateFromBargraphDataForInflux(bargraphData []interface{}) map[string]interface{} {
	if e.Data == nil {
		return map[string]interface{}{}
	}

	updatedFields := make(map[string]interface{})

	for _, item := range bargraphData {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		fieldName, _ := itemMap["Name"].(string)
		statusValue, hasStatus := itemMap["Status"]
		if fieldName == "" || !hasStatus || statusValue == nil {
			continue
		}
		fieldName = strings.ReplaceAll(fieldName, " ", "")
		statusInt := toInt(statusValue)

		if setStructField(e.Data, fieldName, statusInt) {
			itemID := itemMap["ID"]
			titles := toStringMap(itemMap["Titles"])
			descriptions := toStringMap(itemMap["Descriptions"])

			updatedFields[fieldName] = map[string]interface{}{
				"id":             itemID,
				"status":         statusInt,
				"title_en":       mapGetDefault(titles, "en", fieldName),
				"title_ko":       mapGetDefault(titles, "ko", fieldName),
				"title_ja":       mapGetDefault(titles, "ja", fieldName),
				"description_en": mapGetDefault(descriptions, "en", ""),
				"description_ko": mapGetDefault(descriptions, "ko", ""),
				"description_ja": mapGetDefault(descriptions, "ja", ""),
			}
		}
	}
	return updatedFields
}

// --- Helper functions ---

// structToMap converts a struct pointer to a map[string]int using reflection.
func structToMap(s interface{}) map[string]int {
	result := make(map[string]int)
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		result[t.Field(i).Name] = int(v.Field(i).Int())
	}
	return result
}

// setStructField sets a named int field on a struct pointer. Returns true if the field exists and was set.
func setStructField(s interface{}, fieldName string, value int) bool {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	field := v.FieldByName(fieldName)
	if !field.IsValid() || !field.CanSet() {
		return false
	}
	field.SetInt(int64(value))
	return true
}

// toInt converts an interface{} value to int (handles float64 from JSON and int).
func toInt(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	case float32:
		return int(val)
	default:
		return 0
	}
}

// toStringMap safely converts an interface{} to map[string]interface{}.
func toStringMap(v interface{}) map[string]interface{} {
	if v == nil {
		return map[string]interface{}{}
	}
	if m, ok := v.(map[string]interface{}); ok {
		return m
	}
	return map[string]interface{}{}
}

// mapGetDefault gets a string value from a map, returning defaultVal if not found.
func mapGetDefault(m map[string]interface{}, key string, defaultVal string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultVal
}

// --- Compatibility aliases for processors ---

// EventWrapper is an alias for Event (used by diagnosis_processor).
type EventWrapper = Event

// FaultWrapper is an alias for Fault (used by diagnosis_processor).
type FaultWrapper = Fault

// InitPQDefaults sets all PQ fields to -1.
func InitPQDefaults(pq *PQ) {
	pq.VoltagePhaseAngle = -1
	pq.CurrentRMS = -1
	pq.CrestFactor = -1
	pq.Unbalance = -1
	pq.Harmonics = -1
	pq.ZeroSequence = -1
	pq.NegativeSequence = -1
	pq.CurrentPhaseAngle = -1
	pq.PhaseAngle = -1
	pq.PowerFactor = -1
	pq.TotalDemandDistortion = -1
	pq.Power = -1
	pq.VoltageRMS = -1
	pq.DC = -1
	pq.Events = -1
}

// InitFTDefaults sets all FT fields to -1.
func InitFTDefaults(ft *FT) {
	ft.PhaseOrder = -1
	ft.NoLoad = -1
	ft.OverCurrent = -1
	ft.CF = -1
	ft.NoPower = -1
	ft.OverVoltage = -1
	ft.UnderVoltage = -1
	ft.LowFrequency = -1
	ft.VF = -1
}

// InitEVDefaults sets all EV fields to -1.
func InitEVDefaults(ev *EV) {
	ev.TransientCurrentEvent = -1
	ev.OverCurrentEvent = -1
	ev.UnderCurrentEvent = -1
	ev.SagEvent = -1
	ev.SwellEvent = -1
	ev.InterruptionEvent = -1
	ev.TransientVoltageEvent = -1
}
