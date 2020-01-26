package main

import (
	"log"
	"os"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const Tabwidth int = 4

const (
	MainMinX int = 0
	MainMaxX int = 80
	MainMinY int = 0
	MainMaxY int = 15
)

type Viewport interface {
	Render()
	Handler(<-chan ui.Event) Viewport
}

// Selection implements Viewport
type Selection struct {
	title   string
	lessons *widgets.List
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
		return createTyping()
	}
	return self
}

// Typing implements Viewport
type Typing struct {
	title     string
	paragraph *widgets.Paragraph
}

func (self Typing) Render() {
	ui.Render(self.paragraph)
}

func (self Typing) Handler(e <-chan ui.Event) Viewport {
	event := <-e
	text := self.paragraph.Text
	length := len(text)

	switch event.ID {
	case "<C-c>":
		os.Exit(0)
	case "<Escape>":
		return createSelection()
	case "<Space>":
		self.paragraph.Text = text[:length-3] + " " + text[length-3:]
	case "<Tab>":
		self.paragraph.Text = text[:length-3] + strings.Repeat(" ", Tabwidth) + text[length-3:]
	case "<Enter>":
		self.paragraph.Text = text[:length-3] + "\n" + text[length-3:]
	case "<Backspace>":
		if length > 3 {
			self.paragraph.Text = text[:length-4] + text[length-3:]
		}
	default:
		self.paragraph.Text = text[:length-3] + event.ID + text[length-3:]
	}
	return self
}

func createSelection() Selection {
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
	lessons.SetRect(MainMinX, MainMinY, MainMaxX, MainMaxY)
	lessons.SelectedRowStyle = ui.NewStyle(ui.ColorGreen)
	return Selection{
		title:   "Selection",
		lessons: lessons,
	}
}

func createTyping() Typing {
	paragraph := widgets.NewParagraph()
	paragraph.Title = "Paragraph"
	paragraph.Text = "\u2588"
	paragraph.SetRect(MainMinX, MainMinY, MainMaxX, MainMaxY)

	return Typing{
		title:     "Typing",
		paragraph: paragraph,
	}

}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	var view Viewport
	view = createSelection()
	view.Render()

	uiEvents := ui.PollEvents()
	for {
		view = view.Handler(uiEvents)
		view.Render()
	}

}
