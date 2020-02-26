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

// Handler handles ui events
func (typing Typing) Handler(e <-chan ui.Event) (Viewport, error) {
	event := <-e
	text := DropCursor(typing.input.Text)

	switch event.ID {
	case "<C-c>":
		return typing, Quit{}
	case "<Escape>":
		return createSelection(typing.selectionCursorPos), nil
	// TODO: replace ad-hoc text handling
	case "<Space>":
		if text == typing.words[typing.cursorPos] {
			typing.wordStatus[typing.cursorPos] = StatusCorrect
		} else {
			typing.wordStatus[typing.cursorPos] = StatusIncorrect
		}
		updateCpm(text, &typing)
		typing.cursorPos++
		if typing.cursorPos == len(typing.words) {
			// end the game
			duration := time.Since(typing.startTime)
			return CreateScoring(typing.correctCharacters, typing.totalCharacters, duration, typing.selectionCursorPos), nil
		}
		if typing.cursorPos == typing.newline {
			typing.start = typing.newline
			typing.newline = typing.end
			typing.end = typing.newline + CalculateLineBreak(typing.display.Inner.Dx(), typing.words[typing.newline:])
		}
		typing.input.Text = Cursor
	case "<Tab>", "<Enter>":
	case "<Backspace>":
		if len(text) > 0 {
			text := text[:len(text)-1]
			typing.setSubwordStatus(text)
			typing.input.Text = text + Cursor
		}
	default:
		if !typing.started {
			typing.startTime = time.Now()
			typing.started = true
		}
		text = text + event.ID
		typing.setSubwordStatus(text)
		typing.input.Text = text + Cursor
	}
	return typing, nil
}

// Render renders the ui
func (typing Typing) Render() {
	typing.UpdateText()
	ui.Render(typing.display, typing.input)
}

// UpdateText computes the text to be displayed
func (typing Typing) UpdateText() {
	words := make([]string, len(typing.words))
	copy(words, typing.words)
	for i, word := range words {
		switch typing.wordStatus[i] {
		case StatusCorrect:
			words[i] = "[" + word + "](fg:" + CorrectFg + ")"
		case StatusIncorrect:
			words[i] = "[" + word + "](fg:" + FalseFg + ")"
		}
	}
	text := strings.Join(words[typing.start:Min(len(words), typing.newline)], " ")
	text += "\n"
	text += strings.Join(
		words[Min(len(words),
			typing.newline):Min(len(words), typing.end)], " ")

	typing.display.Text = text
}

func (typing Typing) setSubwordStatus(word string) {
	length := Min(len(word), len(typing.words[typing.cursorPos]))
	if word == typing.words[typing.cursorPos][:length] {
		typing.wordStatus[typing.cursorPos] = StatusNeutral
	} else {
		typing.wordStatus[typing.cursorPos] = StatusIncorrect

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
	for range words {
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
