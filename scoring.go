package main

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// Scoring implements Viewport
type Scoring struct {
	title string
	card  *widgets.Paragraph
}

func (self Scoring) Handler(e <-chan ui.Event) (Viewport, error) {
	event := <-e
	switch event.ID {
	case "<C-c>":
		return self, Quit{}
	case "<Enter>":
		return createSelection(), nil
	}
	return self, nil
}

func (self Scoring) Render() {
	ui.Render(self.card)
}

func createScoring(cpm float64) Scoring {
	card := widgets.NewParagraph()
	card.Title = "Scoring Card"
	card.Text = fmt.Sprintf("CPM: %.0f", cpm)

	card.SetRect(MainMinX, MainMinY, MainMaxX, MainMaxY)
	return Scoring{
		title: "Scoring",
		card:  card,
	}
}
