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

func (self Selection) Handler(e <-chan ui.Event) (Viewport, error) {
	event := <-e
	switch event.ID {
	case "<C-c>":
		return self, Quit{}
	case "<Up>", "k":
		self.content.ScrollUp()
	case "<Down>", "j":
		self.content.ScrollDown()
	case "<Enter>":
		self.savedCursorPos = self.content.SelectedRow
		lesson := self.lessons[self.savedCursorPos]

		return createTyping(lesson, self.savedCursorPos), nil
	}
	return self, nil
}

func (self Selection) Render() {
	ui.Render(self.content)
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
		lesson_name := lesson.Title
		content.Rows = append(content.Rows, lesson_name)
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
