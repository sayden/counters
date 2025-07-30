package output

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sayden/counters"
	"github.com/thehivecorporation/log"
)

func CardsToPNG(template *counters.CardsTemplate) error {
	var (
		sheetFileNumber      = 1
		totalCardsToGenerate int
		cardsGenerated       int
		rows                 = template.Rows
		columns              = template.Columns
		cards                = template.Cards
		currentIndex         int
	)

	// count the total amounts to card that will be processed
	for _, card := range template.Cards {
		totalCardsToGenerate += *card.Multiplier
	}

	log.Infof("Generating a total of %d cards in sheets of %d cards", totalCardsToGenerate, template.Rows*template.Columns)

	for {
		// sheet, err := createOutputCanvas(template)
		sheet, err := template.SheetCanvas()
		if err != nil {
			return errors.Wrap(err, "could not create canvas to write cards")
		}

		// cardProc := newCardProcessor(&cardProcessorConfig{
		// 	template:   template,
		// 	cardCanvas: sheet,
		// })

	processing:
		for row := range rows {
			for col := range columns {
				if currentIndex >= len(cards) {
					break processing
				}

				card := cards[currentIndex]
				cardCanvas, err := card.ToCanvas(template)
				if err != nil {
					return errors.Wrap(err, "error trying to create card canvas")
				}
				if err = cardCanvas.SavePNG(fmt.Sprintf("/tmp/%d.png", currentIndex)); err != nil {
					panic(err)
				}

				sheet.DrawImage(cardCanvas.Image(), col*card.Width, row*card.Height)
				// if err := cardProc.processCard(sheet, card, col, row); err != nil {
				// 	return errors.Wrap(err, "error trying to process card")
				// }
				*card.Multiplier--
				cardsGenerated++
				if *card.Multiplier > 0 {
					continue
				}
				currentIndex++
			}
		}

		outputPath := fmt.Sprintf(template.OutputPath, sheetFileNumber)

		if cardsGenerated >= template.Rows*template.Columns {
			// More than one file required to write cards. Expect a Go's format string in the output path
			sheetFileNumber++
		}

		// Save result on a file
		if err = sheet.SavePNG(outputPath); err != nil {
			log.WithField("output_path", outputPath).Error("error trying to save final cards image file")
			return err
		} else {
			log.WithFields(log.Fields{"output_file_path": outputPath, "generated_cards": currentIndex}).Info("Written card file")
		}

		if cardsGenerated >= totalCardsToGenerate {
			break
		}
	}

	return nil
}
