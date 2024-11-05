package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/sayden/counters"
	"github.com/stretchr/testify/assert"
)

func TestCardTemplate(t *testing.T) {
	template := &counters.CardsTemplate{
		OutputPath: "/tmp/card_out_%02d.png",
		Rows:       2,
		Columns:    2,
		Settings: counters.Settings{
			BorderWidth:     floatP(3),
			BorderColorS:    "green",
			Width:           200,
			Height:          280,
			StrokeWidth:     floatP(2),
			FontPath:        "../assets/freesans.ttf",
			BackgroundColor: stringP("white"),
			StrokeColorS:    "white",
			FontColorS:      "black",
		},
		Cards: []counters.Card{
			{
				Settings: counters.Settings{
					Multiplier:      intP(2),
					BackgroundColor: stringP("lightgrey"),
				},
				Areas: []counters.Counter{
					{
						Settings: counters.Settings{
							Height:          50,
							BorderWidth:     floatP(2),
							BorderColorS:    "red",
							BackgroundColor: stringP("black"),
						},
						Texts: []counters.Text{
							{
								Settings: counters.Settings{
									StrokeWidth:   floatP(2),
									FontHeight:    12,
									Position:      11,
									AvoidClipping: true,
								},
								String: "Area text",
							},
						},
					},
				},
				Images: []counters.Image{
					{
						Settings: counters.Settings{Position: 11},
						Path:     "../assets/binoculars.png",
						Scale:    0.2,
					},
				},
				Texts: []counters.Text{
					{
						Settings: counters.Settings{StrokeWidth: floatP(2)},
						String:   "Full card text",
					},
				},
			},
		},
	}

	byt, err := json.Marshal(template)
	if assert.NoError(t, err) {
		newTempl, err := counters.ParseCardTemplate(byt)
		if assert.NoError(t, err) {
			err := CardsToPNG(newTempl)
			assert.NoError(t, err)

			f, err := os.Open(fmt.Sprintf(template.OutputPath, 1))
			if !assert.NoError(t, err) {
				t.FailNow()
			}
			defer os.Remove(f.Name())
			gotBytes, err := os.ReadFile(fmt.Sprintf(template.OutputPath, 1))
			if assert.NoError(t, err) {
				expectedFile, err := os.Open("../testdata/card_template_01.png")
				if assert.NoError(t, err) {
					defer expectedFile.Close()

					expectedBytes, err := io.ReadAll(expectedFile)
					if assert.NoError(t, err) {
						assert.Equal(t, expectedBytes, gotBytes)
					}
				}
			}
		}
	}
}
