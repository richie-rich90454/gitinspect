package repo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCacheSetGet(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gitinspect-cache-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	c := &Cache{Dir: tmpDir}
	key := "testkey123"

	data := []byte("hello cache")
	if err := c.Set(key, data); err != nil {
		t.Fatal(err)
	}

	got, ok := c.Get(key)
	if !ok {
		t.Fatal("cache miss after set")
	}
	if string(got) != string(data) {
		t.Errorf("got %q, want %q", got, data)
	}
}

func TestCacheMiss(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gitinspect-cache-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	c := &Cache{Dir: tmpDir}
	_, ok := c.Get("nonexistent")
	if ok {
		t.Error("expected cache miss for nonexistent key")
	}
}

func TestGenerateKey(t *testing.T) {
	k1 := GenerateKey("https://github.com/foo/bar.git", "abc123")
	k2 := GenerateKey("https://github.com/foo/bar.git", "abc123")
	k3 := GenerateKey("https://github.com/foo/bar.git", "def456")

	if k1 != k2 {
		t.Error("same inputs should produce same key")
	}
	if k1 == k3 {
		t.Error("different inputs should produce different keys")
	}
}

func TestStripCommentsDoubleSlash(t *testing.T) {
	input := "// comment\npackage main\n// another comment\nfunc foo() {}\n"
	got := StripComments(input, ".go")
	if containsLine(got, "// comment") {
		t.Error("expected // comment to be stripped for .go")
	}
	if containsLine(got, "// another comment") {
		t.Error("expected // another comment to be stripped for .go")
	}
	if !containsLine(got, "package main") {
		t.Error("expected 'package main' to be preserved")
	}
	if !containsLine(got, "func foo() {}") {
		t.Error("expected 'func foo() {}' to be preserved")
	}
}

func TestStripCommentsHash(t *testing.T) {
	input := "# comment\nprint('hello')\n# another\nx = 1\n"
	got := StripComments(input, ".py")
	if containsLine(got, "# comment") {
		t.Error("expected # comment to be stripped for .py")
	}
	if !containsLine(got, "print('hello')") {
		t.Error("expected code line to be preserved")
	}
}

func TestStripCommentsDoubleDash(t *testing.T) {
	input := "-- comment\nSELECT 1;\n-- another\n"
	got := StripComments(input, ".sql")
	if containsLine(got, "-- comment") {
		t.Error("expected -- comment to be stripped for .sql")
	}
	if !containsLine(got, "SELECT 1;") {
		t.Error("expected SQL to be preserved")
	}
}

func TestStripCommentsSemicolon(t *testing.T) {
	input := "; comment\nmov eax, 1\n; another\n"
	got := StripComments(input, ".asm")
	if containsLine(got, "; comment") {
		t.Error("expected ; comment to be stripped for .asm")
	}
	if !containsLine(got, "mov eax, 1") {
		t.Error("expected assembly to be preserved")
	}
}

func TestStripCommentsUnknownExt(t *testing.T) {
	input := "// comment\ncode line\n"
	got := StripComments(input, ".xyz")
	if got != input {
		t.Error("expected no stripping for unknown extension")
	}
}

func TestStripCommentsOnlyComments(t *testing.T) {
	input := "// comment1\n// comment2\n// comment3\n"
	got := StripComments(input, ".go")
	if got != "" {
		t.Errorf("expected empty result for comments-only file, got %q", got)
	}
}

func TestStripCommentsMixedContent(t *testing.T) {
	input := "// header\npackage main\n\nimport \"fmt\"\n// TODO: fix\nfunc main() {}\n"
	got := StripComments(input, ".go")
	if containsLine(got, "// header") {
		t.Error("expected // header to be stripped")
	}
	if containsLine(got, "// TODO: fix") {
		t.Error("expected // TODO: fix to be stripped")
	}
	if !containsLine(got, "package main") {
		t.Error("expected 'package main' to be preserved")
	}
	if !containsLine(got, "import \"fmt\"") {
		t.Error("expected import to be preserved")
	}
	if !containsLine(got, "func main() {}") {
		t.Error("expected func main to be preserved")
	}
}

func containsLine(s, line string) bool {
	for _, l := range splitLines(s) {
		if l == line {
			return true
		}
	}
	return false
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		result = append(result, s[start:])
	}
	return result
}

func TestIsBinaryNullByte(t *testing.T) {
	data := []byte("hello\x00world")
	if !isBinary(data) {
		t.Error("expected data with null byte to be binary")
	}
}

func TestIsBinaryEmptyContent(t *testing.T) {
	if isBinary([]byte{}) {
		t.Error("expected empty content to not be binary")
	}
}

func TestIsBinaryLargeBinaryContent(t *testing.T) {
	data := make([]byte, 16000)
	for i := range data {
		if i == 4000 {
			data[i] = 0
		} else {
			data[i] = 'a'
		}
	}
	if !isBinary(data) {
		t.Error("expected large content with null byte to be binary")
	}
}

func TestIsBinaryPlainText(t *testing.T) {
	data := []byte("hello world\nthis is plain text\n")
	if isBinary(data) {
		t.Error("expected plain text to not be binary")
	}
}

func TestReadLocalRepo(t *testing.T) {
	repoPath := filepath.Join("..", "..", "testdata", "sample-repo")
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		t.Fatal(err)
	}

	entries, err := ReadLocalRepo(absPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one file entry")
	}

	found := false
	for _, e := range entries {
		if e.Path == "main.go" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected main.go in entries")
	}
}
