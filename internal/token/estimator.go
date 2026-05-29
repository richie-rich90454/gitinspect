package token

func Estimate(content string) int {
	return (len(content) + 3) / 4
}

func Truncate(content string) string {
	const headSize = 1000
	const tailSize = 500

	if len(content) <= headSize+tailSize+50 {
		return content
	}

	head := content[:headSize]
	tail := content[len(content)-tailSize:]
	return head + "\n... [truncated] ...\n" + tail
}
