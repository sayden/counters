package transform

import (
	"math/rand"

	"github.com/imdario/mergo"
	"github.com/sayden/counters"
)

type EventsToCardsConfig struct {
	Events             []counters.Event
	Images             []string
	BackImageFile      string
	GeneratedImageName string
}

func stringP(s string) *string {
	return &s
}

func intP(i int) *int {
	return &i
}

func floatP(f float64) *float64 {
	return &f
}

func EventsToCards(cfg *EventsToCardsConfig) *counters.CardsTemplate {
	borderWidth := 20.0

	settings := counters.Settings{
		Width:           742,
		Height:          1200,
		Margins:         floatP(20),
		FontHeight:      60,
		FontColorS:      "white",
		BackgroundColor: stringP("gold"),
		BorderWidth:     &borderWidth,
		BorderColorS:    "black",
		StrokeColorS:    "black",
		Alignment:       "center",
		ImageScaling:    "fitHeight",
	}

	titleSettings := counters.Settings{}
	mergo.Merge(&titleSettings, settings)
	titleSettings.StrokeWidth = floatP(7)
	titleSettings.AvoidClipping = true
	titleSettings.Position = 3

	textSettings := counters.Settings{}
	mergo.Merge(&textSettings, settings)
	textSettings.FontColorS = "black"
	textSettings.AvoidClipping = false
	textSettings.FontHeight = 50

	downAreaSettings := counters.Settings{}
	mergo.Merge(&downAreaSettings, settings)

	cards := make([]counters.Card, 0)
	for _, event := range cfg.Events {
		bgImage := cfg.Images[rand.Intn(len(cfg.Images))]

		card := counters.Card{
			Settings: settings,
			Areas: []counters.Counter{
				{Images: []counters.Image{{Path: bgImage}}},
				{
					Settings: downAreaSettings,
					Images:   []counters.Image{{Path: "assets/old_paper.jpg"}},
					Texts: []counters.Text{
						{Settings: textSettings, String: event.Desc},
						{String: event.Title, Settings: titleSettings},
					},
				},
			},
			Images: nil,
		}

		cards = append(cards, card)
	}

	cardTemplate := counters.CardsTemplate{
		Settings:   settings,
		Rows:       7,
		Columns:    10,
		DrawGuides: false,
		Mode:       "tiles",
		OutputPath: cfg.GeneratedImageName,
		Cards:      cards,
	}

	return &cardTemplate
}
