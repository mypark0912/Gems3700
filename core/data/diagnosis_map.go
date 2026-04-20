package data

import "strings"

// Diagnosis status constants
const (
	StatusNoData  = 0
	StatusOK      = 1
	StatusWarning = 2
	StatusInspect = 3
	StatusRepair  = 4
)

// MatchedParameter represents a parameter that triggered an alarm.
type MatchedParameter struct {
	Name            string `json:"name"`
	ConfiguredLevel int    `json:"configured_level"`
	ActualStatus    int    `json:"actual_status"`
	Triggered       bool   `json:"triggered"`
}

// DiagnosisResult holds the result of a diagnosis evaluation.
type DiagnosisResult struct {
	FinalStatus       int                `json:"final_status"`
	StatusName        string             `json:"status_name"`
	MatchedParameters []MatchedParameter `json:"matched_parameters"`
	TotalConfigured   int                `json:"total_configured"`
	TotalMatched      int                `json:"total_matched"`
}

// AlarmStatusMatcher determines the final diagnosis status by matching
// status_info config against API bar_graph responses.
type AlarmStatusMatcher struct{}

// NewAlarmStatusMatcher creates a new AlarmStatusMatcher.
func NewAlarmStatusMatcher() *AlarmStatusMatcher {
	return &AlarmStatusMatcher{}
}

// normalizeName removes spaces and lowercases the name.
func (a *AlarmStatusMatcher) normalizeName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", ""))
}

// ParseConfigList converts a list of config items (each with "name" and "level")
// into a normalized map of name -> level.
func (a *AlarmStatusMatcher) ParseConfigList(configList []map[string]interface{}) map[string]int {
	config := make(map[string]int)
	for _, item := range configList {
		name, ok := item["name"].(string)
		if !ok {
			continue
		}
		level := 0
		switch v := item["level"].(type) {
		case int:
			level = v
		case float64:
			level = int(v)
		}
		config[a.normalizeName(name)] = level
	}
	return config
}

// ParseAPIResponse converts a bar_graph API response list into a normalized
// map of name -> status.
func (a *AlarmStatusMatcher) ParseAPIResponse(barGraph []map[string]interface{}) map[string]int {
	result := make(map[string]int)
	for _, item := range barGraph {
		name, ok := item["Name"].(string)
		if !ok {
			continue
		}
		status := 0
		switch v := item["Status"].(type) {
		case int:
			status = v
		case float64:
			status = int(v)
		}
		result[a.normalizeName(name)] = status
	}
	return result
}

// CalculateFinalStatus compares config levels against actual response statuses.
// Returns the maximum triggered status and a list of matched parameters.
func (a *AlarmStatusMatcher) CalculateFinalStatus(
	statusConfig map[string]int,
	apiResponse map[string]int,
) (int, []MatchedParameter) {
	var matched []MatchedParameter
	maxStatus := 0

	for name, configuredLevel := range statusConfig {
		actualStatus, exists := apiResponse[name]
		if !exists {
			continue
		}
		if actualStatus >= configuredLevel {
			matched = append(matched, MatchedParameter{
				Name:            name,
				ConfiguredLevel: configuredLevel,
				ActualStatus:    actualStatus,
				Triggered:       true,
			})
			if actualStatus > maxStatus {
				maxStatus = actualStatus
			}
		}
	}

	// If nothing matched, return OK (1)
	if len(matched) == 0 {
		return StatusOK, matched
	}

	return maxStatus, matched
}

// GetStatusName converts a status code to its string representation.
func GetStatusName(status int) string {
	switch status {
	case StatusNoData:
		return "NoData"
	case StatusOK:
		return "OK"
	case StatusWarning:
		return "Warning"
	case StatusInspect:
		return "Inspect"
	case StatusRepair:
		return "Repair"
	default:
		return "Unknown"
	}
}

// Diagnose evaluates status_info config against bar_graph API response and
// returns the overall diagnosis result.
func (a *AlarmStatusMatcher) Diagnose(
	statusInfo []map[string]interface{},
	barGraph []map[string]interface{},
) DiagnosisResult {
	// Parse config
	config := a.ParseConfigList(statusInfo)

	// Parse API response
	response := a.ParseAPIResponse(barGraph)

	// If all response statuses are 0, return NoData
	if len(response) > 0 {
		allZero := true
		for _, status := range response {
			if status != 0 {
				allZero = false
				break
			}
		}
		if allZero {
			return DiagnosisResult{
				FinalStatus:       StatusNoData,
				StatusName:        GetStatusName(StatusNoData),
				MatchedParameters: []MatchedParameter{},
				TotalConfigured:   len(config),
				TotalMatched:      0,
			}
		}
	}

	// Calculate final status
	finalStatus, matched := a.CalculateFinalStatus(config, response)

	return DiagnosisResult{
		FinalStatus:       finalStatus,
		StatusName:        GetStatusName(finalStatus),
		MatchedParameters: matched,
		TotalConfigured:   len(config),
		TotalMatched:      len(matched),
	}
}
