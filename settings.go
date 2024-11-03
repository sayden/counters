package counters

import (
	"fmt"
	"image/color"

	"dario.cat/mergo"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"github.com/thehivecorporation/log"
)

type SettingsGetter interface {
	GetSettings() *Settings
}

// Template Settings
//
//	Counter / Card Settings
//		Image / Text Settings
type Settings struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`

	Margins *float64 `json:"margins,omitempty"`

	FontHeight float64 `json:"font_height,omitempty"`
	FontPath   string  `json:"font_path,omitempty"`

	FontColorS string      `json:"font_color,omitempty"`
	FontColor  color.Color `json:"-"`

	BackgroundImage *string     `json:"background_image,omitempty"`
	BackgroundColor *string     `json:"background_color,omitempty"`
	BgColor         color.Color `json:"-"`

	ShadowDistance *int `json:"shadow,omitempty"`
	ShadowSigma    *int `json:"shadow_sigma,omitempty"`

	Rotation *float64 `json:"rotation,omitempty"`

	//CounterTemplate Card specific
	BorderWidth  *float64    `json:"border_width,omitempty"`
	BorderColorS string      `json:"border_color,omitempty"`
	BorderColor  color.Color `json:"-"`

	XShift *float64 `json:"x_shift,omitempty"`
	YShift *float64 `json:"y_shift,omitempty"`

	//Card specific
	Multiplier *int `json:"multiplier,omitempty" default:"1"`

	//Text options specific
	StrokeWidth  *float64    `json:"stroke_width,omitempty"`
	StrokeColorS string      `json:"stroke_color,omitempty"`
	StrokeColor  color.Color `json:"-"`
	Alignment    string      `json:"alignment,omitempty"`

	//Image options specific
	ImageScaling string `json:"image_scaling,omitempty"`

	AvoidClipping bool `json:"avoid_clipping,omitempty"`

	Position int `json:"position,omitempty"`

	Skip bool `json:"skip,omitempty"`

	BackPersistent bool `json:"back_persistent,omitempty"`
}

func (s *Settings) ApplySettingsScaling(scaling float64) {
	s.Width = int(scaling * float64(s.Width))
	s.Height = int(scaling * float64(s.Height))

	if s.Margins != nil {
		*s.Margins *= scaling
		if *s.Margins < 1 {
			*s.Margins = 1
		}
	}

	s.FontHeight *= scaling

	if s.ShadowDistance != nil {
		*s.ShadowDistance = int(scaling * float64(*s.ShadowDistance))
	}

	if s.BorderWidth != nil {
		*s.BorderWidth *= scaling
		if *s.BorderWidth < 1 {
			*s.BorderWidth = 0
		}
	}

	if s.XShift != nil {
		*s.XShift *= scaling
	}
	if s.YShift != nil {
		*s.YShift *= scaling
	}

	if s.StrokeWidth != nil {
		*s.StrokeWidth *= scaling
	}
}

// DrawBackgroundImage draws the background image, if any, on the provided context
func (s *Settings) DrawBackgroundImage(dc *gg.Context) error {
	if s.BackgroundImage == nil {
		return nil
	}
	if s.BackgroundImage != nil && *s.BackgroundImage == "" {
		return nil
	}

	img, err := imaging.Open(*s.BackgroundImage)
	if err != nil {
		log.WithField("image", s.BackgroundImage).Error(err)
		return err
	}

	img = imaging.Resize(img, 0, s.Height, imaging.Gaussian)
	dc.DrawImageAnchored(img, dc.Width()/2, dc.Height()/2, 0.5, 0.5)

	return nil
}

func Merge(d *Settings, s2 Settings, opt ...func(*mergo.Config)) error {
	d.FontColor = nil
	d.StrokeColor = nil
	d.BorderColor = nil
	d.BgColor = nil

	err := mergo.Merge(d, s2, opt...)
	if err != nil {
		return fmt.Errorf("error merging settings: %v", err)
	}

	SetColors(d)

	return nil
}

func Mergev2(d *Settings, s *Settings) error {
	if d.Width == 0 && s.Width == 0 {
		return errors.New("no valid width was found")
	} else if d.Width == 0 {
		d.Width = s.Width
	}

	if d.Height == 0 && s.Height == 0 {
		return errors.New("no valid height was found")
	} else if d.Height == 0 {
		d.Height = s.Height
	}

	if d.Margins == nil {
		if s.Margins != nil {
			d.Margins = s.Margins
		} else {
			zeroFloat := 0.0
			d.Margins = &zeroFloat
		}
	}

	if d.FontHeight == 0 {
		if s.FontHeight != 0 {
			d.FontHeight = s.FontHeight
		}
		// can be zero because it can be an image
	}

	if d.FontPath == "" && s.FontPath == "" {
		// can be empty because it can be an image
	} else {
		d.FontPath = s.FontPath
	}

	if d.FontColor == nil {
		if d.FontColorS == "" {
			if s.FontColor == nil {
				if s.FontColorS == "" {
					//default
					d.FontColor = ColorFromStringOrDefault("black", color.Black)
					d.FontColorS = "black"
				} else {
					d.FontColorS = s.FontColorS
					s.FontColor = ColorFromStringOrDefault(s.FontColorS, color.Black)
					d.FontColor = ColorFromStringOrDefault(s.FontColorS, color.Black)
				}
			} else {
				d.FontColor = s.FontColor
				d.FontColorS = s.FontColorS
			}
		} else {
			d.FontColor = ColorFromStringOrDefault(d.FontColorS, color.Black)
		}
	}

	if d.BgColor == nil {
		if d.BackgroundColor == nil {
			if s.BgColor == nil {
				if s.BackgroundColor == nil {
				} else if *s.BackgroundColor != "" {
					d.BackgroundColor = s.BackgroundColor
					d.BgColor = ColorFromStringOrDefault(*s.BackgroundColor, color.White)
					s.BgColor = ColorFromStringOrDefault(*s.BackgroundColor, color.White)
				}
			} else {
				d.BgColor = s.BgColor
				d.BackgroundColor = s.BackgroundColor
			}
		} else if *d.BackgroundColor != "" {
			d.BgColor = ColorFromStringOrDefault(*d.BackgroundColor, color.White)
		}
	}

	if d.BgColor == nil {
		if d.BackgroundColor == nil {
			if s.BgColor == nil {
				if s.BackgroundColor == nil {
					// default
					d.BgColor = color.White
				} else if *s.BackgroundColor == "" {
					// default
					d.BgColor = color.White
				} else {
					d.BackgroundColor = s.BackgroundColor
					d.BgColor = ColorFromStringOrDefault(*s.BackgroundColor, color.White)
					s.BgColor = ColorFromStringOrDefault(*s.BackgroundColor, color.White)
				}
			} else {
				d.BgColor = s.BgColor
			}
		} else {
			d.BgColor = ColorFromStringOrDefault(*d.BackgroundColor, color.White)
		}
	}

	if d.ShadowDistance == nil {
		if s.ShadowDistance != nil {
			d.ShadowDistance = s.ShadowDistance
		} else {
			zeroInt := 0
			d.ShadowDistance = &zeroInt
		}
	}

	if d.ShadowSigma == nil {
		if s.ShadowSigma != nil {
			d.ShadowSigma = s.ShadowSigma
		} else {
			zeroInt := 0
			d.ShadowSigma = &zeroInt
		}
	}

	if d.Rotation == nil {
		if s.Rotation != nil {
			d.Rotation = s.Rotation
		} else {
			zeroFloat := 0.0
			d.Rotation = &zeroFloat
		}
	}

	if d.BorderWidth == nil {
		if s.BorderWidth != nil {
			d.BorderWidth = s.BorderWidth
		} else {
			zeroFloat := 0.0
			d.BorderWidth = &zeroFloat
		}
	}

	if d.BorderColor == nil {
		if s.BorderColor == nil {
			if s.BorderColorS != "" {
				d.BorderColorS = s.BorderColorS
				d.BorderColor = ColorFromStringOrDefault(s.BorderColorS, color.Black)
				s.BorderColor = ColorFromStringOrDefault(s.BorderColorS, color.Black)
			}
		} else {
			d.BorderColor = s.BorderColor
		}
	}

	if d.XShift == nil {
		if s.XShift != nil {
			d.XShift = s.XShift
		} else {
			zeroFloat := 0.0
			d.XShift = &zeroFloat
		}
	}

	if d.YShift == nil {
		if s.YShift != nil {
			d.YShift = s.YShift
		} else {
			zeroFloat := 0.0
			d.YShift = &zeroFloat
		}
	}

	if d.Multiplier == nil {
		if s.Multiplier != nil {
			d.Multiplier = s.Multiplier
		} else {
			zeroInt := 1
			d.Multiplier = &zeroInt
		}
	}

	if d.StrokeWidth == nil {
		if s.StrokeWidth != nil {
			d.StrokeWidth = s.StrokeWidth
		} else {
			zeroFloat := 0.0
			d.StrokeWidth = &zeroFloat
		}
	}

	if d.StrokeColor == nil {
		if s.StrokeColor == nil {
			if s.StrokeColorS != "" {
				d.StrokeColorS = s.StrokeColorS
				d.StrokeColor = ColorFromStringOrDefault(s.StrokeColorS, color.Black)
				s.StrokeColor = ColorFromStringOrDefault(s.StrokeColorS, color.Black)
			} else {
				// default
				d.StrokeColor = color.Black
			}
		} else {
			d.StrokeColor = s.StrokeColor
		}
	}

	if d.Alignment == "" && s.Alignment != "" {
		d.Alignment = s.Alignment
	}

	if d.ImageScaling == "" && s.ImageScaling != "" {
		d.ImageScaling = s.ImageScaling
	}

	if !d.AvoidClipping && s.AvoidClipping {
		d.AvoidClipping = s.AvoidClipping
	}

	return nil
}
