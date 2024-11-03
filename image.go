package counters

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/thehivecorporation/log"
	"golang.org/x/exp/slices"
)

type Image struct {
	Settings
	Positioner

	Path          string  `json:"path,omitempty"`
	Scale         float64 `json:"scale,omitempty"`
	AvoidCropping bool    `json:"avoid_cropping,omitempty"`
}

func (i Image) GetSettings() *Settings {
	return &i.Settings
}

type Images []Image

func (i *Image) Draw(dc *gg.Context, pos int) error {
	img, err := gg.LoadImage(i.Path)
	if err != nil {
		log.WithField("image", i.Path).Error("error trying to load image in 'Image' item")
		return err
	}

	if !i.AvoidCropping {
		img = CropToContent(img)
	}

	if i.Rotation != nil && *i.Rotation != 0 {
		img = imaging.Rotate(img, *i.Rotation, color.Transparent)
	}

	if i.ImageScaling == "" {
		i.ImageScaling = IMAGE_SCALING_FIT_NONE
	}

	if i.Scale == 0 {
		i.Scale = 1
	}

	switch i.ImageScaling {
	case IMAGE_SCALING_FIT_WIDTH:
		img = imaging.Resize(img, int(math.Ceil(float64(dc.Width())*i.Scale)), 0, imaging.Lanczos)
	case IMAGE_SCALING_FIT_WRAP:
		img = imaging.Resize(img, dc.Width(), dc.Height(), imaging.Lanczos)
	case IMAGE_SCALING_FIT_HEIGHT:
		img = imaging.Resize(img, int(math.Ceil(float64(dc.Height())*i.Scale)), 0, imaging.Lanczos)
	case IMAGE_SCALING_FIT_NONE:
		// Just apply scaling to image dimensions
		if i.Scale != 1 {
			img = imaging.Resize(img, int(math.Ceil(float64(img.Bounds().Dx())*i.Scale)), 0, imaging.Lanczos)
		}
	default:
		// Do nothing, image untouched from original
	}

	x, y, ax, ay, err := i.getObjectPositions(pos, &i.Settings)
	if err != nil {
		return err
	}

	if i.ShadowDistance != nil && *i.ShadowDistance != 0 {
		shadow := getShadowFromImage(img, *i.ShadowDistance, *i.ShadowSigma)
		x1 := math.Floor(x + float64(*i.ShadowDistance))
		y1 := math.Floor(y + float64(*i.ShadowDistance))
		dc.DrawImageAnchored(shadow, int(x1), int(y1), ax, ay)
	}

	dc.DrawImageAnchored(img, int(x), int(y), ax, ay)

	return nil
}

// DrawImagesOnCanvas using the provided height `h` and width `w`
func (images Images) DrawImagesOnCanvas(s *Settings, areaCanvas *gg.Context, w, h int) error {
	// sort the internal objects by position
	slices.SortFunc(images, func(i, j Image) int {
		return i.Position - j.Position
	})

	for _, image := range images {
		Mergev2(&image.Settings, s)

		image.Width = w
		image.Height = h

		if err := image.Draw(areaCanvas, image.Position); err != nil {
			return fmt.Errorf("could not draw image using path '%s': %w", image.Path, err)
		}
	}

	return nil
}

func getShadowFromImage(img image.Image, shadowDistance int, sigma int) image.Image {
	grey := imaging.AdjustBrightness(img, -100)
	w, h := getShadowImageSize(img, shadowDistance, sigma)
	temp := gg.NewContext(w*2, h*2)
	temp.DrawImageAnchored(grey, w/2, h/2, 0.5, 0.5)
	grey = imaging.AdjustBrightness(temp.Image(), 15)
	grey = imaging.Blur(grey, float64(sigma))

	return CropToContent(grey)
}

func getShadowImageSize(img image.Image, shadowDistance int, sigma int) (int, int) {
	rect := img.Bounds()
	w := rect.Dx() + sigma*sigma
	h := rect.Dy() + sigma*sigma

	return w, h
}
