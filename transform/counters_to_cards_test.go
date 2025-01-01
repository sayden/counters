package transform

import (
	"testing"

	"github.com/sayden/counters"
	"github.com/stretchr/testify/assert"
)

func TestCountersToCards(t *testing.T) {
	// Setup test data
	counterTemplate := &counters.CounterTemplate{
		Counters: []counters.Counter{
			{
				Settings: counters.Settings{Multiplier: intP(1)},
				Texts:    counters.Texts{{String: "Test Counter 1"}},
			},
			{
				Metadata:ounters.Extra{SkipCardGeneration: true},
				Settings: counters.Settings{Multiplier: intP(2)},
				Texts:    counters.Texts{{String: "Test Counter 2"}},
			},
		},
	}

	cardTemplate := &counters.CardsTemplate{
		Settings: counters.Settings{
			Width:  800,
			Height: 1200,
		},
	}

	cardTransformer := &QuotesToCardTransformer{
		Quotes:         counters.Quotes{{Origin: "Origin", Quote: "Test Quote"}},
		IndexForTitles: counterTemplate.PositionNumberForFilename,
	}

	cfg := &CountersToCardsConfig{
		CountersTemplate:   counterTemplate,
		CardTemplate:       cardTemplate,
		CounterTransformer: cardTransformer,
	}

	// Call the method
	output, err := cfg.CountersToCards()

	// Assertions
	if assert.NoError(t, err) {
		assert.NotNil(t, output)
		assert.Equal(t, 1, len(output.Cards))
		assert.Equal(t, 2, len(output.Cards[0].Areas))
	}
}
