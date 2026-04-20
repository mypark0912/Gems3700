package data

// GetIntFromInterface safely converts interface{} to int.
func GetIntFromInterface(v interface{}) int {
	switch val := v.(type) {
	case float64:
		return int(val)
	case int:
		return val
	case int64:
		return int(val)
	default:
		return 0
	}
}

// Channel represents a measurement channel with its configuration.
type Channel struct {
	Name          string
	Trend         bool
	Diagnosis     bool
	Alarm         bool
	AlarmComDelay int
	Period        int
	TrendList     []string
	AssetName     string
	AssetType     string
	DemandTrend   int
	DemandPeriod  int
	UseDO         bool
	UseConfStatus bool
	ConfStatus    map[string]interface{}
	AssetDrive    bool
	VoltageType   int
}

// NewChannel creates a new Channel with the given name and sensible defaults.
func NewChannel(name string) *Channel {
	return &Channel{
		Name:       name,
		TrendList:  make([]string, 0),
		ConfStatus: make(map[string]interface{}),
	}
}
