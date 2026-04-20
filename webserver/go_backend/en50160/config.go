// Package en50160 is a Go port of python_backend/utils/en50160_dataMap.py.
// It reads weekly EN50160 report parquet files and produces Chart.js-ready chart
// datasets (time series, histograms, statistics, compliance verdicts).
package en50160

// Range holds a {min, max} pair for compliance limits.
type Range struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// FrequencyLimits mirrors Python self.limits["frequency"].
type FrequencyLimits struct {
	Nominal   float64 `json:"nominal"`
	Limit99_5 Range   `json:"limit_99_5"`
	Limit100  Range   `json:"limit_100"`
}

// VoltageLimits mirrors Python self.limits["voltage"].
// The two Range values are % of nominal and get multiplied by nominal_voltage at runtime.
type VoltageLimits struct {
	NominalPercent float64 `json:"nominal_percent"`
	Limit95        Range   `json:"limit_95"`
	Limit100       Range   `json:"limit_100"`
}

// Limits mirrors Python EN50160_LIMITS_DEFAULT.
type Limits struct {
	Frequency        FrequencyLimits    `json:"frequency"`
	Voltage          VoltageLimits      `json:"voltage"`
	VoltageUnbalance map[string]float64 `json:"voltage_unbalance"`
	THD              map[string]float64 `json:"thd"`
	Flicker          map[string]float64 `json:"flicker"`
	Harmonics        map[string]float64 `json:"harmonics"`
}

// DefaultLimits returns a fresh copy of EN50160_LIMITS_DEFAULT (60Hz system).
func DefaultLimits() *Limits {
	return &Limits{
		Frequency: FrequencyLimits{
			Nominal:   60.0,
			Limit99_5: Range{Min: 59.40, Max: 60.60}, // ±1 %
			Limit100:  Range{Min: 56.40, Max: 62.40}, // +4 %/-6 %
		},
		Voltage: VoltageLimits{
			NominalPercent: 100.0,
			Limit95:        Range{Min: 90.0, Max: 110.0},
			Limit100:       Range{Min: 85.0, Max: 110.0},
		},
		VoltageUnbalance: map[string]float64{"limit_95": 2.0},
		THD:              map[string]float64{"limit_95": 8.0},
		Flicker: map[string]float64{
			"pst_limit":    1.0,
			"plt_limit_95": 1.0,
		},
		Harmonics: map[string]float64{
			"h2": 2.0, "h3": 5.0, "h4": 1.0, "h5": 6.0,
			"h6": 0.5, "h7": 5.0, "h8": 0.5, "h9": 1.5,
			"h10": 0.5, "h11": 3.5, "h12": 0.5, "h13": 3.0,
			"h14": 0.5, "h15": 0.5, "h16": 0.5, "h17": 2.0,
			"h18": 0.5, "h19": 1.5, "h20": 0.5, "h21": 0.5,
			"h22": 0.5, "h23": 1.5, "h24": 0.5, "h25": 1.5,
		},
	}
}

// WeeklyReportConfig mirrors Python WeeklyReportConfig dataclass.
type WeeklyReportConfig struct {
	OutputDir      string
	Bucket         string
	Measurement    string
	Channels       []string
	RetentionWeeks int
}

// DefaultConfig returns a fresh WeeklyReportConfig with Python defaults.
func DefaultConfig() *WeeklyReportConfig {
	return &WeeklyReportConfig{
		OutputDir:      "/usr/local/sv500/reports",
		Bucket:         "ntek30",
		Measurement:    "en10min",
		Channels:       []string{"Main", "Sub"},
		RetentionWeeks: 12,
	}
}
