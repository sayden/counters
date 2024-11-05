package counters

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardTemplate(t *testing.T) {
	template := &CardsTemplate{
		Scaling: 1.5,
		Settings: Settings{
			Width:           800,
			Height:          600,
			StrokeWidth:     floatP(2),
			BackgroundColor: stringP("white"),
			FontColorS:      "black",
			FontPath:        "assets/freesans.ttf",
		},
		Cards: []Card{
			{
				Settings: Settings{
					BackgroundColor: stringP("lightgrey"),
				},
				Areas: []Counter{
					{
						Settings: Settings{
							Height:          50,
							BackgroundColor: stringP("black"),
						},
						Texts: []Text{
							{
								Settings: Settings{StrokeWidth: floatP(1)},
								String:   "Area text",
							},
						},
					},
				},
				Images: []Image{
					{
						Settings: Settings{Position: 11},
						Path:     "assets/binoculars.png",
						Scale:    0.2,
					},
				},
				Texts: []Text{
					{
						Settings: Settings{StrokeWidth: floatP(1)},
						String:   "Full card text",
					},
				},
			},
		},
	}

	byt, err := json.Marshal(template)
	if assert.NoError(t, err) {
		newTempl, err := ParseCardTemplate(byt)
		if assert.NoError(t, err) {
			canvas, err := newTempl.Canvas(&newTempl.Settings, newTempl.Settings.Width, newTempl.Settings.Height)
			if assert.NoError(t, err) {
				assert.NotNil(t, canvas)
			}
		}
	}

}
