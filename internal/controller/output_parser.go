package controller

import "strings"

const (
	outputStartMarker = "---AXON_OUTPUTS_START---"
	outputEndMarker   = "---AXON_OUTPUTS_END---"
)

// ParseOutputs extracts output lines from log data between the
// ---AXON_OUTPUTS_START--- and ---AXON_OUTPUTS_END--- markers.
func ParseOutputs(logData string) []string {
	startIdx := strings.Index(logData, outputStartMarker)
	if startIdx == -1 {
		return nil
	}
	endIdx := strings.Index(logData, outputEndMarker)
	if endIdx == -1 || endIdx <= startIdx {
		return nil
	}

	between := logData[startIdx+len(outputStartMarker) : endIdx]
	between = strings.TrimSpace(between)
	if between == "" {
		return nil
	}

	lines := strings.Split(between, "\n")
	var result []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
