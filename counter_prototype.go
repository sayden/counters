package counters

import (
	"encoding/json"
	"fmt"
	"sync"

	"dario.cat/mergo"
	"github.com/charmbracelet/log"
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

func (p *CounterPrototype) ToCounters(filenamesInUse *sync.Map, sideName string, positionNumberForFilename int) ([]Counter, error) {
	cts := make([]Counter, 0)

	// You can prototype texts and images, so one of the two must be present, get their length
	length, err := p.isLengthConsistent()
	if err != nil {
		byt, _ := json.MarshalIndent(p, "", "  ")
		return nil, fmt.Errorf("error in counter prototype\n%s\n%w", string(byt), err)
	}

	// this is a 2 step proccess for every prototype. First we apply the prototype to the front counter
	// then, using the resulting front counter, we apply the back prototype to it
	for i := 0; i < length; i++ {
		var newCounter Counter
		if err := deepcopy.FromTo(p.Counter, &newCounter); err != nil {
			return nil, err
		}

		if err = p.applyPrototypes(&newCounter, i); err != nil {
			return nil, err
		}

		var counterFilename string
		if p.Extra != nil && p.Extra.TitlePosition != nil {
			counterFilename = newCounter.GetCounterFilename(*p.Extra.TitlePosition, filenamesInUse)
		} else {
			counterFilename = newCounter.GetCounterFilename(positionNumberForFilename, filenamesInUse)
		}
		newCounter.Filename = counterFilename

		if p.Back != nil {
			var tempBackCounter Counter
			if err := deepcopy.FromTo(newCounter, &tempBackCounter); err != nil {
				return nil, err
			}

			backCounter, err := mergeFrontAndBack(&tempBackCounter, p.Back, i)
			if err != nil {
				return nil, err
			}

			backCounter.Filename = newCounter.Extra.Title + "_back.png"

			if sideName != "" {
				err = backCounter.ToVassal(sideName)
				if err != nil {
					log.Warn("could not create vassal piece")
				}
			}

			cts = append(cts, backCounter)
		}

		if sideName != "" {
			err = newCounter.ToVassal(sideName)
			if err != nil {
				log.Warn("could not create vassal piece")
			}
		}
		cts = append(cts, newCounter)
	}

	return cts, nil
}

/*
applyPrototypes applies the text and image prototypes to the given counter at the specified index.
It deep copies the text and image prototypes, updates their string and path values respectively,
and appends them to the counter's texts and images.
*/
func (p *CounterPrototype) applyPrototypes(newCounter *Counter, index int) error {
	if p.TextPrototypes != nil {
		for _, textPrototype := range p.TextPrototypes {
			originalText := Text{}
			if err := deepcopy.FromTo(textPrototype.Text, &originalText); err != nil {
				return err
			}
			originalText.String = textPrototype.StringList[index]
			newCounter.Texts = append(newCounter.Texts, originalText)
		}
	}

	if p.ImagePrototypes != nil {
		for _, imagePrototype := range p.ImagePrototypes {
			originalImage := Image{}
			if err := deepcopy.FromTo(imagePrototype.Image, &originalImage); err != nil {
				return err
			}
			originalImage.Path = imagePrototype.PathList[index]
			newCounter.Images = append(newCounter.Images, originalImage)
		}
	}

	return nil
}

/*
mergeFrontAndBack merges the images and texts from both counters. If the back prototype exists and
it has its own image or text prototypes, they are applied to the back counter, replacing or adding
to the existing images and texts.
*/
func mergeFrontAndBack(frontCounter *Counter, backPrototype *CounterPrototype, index int) (Counter, error) {
	backCounter := backPrototype.Counter
	if err := mergo.Merge(&backCounter, frontCounter); err != nil {
		return Counter{}, err
	}
	backCounter.Extra.Title = frontCounter.Extra.Title + "_back"

	backCounter.Images = mergeImagesOrTexts(frontCounter.Images, backCounter.Images)
	backCounter.Texts = mergeImagesOrTexts(frontCounter.Texts, backCounter.Texts)

	if backPrototype.ImagePrototypes != nil {
		for _, imagePrototype := range backPrototype.ImagePrototypes {
			originalImage := Image{}
			if err := deepcopy.FromTo(imagePrototype.Image, &originalImage); err != nil {
				return Counter{}, err
			}
			originalImage.Path = imagePrototype.PathList[index]

			backCounter.Images = replaceOrAddPrototypes(backCounter.Images, originalImage)
		}
	}

	if backPrototype.TextPrototypes != nil {
		for _, textPrototype := range backPrototype.TextPrototypes {
			originalText := Text{}
			if err := deepcopy.FromTo(textPrototype.Text, &originalText); err != nil {
				return Counter{}, err
			}
			originalText.String = textPrototype.StringList[index]

			backCounter.Texts = replaceOrAddPrototypes(backCounter.Texts, originalText)
		}
	}

	return backCounter, nil
}

/*
mergeImagesOrTexts is used to merge texts or images slices from the front and back of the counter. It
skips items sharing the same position unless the field BackPersistent is set to true, in which case
both items are kept with the front will be on top of the image
*/
func mergeImagesOrTexts[T SettingsGetter](frontPrototypes, backPrototypes []T) []T {
	for _, frontPrototype := range frontPrototypes {
		found := false
		for _, backPrototype := range backPrototypes {
			if frontPrototype.GetSettings().Position == backPrototype.GetSettings().Position {
				found = true
				if frontPrototype.GetSettings().BackPersistent {
					backPrototypes = append(backPrototypes, frontPrototype)
				}
				break
			}
		}

		if !found {
			backPrototypes = append(backPrototypes, frontPrototype)
		}
	}

	if len(frontPrototypes) == 0 {
		return backPrototypes
	}

	return backPrototypes
}

// replaceOrAddPrototypes is used to override texts or images in the back of the counter. It checks
// if a prototype with the same position as newPrototype already exists and replaces it if it does
// or adds it if it doesn't
func replaceOrAddPrototypes[T SettingsGetter](originals []T, newPrototype T) []T {
	found := false
	for j, original := range originals {
		if original.GetSettings().Position == newPrototype.GetSettings().Position {
			originals[j] = newPrototype
			found = true
			break
		}
	}
	if !found {
		originals = append(originals, newPrototype)
	}

	return originals
}

// isLengthConsistent checks if the lengths of text and image prototypes are consistent.
func (p *CounterPrototype) isLengthConsistent() (int, error) {
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
