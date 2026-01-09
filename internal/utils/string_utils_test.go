package utils

import "testing"

func TestSubstr(t *testing.T) {
	tests := []struct {
		input    string
		start    int
		length   int
		expected string
	}{
		{"hello", 0, 5, "hello"},
		{"hello", 0, 2, "he"},
		{"hello", 2, 2, "ll"},
		{"hello", 2, 10, "llo"},
		{"hello", 10, 2, ""},
		{"😊🌍", 0, 1, "😊"},
		{"😊🌍", 1, 1, "🌍"},
	}

	for _, tt := range tests {
		result := Substr(tt.input, tt.start, tt.length)
		if result != tt.expected {
			t.Errorf("Substr(%q, %d, %d) = %q; expected %q", tt.input, tt.start, tt.length, result, tt.expected)
		}
	}
}

func TestIsNullOrEmpty(t *testing.T) {
	if !IsNullOrEmpty("") {
		t.Error("Expected IsNullOrEmpty(\"\") to be true")
	}
	if IsNullOrEmpty("a") {
		t.Error("Expected IsNullOrEmpty(\"a\") to be false")
	}
}
