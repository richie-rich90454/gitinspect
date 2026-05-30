package token

import "testing"

func TestEstimate(t *testing.T) {
	tests := []struct {
		content  string
		expected int
	}{
		{"", 0},
		{"a", 1},
		{"abcd", 1},
		{"abcde", 2},
		{"abcdefgh", 2},
		{"abcdefghi", 3},
	}

	for _, tt := range tests {
		got := Estimate(tt.content)
		if got != tt.expected {
			t.Errorf("Estimate(%q) = %d, want %d", tt.content, got, tt.expected)
		}
	}
}

func TestEstimateSize(t *testing.T) {
	tests := []struct {
		size     int
		expected int
	}{
		{0, 0},
		{1, 1},
		{4, 1},
		{5, 2},
		{8, 2},
		{9, 3},
		{100, 25},
	}

	for _, tt := range tests {
		got := EstimateSize(tt.size)
		if got != tt.expected {
			t.Errorf("EstimateSize(%d) = %d, want %d", tt.size, got, tt.expected)
		}
	}
}

func TestTruncateMultiByteUTF8(t *testing.T) {
	head := "你好世界"
	for len(head) < 1000 {
		head += "你好世界"
	}
	tail := "再见世界"
	for len(tail) < 500 {
		tail += "再见世界"
	}
	content := head + "padding" + tail

	result := Truncate(content)

	for _, r := range result {
		if r == '\ufffd' {
			t.Error("Truncate produced invalid UTF-8 (replacement character found)")
			break
		}
	}
}

func TestTruncateShortMultiByte(t *testing.T) {
	content := "hello"
	result := Truncate(content)
	if result != content {
		t.Error("Truncate should not modify short multi-byte-safe content")
	}
}
