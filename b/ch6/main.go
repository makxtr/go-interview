package main

import "strings"

// CountWordFrequency Input: "The quick brown fox jumps over the lazy dog."
// Output: map[string]int{"the": 2, "quick": 1, "brown": 1, "fox": 1, "jumps": 1, "over": 1, "lazy": 1, "dog": 1}
func CountWordFrequency(text string) map[string]int {
	text = strings.ReplaceAll(strings.ToLower(text), "'", "")
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9')
	})

	res := make(map[string]int, len(words))
	for _, word := range words {
		res[word]++
	}
	return res
}
