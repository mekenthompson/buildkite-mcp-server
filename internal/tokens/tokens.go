package tokens

import "strings"

// EstimateTokens returns an estimate of the number of tokens in the given text.
func EstimateTokens(text string) int {
	words := strings.Fields(text)
	tokenCount := 0

	for _, word := range words {
		// Simple heuristic: longer words typically split into more tokens
		wordLen := len([]rune(word))
		if wordLen <= 4 {
			tokenCount += 1
		} else if wordLen <= 8 {
			tokenCount += 2
		} else {
			tokenCount += (wordLen + 3) / 4
		}
	}

	return tokenCount
}
