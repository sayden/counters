package counters

import (
	"encoding/json"
	"slices"

	"github.com/creasty/defaults"
	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	deepcopy "github.com/qdm12/reprint"
)

// CardsTemplate is the template sheet (usually A4) to place cards on top in grid fashion
type CardsTemplate struct {
	Settings

	Rows    int `json:"rows,omitempty" default:"8"`
	Columns int `json:"columns,omitempty" default:"5"`

	DrawGuides bool `json:"draw_guides,omitempty"`

	// TODO is this field still used? Mode can be 'tiles' or 'template' to generate an A4 sheet
	// like of cards or a single file per card. I don't think it should be removed, because it's
	// necessary for printing or TTS
	Mode string `json:"mode,omitempty" default:"tiles"`

	// TODO Rename this to OutputFolder or the one in counters to OutputPath and update JSON's
	OutputPath string `json:"output_path,omitempty" default:"output_%02d"`

	Scaling float64 `json:"scaling,omitempty" default:"1.0"`

	Cards           []Card `json:"cards"`
	MaxCardsPerFile int    `json:"max_cards_per_file,omitempty"`

	Prototypes map[string]CardPrototype `json:"prototypes,omitempty"`
	Extra      CardsExtra               `json:",omitempty"`
}

// CardsExtra is a container for extra information used in different projects but that they are not
// common to all of them
type CardsExtra struct {
	FactionImage      string  `json:"faction_image,omitempty"`
	FactionImageScale float64 `json:"faction_image_scale,omitempty"`
	BackImage         string  `json:"back_image,omitempty"`
}

func ParseCardTemplate(byt []byte) (*CardsTemplate, error) {
	err := ValidateSchemaBytes[CardsTemplate](byt)
	if err != nil {
		return nil, errors.Wrap(err, "JSON file is not valid")
	}

	t := CardsTemplate{}
	if err = json.Unmarshal(byt, &t); err != nil {
		return nil, err
	}

	if t.Scaling > 0 {
		t.Settings.ApplySettingsScaling(t.Scaling)
	}

	err = t.ApplyCardWaterfallSettings()
	if err != nil {
		return nil, errors.Wrap(err, "could not apply card waterfall settings")
	}

	newTemplate, err := t.ParsePrototype()
	if err != nil {
		return nil, errors.Wrap(err, "could not parse prototype")
	}

	return newTemplate, nil
}

