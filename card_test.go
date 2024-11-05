package counters

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCanvas(t *testing.T) {
	width := 800
	height := 600
	settings := &Settings{
		FontPath:   "assets/freesans.ttf",
		FontHeight: 12,
		BgColor:    color.RGBA{255, 255, 255, 255},
		Width:      width,
		Height:     height,
		FontColorS: "#000000",
	}
	template := &CardsTemplate{
		Settings: *settings,
	}

	canvas, err := template.Canvas(settings, width, height)
	assert.NoError(t, err)
	assert.NotNil(t, canvas)
}

func TestApplyCardScaling(t *testing.T) {
	template := &CardsTemplate{
		Scaling: 8,
		Settings: Settings{
			Width:        60,
			FontHeight:   8,
			Height:       120,
			StrokeColorS: "black",
			StrokeWidth:  floatP(1),
			Margins:      floatP(2),
			FontPath:     "assets/freesans.ttf",
			FontColorS:   "white",
			BorderColorS: "green",
			BorderWidth:  floatP(2),
		},
		Cards: []Card{
			{
				Settings: Settings{
					BackgroundColor: stringP("white"),
				},
				Areas: []Counter{
					{
						Settings: Settings{
							BorderColorS:    "red",
							BorderWidth:     floatP(2),
							Height:          50,
							BackgroundColor: stringP("black"),
						},
						Texts: []Text{
							{
								Settings: Settings{StrokeWidth: floatP(1)},
								String:   "Area text",
							},
						},
						Images: []Image{
							{
								Settings: Settings{Position: 3},
								Path:     "assets/binoculars.png",
								Scale:    0.2,
							},
						},
					},
					{
						Settings: Settings{
							BackgroundColor: stringP("lightgrey"),
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
						Settings: Settings{
							StrokeColorS: "black",
						},
						String: "Full card text",
					},
				},
			},
		},
	}

	template.Settings.ApplySettingsScaling(template.Scaling)
	err := template.ApplyCardWaterfallSettings()
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}

	// f, _ := os.Create("/tmp/card_01.png")
	// defer f.Close()
	// err = template.Cards[0].EncodeImage(f, template)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	t.Skip("TODO: implement test")
	assert.Equal(t, 120, template.Settings.Width)
	assert.Equal(t, 240, template.Settings.Height)
	assert.Equal(t, 120, template.Cards[0].Settings.Width)
	assert.Equal(t, 240, template.Cards[0].Settings.Height)
	assert.Equal(t, 120, template.Cards[0].Areas[0].Width)
	assert.Equal(t, 100, template.Cards[0].Areas[0].Height)
	assert.Equal(t, float64(2), *template.Settings.StrokeWidth)
	assert.Equal(t, 120, template.Cards[0].Images[0].Width)
	assert.Equal(t, 240, template.Cards[0].Images[0].Height)
	assert.Equal(t, 120, template.Cards[0].Texts[0].Settings.Width)
	assert.Equal(t, 240, template.Cards[0].Texts[0].Settings.Height)
	assert.Equal(t, float64(2), *template.Cards[0].Texts[0].Settings.StrokeWidth)

	canvas, err := template.Cards[0].ToCanvas(template)
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}
	if assert.NoError(t, err) {
		TestImageContent(t, "testdata/card_01.png", 5643, canvas)
	}
}
