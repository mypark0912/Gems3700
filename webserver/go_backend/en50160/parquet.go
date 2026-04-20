package en50160

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/parquet-go/parquet-go"
)

// DataFrame is a columnar, pandas-flavoured view of an EN50160 weekly parquet file.
// Every column slice has the same length (number of samples). Metadata stored as
// the "en50160_summary" schema key is decoded into Summary.
type DataFrame struct {
	Timestamps     []time.Time
	FrequencyAvg   []float64
	VoltageL1      []float64
	VoltageL2      []float64
	VoltageL3      []float64
	VoltageUnbal0  []float64
	VoltageThdL1   []float64
	VoltageThdL2   []float64
	VoltageThdL3   []float64
	PstL1          []float64
	PstL2          []float64
	PstL3          []float64
	PltL1          []float64
	PltL2          []float64
	PltL3          []float64
	HarmonicsL1    [][]float64 // [sample][24 harmonics]
	HarmonicsL2    [][]float64
	HarmonicsL3    [][]float64
	Summary        map[string]interface{}
	TotalSamples  int
}

// parquetRow matches the weekly EN50160 parquet schema written by Python/pyarrow.
// Timestamps are stored as microseconds since epoch (pyarrow default for timestamp columns).
type parquetRow struct {
	Timestamp         int64     `parquet:"timestamp,timestamp(microsecond)"`
	FrequencyAvg      float64   `parquet:"frequency_avg"`
	VoltageL1         float64   `parquet:"voltage_l1"`
	VoltageL2         float64   `parquet:"voltage_l2"`
	VoltageL3         float64   `parquet:"voltage_l3"`
	VoltageUnbalance0 float64   `parquet:"voltage_unbalance_0"`
	VoltageThdL1      float64   `parquet:"voltage_thd_l1"`
	VoltageThdL2      float64   `parquet:"voltage_thd_l2"`
	VoltageThdL3      float64   `parquet:"voltage_thd_l3"`
	PstL1             float64   `parquet:"pst_l1"`
	PstL2             float64   `parquet:"pst_l2"`
	PstL3             float64   `parquet:"pst_l3"`
	PltL1             float64   `parquet:"plt_l1"`
	PltL2             float64   `parquet:"plt_l2"`
	PltL3             float64   `parquet:"plt_l3"`
	HarmonicsL1       []float64 `parquet:"harmonics_l1"`
	HarmonicsL2       []float64 `parquet:"harmonics_l2"`
	HarmonicsL3       []float64 `parquet:"harmonics_l3"`
}

// ReadParquet opens an EN50160 weekly parquet file and returns a columnar DataFrame.
// Schema metadata under the "en50160_summary" key is JSON-decoded into df.Summary
// (matching Python's table.schema.metadata[b'en50160_summary'] path).
func ReadParquet(path string) (*DataFrame, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	pf, err := parquet.OpenFile(f, stat.Size())
	if err != nil {
		return nil, fmt.Errorf("open parquet: %w", err)
	}

	df := &DataFrame{}

	// File-level key/value metadata (pyarrow stores custom schema metadata here).
	for _, kv := range pf.Metadata().KeyValueMetadata {
		if kv.Key == "en50160_summary" {
			var summary map[string]interface{}
			if err := json.Unmarshal([]byte(kv.Value), &summary); err == nil {
				df.Summary = summary
			}
			break
		}
	}

	reader := parquet.NewGenericReader[parquetRow](pf)
	defer reader.Close()

	total := int(reader.NumRows())
	if total <= 0 {
		df.TotalSamples = 0
		return df, nil
	}

	rows := make([]parquetRow, total)
	if _, err := reader.Read(rows); err != nil {
		return nil, fmt.Errorf("read rows: %w", err)
	}

	df.Timestamps = make([]time.Time, total)
	df.FrequencyAvg = make([]float64, total)
	df.VoltageL1 = make([]float64, total)
	df.VoltageL2 = make([]float64, total)
	df.VoltageL3 = make([]float64, total)
	df.VoltageUnbal0 = make([]float64, total)
	df.VoltageThdL1 = make([]float64, total)
	df.VoltageThdL2 = make([]float64, total)
	df.VoltageThdL3 = make([]float64, total)
	df.PstL1 = make([]float64, total)
	df.PstL2 = make([]float64, total)
	df.PstL3 = make([]float64, total)
	df.PltL1 = make([]float64, total)
	df.PltL2 = make([]float64, total)
	df.PltL3 = make([]float64, total)
	df.HarmonicsL1 = make([][]float64, total)
	df.HarmonicsL2 = make([][]float64, total)
	df.HarmonicsL3 = make([][]float64, total)

	for i, r := range rows {
		// pyarrow microsecond timestamps.
		df.Timestamps[i] = time.Unix(r.Timestamp/1_000_000, (r.Timestamp%1_000_000)*1_000).UTC()
		df.FrequencyAvg[i] = r.FrequencyAvg
		df.VoltageL1[i] = r.VoltageL1
		df.VoltageL2[i] = r.VoltageL2
		df.VoltageL3[i] = r.VoltageL3
		df.VoltageUnbal0[i] = r.VoltageUnbalance0
		df.VoltageThdL1[i] = r.VoltageThdL1
		df.VoltageThdL2[i] = r.VoltageThdL2
		df.VoltageThdL3[i] = r.VoltageThdL3
		df.PstL1[i] = r.PstL1
		df.PstL2[i] = r.PstL2
		df.PstL3[i] = r.PstL3
		df.PltL1[i] = r.PltL1
		df.PltL2[i] = r.PltL2
		df.PltL3[i] = r.PltL3
		df.HarmonicsL1[i] = r.HarmonicsL1
		df.HarmonicsL2[i] = r.HarmonicsL2
		df.HarmonicsL3[i] = r.HarmonicsL3
	}
	df.TotalSamples = total
	return df, nil
}

// Records converts the DataFrame into a list of generic records, mirroring
// Python df.to_dict(orient='records') with ISO-formatted timestamps.
func (df *DataFrame) Records() []map[string]interface{} {
	if df == nil {
		return nil
	}
	out := make([]map[string]interface{}, df.TotalSamples)
	for i := 0; i < df.TotalSamples; i++ {
		out[i] = map[string]interface{}{
			"timestamp":           df.Timestamps[i].Format("2006-01-02 15:04:05"),
			"frequency_avg":       df.FrequencyAvg[i],
			"voltage_l1":          df.VoltageL1[i],
			"voltage_l2":          df.VoltageL2[i],
			"voltage_l3":          df.VoltageL3[i],
			"voltage_unbalance_0": df.VoltageUnbal0[i],
			"voltage_thd_l1":      df.VoltageThdL1[i],
			"voltage_thd_l2":      df.VoltageThdL2[i],
			"voltage_thd_l3":      df.VoltageThdL3[i],
			"pst_l1":              df.PstL1[i],
			"pst_l2":              df.PstL2[i],
			"pst_l3":              df.PstL3[i],
			"plt_l1":              df.PltL1[i],
			"plt_l2":              df.PltL2[i],
			"plt_l3":              df.PltL3[i],
			"harmonics_l1":        df.HarmonicsL1[i],
			"harmonics_l2":        df.HarmonicsL2[i],
			"harmonics_l3":        df.HarmonicsL3[i],
		}
	}
	return out
}

// timestampLabels returns "YYYY-MM-DD HH:MM" labels used by Chart.js time axes.
func (df *DataFrame) timestampLabels() []string {
	out := make([]string, df.TotalSamples)
	for i, t := range df.Timestamps {
		out[i] = t.Format("2006-01-02 15:04")
	}
	return out
}
