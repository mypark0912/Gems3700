package en50160

import (
	"fmt"
	"math"
	"sort"
)

// Histogram mirrors the dict returned by python _calculate_histogram.
type Histogram struct {
	Bins        []float64 `json:"bins"`
	Counts      []int     `json:"counts"`
	Percentages []float64 `json:"percentages"`
	BinLabels   []string  `json:"bin_labels"`
	BinCenters  []float64 `json:"bin_centers"`
}

// CalculateHistogram ports numpy.histogram with equal-width bins plus the
// derived helpers produced by the Python implementation.
// When values is empty, returns zero-length slices.
func CalculateHistogram(values []float64, bins int, rangeMin, rangeMax float64) *Histogram {
	h := &Histogram{
		Bins:        []float64{},
		Counts:      []int{},
		Percentages: []float64{},
		BinLabels:   []string{},
		BinCenters:  []float64{},
	}
	if len(values) == 0 || bins <= 0 {
		return h
	}

	if rangeMin == rangeMax {
		// numpy widens the range so values fall inside one bin.
		rangeMax = rangeMin + 1
	}
	step := (rangeMax - rangeMin) / float64(bins)

	edges := make([]float64, bins+1)
	for i := 0; i <= bins; i++ {
		edges[i] = rangeMin + step*float64(i)
	}

	counts := make([]int, bins)
	for _, v := range values {
		// numpy.histogram: [edge_i, edge_i+1) except last which is closed.
		if v < rangeMin || v > rangeMax {
			continue
		}
		idx := int((v - rangeMin) / step)
		if idx == bins {
			idx = bins - 1
		}
		if idx >= 0 && idx < bins {
			counts[idx]++
		}
	}

	total := float64(len(values))
	centers := make([]float64, bins)
	labels := make([]string, bins)
	pct := make([]float64, bins)
	for i := 0; i < bins; i++ {
		centers[i] = roundN((edges[i]+edges[i+1])/2, 4)
		labels[i] = fmt.Sprintf("%.2f-%.2f", edges[i], edges[i+1])
		pct[i] = roundN(float64(counts[i])/total*100, 2)
	}

	binsRounded := make([]float64, len(edges))
	for i, e := range edges {
		binsRounded[i] = roundN(e, 4)
	}

	h.Bins = binsRounded
	h.Counts = counts
	h.Percentages = pct
	h.BinLabels = labels
	h.BinCenters = centers
	return h
}

// -----------------------------------------------------------------------------
// numpy-like helpers
// -----------------------------------------------------------------------------

// Min returns the smallest value, or 0 on empty input.
func Min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	m := values[0]
	for _, v := range values[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

// Max returns the largest value, or 0 on empty input.
func Max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	m := values[0]
	for _, v := range values[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

// Mean returns the arithmetic mean, or 0 on empty input.
func Mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// Percentile ports numpy.percentile with default "linear" interpolation.
// p is in [0, 100].
func Percentile(values []float64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sorted := append([]float64(nil), values...)
	sort.Float64s(sorted)

	if p <= 0 {
		return sorted[0]
	}
	if p >= 100 {
		return sorted[len(sorted)-1]
	}

	rank := p / 100 * float64(len(sorted)-1)
	lo := int(math.Floor(rank))
	hi := int(math.Ceil(rank))
	if lo == hi {
		return sorted[lo]
	}
	frac := rank - float64(lo)
	return sorted[lo] + frac*(sorted[hi]-sorted[lo])
}

// CountInRange returns how many values fall within [min, max] inclusive.
func CountInRange(values []float64, min, max float64) int {
	var n int
	for _, v := range values {
		if v >= min && v <= max {
			n++
		}
	}
	return n
}

// CountLE returns how many values are ≤ limit.
func CountLE(values []float64, limit float64) int {
	var n int
	for _, v := range values {
		if v <= limit {
			n++
		}
	}
	return n
}

// DropNaN returns values with NaN filtered out (mirrors pandas df.dropna()).
func DropNaN(values []float64) []float64 {
	out := make([]float64, 0, len(values))
	for _, v := range values {
		if !math.IsNaN(v) {
			out = append(out, v)
		}
	}
	return out
}

func roundN(v float64, digits int) float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	pow := math.Pow10(digits)
	return math.Round(v*pow) / pow
}

// verdict returns "PASS" when pct ≥ target, else "FAIL".
func verdict(pct, target float64) string {
	if pct >= target {
		return "PASS"
	}
	return "FAIL"
}
