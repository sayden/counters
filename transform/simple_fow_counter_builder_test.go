package transform

import (
	"testing"

	"github.com/sayden/counters"
	"github.com/stretchr/testify/assert"
)

func TestToNewCounter(t *testing.T) {
	builder := &SimpleFowCounterBuilder{}

	t.Run("No Public Icon Path", func(t *testing.T) {
		counter := &counters.Counter{}
		newCounter, err := builder.ToNewCounter(counter)
		assert.NoError(t, err)
		assert.Equal(t, counter, newCounter)
	})

	t.Run("With Public Icon Path", func(t *testing.T) {
		counter := &counters.Counter{
			Extra: &counters.Extra{
				PublicIcon: &counters.Image{Path: "assets/binoculars.png", Scale: 1.5}},
			Images: []counters.Image{
				{Path: "assets/stripe.png", Scale: 0.5, Settings: counters.Settings{YShift: 10, XShift: 10}},
				{Path: "assets/binoculars.png", Settings: counters.Settings{Position: 1}}},
		}
		expectedCounter := &counters.Counter{
			Images: []counters.Image{
				{Path: "assets/binoculars.png", Scale: 1.5}},
		}
		newCounter, err := builder.ToNewCounter(counter)
		assert.NoError(t, err)
		assert.Equal(t, expectedCounter, newCounter)
	})
}