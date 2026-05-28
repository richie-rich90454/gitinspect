package token

// Estimate estimates the number of tokens for a given string.
// Uses a simple heuristic: len(content) / 4
func Estimate(content string) int {
	return len(content) / 4
}

// Truncate truncates content to fit within maxTokens, keeping first and last parts.
func Truncate(content string, maxTokens int) (string, bool) {
	estimated := Estimate(content)
	if estimated <= maxTokens {
		return content, false
	}

	// Keep first and last parts
	chars := maxTokens * 4 // 4 chars per token
	if chars <= 0 {
		return "", true
	}

	keepChars := chars / 2
	if keepChars*2 > len(content) {
		return content, true
	}

	return content[:keepChars] + "\n... [TRUNCATED] ...\n" + content[len(content)-keepChars:], true
}
