package output

import (
	"fmt"
	"strings"
)

// ToText converts the snapshot to plain text
func ToText(snapshot Snapshot) (string, error) {
	var sb strings.Builder

	if snapshot.RepoURL != "" {
		sb.WriteString(fmt.Sprintf("# Repository: %s\n", snapshot.RepoURL))
		if snapshot.CommitHash != "" {
			sb.WriteString(fmt.Sprintf("# Commit: %s\n", snapshot.CommitHash))
		}
		sb.WriteString("\n")
	}

	if len(snapshot.Dependencies) > 0 {
		sb.WriteString("## Dependencies:\n")
		for _, d := range snapshot.Dependencies {
			if d.Version != "" {
				sb.WriteString(fmt.Sprintf("- %s: %s\n", d.Name, d.Version))
			} else {
				sb.WriteString(fmt.Sprintf("- %s\n", d.Name))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("## Files (%d total, %d tokens):\n", len(snapshot.Files), snapshot.TotalTokens))
	sb.WriteString("---\n\n")

	for _, f := range snapshot.Files {
		sb.WriteString(fmt.Sprintf("### %s\n", f.Path))
		if f.Truncated {
			sb.WriteString(fmt.Sprintf("(Truncated, ~%d tokens)\n", f.Tokens))
		} else {
			sb.WriteString(fmt.Sprintf("(~%d tokens)\n", f.Tokens))
		}
		sb.WriteString("```\n")
		sb.WriteString(f.Content)
		sb.WriteString("\n```\n\n")
	}

	return sb.String(), nil
}
