package transform

import (
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
)

// CountersToCardsConfig is a configuration to take a counter template and convert it to a card template
type CountersToCardsConfig struct {
	// An incoming template of counters, to be converted to cards
	CountersTemplate *counters.CounterTemplate

	// The template of cards to be used as the base for the conversion
	CardTemplate *counters.CardsTemplate

	// The transformer to be used to convert the counters to cards
	CounterTransformer counters.CounterToCardTransformer
}

// CountersToCards using the provided input writes png images into files in batches of 70 (70 cards per file).
func (cfg *CountersToCardsConfig) CountersToCards() (*counters.CardsTemplate, error) {
	var outputCardTemplate counters.CardsTemplate
	err := mergo.Merge(&outputCardTemplate, cfg.CardTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "error trying to merge card templates")
	}

	// INPUT
	for _, counter := range cfg.CountersTemplate.Counters {
		if counter.Metadata != nil && counter.Metadata.SkipCardGeneration {
			// skip card generation
			continue
		}

		// convert counter to fow card
		card, err := cfg.CounterTransformer.ToNewCard(&counter, cfg.CardTemplate)
		if err != nil {
			return nil, errors.Wrap(err, "error trying to create card")
		}

		if card == nil {
			return nil, errors.New("error creating card")
		}

		for i := 0; i < *counter.Multiplier; i++ {
			outputCardTemplate.Cards = append(outputCardTemplate.Cards, *card)
		}
	}

	return &outputCardTemplate, nil
}
