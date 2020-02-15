package main

import "strings"

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func CalculateLineBreak(lineLength int, words []string) int {
	if len(words) == 0 {
		return 0
	}
	length := 0
	for {
		if line := strings.Join(words[:Min(len(words), length+1)], " "); length < len(words) && len(line) < lineLength {
			length += 1
		} else {
			return length
		}

	}
}

func dropCursor(word string) string {
	return strings.ReplaceAll(word, Cursor, "")
}
