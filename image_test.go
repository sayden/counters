package counters

import (
	"image/color"
	"testing"

	"github.com/fogleman/gg"
	"github.com/stretchr/testify/assert"
)

func TestImageDraw(t *testing.T) {
	t.Run("image_draw_01", func(t *testing.T) {
		testCanvas := gg.NewContext(100, 100)

		testCanvas.Push()
		testCanvas.SetColor(color.Black)
		testCanvas.DrawRectangle(0, 0, float64(100), float64(100))
		testCanvas.Fill()
		testCanvas.Pop()

		testImage := Image{
			Settings: Settings{
				Width:        100,
				Height:       100,
				ImageScaling: IMAGE_SCALING_FIT_NONE,
				// FontPath should not be required for images, but it's expected to be set at the
				// template level
				FontPath: "assets/freesans.ttf",

				// BackgroundImage is only used for the counter and card background, but not for images
				// BackgroundImage: "",
			},
			Path: "assets/old_paper.jpg",
		}

		err := Mergev2(&testImage.Settings, &testImage.Settings)
		assert.NoError(t, err)

		// Test a background image with no scaling. Settings.BackgroundImage cannot be used
		// because it is only used for the counter background
		err = testImage.Draw(testCanvas, 0)
		assert.NoError(t, err)

		// Test image scaling
		testImage.Path = "assets/binoculars.png"
		testImage.ImageScaling = IMAGE_SCALING_FIT_WIDTH
		testImage.ShadowDistance = intP(8)
		testImage.ShadowSigma = intP(3)
		testImage.Scale = 0.5
		testImage.Margins = floatP(5)
		err = testImage.Draw(testCanvas, 11)
		assert.NoError(t, err)

		// Test header sticked to the top
		testImage.Path = "assets/stripe.png"
		testImage.ShadowDistance = intP(0)
		testImage.Scale = 1
		testImage.Margins = floatP(0)
		err = testImage.Draw(testCanvas, 1)
		assert.NoError(t, err)

		TestImageContent(t, "testdata/image_draw_01.png", 13672, testCanvas)
	})
}
