package main

import "strings"

// Min returns the minimum of integers x and y.
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// CalculateLineBreak calculates the correct index for a line break. 
func CalculateLineBreak(lineLength int, words []string) int {
	if len(words) == 0 {
		return 0
	}
	length := 0
	for {
		if line := strings.Join(words[:Min(len(words), length+1)], " "); length < len(words) && len(line) < lineLength {
			length++
		} else {
			return length
		}
	}
}

// DropCursor removes all occurances of `Cursor` in a given word.
func DropCursor(word string) string {
	return strings.ReplaceAll(word, Cursor, "")
}
