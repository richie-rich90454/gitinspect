// Package token provides token estimation and file prioritization.
package token

import "unicode/utf8"

// Estimate returns an approximate token count for the given content.
func Estimate(content string) int {
	return (len(content) + 3) / 4
}

// Truncate shortens content by keeping the head and tail with a truncation marker.
func Truncate(content string) string {
	const headSize = 1000
	const tailSize = 500

	if len(content) <= headSize+tailSize+50 {
		return content
	}

	head := content[:headSize]
	for !utf8.ValidString(head) {
		head = head[:len(head)-1]
	}

	tailStart := len(content) - tailSize
	for tailStart < len(content) && !utf8.RuneStart(content[tailStart]) {
		tailStart++
	}
	tail := content[tailStart:]

	return head + "\n... [truncated] ...\n" + tail
}
