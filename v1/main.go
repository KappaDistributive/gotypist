package main

import (
	ui "github.com/gizak/termui/v3"
)

// Configurations
const (
	Tabwidth   int    = 4
	Cursor     string = "\u2588"
	MainMinX   int    = 0
	MainMaxX   int    = 80
	MainMinY   int    = 0
	MainMaxY   int    = 15
	CorrectFg  string = "green"
	FalseFg    string = "red"
	ConfigDir  string = "/.config/gotypist"
	LessonsDir string = ConfigDir + "/lessons"
)

const (
	StatusNeutral   string = "neutral"
	StatusCorrect   string = "correct"
	StatusIncorrect string = "incorrect"
)

func main() {
	var view Viewport
	initialize(&view)
	defer ui.Close()
	var err error

	// event loop
	uiEvents := ui.PollEvents()
	for {
		view, err = view.Handler(uiEvents)
		if err != nil {
			return
		}
		ui.Clear()
		view.Render()
	}
}
