package filter

import "testing"

func TestMatchAny(t *testing.T) {
	tests := []struct {
		path     string
		includes []string
		excludes []string
		expected bool
	}{
		{"foo.go", nil, nil, true},
		{"foo.go", []string{"**/*.go"}, nil, true},
		{"foo.go", []string{"**/*.go"}, []string{"vendor/**"}, true},
		{"vendor/foo.go", []string{"**/*.go"}, []string{"vendor/**"}, false},
		{"foo.txt", []string{"**/*.go"}, nil, false},
		{"a/b/c.go", []string{"**/*.go"}, nil, true},
		{"a/b/c.go", nil, []string{"**/vendor/**"}, true},
	}

	for _, tt := range tests {
		got := MatchAny(tt.path, tt.includes, tt.excludes)
		if got != tt.expected {
			t.Errorf("MatchAny(%q, %v, %v) = %v, want %v", tt.path, tt.includes, tt.excludes, got, tt.expected)
		}
	}
}
