package transform

import (
	"testing"

	"github.com/sayden/counters"
	"github.com/stretchr/testify/assert"
)

func TestToNewCard(t *testing.T) {
	quotes := []counters.Quote{
		{Quote: "Test Quote 1", Origin: "Author 1"},
	}

	counter := counters.Counter{
		Metadata: &counters.Metadata{
			CardImage: &counters.Image{
				Path:     "assets/binoculars.jpg",
				Settings: counters.Settings{ImageScaling: "fitWidth"},
				Scale:    1.5,
			},
		},
		Images: []counters.Image{
			{
				Settings: counters.Settings{Position: 0},
				Path:     "assets/binoculars.jpg",
			},
		},
	}

	cardTemplate := &counters.CardsTemplate{
		Settings: counters.Settings{FontPath: "assets/freesans.ttf"}}

	builder := &QuotesToCardTransformer{
		Quotes:         quotes,
		IndexForTitles: 0,
	}

	card, err := builder.ToNewCard(&counter, cardTemplate)
	assert.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, "#faebd7", *card.BackgroundColor)
	assert.Equal(t, "black", card.BorderColorS)
	assert.Equal(t, "white", card.FontColorS)
	assert.Equal(t, "fitHeight", card.ImageScaling)
	assert.Equal(t, "black", card.StrokeColorS)
	assert.Equal(t, float64(0), *card.StrokeWidth)
	assert.NotEmpty(t, card.Areas)
	assert.Equal(t, "Test Quote 1", card.Areas[1].Texts[1].String)
	assert.Equal(t, " -Author 1", card.Areas[1].Texts[2].String)
}
