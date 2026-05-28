package output

import (
	"gopkg.in/yaml.v3"
)

// ToYAML converts the snapshot to YAML
func ToYAML(snapshot Snapshot) (string, error) {
	data, err := yaml.Marshal(snapshot)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
