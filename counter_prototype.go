package counters

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	deepcopy "github.com/qdm12/reprint"
)

type CounterPrototype struct {
	Counter
	ImagePrototypes []ImagePrototype  `json:"image_prototypes,omitempty"`
	TextPrototypes  []TextPrototype   `json:"text_prototypes,omitempty"`
	Back            *CounterPrototype `json:"back,omitempty"`
}

type ImagePrototype struct {
	Image
	PathList []string `json:"path_list"`
}

type TextPrototype struct {
	Text
	StringList []string `json:"string_list"`
}

func (p *CounterPrototype) ToCounters() ([]Counter, error) {
	cts := make([]Counter, 0)

	// You can prototype texts and images, so one of the two must be present, get their length
	length, err := p.isExpectedLengthsCorrect()
	if err != nil {
		byt, _ := json.MarshalIndent(p, "", "  ")
		return nil, fmt.Errorf("error in counter prototype\n%s\n%w", string(byt), err)
	}

	for i := 0; i < length; i++ {
		var newCounter Counter
		if err := deepcopy.FromTo(p.Counter, &newCounter); err != nil {
			return nil, err
		}

		if p.TextPrototypes != nil {
			for _, textPrototype := range p.TextPrototypes {
				originalText := Text{}
				if err := deepcopy.FromTo(textPrototype.Text, &originalText); err != nil {
					return nil, err
				}
				originalText.String = textPrototype.StringList[i]
				newCounter.Texts = append(newCounter.Texts, originalText)
			}
		}

		if p.ImagePrototypes != nil {
			for _, imagePrototype := range p.ImagePrototypes {
				originalImage := Image{}
				if err := deepcopy.FromTo(imagePrototype.Image, &originalImage); err != nil {
					return nil, err
				}
				originalImage.Path = imagePrototype.PathList[i]
				newCounter.Images = append(newCounter.Images, originalImage)
			}
		}
		cts = append(cts, newCounter)

		if p.Back != nil {
			backCounter := p.Back.Counter
			images := p.Back.Counter.Images
			texts := p.Back.Counter.Texts
			if err := deepcopy.FromTo(newCounter, &backCounter); err != nil {
				return nil, err
			}
			if backCounter.Extra != nil {
				backCounter.Extra.Title += " back"
			}

			// If 2 images shares position, the one in the back counter will be used
			frontImages := newCounter.Images
			backCounter.Images = images
			for _, frontImage := range frontImages {
				// check if the image's position already exists in images, skip it in such case
				found := false
				for _, backImage := range backCounter.Images {
					if frontImage.Position == backImage.Position {
						found = true
						if frontImage.BackPersistent {
							backCounter.Images = append(backCounter.Images, frontImage)
						}
						break
					}
				}

				if !found {
					backCounter.Images = append(backCounter.Images, frontImage)
				}
			}

			// Do the same with texts
			frontTexts := newCounter.Texts
			backCounter.Texts = texts
			for _, frontText := range frontTexts {
				// check if the image's position already exists in images, skip it in such case
				found := false
				for _, backText := range backCounter.Texts {
					if frontText.Position == backText.Position {
						found = true
						if frontText.BackPersistent {
							backCounter.Texts = append(backCounter.Texts, frontText)
						}
						break
					}
				}

				if !found {
					backCounter.Texts = append(backCounter.Texts, frontText)
				}
			}

			if p.Back.ImagePrototypes != nil {
				for _, imagePrototype := range p.Back.ImagePrototypes {
					originalImage := Image{}
					if err := deepcopy.FromTo(imagePrototype.Image, &originalImage); err != nil {
						return nil, err
					}
					originalImage.Path = imagePrototype.PathList[i]

					// Replace the image in the back counter on the same position
					found := false
					for j, image := range backCounter.Images {
						if image.Position == originalImage.Position {
							backCounter.Images[j] = originalImage
							found = true
							break
						}
					}
					if !found {
						backCounter.Images = append(backCounter.Images, originalImage)
					}
				}
			}

			if p.Back.TextPrototypes != nil {
				for _, textPrototype := range p.Back.TextPrototypes {
					originalText := Text{}
					if err := deepcopy.FromTo(textPrototype.Text, &originalText); err != nil {
						return nil, err
					}
					originalText.String = textPrototype.StringList[i]

					// Replace the text in the back counter on the same position
					found := false
					for j, text := range backCounter.Texts {
						if text.Position == originalText.Position {
							backCounter.Texts[j] = originalText
							found = true
							break
						}
					}
					if !found {
						backCounter.Texts = append(backCounter.Texts, originalText)
					}
				}
			}

			cts = append(cts, backCounter)
		}

	}

	return cts, nil
}

func (p *CounterPrototype) isExpectedLengthsCorrect() (int, error) {
	// find a reference length
	length := p.getTextLength(p.TextPrototypes)

	// if no text prototypes found to use for reference, try with image prototypes
	if length == 0 {
		length = p.getImageLength(p.ImagePrototypes)
	}

	if length == 0 {
		return 0, errors.New("no prototypes found in the counter template")
	}

	lengths := map[string]int{
		"Text prototypes":  p.getTextLength(p.TextPrototypes),
		"Image prototypes": p.getImageLength(p.ImagePrototypes),
	}
	if p.Back != nil {
		lengths["Back text prototypes"] = p.getTextLength(p.Back.TextPrototypes)
		lengths["Back image prototypes"] = p.getImageLength(p.Back.ImagePrototypes)
	}

	for s, l := range lengths {
		if l > 0 && l != length {
			return 0, fmt.Errorf("the number of images and texts prototypes must be the same than "+
				"the reference '%d' in '%s', found %d != %d", length, s, l, length)
		}
	}

	return length, nil
}

func (p *CounterPrototype) getImageLength(ts []ImagePrototype) int {
	if len(ts) > 0 {
		return len(ts[0].PathList)
	}

	return 0
}

func (p *CounterPrototype) getTextLength(ts []TextPrototype) int {
	if len(ts) > 0 {
		return len(ts[0].StringList)
	}

	return 0
}
