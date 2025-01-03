package scripts

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/sayden/counters"
	"github.com/sayden/counters/input"
)

var logger = log.New(os.Stderr)

func main() {

}

func jsonToAsset(inputPath, outputPath string) (err error) {
	// TODO: Check if the input is a CounterTemplate or a CardTemplate
	if err := counters.ValidateSchemaAtPath[counters.CounterTemplate](inputPath); err != nil {
		return err
	}

	counterTemplate, err := input.ReadCounterTemplate(inputPath, outputPath)
	if err != nil {
		return err
	}

	newTemplate, err := counterTemplate.ParsePrototype()
	if err != nil {
		return err
	}

	// Override output path with the one provided in the CLI
	if outputPath != "" {
		logger.Info("Overriding output path", "output_path", outputPath)
		newTemplate.OutputFolder = outputPath
	}

	return nil
}
