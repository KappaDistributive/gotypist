package main

import (
	"fmt"
	"strings"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// Typing implements Viewport
type Typing struct {
	title             string
	input             *widgets.Paragraph
	display           *widgets.Paragraph
	words             []string
	cursorPos         int
	start             int
	newline           int
	end               int
	totalCharacters   int
	started           bool
	startTime         time.Time
	correctCharacters int
}

func (self Typing) Handler(e <-chan ui.Event) (Viewport, error) {
	event := <-e
	text := dropCursor(self.input.Text)
	length := len(text)

	switch event.ID {
	case "<C-c>":
		return self, Quit{}
	case "<Escape>":
		return createSelection(), nil
	// TODO: replace ad-hoc text handling
	case "<Space>":
		updateCpm(text, &self)
		checkWord(text, self.cursorPos, &self.words)
		self.cursorPos += 1
		if self.cursorPos == len(self.words) {
			// end the game
			return createScoring(self.Cpm()), nil
		}
		if self.cursorPos == self.newline {
			self.start = self.newline
			self.newline = self.end
			self.end = self.newline + CalculateLineBreak(self.display.Inner.Dx(), self.words[self.newline:])
		}
		self.input.Text = Cursor
	case "<Tab>", "<Enter>":
	case "<Backspace>":
		if length > 0 {
			self.input.Text = text[:length-1] + Cursor
		}
	default:
		if !self.started {
			self.startTime = time.Now()
			self.started = true
		}
		self.input.Text = text + event.ID + Cursor
	}
	return self, nil
}

func (self Typing) Render() {
	self.display.Text = getDisplayText(self.words, self.start, self.newline, self.end)
	ui.Render(self.display, self.input)
}

func getDisplayText(words []string, start, newline, end int) string {
	length := len(words)
	text := strings.Join(words[start:Min(length, newline)], " ")
	text += "\n"
	text += strings.Join(words[Min(length, newline):Min(length, end)], " ")

	return text
}

func (self Typing) Cpm() float64 {
	return (60. * float64(self.correctCharacters) /
		float64(time.Since(self.startTime).Seconds()))
}

func updateCpm(word string, typing *Typing) {
	correctCharacters := 0
	correctWord := typing.words[typing.cursorPos]
	if word == correctWord {
		correctCharacters = len(word) + 1 // +1 for the space
	}
	typing.correctCharacters += correctCharacters

	cpm := typing.Cpm()
	typing.input.Title = fmt.Sprintf("CPM: %.0f", cpm)
}

func checkWord(word string, cursorPos int, words *[]string) {
	if ref := (*words)[cursorPos]; strings.Trim(word, Cursor) == ref {
		// input is correct
		(*words)[cursorPos] = "[" + ref + "](fg:" + CorrectFg + ")"
	} else {
		// input is false
		(*words)[cursorPos] = "[" + ref + "](fg:" + FalseFg + ")"

	}
}

func createTyping(lesson Lesson) Typing {
	display := widgets.NewParagraph()
	display.Title = lesson.Title
	display.SetRect(MainMinX, MainMinY, MainMaxX, 5)

	input := widgets.NewParagraph()
	input.Title = ""
	input.Text = Cursor
	input.SetRect(MainMinX, 6, MainMaxX, 9)
	words := strings.Split(lesson.Content, " ")
	start := 0
	newline := CalculateLineBreak(display.Inner.Dx(), words)
	end := newline + CalculateLineBreak(display.Inner.Dx(), words[newline:])

	return Typing{
		title:             "Typing",
		input:             input,
		display:           display,
		words:             words,
		cursorPos:         0,
		start:             start,
		newline:           newline,
		end:               end,
		totalCharacters:   len(lesson.Content) + 1, // +1 for the final space
		started:           false,
		startTime:         time.Now(),
		correctCharacters: 0,
	}
}
