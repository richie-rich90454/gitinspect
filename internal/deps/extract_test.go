package deps

import "testing"

func TestExtractGoMod(t *testing.T) {
	content := `module example.com/foo

go 1.22

require (
	github.com/bar/baz v1.2.3
	github.com/qux/quux v0.0.0-20240101000000-abcdef123456
)
`
	deps := Extract("go.mod", content)
	if len(deps) != 2 {
		t.Fatalf("expected 2 deps, got %d: %v", len(deps), deps)
	}
	if deps[0] != "github.com/bar/baz@v1.2.3" {
		t.Errorf("first dep = %q, want github.com/bar/baz@v1.2.3", deps[0])
	}
}

func TestExtractPackageJSON(t *testing.T) {
	content := `{"dependencies":{"react":"^18.0.0"},"devDependencies":{"jest":"^29.0.0"}}`
	deps := Extract("package.json", content)
	if len(deps) != 2 {
		t.Fatalf("expected 2 deps, got %d: %v", len(deps), deps)
	}
}

func TestExtractRequirementsTxt(t *testing.T) {
	content := `flask==2.0
requests>=2.28
# comment
numpy
`
	deps := Extract("requirements.txt", content)
	if len(deps) != 3 {
		t.Fatalf("expected 3 deps, got %d: %v", len(deps), deps)
	}
}

func TestExtractCargoToml(t *testing.T) {
	content := `[dependencies]
serde = "1.0"
tokio = { version = "1.0", features = ["full"] }
`
	deps := Extract("Cargo.toml", content)
	if len(deps) != 2 {
		t.Fatalf("expected 2 deps, got %d: %v", len(deps), deps)
	}
}

func TestExtractGemfile(t *testing.T) {
	content := `source "https://rubygems.org"
gem "rails"
gem "puma", "~> 5.0"
`
	deps := Extract("Gemfile", content)
	if len(deps) != 2 {
		t.Fatalf("expected 2 deps, got %d: %v", len(deps), deps)
	}
}

func TestExtractUnknown(t *testing.T) {
	deps := Extract("unknown.txt", "hello")
	if deps != nil {
		t.Errorf("expected nil for unknown file type")
	}
}
