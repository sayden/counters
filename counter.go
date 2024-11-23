package counters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"strings"
	"sync"
	"text/template"

	"dario.cat/mergo"
	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"github.com/thehivecorporation/log"
)

func init() {
	var err error
	pieceSlotTemplate, err = template.New("piece_text").Parse(Template_NewVassalPiece)
	if err != nil {
		panic(fmt.Errorf("could not parse template string %w", err))
	}
}

var pieceSlotTemplate *template.Template

// Counter is POGO-like holder for data needed for other parts to fill and draw
// a counter in a container
type Counter struct {
	Settings

	SingleStep bool `json:"single_step,omitempty"`
	Frame      bool `json:"frame,omitempty"`

	Images Images `json:"images,omitempty"`
	Texts  Texts  `json:"texts,omitempty"`
	Extra  *Extra `json:"extra,omitempty"`

	// Generate the following counter with 'back' suffix in its filename
	Back *Counter `json:"back,omitempty"`

	Filename      string     `json:"filename,omitempty"`
	PrototypeName string     `json:"-"`
	VassalPiece   *PieceSlot `json:"vassal,omitempty"`
}

type Counters []Counter

// TODO This Extra contains data from all projects
type Extra struct {
	// PublicIcon in a FOW counter is the visible icon for the enemy. Imagine an icon for the back
	// of a block in a Columbia game
	CardImage          *Image `json:"card_image,omitempty"`
	Cost               int    `json:"cost,omitempty"`
	PublicIcon         *Image `json:"public_icon,omitempty"`
	Side               string `json:"side,omitempty"`
	SkipCardGeneration bool   `json:"skip_card_generation,omitempty"`
	Title              string `json:"title,omitempty"`
	TitlePosition      *int   `json:"title_position,omitempty"`
}

type ImageExtraData struct {
	// the path to find the image file
	Path string `json:"path,omitempty"`

	// a percentage of original image's size
	Scale float64 `json:"scale,omitempty"`

	// none, fitHeight, fitWidth or wrap
	ImageScaling string `json:"image_scaling,omitempty"`
}

func (c *Counter) GetTextInPosition(i int) string {
	for _, text := range c.Texts {
		if text.Position == i {
			return text.String
		}
	}

	return ""
}

// filenumber: CounterTemplate.PositionNumberForFilename. So it will always be fixed number
// position: The position of the text in the counter (0-16)
// suffix: A suffix on the file. Constant
func (c *Counter) GenerateCounterFilename(sideName string, position int, filenamesInUse *sync.Map) {
	if c.Filename != "" {
		return
	}

	var b strings.Builder
	var name string
	name = c.GetTextInPosition(position)

	if c.Extra != nil {
		if c.Extra.TitlePosition != nil && *c.Extra.TitlePosition != position {
			name = c.GetTextInPosition(*c.Extra.TitlePosition)
		}
		if name != "" {
			b.WriteString(name + " ")
		}
		// This way, the positional based name will always be the first part of the filename
		// while the manual title will come later. This is useful when using prototypes so that
		// counters with the same positional name are close together in the destination folder
		name = ""

		if c.Extra.Side != "" {
			b.WriteString(c.Extra.Side)
			b.WriteString(" ")
		}

		if c.Extra.Title != "" {
			b.WriteString(c.Extra.Title)
			b.WriteString(" ")
		}
	}

	if name != "" {
		b.WriteString(name + " ")
	}

	if c.PrototypeName != "" {
		b.WriteString(c.PrototypeName + " ")
	}

	res := b.String()

	res = strings.TrimSpace(res)

	c.PrettyName = res

	if sideName != "" {
		res = sideName + "_" + res
	}

	filenumber := 0
	_, isFound := filenamesInUse.Load(res)
	if isFound {
		for {
			tempRes := fmt.Sprintf("%s_%04d", res, filenumber)
			_, isFound = filenamesInUse.Load(tempRes)
			if !isFound {
				res = tempRes
				break
			}
			filenumber++
		}
	}

	if res == "" {
		res = fmt.Sprintf("%04d", filenumber)
	}
	res = strings.TrimSpace(res)

	filenamesInUse.Store(res, true)

	res += ".png"
	c.Filename = res
}

