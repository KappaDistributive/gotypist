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
	title              string
	input              *widgets.Paragraph
	display            *widgets.Paragraph
	words              []string
	wordStatus         []string
	cursorPos          int
	start              int
	newline            int
	end                int
	totalCharacters    int
	started            bool
	startTime          time.Time
	correctCharacters  int
	typedCharaceters   int
	selectionCursorPos int
}

func (self Typing) Handler(e <-chan ui.Event) (Viewport, error) {
	event := <-e
	text := DropCursor(self.input.Text)
	length := len(text)

	switch event.ID {
	case "<C-c>":
		return self, Quit{}
	case "<Escape>":
		return createSelection(self.selectionCursorPos), nil
	// TODO: replace ad-hoc text handling
	case "<Space>":
		if text == self.words[self.cursorPos] {
			self.wordStatus[self.cursorPos] = StatusCorrect
		} else {
			self.wordStatus[self.cursorPos] = StatusIncorrect
		}
		updateCpm(text, &self)
		self.cursorPos += 1
		if self.cursorPos == len(self.words) {
			// end the game
			duration := time.Since(self.startTime)
			return CreateScoring(self.correctCharacters, self.totalCharacters, duration, self.selectionCursorPos), nil
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
			text := text[:length-1]
			self.setSubwordStatus(text)
			if text == self.words[self.cursorPos][:len(text)] {
				self.wordStatus[self.cursorPos] = StatusNeutral
			} else {
				self.wordStatus[self.cursorPos] = StatusIncorrect
			}
			self.input.Text = text + Cursor
		}
	default:
		if !self.started {
			self.startTime = time.Now()
			self.started = true
		}
		text = text + event.ID
		self.setSubwordStatus(text)
		self.input.Text = text + Cursor
	}
	return self, nil
}

func (self Typing) Render() {
	self.UpdateText()
	ui.Render(self.display, self.input)
}

func (self Typing) UpdateText() {
	words := make([]string, len(self.words))
	copy(words, self.words)
	for i, word := range words {
		switch self.wordStatus[i] {
		case StatusCorrect:
			words[i] = "[" + word + "](fg:" + CorrectFg + ")"
		case StatusIncorrect:
			words[i] = "[" + word + "](fg:" + FalseFg + ")"
		}
	}
	text := strings.Join(words[self.start:Min(len(words), self.newline)], " ")
	text += "\n"
	text += strings.Join(
		words[Min(len(words),
			self.newline):Min(len(words), self.end)], " ")

	self.display.Text = text
}

func (self Typing) setSubwordStatus(word string) {
	if word == self.words[self.cursorPos][:len(word)] {
		self.wordStatus[self.cursorPos] = StatusNeutral
	} else {
		self.wordStatus[self.cursorPos] = StatusIncorrect

	}
}

func updateCpm(word string, typing *Typing) {
	correctWord := typing.words[typing.cursorPos]
	if word == correctWord {
		typing.correctCharacters += len(word) + 1 // +1 for the space
	}
	typing.typedCharaceters += len(word) + 1 // +1 for the space

	cpm := Cpm(typing.correctCharacters, time.Since(typing.startTime))
	typing.input.Title = fmt.Sprintf("CPM: %.0f", cpm)
}

func createTyping(lesson Lesson, selectionCursorPos int) Typing {
	display := widgets.NewParagraph()
	display.Title = lesson.Title
	display.SetRect(MainMinX, MainMinY, MainMaxX, 5)

	input := widgets.NewParagraph()
	input.Title = ""
	input.Text = Cursor
	input.SetRect(MainMinX, 6, MainMaxX, 9)
	words := strings.Split(lesson.Content, " ")
	var wordStatus []string
	for _, _ = range words {
		wordStatus = append(wordStatus, StatusNeutral)
	}
	start := 0
	newline := CalculateLineBreak(display.Inner.Dx(), words)
	end := newline + CalculateLineBreak(display.Inner.Dx(), words[newline:])

	return Typing{
		title:              "Typing",
		input:              input,
		display:            display,
		words:              words,
		wordStatus:         wordStatus,
		cursorPos:          0,
		start:              start,
		newline:            newline,
		end:                end,
		totalCharacters:    len(lesson.Content) + 1, // +1 for the final space
		started:            false,
		startTime:          time.Now(),
		correctCharacters:  0,
		selectionCursorPos: selectionCursorPos,
	}
}
