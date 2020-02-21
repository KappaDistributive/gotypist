package main

import (
	"os"

	ui "github.com/gizak/termui/v3"
)

// Viewport provides a generic interface for a ui element
type Viewport interface {
	Render()                                   // Render takes care of rendering the UI.
	Handler(<-chan ui.Event) (Viewport, error) // Handler takes care of the event loop.
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

	*view = createSelection(0)
	(*view).Render()

}
