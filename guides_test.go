package counters

import (
	"image/color"
	"testing"

	"github.com/fogleman/gg"
	"github.com/stretchr/testify/assert"
)

func TestGuides(t *testing.T) {
	sideSize := 300
	testText := Text{
		Settings: Settings{
			FontColor:  color.Black,
			FontHeight: 15,
			Width:      sideSize,
			Height:     sideSize,
			Margins:    floatP(3),
		},
		String: "Guides",
	}

	template := CounterTemplate{
		Settings: Settings{
			FontPath: "assets/freesans.ttf",
			Width:    100,
			Height:   100,
		},
		DrawGuides: true,
		Counters: []Counter{
			{Texts: []Text{testText}},
		},
	}

	parsedTemplate, err := template.ParsePrototype()
	if assert.NoError(t, err) {
		testCanvas := gg.NewContext(sideSize, sideSize)

		testCanvas.Push()
		testCanvas.SetColor(color.White)
		testCanvas.DrawRectangle(0, 0, float64(sideSize), float64(sideSize))
		testCanvas.Fill()
		testCanvas.Pop()

		canvas, err := parsedTemplate.Counters[0].Canvas(true)
		if assert.NoError(t, err) {
			TestImageContent(t, "testdata/guides_01.png", 1709, canvas)
		}
	}

}
