package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"gopkg.in/yaml.v2"
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

// Quit implements error
// Throw this when quitting the app.
type Quit struct{}

func (self Quit) Error() string {
	return "Quit program."
}

type Lesson struct {
	Title   string
	Content string
}

type Viewport interface {
	Render()                                   // Render takes care of rendering the UI.
	Handler(<-chan ui.Event) (Viewport, error) // Handler takes care of the event loop.
}

// Selection implements Viewport
type Selection struct {
	title   string
	lessons *widgets.List
}

func (self Scoring) Handler(e <-chan ui.Event) (Viewport, error) {
	event := <-e
	switch event.ID {
	case "<C-c>":
		return createSelection(), Quit{}
	case "<Enter>":
		return createSelection(), nil
	}
	return self, nil
}

func (self Selection) Render() {
	ui.Render(self.lessons)
}

func (self Selection) Handler(e <-chan ui.Event) (Viewport, error) {
	event := <-e
	switch event.ID {
	case "<C-c>":
		return self, Quit{}
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
		return createTyping(lesson), nil
	}
	return self, nil
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
	text := strings.Join(words[start:Min(length, newline)], " ")
	text += "\n"
	text += strings.Join(words[Min(length, newline):Min(length, end)], " ")

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

func (self Typing) Handler(e <-chan ui.Event) (Viewport, error) {
	event := <-e
	text := self.input.Text
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
			return createSelection(), nil
		}
		if self.cursorPos == self.newline {
			self.start = self.newline
			self.newline = self.end
			self.end = self.newline + CalculateLineBreak(self.display.Inner.Dx(), self.words[self.newline:])
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
	return self, nil
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
	home, err := os.UserHomeDir()
	if err != nil {
		errorHandling(err)
	}

	lessons := widgets.NewList()
	lessons.Title = "Lessons"
	// Create lessons.Rows from the content of the lessons directory.
	lessons.Rows = []string{}
	files, err := ioutil.ReadDir(home + "/.config/gotypist/lessons")
	if err != nil {
		errorHandling(err)
	}
	for _, file := range files {
		lesson_name_splits := strings.Split(file.Name(), ".")
		lesson_name := strings.Join(lesson_name_splits[:len(lesson_name_splits)-1], "")
		lessons.Rows = append(lessons.Rows, lesson_name)
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

	// if no lessons exist, add sample lessons
	createSampleLessons()

	*view = createSelection()
	(*view).Render()

}

func createSampleLessons() {
	home, err := os.UserHomeDir()
	if err != nil {
		errorHandling(err)
	}

	files, err := ioutil.ReadDir(home + "/.config/gotypist/lessons")
	if err != nil {
		errorHandling(err)
	}

	// check whether there are already lessons
	for _, file := range files {
		splits := strings.Split(file.Name(), ".")
		if splits[len(splits)-1] == "yaml" {
			return
		}
	}

	// create sample lessons if no lessons exist
	files, err = ioutil.ReadDir("data/sample_lessons")
	if err != nil {
		errorHandling(err)
	}

	for _, lesson := range files {
		source_file, err := os.Open(fmt.Sprintf("data/sample_lessons/%s", lesson.Name()))
		if err != nil {
			errorHandling(err)
		}
		target_file, err := os.Create(home + fmt.Sprintf("/.config/gotypist/lessons/%s", lesson.Name()))
		if err != nil {
			errorHandling(err)
		}
		io.Copy(target_file, source_file)

		source_file.Close()
		target_file.Close()
	}
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
		view, err := view.Handler(uiEvents)
		if err != nil {
			ui.Close()
			os.Exit(0)
		}

		ui.Clear()
		view.Render()
	}
}
