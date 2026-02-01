package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode/utf8"
)

func main() {
	// Read input from standard input
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := scanner.Text()

		// Call the ReverseString function
		output := ReverseString(input)

		// Print the result
		fmt.Println(output)
	}
}

// ReverseString returns the reversed string of s.
func ReverseString(s string) string {
	res := make([]byte, len(s))
	start := 0
	end := len(s)

	for end > 0 {
		r, size := utf8.DecodeLastRuneInString(s[:end])

		utf8.EncodeRune(res[start:], r)

		start += size
		end -= size
	}

	return string(res)
}
