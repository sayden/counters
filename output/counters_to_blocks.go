package output

import (
	"image/color"
	"os"

	"github.com/charmbracelet/log"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
)

func CountersToBlocks(frontTemplate, backTemplate *counters.CounterTemplate) error {
	var canvas *gg.Context
	if frontTemplate.Mode == counters.TEMPLATE_MODE_TILES {
		return errors.New("tiles cannot be set when producing blocks")
	}

	_ = os.MkdirAll(backTemplate.OutputFolder, 0750)

	counterPos := 0
	row := 0
	// fileNumber := 1
	gs := newGlobalState(frontTemplate)

	//Iterate rows and columns, painting a counter on each
iteration:
	for {
		for columns := 0; columns < frontTemplate.Columns; columns++ {
			if counterPos == len(frontTemplate.Counters) {
				break iteration
			}

			counter := frontTemplate.Counters[counterPos]

			counterCanvas, err := counter.Canvas(frontTemplate.DrawGuides)
			if err != nil {
				return err
			}

			blockCanvas, err := getBlockCanvasFromCounterCanvas(counterCanvas, &counter)
			if err != nil {
				return err
			}

			if backTemplate != nil {
				backCounter := backTemplate.Counters[counterPos]
				backCounterCanvas, err := backCounter.Canvas(backTemplate.DrawGuides)
				if err != nil {
					return err
				}

				addBackCounterToBlockCanvas(backCounterCanvas, blockCanvas)
			}

			iw := imageWriter{canvas: blockCanvas, template: gs.template}
			if err = iw.createFile(&counter, gs); err != nil {
				log.Error("error trying to write counter to file", err)
			}

			counterPos++
		}

		row++
	}

	// Save result on a file
	if frontTemplate.Mode == counters.TEMPLATE_MODE_TILES {
		return canvas.SavePNG(frontTemplate.OutputFolder)
	}

	return nil
}

func addBackCounterToBlockCanvas(backCounter, blockCanvas *gg.Context) {
	canvasWidth := blockCanvas.Width()
	canvasHeight := blockCanvas.Height()
	margin := int(float64(canvasWidth) * 0.0508)

	img := backCounter.Image()
	img = imaging.Rotate(img, 180, color.Black)

	blockCanvas.DrawImage(img, int(canvasWidth-margin-backCounter.Width()), int(canvasHeight-margin-backCounter.Height()))
}

// getBlockCanvasFromCounterCanvas creates a new canvas for a block by scaling up the given counter canvas.
// It sets the background color, loads the specified font, and draws the counter image onto the new canvas.
func getBlockCanvasFromCounterCanvas(counterCanvas *gg.Context, cc *counters.Counter) (*gg.Context, error) {
	const canvasRatio = 2.49
	canvasWidth := float64(counterCanvas.Width()) * canvasRatio
	canvasHeigth := float64(counterCanvas.Height()) * canvasRatio
	dc := gg.NewContext(int(canvasWidth), int(canvasHeigth))
	if err := cc.LoadFontOrDefault(dc); err != nil {
		return nil, err
	}

	dc.Push()
	dc.SetColor(cc.BgColor)
	dc.DrawRectangle(0, 0, canvasWidth, canvasHeigth)
	dc.Fill()
	dc.Pop()

	if cc.FontColorS != "" {
		counters.ColorFromStringOrDefault(cc.FontColorS, cc.BgColor)
	}

	// Draw the counter into the new canvas
	margin := canvasWidth * 0.0508
	bottomMargin := int(canvasHeigth - margin)
	dc.DrawImage(counterCanvas.Image(), int(margin), bottomMargin-cc.Height)

	return dc, nil
}
