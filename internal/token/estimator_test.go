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

func TestTruncate(t *testing.T) {
	short := "hello world"
	result := Truncate(short)
	if result != short {
		t.Errorf("Truncate should not modify short content")
	}

	long := make([]byte, 2000)
	for i := range long {
		long[i] = 'a' + byte(i%26)
	}
	content := string(long)
	result = Truncate(content)
	if len(result) >= len(content) {
		t.Errorf("Truncate should shorten long content")
	}
	if len(result) < 1000 {
		t.Errorf("Truncate should keep at least 1000 chars from head, got %d", len(result))
	}
}
