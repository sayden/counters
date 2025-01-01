package transform

import (
	"testing"

	"github.com/sayden/counters"
	"github.com/stretchr/testify/assert"
)

func TestToNewCounter_SingleStep(t *testing.T) {
	builder := &StepLossBackCounterBuilder{}
	counter := &counters.Counter{SingleStep: true}

	newCounter, err := builder.ToNewCounter(counter)
	assert.NoError(t, err)
	assert.Nil(t, newCounter.Texts)
	assert.Nil(t, newCounter.Images)
}

func TestToNewCounter_ExtraNil(t *testing.T) {
	builder := &StepLossBackCounterBuilder{}
	counter := &counters.Counter{SingleStep: false, Metadata: nil}

	newCounter, err := builder.ToNewCounter(counter)
	assert.NoError(t, err)
	assert.NotNil(t, newCounter.Metadata)
	assert.Equal(t, "back", newCounter.Metadata.Side)
	assert.Equal(t, counters.STRIPE, newCounter.Images[0].Path)
	assert.Equal(t, "fitWidth", newCounter.Images[0].Settings.ImageScaling)
}

func TestToNewCounter_ExtraNotNil(t *testing.T) {
	builder := &StepLossBackCounterBuilder{}
	counter := &counters.Counter{SingleStep: false, Metadata: &counters.Metadata{Side: "front"}}

	newCounter, err := builder.ToNewCounter(counter)
	assert.NoError(t, err)
	assert.NotNil(t, newCounter.Metadata)
	assert.Equal(t, "frontback", newCounter.Metadata.Side)
	assert.Equal(t, counters.STRIPE, newCounter.Images[0].Path)
	assert.Equal(t, "fitWidth", newCounter.Images[0].Settings.ImageScaling)
}
