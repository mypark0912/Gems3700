package data

import (
	"time"
)

// EnergyConsumptionCalculator computes energy consumption from metered values.
type EnergyConsumptionCalculator struct {
	PreviousValues map[string]float64
	CurrentValues  map[string]float64
}

// NewEnergyConsumptionCalculator creates a new calculator instance.
func NewEnergyConsumptionCalculator() *EnergyConsumptionCalculator {
	return &EnergyConsumptionCalculator{
		PreviousValues: make(map[string]float64),
		CurrentValues:  make(map[string]float64),
	}
}

// CalculateConsumption computes the difference between current and previous values
// for each energy parameter. Returns a map of parameter name to consumption value.
func (ec *EnergyConsumptionCalculator) CalculateConsumption(current map[string]float64, previous map[string]float64) map[string]float64 {
	consumption := make(map[string]float64)

	for _, field := range []string{"kwh_import", "kwh_export", "kvarh_import", "kvarh_export", "kvah_import", "kvah_export"} {
		currentVal, hasCurrent := current[field]
		prevVal, hasPrev := previous[field]

		consumptionKey := field + "_consumption"
		if hasCurrent && hasPrev {
			diff := currentVal - prevVal
			if diff < 0 {
				// Handle meter rollover or comm error: set consumption to 0
				diff = 0.0
			}
			consumption[consumptionKey] = diff
		} else {
			consumption[consumptionKey] = 0
		}
	}

	ec.PreviousValues = previous
	ec.CurrentValues = current

	return consumption
}

// EnergyPeriodCalculator handles period-based energy calculations.
type EnergyPeriodCalculator struct {
	PeriodMinutes int
}

// NewEnergyPeriodCalculator creates a calculator with the specified period in minutes.
func NewEnergyPeriodCalculator(periodMinutes int) *EnergyPeriodCalculator {
	return &EnergyPeriodCalculator{
		PeriodMinutes: periodMinutes,
	}
}

// Calculate computes the consumption for the given parsed data over a time period.
func (ep *EnergyPeriodCalculator) Calculate(parsed *ParsedData, startTime, endTime time.Time) (float64, error) {
	// Sum all numeric values in the parsed data as a simple consumption metric.
	var total float64
	values := parsed.ScaledValues
	if len(values) == 0 {
		values = parsed.RawValues
	}
	for _, v := range values {
		switch val := v.(type) {
		case float32:
			total += float64(val)
		case float64:
			total += val
		}
	}
	return total, nil
}

// PeriodBoundary holds the start and end times for a calculation period.
type PeriodBoundary struct {
	Start time.Time
	End   time.Time
}

// GetPeriodBoundaries returns the start and end times for the period containing
// the given timestamp.
func (ep *EnergyPeriodCalculator) GetPeriodBoundaries(timestamp time.Time) PeriodBoundary {
	period := time.Duration(ep.PeriodMinutes) * time.Minute
	if period == 0 {
		period = 15 * time.Minute // default 15-minute period
	}

	// Truncate to the start of the period
	periodSeconds := int64(period.Seconds())
	unixTime := timestamp.Unix()
	periodStart := unixTime - (unixTime % periodSeconds)

	start := time.Unix(periodStart, 0)
	end := start.Add(period)

	return PeriodBoundary{
		Start: start,
		End:   end,
	}
}

// CalculatePeriodComparison compares energy consumption between two periods.
// Returns the difference (current - previous) for each parameter, and the
// percentage change.
func (ep *EnergyPeriodCalculator) CalculatePeriodComparison(
	currentPeriod map[string]float64,
	previousPeriod map[string]float64,
) (diff map[string]float64, percentChange map[string]float64) {
	diff = make(map[string]float64)
	percentChange = make(map[string]float64)

	for key, currentVal := range currentPeriod {
		if prevVal, ok := previousPeriod[key]; ok {
			diff[key] = currentVal - prevVal
			if prevVal != 0 {
				percentChange[key] = ((currentVal - prevVal) / prevVal) * 100.0
			} else if currentVal != 0 {
				percentChange[key] = 100.0
			} else {
				percentChange[key] = 0
			}
		} else {
			diff[key] = currentVal
			if currentVal != 0 {
				percentChange[key] = 100.0
			} else {
				percentChange[key] = 0
			}
		}
	}

	return diff, percentChange
}