func (c *Counter) Canvas(withGuides bool) (*gg.Context, error) {
	SetColors(&c.Settings)
	canvas, err := c.canvas()
	if err != nil {
		return nil, err
	}

	// Draw background image
	if err = c.DrawBackgroundImage(canvas); err != nil {
		return nil, errors.Wrap(err, "error trying to draw background image")
	}

	// Draw images
	if err = c.Images.DrawImagesOnCanvas(&c.Settings, canvas, c.Width, c.Height); err != nil {
		return nil, errors.Wrap(err, "error trying to process image")
	}

	// Draw texts
	if err = c.Texts.DrawTextsOnCanvas(&c.Settings, canvas, c.Width, c.Height); err != nil {
		return nil, errors.Wrap(err, "error trying to draw text")
	}

	// Draw guides
	if withGuides {
		guides, err := DrawGuides(&c.Settings)
		if err != nil {
			return nil, err
		}
		canvas.DrawImage(*guides, 0, 0)
	}

	// Draw borders
	if c.BorderWidth != nil && *c.BorderWidth > 0 {
		c.drawBorders(canvas)
	}

	return canvas, nil
}

func (c *Counter) EncodeCounter(w io.Writer, drawGuides bool) error {

	counterCanvas, err := c.Canvas(drawGuides)
	if err != nil {
		return err
	}

	return counterCanvas.EncodePNG(w)
}

func (a *Counter) drawBorders(canvas *gg.Context) {
	canvas.Push()
	canvas.SetColor(a.Settings.BorderColor)
	canvas.SetLineWidth(*a.Settings.BorderWidth)
	canvas.DrawRectangle(0, 0, float64(a.Settings.Width), float64(a.Settings.Height))
	canvas.Stroke()
	canvas.Pop()
}

func (c *Counter) canvas() (*gg.Context, error) {
	dc := gg.NewContext(c.Width, c.Height)
	if err := dc.LoadFontFace(c.FontPath, c.FontHeight); err != nil {
		log.WithFields(log.Fields{"font": "'" + c.FontPath + "'", "height": c.FontHeight}).Error(err)
		return nil, err
	}

	dc.Push()
	dc.SetColor(c.BgColor)
	dc.DrawRectangle(0, 0, float64(c.Width), float64(c.Height))
	dc.Fill()
	dc.Pop()

	if c.FontColorS != "" && c.FontColor == nil {
		ColorFromStringOrDefault(c.FontColorS, c.FontColor)
	}

	return dc, nil
}

func (c *Counter) ToVassal(sideName string) error {
	if c.Filename == "" {
		return errors.New("vassal: counter filename is empty")
	}
	if sideName == "" {
		return errors.New("vassal: side name is empty")
	}

	backFilename := strings.TrimSuffix(c.Filename, path.Ext(c.Filename)) + "_back.png"
	buf := bytes.NewBufferString("")
	pieceTemp := PieceTemplateData{
		FrontFilename: c.Filename,
		BackFilename:  backFilename,
		FlipName:      sideName,
		PieceName:     c.PrettyName,
		Id:            c.Filename,
	}

	err := pieceSlotTemplate.ExecuteTemplate(buf, "piece_text", pieceTemp)
	if err != nil {
		return fmt.Errorf("could not execute template %w", err)
	}

	data := buf.String()
	data = strings.ReplaceAll(data, "&#x9;", "\t")
	piece := PieceSlot{
		EntryName: c.PrettyName,
		Height:    c.Height,
		Width:     c.Width,
		Data:      data,
	}

	c.VassalPiece = &piece

	return nil
}

func (c *Counter) mergeFrontAndBack() (*Counter, error) {
	if err := mergo.Merge(c.Back, c); err != nil {
		return nil, fmt.Errorf("could not merge back and front counter: %w", err)
	}
	c.Back.Back = nil

	if c.PrettyName == "" {
		byt, _ := json.MarshalIndent(c, "", "  ")
		return nil, fmt.Errorf("PrettyName was empty for counter:\n%s\n", string(byt))
	}

	c.Back.PrettyName = c.PrettyName + "_back"
	c.Back.Filename = strings.TrimSuffix(c.Filename, path.Ext(c.Filename)) + "_back.png"

	c.Back.Images = mergeImagesOrTexts(c.Images, c.Back.Images)
	c.Back.Texts = mergeImagesOrTexts(c.Texts, c.Back.Texts)

	return c.Back, nil
}
