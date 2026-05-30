// Package deps extracts dependency lists from common manifest files.
package deps

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Extract parses dependencies from a manifest file based on its name.
func Extract(filePath, content string) []string {
	base := filepath.Base(filePath)

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
		return nil
	}
}

func extractGoMod(content string) []string {
	var deps []string
	lines := strings.Split(content, "\n")
	inRequire := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "require (") {
			inRequire = true
			continue
		}
		if strings.HasPrefix(line, ")") {
			inRequire = false
			continue
		}
		if inRequire {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				deps = append(deps, fmt.Sprintf("%s@%s", parts[0], parts[1]))
			}
		} else if strings.HasPrefix(line, "require ") {
			parts := strings.Fields(line[8:])
			if len(parts) >= 2 {
				deps = append(deps, fmt.Sprintf("%s@%s", parts[0], parts[1]))
			}
		}
	}
	return deps
}

func extractPackageJSON(content string) []string {
	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal([]byte(content), &pkg); err != nil {
		return nil
	}

	var deps []string
	for name, version := range pkg.Dependencies {
		deps = append(deps, fmt.Sprintf("%s@%s", name, version))
	}
	for name, version := range pkg.DevDependencies {
		deps = append(deps, fmt.Sprintf("%s@%s", name, version))
	}
	return deps
}

func extractCargoToml(content string) []string {
	var cargo struct {
		Dependencies map[string]interface{} `toml:"dependencies"`
	}
	if err := toml.Unmarshal([]byte(content), &cargo); err != nil {
		return nil
	}

	var deps []string
	for name, val := range cargo.Dependencies {
		switch v := val.(type) {
		case string:
			deps = append(deps, fmt.Sprintf("%s@%s", name, v))
		case map[string]interface{}:
			if ver, ok := v["version"].(string); ok {
				deps = append(deps, fmt.Sprintf("%s@%s", name, ver))
			} else {
				deps = append(deps, name)
			}
		default:
			deps = append(deps, name)
		}
	}
	return deps
}

func extractRequirementsTxt(content string) []string {
	var deps []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "-") {
			continue
		}
		deps = append(deps, line)
	}
	return deps
}

func extractGemfile(content string) []string {
	var deps []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "gem ") {
			rest := strings.TrimSpace(line[4:])
			if len(rest) >= 2 && (rest[0] == '\'' || rest[0] == '"') {
				end := strings.IndexByte(rest[1:], rest[0])
				if end != -1 {
					deps = append(deps, rest[1:end+1])
				}
			} else {
				parts := strings.Fields(rest)
				if len(parts) > 0 {
					deps = append(deps, parts[0])
				}
			}
		}
	}
	return deps
}
