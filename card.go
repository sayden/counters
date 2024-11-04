package counters

import (
	"image"
	"image/color"
	"io"
	"math"

	"github.com/fogleman/gg"
	"github.com/pkg/errors"
)

type Card struct {
	Settings
	Areas  []Counter `json:"areas,omitempty"`
	Texts  Texts     `json:"texts,omitempty"`
	Images Images    `json:"images,omitempty"`
}

func (c *Card) Image(template *CardsTemplate) (image.Image, error) {
	cardCanvas, err := c.ToCanvas(template)
	if err != nil {
		return nil, err
	}

	return cardCanvas.Image(), nil
}

func (c *Card) ToCanvas(template *CardsTemplate) (*gg.Context, error) {
	Mergev2(&c.Settings, &template.Settings)
	SetColors(&c.Settings)

	cardCanvas, err := template.Canvas(&c.Settings, template.Width, template.Height)
	if err != nil {
		return nil, err
	}

	if err = c.DrawBackgroundImage(cardCanvas); err != nil {
		return nil, err
	}

	if err = c.Images.DrawImagesOnCanvas(&c.Settings, cardCanvas, c.Width, c.Height); err != nil {
		return nil, err
	}
	// Height when all areas have the same height
	numberOfAreas := len(c.Areas)

	//calculatedAreaHeight := (float64(template.Height) - (template.Margins * 2)) / float64(numberOfAreas)

	areasHeights := c.GetAreasHeights()
	// Process each area on the text
	var y float64
	if template.Margins != nil {
		y = *template.Margins
	}

	for areaIndex, areaCounter := range c.Areas {
		isLastAreaOfCard := areaIndex != numberOfAreas
		c.Areas[areaIndex].Height = int(math.Floor(areasHeights[areaIndex]))

		areaCanvas, err := c.ProcessAreav2(&areaCounter, template, int(math.Floor(areasHeights[areaIndex])), isLastAreaOfCard)
		if err != nil {
			return nil, err
		}

		var x float64
		if template.Margins != nil {
			x = *template.Margins
		}
		if err = c.drawOnCard(template, cardCanvas, areaCanvas, x, y); err != nil {
			return nil, err
		}
		y += areasHeights[areaIndex]
	}

	if err = c.Texts.DrawTextsOnCanvas(&c.Settings, cardCanvas, c.Width, c.Height); err != nil {
		return nil, err
	}

	c.maybeDrawBorders(cardCanvas)

	return cardCanvas, nil
}

func (c *Card) EncodeImage(w io.Writer, t *CardsTemplate) error {
	cardCanvas, err := c.ToCanvas(t)
	if err != nil {
		return err
	}

	if err = cardCanvas.EncodePNG(w); err != nil {
		return err
	}

	return nil
}

func (c *Card) ProcessAreav2(area *Counter, t *CardsTemplate, calculatedAreaHeight int, isLastArea bool) (*gg.Context, error) {
	Mergev2(&area.Settings, &c.Settings)

	margins := 0
	if t.Margins != nil {
		margins = int(*t.Margins * 2)
	}
	area.Width = (t.Width) - margins
	area.Height = calculatedAreaHeight

	areaCanvas, err := area.Canvas(false)
	if err != nil {
		return nil, err
	}

	if !isLastArea && area.Frame {
		c.drawFrame(areaCanvas, *c.BorderWidth, c.BorderColor)
	}

	return areaCanvas, nil
}

func (c *Card) GetAreasHeights() (hs []float64) {
	hs = make([]float64, len(c.Areas))
	availableH := float64(c.Height) - (*c.Margins * 2)
	hasCustomHeight := make([]bool, len(c.Areas))
	totalNonCustomAreas := 0
	for i, area := range c.Areas {
		if area.Height == 0 {
			totalNonCustomAreas++
			continue
		}
		if area.Height != c.Height {
			hasCustomHeight[i] = true
			availableH -= float64(area.Height)
			continue
		}
		totalNonCustomAreas++
	}

	availableSpaceForNonCustom := availableH / float64(totalNonCustomAreas)
	for i, isCustom := range hasCustomHeight {
		if isCustom {
			hs[i] = float64(c.Areas[i].Height)
		} else {
			hs[i] = availableSpaceForNonCustom
		}
	}

	return
}

func (c *Card) drawOnCard(t *CardsTemplate, cardCanvas, areaCanvas *gg.Context, x, y float64) error {
	cardCanvas.DrawImage(areaCanvas.Image(), int(x), int(y))

	if t.DrawGuides {
		guidesImage, err := DrawGuides(&c.Settings)
		if err != nil {
			return errors.Wrap(err, "error tyring to draw guides")
		}
		cardCanvas.DrawImage(*guidesImage, int(x), int(y))
	}

	return nil
}

func (c *Card) drawFrame(areaCanvas *gg.Context, w float64, col color.Color) {
	areaCanvas.Push()
	areaCanvas.SetColor(col)
	areaCanvas.SetLineWidth(w)
	frameX := float64(areaCanvas.Width())
	frameY := float64(areaCanvas.Height())
	areaCanvas.DrawRectangle(0, 0, frameX, frameY)
	areaCanvas.Stroke()
	areaCanvas.Pop()
}

func (c *Card) maybeDrawBorders(cardCanvas *gg.Context) {
	borderColorIsSet := c.Settings.BorderColor != nil
	borderWidthIsSet := c.Settings.BorderWidth != nil && *c.Settings.BorderWidth != 0

	if borderColorIsSet && borderWidthIsSet {
		cardCanvas.Push()
		cardCanvas.SetColor(c.Settings.BorderColor)
		cardCanvas.SetLineWidth(*c.Settings.BorderWidth)
		cardCanvas.DrawRectangle(0, 0, float64(c.Settings.Width), float64(c.Settings.Height))
		cardCanvas.Stroke()
		cardCanvas.Pop()
	}
}
