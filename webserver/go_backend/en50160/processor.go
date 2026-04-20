package en50160

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Processor mirrors Python EN50160ReportProcessor.
type Processor struct {
	Config           *WeeklyReportConfig
	Limits           *Limits
	NominalVoltage   float64
	NominalCurrent   float64
	NominalFrequency float64
}

// NewProcessor creates a processor with Python defaults (60Hz, 22900V, 100A).
func NewProcessor(cfg *WeeklyReportConfig) *Processor {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	return &Processor{
		Config:           cfg,
		Limits:           DefaultLimits(),
		NominalVoltage:   22900.0,
		NominalCurrent:   100.0,
		NominalFrequency: 60.0,
	}
}

// SetLimits mirrors Python set_limits — updates nominals and re-derives the
// frequency ±1 % / +4 %-6 % windows. Pass nil to keep an existing value.
func (p *Processor) SetLimits(voltage, current, frequency *float64) {
	if voltage != nil {
		p.NominalVoltage = *voltage
	}
	if current != nil {
		p.NominalCurrent = *current
	}
	if frequency != nil {
		p.NominalFrequency = *frequency
		p.Limits.Frequency.Nominal = *frequency
		p.Limits.Frequency.Limit99_5 = Range{
			Min: roundN(*frequency*0.99, 2),
			Max: roundN(*frequency*1.01, 2),
		}
		p.Limits.Frequency.Limit100 = Range{
			Min: roundN(*frequency*0.94, 2),
			Max: roundN(*frequency*1.04, 2),
		}
	}
}

// -----------------------------------------------------------------------------
// File discovery + parquet IO
// -----------------------------------------------------------------------------

// ListFiles returns unique date strings (YYYYMMDD, descending) from the reports
// folder for the given channel, matching Python list_files.
func (p *Processor) ListFiles(channel string) ([]string, error) {
	dir := filepath.Join(p.Config.OutputDir, channel)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasPrefix(name, "en50160_weekly_") || !strings.HasSuffix(name, ".parquet") {
			continue
		}
		stem := strings.TrimSuffix(name, ".parquet")
		parts := strings.Split(stem, "_")
		if len(parts) == 0 {
			continue
		}
		date := parts[len(parts)-1]
		if len(date) == 8 && isAllDigits(date) {
			seen[date] = true
		}
	}

	dates := make([]string, 0, len(seen))
	for d := range seen {
		dates = append(dates, d)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(dates)))
	return dates, nil
}

func isAllDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// ReadParquet loads the named file under <outputDir>/<channel>/<filename>.
func (p *Processor) ReadParquet(channel, filename string) (*DataFrame, error) {
	path := filepath.Join(p.Config.OutputDir, channel, filename)
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("file not found: %s", path)
	}
	return ReadParquet(path)
}

// ReadToDict mirrors Python read_to_dict — list of records with ISO timestamps.
func (p *Processor) ReadToDict(channel, filename string) ([]map[string]interface{}, error) {
	df, err := p.ReadParquet(channel, filename)
	if err != nil {
		return nil, err
	}
	return df.Records(), nil
}

// GetAllChartData mirrors Python get_all_chart_data — single bundle with every section.
func (p *Processor) GetAllChartData(channel, filename string) (map[string]interface{}, error) {
	df, err := p.ReadParquet(channel, filename)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"summary":   df.Summary,
		"frequency": p.frequencyFromDF(df),
		"voltage":   p.voltageFromDF(df),
		"thd":       p.thdFromDF(df),
		"unbalance": p.unbalanceFromDF(df),
		"flicker":   p.flickerFromDF(df),
		"harmonics": p.harmonicsFromDF(df),
	}, nil
}

// -----------------------------------------------------------------------------
// Frequency
// -----------------------------------------------------------------------------

// GetFrequencyChartData loads the file and returns Chart.js data for frequency.
func (p *Processor) GetFrequencyChartData(channel, filename string) (map[string]interface{}, error) {
	df, err := p.ReadParquet(channel, filename)
	if err != nil {
		return nil, err
	}
	return p.frequencyFromDF(df), nil
}

