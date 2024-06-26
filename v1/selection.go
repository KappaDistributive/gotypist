package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"gopkg.in/yaml.v2"
)

// Selection implements Viewport
type Selection struct {
	title          string
	lessons        []Lesson
	content        *widgets.List
	savedCursorPos int
}

// Handler manages ui events
func (selection Selection) Handler(e <-chan ui.Event) (Viewport, error) {
	event := <-e

	if event.Type == ui.KeyboardEvent {
		switch event.ID {
		case "<C-c>":
			return selection, Quit{}
		case "<Up>", "k":
			selection.content.ScrollUp()
		case "<Down>", "j":
			selection.content.ScrollDown()
		case "<Enter>":
			selection.savedCursorPos = selection.content.SelectedRow
			lesson := selection.lessons[selection.savedCursorPos]

			return createTyping(lesson, selection.savedCursorPos), nil
		}
	}
	return selection, nil
}

// Render renders the ui
func (selection Selection) Render() {
	selection.savedCursorPos = selection.content.SelectedRow
	lesson := selection.lessons[selection.savedCursorPos]
	selection.content.Title = fmt.Sprintf("Lesson | %v", lesson.Tag)
	ui.Render(selection.content)
}

func createSelection(cursorPos int) Selection {
	home, err := os.UserHomeDir()
	if err != nil {
		errorHandling(err)
	}

	// Create lessons
	content := widgets.NewList()
	content.Title = "Lessons"
	content.Rows = []string{}
	lessons := []Lesson{}

	// Create top 300 words lesson
	data, err := os.ReadFile(home + BagsOfWordsDir + "/en_us.yaml")
	if err != nil {
		errorHandling(err)
	}
	bagOfWords := BagOfWords{}
	if err = yaml.Unmarshal(data, &bagOfWords); err != nil {
		errorHandling(err)
	}
	lesson_text := []string{}
	lesson_length := 0
	var index int
	var new_word string
	for lesson_length < BagsOfWordsLessonLengthInCharacters {
		index = rand.Intn(299)
		new_word = bagOfWords.Words[index]
		lesson_text = append(lesson_text, new_word)
		lesson_length += len(new_word)
	}
	lesson := Lesson{
		Title:   "Top 300 words",
		Content: strings.Join(lesson_text, " "),
		Tag:     DASH_MODE,
	}
	lessons = append(lessons, lesson)
	content.Rows = append(content.Rows, lesson.Title)

	lesson = Lesson{
		Title:   "Random Prose",
		Content: "",
		Tag:     PROSE_MODE,
	}
	lessons = append(lessons, lesson)
	content.Rows = append(content.Rows, lesson.Title)

	// Load lessons from directory
	files, err := os.ReadDir(home + LessonsDir)
	if err != nil {
		errorHandling(err)
	}
	for _, fileinfo := range files {
		data, err := os.ReadFile(home + LessonsDir + "/" + fileinfo.Name())
		if err != nil {
			errorHandling(err)
		}
		lesson := Lesson{}
		if err = yaml.UnmarshalStrict(data, &lesson); err != nil {
			errorHandling(err)
		}
		lesson.Tag = PROSE_MODE
		lessons = append(lessons, lesson)
		content.Rows = append(content.Rows, lesson.Title)
	}

	random_index := rand.Intn(len(lessons)-2) + 2
	lessons[1].Title = "Random Prose - " + lessons[random_index].Title
	lessons[1].Content = lessons[random_index].Content

	content.SetRect(MainMinX, MainMinY, MainMaxX, MainMaxY)
	content.SelectedRowStyle = ui.NewStyle(ui.ColorGreen)
	content.SelectedRow = cursorPos
	return Selection{
		title:          "Selection",
		lessons:        lessons,
		content:        content,
		savedCursorPos: cursorPos,
	}
}
