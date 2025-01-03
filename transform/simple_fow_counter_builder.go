package transform

import "github.com/sayden/counters"

// SimpleFowCounterBuilder applies a default set of modification to the counters like removing numeric values
type SimpleFowCounterBuilder struct{}

func (d *SimpleFowCounterBuilder) ToNewCounter(cc *counters.Counter) (*counters.Counter, error) {
	defer func() { cc.Metadata = counters.Metadata{} }()

	if cc.Metadata.PublicIcon.Path == "" {
		// No public image, no fow counter
		return cc, nil
	}

	cc.Texts = nil

	// Don't copy specific images in counters to the Fow. Take only center image, shield (if any) and air units faction
	// This way you can avoid a fow counter with infantry unit but the brigade/division icon with it. Or the flamethrower
	validFowImagesInCounter := make([]counters.Image, 0)
	for _, image := range cc.Images {
		if image.Position == 0 {
			image.Path = cc.Metadata.PublicIcon.Path
			image.Scale = cc.Metadata.PublicIcon.Scale
			image.YShift = floatP(0)
			image.XShift = floatP(0)
			validFowImagesInCounter = append(validFowImagesInCounter, image)
		}
	}

	cc.Images = validFowImagesInCounter

	return cc, nil
}
