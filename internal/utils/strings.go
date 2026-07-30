package utils

import (
	"fmt"
	"net/url"
	"strings"
)

// IsURL checks if a string is a valid URL
func IsURL(input string) bool {
	parsed, err := url.Parse(input)
	return err == nil && parsed.Scheme != "" && parsed.Host != ""
}

// FormatVersion normalizes a version string with a single leading 'v'.
func FormatVersion(version string) string {
	trimmed := strings.TrimSpace(version)
	trimmed = strings.TrimPrefix(trimmed, "v")
	if trimmed == "" {
		return "v0.0.0"
	}
	return "v" + trimmed
}

// CompareVersions compares two semantic versions (e.g., "1.0.156" vs "1.0.157")
// Returns -1 if v1 < v2, 0 if equal, 1 if v1 > v2
func CompareVersions(v1, v2 string) int {
	// Strip 'v' prefix if present

	if len(v1) > 0 && v1[0] == 'v' {
		v1 = v1[1:]
	}
	if len(v2) > 0 && v2[0] == 'v' {
		v2 = v2[1:]
	}

	// Parse into version parts
	parts1 := SplitVersion(v1)
	parts2 := SplitVersion(v2)

	// Compare each part
	maxLen := max(len(parts2), len(parts1))

	for i := range maxLen {
		var p1, p2 int

		if i < len(parts1) {
			fmt.Sscanf(parts1[i], "%d", &p1)
		}
		if i < len(parts2) {
			fmt.Sscanf(parts2[i], "%d", &p2)
		}

		if p1 < p2 {
			return -1
		} else if p1 > p2 {
			return 1
		}
	}

	return 0
}

// SplitVersion splits a version string into parts (e.g., "1.0.156" -> ["1", "0", "156"])
func SplitVersion(v string) []string {
	var parts []string
	var current string
	for i := 0; i < len(v); i++ {
		if v[i] == '.' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else if v[i] >= '0' && v[i] <= '9' {
			current += string(v[i])
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}