func (p *Processor) frequencyFromDF(df *DataFrame) map[string]interface{} {
	if df == nil || df.TotalSamples == 0 {
		return nil
	}
	limits := p.Limits.Frequency
	values := DropNaN(df.FrequencyAvg)
	labels := df.timestampLabels()

	series := make([]float64, df.TotalSamples)
	for i, v := range df.FrequencyAvg {
		series[i] = roundN(v, 3)
	}

	histogram := CalculateHistogram(values, 30, limits.Limit99_5.Min, limits.Limit99_5.Max)

	in995 := float64(CountInRange(values, limits.Limit99_5.Min, limits.Limit99_5.Max)) / float64(len(values)) * 100
	in100 := float64(CountInRange(values, limits.Limit100.Min, limits.Limit100.Max)) / float64(len(values)) * 100

	return map[string]interface{}{
		"timeseries": map[string]interface{}{"labels": labels, "data": series},
		"histogram":  histogram,
		"statistics": map[string]interface{}{
			"min":                    roundN(Min(values), 3),
			"max":                    roundN(Max(values), 3),
			"avg":                    roundN(Mean(values), 3),
			"in_range_99_5_percent":  roundN(in995, 2),
			"in_range_100_percent":   roundN(in100, 2),
			"result_99_5":            verdict(in995, 99.5),
			"result_100":             verdict(in100, 100.0),
			"total_samples":          len(values),
		},
		"limits": limits,
	}
}

// -----------------------------------------------------------------------------
// Voltage
// -----------------------------------------------------------------------------

// GetVoltageChartData loads the file and returns per-phase voltage chart data.
func (p *Processor) GetVoltageChartData(channel, filename string) (map[string]interface{}, error) {
	df, err := p.ReadParquet(channel, filename)
	if err != nil {
		return nil, err
	}
	return p.voltageFromDF(df), nil
}

func (p *Processor) voltageFromDF(df *DataFrame) map[string]interface{} {
	if df == nil || df.TotalSamples == 0 {
		return nil
	}
	limits := p.Limits.Voltage
	limit95Min := p.NominalVoltage * limits.Limit95.Min / 100
	limit95Max := p.NominalVoltage * limits.Limit95.Max / 100
	limit100Min := p.NominalVoltage * limits.Limit100.Min / 100
	limit100Max := p.NominalVoltage * limits.Limit100.Max / 100

	labels := df.timestampLabels()
	phases := map[string]interface{}{}

	for i, col := range [][]float64{df.VoltageL1, df.VoltageL2, df.VoltageL3} {
		values := DropNaN(col)
		if len(values) == 0 {
			continue
		}
		phaseKey := fmt.Sprintf("L%d", i+1)

		series := make([]float64, df.TotalSamples)
		for j, v := range col {
			series[j] = roundN(v, 1)
		}

		hist := CalculateHistogram(values, 30, limit95Min, limit95Max)
		in95 := float64(CountInRange(values, limit95Min, limit95Max)) / float64(len(values)) * 100
		in100 := float64(CountInRange(values, limit100Min, limit100Max)) / float64(len(values)) * 100

		phases[phaseKey] = map[string]interface{}{
			"timeseries": map[string]interface{}{"labels": labels, "data": series},
			"histogram":  hist,
			"statistics": map[string]interface{}{
				"min":                  roundN(Min(values), 1),
				"max":                  roundN(Max(values), 1),
				"avg":                  roundN(Mean(values), 1),
				"in_range_95_percent":  roundN(in95, 2),
				"in_range_100_percent": roundN(in100, 2),
				"result_95":            verdict(in95, 95.0),
				"result_100":           verdict(in100, 100.0),
			},
		}
	}

	return map[string]interface{}{
		"phases": phases,
		"limits": map[string]interface{}{
			"nominal":   p.NominalVoltage,
			"limit_95":  Range{Min: roundN(limit95Min, 1), Max: roundN(limit95Max, 1)},
			"limit_100": Range{Min: roundN(limit100Min, 1), Max: roundN(limit100Max, 1)},
		},
		"total_samples": df.TotalSamples,
	}
}

// -----------------------------------------------------------------------------
// Voltage unbalance
// -----------------------------------------------------------------------------

// GetUnbalanceChartData loads the file and returns unbalance chart data.
func (p *Processor) GetUnbalanceChartData(channel, filename string) (map[string]interface{}, error) {
	df, err := p.ReadParquet(channel, filename)
	if err != nil {
		return nil, err
	}
	return p.unbalanceFromDF(df), nil
}

