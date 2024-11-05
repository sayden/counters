package output

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/sayden/counters"
	"github.com/sayden/counters/fsops"
	"github.com/sayden/counters/input"
	"github.com/thehivecorporation/log"
)

type VassalConfig struct {
	Csv              string `help:"Input path of the file to read. Be aware that some outputs requires specific inputs."`
	VassalOutputFile string `help:"Name and path of .vmod file to write. The extension .vmod is required"`
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

	CountersToPNG(counterTemplate)

	// Hardcoded output file, user selects the output of the vmod file
	xmlBytes, err := getVassalDataForCounters(counterTemplate, counters.VassalInputXmlFile)
	if err != nil {
		log.WithError(err).Fatal("could not create xml file for vassal")
	}

	if err = os.WriteFile(counters.VassalOutputXmlFile, xmlBytes, 0666); err != nil {
		log.WithError(err).Fatal("no output xml file was generated")
	}

	return WriteZipFileWithFolderContent(cfg.VassalOutputFile, counters.BASE_FOLDER)
}

// getVassalDataForCounters returns the Vassal module data for the counters
func getVassalDataForCounters(t *counters.CounterTemplate, xmlFilepath string) ([]byte, error) {
	var g counters.VassalGameModule
	err := fsops.ReadMarkupFile(xmlFilepath, &g)
	if err != nil {
		return nil, errors.Wrap(err, "error trying to decode content")
	}

	// Piece palette definition
	tw := counters.TabWidget{
		EntryName: "Forces",
		// PanelWidget: counters.PanelWidget{
		// 	ListWidget: make([]counters.ListWidget, 0),
		// 	EntryName:  "Forces",
		// 	NColumns:   "3",
		// 	Scale:      "1.0",
		// 	Text:       "Forces",
		// 	Vert:       "false",
		// 	Fixed:      "false",
		// },
	}

	forces := make(map[string]counters.ListWidget)
	forces["Markers"] = counters.ListWidget{
		EntryName: "Markers",
		PieceSlot: make([]counters.PieceSlot, 0),
		Scale:     "1.0",
		Height:    "215",
		Width:     "562",
		Divider:   "194",
	}

	// originalPieceTemplate := `+/null/prototype;Basic Pieces	emb2;Flip1;128;A;;128;;;128;;;;1;false;0;0;1 TD X HQ.png;Back;true;Flip Layer (Name);;;false;;1;1;false;;;;Description;1.0;;true\	piece;;;1 TD Unit.png;1 TD Unit/	-1\	null;0;0;;1;ppScale;1.0`

	xmlTemplateString := `+/null/prototype;BasicPrototype	piece;;;{{ .Filename }};{{ .PieceName}}/	null;0;0;{{ .Id }};0`
	xmlTemplate, err := template.New("xml").Parse(xmlTemplateString)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse template string")
	}

	// Read special files from images folder to statically load them into the module as pieces
	// (terrain, -1, -2 markers, Disorganized, Spent, OOS, etc,)
	files, err := readFiles(counters.BASE_FOLDER + "/images")
	if err != nil {
		return nil, errors.Wrap(err, "could not read files")
	}

	gpid := 200
	id := 200

	// Load markers into the module
	for _, file := range files {
		buf := bytes.NewBufferString("")
		err = xmlTemplate.ExecuteTemplate(buf, "xml", counters.CSVTemplateData{
			Filename:  file,
			PieceName: file,
			Id:        fmt.Sprintf("%d", id),
		})
		if err != nil {
			return nil, errors.Wrap(err, "error trying to write Vassal xml file using templates")
		}
		id++

		extension := path.Ext(file)
		fileWithoutExtension := strings.TrimSuffix(file, extension)
		sss := strings.Split(fileWithoutExtension, "_")
		name := sss[0]
		if len(sss) > 1 {
			name = sss[1]
		}
		piece := counters.PieceSlot{
			EntryName: name,
			Gpid:      fmt.Sprintf("%d", gpid),
			Height:    t.Height,
			Width:     t.Width,
			Data:      buf.String(),
		}

		temp := forces["Markers"]
		temp.PieceSlot = append(forces["Markers"].PieceSlot, piece)
		forces["Markers"] = temp

		gpid++
	}

	// filenamesInUse := new(sync.Map)

	// Load counters into the module
	for _, counter := range t.Counters {
		buf := bytes.NewBufferString("")
		if err = xmlTemplate.ExecuteTemplate(buf, "xml",
			counters.CSVTemplateData{
				// Filename:  counter.GetCounterFilename("", t.PositionNumberForFilename, filenamesInUse),
				// PieceName: counter.GetCounterFilename("", t.PositionNumberForFilename, filenamesInUse),
				Id: fmt.Sprintf("%d", id),
			},
		); err != nil {
			return nil, errors.Wrap(err, "error trying to write Vassal xml file using templates")
		}
		id++

		piece := counters.PieceSlot{
			EntryName: counter.GetTextInPosition(t.PositionNumberForFilename),
			Gpid:      fmt.Sprintf("%d", gpid),
			Height:    t.Height,
			Width:     t.Width,
			Data:      buf.String(),
		}

		if _, ok := forces[counter.Extra.Side]; !ok {
			forces[counter.Extra.Side] = counters.ListWidget{
				EntryName: counter.Extra.Side,
				PieceSlot: make([]counters.PieceSlot, 0),
				Scale:     "1.0",
				Height:    "215",
				Width:     "562",
				Divider:   "194",
			}
		}

		temp := forces[counter.Extra.Side]
		temp.PieceSlot = append(forces[counter.Extra.Side].PieceSlot, piece)
		forces[counter.Extra.Side] = temp

		gpid++
	}

	// tw.PanelWidget.ListWidget = append(tw.PanelWidget.ListWidget, mapToArray[counters.ListWidget](forces)...)
	tw.ListWidget = append(tw.ListWidget, mapToArray[counters.ListWidget](forces)...)
	g.PieceWindow.TabWidget = tw

	byt, err := xml.MarshalIndent(g, "", "  ")
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal the final game module data")
	}

	return byt, nil
}

func readFiles(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	filenames := make([]string, 0)

	for _, file := range files {
		if !file.IsDir() {
			if file.Name()[0] == '_' {
				filenames = append(filenames, file.Name())
			}
		}
	}

	return filenames, nil
}

func mapToArray[T any](m map[string]T) []T {
	temp := make([]T, len(m))

	i := 0
	for _, item := range m {
		temp[i] = item

		i++
	}

	return temp
}
