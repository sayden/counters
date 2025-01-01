package transform

import "github.com/sayden/counters"

// StepLossBackCounterBuilder is the most basic type of back counter builder, by adding a simple red stripe in the
// background of the counter
type StepLossBackCounterBuilder struct {
	finalOutputCounters []*counters.Counter
}

func (bc *StepLossBackCounterBuilder) ToNewCounter(counter *counters.Counter) (*counters.Counter, error) {
	if counter.SingleStep {
		counter.Texts = nil
		counter.Images = nil
		return counter, nil
	}

	if counter.Metadata == nil {
		counter.Metadata = &counters.Metadata{Side: "back"}
	} else {
		counter.Metadata.Side += "back"
	}

	// ensure the 'stripe' is the first image in the counter and everything else is on top of it
	counter.Images = append(
		[]counters.Image{
			{
				Path: counters.STRIPE, Settings: counters.Settings{ImageScaling: "fitWidth"},
			},
		}, counter.Images...,
	)

	return counter, nil
}
