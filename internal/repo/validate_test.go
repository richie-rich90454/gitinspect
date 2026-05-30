package repo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateRemoteURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"https URL", "https://github.com/user/repo.git", false},
		{"http URL", "http://github.com/user/repo.git", false},
		{"git protocol", "git://github.com/user/repo.git", false},
		{"ssh protocol", "ssh://git@github.com/user/repo.git", false},
		{"git@ shorthand", "git@github.com:user/repo.git", false},
		{"empty URL", "", true},
		{"flag injection", "-flag", true},
		{"local path", "/home/user/repo", true},
		{"relative path", "./repo", true},
		{"ftp URL", "ftp://example.com/repo.git", true},
		{"file URL", "file:///home/user/repo", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRemoteURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRemoteURL(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
		})
	}
}

func TestValidateRemoteURLErrorMessages(t *testing.T) {
	err := ValidateRemoteURL("")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
	if err.Error() != "remote URL must not be empty" {
		t.Errorf("unexpected error message: %q", err.Error())
	}

	err = ValidateRemoteURL("-flag")
	if err == nil {
		t.Fatal("expected error for flag injection")
	}
	if err.Error() != "remote URL must not start with a dash" {
		t.Errorf("unexpected error message: %q", err.Error())
	}

	err = ValidateRemoteURL("ftp://example.com")
	if err == nil {
		t.Fatal("expected error for unsupported scheme")
	}
	if err.Error() != "remote URL must start with one of: https://, http://, git://, ssh://, git@" {
		t.Errorf("unexpected error message: %q", err.Error())
	}
}

func TestValidateLocalPathStrict(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gitinspect-validate-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	tmpFile := filepath.Join(tmpDir, "file.txt")
	if err := os.WriteFile(tmpFile, []byte("hello"), 0o644); err != nil { //nolint:gosec
		t.Fatal(err)
	}

	t.Run("valid directory", func(t *testing.T) {
		if err := ValidateLocalPath(tmpDir, true); err != nil {
			t.Errorf("ValidateLocalPath(%q, true) unexpected error: %v", tmpDir, err)
		}
	})

	t.Run("empty path", func(t *testing.T) {
		err := ValidateLocalPath("", true)
		if err == nil {
			t.Error("expected error for empty path")
		}
	})

	t.Run("path traversal", func(t *testing.T) {
		err := ValidateLocalPath(filepath.Join(tmpDir, "..", "something"), true)
		if err == nil {
			t.Error("expected error for path with ..")
		}
	})

	t.Run("root path", func(t *testing.T) {
		err := ValidateLocalPath("/", true)
		if err == nil {
			t.Error("expected error for root path")
		}
	})

	t.Run("nonexistent path", func(t *testing.T) {
		err := ValidateLocalPath(filepath.Join(tmpDir, "nonexistent"), true)
		if err == nil {
			t.Error("expected error for nonexistent path")
		}
	})

	t.Run("file not directory", func(t *testing.T) {
		err := ValidateLocalPath(tmpFile, true)
		if err == nil {
			t.Error("expected error for file path")
		}
	})
}

func TestValidateLocalPathNonStrict(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gitinspect-validate-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	tmpFile := filepath.Join(tmpDir, "file.txt")
	if err := os.WriteFile(tmpFile, []byte("hello"), 0o644); err != nil { //nolint:gosec
		t.Fatal(err)
	}

	t.Run("valid directory", func(t *testing.T) {
		if err := ValidateLocalPath(tmpDir, false); err != nil {
			t.Errorf("ValidateLocalPath(%q, false) unexpected error: %v", tmpDir, err)
		}
	})

	t.Run("path traversal allowed", func(t *testing.T) {
		parentWithDotDot := filepath.Join(tmpDir, "..")
		if err := ValidateLocalPath(parentWithDotDot, false); err != nil {
			t.Errorf("ValidateLocalPath with .. should be allowed in non-strict mode: %v", err)
		}
	})

	t.Run("empty path", func(t *testing.T) {
		err := ValidateLocalPath("", false)
		if err == nil {
			t.Error("expected error for empty path")
		}
	})

	t.Run("nonexistent path", func(t *testing.T) {
		err := ValidateLocalPath(filepath.Join(tmpDir, "nonexistent"), false)
		if err == nil {
			t.Error("expected error for nonexistent path")
		}
	})

	t.Run("file not directory", func(t *testing.T) {
		err := ValidateLocalPath(tmpFile, false)
		if err == nil {
			t.Error("expected error for file path")
		}
	})
}
