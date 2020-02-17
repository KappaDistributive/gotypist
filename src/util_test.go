package main

import (
	"strings"
	"testing"
)

// TestMin tests whether Min(x, y) returns the correct minima.
func TestMin(t *testing.T) {
	for j := -10; j < 10; j++ {
		for i := -10; i < 10; i++ {
			minimum := Min(i, j)
			if i < j && i != minimum {
				t.Errorf("Min(%d, %d) = %d; want %d", i, j, minimum, i)
			} else if i >= j && j != minimum {
				t.Errorf("Min(%d, %d) = %d; want %d", i, j, minimum, j)

			}
		}
	}
}

// TODO. Verify that CalculateLineBreak indeed works as intended, even in edge cases.
func TestCalculateLineBreak(t *testing.T) {
	testLines := [][]string{
		[]string{"a", "b"},
	}
	testLineLengths := []int{
		1,
	}
	testLineBreaks := []int{
		0,
	}

	for idx, line := range testLines {
		lineBreak := CalculateLineBreak(testLineLengths[idx], line)
		if lineBreak != testLineBreaks[idx] {
			t.Errorf(
				"CalculateLineBreak(%d, %s) = %d; want %d",
				testLineLengths[idx],
				strings.Join(testLines[idx], ", "),
				lineBreak,
				testLineBreaks[idx])
		}
	}
}

// TestDropCursor tests whether DropCursor correctly removes all occurances of `Cursor` in a given word.
func TestDropCursor(t *testing.T) {
	testWords := []string{
		"Hello",
		"Obscure Diagram",
		"Remove " + Cursor,
		"Embedded " + Cursor + " present",
		Cursor,
		Cursor + " " + Cursor,
	}
	cleanedWords := []string{
		"Hello",
		"Obscure Diagram",
		"Remove ",
		"Embedded " + " present",
		"",
		" ",
	}

	for idx, word := range testWords {
		cleaned := DropCursor(word)
		if cleaned != cleanedWords[idx] {
			t.Errorf("dropCusor(%s) = %s; want %s", word, cleaned, cleanedWords[idx])
		}
	}
}
