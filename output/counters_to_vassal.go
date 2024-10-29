package output

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"os"

	"github.com/pkg/errors"
	"github.com/sayden/counters"
	"github.com/sayden/counters/fsops"
)

// GetVassalDataForCounters returns the Vassal module data for the counters
func GetVassalDataForCounters(t *counters.CounterTemplate, xmlFilepath string) ([]byte, error) {
	var g counters.GameModule
	err := fsops.ReadMarkupFile(xmlFilepath, &g)
	if err != nil {
		return nil, errors.Wrap(err, "error trying to decode content")
	}

	tw := counters.TabWidget{
		EntryName:  "Forces",
		ListWidget: make([]counters.ListWidget, 0),
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

	xmlTemplateString := `+/null/prototype;BasicPrototype	piece;;;{{ .Filename }};{{ .PieceName}}/	null;0;0;{{ .Id }};0`
	xmlTemplate, err := template.New("xml").Parse(xmlTemplateString)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse template string")
	}

	// Read special files from images folder to statically load them into the module as pieces (terrain, -1, -2 markers, Disorganized, Spent, OOS, etc,)
	files, err := readFiles(counters.BASE_FOLDER + "/images")
	if err != nil {
		return nil, errors.Wrap(err, "could not read files")
	}

	gpid := 200
	id := 200

	for _, file := range files {
		buf := bytes.NewBufferString("")
		if err = xmlTemplate.ExecuteTemplate(
			buf, "xml", counters.TemplateData{
				Filename:  file, //+1 because file number starts in 1 instead of 0 when they are generated
				PieceName: file,
				Id:        fmt.Sprintf("%d", id),
			},
		); err != nil {
			return nil, errors.Wrap(err, "error trying to write Vassal xml file using templates")
		}
		id++

		piece := counters.PieceSlot{
			EntryName: file,
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

	filenamesInUse := make(map[string]bool)

	for i, counter := range t.Counters {
		buf := bytes.NewBufferString("")
		if err = xmlTemplate.ExecuteTemplate(
			buf, "xml", counters.TemplateData{
				Filename: counter.GetCounterFilename(
					t.PositionNumberForFilename,
					"",
					i+1,
					filenamesInUse,
				), //+1 because file number starts in 1 instead of 0 when they are generated
				PieceName: counter.GetCounterFilename(t.PositionNumberForFilename, "", -1, filenamesInUse),
				Id:        fmt.Sprintf("%d", id),
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
