package main

import (
	"image"
	"image/png"
	"math"
	"os"
	"path"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/log"
	"github.com/fogleman/gg"
)

type Tilemap struct {
	InputFolder  string `help:"Input folder to read" short:"i" required:"true"`
	OutputFolder string `help:"Output folder to write to" short:"o" required:"true"`
}

func (t *Tilemap) Run(ctx *kong.Context) error {
	log.Infof("Creating tilemap from images in %s", t.InputFolder)
	return createImageGrid(t.InputFolder, t.OutputFolder)
}

func createImageGrid(inputFolder, outputPath string) error {
	files, err := os.ReadDir(inputFolder)
	if err != nil {
		return err
	}

	images := make([]image.Image, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			imagePath := filepath.Join(inputFolder, file.Name())
			f, err := os.Open(imagePath)
			if err != nil {
				panic(err)
			}

			img, err := png.Decode(f)
			if err != nil {
				panic(err)
			}

			images = append(images, img)
			f.Close()
		}
	}

	numImages := len(images)
	gridSize := int(math.Ceil(math.Sqrt(float64(numImages))))

	// Determine the maximum dimension of all images
	maxDim := 0
	for _, img := range images {
		bounds := img.Bounds()
		maxDim = max(maxDim, max(bounds.Dx(), bounds.Dy()))
	}

	tileSize := maxDim
	canvasSize := gridSize * tileSize

	dc := gg.NewContext(canvasSize, canvasSize)

	for i, img := range images {
		x := (i % gridSize) * tileSize
		y := (i / gridSize) * tileSize

		// Center the image within its tile
		bounds := img.Bounds()
		offsetX := (tileSize - bounds.Dx()) / 2
		offsetY := (tileSize - bounds.Dy()) / 2
		dc.DrawImage(img, x+offsetX, y+offsetY)
	}

	os.MkdirAll(outputPath, 0755)
	return dc.SavePNG(path.Join(outputPath, "tilemap.png"))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