func (p *Processor) unbalanceFromDF(df *DataFrame) map[string]interface{} {
	if df == nil || df.TotalSamples == 0 {
		return nil
	}
	limit := p.Limits.VoltageUnbalance["limit_95"]

	values := DropNaN(df.VoltageUnbal0)
	if len(values) == 0 {
		return nil
	}
	labels := df.timestampLabels()

	series := make([]float64, df.TotalSamples)
	for i, v := range df.VoltageUnbal0 {
		series[i] = roundN(v, 3)
	}

	hist := CalculateHistogram(values, 25, 0.0, limit*1.5)
	in95 := float64(CountLE(values, limit)) / float64(len(values)) * 100

	return map[string]interface{}{
		"timeseries": map[string]interface{}{"labels": labels, "data": series},
		"histogram":  hist,
		"statistics": map[string]interface{}{
			"max":                 roundN(Max(values), 3),
			"avg":                 roundN(Mean(values), 3),
			"percentile_95":       roundN(Percentile(values, 95), 3),
			"in_range_95_percent": roundN(in95, 2),
			"result":              verdict(in95, 95.0),
		},
		"limits": map[string]interface{}{"limit_95": limit},
	}
}

// -----------------------------------------------------------------------------
// Voltage THD
// -----------------------------------------------------------------------------

// GetTHDChartData loads the file and returns per-phase THD chart data.
func (p *Processor) GetTHDChartData(channel, filename string) (map[string]interface{}, error) {
	df, err := p.ReadParquet(channel, filename)
	if err != nil {
		return nil, err
	}
	return p.thdFromDF(df), nil
}

func (p *Processor) thdFromDF(df *DataFrame) map[string]interface{} {
	if df == nil || df.TotalSamples == 0 {
		return nil
	}
	limit := p.Limits.THD["limit_95"]
	labels := df.timestampLabels()

	phases := map[string]interface{}{}
	for i, col := range [][]float64{df.VoltageThdL1, df.VoltageThdL2, df.VoltageThdL3} {
		values := DropNaN(col)
		if len(values) == 0 {
			continue
		}
		phaseKey := fmt.Sprintf("L%d", i+1)

		series := make([]float64, df.TotalSamples)
		for j, v := range col {
			series[j] = roundN(v, 2)
		}
		hist := CalculateHistogram(values, 25, 0.0, limit*1.25)
		in95 := float64(CountLE(values, limit)) / float64(len(values)) * 100

		phases[phaseKey] = map[string]interface{}{
			"timeseries": map[string]interface{}{"labels": labels, "data": series},
			"histogram":  hist,
			"statistics": map[string]interface{}{
				"max":                 roundN(Max(values), 2),
				"avg":                 roundN(Mean(values), 2),
				"percentile_95":       roundN(Percentile(values, 95), 2),
				"in_range_95_percent": roundN(in95, 2),
				"result":              verdict(in95, 95.0),
			},
		}
	}

	return map[string]interface{}{
		"phases":        phases,
		"limits":        map[string]interface{}{"limit_95": limit},
		"total_samples": df.TotalSamples,
	}
}

// -----------------------------------------------------------------------------
// Flicker (Pst/Plt)
// -----------------------------------------------------------------------------

// GetFlickerChartData loads the file and returns Pst/Plt chart data.
func (p *Processor) GetFlickerChartData(channel, filename string) (map[string]interface{}, error) {
	df, err := p.ReadParquet(channel, filename)
	if err != nil {
		return nil, err
	}
	return p.flickerFromDF(df), nil
}

