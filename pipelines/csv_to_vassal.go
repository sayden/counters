package pipelines

import (
	"os"
	"path"

	"github.com/sayden/counters"
	"github.com/sayden/counters/input"
	"github.com/sayden/counters/output"
	"github.com/thehivecorporation/log"
)

type VassalConfig struct {
	Csv              string `help:"Input path of the file to read. Be aware that some outputs requires specific inputs." required:"true"`
	VassalOutputFile string `help:"Name and path of .vmod file to write. The extension .vmod is required" required:"true"`
	CounterTitle     int    `help:"The title for the counter and the file with the image comes from a column in the CSV file. Define which column here, 0 indexed" default:"3"`
}

// CSVToVassalFile takes a CSV as an input and creates a Vassal module file as output
// It uses the Vassal module stored in the TemplateModule folder as a Vassal prototype to build over it
func CSVToVassalFile(cfg VassalConfig) error {
	// Ensure that the extension of the output file is Vmod
	if path.Ext(cfg.VassalOutputFile) != "vmod" {
		log.Fatal("output file path for vassal must have '.vmod' extension")
	}

	var counterTemplate *counters.CounterTemplate
	var err error
	counterTemplate, err = input.ReadCounterTemplate(cfg.Csv, cfg.VassalOutputFile)
	if err != nil {
		return err
	}

	// Vassal mode forces individual rendering of counters
	counterTemplate.Mode = counters.TEMPLATE_MODE_TEMPLATE
	counterTemplate.OutputFolder = counters.BASE_FOLDER + "/images"
	counterTemplate.PositionNumberForFilename = 3

	output.CountersToPNG(counterTemplate)

	// Hardcoded output file, user selects the output of the vmod file
	xmlBytes, err := output.GetVassalDataForCounters(counterTemplate, counters.VassalInputXmlFile)
	if err != nil {
		log.WithError(err).Fatal("could not create xml file for vassal")
	}

	if err = os.WriteFile(counters.VassalOutputXmlFile, xmlBytes, 0666); err != nil {
		log.WithError(err).Fatal("no output xml file was generated")
	}

	return output.WriteZipFileWithFolderContent(cfg.VassalOutputFile, counters.BASE_FOLDER)
}
