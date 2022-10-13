package utils

import "strings"

// ParseUpload parser of uploaded file
func ParseUpload(data []byte) []string {
	lines := strings.Split(string(data), "\n")
	return lines
}
