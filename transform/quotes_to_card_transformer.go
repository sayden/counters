package transform

import (
	"math/rand"

	"github.com/creasty/defaults"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
)

type QuotesToCardTransformer struct {
	Quotes         []counters.Quote
	IndexForTitles int
}

// TODO It seems that 'destination' field can be omitted
func (w *QuotesToCardTransformer) ToNewCard(cc *counters.Counter, inputCardTemplate *counters.CardsTemplate) (card *counters.Card, err error) {
	//Set a random quote
	quote := w.Quotes[rand.Intn(len(w.Quotes))]

	card, err = w.getCardPrototype()
	if err != nil {
		return nil, errors.Wrap(err, "could not read card prototype")
	}
	mergo.Merge(&card.Settings, inputCardTemplate.Settings)

	mergo.Merge(&cc.Settings, card.Settings)
	card.Areas = w.getCardAreas(cc, quote)

	return card, nil
}

func (w *QuotesToCardTransformer) getCardAreas(cc *counters.Counter, q counters.Quote) []counters.Counter {
	margins := 0.0
	backgroundColor := ""
	cc.Texts = nil
	cc.Margins = &margins
	cc.BackgroundColor = &backgroundColor

	// Modify images in incoming counter
	xShift := 0.0
	yShift := 0.0
	for i, img := range cc.Images {
		cc.Images[i].YShift = &yShift
		cc.Images[i].XShift = &xShift

		switch img.Position {
		case 0:
			// Use a different image
			if cc.Metadata.CardImage.Path != "" {
				cc.Images[i].Path = cc.Metadata.CardImage.Path
			}

			// Change the scaling
			if cc.Metadata.CardImage.ImageScaling != "" {
				cc.Images[i].Path = cc.Metadata.CardImage.ImageScaling
			}

			// Change the scale in case it is too small or too big
			if cc.Metadata.CardImage.Scale != 0 {
				cc.Images[i].Scale = cc.Metadata.CardImage.Scale
			}
		}
	}

	images, texts := w.getDownAreaCounterItems(cc, q)
	return []counters.Counter{
		*cc,
		counters.Counter{
			Images: images,
			Texts:  texts,
			Frame:  true,
		},
	}
}

func (w *QuotesToCardTransformer) getDownAreaCounterItems(cc *counters.Counter, q counters.Quote) ([]counters.Image, []counters.Text) {
	return []counters.Image{{Path: "assets/old_paper.jpg"}},
		[]counters.Text{
			{
				Settings: counters.Settings{
					AvoidClipping: true,
					StrokeWidth:   floatP(5),
					Position:      3,
				},
				String: cc.GetTextInPosition(w.IndexForTitles),
			}, {
				Settings: counters.Settings{
					FontPath:    "assets/freesans.ttf",
					StrokeWidth: floatP(0),
					FontHeight:  20,
					FontColorS:  "black",
				},
				String: q.Quote,
			}, {
				Settings: counters.Settings{
					StrokeWidth:   floatP(0),
					FontHeight:    30,
					FontColorS:    "black",
					YShift:        floatP(-70),
					AvoidClipping: true,
					Position:      9,
				},
				String: " -" + q.Origin,
			},
		}
}

func (q *QuotesToCardTransformer) getCardPrototype() (*counters.Card, error) {
	c := counters.Card{}
	if err := defaults.Set(&c); err != nil {
		return nil, errors.Wrap(err, "could not set defaults to card")
	}

	c.BorderColorS = "black"
	c.BackgroundColor = stringP("#faebd7")
	c.FontColorS = "white"
	c.ImageScaling = "fitHeight"
	c.StrokeColorS = "black"
	c.StrokeWidth = floatP(0)

	return &c, nil
}
