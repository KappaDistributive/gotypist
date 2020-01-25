package main

import (
	"fmt"
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Mode int

const (
	Selection Mode = iota
	Typing
)

func (m Mode) String() string {
	return [...]string{"Selection", "Typing"}[m]
}

func startLesson(lessons *widgets.List, test *widgets.Paragraph) {
	test.Text = fmt.Sprintf("Lesson %d", lessons.SelectedRow+1)
}

type Viewport struct {
	items []ui.Drawable
}

func (v *Viewport) Render() {
	ui.Clear()
	for _, item := range v.items {
		ui.Render(item)
	}
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	lessons := widgets.NewList()
	lessons.Title = "Lessons"
	lessons.Rows = []string{
		"Lesson 1",
		"Lesson 2",
		"Lesson 3",
		"Lesson 4",
		"Lesson 5",
		"Lesson 6",
		"Lesson 7",
		"Lesson 8",
		"Lesson 9",
		"Lesson 10",
	}
	lessons.SetRect(0, 0, 80, 5)
	lessons.SelectedRowStyle = ui.NewStyle(ui.ColorGreen)
	items := []ui.Drawable{lessons}
	selectionUi := Viewport{items}

	selectionUi.Render()

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j":
			lessons.ScrollDown()
		case "k":
			lessons.ScrollUp()
		}
		selectionUi.Render()
	}

}