func (p *Processor) flickerFromDF(df *DataFrame) map[string]interface{} {
	if df == nil || df.TotalSamples == 0 {
		return nil
	}
	pstLimit := p.Limits.Flicker["pst_limit"]
	pltLimit := p.Limits.Flicker["plt_limit_95"]
	labels := df.timestampLabels()

	pst := map[string]interface{}{}
	plt := map[string]interface{}{}

	pstCols := [][]float64{df.PstL1, df.PstL2, df.PstL3}
	pltCols := [][]float64{df.PltL1, df.PltL2, df.PltL3}

	for i := 0; i < 3; i++ {
		phaseKey := fmt.Sprintf("L%d", i+1)

		if pstValues := DropNaN(pstCols[i]); len(pstValues) > 0 {
			series := make([]float64, df.TotalSamples)
			for j, v := range pstCols[i] {
				series[j] = roundN(v, 3)
			}
			hist := CalculateHistogram(pstValues, 25, 0.0, pstLimit*2.0)

			pst[phaseKey] = map[string]interface{}{
				"timeseries": map[string]interface{}{"labels": labels, "data": series},
				"histogram":  hist,
				"statistics": map[string]interface{}{
					"max":           roundN(Max(pstValues), 3),
					"avg":           roundN(Mean(pstValues), 3),
					"percentile_95": roundN(Percentile(pstValues, 95), 3),
				},
			}
		}

		if pltValues := DropNaN(pltCols[i]); len(pltValues) > 0 {
			series := make([]float64, df.TotalSamples)
			for j, v := range pltCols[i] {
				series[j] = roundN(v, 3)
			}
			hist := CalculateHistogram(pltValues, 25, 0.0, pltLimit*2.0)
			inRange := float64(CountLE(pltValues, pltLimit)) / float64(len(pltValues)) * 100

			plt[phaseKey] = map[string]interface{}{
				"timeseries": map[string]interface{}{"labels": labels, "data": series},
				"histogram":  hist,
				"statistics": map[string]interface{}{
					"max":                 roundN(Max(pltValues), 3),
					"avg":                 roundN(Mean(pltValues), 3),
					"percentile_95":       roundN(Percentile(pltValues, 95), 3),
					"in_range_95_percent": roundN(inRange, 2),
					"result":              verdict(inRange, 95.0),
				},
			}
		}
	}

	return map[string]interface{}{
		"pst":           pst,
		"plt":           plt,
		"limits":        map[string]interface{}{"pst": pstLimit, "plt_95": pltLimit},
		"total_samples": df.TotalSamples,
	}
}

// -----------------------------------------------------------------------------
// Harmonics
// -----------------------------------------------------------------------------

// GetHarmonicsTableData loads the file and returns per-phase H2..H25 stats.
func (p *Processor) GetHarmonicsTableData(channel, filename string) (map[string]interface{}, error) {
	df, err := p.ReadParquet(channel, filename)
	if err != nil {
		return nil, err
	}
	return p.harmonicsFromDF(df), nil
}

func (p *Processor) harmonicsFromDF(df *DataFrame) map[string]interface{} {
	if df == nil || df.TotalSamples == 0 {
		return nil
	}
	harmLimits := p.Limits.Harmonics

	phases := map[string]interface{}{}
	harmCols := [3][][]float64{df.HarmonicsL1, df.HarmonicsL2, df.HarmonicsL3}

	for i := 0; i < 3; i++ {
		phaseKey := fmt.Sprintf("L%d", i+1)
		phaseMap := map[string]interface{}{}

		rows := harmCols[i]
		if len(rows) == 0 {
			phases[phaseKey] = phaseMap
			continue
		}

		// 24 harmonics per row: H2..H25.
		for h := 0; h < 24; h++ {
			hNum := h + 2
			hKey := fmt.Sprintf("h%d", hNum)
			limit, ok := harmLimits[hKey]
			if !ok {
				limit = 0.5
			}

			values := make([]float64, 0, len(rows))
			for _, row := range rows {
				if h < len(row) {
					values = append(values, row[h])
				}
			}
			if len(values) == 0 {
				continue
			}

			inRange := float64(CountLE(values, limit)) / float64(len(values)) * 100
			hist := CalculateHistogram(values, 20, 0.0, limit*1.5)

			phaseMap[hKey] = map[string]interface{}{
				"max":                 roundN(Max(values), 3),
				"avg":                 roundN(Mean(values), 3),
				"percentile_95":       roundN(Percentile(values, 95), 3),
				"limit":               limit,
				"in_range_95_percent": roundN(inRange, 2),
				"result":              verdict(inRange, 95.0),
				"histogram":           hist,
			}
		}
		phases[phaseKey] = phaseMap
	}

	return map[string]interface{}{
		"phases":        phases,
		"limits":        harmLimits,
		"total_samples": df.TotalSamples,
	}
}

// WorstResult returns "FAIL" > "N/A" > "PASS" — matches Python _get_worst_result.
func WorstResult(results []string) string {
	for _, r := range results {
		if r == "FAIL" {
			return "FAIL"
		}
	}
	for _, r := range results {
		if r == "N/A" {
			return "N/A"
		}
	}
	return "PASS"
}
