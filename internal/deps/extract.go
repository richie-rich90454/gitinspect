package deps

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

// Dependency represents a single dependency
type Dependency struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// Extract extracts dependencies from various manifest files
func Extract(path string, content string) ([]Dependency, error) {
	base := filepath.Base(path)

	switch base {
	case "go.mod":
		return extractGoMod(content)
	case "package.json":
		return extractPackageJSON(content)
	case "Cargo.toml":
		return extractCargoToml(content)
	case "requirements.txt":
		return extractRequirementsTxt(content)
	case "Gemfile":
		return extractGemfile(content)
	default:
		return nil, nil
	}
}

func extractGoMod(content string) ([]Dependency, error) {
	var deps []Dependency
	lines := strings.Split(content, "\n")
	inRequire := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "require") {
			inRequire = true
			continue
		}
		if strings.HasPrefix(line, ")") {
			inRequire = false
			continue
		}
		if inRequire || (!strings.HasPrefix(line, "module") && !strings.HasPrefix(line, "go ") && !strings.HasPrefix(line, "replace")) {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				deps = append(deps, Dependency{
					Name:    parts[0],
					Version: parts[1],
				})
			}
		}
	}
	return deps, nil
}

func extractPackageJSON(content string) ([]Dependency, error) {
	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal([]byte(content), &pkg); err != nil {
		return nil, err
	}

	var deps []Dependency
	for name, version := range pkg.Dependencies {
		deps = append(deps, Dependency{Name: name, Version: version})
	}
	for name, version := range pkg.DevDependencies {
		deps = append(deps, Dependency{Name: name, Version: version})
	}
	return deps, nil
}

func extractCargoToml(content string) ([]Dependency, error) {
	var deps []Dependency
	lines := strings.Split(content, "\n")
	inDeps := false
	inDevDeps := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.ToLower(line) == "[dependencies]" {
			inDeps = true
			inDevDeps = false
			continue
		}
		if strings.ToLower(line) == "[dev-dependencies]" {
			inDeps = false
			inDevDeps = true
			continue
		}
		if strings.HasPrefix(line, "[") {
			inDeps = false
			inDevDeps = false
			continue
		}
		if inDeps || inDevDeps {
			if idx := strings.Index(line, "="); idx != -1 {
				name := strings.TrimSpace(line[:idx])
				versionPart := strings.TrimSpace(line[idx+1:])
				version := strings.Trim(versionPart, `"'`)
				deps = append(deps, Dependency{Name: name, Version: version})
			}
		}
	}
	return deps, nil
}

func extractRequirementsTxt(content string) ([]Dependency, error) {
	var deps []Dependency
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Handle name==version or name~=version or name>=version etc.
		var name, version string
		if idx := strings.IndexAny(line, "=~<>"); idx != -1 {
			name = strings.TrimSpace(line[:idx])
			version = strings.TrimSpace(line[idx:])
		} else {
			name = line
		}
		deps = append(deps, Dependency{Name: name, Version: version})
	}
	return deps, nil
}

func extractGemfile(content string) ([]Dependency, error) {
	var deps []Dependency
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "gem ") {
			parts := strings.Fields(line[4:])
			if len(parts) > 0 {
				name := strings.Trim(parts[0], `'"`)
				version := ""
				if len(parts) > 1 {
					version = strings.Trim(parts[1], `'"`)
				}
				deps = append(deps, Dependency{Name: name, Version: version})
			}
		}
	}
	return deps, nil
}

// ToString converts dependencies to a string
func ToString(deps []Dependency) string {
	if len(deps) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, d := range deps {
		if d.Version != "" {
			sb.WriteString(fmt.Sprintf("- %s: %s\n", d.Name, d.Version))
		} else {
			sb.WriteString(fmt.Sprintf("- %s\n", d.Name))
		}
	}
	return sb.String()
}
