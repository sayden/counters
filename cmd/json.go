package main

import (
	"os"
	"sync"

	"github.com/alecthomas/kong"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
	"github.com/sayden/counters/fsops"
	"github.com/sayden/counters/input"
	"github.com/sayden/counters/output"
	"github.com/sayden/counters/transform"
)

type JsonOutput struct {
	InputPath   string `help:"Input path of the file to read" short:"i" required:"true"`
	OutputPath  string `help:"Path to the folder to write the JSON" short:"o"`
	OutputType  string `help:"Type of content to produce: back-counters, cards, fow-counters, counters or events" short:"t" default:"back-counters"`
	Destination string `help:"When generating a JSON Template, this contains the destination folder for images inside the template" short:"d"`

	EventsPoolFile          string `help:"A file to take 'events' from"`
	BackImage               string `help:"The image for the back of the cards"`
	CardTemplateFilepath    string `help:"When writing cards from a CSV, a template for those cards must be provided"`
	CounterTemplateFilepath string `help:"When writing counters from a CSV, a template for those counters must be provided"`
	BackgroundImages        string `help:"Path to a folder containing background images for the cards"`
	QuotesFile              string `help:"Path to a JSON file containing quotes for the cards"`
}

func (i *JsonOutput) Run(ctx *kong.Context) error {
	var counterTemplate *counters.CounterTemplate
	var events []counters.Event

	// content could be JSON or CSV
	inputContent, err := fsops.GetExtension(Cli.Json.InputPath)
	if err != nil {
		return err
	}

	byt, err := os.ReadFile(Cli.Json.InputPath)
	if err != nil {
		return errors.Wrap(err, "could not read input file")
	}

	// TODO: Check if the input is a CounterTemplate or a CardTemplate
	if err := counters.ValidateSchemaAtPath[counters.CounterTemplate](Cli.Json.InputPath); err != nil {
		return err
	}

	filenamesInUse := &sync.Map{}

	switch inputContent {
	case counters.FileContent_CSV:
		switch Cli.Json.OutputType {
		case "cards":
			return csvToCards(byt)
		case "counters":
			return csvToCounters(byt)
		}
	case counters.FileContent_JSON:
		inputContent, err := fsops.IdentifyJSONFileContent(byt)
		if err != nil {
			return errors.Wrap(err, "could not identify file content")
		}

		// The input is a JSON file, either a counter template or a list of events (in a special JSON format too)
		switch inputContent {
		case counters.FileContent_CounterTemplate:
			counterTemplate, err = input.ReadCounterTemplate(Cli.Json.InputPath, Cli.Json.OutputPath)
			if err != nil {
				return errors.Wrap(err, "error reading counter template")
			}
		case counters.FileContent_Events:
			events, err = input.JsonEventsToEvents(Cli.Json.InputPath)
			if err != nil {
				return errors.Wrap(err, "could not read events file")
			}
		default:
			return errors.New("combination of options not recognized")
		}

		switch Cli.Json.OutputType {
		case "counters":
			// JSON counters to Counters
			newTempl, err := counterTemplate.ExpandPrototypeCounterTemplate(filenamesInUse)
			if err != nil {
				return errors.Wrap(err, "error trying to convert a counter template into another counter template")
			}
			return output.ToJSONFile(newTempl, Cli.Json.OutputPath)

		case "back-counters":
			// JSON counters to Counters
			newTempl, err := counterTemplate.ExpandPrototypeCounterTemplate(filenamesInUse)
			if err != nil {
				return errors.Wrap(err, "error trying to convert a counter template into another counter template")
			}

			// JSON counters to Back Counters
			ctc := &transform.CountersToCountersConfig{
				OriginalCounterTemplate: newTempl,
				OutputPathInTemplate:    Cli.Json.Destination,
				CounterTransformer:      &transform.StepLossBackCounterBuilder{}}
			finalCounters, err := ctc.CountersToCounters()
			if err != nil {
				return errors.Wrap(err, "error trying to convert a counter template into another counter template")
			}

			return output.ToJSONFile(finalCounters, Cli.Json.OutputPath)
		case "cards":
			// JSON counters to Cards
			return jsonCountersToJsonCards(counterTemplate)
		case "fow-counters":
			// JSON counters to Fow Counters
			return jsonCountersToJsonFowCounters(counterTemplate)
		case "events":
			// FIXME JSON to Events
			return jsonToJsonCardEvents(events)
		}
	}

	return errors.New("'--output-content' not recognized for JSON output")
}
