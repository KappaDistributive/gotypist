package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"gopkg.in/yaml.v2"

	"github.com/KappaDistributive/gotypist/utils"
)

// Configurations
const (
	Tabwidth  int    = 4
	Cursor    string = "\u2588"
	MainMinX  int    = 0
	MainMaxX  int    = 80
	MainMinY  int    = 0
	MainMaxY  int    = 15
	CorrectFg string = "green"
	FalseFg   string = "red"
)

type Lesson struct {
	Title   string
	Content string
}

type Viewport interface {
	Render()                          // Render takes care of rendering the UI.
	Handler(<-chan ui.Event) Viewport // Handler takes care of the event loop.
}

// Selection implements Viewport
type Selection struct {
	title   string
	lessons *widgets.List
}

func (self Scoring) Handler(e <-chan ui.Event) Viewport {
	event := <-e
	switch event.ID {
	case "<C-c>":
		os.Exit(0)
	case "<Enter>":
		return createSelection()
	}
	return self
}

func (self Selection) Render() {
	ui.Render(self.lessons)
}

func (self Selection) Handler(e <-chan ui.Event) Viewport {
	event := <-e
	switch event.ID {
	case "<C-c>":
		os.Exit(0)
	case "<Up>", "k":
		self.lessons.ScrollUp()
	case "<Down>", "j":
		self.lessons.ScrollDown()
	case "<Enter>":
		cursorPos := self.lessons.SelectedRow
		home, err := os.UserHomeDir()
		if err != nil {
			errorHandling(err)
		}
		// TODO: Create test lessons if none exist, yet.
		data, err := ioutil.ReadFile(fmt.Sprintf(home+"/.config/gotypist/lessons/%s.yaml", self.lessons.Rows[cursorPos]))
		if err != nil {
			errorHandling(err)
		}

		lesson := Lesson{}

		if err = yaml.Unmarshal([]byte(data), &lesson); err != nil {
			log.Fatal(err)
		}
		return createTyping(lesson)
	}
	return self
}

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

func getDisplayText(words []string, start, newline, end int) string {
	length := len(words)
	text := strings.Join(words[start:utils.Min(length, newline)], " ")
	text += "\n"
	text += strings.Join(words[utils.Min(length, newline):utils.Min(length, end)], " ")

	return text
}

func (self Typing) Cpm() float64 {
	return (60. * float64(self.correctCharacters) /
		float64(time.Since(self.startTime).Seconds()))
}

func (self Typing) Render() {
	self.display.Text = getDisplayText(self.words, self.start, self.newline, self.end)
	ui.Render(self.display, self.input)
}

func (self Typing) Handler(e <-chan ui.Event) Viewport {
	event := <-e
	text := self.input.Text
	length := len(text)

	switch event.ID {
	case "<C-c>":
		os.Exit(0)
	case "<Escape>":
		return createSelection()
	// TODO: replace ad-hoc text handling
	case "<Space>":
		updateCpm(text, &self)
		checkWord(text, self.cursorPos, &self.words)
		self.cursorPos += 1
		if self.cursorPos == len(self.words) {
			// end the game
			return createScoring(self.Cpm())
			return createSelection()
		}
		if self.cursorPos == self.newline {
			self.start = self.newline
			self.newline = self.end
			self.end = self.newline + utils.CalculateLineBreak(self.display.Inner.Dx(), self.words[self.newline:])
		}
		self.input.Text = Cursor
	case "<Tab>", "<Enter>":
	case "<Backspace>":
		if length > len(Cursor) {
			self.input.Text = text[:length-len(Cursor)-1] + text[length-len(Cursor):]
		}
	default:
		if !self.started {
			self.startTime = time.Now()
			self.started = true
		}
		self.input.Text = text[:length-len(Cursor)] + event.ID + text[length-len(Cursor):]
	}
	return self
}

// Scoring implements Viewport
type Scoring struct {
	title string
	card  *widgets.Paragraph
}

func (self Scoring) Render() {
	ui.Render(self.card)
}

func updateCpm(word string, typing *Typing) {
	correctCharacters := 0
	correctWord := typing.words[typing.cursorPos]
	for pos, char := range correctWord {
		if pos < len(word) && word[pos] == byte(char) {
			correctCharacters += 1
		}
	}
	typing.correctCharacters += correctCharacters + 1 // +1 for the space

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

func createScoring(cpm float64) Scoring {
	card := widgets.NewParagraph()
	card.Title = "Scoring Card"
	card.Text = fmt.Sprintf("CPM: %.0f", cpm)

	card.SetRect(MainMinX, MainMinY, MainMaxX, MainMaxY)
	return Scoring{
		title: "Scoring",
		card:  card,
	}
}

func createSelection() Selection {
	lessons := widgets.NewList()
	lessons.Title = "Lessons"
	// Create lessons.Rows from the content of the lessons directory.
	lessons.Rows = []string{
		"Lesson 1",
		"Lesson 2",
		"Lesson 3",
	}
	lessons.SetRect(MainMinX, MainMinY, MainMaxX, MainMaxY)
	lessons.SelectedRowStyle = ui.NewStyle(ui.ColorGreen)
	return Selection{
		title:   "Selection",
		lessons: lessons,
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
	newline := utils.CalculateLineBreak(display.Inner.Dx(), words)
	end := newline + utils.CalculateLineBreak(display.Inner.Dx(), words[newline:])

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

func initialize(view *Viewport) {
	// initialize ui
	if err := ui.Init(); err != nil {
		errorHandling(err)
	}

	// create config directories
	home, err := os.UserHomeDir()
	if err != nil {
		errorHandling(err)
	}
	// TODO: make config path configurable
	if err := os.MkdirAll(home+"/.config/gotypist/lessons", 0755); err != nil {
		errorHandling(err)
	}

	*view = createSelection()
	(*view).Render()

}

func errorHandling(err error) {
	log.Fatal(err)
}

func main() {
	var view Viewport
	initialize(&view)
	defer ui.Close()

	// event loop
	uiEvents := ui.PollEvents()
	for {
		view = view.Handler(uiEvents)
		ui.Clear()
		view.Render()
	}
}
