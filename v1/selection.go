package main

import (
	"io/ioutil"
	"os"

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
	return selection, nil
}

// Render renders the ui
func (selection Selection) Render() {
	ui.Render(selection.content)
}

func createSelection(cursorPos int) Selection {
	home, err := os.UserHomeDir()
	if err != nil {
		errorHandling(err)
	}

	// Load lessons from directory
	content := widgets.NewList()
	content.Title = "Lessons"
	content.Rows = []string{}
	lessons := []Lesson{}
	files, err := ioutil.ReadDir(home + LessonsDir)
	if err != nil {
		errorHandling(err)
	}
	for _, fileinfo := range files {
		lesson := Lesson{}
		data, err := ioutil.ReadFile(home + LessonsDir + "/" + fileinfo.Name())
		if err != nil {
			errorHandling(err)
		}
		if err = yaml.Unmarshal(data, &lesson); err != nil {
			errorHandling(err)
		}
		lessons = append(lessons, lesson)
		lessonName := lesson.Title
		content.Rows = append(content.Rows, lessonName)
	}

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
