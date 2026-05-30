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
		got, err := MatchAny(tt.path, tt.includes, tt.excludes)
		if err != nil {
			t.Errorf("MatchAny(%q, %v, %v) returned unexpected error: %v", tt.path, tt.includes, tt.excludes, err)
			continue
		}
		if got != tt.expected {
			t.Errorf("MatchAny(%q, %v, %v) = %v, want %v", tt.path, tt.includes, tt.excludes, got, tt.expected)
		}
	}
}
