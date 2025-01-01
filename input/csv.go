package input

import (
	"encoding/csv"
	"image/color"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/creasty/defaults"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
)

// ReadCSVCounters follows a convention to read CSV files into counter templates:
// A-Q Columns in the CSV corresponds to the 0-16 possible position to write into a counter, being 0 the
// horizontal and vertical center
// Column R (17): Is the Side: red, blue, german, etc.
// Column S (18): Is the background color
func ReadCSVCounters(filepath string) (*counters.CounterTemplate, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "error trying to open counter template")
	}
	defer f.Close()

	csvReader := csv.NewReader(f)

	// Skip headers
	_, err = csvReader.Read()
	if err != nil {
		return nil, errors.Wrap(err, "error trying to read the first row on the csv file")
	}

	cnts := make([]counters.Counter, 0)

	for {
		cols, err := csvReader.Read()
		if err != nil {
			break
		}

		cnt := counters.Counter{
			Texts: make([]counters.Text, 0),
		}

		for colIndex, col := range cols {
			if col != "" {
				cell := strings.TrimSpace(col)

				// Side: red, blue, german, russian, etc.
				if colIndex == 17 {
					cnt.Metadata.Side = cell
					continue
				} else if colIndex == 18 {
					//background color
					cnt.BackgroundColor = stringP("#" + cell)
					break
				}

				cnt.Texts = append(
					cnt.Texts, counters.Text{
						Settings: counters.Settings{
							Position:      colIndex,
							AvoidClipping: true,
							FontColorS:    "white",
							StrokeWidth:   floatP(2),
							StrokeColorS:  "black",
						},
						String: cell,
					},
				)
			}
		}

		cnts = append(cnts, cnt)
	}

	//TODO remove hardcoded data
	template := counters.CounterTemplate{
		Rows:                      7,
		Columns:                   7,
		Mode:                      "template",
		OutputFolder:              "/tmp",
		Counters:                  cnts,
		PositionNumberForFilename: 3,
	}

	template.Settings.Width = 50
	template.Settings.Height = 50
	template.Settings.FontHeight = 16
	template.Settings.Margins = floatP(2)
	template.Settings.FontColor = color.White
	template.Settings.ImageScaling = counters.IMAGE_SCALING_FIT_NONE

	if err = template.EnrichTemplate(); err != nil {
		return nil, errors.Wrap(err, "could not enrich template")
	}

	return &template, nil
}

func ReadCSVCards(f io.Reader, template *counters.CardsTemplate) (*counters.CardsTemplate, error) {
	csvReader := csv.NewReader(f)

	headTitles, err := csvReader.Read()
	if err != nil {
		return nil, errors.Wrap(err, "error trying to read the first row on the csv file")
	}

	cards := make([]counters.Card, 0)

	for {
		cols, err := csvReader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		card := counters.Card{
			Areas: make([]counters.Counter, 0),
		}

		// skip leftmost column (multiplier) by now
		for i, col := range cols[1:] {
			if i == 0 && headTitles[1] == "bg_color" {
				card.Settings.BackgroundColor = stringP(col)
				continue
			}

			cell := strings.TrimSpace(col)

			card.Areas = append(card.Areas, counters.Counter{
				Texts: []counters.Text{{
					Settings: counters.Settings{
						Position:      0,
						AvoidClipping: true,
						BgColor:       color.White,
					},
					String: cell,
				}},
			})
		}

		// read the leftmost column now
		multiplier, err := strconv.Atoi(cols[0])
		if err != nil {
			return nil, err
		}

		for i := 0; i < multiplier; i++ {
			cards = append(cards, card)
		}
	}

	template.Cards = cards

	if err := defaults.Set(&template.Settings); err != nil {
		return nil, errors.Wrap(err, "could not set defaults into card template")
	}

	// counters.ApplyCardWaterfallSettings(template)

	// if err = counters.EnrichTemplate(template); err != nil {
	// 	return nil, errors.Wrap(err, "could not enrich template")
	// }

	return template, nil
}
