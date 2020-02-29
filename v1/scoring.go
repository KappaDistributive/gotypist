package main

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// Scoring implements Viewport
type Scoring struct {
	title              string
	card               *widgets.Paragraph
	selectionCursorPos int
}

// Handler manges ui events
func (scoring Scoring) Handler(e <-chan ui.Event) (Viewport, error) {
	event := <-e

	if event.Type == ui.KeyboardEvent {
		switch event.ID {
		case "<C-c>":
			return scoring, Quit{}
		case "<Enter>":
			return createSelection(scoring.selectionCursorPos), nil
		}
	}
	return scoring, nil
}

// Render renders the ui
func (scoring Scoring) Render() {
	ui.Render(scoring.card)
}

// Cpm calculates characters per minute
func Cpm(correctCharacters int, duration time.Duration) float64 {
	return 60.0 * float64(correctCharacters) / float64(duration.Seconds())
}

// Accuracy calculates the ratio of correctly typed characters
func Accuracy(correctCharacters int, typedCharacters int) float64 {
	return float64(correctCharacters) / float64(typedCharacters)
}

// CreateScoring creates a scoring view
func CreateScoring(correctCharacters int, totalCharacters int, duration time.Duration, selectionCursorPos int) Scoring {
	cpm := Cpm(correctCharacters, duration)
	accuracy := Accuracy(correctCharacters, totalCharacters)
	card := widgets.NewParagraph()
	card.Title = "Scoring Card"
	card.Text = fmt.Sprintf("CPM: %.0f\nAccuracy: %.2f%%", cpm, 100.0*accuracy)
	card.SetRect(MainMinX, MainMinY, MainMaxX, MainMaxY)
	return Scoring{
		title:              "Scoring",
		card:               card,
		selectionCursorPos: selectionCursorPos,
	}
}
