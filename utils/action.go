package utils

import (
	"strings"
)

// GetActionInputAsSlice handle the action multiline input as slice
func GetActionInputAsSlice(input string) []string {
	result := make([]string, 0, 5)
	s := strings.Split(input, "\n")
	for _, i := range s {
		if c := strings.TrimSpace(i); c != "" {
			result = append(result, c)
		}
	}
	return result
}
