package main

import (
	ui "github.com/gizak/termui/v3"
)

// Configurations
const (
	Tabwidth                            int    = 4
	Cursor                              string = "\u2588"
	MainMinX                            int    = 0
	MainMaxX                            int    = 80
	MainMinY                            int    = 0
	MainMaxY                            int    = 15
	CorrectFg                           string = "green"
	FalseFg                             string = "red"
	ConfigDir                           string = "/.config/gotypist"
	LessonsDir                          string = ConfigDir + "/lessons"
	BagsOfWordsDir                      string = ConfigDir + "/bags_of_words"
	BagsOfWordsLessonLengthInCharacters int    = 500
)

const (
	// StatusNeutral is the status for neutral text
	StatusNeutral string = "neutral"
	// StatusCorrect is the status for correct text input
	StatusCorrect string = "correct"
	// StatusIncorrect is the status for incorrect text input
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
