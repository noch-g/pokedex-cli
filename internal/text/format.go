package text

import (
	"fmt"
	"io"
	"strings"
)

func ToBold(s string) string {
	return "\033[1m" + s + "\033[0m"
}

func StartFromClearLine(output io.Writer) {
	fmt.Fprintf(output, "\r\033[K")
}

func LongestWordLength(words []string) int {
	max := 0
	for _, word := range words {
		if len(word) > max {
			max = len(word)
		}
	}
	return max
}

func LongestCommonPrefix(words []string) string {
	if len(words) == 0 {
		return ""
	}
	prefix := words[0]

	for _, word := range words[1:] {
		for len(prefix) > 0 && len(word) > 0 && !strings.HasPrefix(word, prefix) {
			prefix = prefix[:len(prefix)-1]
		}
		if prefix == "" {
			return ""
		}
	}
	return prefix
}

func CleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}
