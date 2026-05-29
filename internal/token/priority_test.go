package token

import "testing"

func TestPriorityScore(t *testing.T) {
	tests := []struct {
		path     string
		expected int
	}{
		{"README.md", 4},
		{"readme.txt", 4},
		{"README", 4},
		{"main.go", 3},
		{"cmd/app/main.go", 3},
		{"go.mod", 2},
		{"package.json", 2},
		{"Cargo.toml", 2},
		{"requirements.txt", 2},
		{"Makefile", 1},
		{"Dockerfile", 1},
		{"docker-compose.yml", 1},
		{"utils.go", 0},
		{"foo/bar.txt", 0},
	}

	for _, tt := range tests {
		got := PriorityScore(tt.path)
		if got != tt.expected {
			t.Errorf("PriorityScore(%q) = %d, want %d", tt.path, got, tt.expected)
		}
	}
}

func TestSortByPriority(t *testing.T) {
	files := []string{
		"utils.go",
		"README.md",
		"go.mod",
		"main.go",
		"Dockerfile",
	}
	SortByPriority(files)

	if files[0] != "README.md" {
		t.Errorf("expected README.md first, got %s", files[0])
	}
	if files[1] != "main.go" {
		t.Errorf("expected main.go second, got %s", files[1])
	}
	if files[2] != "go.mod" {
		t.Errorf("expected go.mod third, got %s", files[2])
	}
}
