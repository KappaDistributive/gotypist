package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"gopkg.in/yaml.v2"
)

// Selection implements Viewport
type Selection struct {
	title   string
	lessons *widgets.List
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
			errorHandling(err)
		}
		return createTyping(lesson), nil
	}
	return self, nil
}

func (self Selection) Render() {
	ui.Render(self.lessons)
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
