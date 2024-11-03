package counters

import (
	"fmt"
	"image/color"

	"github.com/danielgtaylor/unistyle"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/pkg/errors"
)

type Text struct {
	Settings
	Positioner

	String string `json:"string,omitempty"`

	Underline bool `json:"underline,omitempty"`

	TextBackgroundColor string      `json:"text_background_color,omitempty"`
	TextBgColor         color.Color `json:"-"`
}

func (t Text) GetSettings() *Settings {
	return &t.Settings
}

type Texts []Text

// DrawTextsOnCanvas draws the texts provided on areaCanvas at positions `w` and `h`
// using the provided Settings
func (texts Texts) DrawTextsOnCanvas(s *Settings, areaCanvas *gg.Context, w, h int) error {
	for _, text := range texts {
		Mergev2(&text.Settings, s)

		text.Width = w
		text.Height = h
		err := text.Draw(areaCanvas, text.Position, &text.Settings)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Text) GetAlignment() gg.Align {
	switch t.Alignment {
	case ALIGMENT_CENTER:
		return gg.AlignCenter
	case ALIGMENT_RIGHT:
		return gg.AlignRight
	default:
		return gg.AlignLeft
	}
}

func (t *Text) Draw(dc *gg.Context, pos int, settings *Settings) error {
	dc.Push()
	defer dc.Pop()

	// Font and font color
	err := dc.LoadFontFace(settings.FontPath, settings.FontHeight)
	if err != nil {
		return fmt.Errorf("could not load font face '%s' with heigth '%f': %w", settings.FontPath, settings.FontHeight, err)
	}
	dc.SetColor(settings.FontColor)

	if t.Underline {
		t.String = unistyle.Underline(t.String, unistyle.UnderlineLine)
	}

	x, _, ax, ay, maxWidthForPosition, err := t.getTextDimensions(pos, settings)
	if err != nil {
		return errors.Wrap(err, "could not get dimensions to draw")
	}
	if settings.AvoidClipping {
		maxWidthForPosition = float64(settings.Width)
	}

	// Just take all the possible space to write the text, later it will be cropped to just the necessary by using the
	// CropToContent function
	temp := gg.NewContext(dc.Width(), dc.Height())

	if err = temp.LoadFontFace(settings.FontPath, settings.FontHeight); err != nil {
		return err
	}
	temp.SetColor(settings.FontColor)

	centerVertical := float64(dc.Height()) / 2

	if *settings.StrokeWidth != 0 {
		drawTextWrappedWithStroke(t.String, *settings.StrokeWidth, centerVertical, 0, 0.5, maxWidthForPosition, temp, settings.FontColor, *settings.StrokeWidth, settings.StrokeColor, t.GetAlignment())
	} else {
		temp.DrawStringWrapped(t.String, 0, centerVertical, 0, 0.5, maxWidthForPosition, 2, t.GetAlignment())
	}

	img := temp.Image()

	if *settings.Rotation != 0 {
		img = imaging.Rotate(img, *settings.Rotation, color.Transparent)
	}

	img = CropToContent(img)

	x, y, ax, ay, err := t.getObjectPositions(pos, settings)
	if err != nil {
		return errors.Wrap(err, "could not get a correct position")
	}

	if *settings.ShadowDistance != 0 {
		shadow := getShadowFromImage(img, *t.ShadowDistance, *t.ShadowSigma)
		shadowCtx := gg.NewContextForImage(shadow)
		shadowCtx.DrawImageAnchored(img, (shadowCtx.Width() / 2), (shadowCtx.Height() / 2), 0.5, 0.5)
		img = shadowCtx.Image()
	}

	if t.TextBackgroundColor != "" {
		bgTemp := gg.NewContext(img.Bounds().Dx()+int(settings.FontHeight*0.40), img.Bounds().Dy()+int(settings.FontHeight*0.40))
		bgTemp.Push()
		bgTemp.SetColor(ColorFromStringOrDefault(t.TextBackgroundColor, t.TextBgColor))
		bgTemp.DrawRectangle(0, 0, float64(dc.Width()), float64(dc.Height()))
		bgTemp.Fill()
		bgTemp.Pop()

		if err != nil {
			return errors.Wrap(err, "could not get dimensions to draw")
		}
		bgTemp.DrawImageAnchored(img, int(settings.FontHeight*0.40/2), int(settings.FontHeight*0.40/2), 0, 0)
		img = bgTemp.Image()
	}

	dc.DrawImageAnchored(img, int(x), int(y), ax, ay)

	return nil
}

func (t *Text) getTextDimensions(pos int, def *Settings) (float64, float64, float64, float64, float64, error) {
	ax, ay, maxWidth, err := t.GetAnchorPointsAndMaxWidth(pos, def)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}

	x, y, err := t.GetXYPosition(pos, def)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}

	return x, y, ax, ay, maxWidth, nil
}

func drawTextWrappedWithStroke(t string, x, y, ax, ay, w float64, temp *gg.Context, textColor color.Color, strokeSize float64, strokeColor color.Color, align gg.Align) {
	temp.Push()
	defer temp.Pop()

	//Draw stroke
	if strokeColor != nil {
		temp.SetColor(strokeColor)
		for dy := -strokeSize; dy <= strokeSize; dy++ {
			for dx := -strokeSize; dx <= strokeSize; dx++ {
				if dx*dx+dy*dy >= strokeSize*strokeSize {
					// give it rounded corners
					continue
				}
				x := x + dx
				y := y + dy
				temp.DrawStringWrapped(t, x, y, ax, ay, w, 2, align)
			}
		}
	}

	//Draw text
	temp.SetColor(textColor)
	temp.DrawStringWrapped(t, x, y, ax, ay, w, 2, align)
}