// ApplyCardWaterfallSettings traverses the cards in the template applying the default settings to
// value that are zero-valued
func (t *CardsTemplate) ApplyCardWaterfallSettings() error {
	SetColors(&t.Settings)

	for cardIdx := range t.Cards {
		err := defaults.Set(&t.Cards[cardIdx].Settings)
		if err != nil {
			return errors.Wrap(err, "could not set defaults for card settings")
		}
		card := &t.Cards[cardIdx]
		if t.Scaling > 0 {
			card.ApplySettingsScaling(t.Scaling)
		}
		err = Mergev2(&card.Settings, &t.Settings)
		if err != nil {
			return err
		}

		for areaIdx := range card.Areas {
			area := &card.Areas[areaIdx]
			if t.Scaling > 0 {
				area.ApplySettingsScaling(t.Scaling)
			}
			err := Mergev2(&area.Settings, &card.Settings)
			if err != nil {
				return err
			}

			for imageIdx := range area.Images {
				image := &area.Images[imageIdx]
				if t.Scaling > 0 {
					image.ApplySettingsScaling(t.Scaling)
				}
				err := Mergev2(&image.Settings, &area.Settings)
				if err != nil {
					return err
				}
			}

			for textIdx := range area.Texts {
				text := &area.Texts[textIdx]
				if t.Scaling > 0 {
					text.ApplySettingsScaling(t.Scaling)
				}
				err := Mergev2(&text.Settings, &area.Settings)
				if err != nil {
					return err
				}
			}
		}

		for imageIdx := range card.Images {
			image := &card.Images[imageIdx]
			if t.Scaling > 0 {
				image.ApplySettingsScaling(t.Scaling)
			}
			err := Mergev2(&image.Settings, &card.Settings)
			if err != nil {
				return err
			}
		}

		for textIdx := range card.Texts {
			text := &card.Texts[textIdx]
			if t.Scaling > 0 {
				text.ApplySettingsScaling(t.Scaling)
			}
			err := Mergev2(&text.Settings, &card.Settings)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *CardsTemplate) SheetCanvas() (*gg.Context, error) {
	width := t.Columns * t.Width
	height := t.Rows * t.Height
	return t.Canvas(&t.Settings, width, height)
}

// Canvas returns a Canvas with attributes (like background color or size)
// taken from `settings`
func (t *CardsTemplate) Canvas(settings *Settings, width, height int) (*gg.Context, error) {
	dc := gg.NewContext(width, height)
	if err := settings.LoadFontOrDefault(dc); err != nil {
		return nil, err
	}

	if settings.BgColor != nil {
		dc.Push()
		dc.SetColor(settings.BgColor)
		dc.DrawRectangle(0, 0, float64(width), float64(height))
		dc.Fill()
		dc.Pop()
	}

	if settings.FontColorS != "" {
		ColorFromStringOrDefault(settings.FontColorS, t.BgColor)
	}

	return dc, nil
}

func (c *CardsTemplate) DuplicateCard(card *CardPrototype, name string) ([]Card, error) {
	// If an area does not contains text or image prototypes, it's a static area, ie.
	// it must be present in every card.
	// If an area contains text of image prototypes, then the number of prototypes
	// reflects the amount of cards that will be generated
	if len(card.Areas) == 0 {
		return nil, errors.New("card has no areas defined")
	}

	totalCards := 0
	var err error
	for _, area := range card.Areas {
		if len(area.TextPrototypes) != 0 {
			if totalCards, err = area.isLengthConsistent(); err != nil {
				return nil, errors.Wrap(err, "could not check length consistency for area")
			} else if totalCards != 0 {
				break
			}
		} else if len(area.ImagePrototypes) != 0 {
			if totalCards, err = area.isLengthConsistent(); err != nil {
				return nil, errors.Wrap(err, "could not check length consistency for area")
			} else if totalCards != 0 {
				break
			}
		}
	}
	if totalCards == 0 {
		return nil, errors.New("card has no areas defined")
	}

	cards := make([]Card, 0, totalCards)
	for cardNumber := range totalCards {
		var cardCopy Card
		if err := deepcopy.FromTo(card.Card, &cardCopy); err != nil {
			return nil, errors.Wrap(err, "could not copy card")
		}
		err := defaults.Set(&cardCopy.Settings)
		if err != nil {
			return nil, errors.Wrap(err, "could not set defaults for card copy")
		}

		staticAreas := make([]Counter, 0, len(card.Areas))
		cardCopy.Areas = make([]Counter, 0, len(card.Areas))
		protoAreas := make([]CounterPrototype, 0, len(card.Areas))
		for _, area := range card.Areas {
			if area.TextPrototypes == nil && area.ImagePrototypes == nil {
				area.Counter.Settings = area.Settings
				area.Counter.Metadata = area.Metadata
				staticAreas = append(staticAreas, area.Counter)
				continue
			}

			protoAreas = append(protoAreas, area)
		}

		for _, protoArea := range protoAreas {
			var newArea Counter
			if err := defaults.Set(&newArea.Settings); err != nil {
				return nil, errors.Wrap(err, "could not set defaults for new area")
			}
			if err := deepcopy.FromTo(protoArea.Counter, &newArea); err != nil {
				return nil, errors.Wrap(err, "could not copy area counter")
			}
			if err := deepcopy.FromTo(protoArea.Counter.Settings, &newArea.Settings); err != nil {
				return nil, errors.Wrap(err, "could not copy area counter")
			}
			if err := deepcopy.FromTo(protoArea.Metadata, &newArea.Metadata); err != nil {
				return nil, errors.Wrap(err, "could not copy area metadata")
			}
			if err = protoArea.applyPrototypes(&newArea, cardNumber); err != nil {
				return nil, err
			}
			staticAreas = append(staticAreas, newArea)
		}

		cardCopy.Areas = append(cardCopy.Areas, staticAreas...)

		// reorder the areas by their ordering
		if err := sortAreas(cardCopy.Areas); err != nil {
			return nil, errors.Wrap(err, "could not reorder areas")
		}

		cards = append(cards, cardCopy)
	}

	return cards, nil
}

func sortAreas(areas []Counter) error {
	if areas == nil || len(areas) == 0 {
		return nil
	}

	slices.SortFunc(areas, func(i Counter, j Counter) int {
		if i.Metadata.Ordering < j.Metadata.Ordering {
			return -1
		} else if i.Metadata.Ordering > j.Metadata.Ordering {
			return 1
		}
		return 0
	})

	return nil
}

func (c *CardsTemplate) ParsePrototype() (*CardsTemplate, error) {
	if c.Prototypes != nil {
		for name, cardProto := range c.Prototypes {
			cards, err := c.DuplicateCard(&cardProto, name)
			if err != nil {
				return nil, errors.Wrapf(err, "could not duplicate card prototype %s", name)
			}
			c.Cards = append(c.Cards, cards...)
		}
	}

	return c, nil
}
