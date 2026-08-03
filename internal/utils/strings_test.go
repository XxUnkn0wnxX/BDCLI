package utils

import (
	"reflect"
	"testing"
)

func TestFormatVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Trims and preserves leading v",
			input:    " v1.2.3 ",
			expected: "v1.2.3",
		},
		{
			name:     "Adds leading v",
			input:    "1.2.3",
			expected: "v1.2.3",
		},
		{
			name:     "Keeps v0.0.0",
			input:    "v0.0.0",
			expected: "v0.0.0",
		},
		{
			name:     "Empty defaults",
			input:    "",
			expected: "v0.0.0",
		},
		{
			name:     "Whitespace defaults",
			input:    "   ",
			expected: "v0.0.0",
		},
		{
			name:     "Bare v defaults",
			input:    "v",
			expected: "v0.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatVersion(tt.input)
			if result != tt.expected {
				t.Errorf("FormatVersion(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSplitVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Standard version",
			input:    "1.0.156",
			expected: []string{"1", "0", "156"},
		},
		{
			name:     "Skips empty segments",
			input:    "1..2",
			expected: []string{"1", "2"},
		},
		{
			name:     "Ignores non-digits",
			input:    "v1.2.3-beta",
			expected: []string{"1", "2", "3"},
		},
		{
			name:     "No digits",
			input:    "beta",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitVersion(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SplitVersion(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		{
			name:     "Equal with v prefix",
			v1:       "v1.2.3",
			v2:       "1.2.3",
			expected: 0,
		},
		{
			name:     "Less than",
			v1:       "1.2.3",
			v2:       "1.2.4",
			expected: -1,
		},
		{
			name:     "Greater than",
			v1:       "2.0.0",
			v2:       "1.9.9",
			expected: 1,
		},
		{
			name:     "Missing parts default to zero",
			v1:       "1.2",
			v2:       "1.2.0",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareVersions(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("CompareVersions(%q, %q) = %d, expected %d", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}
